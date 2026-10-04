package commitment_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitcmd"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/intent"
)

// Each element that verification compares refuses on its own: the milestone, its
// delivery obligations, the evidence version, the gate evidence, and each criterion's
// evidence reference and reviewer assessment.
func TestCommitmentVerificationEvidence(t *testing.T) {
	root := deliveredMilestone(t)
	for _, row := range []struct {
		name      string
		milestone string
		want      []string
		edit      func(*commitment.MilestoneEvidence)
	}{
		{name: "inactive-milestone", milestone: commitmenttest.SuccessorMilestone, want: []string{"is not the active milestone"}},
		{name: "unsupported-version", want: []string{"unsupported version 2"}, edit: func(evidence *commitment.MilestoneEvidence) { evidence.Version = 2 }},
		{name: "gate-incomplete", want: []string{"gate evidence", "is incomplete"}, edit: func(evidence *commitment.MilestoneEvidence) { evidence.Gate = "" }},
		{name: "gate-unresolved", want: []string{"gate evidence", "does not resolve"}, edit: func(evidence *commitment.MilestoneEvidence) { evidence.Gate = "reviews/absent.md" }},
		{name: "evidence-incomplete", want: []string{"B-done", "is incomplete"}, edit: func(evidence *commitment.MilestoneEvidence) { evidence.Results[1].Evidence = "" }},
		{name: "evidence-unresolved", want: []string{"B-done", "does not resolve"}, edit: func(evidence *commitment.MilestoneEvidence) { evidence.Results[1].Evidence = "specs/absent.md" }},
		{name: "assessment-missing", want: []string{"B-done", "has no reviewer outcome assessment"}, edit: func(evidence *commitment.MilestoneEvidence) { evidence.Results[1].Assessment = " " }},
	} {
		t.Run(row.name, func(t *testing.T) {
			milestone := row.milestone
			if milestone == "" {
				milestone = active
			}
			refusesVerification(t, root, milestone, row.edit, row.want...)
		})
	}
	t.Run("undelivered-outcome", func(t *testing.T) {
		partial := commitmenttest.SeedMilestone(t, commitmenttest.MilestoneSpec)
		refusesVerification(t, partial, active, nil, "has undelivered outcomes B")
	})
}

// Each element that a completion proposal must keep refuses on its own: the published
// delivery facts, a recorded completion, the cleared milestone with no activation, the
// verified criteria, and the receipt's milestone.
func TestCommitmentCompletionProposal(t *testing.T) {
	for _, row := range []struct {
		name, want string
		edit       func(*commitment.Policy)
	}{
		{name: "activates-successor", want: "must clear that active milestone and activate none", edit: func(policy *commitment.Policy) {
			policy.ActiveMilestone = commitmenttest.SuccessorMilestone
		}},
		{name: "changes-criteria", want: "changes the criteria", edit: func(policy *commitment.Policy) {
			policy.Milestones[0].Outcomes[1].Criteria[0].Text = "A weaker B obligation."
		}},
		{name: "changes-delivery-fact", want: "cannot change a published delivery fact", edit: func(policy *commitment.Policy) {
			policy.Deliveries[0].Evidence = "other-evidence"
		}},
	} {
		t.Run(row.name, func(t *testing.T) {
			root := deliveredMilestone(t)
			refusesCompletion(t, root, completionProposal(t, root, verified(t, root).ID, row.edit), row.want)
		})
	}
	t.Run("receipt-after-policy-change", func(t *testing.T) {
		root := deliveredMilestone(t)
		verification := verified(t, root)
		commitmenttest.EditPolicy(t, root, func(policy *commitment.Policy) {
			policy.Milestones[1].Outcomes[0].Criteria[0].Text = "The successor obligation changed."
		})
		commitmenttest.Commit(t, root, "change successor")
		refusesCompletion(t, root, completionProposal(t, root, verification.ID, nil), "is stale")
	})
	t.Run("removes-completion", func(t *testing.T) {
		root := deliveredMilestone(t)
		commitmenttest.EditPolicy(t, root, func(policy *commitment.Policy) {
			policy.Completions = []commitment.Completion{{Milestone: active, Verification: "sha256:recorded"}}
			policy.ActiveMilestone = ""
		})
		commitmenttest.Commit(t, root, "record completion")
		refusesCompletion(t, root, completionProposal(t, root, "", func(policy *commitment.Policy) { policy.Completions = nil }), "cannot remove the completion")
	})
	t.Run("receipt-for-another-milestone", func(t *testing.T) {
		root := deliveredMilestone(t)
		other := verified(t, root)
		other.Milestone = commitmenttest.SuccessorMilestone
		other.ID = "sha256:other"
		payload, err := json.Marshal(other)
		if err != nil {
			t.Fatal(err)
		}
		err = intent.Transact(root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
			ledger.MilestoneReceipts = append(ledger.MilestoneReceipts, intent.MilestoneReceipt{ID: other.ID, Payload: string(payload)})
			return ledger, true, nil
		}, nil)
		if err != nil {
			t.Fatal(err)
		}
		refusesCompletion(t, root, completionProposal(t, root, other.ID, nil), "does not verify milestone")
	})
}

// Approval checks the completion receipt again: a receipt that is gone after planning
// refuses approval and stages nothing.
func TestCommitmentCompletionApprovalConsumesReceipt(t *testing.T) {
	root := deliveredMilestone(t)
	proposal := completionProposal(t, root, verified(t, root).ID, nil)
	data, err := os.ReadFile(proposal)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := (commitrepo.Store{Root: root}).Plan(data)
	if err != nil {
		t.Fatal(err)
	}
	err = intent.Transact(root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
		ledger.MilestoneReceipts = nil
		return ledger, true, nil
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	planning := commitmenttest.Planning(t, root)
	before := commitmenttest.MilestoneState(t, root, planning)
	out, code := commitcmd.Command(planning, []string{"approve", "--plan", plan.ID, "--decision", "reviewer-completes-M", "--delayed", "none", "--removed", "none"})
	if code != 1 || !strings.Contains(out, "has no verification receipt") {
		t.Fatalf("approve = (%q, %d), want the missing-receipt refusal", out, code)
	}
	if commitmenttest.MilestoneState(t, root, planning) != before {
		t.Fatal("refused completion approval changed the milestone state")
	}
}

// The policy records a completion only for a known, inactive, fully delivered milestone,
// once, with its verification receipt.
func TestCommitmentCompletionValidation(t *testing.T) {
	delivered, _, err := (commitrepo.Store{Root: deliveredMilestone(t)}).Policy()
	if err != nil {
		t.Fatal(err)
	}
	partial, _, err := (commitrepo.Store{Root: commitmenttest.SeedMilestone(t, commitmenttest.MilestoneSpec)}).Policy()
	if err != nil {
		t.Fatal(err)
	}
	complete := commitment.Completion{Milestone: active, Verification: "sha256:verified"}
	for _, row := range []struct {
		name        string
		policy      commitment.Policy
		active      string
		completions []commitment.Completion
	}{
		{"unknown-milestone", delivered, "", []commitment.Completion{{Milestone: "M9", Verification: complete.Verification}}},
		{"duplicate", delivered, "", []commitment.Completion{complete, complete}},
		{"still-active", delivered, active, []commitment.Completion{complete}},
		{"no-verification", delivered, "", []commitment.Completion{{Milestone: active}}},
		{"undelivered-outcome", partial, "", []commitment.Completion{complete}},
	} {
		t.Run(row.name, func(t *testing.T) {
			policy := row.policy
			policy.ActiveMilestone, policy.Completions = row.active, row.completions
			if err := commitment.Validate(policy); err == nil || !strings.Contains(err.Error(), "invalid completion") {
				t.Fatalf("Validate = %v, want the invalid completion refusal", err)
			}
		})
	}
	delivered.ActiveMilestone, delivered.Completions = "", []commitment.Completion{complete}
	if err := commitment.Validate(delivered); err != nil {
		t.Fatalf("Validate(completed milestone) = %v, want valid", err)
	}
}
