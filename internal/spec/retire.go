package spec

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// RetireNextPrefix opens the last line of a completed retire; the retire listing goes before it.
const RetireNextPrefix = "next: "

// checkoutRefusal is the checkout rule of every spec write leaf. It returns the
// identity-unknown error or the shared primary-checkout refusal, and primary reports the
// second case so a leaf can append its own follow-on step. A linked worktree returns "".
func checkoutRefusal(base string) (refusal string, primary bool) {
	primary, err := git.IsPrimaryCheckout(base)
	if err != nil {
		return toon.Errorf("checkout identity is unknown", "repair Git metadata, then retry from a Bench worktree"), false
	}
	if primary {
		return usage.PrimaryCheckoutRefusal(), true
	}
	return "", false
}

// retireCommand runs `bench spec retire` on its parsed operand. On a merged-implemented spec, it
// deletes the review pickup when present, the tickets, the spec file, and the complete folder,
// and it prints what it removed plus the judgment duties that remain. It never commits and
// never runs the gate; `bench commit` owns commit discipline. Every unsafe input refuses at
// exit 1 without deleting anything. A spec refuses when it is not merged-implemented: staged,
// or implemented only in the working tree and not yet at HEAD. An unknown slug refuses, and so
// does an orphaned review pickup with no spec. Inside a repository, the primary checkout
// refuses too, because Bench write verbs run from a worktree. retireNext renders the next:
// line from the Roadmap: value, read before deletion, and from the delivery at HEAD.
func retireCommand(arg string) (string, int) {
	base := RepoBase()
	if base != "" {
		if refusal, primary := checkoutRefusal(base); refusal != "" {
			if primary {
				// The shared refusal names only the worktree creation. Retire appends its own
				// follow-on step here, at the one call site that needs it, so the operator reads
				// one route and no other verb inherits a retire step.
				refusal += ", then rerun bench spec retire from that worktree"
			}
			return refusal + "\n", 1
		}
	}
	content, resolved, tried, found, err := Resolve(base, arg)
	if err != nil {
		folder := filepath.Join(base, filepath.FromSlash(LiveSpecPath(arg)))
		if resolved == folder {
			if residue, ok := folderResidue(base, arg); ok {
				return toon.Errorf("incomplete retired spec folder: "+residue,
					"remove "+residue+" by hand after reviewing its residue; retire will not auto-clean it") + "\n", 1
			}
		}
		return toon.Errorf(fmt.Sprintf("spec not readable: %s: %v", resolved, err), "check file permissions") + "\n", 1
	}
	if !found {
		// A review pickup with no spec is a reviewer judgment, never an auto-clean.
		if orphan, ok := orphanPickup(base, arg); ok {
			return toon.Errorf("orphaned review pickup: "+orphan+" has no spec",
				"delete "+orphan+" by hand if its residual risk is accepted; retire never auto-cleans a pickup") + "\n", 1
		}
		return toon.Errorf("spec not found: "+strings.Join(tried, ", "), "pass a merged spec path or slug") + "\n", 1
	}
	// Merged-implemented reuses AwaitsRetirement, the status detector, and also requires
	// the marker at HEAD. An implemented-but-uncommitted spec means the finishing commit
	// has not landed, so retiring it would delete work git cannot recover.
	if !AwaitsRetirement(content) {
		return toon.Errorf("spec not merged-implemented: "+RelTo(base, resolved),
			"retire only a spec whose Status line reads implemented") + "\n", 1
	}
	if !implementedAtHEAD(base, resolved) {
		return toon.Errorf("spec implemented in the working tree but not at HEAD: "+RelTo(base, resolved),
			"commit the finishing flip before retiring") + "\n", 1
	}
	plan := retirementPlan(base, resolved)
	if blocked := blockedRemovalTargets(base, plan); len(blocked) > 0 {
		return toon.Errorf("retirement plan is not writable: "+strings.Join(blocked, ", "), "check file permissions") + "\n", 1
	}
	// Deletion order leaves each recoverable interrupt state with a spec file, never an
	// orphaned review pickup. Once the spec file is gone, remaining folder content is terminal.
	var b strings.Builder
	slug := SlugOf(resolved)
	for _, target := range plan {
		if err := target.remove(); err != nil {
			return b.String() + toon.Errorf(fmt.Sprintf("remove %s: %v", RelTo(base, target.path), err), "check file permissions") + "\n", 1
		}
		switch target.path {
		case filepath.Join(base, "reviews", slug+".md"):
			fmt.Fprintf(&b, "retired: %s\n", RelTo(base, target.path))
		case filepath.Dir(resolved):
			fmt.Fprintf(&b, "retired: %s\n", RelTo(base, target.path))
		case resolved:
			if filepath.Base(resolved) != "spec.md" {
				fmt.Fprintf(&b, "retired: %s\n", RelTo(base, target.path))
			}
		}
	}
	_, roadmapID := metadata(content)
	b.WriteString(retireNext(base, resolved, roadmapID, slug))
	return b.String(), 0
}

// removalTarget is one planned filesystem removal. The plan is shared by the
// writability preflight and the apply loop so no target can be deleted without first
// being checked.
type removalTarget struct {
	path      string
	recursive bool
}

func (target removalTarget) remove() error {
	if target.recursive {
		return os.RemoveAll(target.path)
	}
	return os.Remove(target.path)
}

// retirementPlan resolves every target before the first removal. A folder spec's
// tickets, spec file, and folder are distinct planned removals because each has its
// own filesystem permission boundary.
func retirementPlan(base, resolved string) []removalTarget {
	plan := []removalTarget{}
	slug := SlugOf(resolved)
	if pickup := filepath.Join(base, "reviews", slug+".md"); fileExists(pickup) {
		plan = append(plan, removalTarget{path: pickup})
	}
	if filepath.Base(resolved) != "spec.md" {
		return append(plan, removalTarget{path: resolved})
	}
	folder := filepath.Dir(resolved)
	return append(plan,
		removalTarget{path: filepath.Join(folder, "tickets"), recursive: true},
		removalTarget{path: resolved},
		removalTarget{path: folder, recursive: true},
	)
}

// blockedRemovalTargets probes each target's parent with a temporary file. Deletion
// permission belongs to the parent directory, so this establishes the same local
// write capability without touching a planned target. Recursive targets additionally
// probe every real descendant directory because RemoveAll deletes entries there. The
// returned rendering is sorted to keep a refusal stable even when apply order differs.
func blockedRemovalTargets(base string, plan []removalTarget) []string {
	blocked := []string{}
	for _, target := range plan {
		err := writableDirectory(filepath.Dir(target.path))
		if err == nil && target.recursive {
			err = writableRemovalTree(target.path)
		}
		if err != nil {
			blocked = append(blocked, RelTo(base, target.path))
		}
	}
	sort.Strings(blocked)
	return blocked
}

// writableDirectory confirms that entries can be created and removed in dir without
// retaining a filesystem change. That is the capability a removal in dir requires.
func writableDirectory(dir string) error {
	probe, err := os.CreateTemp(dir, ".bench-retire-preflight-")
	if err != nil {
		return err
	}
	if err := probe.Close(); err != nil {
		_ = os.Remove(probe.Name())
		return err
	}
	return os.Remove(probe.Name())
}

// writableRemovalTree checks the directories RemoveAll will empty. Lstat keeps a
// symlink as a leaf: preflight never follows it into a directory outside the plan.
func writableRemovalTree(path string) error {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return err
	}
	if err := writableDirectory(path); err != nil {
		return err
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		child := filepath.Join(path, entry.Name())
		info, err := os.Lstat(child)
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
			continue
		}
		if err := writableRemovalTree(child); err != nil {
			return err
		}
	}
	return nil
}

// roadmapRe matches a well-formed Roadmap: value: `FT` followed by one or more ASCII
// digits, nothing else. metadata already trims surrounding whitespace, including NBSP.
var roadmapRe = regexp.MustCompile(`^FT[0-9]+$`)

// roadmapRemainder renders the retire next: line's board-remainder clause. A
// well-formed roadmapID names the exact ROADMAP.md row, and, when
// roadmap/<roadmapID>.md exists as a regular file under base, names that path too.
// Any other value, including empty, falls back to the generic `FT<n>` placeholder
// naming both the row and its detail file.
func roadmapRemainder(base, roadmapID string) string {
	if !roadmapRe.MatchString(roadmapID) {
		return " FT<n> and its roadmap/FT<n>.md detail file"
	}
	detail := filepath.Join(base, "roadmap", roadmapID+".md")
	if fileExists(detail) {
		return fmt.Sprintf(" %s and its roadmap/%s.md detail file", roadmapID, roadmapID)
	}
	return " " + roadmapID
}

// folderResidue identifies the terminal interrupted-retire state: a folder remains but
// its spec file is gone. The function cannot infer which residual files a reviewer
// accepts.
func folderResidue(base, arg string) (string, bool) {
	if strings.ContainsRune(arg, '/') {
		return "", false
	}
	folder := filepath.Join(base, "specs", strings.TrimSuffix(arg, ".md"))
	fi, err := os.Stat(folder)
	if err != nil || !fi.IsDir() {
		return "", false
	}
	if _, err := os.Stat(filepath.Join(folder, "spec.md")); !os.IsNotExist(err) {
		return "", false
	}
	return RelTo(base, folder), true
}

// implementedAtHEAD reports whether the spec's content at HEAD carries the retirement
// marker — the "finishing commit has landed" guard. It reads the blob through git and
// discards stderr, so a spec absent from HEAD or still staged there reads as false.
func implementedAtHEAD(base, resolved string) bool {
	content, err := headBlob(base, filepath.ToSlash(RelTo(base, resolved)))
	if err != nil {
		return false
	}
	return AwaitsRetirement(content)
}

// orphanPickup reports the repo-relative reviews/<slug>.md when a review pickup exists for a
// slug whose spec did not resolve — the orphaned-pickup refusal names it.
func orphanPickup(base, arg string) (string, bool) {
	pickup := filepath.Join(base, "reviews", SlugOf(arg)+".md")
	if fileExists(pickup) {
		return RelTo(base, pickup), true
	}
	return "", false
}
