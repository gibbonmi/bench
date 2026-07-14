package runtime

import (
	"encoding/json"
	"github.com/gibbonmi/bench/internal/contract"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRuntimeShiftLoopContracts(t *testing.T) {
	t.Parallel()
	contract.RunParallel(t, "bench shift gated-loop contract", testShiftGatedLoop)
	contract.RunParallel(t, "bench shift worktree-isolation contract", testShiftWorktreeIsolation)
	contract.RunParallel(t, "bench shift stage-touched contract", testShiftStageTouched)
	contract.RunParallel(t, "bench shift red-rollback isolation contract", testShiftRedRollback)
	contract.RunParallel(t, "bench shift commit-failure contract", testShiftCommitFailure)
	contract.RunParallel(t, "bench shift touched-scope structure contract", testShiftTouchedScopeStructure)
	contract.RunParallel(t, "bench shift refactor no-op contract", testShiftRefactorNoop)
	contract.RunParallel(t, "bench shift interrupt cleanup contract", testShiftInterruptCleanup)
	contract.RunParallel(t, "bench shift gate-interrupt cleanup contract", testShiftGateInterruptCleanup)
	contract.RunParallel(t, "bench shift done.sh early-completion contract", testShiftDoneEarlyCompletion)
	contract.RunParallel(t, "bench shift scratch-survival contract", testShiftScratchSurvival)
	contract.RunParallel(t, "bench shift refactor-prompt scope contract", testShiftRefactorPromptScope)
	contract.RunParallel(t, "bench shift no-objective usage contract", testShiftNoObjective)
	contract.RunParallel(t, "bench shift whitespace-only objective usage contract", testShiftWhitespaceObjective)
	contract.RunParallel(t, "bench shift control-byte objective usage contract", testShiftControlByteObjective)
	contract.RunParallel(t, "bench shift invalid env usage contract", testShiftInvalidEnvUsage)
	contract.RunParallel(t, "bench shift empty env default contract", testShiftEmptyEnvDefault)
	contract.RunParallel(t, "bench shift complete result contract", testShiftCompleteResult)
	contract.RunParallel(t, "bench shift adapter-failure zero-commit contract", testShiftAdapterFailureZeroCommit)
	contract.RunParallel(t, "bench shift adapter spawn-failure contract", testShiftAdapterSpawnFailure)
	contract.RunParallel(t, "bench shift adapter-failure after-commit contract", testShiftAdapterFailureAfterCommit)
	contract.RunParallel(t, "bench shift cap-exhaustion contract", testShiftCapExhaustion)
	contract.RunParallel(t, "bench shift no-op done contract", testShiftNoOpDone)
	contract.RunParallel(t, "bench shift distinct recovery per re-run contract", testShiftDistinctRecoveryPerRun)
	contract.RunParallel(t, "bench shift locked fallback worktree survives resume contract", testShiftLockedFallbackWorktreeSurvivesResume)
	contract.RunParallel(t, "bench shift wall deadline contract", testShiftWallDeadlineKillsAdapter)
	contract.RunParallel(t, "bench shift wall deadline mid-gate contract", testShiftWallDeadlineCancelsRunningGate)
	contract.RunParallel(t, "bench shift intent and status recovery surfacing contract", testShiftIntentAndStatusSurfaceRecovery)
	contract.RunParallel(t, "bench shift interrupt-with-mutation recovery contract", testShiftInterruptWithMutationRecovers)
}

func testShiftGatedLoop(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	home := t.TempDir()
	beforeBranch := strings.TrimSpace(f.Git("branch", "--show-current").Stdout)
	beforeStatus := f.Git("status", "--porcelain").Stdout

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": "true", "BENCH_MAX_ITERS": "1", "BENCH_HOME": home}, "shift", "noop")

	// A no-op adapter commits nothing: honest taxonomy is no-op/4, not complete/0 —
	// the loop still runs its normal "shift done" summary and cleanup.
	probe.RequireExit(4)
	probe.RequireContains(probe.Stdout, "shift done")
	branch := shiftBranch(t, probe.Stdout)
	requireEqual(t, strings.TrimSpace(f.Git("branch", "--show-current").Stdout), beforeBranch, "bench shift changed the main checkout branch")
	requireEqual(t, f.Git("status", "--porcelain").Stdout, beforeStatus, "bench shift dirtied the main checkout")
	f.Git("rev-parse", "--verify", branch)
	requireNoWorktreeBranch(t, f, branch)
	requireNoLease(t, home)
}

func testShiftWorktreeIsolation(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable("agent", "#!/usr/bin/env bash\nprintf 'shifted\\n' > shifted.txt\n")
	f.CommitAll("agent")
	home := t.TempDir()
	beforeBranch := strings.TrimSpace(f.Git("branch", "--show-current").Stdout)
	beforeHead := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
	beforeStatus := f.Git("status", "--porcelain").Stdout

	// Two iterations of headroom so the loop stops clean (no-change, complete) rather
	// than exhausting the cap (incomplete) — this test's intent is isolation of a
	// completed shift, not the cap boundary.
	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "2", "BENCH_HOME": home}, "shift", "isolate")

	probe.RequireExit(0)
	probe.RequireContains(probe.Stdout, "1 committed iteration(s)")
	branch := shiftBranch(t, probe.Stdout)
	requireEqual(t, strings.TrimSpace(f.Git("branch", "--show-current").Stdout), beforeBranch, "shift changed the main checkout branch")
	requireEqual(t, strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout), beforeHead, "shift moved main checkout HEAD")
	requireEqual(t, f.Git("status", "--porcelain").Stdout, beforeStatus, "shift dirtied the main checkout")
	f.Git("cat-file", "-e", branch+":shifted.txt")
	requireEqual(t, strings.TrimSpace(f.Git("config", "branch."+branch+".benchBase").Stdout), beforeHead, "shift did not record branch.<name>.benchBase")
	requireEqual(t, strings.TrimSpace(f.Git("rev-list", "--count", beforeHead+".."+branch).Stdout), "1", "shift branch has the wrong commit count")
	requireNoWorktreeBranch(t, f, branch)
	requireNoLease(t, home)
}

func testShiftStageTouched(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nprintf cache > gate-artifact.txt\nexit 0\n")
	f.WriteExecutable("agent", `#!/usr/bin/env bash
n=0; [ -f count ] && n="$(cat count)"
n=$((n+1)); printf '%s\n' "$n" > count
printf 'work\n' > "step $n [a].txt"
`)
	f.CommitAll("agent")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "2", "BENCH_HOME": home}, "shift", "stage-touched")

	probe.RequireExit(1)
	probe.RequireContains(probe.Stdout, "snapshotting iteration 1")
	// The released pool worktree must never be advertised as the preservation site.
	probe.RequireNotContains(probe.Stdout, "preserved in")
	branch := shiftBranchFromStart(t, probe.Stdout)
	ref := "refs/bench/recovery/" + branch
	f.Git("show-ref", "--verify", ref)
	tree := f.Git("ls-tree", "-r", "--name-only", ref).Stdout
	for _, path := range []string{"step 1 [a].txt", "gate-artifact.txt"} {
		if !strings.Contains(tree, path) {
			t.Fatalf("recovery snapshot did not preserve %s:\n%s", path, tree)
		}
	}
	requireShiftResult(t, probe.Stdout, "failed,\"1\","+branch+",\"0\",\"1\",\"ref:"+ref+"\",")
	probe.RequireContains(probe.Stdout, "recovery: ref:"+ref)
	requireNoLease(t, home)
}

func testShiftRedRollback(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 1\n")
	f.WriteExecutable("agent", "#!/usr/bin/env bash\nprintf 'red\\n' > red.txt\n")
	f.CommitAll("agent")
	home := t.TempDir()
	beforeBranch := strings.TrimSpace(f.Git("branch", "--show-current").Stdout)
	beforeHead := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
	beforeStatus := f.Git("status", "--porcelain").Stdout

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "1", "BENCH_HOME": home}, "shift", "red-rollback")

	probe.RequireExit(1)
	probe.RequireContains(probe.Stdout, "gate failed")
	requireEqual(t, strings.TrimSpace(f.Git("branch", "--show-current").Stdout), beforeBranch, "red shift changed the main checkout branch")
	requireEqual(t, strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout), beforeHead, "red shift moved main checkout HEAD")
	requireEqual(t, f.Git("status", "--porcelain").Stdout, beforeStatus, "red shift dirtied the main checkout")
	// Row 1: a red-gate iteration snapshots the dirty tree to refs/bench/recovery/<branch>
	// and releases the pool worktree — replacing the old preserve-worktree special case.
	branch := shiftBranchFromStart(t, probe.Stdout)
	ref := "refs/bench/recovery/" + branch
	f.Git("show-ref", "--verify", ref)
	tree := f.Git("ls-tree", "-r", "--name-only", ref).Stdout
	if !strings.Contains(tree, "red.txt") {
		t.Fatalf("recovery snapshot did not preserve red.txt:\n%s", tree)
	}
	// Row 1,5: the shift_result recovery cell and the printed location both name the ref.
	requireShiftResult(t, probe.Stdout, "failed,\"1\","+branch+",\"0\",\"1\",\"ref:"+ref+"\",")
	probe.RequireContains(probe.Stdout, "recovery: ref:"+ref)
	requireNoLease(t, home)
}

func testShiftCommitFailure(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable("agent", "#!/usr/bin/env bash\nprintf 'work\\n' > work.txt\n")
	f.CommitAll("agent")
	// Row 2: the snapshot is built with an explicit synthetic identity under plumbing, so
	// it survives exactly where an ordinary commit blocks — both a missing Git identity
	// and a failing commit hook. Written after the fixture's own setup commit, so it
	// blocks only the shift's iteration commit, not fixture setup.
	f.WriteExecutable(".git/hooks/pre-commit", "#!/usr/bin/env bash\nexit 1\n")
	home := t.TempDir()
	cleanHome := t.TempDir()
	cleanXDG := t.TempDir()

	probe := f.BenchEnv(map[string]string{
		"HOME": cleanHome, "XDG_CONFIG_HOME": cleanXDG, "GIT_CONFIG_NOSYSTEM": "1",
		"GIT_AUTHOR_NAME": "", "GIT_AUTHOR_EMAIL": "", "GIT_COMMITTER_NAME": "", "GIT_COMMITTER_EMAIL": "",
		"BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "1", "BENCH_HOME": home,
	}, "shift", "no-author")

	// Zero commits landed before the failure: honest taxonomy is failed/1.
	probe.RequireExit(1)
	probe.RequireContains(probe.Stderr, "could not commit iteration 1")
	probe.RequireNotContains(probe.Stdout, "1 committed iteration(s)")
	branch := shiftBranchFromStart(t, probe.Stdout)
	ref := "refs/bench/recovery/" + branch
	f.Git("show-ref", "--verify", ref)
	tree := f.Git("ls-tree", "-r", "--name-only", ref).Stdout
	if !strings.Contains(tree, "work.txt") {
		t.Fatalf("recovery snapshot did not survive the identity/hook failure that blocked the commit:\n%s", tree)
	}
	requireNoLease(t, home)
}

func testShiftTouchedScopeStructure(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteFile("preexisting.py", strings.Repeat("x = \n", 401))
	f.CommitAll("preexisting debt")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": "true", "BENCH_MAX_ITERS": "1", "BENCH_REFACTOR_ITERS": "1", "BENCH_HOME": home}, "shift", "noop")

	// A no-op adapter commits nothing: honest taxonomy is no-op/4, not complete/0.
	probe.RequireExit(4)
	probe.RequireContains(probe.Stdout, "shift done")
	probe.RequireNotContains(probe.Stdout, "refactor phase")
}

func testShiftRefactorNoop(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable("agent", `#!/usr/bin/env bash
if [ ! -f made-big ]; then
  seq 401 | sed 's/^/x = /' > touched.py
  : > made-big
fi
`)
	f.CommitAll("agent")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "1", "BENCH_REFACTOR_ITERS": "3", "BENCH_HOME": home}, "shift", "make-big")

	// One committed iteration then the cap is exhausted: honest taxonomy is
	// incomplete/3, not complete/0. This test's intent is the refactor no-op output,
	// which is unaffected.
	probe.RequireExit(3)
	probe.RequireContains(probe.Stdout, "refactor phase")
	probe.RequireContains(probe.Stdout, "refactor 1 made no staged change")
	probe.RequireNotContains(probe.Stdout, "refactor 2/")
	probe.RequireNotContains(probe.Stdout, "refactor 1 committed")
	probe.RequireNotContains(probe.Stdout, "/improve-codebase-architecture")
	branch := shiftBranch(t, probe.Stdout)
	requireEqual(t, strings.TrimSpace(f.Git("rev-list", "--count", "HEAD.."+branch).Stdout), "1", "no-op refactor created an unexpected commit")
}

func testShiftInterruptCleanup(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable("interrupt-agent", "#!/usr/bin/env bash\nkill -INT \"$PPID\"\nexit 130\n")
	f.CommitAll("agent")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": filepath.Join(f.Root, "interrupt-agent"), "BENCH_MAX_ITERS": "2", "BENCH_HOME": home}, "shift", "interrupt")

	probe.RequireExit(130)
	probe.RequireContains(probe.Stdout, "interrupted")
	if f.Exists(".bench-objective") {
		t.Fatal("interrupted shift left .bench-objective")
	}
	if f.Exists(".bench-notes.md") {
		t.Fatal("interrupted shift left .bench-notes.md")
	}
	requireNoLease(t, home)
	waitSeconds(t, 1)
	// The follow-up shift is a proof the pool/lease recovered after the interrupt, not
	// a claim about this shift's own outcome — a no-op adapter commits nothing, so the
	// honest taxonomy is no-op/4.
	f.BenchEnv(map[string]string{"BENCH_AGENT": "true", "BENCH_MAX_ITERS": "1", "BENCH_HOME": home}, "shift", "after-interrupt").RequireExit(4)
}

func testShiftGateInterruptCleanup(t *testing.T) {
	f := shiftFixture(t, `#!/usr/bin/env bash
trap 'exit 130' INT
if [ ! -f "$BENCH_TEST_STATE/gate-interrupted-once" ]; then
  : > "$BENCH_TEST_STATE/gate-interrupted-once"
  kill -INT "$PPID"
  sleep 2
  printf 'late\n' > late-gate-write.txt
fi
exit 0
`)
	home := t.TempDir()
	state := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_TEST_STATE": state, "BENCH_AGENT": "true", "BENCH_MAX_ITERS": "2", "BENCH_HOME": home}, "shift", "gate-interrupt")

	probe.RequireExit(130)
	wt := shiftWorktree(t, probe.Stdout)
	waitSeconds(t, 3)
	if _, err := os.Stat(filepath.Join(wt, "late-gate-write.txt")); err == nil {
		t.Fatal("gate child kept running after cancellation and dirtied the pooled worktree")
	}
	// Nothing beyond scratch is dirty here (BENCH_AGENT=true never mutates), so the
	// uniform rule releases the worktree normally instead of retaining it — replacing the
	// old assertion that the scratch survived in a permanently-preserved worktree.
	requireShiftResult(t, probe.Stdout, "interrupted,\"130\","+shiftBranchFromStart(t, probe.Stdout)+",\"0\",\"1\",none,")
	requireNoLease(t, home)
}

func testShiftDoneEarlyCompletion(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable(".bench/done.sh", "#!/usr/bin/env bash\n[ -f step1.txt ]\n")
	f.WriteExecutable("agent", `#!/usr/bin/env bash
n=0; [ -f count ] && n="$(cat count)"
n=$((n + 1))
printf '%s\n' "$n" > count
printf '%s\n' "$n" > "step$n.txt"
`)
	f.CommitAll("agent")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "3", "BENCH_HOME": home}, "shift", "done-early")

	probe.RequireExit(0)
	probe.RequireContains(probe.Stdout, "objective met.")
	probe.RequireNotContains(probe.Stdout, "iteration 2/3")
	probe.RequireContains(probe.Stdout, "1 committed iteration(s)")
	branch := shiftBranch(t, probe.Stdout)
	f.Git("cat-file", "-e", branch+":step1.txt")
}

func testShiftScratchSurvival(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\n[ ! -f junk.txt ]\n")
	f.WriteExecutable("agent", `#!/usr/bin/env bash
n=0; [ -f "$BENCH_TEST_STATE/count" ] && n="$(cat "$BENCH_TEST_STATE/count")"
n=$((n+1)); printf '%s\n' "$n" > "$BENCH_TEST_STATE/count"
if [ "$n" = 1 ]; then
  printf 'tried A, broke gate\n' >> .bench-notes.md
  printf 'junk\n' > junk.txt
else
  [ -f .bench-notes.md ] && grep -q 'tried A' .bench-notes.md && printf 'notes-survived\n' >> "$BENCH_TEST_STATE/report"
  [ -f .bench-objective ] && printf 'objective-survived\n' >> "$BENCH_TEST_STATE/report"
  printf 'ok\n' > done.txt
fi
`)
	f.CommitAll("agent")
	home := t.TempDir()
	state := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_TEST_STATE": state, "BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "2", "BENCH_HOME": home}, "shift", "survive")

	probe.RequireExit(1)
	// The scratch-survives-a-red-iteration mechanic (a later iteration inside one run
	// still reads what a prior one learned) is rollback's and is unchanged by FT79 — it
	// never fires in this single-red-iteration fixture, since a red main-loop gate exits
	// the shift rather than retrying. What changes here is where the final red exit's
	// evidence lands: a recovery snapshot of the agent's mutation, scratch excluded by
	// design, not the physical (now-released) worktree.
	branch := shiftBranchFromStart(t, probe.Stdout)
	ref := "refs/bench/recovery/" + branch
	f.Git("show-ref", "--verify", ref)
	tree := f.Git("ls-tree", "-r", "--name-only", ref).Stdout
	if !strings.Contains(tree, "junk.txt") {
		t.Fatalf("recovery snapshot did not preserve junk.txt:\n%s", tree)
	}
	for _, path := range []string{".bench-notes.md", ".bench-objective"} {
		if strings.Contains(tree, path) {
			t.Fatalf("recovery snapshot rode scratch %s into the tree:\n%s", path, tree)
		}
	}
	requireNoLease(t, home)
}

func testShiftRefactorPromptScope(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable("agent", `#!/usr/bin/env bash
printf '%s\n@@@@\n' "$1" >> "$BENCH_TEST_PROMPTS"
if [ ! -f made-big ]; then seq 401 | sed 's/^/x = /' > touched.py; : > made-big; fi
`)
	f.CommitAll("agent")
	home := t.TempDir()
	prompts := filepath.Join(t.TempDir(), "prompts.txt")

	// One committed iteration then the cap is exhausted: honest taxonomy is
	// incomplete/3, not complete/0. This test's intent is the refactor prompt scope,
	// which is unaffected.
	f.BenchEnv(map[string]string{"BENCH_TEST_PROMPTS": prompts, "BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "1", "BENCH_REFACTOR_ITERS": "1", "BENCH_HOME": home}, "shift", "make-big").RequireExit(3)

	data, err := os.ReadFile(prompts)
	if err != nil {
		t.Fatalf("read prompts: %v", err)
	}
	parts := strings.Split(string(data), "@@@@\n")
	if len(parts) < 2 {
		t.Fatalf("adapter prompts missing delimiter:\n%s", data)
	}
	refactor := parts[len(parts)-2]
	if !strings.Contains(refactor, "touched.py") {
		t.Fatal("refactor prompt does not name the flagged touched files")
	}
	if strings.Contains(refactor, "Run `bench structure` to see the flagged files") {
		t.Fatal("refactor prompt still points at repo-wide structure output")
	}
}

// testShiftDistinctRecoveryPerRun covers row 3,18: two preserving shifts in one fixture
// get distinct branches and distinct recovery refs, and the first ref still resolves to
// its original commit after the second shift runs. Branch names are per-second, so the
// two runs wait apart.
func testShiftDistinctRecoveryPerRun(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 1\n")
	f.WriteExecutable("agent", "#!/usr/bin/env bash\nprintf 'red\\n' > red.txt\n")
	f.CommitAll("agent")
	home := t.TempDir()
	env := map[string]string{"BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "1", "BENCH_HOME": home}

	first := f.BenchEnv(env, "shift", "first red")
	first.RequireExit(1)
	firstBranch := shiftBranchFromStart(t, first.Stdout)
	firstRef := "refs/bench/recovery/" + firstBranch
	firstOID := strings.TrimSpace(f.Git("rev-parse", "--verify", firstRef).Stdout)

	waitSeconds(t, 1)

	second := f.BenchEnv(env, "shift", "second red")
	second.RequireExit(1)
	secondBranch := shiftBranchFromStart(t, second.Stdout)
	secondRef := "refs/bench/recovery/" + secondBranch

	if firstBranch == secondBranch {
		t.Fatalf("two shifts in one fixture reused the same branch: %s", firstBranch)
	}
	if firstRef == secondRef {
		t.Fatalf("two shifts in one fixture reused the same recovery ref: %s", firstRef)
	}
	f.Git("show-ref", "--verify", secondRef)
	if got := strings.TrimSpace(f.Git("rev-parse", "--verify", firstRef).Stdout); got != firstOID {
		t.Fatalf("the second shift disturbed the first shift's recovery ref: got %s want %s", got, firstOID)
	}
}

// testShiftLockedFallbackWorktreeSurvivesResume covers row 4,18's resume half: a locked
// worktree with no owner marker — exactly what the retain-and-lock fallback leaves
// behind when a snapshot itself fails — classifies as ReasonUnexpectedLock and is
// retained, never removed, by bench resume-clean.
func testShiftLockedFallbackWorktreeSurvivesResume(t *testing.T) {
	f := onMainFixture(t)
	wt := filepath.Join(t.TempDir(), "fallback-worktree")
	f.Git("worktree", "add", "-q", "--detach", wt)
	f.Git("worktree", "lock", "--reason", "bench shift recovery: snapshot failed", wt)

	out := f.Bench("resume-clean")
	out.RequireExit(0)

	porcelain := f.Git("worktree", "list", "--porcelain").Stdout
	if !strings.Contains(porcelain, "worktree "+wt) {
		t.Fatalf("resume-clean removed the locked fallback worktree:\n%s", porcelain)
	}
	if !strings.Contains(porcelain, "locked") {
		t.Fatalf("resume-clean unlocked the fallback worktree:\n%s", porcelain)
	}
}

// testShiftWallDeadlineKillsAdapter covers row 8: a sleeping adapter under a tiny
// BENCH_MAX_WALL is killed, its prior mutation is snapshotted, the shift exits
// incomplete/3 with a deadline detail, and the adapter does not keep running past exit.
func testShiftWallDeadlineKillsAdapter(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable("agent", "#!/usr/bin/env bash\nprintf 'partial\\n' > partial.txt\nsleep 3\nprintf 'late\\n' > late.txt\n")
	f.CommitAll("agent")
	home := t.TempDir()

	start := time.Now()
	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "2", "BENCH_MAX_WALL": "1s", "BENCH_HOME": home}, "shift", "wall deadline")
	elapsed := time.Since(start)

	probe.RequireExit(3)
	if elapsed > 15*time.Second {
		t.Fatalf("shift did not honor the wall deadline promptly: %s elapsed", elapsed)
	}
	probe.RequireContains(probe.Stdout, "deadline")
	wt := shiftWorktree(t, probe.Stdout)
	branch := shiftBranchFromStart(t, probe.Stdout)
	waitSeconds(t, 4)
	if _, err := os.Stat(filepath.Join(wt, "late.txt")); err == nil {
		t.Fatal("adapter kept running after the wall deadline and dirtied the pooled worktree")
	}
	ref := "refs/bench/recovery/" + branch
	f.Git("show-ref", "--verify", ref)
	requireNoLease(t, home)
}

// testShiftWallDeadlineCancelsRunningGate covers row 8's other half: a wall deadline
// that fires while the gate itself is running cancels the gate instead of waiting it
// out — row 8's sleeping-adapter case only exercises the wall against a hung adapter,
// so a wall that never wires cancel-gate would pass that row yet hang here.
func testShiftWallDeadlineCancelsRunningGate(t *testing.T) {
	f := shiftFixture(t, `#!/usr/bin/env bash
if [ ! -f "$BENCH_TEST_STATE/gate-ran-once" ]; then
  : > "$BENCH_TEST_STATE/gate-ran-once"
  sleep 5
  printf 'late\n' > "$BENCH_TEST_STATE/late-gate-finished"
fi
exit 0
`)
	f.WriteExecutable("agent", "#!/usr/bin/env bash\nprintf 'work\\n' > work.txt\n")
	f.CommitAll("agent")
	home := t.TempDir()
	state := t.TempDir()

	start := time.Now()
	probe := f.BenchEnv(map[string]string{
		"BENCH_TEST_STATE": state, "BENCH_AGENT": filepath.Join(f.Root, "agent"),
		"BENCH_MAX_ITERS": "2", "BENCH_MAX_WALL": "1s", "BENCH_HOME": home,
	}, "shift", "wall mid gate")
	elapsed := time.Since(start)

	probe.RequireExit(3)
	if elapsed > 15*time.Second {
		t.Fatalf("shift did not honor the wall deadline while the gate was running: %s elapsed", elapsed)
	}
	probe.RequireContains(probe.Stdout, "deadline")
	branch := shiftBranchFromStart(t, probe.Stdout)
	waitSeconds(t, 5)
	if _, err := os.Stat(filepath.Join(state, "late-gate-finished")); err == nil {
		t.Fatal("gate kept running past the wall deadline instead of being cancelled")
	}
	ref := "refs/bench/recovery/" + branch
	f.Git("show-ref", "--verify", ref)
	tree := f.Git("ls-tree", "-r", "--name-only", ref).Stdout
	if !strings.Contains(tree, "work.txt") {
		t.Fatalf("recovery snapshot did not preserve the adapter's mutation before the gate-cancelling deadline:\n%s", tree)
	}
	requireNoLease(t, home)
}

// testShiftIntentAndStatusSurfaceRecovery covers row 17: a preserving shift's outcome
// and recovery pointer land on its intent-ledger entry, and bench status renders the
// pointer, so a preserved shift is discoverable after the terminal is gone.
func testShiftIntentAndStatusSurfaceRecovery(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 1\n")
	f.WriteExecutable("agent", "#!/usr/bin/env bash\nprintf 'work\\n' > work.txt\n")
	f.CommitAll("agent")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "1", "BENCH_HOME": home}, "shift", "surface recovery")
	probe.RequireExit(1)
	branch := shiftBranchFromStart(t, probe.Stdout)
	pointer := "ref:refs/bench/recovery/" + branch

	ledgerPath := filepath.Join(gitDir(t, f), "bench-intent.json")
	data, err := os.ReadFile(ledgerPath)
	if err != nil {
		t.Fatalf("read intent ledger: %v", err)
	}
	var ledger struct {
		Entries []struct{ Objective, Outcome, Recovery string }
	}
	if err := json.Unmarshal(data, &ledger); err != nil {
		t.Fatalf("decode intent ledger: %v", err)
	}
	found := false
	for _, e := range ledger.Entries {
		if e.Objective != "surface recovery" {
			continue
		}
		found = true
		if e.Outcome != "failed" || e.Recovery != pointer {
			t.Fatalf("intent entry = %+v, want outcome=failed recovery=%s", e, pointer)
		}
	}
	if !found {
		t.Fatal("intent ledger is missing the shift entry")
	}

	status := f.BenchEnv(map[string]string{"BENCH_HOME": home}, "status")
	status.RequireExit(0)
	status.RequireContains(status.Stdout, pointer)
}

// testShiftInterruptWithMutationRecovers extends the interrupt contracts: an interrupt
// that lands after the adapter has already mutated the tree still snapshots the
// mutation (scratch excluded) before releasing the pool worktree.
func testShiftInterruptWithMutationRecovers(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable("interrupt-agent", "#!/usr/bin/env bash\nprintf 'work\\n' > work.txt\nkill -INT \"$PPID\"\nexit 130\n")
	f.CommitAll("agent")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": filepath.Join(f.Root, "interrupt-agent"), "BENCH_MAX_ITERS": "1", "BENCH_HOME": home}, "shift", "interrupt with mutation")

	probe.RequireExit(130)
	branch := shiftBranchFromStart(t, probe.Stdout)
	ref := "refs/bench/recovery/" + branch
	f.Git("show-ref", "--verify", ref)
	tree := f.Git("ls-tree", "-r", "--name-only", ref).Stdout
	if !strings.Contains(tree, "work.txt") {
		t.Fatalf("interrupted shift did not snapshot the adapter's mutation:\n%s", tree)
	}
	for _, path := range []string{".bench-notes.md", ".bench-objective"} {
		if strings.Contains(tree, path) {
			t.Fatalf("interrupted shift's snapshot rode scratch %s into the tree:\n%s", path, tree)
		}
	}
	requireNoLease(t, home)
}

// requireNoShiftBranch asserts no branch matching bench/shift-* exists — the signal
// that a usage failure exited before the loop created anything.
func requireNoShiftBranch(t *testing.T, f contract.Fixture) {
	t.Helper()
	if out := f.Git("for-each-ref", "refs/heads/bench/shift-*").Stdout; strings.TrimSpace(out) != "" {
		t.Fatalf("usage failure created a shift branch:\n%s", out)
	}
}

// requireShiftResult asserts the shift_result TOON block's row equals wantRow exactly —
// header plus the one data row, in the pinned field order.
func requireShiftResult(t *testing.T, stdout, wantRow string) {
	t.Helper()
	header := "shift_result[1]{outcome,exit,branch,committed,iterations_used,recovery,detail}:"
	if !strings.Contains(stdout, header) {
		t.Fatalf("shift_result block missing or malformed header:\n%s", stdout)
	}
	if !strings.Contains(stdout, wantRow) {
		t.Fatalf("shift_result row = %q not found in:\n%s", wantRow, stdout)
	}
}

// testShiftNoObjective covers row 9: an empty objective exits 2 before acquiring
// anything — the "improve the codebase" default is gone.
func testShiftNoObjective(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": "true", "BENCH_HOME": home}, "shift")

	probe.RequireExit(2)
	probe.RequireContains(probe.Stdout, "usage,\"2\",none,\"0\",\"0\",none,")
	requireNoLease(t, home)
	requireNoShiftBranch(t, f)
}

// testShiftWhitespaceObjective covers row 9's other half: a whitespace-only objective
// exits 2 before acquiring anything, same as an empty one — a bare `== ""` guard would
// let this slip through as a legitimate objective.
func testShiftWhitespaceObjective(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": "true", "BENCH_HOME": home}, "shift", "   ")

	probe.RequireExit(2)
	probe.RequireContains(probe.Stdout, "usage,\"2\",none,\"0\",\"0\",none,")
	requireNoLease(t, home)
	requireNoShiftBranch(t, f)
}

// testShiftControlByteObjective covers row 10: an objective carrying a control byte
// (ESC, 0x1b) exits 2 at entry, before it can reach the ledger or the TOON emitter.
func testShiftControlByteObjective(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": "true", "BENCH_HOME": home}, "shift", "bad\x1bobjective")

	probe.RequireExit(2)
	probe.RequireContains(probe.Stdout, "usage,\"2\",none,\"0\",\"0\",none,")
	requireNoLease(t, home)
	requireNoShiftBranch(t, f)
}

// testShiftInvalidEnvUsage covers row 11: BENCH_MAX_ITERS, BENCH_MAX_WALL set to
// values outside the accepted range each exit 2, naming the variable and range, before
// any lease is acquired.
func testShiftInvalidEnvUsage(t *testing.T) {
	cases := []struct {
		name string
		env  map[string]string
		want string // substring stderr must contain: the variable and its range
	}{
		{"zero iters", map[string]string{"BENCH_MAX_ITERS": "0"}, "BENCH_MAX_ITERS"},
		{"non-integer iters", map[string]string{"BENCH_MAX_ITERS": "abc"}, "BENCH_MAX_ITERS"},
		{"wall over bound", map[string]string{"BENCH_MAX_WALL": "48h"}, "BENCH_MAX_WALL"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
			home := t.TempDir()
			env := map[string]string{"BENCH_AGENT": "true", "BENCH_HOME": home}
			for k, v := range tc.env {
				env[k] = v
			}

			probe := f.BenchEnv(env, "shift", "objective")

			probe.RequireExit(2)
			probe.RequireContains(probe.Stderr, tc.want)
			requireNoLease(t, home)
			requireNoShiftBranch(t, f)
		})
	}
}

// testShiftEmptyEnvDefault covers row 11's other half: BENCH_MAX_ITERS explicitly set
// to the empty string behaves exactly like unset — the default cap applies and no
// usage error fires.
func testShiftEmptyEnvDefault(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	home := t.TempDir()

	// No adapter mutation ever lands, so the honest end state is no-op/4 either way —
	// the point of this row is that "" behaves as unset (no exit 2), not that the
	// default numeric cap is visible in the output.
	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": "true", "BENCH_MAX_ITERS": "", "BENCH_HOME": home}, "shift", "objective")

	probe.RequireExit(4)
}

// testShiftCompleteResult covers rows 6 and 16: a green gate, a committing adapter, and
// a .bench/done.sh that passes once the work has landed stop the loop clean — exit 0,
// and the shift_result block carries all seven fields with outcome complete, committed
// 1, recovery none.
func testShiftCompleteResult(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable(".bench/done.sh", "#!/usr/bin/env bash\n[ -f landed.txt ]\n")
	f.WriteExecutable("agent", "#!/usr/bin/env bash\nprintf 'work\\n' > landed.txt\n")
	f.CommitAll("agent")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "3", "BENCH_HOME": home}, "shift", "land it")

	probe.RequireExit(0)
	branch := shiftBranch(t, probe.Stdout)
	requireShiftResult(t, probe.Stdout, "complete,\"0\","+branch+",\"1\",\"1\",none,")
}

// testShiftAdapterFailureZeroCommit covers rows 6 and 7: BENCH_AGENT=/bin/false with a
// green gate stops the loop as failed/1 with zero commits — today this exits 0
// "objective likely met"; the whole point of this row is that it no longer does.
func testShiftAdapterFailureZeroCommit(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": "/bin/false", "BENCH_MAX_ITERS": "2", "BENCH_HOME": home}, "shift", "objective")

	probe.RequireExit(1)
	branch := shiftBranchFromStart(t, probe.Stdout)
	requireShiftResult(t, probe.Stdout, "failed,\"1\","+branch+",\"0\",")
}

// testShiftAdapterSpawnFailure covers row 7's other half: an adapter that passes the
// preflight (it is a regular executable file) but fails to exec at iteration time — here
// because its shebang names an interpreter that does not exist — stops the loop as
// failed/1, not no-op/4. runAdapter used to swallow cmd.Start's error, so a spawn
// failure read as a clean "nothing happened" no-op; this proves it now propagates as
// real evidence of a broken adapter.
func testShiftAdapterSpawnFailure(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable("agent", "#!/no/such/interpreter\necho hi\n")
	f.CommitAll("agent")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "1", "BENCH_HOME": home}, "shift", "objective")

	probe.RequireExit(1)
	branch := shiftBranchFromStart(t, probe.Stdout)
	requireShiftResult(t, probe.Stdout, "failed,\"1\","+branch+",\"0\",")
	requireNoLease(t, home)
}

// testShiftAdapterFailureAfterCommit covers row 7: an adapter that writes a file then
// exits nonzero on every run commits its first iteration's mutation (gate green) and
// then stops the loop as incomplete/3 — folding this into failed/1 would understate the
// committed work.
func testShiftAdapterFailureAfterCommit(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable("agent", `#!/usr/bin/env bash
n=0; [ -f count ] && n="$(cat count)"
n=$((n+1)); printf '%s\n' "$n" > count
printf 'work %s\n' "$n" > "work-$n.txt"
exit 1
`)
	f.CommitAll("agent")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "3", "BENCH_HOME": home}, "shift", "objective")

	probe.RequireExit(3)
	branch := shiftBranchFromStart(t, probe.Stdout)
	requireShiftResult(t, probe.Stdout, "incomplete,\"3\","+branch+",\"1\",")
	requireEqual(t, strings.TrimSpace(f.Git("rev-list", "--count", "HEAD.."+branch).Stdout), "1", "adapter-failure-after-commit shift committed the wrong number of iterations")
}

// testShiftCapExhaustion covers row 12: an always-dirty adapter runs to the iteration
// cap, leaving the committed work on the branch and exiting incomplete/3 — today cap
// exhaustion exits 0, reading as done.
func testShiftCapExhaustion(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable("agent", `#!/usr/bin/env bash
n=0; [ -f count ] && n="$(cat count)"
n=$((n+1)); printf '%s\n' "$n" > count
printf 'work %s\n' "$n" > "work-$n.txt"
`)
	f.CommitAll("agent")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "2", "BENCH_HOME": home}, "shift", "objective")

	probe.RequireExit(3)
	branch := shiftBranch(t, probe.Stdout)
	requireShiftResult(t, probe.Stdout, "incomplete,\"3\","+branch+",\"2\",\"2\",none,")
	requireEqual(t, strings.TrimSpace(f.Git("rev-list", "--count", "HEAD.."+branch).Stdout), "2", "cap-exhaustion shift committed the wrong number of iterations")
}

// testShiftNoOpDone covers row 13: a no-op adapter with a .bench/done.sh that always
// passes commits zero iterations and exits no-op/4, with a detail naming the predicate
// — complete-on-predicate would exit 0 and hide that nothing was done.
func testShiftNoOpDone(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable(".bench/done.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.CommitAll("done.sh")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": "true", "BENCH_MAX_ITERS": "2", "BENCH_HOME": home}, "shift", "objective")

	probe.RequireExit(4)
	probe.RequireContains(probe.Stdout, "predicate")
}
