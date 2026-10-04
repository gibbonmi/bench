package commitmenttest

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/gate/greenmarker"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/landing/published"
)

// MilestoneSpec is the spec that SeedMilestone approves for the delivery outcome.
const MilestoneSpec = "specs/x/spec.md"

// MilestoneRecord is the review record that MilestoneSpec retains as completion evidence.
const MilestoneRecord = "reviews/x.md"

// SuccessorMilestone is the planned milestone that follows the SeedMilestone milestone.
const SuccessorMilestone = "M2"

// SeedMilestone commits the SeedClosure milestone M with a planned successor. The delivery
// outcome owns FT1 through MilestoneSpec, and outcome B owns FT2 through the tickets-only
// folder. It then publishes the delivery of each published deliverable in order and
// returns the repository. With both deliverables published, every outcome of M is
// delivered and no row of M remains on the board.
func SeedMilestone(t testing.TB, published ...string) string {
	t.Helper()
	root := gittest.RepoOnBranch(t, "main")
	Write(t, root, MilestoneSpec, "# x\n\nStatus: staged\n")
	Write(t, root, MilestoneRecord, "record\n")
	SeedClosure(t, root, MilestoneSpec)
	identity := WriteTickets(t, root)
	EditPolicy(t, root, func(policy *commitment.Policy) {
		policy.Milestones[0].Outcomes[1].Deliverables = []commitment.DeliveryBinding{TicketsBinding(identity, "FT2")}
		policy.Milestones = append(policy.Milestones, commitment.Milestone{ID: SuccessorMilestone, Outcomes: []commitment.Outcome{outcome("C")}})
	})
	Commit(t, root, "approve milestone deliveries")
	for _, deliverable := range published {
		Publish(t, root, deliverable)
	}
	return root
}

// TicketEvidence is the native evidence of a tickets-only delivery on the published tree:
// the policy whose delivery fact records it. The publication removes the folder itself.
const TicketEvidence = commitment.PolicyPath

// Evidence writes a complete evidence document for the active milestone of the published
// policy in root and returns its path. Each result is verified at the current main
// revision for its current criterion identity. A delivery outcome criterion cites
// MilestoneRecord, and every other criterion cites TicketEvidence. The gate evidence is
// MilestoneRecord. A non-nil edit changes the document before it is written.
func Evidence(t testing.TB, root string, edit func(*commitment.MilestoneEvidence)) string {
	t.Helper()
	policy, exists, err := (commitrepo.Store{Root: root}).Policy()
	if err != nil || !exists {
		t.Fatalf("published policy = %v, %v; want a policy", exists, err)
	}
	evidence := commitment.MilestoneEvidence{Version: 1, Revision: gittest.Output(t, root, "rev-parse", "main"), Gate: MilestoneRecord}
	for _, milestone := range policy.Milestones {
		if milestone.ID != policy.ActiveMilestone {
			continue
		}
		for _, outcome := range milestone.Outcomes {
			reference := TicketEvidence
			if outcome.ID == DeliveryOutcome {
				reference = MilestoneRecord
			}
			for _, criterion := range outcome.Criteria {
				evidence.Results = append(evidence.Results, commitment.CriterionResult{Criterion: criterion.ID, Identity: commitment.CriterionIdentity(criterion), Result: commitment.ResultVerified, Evidence: reference, Assessment: "The reviewer confirmed the " + outcome.ID + " outcome."})
			}
		}
	}
	if edit != nil {
		edit(&evidence)
	}
	data, err := json.Marshal(evidence)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "evidence.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// MilestoneState returns the milestone state in root that a refused verification or
// completion must keep: the published policy, or its absence, the local intent ledger, and
// the working policy of each checkout.
func MilestoneState(t testing.TB, root string, checkouts ...string) string {
	t.Helper()
	state := gittest.Output(t, root, "ls-tree", "main", "--", commitment.PolicyPath)
	if state != "" {
		state = gittest.Output(t, root, "show", "main:"+commitment.PolicyPath)
	}
	address, err := intent.Address(root)
	if err != nil {
		t.Fatal(err)
	}
	paths := []string{address}
	for _, checkout := range checkouts {
		paths = append(paths, filepath.Join(checkout, filepath.FromSlash(commitment.PolicyPath)))
	}
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			t.Fatal(err)
		}
		state += "\x00" + string(data)
	}
	return state
}

// Publish commits on main the tree that the verified landing of deliverable publishes
// from the current head, checks that tree out, and advances the project-green marker to
// the commit as the landing does.
func Publish(t testing.TB, root, deliverable string) {
	t.Helper()
	head := gittest.Output(t, root, "rev-parse", "HEAD")
	base := gittest.Output(t, root, "rev-parse", "HEAD^{tree}")
	tree, err := published.Tree(root, base, deliverable, head)
	if err != nil || tree == base {
		t.Fatalf("published.Tree(%s) = %s, %v; want a delivery publication", deliverable, tree, err)
	}
	gittest.Output(t, root, "reset", "-q", "--hard", gittest.Output(t, root, "commit-tree", tree, "-p", head, "-m", "publish "+deliverable))
	MarkGreen(t, root)
}

// MarkGreen advances the project-green marker of main from its prior commit to the tip of
// main, as a landing does once its publication commits.
func MarkGreen(t testing.TB, root string) {
	t.Helper()
	prior, _, err := greenmarker.Read(root, "main")
	if err != nil {
		t.Fatal(err)
	}
	if err := greenmarker.Advance(root, "main", gittest.Output(t, root, "rev-parse", "main"), prior); err != nil {
		t.Fatal(err)
	}
}
