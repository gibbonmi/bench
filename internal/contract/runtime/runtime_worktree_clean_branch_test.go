package runtime

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestRuntimeWorktreeCleanBranchContracts(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	contract.RunParallel(t, "bench worktree clean --discard-branch plan-then-apply contract", testRuntimeWorktreeCleanDiscardBranch)
	contract.RunParallel(t, "bench worktree clean --discard-branch refusal contract", testRuntimeWorktreeCleanDiscardBranchRefusals)
}

// squashLandedCheckout returns one owned assignment whose two commits were composed into a
// single commit on main, with the branch ref it left behind. Ancestry, merge detection, and
// patch-equivalence all miss that shape, so `bench worktree clean` classifies the branch as
// unmerged while the work is already in — the residue [RW1] removes.
func squashLandedCheckout(t *testing.T, f contract.Fixture, env map[string]string, request string) (string, string) {
	t.Helper()
	created := f.BenchEnv(env, "worktree", "create", "--request", request, "--label", request)
	created.RequireExit(0)
	path := worktreeCreatePath(t, created.Stdout)
	branch := strings.TrimSpace(contract.RunAt(t, f, path, nil, "git", "symbolic-ref", "HEAD").Stdout)
	start := strings.TrimSpace(contract.RunAt(t, f, path, nil, "git", "rev-parse", "HEAD").Stdout)
	for _, name := range []string{"one.txt", "two.txt"} {
		contract.WriteFileAbs(t, filepath.Join(path, name), name+"\n")
		contract.RunAt(t, f, path, nil, "git", "add", name).RequireExit(0)
		contract.RunAt(t, f, path, nil, "git", "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", name).RequireExit(0)
	}
	f.Git("cherry-pick", "--no-commit", start+".."+strings.TrimPrefix(branch, "refs/heads/"))
	f.Git("-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "squashed")
	return path, branch
}

// [RW1] [RW2]
func testRuntimeWorktreeCleanDiscardBranch(t *testing.T) {
	f := onMainFixture(t)
	env := map[string]string{"BENCH_HOME": filepath.Join(f.Root, ".bench-home")}
	path, branch := squashLandedCheckout(t, f, env, "discard-branch")
	short := strings.TrimPrefix(branch, "refs/heads/")

	unproven := f.BenchEnv(env, "worktree", "clean", path)
	unproven.RequireExit(0)
	contract.RequireNotContains(t, unproven.Stdout, "discards branch")

	plan := f.BenchEnv(env, "worktree", "clean", "--discard-branch", path)
	plan.RequireExit(0)
	contract.RequireContains(t, plan.Stdout, "discards branch "+branch)
	contract.RequireContains(t, plan.Stdout, ",remove,clean,")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("planning removed the checkout: %v", err)
	}
	if !headExists(f, short) {
		t.Fatal("planning deleted the assignment branch")
	}

	apply := f.BenchEnv(env, "worktree", "clean", "--discard-branch", path, "--apply", cleanupFingerprint(t, plan.Stdout))
	apply.RequireExit(0)
	contract.RequireContains(t, apply.Stdout, ",removed,")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("apply left the checkout behind: %v", err)
	}
	if headExists(f, short) {
		t.Fatalf("apply left assignment branch %s behind", branch)
	}
	if refs := strings.TrimSpace(f.Git("for-each-ref", "--format=%(refname)", "refs/bench/recovery/").Stdout); refs != "" {
		t.Fatalf("apply left recovery refs behind: %q", refs)
	}
}

// [RW4] The override authorizes discarding a landed-by-the-operator's-word payload and
// nothing else: a foreign registration and the primary checkout refuse for exactly the
// reasons they refuse without it, and neither loses its branch.
func testRuntimeWorktreeCleanDiscardBranchRefusals(t *testing.T) {
	f := onMainFixture(t)
	env := map[string]string{"BENCH_HOME": filepath.Join(f.Root, ".bench-home")}
	foreign := filepath.Join(t.TempDir(), "foreign locked")
	f.Git("worktree", "add", "-q", "-b", "foreign-locked", foreign, "HEAD")
	f.Git("worktree", "lock", "--reason", "foreign", foreign)

	locked := f.BenchEnv(env, "worktree", "clean", "--discard-branch", foreign)
	locked.RequireExit(0)
	contract.RequireContains(t, locked.Stdout, ",retain,")
	contract.RequireContains(t, locked.Stdout, "foreign or unexpected lock is retained")
	contract.RequireNotContains(t, locked.Stdout, "discards branch")
	f.BenchEnv(env, "worktree", "clean", "--discard-branch", foreign, "--apply", cleanupFingerprint(t, locked.Stdout)).RequireExit(0)
	if _, err := os.Stat(foreign); err != nil {
		t.Fatalf("apply removed a foreign worktree: %v", err)
	}
	if !headExists(f, "foreign-locked") {
		t.Fatal("apply deleted a foreign branch")
	}

	primary := f.BenchEnv(env, "worktree", "clean", "--discard-branch", f.Root)
	primary.RequireExit(0)
	contract.RequireContains(t, primary.Stdout, "primary checkout is never removable")
	f.BenchEnv(env, "worktree", "clean", "--discard-branch", f.Root, "--apply", cleanupFingerprint(t, primary.Stdout)).RequireExit(0)
	if !headExists(f, "main") {
		t.Fatal("apply deleted the default branch")
	}
}
