package gate

// The PhasesCommand surface: signal handling across the process boundary, the
// benchkit phase table, and shellcheck file expansion. The runner engine's own
// tests (concurrency, output shape, exit codes, cancel) live in phases_test.go.

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
)

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

func phaseNames(phases []Phase) []string {
	names := make([]string, 0, len(phases))
	for _, phase := range phases {
		names = append(names, phase.Name)
	}
	return names
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
