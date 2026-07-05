package gate

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"
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
		if strings.TrimSpace(stderr.String()) != "gate: not in a git repo" {
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

func TestPhasesCommandSignalCancelsRunningPhaseGroups(t *testing.T) {
	root := t.TempDir()
	pidfile := filepath.Join(root, "sleep.pid")
	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("resolve test executable: %v", err)
	}
	cmd := exec.Command(exe, "-test.run=^TestPhasesCommandSignalHelper$", "--", root)
	cmd.Env = append(os.Environ(),
		"BENCH_TEST_PHASES_SIGNAL_HELPER=1",
		"BENCH_TEST_PHASES_PIDFILE="+pidfile,
		"BENCH_TEST_PHASES_ROOT="+root,
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("start helper: %v", err)
	}
	t.Cleanup(func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	})

	pid := waitForPIDFile(t, pidfile)
	t.Cleanup(func() { _ = syscall.Kill(pid, syscall.SIGKILL) })
	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		t.Fatalf("send SIGINT to helper: %v", err)
	}
	err = cmd.Wait()
	if err == nil {
		t.Fatalf("helper exited 0, want 130; stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok || exitErr.ExitCode() != 130 {
		t.Fatalf("helper exit = %v, want code 130; stdout=%q stderr=%q", err, stdout.String(), stderr.String())
	}
	waitForProcessExit(t, pid)
}

func TestPhasesCommandSignalHelper(t *testing.T) {
	if os.Getenv("BENCH_TEST_PHASES_SIGNAL_HELPER") != "1" {
		return
	}
	root := os.Getenv("BENCH_TEST_PHASES_ROOT")
	pidfile := os.Getenv("BENCH_TEST_PHASES_PIDFILE")
	benchkitPhasesForCommand = func(root, kit string) []Phase {
		return []Phase{{
			Name: "slow",
			Argv: []string{"bash", "-c", `sleep 30 & echo $! > "$1"; wait`, "bash", pidfile},
		}}
	}
	os.Exit(PhasesCommand([]string{root}, os.Stdout, os.Stderr))
}

func TestPhaseTable(t *testing.T) {
	root := "/tmp/root with spaces"
	kit := "/tmp/kit"
	phases := BenchkitPhases(root, kit)
	if len(phases) != 4 {
		t.Fatalf("BenchkitPhases len = %d, want 4: %#v", len(phases), phases)
	}
	if got, want := phaseNames(phases), []string{"conformance", "contract", "shellcheck", "canary"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("phase names = %#v, want %#v", got, want)
	}
	contractCount := 0
	for _, phase := range phases {
		for _, arg := range phase.Argv {
			if arg == "./internal/contract/..." {
				contractCount++
			}
		}
	}
	if contractCount != 1 {
		t.Fatalf("contract subtree argv count = %d, want exactly 1", contractCount)
	}
	if got := phaseNames(phasesForMode(phases, innerMode)); strings.Contains(strings.Join(got, ","), "canary") {
		t.Fatalf("inner phases include canary: %#v", got)
	}
}

func TestShellcheckPhaseExpandsHookAndLibShellFiles(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "bin", "bench.sh"), "#!/usr/bin/env bash\n")
	writeFile(t, filepath.Join(root, ".bench", "hooks", "z.sh"), "#!/bin/sh\n")
	writeFile(t, filepath.Join(root, ".bench", "hooks", "a.sh"), "#!/bin/sh\n")
	writeFile(t, filepath.Join(root, ".bench", "hooks", "README"), "not shell\n")
	writeFile(t, filepath.Join(root, ".bench", "lib", "resolve-bench.sh"), "#!/bin/sh\n")
	writeFile(t, filepath.Join(root, ".bench", "lib", "notes.txt"), "not shell\n")

	var shellcheck Phase
	for _, phase := range BenchkitPhases("/tmp/root-under-grade", root) {
		if phase.Name == "shellcheck" {
			shellcheck = phase
			break
		}
	}
	want := []string{
		"shellcheck",
		"-S",
		"warning",
		"bin/bench.sh",
		".bench/hooks/a.sh",
		".bench/hooks/z.sh",
		".bench/lib/resolve-bench.sh",
	}
	if !reflect.DeepEqual(shellcheck.Argv, want) {
		t.Fatalf("shellcheck argv = %#v, want %#v", shellcheck.Argv, want)
	}
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

func phaseNames(phases []Phase) []string {
	names := make([]string, 0, len(phases))
	for _, phase := range phases {
		names = append(names, phase.Name)
	}
	return names
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

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
