package gate

// The launch of a background grade: after a lane pass publishes a tip in a kit worktree, the
// launch detaches the lane pass of that tip, which passes its grade to the machine queue. The
// queue runs one grade at a time through the wrapper of its worktree as a detached session, and
// no launch waits for a pass or a grade. When a grade ends, its own process starts the oldest
// queued grade and then exits.

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
	"github.com/gibbonmi/bench/internal/landing/published"
	"github.com/gibbonmi/bench/internal/subprocess"
)

// StartBackground starts the lane pass of tip, the commit that a lane pass just published in the
// worktree root, for the machine queue under the Bench home home. kit is the kit that the verb
// read. A pass starts only in a background worktree. The verb only detaches the pass, so it
// derives no grade key, takes no queue lock, and stops no grade. The pass runs as the executable
// that the wrapper resolves, so its key is the key of the grade. The publication stands whatever
// the launch does, so a launch that fails prints one line to stderr and the verb keeps its exit.
func StartBackground(home, kit, root, tip string, stderr io.Writer) {
	if !backgroundWorktree(kit, root) {
		return
	}
	if _, err := startGrade(root, tip, backgroundPassEnv+"="+home); err != nil {
		fmt.Fprintf(stderr, "background grade: not queued: %v\n", err)
	}
}

// backgroundWorktree reports whether root, under the kit that the verb read, is a worktree that
// background grades run in: the worktree of an active assignment for which KitForRoot answers root
// itself, which holds its executable wrapper, because a grade runs as that wrapper. The checks
// read no tree, so a worktree of another repository pays only for them.
func backgroundWorktree(kit, root string) bool {
	if KitForRoot(kit, root) != root || !RealFS().Executable(backgroundWrapper(root)) {
		return false
	}
	_, assigned := intent.AssignmentForWorktree(root)
	return assigned
}

// passLane applies the lane pass of tip in root to the queue under home. The pass queues nothing
// for a key that a fresh green or a live claim already holds. A running grade that the pass
// supersedes stops before the queue moves on. A claim that does not read retains its grade, so
// the pass queues nothing. A pass runs after the verb that published its tip, so a later pass of
// root can run first: a tip that is not the head of root drops out, and it replaces no grade of a
// newer tip. A dropped pass still moves the queue on. A landing subject whose landing plan
// selects the reader set, and a landing that its composition refuses, queue nothing, because
// the landing runs no whole gate that a grade could answer.
func passLane(home, root, tip string) error {
	key, settled := "", false
	switch subject, derived := backgroundSubjectOf(root, tip); {
	case errors.As(derived, new(published.ConflictError)):
		return advanceQueue(home, 0)
	case derived != nil:
	case subject.landing != nil:
		evidence, reader, inspected := subject.landing.inspect(subject.tree)
		if inspected == nil && reader {
			return advanceQueue(home, 0)
		}
		key, settled = evidence.Grade, evidence.ReusableGreen
	default:
		evidence := InspectTree(root, subject.tree)
		key, settled = evidence.Grade, evidence.ReusableGreen
	}
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
		commit, err := benchgit.ResolveCommit(root, tip)
		if head, headErr := benchgit.ResolveCommit(root, "HEAD"); err == nil && headErr == nil && commit == head {
			superseded = queue.Pass(backgroundstate.Entry{Key: key, Worktree: root, Tip: commit}, settled)
		}
		return errors.Join(reaped, startOldest(home, queue))
	})
	if err != nil || superseded == nil {
		return err
	}
	return errors.Join(stopGrade(*superseded), advanceQueue(home, 0))
}

// StopWorktreeGrades takes each worktree at dir or under it out of the machine queue under home,
// for a verb that removes that worktree or lands it as a stack member, before the verb acts. It
// drops each queued grade of such a worktree, and it stops the running grade of such a worktree
// and waits for its exit. The queue then starts its oldest queued grade. A queue that holds no
// grade of dir needs no lock.
func StopWorktreeGrades(home, dir string) error {
	// The withdrawal first runs on the queue as read, which it does not write.
	queue, err := backgroundstate.ReadQueue(home)
	if queued := len(queue.Queued); err != nil || (queue.Withdraw(dir) == nil && len(queue.Queued) == queued) {
		return err
	}
	var stopped *backgroundstate.Slot
	err = updateQueue(home, func(queue *backgroundstate.Queue) error {
		reaped := reapRunning(queue)
		stopped = queue.Withdraw(dir)
		return reaped
	})
	if err != nil {
		return err
	}
	if stopped != nil {
		err = stopGrade(*stopped)
	}
	return errors.Join(err, advanceQueue(home, 0))
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
	ctx, cancel := bounds.Context(context.Background(), bounds.FixedWindow(backgroundstate.LockWindow))
	defer cancel()
	return backgroundstate.UpdateQueue(ctx, durablefile.Native(), home, update)
}

// reapRunning frees the slot of a running grade whose owner is not live. Unless the grade recorded
// its own end, its worktree records the end of the grade through WriteEnd: stopped when a stop
// marked that grade, and aborted otherwise. A slot whose liveness does not read retains its grade.
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
	_, err = backgroundstate.WriteEnd(durablefile.Native(), gitdir, dead.Claim(), backgroundstate.Record{Outcome: backgroundstate.Aborted, Key: dead.Key, Tip: dead.Tip, Log: backgroundstate.LogPath(gitdir)})
	return err
}

// stopGrade stops the running grade of slot through the one stop step of backgroundstate, which
// marks the grade before its signal, so the grade names its own end stopped.
func stopGrade(slot backgroundstate.Slot) error {
	return backgroundstate.StopSlot(durablefile.Native(), slot)
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
		pid, err := startGrade(next.Worktree, next.Tip, backgroundQueueEnv+"="+home)
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

// startGrade starts the background entry of tip in root through the lifetime owner, as a detached
// session, and answers its pid. marker is the one environment entry that names the role of the
// entry: the queue of a grade, which the grade advances when it ends, or the queue of a lane pass.
// The owned group is only the detach helper, which exits once the entry starts, so the teardown
// window of a process bounds it, and a cancellation kills it at once.
func startGrade(root, tip, marker string) (int, error) {
	ctx, cancel := bounds.Context(context.Background(), bounds.FixedWindow(bounds.ProcessFinalWindow))
	defer cancel()
	cmd := exec.Command(backgroundWrapper(root), backgroundVerb, root, tip)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), marker)
	detached, err := subprocess.Detach(cmd)
	if err != nil {
		return 0, err
	}
	return detached.PID(subprocess.RunInherited(ctx, detached.Cmd, subprocess.Policy{}, detached.Streams, subprocess.Interrupt(syscall.SIGKILL), subprocess.Registration{}))
}

// backgroundWrapper answers the wrapper of the worktree root.
func backgroundWrapper(root string) string {
	return filepath.Join(root, "bin", "bench.sh")
}
