package adopt

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

// Unlink reverses `bench link`. It consumes the link manifest and does the following,
// sparing anything the user has made theirs:
//   - removes every managed file whose fingerprint still matches
//   - prunes managed directories left empty
//   - strips the fenced AGENTS.md block
//   - removes the pre-push hook when it carries the managed marker
//
// A file edited since link is kept and reported. The manifest is removed last and only
// when nothing was refused, so a partial run stays resumable. `--dry-run` computes the
// same verdicts but writes nothing. With no manifest to consume, or one it cannot read,
// unlink fails loudly (exit 1). It never silently no-ops on a repo it cannot account for.
func Unlink(args []string, stdout, stderr io.Writer) int {
	dryRun := false
	for _, a := range args {
		switch a {
		case "--dry-run":
			dryRun = true
		default:
			fmt.Fprintln(stderr, "usage: bench unlink [--dry-run]")
			return 2
		}
	}
	root, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, toon.NotInRepo())
		return 1
	}

	manifestPath := filepath.Join(root, ".bench", "link-manifest.tsv")
	// Explicit absence guard: the shared reader maps an absent manifest to an empty one with
	// no error. unlink must not inherit that false-empty and exit 0 on a repo it cannot
	// account for. Absent or unreadable is a loud exit 1 with the documented manual path.
	if _, err := os.Stat(manifestPath); err != nil {
		fmt.Fprintln(stderr, "bench unlink: no link manifest at .bench/link-manifest.tsv; nothing to consume.")
		fmt.Fprintln(stderr, "  a repo linked before manifests must be uninstalled by hand (see the README uninstall section).")
		return 1
	}
	m, err := ReadManifest(manifestPath)
	if err != nil {
		fmt.Fprintf(stderr, "bench unlink: cannot read link manifest at .bench/link-manifest.tsv: %v\n", err)
		return 1
	}

	plan, err := planUnlink(root, m, dryRun)
	if err != nil {
		fmt.Fprintf(stderr, "bench unlink: %s: %v\n", root, err)
		return 1
	}
	if err := writeUnlinkReport(stdout, root, dryRun, plan); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if plan.manifestKept {
		return 3
	}
	return 0
}

// unlinkPlan is the resolved verdict set for one unlink walk: the ordering the report and the
// manifest-last rule both read. removed/keptModified/refused hold manifest rels; emptyDirs is
// the swept directory count; agentsAction/hookAction/manifestAction are the bespoke lines.
type unlinkPlan struct {
	removed      []string
	keptModified []string
	refused      []string
	emptyDirs    int
	agentsAction string
	hookAction   string
	manifestKept bool
}

func planUnlink(root string, m Manifest, dryRun bool) (unlinkPlan, error) {
	var p unlinkPlan

	rels := make([]string, 0, len(m.hashes))
	for rel := range m.hashes {
		rels = append(rels, rel)
	}
	sort.Strings(rels)

	removedSet := map[string]bool{}
	var candidateDirs []string
	for _, rel := range rels {
		full, ok := resolveInside(root, rel)
		if !ok || hasSymlinkParent(root, rel) {
			// A hand-edited row that escapes the repo (absolute path, `..` traversal, or a
			// symlinked parent) is refused and reported, never removed.
			p.refused = append(p.refused, rel)
			continue
		}
		if _, err := os.Lstat(full); os.IsNotExist(err) {
			// Already absent — no action, but its parents are still sweep candidates.
			collectAncestors(&candidateDirs, root, full)
			continue
		}
		fp, err := fingerprintPath(full)
		if err != nil {
			// A directory or other non-regular target cannot be a managed file: skip and
			// report it rather than delete it.
			p.refused = append(p.refused, rel)
			continue
		}
		if fp != m.hashes[rel] {
			// Edited since link — the user's content now: keep and report.
			p.keptModified = append(p.keptModified, rel)
			continue
		}
		if !dryRun {
			if err := os.Remove(full); err != nil {
				p.refused = append(p.refused, rel)
				continue
			}
		}
		p.removed = append(p.removed, rel)
		removedSet[full] = true
		collectAncestors(&candidateDirs, root, full)
	}

	p.emptyDirs = sweepEmptyDirs(candidateDirs, removedSet, dryRun)
	p.agentsAction = stripAgentsForUnlink(root, dryRun, &p)
	hookAction, err := removeManagedHook(root, dryRun)
	if err != nil {
		return unlinkPlan{}, err
	}
	p.hookAction = hookAction

	// Manifest removed last, and only when nothing was refused, so a partial run leaves the
	// residual managed state tracked for a follow-up. Dry-run never removes it.
	if len(p.keptModified)+len(p.refused) > 0 {
		p.manifestKept = true
	} else if !dryRun {
		_ = os.Remove(filepath.Join(root, ".bench", "link-manifest.tsv"))
	}
	return p, nil
}

// resolveInside resolves a manifest rel against the repo root and returns the absolute path
// only when it stays strictly inside the tree. Absolute paths and `..` traversal are refused,
// so a corrupted manifest can never name a removal target outside the repo.
func resolveInside(root, rel string) (string, bool) {
	if rel == "" || filepath.IsAbs(filepath.FromSlash(rel)) {
		return "", false
	}
	clean := filepath.Clean(filepath.FromSlash(rel))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", false
	}
	rootClean := filepath.Clean(root)
	full := filepath.Join(rootClean, clean)
	back, err := filepath.Rel(rootClean, full)
	if err != nil || back == ".." || strings.HasPrefix(back, ".."+string(filepath.Separator)) {
		return "", false
	}
	return full, true
}

// collectAncestors appends every directory between full's parent and the repo root (exclusive)
// to dirs, so the empty-directory sweep can prune a fully-emptied managed subtree deepest-first.
func collectAncestors(dirs *[]string, root, full string) {
	rootClean := filepath.Clean(root)
	for d := filepath.Dir(full); ; {
		dc := filepath.Clean(d)
		if dc == rootClean {
			return
		}
		back, err := filepath.Rel(rootClean, dc)
		if err != nil || back == "." || strings.HasPrefix(back, "..") {
			return
		}
		*dirs = append(*dirs, dc)
		parent := filepath.Dir(dc)
		if parent == dc {
			return
		}
		d = parent
	}
}

// sweepEmptyDirs prunes managed directories left empty, deepest-first with a
// non-recursive rmdir that only succeeds on an empty directory. Any directory still
// holding a kept-modified file or a user artifact survives untouched. It runs the same
// simulation for both modes. A directory is empty when every current entry is itself
// already gone, a removed file or an already-swept child. For dry-run this consults
// removedSet since nothing was written, and for a real run it matches the on-disk reality
// after removals. Returns the count.
func sweepEmptyDirs(candidateDirs []string, removedSet map[string]bool, dryRun bool) int {
	seen := map[string]bool{}
	dirs := candidateDirs[:0:0]
	for _, d := range candidateDirs {
		if !seen[d] {
			seen[d] = true
			dirs = append(dirs, d)
		}
	}
	sort.Slice(dirs, func(i, j int) bool {
		di, dj := strings.Count(dirs[i], string(filepath.Separator)), strings.Count(dirs[j], string(filepath.Separator))
		if di != dj {
			return di > dj
		}
		return dirs[i] > dirs[j]
	})
	count := 0
	for _, d := range dirs {
		entries, err := os.ReadDir(d)
		if err != nil {
			continue
		}
		empty := true
		for _, e := range entries {
			if !removedSet[filepath.Join(d, e.Name())] {
				empty = false
				break
			}
		}
		if !empty {
			continue
		}
		if !dryRun {
			if err := os.Remove(d); err != nil {
				continue
			}
		}
		removedSet[d] = true
		count++
	}
	return count
}

// stripAgentsForUnlink removes the fenced Bench block from AGENTS.md while preserving the
// user's surrounding prose. When stripping leaves the file whitespace-only, the case
// where link created it with no user content, the file is removed, mirroring link's
// create-if-absent symmetry. A malformed managed block is left in place and counted as a
// refusal, so the manifest survives for a manual fix. AGENTS.md is bespoke, not a
// manifest row.
func stripAgentsForUnlink(root string, dryRun bool, p *unlinkPlan) string {
	path := filepath.Join(root, "AGENTS.md")
	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	stripped, serr := StripAgentsBlock(string(content))
	if serr != nil {
		p.refused = append(p.refused, "AGENTS.md")
		return "kept AGENTS.md (its managed block could not be parsed)"
	}
	if stripped == string(content) {
		return ""
	}
	if strings.TrimSpace(stripped) == "" {
		if !dryRun {
			_ = os.Remove(path)
		}
		return "AGENTS.md (removed - no user prose remained)"
	}
	if !dryRun {
		_ = os.WriteFile(path, []byte(stripped), 0o644)
	}
	return "AGENTS.md (managed block stripped, prose kept)"
}

// removeManagedHook removes the pre-push hook only when it carries the managed marker. It
// resolves the effective hooks directory the same way link does, honoring core.hooksPath.
// A hook without the marker is a foreign hook, left in place; its presence is not a
// refusal because it was never Bench's. The hook is bespoke, not a manifest row.
func removeManagedHook(root string, dryRun bool) (string, error) {
	hooks, err := hooksDir(root)
	if err != nil {
		return "", err
	}
	path := filepath.Join(hooks, "pre-push")
	content, err := os.ReadFile(path)
	if err != nil || !strings.Contains(string(content), PrePushMarker) {
		return "", nil
	}
	if !dryRun {
		_ = os.Remove(path)
	}
	return "pre-push hook (managed - removed)", nil
}

func writeUnlinkReport(w io.Writer, root string, dryRun bool, p unlinkPlan) error {
	verb := "removed"
	keep := "kept (modified)"
	if dryRun {
		fmt.Fprintf(w, "bench unlink --dry-run: plan for %s (no changes written)\n", root)
		verb = "would remove"
		keep = "would keep (modified)"
	} else {
		fmt.Fprintf(w, "bench unlink: removed Bench from %s\n", root)
	}
	for _, rel := range p.removed {
		fmt.Fprintf(w, "  %s: %s\n", verb, rel)
	}
	if p.agentsAction != "" {
		fmt.Fprintf(w, "  %s %s\n", verb, p.agentsAction)
	}
	if p.hookAction != "" {
		fmt.Fprintf(w, "  %s %s\n", verb, p.hookAction)
	}
	if p.emptyDirs > 0 {
		fmt.Fprintf(w, "  %s %d empty %s\n", verb, p.emptyDirs, pluralDir(p.emptyDirs))
	}
	for _, rel := range p.keptModified {
		fmt.Fprintf(w, "  %s: %s\n", keep, rel)
	}
	for _, rel := range p.refused {
		fmt.Fprintf(w, "  refused: %s\n", rel)
	}
	// The prose above serves humans; the block below exposes the exact same verdict data
	// once for automation. Link and unlink deliberately share toon.Table rather than growing
	// bespoke parsers for their partial outcomes.
	verdicts := make([]lifecycleVerdict, 0, len(p.keptModified)+len(p.refused))
	for _, rel := range p.keptModified {
		verdicts = append(verdicts, lifecycleVerdict{rel, "modified"})
	}
	for _, rel := range p.refused {
		verdicts = append(verdicts, lifecycleVerdict{rel, "refused"})
	}
	if len(verdicts) > 0 {
		block, err := renderVerdicts("residuals", verdicts)
		if err != nil {
			return err
		}
		fmt.Fprint(w, block)
	}
	if len(p.removed) == 0 && p.agentsAction == "" && p.hookAction == "" && len(p.keptModified) == 0 && len(p.refused) == 0 {
		fmt.Fprintln(w, "  nothing to remove (manifest had no managed rows)")
	}
	switch {
	case p.manifestKept:
		fmt.Fprintln(w, "  kept link manifest (refusals present; re-run to finish)")
	case dryRun:
		fmt.Fprintln(w, "  would remove link manifest")
	default:
		fmt.Fprintln(w, "  removed link manifest")
	}
	return nil
}

func pluralDir(n int) string {
	if n == 1 {
		return "directory"
	}
	return "directories"
}
