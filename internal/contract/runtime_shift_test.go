package contract

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestRuntimeShiftContracts(t *testing.T) {
	t.Run("bench shift gated-loop contract", testShiftGatedLoop)
	t.Run("bench shift worktree-isolation contract", testShiftWorktreeIsolation)
	t.Run("bench shift stage-touched contract", testShiftStageTouched)
	t.Run("bench shift red-rollback isolation contract", testShiftRedRollback)
	t.Run("bench shift commit-failure contract", testShiftCommitFailure)
	t.Run("bench shift touched-scope structure contract", testShiftTouchedScopeStructure)
	t.Run("bench shift refactor no-op contract", testShiftRefactorNoop)
	t.Run("bench shift interrupt cleanup contract", testShiftInterruptCleanup)
	t.Run("bench shift gate-interrupt cleanup contract", testShiftGateInterruptCleanup)
	t.Run("bench shift done.sh early-completion contract", testShiftDoneEarlyCompletion)
	t.Run("bench shift scratch-survival contract", testShiftScratchSurvival)
	t.Run("bench shift refactor-prompt scope contract", testShiftRefactorPromptScope)
	t.Run("bench shift adapter preflight contract", testShiftAdapterPreflight)
	t.Run("bench shift adapter single-argument contract", testShiftAdapterSingleArgument)
	t.Run("reference adapter files contract", testReferenceAdapterFiles)
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

	probe.RequireExit(0)
	probe.RequireContains(probe.Stdout, "2 committed iteration(s)")
	branch := shiftBranch(t, probe.Stdout)
	f.Git("cat-file", "-e", branch+":step 1 [a].txt")
	f.Git("cat-file", "-e", branch+":step 2 [a].txt")
	if f.GitAllow("cat-file", "-e", branch+":gate-artifact.txt").ExitCode == 0 {
		t.Fatal("gate byproduct rode into an iteration commit")
	}
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

	probe.RequireExit(0)
	probe.RequireContains(probe.Stdout, "red gate")
	branch := shiftBranch(t, probe.Stdout)
	requireEqual(t, strings.TrimSpace(f.Git("branch", "--show-current").Stdout), beforeBranch, "red shift changed the main checkout branch")
	requireEqual(t, strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout), beforeHead, "red shift moved main checkout HEAD")
	requireEqual(t, f.Git("status", "--porcelain").Stdout, beforeStatus, "red shift dirtied the main checkout")
	if f.GitAllow("cat-file", "-e", branch+":red.txt").ExitCode == 0 {
		t.Fatal("red shift preserved rolled-back work")
	}
	requireEqual(t, strings.TrimSpace(f.Git("rev-list", "--count", beforeHead+".."+branch).Stdout), "0", "red shift branch gained a commit")
	requireNoLease(t, home)
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

	probe.RequireExit(0)
	report, err := os.ReadFile(filepath.Join(state, "report"))
	if err != nil {
		t.Fatalf("read scratch-survival report: %v", err)
	}
	if !strings.Contains(string(report), "notes-survived") {
		t.Fatal("red rollback wiped .bench-notes.md")
	}
	if !strings.Contains(string(report), "objective-survived") {
		t.Fatal("red rollback wiped .bench-objective")
	}
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

func testShiftAdapterPreflight(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	home := t.TempDir()

	unset := f.BenchEnvSpec(Env{"BENCH_AGENT": nil, "BENCH_HOME": strPtr(home)}, "shift", "probe")
	if unset.ExitCode == 0 {
		t.Fatal("shift with no BENCH_AGENT succeeded; should error")
	}
	unset.RequireContains(unset.Stderr, "BENCH_AGENT")
	if !regexp.MustCompile(`(?i)configure.*adapter|adapter.*configure`).MatchString(unset.Stderr) {
		t.Fatalf("unconfigured-adapter error is not a configure-your-adapter message:\n%s", unset.Stderr)
	}
	unset.RequireNotContains(unset.Stdout, "iteration 1/")

	empty := f.BenchEnv(map[string]string{"BENCH_AGENT": "", "BENCH_HOME": home}, "shift", "probe")
	if empty.ExitCode == 0 {
		t.Fatal("shift with empty BENCH_AGENT succeeded; should error")
	}
	empty.RequireContains(empty.Stderr, "BENCH_AGENT")
	if !regexp.MustCompile(`(?i)configure.*adapter|adapter.*configure`).MatchString(empty.Stderr) {
		t.Fatalf("empty-adapter error is not a configure-your-adapter message:\n%s", empty.Stderr)
	}
	empty.RequireNotContains(empty.Stdout, "iteration 1/")

	missing := f.BenchEnv(map[string]string{"BENCH_AGENT": "/no/such/adapter", "BENCH_HOME": home}, "shift", "probe")
	if missing.ExitCode == 0 {
		t.Fatal("shift with a missing adapter path succeeded; should error")
	}
	missing.RequireContains(missing.Stderr, "not executable")
	missing.RequireNotContains(missing.Stdout, "iteration 1/")

	keyword := f.BenchEnv(map[string]string{"BENCH_AGENT": "if", "BENCH_HOME": home}, "shift", "probe")
	if keyword.ExitCode == 0 {
		t.Fatal("shift with a shell-keyword adapter succeeded; should error")
	}
	keyword.RequireContains(keyword.Stderr, "not executable")
}

func testShiftAdapterSingleArgument(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable("adapter", `#!/usr/bin/env bash
{
  printf 'argc=%s\n' "$#"
  printf 'shift_env=%s\n' "${BENCH_SHIFT:-unset}"
  printf '%s\n@@@@\n' "$1"
} >> "$BENCH_TEST_RECORD"
`)
	f.CommitAll("adapter")
	home := t.TempDir()
	record := filepath.Join(t.TempDir(), "record.txt")

	f.BenchEnv(map[string]string{"BENCH_TEST_RECORD": record, "BENCH_AGENT": filepath.Join(f.Root, "adapter"), "BENCH_MAX_ITERS": "1", "BENCH_HOME": home}, "shift", "adapter-arg-probe").RequireExit(0)

	data, err := os.ReadFile(record)
	if err != nil {
		t.Fatalf("adapter was never invoked: %v", err)
	}
	text := string(data)
	for _, needle := range []string{
		"argc=1",
		"shift_env=1",
		"adapter-arg-probe",
		"You are one iteration of a Bench shift",
		"decides if it counts",
	} {
		if !strings.Contains(text, needle) {
			t.Fatalf("adapter record missing %q:\n%s", needle, text)
		}
	}
	if regexp.MustCompile(`(?m)^-p$`).MatchString(text) {
		t.Fatal("loop still passes the Claude-specific -p flag")
	}
}

func testReferenceAdapterFiles(t *testing.T) {
	root := KitRoot(t)
	for _, adapter := range []string{"claude", "codex", "opencode"} {
		path := filepath.Join(root, ".bench", "adapters", adapter)
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("reference adapter missing: .bench/adapters/%s", adapter)
		}
		if info.Mode()&0o111 == 0 {
			t.Fatalf("reference adapter not executable: .bench/adapters/%s", adapter)
		}
		probe := NewFixture(t, WithNoRepo()).Run("bash", "-n", path)
		probe.RequireExit(0)
		text, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read adapter %s: %v", adapter, err)
		}
		if !regexp.MustCompile(`(?m)^exec `).Match(text) {
			t.Fatalf("reference adapter %s does not exec its harness", adapter)
		}
		if !strings.Contains(string(text), `"$1"`) {
			t.Fatalf("reference adapter %s does not pass the prompt as $1", adapter)
		}
	}
	requireFileContains(t, filepath.Join(root, ".bench", "adapters", "claude"), `claude -p -- "$1"`, "claude adapter does not map the prompt to claude -p")
	requireFileContains(t, filepath.Join(root, ".bench", "adapters", "codex"), `codex exec -- "$1"`, "codex adapter does not map the prompt to codex exec")
	requireFileContains(t, filepath.Join(root, ".bench", "adapters", "opencode"), `opencode run -- "$1"`, "opencode adapter does not map the prompt to opencode run")
}

func shiftFixture(t *testing.T, gate string) Fixture {
	t.Helper()
	f := NewFixture(t)
	f.Env["GIT_AUTHOR_NAME"] = "Bench"
	f.Env["GIT_AUTHOR_EMAIL"] = "bench@local"
	f.Env["GIT_COMMITTER_NAME"] = "Bench"
	f.Env["GIT_COMMITTER_EMAIL"] = "bench@local"
	f.WriteExecutable(".bench/gate.sh", gate)
	f.CommitAll("init")
	return f
}

func shiftBranch(t *testing.T, output string) string {
	t.Helper()
	re := regexp.MustCompile(`(?m)^■ shift done: ([^,]+),`)
	m := re.FindStringSubmatch(output)
	if len(m) != 2 {
		t.Fatalf("shift summary did not name the branch:\n%s", output)
	}
	return m[1]
}

func shiftWorktree(t *testing.T, output string) string {
	t.Helper()
	re := regexp.MustCompile(`(?m)^  worktree: (.+)$`)
	m := re.FindStringSubmatch(output)
	if len(m) != 2 {
		t.Fatalf("shift did not report its worktree:\n%s", output)
	}
	return m[1]
}

func requireNoWorktreeBranch(t *testing.T, f Fixture, branch string) {
	t.Helper()
	out := f.Git("worktree", "list", "--porcelain").Stdout
	if strings.Contains(out, "branch refs/heads/"+branch) {
		t.Fatalf("released worktree still holds the shift branch:\n%s", out)
	}
}

func requireNoLease(t *testing.T, home string) {
	t.Helper()
	var found []string
	err := filepath.WalkDir(home, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Name() == "bench-lease" {
			found = append(found, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk BENCH_HOME for leases: %v", err)
	}
	if len(found) > 0 {
		t.Fatalf("shift worktree lease was not released: %s", strings.Join(found, ", "))
	}
}

func requireEqual(t *testing.T, got, want, msg string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: got %q, want %q", msg, got, want)
	}
}

func runGitAt(t *testing.T, root string, args ...string) string {
	t.Helper()
	f := Fixture{t: t, Root: root, Env: isolatedEnv(t, root)}
	return f.Git(args...).Stdout
}

func waitSeconds(t *testing.T, seconds int) {
	t.Helper()
	time.Sleep(time.Duration(seconds) * time.Second)
}

func strPtr(s string) *string {
	return &s
}

func requireFileContains(t *testing.T, path, needle, msg string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if !strings.Contains(string(data), needle) {
		t.Fatal(msg)
	}
}
