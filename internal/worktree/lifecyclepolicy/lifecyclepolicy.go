// Package lifecyclepolicy owns the pure lifecycle decisions of the worktree
// commands: ownership and marker agreement, lease liveness, ordered cleanup
// eligibility, assignment age, ignored-output authorization, preservation, and
// the resulting action. The parent package translates Git, filesystem, and
// process state into these typed facts once at its effect boundary. The
// decisions here read only the supplied facts and return a verdict the parent
// projects; the package performs no effects, reads no ambient process state,
// and starts no descendants. The source census test enforces that boundary.
package lifecyclepolicy

import (
	"strconv"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/intent"
)

// Action names what a cleanup plan will do, or why it will do nothing.
type Action string

// Reason names why a plan retains rather than removes.
type Reason string

const (
	ActionRetain         Action = "retain"
	ActionRemove         Action = "remove"
	ActionRecoverRemove  Action = "recover-remove"
	ActionDiscardRemove  Action = "discard-remove"
	ActionRemoved        Action = "removed"
	ActionError          Action = "error"
	ReasonForeign        Reason = "foreign"
	ReasonActive         Reason = "active"
	ReasonLiveLease      Reason = "live-lease"
	ReasonUnmerged       Reason = "unmerged"
	ReasonIgnored        Reason = "ignored"
	ReasonMalformed      Reason = "malformed"
	ReasonUncertain      Reason = "uncertain"
	ReasonUnexpectedLock Reason = "unexpected-lock"
	ReasonOrphaned       Reason = "orphaned"
	ReasonDirty          Reason = "dirty"
	ReasonLanded         Reason = "landed"
)

// ActionReleaseRemove is the release path's removal action.
const ActionReleaseRemove Action = "release-remove"

// ActionReleaseLeftover releases one assignment's registration and ledger entry
// while the bytes at its path stay exactly where they are. It is deliberately
// outside Removes(). Nothing proves what those bytes are: no checkout answers
// for them, and no recovery ref holds them. Disposing of them stays with the
// path-addressed clean surface, whose inventory is size-bounded.
const ActionReleaseLeftover Action = "release-leftover"

// Removes reports whether an action still has a removal ahead of it. It is not
// a refusal, an invocation error, or a completed transaction.
func (action Action) Removes() bool {
	return action == ActionRemove || action == ActionRecoverRemove || action == ActionDiscardRemove || action == ActionReleaseRemove
}

// Preserves reports whether executing a plan with this action, tracked state,
// and registration shape would write work to a recovery ref before removing the
// checkout. The execution and the planners that must not reach it read the same
// predicate. A plan can never be preserving to one and not to the other. A
// detached registration counts whatever its tree holds: the checkout's HEAD is
// the only thing naming its commits, so the removal would strand them.
func Preserves(action Action, tracked string, registrationDetached bool) bool {
	return action == ActionRecoverRemove ||
		(action == ActionDiscardRemove && tracked != "clean") ||
		registrationDetached
}

// LeaseState is the liveness verdict over one lease file's translated content.
type LeaseState string

const (
	LeaseLive    LeaseState = "live"
	LeaseDead    LeaseState = "dead"
	LeaseUnknown LeaseState = "unknown"
)

// UnknownLeaseReason is the refusal detail an unknown lease state renders.
const UnknownLeaseReason = "assignment lease state is unknown"

// LeaseTimeLayout is the UTC instant format a lease line records.
const LeaseTimeLayout = "2006-01-02T15:04:05Z"

// LeaseOwnerPID parses a well-formed lease's recorded owner pid. A lease is
// well formed only as one "<pid> <utc-time>\n" line whose pid and instant
// round-trip exactly; everything else reports false so consumers fail closed.
func LeaseOwnerPID(content []byte) (int, bool) {
	if len(content) == 0 || content[len(content)-1] != '\n' || strings.Count(string(content), "\n") != 1 {
		return 0, false
	}
	fields := strings.Split(string(content[:len(content)-1]), " ")
	if len(fields) != 2 {
		return 0, false
	}
	pid, err := strconv.Atoi(fields[0])
	if err != nil || pid <= 0 || strconv.Itoa(pid) != fields[0] {
		return 0, false
	}
	stamp, err := time.Parse(LeaseTimeLayout, fields[1])
	return pid, err == nil && stamp.Format(LeaseTimeLayout) == fields[1]
}

// Reclaimable requires a dead recorded pid or a lease aged past staleAfter,
// the caller-supplied window that separates a crashed legacy lease from a
// fresh writer mid-claim. The liveness probe and the window are supplied
// facts, so this decision never signals a process or owns a bound itself.
func Reclaimable(content []byte, mtime, now time.Time, alive func(int) bool, staleAfter time.Duration) bool {
	if pid, ok := LeaseOwnerPID(content); ok {
		return !alive(pid)
	}
	return now.Sub(mtime) > staleAfter
}

// NestedState names what a checkout's nested and embedded repository scan found.
type NestedState string

const (
	NestedClean         NestedState = "clean"
	NestedDirty         NestedState = "dirty"
	NestedEmbeddedClean NestedState = "embedded-clean"
	NestedEmbeddedDirty NestedState = "embedded-dirty"
	NestedUnknown       NestedState = "unknown"
)

// Orphaned reports whether an assignment has been abandoned by the session that
// cut it. Age is the whole test: nothing records liveness for a request-created
// worktree. staleAfter, the caller-supplied window, is the only thing
// separating a long-running one from residue. Every consumer must ask this one
// question, so the window has a single meaning.
//
// An absent stamp is aged, because a record written before the field existed
// carries none and would otherwise be immortal. A stamp the reading host's
// clock has not reached yet is not aged, so skew cannot manufacture an orphan.
// An unparseable stamp is unknown age rather than infinite age.
// ValidateAssignment rejects one on every ledger read, so a record reaching
// here with one never came from the ledger.
func Orphaned(a intent.Assignment, now time.Time, staleAfter time.Duration) bool {
	if a.State != intent.StateActive {
		return false
	}
	if a.CreatedAt == nil {
		return true
	}
	created, err := time.Parse(time.RFC3339, *a.CreatedAt)
	if err != nil {
		return false
	}
	return now.Sub(created) > staleAfter
}

// Residual reports whether a record preserves no work and is therefore safe to
// compact. Its recovery set is the single source of that judgment: an empty set
// means residue. A non-empty set means preserved work that must never be
// silently discarded. Both the release reconcile and the resume sweep consult
// this one predicate.
func Residual(a intent.Assignment) bool { return len(a.Recovery) == 0 }

// LandednessKind names which of the four ways the explicit planner can know a
// branch's relationship to the default ref. The branch may be never resolvable
// (detached), resolvable but never asked (no default), asked and refused by the
// query itself, or asked and answered.
type LandednessKind int

const (
	LandednessDetached LandednessKind = iota
	LandednessUnknownNoDefault
	LandednessUnknownError
	LandednessProven
)

// Landedness carries a branch's proven relationship to the default ref as typed
// fields. Every decision reads those. String renders the "true:ancestry" /
// "unknown:<err>" wire form, whose only consumer is the explicit fingerprint,
// which hashes it as evidence.
type Landedness struct {
	Kind      LandednessKind
	Err       string
	Landed    bool
	ByContent bool
}

// ProvenLanded reports whether the query was actually asked and answered yes.
// That is the one reading that authorizes acting on a landing, rather than
// detachment, an absent default ref, a failed query, or a proven non-landing.
func (l Landedness) ProvenLanded() bool {
	return l.Kind == LandednessProven && l.Landed
}

func (l Landedness) String() string {
	switch l.Kind {
	case LandednessDetached:
		return "detached"
	case LandednessUnknownNoDefault:
		return "unknown"
	case LandednessUnknownError:
		return "unknown:" + l.Err
	case LandednessProven:
		proof := "ancestry"
		if l.ByContent {
			proof = "patch"
		}
		return boolLabel(l.Landed) + ":" + proof
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

// RecoveryLookupKind names whether a recovery ref still needs a git query, and
// which query, once the verdict has decided a ref is needed at all. Predicting
// a ref is not itself a decision: the parent performs the named lookup exactly
// once, after the eligibility call, using whichever kind the verdict returns.
type RecoveryLookupKind int

const (
	RecoveryNoLookup RecoveryLookupKind = iota
	RecoveryLookupOwned
	RecoveryLookupForeign
)

// ExplicitFacts carries every typed fact the explicit planner gathers before it
// can answer "ours and safe to remove". Each field is evidence, not a
// conclusion: the same facts feed DecideExplicit whether the eventual verdict
// retains, removes, recovers, or discards. A field left at its zero value was
// never gathered, because an earlier fact already made it inapplicable. The
// assignment ledger, for one, is read only when a marker validated.
type ExplicitFacts struct {
	RegistrationBranchRef  string
	RegistrationLockReason string
	RegistrationLocked     bool
	RegistrationDetached   bool

	MarkerPresent       bool
	MarkerErr           error
	AssignmentLedgerErr error
	AssignmentAmbiguous bool
	MatchedAssignment   *intent.Assignment
	ForeignAssignment   *intent.Assignment
	// AssignmentLockReason is the Bench lock reason the matched assignment's
	// bundle prescribes, rendered at the boundary. The decision only compares
	// it to the registration's recorded reason.
	AssignmentLockReason string

	LeasePresent bool
	LeaseState   LeaseState
	LeaseStatErr error

	NestedState NestedState
	NestedErr   error

	BuildOutputErr   error
	IgnoredErr       error
	IgnoredOverLimit bool
	IgnoredCount     int
	DeclaredIgnored  bool
	DiscardIgnored   bool

	InitialTracked string

	HeadDetached    bool
	DefaultKnown    bool
	LandedOK        bool
	LandedByContent bool
	LandedErr       error
	HeadRef, Head   string

	UnsafeTarget bool
}

// ExplicitVerdict is the decided answer plus every piece of evidence the
// explicit planner needs to project it onto its plan. For a RecoverRemove or
// DiscardRemove verdict, it also finishes resolving the recovery ref. Every
// field is the one-time output of a single DecideExplicit call. No consumer
// decides over it again.
type ExplicitVerdict struct {
	Action     Action
	ReasonCode Reason
	Reason     string

	Owned      bool
	Assignment *intent.Assignment

	Tracked string

	Landed       Landedness
	DeleteBranch bool
	BranchRef    string
	BranchOID    string

	RecoveryLookup RecoveryLookupKind
	Recovery       string
}

// DecideExplicit is the single place the explicit ordered eligibility decision
// is made. That decision covers ownership and marker agreement, assignment
// join, lock agreement, lease liveness, and nested or embedded repository
// state. It also covers ignored-residue authorization, the
// recover/discard-remove promotion, and the final unsafe-target override. The
// order and the last-write-wins collisions between blocks are pinned by
// TestExplicitEligibilityOutcomeMatrix and must not change here without that
// characterization moving first.
func DecideExplicit(f ExplicitFacts) ExplicitVerdict {
	v := ExplicitVerdict{Action: ActionRemove, Tracked: f.InitialTracked, Recovery: "none"}

	switch {
	case f.MarkerPresent:
		switch {
		case f.MarkerErr != nil:
			v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonMalformed, f.MarkerErr.Error()
		case f.AssignmentLedgerErr != nil:
			v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonMalformed, "assignment ledger is unreadable"
		case f.AssignmentAmbiguous:
			v.Assignment = f.MatchedAssignment
			v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonMalformed, "registration has ambiguous assignments"
		case f.MatchedAssignment == nil:
			v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonMalformed, "owner marker has no matching assignment"
		default:
			v.Owned = true
			v.Assignment = f.MatchedAssignment
			if f.RegistrationBranchRef != f.MatchedAssignment.Branch {
				v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonUncertain, "assignment does not match current branch"
			} else if f.RegistrationLockReason != f.AssignmentLockReason {
				v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonUnexpectedLock, "assignment does not match current Bench lock"
			}
		}
	case f.RegistrationLocked:
		v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonUnexpectedLock, "foreign or unexpected lock is retained"
	default:
		v.Assignment = f.ForeignAssignment
	}

	if f.LeasePresent {
		switch f.LeaseState {
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
			v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonUncertain, UnknownLeaseReason
		}
	} else if f.LeaseStatErr != nil && v.Owned {
		v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonUncertain, UnknownLeaseReason
	}

	if f.NestedErr != nil {
		v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonUncertain, "nested repository state is unknown"
	} else if f.NestedState == NestedDirty {
		v.Tracked = "nested-dirty"
		v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonUncertain, "nested repository or submodule is dirty"
	} else if f.NestedState == NestedEmbeddedClean || f.NestedState == NestedEmbeddedDirty {
		v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonUncertain, "embedded repository is retained"
	}

	if f.BuildOutputErr != nil {
		v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonMalformed, "build-output declaration is malformed"
	} else if f.IgnoredErr != nil {
		v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonUncertain, "ignored inventory is uncertain"
	} else if f.IgnoredOverLimit {
		v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonIgnored, "ignored inventory exceeds the destructive limit"
	} else if f.IgnoredCount > 0 && !f.DiscardIgnored && !f.DeclaredIgnored {
		v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonIgnored, "ignored residuals require --discard-ignored"
	}

	switch {
	case f.HeadDetached:
		v.Landed = Landedness{Kind: LandednessDetached}
	case !f.DefaultKnown:
		v.Landed = Landedness{Kind: LandednessUnknownNoDefault}
	case f.LandedErr != nil:
		v.Landed = Landedness{Kind: LandednessUnknownError, Err: f.LandedErr.Error()}
	default:
		v.Landed = Landedness{Kind: LandednessProven, Landed: f.LandedOK, ByContent: f.LandedByContent}
		if f.LandedOK {
			v.DeleteBranch = true
			v.BranchRef, v.BranchOID = f.HeadRef, f.Head
		}
	}

	if v.Action != ActionRetain && (v.Tracked != "clean" || f.RegistrationDetached || (v.Assignment != nil && len(v.Assignment.Recovery) > 0)) {
		v.Action = ActionRecoverRemove
		switch {
		case v.Owned && v.Assignment != nil:
			if len(v.Assignment.Recovery) > 0 {
				v.Recovery = v.Assignment.Recovery[0].Ref
			} else {
				v.RecoveryLookup = RecoveryLookupOwned
			}
		case v.Assignment != nil:
			v.Recovery = v.Assignment.Recovery[0].Ref
		default:
			v.RecoveryLookup = RecoveryLookupForeign
		}
	}

	if v.Action != ActionRetain && f.IgnoredCount > 0 && (f.DiscardIgnored || f.DeclaredIgnored) {
		v.Action = ActionDiscardRemove
	}

	if f.UnsafeTarget {
		v.Action, v.ReasonCode, v.Reason = ActionRetain, ReasonUncertain, "target contains unsafe control bytes"
	}

	return v
}

// ExplicitOutcome is the slice of an already-decided explicit plan the
// automatic decision reads. The parent projects it from its plan once, so the
// automatic decision never reaches back into effectful plan state.
type ExplicitOutcome struct {
	Action     Action
	ReasonCode Reason
	Reason     string

	HasAssignment        bool
	Owned                bool
	AssignmentID         string
	AssignmentState      intent.AssignmentState
	Tracked              string
	RegistrationDetached bool
	Landed               Landedness
}

// AutomaticFacts carries the explicit planner's own outcome plus every
// automatic-specific fact the automatic planner gathers on top of it.
// DecideAutomatic uses these to answer "is this eligible for unattended,
// automatic cleanup, and if not, why not". ExplicitErr and Explicit are
// mutually exclusive with the MissingBranch* fields. The MissingBranch* fields
// are gathered only when explicit planning itself failed and the automatic
// planner salvages one narrow fact from that failure. Every other field is
// gathered only when explicit planning succeeded, mirroring the planner's own
// conditional evidence-gathering. A field left at its zero value simply was
// never gathered, because an earlier fact already made it inapplicable.
type AutomaticFacts struct {
	ExplicitErr error
	Explicit    ExplicitOutcome

	// MissingBranchAssignmentID is non-empty when explicit planning failed but
	// an active assignment with an unresolvable branch was salvaged.
	MissingBranchAssignmentID string
	MissingBranchLiveLease    bool

	LiveLease       bool
	Landed          bool
	OrphanedActive  bool
	RecoveryMatches bool
}

// AutomaticVerdict is the decided answer for the automatic, unattended cleanup
// path. Action, ReasonCode, and Reason are the decision itself. AssignmentID
// echoes the one piece of evidence the automatic planner projects onto its
// plan's Assignment field that DecideExplicit never sets. It stays empty on
// every branch that reaches its verdict without a verified assignment to name,
// so the projection leaves the plan's Assignment untouched there.
type AutomaticVerdict struct {
	Action       Action
	ReasonCode   Reason
	Reason       string
	AssignmentID string
}

// DecideAutomatic is the single place the automatic ordered eligibility
// decision is made. That decision covers the explicit-planning-error salvage
// (an active assignment whose branch cannot be resolved), the live-lease
// override, and retain-passthrough with a landed-reason swap. It also covers
// the foreign/unowned refusal, the not-cleanup-pending reasons (with the
// orphaned-age override), the recovery-metadata-match check, the
// unknown/unmerged landedness checks, and the final preservation refusal. The
// order and every message are pinned by TestAutomaticEligibilityOutcomeMatrix
// and must not change here without that characterization moving first.
func DecideAutomatic(f AutomaticFacts) AutomaticVerdict {
	if f.ExplicitErr != nil {
		if f.MissingBranchAssignmentID != "" {
			v := AutomaticVerdict{Action: ActionRetain, ReasonCode: ReasonActive, Reason: "assignment landedness is unknown", AssignmentID: f.MissingBranchAssignmentID}
			if f.MissingBranchLiveLease {
				v.ReasonCode, v.Reason = ReasonLiveLease, "assignment has a live lease"
			}
			return v
		}
		reason := ReasonUncertain
		if strings.Contains(f.ExplicitErr.Error(), "assignment") || strings.Contains(f.ExplicitErr.Error(), "intent ledger") {
			reason = ReasonMalformed
		}
		return AutomaticVerdict{Action: ActionRetain, ReasonCode: reason, Reason: f.ExplicitErr.Error()}
	}

	plan := f.Explicit
	if plan.HasAssignment && f.LiveLease {
		return AutomaticVerdict{Action: ActionRetain, ReasonCode: ReasonLiveLease, Reason: "assignment has a live lease"}
	}
	if plan.Action == ActionRetain {
		if plan.HasAssignment && f.Landed {
			return AutomaticVerdict{Action: ActionRetain, ReasonCode: ReasonLanded, Reason: "assignment branch has landed"}
		}
		return AutomaticVerdict{Action: plan.Action, ReasonCode: plan.ReasonCode, Reason: plan.Reason}
	}
	if !plan.HasAssignment || !plan.Owned {
		return AutomaticVerdict{Action: ActionRetain, ReasonCode: ReasonForeign, Reason: "registration is not a verified owned assignment"}
	}

	assignmentID := plan.AssignmentID
	if plan.AssignmentState != intent.StateCleanupPending {
		reason := ReasonUncertain
		if plan.AssignmentState == intent.StateActive {
			if f.Landed {
				reason = ReasonLanded
			} else {
				reason = ReasonActive
			}
			if reason == ReasonActive && f.OrphanedActive {
				reason = ReasonOrphaned
			}
		}
		return AutomaticVerdict{Action: ActionRetain, ReasonCode: reason, Reason: "assignment is not cleanup-pending", AssignmentID: assignmentID}
	}
	if !f.RecoveryMatches {
		return AutomaticVerdict{Action: ActionRetain, ReasonCode: ReasonMalformed, Reason: "assignment recovery metadata does not match refs", AssignmentID: assignmentID}
	}
	if plan.Landed.Kind == LandednessUnknownNoDefault || plan.Landed.Kind == LandednessUnknownError {
		return AutomaticVerdict{Action: ActionRetain, ReasonCode: ReasonUncertain, Reason: "assignment landedness is unknown", AssignmentID: assignmentID}
	}
	if !plan.Landed.ProvenLanded() {
		return AutomaticVerdict{Action: ActionRetain, ReasonCode: ReasonUnmerged, Reason: "assignment branch has not landed", AssignmentID: assignmentID}
	}
	// The automatic path authors no preservation refs: it runs unattended at every session
	// start and through every release, and the standing cleaner sweeps the namespace such a
	// ref would live in, so preserving there would write work nothing can hand back.
	// Disposing of the checkout stays with the operator's explicit path-addressed clean.
	if retain, action, reasonCode, reason := AutomaticPreservation(plan, "automatic cleanup does not preserve uncommitted work"); retain {
		return AutomaticVerdict{Action: action, ReasonCode: reasonCode, Reason: reason, AssignmentID: assignmentID}
	}
	return AutomaticVerdict{Action: plan.Action, ReasonCode: plan.ReasonCode, Reason: plan.Reason, AssignmentID: assignmentID}
}

// AutomaticPreservation is the one place automatic-flavored policy decides
// whether a plan must be retained to avoid stranding uncommitted work, and what
// Action/ReasonCode that refusal carries. DecideAutomatic's own dirty-refusal
// branch and the landed-set's retained preservation both call this. Neither
// writes its own ActionRetain/ReasonDirty literal, so the two routes can never
// derive "would removing this strand uncommitted work" differently. Each caller
// still supplies its own operator-facing message for its own command surface.
func AutomaticPreservation(o ExplicitOutcome, message string) (retain bool, action Action, reasonCode Reason, reason string) {
	if !Preserves(o.Action, o.Tracked, o.RegistrationDetached) {
		return false, o.Action, o.ReasonCode, o.Reason
	}
	return true, ActionRetain, ReasonDirty, message
}
