package backgroundstate

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/durablefile"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/subprocess"
)

// The windows of the queue steps. A queue holder holds the lock at most for the start of one
// grade, which the teardown window of a process bounds. A stopped grade cancels its oracle
// group and then tears down its own resources.
const (
	LockWindow = 2 * bounds.ProcessFinalWindow
	StopWindow = bounds.ProcessGroupCancelGrace + bounds.ProcessFinalWindow
)

// stopFile is the stop marker beside the background record of a worktree. It holds the claim of
// the grade that a stopper asked to end.
const stopFile = "bench-background-grade.stop"

// Stop is the one stop step of a background grade. It asks the owner of claim to end, and it
// returns only after that owner has exited, or when ctx ends. Before the signal, it writes the
// stop marker of claim beside the background record of its worktree, so WriteEnd names the end of
// that owner stopped, whenever that end comes. The signal goes through the lifetime owner in
// internal/subprocess. A claim that is not live has no owner to stop, so Stop sends nothing: a
// pid of an earlier boot can name another process. A worktree whose git directory does not
// resolve keeps no record, so it takes no marker. A marker that fails does not keep the grade.
func Stop(ctx context.Context, replacer durablefile.Replacer, claim Claim) error {
	live, err := claim.Live()
	if err != nil || !live {
		return err
	}
	var marked error
	if gitdir, err := benchgit.AdminDir(claim.Worktree); err == nil {
		marked = writeRecord(replacer, stopPath(gitdir), claim)
	}
	return errors.Join(marked, subprocess.Stop(ctx, claim.Owner, bounds.ProcessPoll))
}

// StopSlot stops the running grade of slot through Stop within the stop window.
func StopSlot(replacer durablefile.Replacer, slot Slot) error {
	ctx, cancel := bounds.Context(context.Background(), bounds.FixedWindow(StopWindow))
	defer cancel()
	return Stop(ctx, replacer, slot.Claim())
}

// StopAll drops every queued grade of the machine queue under the Bench home home through
// replacer, and then stops its running grade through StopSlot. The lock wait ends with the lock
// window. A queue that holds no grade needs no lock. The stopped grade frees its own slot when it
// ends, and the next queue step frees the slot of a grade whose owner died.
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
	return StopSlot(replacer, *running)
}

// WriteEnd installs record, the end that the grade of owner reached, as the background record of
// the worktree git directory gitdir, and it answers the record that it installed. A grade writes
// its own end, and the next queue step writes the end of a grade whose owner died. An aborted end
// names stopped when the stop marker in gitdir names owner, because a stop came first. A record
// with no key takes the key of the marker, and with no key it stays aborted. A green or a red
// keeps its outcome. The log of the running grade replaces the log together with the record, so
// the log that a record names holds the output of its own grade. The marker of owner is removed
// after the record.
func WriteEnd(replacer durablefile.Replacer, gitdir string, owner Claim, record Record) (Record, error) {
	marker, marked := stopOf(gitdir, owner)
	if marked && record.Outcome == Aborted {
		if record.Key == "" {
			record.Key = marker.Key
		}
		if record.Key != "" {
			record.Outcome = Stopped
		}
	}
	if err := record.valid(); err != nil {
		return record, err
	}
	if err := os.Rename(RunningLogPath(gitdir), LogPath(gitdir)); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return record, err
	}
	if err := WriteRecord(replacer, gitdir, record); err != nil || !marked {
		return record, err
	}
	if err := os.Remove(stopPath(gitdir)); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return record, err
	}
	return record, nil
}

// stopOf answers the stop marker in the worktree git directory gitdir, and whether it names the
// process of owner: the same worktree, pid, and boot. The key is not compared, because a grade
// that ends before it knows its key holds none. An absent marker and a marker that does not read
// name no stop, so the end keeps the outcome that the grade reached.
func stopOf(gitdir string, owner Claim) (Claim, bool) {
	var marker Claim
	err := readRecord(stopPath(gitdir), &marker, func() error {
		if err := validKey(marker.Key); err != nil {
			return err
		}
		return validOwner(marker.Owner)
	})
	if err != nil || marker == (Claim{}) || marker.Owner != owner.Owner || marker.Boot != owner.Boot || marker.Worktree != owner.Worktree {
		return Claim{}, false
	}
	return marker, true
}

// stopPath answers the path of the stop marker in the worktree git directory gitdir.
func stopPath(gitdir string) string { return filepath.Join(gitdir, stopFile) }
