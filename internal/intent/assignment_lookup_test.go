package intent

import (
	"os"
	"path/filepath"
	"testing"
)

// TestAssignmentForWorktreeMatchesThroughASymlink covers the hostile-path row: a tree
// reached through a link is the same tree, and the canonical form is what says so. A
// comparison over path strings would leave the linked spelling unowned, and the caller
// would then refuse a checkout that has a live assignment.
func TestAssignmentForWorktreeMatchesThroughASymlink(t *testing.T) {
	root := newRepo(t)
	tree := filepath.Join(t.TempDir(), "tree")
	runGit(t, root, "worktree", "add", "-q", "--detach", tree, "HEAD")
	want := putActive(t, root, tree, "lookup-symlink")

	link := filepath.Join(t.TempDir(), "link")
	if err := os.Symlink(tree, link); err != nil {
		t.Fatalf("Symlink: %v", err)
	}
	got, owned := AssignmentForWorktree(link)
	if !owned || got.ID != want.ID {
		t.Fatalf("AssignmentForWorktree(link) = (%q, %t), want (%q, true)", got.ID, owned, want.ID)
	}
}

// TestAssignmentForWorktreeIgnoresARetiredAssignment pins the state predicate. Only an
// active assignment owns a tree; a cleanup-pending record names a tree whose phase is
// over, and adopting it would let a released request keep a live section.
func TestAssignmentForWorktreeIgnoresARetiredAssignment(t *testing.T) {
	root := newRepo(t)
	tree := filepath.Join(t.TempDir(), "tree")
	runGit(t, root, "worktree", "add", "-q", "--detach", tree, "HEAD")
	retired := putActive(t, root, tree, "lookup-retired")
	retired.State = StateCleanupPending
	if err := PutAssignment(root, retired); err != nil {
		t.Fatalf("PutAssignment: %v", err)
	}
	if got, owned := AssignmentForWorktree(tree); owned {
		t.Fatalf("AssignmentForWorktree = (%q, true), want no owner", got.ID)
	}
}

// TestAssignmentForWorktreeAnswersNoOwnerOutsideThePool pins the unowned answer. A stray
// checkout is a state to report, not a failure to raise.
func TestAssignmentForWorktreeAnswersNoOwnerOutsideThePool(t *testing.T) {
	root := newRepo(t)
	if got, owned := AssignmentForWorktree(root); owned {
		t.Fatalf("AssignmentForWorktree(primary) = (%q, true), want no owner", got.ID)
	}
}

// putActive registers one active assignment owning tree and returns the record.
func putActive(t *testing.T, root, tree, token string) Assignment {
	t.Helper()
	const owner = "0000000000000000000000000000000a"
	const id = "0000000000000000000000000000000b"
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
		State:        StateActive,
	}
	if err := PutAssignment(root, a); err != nil {
		t.Fatalf("PutAssignment: %v", err)
	}
	return a
}

func headOID(t *testing.T, root string) string {
	t.Helper()
	out, err := os.ReadFile(filepath.Join(root, ".git", "refs", "heads", "main"))
	if err != nil {
		t.Fatalf("read main ref: %v", err)
	}
	return string(out[:40])
}
