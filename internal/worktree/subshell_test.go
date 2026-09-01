package worktree

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/intent"
)

func TestSubshellNormalExitReleasesItsAssignment(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := t.TempDir()
	pwd := filepath.Join(t.TempDir(), "shell-pwd")
	shell := writeSubshellScript(t, "printf %s \"$PWD\" > \"$BENCH_SUBSHELL_PWD\"\n")
	environ := append(os.Environ(), "BENCH_SUBSHELL_PWD="+pwd)

	var stdout, stderr bytes.Buffer
	code := subshellAt(root, home, shell, environ, nil, strings.NewReader(""), &stdout, &stderr)
	requireTest(t, code == 0, "subshell exit = %d, stderr %q", code, stderr.String())
	assignments, err := intent.Assignments(root)
	requireTest(t, err == nil && len(assignments) == 0, "subshell assignments = %#v, %v; want no released owner", assignments, err)
	path := strings.TrimPrefix(strings.SplitN(stderr.String(), "  ", 2)[0], "🪵 worktree: ")
	requireTest(t, filepath.IsAbs(path), "subshell announcement path = %q, want literal absolute path", path)
	childDir, err := os.ReadFile(pwd)
	requireTest(t, err == nil && string(childDir) == path, "subshell child directory = %q, %v; want announced worktree %q", childDir, err, path)
}

const subshellSignalHelperEnv = "BENCH_SUBSHELL_SIGNAL_HELPER"

func TestSubshellSignalsLeaveAReclaimableLease(t *testing.T) {
	t.Parallel()
	if root := os.Getenv(subshellSignalHelperEnv); root != "" {
		home := os.Getenv("BENCH_SUBSHELL_HOME")
		signalValue := syscall.Signal(mustSubshellSignal(t, os.Getenv("BENCH_SUBSHELL_SIGNAL")))
		var stdout, stderr bytes.Buffer
		code := subshellAt(root, home, subshellShell(), os.Environ(), nil, strings.NewReader(""), &stdout, &stderr)
		requireTest(t, code == 128+int(signalValue), "signalled subshell exit = %d, want %d (stderr %q)", code, 128+int(signalValue), stderr.String())
		return
	}

	for _, signalValue := range []syscall.Signal{syscall.SIGINT, syscall.SIGTERM} {
		t.Run(signalValue.String(), func(t *testing.T) {
			root := newWorktreeRepo(t)
			home := t.TempDir()
			ready := filepath.Join(t.TempDir(), "ready")
			shell := writeSubshellScript(t, "printf ready > \"$BENCH_SUBSHELL_READY\"\nwhile :; do sleep 1; done\n")
			cmd := descendant(t, os.Args[0], "-test.run=^TestSubshellSignalsLeaveAReclaimableLease$", "-test.v")
			cmd.Env = append(os.Environ(), subshellSignalHelperEnv+"="+root, "BENCH_SUBSHELL_HOME="+home, "BENCH_SUBSHELL_READY="+ready, "BENCH_SUBSHELL_SIGNAL="+strconv.Itoa(int(signalValue)), "SHELL="+shell)
			requireTest(t, cmd.Start() == nil, "start subshell signal helper")
			waitForSubshellFile(t, ready)
			requireTest(t, cmd.Process.Signal(signalValue) == nil, "signal subshell helper with %s", signalValue)
			done := make(chan error, 1)
			go func() { done <- cmd.Wait() }()
			select {
			case err := <-done:
				requireTest(t, err == nil, "subshell signal helper: %v", err)
			case <-time.After(5 * time.Second):
				t.Fatal("subshell signal helper did not exit")
			}

			assignments, err := intent.Assignments(root)
			requireTest(t, err == nil && len(assignments) == 1, "signalled subshell assignments = %#v, %v; want one retained owner", assignments, err)
			lease, err := LeaseFile(assignments[0].Worktree)
			requireTest(t, err == nil && ProbeLease(lease) == LeaseDead, "signalled subshell lease = %q, %v; want a dead reclaimable lease", lease, err)
			plan, err := PlanExplicit(root, assignments[0].Worktree)
			requireTest(t, err == nil && plan.Action == ActionRemove, "signalled subshell plan = %#v, %v; want reclaimable removal", plan, err)
			var stderr bytes.Buffer
			code := ReleaseCommand(root, home, []string{"--request", assignments[0].RequestToken, assignments[0].Worktree}, io.Discard, &stderr)
			requireTest(t, code == 0, "release reclaimable signal state = %d, stderr %q", code, stderr.String())
		})
	}
}

func writeSubshellScript(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "shell")
	mustWrite(t, path, []byte("#!/bin/sh\n"+body), 0o755)
	return path
}

func waitForSubshellFile(t *testing.T, path string) {
	t.Helper()
	deadline := time.After(5 * time.Second)
	tick := time.NewTicker(10 * time.Millisecond)
	defer tick.Stop()
	for {
		if _, err := os.Stat(path); err == nil {
			return
		}
		select {
		case <-deadline:
			t.Fatalf("subshell did not create readiness file %s", path)
		case <-tick.C:
		}
	}
}

func mustSubshellSignal(t *testing.T, value string) int {
	t.Helper()
	var signalValue int
	var err error
	if signalValue, err = strconv.Atoi(value); err != nil {
		t.Fatal(err)
	}
	return signalValue
}
