package subprocess

import (
	"errors"
	"io"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// Window starts one caller-classified window. Its channel receives when the window
// expires, and stop releases its timer. A nil channel never expires. The duration policy
// owner builds every window, so this package holds no duration of its own.
type Window func() (expired <-chan time.Time, stop func() bool)

// NormalMode selects what the owner does with the group after its leader exits on its own.
type NormalMode uint8

const (
	// Preserve leaves the group unsignaled after a normal leader exit.
	Preserve NormalMode = iota
	// Drain kills the remaining group after a normal leader exit and observes its absence.
	Drain
)

// Policy is one caller's lifetime policy. A nil Grace sends KILL with the cancellation, and
// a nil Final expires at once. A nil OutputDelay waits for owned streams until cancellation.
// A nil Poll observes the group only when another event arrives.
type Policy struct {
	Normal      NormalMode
	Grace       Window
	Final       Window
	OutputDelay Window
	Poll        Window
}

// Cause names why a caller requested shutdown.
type Cause uint8

const (
	CauseNone Cause = iota
	CauseInterruption
	CauseDeadline
)

// Output selects one child output stream. Capture keeps the bytes in the outcome. Otherwise
// an *os.File reaches the child directly, any other Writer receives the bytes through an
// owned pipe, and a nil Writer discards them.
type Output struct {
	Writer  io.Writer
	Capture bool
}

// Streams selects the child's streams; the owner replaces any streams set on the command.
// An *os.File Stdin reaches the child directly, and any other Reader feeds an owned pipe.
// MergeStderr gives stderr the stdout endpoint, so one consumer reads both in order.
type Streams struct {
	Stdin       io.Reader
	Stdout      Output
	Stderr      Output
	MergeStderr bool
}

// Registration names the resources that a launch protects before Start: the inherited union
// plus the caller's local resources. The child receives that union in its environment.
type Registration struct {
	Protocol  *protection.Protocol
	Resources []protection.Descriptor
}

var errStartedTwice = errors.New("subprocess: lifetime already started")

// Lifetime owns one command's process group from Start until its last worker completes.
// No caller calls Wait, Run, or a context Cancel on the owned command.
type Lifetime struct {
	cmd          *exec.Cmd
	policy       Policy
	streams      Streams
	registration Registration
	// start, wait, and kill are the real calls; tests hold them to order a race or fail one.
	start, wait func(*exec.Cmd) error
	kill        func(int, syscall.Signal) error

	mu       sync.Mutex
	started  bool
	request  Cancellation
	notify   chan struct{}
	returned chan struct{}
	outcome  Outcome
	complete *Completion
}

// NewLifetime prepares cmd under policy, streams, and registration without starting it.
func NewLifetime(cmd *exec.Cmd, policy Policy, streams Streams, registration Registration) *Lifetime {
	return &Lifetime{
		cmd: cmd, policy: policy, streams: streams, registration: registration,
		start:    (*exec.Cmd).Start,
		wait:     (*exec.Cmd).Wait,
		kill:     syscall.Kill,
		notify:   make(chan struct{}, 1),
		returned: make(chan struct{}),
		complete: &Completion{done: make(chan struct{})},
	}
}

// Start queues the one startup worker and returns. Acceptance proves no launch: Wait
// reports whether the child started.
func (l *Lifetime) Start() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.started {
		return errStartedTwice
	}
	l.started = true
	c := newController(l)
	// The controller now owns the caller's streams; the handle keeps no reference to them.
	l.streams = Streams{}
	go c.run()
	return nil
}

// Cancel requests shutdown with signal. The first request wins. A request before the
// startup worker launches prevents the launch; a later one reaches the eventual child.
func (l *Lifetime) Cancel(cause Cause, signal syscall.Signal) {
	l.mu.Lock()
	if l.request.Cause == CauseNone {
		l.request = Cancellation{Cause: cause, Signal: signal}
	}
	l.mu.Unlock()
	select {
	case l.notify <- struct{}{}:
	default:
	}
}

// Wait returns the outcome at completion or at the final deadline, whichever comes first.
// The completion handle publishes a new outcome after every unfinished worker completes.
// A lifetime that never started returns a start error and no completion handle.
func (l *Lifetime) Wait() (Outcome, *Completion) {
	l.mu.Lock()
	started := l.started
	l.mu.Unlock()
	if !started {
		return Outcome{Process: Process{StartErr: errors.New("subprocess: lifetime never started")}}, nil
	}
	<-l.returned
	return l.outcome, l.complete
}

func (l *Lifetime) pendingRequest() Cancellation {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.request
}
