package bounds_test

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// TestBoundsCancel holds the runner's cancellation policy: KILL with no catchable signal,
// and the existing canceled and timeout statuses.
func TestBoundsCancel(t *testing.T) {
	for _, interrupt := range []bool{true, false} {
		dir := t.TempDir()
		cmd := exec.Command("sh", "-c", `trap "echo TERM >> signals" TERM; echo $$ > ready; while :; do sleep 0.05; done`)
		cmd.Dir = dir
		parent, cancel := context.WithCancel(context.Background())
		limit := bounds.TestDeadline(0)
		if !interrupt {
			limit = 300 * time.Millisecond
		}
		done := make(chan bounds.ProcessResult, 1)
		go func() { done <- bounds.Run(parent, limit, cmd) }()
		processtest.WaitPID(t, filepath.Join(dir, "ready"))
		if interrupt {
			cancel()
		}
		var result bounds.ProcessResult
		select {
		case result = <-done:
		case <-time.After(bounds.TestDeadline(bounds.ProcessFinalWindow)):
			t.Fatal(bounds.TestTimeoutVerdict("the canceled runner", bounds.TestDeadline(bounds.ProcessFinalWindow)))
		}
		cancel()
		want, wantErr := bounds.ProcessTimeout, context.DeadlineExceeded
		if interrupt {
			want, wantErr = bounds.ProcessCanceled, context.Canceled
		}
		signals, _ := os.ReadFile(filepath.Join(dir, "signals"))
		if result.Status != want || !errors.Is(result.Err, wantErr) || len(signals) != 0 {
			t.Fatalf("interrupt=%t: status=%s err=%v catchable signals=%q, want %s with immediate KILL", interrupt, result.Status, result.Err, signals, want)
		}
	}
}

// TestBoundsMalformedResources holds that a malformed inherited resource vector refuses the run
// before any child starts.
func TestBoundsMalformedResources(t *testing.T) {
	t.Setenv(protection.Env, "{not a vector")
	cmd := exec.Command("true")
	result := bounds.Run(context.Background(), bounds.TestDeadline(bounds.ProcessFinalWindow), cmd)
	if result.Status != bounds.ProcessStart || result.Err == nil || cmd.Process != nil {
		t.Fatalf("a malformed %s started %v: status=%s err=%v", protection.Env, cmd.Process, result.Status, result.Err)
	}
}

// TestBoundsFixedShutdown runs under the unbounded-verdict switch. An escaped holder keeps
// the output pipe open, so only a fixed final window lets the canceled runner return.
func TestBoundsFixedShutdown(t *testing.T) {
	t.Setenv(bounds.UnboundedWaitsEnv, "1")
	holder := processtest.NewHolder(t)
	cmd := exec.Command("sh", "-c", `"$0" & exec sleep 30`, holder.Command)
	cmd.Env = append(os.Environ(), holder.Env)
	parent, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan bounds.ProcessResult, 1)
	go func() { done <- bounds.Run(parent, bounds.VerdictWindow(bounds.ProviderTimeout), cmd) }()
	holder.PID(t)
	cancel()
	select {
	case result := <-done:
		if result.Status != bounds.ProcessCanceled {
			t.Fatalf("canceled runner status = %s: %v", result.Status, result.Err)
		}
	case <-time.After(bounds.TestDeadline(bounds.ProcessFinalWindow)):
		t.Fatal(bounds.TestTimeoutVerdict("teardown under the unbounded-verdict switch", bounds.TestDeadline(bounds.ProcessFinalWindow)))
	}
}
