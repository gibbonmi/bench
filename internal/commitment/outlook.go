package commitment

import (
	"strings"

	"github.com/gibbonmi/bench/internal/intent"
)

// The outlook states. Each state names the decision that the next action needs.
const (
	// OutlookAdoptionRequired means no policy is published. An old recommended sequence
	// stays unapproved input, and the next action plans the initial commitment.
	OutlookAdoptionRequired = "adoption-required"
	// OutlookUnreadable means the published policy or the local runtime record did not read.
	OutlookUnreadable = "unreadable"
	// OutlookInactive means a policy exists, but no milestone is active.
	OutlookInactive = "no-active-milestone"
	// OutlookEligible means one committed outcome can start now.
	OutlookEligible = "eligible"
	// OutlookActive means the claimed outcomes and the open legacy continuations hold every
	// slot that admission allows.
	OutlookActive = "active"
	// OutlookAllBlocked means no remaining outcome can start, and none is active.
	OutlookAllBlocked = "all-blocked"
	// OutlookDelivered means every outcome of the active milestone is delivered.
	OutlookDelivered = "delivered"
)

// The commitment operations that an outlook names as its next action.
const (
	OperationPlan   = "plan"
	OperationShow   = "show"
	OperationStart  = "start"
	OperationVerify = "verify"
)

// Outlook is the one commitment projection that the roadmap, status, and dashboard
// readers render. Project decides it from the policy and the runtime state, and the
// eligible outcome comes from the admission decision itself. A reader therefore never
// selects work that a start refuses.
type Outlook struct {
	State     string
	Milestone string
	// Next is the first unclaimed remaining outcome that admission accepts now.
	Next string
	// Deliverable is the path of the one undelivered deliverable that Next approves. It is
	// empty when Next approves no deliverable or more than one.
	Deliverable string
	// Active holds the remaining outcomes with a runtime claim, in milestone order.
	Active []string
	// Blocked holds each remaining outcome with a recorded blocker, in milestone order.
	Blocked []intent.OutcomeBlocker
	// Waiting holds the other remaining outcomes, in milestone order. Order, a dependency,
	// or the active slot holds each one back.
	Waiting []string
	// Operation names the commitment operation of the next action. It is empty while the
	// active outcomes continue.
	Operation string
	// Command is the complete next-action command. The CLI adapter fills it from its form
	// table, so the domain names no command grammar.
	Command string
}

// Project decides the outlook of policy and state. A nil policy is an absent policy.
func Project(policy *Policy, state intent.CommitmentState) Outlook {
	if policy == nil {
		return Outlook{State: OutlookAdoptionRequired, Operation: OperationPlan}
	}
	selection := Selection(*policy)
	if selection.Milestone == "" {
		return Outlook{State: OutlookInactive, Operation: OperationPlan}
	}
	outlook := Outlook{Milestone: selection.Milestone}
	remaining := Remaining(*policy)
	if len(remaining) == 0 {
		outlook.State, outlook.Operation = OutlookDelivered, OperationVerify
		return outlook
	}
	active := claimed(state)
	reasons := map[string]string{}
	for _, blocker := range state.Blockers {
		reasons[blocker.Outcome] = blocker.Reason
	}
	for _, id := range remaining {
		reason, blocked := reasons[id]
		switch {
		case active[id]:
			outlook.Active = append(outlook.Active, id)
		case blocked:
			outlook.Blocked = append(outlook.Blocked, intent.OutcomeBlocker{Outcome: id, Reason: reason})
		case outlook.Next == "" && eligible(*policy, state, id) == nil:
			outlook.Next = id
		default:
			outlook.Waiting = append(outlook.Waiting, id)
		}
	}
	switch {
	case outlook.Next != "":
		outlook.State, outlook.Operation = OutlookEligible, OperationStart
		outlook.Deliverable = soleDeliverable(*policy, outlook.Next)
	case len(outlook.Active) != 0 || len(OpenContinuations(*policy, state)) != 0:
		outlook.State = OutlookActive
	default:
		outlook.State, outlook.Operation = OutlookAllBlocked, OperationPlan
	}
	return outlook
}

// soleDeliverable returns the path of the one undelivered deliverable that the active
// outcome id approves, or the empty string when it approves none or several.
func soleDeliverable(policy Policy, id string) string {
	outcome, err := ActiveOutcome(policy, id)
	if err != nil {
		return ""
	}
	_, open := settle(policy).unsettled(outcome)
	if len(open) != 1 {
		return ""
	}
	return open[0].Source.Path
}

// Unreadable is the outlook of a policy or runtime record that did not read. Its next
// action shows the refusal that names the cause.
func Unreadable() Outlook {
	return Outlook{State: OutlookUnreadable, Operation: OperationShow}
}

// OutlookTable and BlockerTable name the two tables that every TOON reader prints for an
// outlook, so the roadmap readers and the commitment command show one shape.
const (
	OutlookTable = "commitment_outlook"
	BlockerTable = "commitment_blockers"
)

// OutlookFields names the cells of the outlook table's one row.
var OutlookFields = []string{"state", "active_milestone", "next_outcome", "deliverable", "active", "blocked", "waiting", "command"}

// Cells returns the outlook's row for the OutlookFields table. A list cell joins its
// outcome identities with one space.
func (outlook Outlook) Cells() []string {
	return []string{outlook.State, outlook.Milestone, outlook.Next, outlook.Deliverable, strings.Join(outlook.Active, " "), strings.Join(outlook.BlockedOutcomes(), " "), strings.Join(outlook.Waiting, " "), outlook.Command}
}

// BlockedOutcomes returns the identity of each blocked outcome, in milestone order.
func (outlook Outlook) BlockedOutcomes() []string {
	ids := make([]string, len(outlook.Blocked))
	for i, blocker := range outlook.Blocked {
		ids[i] = blocker.Outcome
	}
	return ids
}

// BlockerFields names the cells of the blocker table that follows the outlook table.
var BlockerFields = []string{"outcome", "reason"}

// BlockerCells returns one row for each blocked outcome.
func (outlook Outlook) BlockerCells() [][]string {
	rows := make([][]string, len(outlook.Blocked))
	for i, blocker := range outlook.Blocked {
		rows[i] = []string{blocker.Outcome, blocker.Reason}
	}
	return rows
}
