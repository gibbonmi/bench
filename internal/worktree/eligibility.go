package worktree

import (
	"github.com/gibbonmi/bench/internal/worktree/lifecyclepolicy"
)

// The ordered eligibility decisions live in the pure child package
// internal/worktree/lifecyclepolicy: the single answer to "is this worktree
// ours and safe to remove, and if not, why". PlanExplicitWithOptions in
// subshell.go gathers every Git and filesystem fact, builds one
// lifecyclepolicy.ExplicitFacts value from what it gathered, and calls
// lifecyclepolicy.DecideExplicit exactly once. PlanAutomatic in classifier.go
// layers its own stricter reading on top: it calls PlanExplicit first, gathers
// the automatic-specific facts lifecyclepolicy.DecideAutomatic needs, and calls
// it exactly once. This file holds the projections from a CleanupPlan into a
// policy input.
//
// Nothing outside the policy package orders or selects an eligibility action or
// reason before execution. subshell.go and classifier.go only project the
// returned verdict onto the operator-facing CleanupPlan. This includes any
// lookup (a recovery ref prediction) that needs I/O the decision itself must
// not perform.

// explicitOutcome projects an already-decided plan into the typed slice the
// automatic policy decision reads. It is a translation, not a decision: every
// field is copied evidence.
func explicitOutcome(plan CleanupPlan) lifecyclepolicy.ExplicitOutcome {
	outcome := lifecyclepolicy.ExplicitOutcome{
		Action:               plan.Action,
		ReasonCode:           plan.ReasonCode,
		Reason:               plan.Reason,
		HasAssignment:        plan.assignment != nil,
		Owned:                plan.owned,
		Tracked:              plan.Tracked,
		RegistrationDetached: plan.registration.Detached,
		Landed:               plan.landedTyped,
	}
	if plan.assignment != nil {
		outcome.AssignmentID = plan.assignment.ID
		outcome.AssignmentState = plan.assignment.State
	}
	return outcome
}

// automaticPreservationVerdict is the parent form of the policy's
// AutomaticPreservation over a full plan.
func automaticPreservationVerdict(plan CleanupPlan, message string) (retain bool, action CleanupAction, reasonCode CleanupReason, reason string) {
	return lifecyclepolicy.AutomaticPreservation(explicitOutcome(plan), message)
}
