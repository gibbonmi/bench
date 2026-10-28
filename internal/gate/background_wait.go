package gate

// The wait of a chunk checkpoint or a landing on a background grade. A foreground chunk checkpoint
// that finds no green waits on a live claim of its own key before it takes the checkout lock, and a
// landing waits on a live claim of its own key before it takes the destination lease. Before its
// own run, each clears the background grades that would run beside it.

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

// landingWaitKey carries the landing wait of a landing gate run.
type landingWaitKey struct{}

// landingWait is the wait of a landing before its destination lease: the Bench home of the machine
// queue, the linked worktree of the landing source, and the budget that the wait left the run.
type landingWait struct {
	home, worktree string
	budget         gateBudget
}

// AwaitLanding waits, for the landing of tip from the linked worktree under the Bench home home and
// the kit kit, until no live claim holds the grade key of the tree that the landing publishes onto
// the destination at this call. The landing calls it before it takes the destination lease, so the
// wait blocks no other landing. A worktree that is not a background worktree runs no background
// grade, so its landing derives no subject, waits on nothing, and clears nothing. A key that a
// fresh green holds needs no wait, and a key that does not derive holds none. The answered context
// carries the wait to the gate run of the landing, which then spends the rest of the gate timeout
// and clears the queue for its own run. The landing calls the answered release when the run no
// longer needs the budget.
func AwaitLanding(ctx context.Context, home, kit, worktree, tip string, stderr io.Writer) (context.Context, func()) {
	if !backgroundWorktree(kit, worktree) {
		return ctx, func() {}
	}
	wait := landingWait{home: home, worktree: worktree}
	common, err := benchgit.CommonDir(worktree)
	if subject, derived := backgroundSubjectOf(worktree, tip); err == nil && derived == nil && subject.landing != nil {
		if evidence, _, inspected := subject.landing.inspect(subject.tree); inspected == nil && evidence.Grade != "" && !evidence.ReusableGreen {
			wait.budget, _, _ = awaitKey(ctx, worktree, common, evidence.Grade, stderr)
		}
	}
	return context.WithValue(ctx, landingWaitKey{}, wait), wait.budget.release
}

// landingWaitOf answers the landing wait that ctx carries, and false when ctx carries none.
func landingWaitOf(ctx context.Context) (landingWait, bool) {
	wait, waited := ctx.Value(landingWaitKey{}).(landingWait)
	return wait, waited
}

// awaitClaim waits, for a foreground chunk checkpoint, until no live claim holds the grade key of
// plan in the git common directory of root. A landing waited before its lease, so its run answers
// the budget of that wait and does not wait again.
func awaitClaim(ctx context.Context, root string, plan subject, stderr io.Writer) (gateBudget, Result, bool) {
	if wait, landing := landingWaitOf(ctx); landing {
		return wait.budget, Result{}, false
	}
	common, err := benchgit.CommonDir(root)
	if !foregroundChunk(ctx) || err != nil {
		return gateBudget{}, Result{}, false
	}
	return awaitKey(ctx, root, common, evidenceName(plan), stderr)
}

// awaitKey waits until no live claim holds key in the git common directory common of root. When
// the wait starts on the claim of an owner, the I/O class of that grade returns to best-effort. The
// wait spends the gate timeout, so it answers the budget that the run keeps. A wait that the
// timeout or an interrupt ends answers its result, as a run that ends before its oracle does. A
// claim that does not read holds no wait.
func awaitKey(ctx context.Context, root, common, key string, stderr io.Writer) (gateBudget, Result, bool) {
	promoted := 0
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
// root, or of a landing of plan from the worktree of its wait. A queued grade of its key leaves the
// queue, the queued grade of the worktree that holds another key leaves it, and the running grade
// of the worktree that holds another key stops, and the run waits for its exit. A lane pass keeps a
// grade whose key the latest tip of another worktree holds, but the run does not, because two runs
// of one worktree never write one delta record at once, and a landing grade beside the landing
// gate doubles the machine load. A step that fails prints one line, because a background grade
// never blocks a run.
func clearForRun(ctx context.Context, root string, plan subject, stderr io.Writer) {
	home := benchhome.Dir()
	if wait, landing := landingWaitOf(ctx); landing {
		home, root = wait.home, wait.worktree
	} else if !foregroundChunk(ctx) {
		return
	}
	key := evidenceName(plan)
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
