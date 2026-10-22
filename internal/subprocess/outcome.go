package subprocess

import (
	"syscall"

	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// Outcome is one immutable lifetime snapshot. Its parts stay separate, so a cleanup
// failure cannot erase the child result and a child exit cannot certify its streams.
type Outcome struct {
	Process      Process
	Cancellation Cancellation
	Streams      StreamSet
	Cleanup      Cleanup
}

// Process holds the launch and wait facts. Code is the raw exit code, or -1 when Signal
// ended the child. Waited is false while the one waiter is still pending.
type Process struct {
	Started  bool
	StartErr error
	PID      int
	PGID     int
	Waited   bool
	Code     int
	Signal   syscall.Signal
	WaitErr  error
}

// Cancellation is the first shutdown request, or CauseNone when no caller asked.
type Cancellation struct {
	Cause  Cause
	Signal syscall.Signal
}

// StreamStatus is the terminal state of one owned stream. StreamNone marks a stream the
// owner holds no endpoint for: a caller file, a discarded stream, or a merged stderr.
type StreamStatus string

const (
	StreamNone     StreamStatus = "none"
	StreamComplete StreamStatus = "complete"
	StreamPending  StreamStatus = "pending"
	StreamFailed   StreamStatus = "failed"
	StreamClosed   StreamStatus = "forcibly-closed"
)

// Stream is one stream's status, its error, and its captured bytes. Data stays nil until
// the stream's worker has finished.
type Stream struct {
	Status StreamStatus
	Err    error
	Data   []byte
}

// StreamSet holds the three child streams.
type StreamSet struct {
	Stdin, Stdout, Stderr Stream
}

// CleanupStatus is the teardown verdict. CleanupNotRequested means the lifetime ended
// without a shutdown request, a drain, or a forced stream closure.
type CleanupStatus string

const (
	CleanupNotRequested CleanupStatus = "not-requested"
	CleanupComplete     CleanupStatus = "complete"
	CleanupIncomplete   CleanupStatus = "incomplete"
)

// Obligation names one teardown duty that was still unresolved at the snapshot.
type Obligation string

const (
	ObligationStart  Obligation = "start"
	ObligationWaiter Obligation = "waiter"
	ObligationGroup  Obligation = "group"
	ObligationStdin  Obligation = "stdin"
	ObligationStdout Obligation = "stdout"
	ObligationStderr Obligation = "stderr"
)

// Cleanup is the teardown verdict, each unresolved obligation once, the last group
// observation, and the first signal error that does not prove the group absent.
type Cleanup struct {
	Status      CleanupStatus
	Obligations []Obligation
	Group       protection.Observation
	SignalErr   error
}

// snapshot copies the controller state into an immutable outcome.
func (c *controller) snapshot() Outcome {
	o := Outcome{Cancellation: c.request}
	o.Process = Process{
		Started: c.started, StartErr: c.startErr, PID: c.pid, PGID: c.pid,
		Waited: c.started && !c.waitPending, WaitErr: c.waitErr,
	}
	if c.state != nil {
		o.Process.Code = c.state.ExitCode()
		if status, ok := c.state.Sys().(syscall.WaitStatus); ok && status.Signaled() {
			o.Process.Signal = status.Signal()
		}
	}
	o.Streams = StreamSet{Stdin: Stream{Status: StreamNone}, Stdout: Stream{Status: StreamNone}, Stderr: Stream{Status: StreamNone}}
	slots := map[Obligation]*Stream{ObligationStdin: &o.Streams.Stdin, ObligationStdout: &o.Streams.Stdout, ObligationStderr: &o.Streams.Stderr}
	var open []Obligation
	if c.startPending {
		open = append(open, ObligationStart)
	}
	if c.started && c.waitPending {
		open = append(open, ObligationWaiter)
	}
	if c.needGroup && c.group.State != protection.GroupAbsent {
		open = append(open, ObligationGroup)
	}
	for _, p := range c.pipes {
		*slots[p.obligation] = p.stream(c.closed)
		if !p.done {
			open = append(open, p.obligation)
		}
	}
	o.Cleanup = Cleanup{Status: CleanupNotRequested, Obligations: open, Group: c.group, SignalErr: c.signalErr}
	switch {
	case len(open) > 0:
		o.Cleanup.Status = CleanupIncomplete
	case c.stopping || c.needGroup || c.closed:
		o.Cleanup.Status = CleanupComplete
	}
	return o
}

// Completion publishes the outcome after every unfinished worker completes.
type Completion struct {
	done    chan struct{}
	outcome Outcome
}

// Done closes when the completion outcome is ready.
func (c *Completion) Done() <-chan struct{} { return c.done }

// Outcome blocks until Done and returns the completion snapshot.
func (c *Completion) Outcome() Outcome {
	<-c.done
	return c.outcome
}
