package commitment_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/gittest"
)

// headSpec is the spec path that a worktree HEAD holds and main lacks. headRequest is the
// request of that worktree's assignment, and headOutcome is the outcome that binds the spec.
const (
	headSpec    = "specs/head/spec.md"
	headRequest = "spec-stage"
	headOutcome = "A"
)

// headWorktree commits a main whose active milestone M1 holds headOutcome and then the
// other outcomes, and returns the root and one assignment worktree of it.
func headWorktree(t *testing.T, others ...string) (root, worktree string) {
	t.Helper()
	root = commitmenttest.Repo(t, policy([]commitment.Milestone{milestone("M1", append([]string{headOutcome}, others...)...)}, "M1"))
	return root, commitmenttest.Assignment(t, root, headRequest)
}

// headCommit commits body at headSpec in worktree.
func headCommit(t *testing.T, worktree, body string) {
	t.Helper()
	commitmenttest.Write(t, worktree, headSpec, body)
	commitmenttest.Commit(t, worktree, "stage the spec")
}

// headBinding binds the spec bytes body at headSpec with no obligation.
func headBinding(body string) commitment.DeliveryBinding {
	return commitment.DeliveryBinding{Source: commitment.SourceBinding{ID: "A.spec", Path: headSpec, Identity: commitment.Identity([]byte(body))}}
}

// Plan accepts the binding of a staged spec that the worktree HEAD holds and main lacks.
// Approve then writes the bound policy into the worktree's policy file.
func TestHeadBindingPlansAStagedSpecThatHEADHolds(t *testing.T) {
	root, worktree := headWorktree(t)
	headCommit(t, worktree, commitmenttest.StagedBody)
	proposal := commitmenttest.Bound(t, root, headBinding(commitmenttest.StagedBody))
	store := commitrepo.Store{Root: worktree}
	plan, err := store.Plan(proposal)
	if err != nil {
		t.Fatalf("Plan(a binding of the spec that HEAD holds) = %v, want acceptance", err)
	}
	if _, err := store.Approve(plan.ID, "reviewer sign-off", plan.Effects.Delayed, plan.Effects.Removed); err != nil {
		t.Fatalf("Approve(the HEAD binding) = %v, want the bound policy staged", err)
	}
	staged, err := os.ReadFile(filepath.Join(worktree, filepath.FromSlash(commitment.PolicyPath)))
	if err != nil || !bytes.Equal(staged, proposal) {
		t.Fatalf("worktree policy = (%q, %v), want the bound policy %q", staged, err, proposal)
	}
}

// The worktree inventory lists the staged spec that HEAD holds. The row carries the spec's
// Roadmap value as its scope and the identity that plan accepts for a binding of the spec.
func TestHeadBindingInventoryListsAHeadHeldSpec(t *testing.T) {
	const row = "FT7"
	body := commitmenttest.StagedBody + "Roadmap: " + row + "\n"
	root, worktree := headWorktree(t)
	headCommit(t, worktree, body)
	store := commitrepo.Store{Root: worktree}
	items, err := store.Inventory()
	if err != nil {
		t.Fatal(err)
	}
	var held []commitrepo.InventoryItem
	for _, item := range items {
		if item.Kind == "deliverable" && item.ID == headSpec {
			held = append(held, item)
		}
	}
	if len(held) != 1 {
		t.Fatalf("worktree inventory = %+v, want one deliverable row for %s", items, headSpec)
	}
	if held[0].Scope != row {
		t.Errorf("scope of %s = %q, want its Roadmap value %q", headSpec, held[0].Scope, row)
	}
	if want := commitment.Identity([]byte(body)); held[0].Identity != want {
		t.Errorf("identity of %s = %q, want the identity %q of its HEAD bytes", headSpec, held[0].Identity, want)
	}
	binding := commitment.DeliveryBinding{Source: commitment.SourceBinding{ID: "A.spec", Path: held[0].ID, Identity: held[0].Identity}}
	if _, err := store.Plan(commitmenttest.Bound(t, root, binding)); err != nil {
		t.Errorf("Plan(the binding of the inventory row) = %v, want acceptance", err)
	}
}

// Plan accepts the binding of a tickets-only folder that the worktree HEAD holds and main
// lacks.
func TestHeadBindingPlansATicketsOnlyFolder(t *testing.T) {
	root, worktree := headWorktree(t)
	identity := commitmenttest.WriteTickets(t, worktree)
	commitmenttest.Commit(t, worktree, "stage the tickets-only folder")
	if _, err := (commitrepo.Store{Root: worktree}).Plan(commitmenttest.Bound(t, root, commitmenttest.TicketsBinding(identity))); err != nil {
		t.Fatalf("Plan(a binding of the tickets-only folder that HEAD holds) = %v, want acceptance", err)
	}
}

// Plan refuses a new binding unless HEAD holds the staged spec, and the refusal writes
// nothing. Each row holds the spec elsewhere: only in the working tree, only on main, or at
// HEAD with no Status: staged line.
func TestHeadBindingRefusesAnUnheldDeliverable(t *testing.T) {
	const unstaged = "# Unstaged delivery\n"
	refused := "commitment source " + strconv.Quote(headBinding("").Source.ID) + " refused: "
	for _, row := range []struct {
		name, body, want string
		hold             func(t *testing.T, root, worktree string)
	}{
		{name: "working-tree-only", body: commitmenttest.StagedBody, want: refused, hold: func(t *testing.T, _, worktree string) {
			commitmenttest.Write(t, worktree, headSpec, commitmenttest.StagedBody)
		}},
		{name: "main-only", body: commitmenttest.StagedBody, want: refused, hold: func(t *testing.T, root, _ string) {
			commitmenttest.Write(t, root, headSpec, commitmenttest.StagedBody)
			commitmenttest.Commit(t, root, "main stages the spec")
		}},
		{name: "not-staged", body: unstaged, want: "deliverable " + strconv.Quote(headSpec) + " is not staged", hold: func(t *testing.T, _, worktree string) {
			headCommit(t, worktree, unstaged)
		}},
	} {
		t.Run(row.name, func(t *testing.T) {
			root, worktree := headWorktree(t)
			row.hold(t, root, worktree)
			before := commitmenttest.MilestoneState(t, root, worktree)
			_, err := (commitrepo.Store{Root: worktree}).Plan(commitmenttest.Bound(t, root, headBinding(row.body)))
			if err == nil || !strings.Contains(err.Error(), row.want) || commitmenttest.MilestoneState(t, root, worktree) != before {
				t.Fatalf("Plan = %v, want a refusal naming %q that writes nothing", err, row.want)
			}
		})
	}
}

// Every source other than a new deliverable reads the default-branch tip. A worktree edit
// of a row therefore cannot pass as current, and a binding that main added after the
// worktree's HEAD stays valid there.
func TestHeadBindingKeepsTheDefaultTipForEveryOtherSource(t *testing.T) {
	t.Run("changed-row", func(t *testing.T) {
		const edited = "**FT1 — A**\n\nAn obligation that the worktree changes.\n"
		root := commitmenttest.Repo(t, policy([]commitment.Milestone{milestone("M1", headOutcome)}, "M1"))
		commitmenttest.SeedProtected(t, root)
		commitmenttest.Commit(t, root, "protect the rows")
		current, _, err := commitrepo.Store{Root: root}.Policy()
		if err != nil {
			t.Fatal(err)
		}
		published := current.Milestones[0].Outcomes[0].Sources[0]
		worktree := commitmenttest.Assignment(t, root, headRequest)
		commitmenttest.Write(t, worktree, published.Path, edited)
		commitmenttest.Commit(t, worktree, "change the row")
		changed := commitment.Identity([]byte(edited))
		current.Milestones[0].Outcomes[0].Sources[0].Identity = changed
		_, err = commitrepo.Store{Root: worktree}.Plan(mustBytes(t, current))
		want := "commitment source " + strconv.Quote(published.ID) + " changed from " + changed + " to " + published.Identity
		if err == nil || err.Error() != want {
			t.Fatalf("Plan(the row identity that HEAD holds) = %v, want %q", err, want)
		}
	})
	t.Run("retained-binding", func(t *testing.T) {
		const mainSpec = "specs/main/spec.md"
		root, worktree := headWorktree(t, "B")
		headCommit(t, worktree, commitmenttest.StagedBody)
		commitmenttest.Write(t, root, mainSpec, commitmenttest.StagedBody)
		commitmenttest.EditPolicy(t, root, func(policy *commitment.Policy) {
			policy.Milestones[0].Outcomes[1].Deliverables = []commitment.DeliveryBinding{{Source: commitment.SourceBinding{ID: "B.spec", Path: mainSpec, Identity: commitment.Identity([]byte(commitmenttest.StagedBody))}}}
		})
		commitmenttest.Commit(t, root, "main binds the B spec")
		if _, err := (commitrepo.Store{Root: worktree}).Plan(commitmenttest.Bound(t, root, headBinding(commitmenttest.StagedBody))); err != nil {
			t.Fatalf("Plan(a policy that keeps the binding that main added after HEAD) = %v, want acceptance", err)
		}
	})
}

// Plan reads each new deliverable at the checkout HEAD. Each row binds the spec that HEAD
// holds while main holds no policy, or while main binds that path at another identity.
func TestHeadBindingReadsEachNewDeliverableAtHEAD(t *testing.T) {
	t.Run("first-policy", func(t *testing.T) {
		root := gittest.RepoOnBranch(t, "main")
		commitmenttest.Write(t, root, "README.md", "# Fixture\n")
		commitmenttest.Commit(t, root, "main holds no policy")
		worktree := commitmenttest.Assignment(t, root, headRequest)
		headCommit(t, worktree, commitmenttest.StagedBody)
		first := policy([]commitment.Milestone{milestone("M1", headOutcome)}, "M1")
		first.Milestones[0].Outcomes[0].Deliverables = []commitment.DeliveryBinding{headBinding(commitmenttest.StagedBody)}
		if _, err := (commitrepo.Store{Root: worktree}).Plan(mustBytes(t, first)); err != nil {
			t.Fatalf("Plan(the first policy, which binds the spec that HEAD holds) = %v, want acceptance", err)
		}
	})
	t.Run("rebinding", func(t *testing.T) {
		const restaged = commitmenttest.StagedBody + "A body that main does not bind.\n"
		root, worktree := headWorktree(t)
		commitmenttest.Write(t, root, headSpec, commitmenttest.StagedBody)
		commitmenttest.EditPolicy(t, root, func(policy *commitment.Policy) {
			policy.Milestones[0].Outcomes[0].Deliverables = []commitment.DeliveryBinding{headBinding(commitmenttest.StagedBody)}
		})
		commitmenttest.Commit(t, root, "main binds the spec")
		headCommit(t, worktree, restaged)
		rebound, _, err := commitrepo.Store{Root: root}.Policy()
		if err != nil {
			t.Fatal(err)
		}
		rebound.Milestones[0].Outcomes[0].Deliverables[0] = headBinding(restaged)
		if _, err := (commitrepo.Store{Root: worktree}).Plan(mustBytes(t, rebound)); err != nil {
			t.Fatalf("Plan(a binding of the path that main binds, at the identity that HEAD holds) = %v, want acceptance", err)
		}
	})
}

// Before the landing, start reads only the default-branch policy. It refuses the spec that
// the worktree's committed policy binds, and it records nothing.
func TestHeadBindingStartRefusesBeforeTheLanding(t *testing.T) {
	root, worktree := headWorktree(t)
	headCommit(t, worktree, commitmenttest.StagedBody)
	store := commitrepo.Store{Root: worktree}
	plan, err := store.Plan(commitmenttest.Bound(t, root, headBinding(commitmenttest.StagedBody)))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Approve(plan.ID, "reviewer sign-off", plan.Effects.Delayed, plan.Effects.Removed); err != nil {
		t.Fatal(err)
	}
	commitmenttest.Commit(t, worktree, "commit the approved binding")
	before := commitmenttest.MilestoneState(t, root, worktree)
	err = store.Start(headOutcome, headRequest, headSpec)
	want := "deliverable " + strconv.Quote(headSpec) + " is not approved for outcome " + strconv.Quote(headOutcome)
	if err == nil || !strings.Contains(err.Error(), want) || commitmenttest.MilestoneState(t, root, worktree) != before {
		t.Fatalf("Start(the HEAD-bound spec) = %v, want a refusal naming %q that records nothing", err, want)
	}
}
