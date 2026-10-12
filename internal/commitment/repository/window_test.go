package repository_test

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/refusalroute"
)

// windowRepo is the Staged repository of outcomes. A parallel grant lets the outcomes start
// together, so each one can hold its own open build.
func windowRepo(t *testing.T, outcomes ...string) string {
	t.Helper()
	root := commitmenttest.Staged(t, outcomes...)
	if len(outcomes) > 1 {
		commitmenttest.EditPolicy(t, root, func(policy *commitment.Policy) {
			policy.ParallelGrants = []commitment.ParallelGrant{{Outcomes: outcomes}}
		})
		commitmenttest.Commit(t, root, "grant the outcomes in parallel")
	}
	return root
}

// approvedSpec is the staged spec that the policy at root approves for outcome.
func approvedSpec(t *testing.T, root, outcome string) string {
	t.Helper()
	policy, _, err := commitrepo.Store{Root: root}.Policy()
	if err != nil {
		t.Fatal(err)
	}
	approved, err := commitment.ActiveOutcome(policy, outcome)
	if err != nil {
		t.Fatal(err)
	}
	return approved.Deliverables[0].Source.Path
}

// startBuild creates the assignment that request names and starts outcome with deliverable.
// It returns the assignment's publication of deliverable to main.
func startBuild(t *testing.T, root, request, outcome, deliverable string) commitrepo.Publication {
	t.Helper()
	path := commitmenttest.Assignment(t, root, request)
	if err := (commitrepo.Store{Root: path}).Start(outcome, request, deliverable); err != nil {
		t.Fatal(err)
	}
	source := windowSource(t, root, path)
	source.Deliverable = deliverable
	return source
}

// windowSource is the spec-less publication to main of the assignment that owns worktree.
func windowSource(t *testing.T, root, worktree string) commitrepo.Publication {
	t.Helper()
	source := publication(t, root, worktree)
	source.Destination = leasedDestination
	return source
}

// editCommitment applies edit to the commitment state of the ledger at root.
func editCommitment(t *testing.T, root string, edit func(*intent.CommitmentState)) {
	t.Helper()
	if err := intent.Transact(root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
		state := intent.CommitmentState{}
		if ledger.Commitment != nil {
			state = *ledger.Commitment
		}
		edit(&state)
		ledger.Commitment = &state
		return ledger, true, nil
	}, nil); err != nil {
		t.Fatal(err)
	}
}

// putGap replaces the chunk gap record of gap's outcome, in the active milestone of root, as
// the chunk checkpoint and the chunk dispatch write it.
func putGap(t *testing.T, root string, gap intent.ChunkGap) {
	t.Helper()
	policy, _, err := commitrepo.Store{Root: root}.Policy()
	if err != nil {
		t.Fatal(err)
	}
	gap.Milestone = policy.ActiveMilestone
	editCommitment(t, root, func(state *intent.CommitmentState) {
		state.Gaps = append(slices.DeleteFunc(slices.Clone(state.Gaps), func(held intent.ChunkGap) bool {
			return held.Outcome == gap.Outcome
		}), gap)
	})
}

// requireWindowRefusal fails unless err is the window refusal of an open build that names
// want: the build, its deliverable, and the gap to wait for.
func requireWindowRefusal(t *testing.T, err error, want string) {
	t.Helper()
	var raised refusalroute.Raised
	if !errors.As(err, &raised) || raised.Name != refusalroute.LandBuildWindow || !strings.Contains(err.Error(), "landing destination "+leasedDestination+" is held by open build "+want) {
		t.Fatalf("window = %v, want the %s refusal naming open build %s", err, refusalroute.LandBuildWindow, want)
	}
}

// TestWindowNamesTheGapToWaitFor is BW48 and BW49. A closed gap record names the
// checkpoint of its own chunk, and a closed record of the last chunk names the completion.
func TestWindowNamesTheGapToWaitFor(t *testing.T) {
	root := windowRepo(t, "A")
	spec := approvedSpec(t, root, "A")
	startBuild(t, root, "build", "A", spec)
	unbound := windowSource(t, root, commitmenttest.Assignment(t, root, "unbound"))
	for _, row := range []struct {
		gap  intent.ChunkGap
		want string
	}{
		{gap: intent.ChunkGap{Outcome: "A", Chunk: "C2"}, want: "the checkpoint of chunk C2"},
		{gap: intent.ChunkGap{Outcome: "A", Chunk: "C3", Last: true}, want: "its completion landing"},
	} {
		putGap(t, root, row.gap)
		requireWindowRefusal(t, commitrepo.Store{Root: root}.AdmitWindow(unbound), "A ("+spec+"); wait for "+row.want)
	}
}

// TestOpenBuildsDeriveFromActiveOwners is BW54 and BW55. One active bound assignment keeps
// the build open. A started tickets-only deliverable opens no build.
func TestOpenBuildsDeriveFromActiveOwners(t *testing.T) {
	t.Run("one-active-owner", func(t *testing.T) {
		root := windowRepo(t, "A")
		spec := approvedSpec(t, root, "A")
		first := startBuild(t, root, "first", "A", spec)
		startBuild(t, root, "second", "A", spec)
		unbound := windowSource(t, root, commitmenttest.Assignment(t, root, "unbound"))
		moveAssignment(t, root, first.Request, intent.StateComplete)
		requireWindowRefusal(t, commitrepo.Store{Root: root}.AdmitWindow(unbound), "A ("+spec+"); wait for its first chunk checkpoint")
	})
	t.Run("tickets-only", func(t *testing.T) {
		root := windowRepo(t, "A")
		identity := commitmenttest.WriteTickets(t, root)
		commitmenttest.EditPolicy(t, root, func(policy *commitment.Policy) {
			delivery := &policy.Milestones[0].Outcomes[0]
			delivery.Deliverables = append(delivery.Deliverables, commitmenttest.TicketsBinding(identity))
		})
		commitmenttest.Commit(t, root, "approve the tickets-only folder")
		startBuild(t, root, "light", "A", commitmenttest.TicketsFolder)
		unbound := windowSource(t, root, commitmenttest.Assignment(t, root, "unbound"))
		if err := (commitrepo.Store{Root: root}).AdmitWindow(unbound); err != nil {
			t.Fatalf("window over a started tickets-only deliverable = %v, want no open build", err)
		}
	})
}

// TestWindowNeedsAGapOfEveryOtherBuild is BW58. A landing inside the gap of build A still
// meets build B. An override for the landing's assignment lets it through both builds.
func TestWindowNeedsAGapOfEveryOtherBuild(t *testing.T) {
	root := windowRepo(t, "A", "B")
	startBuild(t, root, "build-a", "A", approvedSpec(t, root, "A"))
	specB := approvedSpec(t, root, "B")
	startBuild(t, root, "build-b", "B", specB)
	unbound := windowSource(t, root, commitmenttest.Assignment(t, root, "unbound"))
	putGap(t, root, intent.ChunkGap{Outcome: "A", Chunk: "C1", Open: true})
	store := commitrepo.Store{Root: root}
	requireWindowRefusal(t, store.AdmitWindow(unbound), "B ("+specB+"); wait for its first chunk checkpoint")
	editCommitment(t, root, func(state *intent.CommitmentState) {
		state.Overrides = append(state.Overrides, intent.LandingOverride{Assignment: unbound.Assignment, Reason: "back out a landed change"})
	})
	if err := store.AdmitWindow(unbound); err != nil {
		t.Fatalf("window for an overridden assignment = %v, want admission through both builds", err)
	}
}

// TestStartDropsAStaleGap is BW96. A start claims the outcome again, so a gap record that
// survived the earlier claim cannot let a landing through the new build.
func TestStartDropsAStaleGap(t *testing.T) {
	root := windowRepo(t, "A")
	spec := approvedSpec(t, root, "A")
	build := startBuild(t, root, "build", "A", spec)
	unbound := windowSource(t, root, commitmenttest.Assignment(t, root, "unbound"))
	store := commitrepo.Store{Root: root}
	putGap(t, root, intent.ChunkGap{Outcome: "A", Chunk: "C1", Open: true})
	if err := store.Block("A", "waiting for the reviewer"); err != nil {
		t.Fatal(err)
	}
	if err := store.Unblock("A"); err != nil {
		t.Fatal(err)
	}
	// The open record stays from the earlier claim, so only the start can drop it.
	putGap(t, root, intent.ChunkGap{Outcome: "A", Chunk: "C1", Open: true})
	if err := (commitrepo.Store{Root: build.Worktree}).Start("A", "build", spec); err != nil {
		t.Fatal(err)
	}
	requireWindowRefusal(t, store.AdmitWindow(unbound), "A ("+spec+"); wait for its first chunk checkpoint")
}

// TestClosersDropTheGapRecord is BW97. A block and a completion each drop the claim, and
// each drops the outcome's gap record with it.
func TestClosersDropTheGapRecord(t *testing.T) {
	for _, row := range []struct {
		name  string
		close func(t *testing.T, root string)
	}{
		{name: "block", close: func(t *testing.T, root string) {
			if err := (commitrepo.Store{Root: root}).Block(commitmenttest.DeliveryOutcome, "waiting for the reviewer"); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "completion", close: func(t *testing.T, root string) {
			commitmenttest.Publish(t, root, commitmenttest.MilestoneSpec)
			if err := (commitrepo.Store{Root: root}).ReconcileDelivered(); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(row.name, func(t *testing.T) {
			root := commitmenttest.SeedMilestone(t)
			startBuild(t, root, "build", commitmenttest.DeliveryOutcome, commitmenttest.MilestoneSpec)
			putGap(t, root, intent.ChunkGap{Outcome: commitmenttest.DeliveryOutcome, Chunk: "C1", Open: true})
			row.close(t, root)
			ledger, err := intent.Read(root)
			if err != nil {
				t.Fatal(err)
			}
			if ledger.Commitment != nil && len(ledger.Commitment.Gaps) != 0 {
				t.Fatalf("gap records after the %s = %+v, want none", row.name, ledger.Commitment.Gaps)
			}
		})
	}
}
