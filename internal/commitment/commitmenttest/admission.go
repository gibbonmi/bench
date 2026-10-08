package commitmenttest

import (
	"encoding/json"
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

// Reworded is plan input that changes only the text of the first criterion in the policy
// that root's default branch holds, to text.
func Reworded(t testing.TB, root, text string) []byte {
	t.Helper()
	policy, _, err := (commitrepo.Store{Root: root}).Policy()
	if err != nil {
		t.Fatal(err)
	}
	policy.Milestones[0].Outcomes[0].Criteria[0].Text = text
	data, err := commitment.Bytes(policy)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// LegacyRun is the request of the one run that Legacy creates, and LegacyScope is the
// path that the run's branch holds, its existing scope.
const (
	LegacyRun   = "legacy"
	LegacyScope = "owned.txt"
)

// Legacy creates a repository with no policy, whose board lists one row, and one run
// whose branch holds LegacyScope. It returns the root and the run's checkout.
func Legacy(t testing.TB) (root, run string) {
	t.Helper()
	root = gittest.RepoOnBranch(t, "main")
	Write(t, root, "ROADMAP.md", oldBoard)
	Commit(t, root, "board")
	run = Assignment(t, root, LegacyRun)
	Write(t, run, LegacyScope, "owned\n")
	Commit(t, run, "legacy scope")
	return root, run
}

// Adoption is the plan input of an initial adoption: one policy and the runs it lists.
type Adoption struct {
	commitment.Policy
	Continuations []intent.LegacyContinuation `json:"continuations,omitempty"`
}

// LegacyAdoption is the adoption whose active milestone M1 holds outcome A and that lists
// the Legacy run under the id run, with LegacyScope as its scope.
func LegacyAdoption(run string) Adoption {
	policy := commitment.Policy{Version: 1, ActiveMilestone: "M1", Milestones: []commitment.Milestone{{ID: "M1", Outcomes: []commitment.Outcome{{ID: "A", Criteria: []commitment.Criterion{{ID: "A.done", Text: "The outcome is delivered."}}}}}}}
	return Adoption{policy, []intent.LegacyContinuation{{Assignment: run, Request: intent.RequestDigest(LegacyRun), Scope: []string{LegacyScope}}}}
}

// Encode is the adoption as plan input.
func (adoption Adoption) Encode(t testing.TB) []byte {
	t.Helper()
	data, err := json.Marshal(adoption)
	if err != nil {
		t.Fatal(err)
	}
	return data
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
