package subprocess

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// CancelSignals is the one source for the termination signals a Bench command
// traps. Every production owner registers this set rather than its own literal,
// because a command that detaches a child into its own process group
// (SysProcAttr.Setpgid) has no cleanup path but its own handler: the detached
// group never receives the terminal's SIGINT or SIGHUP, so a signal the owner
// declines to trap leaks the whole group for as long as the child runs. Session
// and harness teardown sends SIGTERM, so an owner trapping SIGINT alone is the
// leak, not the exception.
//
// SIGKILL is absent because it cannot be trapped; a SIGKILLed owner still leaks
// its group.
var CancelSignals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGHUP}

// NotifyCancel derives a context cancelled by any of CancelSignals. Callers that
// need a channel rather than a context spread CancelSignals into signal.Notify
// directly, so both shapes still read the set from here.
func NotifyCancel(parent context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(parent, CancelSignals...)
}
