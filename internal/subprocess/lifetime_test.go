package subprocess_test

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
)

func TestLifetimePreCanceled(t *testing.T) {
	cmd := script(t, "touch started", nil)
	l := subprocess.NewLifetime(cmd, policy(0, final, 0), subprocess.Streams{})
	var launches int
	subprocess.Hold(l, func(real func(*exec.Cmd) error, cmd *exec.Cmd) error { launches++; return real(cmd) }, nil)
	l.Cancel(subprocess.CauseInterruption, syscall.SIGKILL)
	o, _ := waitWithin(t, started(t, l), final)
	if o.Process.Started || !errors.Is(o.Process.StartErr, subprocess.ErrNoStart) || launches != 0 {
		t.Fatalf("pre-canceled lifetime launched: process=%+v launches=%d", o.Process, launches)
	}
	if _, err := os.Stat(filepath.Join(cmd.Dir, "started")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("start marker exists after a pre-start cancellation: %v", err)
	}
}

func TestLifetimeStartingCancel(t *testing.T) {
	l := subprocess.NewLifetime(script(t, `trap "" TERM; exec sleep 30`, nil), policy(grace, final, 0), subprocess.Streams{})
	entered, release := holdStart(t, l)
	started(t, l)
	<-entered
	l.Cancel(subprocess.CauseInterruption, syscall.SIGTERM)
	release()
	o, _ := waitWithin(t, l, grace+final)
	if !o.Process.Started || !o.Process.Waited || o.Process.Signal == 0 || alive(o.Process.PID) {
		t.Fatalf("cancellation during Start missed the eventual child: %+v", o.Process)
	}
}

func TestLifetimeBlockedStart(t *testing.T) {
	l := subprocess.NewLifetime(script(t, "exec sleep 30", nil), policy(0, final, 0), subprocess.Streams{})
	entered, _ := holdStart(t, l)
	started(t, l)
	<-entered
	l.Cancel(subprocess.CauseInterruption, syscall.SIGKILL)
	o, _ := waitWithin(t, l, final)
	if o.Process.Started {
		t.Fatalf("held start reported a launch: %+v", o.Process)
	}
	requireIncomplete(t, o, subprocess.ObligationStart)
}

// TestLifetimeLateStart releases the held Start after the TERM grace and the final window
// expired. Only an immediate KILL ends the TERM-resistant child, because no grace remains.
func TestLifetimeLateStart(t *testing.T) {
	l := subprocess.NewLifetime(script(t, `trap "" TERM; exec sleep 30`, nil), policy(grace, final, 0), subprocess.Streams{})
	entered, release := holdStart(t, l)
	waits := countWaits(l)
	started(t, l)
	<-entered
	l.Cancel(subprocess.CauseInterruption, syscall.SIGTERM)
	_, completion := waitWithin(t, l, grace+final)
	release()
	o := completeWithin(t, completion, final)
	if !o.Process.Waited || o.Process.Signal != syscall.SIGKILL || waits.Load() != 1 || alive(o.Process.PID) {
		t.Fatalf("late start escaped shutdown: process=%+v waits=%d", o.Process, waits.Load())
	}
}

func TestLifetimeSpawnFailure(t *testing.T) {
	l := subprocess.NewLifetime(exec.Command("/definitely/missing/bench-child"), policy(0, final, 0), subprocess.Streams{})
	waits := countWaits(l)
	o, _ := waitWithin(t, started(t, l), final)
	if o.Process.Started || o.Process.StartErr == nil || waits.Load() != 0 {
		t.Fatalf("failed start: process=%+v waits=%d", o.Process, waits.Load())
	}
}

func TestLifetimeSingleWaiter(t *testing.T) {
	for range 20 {
		l := subprocess.NewLifetime(script(t, "exit 0", nil), policy(0, final, 0), subprocess.Streams{})
		var waits atomic.Int32
		subprocess.Hold(l, nil, func(real func(*exec.Cmd) error, cmd *exec.Cmd) error {
			waits.Add(1)
			go l.Cancel(subprocess.CauseInterruption, syscall.SIGKILL)
			return real(cmd)
		})
		_, completion := waitWithin(t, started(t, l), final)
		completeWithin(t, completion, final)
		if waits.Load() != 1 {
			t.Fatalf("one launch created %d waiters", waits.Load())
		}
	}
}

func TestLifetimeNonzero(t *testing.T) {
	streams := subprocess.Streams{Stdout: capture, Stderr: capture}
	l := subprocess.NewLifetime(script(t, "printf out; printf err >&2; exit 7", nil), policy(0, final, 0), streams)
	o, _ := waitWithin(t, started(t, l), final)
	var exit *exec.ExitError
	if !o.Process.Waited || o.Process.Code != 7 || !errors.As(o.Process.WaitErr, &exit) || o.Cleanup.Status != subprocess.CleanupNotRequested {
		t.Fatalf("exit seven projected as %+v cleanup=%+v", o.Process, o.Cleanup)
	}
	if string(o.Streams.Stdout.Data) != "out" || string(o.Streams.Stderr.Data) != "err" || o.Streams.Stdout.Status != subprocess.StreamComplete {
		t.Fatalf("exit seven streams = %+v", o.Streams)
	}
}

func TestLifetimeExitAndIncomplete(t *testing.T) {
	holder := processtest.NewHolder(t)
	cmd := script(t, `trap 'exit 7' TERM; "$0" & while :; do sleep 0.05; done`, holder)
	l := started(t, subprocess.NewLifetime(cmd, policy(grace, final, 0), subprocess.Streams{Stdout: capture}))
	holder.PID(t)
	l.Cancel(subprocess.CauseInterruption, syscall.SIGTERM)
	o, _ := waitWithin(t, l, grace+final)
	if !o.Process.Waited || o.Process.Code != 7 {
		t.Fatalf("cleanup failure replaced the observed exit: %+v", o.Process)
	}
	requireIncomplete(t, o, subprocess.ObligationStdout)
}

func TestLifetimeDelayedCapture(t *testing.T) {
	holder := processtest.NewHolder(t)
	cmd := script(t, `printf early; "$0" & exit 0`, holder)
	l := started(t, subprocess.NewLifetime(cmd, policy(0, final, 0), subprocess.Streams{Stdout: capture}))
	holder.PID(t)
	l.Cancel(subprocess.CauseInterruption, syscall.SIGKILL)
	o, completion := waitWithin(t, l, final)
	if o.Streams.Stdout.Status != subprocess.StreamPending || o.Streams.Stdout.Data != nil {
		t.Fatalf("pending capture exposed its buffer: %+v", o.Streams.Stdout)
	}
	done := completeWithin(t, completion, final)
	if done.Streams.Stdout.Status != subprocess.StreamClosed || string(done.Streams.Stdout.Data) != "early" {
		t.Fatalf("completed capture = %+v", done.Streams.Stdout)
	}
}

func TestLifetimeCallerFile(t *testing.T) {
	cmd := script(t, "printf child", nil)
	file, err := os.Create(filepath.Join(cmd.Dir, "caller.out"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	l := subprocess.NewLifetime(cmd, policy(0, final, 0), subprocess.Streams{Stdout: subprocess.Output{Writer: file}})
	o, _ := waitWithin(t, started(t, l), final)
	if _, err := file.WriteString("caller"); err != nil || o.Streams.Stdout.Status != subprocess.StreamNone {
		t.Fatalf("caller file after completion: write=%v stream=%+v", err, o.Streams.Stdout)
	}
	if data, _ := os.ReadFile(file.Name()); string(data) != "childcaller" {
		t.Fatalf("caller file holds %q", data)
	}
}

func TestLifetimeForcedClose(t *testing.T) {
	holder := processtest.NewHolder(t)
	cmd := script(t, `printf early; "$0" & exit 0`, holder)
	l := started(t, subprocess.NewLifetime(cmd, policy(0, final, grace), subprocess.Streams{Stdout: capture}))
	o, _ := waitWithin(t, l, grace+final)
	if o.Streams.Stdout.Status != subprocess.StreamClosed || string(o.Streams.Stdout.Data) != "early" {
		t.Fatalf("forced closure projected %+v", o.Streams.Stdout)
	}
	// The owner closed a stream itself, so the cleanup was requested.
	if o.Cleanup.Status != subprocess.CleanupComplete {
		t.Fatalf("forced closure cleanup = %+v, want complete", o.Cleanup)
	}
}

func TestLifetimeBlockedWriter(t *testing.T) {
	streams := subprocess.Streams{Stdout: subprocess.Output{Writer: newBlockingWriter(t)}}
	l := started(t, subprocess.NewLifetime(script(t, "printf x", nil), policy(0, final, grace), streams))
	o, _ := waitWithin(t, l, grace+final)
	if o.Streams.Stdout.Status != subprocess.StreamPending {
		t.Fatalf("blocked writer stream = %+v", o.Streams.Stdout)
	}
	requireIncomplete(t, o, subprocess.ObligationStdout)
}

func TestLifetimeResistantGroup(t *testing.T) {
	cmd := script(t, `trap "" TERM; sleep 30 & echo $! > descendant; wait`, nil)
	l := started(t, subprocess.NewLifetime(cmd, policy(grace, final, 0), subprocess.Streams{}))
	descendant := processtest.WaitPID(t, filepath.Join(cmd.Dir, "descendant"))
	l.Cancel(subprocess.CauseInterruption, syscall.SIGTERM)
	o, _ := waitWithin(t, l, grace+final)
	if o.Cleanup.Status != subprocess.CleanupComplete || o.Process.Signal != syscall.SIGKILL || alive(descendant) {
		t.Fatalf("TERM-resistant group after grace: process=%+v cleanup=%+v", o.Process, o.Cleanup)
	}
}

func TestLifetimePendingWaiter(t *testing.T) {
	l := subprocess.NewLifetime(script(t, "exec sleep 30", nil), policy(0, final, 0), subprocess.Streams{})
	waiting, release := make(chan struct{}), make(chan struct{})
	subprocess.Hold(l, nil, func(real func(*exec.Cmd) error, cmd *exec.Cmd) error {
		close(waiting)
		<-release
		return real(cmd)
	})
	t.Cleanup(func() { close(release) })
	started(t, l)
	<-waiting
	l.Cancel(subprocess.CauseInterruption, syscall.SIGKILL)
	o, _ := waitWithin(t, l, final)
	if o.Process.Waited {
		t.Fatalf("held waiter reported a wait: %+v", o.Process)
	}
	// The unreaped leader keeps its group present, so the group stays unresolved too.
	requireIncomplete(t, o, subprocess.ObligationWaiter, subprocess.ObligationGroup)
}

func TestLifetimeInheritedDescriptor(t *testing.T) {
	holder := processtest.NewHolder(t)
	streams := subprocess.Streams{Stdout: subprocess.Output{Writer: io.Discard}}
	l := started(t, subprocess.NewLifetime(script(t, `"$0" & exec sleep 30`, holder), policy(0, final, 0), streams))
	pid := holder.PID(t)
	l.Cancel(subprocess.CauseDeadline, syscall.SIGKILL)
	o, completion := waitWithin(t, l, final)
	requireIncomplete(t, o, subprocess.ObligationStdout)
	done := completeWithin(t, completion, final)
	if done.Streams.Stdout.Status != subprocess.StreamClosed || !alive(pid) {
		t.Fatalf("decoder after owned closure = %+v, holder alive=%t", done.Streams.Stdout, alive(pid))
	}
}

func TestLifetimeNormalPreserve(t *testing.T) {
	cmd := script(t, "sleep 30 >/dev/null 2>&1 & echo $! > descendant", nil)
	l := started(t, subprocess.NewLifetime(cmd, policy(0, final, 0), subprocess.Streams{Stdout: capture}))
	o, _ := waitWithin(t, l, final)
	descendant := processtest.WaitPID(t, filepath.Join(cmd.Dir, "descendant"))
	if o.Cleanup.Status != subprocess.CleanupNotRequested || !alive(descendant) {
		t.Fatalf("preserve mode signaled a surviving group: cleanup=%+v alive=%t", o.Cleanup, alive(descendant))
	}
}

// TestLifetimeSingleFinalWindow fires the one final window only after the owner took the
// waiter result, so the outcome names every obligation that remains after it.
func TestLifetimeSingleFinalWindow(t *testing.T) {
	holder := processtest.NewHolder(t)
	fire, starts := make(chan time.Time, 1), 0
	polled := make(chan struct{}, 1)
	p := policy(0, final, 0)
	p.Final = func() (<-chan time.Time, func() bool) { starts++; return fire, func() bool { return true } }
	// An unreaped member keeps the group present, so the owner starts a poll only after it takes the waiter result.
	p.Poll = func() (<-chan time.Time, func() bool) {
		select {
		case polled <- struct{}{}:
		default:
		}
		return make(chan time.Time), func() bool { return true }
	}
	cmd := script(t, `"$0" & echo $$ > leader; exec sleep 30`, holder)
	l := started(t, subprocess.NewLifetime(cmd, p, subprocess.Streams{Stdout: capture}))
	holder.PID(t)
	joinGroup(t, processtest.WaitPID(t, filepath.Join(cmd.Dir, "leader")))
	l.Cancel(subprocess.CauseDeadline, syscall.SIGKILL)
	select {
	case <-polled:
	case <-time.After(bounds.TestDeadline(final)):
		t.Fatal(bounds.TestTimeoutVerdict("the owner to take the waiter result", bounds.TestDeadline(final)))
	}
	fire <- time.Time{}
	o, _ := waitWithin(t, l, final)
	if starts != 1 {
		t.Fatalf("teardown started %d final windows, want one shared window", starts)
	}
	// No waiter obligation remains: the owner took the waiter result before the window expired.
	requireIncomplete(t, o, subprocess.ObligationGroup, subprocess.ObligationStdout)
}

// joinGroup adds a member to group pgid. The test reaps it only at cleanup, so after KILL
// its zombie keeps the group present.
func joinGroup(t *testing.T, pgid int) {
	t.Helper()
	member := exec.Command("sleep", "30")
	member.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: pgid}
	if err := member.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = member.Process.Kill()
		_ = member.Wait()
	})
}

func TestLifetimeWorkerOwnership(t *testing.T) {
	collected := make(chan struct{})
	sentinel := &closeableSentinel{}
	runtime.AddCleanup(sentinel, func(done chan struct{}) { close(done) }, collected)
	streams := subprocess.Streams{Stdout: subprocess.Output{Writer: newBlockingWriter(t)}, Stderr: subprocess.Output{Writer: sentinel}}
	l := started(t, subprocess.NewLifetime(script(t, "printf x; printf y >&2", nil), policy(0, final, grace), streams))
	o, _ := waitWithin(t, l, grace+final)
	if o.Streams.Stdout.Status != subprocess.StreamPending || o.Streams.Stderr.Status != subprocess.StreamComplete {
		t.Fatalf("streams = %+v", o.Streams)
	}
	sentinel, l = nil, nil
	for deadline := time.Now().Add(bounds.TestDeadline(0)); time.Now().Before(deadline); {
		runtime.GC()
		select {
		case <-collected:
			return
		default:
		}
		<-time.After(10 * time.Millisecond)
	}
	t.Fatal("the unfinished owner retains a completed caller writer")
}

func TestLifetimeBlockedInput(t *testing.T) {
	input, feed := io.Pipe()
	t.Cleanup(func() { _ = feed.Close() })
	l := started(t, subprocess.NewLifetime(script(t, "exit 0", nil), policy(0, final, grace), subprocess.Streams{Stdin: input}))
	o, _ := waitWithin(t, l, grace+final)
	if !o.Process.Waited || o.Streams.Stdin.Status != subprocess.StreamPending {
		t.Fatalf("blocked stdin: process=%+v stdin=%+v", o.Process, o.Streams.Stdin)
	}
	requireIncomplete(t, o, subprocess.ObligationStdin)
}

// closeableSentinel is an unrelated caller writer whose collection proves release.
type closeableSentinel struct{ written []byte }

func (s *closeableSentinel) Write(p []byte) (int, error) {
	s.written = append(s.written, p...)
	return len(p), nil
}

func (s *closeableSentinel) Close() error { return nil }
