package worktree

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
)

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
// site: the path itself resolves to a directory (Lstat and Stat both succeed, since
// traversing to a path needs permission on its parents, not on the path itself), but
// listing the .git entry inside it fails, because the directory's own mode denies
// search access into it. Root bypasses that mode entirely, so the fixture is skipped
// under root rather than silently passing without exercising the site.
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
