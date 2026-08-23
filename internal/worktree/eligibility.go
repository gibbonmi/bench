package worktree

import (
	"github.com/gibbonmi/bench/internal/worktree/lifecyclepolicy"
)

// The ordered eligibility decisions live in the pure child package
// internal/worktree/lifecyclepolicy: the single answer to "is this worktree
// ours and safe to remove, and if not, why". This file is the parent's seam to
// it. PlanExplicitWithOptions in subshell.go gathers every Git and filesystem
// fact, builds one explicitFacts value from what it gathered, and calls
// decideExplicit exactly once. PlanAutomatic in classifier.go layers its own
// stricter reading on top: it calls PlanExplicit first, gathers the
// automatic-specific facts decideAutomatic needs, and calls decideAutomatic
// exactly once.
//
// Nothing outside the policy package orders or selects an eligibility action or
// reason before execution. subshell.go and classifier.go only project the
// returned verdict onto the operator-facing CleanupPlan. This includes any
// lookup (a recovery ref prediction) that needs I/O the decision itself must
// not perform.

type landedness = lifecyclepolicy.Landedness

const (
	landednessDetached         = lifecyclepolicy.LandednessDetached
	landednessUnknownNoDefault = lifecyclepolicy.LandednessUnknownNoDefault
	landednessUnknownError     = lifecyclepolicy.LandednessUnknownError
	landednessProven           = lifecyclepolicy.LandednessProven
)

type explicitFacts = lifecyclepolicy.ExplicitFacts
type explicitVerdict = lifecyclepolicy.ExplicitVerdict
type automaticFacts = lifecyclepolicy.AutomaticFacts
type automaticVerdict = lifecyclepolicy.AutomaticVerdict

const (
	recoveryNoLookup      = lifecyclepolicy.RecoveryNoLookup
	recoveryLookupOwned   = lifecyclepolicy.RecoveryLookupOwned
	recoveryLookupForeign = lifecyclepolicy.RecoveryLookupForeign
)

func decideExplicit(f explicitFacts) explicitVerdict { return lifecyclepolicy.DecideExplicit(f) }

func decideAutomatic(f automaticFacts) automaticVerdict { return lifecyclepolicy.DecideAutomatic(f) }

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
