package worktree

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

// decayToHusk leaves the assignment directory in place, holding content, with its git
// metadata entry gone: the shape a partially removed checkout leaves behind.
func decayToHusk(t *testing.T, path string) {
	t.Helper()
	mustWrite(t, filepath.Join(path, "keep.txt"), []byte("husk bytes\n"), 0o644)
	mustRemove(t, filepath.Join(path, ".git"))
}

// TestReleaseRegistrationSkipsUnrelatedSpecialControlRecords plants a FIFO where an
// unrelated private administration entry keeps its gitdir record. Releasing one
// registration sweeps the whole pool, so a stranger's record is reached too, and reading
// this one never returns.
//
// The release is driven directly rather than through a cleanup that reaches it.
// `git worktree list` reads the same pool and blocks on the same FIFO.
// No caller therefore ever gets far enough to show what this sweep does with it.
// The record is planted after the assignment's own registration is known, which is also
// the only window a live repository has for one.
func TestReleaseRegistrationSkipsUnrelatedSpecialControlRecords(t *testing.T) {
	t.Parallel()
	root, creation, _ := newOwnedAssignment(t, "husk-special-record")
	decayToHusk(t, creation.Path)
	common := gitOutput(t, root, "rev-parse", "--path-format=absolute", "--git-common-dir")
	stranger := filepath.Join(filepath.Clean(common), "worktrees", "stranger")
	mustMkdirAll(t, stranger, 0o700)
	mustNoError(t, syscall.Mkfifo(filepath.Join(stranger, "gitdir"), 0o600))

	done := make(chan error, 1)
	go func() { done <- releaseRegistration(root, creation.Path) }()
	select {
	case releaseErr := <-done:
		mustNoError(t, releaseErr)
	case <-time.After(noWriterDeadline):
		t.Fatal("releaseRegistration blocked on an unrelated special control record")
	}
	info, statErr := os.Lstat(filepath.Join(stranger, "gitdir"))
	requireTest(t, statErr == nil && info.Mode()&os.ModeNamedPipe != 0,
		"the release disturbed the unrelated administration entry: %v, %v", info, statErr)
	// Retiring the planted record before asking git anything: it blocks on this FIFO too.
	mustRemove(t, stranger)
	registrations := gitOutput(t, root, "worktree", "list", "--porcelain")
	requireTest(t, !strings.Contains(registrations, "worktree "+creation.Path),
		"registration survived the release:\n%s", registrations)
}
