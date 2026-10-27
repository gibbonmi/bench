package backgroundstate

import (
	"path/filepath"

	"github.com/gibbonmi/bench/internal/durablefile"
)

// queueDir is the directory of the machine queue under the Bench home, so one queue serves
// every repository of the machine.
const (
	queueDir  = "background-grade"
	queueFile = "queue"
)

// Entry is one grade of the queue: its grade key, its worktree, and the tip that it grades.
type Entry struct {
	Key      string `json:"key"`
	Worktree string `json:"worktree"`
	Tip      string `json:"tip"`
}

// Queue is the machine queue. Running is the slot, which holds the one running grade or
// nothing, and Queued holds the waiting grades in their order.
type Queue struct {
	Running *Entry  `json:"running"`
	Queued  []Entry `json:"queued"`
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
	return nil
}

func (e Entry) valid() error {
	if err := pathName("key", e.Key); err != nil {
		return err
	}
	return required("worktree", e.Worktree, "tip", e.Tip)
}

// QueuePath answers the path of the machine queue under the Bench home home.
func QueuePath(home string) string {
	return filepath.Join(home, queueDir, queueFile)
}
