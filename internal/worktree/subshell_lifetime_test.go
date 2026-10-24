package worktree

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

const shellSignalHelperEnv = "BENCH_SHELL_SIGNAL_HELPER"

// TestShellSignal covers PL82. A TERM reaches the shell's group through the lifetime, and
// the shutdown completes: no member of the group survives. The interrupted session still
// keeps its assignment and its lease file for a reclaim.
func TestShellSignal(t *testing.T) {
	t.Parallel()
	if root := os.Getenv(shellSignalHelperEnv); root != "" {
		var stdout, stderr bytes.Buffer
		code := subshellAt(root, newAmbient(os.Getenv("BENCH_SUBSHELL_HOME"), &stderr), subshellShell(), os.Environ(), nil, strings.NewReader(""), &stdout, &stderr)
		requireTest(t, code == 128+int(syscall.SIGTERM), "signalled shell exit = %d (stderr %q)", code, stderr.String())
		return
	}
	root, home := newWorktreeRepo(t), t.TempDir()
	ready := filepath.Join(t.TempDir(), "ready")
	shell := writeSubshellScript(t, "printf %s $$ > \"$BENCH_SUBSHELL_READY.tmp\"\nmv \"$BENCH_SUBSHELL_READY.tmp\" \"$BENCH_SUBSHELL_READY\"\nwhile :; do sleep 1; done\n")
	cmd := descendant(t, os.Args[0], "-test.run=^TestShellSignal$", "-test.v")
	cmd.Env = append(os.Environ(), shellSignalHelperEnv+"="+root, "BENCH_SUBSHELL_HOME="+home, "BENCH_SUBSHELL_READY="+ready, "SHELL="+shell)
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	requireTest(t, cmd.Start() == nil, "start the shell signal helper")
	waitForSubshellFile(t, ready)
	raw, err := os.ReadFile(ready)
	requireTest(t, err == nil, "read the shell PID: %v", err)
	group, err := strconv.Atoi(string(raw))
	requireTest(t, err == nil && group > 0, "shell PID %q: %v", raw, err)
	requireTest(t, cmd.Process.Signal(syscall.SIGTERM) == nil, "signal the shell helper")
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	// The shell dies at the forwarded signal, so the helper's exit contains no grace.
	window := bounds.TestDeadline(bounds.ProcessFinalWindow)
	select {
	case err := <-done:
		requireTest(t, err == nil, "shell helper: %v\n%s", err, &output)
	case <-time.After(window):
		t.Fatal(bounds.TestTimeoutVerdict("the shell signal helper to exit", window))
	}

	requireTest(t, errors.Is(syscall.Kill(-group, 0), syscall.ESRCH), "the shell group %d survived its interrupted shutdown", group)
	assignments, err := intent.Assignments(root)
	requireTest(t, err == nil && len(assignments) == 1, "interrupted shell assignments = %#v, %v; want the one retained owner", assignments, err)
	lease, err := LeaseFile(assignments[0].Worktree)
	requireTest(t, err == nil, "lease path: %v", err)
	_, statErr := os.Stat(lease)
	requireTest(t, statErr == nil, "the interrupted shell's lease is gone: %v", statErr)
}

// TestShellRefusesAMalformedVector covers the shell's launch refusal. A malformed inherited
// vector starts no shell, and the session releases its assignment as a failed start.
func TestShellRefusesAMalformedVector(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	marker := filepath.Join(t.TempDir(), "started")
	shell := writeSubshellScript(t, "touch \""+marker+"\"\n")
	var stdout, stderr bytes.Buffer
	subshellAt(root, newAmbient(t.TempDir(), &stderr), shell, append(os.Environ(), protection.Env+"=not-json"), nil, strings.NewReader(""), &stdout, &stderr)
	_, statErr := os.Stat(marker)
	requireTest(t, os.IsNotExist(statErr) && strings.Contains(stderr.String(), protection.Env), "malformed vector shell marker %v, stderr %q; want a refused start", statErr, stderr.String())
	assignments, err := intent.Assignments(root)
	requireTest(t, err == nil && len(assignments) == 0, "refused shell assignments = %#v, %v; want the released owner", assignments, err)
}
