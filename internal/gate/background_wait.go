package gate

// The wait of a chunk checkpoint on a background grade. A foreground chunk checkpoint that finds no
// green waits on a live claim of its own key before it takes the checkout lock, and before its own
// run it clears the background grades that would run beside it.

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
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

// awaitClaim waits, for a foreground chunk checkpoint, until no live claim holds the grade key of
// plan in the git common directory of root. When the wait starts on the claim of an owner, the I/O
// class of that grade returns to best-effort. The wait counts against the gate timeout, so it
// answers the timeout that the run keeps. A wait that the timeout or an interrupt ends answers its
// result, as a run that ends before its oracle does. A claim that does not read holds no wait.
func awaitClaim(ctx context.Context, root string, plan subject, stderr io.Writer) (time.Duration, Result, bool) {
	common, err := benchgit.CommonDir(root)
	if !foregroundChunk(ctx) || err != nil {
		return gateTimeout, Result{}, false
	}
	key, start, promoted := evidenceName(plan), time.Now(), 0
	waitCtx, cancel := bounds.ContextCause(ctx, gateTimeout, errGateTimeout)
	defer cancel()
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
			return 0, Result{GateExit: code, ActionExit: code, Inspection: inspectAt(root, time.Now().UTC())}, true
		case <-expired:
		}
	}
	if promoted == 0 || gateTimeout == bounds.Unbounded {
		return gateTimeout, Result{}, false
	}
	return gateTimeout - time.Since(start), Result{}, false
}

// clearForRun clears the machine queue for the own run of a foreground chunk checkpoint of plan in
// root. A queued grade of its key leaves the queue, the queued grade of root that holds another key
// leaves it, and the running grade of root that holds another key stops, and the checkpoint waits
// for its exit. The lane pass owns that supersede policy, so a grade whose key the latest tip of
// another worktree holds stays. A queue that holds no grade of root and none of its key is not
// written. A step that fails prints one line, because a background grade never blocks a run.
func clearForRun(ctx context.Context, root string, plan subject, stderr io.Writer) {
	if !foregroundChunk(ctx) {
		return
	}
	home, key := benchhome.Dir(), evidenceName(plan)
	clears := func(entry backgroundstate.Entry) bool { return entry.Worktree == root || entry.Key == key }
	queue, err := backgroundstate.ReadQueue(home)
	if err == nil && (queue.Running == nil || queue.Running.Worktree != root) && !slices.ContainsFunc(queue.Queued, clears) {
		return
	}
	source, _ := ctx.Value(deltaSourceKey{}).(chunkSource)
	var superseded *backgroundstate.Slot
	if err == nil {
		err = updateQueue(home, func(queue *backgroundstate.Queue) error {
			reaped := reapRunning(queue)
			superseded = queue.Pass(backgroundstate.Entry{Key: key, Worktree: root, Tip: source.tip}, true)
			queue.Queued = slices.DeleteFunc(queue.Queued, func(entry backgroundstate.Entry) bool { return entry.Key == key })
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
