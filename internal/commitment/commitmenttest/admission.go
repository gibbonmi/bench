package commitmenttest

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/intent"
)

// DeliveryOutcome is the one outcome that SeedAdmission approves.
const DeliveryOutcome = "delivery"

// SeedAdmission writes the approved policy for one existing spec before the fixture's base commit.
func SeedAdmission(t testing.TB, root, deliverable string) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, deliverable))
	if err != nil {
		t.Fatal(err)
	}
	source := commitment.SourceBinding{ID: "spec", Path: deliverable, Identity: commitment.Identity(data)}
	outcome := commitment.Outcome{ID: DeliveryOutcome, Criteria: []commitment.Criterion{{ID: "accepted", Text: "The fixture delivery is accepted."}}, Sources: []commitment.SourceBinding{}, Deliverables: []commitment.DeliveryBinding{{Source: source}}}
	WritePolicy(t, root, commitment.Policy{Version: 1, ActiveMilestone: "fixture", Milestones: []commitment.Milestone{{ID: "fixture", Outcomes: []commitment.Outcome{outcome}}}})
}

// Register records an active assignment for a command fixture that owns its checkout directly.
func Register(t testing.TB, root, request string) {
	t.Helper()
	id, owner := intent.RequestDigest(request)[:32], strings.Repeat("b", 32)
	err := intent.PutAssignment(root, intent.Assignment{Schema: intent.AssignmentRecordSchema, ID: id, OwnerID: owner, Request: intent.RequestDigest(request), Label: request, Start: gittest.Output(t, root, "rev-parse", "HEAD"), Branch: intent.AssignmentBranchRef(owner, id), Worktree: root, State: intent.StateActive})
	if err != nil {
		t.Fatal(err)
	}
}

// Admit exercises the real start owner for a fixture's already registered assignment.
func Admit(t testing.TB, root, request, deliverable string) {
	t.Helper()
	if err := (commitrepo.Store{Root: root}).Start(DeliveryOutcome, request, deliverable); err != nil {
		t.Fatal(err)
	}
}

// Rebind drops the binding of root's assignment and starts it again. A fixture calls it
// after it commits a re-approved edit of the bound deliverable.
func Rebind(t testing.TB, root, request, deliverable string) {
	t.Helper()
	err := intent.Transact(root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
		owners := intent.AssignmentsOwning(ledger.Assignments, root)
		if len(owners) != 1 || ledger.Commitment == nil {
			return ledger, false, nil
		}
		state := *ledger.Commitment
		state.Bindings = slices.DeleteFunc(slices.Clone(state.Bindings), func(binding intent.DeliveryBinding) bool {
			return binding.Assignment == owners[0].ID
		})
		ledger.Commitment = &state
		return ledger, true, nil
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	Admit(t, root, request, deliverable)
}
