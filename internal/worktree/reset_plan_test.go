package worktree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/otelrecord"
)

func TestResetPlanReportsTheDirtyCheckoutAndWritesNothing(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-plan")
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("dirty\n"), 0o644)
	before, err := git.Raw("-C", creation.Path, "status", "--porcelain=v1", "-z", "--untracked-files=all")
	mustNoError(t, err)
	code, out, errout := runReset(t, root, home, "--to", creation.Assignment.Start, creation.Assignment.ID)
	requireTest(t, code == 0, "plan = %d %s %s", code, out, errout)
	for _, cell := range []string{"reset_plan{", "worktree=" + creation.Assignment.ID, "mode=reset", "action=reset",
		"checkpoint=" + creation.Assignment.Start, "head=" + creation.Assignment.Start, "ref=" + creation.Assignment.Branch,
		"tip=" + creation.Assignment.Start, "tracked=dirty", "lock=ok", "preserve=envelope", "fingerprint="} {
		requireTest(t, strings.Contains(out, cell), "plan lacks %s: %s", cell, out)
	}
	after, err := git.Raw("-C", creation.Path, "status", "--porcelain=v1", "-z", "--untracked-files=all")
	mustNoError(t, err)
	requireTest(t, string(after) == string(before), "plan changed checkout status")
	requireTest(t, strings.Contains(out, "next=bench worktree reset --to "), "plan lacks the apply command: %s", out)
}

func TestResetPlanNamesTheApplyCommand(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-next")
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("dirty\n"), 0o644)
	fingerprint := resetFingerprint(t, root, home, creation.Assignment.Start, creation.Assignment.ID)
	code, out, errout := runReset(t, root, home, "--to", creation.Assignment.Start, creation.Assignment.ID)
	want := "next=bench worktree reset --to " + creation.Assignment.Start + " " + creation.Assignment.ID + " --apply " + fingerprint
	requireTest(t, code == 0 && strings.Contains(out, want), "apply command = %d %s %s; want %s", code, out, errout, want)
}

func TestResetPlanReportsNothingToReset(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-clean")
	code, out, errout := runReset(t, root, home, "--to", creation.Assignment.Start, creation.Assignment.ID)
	requireTest(t, code == 0 && strings.Contains(out, "action=none") && strings.Contains(out, "fingerprint=none"),
		"clean plan = %d %s %s", code, out, errout)
	requireTest(t, strings.Contains(out, "tracked=clean") && strings.Contains(out, "preserve=none") &&
		!strings.Contains(out, "next=") && !strings.Contains(out, "reset_paths"), "idle plan = %s", out)
}

func TestResetPlanListsEveryAffectedPath(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-paths")
	commitInWorktree(t, creation.Path, ".gitignore", "ignored\n", "ignore output")
	mustWrite(t, filepath.Join(creation.Path, "staged"), []byte("staged\n"), 0o644)
	gitRun(t, creation.Path, "add", "staged")
	mustWrite(t, filepath.Join(creation.Path, "tracked.txt"), []byte("unstaged\n"), 0o644)
	mustWrite(t, filepath.Join(creation.Path, "untracked"), []byte("new\n"), 0o644)
	mustWrite(t, filepath.Join(creation.Path, "ignored"), []byte("output\n"), 0o644)
	gitRun(t, creation.Path, "mv", "README.md", "renamed")
	code, out, errout := runReset(t, root, home, "--to", creation.Assignment.Start, creation.Assignment.ID)
	requireTest(t, code == 0 && strings.Contains(out, "reset_paths[4]{path,status}:"), "paths = %d %s %s", code, out, errout)
	for _, path := range []string{"staged", "tracked.txt", "untracked", "renamed"} {
		requireTest(t, strings.Contains(out, "\n  "+path+","), "missing path %s: %s", path, out)
	}
	requireTest(t, !strings.Contains(out, "README.md") && !strings.Contains(out, "ignored"), "unexpected path row: %s", out)
}

func TestResetPlanSanitizesAHostilePath(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-control-path")
	mustWrite(t, filepath.Join(creation.Path, "bad\x1bname"), []byte("new\n"), 0o644)
	code, out, errout := runReset(t, root, home, "--to", creation.Assignment.Start, creation.Assignment.ID)
	requireTest(t, code == 0 && strings.Contains(out, "reset_paths[1]") && !strings.ContainsRune(out, '\x1b'),
		"hostile path = %d %s %s", code, out, errout)
	requireTest(t, strings.Count(out, "reset_plan{") == 1 && len(strings.Split(strings.TrimSpace(out), "\n")) == 3,
		"hostile path split the plan: %q", out)
}

func TestResetPlanReachesAShiftBranchCheckout(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-shift")
	gitRun(t, creation.Path, "switch", "-c", "bench/shift-reset-plan")
	commitInWorktree(t, creation.Path, "shift", "work\n", "shift work")
	code, out, errout := runReset(t, root, home, "--to", creation.Assignment.Start, creation.Assignment.ID)
	requireTest(t, code == 0 && strings.Contains(out, "ref=refs/heads/bench/shift-reset-plan") &&
		strings.Contains(out, "action=reset"), "shift plan = %d %s %s", code, out, errout)
}

func TestResetPlanReachesADriftedLock(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-lock")
	gitRun(t, root, "worktree", "unlock", creation.Path)
	gitRun(t, root, "worktree", "lock", "--reason", "shift retained", creation.Path)
	code, out, errout := runReset(t, root, home, "--to", creation.Assignment.Start, creation.Assignment.ID)
	requireTest(t, code == 0 && strings.Contains(out, "lock=repair") && strings.Contains(out, "action=reset") &&
		strings.Contains(out, "preserve=none"), "lock plan = %d %s %s", code, out, errout)
}

func TestResetResolvesTheCheckpointInTheTarget(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-symbolic")
	commitInWorktree(t, creation.Path, "ahead", "ahead\n", "ahead")
	head := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	requireTest(t, gitOutput(t, root, "rev-parse", "HEAD") != head, "fixture root and target share a head")
	code, out, errout := runReset(t, root, home, "--to", "HEAD", creation.Assignment.ID)
	requireTest(t, code == 0 && strings.Contains(out, "checkpoint="+head) && strings.Contains(out, "action=none"),
		"symbolic checkpoint = %d %s %s", code, out, errout)
}

func TestResetPlanOpensNoSpan(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-no-span")
	code, out, errout := runReset(t, root, home, "--to", creation.Assignment.Start, creation.Assignment.ID)
	requireTest(t, code == 0, "plan = %d %s %s", code, out, errout)
	_, err := os.Stat(otelrecord.Path(home, root))
	requireTest(t, os.IsNotExist(err), "plan created a span record: %v", err)
}

func TestResetSeamIsRegistered(t *testing.T) {
	t.Parallel()
	for _, entry := range otelrecord.Registry {
		if entry.Seam == "worktree.reset" && entry.Package == "internal/worktree" && entry.Function == "beginVerbSpan" {
			return
		}
	}
	t.Fatal("reset seam is absent from the registry")
}
