package intent

import (
	"os"
	"path/filepath"
	"testing"
)

// TestAssignmentsOwningAnswersEveryStateInLedgerOrder pins the any-state contract the
// selector's refusals depend on. A match that dropped a retired row would leave the
// state refusal with nothing to name, and the operator would read "unassigned" for a
// tree the ledger records. The query arrives through a link, because the canonical
// form is what makes two spellings of one tree compare equal.
func TestAssignmentsOwningAnswersEveryStateInLedgerOrder(t *testing.T) {
	root := newRepo(t)
	tree := filepath.Join(t.TempDir(), "tree")
	runGit(t, root, "worktree", "add", "-q", "--detach", tree, "HEAD")
	retired := putOwner(t, root, tree, "owning-retired", "0000000000000000000000000000000c", StateCleanupPending)
	active := putOwner(t, root, tree, "owning-active", "0000000000000000000000000000000d", StateActive)

	link := filepath.Join(t.TempDir(), "link")
	if err := os.Symlink(tree, link); err != nil {
		t.Fatalf("Symlink: %v", err)
	}
	assignments, err := Assignments(root)
	if err != nil {
		t.Fatalf("Assignments: %v", err)
	}

	owning := AssignmentsOwning(assignments, link)
	if len(owning) != 2 {
		t.Fatalf("AssignmentsOwning returned %d rows, want 2", len(owning))
	}
	if owning[0].ID != retired.ID || owning[1].ID != active.ID {
		t.Fatalf("AssignmentsOwning = (%q, %q), want (%q, %q) in ledger order", owning[0].ID, owning[1].ID, retired.ID, active.ID)
	}
}

// TestAssignmentsOwningResolvesTheRecordedWorktree pins the row side of the compare. A
// ledger written through a linked pool home records the linked spelling, and a caller
// then arrives with the resolved path. A match that compared the recorded string raw
// would leave that row unowned, and the selector would answer unassigned for a tree the
// ledger owns.
func TestAssignmentsOwningResolvesTheRecordedWorktree(t *testing.T) {
	root := newRepo(t)
	tree := filepath.Join(t.TempDir(), "tree")
	runGit(t, root, "worktree", "add", "-q", "--detach", tree, "HEAD")
	link := filepath.Join(t.TempDir(), "link")
	if err := os.Symlink(tree, link); err != nil {
		t.Fatalf("Symlink: %v", err)
	}
	recorded := putOwner(t, root, link, "owning-linked-record", "0000000000000000000000000000000e", StateActive)

	assignments, err := Assignments(root)
	if err != nil {
		t.Fatalf("Assignments: %v", err)
	}
	resolved, err := filepath.EvalSymlinks(tree)
	if err != nil {
		t.Fatalf("EvalSymlinks: %v", err)
	}

	owning := AssignmentsOwning(assignments, resolved)
	if len(owning) != 1 || owning[0].ID != recorded.ID {
		t.Fatalf("AssignmentsOwning(%q) returned %d rows, want one row recorded at %q", resolved, len(owning), link)
	}
}

// putOwner registers one assignment owning tree under the id and state the caller names.
func putOwner(t *testing.T, root, tree, token, id string, state AssignmentState) Assignment {
	t.Helper()
	const owner = "0000000000000000000000000000000a"
	a := Assignment{
		Schema:       AssignmentRecordSchema,
		ID:           id,
		OwnerID:      owner,
		Request:      RequestDigest(token),
		RequestToken: token,
		Label:        token,
		Start:        headOID(t, root),
		Branch:       AssignmentBranchRef(owner, id),
		Worktree:     filepath.Clean(tree),
		State:        state,
	}
	if err := PutAssignment(root, a); err != nil {
		t.Fatalf("PutAssignment: %v", err)
	}
	return a
}
