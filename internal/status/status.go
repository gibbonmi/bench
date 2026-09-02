// Package status ports `bench status`: the ambient dashboard a session-start hook
// consumes verbatim. It is the composition crux of the slice. Every sibling query
// package, maps, structure, worktree, roadmap, contributes one signal. This package
// renders them into the single severity-sorted board the shell renderer produced.
// The output format is stable because a hook parses it: the lead line, the
// fixed-width rows, and the `+N more` overflow.
//
// One rule per signal lives in its own package; status only orders the signals by
// severity and formats them. The specs housekeeping signals are counted here.
// retirementCount scans specs/ and applies the shared spec.AwaitsRetirement predicate.
// orphanedPickupCount pairs reviews/ against specs/. roadmapReconcileCounts classifies
// ROADMAP.md's spec-path tokens against the tree. The merged-implemented predicate
// itself is spec.AwaitsRetirement, one source shared with `bench spec retire`.
package status

import (
	"bytes"
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
	"github.com/gibbonmi/bench/internal/census"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/landing"
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

// grammar is the declared argument shape usage.Parse enforces for this subcommand.
// Arity, flag recognition, `--`, and help all come from there rather than a local switch.
var grammar = usage.Grammar{
	Cmd:  "bench status",
	Help: "usage: bench status [--all] [--route [--harness " + HarnessChoices() + "]]",
	Flags: []usage.Flag{
		{Name: "--all"},
		{Name: "--route"},
		{Name: "--harness", HasValue: true, NoEmptyValue: true},
	},
}

// row is one dashboard signal: a severity (the sort/lead key), and the signal/detail/
// action triple the shell packed into a `sev|signal|detail|action` line.
type row struct {
	sev            int
	signal, detail string
	action         statusAction
}

type actionKind uint8

// StepSeparator joins the steps of a board action that names a sequence rather than one
// command. The action grammar rejects it, so sequences never become invocable routes.
const StepSeparator = " / "

const (
	actionNone actionKind = iota
	actionPhase
	actionBench
	actionGit
)

type actionID uint8

const (
	noAction actionID = iota
	setupAction
	gateAction
	freshGateAction
	statusAllAction
	benchWorktreeListAction
	cleanWorktreeAction
	cleanUnclaimedWorktreeAction
	linkAction
	mapsAction
	roadmapAction
	structureAction
	retireSpecAction
	closeTicketsAction
	handoffAction
	gitPushAction
	gitStatusAction
	gitWorktreeListAction
	debugPhaseAction
	drainPhaseAction
	finalCheckPhaseAction
	implementSpecPhaseAction
	shapeIdeaPhaseAction
	writeSpecPhaseAction
	otherPhaseAction
	actionCount
)

type argumentShape uint8

const (
	noArgument argumentShape = iota
	oneWordArgument
	optionalSpecPath
	optionalDecisionPath
	anyPhaseCommand
)

type actionDefinition struct {
	kind     actionKind
	command  string
	argument argumentShape
}

var actionDefinitions = [actionCount]actionDefinition{
	setupAction:                  {kind: actionBench, command: "bench setup"},
	gateAction:                   {kind: actionBench, command: "bench gate"},
	freshGateAction:              {kind: actionBench, command: "bench gate --fresh"},
	statusAllAction:              {kind: actionBench, command: "bench status --all"},
	benchWorktreeListAction:      {kind: actionBench, command: "bench worktree list"},
	cleanWorktreeAction:          {kind: actionBench, command: "bench worktree clean", argument: oneWordArgument},
	cleanUnclaimedWorktreeAction: {kind: actionBench, command: "bench worktree clean --discard-branch --unclaimed --apply-current"},
	linkAction:                   {kind: actionBench, command: "bench link"},
	mapsAction:                   {kind: actionBench, command: "bench maps"},
	roadmapAction:                {kind: actionBench, command: "bench roadmap"},
	structureAction:              {kind: actionBench, command: "bench structure"},
	retireSpecAction:             {kind: actionBench, command: "bench spec retire", argument: oneWordArgument},
	closeTicketsAction:           {kind: actionBench, command: "bench worktree land --spec", argument: oneWordArgument},
	handoffAction:                {kind: actionBench, command: "bench handoff"},
	gitPushAction:                {kind: actionGit, command: "git push"},
	gitStatusAction:              {kind: actionGit, command: "git status"},
	gitWorktreeListAction:        {kind: actionGit, command: "git worktree list"},
	debugPhaseAction:             {kind: actionPhase, command: "/bench-debug"},
	drainPhaseAction:             {kind: actionPhase, command: "/bench-drain"},
	finalCheckPhaseAction:        {kind: actionPhase, command: "/bench-final-check"},
	implementSpecPhaseAction: {
		kind: actionPhase, command: "/bench-implement-spec", argument: optionalSpecPath,
	},
	shapeIdeaPhaseAction: {kind: actionPhase, command: "/bench-shape-idea"},
	writeSpecPhaseAction: {
		kind: actionPhase, command: "/bench-write-spec", argument: optionalDecisionPath,
	},
	otherPhaseAction: {kind: actionPhase, argument: anyPhaseCommand},
}

type statusAction struct {
	id       actionID
	argument string
	advisory string
}

func parseAction(text string) statusAction {
	if strings.Contains(text, StepSeparator) {
		return statusAction{advisory: text}
	}
	for id := actionID(1); id < actionCount; id++ {
		if argument, ok := actionDefinitions[id].match(text); ok {
			return statusAction{id: id, argument: argument}
		}
	}
	return statusAction{advisory: text}
}

func (definition actionDefinition) match(text string) (string, bool) {
	if definition.command == "" && definition.argument != anyPhaseCommand {
		return "", false
	}
	switch definition.argument {
	case noArgument:
		if definition.kind == actionPhase {
			return "", text == definition.command
		}
		return "", strings.HasPrefix(text, strings.Fields(definition.command)[0]+" ") &&
			strings.Join(strings.Fields(text), " ") == definition.command
	case oneWordArgument:
		parts := strings.Fields(text)
		command := strings.Fields(definition.command)
		if !strings.HasPrefix(text, command[0]+" ") || len(parts) != len(command)+1 {
			return "", false
		}
		for i := range command {
			if parts[i] != command[i] {
				return "", false
			}
		}
		return parts[len(parts)-1], true
	case optionalSpecPath:
		return matchOptionalPath(text, definition.command, "specs/", "/spec.md")
	case optionalDecisionPath:
		return matchOptionalPath(text, definition.command, "decisions/", ".md")
	case anyPhaseCommand:
		command, argument, hasArgument := strings.Cut(text, " ")
		if command == harnessPrefix[HarnessClaude] || !strings.HasPrefix(command, harnessPrefix[HarnessClaude]) ||
			strings.Contains(strings.TrimPrefix(command, harnessPrefix[HarnessClaude]), "/") {
			return "", false
		}
		if hasArgument && !((strings.HasPrefix(argument, "specs/") && strings.HasSuffix(argument, "/spec.md")) ||
			(strings.HasPrefix(argument, "decisions/") && strings.HasSuffix(argument, ".md"))) {
			return "", false
		}
		return text, true
	}
	return "", false
}

func matchOptionalPath(text, command, prefix, suffix string) (string, bool) {
	if text == command {
		return "", true
	}
	argument, ok := strings.CutPrefix(text, command+" ")
	if !ok || !strings.HasPrefix(argument, prefix) || !strings.HasSuffix(argument, suffix) {
		return "", false
	}
	return argument, true
}

func commandAction(id actionID) statusAction {
	return statusAction{id: id}
}

func commandActionWithArgument(id actionID, argument string) statusAction {
	return statusAction{id: id, argument: argument}
}

func advisoryAction(text string) statusAction {
	return statusAction{advisory: text}
}

func (action statusAction) render() string {
	if action.id == noAction {
		return action.advisory
	}
	definition := actionDefinitions[action.id]
	if definition.argument == anyPhaseCommand {
		return action.argument
	}
	if action.argument == "" {
		return definition.command
	}
	return definition.command + " " + action.argument
}

func (action statusAction) invocable() bool {
	return action.id.kind() != actionNone
}

func (id actionID) kind() actionKind {
	if id <= noAction || id >= actionCount {
		return actionNone
	}
	return actionDefinitions[id].kind
}

// IsInvocable reports whether action is one command accepted by the status action grammar.
func IsInvocable(text string) bool {
	return parseAction(text).invocable()
}

// Signal is one ambient-board row exposed as structured data — the severity sort key
// plus the signal/detail/action triple. Action remains the rendered string contract while
// actionID carries the producer's definition into routing. It is the shared
// accessor the text board and the human dashboard both consume, so the two views cannot
// rank or drop signals differently.
type Signal struct {
	Severity int
	Name     string
	Detail   string
	Action   string
	actionID actionID
}

func newSignal(severity int, name, detail string, action statusAction) Signal {
	return Signal{Severity: severity, Name: name, Detail: detail, Action: action.render(), actionID: action.id}
}

func (s Signal) invocable() bool {
	return s.actionID.kind() != actionNone
}

// GateInfo is the structured gate-verdict cache read, shared by the status board and the
// dashboard so neither re-parses `<git-dir>/bench-last-gate`. Present is false when no
// cache file exists. Stale marks a ready verdict the gate no longer stands behind for
// this subject: a record whose tree or oracle has drifted. It also covers an exact-tip
// non-reusable green the gate cannot compose as a whole-tree verdict. Partition carries
// narrowness at component granularity: non-nil for a partial verdict that graded only
// the components whose inputs moved, nil for a full record.
// Status/CachedTree/WorkTree/Timestamp carry the raw fields for a human view; the board
// reduces them to its severity rows.
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
// ascending, the one severity ladder `bench status` renders. render, the text board, and
// the dashboard gatherer both call this. A signal added or reordered here reaches both
// surfaces from one source.
func Signals(root string) []Signal {
	return SignalsWith(root, Query{})
}

// SignalsWith gathers the ambient board under root using the supplied query. It
// resolves the Bench home at this boundary and passes it down, so no signal owner
// reads the environment itself.
func SignalsWith(root string, query Query) []Signal {
	return signalsWith(root, query, worktree.Home())
}

// signalsWith is SignalsWith with the Bench home already resolved. The render path
// resolves the home once and shares it with the --all expanders.
func signalsWith(root string, query Query, home string) []Signal {
	var rows []row

	rows = appendSetup(rows, root)
	rows = appendGate(rows, root)
	rows = appendGit(rows, root, query)
	rows = appendWorktree(rows, root)
	rows = appendIntent(rows, root)
	rows = appendGuards(rows, root)
	rows = appendCensus(rows, root, home)
	rows = appendStagedSpecs(rows, root)
	rows = appendDrain(rows, root)
	rows = appendStructure(rows, root)
	rows = appendMaps(rows, root)
	rows = appendRetirement(rows, root)
	rows = appendOrphanedPickup(rows, root)
	rows = appendRoadmapReconcile(rows, root)
	rows = appendTicketsOnly(rows, root)
	rows = appendHandoff(rows, root)

	// Ascending numeric sort by severity; each severity is unique, so ordering is
	// fully determined and the min-severity row leads.
	sort.SliceStable(rows, func(i, j int) bool { return rows[i].sev < rows[j].sev })

	out := make([]Signal, len(rows))
	for i, r := range rows {
		out[i] = newSignal(r.sev, r.signal, r.detail, r.action)
	}
	return out
}

func appendSetup(rows []row, root string) []row {
	info, err := os.Stat(filepath.Join(root, ".bench"))
	if err != nil {
		if !os.IsNotExist(err) {
			return rows
		}
		return append(rows, row{0, "setup", "no .bench/", commandAction(setupAction)})
	}
	if !info.IsDir() {
		return append(rows, row{0, "setup", "no .bench/", commandAction(setupAction)})
	}
	return rows
}

func appendStagedSpecs(rows []row, root string) []row {
	n, slug := stagedSpecCount(root)
	if n == 0 {
		return rows
	}
	command := commandAction(implementSpecPhaseAction)
	if n == 1 {
		command = commandActionWithArgument(implementSpecPhaseAction, "specs/"+slug+"/spec.md")
	}
	return append(rows, row{4, "specs", fmt.Sprintf("%d staged spec(s)", n), command})
}

func stagedSpecCount(root string) (int, string) {
	facts, err := spec.Facts(root)
	if err != nil {
		return 0, ""
	}
	n, slug := 0, ""
	for _, fact := range facts {
		if fact.Status != "staged" {
			continue
		}
		n++
		slug = fact.Slug
	}
	return n, slug
}

// Command implements `bench status`. It renders the ambient board by default, its full
// form with --all, or the one-row route for one named harness.
func Command(args []string) (string, int) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		if code == 2 && (usage.FlagPresent(grammar, args, "--route") || usage.FlagPresent(grammar, args, "--harness")) {
			return grammar.Help + "\n", code
		}
		return line + "\n", code
	}
	_, all := parsed.Flags["--all"]
	_, route := parsed.Flags["--route"]
	harness := HarnessClaude
	if value, present := parsed.Flags["--harness"]; present {
		if !route || !ValidHarness(value) {
			return grammar.Help + "\n", 2
		}
		harness = value
	}
	if all && route {
		return grammar.Help + "\n", 2
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	if route {
		return renderRoute(RouteFor(root, Signals(root), harness))
	}
	return render(root, all), 0
}

func renderRoute(route RouteResult) (string, int) {
	out, err := toon.Table("next", []string{"state", "why", "command"}, [][]string{{
		sanitize.Controls(route.Lead.Name),
		sanitize.Controls(route.Lead.Detail),
		sanitize.Controls(route.Lead.Action),
	}})
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	var b strings.Builder
	b.WriteString(out)
	b.WriteString("also: ")
	if len(route.RunnersUp) == 0 {
		b.WriteString("none\n")
		return b.String(), 0
	}
	for i, runnerUp := range route.RunnersUp {
		if i > 0 {
			b.WriteString("; ")
		}
		fmt.Fprintf(&b, "%s (%s) → %s", sanitize.Controls(runnerUp.Name), sanitize.Controls(runnerUp.Detail), sanitize.Controls(runnerUp.Action))
	}
	b.WriteByte('\n')
	return b.String(), 0
}

// render gathers every signal under root, sorts ascending by severity, and formats the
// board. This is the byte-for-byte counterpart of the shell `status()` renderer. When all
// is false it applies the five-row budget and appends the overflow line. When all is true,
// `bench status --all`, it prints every row and emits no overflow line. The SessionStart
// hook calls with all=false so the ambient surface stays bounded.
func render(root string, all bool) string {
	home := worktree.Home()
	signals := signalsWith(root, Query{}, home)
	if all {
		signals = expandIntentSignals(root, signals)
		signals = expandCensusSignals(root, home, signals)
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
			return append(rows, row{1, "gate", "locked-pending", advisoryAction("")})
		}
		return append(rows, row{2, "gate", "interrupted-pending", commandAction(gateAction)})
	}
	if gv.State == string(gate.Invalid) {
		return append(rows, row{3, "gate", "invalid verdict", commandAction(gateAction)})
	}
	if gv.State == string(gate.Unavailable) {
		return append(rows, row{3, "gate", "verdict unavailable", commandAction(freshGateAction)})
	}
	// Staleness outranks the verdict the record carries. A red the gate has retired names
	// a tree the reader has left. "Fix before commit" would then send them after work that
	// is no longer in the tree. Only a verdict the gate still stands behind reaches the
	// rows below.
	if gv.Stale {
		// An exact-tip narrow verdict stays distinct from drift even when the gate no
		// longer composes it as a whole-tree green. A moved tree is drift instead.
		if gv.CachedTree == gv.WorkTree && (gv.Partition != nil || gv.CheckPartition != nil) {
			detail := partialGreenDetail(gv.Partition, gv.CheckPartition)
			return append(rows, row{7, "gate", detail, commandAction(freshGateAction)})
		}
		detail, command := staleGateDetailAction(root, gv.CachedTree, gv.WorkTree, gv.Reason)
		return append(rows, row{7, "gate", detail, command})
	}
	if gv.Status == "timeout" {
		return append(rows, row{0, "gate", "timeout", commandAction(freshGateAction)})
	}
	if gv.Status == "red" {
		return append(rows, row{0, "gate", "red", commandAction(debugPhaseAction)})
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
	// A drifted record is stale whatever verdict it carries: the gate has already decided it
	// describes another tree or another oracle. A green the gate will not reuse is stale on
	// the reuse rule alone. A composed whole-tree green can only rescue the second case: it
	// holds exactly when the record matches the subject, which is when drift is absent.
	nonReusableGreen := in.State == gate.Ready && in.Status == "green" && !in.ReusableGreen
	gi.Stale = nonReusableGreen || (in.State == gate.Ready && in.Drifted)
	if gi.Stale && in.CachedTree == in.CurrentTree && gate.ComposedGreen(root) {
		gi.Stale = false
	}
	return gi
}

// staleGateDetailAction states the staleness the reader must act on. The gate's reason
// joins the two trees when it has one: matching trees otherwise read as a staleness with
// no cause.
func staleGateDetailAction(root, cachedTree, currentTree, reason string) (detail string, action statusAction) {
	detail = fmt.Sprintf("stale (gated tree %s, work tree %s", Short(cachedTree), Short(currentTree))
	if reason != "" {
		detail += "; " + reason
	}
	return detail + ")", commandAction(gateAction)
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

// appendGit adds the uncommitted/unpushed signal (sev 1). dirty is the porcelain status;
// ahead is the upstream-relative commit list, read only when an upstream is configured.
func appendGit(rows []row, root string, query Query) []row {
	fact, err := git.LandedState(root, query.ExcludeDirtyPaths...)
	if err != nil {
		return append(rows, row{1, "git", "git state unavailable", commandAction(gitStatusAction)})
	}
	var details []string
	unclaimedUniqueBranches, claimedUniqueBranches := 0, 0
	uniqueRefs := make(map[string]bool, len(fact.UniqueBranchNames))
	for _, branch := range fact.UniqueBranchNames {
		uniqueRefs["refs/heads/"+branch] = true
	}
	if fact.DirtyPaths > 0 {
		details = append(details, Plural(fact.DirtyPaths, "dirty path", "dirty paths"))
	}
	if fact.UnpushedCommits > 0 {
		details = append(details, Plural(fact.UnpushedCommits, "unpushed commit", "unpushed commits"))
	}
	unclaimedRefs, unclaimedErr := worktree.UnclaimedAssignmentBranchRefs(root)
	if unclaimedErr == nil {
		for _, ref := range unclaimedRefs {
			if uniqueRefs[ref] {
				unclaimedUniqueBranches++
			}
		}
	}
	if assignments, err := intent.Assignments(root); err == nil {
		for _, assignment := range assignments {
			if uniqueRefs[assignment.Branch] {
				claimedUniqueBranches++
			}
		}
	}
	if unclaimedUniqueBranches > 0 {
		details = append(details, Plural(unclaimedUniqueBranches, "unclaimed assignment branch", "unclaimed assignment branches"))
	}
	ordinaryUniqueBranches := fact.UniqueBranches - unclaimedUniqueBranches - claimedUniqueBranches
	if ordinaryUniqueBranches > 0 {
		details = append(details, Plural(ordinaryUniqueBranches, "unique branch", "unique branches"))
	}
	if len(details) == 0 {
		return rows
	}
	command := commandAction(gitPushAction)
	if fact.DirtyPaths > 0 {
		command = commandAction(finalCheckPhaseAction)
	}
	if unclaimedUniqueBranches > 0 {
		command = commandAction(cleanUnclaimedWorktreeAction)
	}
	return append(rows, row{1, "git", strings.Join(details, ", "), command})
}

// objectiveDisplay resolves the human-facing objective text for an intent entry. The
// ledger no longer stores objective text. A live shift keeps its full objective only in
// <worktree>/.bench-objective (mode 0600), which exists for exactly as long as the entry
// is live. For a shift entry with a recorded worktree the renderer reads it back and
// previews it. When the file or the worktree is gone, the normal end state, it degrades
// to the entry key and never propagates the read error. Every other kind (claude-agent,
// worktree) has no such file and renders the key alone.
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
		return append(rows, row{2, "intent", "intent ledger unavailable", commandAction(statusAllAction)})
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
	return append(rows, row{2, "intent", detail, commandAction(statusAllAction)})
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
		out = append(out, newSignal(2, "intent", strings.Join(parts, " "), commandAction(statusAllAction)))
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
		// unknowable, not zero. It surfaces the git failure itself as the row, rather
		// than falling through to an empty-looking board, the false-empty class FT29 swept.
		var typed git.WorktreeFailure
		if errors.As(err, &typed) {
			return append(rows, row{2, "worktree", typed.Error(), commandAction(benchWorktreeListAction)})
		}
		return append(rows, row{2, "worktree", fmt.Sprintf("git worktree list failed: %v", err), commandAction(gitWorktreeListAction)})
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
		rows = append(rows, row{2, "worktree", Plural(outOfPool, "out-of-pool worktree", "out-of-pool worktrees"), commandActionWithArgument(cleanWorktreeAction, "<path>")})
	}
	if leased > 0 {
		rows = append(rows, row{2, "worktree", Plural(leased, "leased pool worktree", "leased pool worktrees"), commandAction(benchWorktreeListAction)})
	}
	return rows
}

// Plural renders a count with the unit its number takes. It is exported so a surface
// rendering board-derived counts states them in the board's own words. That avoids a
// second phrasing of the same fact.
func Plural(n int, one, many string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s", one)
	}
	return fmt.Sprintf("%d %s", n, many)
}

// appendGuards adds the pre-push backstop signal (sev 3), ranked just below the worktree
// signals and above the drain row. git does not clone hooks, so a fresh clone silently
// loses the harness-independent default-branch backstop. This surfaces the gap ambiently
// rather than only under `bench doctor`. It fires only on the primary checkout of a
// routed repo, where `.bench/lines.env` is present. Pool and linked worktrees share the
// main `.git` and must not double-report the same hook. It stays quiet when its managed
// bytes are current; the remedy is bench link.
func appendGuards(rows []row, root string) []row {
	primary, err := git.IsPrimaryCheckout(root)
	if err != nil || !primary {
		return rows
	}
	if _, err := os.Stat(filepath.Join(root, ".bench", "lines.env")); err != nil {
		return rows
	}
	health := adopt.InspectPrePush(root)
	if health.State == adopt.PrePushManaged && health.Currency == adopt.PrePushCurrent {
		return rows
	}
	return append(rows, row{3, "guards", prePushDetail(health), commandAction(linkAction)})
}

// appendCensus adds the raw-call signal (sev 3), which ranks beside the guards row. It
// joins the census counts to the ledger's active assignments, so a count with no active
// entry names a worktree the reviewer has already released and never reaches the board.
// The row names no action, because no command is the remedy: the drain reads the heads
// and proposes the Bench form.
func appendCensus(rows []row, root, home string) []row {
	calls, worktrees := censusTotals(root, home)
	if calls == 0 {
		return rows
	}
	detail := Plural(calls, "raw call", "raw calls") + " across " + Plural(worktrees, "worktree", "worktrees")
	return append(rows, row{3, "census", detail, advisoryAction("")})
}

// censusTotals returns the raw-call count and the number of active assignments that
// hold at least one record. A census read failure and a ledger read failure both
// return zero, because the census is ambient evidence and never a board failure.
func censusTotals(root, home string) (calls, worktrees int) {
	for _, entry := range activeCensusCounts(root, home) {
		calls += entry.count
		worktrees++
	}
	return calls, worktrees
}

// censusEntry is one active assignment's identity, sanitized label, and record count.
type censusEntry struct {
	id, label string
	count     int
}

// activeCensusCounts pairs each active assignment with its record count, and drops
// every pair whose count is zero. The row and its --all expansion read the same join,
// so the sum and the per-worktree lines cannot disagree. The result is ordered by
// label, then by assignment id, so two worktrees that share a label stay two rows.
func activeCensusCounts(root, home string) []censusEntry {
	counts, err := census.Counts(home, root)
	if err != nil || len(counts) == 0 {
		return nil
	}
	assignments, err := intent.Assignments(root)
	if err != nil {
		return nil
	}
	var active []censusEntry
	for _, assignment := range assignments {
		if assignment.State != intent.StateActive {
			continue
		}
		if count := counts[assignment.ID]; count > 0 {
			active = append(active, censusEntry{assignment.ID, sanitize.Controls(assignment.Label), count})
		}
	}
	sort.SliceStable(active, func(i, j int) bool {
		if active[i].label != active[j].label {
			return active[i].label < active[j].label
		}
		return active[i].id < active[j].id
	})
	return active
}

// expandCensusSignals replaces the summed census row with one row per active
// assignment. It is the sibling of expandIntentSignals: both run only under --all, so
// the default board keeps its five-row budget however many worktrees are live.
func expandCensusSignals(root, home string, signals []Signal) []Signal {
	active := activeCensusCounts(root, home)
	if len(active) == 0 {
		return signals
	}
	out := make([]Signal, 0, len(signals)+len(active))
	for _, signal := range signals {
		if signal.Name != "census" {
			out = append(out, signal)
		}
	}
	for _, entry := range active {
		detail := entry.label + " " + Plural(entry.count, "raw call", "raw calls")
		out = append(out, newSignal(3, "census", detail, advisoryAction("")))
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Severity < out[j].Severity })
	return out
}

// prePushDetail names the pre-push gap the guards row reports. It mirrors the adopt
// classifier's states, so the ambient signal and the doctor row describe the same
// condition.
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

// appendDrain adds the capture-drain signal (sev 4): parked ideas in capture/IDEAS.md plus open
// journal headings, one combined row pointing at the single maintenance phase. The
// counts are roadmap.DrainCounts — the same counters `bench roadmap` reports. The
// learnings component shows only at or above the floor (env BENCH_LEARNINGS_FLOOR,
// default 1); parked ideas always count.
//
// Each capture source carries its own readability state. A source whose read failed,
// the FileState.Failed test appendMaps also applies, renders as toon.UnknownCell's
// explicit `unknown (<path> is <state>)` segment instead of a fabricated 0. The good
// source's count still renders alongside it: one source failing must not hide the
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
		return append(rows, row{4, "drain", fmt.Sprintf("%d idea(s), %d open learning(s), %d pending retro(s)", ideas, open, retroCount), commandAction(drainPhaseAction)})
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
	return append(rows, row{4, "drain", strings.Join(parts, ", "), commandAction(drainPhaseAction)})
}

// appendStructure adds the structural-debt signal (sev 5) when the violation count is positive.
func appendStructure(rows []row, root string) []row {
	if n := structure.ViolationCount(root); n > 0 {
		return append(rows, row{5, "structure", fmt.Sprintf("%d issue(s)", n), commandAction(structureAction)})
	}
	return rows
}

// appendMaps adds the unresolved-decision-map signal (sev 6): a positive count when
// the scan ran cleanly. It adds an explicit unknown row naming the decisions/ read
// failure when the scan did not run cleanly. A scan that could not run must never
// render as zero unresolved.
func appendMaps(rows []row, root string) []row {
	n, ready, state := maps.ActiveCounts(root)
	if state.Failed() {
		return append(rows, row{6, "decisions", toon.UnknownCell(maps.DecisionsDir, state), commandAction(mapsAction)})
	}
	if n > 0 {
		return append(rows, row{6, "decisions", fmt.Sprintf("%d unresolved map(s)", n), commandAction(shapeIdeaPhaseAction)})
	}
	if ready > 0 {
		command := commandAction(writeSpecPhaseAction)
		if ready == 1 {
			candidates, err := maps.DiscoverDecisionMapCandidates(root)
			if err != nil {
				return rows
			}
			for _, candidate := range candidates {
				if !candidate.Compiled {
					command = commandActionWithArgument(writeSpecPhaseAction, candidate.Path)
					break
				}
			}
		}
		return append(rows, row{6, "decisions", fmt.Sprintf("%d ready map(s)", ready), command})
	}
	return rows
}

// appendRetirement adds the merged-spec-awaiting-retirement signal (sev 8), but only on
// the default branch — a topic branch's spec is still in flight, not awaiting retirement.
func appendRetirement(rows []row, root string) []row {
	// Audit #5: tolerate. An unreadable branch or an unresolvable default reads as "not
	// the default branch". This skips this advisory housekeeping signal; it is non-fatal
	// on the ambient board.
	cur, _ := git.Output("-C", root, "rev-parse", "--abbrev-ref", "HEAD")
	if def, ok := git.ResolvedDefault(root); !ok || cur != def {
		return rows
	}
	if n := retirementCount(root); n > 0 {
		return append(rows, row{8, "specs", fmt.Sprintf("%d merged spec(s) awaiting retirement", n), commandActionWithArgument(retireSpecAction, "<slug>")})
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
		return append(rows, row{9, "reviews", Plural(n, "orphaned review pickup", "orphaned review pickups"), advisoryAction("")})
	}
	return rows
}

// appendRoadmapReconcile adds the roadmap-reconcile signal (sev 9). It fires on a
// ROADMAP.md row naming a specs/<slug>/spec.md whose work has already shipped. The spec
// file is merged-implemented, spec.AwaitsRetirement, or was retired out of the tree
// entirely. Such a row has outlived the drain that should have removed it.
//
// Like appendRetirement, it fires only on the default branch. A topic branch's roadmap
// is mid-build, so a row there names in-flight work, not a shipped-work leak. Severity
// 10 ranks it below the housekeeping rows, retirement 8 and orphaned-pickup 9, and far
// below gate/git. So it never displaces a red-gate or dirty-tree row in the budget.
func appendRoadmapReconcile(rows []row, root string) []row {
	// Audit #6: tolerate. As in appendRetirement, an unreadable branch or an unresolvable
	// default reads as "not the default branch". This skips this advisory reconcile signal.
	cur, _ := git.Output("-C", root, "rev-parse", "--abbrev-ref", "HEAD")
	if def, ok := git.ResolvedDefault(root); !ok || cur != def {
		return rows
	}
	merged, dangling, state := roadmapReconcileCounts(root)
	if state.Failed() {
		return append(rows, row{10, "roadmap", toon.UnknownCell(roadmap.RoadmapFile, state), commandAction(drainPhaseAction)})
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
	return append(rows, row{10, "roadmap", strings.Join(details, ", "), commandAction(drainPhaseAction)})
}

// appendTicketsOnly adds the tickets-only residue signal (sev 11): a specs/<slug>/ that a
// light-path landing wrote and never consumed. It names the command that closes one, so
// the row is actionable without a lookup. It fires only on a nonzero count: residue is
// debt, and a clean tree must not spend a row of the five-row budget saying so. It closes
// the housekeeping band below retirement (8) and orphaned pickups (9), so a count of
// residue never displaces a more urgent row.
//
// The tickets-only shape comes from landing.TicketsOnlyFolders, the same predicate the
// landing's `--spec` close consumes, so the row routes to `bench worktree land --spec
// <slug>`. It names no request, base, or tip, because the primary checkout does not know
// the in-flight worktree. An unreadable specs/ counts nothing, the posture the other
// advisory housekeeping rows take.
func appendTicketsOnly(rows []row, root string) []row {
	slugs, err := landing.TicketsOnlyFolders(root)
	if err != nil || len(slugs) == 0 {
		return rows
	}
	detail := Plural(len(slugs), "tickets-only spec folder", "tickets-only spec folders")
	return append(rows, row{11, "specs", detail, commandActionWithArgument(closeTicketsAction, "<slug>")})
}

// retirementCount counts live folder specs that spec.AwaitsRetirement marks, a merged
// spec awaiting retirement. Absent `specs/` returns 0. The unfenced-marker predicate is
// spec.AwaitsRetirement's one source. Every read goes through the classifier: this
// signal is advisory and stays quiet about a spec it could not read. A FIFO parked at
// specs/*/spec.md never yields EOF. A board that blocks forever is worse than one
// missing a housekeeping row, since `bench status` is what the SessionStart hook runs.
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

// roadmapReconcileCounts scans ROADMAP.md's index plus every roadmap/ row-file body for
// live spec-path tokens. A row's spec path lives in its row file now, not the index
// line. It classifies each distinct path against the tree. A missing file is a dangling
// row, the spec retired but its roadmap row survived. A present file that
// spec.AwaitsRetirement marks is a merged row, the work shipped but the drain missed it.
// A present, still-staged spec is the normal open-work state and counts nothing.
//
// Absent or empty ROADMAP.md returns 0, 0, bounds.StateParsed, the ordinary
// quiet-roadmap posture. A ROADMAP.md whose read failed reports that state instead, so
// appendRoadmapReconcile renders the failed read as unknown rather than a fabricated
// clean board. Each named live spec goes through the classifier too, so a FIFO cannot
// block the board. A spec path that yields no content at all is the dangling case,
// whether nothing is there or what is there could not be read. The merged predicate is
// spec.AwaitsRetirement, the same one source the retirement counter applies.
func roadmapReconcileCounts(root string) (merged, dangling int, state bounds.FileState) {
	tree := roadmap.LoadTree(root)
	switch {
	case tree.Index.State.Failed():
		return 0, 0, tree.Index.State
	case tree.Index.State == bounds.StateAbsent || tree.Index.State == bounds.StateEmpty:
		return 0, 0, bounds.StateParsed
	}
	var content bytes.Buffer
	content.Write(tree.Index.Data)
	for _, file := range tree.Files {
		if file.State != bounds.StateParsed && file.State != bounds.StateEmpty {
			continue
		}
		content.WriteByte('\n')
		content.Write(file.Data)
	}
	for _, slug := range spec.LiveSpecSlugs(content.Bytes()) {
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

// Short returns the first up-to-7 bytes of s, the shell `${var:0:7}` tree-prefix slice,
// guarding a short or "none" hash so the slice never panics. It is exported as the one
// source of the prefix width. Every surface that renders a tree or commit reference
// cuts it at the same place.
func Short(s string) string {
	return s[:min(7, len(s))]
}
