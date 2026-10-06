package commitment_test

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitcmd"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/intent"
)

func TestCommitmentRemovalBindsRemovedSource(t *testing.T) {
	current := policy([]commitment.Milestone{milestone("M1", "A", "B")}, "M1")
	original := "original B obligation\n"
	current.Milestones[0].Outcomes[1].Sources = []commitment.SourceBinding{{ID: "FT2", Path: "obligation.md", Identity: commitment.Identity([]byte(original))}}
	root := commitmenttest.Repo(t, current)
	commitmenttest.Write(t, root, "obligation.md", original)
	commitmenttest.Commit(t, root, "publish obligation")
	planning := commitmenttest.Planning(t, root)
	proposal, err := commitment.Bytes(policy([]commitment.Milestone{milestone("M1", "A")}, "M1"))
	if err != nil {
		t.Fatal(err)
	}
	plan, err := (commitrepo.Store{Root: root}).Plan(proposal)
	if err != nil {
		t.Fatal(err)
	}
	commitmenttest.Write(t, root, "obligation.md", "changed B obligation\n")
	commitmenttest.Commit(t, root, "change removed source")
	before, err := os.ReadFile(filepath.Join(planning, ".bench", "commitment.json"))
	if err != nil {
		t.Fatal(err)
	}
	receiptPath, err := intent.Address(planning)
	if err != nil {
		t.Fatal(err)
	}
	receiptBefore, err := os.ReadFile(receiptPath)
	if err != nil {
		t.Fatal(err)
	}
	out, code := commitcmd.Command(planning, []string{"approve", "--plan", plan.ID, "--decision", "decision", "--delayed", "none", "--removed", "B"})
	after, err := os.ReadFile(filepath.Join(planning, ".bench", "commitment.json"))
	if err != nil {
		t.Fatal(err)
	}
	receiptAfter, err := os.ReadFile(receiptPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(receiptAfter) != string(receiptBefore) {
		t.Fatal("refused source approval changed receipts")
	}
	if code != 1 || !strings.Contains(out, "source") || string(before) != string(after) {
		t.Fatalf("changed removed source approval=(%q,%d), policy unchanged=%v; want source refusal without staging", out, code, string(before) == string(after))
	}
}

// A plan that keeps an unsettled deliverable binds it, so a later main commit that deletes
// that tickets-only folder refuses the approval.
func TestCommitmentApprovalBindsKeptDeliverable(t *testing.T) {
	root := commitmenttest.SeedMilestone(t)
	store := commitrepo.Store{Root: root}
	current, _, err := store.Policy()
	if err != nil {
		t.Fatal(err)
	}
	kept := current.Milestones[0].Outcomes[1].Deliverables[0].Source
	current.Milestones = append(current.Milestones, milestone("M3", "D"))
	proposal, err := commitment.Bytes(current)
	if err != nil {
		t.Fatal(err)
	}
	planning := commitmenttest.Planning(t, root)
	plan, err := store.Plan(proposal)
	if err != nil {
		t.Fatal(err)
	}
	commitmenttest.RemoveTickets(t, root)
	staged, err := (commitrepo.Store{Root: planning}).Approve(plan.ID, "decision", nil, nil)
	if staged || err == nil || !strings.Contains(err.Error(), "commitment source "+strconv.Quote(kept.ID)) {
		t.Fatalf("Approve after the kept folder was deleted = (%v, %v), want the refusal of source %q", staged, err, kept.ID)
	}
}
