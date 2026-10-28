package backgroundstate

import (
	"context"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/subprocess"
)

// The windows of the queue steps. A queue holder holds the lock at most for the start of one
// grade, which the teardown window of a process bounds. A stopped grade cancels its oracle
// group and then tears down its own resources.
const (
	LockWindow = 2 * bounds.ProcessFinalWindow
	StopWindow = bounds.ProcessGroupCancelGrace + bounds.ProcessFinalWindow
)

// Stop asks the owner of claim to end and returns only after that owner has exited, or when
// ctx ends. The signal goes through the lifetime owner in internal/subprocess. A claim that is
// not live has no owner to stop, so Stop sends nothing: a pid of an earlier boot can name
// another process.
func Stop(ctx context.Context, claim Claim) error {
	live, err := claim.Live()
	if err != nil || !live {
		return err
	}
	return subprocess.Stop(ctx, claim.Owner, bounds.ProcessPoll)
}

// StopAll drops every queued grade of the machine queue under the Bench home home through
// replacer, and then stops its running grade and returns only after that grade has exited. The
// lock wait ends with the lock window and the stop wait with the stop window. A queue that holds
// no grade needs no lock. The stopped grade frees its own slot when it ends, and the next queue
// step frees the slot of a grade whose owner died.
func StopAll(replacer durablefile.Replacer, home string) error {
	queue, err := ReadQueue(home)
	if err != nil || (queue.Running == nil && len(queue.Queued) == 0) {
		return err
	}
	lockCtx, cancelLock := bounds.Context(context.Background(), bounds.FixedWindow(LockWindow))
	defer cancelLock()
	var running *Slot
	err = UpdateQueue(lockCtx, replacer, home, func(queue *Queue) error {
		queue.Queued, running = nil, queue.Running
		return nil
	})
	if err != nil || running == nil {
		return err
	}
	stopCtx, cancelStop := bounds.Context(context.Background(), bounds.FixedWindow(StopWindow))
	defer cancelStop()
	return Stop(stopCtx, running.Claim())
}
