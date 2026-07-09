package gate

// The runner engine: concurrency, output shape, exit codes, cancellation, and
// optional/required start behavior. The PhasesCommand surface (signal handling,
// the phase table, shellcheck expansion) is tested in phases_command_test.go.

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/toon"
)

func TestRunnerRunsPhasesConcurrently(t *testing.T) {
	root := t.TempDir()
	phases := []Phase{
		fakePhase("alpha", "printf 'alpha\\n'; sleep 0.25"),
		fakePhase("bravo", "printf 'bravo\\n'; sleep 0.25"),
		fakePhase("charlie", "printf 'charlie\\n'; sleep 0.25"),
		fakePhase("delta", "printf 'delta\\n'; sleep 0.25"),
	}

	var stdout, stderr bytes.Buffer
	start := time.Now()
	rc := runPhases(context.Background(), root, phases, outerMode, &stdout, &stderr)
	elapsed := time.Since(start)

	if rc != 0 {
		t.Fatalf("runPhases rc = %d, stderr:\n%s", rc, stderr.String())
	}
	if elapsed >= 750*time.Millisecond {
		t.Fatalf("runPhases took %v, want concurrent execution under serial sum", elapsed)
	}
	out := stdout.String() + stderr.String()
	for _, marker := range []string{"[alpha] alpha", "[bravo] bravo", "[charlie] charlie", "[delta] delta"} {
		if !strings.Contains(out, marker) {
			t.Fatalf("output missing %q:\n%s", marker, out)
		}
	}
}

func TestRunnerPrefixesAndKeepsLinesIntact(t *testing.T) {
	root := t.TempDir()
	phases := []Phase{
		fakePhase("alpha", "printf 'alpha'; sleep 0.05; printf ' one\\nalpha two\\n'"),
		fakePhase("bravo", "printf 'bravo'; sleep 0.05; printf ' one\\nbravo two\\n'"),
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phases, outerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("runPhases rc = %d, stderr:\n%s", rc, stderr.String())
	}

	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	got := map[string]bool{}
	for _, line := range lines {
		if strings.HasPrefix(line, "phase ") || strings.HasPrefix(line, "gate: ") {
			continue
		}
		got[line] = true
		if !strings.HasPrefix(line, "[alpha] ") && !strings.HasPrefix(line, "[bravo] ") {
			t.Fatalf("unprefixed phase line %q in output:\n%s", line, stdout.String())
		}
	}
	for _, want := range []string{
		"[alpha] alpha one",
		"[alpha] alpha two",
		"[bravo] bravo one",
		"[bravo] bravo two",
	} {
		if !got[want] {
			t.Fatalf("missing intact line %q; got %#v\nfull output:\n%s", want, got, stdout.String())
		}
	}
}

func TestRunnerFinalLineAndExitCodes(t *testing.T) {
	t.Run("green", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		rc := runPhases(context.Background(), t.TempDir(), []Phase{fakePhase("ok", "true")}, outerMode, &stdout, &stderr)
		if rc != 0 {
			t.Fatalf("runPhases rc = %d, stderr:\n%s", rc, stderr.String())
		}
		if !strings.HasSuffix(stdout.String(), "gate: green\n") {
			t.Fatalf("stdout final line = %q, want gate: green", stdout.String())
		}
	})

	t.Run("red", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		rc := runPhases(context.Background(), t.TempDir(), []Phase{fakePhase("bad", "exit 7")}, outerMode, &stdout, &stderr)
		if rc != 1 {
			t.Fatalf("runPhases rc = %d, want 1", rc)
		}
		if !strings.HasSuffix(stderr.String(), "gate: red\n") {
			t.Fatalf("stderr final line = %q, want gate: red", stderr.String())
		}
	})

	t.Run("not in repo", func(t *testing.T) {
		cwd := mustGetwd(t)
		t.Cleanup(func() { mustChdir(t, cwd) })
		mustChdir(t, t.TempDir())

		var stdout, stderr bytes.Buffer
		rc := PhasesCommand(nil, &stdout, &stderr)
		if rc != 3 {
			t.Fatalf("PhasesCommand rc = %d, want 3; stderr:\n%s", rc, stderr.String())
		}
		if strings.TrimSpace(stderr.String()) != toon.NotInRepo() {
			t.Fatalf("stderr = %q, want gate not-in-repo line", stderr.String())
		}
	})
}

func TestRunnerAggregatesAllPhases(t *testing.T) {
	root := t.TempDir()
	phases := []Phase{
		fakePhase("first", "printf 'first\\n'"),
		fakePhase("bad", "printf 'bad\\n'; exit 1"),
		fakePhase("third", "printf 'third\\n'"),
		fakePhase("fourth", "printf 'fourth\\n'"),
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phases, outerMode, &stdout, &stderr)
	if rc != 1 {
		t.Fatalf("runPhases rc = %d, want 1", rc)
	}
	out := stdout.String() + stderr.String()
	for _, marker := range []string{"[first] first", "[bad] bad", "[third] third", "[fourth] fourth"} {
		if !strings.Contains(out, marker) {
			t.Fatalf("output missing %q:\n%s", marker, out)
		}
	}
}

func TestRunnerInnerModeByteShape(t *testing.T) {
	root := t.TempDir()
	phases := []Phase{
		fakePhase("conformance", "printf 'conformance\\n'"),
		fakePhase("contract", "printf 'contract\\n'"),
		fakePhase("shellcheck", "printf 'shellcheck\\n'"),
		fakePhase("canary", "printf 'canary\\n'"),
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phasesForMode(phases, innerMode), innerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("runPhases rc = %d, stderr:\n%s", rc, stderr.String())
	}
	if got, want := stdout.String(), "conformance\ncontract\nshellcheck\ngate: green\n"; got != want {
		t.Fatalf("inner output = %q, want %q", got, want)
	}
	if strings.Contains(stdout.String()+stderr.String(), "[") {
		t.Fatalf("inner output gained outer prefix bytes:\nstdout=%q\nstderr=%q", stdout.String(), stderr.String())
	}
}

func TestRunnerSerialPhaseCompletesBeforeConcurrentPhasesStart(t *testing.T) {
	root := t.TempDir()
	marker := filepath.Join(root, "built")
	phases := []Phase{
		{
			Name:   "build",
			Argv:   []string{"bash", "-c", `sleep 0.1; touch "$1"`, "bash", marker},
			Serial: true,
		},
		// Each reader phase goes red unless the serial phase already finished:
		// a runner that launches everything concurrently starts these ~100ms
		// before the marker exists.
		{Name: "alpha", Argv: []string{"bash", "-c", `test -f "$1"`, "bash", marker}},
		{Name: "bravo", Argv: []string{"bash", "-c", `test -f "$1"`, "bash", marker}},
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phases, outerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("runPhases rc = %d; a reader phase started before the serial phase finished\nstdout=%q\nstderr=%q", rc, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "phase build: green") {
		t.Fatalf("serial phase missing from summaries:\n%s", stdout.String())
	}
}

func TestRunnerSerialRedFailsFast(t *testing.T) {
	root := t.TempDir()
	leak := filepath.Join(root, "leak")
	phases := []Phase{
		{Name: "build", Argv: []string{"bash", "-c", "exit 3"}, Serial: true},
		{Name: "alpha", Argv: []string{"bash", "-c", `touch "$1"`, "bash", leak}},
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phases, outerMode, &stdout, &stderr)
	if rc != 1 {
		t.Fatalf("runPhases rc = %d, want 1", rc)
	}
	out := stdout.String() + stderr.String()
	if !strings.Contains(out, "phase build: red (exit 3)") || !strings.Contains(out, "gate: red") {
		t.Fatalf("serial red not attributed:\n%s", out)
	}
	if _, err := os.Stat(leak); err == nil {
		t.Fatalf("phase after failed serial phase still ran (it would grade a stale or absent binary)")
	}
}

func TestRunnerSerialRedFailsFastInnerMode(t *testing.T) {
	root := t.TempDir()
	leak := filepath.Join(root, "leak")
	phases := []Phase{
		{Name: "build", Argv: []string{"bash", "-c", "exit 3"}, Serial: true},
		{Name: "alpha", Argv: []string{"bash", "-c", `touch "$1"`, "bash", leak}},
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phases, innerMode, &stdout, &stderr)
	if rc != 1 {
		t.Fatalf("runPhases rc = %d, want 1", rc)
	}
	if !strings.HasSuffix(stderr.String(), "gate: red\n") {
		t.Fatalf("stderr final line = %q, want gate: red", stderr.String())
	}
	if _, err := os.Stat(leak); err == nil {
		t.Fatalf("inner-mode phase after failed serial phase still ran")
	}
}

func TestRunnerCancelKillsGroup(t *testing.T) {
	root := t.TempDir()
	pidfile := filepath.Join(root, "sleep.pid")
	phase := Phase{
		Name: "slow",
		Argv: []string{"bash", "-c", `sleep 30 & echo $! > "$1"; wait`, "bash", pidfile},
	}
	ctx, cancel := context.WithCancel(context.Background())
	var stdout, stderr bytes.Buffer
	done := make(chan int, 1)
	go func() {
		done <- runPhases(ctx, root, []Phase{phase}, outerMode, &stdout, &stderr)
	}()

	pid := waitForPIDFile(t, pidfile)
	t.Cleanup(func() { _ = syscall.Kill(pid, syscall.SIGKILL) })
	cancel()

	select {
	case rc := <-done:
		if rc != 130 {
			t.Fatalf("runPhases rc = %d, want 130; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
		}
	case <-time.After(5 * time.Second):
		t.Fatal("runPhases did not return after cancellation")
	}
	waitForProcessExit(t, pid)
}

func TestRunnerRootWithSpace(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "root with space")
	if err := os.Mkdir(root, 0o755); err != nil {
		t.Fatalf("mkdir root with space: %v", err)
	}
	phase := Phase{
		Name: "pwd",
		Argv: []string{"bash", "-c", `test "$PWD" = "$1"; printf 'root=%s\n' "$1"`, "bash", root},
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, []Phase{phase}, outerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("runPhases rc = %d, stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "[pwd] root="+root) {
		t.Fatalf("stdout does not show literal root argv:\n%s", stdout.String())
	}
}

func TestRunnerShellcheckAbsentSkips(t *testing.T) {
	root := t.TempDir()
	phase := Phase{Name: "shellcheck", Argv: []string{"definitely-not-installed-shellcheck-for-bench-test"}, Optional: true}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, []Phase{phase}, outerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("runPhases rc = %d, want skip to stay green; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
	}
	out := stdout.String() + stderr.String()
	if strings.Contains(out, "gate: red") || strings.Contains(out, "shellcheck reported issues") {
		t.Fatalf("optional missing shellcheck looked red:\n%s", out)
	}
	if !strings.Contains(out, "phase shellcheck: skipped") {
		t.Fatalf("missing shellcheck skip summary:\n%s", out)
	}
}

func TestRunnerOptionalBrokenSymlinkSkips(t *testing.T) {
	root := t.TempDir()
	bin := filepath.Join(root, "bin")
	if err := os.Mkdir(bin, 0o755); err != nil {
		t.Fatalf("mkdir bin: %v", err)
	}
	interpreters := filepath.Join(root, "interpreters")
	if err := os.Mkdir(interpreters, 0o755); err != nil {
		t.Fatalf("mkdir interpreters: %v", err)
	}
	brokenInterpreter := filepath.Join(interpreters, "sh")
	if err := os.Symlink(filepath.Join(root, "missing-sh"), brokenInterpreter); err != nil {
		t.Fatalf("symlink broken interpreter: %v", err)
	}
	shellcheck := filepath.Join(bin, "shellcheck")
	if err := os.WriteFile(shellcheck, []byte("#!"+brokenInterpreter+"\n"), 0o755); err != nil {
		t.Fatalf("write shellcheck shim: %v", err)
	}
	t.Setenv("PATH", bin)
	phase := Phase{Name: "shellcheck", Argv: []string{"shellcheck"}, Optional: true}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, []Phase{phase}, outerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("runPhases rc = %d, want broken optional executable skipped; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
	}
	if out := stdout.String() + stderr.String(); !strings.Contains(out, "phase shellcheck: skipped") || strings.Contains(out, "gate: red") {
		t.Fatalf("broken optional executable did not produce a green skip:\n%s", out)
	}
}

func TestRunnerShellcheckAbsentSkipVerdictNamesNotInstalled(t *testing.T) {
	root := t.TempDir()
	phase := Phase{Name: "shellcheck", Argv: []string{"definitely-absent-shellcheck-binary-for-bench-test"}, Optional: true}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, []Phase{phase}, outerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("absent optional binary went red rc=%d; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
	}
	if out := stdout.String() + stderr.String(); !strings.Contains(out, "phase shellcheck: skipped (not installed)") {
		t.Fatalf("absent shellcheck skip verdict is silent about why (want 'skipped (not installed)'):\n%s", out)
	}
}

func TestRunnerOptionalUnexecutableStubGoesRed(t *testing.T) {
	root := t.TempDir()
	bin := filepath.Join(root, "bin")
	if err := os.Mkdir(bin, 0o755); err != nil {
		t.Fatalf("mkdir bin: %v", err)
	}
	stub := filepath.Join(bin, "shellcheck")
	if err := os.WriteFile(stub, []byte("#!/usr/bin/env bash\nexit 0\n"), 0o000); err != nil {
		t.Fatalf("write unexecutable stub: %v", err)
	}
	t.Setenv("PATH", bin)
	phase := Phase{Name: "shellcheck", Argv: []string{"shellcheck"}, Optional: true}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, []Phase{phase}, outerMode, &stdout, &stderr)
	if rc != 1 {
		t.Fatalf("present-but-unexecutable shellcheck did not go red rc=%d; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
	}
	out := stdout.String() + stderr.String()
	if strings.Contains(out, "skipped") {
		t.Fatalf("present-but-unexecutable shellcheck was masked as a skip:\n%s", out)
	}
	if !strings.Contains(out, "phase shellcheck: red") || !strings.Contains(out, "permission denied") {
		t.Fatalf("red verdict does not name the exec failure:\n%s", out)
	}
}

func TestRunnerRequiredStartFailureRed(t *testing.T) {
	root := t.TempDir()
	phase := Phase{Name: "required", Argv: []string{filepath.Join(root, "missing-required")}}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, []Phase{phase}, outerMode, &stdout, &stderr)
	if rc != 1 {
		t.Fatalf("runPhases rc = %d, want required start failure red; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
	}
	out := stdout.String() + stderr.String()
	if !strings.Contains(out, "phase required: red") || !strings.Contains(out, "gate: red") {
		t.Fatalf("required start failure did not stay red:\n%s", out)
	}
}

func fakePhase(name, script string) Phase {
	return Phase{Name: name, Argv: []string{"bash", "-c", script}}
}

func mustGetwd(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get cwd: %v", err)
	}
	return cwd
}

func mustChdir(t *testing.T, dir string) {
	t.Helper()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir %s: %v", dir, err)
	}
}

func waitForPIDFile(t *testing.T, path string) int {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		data, err := os.ReadFile(path)
		if err == nil {
			var pid int
			if _, scanErr := fmt.Sscanf(strings.TrimSpace(string(data)), "%d", &pid); scanErr == nil && pid > 0 {
				return pid
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("pid file %s was not written", path)
	return 0
}

func waitForProcessExit(t *testing.T, pid int) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		err := syscall.Kill(pid, 0)
		if errors.Is(err, syscall.ESRCH) {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("process %d still exists after cancellation", pid)
}
