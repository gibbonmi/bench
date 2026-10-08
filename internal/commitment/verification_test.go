package commitment_test

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitcmd"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/refusalroute/routetest"
	"github.com/gibbonmi/bench/internal/sanitize"
)

const active = commitmenttest.ClosureMilestone

// deliveredMilestone returns a repository whose active milestone has every outcome
// delivered through the real closure: one spec and one tickets-only delivery.
func deliveredMilestone(t *testing.T) string {
	t.Helper()
	return commitmenttest.SeedMilestone(t, commitmenttest.MilestoneSpec, commitmenttest.TicketsFolder)
}

func verifyCommand(root, milestone, evidence string) (string, int) {
	return commitcmd.Command(root, []string{"verify", "--milestone", milestone, "--evidence", evidence})
}

// refusesVerification proves that verification of milestone refuses the evidence that
// edit produces, names each want fragment, and keeps the milestone state. RR47: the next
// cell keeps the verification of the same milestone, with the evidence the operator corrects.
func refusesVerification(t *testing.T, root, milestone string, edit func(*commitment.MilestoneEvidence), want ...string) {
	t.Helper()
	evidence := commitmenttest.Evidence(t, root, edit)
	before := commitmenttest.MilestoneState(t, root)
	out, code := verifyCommand(root, milestone, evidence)
	for _, fragment := range want {
		if code != 1 || !strings.Contains(out, fragment) {
			t.Fatalf("verify = (%q, %d), want a refusal naming %q", out, code, fragment)
		}
	}
	retry := "bench commitment verify --milestone " + sanitize.ShellQuote(milestone) + " --evidence <file>"
	if next, printed := routetest.NextCell(out); !printed || next != retry {
		t.Fatalf("verify = %q, want the next cell %q", out, retry)
	}
	if commitmenttest.MilestoneState(t, root) != before {
		t.Fatal("refused verification changed the milestone state")
	}
}

// verified records the receipt of the complete evidence for the active milestone of root
// through the command and returns it.
func verified(t *testing.T, root string) commitment.Verification {
	t.Helper()
	evidence := commitmenttest.Evidence(t, root, nil)
	if out, code := verifyCommand(root, active, evidence); code != 0 || !strings.Contains(out, "commitment_verification[1]") || !strings.Contains(out, "bench commitment plan --input <file>") {
		t.Fatalf("verify = (%q, %d), want a receipt that names the completion proposal", out, code)
	}
	data, err := os.ReadFile(evidence)
	if err != nil {
		t.Fatal(err)
	}
	// The identical replay records nothing new and returns the recorded receipt.
	verification, err := (commitrepo.Store{Root: root}).Verify(active, data)
	if err != nil {
		t.Fatal(err)
	}
	// The receipt stores the examined revision, the project-green marker there, and the
	// object of the gate evidence and of each criterion's evidence at that revision.
	document, err := commitment.ParseEvidence(data)
	if err != nil {
		t.Fatal(err)
	}
	revision := gittest.Output(t, root, "rev-parse", "main")
	object := func(reference string) string { return gittest.Output(t, root, "rev-parse", "main:"+reference) }
	if verification.Revision != revision || verification.Green != revision || verification.Gate != object(document.Gate) || len(verification.Results) != len(document.Results) {
		t.Fatalf("receipt = %+v, want revision, marker, and gate object of %s at %s", verification, document.Gate, revision)
	}
	for i, result := range verification.Results {
		if result.Result != document.Results[i] || result.Object != object(document.Results[i].Evidence) {
			t.Fatalf("receipt result %d = %+v, want %+v at the object of %s", i, result, document.Results[i], document.Results[i].Evidence)
		}
	}
	return verification
}

// completionProposal writes the published policy of root with its active milestone
// completed by verification and cleared, changes it with a non-nil edit, and returns its
// path.
func completionProposal(t *testing.T, root, verification string, edit func(*commitment.Policy)) string {
	t.Helper()
	policy, _, err := (commitrepo.Store{Root: root}).Policy()
	if err != nil {
		t.Fatal(err)
	}
	policy.Completions = append(slices.Clone(policy.Completions), commitment.Completion{Milestone: policy.ActiveMilestone, Verification: verification})
	policy.ActiveMilestone = ""
	if edit != nil {
		edit(&policy)
	}
	data, err := commitment.Bytes(policy)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "completion.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// refusesCompletion proves that the completion proposal at path refuses to plan with want
// and keeps the milestone state, the planning policy included.
func refusesCompletion(t *testing.T, root, path, want string) {
	t.Helper()
	planning := commitmenttest.Planning(t, root)
	before := commitmenttest.MilestoneState(t, root, planning)
	out, code := commitcmd.Command(planning, []string{"plan", "--input", path})
	if code != 1 || !strings.Contains(out, want) {
		t.Fatalf("completion plan = (%q, %d), want a refusal naming %q", out, code, want)
	}
	if commitmenttest.MilestoneState(t, root, planning) != before {
		t.Fatal("refused completion changed the milestone state")
	}
}

// A milestone whose roadmap rows are all closed is not complete. Without criterion
// results verification refuses, and without a verification receipt completion refuses.
func TestCommitmentEmptyRowsNotComplete(t *testing.T) {
	root := deliveredMilestone(t)
	board := gittest.Output(t, root, "show", "main:ROADMAP.md")
	if strings.Contains(board, "FT1") || strings.Contains(board, "FT2") {
		t.Fatalf("board = %q, want every row of the milestone closed", board)
	}
	if listed := gittest.Output(t, root, "ls-tree", "-r", "--name-only", "main", "--", commitmenttest.TicketsFolder); listed != "" {
		t.Fatalf("published tickets-only folder = %q, want the folder closed by its publication", listed)
	}
	t.Run("no-criterion-results", func(t *testing.T) {
		refusesVerification(t, root, active, func(evidence *commitment.MilestoneEvidence) { evidence.Results = nil }, "has no result")
	})
	t.Run("no-verification-receipt", func(t *testing.T) {
		refusesCompletion(t, root, completionProposal(t, root, "sha256:unverified", nil), "has no verification receipt")
	})
}

// One unmet criterion refuses verification although the green gate evidence resolves.
func TestCommitmentUnmetCriterion(t *testing.T) {
	root := deliveredMilestone(t)
	refusesVerification(t, root, active, func(evidence *commitment.MilestoneEvidence) {
		evidence.Results[1].Result = commitment.ResultUnmet
	}, "B-done", "is unmet")
}

// Complete current evidence and an explicit approval complete the milestone. The
// completion clears the active milestone and leaves the planned successor inactive, and
// the staged policy commits under the exact approval.
func TestCommitmentMilestoneCompletion(t *testing.T) {
	root := deliveredMilestone(t)
	verification := verified(t, root)
	planning := commitmenttest.Planning(t, root)
	proposal := completionProposal(t, root, verification.ID, nil)
	out, code := commitcmd.Command(planning, []string{"plan", "--input", proposal})
	if code != 0 || !strings.Contains(out, "completed,"+active) || strings.Contains(out, "activated,") {
		t.Fatalf("completion plan = (%q, %d), want the completed effect and no activation", out, code)
	}
	data, err := os.ReadFile(proposal)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := (commitrepo.Store{Root: root}).Plan(data)
	if err != nil {
		t.Fatal(err)
	}
	published := gittest.Output(t, root, "show", "main:"+commitment.PolicyPath)
	out, code = commitcmd.Command(planning, []string{"approve", "--plan", plan.ID, "--decision", "reviewer-completes-M", "--delayed", "none", "--removed", "none"})
	if code != 0 {
		t.Fatalf("approve = (%q, %d), want the completion staged", out, code)
	}
	staged, err := os.ReadFile(filepath.Join(planning, filepath.FromSlash(commitment.PolicyPath)))
	if err != nil {
		t.Fatal(err)
	}
	policy, err := commitment.Parse(staged)
	if err != nil {
		t.Fatal(err)
	}
	want := []commitment.Completion{{Milestone: active, Verification: verification.ID}}
	if policy.ActiveMilestone != "" || !slices.Equal(policy.Completions, want) {
		t.Fatalf("staged policy active=%q completions=%+v, want no active milestone and %+v", policy.ActiveMilestone, policy.Completions, want)
	}
	if gittest.Output(t, root, "show", "main:"+commitment.PolicyPath) != published {
		t.Fatal("approval changed the published policy before publication")
	}
	commitmenttest.Commit(t, planning, "complete milestone")
	if err := (commitrepo.Store{Root: planning}).AuthorizeCandidate(gittest.Output(t, planning, "rev-parse", "HEAD^{tree}")); err != nil {
		t.Fatalf("AuthorizeCandidate = %v, want the approved completion admitted", err)
	}
}

// A changed criterion or revision invalidates evidence and its receipt.
func TestCommitmentStaleEvidence(t *testing.T) {
	rewordB := func(t *testing.T, root string) {
		commitmenttest.EditPolicy(t, root, func(policy *commitment.Policy) {
			policy.Milestones[0].Outcomes[1].Criteria[0].Text = "The B obligation is satisfied for every reader."
		})
		commitmenttest.Commit(t, root, "reword criterion")
		commitmenttest.MarkGreen(t, root)
	}
	t.Run("criterion-identity", func(t *testing.T) {
		root := deliveredMilestone(t)
		policy, _, err := (commitrepo.Store{Root: root}).Policy()
		if err != nil {
			t.Fatal(err)
		}
		examined := commitment.CriterionIdentity(policy.Milestones[0].Outcomes[1].Criteria[0])
		rewordB(t, root)
		refusesVerification(t, root, active, func(evidence *commitment.MilestoneEvidence) { evidence.Results[1].Identity = examined }, "B-done", "result is stale")
	})
	t.Run("revision", func(t *testing.T) {
		root := deliveredMilestone(t)
		earlier := gittest.Output(t, root, "rev-parse", "main~1")
		refusesVerification(t, root, active, func(evidence *commitment.MilestoneEvidence) { evidence.Revision = earlier }, "is not the published revision")
	})
	t.Run("receipt-after-criterion-change", func(t *testing.T) {
		root := deliveredMilestone(t)
		verification := verified(t, root)
		rewordB(t, root)
		refusesCompletion(t, root, completionProposal(t, root, verification.ID, nil), "is stale")
	})
}

// A missing, duplicate, blocked, or unknown criterion result each refuses
// verification and keeps the milestone state.
func TestCommitmentCriterionCoverage(t *testing.T) {
	root := deliveredMilestone(t)
	for _, row := range []struct {
		name string
		want []string
		edit func(*commitment.MilestoneEvidence)
	}{
		{"missing", []string{"B-done", "has no result"}, func(evidence *commitment.MilestoneEvidence) { evidence.Results = evidence.Results[:1] }},
		{"duplicate", []string{"delivery-done", "has a duplicate result"}, func(evidence *commitment.MilestoneEvidence) {
			evidence.Results = append(evidence.Results, evidence.Results[0])
		}},
		{"blocked", []string{"B-done", "is blocked"}, func(evidence *commitment.MilestoneEvidence) { evidence.Results[1].Result = commitment.ResultBlocked }},
		{"unknown-criterion", []string{"Z-done", "is unknown to milestone"}, func(evidence *commitment.MilestoneEvidence) {
			evidence.Results = append(evidence.Results, commitment.CriterionResult{Criterion: "Z-done", Identity: "sha256:z", Result: commitment.ResultVerified, Evidence: commitmenttest.MilestoneRecord, Assessment: "Confirmed."})
		}},
		{"unknown-result", []string{"B-done", "has unknown result", "passed"}, func(evidence *commitment.MilestoneEvidence) { evidence.Results[1].Result = "passed" }},
	} {
		t.Run(row.name, func(t *testing.T) { refusesVerification(t, root, active, row.edit, row.want...) })
	}
}
