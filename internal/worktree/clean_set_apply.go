// Cleanup-set apply: what every selection mode shares once a caller asks to apply a plan.
// The preflight each mode runs before its first transaction, the outcome rows a stopped
// apply reports, and the refusal a stale plan renders all live here, so the explicit set,
// the landed selector, and the unclaimed selector cannot answer the same question in three
// different ways.

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
// Its placement is unsettled. lifecyclepolicy owns the other action values, and this one
// sits here because moving it grows two files that are already over their size budget. The
// reviewer holds that decision.
const ActionNotAttempted CleanupAction = "not-attempted"

const notAttemptedDetail = "not attempted; an earlier target in this set did not complete"

// notAttemptedPlan reports one selected target the apply never reached. The ignored preview
// belongs to the plan that offered the removal, so an unstarted row keeps the count in its
// summary and leaves the path table to that plan.
func notAttemptedPlan(plan CleanupPlan) CleanupPlan {
	plan.Ignored.Shown = 0
	plan.Action, plan.ReasonCode, plan.Reason = ActionNotAttempted, "", notAttemptedDetail
	return plan
}

// notAttemptedPlans appends one outcome per member a stopped apply never reached. Both
// selection modes report this the same way; only the row type each one slices differs.
//
// A member the plan did not mark removable keeps the verdict the plan gave it. The apply was
// never going to touch that member, so its retain or refusal authority is what happened to
// it, and calling it unstarted would erase the reason it was spared.
func notAttemptedPlans(plans, unreached []CleanupPlan) []CleanupPlan {
	for _, plan := range unreached {
		if plan.Action.Removes() {
			plan = notAttemptedPlan(plan)
		}
		plans = append(plans, plan)
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

// cleanArguments is the head of every clean command this package renders: the verb, the
// modifiers the plan answered under, then the mode's own selectors. It is also the exact
// re-plan command a refusal offers, so recovery keeps the scope the plan was asked for.
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
// re-plans the same selection. Every mode's stale result renders through this one name.
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
func preflightExplicitSet(j joins, root string, set explicitCleanupSet, options CleanupOptions) error {
	for _, planned := range set.rows {
		if !planned.plan.Action.Removes() {
			continue
		}
		if _, err := requalifyExplicitRow(j, root, planned, options); err != nil {
			return err
		}
	}
	return nil
}

func preflightLandedSet(j joins, root string, set landedCleanupSet, options CleanupOptions, scope string) error {
	for _, planned := range set.rows {
		if !planned.plan.Action.Removes() {
			continue
		}
		if _, err := requalifyLandedRow(j, root, planned, options, scope); err != nil {
			return err
		}
	}
	return nil
}
