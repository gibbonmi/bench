package responsebound

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
)

// Source is one complete file that an exact read pages. Its digest covers every byte of the
// file. A bounded source holds its complete bytes; a stream source holds the open file and
// reads each page window from it, so an artifact larger than any control record pages with
// bounded memory. Close releases the file.
type Source struct {
	file   *os.File
	info   fs.FileInfo
	path   string
	data   []byte
	stream bool
	// Digest is the SHA-256 of the complete source bytes, and Size is their count.
	Digest [sha256.Size]byte
	Size   int64
}

// SourceError is the file state that refuses an exact read before any content.
type SourceError struct {
	State  bounds.FileState
	Reason string
}

func (e *SourceError) Error() string { return string(e.State) + ": " + e.Reason }

// changedReason is the refusal of a file whose identity or size moved during the read.
const changedReason = "the file changed during the read"

// OpenSource opens the file at path for one exact read through the no-follow opener. A
// bounded source reads its complete bytes within ControlRecordLimit. A stream source hashes
// the complete file and keeps it open for the page window. A missing file is StateAbsent, an
// operational refusal distinct from a present empty file.
func OpenSource(path string, stream bool) (*Source, error) {
	file, classified := bounds.OpenNoFollow(path)
	if file == nil {
		reason := classified.Reason
		if classified.State == bounds.StateAbsent {
			reason = "no file at the path"
		}
		return nil, &SourceError{State: classified.State, Reason: reason}
	}
	info, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return nil, &SourceError{State: bounds.StateUnreadable, Reason: err.Error()}
	}
	source := &Source{file: file, info: info, path: path, stream: stream}
	hash := sha256.New()
	if stream {
		source.Size, err = io.Copy(hash, file)
	} else {
		read := bounds.Read(file, bounds.ControlRecordLimit)
		err = read.Err
		source.data = read.Data
		source.Size = int64(len(read.Data))
		hash.Write(read.Data)
	}
	if err != nil {
		_ = file.Close()
		return nil, &SourceError{State: bounds.StateUnreadable, Reason: err.Error()}
	}
	copy(source.Digest[:], hash.Sum(nil))
	if err := source.Unchanged(); err != nil {
		_ = file.Close()
		return nil, err
	}
	return source, nil
}

// ReadAt reads the page window at off from the complete source.
func (s *Source) ReadAt(p []byte, off int64) (int, error) {
	if s.stream {
		return s.file.ReadAt(p, off)
	}
	return bytes.NewReader(s.data).ReadAt(p, off)
}

// Bytes answers the complete bytes of a bounded source, and nil for a stream source.
func (s *Source) Bytes() []byte { return s.data }

// Unchanged refuses a source whose path now names another file, or whose open file moved in
// size or modification time since it opened. A reader calls it after its page is assembled,
// so a replaced or modified file never returns a partial page.
func (s *Source) Unchanged() error {
	now, err := s.file.Stat()
	leaf, leafErr := os.Lstat(s.path)
	if err != nil || leafErr != nil || !os.SameFile(s.info, leaf) || now.Size() != s.Size || !now.ModTime().Equal(s.info.ModTime()) {
		return &SourceError{State: bounds.StateUnreadable, Reason: changedReason}
	}
	return nil
}

// Close releases the open file.
func (s *Source) Close() error { return s.file.Close() }

// ParseLineRange reads a `<first>:<last>` physical-line range: decimal one-based bounds with
// first no greater than last.
func ParseLineRange(value string) (int, int, error) {
	firstText, lastText, found := strings.Cut(value, ":")
	first, firstErr := decimal(firstText)
	last, lastErr := decimal(lastText)
	if !found || firstErr != nil || lastErr != nil || first < 1 || first > last {
		return 0, 0, fmt.Errorf("line range %q is not <first>:<last> with ordered positive bounds", value)
	}
	return int(first), int(last), nil
}

// SelectLines answers physical lines first through last of data, inclusive and one-based,
// each with its original terminator. A physical line ends after its newline byte, and a last
// line without one still counts. An empty source has no lines, so a range beyond the last
// line refuses rather than clipping.
func SelectLines(data []byte, first, last int) ([]byte, error) {
	lines := bytes.SplitAfter(data, []byte("\n"))
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	if first < 1 || first > last || last > len(lines) {
		return nil, fmt.Errorf("line range %d:%d is outside the %d physical lines", first, last, len(lines))
	}
	return bytes.Join(lines[first-1:last], nil), nil
}
