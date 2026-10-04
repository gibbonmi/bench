package repository_test

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
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
