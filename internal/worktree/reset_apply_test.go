package worktree

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/otelrecord"
)

func TestResetApplyRefusesAConsumedFingerprint(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-consumed")
	mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("dirty\n"), 0o644)
	plan := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	args := []string{"--to", f.creation.Assignment.Start, f.creation.Assignment.ID, "--apply", fingerprint}
	result := runVerb(t, verbReset, f.call(args...))
	requireTest(t, result.exit == 0, "first apply = %d %s %s", result.exit, result.stdout, result.stderr)
	prefix := intent.ResetRefPrefix(f.creation.Assignment.OwnerID, f.creation.Assignment.ID)
	before := gitOutput(t, f.root, "for-each-ref", "--format=%(refname) %(objectname)", prefix)
	result = runVerb(t, verbReset, f.call(args...))
	requireTest(t, result.exit == 1 && strings.Contains(result.stdout, "reset plan is stale") &&
		gitOutput(t, f.root, "for-each-ref", "--format=%(refname) %(objectname)", prefix) == before, "consumed apply = %d %s %s", result.exit, result.stdout, result.stderr)
}

func TestResetApplyLeavesTheAssignmentRecordUnchanged(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-ledger")
	before, err := intent.Assignments(f.root)
	mustNoError(t, err)
	mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("dirty\n"), 0o644)
	plan := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	result := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID, "--apply", fingerprint))
	requireTest(t, result.exit == 0, "apply = %d %s %s", result.exit, result.stdout, result.stderr)
	after, err := intent.Assignments(f.root)
	mustNoError(t, err)
	requireTest(t, reflect.DeepEqual(before, after) && len(after) == 1 && after[0].State == intent.StateActive && len(after[0].Recovery) == 0, "assignment record changed: %#v", after)
}

func TestResetApplyRecordsOneVerbSpan(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-span")
	mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("dirty\n"), 0o644)
	plan := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	result := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID, "--apply", fingerprint))
	requireTest(t, result.exit == 0, "apply = %d %s %s", result.exit, result.stdout, result.stderr)
	spans, err := otelrecord.ReadSpans(f.home, f.root)
	mustNoError(t, err)
	requireTest(t, len(spans) == 1 && spans[0].Seam == "worktree.reset" && spans[0].Attributes[otelrecord.AttrSubjectID] == f.creation.Assignment.ID, "reset spans = %#v", spans)
}

func TestResetApplyExitsThreeWhenTheMoveDidNotLand(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-silent-move")
	mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("dirty\n"), 0o644)
	j := defaultJoins()
	j.resetMove = func(string, string, string) error { return nil }
	plan := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	result := runVerb(t, verbReset, verbCall{root: f.root, home: f.home, joins: &j, args: []string{"--to", f.creation.Assignment.Start, f.creation.Assignment.ID, "--apply", fingerprint}})
	ref := intent.ResetRefPrefix(f.creation.Assignment.OwnerID, f.creation.Assignment.ID) + "1"
	requireTest(t, result.exit == 3 && strings.Contains(result.stdout, "preserved="+ref), "silent move = %d %s %s", result.exit, result.stdout, result.stderr)
}

func TestResetApplyExitsThreeOnAMoveFault(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-move-fault")
	mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("dirty\n"), 0o644)
	j := defaultJoins()
	j.resetMove = func(string, string, string) error { return errors.New("move failed") }
	plan := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	result := runVerb(t, verbReset, verbCall{root: f.root, home: f.home, joins: &j, args: []string{"--to", f.creation.Assignment.Start, f.creation.Assignment.ID, "--apply", fingerprint}})
	ref := intent.ResetRefPrefix(f.creation.Assignment.OwnerID, f.creation.Assignment.ID) + "1"
	requireTest(t, result.exit == 3 && strings.Contains(result.stdout, "preserved="+ref) && strings.Contains(result.stdout, "next=bench worktree reset --restore "+ref+" "+f.creation.Assignment.ID), "move fault = %d %s %s", result.exit, result.stdout, result.stderr)
	_, ok := readRecoveryManifest(f.root, ref)
	requireTest(t, ok, "move fault lost the envelope")
}

func TestResetApplyTakesTheCleanupLock(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-lock-attempt")
	mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("dirty\n"), 0o644)
	j := defaultJoins()
	var attempts []string
	j.cleanupLockAttempt = func(target string) { attempts = append(attempts, target) }
	result := runVerb(t, verbReset, verbCall{root: f.root, home: f.home, joins: &j, args: []string{"--to", f.creation.Assignment.Start, f.creation.Assignment.ID}})
	requireTest(t, result.exit == 0 && len(attempts) == 0, "plan took cleanup lock: %d %s %s %#v", result.exit, result.stdout, result.stderr, attempts)
	plan := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	result = runVerb(t, verbReset, verbCall{root: f.root, home: f.home, joins: &j, args: []string{"--to", f.creation.Assignment.Start, f.creation.Assignment.ID, "--apply", fingerprint}})
	requireTest(t, result.exit == 0 && len(attempts) == 1 && attempts[0] == f.creation.Path, "apply lock = %d %s %s %#v", result.exit, result.stdout, result.stderr, attempts)
}

func TestResetApplyRefusesAStalePlan(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-stale")
	mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("first\n"), 0o644)
	plan := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("second\n"), 0o644)
	plan = runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	wanted := plan.mustFingerprint(t)
	status := gitOutput(t, f.creation.Path, "status", "--porcelain=v1")
	result := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID, "--apply", fingerprint))
	requireTest(t, result.exit == 1 && strings.Contains(result.stdout, "reset plan is stale") && strings.Contains(result.stdout, "wanted="+wanted) &&
		strings.Contains(result.stdout, "next=bench worktree reset --to "+f.creation.Assignment.Start+" "+f.creation.Assignment.ID), "stale apply = %d %s %s", result.exit, result.stdout, result.stderr)
	requireTest(t, gitOutput(t, f.root, "for-each-ref", "--format=%(refname)", intent.ResetRefPrefix(f.creation.Assignment.OwnerID, f.creation.Assignment.ID)) == "" &&
		gitOutput(t, f.creation.Path, "status", "--porcelain=v1") == status, "stale apply wrote a ref or moved the checkout")
}

func TestResetApplyRefusesAnUnverifiedEnvelope(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-unverified")
	commitInWorktree(t, f.creation.Path, "ahead", "ahead\n", "ahead")
	mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("dirty\n"), 0o644)
	head := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	status := gitOutput(t, f.creation.Path, "status", "--porcelain=v1")
	j := defaultJoins()
	j.resetEnvelope = func(root string, plan resetPlan) (intent.Recovery, error) {
		envelope, err := writeResetEnvelope(root, plan)
		if err != nil {
			return envelope, err
		}
		tree := gitOutput(t, root, "rev-parse", envelope.Root+"^{tree}")
		broken, err := commitTree(root, tree, nil, "dropped payload\n")
		if err != nil {
			return envelope, err
		}
		gitRun(t, root, "update-ref", envelope.Ref, broken, envelope.Root)
		envelope.Root = broken
		return envelope, nil
	}
	plan := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	result := runVerb(t, verbReset, verbCall{root: f.root, home: f.home, joins: &j, args: []string{"--to", f.creation.Assignment.Start, f.creation.Assignment.ID, "--apply", fingerprint}})
	requireTest(t, gitOutput(t, f.creation.Path, "rev-parse", "HEAD") == head && gitOutput(t, f.creation.Path, "status", "--porcelain=v1") == status,
		"unverified envelope moved the checkout: %d %s %s", result.exit, result.stdout, result.stderr)
	requireTest(t, result.exit == 1 && strings.Contains(result.stdout, "reset envelope failed verification"), "unverified envelope = %d %s %s", result.exit, result.stdout, result.stderr)
}

func TestResetApplyPreservesTheLayersAndMovesTheCheckout(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-layers")
	commitInWorktree(t, f.creation.Path, ".gitignore", "untracked-dir/output\nignored-dir/\n", "ignore outputs")
	checkpoint := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	mustMkdirAll(t, filepath.Join(f.creation.Path, "untracked-dir"), 0o755)
	mustMkdirAll(t, filepath.Join(f.creation.Path, "ignored-dir"), 0o755)
	for _, path := range []string{"untracked-dir/keep.txt", "untracked-dir/output", "ignored-dir/output"} {
		mustWrite(t, filepath.Join(f.creation.Path, path), []byte(path+"\n"), 0o644)
	}
	mustNoError(t, os.Symlink("README.md", filepath.Join(f.creation.Path, "link.txt")))
	mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("staged\n"), 0o644)
	gitRun(t, f.creation.Path, "add", "README.md")
	mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("working\n"), 0o644)
	before, _, err := captureLayers(f.root, f.creation.Path, true, checkpoint)
	mustNoError(t, err)
	expected, ok := readRecoveryManifest(f.root, before)
	requireTest(t, ok, "fixture envelope is unreadable")
	plan := runVerb(t, verbReset, f.call("--to", checkpoint, f.creation.Assignment.ID))
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	result := runVerb(t, verbReset, f.call("--to", checkpoint, f.creation.Assignment.ID, "--apply", fingerprint))
	requireTest(t, result.exit == 0, "apply = %d %s %s", result.exit, result.stdout, result.stderr)
	ref := intent.ResetRefPrefix(f.creation.Assignment.OwnerID, f.creation.Assignment.ID) + "1"
	actual, ok := readRecoveryManifest(f.root, ref)
	requireTest(t, ok && actual.Tip == checkpoint, "reset manifest = %#v, readable=%v", actual, ok)
	for _, layer := range []string{"working", "staged"} {
		want := gitOutput(t, f.root, "rev-parse", expected.Layers[layer]+"^{tree}")
		got := gitOutput(t, f.root, "rev-parse", actual.Layers[layer]+"^{tree}")
		requireTest(t, got == want, "%s layer differs: %s != %s", layer, got, want)
	}
	requireTest(t, strings.Contains(result.stdout, "preserved="+ref) && strings.Contains(result.stdout, "restore=bench worktree reset --restore "+ref+" "+f.creation.Assignment.ID), "reset record = %s", result.stdout)
	for _, path := range []string{"untracked-dir/keep.txt", "link.txt"} {
		_, err := os.Lstat(filepath.Join(f.creation.Path, path))
		requireTest(t, os.IsNotExist(err), "untracked %s survives: %v", path, err)
	}
	for _, path := range []string{"untracked-dir/output", "ignored-dir/output"} {
		body, err := os.ReadFile(filepath.Join(f.creation.Path, path))
		mustNoError(t, err)
		requireTest(t, string(body) == path+"\n", "ignored bytes changed: %s", path)
	}
	requireTest(t, gitOutput(t, f.creation.Path, "rev-parse", "HEAD") == checkpoint &&
		gitOutput(t, f.root, "rev-parse", f.creation.Assignment.Branch) == checkpoint &&
		gitOutput(t, f.creation.Path, "symbolic-ref", "HEAD") == f.creation.Assignment.Branch &&
		gitOutput(t, f.creation.Path, "status", "--porcelain=v1") == "", "apply did not leave the checkpoint clean and attached")
}

func TestResetRefusesAnIgnoredCollision(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-collision")
	commitInWorktree(t, f.creation.Path, "output", "tracked output\n", "track output")
	commitInWorktree(t, f.creation.Path, "build", "tracked build\n", "track build")
	mustMkdirAll(t, filepath.Join(f.creation.Path, "out"), 0o755)
	commitInWorktree(t, f.creation.Path, "out/a.c", "tracked source\n", "track a directory")
	checkpoint := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	gitRun(t, f.creation.Path, "rm", "-q", "output", "build", "out/a.c")
	commitInWorktree(t, f.creation.Path, ".gitignore", "output\nbuild/\nout\nother.log\n", "ignore the outputs")
	mustMkdirAll(t, filepath.Join(f.creation.Path, "build"), 0o755)
	ignored := map[string]string{"output": "ignored output\n", "build/inner": "ignored build\n", "out": "ignored file over a directory\n", "other.log": "ignored log\n"}
	for path, body := range ignored {
		mustWrite(t, filepath.Join(f.creation.Path, path), []byte(body), 0o644)
	}
	head := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	result := runVerb(t, verbReset, f.call("--to", checkpoint, f.creation.Assignment.ID))
	requireTest(t, result.exit == 1 && strings.Contains(result.stdout, "refused{detail=ignored content would be overwritten}"), "collision plan = %d %s %s", result.exit, result.stdout, result.stderr)
	requireTest(t, strings.Contains(result.stdout, "refusal_paths[3]{path}:\n  build/inner\n  out\n  output\n"), "collision paths = %s", result.stdout)
	for path, body := range ignored {
		got, err := os.ReadFile(filepath.Join(f.creation.Path, path))
		mustNoError(t, err)
		requireTest(t, string(got) == body, "ignored bytes changed at %s: %q", path, got)
	}
	requireTest(t, gitOutput(t, f.creation.Path, "rev-parse", "HEAD") == head &&
		gitOutput(t, f.root, "for-each-ref", "--format=%(refname)", intent.ResetRefPrefix(f.creation.Assignment.OwnerID, f.creation.Assignment.ID)) == "",
		"collision refusal moved the checkout or wrote a ref")
}

func TestResetApplyKeepsIgnoredBytesAcrossAnIgnoreRuleChange(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-ignore-drift")
	commitInWorktree(t, f.creation.Path, ".gitignore", "build/\n", "ignore the build directory")
	mustMkdirAll(t, filepath.Join(f.creation.Path, "build"), 0o755)
	mustWrite(t, filepath.Join(f.creation.Path, "build/output"), []byte("build output\n"), 0o644)
	plan := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	result := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID, "--apply", fingerprint))
	ref := intent.ResetRefPrefix(f.creation.Assignment.OwnerID, f.creation.Assignment.ID) + "1"
	requireTest(t, result.exit == 3 && strings.Contains(result.stdout, "preserved="+ref) && strings.Contains(result.stdout, "next=bench worktree reset --restore "+ref+" "+f.creation.Assignment.ID),
		"ignore drift = %d %s %s", result.exit, result.stdout, result.stderr)
	body, err := os.ReadFile(filepath.Join(f.creation.Path, "build/output"))
	mustNoError(t, err)
	requireTest(t, string(body) == "build output\n", "ignored bytes changed: %q", body)
	requireTest(t, gitOutput(t, f.creation.Path, "rev-parse", "HEAD") == f.creation.Assignment.Start &&
		gitOutput(t, f.creation.Path, "status", "--porcelain=v1") == "?? build/", "checkpoint state = %s", gitOutput(t, f.creation.Path, "status", "--porcelain=v1"))
	plan = runVerb(t, verbReset, f.call("--restore", ref, f.creation.Assignment.ID))
	requireTest(t, isRestorePlanOf(plan, ref), "restore plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint = plan.mustFingerprint(t)
	result = runVerb(t, verbReset, f.call("--restore", ref, f.creation.Assignment.ID, "--apply", fingerprint))
	second := intent.ResetRefPrefix(f.creation.Assignment.OwnerID, f.creation.Assignment.ID) + "2"
	requireTest(t, result.exit == 0 && strings.Contains(result.stdout, "preserved="+second), "named restore = %d %s %s", result.exit, result.stdout, result.stderr)
	manifest, ok := readRecoveryManifest(f.root, second)
	requireTest(t, ok && gitOutput(t, f.root, "show", manifest.Layers["working"]+":build/output") == "build output", "second envelope lost the drifted bytes: %#v", manifest)
}

func TestResetApplyExitsThreeWithoutAnEnvelope(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-fault-no-envelope")
	gitRun(t, f.creation.Path, "switch", "--detach", "HEAD")
	j := defaultJoins()
	j.resetMove = func(string, string, string) error { return errors.New("move failed") }
	plan := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	result := runVerb(t, verbReset, verbCall{root: f.root, home: f.home, joins: &j, args: []string{"--to", f.creation.Assignment.Start, f.creation.Assignment.ID, "--apply", fingerprint}})
	requireTest(t, result.exit == 3 && strings.Contains(result.stdout, "preserved=none") &&
		strings.Contains(result.stdout, "next=bench worktree reset --to "+f.creation.Assignment.Start+" "+f.creation.Assignment.ID+"}"), "fault without envelope = %d %s %s", result.exit, result.stdout, result.stderr)
}
