package repository_test

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

// Each row branches an unbound assignment from main, moves main past the branch point,
// and then commits the row's own change on the branch. A main move is main's, so a branch
// that leaves the moved paths alone is admitted. A branch change still grades as today.
func TestStaleBaseCandidate(t *testing.T) {
	t.Parallel()
	light := steps(ticket(lightTicket, "change.go"), write("change.go"))
	for _, row := range []struct {
		name  string
		move  func(*testing.T, string)
		steps []func(*testing.T, string)
		want  string
	}{
		{name: "main-lands-production", move: landProduction, steps: steps(func(t *testing.T, worktree string) {
			commitmenttest.Write(t, worktree, "capture/notes.md", "A planning note.\n")
		})},
		{name: "main-approves-policy-edit", move: approvePolicyEdit, steps: light},
		// The FT290 case: main closes a delivery, which changes the policy, the
		// recommended sequence, and the row owners.
		{name: "main-closes-delivery", move: closeDelivery, steps: light},
		{name: "main-pins-new-row", move: pinNewRow, steps: light},
		{name: "branch-edits-policy", move: landProduction, steps: append(light, func(t *testing.T, worktree string) {
			commitmenttest.EditPolicy(t, worktree, func(policy *commitment.Policy) {
				policy.Milestones[0].Outcomes[0].Criteria[0].Text = "A criterion the branch proposes."
			})
		}), want: "candidate policy has no exact approval"},
		{name: "branch-changes-sequence", move: landProduction, steps: steps(func(t *testing.T, worktree string) {
			commitmenttest.Write(t, worktree, "ROADMAP.md", strings.Replace(boardIndex(), "1. "+commitmenttest.DeliveryOutcome+"\n2. B\n", "1. B\n2. "+commitmenttest.DeliveryOutcome+"\n", 1))
		}), want: "protected recommended sequence changed"},
		{name: "branch-deletes-pinned-row", move: landProduction, steps: append(light, remove("roadmap/FT2.md")), want: "candidate changes protected commitment"},
		{name: "branch-writes-unticketed-production", move: landProduction, steps: steps(write("change.go")), want: unbound},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			root, _ := lightPathRepo(t)
			worktree := commitmenttest.Assignment(t, root, row.name)
			row.move(t, root)
			for _, step := range row.steps {
				step(t, worktree)
			}
			commitmenttest.Commit(t, worktree, row.name)
			err := (commitrepo.Store{Root: worktree}).AuthorizeCandidate(gittest.Output(t, worktree, "rev-parse", "HEAD^{tree}"))
			if row.want == "" {
				if err != nil {
					t.Fatalf("AuthorizeCandidate = %v, want admission", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), row.want) {
				t.Fatalf("AuthorizeCandidate = %v, want a refusal naming %q", err, row.want)
			}
		})
	}
}

// Each row commits its own change in an unbound assignment and then leaves a merge state
// there. The state is a merge of main, of a sibling commit that main lacks, or of both. One
// row moves main again after its merge, and another row leaves a directory at MERGE_HEAD. A
// row whose pending step returns a wanted path wants a refusal that names that path and not
// the absent path. Every other row wants admission.
func TestPendingMergeCandidate(t *testing.T) {
	t.Parallel()
	const siblingPath = "sibling.go"
	light := steps(ticket(lightTicket, "change.go"), write("change.go"))
	merge := func(t *testing.T, worktree string, revisions ...string) {
		gittest.Output(t, worktree, append([]string{"merge", "--no-commit", "--no-ff"}, revisions...)...)
	}
	// moveMain lands production on main and returns the path of that move.
	moveMain := func(t *testing.T, root string) string {
		landProduction(t, root)
		return gittest.Output(t, root, "diff-tree", "--no-commit-id", "--name-only", "-r", "HEAD")
	}
	mergeMain := func(t *testing.T, root, worktree string) (string, string) {
		landProduction(t, root)
		merge(t, worktree, "main")
		return "", ""
	}
	// sibling commits siblingPath in a sibling assignment from the current main and returns
	// that commit.
	sibling := func(t *testing.T, root string) string {
		worktree := commitmenttest.Assignment(t, root, "sibling")
		write(siblingPath)(t, worktree)
		commitmenttest.Commit(t, worktree, "a sibling lands production")
		return gittest.Output(t, worktree, "rev-parse", "HEAD")
	}
	for _, row := range []struct {
		name    string
		steps   []func(*testing.T, string)
		pending func(t *testing.T, root, worktree string) (want, absent string)
	}{
		{name: "light-path-merges-main", steps: light, pending: mergeMain},
		{name: "planning-merges-main", steps: steps(func(t *testing.T, worktree string) {
			commitmenttest.Write(t, worktree, "capture/notes.md", "A planning note.\n")
		}), pending: mergeMain},
		// The merged main is older than the tip, so the later move of main stays main's.
		{name: "light-path-merges-older-main", steps: light, pending: func(t *testing.T, root, worktree string) (string, string) {
			mergeMain(t, root, worktree)
			write("later.go")(t, root)
			commitmenttest.Commit(t, root, "main lands a later production fix")
			return "", ""
		}},
		{name: "light-path-merges-sibling", steps: light, pending: func(t *testing.T, root, worktree string) (string, string) {
			merge(t, worktree, sibling(t, root))
			return siblingPath, ""
		}},
		// Main is the last merge head, so the reference must read beyond the first head to count
		// main's change as main's.
		{name: "light-path-merges-sibling-and-main", steps: light, pending: func(t *testing.T, root, worktree string) (string, string) {
			head := sibling(t, root)
			moved := moveMain(t, root)
			merge(t, worktree, head, "main")
			return siblingPath, moved
		}},
		{name: "unreadable-merge-state", steps: light, pending: func(t *testing.T, root, worktree string) (string, string) {
			moved := moveMain(t, root)
			path, err := git.AdminPath(worktree, "MERGE_HEAD")
			if err != nil {
				t.Fatal(err)
			}
			if err := os.Mkdir(path, 0o755); err != nil {
				t.Fatal(err)
			}
			return moved, ""
		}},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			root, _ := lightPathRepo(t)
			worktree := commitmenttest.Assignment(t, root, row.name)
			for _, step := range row.steps {
				step(t, worktree)
			}
			commitmenttest.Commit(t, worktree, row.name)
			want, absent := row.pending(t, root, worktree)
			err := (commitrepo.Store{Root: worktree}).AuthorizeCandidate(gittest.Output(t, worktree, "write-tree"))
			if want == "" {
				if err != nil {
					t.Fatalf("AuthorizeCandidate = %v, want admission", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), strconv.Quote(want)) || (absent != "" && strings.Contains(err.Error(), strconv.Quote(absent))) {
				t.Fatalf("AuthorizeCandidate = %v, want a refusal naming %q and not %q", err, want, absent)
			}
		})
	}
}

// boardIndex is the board index that lightPathRepo commits on main.
func boardIndex() string {
	return strings.Replace(commitmenttest.ClosureIndex(), "## Recommended", unpinnedRow+"## Recommended", 1)
}

// landProduction commits a production file on main.
func landProduction(t *testing.T, root string) {
	commitmenttest.Write(t, root, "fix.go", "package fixture\n")
	commitmenttest.Commit(t, root, "main lands a production fix")
}

// approvePolicyEdit commits an edit of one criterion of the policy on main.
func approvePolicyEdit(t *testing.T, root string) {
	commitmenttest.EditPolicy(t, root, func(policy *commitment.Policy) {
		policy.Milestones[0].Outcomes[0].Criteria[0].Text = "A criterion that main approved later."
	})
	commitmenttest.Commit(t, root, "main approves a policy edit")
}

// closeDelivery commits on main the closure of the delivery outcome's FT1: the delivery
// fact, the implemented spec, the board without FT1 and its sequence entry, and no FT1
// detail owner.
func closeDelivery(t *testing.T, root string) {
	source := gittest.Output(t, root, "rev-parse", "HEAD")
	evidence := gittest.Output(t, root, "rev-parse", source+":"+commitmenttest.MilestoneRecord)
	commitmenttest.EditPolicy(t, root, func(policy *commitment.Policy) {
		binding := policy.Milestones[0].Outcomes[0].Deliverables[0].Source
		policy.Deliveries = append(policy.Deliveries, commitment.DeliveryFact{Milestone: commitmenttest.ClosureMilestone, Outcome: commitmenttest.DeliveryOutcome, Binding: binding.ID, Identity: binding.Identity, Source: source, Evidence: evidence})
	})
	commitmenttest.Write(t, root, commitmenttest.MilestoneSpec, "# x\n\nStatus: implemented\n")
	commitmenttest.Write(t, root, "ROADMAP.md", "# Roadmap\n\n## Parked\n\n**FT2 — B**\n\n"+unpinnedRow+"## Recommended sequence\n\n1. B\n")
	if err := os.Remove(filepath.Join(root, "roadmap", "FT1.md")); err != nil {
		t.Fatal(err)
	}
	commitmenttest.Commit(t, root, "main closes the FT1 delivery")
}

// pinNewRow commits on main a new row FT3 that the delivery outcome pins, with its detail
// owner and its board row.
func pinNewRow(t *testing.T, root string) {
	const body = "**FT3 — " + commitmenttest.DeliveryOutcome + "**\n\nKeep the later obligation.\n"
	commitmenttest.Write(t, root, "roadmap/FT3.md", body)
	commitmenttest.EditPolicy(t, root, func(policy *commitment.Policy) {
		delivery := &policy.Milestones[0].Outcomes[0]
		delivery.Sources = append(delivery.Sources, commitment.SourceBinding{ID: "FT3", Path: "roadmap/FT3.md", Identity: commitment.Identity([]byte(body))})
	})
	commitmenttest.Write(t, root, "ROADMAP.md", strings.Replace(boardIndex(), "**FT2 — B**", "**FT3 — "+commitmenttest.DeliveryOutcome+"**\n\n**FT2 — B**", 1))
	commitmenttest.Commit(t, root, "main pins FT3")
}
