package commitmenttest

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/spec"
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
		policy.Milestones[0].Outcomes[1].Deliverables = []commitment.DeliveryBinding{{Source: commitment.SourceBinding{ID: "tickets", Path: TicketsFolder, Identity: identity}, Obligations: []string{"FT2"}}}
		policy.Milestones = append(policy.Milestones, commitment.Milestone{ID: SuccessorMilestone, Outcomes: []commitment.Outcome{outcome("C")}})
	})
	Commit(t, root, "approve milestone deliveries")
	for _, deliverable := range published {
		Publish(t, root, deliverable)
	}
	return root
}

// TicketEvidence is the ticket that the tickets-only folder holds, cited as native evidence.
const TicketEvidence = TicketsFolder + "/tickets/one.md"

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
// completion must keep: the published policy, the local intent ledger, and the working
// policy of each checkout.
func MilestoneState(t testing.TB, root string, checkouts ...string) string {
	t.Helper()
	state := gittest.Output(t, root, "show", "main:"+commitment.PolicyPath)
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

// Publish commits on main the publication of the verified delivery of deliverable from
// the current head: the exact closure that the commitment owner derives, and the status
// flip of a delivered spec.
func Publish(t testing.TB, root, deliverable string) {
	t.Helper()
	head := gittest.Output(t, root, "rev-parse", "HEAD")
	edits, err := (commitrepo.Store{Root: root}).Closure(head, commitrepo.Delivery{Spec: deliverable, Source: head})
	if err != nil || len(edits) == 0 {
		t.Fatalf("Closure(%s) = %+v, %v; want a delivery closure", deliverable, edits, err)
	}
	for _, edit := range edits {
		if edit.Delete {
			gittest.Output(t, root, "rm", "-q", "--", edit.Path)
			continue
		}
		Write(t, root, edit.Path, string(edit.Data))
	}
	if spec.IsLiveSpecPath(deliverable) {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(deliverable)))
		if err != nil {
			t.Fatal(err)
		}
		flipped, err := spec.Implemented(data)
		if err != nil {
			t.Fatal(err)
		}
		Write(t, root, deliverable, string(flipped))
	}
	Commit(t, root, "publish "+deliverable)
}
