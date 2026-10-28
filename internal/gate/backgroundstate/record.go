package backgroundstate

import (
	"fmt"
	"path/filepath"
	"slices"

	"github.com/gibbonmi/bench/internal/durablefile"
)

// recordFile is the background record in the git directory of one worktree. It sits beside
// the gate cache of that worktree, so it names the latest background grade of that checkout.
const recordFile = "bench-background-grade"

// logFile is the background log beside the background record. It holds the output of the
// grade that the background record names. runningLogFile beside it holds the output of the
// grade that runs, until the end of that grade replaces the log with it.
const (
	logFile        = "bench-background-grade.log"
	runningLogFile = "bench-background-grade.log.running"
)

// LogPath answers the path of the background log of the worktree git directory gitdir.
func LogPath(gitdir string) string { return filepath.Join(gitdir, logFile) }

// RunningLogPath answers the path of the log of the running grade of the worktree git
// directory gitdir.
func RunningLogPath(gitdir string) string { return filepath.Join(gitdir, runningLogFile) }

// Outcome is the end of one background grade.
type Outcome string

// The outcomes of a background grade.
const (
	Green   Outcome = "green"
	Red     Outcome = "red"
	Aborted Outcome = "aborted"
	Stopped Outcome = "stopped"
)

// outcomes is the closed set of outcomes that a record can name.
var outcomes = []Outcome{Green, Red, Aborted, Stopped}

// Record is the latest background grade of a worktree: its outcome, its grade key, the tip
// that it graded, and the path of its log. An aborted record can name no key.
type Record struct {
	Outcome Outcome `json:"outcome"`
	Key     string  `json:"key"`
	Tip     string  `json:"tip"`
	Log     string  `json:"log"`
}

// WriteRecord installs record as the background record of the worktree git directory gitdir.
func WriteRecord(replacer durablefile.Replacer, gitdir string, record Record) error {
	if err := record.valid(); err != nil {
		return err
	}
	return writeRecord(replacer, RecordPath(gitdir), record)
}

// RecordPath answers the path of the background record of the worktree git directory gitdir.
func RecordPath(gitdir string) string {
	return filepath.Join(gitdir, recordFile)
}

// ReadRecord reads the background record of the worktree git directory gitdir. An absent
// record is the empty record.
func ReadRecord(gitdir string) (Record, error) {
	var record Record
	if err := readRecord(RecordPath(gitdir), &record, func() error { return record.valid() }); err != nil {
		return Record{}, err
	}
	return record, nil
}

func (r Record) valid() error {
	if !slices.Contains(outcomes, r.Outcome) {
		return fmt.Errorf("outcome %q is not one of %v", r.Outcome, outcomes)
	}
	if r.Outcome == Aborted && r.Key == "" {
		return required("tip", r.Tip, "log", r.Log)
	}
	return required("key", r.Key, "tip", r.Tip, "log", r.Log)
}
