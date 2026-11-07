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
// union plus the resources of local, under the protocol of local. A nil protocol means the
// installed protocol. The cancellation of ctx requests the shutdown that cancel answers. The
// error is an unreadable inherited vector or a refused start; every later fact is in the outcome.
func RunInherited(ctx context.Context, cmd *exec.Cmd, policy Policy, streams Streams, cancel CancelRequest, local Registration) (Outcome, error) {
	inherited, err := protection.Inherited()
	if err != nil {
		return Outcome{}, err
	}
	if local.Protocol == nil {
		local.Protocol = protection.Installed()
	}
	lifetime := NewLifetime(cmd, policy, streams, Registration{Protocol: local.Protocol, Resources: append(inherited, local.Resources...)})
	if err := lifetime.Start(); err != nil {
		return Outcome{}, err
	}
	stop := context.AfterFunc(ctx, func() { lifetime.Cancel(cancel(ctx)) })
	defer stop()
	outcome, _ := lifetime.Wait()
	return outcome, nil
}
