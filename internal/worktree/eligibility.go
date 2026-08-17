package worktree

import (
	"strings"

	"github.com/gibbonmi/bench/internal/intent"
)

// This file owns the ordered eligibility decisions: the single answer to "is this worktree
// ours and safe to remove, and if not, why", for every consumer that decides one.
// PlanExplicitWithOptions in subshell.go gathers every Git and filesystem fact, builds one
// explicitFacts value from what it gathered, and calls decideExplicit exactly once.
// PlanAutomatic in classifier.go layers its own stricter reading on top: it calls
// PlanExplicit first, gathers the automatic-specific facts decideAutomatic needs, and
// calls decideAutomatic exactly once. Nothing outside this
// file orders or selects an eligibility action or reason before execution; subshell.go and
// classifier.go only project the returned verdict onto the operator-facing CleanupPlan,
// including any lookup (a recovery ref prediction) that needs I/O the decision itself must
// not perform.

// landednessKind names which of the four ways PlanExplicitWithOptions can know a branch's
// relationship to the default ref: never resolvable (detached), resolvable but never asked
// (no default), asked and refused by the query itself, or asked and answered.
type landednessKind int

const (
	landednessDetached landednessKind = iota
	landednessUnknownNoDefault
	landednessUnknownError
	landednessProven
)

// landedness carries a branch's proven relationship to the default ref as typed fields;
// every decision reads those. String renders the "true:ancestry" / "unknown:<err>" wire
// form, whose only consumer is the explicit fingerprint, which hashes it as evidence.
type landedness struct {
	kind      landednessKind
	err       string
	landed    bool
	byContent bool
}

// provenLanded reports whether the query was actually asked and answered yes — the one
// reading that authorizes acting on a landing, as opposed to detachment, an absent default
// ref, a failed query, or a proven non-landing.
func (l landedness) provenLanded() bool {
	return l.kind == landednessProven && l.landed
}

func (l landedness) String() string {
	switch l.kind {
	case landednessDetached:
		return "detached"
	case landednessUnknownNoDefault:
		return "unknown"
	case landednessUnknownError:
		return "unknown:" + l.err
	case landednessProven:
		proof := "ancestry"
		if l.byContent {
			proof = "patch"
		}
		return boolLabel(l.landed) + ":" + proof
	default:
		return "unknown"
	}
}

func boolLabel(ok bool) string {
	if ok {
		return "true"
	}
	return "false"
}

// recoveryLookupKind names whether a recovery ref still needs a git query, and which query,
// once the verdict has decided a ref is needed at all. Predicting a ref is not itself a
// decision — nextRecoveryRef and predictedForeignRef stay in subshell.go, called exactly
// once, after the eligibility call, using whichever kind the verdict returns.
type recoveryLookupKind int

const (
	recoveryNoLookup recoveryLookupKind = iota
	recoveryLookupOwned
	recoveryLookupForeign
)

// explicitFacts carries every typed fact PlanExplicitWithOptions gathers before it can
// answer "ours and safe to remove". Each field is evidence, not a conclusion: the same
// facts feed decideExplicit whether the eventual verdict retains, removes, recovers, or
// discards, and a field left at its zero value was never gathered because an earlier fact
// already made it inapplicable — the assignment ledger, for one, is read only when a
// marker validated.
type explicitFacts struct {
	registrationBranchRef  string
	registrationLockReason string
	registrationLocked     bool
	registrationDetached   bool

	markerPresent       bool
	markerErr           error
	assignmentLedgerErr error
	assignmentAmbiguous bool
	matchedAssignment   *intent.Assignment
	foreignAssignment   *intent.Assignment

	leasePresent bool
	leaseState   LeaseState
	leaseStatErr error

	nestedState nestedState
	nestedErr   error

	buildOutputErr   error
	ignoredErr       error
	ignoredOverLimit bool
	ignoredCount     int
	declaredIgnored  bool
	discardIgnored   bool

	initialTracked string

	headDetached    bool
	defaultKnown    bool
	landedOK        bool
	landedByContent bool
	landedErr       error
	headRef, head   string

	unsafeTarget bool
}

// explicitVerdict is the decided answer plus every piece of evidence PlanExplicitWithOptions
// needs to project it onto CleanupPlan and, for a RecoverRemove or DiscardRemove verdict,
// to finish resolving the recovery ref. Every field is the one-time output of a single
// decideExplicit call; no consumer decides over it again.
type explicitVerdict struct {
	Action     CleanupAction
	ReasonCode CleanupReason
	Reason     string

	Owned      bool
	Assignment *intent.Assignment

	Tracked string

	Landed       landedness
	DeleteBranch bool
	BranchRef    string
	BranchOID    string

	RecoveryLookup recoveryLookupKind
	Recovery       string
}

// decideExplicit is the single place the explicit ordered eligibility decision is made:
// ownership and marker agreement, assignment join, lock agreement, lease liveness, nested
// or embedded repository state, ignored-residue authorization, the recover/discard-remove
// promotion, and the final unsafe-target override. The order and the last-write-wins
// collisions between blocks are pinned by TestExplicitEligibilityOutcomeMatrix and must not
// change here without that characterization moving first.
func decideExplicit(f explicitFacts) explicitVerdict {
	v := explicitVerdict{Action: ActionRemove, Tracked: f.initialTracked, Recovery: "none"}

	switch {
	case f.markerPresent:
		switch {
		case f.markerErr != nil:
			v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonMalformed, f.markerErr.Error()
		case f.assignmentLedgerErr != nil:
			v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonMalformed, "assignment ledger is unreadable"
		case f.assignmentAmbiguous:
			v.Assignment = f.matchedAssignment
			v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonMalformed, "registration has ambiguous assignments"
		case f.matchedAssignment == nil:
			v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonMalformed, "owner marker has no matching assignment"
		default:
			v.Owned = true
			v.Assignment = f.matchedAssignment
			if f.registrationBranchRef != f.matchedAssignment.Branch {
				v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonUncertain, "assignment does not match current branch"
			} else if f.registrationLockReason != lockReason(*f.matchedAssignment) {
				v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonUnexpectedLock, "assignment does not match current Bench lock"
			}
		}
	case f.registrationLocked:
		v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonUnexpectedLock, "foreign or unexpected lock is retained"
	default:
		v.Assignment = f.foreignAssignment
	}

	if f.leasePresent {
		switch f.leaseState {
		case LeaseLive:
			if v.Owned {
				v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonLiveLease, "assignment has a live lease"
			} else {
				v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonLiveLease, "unowned assignment has an ambiguous lease"
			}
		case LeaseDead:
			if !v.Owned {
				v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonLiveLease, "unowned assignment has an ambiguous lease"
			}
		case LeaseUnknown:
			v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonUncertain, unknownLeaseReason
		}
	} else if f.leaseStatErr != nil && v.Owned {
		v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonUncertain, unknownLeaseReason
	}

	if f.nestedErr != nil {
		v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonUncertain, "nested repository state is unknown"
	} else if f.nestedState == nestedDirty {
		v.Tracked = "nested-dirty"
		v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonUncertain, "nested repository or submodule is dirty"
	} else if f.nestedState == nestedEmbeddedClean || f.nestedState == nestedEmbeddedDirty {
		v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonUncertain, "embedded repository is retained"
	}

	if f.buildOutputErr != nil {
		v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonMalformed, "build-output declaration is malformed"
	} else if f.ignoredErr != nil {
		v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonUncertain, "ignored inventory is uncertain"
	} else if f.ignoredOverLimit {
		v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonIgnored, "ignored inventory exceeds the destructive limit"
	} else if f.ignoredCount > 0 && !f.discardIgnored && !f.declaredIgnored {
		v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonIgnored, "ignored residuals require --discard-ignored"
	}

	switch {
	case f.headDetached:
		v.Landed = landedness{kind: landednessDetached}
	case !f.defaultKnown:
		v.Landed = landedness{kind: landednessUnknownNoDefault}
	case f.landedErr != nil:
		v.Landed = landedness{kind: landednessUnknownError, err: f.landedErr.Error()}
	default:
		v.Landed = landedness{kind: landednessProven, landed: f.landedOK, byContent: f.landedByContent}
		if f.landedOK {
			v.DeleteBranch = true
			v.BranchRef, v.BranchOID = f.headRef, f.head
		}
	}

	if v.Action != ActionRetain && (v.Tracked != "clean" || f.registrationDetached || (v.Assignment != nil && len(v.Assignment.Recovery) > 0)) {
		v.Action = ActionRecoverRemove
		switch {
		case v.Owned && v.Assignment != nil:
			if len(v.Assignment.Recovery) > 0 {
				v.Recovery = v.Assignment.Recovery[0].Ref
			} else {
				v.RecoveryLookup = recoveryLookupOwned
			}
		case v.Assignment != nil:
			v.Recovery = v.Assignment.Recovery[0].Ref
		default:
			v.RecoveryLookup = recoveryLookupForeign
		}
	}

	if v.Action != ActionRetain && f.ignoredCount > 0 && (f.discardIgnored || f.declaredIgnored) {
		v.Action = ActionDiscardRemove
	}

	if f.unsafeTarget {
		v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonUncertain, "target contains unsafe control bytes"
	}

	return v
}

// automaticFacts carries PlanExplicit's own result plus every automatic-specific fact
// PlanAutomatic gathers on top of it before decideAutomatic can answer "is this eligible for
// unattended, automatic cleanup, and if not, why not". explicitErr and explicit are mutually
// exclusive with the missingBranch* fields: the missingBranch* fields are gathered only when
// PlanExplicit itself failed and PlanAutomatic salvages one narrow fact from that failure;
// every other field is gathered only when PlanExplicit succeeded, mirroring PlanAutomatic's
// own conditional evidence-gathering — a field left at its zero value simply was never
// gathered because an earlier fact already made it inapplicable.
type automaticFacts struct {
	explicitErr error
	explicit    CleanupPlan

	missingBranchAssignment *intent.Assignment
	missingBranchLiveLease  bool

	liveLease       bool
	landed          bool
	orphanedActive  bool
	recoveryMatches bool
}

// automaticVerdict is the decided answer for the automatic, unattended cleanup path. Action,
// ReasonCode, and Reason are the decision itself; AssignmentID echoes the one piece of
// evidence PlanAutomatic projects onto CleanupPlan.Assignment that decideExplicit never sets
// (subshell.go never assigns it). It stays empty on every branch that reaches its verdict
// without a verified assignment to name, so the projection leaves CleanupPlan.Assignment
// untouched there.
type automaticVerdict struct {
	Action       CleanupAction
	ReasonCode   CleanupReason
	Reason       string
	AssignmentID string
}

// decideAutomatic is the single place the automatic ordered eligibility decision is made:
// the explicit-planning-error salvage (an active assignment whose branch cannot be
// resolved), the live-lease override, retain-passthrough with a landed-reason swap, the
// foreign/unowned refusal, the not-cleanup-pending reasons (with the orphaned-age override),
// the recovery-metadata-match check, the unknown/unmerged landedness checks, and the final
// preservation refusal. The order and every message are pinned by
// TestAutomaticEligibilityOutcomeMatrix and must not change here without that
// characterization moving first.
func decideAutomatic(f automaticFacts) automaticVerdict {
	if f.explicitErr != nil {
		if f.missingBranchAssignment != nil {
			v := automaticVerdict{Action: ActionRetain, ReasonCode: ReasonActive, Reason: "assignment landedness is unknown", AssignmentID: f.missingBranchAssignment.ID}
			if f.missingBranchLiveLease {
				v.ReasonCode, v.Reason = ReasonLiveLease, "assignment has a live lease"
			}
			return v
		}
		reason := ReasonUncertain
		if strings.Contains(f.explicitErr.Error(), "assignment") || strings.Contains(f.explicitErr.Error(), "intent ledger") {
			reason = ReasonMalformed
		}
		return automaticVerdict{Action: ActionRetain, ReasonCode: reason, Reason: f.explicitErr.Error()}
	}

	plan := f.explicit
	if plan.assignment != nil && f.liveLease {
		return automaticVerdict{Action: ActionRetain, ReasonCode: ReasonLiveLease, Reason: "assignment has a live lease"}
	}
	if plan.Action == ActionRetain {
		if plan.assignment != nil && f.landed {
			return automaticVerdict{Action: ActionRetain, ReasonCode: ReasonLanded, Reason: "assignment branch has landed"}
		}
		return automaticVerdict{Action: plan.Action, ReasonCode: plan.ReasonCode, Reason: plan.Reason}
	}
	if plan.assignment == nil || !plan.owned {
		return automaticVerdict{Action: ActionRetain, ReasonCode: ReasonForeign, Reason: "registration is not a verified owned assignment"}
	}

	assignmentID := plan.assignment.ID
	if plan.assignment.State != intent.StateCleanupPending {
		reason := ReasonUncertain
		if plan.assignment.State == intent.StateActive {
			if f.landed {
				reason = ReasonLanded
			} else {
				reason = ReasonActive
			}
			if reason == ReasonActive && f.orphanedActive {
				reason = ReasonOrphaned
			}
		}
		return automaticVerdict{Action: ActionRetain, ReasonCode: reason, Reason: "assignment is not cleanup-pending", AssignmentID: assignmentID}
	}
	if !f.recoveryMatches {
		return automaticVerdict{Action: ActionRetain, ReasonCode: ReasonMalformed, Reason: "assignment recovery metadata does not match refs", AssignmentID: assignmentID}
	}
	if plan.landedTyped.kind == landednessUnknownNoDefault || plan.landedTyped.kind == landednessUnknownError {
		return automaticVerdict{Action: ActionRetain, ReasonCode: ReasonUncertain, Reason: "assignment landedness is unknown", AssignmentID: assignmentID}
	}
	if !plan.landedTyped.provenLanded() {
		return automaticVerdict{Action: ActionRetain, ReasonCode: ReasonUnmerged, Reason: "assignment branch has not landed", AssignmentID: assignmentID}
	}
	// The automatic path authors no preservation refs: it runs unattended at every session
	// start and through every release, and the standing cleaner sweeps the namespace such a
	// ref would live in, so preserving there would write work nothing can hand back.
	// Disposing of the checkout stays with the operator's explicit path-addressed clean.
	if retain, action, reasonCode, reason := automaticPreservationVerdict(plan, "automatic cleanup does not preserve uncommitted work"); retain {
		return automaticVerdict{Action: action, ReasonCode: reasonCode, Reason: reason, AssignmentID: assignmentID}
	}
	return automaticVerdict{Action: plan.Action, ReasonCode: plan.ReasonCode, Reason: plan.Reason, AssignmentID: assignmentID}
}

// automaticPreservationVerdict is the one place automatic-flavored policy decides whether
// a plan must be retained to avoid stranding uncommitted work, and what Action/ReasonCode
// that refusal carries. decideAutomatic's own dirty-refusal branch and the landed-set's
// retainForLandedPreservation (clean_landed.go) both call this — neither writes its own
// ActionRetain/ReasonDirty literal — so the two routes can never derive "would removing
// this strand uncommitted work" differently. Each caller still supplies its own
// operator-facing message for its own command surface.
func automaticPreservationVerdict(plan CleanupPlan, message string) (retain bool, action CleanupAction, reasonCode CleanupReason, reason string) {
	if !plan.preserves() {
		return false, plan.Action, plan.ReasonCode, plan.Reason
	}
	return true, ActionRetain, ReasonDirty, message
}
