package commitment_test

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitcmd"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/intent"
)

// An open listed continuation holds the default active slot. A start for the first outcome
// refuses until a parallel grant names that continuation's run, the delivery of its scope
// prunes it, or its run is no longer active. A run that ends frees the slot at once, with
// no reconciliation first. In each case the projection reads the same decision as
// admission: it offers the start only when admission accepts it, and it never reads
// all-blocked.
func TestCommitmentContinuationOccupiesSlot(t *testing.T) {
	for _, row := range []struct {
		name  string
		grant string
		end   func(t *testing.T, root string, run intent.Assignment)
		admit bool
	}{
		{name: "open"},
		{name: "granted", grant: "legacy", admit: true},
		{name: "other-run", grant: "other"},
		{name: "delivered", admit: true},
		{name: "run-complete", end: func(t *testing.T, root string, run intent.Assignment) {
			run.State = intent.StateComplete
			if err := intent.PutAssignment(root, run); err != nil {
				t.Fatal(err)
			}
		}, admit: true},
		{name: "run-purged", end: func(t *testing.T, root string, run intent.Assignment) {
			if _, err := intent.PurgeAssignments(root, func(kept intent.Assignment) bool { return kept.ID != run.ID }); err != nil {
				t.Fatal(err)
			}
		}, admit: true},
	} {
		t.Run(row.name, func(t *testing.T) {
			root, first, _ := admissionRepo(t, false, false)
			runs := map[string]intent.Assignment{}
			for _, request := range []string{"legacy", "other"} {
				path := commitmenttest.Assignment(t, root, request)
				ledger, err := intent.Read(root)
				if err != nil {
					t.Fatal(err)
				}
				owners := intent.AssignmentsOwning(ledger.Assignments, path)
				if len(owners) != 1 {
					t.Fatalf("%s owners = %+v", request, owners)
				}
				runs[request] = owners[0]
			}
			legacy := runs["legacy"]
			if err := intent.Transact(root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
				ledger.Commitment = &intent.CommitmentState{Continuations: []intent.LegacyContinuation{{Assignment: legacy.ID, Request: legacy.Request, Scope: []string{deliverable("C")}}}}
				return ledger, true, nil
			}, nil); err != nil {
				t.Fatal(err)
			}
			if row.end != nil {
				row.end(t, root, legacy)
			}
			if row.grant != "" {
				commitmenttest.EditPolicy(t, root, func(policy *commitment.Policy) {
					policy.ParallelGrants = []commitment.ParallelGrant{{Outcomes: []string{"A"}, Continuations: []string{runs[row.grant].ID}}}
				})
				commitmenttest.Commit(t, root, "grant beside a continuation")
			}
			if row.name == "delivered" {
				commitmenttest.EditPolicy(t, root, func(policy *commitment.Policy) {
					delivered, _, err := commitment.Deliver(*policy, deliverable("C"), "source", "evidence")
					if err != nil {
						t.Fatal(err)
					}
					*policy = delivered
				})
				commitmenttest.Commit(t, root, "deliver the continuation scope")
				if err := (commitrepo.Store{Root: root}).ReconcileDelivered(); err != nil {
					t.Fatal(err)
				}
				state, err := intent.Read(root)
				if err != nil {
					t.Fatal(err)
				}
				if len(state.Commitment.Continuations) != 0 {
					t.Fatalf("delivered scope kept its continuation: %+v", state.Commitment)
				}
			}
			outlook := commitcmd.Outlook(root)
			out, code := startOutcome(first, "A", "first")
			if row.admit {
				if code != 0 || outlook.State != commitment.OutlookEligible || outlook.Next != "A" {
					t.Fatalf("start = (%s,%d), outlook = %+v; want A admitted and offered", out, code, outlook)
				}
				return
			}
			if code != 1 || !strings.Contains(out, "legacy continuation is active") || outlook.State != commitment.OutlookActive || outlook.Next != "" {
				t.Fatalf("start = (%s,%d), outlook = %+v; want the continuation to hold the slot", out, code, outlook)
			}
		})
	}
}
