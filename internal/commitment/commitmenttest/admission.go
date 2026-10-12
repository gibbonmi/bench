package commitmenttest

import (
	"bytes"
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

// StagedInParallel is the Staged repository of outcomes. With more than one outcome, a
// committed parallel grant lets the outcomes start together, so each one can hold its own
// open build.
func StagedInParallel(t testing.TB, outcomes ...string) string {
	t.Helper()
	root := Staged(t, outcomes...)
	if len(outcomes) > 1 {
		EditPolicy(t, root, func(policy *commitment.Policy) {
			policy.ParallelGrants = []commitment.ParallelGrant{{Outcomes: outcomes}}
		})
		Commit(t, root, "grant the outcomes in parallel")
	}
	return root
}

// BindBuild starts the delivery outcome with deliverable from a new linked checkout of root
// that the assignment of request owns. It returns that checkout. The policy at main must
// approve deliverable.
func BindBuild(t testing.TB, root, request, deliverable string) string {
	t.Helper()
	build := Assignment(t, root, request)
	Admit(t, build, request, deliverable)
	return build
}

// EditCommitment applies edit to the commitment state of the ledger at root.
func EditCommitment(t testing.TB, root string, edit func(*intent.CommitmentState)) {
	t.Helper()
	err := intent.Transact(root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
		state := intent.CommitmentState{}
		if ledger.Commitment != nil {
			state = *ledger.Commitment
		}
		edit(&state)
		ledger.Commitment = &state
		return ledger, true, nil
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
}

// PutGap replaces the chunk gap record of gap's outcome in the ledger at root, as the chunk
// checkpoint and the chunk dispatch write it. An empty milestone takes the policy's active
// milestone.
func PutGap(t testing.TB, root string, gap intent.ChunkGap) {
	t.Helper()
	if gap.Milestone == "" {
		policy, _, err := (commitrepo.Store{Root: root}).Policy()
		if err != nil {
			t.Fatal(err)
		}
		gap.Milestone = policy.ActiveMilestone
	}
	EditCommitment(t, root, func(state *intent.CommitmentState) {
		state.Gaps = append(slices.DeleteFunc(slices.Clone(state.Gaps), func(held intent.ChunkGap) bool {
			return held.Outcome == gap.Outcome
		}), gap)
	})
}

// OpenChunkGaps records an open gap of chunk C1 for each claimed outcome of the ledger at
// root, as a green chunk checkpoint records it. A spec-less landing from a bound source
// passes the build window only in a gap.
func OpenChunkGaps(t testing.TB, root string) {
	t.Helper()
	ledger, err := intent.Read(root)
	if err != nil || ledger.Commitment == nil {
		t.Fatalf("ledger at %s = %+v, %v; want claimed outcomes", root, ledger.Commitment, err)
	}
	for _, claim := range ledger.Commitment.Claims {
		PutGap(t, root, intent.ChunkGap{Milestone: claim.Milestone, Outcome: claim.Outcome, Chunk: "C1", Open: true})
	}
}

// LedgerPath is the intent ledger file of the repository at root.
func LedgerPath(t testing.TB, root string) string {
	t.Helper()
	path, err := intent.Address(root)
	if err != nil {
		t.Fatal(err)
	}
	return path
}

// RequireLedgerKept runs act and fails unless the ledger at root keeps its exact bytes.
func RequireLedgerKept(t testing.TB, root string, act func()) {
	t.Helper()
	path := LedgerPath(t, root)
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	act()
	if after, err := os.ReadFile(path); err != nil || !bytes.Equal(before, after) {
		t.Fatalf("ledger = %s (%v), want the bytes %s", after, err, before)
	}
}

// ClaimedGaps is the chunk gap list of the ledger at root, and the one claim of that ledger.
func ClaimedGaps(t testing.TB, root string) ([]intent.ChunkGap, intent.OutcomeClaim) {
	t.Helper()
	ledger, err := intent.Read(root)
	if err != nil || ledger.Commitment == nil || len(ledger.Commitment.Claims) != 1 {
		t.Fatalf("ledger at %s = %+v, %v; want one claim", root, ledger.Commitment, err)
	}
	return ledger.Commitment.Gaps, ledger.Commitment.Claims[0]
}

// Leases is the landing lease list of the ledger at root. It fails no test, so a
// goroutine at a gate barrier can read it.
func Leases(root string) ([]intent.LandingLease, error) {
	ledger, err := intent.Read(root)
	if err != nil || ledger.Commitment == nil {
		return nil, err
	}
	return ledger.Commitment.Leases, nil
}

// Overrides is the reviewer override list of the ledger at root.
func Overrides(root string) ([]intent.LandingOverride, error) {
	ledger, err := intent.Read(root)
	if err != nil || ledger.Commitment == nil {
		return nil, err
	}
	return ledger.Commitment.Overrides, nil
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
