package commitment_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/intent"
)

func TestCommitmentStaleApproval(t *testing.T) {
	current := policy([]commitment.Milestone{milestone("M1", "A")}, "M1")
	root := commitmenttest.Repo(t, current)
	store := commitrepo.Store{Root: root}
	proposal, err := commitment.Bytes(policy([]commitment.Milestone{milestone("M1", "A", "B")}, "M1"))
	if err != nil {
		t.Fatal(err)
	}
	plan, err := store.Plan(proposal)
	if err != nil {
		t.Fatal(err)
	}
	commitmenttest.WritePolicy(t, root, policy([]commitment.Milestone{milestone("M1", "C")}, "M1"))
	commitmenttest.Commit(t, root, "replace commitment")

	_, err = store.Approve(plan.ID, "decision-1", plan.Effects.Delayed, plan.Effects.Removed)
	if err == nil || !strings.Contains(err.Error(), "predecessor") {
		t.Fatalf("Approve() error = %v, want changed predecessor refusal", err)
	}
}

func TestCommitmentApprovalReplay(t *testing.T) {
	current := policy([]commitment.Milestone{milestone("M1", "A")}, "M1")
	root := commitmenttest.Repo(t, current)
	store := commitrepo.Store{Root: root}
	proposal, err := commitment.Bytes(policy([]commitment.Milestone{milestone("M1", "A", "B")}, "M1"))
	if err != nil {
		t.Fatal(err)
	}
	plan, err := store.Plan(proposal)
	if err != nil {
		t.Fatal(err)
	}
	changed, err := store.Approve(plan.ID, "decision-1", plan.Effects.Delayed, plan.Effects.Removed)
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("first approval reported no change")
	}
	roadmap, err := os.ReadFile(filepath.Join(root, "ROADMAP.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(roadmap), "## Recommended sequence\n\n1. A\n2. B\n") {
		t.Fatalf("ROADMAP.md = %q, want projected A then B", roadmap)
	}
	policyBefore, err := os.ReadFile(filepath.Join(root, ".bench", "commitment.json"))
	if err != nil {
		t.Fatal(err)
	}
	receiptPath, err := intent.Address(root)
	if err != nil {
		t.Fatal(err)
	}
	receiptBefore, err := os.ReadFile(receiptPath)
	if err != nil {
		t.Fatal(err)
	}

	changed, err = store.Approve(plan.ID, "decision-1", plan.Effects.Delayed, plan.Effects.Removed)
	if err != nil {
		t.Fatal(err)
	}
	if changed {
		t.Fatal("identical replay reported a change")
	}
	policyAfter, _ := os.ReadFile(filepath.Join(root, ".bench", "commitment.json"))
	receiptAfter, _ := os.ReadFile(receiptPath)
	if !bytes.Equal(policyAfter, policyBefore) || !bytes.Equal(receiptAfter, receiptBefore) {
		t.Fatal("identical replay changed policy or receipt bytes")
	}
}

func TestCommitmentApprovalRequiresExactEffects(t *testing.T) {
	current := policy([]commitment.Milestone{milestone("M1", "A", "B")}, "M1")
	root := commitmenttest.Repo(t, current)
	store := commitrepo.Store{Root: root}
	proposal, err := commitment.Bytes(policy([]commitment.Milestone{milestone("M1", "A")}, "M1"))
	if err != nil {
		t.Fatal(err)
	}
	plan, err := store.Plan(proposal)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := store.Approve(plan.ID, "decision-1", nil, nil); err == nil || !strings.Contains(err.Error(), "operands") {
		t.Fatalf("Approve() error = %v, want missing-effect refusal", err)
	}
	ledger, err := intent.Read(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, receipt := range ledger.CommitmentReceipts {
		if receipt.Plan == plan.ID && receipt.Approved {
			t.Fatal("missing-effect refusal marked the plan approved")
		}
	}
}

func TestCommitmentUnrelatedAdvance(t *testing.T) {
	current := policy([]commitment.Milestone{milestone("M1", "A")}, "M1")
	root := commitmenttest.Repo(t, current)
	store := commitrepo.Store{Root: root}
	proposal, err := commitment.Bytes(policy([]commitment.Milestone{milestone("M1", "A", "B")}, "M1"))
	if err != nil {
		t.Fatal(err)
	}
	plan, err := store.Plan(proposal)
	if err != nil {
		t.Fatal(err)
	}
	commitmenttest.Write(t, root, "notes.txt", "unrelated\n")
	commitmenttest.Commit(t, root, "unrelated advance")

	changed, err := store.Approve(plan.ID, "decision-1", plan.Effects.Delayed, plan.Effects.Removed)
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("approval after unrelated advance reported no change")
	}
}

func TestCommitmentConcurrentReplacementAdmitsOnePlan(t *testing.T) {
	current := policy([]commitment.Milestone{milestone("M1", "A")}, "M1")
	root := commitmenttest.Repo(t, current)
	store := commitrepo.Store{Root: root}
	plans := make([]commitment.Plan, 2)
	for i, outcome := range []string{"B", "C"} {
		proposal, err := commitment.Bytes(policy([]commitment.Milestone{milestone("M1", outcome)}, "M1"))
		if err != nil {
			t.Fatal(err)
		}
		plans[i], err = store.Plan(proposal)
		if err != nil {
			t.Fatal(err)
		}
	}

	start := make(chan struct{})
	errs := make([]error, len(plans))
	changed := make([]bool, len(plans))
	var group sync.WaitGroup
	for i := range plans {
		group.Add(1)
		go func(i int) {
			defer group.Done()
			<-start
			changed[i], errs[i] = store.Approve(plans[i].ID, "decision", plans[i].Effects.Delayed, plans[i].Effects.Removed)
		}(i)
	}
	close(start)
	group.Wait()

	admitted := 0
	refused := 0
	for i := range plans {
		if changed[i] && errs[i] == nil {
			admitted++
		}
		if errs[i] != nil && strings.Contains(errs[i].Error(), "predecessor changed") {
			refused++
		}
	}
	if admitted != 1 || refused != 1 {
		t.Fatalf("concurrent approvals changed = %v, errors = %v; want one admission and one predecessor refusal", changed, errs)
	}
}

func TestCommitmentPublishedSourceIdentity(t *testing.T) {
	current := policy([]commitment.Milestone{milestone("M1", "A")}, "M1")
	root := commitmenttest.Repo(t, current)
	commitmenttest.Write(t, root, "obligation.md", "original obligation\n")
	commitmenttest.Commit(t, root, "publish source")
	proposed := policy([]commitment.Milestone{milestone("M1", "A", "B")}, "M1")
	proposed.Milestones[0].Outcomes[0].Sources = []commitment.SourceBinding{{ID: "FT1", Path: "obligation.md", Identity: commitment.Identity([]byte("original obligation\n"))}}
	payload, err := commitment.Bytes(proposed)
	if err != nil {
		t.Fatal(err)
	}
	store := commitrepo.Store{Root: root}
	plan, err := store.Plan(payload)
	if err != nil {
		t.Fatal(err)
	}
	commitmenttest.Write(t, root, "obligation.md", "changed published obligation\n")
	commitmenttest.Commit(t, root, "change published source")
	commitmenttest.Write(t, root, "obligation.md", "original obligation\n")
	if _, err := store.Approve(plan.ID, "decision", nil, nil); err == nil || !strings.Contains(err.Error(), "source") {
		t.Fatalf("approval with stale working copy = %v, want changed published source refusal", err)
	}
}
