package otelrecord

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/poolkey"
)

// The record's files inside the repository's record directory. The live segment takes
// every append. A rotation renames it to a sealed segment named by a sequence, and only
// the rotation takes the lock file.
const (
	recordFile       = "traces.jsonl"
	rotationLockFile = "traces.lock"

	sealedPrefix = "traces-"
	sealedSuffix = ".jsonl"
	// sealedDigits pads the sequence, so name order is sequence order.
	sealedDigits = 20
)

// sealedName formats the sealed segment name of one sequence. No clock enters the name,
// so two rotations in one instant never share a name and a clock step never reorders them.
func sealedName(sequence uint64) string {
	return fmt.Sprintf("%s%0*d%s", sealedPrefix, sealedDigits, sequence, sealedSuffix)
}

// parseSealedName returns the sequence a sealed segment name carries, and reports whether
// the name is a sealed segment name at all.
func parseSealedName(name string) (uint64, bool) {
	digits, ok := strings.CutPrefix(name, sealedPrefix)
	if !ok {
		return 0, false
	}
	if digits, ok = strings.CutSuffix(digits, sealedSuffix); !ok || len(digits) != sealedDigits {
		return 0, false
	}
	sequence, err := strconv.ParseUint(digits, 10, 64)
	return sequence, err == nil
}

// sealedSequences returns the sequences of the sealed segments in dir, lowest first.
func sealedSequences(dir string) ([]uint64, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var sequences []uint64
	// ReadDir sorts by name, and name order is sequence order.
	for _, entry := range entries {
		if sequence, ok := parseSealedName(entry.Name()); ok {
			sequences = append(sequences, sequence)
		}
	}
	return sequences, nil
}

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

// memoryDir holds the retained shift notes below the record directory. The notes are
// model prose, so they live beside the record and never inside a record line.
const memoryDir = "memory"

// MemoryDir returns the directory that holds root's retained shift notes.
func MemoryDir(home, root string) string {
	return filepath.Join(Dir(home, root), memoryDir)
}

// RetainMemory keeps body as one memory file of root's record, below home or the resolved
// Bench home when home is empty, and returns the file's SHA-256 digest. The store keeps the
// newest RecordMemoryRetained files.
func RetainMemory(home, root, traceID string, body []byte) (string, error) {
	home, ok := recordHome(home)
	if !ok {
		return "", errors.New("a test binary may not write below the fallback Bench home")
	}
	return retainMemory(home, root, traceID, body, time.Now(), bounds.RecordMemoryRetained)
}

// retainMemory is RetainMemory with the instant and the retained count given. The file name
// is the UTC instant and the trace id, so name order is time order. The path is graded as
// the record path is, and the file is created 0600 and never opened through a link.
func retainMemory(home, root, traceID string, body []byte, now time.Time, retained int) (string, error) {
	dir := MemoryDir(home, root)
	file := filepath.Join(dir, now.UTC().Format("20060102T150405.000000000Z")+"-"+traceID+".md")
	if err := (&Writer{home: home}).gradeRecordPath(file); err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", fmt.Errorf("create memory directory: %w", err)
	}
	out, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_EXCL|syscall.O_NOFOLLOW, 0o600)
	if err != nil {
		return "", fmt.Errorf("create memory file: %w", err)
	}
	_, err = out.Write(body)
	if err = errors.Join(err, out.Close()); err != nil {
		return "", fmt.Errorf("write memory file: %w", err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("list memory directory: %w", err)
	}
	// ReadDir sorts by name, and name order is time order.
	for index := 0; index < len(entries)-retained; index++ {
		if err := os.Remove(filepath.Join(dir, entries[index].Name())); err != nil {
			return "", fmt.Errorf("prune memory file: %w", err)
		}
	}
	sum := sha256.Sum256(body)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

// Writer appends encoded spans to one repository's record file. The caller resolves
// the Bench home and passes it in, the census form, so the package reads no
// environment variable itself.
type Writer struct {
	home     string
	dir      string
	limit    int64
	retained int
}

// NewWriter returns the writer for root's record below an explicitly resolved home.
// The writer opens the file on each append, so two writers share no state. The home
// is kept because the record path is graded against it on each append.
func NewWriter(home, root string) *Writer {
	return newWriter(home, root, bounds.RecordSegmentLimit, bounds.RecordSegmentsRetained)
}

// newWriter is NewWriter with the segment limit and the retained count given.
func newWriter(home, root string, limit int64, retained int) *Writer {
	return &Writer{home: home, dir: Dir(home, root), limit: limit, retained: retained}
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
// An append that would take the live segment past the limit first seals it. A failed
// rotation still appends the line. The writer returns every failure and swallows none;
// the caller decides whether a failed record changes its outcome.
func (w *Writer) Append(line []byte) error {
	record := filepath.Join(w.dir, recordFile)
	if err := w.gradeRecordPath(record); err != nil {
		return err
	}
	if err := os.MkdirAll(w.dir, 0o700); err != nil {
		return fmt.Errorf("create record directory: %w", err)
	}
	data := append(append([]byte{}, line...), '\n')
	rotated := w.rotate(record, int64(len(data)))
	file, err := os.OpenFile(record, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return errors.Join(rotated, fmt.Errorf("open seam record: %w", err))
	}
	defer file.Close()
	if _, err := file.Write(data); err != nil {
		return errors.Join(rotated, fmt.Errorf("append seam record: %w", err))
	}
	return rotated
}

// full reports whether an append of incoming bytes takes the live segment past the
// limit. An empty or absent segment is never full, so one oversized line still lands.
func (w *Writer) full(record string, incoming int64) (bool, error) {
	info, err := os.Lstat(record)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("stat seam record: %w", err)
	}
	return info.Size() > 0 && info.Size()+incoming > w.limit, nil
}

// rotate seals a full live segment under the next sequence and prunes the lowest
// sequences down to the retained count. Only the rotation takes the lock, and it never
// waits for it: a writer that finds the lock held appends to the live segment, and a
// later append rotates.
func (w *Writer) rotate(record string, incoming int64) error {
	if full, err := w.full(record, incoming); !full || err != nil {
		return err
	}
	lockPath := filepath.Join(w.dir, rotationLockFile)
	if err := w.gradeRecordPath(lockPath); err != nil {
		return err
	}
	lock, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR|syscall.O_NOFOLLOW, 0o600)
	if err != nil {
		return fmt.Errorf("open rotation lock: %w", err)
	}
	defer lock.Close()
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		if errors.Is(err, syscall.EWOULDBLOCK) {
			return nil
		}
		return fmt.Errorf("take rotation lock: %w", err)
	}
	defer syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)

	// Another writer can have rotated between the first check and the lock.
	if full, err := w.full(record, incoming); !full || err != nil {
		return err
	}
	sequences, err := sealedSequences(w.dir)
	if err != nil {
		return fmt.Errorf("list sealed segments: %w", err)
	}
	// The listing holds every sealed name, planted or not, and only a lock holder seals, so
	// one more than the highest sequence is a free name. A highest sequence at the maximum
	// has no successor, and a wrapped sequence would rename over a sealed segment.
	next := uint64(1)
	if len(sequences) > 0 {
		highest := sequences[len(sequences)-1]
		if highest == math.MaxUint64 {
			return fmt.Errorf("seal seam record: %s holds the last sequence", sealedName(highest))
		}
		next = highest + 1
	}
	if err := os.Rename(record, filepath.Join(w.dir, sealedName(next))); err != nil {
		return fmt.Errorf("seal seam record: %w", err)
	}
	return w.prune(append(sequences, next))
}

// prune removes the lowest sequences until the retained count of sealed segments
// remains. It orders by sequence alone, so a segment's modification time never decides.
func (w *Writer) prune(sequences []uint64) error {
	var failures []error
	for len(sequences) > w.retained {
		if err := os.Remove(filepath.Join(w.dir, sealedName(sequences[0]))); err != nil && !errors.Is(err, os.ErrNotExist) {
			failures = append(failures, fmt.Errorf("prune sealed segment: %w", err))
		}
		sequences = sequences[1:]
	}
	return errors.Join(failures...)
}

// segments returns root's record segments in read order: the sealed segments by
// sequence, then the live segment. Each path passes the writer's grade before a reader
// opens it, so a planted link or special file refuses the read rather than misdirecting
// or blocking it.
func segments(home, root string) ([]string, error) {
	w := NewWriter(home, root)
	live := filepath.Join(w.dir, recordFile)
	if err := w.gradeRecordPath(live); err != nil {
		return nil, err
	}
	sequences, err := sealedSequences(w.dir)
	if err != nil {
		return nil, fmt.Errorf("open seam record: %w", err)
	}
	paths := make([]string, 0, len(sequences)+1)
	for _, sequence := range sequences {
		paths = append(paths, filepath.Join(w.dir, sealedName(sequence)))
	}
	paths = append(paths, live)
	for _, path := range paths {
		if err := w.gradeRecordPath(path); err != nil {
			return nil, err
		}
	}
	return paths, nil
}
