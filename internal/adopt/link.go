package adopt

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	kitpayload "github.com/gibbonmi/bench"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

type planEntry struct {
	src, rel, kind, content string
}

// distGitignore is written into a linked repo's .bench/dist/.gitignore. It ignores the
// arch-specific binary link copies next to it (.bench/dist/bench) so a consumer who
// commits .bench/dist/ can't hand a different-arch teammate a broken CLI and a silently
// fail-open Stop hook. A single unanchored line — not "/bench" — and it must not ignore
// itself: the ignore file is meant to be committed and travel with the repo.
const distGitignore = "bench\n"

func Link(args []string, stdout, stderr io.Writer, version string) int {
	mode := "copy"
	if len(args) > 0 {
		mode = args[0]
	}
	if mode != "copy" && mode != "symlink" {
		fmt.Fprintln(stderr, "usage: bench link [copy|symlink]")
		return 2
	}
	root, err := git.Root()
	if err != nil {
		// link's whole job is to create the Bench linkage, so the shared AXI
		// toon.NotInRepo() ("run inside a Bench-linked repo") is nonsensical here —
		// point at the actual remedy instead. Every genuine AXI query command keeps
		// the shared phrasing; this is link's own message, not a second source of it.
		fmt.Fprintln(stderr, toon.Errorf("not a git repository", "run inside a git repository (e.g. run 'git init' first)"))
		return 1
	}
	kit := kitDir()
	if isEphemeralKit(kit) && mode == "symlink" {
		mode = "copy"
		fmt.Fprintln(stdout, "(running from an ephemeral package cache - using copy mode so files don't dangle)")
	}
	plan, err := buildLinkPlan(kit)
	if err != nil {
		fmt.Fprintln(stderr, toon.Errorf("consumer payload allowlist rejected", err.Error()))
		return 1
	}
	result, _ := transactionalLink(root, kit, mode, version, plan, stdout, stderr)
	if result != 0 && result != 3 {
		return result
	}
	fmt.Fprintf(stdout, "linked Bench into %s (mode: %s).\n", root, mode)
	fmt.Fprintln(stdout, "  instructions: AGENTS.md managed block -> .bench/BENCH.md")
	fmt.Fprintln(stdout, "  portable surface: .agents/{skills,commands}")
	fmt.Fprintln(stdout, "  adapters: .claude/settings.json, .codex/hooks.json, .claude/{skills,commands} -> .agents")
	fmt.Fprintln(stdout, "  enforcement: shared .bench/hooks + git pre-push guard + the bench shift loop")
	fmt.Fprintln(stdout, "Run 'bench init' next to scaffold .bench/gate.sh.")
	return result
}

func isEphemeralKit(kit string) bool {
	return strings.Contains(kit, "_npx") ||
		strings.Contains(kit, "/dlx-") ||
		strings.Contains(kit, "npm-cache") ||
		strings.Contains(kit, "/.npm/_cacache/")
}

// buildLinkPlan is the sole reader of the consumer-payload allowlist for the link
// destination: every row it writes, and the audience filter that keeps kit-only rows
// out of every destination (including the .claude/ adapter mirrors), come from the
// root-level kitpayload package (embedding .bench/consumer-payload.json) rather than a
// second hand-listed plan. A consumer row's destination is its source path unchanged,
// except the top-level bin/ tree which lands under the linked repo's .bench/bin/ so a
// consumer's own bin/ stays untouched.
func buildLinkPlan(kit string) ([]planEntry, error) {
	rows, err := kitpayload.PayloadRows()
	if err != nil {
		return nil, err
	}
	excluded := kitpayload.PayloadKitOnlyPrefixes(rows)

	var plan []planEntry
	if bin := currentExecutablePath(); bin != "" {
		plan = append(plan, planEntry{src: bin, rel: ".bench/dist/bench", kind: "file"})
	}
	plan = append(plan, planEntry{rel: ".bench/dist/.gitignore", kind: "inline", content: distGitignore})

	for _, row := range kitpayload.PayloadConsumerRows(rows) {
		dest, ok := linkDestination(row.Source)
		if !ok {
			continue
		}
		src := filepath.Join(kit, filepath.FromSlash(row.Source))
		if !row.Tree {
			// A row absent from this kit checkout is skipped, not a hard failure:
			// a minimal or partial kit tree (a stripped fixture, an in-progress
			// checkout) omits assets it does not carry rather than promising ones
			// it cannot deliver. transactionalLink still catches a source that
			// disappears between planning and promotion.
			if _, statErr := os.Stat(src); statErr != nil {
				continue
			}
			plan = append(plan, planEntry{src: src, rel: dest, kind: "file"})
			continue
		}
		entries, err := treeEntries(src, dest, "file", row.Source, excluded)
		if err != nil {
			return nil, err
		}
		plan = append(plan, entries...)
	}

	adapterCommands, err := treeEntries(filepath.Join(kit, ".agents", "commands"), ".claude/commands", "adapter", ".agents/commands", excluded)
	if err != nil {
		return nil, err
	}
	plan = append(plan, adapterCommands...)

	adapterSkills, err := claudeSkillsEntries(kit, excluded)
	if err != nil {
		return nil, err
	}
	plan = append(plan, adapterSkills...)
	return plan, nil
}

// linkDestination maps an allowlist row's kit-relative source path to the path it is
// written at inside a linked repo. Every row's destination equals its source, except
// the top-level bin/ tree (npm's install-time surface), which is written under the
// linked repo's own .bench/bin/ instead of colliding with a consumer's bin/. A row
// outside the linked destinations this function names (README.md, CHANGELOG.md,
// projects/*.md) is not part of the link plan at all — it ships only in the npm
// tarball, via the same allowlist read by the release-evidence and package.json
// derivations.
func linkDestination(src string) (string, bool) {
	switch {
	case strings.HasPrefix(src, "bin/"):
		return ".bench/" + src, true
	case strings.HasPrefix(src, ".bench/"), strings.HasPrefix(src, ".claude/"), strings.HasPrefix(src, ".codex/"), strings.HasPrefix(src, ".agents/"):
		return src, true
	default:
		return "", false
	}
}

// treeEntries walks srcRoot (the on-disk directory for the allowlist row named by
// srcRel) and returns one planEntry per regular file, skipping any path the allowlist
// marks kit-only. A non-regular file (FIFO, device, socket) is refused by name instead
// of being silently skipped or blocking a later open — the kit's own source tree must
// never be able to wedge the plan builder. A symbolic link is refused on the same path
// and named as one: following it would copy bytes the allowlist never named, from
// wherever the link points, so the payload stays limited to files that are what they
// appear to be.
func treeEntries(srcRoot, destRoot, kind, srcRel string, excludedPrefixes []string) ([]planEntry, error) {
	info, err := os.Stat(srcRoot)
	if err != nil || !info.IsDir() {
		return nil, nil
	}
	var entries []planEntry
	walkErr := filepath.WalkDir(srcRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, relErr := filepath.Rel(srcRoot, path)
		if relErr != nil {
			return relErr
		}
		rel = filepath.ToSlash(rel)
		sourcePath := srcRel + "/" + rel
		if kitpayload.PayloadExcluded(sourcePath, excludedPrefixes) {
			return nil
		}
		fileInfo, infoErr := d.Info()
		if infoErr != nil {
			return infoErr
		}
		if !fileInfo.Mode().IsRegular() {
			kind := "non-regular file"
			if fileInfo.Mode()&os.ModeSymlink != 0 {
				kind = "symbolic link"
			}
			return fmt.Errorf("consumer payload tree %q contains a %s: %s", srcRel, kind, sourcePath)
		}
		entries = append(entries, planEntry{src: path, rel: filepath.ToSlash(filepath.Join(destRoot, rel)), kind: kind})
		return nil
	})
	if walkErr != nil {
		return nil, walkErr
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].rel < entries[j].rel })
	return entries, nil
}

// claudeSkillsEntries mirrors .agents/skills into .claude/skills for skills that have
// no same-named command (a skill with a command already reaches Claude Code through the
// command mirror). It reuses treeEntries for the walk, the kit-only exclusion, and the
// non-regular-file refusal, then drops the skills a command mirror already covers.
func claudeSkillsEntries(kit string, excludedPrefixes []string) ([]planEntry, error) {
	entries, err := treeEntries(filepath.Join(kit, ".agents", "skills"), ".claude/skills", "adapter", ".agents/skills", excludedPrefixes)
	if err != nil {
		return nil, err
	}
	filtered := entries[:0]
	for _, e := range entries {
		top := strings.SplitN(strings.TrimPrefix(e.rel, ".claude/skills/"), "/", 2)[0]
		if _, err := os.Stat(filepath.Join(kit, ".agents", "commands", top+".md")); err == nil {
			continue
		}
		filtered = append(filtered, e)
	}
	return filtered, nil
}

func currentExecutablePath() string {
	path, err := os.Executable()
	if err != nil {
		return ""
	}
	info, err := os.Stat(path)
	if err != nil || info.IsDir() || info.Size() == 0 || info.Mode().Perm()&0o111 == 0 {
		return ""
	}
	return path
}

func validateAgentsPath(path string) error {
	content, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return validateAgentsContent(string(content))
}

func manifestHash(root, rel string) string {
	m, err := ReadManifest(filepath.Join(root, ".bench", "link-manifest.tsv"))
	if err != nil {
		return ""
	}
	return m.Hash(rel)
}

func hasSymlinkParent(root, rel string) bool {
	dir := filepath.Dir(rel)
	if dir == "." {
		return false
	}
	path := root
	for _, part := range strings.Split(filepath.ToSlash(dir), "/") {
		path = filepath.Join(path, part)
		if info, err := os.Lstat(path); err == nil && info.Mode()&os.ModeSymlink != 0 {
			return true
		}
	}
	return false
}

func AdapterTarget(rel string) (string, bool) {
	var sourceRel string
	switch {
	case strings.HasPrefix(rel, ".claude/commands/"):
		sourceRel = ".agents/commands/" + strings.TrimPrefix(rel, ".claude/commands/")
	case strings.HasPrefix(rel, ".claude/skills/"):
		sourceRel = ".agents/skills/" + strings.TrimPrefix(rel, ".claude/skills/")
	default:
		return "", false
	}
	destDir := filepath.ToSlash(filepath.Dir(rel))
	parts := strings.Split(destDir, "/")
	ups := strings.Repeat("../", len(parts))
	return ups + sourceRel, true
}

func copyFile(src, dest string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dest, data, info.Mode().Perm())
}
