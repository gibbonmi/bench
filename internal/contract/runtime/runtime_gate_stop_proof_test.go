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
)

type stopProofRun struct {
	f                                                 contract.Fixture
	gitdir, cache, count, wrapperExit, snapshot, kind string
	marker, started, gatePID, childPID                string
}

func proveStopResult(t *testing.T, variant string) {
	run := newStopProofRun(t, variant)
	var held *os.File
	switch variant {
	case "unarmed", "active":
		if err := syscall.Mkfifo(run.cache, 0o600); err != nil {
			t.Fatal(err)
		}
	case "lock-open":
		if err := os.Mkdir(filepath.Join(run.gitdir, "bench-gate.lock"), 0o700); err != nil {
			t.Fatal(err)
		}
	case "lock-acquire":
		held = holdStopProofLock(t, run)
		defer held.Close()
	case "pending-persistence":
		if err := os.Mkdir(run.cache, 0o700); err != nil {
			t.Fatal(err)
		}
	}

	env := map[string]string{"BENCH_SHIFT": "1"}
	stdin := "{}\n"
	if variant == "unarmed" {
		delete(env, "BENCH_SHIFT")
	}
	if variant == "active" {
		stdin = `{"stop_hook_active":true}` + "\n"
	}
	if variant == "subject-build" {
		env["PATH"] = story5FailingGitPath(t) + string(os.PathListSeparator) + os.Getenv("PATH")
	}

	var probe contract.Probe
	if variant == "cancellation" {
		probe = runCancelledStopProof(t, run, env, stdin)
	} else if variant == "unarmed" || variant == "active" {
		probe = runBoundedStopProof(t, run, env, stdin)
	} else {
		probe = runStopHook(t, run.f, env, stdin)
	}
	wantExit := 2
	if variant == "green" || variant == "unarmed" || variant == "active" {
		wantExit = 0
	}
	if probe.ExitCode != wantExit {
		t.Fatalf("%s Stop exit = %d, want %d\n%s%s", variant, probe.ExitCode, wantExit, probe.Stdout, probe.Stderr)
	}

	armed := variant != "unarmed" && variant != "active"
	wantRuns := 0
	if armed {
		wantRuns = 1
	}
	if got := len(contract.NonEmptyLines(readOptionalStopProof(t, run.count))); got != wantRuns {
		t.Fatalf("%s standalone wrapper runs = %d, want %d", variant, got, wantRuns)
	}
	if armed {
		wantWrapperExit := map[string]int{
			"green": 0, "red": 23, "lock-open": 1, "lock-acquire": 1,
			"pending-persistence": 1, "final-persistence": 1, "subject-build": 1,
			"drift": 1, "cancellation": 130, "start-failure": 1, "no-gate": 3,
		}[variant]
		if got := parseStopProofPID(t, run.wrapperExit); got != wantWrapperExit {
			t.Fatalf("%s standalone wrapper exit = %d, want %d", variant, got, wantWrapperExit)
		}
		requireStopWrapperVerdictPreserved(t, run)
	} else {
		if info, err := os.Lstat(run.cache); err != nil || info.Mode()&os.ModeNamedPipe == 0 {
			t.Fatalf("%s Stop inspected or rewrote guard FIFO: info=%v err=%v", variant, info, err)
		}
	}

	wantMessage := map[string]string{
		"red": "gate-red-23", "lock-open": "gate lock unavailable",
		"lock-acquire":        "gate execution already in progress",
		"pending-persistence": "gate pending persistence failed",
		"final-persistence":   "gate final persistence failed",
		"subject-build":       "gate subject unavailable",
		"drift":               "gate subject changed during execution",
		"cancellation":        "gate-cancelled-130",
		"start-failure":       "BLOCKED: the gate is red",
		"no-gate":             "no gate found",
	}[variant]
	if wantMessage != "" && !strings.Contains(probe.Stdout+probe.Stderr, wantMessage) {
		t.Fatalf("%s Stop result missing %q:\n%s%s", variant, wantMessage, probe.Stdout, probe.Stderr)
	}
	if !armed || variant == "lock-open" || variant == "lock-acquire" || variant == "pending-persistence" || variant == "subject-build" || variant == "start-failure" || variant == "no-gate" {
		if _, err := os.Stat(run.marker); err == nil {
			t.Fatalf("%s Stop ran the gate unexpectedly", variant)
		}
	}
}

func newStopProofRun(t *testing.T, variant string) stopProofRun {
	t.Helper()
	f := copiedCLIHookFixture(t, true)
	manifest := `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}` + "\n"
	gate := "#!/usr/bin/env bash\nprintf 'gate-ran\\n' > .git/r16-gate-marker\nexit 0\n"
	switch variant {
	case "red":
		gate = "#!/usr/bin/env bash\nprintf 'gate-ran\\n' > .git/r16-gate-marker\nprintf 'gate-red-23\\n'\nexit 23\n"
	case "final-persistence":
		gate = "#!/usr/bin/env bash\nprintf 'gate-ran\\n' > .git/r16-gate-marker\nrm -f .git/bench-last-gate\nmkdir .git/bench-last-gate\nexit 0\n"
	case "drift":
		gate = "#!/usr/bin/env bash\nprintf 'gate-ran\\n' > .git/r16-gate-marker\nprintf drift >> tracked.txt\nexit 0\n"
		manifest = `{"schema":1,"closure":"local","environment":[],"paths":["tracked.txt"],"tools":[]}` + "\n"
	case "cancellation":
		gate = `#!/usr/bin/env bash
echo $$ > .git/r16-gate-pid
sleep 300 & child=$!
echo "$child" > .git/r16-child-pid
trap 'kill "$child" 2>/dev/null || true; wait "$child" 2>/dev/null || true; printf "gate-cancelled-130\n"; exit 130' INT TERM
touch .git/r16-gate-started
wait "$child"
`
	case "start-failure":
		gate = "#!/definitely/missing\n"
	case "no-gate":
		gate = ""
	}
	if gate != "" {
		f.WriteExecutable(".bench/gate.sh", gate)
		f.WriteFile(".bench/gate-inputs.json", manifest)
	}
	f.WriteFile("tracked.txt", "base\n")
	f.CommitAll("R16 Stop fixture")
	gitdir := gitDir(t, f)
	run := stopProofRun{
		f: f, gitdir: gitdir, cache: filepath.Join(gitdir, "bench-last-gate"),
		count: filepath.Join(gitdir, "r16-wrapper-runs"), wrapperExit: filepath.Join(gitdir, "r16-wrapper-exit"),
		snapshot: filepath.Join(gitdir, "r16-wrapper-cache"),
		kind:     filepath.Join(gitdir, "r16-wrapper-kind"), marker: filepath.Join(gitdir, "r16-gate-marker"),
		started: filepath.Join(gitdir, "r16-gate-started"), gatePID: filepath.Join(gitdir, "r16-gate-pid"),
		childPID: filepath.Join(gitdir, "r16-child-pid"),
	}
	installStopCountingWrapper(t, run)
	return run
}

func installStopCountingWrapper(t *testing.T, run stopProofRun) {
	t.Helper()
	real := benchPath(t)
	script := fmt.Sprintf(`#!/usr/bin/env bash
if [[ "${1:-}" == gate ]]; then
  printf 'run\n' >> %q
  rc=0
  bash %q "$@" || rc=$?
  printf '%%s\n' "$rc" > %q
  if [[ -f %q ]]; then
    /bin/cp %q %q
    printf 'file\n' > %q
  elif [[ -d %q ]]; then
    printf 'dir\n' > %q
  else
    printf 'absent\n' > %q
  fi
  exit "$rc"
fi
exec bash %q "$@"
`, run.count, real, run.wrapperExit, run.cache, run.cache, run.snapshot, run.kind, run.cache, run.kind, run.kind, real)
	run.f.WriteExecutable("bin/bench.sh", script)
}

func holdStopProofLock(t *testing.T, run stopProofRun) *os.File {
	t.Helper()
	held, err := os.OpenFile(filepath.Join(run.gitdir, "bench-gate.lock"), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if err := syscall.Flock(int(held.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		held.Close()
		t.Fatal(err)
	}
	return held
}

func requireStopWrapperVerdictPreserved(t *testing.T, run stopProofRun) {
	t.Helper()
	switch strings.TrimSpace(readOptionalStopProof(t, run.kind)) {
	case "file":
		want := readOptionalStopProof(t, run.snapshot)
		if got := readOptionalStopProof(t, run.cache); got != want {
			t.Fatalf("Stop performed a second verdict write\nwrapper=%q\nafter=%q", want, got)
		}
	case "dir":
		if info, err := os.Stat(run.cache); err != nil || !info.IsDir() {
			t.Fatalf("Stop replaced wrapper-owned verdict directory: info=%v err=%v", info, err)
		}
	case "absent":
		if _, err := os.Lstat(run.cache); !os.IsNotExist(err) {
			t.Fatalf("Stop wrote after wrapper left no verdict: %v", err)
		}
	default:
		t.Fatalf("standalone wrapper did not snapshot its verdict outcome")
	}
}

func readOptionalStopProof(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return ""
	}
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func runBoundedStopProof(t *testing.T, run stopProofRun, env map[string]string, stdin string) contract.Probe {
	t.Helper()
	var out bytes.Buffer
	cmd := exec.Command("bash", filepath.Join(contract.SubjectRoot(t), ".bench", "hooks", "stop.sh"))
	cmd.Dir, cmd.Env, cmd.Stdin = run.f.Root, surfaceEnv(run.f, env), strings.NewReader(stdin)
	cmd.Stdout, cmd.Stderr = &out, &out
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		exit := 0
		if ee, ok := err.(*exec.ExitError); ok {
			exit = ee.ExitCode()
		} else if err != nil {
			exit = 1
		}
		return contract.Probe{ExitCode: exit, Stdout: out.String()}
	case <-time.After(3 * time.Second):
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		<-done
		return contract.Probe{ExitCode: -1, Stdout: out.String(), TimedOut: true}
	}
}

func runCancelledStopProof(t *testing.T, run stopProofRun, env map[string]string, stdin string) contract.Probe {
	t.Helper()
	var out bytes.Buffer
	cmd := exec.Command("bash", filepath.Join(contract.SubjectRoot(t), ".bench", "hooks", "stop.sh"))
	cmd.Dir, cmd.Env, cmd.Stdin = run.f.Root, surfaceEnv(run.f, env), strings.NewReader(stdin)
	cmd.Stdout, cmd.Stderr = &out, &out
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) }()
	waitStopProofPath(t, run.started, cmd)
	gatePID := parseStopProofPID(t, run.gatePID)
	childPID := parseStopProofPID(t, run.childPID)
	defer func() {
		_ = syscall.Kill(-gatePID, syscall.SIGKILL)
		_ = syscall.Kill(childPID, syscall.SIGKILL)
	}()
	if err := syscall.Kill(-gatePID, syscall.SIGINT); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	var err error
	select {
	case err = <-done:
	case <-time.After(5 * time.Second):
		t.Fatalf("cancelled Stop did not exit: %s", out.String())
	}
	exit := 0
	if ee, ok := err.(*exec.ExitError); ok {
		exit = ee.ExitCode()
	} else if err != nil {
		exit = 1
	}
	deadline := time.Now().Add(3 * time.Second)
	for syscall.Kill(childPID, 0) == nil && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if syscall.Kill(childPID, 0) == nil {
		t.Fatalf("cancelled Stop left gate child %d alive", childPID)
	}
	return contract.Probe{ExitCode: exit, Stdout: out.String()}
}

func waitStopProofPath(t *testing.T, path string, cmd *exec.Cmd) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	t.Fatalf("Stop proof did not reach %s", filepath.Base(path))
}

func parseStopProofPID(t *testing.T, path string) int {
	t.Helper()
	var pid int
	if _, err := fmt.Sscanf(strings.TrimSpace(readOptionalStopProof(t, path)), "%d", &pid); err != nil {
		t.Fatal(err)
	}
	return pid
}
