package subprocess_test

import (
	"errors"
	"os"
	"os/exec"
	"slices"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
)

const (
	grace = 100 * time.Millisecond
	final = 300 * time.Millisecond
)

// script is a real sh child in a private directory. $0 names the holder command.
func script(t *testing.T, body string, holder *processtest.Holder) *exec.Cmd {
	t.Helper()
	command, env := "", os.Environ()
	if holder != nil {
		command, env = holder.Command, append(env, holder.Env)
	}
	cmd := exec.Command("sh", "-c", body, command)
	cmd.Dir, cmd.Env = t.TempDir(), env
	return cmd
}

// policy builds real windows; a zero grace sends KILL with the cancellation.
func policy(grace, final, delay time.Duration) subprocess.Policy {
	p := subprocess.Policy{Final: bounds.ProcessWindow(final), Poll: bounds.ProcessPoll}
	if grace > 0 {
		p.Grace = bounds.ProcessWindow(grace)
	}
	if delay > 0 {
		p.OutputDelay = bounds.ProcessWindow(delay)
	}
	return p
}

var capture = subprocess.Output{Capture: true}

// started starts l and fails the test when Start refuses.
func started(t *testing.T, l *subprocess.Lifetime) *subprocess.Lifetime {
	t.Helper()
	if err := l.Start(); err != nil {
		t.Fatal(err)
	}
	return l
}

// waitWithin fails the test unless Wait returns inside the deadline derived from inner.
func waitWithin(t *testing.T, l *subprocess.Lifetime, inner time.Duration) (subprocess.Outcome, *subprocess.Completion) {
	t.Helper()
	type waited struct {
		outcome    subprocess.Outcome
		completion *subprocess.Completion
	}
	done := make(chan waited, 1)
	go func() {
		outcome, completion := l.Wait()
		done <- waited{outcome, completion}
	}()
	select {
	case got := <-done:
		return got.outcome, got.completion
	case <-time.After(bounds.TestDeadline(inner)):
		t.Fatal(bounds.TestTimeoutVerdict("the lifetime to return", bounds.TestDeadline(inner)))
	}
	return subprocess.Outcome{}, nil
}

// completeWithin fails the test unless the completion publishes inside the derived deadline.
func completeWithin(t *testing.T, c *subprocess.Completion, inner time.Duration) subprocess.Outcome {
	t.Helper()
	select {
	case <-c.Done():
		return c.Outcome()
	case <-time.After(bounds.TestDeadline(inner)):
		t.Fatal(bounds.TestTimeoutVerdict("the lifetime completion", bounds.TestDeadline(inner)))
	}
	return subprocess.Outcome{}
}

// holdStart pauses the real launch until the returned release runs. entered closes when
// the startup worker reaches the launch.
func holdStart(t *testing.T, l *subprocess.Lifetime) (entered <-chan struct{}, release func()) {
	t.Helper()
	reached, gate := make(chan struct{}), make(chan struct{})
	subprocess.Hold(l, func(real func(*exec.Cmd) error, cmd *exec.Cmd) error {
		close(reached)
		<-gate
		return real(cmd)
	}, nil)
	var once atomic.Bool
	release = func() {
		if once.CompareAndSwap(false, true) {
			close(gate)
		}
	}
	t.Cleanup(release)
	return reached, release
}

// countWaits counts the real waiters that l creates.
func countWaits(l *subprocess.Lifetime) *atomic.Int32 {
	var waits atomic.Int32
	subprocess.Hold(l, nil, func(real func(*exec.Cmd) error, cmd *exec.Cmd) error {
		waits.Add(1)
		return real(cmd)
	})
	return &waits
}

// requireIncomplete fails unless the outcome names exactly these unresolved obligations.
func requireIncomplete(t *testing.T, o subprocess.Outcome, want ...subprocess.Obligation) {
	t.Helper()
	if o.Cleanup.Status != subprocess.CleanupIncomplete || !slices.Equal(o.Cleanup.Obligations, want) {
		t.Fatalf("cleanup = %s %v, want incomplete %v", o.Cleanup.Status, o.Cleanup.Obligations, want)
	}
}

// alive reports whether pid still answers a probe.
func alive(pid int) bool { return syscall.Kill(pid, 0) == nil }

// blockingWriter blocks every write until the test ends.
type blockingWriter chan struct{}

func newBlockingWriter(t *testing.T) blockingWriter {
	w := make(blockingWriter)
	t.Cleanup(func() { close(w) })
	return w
}

func (w blockingWriter) Write(p []byte) (int, error) {
	<-w
	return 0, errors.New("blocking writer released")
}
