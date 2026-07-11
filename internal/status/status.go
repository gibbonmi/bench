// Package status ports `bench status`: the ambient dashboard a session-start hook
// consumes verbatim. It is the composition crux of the slice — every sibling query
// package (maps, structure, worktree, roadmap) contributes one signal, and this
// package renders them into the single severity-sorted board the shell renderer
// produced. The output format is stable because a hook parses it: the lead line,
// the fixed-width rows, and the `+N more` overflow.
//
// One rule per signal lives in its own package; status only orders the signals by
// severity and formats them. The specs housekeeping signals are counted here —
// retirementCount over specs/ (applying the shared spec.AwaitsRetirement predicate),
// orphanedPickupCount pairing reviews/ against specs/, and roadmapReconcileCounts
// classifying ROADMAP.md's spec-path tokens against the tree — but the merged-implemented
// predicate itself is spec.AwaitsRetirement, one source shared with `bench spec retire`.
package status

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/adopt"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/maps"
	"github.com/gibbonmi/bench/internal/roadmap"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/structure"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/worktree"
)

var captureOnlyStalePaths = map[string]bool{
	".bench-notes.md": true,
	"IDEAS.md":        true,
	"ROADMAP.md":      true,
}

// row is one dashboard signal: a severity (the sort/lead key), and the signal/detail/
// action triple the shell packed into a `sev|signal|detail|action` line.
type row struct {
	sev            int
	signal, detail string
	action         string
}

// Signal is one ambient-board row exposed as structured data — the severity sort key
// plus the signal/detail/action triple. It is the shared accessor the text board and the
// human dashboard both consume, so the two views cannot rank or drop signals differently.
type Signal struct {
	Severity int
	Name     string
	Detail   string
	Action   string
}

// GateInfo is the structured gate-verdict cache read, shared by the status board and the
// dashboard so neither re-parses `<git-dir>/bench-last-gate`. Present is false when no
// cache file exists; Stale marks a verdict whose cached tree no longer matches the work
// tree (or whose line is untrusted). Status/CachedTree/WorkTree/Timestamp carry the raw
// fields for a human view; the board reduces them to its severity rows.
type GateInfo struct {
	Present    bool
	Status     string
	CachedTree string
	WorkTree   string
	Stale      bool
	Timestamp  string
}

// Signals gathers every ambient signal under root and returns them severity-sorted
// ascending — the one severity ladder `bench status` renders. render (the text board) and
// the dashboard gatherer both call this, so a signal added or reordered here reaches both
// surfaces from one source.
func Signals(root string) []Signal {
	var rows []row

	rows = appendGate(rows, root)
	rows = appendGit(rows, root)
	rows = appendWorktree(rows, root)
	rows = appendGuards(rows, root)
	rows = appendDrain(rows, root)
	rows = appendStructure(rows, root)
	rows = appendMaps(rows, root)
	rows = appendRetirement(rows, root)
	rows = appendOrphanedPickup(rows, root)
	rows = appendRoadmapReconcile(rows, root)

	// Ascending numeric sort by severity; each severity is unique, so ordering is
	// fully determined and the min-severity row leads.
	sort.SliceStable(rows, func(i, j int) bool { return rows[i].sev < rows[j].sev })

	out := make([]Signal, len(rows))
	for i, r := range rows {
		out[i] = Signal{Severity: r.sev, Name: r.signal, Detail: r.detail, Action: r.action}
	}
	return out
}

// Command implements `bench status`. It composes every sibling signal into the ambient
// board and returns it with exit 0. `--all` lifts the five-row budget and prints every
// signal; `-h/--help` prints usage (exit 0); an unknown argument — including `--all`
// with any trailing token — is a usage error (exit 2); outside a repo is the structured
// error (exit 1).
func Command(args []string) (string, int) {
	all := false
	switch {
	case len(args) == 0:
	case args[0] == "-h" || args[0] == "--help":
		return "usage: bench status [--all]\n", 0
	case args[0] == "--all" && len(args) == 1:
		all = true
	default:
		return toon.Usage("bench status", args[0]) + "\n", 2
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	return render(root, all), 0
}

// render gathers every signal under root, sorts ascending by severity, and formats the
// board. This is the byte-for-byte counterpart of the shell `status()` renderer. When all
// is false it applies the five-row budget and appends the overflow line; when all is true
// (`bench status --all`) it prints every row and emits no overflow line. The SessionStart
// hook calls with all=false so the ambient surface stays bounded.
func render(root string, all bool) string {
	signals := Signals(root)

	var b strings.Builder
	if len(signals) == 0 {
		b.WriteString("bench: clean — nothing pending\n")
		return b.String()
	}

	lead := signals[0]
	fmt.Fprintf(&b, "▶ %s  (%s)\n", lead.Action, lead.Name)
	for i, r := range signals {
		if all || i < 5 {
			fmt.Fprintf(&b, "  %-10s %-30s → %s\n", r.Name, r.Detail, r.Action)
		}
	}
	if !all && len(signals) > 5 {
		fmt.Fprintf(&b, "  +%d more (bench status --all)\n", len(signals)-5)
	}
	return b.String()
}

// appendGate reads the gate verdict cache `<git-dir>/bench-last-gate` (line 0:
// `<status> <tree> …`). If the cached tree differs from the working tree, the verdict is
// stale (sev 7); else a `red` verdict is a fail-before-commit signal (sev 0). No cache
// file → no gate row.
func appendGate(rows []row, root string) []row {
	gv := GateVerdict(root)
	if !gv.Present {
		return rows
	}
	if gv.Stale {
		detail, action := staleGateDetailAction(root, gv.CachedTree, gv.WorkTree)
		return append(rows, row{7, "gate", detail, action})
	}
	if gv.Status == "red" {
		return append(rows, row{0, "gate", "red", "fix before commit"})
	}
	return rows
}

// GateVerdict reads the gate cache `<git-dir>/bench-last-gate` (line 0:
// `<status> <tree> <timestamp>`) into structured data. It is the one gate-cache reader —
// appendGate reduces it to a severity row, the dashboard renders its raw fields — so the
// stale/red honesty is computed once. A missing cache file or unresolvable git dir yields
// Present=false; an untrusted line, or a cached tree that differs from the work tree,
// yields Stale=true. WorkTree is always the current tree hash when a cache file exists.
func GateVerdict(root string) GateInfo {
	gitDir, err := git.Output("-C", root, "rev-parse", "--absolute-git-dir")
	if err != nil {
		return GateInfo{}
	}
	data, err := os.ReadFile(filepath.Join(gitDir, git.GateCacheFile))
	if err != nil {
		return GateInfo{}
	}
	first := string(data)
	if i := strings.IndexByte(first, '\n'); i >= 0 {
		first = first[:i]
	}
	fields := strings.Fields(first)
	tree := git.TreeHash(root)
	gi := GateInfo{Present: true, WorkTree: tree}
	if len(fields) > 0 {
		gi.Status = fields[0]
	}
	if len(fields) > 1 {
		gi.CachedTree = fields[1]
	}
	if len(fields) > 2 {
		gi.Timestamp = fields[2]
	}
	if !trustedGateCache(fields, tree) || gi.CachedTree != tree {
		gi.Stale = true
	}
	return gi
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
	// Audit #3/#4 — tolerate: this is an ambient advisory board the SessionStart hook
	// consumes, so a git failure must degrade the git row, not crash the hook. dirty drops
	// its porcelain error; ahead is read only after an OK-checked upstream and degrades to
	// no ahead-count.
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

// appendWorktree adds separate worktree signals (sev 2) for out-of-pool worktrees,
// leased pool entries, and orphaned harness scratch branches (`worktree-*` refs no
// registered worktree still holds). The repo root and warm pooled entries are expected
// state, not signals.
func appendWorktree(rows []row, root string) []row {
	registered, err := worktree.ClassifyRegisteredWorktrees(root)
	if err != nil {
		// A classify failure means the pool/leased/out-of-pool counts below are
		// unknowable, not zero: surface the git failure itself as the row rather than
		// falling through to an empty-looking board (the false-empty class FT29 swept).
		return append(rows, row{2, "worktree", fmt.Sprintf("git worktree list failed: %v", err), "investigate the git failure, then re-run bench status"})
	}
	outOfPool, leased := 0, 0
	for _, wt := range registered {
		switch wt.Class {
		case worktree.ClassOutOfPool:
			outOfPool++
		case worktree.ClassPoolLease:
			leased++
		}
	}
	if outOfPool > 0 {
		rows = append(rows, row{2, "worktree", plural(outOfPool, "out-of-pool worktree", "out-of-pool worktrees"), "clean up (bench worktree clean)"})
	}
	if leased > 0 {
		rows = append(rows, row{2, "worktree", plural(leased, "leased pool worktree", "leased pool worktrees"), "resume leased worktree"})
	}
	orphans := worktree.OrphanedDelegateBranches(root)
	if len(orphans) == 0 {
		return rows
	}
	// The action splits by the sweep's own landed proof, so following it always changes
	// something: landed branches disappear under `bench worktree clean`, while a kept
	// branch would survive it — the honest action there is a hand inspection. With an
	// unresolvable default no orphan is classifiable; the combined row stands and the
	// recommended sweep refuses loudly with the reason.
	def, ok := worktree.ResolvedDefault(root)
	if !ok {
		return append(rows, row{2, "worktree", plural(len(orphans), "orphaned worktree branch", "orphaned worktree branches"), "bench worktree clean"})
	}
	landed, kept := 0, 0
	for _, branch := range orphans {
		if isLanded, _ := worktree.LandedInDefault(root, branch, def); isLanded {
			landed++
		} else {
			kept++
		}
	}
	if landed > 0 {
		rows = append(rows, row{2, "worktree", plural(landed, "orphaned worktree branch", "orphaned worktree branches"), "bench worktree clean"})
	}
	if kept > 0 {
		rows = append(rows, row{2, "worktree", plural(kept, "un-landed salvage branch", "un-landed salvage branches"), "inspect salvage branch(es) — bench worktree clean keeps them"})
	}
	return rows
}

func plural(n int, one, many string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s", one)
	}
	return fmt.Sprintf("%d %s", n, many)
}

// appendGuards adds the pre-push backstop signal (sev 3), ranked just below the worktree
// signals and above the drain row. git does not clone hooks, so a fresh clone silently loses
// the harness-independent default-branch backstop; this surfaces the gap ambiently rather
// than only under `bench doctor`. It fires only on the primary checkout of a routed repo
// (`.bench/lines.env` present) — pool and linked worktrees share the main `.git` and must not
// double-report the same hook — and stays quiet when the hook is bench-managed. Remedy: bench link.
func appendGuards(rows []row, root string) []row {
	if !isPrimaryCheckout(root) {
		return rows
	}
	if _, err := os.Stat(filepath.Join(root, ".bench", "lines.env")); err != nil {
		return rows
	}
	st := adopt.ClassifyPrePush(root)
	if st.State == adopt.PrePushManaged {
		return rows
	}
	return append(rows, row{3, "guards", prePushDetail(st.State), "bench link"})
}

// prePushDetail names the pre-push gap the guards row reports, mirroring the adopt classifier's
// states so the ambient signal and the doctor row describe the same condition.
func prePushDetail(state adopt.PrePushState) string {
	switch state {
	case adopt.PrePushForeign:
		return "pre-push not bench-managed"
	case adopt.PrePushDiverted:
		return "pre-push diverted (core.hooksPath)"
	default: // PrePushAbsent
		return "pre-push missing"
	}
}

// isPrimaryCheckout reports whether root is the repository's primary checkout — the git dir
// equals the git-common-dir. A linked or pool worktree's git dir is `.git/worktrees/<name>`
// while its common dir is the main `.git`, so the two differ there. The same test the
// worktree classifier's canonicalRoot uses. An undeterminable repo returns false, so the
// low-noise signal stays quiet rather than double-reporting.
func isPrimaryCheckout(root string) bool {
	gitDir, err1 := git.Output("-C", root, "rev-parse", "--path-format=absolute", "--git-dir")
	common, err2 := git.Output("-C", root, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err1 != nil || err2 != nil {
		return false
	}
	return filepath.Clean(gitDir) == filepath.Clean(common)
}

// appendDrain adds the capture-drain signal (sev 4): parked ideas in IDEAS.md plus open
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
	return append(rows, row{4, "drain", fmt.Sprintf("%d idea(s), %d open learning(s)", ideas, open), "/bench-what-next"})
}

// appendStructure adds the structural-debt signal (sev 5) when the violation count is positive.
func appendStructure(rows []row, root string) []row {
	if n := structure.ViolationCount(root); n > 0 {
		return append(rows, row{5, "structure", fmt.Sprintf("%d issue(s)", n), "split (craft-seams)"})
	}
	return rows
}

// appendMaps adds the unresolved-decision-map signal (sev 6) when the count is positive.
func appendMaps(rows []row, root string) []row {
	if n := maps.UnresolvedCount(root); n > 0 {
		return append(rows, row{6, "decisions", fmt.Sprintf("%d unresolved map(s)", n), "/bench-shape-idea"})
	}
	return rows
}

// appendRetirement adds the merged-spec-awaiting-retirement signal (sev 8), but only on
// the default branch — a topic branch's spec is still in flight, not awaiting retirement.
func appendRetirement(rows []row, root string) []row {
	// Audit #5 — tolerate: a failure reads as "not the default branch", skipping this
	// advisory housekeeping signal; non-fatal on the ambient board.
	cur, _ := git.Output("-C", root, "rev-parse", "--abbrev-ref", "HEAD")
	if cur != git.DefaultBranch(root) {
		return rows
	}
	if n := retirementCount(root); n > 0 {
		return append(rows, row{8, "specs", fmt.Sprintf("%d merged spec(s) awaiting retirement", n), "bench spec retire <slug>"})
	}
	return rows
}

// appendOrphanedPickup adds the orphaned-review-pickup signal (sev 9): a reviews/<slug>.md
// with no matching specs/<slug>.md is a pickup file that escaped its lifecycle. It ranks
// with the housekeeping rows, just below the retirement signal, so it never displaces the
// gate/git rows in the budget. A paired pickup (its spec still present) is expected state.
func appendOrphanedPickup(rows []row, root string) []row {
	if n := orphanedPickupCount(root); n > 0 {
		return append(rows, row{9, "reviews", plural(n, "orphaned review pickup", "orphaned review pickups"), "promote or delete by hand"})
	}
	return rows
}

// appendRoadmapReconcile adds the roadmap-reconcile signal (sev 9): a ROADMAP.md row naming a
// specs/<slug>.md whose work has already shipped — the spec file is merged-implemented
// (spec.AwaitsRetirement) or was retired out of the tree entirely — has outlived the drain that
// should have removed it. Like appendRetirement, it fires only on the default branch: a topic
// branch's roadmap is mid-build, so a row there names in-flight work, not a shipped-work leak.
// Severity 10 ranks it below the housekeeping rows (retirement 8, orphaned-pickup 9) and far
// below gate/git, so it never displaces a red-gate or dirty-tree row in the budget.
func appendRoadmapReconcile(rows []row, root string) []row {
	// Audit #6 — tolerate: as in appendRetirement, an unreadable branch reads as "not the
	// default branch" and skips this advisory reconcile signal.
	cur, _ := git.Output("-C", root, "rev-parse", "--abbrev-ref", "HEAD")
	if cur != git.DefaultBranch(root) {
		return rows
	}
	merged, dangling := roadmapReconcileCounts(root)
	if merged == 0 && dangling == 0 {
		return rows
	}
	var details []string
	if merged > 0 {
		details = append(details, plural(merged, "row for merged work", "rows for merged work"))
	}
	if dangling > 0 {
		details = append(details, plural(dangling, "row names a retired spec", "rows name a retired spec"))
	}
	return append(rows, row{10, "roadmap", strings.Join(details, ", "), "/bench-what-next"})
}

// retirementCount counts specs/*.md files that spec.AwaitsRetirement marks — a merged spec
// awaiting retirement. Absent `specs/` → 0. The unfenced-marker predicate is
// spec.AwaitsRetirement's one source.
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
		if spec.AwaitsRetirement(content) {
			n++
		}
	}
	return n
}

// orphanedPickupCount counts reviews/*.md files with no matching specs/<slug>.md — a review
// pickup whose spec retired first or was never present. Absent `reviews/` → 0. Hidden and
// non-.md entries are skipped, mirroring the retirementCount dir-walk.
func orphanedPickupCount(root string) int {
	entries, err := os.ReadDir(filepath.Join(root, "reviews"))
	if err != nil {
		return 0
	}
	n := 0
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		if info, err := os.Stat(filepath.Join(root, "specs", e.Name())); err == nil && !info.IsDir() {
			continue // paired: its spec is still present
		}
		n++
	}
	return n
}

// roadmapReconcileRe matches a spec-path token in a ROADMAP.md row: `specs/<slug>.md` with a
// kebab/alnum slug. It matches the bare path inside backticks/bold since the markdown decoration
// isn't part of the token, and the char class excludes `<>` so a literal `specs/<slug>.md`
// placeholder in the header prose can't false-fire.
// roadmapReconcileCounts scans ROADMAP.md for specs/<slug>.md path tokens and classifies each
// distinct path against the tree: a missing file is a dangling row (the spec retired but its
// roadmap row survived); a present file that spec.AwaitsRetirement marks is a merged row (the
// work shipped but the drain missed it). A present, still-staged spec is the normal open-work
// state and counts nothing. Absent ROADMAP.md → 0, 0. The merged predicate is
// spec.AwaitsRetirement, the same one source the retirement counter applies.
func roadmapReconcileCounts(root string) (merged, dangling int) {
	data, err := os.ReadFile(filepath.Join(root, "ROADMAP.md"))
	if err != nil {
		return 0, 0
	}
	seen := map[string]bool{}
	for _, slug := range roadmap.SpecSlugs(data) {
		path := "specs/" + slug + ".md"
		if seen[path] {
			continue
		}
		seen[path] = true
		content, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			dangling++
			continue
		}
		if spec.AwaitsRetirement(content) {
			merged++
		}
	}
	return merged, dangling
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
