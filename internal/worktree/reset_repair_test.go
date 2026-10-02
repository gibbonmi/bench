package worktree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// commitIndexLockLane commits a lane at the merge target that leaves a stale index.lock in
// the target's administration directory. The lane passes, so the merge publishes, and the
// reconcile after the publication cannot take the index lock. It returns the lock path.
func commitIndexLockLane(t *testing.T, path string) string {
	t.Helper()
	lock := mustAdminPath(t, path, "index.lock")
	commitLaneManifest(t, path, gate.Phase{Name: "unit", Argv: []string{"sh", "-c", ": > " + sanitize.ShellQuote(lock)}})
	return lock
}

func TestResetApplyReconcilesAnUnfinishedMerge(t *testing.T) {
	t.Parallel()
	f := mergeFixture(t, "reset-merge")
	creation := f.created[0]
	commitInWorktree(t, creation.Path, "target.txt", "target\n", "target work")
	lock := commitIndexLockLane(t, creation.Path)
	incoming := commitOnDefault(t, f.root, "incoming.txt", "incoming\n")
	merged := runVerb(t, verbMerge, f.merge("--from", incoming, creation.Assignment.ID))
	requireTest(t, merged.exit == 3, "fixture merge = %d %s %s", merged.exit, merged.stdout, merged.stderr)
	tip := gitOutput(t, f.root, "rev-parse", creation.Assignment.Branch)
	requireTest(t, gitOutput(t, creation.Path, "status", "--porcelain=v1") != "", "fixture merge was reconciled")
	// The operator removes the stale lock before the repair, because the reset needs the
	// index lock too.
	mustRemove(t, lock)
	plan := runVerb(t, verbReset, f.call("--to", tip, creation.Assignment.ID))
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	result := runVerb(t, verbReset, f.call("--to", tip, creation.Assignment.ID, "--apply", fingerprint))
	requireTest(t, result.exit == 0 && gitOutput(t, creation.Path, "rev-parse", "HEAD") == tip && gitOutput(t, creation.Path, "status", "--porcelain=v1") == "", "merge repair = %d %s %s", result.exit, result.stdout, result.stderr)
	body, err := os.ReadFile(filepath.Join(creation.Path, "incoming.txt"))
	mustNoError(t, err)
	requireTest(t, string(body) == "incoming\n", "merge repair lost incoming bytes")
}

func TestResetApplyRepairsAShiftBranchCheckout(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-shift-repair")
	gitRun(t, f.creation.Path, "switch", "-c", "bench/shift-reset-repair")
	commitInWorktree(t, f.creation.Path, "shift", "shift work\n", "shift work")
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	gitRun(t, f.root, "worktree", "unlock", f.creation.Path)
	gitRun(t, f.root, "worktree", "lock", "--reason", "bench shift recovery: iteration failed", f.creation.Path)
	plan := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	result := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID, "--apply", fingerprint))
	requireTest(t, result.exit == 0, "shift repair = %d %s %s", result.exit, result.stdout, result.stderr)
	_, err := resolveAssignment(f.root, f.creation.Assignment.ID)
	mustNoError(t, err)
	requireTest(t, gitOutput(t, f.root, "rev-parse", "refs/heads/bench/shift-reset-repair") == tip &&
		gitOutput(t, f.creation.Path, "symbolic-ref", "HEAD") == f.creation.Assignment.Branch, "shift repair lost the shift branch or left HEAD off assignment")
}

func TestResetApplyKeepsTheTipOfADetachedCheckout(t *testing.T) {
	t.Parallel()
	checkResetOffBranchTip(t, false)
}

func TestResetApplyKeepsTheTipOfAShiftBranchCheckout(t *testing.T) {
	t.Parallel()
	checkResetOffBranchTip(t, true)
}

func checkResetOffBranchTip(t *testing.T, shift bool) {
	t.Helper()
	f := newOwnedAssignment(t, "reset-off-branch-tip")
	commitInWorktree(t, f.creation.Path, "ahead", "ahead\n", "ahead")
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	if shift {
		gitRun(t, f.creation.Path, "switch", "-c", "bench/shift-reset-tip", f.creation.Assignment.Start)
	} else {
		gitRun(t, f.creation.Path, "switch", "--detach", f.creation.Assignment.Start)
	}
	result := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, result.exit == 0 && strings.Contains(result.stdout, "preserve=envelope"), "detached plan = %d %s %s", result.exit, result.stdout, result.stderr)
	plan := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	result = runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID, "--apply", fingerprint))
	requireTest(t, result.exit == 0, "detached apply = %d %s %s", result.exit, result.stdout, result.stderr)
	ref := intent.ResetRefPrefix(f.creation.Assignment.OwnerID, f.creation.Assignment.ID) + "1"
	parents := strings.Fields(gitOutput(t, f.root, "show", "-s", "--format=%P", ref))
	requireTest(t, strings.Contains(" "+strings.Join(parents, " ")+" ", " "+tip+" "), "root lost detached checkout's tip: %#v", parents)
}

func TestResetApplyKeepsTheRewoundTipReachable(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-rewind")
	commitInWorktree(t, f.creation.Path, "ahead", "ahead\n", "ahead")
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	plan := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	result := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID, "--apply", fingerprint))
	requireTest(t, result.exit == 0, "rewind = %d %s %s", result.exit, result.stdout, result.stderr)
	manifest, ok := readRecoveryManifest(f.root, intent.ResetRefPrefix(f.creation.Assignment.OwnerID, f.creation.Assignment.ID)+"1")
	requireTest(t, ok && gitOutput(t, f.root, "rev-parse", f.creation.Assignment.Branch) == f.creation.Assignment.Start, "rewind failed or envelope absent")
	requireTest(t, gitOutput(t, f.root, "show", "-s", "--format=%P", manifest.Layers["working"]) == tip, "working payload lost previous tip")
}

func TestResetApplyRestoresTheExactLock(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-exact-lock")
	gitRun(t, f.root, "worktree", "unlock", f.creation.Path)
	gitRun(t, f.root, "worktree", "lock", "--reason", "foreign lock", f.creation.Path)
	plan := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	result := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID, "--apply", fingerprint))
	requireTest(t, result.exit == 0 && strings.Contains(result.stdout, "preserved=none"), "lock repair = %d %s %s", result.exit, result.stdout, result.stderr)
	_, err := resolveAssignment(f.root, f.creation.Assignment.ID)
	mustNoError(t, err)
}

func TestResetApplyReattachesWithoutAnEnvelope(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-reattach")
	gitRun(t, f.creation.Path, "switch", "--detach", "HEAD")
	plan := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	result := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID, "--apply", fingerprint))
	requireTest(t, result.exit == 0 && strings.Contains(result.stdout, "preserved=none") && strings.Contains(result.stdout, "restore=none"), "reattach = %d %s %s", result.exit, result.stdout, result.stderr)
	requireTest(t, gitOutput(t, f.creation.Path, "symbolic-ref", "HEAD") == f.creation.Assignment.Branch &&
		gitOutput(t, f.root, "for-each-ref", "--format=%(refname)", intent.ResetRefPrefix(f.creation.Assignment.OwnerID, f.creation.Assignment.ID)) == "", "pure repair wrote an envelope or stayed detached")
}
