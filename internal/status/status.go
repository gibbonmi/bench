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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/adopt"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/learnings"
	"github.com/gibbonmi/bench/internal/maps"
	"github.com/gibbonmi/bench/internal/retros"
	"github.com/gibbonmi/bench/internal/roadmap"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/shift"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/structure"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
	"github.com/gibbonmi/bench/internal/worktree"
)

// grammar is the declared argument shape usage.Parse enforces for this subcommand —
// arity, flag recognition, `--`, and help all come from there rather than a local switch.
var grammar = usage.Grammar{
	Cmd:   "bench status",
	Help:  "usage: bench status [--all]",
	Flags: []usage.Flag{{Name: "--all"}},
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
// cache file exists; Stale primarily marks a ready green verdict whose cached tree
// differs from the work tree. It also marks an exact-tip non-reusable green the gate
// cannot compose as a whole-tree verdict. Partition carries narrowness at component
// granularity: non-nil for a partial verdict that graded only the components whose
// inputs moved, nil for a full record. Status/CachedTree/WorkTree/Timestamp carry the raw
// fields for a human view; the board reduces them to its severity rows.
type GateInfo struct {
	Present        bool
	State          string
	PendingStatus  string
	Status         string
	CachedTree     string
	WorkTree       string
	Stale          bool
	Partition      *gate.Partition
	CheckPartition *gate.CheckPartition
	Timestamp      string
	Reason         string
	CacheBytes     int
}

// Query controls a status read without changing the board's shared signal ordering.
type Query struct {
	ExcludeDirtyPaths []string
}

// Signals gathers every ambient signal under root and returns them severity-sorted
// ascending — the one severity ladder `bench status` renders. render (the text board) and
// the dashboard gatherer both call this, so a signal added or reordered here reaches both
// surfaces from one source.
func Signals(root string) []Signal {
	return SignalsWith(root, Query{})
}

// SignalsWith gathers the ambient board under root using the supplied query.
func SignalsWith(root string, query Query) []Signal {
	var rows []row

	rows = appendGate(rows, root)
	rows = appendGit(rows, root, query)
	rows = appendWorktree(rows, root)
	rows = appendIntent(rows, root)
	rows = appendGuards(rows, root)
	rows = appendDrain(rows, root)
	rows = appendStructure(rows, root)
	rows = appendMaps(rows, root)
	rows = appendRetirement(rows, root)
	rows = appendOrphanedPickup(rows, root)
	rows = appendRoadmapReconcile(rows, root)
	rows = appendHandoff(rows, root)

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
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return line + "\n", code
	}
	_, all := parsed.Flags["--all"]
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
	if all {
		signals = expandIntentSignals(root, signals)
	}

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

// appendGate projects the gate owner's typed inspection onto the existing severity
// ladder. It does not parse, repair, or execute the oracle.
func appendGate(rows []row, root string) []row {
	return appendGateInfo(rows, GateVerdict(root), root)
}

func appendGateInfo(rows []row, gv GateInfo, root string) []row {
	if !gv.Present {
		return rows
	}
	if gv.State == string(gate.Pending) {
		if gv.PendingStatus == "locked-pending" {
			return append(rows, row{1, "gate", "locked-pending", "wait for live gate owner"})
		}
		return append(rows, row{2, "gate", "interrupted-pending", "re-run the gate"})
	}
	if gv.State == string(gate.Invalid) {
		return append(rows, row{3, "gate", "invalid verdict", "re-run the gate"})
	}
	if gv.State == string(gate.Unavailable) {
		return append(rows, row{3, "gate", "verdict unavailable", "inspect gate state"})
	}
	if gv.Status == "timeout" {
		return append(rows, row{0, "gate", "timeout", "inspect the hang and re-run the gate"})
	}
	if gv.Status == "red" {
		return append(rows, row{0, "gate", "red", "fix before commit"})
	}
	if gv.Stale {
		// An exact-tip narrow verdict stays distinct from drift even when the gate no
		// longer composes it as a whole-tree green. A moved tree is drift instead.
		if gv.CachedTree == gv.WorkTree && (gv.Partition != nil || gv.CheckPartition != nil) {
			detail := partialGreenDetail(gv.Partition, gv.CheckPartition)
			return append(rows, row{7, "gate", detail, "bench gate --fresh for a whole-tree verdict"})
		}
		detail, action := staleGateDetailAction(root, gv.CachedTree, gv.WorkTree)
		return append(rows, row{7, "gate", detail, action})
	}
	return rows
}

// GateVerdict adapts gate.Inspect for status, dashboard, and roadmap consumers. The
// reusable-green predicate and every cache classification remain owned by gate.
func GateVerdict(root string) GateInfo {
	in := gate.Inspect(root)
	if in.State == gate.Absent {
		return GateInfo{}
	}
	gi := GateInfo{Present: true, State: string(in.State), PendingStatus: in.PendingStatus, Status: in.Status, CachedTree: in.CachedTree, WorkTree: in.CurrentTree, Reason: in.Reason, CacheBytes: in.CacheBytes}
	if !in.RecordedAt.IsZero() {
		gi.Timestamp = in.RecordedAt.Format(time.RFC3339)
	}
	gi.Partition = in.Partition
	gi.CheckPartition = in.CheckPartition
	nonReusableGreen := in.State == gate.Ready && in.Status == "green" && !in.ReusableGreen
	gi.Stale = nonReusableGreen
	if nonReusableGreen && in.CachedTree == in.CurrentTree && gate.ComposedGreen(root) {
		gi.Stale = false
	}
	return gi
}

func staleGateDetailAction(root, cachedTree, currentTree string) (detail, action string) {
	return fmt.Sprintf("stale (gated tree %s, work tree %s)", Short(cachedTree), Short(currentTree)), "re-run the gate"
}

func skippedComponentNames(p *gate.Partition) []string {
	return componentNames(p.Skipped)
}

func partialGreenDetail(partition *gate.Partition, checks *gate.CheckPartition) string {
	var details []string
	if partition != nil {
		details = append(details, "skipped: "+strings.Join(skippedComponentNames(partition), ", "))
	}
	if checks != nil {
		details = append(details, "inherited checks: "+strings.Join(componentNames(checks.Inherited), ", "))
	}
	return "partial green (" + strings.Join(details, "; ") + ")"
}

func componentNames(components []gate.ComponentSkip) []string {
	names := make([]string, len(components))
	for i, component := range components {
		names[i] = component.Component
	}
	return names
}

// StepSeparator joins the steps of a board action that names a sequence rather than one
// command. It is exported as the one source of that join: a reader deciding whether an
// action is a single invocation recognizes the separator instead of restating it, so the
// producer and the recognizer cannot drift apart.
const StepSeparator = " / "

// appendGit adds the uncommitted/unpushed signal (sev 1). dirty is the porcelain status;
// ahead is the upstream-relative commit list, read only when an upstream is configured.
func appendGit(rows []row, root string, query Query) []row {
	fact, err := git.LandedState(root, query.ExcludeDirtyPaths...)
	if err != nil {
		return append(rows, row{1, "git", "git state unavailable", "investigate local git state"})
	}
	var details, actions []string
	if fact.DirtyPaths > 0 {
		details = append(details, Plural(fact.DirtyPaths, "dirty path", "dirty paths"))
		actions = append(actions, "commit on green")
	}
	if fact.UnpushedCommits > 0 {
		details = append(details, Plural(fact.UnpushedCommits, "unpushed commit", "unpushed commits"))
		actions = append(actions, "/bench-final-check")
	}
	if fact.UniqueBranches > 0 {
		details = append(details, Plural(fact.UniqueBranches, "unique branch", "unique branches"))
		actions = append(actions, "push")
	}
	if len(details) == 0 {
		return rows
	}
	return append(rows, row{1, "git", strings.Join(details, ", "), strings.Join(actions, StepSeparator)})
}

// objectiveDisplay resolves the human-facing objective text for an intent entry. The
// ledger no longer stores objective text; a live shift keeps its full objective only in
// <worktree>/.bench-objective (mode 0600), which exists for exactly as long as the entry is
// live. For a shift entry with a recorded worktree the renderer reads it back and previews
// it; when the file or the worktree is gone — the normal end state — it degrades to the
// entry key and never propagates the read error. Every other kind (claude-agent, worktree)
// has no such file and renders the key alone.
func objectiveDisplay(entry intent.Entry) string {
	if entry.Kind == intent.KindShift && entry.Worktree != "" {
		if data, err := os.ReadFile(filepath.Join(entry.Worktree, ".bench-objective")); err == nil {
			if text := strings.TrimRight(string(data), "\n"); text != "" {
				return sanitize.Preview(text)
			}
		}
	}
	return sanitize.Preview(entry.Key)
}

func appendIntent(rows []row, root string) []row {
	live, err := intent.Snapshot(root)
	if err != nil {
		return append(rows, row{2, "intent", "intent ledger unavailable", "inspect shared git intent ledger"})
	}
	if len(live) == 0 {
		return rows
	}
	correlated, uncorrelated := 0, 0
	for _, entry := range live {
		if entry.Worktree == "" && entry.Branch == "" && entry.Kind == intent.KindClaudeAgent {
			uncorrelated++
		} else {
			correlated++
		}
	}
	detail := fmt.Sprintf("%d correlated, %d uncorrelated; oldest: %s", correlated, uncorrelated, objectiveDisplay(live[0]))
	if r := live[0].Recovery; r != "" && r != shift.RecoveryNone {
		detail += "; recovery: " + sanitize.Preview(r)
	}
	return append(rows, row{2, "intent", detail, "resume interrupted work"})
}

func expandIntentSignals(root string, signals []Signal) []Signal {
	live, err := intent.Snapshot(root)
	if err != nil || len(live) == 0 {
		return signals
	}
	out := make([]Signal, 0, len(signals)+len(live))
	for _, signal := range signals {
		if signal.Name != "intent" {
			out = append(out, signal)
		}
	}
	for _, entry := range live {
		parts := []string{string(entry.Kind)}
		if entry.Worktree != "" {
			parts = append(parts, "path="+sanitize.Preview(entry.Worktree))
		}
		if entry.Branch != "" {
			parts = append(parts, "branch="+sanitize.Preview(entry.Branch))
		}
		if entry.Recovery != "" && entry.Recovery != shift.RecoveryNone {
			parts = append(parts, "recovery="+sanitize.Preview(entry.Recovery))
		}
		parts = append(parts, "objective="+objectiveDisplay(entry))
		out = append(out, Signal{Severity: 2, Name: "intent", Detail: strings.Join(parts, " "), Action: "resume interrupted work"})
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Severity < out[j].Severity })
	return out
}

// appendWorktree adds separate worktree signals (sev 2) for out-of-pool worktrees
// and leased pool entries. The repo root and warm pooled entries are expected state,
// not signals. Branch names are not ownership evidence and do not produce signals.
func appendWorktree(rows []row, root string) []row {
	registered, err := worktree.ClassifyRegisteredWorktrees(root)
	if err != nil {
		// A classify failure means the pool/leased/out-of-pool counts below are
		// unknowable, not zero: surface the git failure itself as the row rather than
		// falling through to an empty-looking board (the false-empty class FT29 swept).
		var typed git.WorktreeFailure
		if errors.As(err, &typed) {
			return append(rows, row{2, "worktree", typed.Error(), typed.WorktreeAction()})
		}
		return append(rows, row{2, "worktree", fmt.Sprintf("git worktree list failed: %v", err), "run git worktree list and retry"})
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
		rows = append(rows, row{2, "worktree", Plural(outOfPool, "out-of-pool worktree", "out-of-pool worktrees"), "inspect exact worktree (bench worktree clean <path>)"})
	}
	if leased > 0 {
		rows = append(rows, row{2, "worktree", Plural(leased, "leased pool worktree", "leased pool worktrees"), "resume leased worktree"})
	}
	return rows
}

// Plural renders a count with the unit its number takes. It is exported so a surface
// rendering board-derived counts states them in the board's own words rather than in a
// second phrasing of the same fact.
func Plural(n int, one, many string) string {
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
// double-report the same hook — and stays quiet when its managed bytes are current. Remedy: bench link.
func appendGuards(rows []row, root string) []row {
	if !isPrimaryCheckout(root) {
		return rows
	}
	if _, err := os.Stat(filepath.Join(root, ".bench", "lines.env")); err != nil {
		return rows
	}
	health := adopt.InspectPrePush(root)
	if health.State == adopt.PrePushManaged && health.Currency == adopt.PrePushCurrent {
		return rows
	}
	return append(rows, row{3, "guards", prePushDetail(health), "bench link"})
}

// prePushDetail names the pre-push gap the guards row reports, mirroring the adopt classifier's
// states so the ambient signal and the doctor row describe the same condition.
func prePushDetail(health adopt.PrePushHealth) string {
	if health.State == adopt.PrePushManaged && health.Currency == adopt.PrePushStale {
		return "pre-push stale"
	}
	switch health.State {
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
	common, err2 := git.CommonDir(root)
	if err1 != nil || err2 != nil {
		return false
	}
	return filepath.Clean(gitDir) == filepath.Clean(common)
}

// appendDrain adds the capture-drain signal (sev 4): parked ideas in capture/IDEAS.md plus open
// journal headings, one combined row pointing at the single maintenance phase. The
// counts are roadmap.DrainCounts — the same counters `bench roadmap` reports. The
// learnings component shows only at or above the floor (env BENCH_LEARNINGS_FLOOR,
// default 1); parked ideas always count.
//
// Each capture source carries its own readability state. A source whose read failed
// (the FileState.Failed test appendMaps also applies) renders as toon.UnknownCell's
// explicit `unknown (<path> is <state>)` segment instead of a fabricated 0, and the
// good source's count still renders alongside it: one source failing must not hide the
// other's number.
// The row only disappears when every source reads cleanly and every count is zero.
func appendDrain(rows []row, root string) []row {
	drain := roadmap.DrainCounts(root)
	ideas, ideasState, open, learningsState, retroCount, retrosState := drain.Ideas, drain.IdeasState, drain.OpenLearnings, drain.LearningsState, drain.Retros, drain.RetrosState
	if open < learningsFloor() {
		open = 0
	}
	ideasFailed := ideasState.Failed()
	learningsFailed := learningsState.Failed()
	retrosFailed := retrosState.Failed()
	if !ideasFailed && !learningsFailed && !retrosFailed {
		if ideas == 0 && open == 0 && retroCount == 0 {
			return rows
		}
		return append(rows, row{4, "drain", fmt.Sprintf("%d idea(s), %d open learning(s), %d pending retro(s)", ideas, open, retroCount), "/bench-what-next"})
	}
	var parts []string
	if ideasFailed {
		parts = append(parts, toon.UnknownCell(roadmap.IdeasFile, ideasState))
	} else {
		parts = append(parts, fmt.Sprintf("%d idea(s)", ideas))
	}
	if learningsFailed {
		parts = append(parts, toon.UnknownCell(learnings.JournalPath, learningsState))
	} else {
		parts = append(parts, fmt.Sprintf("%d open learning(s)", open))
	}
	if retrosFailed {
		parts = append(parts, toon.UnknownCell(retros.Directory+"/", retrosState))
	} else {
		parts = append(parts, fmt.Sprintf("%d pending retro(s)", retroCount))
	}
	return append(rows, row{4, "drain", strings.Join(parts, ", "), "/bench-what-next"})
}

// appendStructure adds the structural-debt signal (sev 5) when the violation count is positive.
func appendStructure(rows []row, root string) []row {
	if n := structure.ViolationCount(root); n > 0 {
		return append(rows, row{5, "structure", fmt.Sprintf("%d issue(s)", n), "split (craft-seams)"})
	}
	return rows
}

// appendMaps adds the unresolved-decision-map signal (sev 6): a positive count when
// the scan ran cleanly, or an explicit unknown row naming the decisions/ read failure
// when it did not — a scan that could not run must never render as zero unresolved.
func appendMaps(rows []row, root string) []row {
	n, state := maps.UnresolvedCount(root)
	if state.Failed() {
		return append(rows, row{6, "decisions", toon.UnknownCell(maps.DecisionsDir, state), "investigate decisions/ (bench maps)"})
	}
	if n > 0 {
		return append(rows, row{6, "decisions", fmt.Sprintf("%d unresolved map(s)", n), "/bench-shape-idea"})
	}
	return rows
}

// appendRetirement adds the merged-spec-awaiting-retirement signal (sev 8), but only on
// the default branch — a topic branch's spec is still in flight, not awaiting retirement.
func appendRetirement(rows []row, root string) []row {
	// Audit #5 — tolerate: an unreadable branch or an unresolvable default reads as "not
	// the default branch", skipping this advisory housekeeping signal; non-fatal on the
	// ambient board.
	cur, _ := git.Output("-C", root, "rev-parse", "--abbrev-ref", "HEAD")
	if def, ok := git.ResolvedDefault(root); !ok || cur != def {
		return rows
	}
	if n := retirementCount(root); n > 0 {
		return append(rows, row{8, "specs", fmt.Sprintf("%d merged spec(s) awaiting retirement", n), "bench spec retire <slug>"})
	}
	return rows
}

// appendOrphanedPickup adds the orphaned-review-pickup signal (sev 9): a reviews/<slug>.md
// with no matching live specs/<slug>/spec.md is a pickup file that
// escaped its lifecycle. It ranks
// with the housekeeping rows, just below the retirement signal, so it never displaces the
// gate/git rows in the budget. A paired pickup (its spec still present) is expected state.
func appendOrphanedPickup(rows []row, root string) []row {
	if n := orphanedPickupCount(root); n > 0 {
		return append(rows, row{9, "reviews", Plural(n, "orphaned review pickup", "orphaned review pickups"), "promote or delete by hand"})
	}
	return rows
}

// appendRoadmapReconcile adds the roadmap-reconcile signal (sev 9): a ROADMAP.md row naming a
// specs/<slug>/spec.md whose work has already shipped — the spec file is merged-implemented
// (spec.AwaitsRetirement) or was retired out of the tree entirely — has outlived the drain that
// should have removed it. Like appendRetirement, it fires only on the default branch: a topic
// branch's roadmap is mid-build, so a row there names in-flight work, not a shipped-work leak.
// Severity 10 ranks it below the housekeeping rows (retirement 8, orphaned-pickup 9) and far
// below gate/git, so it never displaces a red-gate or dirty-tree row in the budget.
func appendRoadmapReconcile(rows []row, root string) []row {
	// Audit #6 — tolerate: as in appendRetirement, an unreadable branch or an unresolvable
	// default reads as "not the default branch" and skips this advisory reconcile signal.
	cur, _ := git.Output("-C", root, "rev-parse", "--abbrev-ref", "HEAD")
	if def, ok := git.ResolvedDefault(root); !ok || cur != def {
		return rows
	}
	merged, dangling, state := roadmapReconcileCounts(root)
	if state.Failed() {
		return append(rows, row{10, "roadmap", toon.UnknownCell(roadmap.RoadmapFile, state), "/bench-what-next"})
	}
	if merged == 0 && dangling == 0 {
		return rows
	}
	var details []string
	if merged > 0 {
		details = append(details, Plural(merged, "row for merged work", "rows for merged work"))
	}
	if dangling > 0 {
		details = append(details, Plural(dangling, "row names a retired spec", "rows name a retired spec"))
	}
	return append(rows, row{10, "roadmap", strings.Join(details, ", "), "/bench-what-next"})
}

// retirementCount counts live folder specs that spec.AwaitsRetirement marks — a merged spec
// awaiting retirement. Absent `specs/` → 0. The unfenced-marker predicate is
// spec.AwaitsRetirement's one source. Every read goes through the classifier: this signal
// is advisory and stays quiet about a spec it could not read, but a FIFO parked at
// specs/*/spec.md never yields EOF, and a board that blocks forever is worse than one missing a
// housekeeping row — `bench status` is what the SessionStart hook runs.
func retirementCount(root string) int {
	dir := filepath.Join(root, "specs")
	cd := bounds.ClassifyDir(dir)
	if cd.State != bounds.StateParsed {
		return 0
	}
	n := 0
	for _, e := range cd.Entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		if !e.IsDir() {
			continue
		}
		path = filepath.Join(path, "spec.md")
		c := bounds.Classify(path, bounds.ControlRecordLimit)
		if c.State != bounds.StateParsed {
			continue
		}
		if spec.AwaitsRetirement(c.Data) {
			n++
		}
	}
	return n
}

// orphanedPickupCount counts reviews/*.md files with no matching live specs/<slug>/spec.md — a review
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
		slug := strings.TrimSuffix(e.Name(), ".md")
		folder := filepath.Join(root, spec.LiveSpecPath(slug))
		if info, err := os.Stat(folder); err == nil && !info.IsDir() {
			continue // paired: its folder spec is still present
		}
		n++
	}
	return n
}

// roadmapReconcileCounts scans ROADMAP.md for live spec-path tokens and classifies each
// distinct path against the tree: a missing file is a dangling row (the spec retired but its
// roadmap row survived); a present file that spec.AwaitsRetirement marks is a merged row (the
// work shipped but the drain missed it). A present, still-staged spec is the normal open-work
// state and counts nothing. Absent or empty ROADMAP.md → 0, 0, bounds.StateParsed, the ordinary
// quiet-roadmap posture; a ROADMAP.md whose read failed reports that state instead, so
// appendRoadmapReconcile renders the failed read as unknown rather than a fabricated clean
// board. Each named live spec goes through the classifier too, so a FIFO cannot
// block the board; a spec path that yields no content at all is the dangling case, whether
// nothing is there or what is there could not be read. The merged predicate is
// spec.AwaitsRetirement, the same one source the retirement counter applies.
func roadmapReconcileCounts(root string) (merged, dangling int, state bounds.FileState) {
	c := bounds.Classify(filepath.Join(root, roadmap.RoadmapFile), bounds.ControlRecordLimit)
	switch {
	case c.State.Failed():
		return 0, 0, c.State
	case c.State == bounds.StateAbsent || c.State == bounds.StateEmpty:
		return 0, 0, bounds.StateParsed
	}
	for _, slug := range spec.LiveSpecSlugs(c.Data) {
		folderPath := spec.LiveSpecPath(slug)
		sc := bounds.Classify(filepath.Join(root, folderPath), bounds.ControlRecordLimit)
		if sc.State == bounds.StateAbsent || sc.State.Failed() {
			dangling++
			continue
		}
		if spec.AwaitsRetirement(sc.Data) {
			merged++
		}
	}
	return merged, dangling, bounds.StateParsed
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

// Short returns the first up-to-7 bytes of s (the shell `${var:0:7}` tree-prefix slice),
// guarding a short or "none" hash so the slice never panics. It is exported as the one
// source of the prefix width, so every surface that renders a tree or commit reference
// cuts it at the same place.
func Short(s string) string {
	return s[:min(7, len(s))]
}
