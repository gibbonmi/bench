package otelrecord

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/poolkey"
)

// recordFile names the record inside the repository's record directory.
const recordFile = "traces.jsonl"

// Dir returns the directory that holds one repository's seam record. The directory
// is a sibling of the census directory, so both records key by the same repository
// identity.
func Dir(home, root string) string {
	return filepath.Join(home, "otel", poolkey.Key(root))
}

// Path returns the record file for one repository below an explicitly resolved home.
func Path(home, root string) string {
	return filepath.Join(Dir(home, root), recordFile)
}

// Writer appends encoded spans to one repository's record file. The caller resolves
// the Bench home and passes it in, the census form, so the package reads no
// environment variable itself.
type Writer struct {
	home string
	dir  string
}

// NewWriter returns the writer for root's record below an explicitly resolved home.
// The writer opens the file on each append, so two writers share no state. The home
// is kept because the record path is graded against it on each append.
func NewWriter(home, root string) *Writer {
	return &Writer{home: home, dir: Dir(home, root)}
}

// gradeRecordPath refuses a record path that the appender must not follow or open.
// Two failures live here. A symlink at any level below the home redirects the record
// outside the home, because os.MkdirAll follows a link at a parent level as readily as
// at the leaf. A non-regular file at the record path — a FIFO or a device — blocks the
// open, so every recorded verb would hang on its first span.
func (w *Writer) gradeRecordPath(file string) error {
	for _, level := range levelsBelow(w.home, file) {
		info, err := os.Lstat(level)
		if err != nil {
			continue
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("record path %s is a symlink", level)
		}
	}
	info, err := os.Lstat(file)
	if err != nil || info.Mode().IsRegular() {
		return nil
	}
	return fmt.Errorf("record path %s is not a regular file", file)
}

// levelsBelow returns each path level from home down to path, outermost first. The
// home itself is the operator's own directory and is not graded: the writer owns what
// it creates below the home, not the home. A path outside the home grades whole.
func levelsBelow(home, path string) []string {
	var levels []string
	for current := path; current != home; current = filepath.Dir(current) {
		levels = append(levels, current)
		if parent := filepath.Dir(current); parent == current {
			break
		}
	}
	for left, right := 0, len(levels)-1; left < right; left, right = left+1, right-1 {
		levels[left], levels[right] = levels[right], levels[left]
	}
	return levels
}

// Append writes one encoded record line. Encode returns a line with no terminator,
// so the writer owns the newline. Each append is one synchronous O_APPEND write with
// no buffer and no background worker, which keeps concurrent writers' lines intact.
// The writer returns every failure and swallows none; the caller decides whether a
// failed record changes its outcome.
func (w *Writer) Append(line []byte) error {
	record := filepath.Join(w.dir, recordFile)
	if err := w.gradeRecordPath(record); err != nil {
		return err
	}
	if err := os.MkdirAll(w.dir, 0o700); err != nil {
		return fmt.Errorf("create record directory: %w", err)
	}
	file, err := os.OpenFile(record, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("open seam record: %w", err)
	}
	defer file.Close()
	if _, err := file.Write(append(append([]byte{}, line...), '\n')); err != nil {
		return fmt.Errorf("append seam record: %w", err)
	}
	return nil
}
