package gate

// The wait of a chunk checkpoint on a background grade. A foreground chunk checkpoint that finds no
// green waits on a live claim of its own key before it takes the checkout lock, and before its own
// run it clears the background grades that would run beside it.

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
	benchgit "github.com/gibbonmi/bench/internal/git"
)

// foregroundChunk reports whether ctx carries a chunk checkpoint that no background grade runs.
func foregroundChunk(ctx context.Context) bool {
	checkpoint, _ := ctx.Value(checkpointKey{}).(Checkpoint)
	return checkpoint.Chunk != "" && checkpoint.Background == nil
}

// gateBudget is the gate timeout that a wait on a background grade started before the checkout
// lock. The zero budget is a run that did not wait.
type gateBudget struct {
	waited context.Context
	cancel context.CancelFunc
}

// release ends the budget once the run no longer needs it.
func (b gateBudget) release() {
	if b.cancel != nil {
		b.cancel()
	}
}

// runContext answers the context of the oracle run under ctx: the gate timeout, which the end of
// the budget of a wait ends earlier, with the cause of that end. So a wait counts against the
// gate timeout of its run, and a run after a wait ends with exit 124 at the same deadline.
func (b gateBudget) runContext(ctx context.Context) (context.Context, context.CancelFunc) {
	run, cancel := bounds.ContextCause(ctx, gateTimeout, errGateTimeout)
	if b.waited == nil {
		return run, cancel
	}
	spent, end := context.WithCancelCause(run)
	stop := context.AfterFunc(b.waited, func() { end(context.Cause(b.waited)) })
	return spent, func() { stop(); end(nil); cancel() }
}

// awaitClaim waits, for a foreground chunk checkpoint, until no live claim holds the grade key of
// plan in the git common directory of root. When the wait starts on the claim of an owner, the I/O
// class of that grade returns to best-effort. The wait spends the gate timeout, so it answers the
// budget that the run keeps. A wait that the timeout or an interrupt ends answers its result, as a
// run that ends before its oracle does. A claim that does not read holds no wait.
func awaitClaim(ctx context.Context, root string, plan subject, stderr io.Writer) (gateBudget, Result, bool) {
	common, err := benchgit.CommonDir(root)
	if !foregroundChunk(ctx) || err != nil {
		return gateBudget{}, Result{}, false
	}
	key, promoted := evidenceName(plan), 0
	waitCtx, cancel := bounds.ContextCause(ctx, gateTimeout, errGateTimeout)
	budget := gateBudget{waited: waitCtx, cancel: cancel}
	for {
		claim, err := backgroundstate.ReadClaim(common, key)
		if err != nil || claim.Key == "" {
			break
		}
		if live, err := claim.Live(); err != nil || !live {
			break
		}
		if claim.Owner != promoted {
			promoted = claim.Owner
			if err := promoteIO(promoted); err != nil {
				fmt.Fprintf(stderr, "background grade: I/O class not promoted: %v\n", err)
			}
		}
		expired, stop := bounds.ProcessPoll()
		select {
		case <-waitCtx.Done():
			stop()
			code, _, _ := interruptedExit(ctx, waitCtx, 130, stderr)
			return budget, Result{GateExit: code, ActionExit: code, Inspection: inspectAt(root, time.Now().UTC())}, true
		case <-expired:
		}
	}
	if promoted == 0 {
		cancel()
		return gateBudget{}, Result{}, false
	}
	return budget, Result{}, false
}

// clearForRun clears the machine queue for the own run of a foreground chunk checkpoint of plan in
// root. A queued grade of its key leaves the queue, the queued grade of root that holds another key
// leaves it, and the running grade of root that holds another key stops, and the checkpoint waits
// for its exit. A lane pass keeps a grade whose key the latest tip of another worktree holds, but
// the checkpoint does not, because two runs of root never write one delta record at once. A step
// that fails prints one line, because a background grade never blocks a run.
func clearForRun(ctx context.Context, root string, plan subject, stderr io.Writer) {
	if !foregroundChunk(ctx) {
		return
	}
	home, key := benchhome.Dir(), evidenceName(plan)
	// The clear first runs on the queue as read, which it does not write: a queue that the clear
	// leaves as it is needs no lock.
	queue, err := backgroundstate.ReadQueue(home)
	if queued := len(queue.Queued); err == nil && queue.ClearForCheckpoint(key, root) == nil && len(queue.Queued) == queued {
		return
	}
	var superseded *backgroundstate.Slot
	if err == nil {
		err = updateQueue(home, func(queue *backgroundstate.Queue) error {
			reaped := reapRunning(queue)
			superseded = queue.ClearForCheckpoint(key, root)
			return reaped
		})
	}
	if err == nil && superseded != nil {
		err = errors.Join(stopGrade(*superseded), advanceQueue(home, 0))
	}
	if err != nil {
		fmt.Fprintf(stderr, "background grade: not cleared: %v\n", err)
	}
}
