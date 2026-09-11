package worktree

import (
	"path/filepath"
	"strings"
	"testing"
)

func resetFingerprint(t *testing.T, root, home, checkpoint, target string) string {
	t.Helper()
	code, out, errout := runReset(t, root, home, "--to", checkpoint, target)
	requireTest(t, code == 0, "fingerprint plan = %d %s %s", code, out, errout)
	_, rest, found := strings.Cut(out, "fingerprint=")
	requireTest(t, found, "plan lacks fingerprint: %s", out)
	value, _, found := strings.Cut(rest, "}")
	requireTest(t, found && value != "", "invalid fingerprint cell: %s", out)
	return value
}

func TestResetFingerprintTracksTheCheckout(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-fingerprint")
	commitInWorktree(t, creation.Path, "ahead", "ahead\n", "ahead")
	before := resetFingerprint(t, root, home, creation.Assignment.Start, creation.Assignment.ID)
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("dirty\n"), 0o644)
	after := resetFingerprint(t, root, home, creation.Assignment.Start, creation.Assignment.ID)
	requireTest(t, before != after, "tracked edit did not change fingerprint: %s", after)
}

func TestResetFingerprintTracksTheContent(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-content")
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("first\n"), 0o644)
	before := resetFingerprint(t, root, home, creation.Assignment.Start, creation.Assignment.ID)
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("second\n"), 0o644)
	after := resetFingerprint(t, root, home, creation.Assignment.Start, creation.Assignment.ID)
	requireTest(t, before != after, "content edit did not change fingerprint: %s", after)
}

func TestResetFingerprintIgnoresIgnoredFiles(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-ignored")
	commitInWorktree(t, creation.Path, ".gitignore", "output\n", "ignore output")
	before := resetFingerprint(t, root, home, creation.Assignment.Start, creation.Assignment.ID)
	mustWrite(t, filepath.Join(creation.Path, "output"), []byte("build\n"), 0o644)
	after := resetFingerprint(t, root, home, creation.Assignment.Start, creation.Assignment.ID)
	requireTest(t, before == after, "ignored output changed fingerprint: %s != %s", before, after)
}

func TestResetFingerprintTracksTheTip(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-tip")
	commitInWorktree(t, creation.Path, "ahead", "one\n", "ahead")
	before := resetFingerprint(t, root, home, creation.Assignment.Start, creation.Assignment.ID)
	commitInWorktree(t, creation.Path, "ahead", "two\n", "ahead again")
	after := resetFingerprint(t, root, home, creation.Assignment.Start, creation.Assignment.ID)
	requireTest(t, before != after, "new commit did not change fingerprint")
}

func TestResetFingerprintTracksTheRef(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-ref")
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("dirty\n"), 0o644)
	before := resetFingerprint(t, root, home, creation.Assignment.Start, creation.Assignment.ID)
	gitRun(t, creation.Path, "switch", "--detach", "HEAD")
	after := resetFingerprint(t, root, home, creation.Assignment.Start, creation.Assignment.ID)
	requireTest(t, before != after, "detach did not change fingerprint")
}

func TestResetFingerprintTracksTheCheckpoint(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-checkpoint")
	commitInWorktree(t, creation.Path, "ahead", "ahead\n", "ahead")
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("dirty\n"), 0o644)
	before := resetFingerprint(t, root, home, creation.Assignment.Start, creation.Assignment.ID)
	after := resetFingerprint(t, root, home, gitOutput(t, creation.Path, "rev-parse", "HEAD"), creation.Assignment.ID)
	requireTest(t, before != after, "checkpoint did not change fingerprint")
}
