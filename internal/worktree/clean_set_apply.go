// Cleanup-set apply: what every selection mode shares once a caller asks to apply a plan.
// The outcome rows a stopped apply reports and the refusal a stale plan renders live here,
// so the explicit set, the landed selector, and the unclaimed selector cannot answer the
// same question in three different ways.
//
// The preflight that requalifies each member before the first transaction lives here for the
// two set modes. The unclaimed selector has no per-member preflight: it re-plans the whole
// branch selection inside applyUnclaimedAssignmentSet and refuses on any difference.

package worktree

import (
	"errors"
	"fmt"
	"io"

	"github.com/gibbonmi/bench/internal/axi"
)

// ActionNotAttempted is the outcome of a member a stopped set apply never started. Only an
// apply reports it, and it says what did not happen, so it is never a planned removal claim.
// It stays outside Removes() for that reason: the member still stands exactly as the
// approved plan found it, and it has no removal ahead of it.
//
// lifecyclepolicy owns the other action values. This one sits here because moving it grows
// two files that are already over their size budget.
const ActionNotAttempted CleanupAction = "not-attempted"

// The three reasons a selected member goes unstarted. A member the apply reached and stopped
// short of names the earlier failure; a member refused before the first transaction names the
// qualification, because at that point nothing had been attempted at all; and a member whose
// own requalification refused it names that, because nothing before it failed.
const (
	notAttemptedDetail = "not attempted; an earlier target in this set did not complete"
	notQualifiedDetail = "not attempted; the set did not qualify before the first removal"
	driftedDetail      = "not attempted; this target no longer matches the approved plan"
)

// notAttemptedPlan reports one selected target the apply never reached. The ignored preview
// belongs to the plan that offered the removal, so an unstarted row keeps the count in its
// summary and leaves the path table to that plan.
func notAttemptedPlan(plan CleanupPlan, detail string) CleanupPlan {
	plan.Ignored.Shown = 0
	plan.Action, plan.ReasonCode, plan.Reason = ActionNotAttempted, "", detail
	return plan
}

// StepMemberRequalify is the window between the set preflight and one member's own
// requalification. Earlier members' transactions have already run by then, so a concurrent
// writer has had real time to change a later member since the preflight cleared it. That is
// the race each member's own recheck exists for, and a test stands in this window to open it.
// The token sits here rather than in ownership.go, which is over its size budget and outside
// this spec's ownership fence.
const StepMemberRequalify LifecycleStep = "member-requalify"

// unstartedPlan is the outcome of one selected member no transaction touched. A member the
// plan marked removable reports that the removal never started. A member the plan did not
// mark removable keeps the verdict it has: the apply was never going to touch it, so its
// retain or refusal authority is what happened to it, and calling it unstarted would erase
// the reason it was spared. Every unstarted row in this package reads that rule here.
func unstartedPlan(plan CleanupPlan, detail string) CleanupPlan {
	if plan.Action.Removes() {
		return notAttemptedPlan(plan, detail)
	}
	return plan
}

// notAttemptedPlans appends one outcome per member a stopped apply never reached. Both
// selection modes report this the same way; only the row type each one slices differs.
func notAttemptedPlans(plans, unreached []CleanupPlan, detail string) []CleanupPlan {
	for _, plan := range unreached {
		plans = append(plans, unstartedPlan(plan, detail))
	}
	return plans
}

// preflightOutcomes reports a selection the preflight refused before any transaction opened.
// Drift refuses the whole approved set, so every member reads as unqualified. Any other
// fault belongs to the member whose requalification raised it, and that member's row carries
// the reason; without this the fault reaches neither a row nor stderr.
func preflightOutcomes(plans, rows []CleanupPlan, offender int, err error) []CleanupPlan {
	if errors.Is(err, errStaleFingerprint) {
		return notAttemptedPlans(plans, rows, notQualifiedDetail)
	}
	for i, plan := range rows {
		if i == offender {
			plans = append(plans, faultedPlan(plan, CleanupPlan{}, err))
			continue
		}
		plans = append(plans, unstartedPlan(plan, notQualifiedDetail))
	}
	return plans
}

// faultedPlan names the target whose transaction did not finish. A fault can return before
// the transaction has a plan of its own, so the approved row supplies the identity the
// outcome still has to report.
func faultedPlan(planned, applied CleanupPlan, err error) CleanupPlan {
	if applied.Target == "" {
		applied = planned
	}
	applied.Action, applied.ReasonCode, applied.Reason = ActionError, "", err.Error()
	return applied
}

// requalifiedOutcome is the row for a member its own requalify refused, after the preflight
// had already passed. Drift is not a fault: the member re-planned, and the fresh verdict is
// what it is now. No transaction opened on it either way, so the fresh row goes through the
// same unstarted rule as any other member the apply did not touch — a member that drifted to
// retained reports retained, and one still marked removable reports that it never started.
// Any other error is a fault, and the row carries its reason instead.
func requalifiedOutcome(planned, current CleanupPlan, err error) CleanupPlan {
	if errors.Is(err, errStaleFingerprint) {
		return unstartedPlan(current, driftedDetail)
	}
	return faultedPlan(planned, current, err)
}

// staleSetPlan is the refusal row a set apply prints when the plan it carries no longer
// describes the repository. The landed set and the explicit set share it, so their two
// stale refusals cannot drift apart. The unclaimed set builds its own row instead, because
// it reports branch refs under that mode's own tracked and ignored spellings.
func staleSetPlan(fingerprint string) CleanupPlan {
	return CleanupPlan{
		Target: "unknown", Action: ActionError, Tracked: "unknown", ignoredSummary: "unknown",
		Recovery: "none", Fingerprint: fingerprint, Reason: errStaleFingerprint.Error(),
	}
}

// cleanupModifierFlags names the discard modifiers one invocation carried. Every rendered
// command reads them from here, so no surface can advertise a command that asks a different
// question than the plan answered.
func cleanupModifierFlags(options CleanupOptions) []string {
	var flags []string
	if options.DiscardIgnored {
		flags = append(flags, "--discard-ignored")
	}
	if options.DiscardBranch {
		flags = append(flags, "--discard-branch")
	}
	if options.Full {
		flags = append(flags, "--full")
	}
	return flags
}

// cleanArguments is the head of every clean command this package renders that carries the
// caller's modifiers: the verb, those modifiers, then the mode's own selectors. It is also the
// exact re-plan command a refusal offers, so recovery keeps the scope the plan was asked for.
// A command that resolves one named path answers no modifiers and builds its own operands.
func cleanArguments(options CleanupOptions, selectors ...string) []axi.InvocationArgument {
	arguments := []axi.InvocationArgument{axi.KnownArgument("worktree"), axi.KnownArgument("clean")}
	for _, selector := range append(cleanupModifierFlags(options), selectors...) {
		arguments = append(arguments, axi.KnownArgument(selector))
	}
	return arguments
}

// renderStale prints a set refusal and the exact command that re-plans the same selection.
// The re-plan carries no digest, because a fresh plan is what the refused caller lacks.
func renderStale(stdout io.Writer, rows []CleanupPlan, replan []axi.InvocationArgument) error {
	if err := renderCleanups(stdout, rows); err != nil {
		return err
	}
	help, err := axi.RenderHelp([]axi.Action{
		axi.ExecutableInvocation("re-plan the refused worktree selection", replan...),
	})
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(stdout, help)
	return err
}

// staleRows puts the row that names the rejected digest ahead of the rows the caller has.
// The refusal before an apply and the refusal from inside one read the composition here, so
// the two cannot describe the same repository differently.
func staleRows(fingerprint string, rows []CleanupPlan) []CleanupPlan {
	return append([]CleanupPlan{staleSetPlan(fingerprint)}, rows...)
}

// renderStaleSet is a stale refusal in full: the composed rows, then the exact command that
// re-plans the same selection. It serves the two refusals a mode raises before it applies;
// an apply-time refusal composes the same rows through staleRows and renders them through
// applyOutcomes. The composition itself is single-sourced in staleRows either way.
func renderStaleSet(stdout io.Writer, fingerprint string, rows []CleanupPlan, replan []axi.InvocationArgument) error {
	return renderStale(stdout, staleRows(fingerprint, rows), replan)
}

// applyOutcomes prints what one apply produced. A stale refusal replaces the outcome rows
// with the refusal form, so the reader sees which plan the repository no longer describes
// and the command that re-plans the same selection. Every selection mode reads that switch
// here, and each supplies the refusal rows its own surface spells.
func applyOutcomes(stdout io.Writer, plans, stale []CleanupPlan, err error, replan []axi.InvocationArgument) error {
	if !errors.Is(err, errStaleFingerprint) {
		return renderCleanups(stdout, plans)
	}
	return renderStale(stdout, stale, replan)
}

// preflightExplicitSet and preflightLandedSet requalify every removable member before their
// mode opens its first transaction. Each transaction keeps its own recheck, but that recheck
// only runs when the set reaches that member, with the earlier members already gone. These
// passes are what stop an avoidable removal when a later member has drifted before the apply
// began. Each one reuses its mode's own single row proof, so neither derives drift twice.
//
// Each returns the index of the member that refused, so a fault that is not drift can name
// itself in that member's row instead of disappearing behind an unstarted detail.
func preflightExplicitSet(j joins, root string, set explicitCleanupSet, options CleanupOptions) (int, error) {
	for i, planned := range set.rows {
		if !planned.plan.Action.Removes() {
			continue
		}
		if _, err := requalifyExplicitRow(j, root, planned, options); err != nil {
			return i, err
		}
	}
	return -1, nil
}

func preflightLandedSet(j joins, root string, set landedCleanupSet, options CleanupOptions, scope string) (int, error) {
	for i, planned := range set.rows {
		if !planned.plan.Action.Removes() {
			continue
		}
		if _, err := requalifyLandedRow(j, root, planned, options, scope); err != nil {
			return i, err
		}
	}
	return -1, nil
}
