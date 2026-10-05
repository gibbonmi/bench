package worktree

import (
	"errors"
	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

// The lookup outcomes here are authored apart from TreeTarget: the spec fixes the exact-label
// rule, the precedence of one active match, and the path-shape outcome.

// retireAssignment records assignment as complete, so it is no longer active.
func retireAssignment(t *testing.T, root string, assignment intent.Assignment) {
	t.Helper()
	assignment.State = intent.StateComplete
	mustNoError(t, intent.PutAssignment(root, assignment))
}

// TestTreeTargetPrefersOneActiveLabel: one active match wins over an inactive match of the
// same label. With no active match, two inactive matches are ambiguous.
func TestTreeTargetPrefersOneActiveLabel(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	retired := mustCreate(t, root, home, "tree-target-retired", "alpha")
	retireAssignment(t, root, retired.Assignment)
	active := mustCreate(t, root, home, "tree-target-active", "alpha")
	path, err := TreeTarget(root, "alpha")
	requireTest(t, err == nil && path == active.Assignment.Worktree, "TreeTarget(alpha) = (%q, %v), want the active worktree %q", path, err, active.Assignment.Worktree)
	retireAssignment(t, root, active.Assignment)
	var ambiguous ambiguousTargetError
	_, err = TreeTarget(root, "alpha")
	requireTest(t, errors.As(err, &ambiguous) && len(ambiguous.IDs) == 2, "TreeTarget(alpha) with two inactive matches = %v, want the ambiguity of both ids", err)
}

// TestTreeTargetTakesOnlyAnExactLabel: an id or a prefix of an active assignment names no
// tree target. A value that names no label is a path shape when the worktree target
// grammar reads it as a path, and it is unassigned otherwise.
func TestTreeTargetTakesOnlyAnExactLabel(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	alpha := mustCreate(t, root, filepath.Join(root, ".bench-home"), "tree-target-exact", "alphabetical")
	for _, row := range []struct {
		value string
		want  error
	}{
		{alpha.Assignment.ID, errTargetUnassigned},
		{alpha.Assignment.ID[:minOperandPrefix], errTargetUnassigned},
		{"alphabet", errTargetUnassigned},
		{alpha.Assignment.Worktree, ErrTreeTargetPath},
		{"./alphabetical", ErrTreeTargetPath},
	} {
		if _, err := TreeTarget(root, row.value); !errors.Is(err, row.want) {
			t.Errorf("TreeTarget(%q) = %v, want %v", row.value, err, row.want)
		}
	}
}

func TestCommitmentSiblingBinding(t *testing.T) {
	t.Parallel()
	f := specLessLandingFixture(t, "integration")
	policy, _, err := (commitrepo.Store{Root: f.root}).Policy()
	if err != nil {
		t.Fatal(err)
	}
	other := commitment.Outcome{ID: "B", Criteria: []commitment.Criterion{{ID: "B-accepted", Text: "B is accepted."}}, Deliverables: policy.Milestones[0].Outcomes[0].Deliverables}
	policy.Milestones[0].Outcomes = append(policy.Milestones[0].Outcomes, other)
	policy.ParallelGrants = []commitment.ParallelGrant{{Outcomes: []string{"delivery", "B"}}}
	commitmenttest.WritePolicy(t, f.root, policy)
	gitRun(t, f.root, "add", ".bench/commitment.json")
	gitRun(t, f.root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "approve B in parallel")
	result := runVerb(t, verbCreate, verbCall{root: f.root, home: f.home, args: []string{"--request", "ticket", "--label", "ticket", "--from", f.creation.Assignment.ID}})
	if result.exit != 0 {
		t.Fatalf("create = %d: %s", result.exit, result.stderr)
	}
	a, found, err := intent.FindAssignmentForRequest(f.root, "ticket")
	if err != nil || !found {
		t.Fatalf("assignment = %v, %v", found, err)
	}
	before, err := intent.Read(f.root)
	if err != nil {
		t.Fatal(err)
	}
	err = (commitrepo.Store{Root: a.Worktree}).Start("B", "ticket", "specs/x/spec.md")
	if err == nil || !strings.Contains(err.Error(), "already bound") {
		t.Fatalf("inherited sibling switched to B: %v", err)
	}
	after, err := intent.Read(f.root)
	if err != nil || !reflect.DeepEqual(before, after) {
		t.Fatalf("refusal changed admission: %v", err)
	}
}

func TestCommitmentRefusalOrder(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"usage", "unknown", "identity", "commitment"} {
		t.Run(kind, func(t *testing.T) {
			f := specLessLandingFixture(t, "integration")
			if err := (commitrepo.Store{Root: f.root}).Block("delivery", "waiting"); err != nil {
				t.Fatal(err)
			}
			args := []string{"--request", "ticket", "--label", "ticket", "--from", f.creation.Assignment.ID}
			want := "blocked"
			switch kind {
			case "usage":
				args = append(args, "--refresh")
				want = "--from with --refresh"
			case "unknown":
				args[len(args)-1] = "unknown-assignment"
				want = "no active assignment"
			case "identity":
				gitRun(t, f.creation.Path, "checkout", "--detach", "-q", "HEAD")
				want = "not on its assignment branch"
			}
			before, err := intent.Read(f.root)
			if err != nil {
				t.Fatal(err)
			}
			registered := gitOutput(t, f.root, "worktree", "list", "--porcelain")
			result := runVerb(t, verbCreate, verbCall{root: f.root, home: f.home, args: args})
			if result.exit == 0 || !strings.Contains(result.stderr, want) {
				t.Fatalf("precedence = %d %s %s; want %s", result.exit, result.stdout, result.stderr, want)
			}
			after, err := intent.Read(f.root)
			if err != nil || !reflect.DeepEqual(before, after) || registered != gitOutput(t, f.root, "worktree", "list", "--porcelain") {
				t.Fatalf("refusal changed ownership: %v", err)
			}
		})
	}
}
