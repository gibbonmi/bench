package worktree

import (
	"errors"
	"path/filepath"
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
