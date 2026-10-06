package commitment_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitcmd"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/intent"
)

func TestCommitmentActivate(t *testing.T) {
	policy := commitment.Policy{
		Version:         1,
		ActiveMilestone: "M1",
		Milestones: []commitment.Milestone{{
			ID: "M1",
			Outcomes: []commitment.Outcome{
				{ID: "A"},
				{ID: "B"},
			},
		}},
	}

	got := commitment.Selection(policy)
	want := commitment.Projection{Milestone: "M1", Outcomes: []string{"A", "B"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Selection() = %#v, want %#v", got, want)
	}
}

func TestCommitmentLiteralInput(t *testing.T) {
	root := commitmenttest.Repo(t, policy([]commitment.Milestone{milestone("M1", "A")}, "M1"))
	proposal, err := commitment.Bytes(policy([]commitment.Milestone{milestone("M1", "A", "B")}, "M1"))
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "proposal [*].json")
	if err := os.WriteFile(path, proposal, 0o644); err != nil {
		t.Fatal(err)
	}

	out, code := commitcmd.Command(root, []string{"plan", "--input", path})
	if code != 0 || !strings.Contains(out, "commitment_plan[1]") {
		t.Fatalf("Command(plan literal) = (%q, %d), want plan", out, code)
	}
}

func TestCommitmentInputFraming(t *testing.T) {
	root := commitmenttest.Repo(t, policy([]commitment.Milestone{milestone("M1", "A")}, "M1"))
	proposal, err := commitment.Bytes(policy([]commitment.Milestone{milestone("M1", "A", "B")}, "M1"))
	if err != nil {
		t.Fatal(err)
	}
	withNewline := filepath.Join(root, "with-newline.json")
	withoutNewline := filepath.Join(root, "without-newline.json")
	if err := os.WriteFile(withNewline, proposal, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(withoutNewline, proposal[:len(proposal)-1], 0o644); err != nil {
		t.Fatal(err)
	}

	withOutput, withCode := commitcmd.Command(root, []string{"plan", "--input", withNewline})
	withoutOutput, withoutCode := commitcmd.Command(root, []string{"plan", "--input", withoutNewline})
	if withCode != 0 || withoutCode != 0 || withOutput != withoutOutput {
		t.Fatalf("framed plan = (%q, %d), unframed plan = (%q, %d)", withOutput, withCode, withoutOutput, withoutCode)
	}
}

// The plan command refuses a binding with an empty obligation list before it records a
// receipt. The row and the staged spec exist at main, so only the refusal stops the plan.
func TestCommitmentPlanCommandRefusesEmptyObligations(t *testing.T) {
	root := commitmenttest.Repo(t, policy([]commitment.Milestone{milestone("M1", "A")}, "M1"))
	proposed := policy([]commitment.Milestone{obligationFreeMilestone(milestone("M1", "A"), "FT1")}, "M1")
	outcome := &proposed.Milestones[0].Outcomes[0]
	for _, source := range []*commitment.SourceBinding{&outcome.Sources[0], &outcome.Deliverables[0].Source} {
		commitmenttest.Write(t, root, source.Path, commitmenttest.StagedBody)
		source.Identity = commitment.Identity([]byte(commitmenttest.StagedBody))
	}
	commitmenttest.Commit(t, root, "write the row and the staged spec")
	data, err := json.Marshal(proposed)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "proposal.json")
	if err := os.WriteFile(path, []byte(strings.Replace(string(data), `"source":{`, `"obligations":[],"source":{`, 1)), 0o644); err != nil {
		t.Fatal(err)
	}

	out, code := commitcmd.Command(root, []string{"plan", "--input", path})
	if code != 1 || !strings.Contains(out, "names no obligation") {
		t.Errorf("Command(plan empty obligations) = (%q, %d), want the obligation refusal", out, code)
	}
	ledger, err := intent.Read(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(ledger.CommitmentReceipts) != 0 {
		t.Errorf("refused plan recorded receipts %v", ledger.CommitmentReceipts)
	}
}

func TestCommitmentGrammar(t *testing.T) {
	root := commitmenttest.Repo(t, policy([]commitment.Milestone{milestone("M1", "A")}, "M1"))
	out, code := commitcmd.Command(root, []string{"plan", "--input", "approve"})
	if code != 1 || !strings.Contains(out, "bench commitment plan refused") || strings.Contains(out, "approval") {
		t.Fatalf("Command(plan --input approve) = (%q, %d), want plan refusal", out, code)
	}
}

func TestCommitmentControlInput(t *testing.T) {
	root := commitmenttest.Repo(t, policy([]commitment.Milestone{milestone("M1", "A")}, "M1"))
	proposal, err := commitment.Bytes(policy([]commitment.Milestone{milestone("M1", "A", "B")}, "M1"))
	if err != nil {
		t.Fatal(err)
	}
	plan, err := (commitrepo.Store{Root: root}).Plan(proposal)
	if err != nil {
		t.Fatal(err)
	}
	planning := commitmenttest.Planning(t, root)
	receiptPath, err := intent.Address(root)
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(receiptPath)
	if err != nil {
		t.Fatal(err)
	}

	out, code := commitcmd.Command(planning, []string{"approve", "--plan", plan.ID, "--decision", "bad\nreference", "--delayed", "none", "--removed", "none"})
	if code != 1 || !strings.Contains(out, "decision") {
		t.Fatalf("Command(approve control) = (%q, %d), want refusal", out, code)
	}
	after, err := os.ReadFile(receiptPath)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(after, before) {
		t.Fatal("control-bearing decision changed receipt bytes")
	}
}

func TestCommitmentUnsafeInput(t *testing.T) {
	root := commitmenttest.Repo(t, validPolicy())
	regular := filepath.Join(root, "policy.json")
	payload, err := commitment.Bytes(validPolicy())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(regular, payload, 0o644); err != nil {
		t.Fatal(err)
	}
	live := filepath.Join(root, "live.json")
	if err := os.Symlink(regular, live); err != nil {
		t.Fatal(err)
	}
	dangling := filepath.Join(root, "dangling.json")
	if err := os.Symlink(filepath.Join(root, "missing.json"), dangling); err != nil {
		t.Fatal(err)
	}
	fifo := filepath.Join(root, "policy.fifo")
	if err := syscall.Mkfifo(fifo, 0o600); err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{live, dangling, fifo} {
		out, code := commitcmd.Command(root, []string{"plan", "--input", path})
		if code != 1 || !strings.Contains(out, "refused") || !strings.Contains(out, "input") {
			t.Fatalf("Command(plan %q) = (%q, %d), want refusal", path, out, code)
		}
	}
}

func TestCommitmentPlannedMilestone(t *testing.T) {
	policy := commitment.Policy{
		Version:         1,
		ActiveMilestone: "M1",
		Milestones: []commitment.Milestone{
			{ID: "M1", Outcomes: []commitment.Outcome{{ID: "A"}}},
			{ID: "M2", Outcomes: []commitment.Outcome{{ID: "B"}}},
		},
	}

	got := commitment.Selection(policy)
	want := commitment.Projection{Milestone: "M1", Outcomes: []string{"A"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Selection() = %#v, want %#v", got, want)
	}
}

func TestCommitmentApprovalRequiresOwnedWorktree(t *testing.T) {
	root := commitmenttest.Repo(t, policy([]commitment.Milestone{milestone("M1", "A")}, "M1"))
	payload, err := commitment.Bytes(policy([]commitment.Milestone{milestone("M1", "A", "B")}, "M1"))
	if err != nil {
		t.Fatal(err)
	}
	plan, err := (commitrepo.Store{Root: root}).Plan(payload)
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(filepath.Join(root, ".bench", "commitment.json"))
	if err != nil {
		t.Fatal(err)
	}
	out, code := commitcmd.Command(root, []string{"approve", "--plan", plan.ID, "--decision", "decision", "--delayed", "none", "--removed", "none"})
	after, err := os.ReadFile(filepath.Join(root, ".bench", "commitment.json"))
	if err != nil {
		t.Fatal(err)
	}
	if code != 1 || !strings.Contains(out, "owned planning worktree") || !reflect.DeepEqual(before, after) {
		t.Fatalf("unowned approval = (%q, %d), policy unchanged=%v", out, code, reflect.DeepEqual(before, after))
	}
}
