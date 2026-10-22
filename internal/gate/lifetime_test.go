package gate

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/subprocess/processtest"
)

// signalLogScript logs the first catchable signal to $1 and exits. It publishes its PID to $2
// once both traps are installed. A KILL opener leaves the log absent.
const signalLogScript = `trap 'echo INT >"$1"; exit 0' INT
trap 'echo TERM >"$1"; exit 0' TERM
echo $$ >"$2.partial" && mv "$2.partial" "$2"
while :; do sleep 0.05; done`

// readyDescendantScript leaves a descendant in its group with no shared stream. The leader
// exits zero once the descendant publishes its PID to $1.
const readyDescendantScript = `sh -c 'echo $$ >"$1.partial" && mv "$1.partial" "$1" && exec sleep 300' sh "$1" </dev/null >/dev/null 2>&1 &
while [ ! -s "$1" ]; do sleep 0.01; done`

// runSignalFixture runs the signal fixture as a real phase and cancels the run with cause
// once the fixture is ready. It returns the phase result and the logged opener.
func runSignalFixture(t *testing.T, cause error) (phaseResult, string) {
	t.Helper()
	dir := t.TempDir()
	log, pid := filepath.Join(dir, "signal.log"), filepath.Join(dir, "child.pid")
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)
	phase := Phase{Name: "signal", Argv: []string{"sh", "-c", signalLogScript, "sh", log, pid}}
	done := make(chan phaseResult, 1)
	go func() { done <- runPhase(ctx, dir, phase, io.Discard, io.Discard) }()
	processtest.WaitPID(t, pid)
	cancel(cause)
	select {
	case result := <-done:
		opener, _ := os.ReadFile(log)
		return result, strings.TrimSpace(string(opener))
	case <-time.After(failureOutcomeDeadline()):
		t.Fatal("the canceled phase did not return")
		return phaseResult{}, ""
	}
}

func TestGateInterruptOpensWithINTAndExits130(t *testing.T) {
	t.Parallel()
	result, opener := runSignalFixture(t, context.Canceled)
	if opener != "INT" || result.Code != 130 || !result.Interrupted {
		t.Fatalf("interrupted phase opener = %q, exit = %d, interrupted = %v, want INT, 130, true", opener, result.Code, result.Interrupted)
	}
}

func TestGateTimeoutOpensWithTERMAndExits124(t *testing.T) {
	t.Parallel()
	result, opener := runSignalFixture(t, errGateTimeout)
	if opener != "TERM" || result.Code != 124 || !result.Interrupted {
		t.Fatalf("timed-out phase opener = %q, exit = %d, interrupted = %v, want TERM, 124, true", opener, result.Code, result.Interrupted)
	}
}

func TestGateNormalDrainRemovesReadyDescendant(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "descendant.pid")
	phase := Phase{Name: "drain", Argv: []string{"sh", "-c", readyDescendantScript, "sh", path}}
	result := runPhase(context.Background(), dir, phase, io.Discard, io.Discard)
	if result.Code != 0 || result.StartErr != nil {
		t.Fatalf("normal phase result = %#v, want green", result)
	}
	if pid := processtest.WaitPID(t, path); !errors.Is(syscall.Kill(pid, 0), syscall.ESRCH) {
		t.Fatalf("descendant %d survived normal completion", pid)
	}
}
