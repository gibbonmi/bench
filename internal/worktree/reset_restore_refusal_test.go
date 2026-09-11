package worktree

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

func TestResetRestoreRefusesAForeignRef(t *testing.T) {
	t.Parallel()
	root, creation, home, ref := restoreFixture(t)
	foreign := intent.ResetRefPrefix(creation.Assignment.OwnerID, "other") + "1"
	gitRun(t, root, "update-ref", foreign, ref)
	code, out, errout := runReset(t, root, home, "--restore", foreign, creation.Assignment.ID)
	requireTest(t, code == 1 && strings.Contains(out, "restore ref is not this assignment's"), "foreign ref = %d %s %s", code, out, errout)
}

func TestResetRestoreRefusesAnUnparentedPayload(t *testing.T) {
	t.Parallel()
	root, creation, home, ref := restoreFixture(t)
	tree := gitOutput(t, root, "rev-parse", ref+"^{tree}")
	broken, err := commitTree(root, tree, nil, "missing parents\n")
	mustNoError(t, err)
	gitRun(t, root, "update-ref", ref, broken)
	code, out, errout := runReset(t, root, home, "--restore", ref, creation.Assignment.ID)
	requireTest(t, code == 1 && strings.Contains(out, "reset envelope does not verify"), "unparented payload = %d %s %s", code, out, errout)
}

func TestResetRestoreRefusesAForeignTip(t *testing.T) {
	t.Parallel()
	root, creation, home, ref := restoreFixture(t)
	foreign, err := commitTree(root, gitOutput(t, root, "rev-parse", "HEAD^{tree}"), nil, "foreign root\n")
	mustNoError(t, err)
	rootOID, _, _, err := captureLayers(root, creation.Path, true, foreign)
	mustNoError(t, err)
	gitRun(t, root, "update-ref", ref, rootOID)
	code, out, errout := runReset(t, root, home, "--restore", ref, creation.Assignment.ID)
	requireTest(t, code == 1 && strings.Contains(out, "envelope tip is not reached by the assignment start") &&
		strings.Contains(out, "wanted="+creation.Assignment.Start), "foreign tip = %d %s %s", code, out, errout)
}

func TestResetRestoreRefusesAStalePlan(t *testing.T) {
	t.Parallel()
	root, creation, home, ref := restoreFixture(t)
	fingerprint := restoreFingerprint(t, root, home, ref, creation.Assignment.ID)
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("later\n"), 0o644)
	wanted := restoreFingerprint(t, root, home, ref, creation.Assignment.ID)
	before := gitOutput(t, root, "for-each-ref", "--format=%(refname) %(objectname)", intent.ResetRefPrefix(creation.Assignment.OwnerID, creation.Assignment.ID))
	status := gitOutput(t, creation.Path, "status", "--porcelain=v1")
	code, out, errout := runReset(t, root, home, "--restore", ref, creation.Assignment.ID, "--apply", fingerprint)
	requireTest(t, code == 1 && strings.Contains(out, "reset plan is stale") && strings.Contains(out, "wanted="+wanted) &&
		strings.Contains(out, "next=bench worktree reset --restore "+ref+" "+creation.Assignment.ID), "stale restore = %d %s %s", code, out, errout)
	requireTest(t, gitOutput(t, root, "for-each-ref", "--format=%(refname) %(objectname)", intent.ResetRefPrefix(creation.Assignment.OwnerID, creation.Assignment.ID)) == before &&
		gitOutput(t, creation.Path, "status", "--porcelain=v1") == status && gitOutput(t, creation.Path, "rev-parse", "HEAD") == creation.Assignment.Start, "stale restore changed the checkout or refs")
}

func TestResetRestoreRefusesARefOutsideTheNamespace(t *testing.T) {
	t.Parallel()
	root, creation, home, ref := restoreFixture(t)
	gitRun(t, root, "update-ref", "refs/bench/green/foreign", ref)
	code, out, errout := runReset(t, root, home, "--restore", "refs/bench/green/foreign", creation.Assignment.ID)
	requireTest(t, code == 1 && strings.Contains(out, "restore ref is not this assignment's"), "outside namespace = %d %s %s", code, out, errout)
}

func TestResetRestoreRefusesAMissingRef(t *testing.T) {
	t.Parallel()
	root, creation, home, _ := restoreFixture(t)
	ref := intent.ResetRefPrefix(creation.Assignment.OwnerID, creation.Assignment.ID) + "999"
	code, out, errout := runReset(t, root, home, "--restore", ref, creation.Assignment.ID)
	requireTest(t, code == 1 && strings.Contains(out, "reset envelope does not verify"), "missing ref = %d %s %s", code, out, errout)
}

func TestResetRestoreRefusesAMissingManifest(t *testing.T) {
	t.Parallel()
	root, creation, home, ref := restoreFixture(t)
	gitRun(t, root, "update-ref", ref, creation.Assignment.Start)
	code, out, errout := runReset(t, root, home, "--restore", ref, creation.Assignment.ID)
	requireTest(t, code == 1 && strings.Contains(out, "reset envelope does not verify"), "missing manifest = %d %s %s", code, out, errout)
}

func TestResetRefusesBothModes(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	mustWrite(t, filepath.Join(root, ".git", intent.Filename), []byte("invalid ledger"), 0o644)
	code, out, errout := runReset(t, root, t.TempDir(), "--to", "HEAD", "--restore", "refs/bench/reset/unknown", "unknown")
	requireTest(t, code == 2 && out == "" && strings.Contains(errout, "usage: bench worktree reset"), "both modes = %d %s %s", code, out, errout)
}

func TestResetRestoreRefusesAControlByteBeforeTheLedger(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	mustWrite(t, filepath.Join(root, ".git", intent.Filename), []byte("invalid ledger"), 0o644)
	code, out, errout := runReset(t, root, t.TempDir(), "--restore", "refs/bench/reset/\x1b", "unknown")
	requireTest(t, code == 1 && strings.Contains(out, "--restore contains control characters"), "control restore = %d %s %s", code, out, errout)
}
