package worktree

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/capability"
)

// noWriterDeadline bounds every read this package must decide without: a FIFO planted at
// an assignment path or in a discovered control record has no writer, so a read of one
// never returns. Running the subject off the test goroutine turns that into the waiting
// test's own failure instead of a package-wide timeout.
const noWriterDeadline = 15 * time.Second

// TestClassifyPathShapeUnknownFileAsParent pins the first ShapeUnknown return site: the
// path's own Lstat fails for a reason other than absence. A regular file makes an
// impossible parent for any further component, so the child never resolves.
func TestClassifyPathShapeUnknownFileAsParent(t *testing.T) {
	parent := filepath.Join(t.TempDir(), "not-a-directory")
	mustWrite(t, parent, []byte("not a directory\n"), 0o644)
	path := filepath.Join(parent, "child")

	shape, err := ClassifyPathShape(path)
	requireTest(t, shape == ShapeUnknown && err != nil,
		"ClassifyPathShape over a file-as-parent path = %v, %v; want %v and a non-nil error", shape, err, ShapeUnknown)
}

// TestClassifyPathShapeUnknownSelfSymlink pins the second ShapeUnknown return site: the
// path's Lstat succeeds (the symlink entry itself exists) but the follow-on Stat fails
// resolving it. A symlink pointing at itself is a loop no resolution terminates.
func TestClassifyPathShapeUnknownSelfSymlink(t *testing.T) {
	path := filepath.Join(t.TempDir(), "self-loop")
	if err := os.Symlink(path, path); err != nil {
		capability.Capability(t, capability.Symlink, "cannot create a symlink: "+err.Error())
	}

	shape, err := ClassifyPathShape(path)
	requireTest(t, shape == ShapeUnknown && err != nil,
		"ClassifyPathShape over a self-symlink = %v, %v; want %v and a non-nil error", shape, err, ShapeUnknown)
}

// TestClassifyPathShapeUnknownUnreadableGitEntry pins the third ShapeUnknown return
// site. The path itself resolves to a directory: Lstat and Stat both succeed.
// Traversing to a path needs permission on its parents, not on the path itself.
// But listing the .git entry inside it fails, because the directory's own mode
// denies search access into it. Root bypasses that mode entirely, so the fixture
// is skipped under root rather than silently passing without exercising the site.
func TestClassifyPathShapeUnknownUnreadableGitEntry(t *testing.T) {
	if os.Geteuid() == 0 {
		capability.Capability(t, capability.Privilege, "root bypasses directory permissions; cannot deny search access to simulate the site")
	}
	path := t.TempDir()
	t.Cleanup(func() { _ = os.Chmod(path, 0o700) })
	mustNoError(t, os.Chmod(path, 0o000))

	shape, err := ClassifyPathShape(path)
	requireTest(t, shape == ShapeUnknown && err != nil,
		"ClassifyPathShape over an unreadable directory = %v, %v; want %v and a non-nil error", shape, err, ShapeUnknown)
}

// TestClassifyPathShapeRefusesSpecialGitEntry plants a no-writer FIFO at the checkout's
// .git entry. This is the exact shape a fail-open classifier would hand to `git -C
// <path> rev-parse` without ever finishing. Git opens .git to follow a gitfile
// pointer. ClassifyPathShape must decide by shape alone, so it runs off the test
// goroutine. It fails the moment it misses the deadline instead of wedging the suite.
func TestClassifyPathShapeRefusesSpecialGitEntry(t *testing.T) {
	_, creation := newOwnedAssignment(t, "special-git-fifo")
	mustRemove(t, filepath.Join(creation.Path, ".git"))
	mustNoError(t, syscall.Mkfifo(filepath.Join(creation.Path, ".git"), 0o600))

	type shapeOutcome struct {
		shape PathShape
		err   error
	}
	shapeDone := make(chan shapeOutcome, 1)
	go func() {
		shape, err := ClassifyPathShape(creation.Path)
		shapeDone <- shapeOutcome{shape, err}
	}()
	select {
	case got := <-shapeDone:
		requireTest(t, got.err == nil && got.shape == ShapeSpecialMetadata,
			"ClassifyPathShape over FIFO .git = %v, %v; want %v, <nil>", got.shape, got.err, ShapeSpecialMetadata)
	case <-time.After(noWriterDeadline):
		t.Fatal("ClassifyPathShape blocked reading a no-writer FIFO .git entry")
	}

	info, err := os.Lstat(filepath.Join(creation.Path, ".git"))
	requireTest(t, err == nil && info.Mode()&os.ModeNamedPipe != 0,
		"the FIFO .git entry was disturbed: %v, %v", info, err)
}
