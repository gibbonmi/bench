package worktree

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/intent"
)

const explicitSetFingerprintVersion = "bench-explicit-set/v1"

// cleanSelection is one parsed `bench worktree clean` invocation: the single selection
// mode it names, the modifiers it was asked under, and the apply form it carries. The
// grammar lives here rather than in the verb, so the verb reads one answered question.
type cleanSelection struct {
	options      CleanupOptions
	target       string
	targets      []string
	fingerprint  string
	landed       bool
	unclaimed    bool
	applyCurrent bool
}

// parseCleanSelection reads the clean grammar and reports whether the invocation is
// valid. Every operand after a flag that takes one stays data: the parser consumes it by
// position and never re-reads it as a flag.
func parseCleanSelection(args []string) (cleanSelection, bool) {
	selection := cleanSelection{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--discard-ignored":
			selection.options.DiscardIgnored = true
		case "--discard-branch":
			selection.options.DiscardBranch = true
		case "--full":
			selection.options.Full = true
		case "--unclaimed":
			if selection.unclaimed {
				return selection, false
			}
			selection.unclaimed, selection.options.Unclaimed = true, true
		case "--landed":
			if selection.landed {
				return selection, false
			}
			selection.landed = true
		case "--target":
			if i+1 >= len(args) || args[i+1] == "" {
				return selection, false
			}
			i++
			selection.targets = append(selection.targets, args[i])
		case "--apply":
			if i+1 >= len(args) || selection.fingerprint != "" {
				return selection, false
			}
			i++
			selection.fingerprint = args[i]
		case "--apply-current":
			if selection.applyCurrent {
				return selection, false
			}
			selection.applyCurrent = true
		case "--":
			if i+1 >= len(args) || selection.target != "" {
				return selection, false
			}
			i++
			selection.target = args[i]
		default:
			if selection.target != "" || strings.HasPrefix(args[i], "-") {
				return selection, false
			}
			selection.target = args[i]
		}
	}
	return selection, selection.valid()
}

// valid holds the rules that read across operands: exactly one selection mode, the
// unclaimed mode's own modifier eligibility, and the two apply forms. The explicit set
// carries the ordinary discard modifiers and implies neither of them.
func (selection cleanSelection) valid() bool {
	modes := 0
	for _, named := range []bool{selection.target != "", selection.landed, selection.unclaimed, len(selection.targets) > 0} {
		if named {
			modes++
		}
	}
	if modes != 1 {
		return false
	}
	if selection.unclaimed && (!selection.options.DiscardBranch || selection.options.DiscardIgnored || selection.options.Full) {
		return false
	}
	if selection.applyCurrent && (!selection.unclaimed || selection.fingerprint != "") {
		return false
	}
	return selection.fingerprint == "" || wellFormedFingerprintOrPrefix(selection.fingerprint)
}

// explicitCleanupRow is one selected target: the cleanup identity its operand resolved
// to, and the plan the existing explicit planner produced for that identity's checkout.
type explicitCleanupRow struct {
	assignment intent.Assignment
	plan       CleanupPlan
	// targetFingerprint is the digest the explicit planner produced for this target alone.
	// The row's own plan carries the set digest for output, while the locked transaction
	// still validates the target against the state this one binds.
	targetFingerprint string
}

// explicitCleanupSet is the complete resolved selection. A selection failure leaves the
// fingerprint empty, so the plan reports every outcome and authorizes no removal.
type explicitCleanupSet struct {
	rows        []explicitCleanupRow
	failures    []CleanupPlan
	fingerprint string
}

// planExplicitSet resolves every operand through the assignment identity resolver that
// the other target-taking verbs share, collapses aliases by that identity, and plans each
// member through the same explicit planner the single-target form calls. One unresolved
// or ambiguous operand makes the whole selection unapplicable.
func planExplicitSet(j joins, root string, targets []string, options CleanupOptions) explicitCleanupSet {
	set := explicitCleanupSet{}
	selected := make(map[string]bool, len(targets))
	for _, target := range targets {
		assignment, err := resolveAssignment(root, target)
		if err != nil {
			set.failures = append(set.failures, selectionFailurePlan(target, err))
			continue
		}
		if selected[assignment.ID] {
			continue
		}
		selected[assignment.ID] = true
		row, planErr := planExplicitCleanupRow(j, root, assignment, options)
		if planErr != nil {
			set.failures = append(set.failures, selectionFailurePlan(target, planErr))
			continue
		}
		set.rows = append(set.rows, row)
	}
	sort.Slice(set.rows, func(a, b int) bool { return set.rows[a].assignment.ID < set.rows[b].assignment.ID })
	if len(set.failures) == 0 {
		set.fingerprint = fingerprintExplicitSet(set.rows, options)
	}
	for i := range set.rows {
		set.rows[i].plan.Fingerprint = unapplicableFingerprint
		if set.fingerprint != "" {
			set.rows[i].plan.Fingerprint = set.fingerprint
		}
	}
	return set
}

// planExplicitCleanupRow reads one member through the planner the single-target form uses.
// The set plan and every pre-transaction re-plan call it, so a member cannot become
// removable through a route the shared fingerprint never described.
func planExplicitCleanupRow(j joins, root string, assignment intent.Assignment, options CleanupOptions) (explicitCleanupRow, error) {
	plan, err := planExplicitWith(j, root, assignment.Worktree, options)
	if err != nil {
		return explicitCleanupRow{}, err
	}
	plan.Assignment = assignment.ID
	return explicitCleanupRow{assignment: assignment, plan: plan, targetFingerprint: plan.Fingerprint}, nil
}

// sameExplicitCleanupTuple compares the facts a removal reads about one target. It is
// deliberately narrower than the per-target fingerprint, which also binds repository-wide
// ledger evidence: an earlier member's own removal rewrites that evidence, and that must
// not read as drift in a later member of the same approved set.
func sameExplicitCleanupTuple(a, b explicitCleanupRow) bool {
	return a.assignment.ID == b.assignment.ID &&
		a.plan.Target == b.plan.Target &&
		a.plan.Action == b.plan.Action &&
		a.plan.Tracked == b.plan.Tracked &&
		a.plan.Recovery == b.plan.Recovery &&
		a.plan.Ignored.Digest == b.plan.Ignored.Digest &&
		a.plan.deleteBranch == b.plan.deleteBranch &&
		a.plan.branchOID == b.plan.branchOID
}

// unapplicableFingerprint is what a row carries when no apply can name the set. The
// invalid-invocation row reads it from here too, so one unapplicable plan carries the same
// spelling wherever the command produces it.
const unapplicableFingerprint = "none"

// selectionFailurePlan reports one operand the resolver could not turn into exactly one
// cleanup identity. The operand stays data: it reaches the reader through the same output
// sink every target passes, and it names no removal.
func selectionFailurePlan(target string, err error) CleanupPlan {
	return CleanupPlan{
		Target: target, Action: ActionError, ReasonCode: ReasonUncertain, Reason: err.Error(),
		Tracked: "unknown", ignoredSummary: "unknown", Recovery: "none", Fingerprint: unapplicableFingerprint, unresolved: true,
	}
}

// fingerprintExplicitSet binds complete membership and each member's own removal-relevant
// state. The per-target digest is that state: the explicit planner already derives it from
// every fact a removal reads, so the set never derives the same tuple a second time.
func fingerprintExplicitSet(rows []explicitCleanupRow, options CleanupOptions) string {
	parts := [][]byte{
		[]byte(explicitSetFingerprintVersion),
		[]byte(strconv.FormatBool(options.DiscardIgnored)),
		[]byte(strconv.FormatBool(options.DiscardBranch)),
		[]byte(strconv.FormatBool(options.Full)),
	}
	for _, row := range rows {
		parts = append(parts, []byte(row.assignment.ID), []byte(row.plan.Target), []byte(row.plan.Action), []byte(row.targetFingerprint))
	}
	return fingerprintParts(parts...)
}

// plans is every selection outcome in one order: the resolved members by canonical
// identity, then each operand the resolver refused.
func (set explicitCleanupSet) plans() []CleanupPlan {
	plans := make([]CleanupPlan, 0, len(set.rows)+len(set.failures))
	for _, row := range set.rows {
		plans = append(plans, row.plan)
	}
	return append(plans, set.failures...)
}

// applyArguments renders the set's own apply invocation: everything its re-plan command
// names, then the digest this plan authorizes.
func (set explicitCleanupSet) applyArguments(options CleanupOptions) []axi.InvocationArgument {
	arguments := cleanArguments(options, set.targetSelectors()...)
	return append(arguments, axi.KnownArgument("--apply"), axi.KnownArgument(set.fingerprint))
}

// targetSelectors names each member by its canonical assignment identity. The identity
// replaces the operand the caller typed, so no operand text reaches a command line.
func (set explicitCleanupSet) targetSelectors() []string {
	selectors := make([]string, 0, 2*len(set.rows))
	for _, row := range set.rows {
		selectors = append(selectors, "--target", row.assignment.ID)
	}
	return selectors
}

// renderExplicitSet prints one row per selection outcome and, when the set is applicable
// and plans a removal, the exact invocation that applies the whole selection.
func renderExplicitSet(stdout io.Writer, set explicitCleanupSet, options CleanupOptions) error {
	if err := renderCleanups(stdout, set.plans()); err != nil || set.fingerprint == "" {
		return err
	}
	removes := false
	for _, row := range set.rows {
		removes = removes || row.plan.Action.Removes()
	}
	if !removes {
		return nil
	}
	help, err := axi.RenderHelp([]axi.Action{axi.ExecutableInvocation("apply the explicit worktree set", set.applyArguments(options)...)})
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(stdout, help)
	return err
}

// requalifyExplicitRow re-plans one member against the repository as it stands now and
// refuses a member the approved set no longer describes. The preflight and the apply loop
// read the same proof, so neither can admit a member the other would refuse.
func requalifyExplicitRow(j joins, root string, planned explicitCleanupRow, options CleanupOptions) (explicitCleanupRow, error) {
	current, err := planExplicitCleanupRow(j, root, planned.assignment, options)
	if err != nil {
		return current, err
	}
	if !sameExplicitCleanupTuple(planned, current) {
		return current, errStaleFingerprint
	}
	return current, nil
}

// explicitRowPlans is the plan each selected member carries, in selection order.
func explicitRowPlans(rows []explicitCleanupRow) []CleanupPlan {
	plans := make([]CleanupPlan, 0, len(rows))
	for _, row := range rows {
		plans = append(plans, row.plan)
	}
	return plans
}

// applyExplicitSet runs every removable member through the existing per-target locked
// transaction, so each member keeps its own authority checks, receipt, and recovery. A
// member the plan retained passes through as the plan reported it. Each member re-plans
// immediately before its own transaction, which gives that transaction the digest the
// repository answers to now. The first member that cannot finish stops the set, and every
// member the set never started reports its own unstarted outcome.
func applyExplicitSet(j joins, root string, set explicitCleanupSet, options CleanupOptions) ([]CleanupPlan, error) {
	plans := make([]CleanupPlan, 0, len(set.rows))
	if offender, err := preflightExplicitSet(j, root, set, options); err != nil {
		return preflightOutcomes(plans, explicitRowPlans(set.rows), offender, err), err
	}
	for i, planned := range set.rows {
		if !planned.plan.Action.Removes() {
			plans = append(plans, planned.plan)
			continue
		}
		unreached := func() []CleanupPlan { return explicitRowPlans(set.rows[i+1:]) }
		current, err := requalifyExplicitRow(j, root, planned, options)
		if err != nil {
			plans = append(plans, requalifiedOutcome(planned.plan, current.plan, err))
			return notAttemptedPlans(plans, unreached(), notAttemptedDetail), err
		}
		applied, applyErr := applyExplicitWith(j, root, current.plan.Target, current.targetFingerprint, options)
		if applyErr != nil {
			plans = append(plans, faultedPlan(current.plan, applied, applyErr))
			return notAttemptedPlans(plans, unreached(), notAttemptedDetail), applyErr
		}
		plans = append(plans, applied)
	}
	return plans, nil
}

// cleanExplicitSet is the explicit-target selection mode. A bare call plans the complete
// set, and an apply call validates the whole selection before the first transaction runs.
func cleanExplicitSet(j joins, root string, selection cleanSelection, stdout, stderr io.Writer) int {
	set := planExplicitSet(j, root, selection.targets, selection.options)
	if selection.fingerprint == "" {
		if err := renderExplicitSet(stdout, set, selection.options); err != nil {
			fmt.Fprintf(stderr, "bench worktree clean: %v\n", err)
			return 1
		}
		if set.fingerprint == "" {
			return 1
		}
		return 0
	}
	// An unresolved selection has no scope to re-plan: the operands, not the digest, are
	// what the caller has to correct, so that refusal renders no recovery command.
	if set.fingerprint == "" {
		_ = renderCleanups(stdout, staleRows(selection.fingerprint, set.plans()))
		return 1
	}
	replan := cleanArguments(selection.options, set.targetSelectors()...)
	if !matchesFingerprint(set.fingerprint, selection.fingerprint) {
		_ = renderStaleSet(stdout, selection.fingerprint, set.plans(), replan)
		return 1
	}
	plans, applyErr := applyExplicitSet(j, root, set, selection.options)
	if err := applyOutcomes(stdout, plans, staleRows(selection.fingerprint, plans), applyErr, replan); err != nil {
		fmt.Fprintf(stderr, "bench worktree clean: %v\n", err)
		return 1
	}
	if applyErr != nil {
		return 1
	}
	return 0
}
