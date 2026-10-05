package repository_test

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/gittest"
)

// A published delivery settles what it delivered: its closure deleted the satisfied row,
// and its publication flipped the delivered spec. A later plan binds only the unsettled
// sources, so it plans, and a changed unsettled source still refuses it.
func TestCommitmentPlanAfterDelivery(t *testing.T) {
	root := commitmenttest.SeedMilestone(t, commitmenttest.MilestoneSpec)
	store := commitrepo.Store{Root: root}
	policy, _, err := store.Policy()
	if err != nil {
		t.Fatal(err)
	}
	successor := &policy.Milestones[1]
	successor.Outcomes = append(successor.Outcomes, commitment.Outcome{ID: "D", Criteria: []commitment.Criterion{{ID: "D-done", Text: "The D obligation is satisfied."}}})
	proposal, err := commitment.Bytes(policy)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Plan(proposal); err != nil {
		t.Fatalf("Plan after delivery = %v, want the settled row and spec unbound", err)
	}
	commitmenttest.Write(t, root, "roadmap/FT2.md", "Changed B obligation.\n")
	commitmenttest.Commit(t, root, "change open row")
	if _, err := store.Plan(proposal); err == nil || !strings.Contains(err.Error(), `"FT2" changed`) {
		t.Fatalf("Plan with a changed open row = %v, want the changed source refusal", err)
	}
}

// A plan does not bind a deliverable that its proposal drops or removes, or the old
// identity of a deliverable that it approves again. A later main commit that deletes or
// changes that deliverable therefore leaves the plan, its approval, and its commit valid.
func TestCommitmentPlanSurvivesRemovedDeliverable(t *testing.T) {
	for _, row := range []struct {
		name    string
		change  func(testing.TB, string)
		propose func(t *testing.T, root string, policy *commitment.Policy)
		removed []string
	}{
		{name: "dropped tickets binding", change: commitmenttest.RemoveTickets, propose: func(_ *testing.T, _ string, policy *commitment.Policy) {
			policy.Milestones[0].Outcomes[1].Deliverables = nil
		}},
		{name: "removed outcome", change: commitmenttest.RemoveTickets, propose: func(_ *testing.T, _ string, policy *commitment.Policy) {
			policy.Milestones[0].Outcomes = policy.Milestones[0].Outcomes[:1]
		}, removed: []string{"B"}},
		{name: "changed staged spec", change: func(t testing.TB, root string) {
			commitmenttest.Write(t, root, commitmenttest.MilestoneSpec, "# x\n\nStatus: staged\n\nA changed scope.\n")
			commitmenttest.Commit(t, root, "change staged spec")
		}, propose: func(t *testing.T, root string, policy *commitment.Policy) {
			identity, err := commitrepo.SourceIdentity(root, "main", commitmenttest.MilestoneSpec)
			if err != nil {
				t.Fatal(err)
			}
			policy.Milestones[0].Outcomes[0].Deliverables[0].Source.Identity = identity
		}},
	} {
		t.Run(row.name, func(t *testing.T) {
			root := commitmenttest.SeedMilestone(t)
			row.change(t, root)
			store := commitrepo.Store{Root: root}
			policy, _, err := store.Policy()
			if err != nil {
				t.Fatal(err)
			}
			row.propose(t, root, &policy)
			proposal, err := commitment.Bytes(policy)
			if err != nil {
				t.Fatal(err)
			}
			planning := commitmenttest.Planning(t, root)
			plan, err := store.Plan(proposal)
			if err != nil {
				t.Fatalf("Plan = %v, want the changed deliverable unbound", err)
			}
			if _, err := (commitrepo.Store{Root: planning}).Approve(plan.ID, "decision", nil, row.removed); err != nil {
				t.Fatalf("Approve = %v, want the plan staged", err)
			}
			commitmenttest.Commit(t, planning, "approved plan")
			if err := (commitrepo.Store{Root: planning}).AuthorizeCandidate(gittest.Output(t, planning, "rev-parse", "HEAD^{tree}")); err != nil {
				t.Fatalf("AuthorizeCandidate = %v, want the approved plan admitted", err)
			}
		})
	}
}
