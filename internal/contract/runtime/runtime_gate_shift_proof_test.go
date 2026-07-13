package runtime

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/intent"
)

type shiftProofRun struct {
	f       contract.Fixture
	home    string
	pool    runtimePoolWorktrees
	base    string
	variant string
}

func proveShiftResult(t *testing.T, variant string) {
	run := newShiftProofRun(t, variant)
	var held *os.File
	if variant == "lock" {
		held = holdShiftGateLock(t, run)
		defer held.Close()
	}
	if variant == "persistence" {
		cache := filepath.Join(shiftProofGitDir(t, run), "bench-last-gate")
		if err := os.Mkdir(cache, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	if variant == "cancellation" {
		proveCancelledShift(t, run)
		return
	}

	probe := run.f.BenchEnv(run.env(), "shift", "R15-"+variant)
	if probe.ExitCode != 1 {
		t.Fatalf("%s shift exit = %d, want 1\n%s%s", variant, probe.ExitCode, probe.Stdout, probe.Stderr)
	}
	want := map[string]string{
		"red":         "gate failed — preserving iteration 1",
		"lock":        "gate execution already in progress",
		"persistence": "gate pending persistence failed",
		"drift":       "gate subject changed during execution",
	}[variant]
	if !strings.Contains(probe.Stdout+probe.Stderr, want) {
		t.Fatalf("%s shift result missing %q:\n%s%s", variant, want, probe.Stdout, probe.Stderr)
	}
	requirePreservedShift(t, run, probe.Stdout)
}

func newShiftProofRun(t *testing.T, variant string) shiftProofRun {
	t.Helper()
	gate := "#!/usr/bin/env bash\nexit 23\n"
	if variant == "drift" {
		gate = "#!/usr/bin/env bash\nprintf drift >> tracked.txt\nexit 0\n"
	}
	if variant == "cancellation" {
		gate = `#!/usr/bin/env bash
echo $$ > "$BENCH_TEST_STATE/gate-pgid"
: > "$BENCH_TEST_STATE/gate-started"
while :; do sleep .05; done
`
	}
	f := shiftFixture(t, gate)
	f.WriteExecutable("agent", "#!/usr/bin/env bash\nprintf 'charged\\n' > charged.txt\n")
	if variant == "drift" {
		f.WriteFile("tracked.txt", "base\n")
		f.WriteFile(".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":["BENCH_TEST_PROMPTS","BENCH_TEST_STATE"],"paths":["tracked.txt"],"tools":[]}`+"\n")
	}
	f.CommitAll("R15 shift fixture")
	home := t.TempDir()
	return shiftProofRun{f: f, home: home, pool: addRuntimePoolWorktrees(t, f, home), base: headSha(f), variant: variant}
}

func (r shiftProofRun) env() map[string]string {
	env := map[string]string{
		"BENCH_AGENT":     filepath.Join(r.f.Root, "agent"),
		"BENCH_HOME":      r.home,
		"BENCH_MAX_ITERS": "1",
	}
	if r.variant == "cancellation" {
		env["BENCH_TEST_STATE"] = filepath.Join(r.home, "cancel-state")
	}
	return env
}

func shiftProofGitDir(t *testing.T, run shiftProofRun) string {
	t.Helper()
	return strings.TrimSpace(contract.RunAt(t, run.f, run.pool.Warm, nil, "git", "rev-parse", "--absolute-git-dir").Stdout)
}

func holdShiftGateLock(t *testing.T, run shiftProofRun) *os.File {
	t.Helper()
	lock := filepath.Join(shiftProofGitDir(t, run), "bench-gate.lock")
	held, err := os.OpenFile(lock, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if err := syscall.Flock(int(held.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		held.Close()
		t.Fatal(err)
	}
	return held
}

func proveCancelledShift(t *testing.T, run shiftProofRun) {
	t.Helper()
	state := run.env()["BENCH_TEST_STATE"]
	if err := os.MkdirAll(state, 0o700); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	cmd := exec.Command("bash", benchPath(t), "shift", "R15-cancellation")
	cmd.Dir = run.f.Root
	cmd.Env = surfaceEnv(run.f, run.env())
	cmd.Stdout, cmd.Stderr = &out, &out
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	started := filepath.Join(state, "gate-started")
	deadline := time.Now().Add(5 * time.Second)
	for {
		if _, err := os.Stat(started); err == nil {
			break
		}
		if time.Now().After(deadline) {
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			t.Fatalf("shift gate did not start: %s", out.String())
		}
		time.Sleep(10 * time.Millisecond)
	}
	var gatePGID int
	if _, err := fmt.Sscanf(strings.TrimSpace(string(mustReadRuntime(t, filepath.Join(state, "gate-pgid")))), "%d", &gatePGID); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGINT); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		_ = syscall.Kill(-gatePGID, syscall.SIGKILL)
		t.Fatalf("cancelled shift did not exit: %s", out.String())
	}
	status, ok := cmd.ProcessState.Sys().(syscall.WaitStatus)
	if !ok || status.ExitStatus() != 130 {
		t.Fatalf("cancelled shift status = %v, want exit 130: %s", cmd.ProcessState.Sys(), out.String())
	}
	deadline = time.Now().Add(3 * time.Second)
	for syscall.Kill(-gatePGID, 0) == nil && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if syscall.Kill(-gatePGID, 0) == nil {
		t.Fatalf("cancelled gate process group %d survived", gatePGID)
	}
	requirePreservedShift(t, run, out.String())
}

func requirePreservedShift(t *testing.T, run shiftProofRun, output string) {
	t.Helper()
	worktree := shiftWorktree(t, output)
	if worktree != run.pool.Warm {
		t.Fatalf("shift acquired %q, want seeded warm worktree %q", worktree, run.pool.Warm)
	}
	branch := shiftBranchFromStart(t, output)
	requireRegisteredWorktree(t, run.f, worktree)
	lease := strings.TrimSpace(contract.RunAt(t, run.f, worktree, nil, "git", "rev-parse", "--git-path", "bench-lease").Stdout)
	if data, err := os.ReadFile(lease); err != nil || strings.TrimSpace(string(data)) == "" {
		t.Fatalf("failed shift lease not preserved: bytes=%q err=%v", data, err)
	}
	if got := strings.TrimSpace(runGitAt(t, worktree, "rev-parse", "HEAD")); got != run.base {
		t.Fatalf("failed shift moved HEAD to %s, want %s", got, run.base)
	}
	status := runGitAt(t, worktree, "status", "--porcelain=v1")
	if !strings.Contains(status, "charged.txt") || !strings.Contains(status, ".bench-objective") || !strings.Contains(status, ".bench-notes.md") {
		t.Fatalf("failed shift destructively cleaned agent/scratch changes:\n%s", status)
	}
	if got := strings.TrimSpace(runGitAt(t, worktree, "branch", "--show-current")); got != branch {
		t.Fatalf("failed shift branch = %q, want %q", got, branch)
	}
	entries, err := intent.Snapshot(run.f.Root)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.Kind == intent.KindShift && entry.Objective == "R15-"+run.variant && entry.Worktree == worktree && entry.Branch == branch {
			return
		}
	}
	t.Fatalf("failed shift lost correlated intent for branch=%s worktree=%s: %#v", branch, worktree, entries)
}

func shiftBranchFromStart(t *testing.T, output string) string {
	t.Helper()
	const prefix = "▶ shift on "
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, prefix) {
			branch, _, ok := strings.Cut(strings.TrimPrefix(line, prefix), " — objective:")
			if ok && branch != "" {
				return branch
			}
		}
	}
	t.Fatalf("shift start did not name branch:\n%s", output)
	return ""
}
