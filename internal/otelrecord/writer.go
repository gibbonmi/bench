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
	dir string
}

// NewWriter returns the writer for root's record below an explicitly resolved home.
// The writer opens the file on each append, so two writers share no state.
func NewWriter(home, root string) *Writer {
	return &Writer{dir: Dir(home, root)}
}

// Append writes one encoded record line. Encode returns a line with no terminator,
// so the writer owns the newline. Each append is one synchronous O_APPEND write with
// no buffer and no background worker, which keeps concurrent writers' lines intact.
// The writer returns every failure and swallows none; the caller decides whether a
// failed record changes its outcome.
func (w *Writer) Append(line []byte) error {
	if info, err := os.Lstat(w.dir); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("record directory %s is a symlink", w.dir)
	}
	if err := os.MkdirAll(w.dir, 0o700); err != nil {
		return fmt.Errorf("create record directory: %w", err)
	}
	file, err := os.OpenFile(filepath.Join(w.dir, recordFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("open seam record: %w", err)
	}
	defer file.Close()
	if _, err := file.Write(append(append([]byte{}, line...), '\n')); err != nil {
		return fmt.Errorf("append seam record: %w", err)
	}
	return nil
}
