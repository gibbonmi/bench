package worktree

import "github.com/gibbonmi/bench/internal/intent"

// This file owns the ordered explicit-eligibility decision: the single answer to "is this
// worktree ours and safe to remove, and if not, why". PlanExplicitWithOptions in
// subshell.go gathers every Git and filesystem fact unconditionally, exactly as it always
// has, builds one explicitFacts value from what it gathered, and calls decideExplicit
// exactly once. Nothing outside this file orders or selects an eligibility action or
// reason before execution; subshell.go only projects the returned explicitVerdict onto the
// operator-facing CleanupPlan, including any lookup (a recovery ref prediction) that needs
// I/O the decision itself must not perform.
//
// PlanAutomatic (classifier.go) still decides its own stricter reading directly; routing it
// through this module is a later, separately reviewed migration.

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

// landedness is the typed replacement for the "true:ancestry" / "unknown:<err>" strings
// PlanAutomatic still parses by prefix; String reproduces that exact wire format so the
// out-of-scope consumer keeps working unchanged while every in-scope decision reads the
// typed fields instead.
type landedness struct {
	kind      landednessKind
	err       string
	landed    bool
	byContent bool
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
// discards, and a field left at its zero value simply was never gathered because an
// earlier fact already made it inapplicable (mirroring the original code's own
// conditional evidence-gathering, e.g. the assignment ledger is read only when a marker
// validated).
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
// to finish resolving the recovery ref. It is not a renamed plan: nothing here has been
// mutated by a second decision, it is the one-time output of decideExplicit.
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
