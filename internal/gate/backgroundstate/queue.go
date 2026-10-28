package backgroundstate

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"syscall"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/durablefile"
)

// queueDir is the directory of the machine queue under the Bench home, so one queue serves
// every repository of the machine. The queue lock sits beside the queue.
const (
	queueDir  = "background-grade"
	queueFile = "queue"
	queueLock = "queue.lock"
)

// Entry is one grade of the queue: its grade key, its worktree, and the tip that it grades. A
// tip whose subject gives no grade key has the empty key, and its grade records aborted.
type Entry struct {
	Key      string `json:"key"`
	Worktree string `json:"worktree"`
	Tip      string `json:"tip"`
}

// Queue is the machine queue. Running is the slot, which holds the one running grade or
// nothing, and Queued holds the waiting grades in the order of their lane passes. Latest maps
// each worktree to the grade key of its latest lane-passed tip, while the queue holds that key.
type Queue struct {
	Running *Slot             `json:"running"`
	Queued  []Entry           `json:"queued"`
	Latest  map[string]string `json:"latest,omitempty"`
}

// Pass applies one lane pass of entry. settled reports that a fresh green or a live claim
// already holds the key of entry. The pass drops each other queued grade of the worktree of
// entry, and it queues entry unless its key is settled, running, or queued. It answers the
// running grade that the pass supersedes, which the caller stops, or nil. A grade whose key is
// the latest key of another worktree is never dropped and never superseded. A tip with no key
// supersedes no running grade, and a running grade with no key ends with no oracle.
func (q *Queue) Pass(entry Entry, settled bool) *Slot {
	if q.Latest == nil {
		q.Latest = map[string]string{}
	}
	q.Latest[entry.Worktree] = entry.Key
	q.Queued = slices.DeleteFunc(q.Queued, func(queued Entry) bool {
		if queued.Worktree != entry.Worktree || (entry.Key != "" && queued.Key == entry.Key) {
			return false
		}
		return !q.wanted(queued.Key, entry.Worktree)
	})
	if !settled && !q.holds(entry.Key) {
		q.Queued = append(q.Queued, entry)
	}
	if r := q.Running; r != nil && r.Worktree == entry.Worktree && r.Key != "" && entry.Key != "" && r.Key != entry.Key && !q.wanted(r.Key, entry.Worktree) {
		return r
	}
	return nil
}

// wanted reports whether a worktree other than except has key as the key of its latest tip.
// The empty key is no grade, so no worktree wants it.
func (q *Queue) wanted(key, except string) bool {
	for worktree, latest := range q.Latest {
		if worktree != except && latest == key && key != "" {
			return true
		}
	}
	return false
}

// holds reports whether the running grade or a queued grade has key. The empty key is no
// grade, so the queue holds none of it.
func (q *Queue) holds(key string) bool {
	if key == "" {
		return false
	}
	if q.Running != nil && q.Running.Key == key {
		return true
	}
	return slices.ContainsFunc(q.Queued, func(queued Entry) bool { return queued.Key == key })
}

// prune drops each latest key that the queue no longer holds, so Latest names only the work
// that a worktree still waits for.
func (q *Queue) prune() {
	for worktree, key := range q.Latest {
		if !q.holds(key) {
			delete(q.Latest, worktree)
		}
	}
	if len(q.Latest) == 0 {
		q.Latest = nil
	}
}

// UpdateQueue applies update to the machine queue under the Bench home home, under the queue
// lock, and installs the result through replacer. The lock wait ends with ctx. A queue that
// does not read retains its grades, so update does not run. The queue that update leaves is
// installed even when update answers an error, because update can start a grade before a later
// step fails.
func UpdateQueue(ctx context.Context, replacer durablefile.Replacer, home string, update func(*Queue) error) (err error) {
	release, err := lockQueue(ctx, home)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, release()) }()
	queue, err := ReadQueue(home)
	if err != nil {
		return err
	}
	updated := update(&queue)
	queue.prune()
	return errors.Join(updated, WriteQueue(replacer, home, queue))
}

// lockQueue takes the lock of the machine queue under the Bench home home and answers its
// release. While another process holds the lock, each expiry of the process poll tries again,
// until ctx ends.
func lockQueue(ctx context.Context, home string) (func() error, error) {
	if err := ensureDir(filepath.Join(home, queueDir)); err != nil {
		return nil, err
	}
	lock, err := os.OpenFile(filepath.Join(home, queueDir, queueLock), os.O_CREATE|os.O_RDWR, recordMode)
	if err != nil {
		return nil, err
	}
	for {
		err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err == nil {
			return lock.Close, nil
		}
		if !errors.Is(err, syscall.EWOULDBLOCK) {
			return nil, errors.Join(err, lock.Close())
		}
		expired, stop := bounds.ProcessPoll()
		select {
		case <-ctx.Done():
			stop()
			return nil, errors.Join(ctx.Err(), lock.Close())
		case <-expired:
		}
	}
}

// WriteQueue installs queue as the machine queue under the Bench home home.
func WriteQueue(replacer durablefile.Replacer, home string, queue Queue) error {
	if err := queue.valid(); err != nil {
		return err
	}
	if err := ensureDir(filepath.Join(home, queueDir)); err != nil {
		return err
	}
	return writeRecord(replacer, QueuePath(home), queue)
}

// ReadQueue reads the machine queue under the Bench home home. An absent queue is the empty
// queue.
func ReadQueue(home string) (Queue, error) {
	var queue Queue
	if err := readRecord(QueuePath(home), &queue, func() error { return queue.valid() }); err != nil {
		return Queue{}, err
	}
	return queue, nil
}

func (q Queue) valid() error {
	if q.Running != nil {
		if err := q.Running.valid(); err != nil {
			return err
		}
	}
	for _, entry := range q.Queued {
		if err := entry.valid(); err != nil {
			return err
		}
	}
	for worktree, key := range q.Latest {
		if err := validKey(key); err != nil {
			return err
		}
		if err := required("worktree", worktree); err != nil {
			return err
		}
	}
	return nil
}

// valid refuses an entry with no worktree or no tip, and a key that validKey refuses.
func (e Entry) valid() error {
	if err := validKey(e.Key); err != nil {
		return err
	}
	return required("worktree", e.Worktree, "tip", e.Tip)
}

// validKey refuses a key that is not one path name. The empty key is the key of a tip whose
// subject gives none.
func validKey(key string) error {
	if key == "" {
		return nil
	}
	return pathName("key", key)
}

// QueuePath answers the path of the machine queue under the Bench home home.
func QueuePath(home string) string {
	return filepath.Join(home, queueDir, queueFile)
}
