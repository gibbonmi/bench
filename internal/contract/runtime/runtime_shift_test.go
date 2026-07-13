package runtime

import (
	"github.com/gibbonmi/bench/internal/contract"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
}

func testShiftGatedLoop(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	home := t.TempDir()
	beforeBranch := strings.TrimSpace(f.Git("branch", "--show-current").Stdout)
	beforeStatus := f.Git("status", "--porcelain").Stdout

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": "true", "BENCH_MAX_ITERS": "1", "BENCH_HOME": home}, "shift", "noop")

	probe.RequireExit(0)
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

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "1", "BENCH_HOME": home}, "shift", "isolate")

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
	probe.RequireContains(probe.Stdout, "preserving iteration 1")
	wt := shiftWorktree(t, probe.Stdout)
	for _, path := range []string{"step 1 [a].txt", "gate-artifact.txt"} {
		if _, err := os.Stat(filepath.Join(wt, path)); err != nil {
			t.Fatalf("failed gate did not preserve %s: %v", path, err)
		}
	}
	requireRegisteredWorktree(t, f, wt)
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
	wt := shiftWorktree(t, probe.Stdout)
	requireEqual(t, strings.TrimSpace(f.Git("branch", "--show-current").Stdout), beforeBranch, "red shift changed the main checkout branch")
	requireEqual(t, strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout), beforeHead, "red shift moved main checkout HEAD")
	requireEqual(t, f.Git("status", "--porcelain").Stdout, beforeStatus, "red shift dirtied the main checkout")
	if _, err := os.Stat(filepath.Join(wt, "red.txt")); err != nil {
		t.Fatalf("red shift did not preserve work: %v", err)
	}
	requireRegisteredWorktree(t, f, wt)
}

func testShiftCommitFailure(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable("agent", "#!/usr/bin/env bash\nprintf 'work\\n' > work.txt\n")
	f.CommitAll("agent")
	home := t.TempDir()
	cleanHome := t.TempDir()
	cleanXDG := t.TempDir()

	probe := f.BenchEnv(map[string]string{
		"HOME": cleanHome, "XDG_CONFIG_HOME": cleanXDG, "GIT_CONFIG_NOSYSTEM": "1",
		"GIT_AUTHOR_NAME": "", "GIT_AUTHOR_EMAIL": "", "GIT_COMMITTER_NAME": "", "GIT_COMMITTER_EMAIL": "",
		"BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "1", "BENCH_HOME": home,
	}, "shift", "no-author")

	if probe.ExitCode == 0 {
		t.Fatalf("shift with no git author config succeeded\nstdout:\n%s\nstderr:\n%s", probe.Stdout, probe.Stderr)
	}
	probe.RequireContains(probe.Stderr, "could not commit iteration 1")
	probe.RequireNotContains(probe.Stdout, "1 committed iteration(s)")
	requireNoLease(t, home)
}

func testShiftTouchedScopeStructure(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteFile("preexisting.py", strings.Repeat("x = \n", 401))
	f.CommitAll("preexisting debt")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": "true", "BENCH_MAX_ITERS": "1", "BENCH_REFACTOR_ITERS": "1", "BENCH_HOME": home}, "shift", "noop")

	probe.RequireExit(0)
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

	probe.RequireExit(0)
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

	if probe.ExitCode == 0 {
		t.Fatalf("interrupted shift exited successfully\nstdout:\n%s\nstderr:\n%s", probe.Stdout, probe.Stderr)
	}
	if f.Exists(".bench-objective") {
		t.Fatal("interrupted shift left .bench-objective")
	}
	if f.Exists(".bench-notes.md") {
		t.Fatal("interrupted shift left .bench-notes.md")
	}
	requireNoLease(t, home)
	waitSeconds(t, 1)
	f.BenchEnv(map[string]string{"BENCH_AGENT": "true", "BENCH_MAX_ITERS": "1", "BENCH_HOME": home}, "shift", "after-interrupt").RequireExit(0)
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

	if probe.ExitCode == 0 {
		t.Fatalf("gate-interrupted shift exited successfully\nstdout:\n%s\nstderr:\n%s", probe.Stdout, probe.Stderr)
	}
	wt := shiftWorktree(t, probe.Stdout)
	requireNoLease(t, home)
	waitSeconds(t, 3)
	if _, err := os.Stat(filepath.Join(wt, "late-gate-write.txt")); err == nil {
		t.Fatal("gate child kept running after lease release and dirtied the pooled worktree")
	}
	if dirty := runGitAt(t, wt, "status", "--porcelain"); dirty != "" {
		t.Fatalf("gate-interrupted pooled worktree is dirty:\n%s", dirty)
	}
	follow := f.BenchEnv(map[string]string{"BENCH_TEST_STATE": state, "BENCH_AGENT": "true", "BENCH_MAX_ITERS": "1", "BENCH_HOME": home}, "shift", "after-gate-interrupt")
	follow.RequireExit(0)
	follow.RequireContains(follow.Stdout, "shift done")
	follow.RequireContains(follow.Stdout, "worktree: "+wt)
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
	wt := shiftWorktree(t, probe.Stdout)
	for _, path := range []string{".bench-notes.md", ".bench-objective", "junk.txt"} {
		if _, err := os.Stat(filepath.Join(wt, path)); err != nil {
			t.Fatalf("red gate did not preserve %s: %v", path, err)
		}
	}
	requireRegisteredWorktree(t, f, wt)
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

	f.BenchEnv(map[string]string{"BENCH_TEST_PROMPTS": prompts, "BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "1", "BENCH_REFACTOR_ITERS": "1", "BENCH_HOME": home}, "shift", "make-big").RequireExit(0)

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
