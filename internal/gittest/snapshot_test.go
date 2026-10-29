package gittest

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestSnapshotSpecialEntry replaces a tracked regular file with a writerless FIFO. A
// read of that FIFO blocks, so the snapshot must refuse the entry before it reads it.
func TestSnapshotSpecialEntry(t *testing.T) {
	root := RepoOnBranch(t, "main")
	tracked := filepath.Join(root, "tracked")
	writeKitCopyFile(t, root, "tracked", "tracked bytes", 0o644)
	MustRun(t, root, "add", "tracked")
	Commit(t, root, "fixture")
	// FIFOWorktreeAdmin owns the FIFO plant and its capability skip; the rename moves its
	// FIFO onto the tracked path within the same file system.
	if err := os.Rename(FIFOWorktreeAdmin(t, root, "snapshot"), tracked); err != nil {
		t.Fatal(err)
	}
	refused := make(chan error, 1)
	go func() { refused <- CommitWorkingTree(root, t.TempDir()) }()
	select {
	case err := <-refused:
		if err == nil || !strings.Contains(err.Error(), "not regular: tracked") {
			t.Fatalf("snapshot of a FIFO entry = %v; want the special-entry refusal", err)
		}
	case <-time.After(bounded):
		t.Fatal("snapshot blocked on the FIFO entry")
	}
}

// TestSnapshotCommandFailure gives each snapshot function a tree with nothing to commit,
// so its git commit fails. Each function returns the process error to its caller.
func TestSnapshotCommandFailure(t *testing.T) {
	for name, snapshot := range map[string]func() error{
		"CommitAll":         func() error { return CommitAll(RepoOnBranch(t, "main"), "nothing") },
		"CommitWorkingTree": func() error { return CommitWorkingTree(Repo(t), t.TempDir()) },
	} {
		err := snapshot()
		var exit *exec.ExitError
		if !errors.As(err, &exit) || !strings.Contains(err.Error(), "commit") {
			t.Fatalf("%s with nothing to commit = %v; want the git commit process error", name, err)
		}
	}
}
