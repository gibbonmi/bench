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
	f := newOwnedAssignment(t, "reset-plan")
	mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("dirty\n"), 0o644)
	before, err := git.Raw("-C", f.creation.Path, "status", "--porcelain=v1", "-z", "--untracked-files=all")
	mustNoError(t, err)
	result := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, result.exit == 0, "plan = %d %s %s", result.exit, result.stdout, result.stderr)
	for _, cell := range []string{"reset_plan{", "worktree=" + f.creation.Assignment.ID, "mode=reset", "action=reset",
		"checkpoint=" + f.creation.Assignment.Start, "head=" + f.creation.Assignment.Start, "ref=" + f.creation.Assignment.Branch,
		"tip=" + f.creation.Assignment.Start, "tracked=dirty", "lock=ok", "preserve=envelope", "fingerprint="} {
		requireTest(t, strings.Contains(result.stdout, cell), "plan lacks %s: %s", cell, result.stdout)
	}
	after, err := git.Raw("-C", f.creation.Path, "status", "--porcelain=v1", "-z", "--untracked-files=all")
	mustNoError(t, err)
	requireTest(t, string(after) == string(before), "plan changed checkout status")
	requireTest(t, strings.Contains(result.stdout, "next=bench worktree reset --to "), "plan lacks the apply command: %s", result.stdout)
}

func TestResetPlanNamesTheApplyCommand(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-next")
	mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("dirty\n"), 0o644)
	plan := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	result := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	want := "next=bench worktree reset --to " + f.creation.Assignment.Start + " " + f.creation.Assignment.ID + " --apply " + fingerprint
	requireTest(t, result.exit == 0 && strings.Contains(result.stdout, want), "apply command = %d %s %s; want %s", result.exit, result.stdout, result.stderr, want)
}

func TestResetPlanReportsNothingToReset(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-clean")
	result := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, result.exit == 0 && strings.Contains(result.stdout, "action=none") && strings.Contains(result.stdout, "fingerprint=none"),
		"clean plan = %d %s %s", result.exit, result.stdout, result.stderr)
	requireTest(t, strings.Contains(result.stdout, "tracked=clean") && strings.Contains(result.stdout, "preserve=none") &&
		!strings.Contains(result.stdout, "next=") && !strings.Contains(result.stdout, "reset_paths"), "idle plan = %s", result.stdout)
}

func TestResetPlanListsEveryAffectedPath(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-paths")
	commitInWorktree(t, f.creation.Path, ".gitignore", "ignored\n", "ignore output")
	mustWrite(t, filepath.Join(f.creation.Path, "staged"), []byte("staged\n"), 0o644)
	gitRun(t, f.creation.Path, "add", "staged")
	mustWrite(t, filepath.Join(f.creation.Path, "tracked.txt"), []byte("unstaged\n"), 0o644)
	mustWrite(t, filepath.Join(f.creation.Path, "untracked"), []byte("new\n"), 0o644)
	mustWrite(t, filepath.Join(f.creation.Path, "ignored"), []byte("output\n"), 0o644)
	gitRun(t, f.creation.Path, "mv", "README.md", "renamed")
	result := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, result.exit == 0 && strings.Contains(result.stdout, resetPathsTable+"[4]{path,status}:"), "paths = %d %s %s", result.exit, result.stdout, result.stderr)
	for _, path := range []string{"staged", "tracked.txt", "untracked", "renamed"} {
		requireTest(t, strings.Contains(result.stdout, "\n  "+path+","), "missing path %s: %s", path, result.stdout)
	}
	requireTest(t, !strings.Contains(result.stdout, "README.md") && !strings.Contains(result.stdout, "ignored"), "unexpected path row: %s", result.stdout)
}

func TestResetPlanSanitizesAHostilePath(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-control-path")
	mustWrite(t, filepath.Join(f.creation.Path, "bad\x1bname"), []byte("new\n"), 0o644)
	result := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, result.exit == 0 && strings.Contains(result.stdout, resetPathsTable+"[1]") && !strings.ContainsRune(result.stdout, '\x1b'),
		"hostile path = %d %s %s", result.exit, result.stdout, result.stderr)
	requireTest(t, strings.Count(result.stdout, "reset_plan{") == 1 && len(strings.Split(strings.TrimSpace(result.stdout), "\n")) == 3,
		"hostile path split the plan: %q", result.stdout)
}

func TestResetPlanReachesAShiftBranchCheckout(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-shift")
	gitRun(t, f.creation.Path, "switch", "-c", "bench/shift-reset-plan")
	commitInWorktree(t, f.creation.Path, "shift", "work\n", "shift work")
	result := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, result.exit == 0 && strings.Contains(result.stdout, "ref=refs/heads/bench/shift-reset-plan") &&
		strings.Contains(result.stdout, "action=reset"), "shift plan = %d %s %s", result.exit, result.stdout, result.stderr)
}

func TestResetPlanReachesADriftedLock(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-lock")
	gitRun(t, f.root, "worktree", "unlock", f.creation.Path)
	gitRun(t, f.root, "worktree", "lock", "--reason", "shift retained", f.creation.Path)
	result := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, result.exit == 0 && strings.Contains(result.stdout, "lock=repair") && strings.Contains(result.stdout, "action=reset") &&
		strings.Contains(result.stdout, "preserve=none"), "lock plan = %d %s %s", result.exit, result.stdout, result.stderr)
}

func TestResetResolvesTheCheckpointInTheTarget(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-symbolic")
	commitInWorktree(t, f.creation.Path, "ahead", "ahead\n", "ahead")
	head := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	requireTest(t, gitOutput(t, f.root, "rev-parse", "HEAD") != head, "fixture root and target share a head")
	result := runVerb(t, verbReset, f.call("--to", "HEAD", f.creation.Assignment.ID))
	requireTest(t, result.exit == 0 && strings.Contains(result.stdout, "checkpoint="+head) && strings.Contains(result.stdout, "action=none"),
		"symbolic checkpoint = %d %s %s", result.exit, result.stdout, result.stderr)
}

func TestResetPlanOpensNoSpan(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-no-span")
	result := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, result.exit == 0, "plan = %d %s %s", result.exit, result.stdout, result.stderr)
	_, err := os.Stat(otelrecord.Path(f.home, f.root))
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
