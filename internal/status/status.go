// Package status ports `bench status`: the ambient dashboard a session-start hook
// consumes verbatim. It is the composition crux of the slice — every sibling query
// package (maps, structure, worktree, roadmap) contributes one signal, and this
// package renders them into the single severity-sorted board the shell renderer
// produced. The output format is stable because a hook parses it: the lead line,
// the fixed-width rows, and the `+N more` overflow.
//
// One rule per signal lives in its own package; status only orders the signals by
// severity and formats them. The merged-spec retirement counter (retirementCount) is
// the sole parser owned here, because no sibling surfaces it.
package status

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/maps"
	"github.com/gibbonmi/bench/internal/roadmap"
	"github.com/gibbonmi/bench/internal/structure"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/worktree"
)

// retireRe matches the unfenced retirement marker: a `Status:` line whose sole value is
// `implemented`, tab/space separated, with only whitespace trailing — the exact awk
// regex `^Status:[ \t]+implemented[ \t]*$`. We scan line by line, so `$` is the line end.
var retireRe = regexp.MustCompile(`^Status:[ \t]+implemented[ \t]*$`)

var captureOnlyStalePaths = map[string]bool{
	".bench-notes.md": true,
	"ROADMAP.md":      true,
}

// row is one dashboard signal: a severity (the sort/lead key), and the signal/detail/
// action triple the shell packed into a `sev|signal|detail|action` line.
type row struct {
	sev            int
	signal, detail string
	action         string
}

// Command implements `bench status`. It composes every sibling signal into the ambient
// board and returns it with exit 0. `-h/--help` prints usage (exit 0); an unknown
// argument is a usage error (exit 2); outside a repo is the structured error (exit 1).
func Command(args []string) (string, int) {
	switch {
	case len(args) == 0:
	case args[0] == "-h" || args[0] == "--help":
		return "usage: bench status\n", 0
	default:
		return toon.Usage("bench status", args[0]) + "\n", 2
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	return render(root), 0
}

// render gathers every signal under root, sorts ascending by severity, and formats the
// board. This is the byte-for-byte counterpart of the shell `status()` renderer.
func render(root string) string {
	var rows []row

	rows = appendGate(rows, root)
	rows = appendGit(rows, root)
	rows = appendWorktree(rows, root)
	rows = appendDrain(rows, root)
	rows = appendStructure(rows, root)
	rows = appendMaps(rows, root)
	rows = appendRetirement(rows, root)

	// Ascending numeric sort by severity; each severity is unique, so ordering is
	// fully determined and the min-severity row leads.
	sort.SliceStable(rows, func(i, j int) bool { return rows[i].sev < rows[j].sev })

	var b strings.Builder
	if len(rows) == 0 {
		b.WriteString("bench: clean — nothing pending\n")
		return b.String()
	}

	lead := rows[0]
	fmt.Fprintf(&b, "▶ %s  (%s)\n", lead.action, lead.signal)
	for i, r := range rows {
		if i < 5 {
			fmt.Fprintf(&b, "  %-10s %-30s → %s\n", r.signal, r.detail, r.action)
		}
	}
	if len(rows) > 5 {
		fmt.Fprintf(&b, "  +%d more\n", len(rows)-5)
	}
	return b.String()
}

// appendGate reads the gate verdict cache `<git-dir>/bench-last-gate` (line 0:
// `<status> <tree> …`). If the cached tree differs from the working tree, the verdict is
// stale (sev 6); else a `red` verdict is a fail-before-commit signal (sev 0). No cache
// file → no gate row.
func appendGate(rows []row, root string) []row {
	gitDir, err := git.Output("-C", root, "rev-parse", "--absolute-git-dir")
	if err != nil {
		return rows
	}
	data, err := os.ReadFile(filepath.Join(gitDir, "bench-last-gate"))
	if err != nil {
		return rows
	}
	first := string(data)
	if i := strings.IndexByte(first, '\n'); i >= 0 {
		first = first[:i]
	}
	fields := strings.Fields(first)
	tree := git.TreeHash(root)
	if !trustedGateCache(fields, tree) {
		ctree := ""
		if len(fields) > 1 {
			ctree = fields[1]
		}
		detail, action := staleGateDetailAction(root, ctree, tree)
		return append(rows, row{6, "gate", detail, action})
	}
	cstatus, ctree := fields[0], fields[1]
	if ctree != tree {
		detail, action := staleGateDetailAction(root, ctree, tree)
		return append(rows, row{6, "gate", detail, action})
	}
	if cstatus == "red" {
		return append(rows, row{0, "gate", "red", "fix before commit"})
	}
	return rows
}

func trustedGateCache(fields []string, currentTree string) bool {
	if len(fields) < 3 {
		return false
	}
	if fields[0] != "green" && fields[0] != "red" {
		return false
	}
	return fields[1] != "" && fields[1] != "none" && currentTree != "" && currentTree != "none"
}

func staleGateDetailAction(root, cachedTree, currentTree string) (detail, action string) {
	detail = fmt.Sprintf("stale (gated tree %s, work tree %s)", short(cachedTree), short(currentTree))
	action = "re-run the gate"
	paths, ok := git.ChangedPathsBetweenTrees(root, cachedTree, currentTree)
	if !ok || !captureOnlyDrift(paths) {
		return detail, action
	}
	return "stale (capture-only drift)", "re-run when convenient"
}

func captureOnlyDrift(paths []string) bool {
	if len(paths) == 0 {
		return false
	}
	for _, path := range paths {
		if !captureOnlyStalePaths[path] {
			return false
		}
	}
	return true
}

// appendGit adds the uncommitted/unpushed signal (sev 1). dirty is the porcelain status;
// ahead is the upstream-relative commit list, read only when an upstream is configured.
func appendGit(rows []row, root string) []row {
	dirty, _ := git.Output("-C", root, "status", "--porcelain")
	ahead := ""
	if git.OK("-C", root, "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}") {
		ahead, _ = git.Output("-C", root, "log", "--oneline", "@{u}..HEAD")
	}
	if dirty == "" && ahead == "" {
		return rows
	}
	detail := "uncommitted + unpushed"
	if dirty != "" && ahead == "" {
		detail = "uncommitted changes"
	} else if dirty == "" && ahead != "" {
		detail = "unpushed commits"
	}
	return append(rows, row{1, "git", detail, "commit on green / push"})
}

// appendWorktree counts worktree cleanup signals (sev 2): out-of-pool worktrees, leased
// pool entries, and orphaned harness scratch branches (`worktree-*` refs no registered
// worktree still holds). The repo root itself and a warm pooled entry (under the pool,
// no lease file on disk) are expected state, not a signal, and are skipped.
func appendWorktree(rows []row, root string) []row {
	pool := worktree.Pool(root)
	out, _ := git.Output("-C", root, "worktree", "list", "--porcelain")
	wtc := 0
	for _, line := range strings.Split(out, "\n") {
		if !strings.HasPrefix(line, "worktree ") {
			continue
		}
		wpath := line[len("worktree "):]
		if wpath == root {
			continue
		}
		if strings.HasPrefix(wpath, pool+"/") {
			lease, _ := worktree.LeaseFile(wpath)
			if !isRegularFile(lease) {
				continue
			}
		}
		wtc++
	}
	orphans := worktree.OrphanedDelegateBranches(root)
	if wtc == 0 {
		if len(orphans) == 0 {
			return rows
		}
		return append(rows, row{2, "worktree", plural(len(orphans), "orphaned worktree branch", "orphaned worktree branches"), "delete scratch branch"})
	}
	detail := fmt.Sprintf("%d active worktree(s)", wtc)
	if len(orphans) > 0 {
		detail += ", " + plural(len(orphans), "orphaned branch", "orphaned branches")
	}
	return append(rows, row{2, "worktree", detail, "resume or clean up (bench worktree)"})
}

func plural(n int, one, many string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s", one)
	}
	return fmt.Sprintf("%d %s", n, many)
}

// appendDrain adds the capture-drain signal (sev 3): parked ideas in IDEAS.md plus open
// journal headings, one combined row pointing at the single maintenance phase. The
// counts are roadmap.DrainCounts — the same counters `bench roadmap` reports. The
// learnings component shows only at or above the floor (env BENCH_LEARNINGS_FLOOR,
// default 1); parked ideas always count.
func appendDrain(rows []row, root string) []row {
	ideas, open := roadmap.DrainCounts(root)
	if open < learningsFloor() {
		open = 0
	}
	if ideas == 0 && open == 0 {
		return rows
	}
	return append(rows, row{3, "drain", fmt.Sprintf("%d idea(s), %d open learning(s)", ideas, open), "/bench-what-next"})
}

// appendStructure adds the structural-debt signal (sev 4) when the violation count is positive.
func appendStructure(rows []row, root string) []row {
	if n := structure.ViolationCount(root); n > 0 {
		return append(rows, row{4, "structure", fmt.Sprintf("%d issue(s)", n), "split (craft-seams)"})
	}
	return rows
}

// appendMaps adds the unresolved-decision-map signal (sev 5) when the count is positive.
func appendMaps(rows []row, root string) []row {
	if n := maps.UnresolvedCount(root); n > 0 {
		return append(rows, row{5, "decisions", fmt.Sprintf("%d unresolved map(s)", n), "craft-grill → /bench-write-spec"})
	}
	return rows
}

// appendRetirement adds the merged-spec-awaiting-retirement signal (sev 7), but only on
// the default branch — a topic branch's spec is still in flight, not awaiting retirement.
func appendRetirement(rows []row, root string) []row {
	cur, _ := git.Output("-C", root, "rev-parse", "--abbrev-ref", "HEAD")
	if cur != git.DefaultBranch(root) {
		return rows
	}
	if n := retirementCount(root); n > 0 {
		return append(rows, row{7, "specs", fmt.Sprintf("%d merged spec(s) awaiting retirement", n), "promote-then-delete (spec-retire)"})
	}
	return rows
}

// retirementCount counts specs/*.md files carrying an unfenced `Status: implemented`
// marker — a merged spec awaiting retirement. Absent `specs/` → 0. Per the awk logic each
// line is CRLF-stripped; a line whose first three bytes are a code fence toggles fence
// state and is skipped; lines inside a fence are skipped; the first line matching the
// retirement regex marks the file (counted once, scan stops).
func retirementCount(root string) int {
	dir := filepath.Join(root, "specs")
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return 0
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	n := 0
	for _, e := range entries {
		// specs/*.md without dotglob: skip directories, non-.md, and hidden names.
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		content, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		if fileAwaitsRetirement(content) {
			n++
		}
	}
	return n
}

// fileAwaitsRetirement is the per-file awk predicate: an unfenced retirement marker.
func fileAwaitsRetirement(content []byte) bool {
	inFence := false
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSuffix(line, "\r")
		if len(line) >= 3 && line[:3] == "```" {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		if retireRe.MatchString(line) {
			return true
		}
	}
	return false
}

// learningsFloor reads BENCH_LEARNINGS_FLOOR, defaulting to 1 when unset, empty, or not
// an integer — the shell `${BENCH_LEARNINGS_FLOOR:-1}` default.
func learningsFloor() int {
	v := os.Getenv("BENCH_LEARNINGS_FLOOR")
	if v == "" {
		return 1
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 1
	}
	return n
}

// short returns the first up-to-7 bytes of s (the shell `${var:0:7}` tree-prefix slice),
// guarding a short or "none" hash so the slice never panics.
func short(s string) string {
	return s[:min(7, len(s))]
}

// isRegularFile reports whether path exists and is a regular file (the shell `-f` test).
func isRegularFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}
