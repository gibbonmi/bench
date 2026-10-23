package subprocess

import (
	"context"
	"os/exec"
	"syscall"

	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// CancelRequest answers the shutdown cause and signal for a cancelled context.
type CancelRequest func(ctx context.Context) (Cause, syscall.Signal)

// Interrupt answers an interruption that sends signal, whatever cancelled the context.
func Interrupt(signal syscall.Signal) CancelRequest {
	return func(context.Context) (Cause, syscall.Signal) { return CauseInterruption, signal }
}

// RunInherited runs cmd under policy and streams as one group that registers on the inherited
// union plus local under the installed protocol. The cancellation of ctx requests the shutdown
// that cancel answers. The error is an unreadable inherited vector or a refused start; every
// later fact is in the outcome.
func RunInherited(ctx context.Context, cmd *exec.Cmd, policy Policy, streams Streams, cancel CancelRequest, local ...protection.Descriptor) (Outcome, error) {
	inherited, err := protection.Inherited()
	if err != nil {
		return Outcome{}, err
	}
	lifetime := NewLifetime(cmd, policy, streams, Registration{Protocol: protection.Installed(), Resources: append(inherited, local...)})
	if err := lifetime.Start(); err != nil {
		return Outcome{}, err
	}
	stop := context.AfterFunc(ctx, func() { lifetime.Cancel(cancel(ctx)) })
	defer stop()
	outcome, _ := lifetime.Wait()
	return outcome, nil
}
