package subprocess_test

import (
	"context"
	"errors"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
)

// poll is the observation window of each stop in these tests.
var poll = bounds.ProcessPoll

// TestStopRefusesAPidThatNamesNoOneProcess guards the signal target. Zero and one select the
// caller's own group and every process, and a negative pid selects a group, so neither Alive
// nor Stop may send to them.
func TestStopRefusesAPidThatNamesNoOneProcess(t *testing.T) {
	for _, pid := range []int{-7, 0, 1} {
		if _, err := subprocess.Alive(pid); err == nil {
			t.Errorf("Alive(%d) answered with no error", pid)
		}
		if err := subprocess.Stop(context.Background(), pid, poll); err == nil {
			t.Errorf("Stop(%d) answered with no error", pid)
		}
	}
}

// TestStopEndsWithItsContext proves that Stop is bounded by its context. The process ignores
// SIGTERM, so Stop returns the context error while the process still lives.
func TestStopEndsWithItsContext(t *testing.T) {
	ready := filepath.Join(t.TempDir(), "ready")
	cmd := exec.Command("sh", "-c", `trap "" TERM; echo $$ > "$1.partial" && mv "$1.partial" "$1"; exec sleep 300`, "sh", ready)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	go func() { _ = cmd.Wait(); close(done) }()
	pid := processtest.WaitPID(t, ready)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	if err := subprocess.Stop(ctx, pid, poll); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Stop of a process that ignores SIGTERM = %v, want the context deadline", err)
	}
	if alive, err := subprocess.Alive(pid); err != nil || !alive {
		t.Fatalf("Alive after the bounded stop = %v, %v; want a live process", alive, err)
	}
	select {
	case <-done:
		t.Fatal("the process exited although it ignores SIGTERM")
	default:
	}
}
