package worktree

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

func restoreFixture(t *testing.T) restoredAssignment {
	t.Helper()
	root, creation, home := newOwnedAssignment(t, "restore")
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("staged\n"), 0o644)
	gitRun(t, creation.Path, "add", "README.md")
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("working\n"), 0o644)
	mustWrite(t, filepath.Join(creation.Path, "untracked"), []byte("untracked\n"), 0o644)
	plan := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID}})
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	apply := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID, "--apply", plan.mustFingerprint(t)}})
	requireTest(t, apply.exit == 0, "fixture reset = %d %s %s", apply.exit, apply.stdout, apply.stderr)
	return restoredAssignment{root: root, creation: creation, home: home, ref: intent.ResetRefPrefix(creation.Assignment.OwnerID, creation.Assignment.ID) + "1"}
}

// isRestorePlanOf reports whether result is a successful restore plan for the envelope
// at ref.
func isRestorePlanOf(result verbResult, ref string) bool {
	return result.exit == 0 && strings.Contains(result.stdout, "mode=restore") && strings.Contains(result.stdout, "envelope="+ref)
}

func TestResetRestoreReturnsThePreservedStateByteExact(t *testing.T) {
	t.Parallel()
	f := restoreFixture(t)
	plan := runVerb(t, verbReset, f.call("--restore", f.ref, f.creation.Assignment.ID))
	requireTest(t, isRestorePlanOf(plan, f.ref), "restore plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	result := runVerb(t, verbReset, f.call("--restore", f.ref, f.creation.Assignment.ID, "--apply", plan.mustFingerprint(t)))
	requireTest(t, result.exit == 0 && strings.Contains(result.stdout, "mode=restore"), "restore = %d %s %s", result.exit, result.stdout, result.stderr)
	requireTest(t, gitOutput(t, f.creation.Path, "rev-parse", "HEAD") == f.creation.Assignment.Start &&
		gitOutput(t, f.creation.Path, "symbolic-ref", "HEAD") == f.creation.Assignment.Branch, "restore changed the base or attachment")
	requireTest(t, gitOutput(t, f.creation.Path, "status", "--porcelain=v1") == "MM README.md\n?? untracked", "restore lost staged or untracked status")
	requireTest(t, gitOutput(t, f.creation.Path, "show", ":README.md") == "staged", "restore lost staged bytes")
	assertRestoreLayers(t, f.root, f.creation.Path, f.ref)
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
	f := restoreFixture(t)
	j := defaultJoins()
	j.resetLayers = func(path string, manifest recoveryManifest) error {
		if err := restoreResetLayers(path, manifest); err != nil {
			return err
		}
		mustWrite(t, filepath.Join(path, "README.md"), []byte("wrong layer\n"), 0o644)
		return nil
	}
	plan := runVerb(t, verbReset, f.call("--restore", f.ref, f.creation.Assignment.ID))
	requireTest(t, isRestorePlanOf(plan, f.ref), "restore plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	call := f.call("--restore", f.ref, f.creation.Assignment.ID, "--apply", plan.mustFingerprint(t))
	call.joins = &j
	result := runVerb(t, verbReset, call)
	requireTest(t, result.exit == 3 && strings.Contains(result.stdout, "next=bench worktree reset --restore "+f.ref+" "+f.creation.Assignment.ID), "recapture mismatch = %d %s %s", result.exit, result.stdout, result.stderr)
}

func TestResetRestoreExitsThreeOnALayerFault(t *testing.T) {
	t.Parallel()
	f := restoreFixture(t)
	mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("current\n"), 0o644)
	j := defaultJoins()
	j.resetLayers = func(string, recoveryManifest) error { return errors.New("layer failed") }
	plan := runVerb(t, verbReset, f.call("--restore", f.ref, f.creation.Assignment.ID))
	requireTest(t, isRestorePlanOf(plan, f.ref), "restore plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	call := f.call("--restore", f.ref, f.creation.Assignment.ID, "--apply", plan.mustFingerprint(t))
	call.joins = &j
	result := runVerb(t, verbReset, call)
	preserved := intent.ResetRefPrefix(f.creation.Assignment.OwnerID, f.creation.Assignment.ID) + "2"
	requireTest(t, result.exit == 3 && strings.Contains(result.stdout, "preserved="+preserved) &&
		strings.Contains(result.stdout, "next=bench worktree reset --restore "+preserved+" "+f.creation.Assignment.ID), "layer fault = %d %s %s", result.exit, result.stdout, result.stderr)
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
	plan := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID}})
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	result := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID, "--apply", plan.mustFingerprint(t)}})
	requireTest(t, result.exit == 0, "fixture shift reset = %d %s %s", result.exit, result.stdout, result.stderr)
	ref := intent.ResetRefPrefix(creation.Assignment.OwnerID, creation.Assignment.ID) + "1"
	plan = runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--restore", ref, creation.Assignment.ID}})
	requireTest(t, isRestorePlanOf(plan, ref), "restore plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	result = runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--restore", ref, creation.Assignment.ID, "--apply", plan.mustFingerprint(t)}})
	requireTest(t, result.exit == 0 && gitOutput(t, root, "rev-parse", creation.Assignment.Branch) == tip &&
		gitOutput(t, creation.Path, "rev-parse", "HEAD") == base && gitOutput(t, creation.Path, "rev-parse", "--abbrev-ref", "HEAD") == "HEAD", "off-branch restore = %d %s %s", result.exit, result.stdout, result.stderr)
	requireTest(t, strings.Contains(result.stdout, "ref=detached") && strings.Contains(result.stdout, "next=bench worktree reset --to "+tip+" "+creation.Assignment.ID), "restore lacks reattach command: %s", result.stdout)
	requireTest(t, gitOutput(t, root, "rev-parse", "refs/heads/bench/shift-restore") == base, "restore changed the shift branch")
	assertRestoreLayers(t, root, creation.Path, ref)
}

func TestResetRestorePreservesTheCurrentStateFirst(t *testing.T) {
	t.Parallel()
	f := restoreFixture(t)
	mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("current staged\n"), 0o644)
	gitRun(t, f.creation.Path, "add", "README.md")
	mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("current working\n"), 0o644)
	plan := runVerb(t, verbReset, f.call("--restore", f.ref, f.creation.Assignment.ID))
	requireTest(t, isRestorePlanOf(plan, f.ref), "restore plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	result := runVerb(t, verbReset, f.call("--restore", f.ref, f.creation.Assignment.ID, "--apply", plan.mustFingerprint(t)))
	preserved := intent.ResetRefPrefix(f.creation.Assignment.OwnerID, f.creation.Assignment.ID) + "2"
	requireTest(t, result.exit == 0 && strings.Contains(result.stdout, "preserved="+preserved), "preserving restore = %d %s %s", result.exit, result.stdout, result.stderr)
	plan = runVerb(t, verbReset, f.call("--restore", preserved, f.creation.Assignment.ID))
	requireTest(t, isRestorePlanOf(plan, preserved), "restore plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	result = runVerb(t, verbReset, f.call("--restore", preserved, f.creation.Assignment.ID, "--apply", plan.mustFingerprint(t)))
	requireTest(t, result.exit == 0, "restore current = %d %s %s", result.exit, result.stdout, result.stderr)
	requireTest(t, gitOutput(t, f.creation.Path, "show", ":README.md") == "current staged", "current staged bytes were lost")
	assertRestoreLayers(t, f.root, f.creation.Path, preserved)
}

func TestResetRestoreOmitsReattachWhenAttached(t *testing.T) {
	t.Parallel()
	f := restoreFixture(t)
	plan := runVerb(t, verbReset, f.call("--restore", f.ref, f.creation.Assignment.ID))
	requireTest(t, isRestorePlanOf(plan, f.ref), "restore plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	result := runVerb(t, verbReset, f.call("--restore", f.ref, f.creation.Assignment.ID, "--apply", plan.mustFingerprint(t)))
	requireTest(t, result.exit == 0 && !strings.Contains(result.stdout, "next=") && strings.Contains(result.stdout, "ref="+f.creation.Assignment.Branch), "attached restore = %d %s %s", result.exit, result.stdout, result.stderr)
}
