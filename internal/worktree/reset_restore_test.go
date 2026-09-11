package worktree

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

func restoreFixture(t *testing.T) (string, Creation, string, string) {
	t.Helper()
	root, creation, home := newOwnedAssignment(t, "restore")
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("staged\n"), 0o644)
	gitRun(t, creation.Path, "add", "README.md")
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("working\n"), 0o644)
	mustWrite(t, filepath.Join(creation.Path, "untracked"), []byte("untracked\n"), 0o644)
	fingerprint := resetFingerprint(t, root, home, creation.Assignment.Start, creation.Assignment.ID)
	code, out, errout := runReset(t, root, home, "--to", creation.Assignment.Start, creation.Assignment.ID, "--apply", fingerprint)
	requireTest(t, code == 0, "fixture reset = %d %s %s", code, out, errout)
	return root, creation, home, intent.ResetRefPrefix(creation.Assignment.OwnerID, creation.Assignment.ID) + "1"
}

func restoreFingerprint(t *testing.T, root, home, ref, target string) string {
	t.Helper()
	code, out, errout := runReset(t, root, home, "--restore", ref, target)
	requireTest(t, code == 0 && strings.Contains(out, "mode=restore") && strings.Contains(out, "envelope="+ref), "restore plan = %d %s %s", code, out, errout)
	_, rest, found := strings.Cut(out, "fingerprint=")
	value, _, ended := strings.Cut(rest, "}")
	requireTest(t, found && ended && value != "" && value != "none", "restore fingerprint = %s", out)
	return value
}

func TestResetRestoreReturnsThePreservedStateByteExact(t *testing.T) {
	t.Parallel()
	root, creation, home, ref := restoreFixture(t)
	fingerprint := restoreFingerprint(t, root, home, ref, creation.Assignment.ID)
	code, out, errout := runReset(t, root, home, "--restore", ref, creation.Assignment.ID, "--apply", fingerprint)
	requireTest(t, code == 0 && strings.Contains(out, "mode=restore"), "restore = %d %s %s", code, out, errout)
	requireTest(t, gitOutput(t, creation.Path, "rev-parse", "HEAD") == creation.Assignment.Start &&
		gitOutput(t, creation.Path, "symbolic-ref", "HEAD") == creation.Assignment.Branch, "restore changed the base or attachment")
	requireTest(t, gitOutput(t, creation.Path, "status", "--porcelain=v1") == "MM README.md\n?? untracked", "restore lost staged or untracked status")
	requireTest(t, gitOutput(t, creation.Path, "show", ":README.md") == "staged", "restore lost staged bytes")
	assertRestoreLayers(t, root, creation.Path, ref)
}

func assertRestoreLayers(t *testing.T, root, path, ref string) {
	t.Helper()
	expected, ok := readRecoveryManifest(root, ref)
	requireTest(t, ok, "expected envelope unreadable")
	recaptured, _, err := captureLayers(root, path, true, expected.Tip)
	mustNoError(t, err)
	actual, ok := readRecoveryManifest(root, recaptured)
	requireTest(t, ok && len(actual.Layers) == len(expected.Layers), "restore layer set = %#v, want %#v", actual, expected)
	for layer, payload := range expected.Layers {
		requireTest(t, gitOutput(t, root, "rev-parse", actual.Layers[layer]+"^{tree}") == gitOutput(t, root, "rev-parse", payload+"^{tree}"), "restored %s differs", layer)
	}
}

func TestResetRestoreExitsThreeOnARecaptureMismatch(t *testing.T) {
	t.Parallel()
	root, creation, home, ref := restoreFixture(t)
	j := defaultJoins()
	j.resetLayers = func(path string, manifest recoveryManifest) error {
		if err := restoreResetLayers(path, manifest); err != nil {
			return err
		}
		mustWrite(t, filepath.Join(path, "README.md"), []byte("wrong layer\n"), 0o644)
		return nil
	}
	fingerprint := restoreFingerprint(t, root, home, ref, creation.Assignment.ID)
	code, out, errout := runResetWith(t, j, root, home, "--restore", ref, creation.Assignment.ID, "--apply", fingerprint)
	requireTest(t, code == 3 && strings.Contains(out, "next=bench worktree reset --restore "+ref+" "+creation.Assignment.ID), "recapture mismatch = %d %s %s", code, out, errout)
}

func TestResetRestoreExitsThreeOnALayerFault(t *testing.T) {
	t.Parallel()
	root, creation, home, ref := restoreFixture(t)
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("current\n"), 0o644)
	j := defaultJoins()
	j.resetLayers = func(string, recoveryManifest) error { return errors.New("layer failed") }
	fingerprint := restoreFingerprint(t, root, home, ref, creation.Assignment.ID)
	code, out, errout := runResetWith(t, j, root, home, "--restore", ref, creation.Assignment.ID, "--apply", fingerprint)
	preserved := intent.ResetRefPrefix(creation.Assignment.OwnerID, creation.Assignment.ID) + "2"
	requireTest(t, code == 3 && strings.Contains(out, "preserved="+preserved) &&
		strings.Contains(out, "next=bench worktree reset --restore "+preserved+" "+creation.Assignment.ID), "layer fault = %d %s %s", code, out, errout)
}

func TestResetRestoreReturnsAnOffBranchCapture(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "restore-off-branch")
	commitInWorktree(t, creation.Path, "assignment", "assignment\n", "assignment tip")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	gitRun(t, creation.Path, "switch", "-c", "bench/shift-restore", creation.Assignment.Start)
	commitInWorktree(t, creation.Path, "shift", "shift\n", "shift tip")
	base := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("shift dirty\n"), 0o644)
	fingerprint := resetFingerprint(t, root, home, creation.Assignment.Start, creation.Assignment.ID)
	code, out, errout := runReset(t, root, home, "--to", creation.Assignment.Start, creation.Assignment.ID, "--apply", fingerprint)
	requireTest(t, code == 0, "fixture shift reset = %d %s %s", code, out, errout)
	ref := intent.ResetRefPrefix(creation.Assignment.OwnerID, creation.Assignment.ID) + "1"
	fingerprint = restoreFingerprint(t, root, home, ref, creation.Assignment.ID)
	code, out, errout = runReset(t, root, home, "--restore", ref, creation.Assignment.ID, "--apply", fingerprint)
	requireTest(t, code == 0 && gitOutput(t, root, "rev-parse", creation.Assignment.Branch) == tip &&
		gitOutput(t, creation.Path, "rev-parse", "HEAD") == base && gitOutput(t, creation.Path, "rev-parse", "--abbrev-ref", "HEAD") == "HEAD", "off-branch restore = %d %s %s", code, out, errout)
	requireTest(t, strings.Contains(out, "ref=detached") && strings.Contains(out, "next=bench worktree reset --to "+tip+" "+creation.Assignment.ID), "restore lacks reattach command: %s", out)
	requireTest(t, gitOutput(t, root, "rev-parse", "refs/heads/bench/shift-restore") == base, "restore changed the shift branch")
	assertRestoreLayers(t, root, creation.Path, ref)
}

func TestResetRestorePreservesTheCurrentStateFirst(t *testing.T) {
	t.Parallel()
	root, creation, home, ref := restoreFixture(t)
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("current staged\n"), 0o644)
	gitRun(t, creation.Path, "add", "README.md")
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("current working\n"), 0o644)
	fingerprint := restoreFingerprint(t, root, home, ref, creation.Assignment.ID)
	code, out, errout := runReset(t, root, home, "--restore", ref, creation.Assignment.ID, "--apply", fingerprint)
	preserved := intent.ResetRefPrefix(creation.Assignment.OwnerID, creation.Assignment.ID) + "2"
	requireTest(t, code == 0 && strings.Contains(out, "preserved="+preserved), "preserving restore = %d %s %s", code, out, errout)
	fingerprint = restoreFingerprint(t, root, home, preserved, creation.Assignment.ID)
	code, out, errout = runReset(t, root, home, "--restore", preserved, creation.Assignment.ID, "--apply", fingerprint)
	requireTest(t, code == 0, "restore current = %d %s %s", code, out, errout)
	requireTest(t, gitOutput(t, creation.Path, "show", ":README.md") == "current staged", "current staged bytes were lost")
	assertRestoreLayers(t, root, creation.Path, preserved)
}

func TestResetRestoreOmitsReattachWhenAttached(t *testing.T) {
	t.Parallel()
	root, creation, home, ref := restoreFixture(t)
	fingerprint := restoreFingerprint(t, root, home, ref, creation.Assignment.ID)
	code, out, errout := runReset(t, root, home, "--restore", ref, creation.Assignment.ID, "--apply", fingerprint)
	requireTest(t, code == 0 && !strings.Contains(out, "next=") && strings.Contains(out, "ref="+creation.Assignment.Branch), "attached restore = %d %s %s", code, out, errout)
}
