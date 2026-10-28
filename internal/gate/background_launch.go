package gate

// The launch of a background grade: after a lane pass publishes a tip in a kit worktree, the
// launch passes the grade of that tip to the machine queue. The queue runs one grade at a time
// through the wrapper of its worktree as a detached session, and no launch waits for a grade.
// When a grade ends, its own process starts the oldest queued grade and then exits.

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/subprocess"
)

// The windows of the queue steps. A queue holder holds the lock at most for the start of one
// grade, which the teardown window of a process bounds. A stopped grade cancels its oracle
// group and then tears down its own resources.
const (
	queueLockWindow = 2 * bounds.ProcessFinalWindow
	gradeStopWindow = bounds.ProcessGroupCancelGrace + bounds.ProcessFinalWindow
)

// StartBackground passes the grade of tip, which a lane pass just published in the worktree
// root, to the machine queue under the Bench home home. kit is the kit that the verb read. A
// grade queues only in the worktree of an active assignment for which KitForRoot answers root
// itself, and only when root holds its executable wrapper, because the grade runs as that
// wrapper. The publication stands whatever the launch does, so a launch that fails prints one
// line to stderr and the verb keeps its exit.
func StartBackground(home, kit, root, tip string, stderr io.Writer) {
	if KitForRoot(kit, root) != root || !RealFS().Executable(backgroundWrapper(root)) {
		return
	}
	if _, assigned := intent.AssignmentForWorktree(root); !assigned {
		return
	}
	if err := passLane(home, root, tip); err != nil {
		fmt.Fprintf(stderr, "background grade: not queued: %v\n", err)
	}
}

// passLane applies the lane pass of tip in root to the queue under home. The pass queues nothing
// for a key that a fresh green or a live claim already holds. A running grade that the pass
// supersedes stops before the queue moves on. A claim that does not read retains its grade, so
// the pass queues nothing.
func passLane(home, root, tip string) error {
	key, settled := gradeKey(root, tip)
	if key != "" && !settled {
		common, err := benchgit.CommonDir(root)
		if err != nil {
			return err
		}
		claim, err := backgroundstate.ReadClaim(common, key)
		if err != nil {
			return err
		}
		if settled, err = claim.Live(); err != nil {
			return err
		}
	}
	var superseded *backgroundstate.Slot
	err := updateQueue(home, func(queue *backgroundstate.Queue) error {
		reaped := reapRunning(queue)
		superseded = queue.Pass(backgroundstate.Entry{Key: key, Worktree: root, Tip: tip}, settled)
		return errors.Join(reaped, startOldest(home, queue))
	})
	if err != nil || superseded == nil {
		return err
	}
	return errors.Join(stopGrade(*superseded), advanceQueue(home, 0))
}

// advanceQueue moves the queue under home on after the grade of process ended ends: it frees
// the slot of that grade, reaps a running grade whose owner died, and starts the oldest queued
// grade when the slot is free.
func advanceQueue(home string, ended int) error {
	return updateQueue(home, func(queue *backgroundstate.Queue) error {
		if queue.Running != nil && queue.Running.Owner == ended {
			queue.Running = nil
		}
		return errors.Join(reapRunning(queue), startOldest(home, queue))
	})
}

// updateQueue applies update to the queue under home within the lock window.
func updateQueue(home string, update func(*backgroundstate.Queue) error) error {
	ctx, cancel := bounds.Context(context.Background(), bounds.FixedWindow(queueLockWindow))
	defer cancel()
	return backgroundstate.UpdateQueue(ctx, durablefile.Native(), home, update)
}

// reapRunning frees the slot of a running grade whose owner is not live. Its worktree records
// aborted unless the grade recorded its own end. A slot whose liveness does not read retains
// its grade.
func reapRunning(queue *backgroundstate.Queue) error {
	if queue.Running == nil {
		return nil
	}
	if live, err := queue.Running.Live(); err != nil || live {
		return err
	}
	dead := *queue.Running
	queue.Running = nil
	gitdir, record, err := endedRecord(dead)
	if err != nil || gitdir == "" || record.Tip == dead.Tip {
		return err
	}
	return backgroundstate.WriteRecord(durablefile.Native(), gitdir, backgroundstate.Record{Outcome: backgroundstate.Aborted, Key: dead.Key, Tip: dead.Tip, Log: backgroundstate.LogPath(gitdir)})
}

// stopGrade stops the running grade of slot within the stop window. The grade records its own
// end before it exits, so the stop then names that end stopped. A grade that reached its verdict
// before the stop keeps that verdict.
func stopGrade(slot backgroundstate.Slot) error {
	ctx, cancel := bounds.Context(context.Background(), bounds.FixedWindow(gradeStopWindow))
	defer cancel()
	if err := backgroundstate.Stop(ctx, slot.Claim()); err != nil {
		return err
	}
	gitdir, record, err := endedRecord(slot)
	if err != nil || gitdir == "" || record.Tip != slot.Tip || record.Outcome == backgroundstate.Green || record.Outcome == backgroundstate.Red {
		return err
	}
	if record.Key == "" {
		record.Key = slot.Key
	}
	record.Outcome = backgroundstate.Stopped
	return backgroundstate.WriteRecord(durablefile.Native(), gitdir, record)
}

// endedRecord answers the git directory of the worktree of slot and its background record,
// after the grade of slot ended. A worktree that is gone answers no git directory, because it
// keeps no record.
func endedRecord(slot backgroundstate.Slot) (string, backgroundstate.Record, error) {
	gitdir, err := benchgit.AdminDir(slot.Worktree)
	if err != nil {
		return "", backgroundstate.Record{}, nil
	}
	record, err := backgroundstate.ReadRecord(gitdir)
	return gitdir, record, err
}

// startOldest starts the oldest queued grade while the slot of queue is free. A grade that does
// not start leaves the queue, and the next one tries.
func startOldest(home string, queue *backgroundstate.Queue) error {
	var errs []error
	for queue.Running == nil && len(queue.Queued) > 0 {
		next := queue.Queued[0]
		queue.Queued = queue.Queued[1:]
		pid, err := startGrade(home, next.Worktree, next.Tip)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", next.Worktree, err))
			continue
		}
		slot, err := backgroundstate.SlotFor(next, pid)
		if err != nil {
			return errors.Join(append(errs, err)...)
		}
		queue.Running = &slot
	}
	return errors.Join(errs...)
}

// startGrade starts the background grade of tip in root through the lifetime owner, as a
// detached session, and answers its pid. The grade learns the queue under home, so it starts the
// next grade when it ends. The owned group is only the detach helper, which exits once the grade
// starts, so the teardown window of a process bounds it, and a cancellation kills it at once.
func startGrade(home, root, tip string) (int, error) {
	ctx, cancel := bounds.Context(context.Background(), bounds.FixedWindow(bounds.ProcessFinalWindow))
	defer cancel()
	cmd := exec.Command(backgroundWrapper(root), backgroundVerb, root, tip)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), backgroundQueueEnv+"="+home)
	detached, err := subprocess.Detach(cmd)
	if err != nil {
		return 0, err
	}
	return detached.PID(subprocess.RunInherited(ctx, detached.Cmd, subprocess.Policy{}, detached.Streams, subprocess.Interrupt(syscall.SIGKILL)))
}

// backgroundWrapper answers the wrapper of the worktree root.
func backgroundWrapper(root string) string {
	return filepath.Join(root, "bin", "bench.sh")
}
