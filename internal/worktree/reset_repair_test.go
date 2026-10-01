package worktree

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

func TestResetApplyReconcilesAnUnfinishedMerge(t *testing.T) {
	t.Parallel()
	j, root, home, _, created := mergeFixture(t, "reset-merge")
	creation := created[0]
	commitInWorktree(t, creation.Path, "target.txt", "target\n", "target work")
	incoming := commitOnDefault(t, root, "incoming.txt", "incoming\n")
	j.mergeReconcile = func(string, string) error { return errors.New("reconcile fault") }
	code, out, errout := runMerge(t, j, root, home, "--from", incoming, creation.Assignment.ID)
	requireTest(t, code == 3, "fixture merge = %d %s %s", code, out, errout)
	tip := gitOutput(t, root, "rev-parse", creation.Assignment.Branch)
	requireTest(t, gitOutput(t, creation.Path, "status", "--porcelain=v1") != "", "fixture merge was reconciled")
	plan := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", tip, creation.Assignment.ID}})
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	result := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", tip, creation.Assignment.ID, "--apply", fingerprint}})
	requireTest(t, result.exit == 0 && gitOutput(t, creation.Path, "rev-parse", "HEAD") == tip && gitOutput(t, creation.Path, "status", "--porcelain=v1") == "", "merge repair = %d %s %s", result.exit, result.stdout, result.stderr)
	body, err := os.ReadFile(filepath.Join(creation.Path, "incoming.txt"))
	mustNoError(t, err)
	requireTest(t, string(body) == "incoming\n", "merge repair lost incoming bytes")
}

func TestResetApplyRepairsAShiftBranchCheckout(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-shift-repair")
	gitRun(t, creation.Path, "switch", "-c", "bench/shift-reset-repair")
	commitInWorktree(t, creation.Path, "shift", "shift work\n", "shift work")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	gitRun(t, root, "worktree", "unlock", creation.Path)
	gitRun(t, root, "worktree", "lock", "--reason", "bench shift recovery: iteration failed", creation.Path)
	plan := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID}})
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	result := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID, "--apply", fingerprint}})
	requireTest(t, result.exit == 0, "shift repair = %d %s %s", result.exit, result.stdout, result.stderr)
	_, err := resolveAssignment(root, creation.Assignment.ID)
	mustNoError(t, err)
	requireTest(t, gitOutput(t, root, "rev-parse", "refs/heads/bench/shift-reset-repair") == tip &&
		gitOutput(t, creation.Path, "symbolic-ref", "HEAD") == creation.Assignment.Branch, "shift repair lost the shift branch or left HEAD off assignment")
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
	root, creation, home := newOwnedAssignment(t, "reset-off-branch-tip")
	commitInWorktree(t, creation.Path, "ahead", "ahead\n", "ahead")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	if shift {
		gitRun(t, creation.Path, "switch", "-c", "bench/shift-reset-tip", creation.Assignment.Start)
	} else {
		gitRun(t, creation.Path, "switch", "--detach", creation.Assignment.Start)
	}
	result := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID}})
	requireTest(t, result.exit == 0 && strings.Contains(result.stdout, "preserve=envelope"), "detached plan = %d %s %s", result.exit, result.stdout, result.stderr)
	plan := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID}})
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	result = runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID, "--apply", fingerprint}})
	requireTest(t, result.exit == 0, "detached apply = %d %s %s", result.exit, result.stdout, result.stderr)
	ref := intent.ResetRefPrefix(creation.Assignment.OwnerID, creation.Assignment.ID) + "1"
	parents := strings.Fields(gitOutput(t, root, "show", "-s", "--format=%P", ref))
	requireTest(t, strings.Contains(" "+strings.Join(parents, " ")+" ", " "+tip+" "), "root lost detached checkout's tip: %#v", parents)
}

func TestResetApplyKeepsTheRewoundTipReachable(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-rewind")
	commitInWorktree(t, creation.Path, "ahead", "ahead\n", "ahead")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	plan := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID}})
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	result := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID, "--apply", fingerprint}})
	requireTest(t, result.exit == 0, "rewind = %d %s %s", result.exit, result.stdout, result.stderr)
	manifest, ok := readRecoveryManifest(root, intent.ResetRefPrefix(creation.Assignment.OwnerID, creation.Assignment.ID)+"1")
	requireTest(t, ok && gitOutput(t, root, "rev-parse", creation.Assignment.Branch) == creation.Assignment.Start, "rewind failed or envelope absent")
	requireTest(t, gitOutput(t, root, "show", "-s", "--format=%P", manifest.Layers["working"]) == tip, "working payload lost previous tip")
}

func TestResetApplyRestoresTheExactLock(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-exact-lock")
	gitRun(t, root, "worktree", "unlock", creation.Path)
	gitRun(t, root, "worktree", "lock", "--reason", "foreign lock", creation.Path)
	plan := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID}})
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	result := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID, "--apply", fingerprint}})
	requireTest(t, result.exit == 0 && strings.Contains(result.stdout, "preserved=none"), "lock repair = %d %s %s", result.exit, result.stdout, result.stderr)
	_, err := resolveAssignment(root, creation.Assignment.ID)
	mustNoError(t, err)
}

func TestResetApplyReattachesWithoutAnEnvelope(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-reattach")
	gitRun(t, creation.Path, "switch", "--detach", "HEAD")
	plan := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID}})
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	result := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID, "--apply", fingerprint}})
	requireTest(t, result.exit == 0 && strings.Contains(result.stdout, "preserved=none") && strings.Contains(result.stdout, "restore=none"), "reattach = %d %s %s", result.exit, result.stdout, result.stderr)
	requireTest(t, gitOutput(t, creation.Path, "symbolic-ref", "HEAD") == creation.Assignment.Branch &&
		gitOutput(t, root, "for-each-ref", "--format=%(refname)", intent.ResetRefPrefix(creation.Assignment.OwnerID, creation.Assignment.ID)) == "", "pure repair wrote an envelope or stayed detached")
}
