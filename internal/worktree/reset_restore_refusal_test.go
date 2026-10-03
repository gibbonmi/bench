package worktree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

func TestResetRestoreRefusesAForeignRef(t *testing.T) {
	t.Parallel()
	f := restoreFixture(t)
	foreign := intent.ResetRefPrefix(f.creation.Assignment.OwnerID, "other") + "1"
	gitRun(t, f.root, "update-ref", foreign, f.ref)
	result := runVerb(t, verbReset, f.call("--restore", foreign, f.creation.Assignment.ID))
	requireTest(t, result.exit == 1 && strings.Contains(result.stdout, "restore ref is not this assignment's"), "foreign ref = %d %s %s", result.exit, result.stdout, result.stderr)
}

func TestResetRestoreRefusesAnUnparentedPayload(t *testing.T) {
	t.Parallel()
	f := restoreFixture(t)
	tree := gitOutput(t, f.root, "rev-parse", f.ref+"^{tree}")
	broken, err := commitTree(f.root, tree, nil, "missing parents\n")
	mustNoError(t, err)
	gitRun(t, f.root, "update-ref", f.ref, broken)
	result := runVerb(t, verbReset, f.call("--restore", f.ref, f.creation.Assignment.ID))
	requireTest(t, result.exit == 1 && strings.Contains(result.stdout, "reset envelope does not verify"), "unparented payload = %d %s %s", result.exit, result.stdout, result.stderr)
}

func TestResetRestoreRefusesAForeignTip(t *testing.T) {
	t.Parallel()
	f := restoreFixture(t)
	foreign, err := commitTree(f.root, gitOutput(t, f.root, "rev-parse", "HEAD^{tree}"), nil, "foreign root\n")
	mustNoError(t, err)
	rootOID, _, err := captureLayers(f.root, f.creation.Path, true, foreign)
	mustNoError(t, err)
	gitRun(t, f.root, "update-ref", f.ref, rootOID)
	result := runVerb(t, verbReset, f.call("--restore", f.ref, f.creation.Assignment.ID))
	requireTest(t, result.exit == 1 && strings.Contains(result.stdout, "envelope tip is not reached by the assignment start") &&
		strings.Contains(result.stdout, "wanted="+f.creation.Assignment.Start), "foreign tip = %d %s %s", result.exit, result.stdout, result.stderr)
}

func TestResetRestoreRefusesAStalePlan(t *testing.T) {
	t.Parallel()
	f := restoreFixture(t)
	plan := runVerb(t, verbReset, f.call("--restore", f.ref, f.creation.Assignment.ID))
	requireTest(t, isRestorePlanOf(plan, f.ref), "restore plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("later\n"), 0o644)
	plan = runVerb(t, verbReset, f.call("--restore", f.ref, f.creation.Assignment.ID))
	requireTest(t, isRestorePlanOf(plan, f.ref), "restore plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	wanted := plan.mustFingerprint(t)
	before := gitOutput(t, f.root, "for-each-ref", "--format=%(refname) %(objectname)", intent.ResetRefPrefix(f.creation.Assignment.OwnerID, f.creation.Assignment.ID))
	status := gitOutput(t, f.creation.Path, "status", "--porcelain=v1")
	result := runVerb(t, verbReset, f.call("--restore", f.ref, f.creation.Assignment.ID, "--apply", fingerprint))
	requireTest(t, result.exit == 1 && strings.Contains(result.stdout, "reset plan is stale") && strings.Contains(result.stdout, "wanted="+wanted) &&
		strings.Contains(result.stdout, "next=bench worktree reset --restore "+f.ref+" "+f.creation.Assignment.ID), "stale restore = %d %s %s", result.exit, result.stdout, result.stderr)
	requireTest(t, gitOutput(t, f.root, "for-each-ref", "--format=%(refname) %(objectname)", intent.ResetRefPrefix(f.creation.Assignment.OwnerID, f.creation.Assignment.ID)) == before &&
		gitOutput(t, f.creation.Path, "status", "--porcelain=v1") == status && gitOutput(t, f.creation.Path, "rev-parse", "HEAD") == f.creation.Assignment.Start, "stale restore changed the checkout or refs")
}

func TestResetRestoreRefusesARefOutsideTheNamespace(t *testing.T) {
	t.Parallel()
	f := restoreFixture(t)
	gitRun(t, f.root, "update-ref", "refs/bench/green/foreign", f.ref)
	result := runVerb(t, verbReset, f.call("--restore", "refs/bench/green/foreign", f.creation.Assignment.ID))
	requireTest(t, result.exit == 1 && strings.Contains(result.stdout, "restore ref is not this assignment's"), "outside namespace = %d %s %s", result.exit, result.stdout, result.stderr)
}

func TestResetRestoreRefusesAMissingRef(t *testing.T) {
	t.Parallel()
	f := restoreFixture(t)
	ref := intent.ResetRefPrefix(f.creation.Assignment.OwnerID, f.creation.Assignment.ID) + "999"
	result := runVerb(t, verbReset, f.call("--restore", ref, f.creation.Assignment.ID))
	requireTest(t, result.exit == 1 && strings.Contains(result.stdout, "reset envelope does not verify"), "missing ref = %d %s %s", result.exit, result.stdout, result.stderr)
}

func TestResetRestoreRefusesAMissingManifest(t *testing.T) {
	t.Parallel()
	f := restoreFixture(t)
	gitRun(t, f.root, "update-ref", f.ref, f.creation.Assignment.Start)
	result := runVerb(t, verbReset, f.call("--restore", f.ref, f.creation.Assignment.ID))
	requireTest(t, result.exit == 1 && strings.Contains(result.stdout, "reset envelope does not verify"), "missing manifest = %d %s %s", result.exit, result.stdout, result.stderr)
}

func TestResetRefusesBothModes(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	mustWrite(t, filepath.Join(root, ".git", intent.Filename), []byte("invalid ledger"), 0o644)
	result := runVerb(t, verbReset, verbCall{root: root, home: t.TempDir(), args: []string{"--to", "HEAD", "--restore", "refs/bench/reset/unknown", "unknown"}})
	requireTest(t, result.exit == 2 && result.stdout == "" && strings.Contains(result.stderr, "usage: bench worktree reset"), "both modes = %d %s %s", result.exit, result.stdout, result.stderr)
}

func TestResetRestoreRefusesAControlByteBeforeTheLedger(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	mustWrite(t, filepath.Join(root, ".git", intent.Filename), []byte("invalid ledger"), 0o644)
	result := runVerb(t, verbReset, verbCall{root: root, home: t.TempDir(), args: []string{"--restore", "refs/bench/reset/\x1b", "unknown"}})
	requireTest(t, result.exit == 1 && strings.Contains(result.stdout, "--restore contains control characters"), "control restore = %d %s %s", result.exit, result.stdout, result.stderr)
}

func TestResetRestoreRefusesAStaleIndex(t *testing.T) {
	t.Parallel()
	f := restoreFixture(t)
	mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("staged one\n"), 0o644)
	gitRun(t, f.creation.Path, "add", "README.md")
	mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("working\n"), 0o644)
	plan := runVerb(t, verbReset, f.call("--restore", f.ref, f.creation.Assignment.ID))
	requireTest(t, isRestorePlanOf(plan, f.ref), "restore plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("staged two\n"), 0o644)
	gitRun(t, f.creation.Path, "add", "README.md")
	mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("working\n"), 0o644)
	result := runVerb(t, verbReset, f.call("--restore", f.ref, f.creation.Assignment.ID, "--apply", fingerprint))
	requireTest(t, result.exit == 1 && strings.Contains(result.stdout, "reset plan is stale") && gitOutput(t, f.creation.Path, "show", ":README.md") == "staged two",
		"stale staged restore = %d %s %s", result.exit, result.stdout, result.stderr)
}

func TestResetRestoreRefusesAnIgnoredCollision(t *testing.T) {
	t.Parallel()
	f := restoreFixture(t)
	commitInWorktree(t, f.creation.Path, ".gitignore", "untracked\n", "ignore the untracked path")
	head := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	mustWrite(t, filepath.Join(f.creation.Path, "untracked"), []byte("ignored now\n"), 0o644)
	result := runVerb(t, verbReset, f.call("--restore", f.ref, f.creation.Assignment.ID))
	requireTest(t, result.exit == 1 && strings.Contains(result.stdout, "refused{detail=ignored content would be overwritten}") &&
		strings.Contains(result.stdout, refusalPathsTable+"[1]{path}:\n  untracked\n"), "collision restore = %d %s %s", result.exit, result.stdout, result.stderr)
	body, err := os.ReadFile(filepath.Join(f.creation.Path, "untracked"))
	mustNoError(t, err)
	requireTest(t, string(body) == "ignored now\n" && gitOutput(t, f.creation.Path, "rev-parse", "HEAD") == head &&
		gitOutput(t, f.root, "rev-parse", f.creation.Assignment.Branch) == head, "collision refusal changed the ignored bytes or moved the checkout")
}

func TestResetRestoreRefusesACollisionWithTheEnvelopeTip(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "restore-tip-collision")
	commitInWorktree(t, f.creation.Path, "gen", "generated\n", "track gen on the tip")
	gitRun(t, f.creation.Path, "switch", "--detach", f.creation.Assignment.Start)
	mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("dirty\n"), 0o644)
	plan := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	result := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID, "--apply", plan.mustFingerprint(t)))
	requireTest(t, result.exit == 0, "fixture reset = %d %s %s", result.exit, result.stdout, result.stderr)
	ref := intent.ResetRefPrefix(f.creation.Assignment.OwnerID, f.creation.Assignment.ID) + "1"
	commitInWorktree(t, f.creation.Path, ".gitignore", "gen\n", "ignore gen")
	mustWrite(t, filepath.Join(f.creation.Path, "gen"), []byte("ignored now\n"), 0o644)
	result = runVerb(t, verbReset, f.call("--restore", ref, f.creation.Assignment.ID))
	requireTest(t, result.exit == 1 && strings.Contains(result.stdout, "ignored content would be overwritten") && strings.Contains(result.stdout, refusalPathsTable+"[1]{path}:\n  gen\n"),
		"tip collision = %d %s %s", result.exit, result.stdout, result.stderr)
}
