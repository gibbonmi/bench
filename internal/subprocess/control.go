package subprocess

import (
	"errors"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// ErrNoStart is the start error of a lifetime whose cancellation preceded its launch.
var ErrNoStart = errors.New("subprocess: canceled before start")

// timer holds one running window.
type timer struct {
	c    <-chan time.Time
	stop func() bool
}

// begin starts w, and a nil w expires at once.
func (t *timer) begin(w Window) {
	t.end()
	if w == nil {
		fired := make(chan time.Time, 1)
		fired <- time.Time{}
		t.c = fired
		return
	}
	t.c, t.stop = w()
}

func (t *timer) end() {
	if t.stop != nil {
		t.stop()
	}
	t.c, t.stop = nil, nil
}

type streamResult struct {
	p   *pipe
	err error
}

// controller is the one goroutine that owns the lifetime's state. The startup worker, the
// waiter, and the stream workers report to it, so no control path blocks on any of them.
type controller struct {
	l      *Lifetime
	cmd    *exec.Cmd
	policy Policy
	pipes  []*pipe

	startDone  chan error
	waitDone   chan error
	streamDone chan streamResult

	startPending, started, waitPending bool
	startErr, waitErr                  error
	pid                                int
	state                              *os.ProcessState
	// The startup worker writes reg and identityErr before it reports on startDone.
	reg         *protection.Registration
	identityErr error

	request                     Cancellation
	stopping, forced, needGroup bool
	closed, expired, returned   bool
	group                       protection.Observation
	signalErr                   error

	grace, final, delay, poll timer
}

// newController prepares the command and its owned pipes. It keeps the caller's platform
// attributes and adds only the group setup.
func newController(l *Lifetime) *controller {
	c := &controller{
		l: l, cmd: l.cmd, policy: l.policy,
		startDone:  make(chan error, 1),
		waitDone:   make(chan error, 1),
		streamDone: make(chan streamResult),
	}
	if c.cmd.SysProcAttr == nil {
		c.cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	c.cmd.SysProcAttr.Setpgid = true
	pipes, err := openStreams(c.cmd, l.streams)
	c.pipes = pipes
	if err != nil {
		c.startDone <- err
	}
	return c
}

func (c *controller) run() {
	c.startPending = true
	if len(c.startDone) == 0 {
		go c.launch(c.l.start)
	}
	for !c.returned || !c.settled() {
		if !c.returned && c.resolved() {
			c.publish()
			continue
		}
		select {
		case err := <-c.startDone:
			c.onStart(err)
		case err := <-c.waitDone:
			c.onWait(err)
		case r := <-c.streamDone:
			r.p.finish(r.err)
		case <-c.l.notify:
			c.onCancel()
		case <-c.grace.c:
			c.grace.end()
			c.force()
		case <-c.final.c:
			c.final.end()
			c.expired = true
			if !c.returned {
				c.publish()
			}
			// Owned closure releases consumers that an escaped descendant still feeds.
			c.closeOwned()
		case <-c.delay.c:
			c.delay.end()
			c.closeOwned()
			c.beginFinal()
		case <-c.poll.c:
			c.poll.end()
			c.observe()
		}
	}
	for _, t := range []*timer{&c.grace, &c.final, &c.delay, &c.poll} {
		t.end()
	}
	if c.needGroup && c.group.State != protection.GroupAbsent {
		c.observe()
		c.poll.end()
	}
	c.l.complete.outcome = c.finished()
	close(c.l.complete.done)
}

// launch is the one startup worker. It never holds the lifetime lock across Start.
func (c *controller) launch(start func(*exec.Cmd) error) {
	c.startDone <- c.register(start)
}

// register publishes pending protection before Start and the actual identity after it. A
// proved no-start closes the registration, so a canceled or failed launch keeps no user.
func (c *controller) register(start func(*exec.Cmd) error) error {
	if c.l.pendingRequest().Cause != CauseNone {
		return ErrNoStart
	}
	if resources := c.l.registration.Resources; len(resources) > 0 {
		reg, err := c.l.registration.Protocol.Prepare(resources)
		if c.reg = reg; err != nil {
			return err
		}
		c.cmd.Env = protection.Environ(c.cmd.Environ(), reg.Resources())
	}
	err := ErrNoStart
	if c.l.pendingRequest().Cause == CauseNone {
		err = start(c.cmd)
	}
	switch {
	case c.reg == nil:
		return err
	case err != nil:
		// The disposition keeps any persistence error of the no-start record.
		_ = c.reg.Finish(protection.Terminal{NoStart: true})
		return err
	}
	c.identityErr = c.reg.Started(c.cmd.Process.Pid, c.cmd.Process.Pid)
	return nil
}

func (c *controller) onStart(err error) {
	c.startPending = false
	if errors.Is(err, ErrNoStart) {
		c.stopping, c.request = true, c.l.pendingRequest()
	}
	if err != nil {
		c.startErr = err
		for _, p := range c.pipes {
			_ = p.child.Close()
			_ = p.owner.Close()
		}
		c.pipes = nil
		return
	}
	c.started, c.pid = true, c.cmd.Process.Pid
	for _, p := range c.pipes {
		_ = p.child.Close()
		go func() { c.streamDone <- streamResult{p, p.work()} }()
	}
	c.waitPending = true
	go func() { c.waitDone <- c.l.wait(c.cmd) }()
	if c.identityErr != nil && !c.stopping {
		// The records cannot name this child, so the live owner stops it with its own handle.
		c.stopping, c.needGroup = true, true
		c.force()
		return
	}
	if !c.stopping {
		return
	}
	c.needGroup = true
	if c.forced {
		// The remembered request already passed its grace, so the late child gets KILL.
		c.signal(syscall.SIGKILL)
		return
	}
	c.signal(c.request.Signal)
}

func (c *controller) onWait(err error) {
	c.waitPending, c.waitErr, c.state = false, err, c.cmd.ProcessState
	if !c.stopping && c.policy.Normal == Drain {
		c.needGroup = true
		c.signal(syscall.SIGKILL)
		c.beginFinal()
	}
	if !c.stopping && c.policy.Normal == Preserve && c.policy.OutputDelay != nil && !c.pipesDone() {
		c.delay.begin(c.policy.OutputDelay)
	}
	if c.needGroup {
		c.observe()
	}
}

func (c *controller) onCancel() {
	request := c.l.pendingRequest()
	if c.stopping || request.Cause == CauseNone {
		return
	}
	c.stopping, c.request = true, request
	c.delay.end()
	if c.started {
		c.needGroup = true
	}
	if request.Signal == syscall.SIGKILL || c.policy.Grace == nil {
		c.force()
		return
	}
	if c.started {
		c.signal(request.Signal)
	}
	c.grace.begin(c.policy.Grace)
}

// force sends KILL to a started group and starts the one final window.
func (c *controller) force() {
	if c.forced {
		return
	}
	c.forced = true
	c.grace.end()
	if c.started {
		c.signal(syscall.SIGKILL)
	}
	c.beginFinal()
	if c.needGroup {
		c.observe()
	}
}

// beginFinal starts the final window once. Every later obligation shares that deadline.
func (c *controller) beginFinal() {
	if c.final.c != nil || c.expired {
		return
	}
	c.final.begin(c.policy.Final)
}

// signal sends sig to the group. An answer that proves absence resolves the group; any
// other error stays evidence.
func (c *controller) signal(sig syscall.Signal) {
	err := c.l.kill(-c.pid, sig)
	switch answer := protection.Classify(err); {
	case answer.State == protection.GroupAbsent:
		c.group = answer
	case err != nil && c.signalErr == nil:
		c.signalErr = err
	}
}

// observe probes the group once its leader is reaped; an unreaped leader keeps it present.
func (c *controller) observe() {
	if !c.started || c.waitPending {
		return
	}
	c.group = protection.ObserveGroup(c.pid)
	if c.group.State != protection.GroupAbsent && c.policy.Poll != nil && c.poll.c == nil {
		c.poll.begin(c.policy.Poll)
	}
}

func (c *controller) closeOwned() {
	c.closed = true
	for _, p := range c.pipes {
		if !p.done {
			_ = p.owner.Close()
		}
	}
}

func (c *controller) pipesDone() bool {
	for _, p := range c.pipes {
		if !p.done {
			return false
		}
	}
	return true
}

// settled reports that no worker is still running.
func (c *controller) settled() bool {
	return !c.startPending && !c.waitPending && c.pipesDone()
}

// resolved reports that every obligation, including a required group absence, is met.
func (c *controller) resolved() bool {
	return c.settled() && (!c.needGroup || c.group.State == protection.GroupAbsent)
}

func (c *controller) publish() {
	c.returned = true
	c.l.outcome = c.finished()
	close(c.l.returned)
}

// finished publishes a started launch's terminal facts to its registration and returns the
// snapshot with the resulting disposition.
func (c *controller) finished() Outcome {
	o := c.snapshot()
	if c.started && c.reg != nil {
		// The disposition keeps any persistence error.
		_ = c.reg.Finish(terminal(o))
		o.Protection = c.reg.Disposition()
	}
	return o
}
