package bounds

import (
	"fmt"
	"io/fs"
	"os"
	"unicode/utf8"
)

// FileState is the vocabulary every control-record read reports itself in. Absence is
// the only authoritative empty state: every other value says something went wrong, so a
// surface can never mistake a failed read for a repository with nothing to report.
type FileState string

const (
	StateAbsent     FileState = "absent"
	StateEmpty      FileState = "empty"
	StateParsed     FileState = "parsed"
	StateMalformed  FileState = "malformed"
	StateUnreadable FileState = "unreadable"
	StateWrongType  FileState = "wrong-type"
	// StateUnsupportedSchema belongs to the parsers, not to this package: it means the
	// bytes read cleanly but their structure is not the document the consumer expects,
	// and only a parser owns that predicate. It is declared here so every parser names
	// one vocabulary.
	StateUnsupportedSchema FileState = "unsupported-schema"
)

// Failed reports whether a state means the read yielded nothing a consumer may trust:
// the path could not be opened, it was not the kind of thing a control record is, or its
// bytes are not valid UTF-8. Absence and emptiness are authoritative answers rather than
// failures, so they are not included. Every surface's fail-closed exit and every
// dashboard `unknown` row keys on this one predicate, so a new state cannot be honored
// by some surfaces and fall through the default branch of others. StateUnsupportedSchema
// is deliberately outside it: the classifier never returns it, and each parser decides
// what an unrecognized shape costs its own command.
func (s FileState) Failed() bool {
	switch s {
	case StateUnreadable, StateWrongType, StateMalformed:
		return true
	}
	return false
}

// Classified is one path's state. Data carries the bytes whenever the read completed,
// Stream exposes the underlying read outcome so an oversized record stays distinguishable
// from a permission failure, and Reason carries the underlying diagnostic for every state
// that reports a problem.
type Classified struct {
	State  FileState
	Stream ReadStatus
	Data   []byte
	Reason string
}

// ClassifiedDir is one directory's state, in the same vocabulary. An unreadable directory
// and an empty one both yield no entries, so only State separates them.
type ClassifiedDir struct {
	State   FileState
	Entries []fs.DirEntry
	Reason  string
}

// Classify reports the state of a control record at path, reading at most limit bytes.
func Classify(path string, limit int64) Classified {
	info, state, reason := resolve(path)
	if state != "" {
		return Classified{State: state, Reason: reason}
	}
	return gradeBytes(path, info, limit)
}

// ClassifyDir reports the state of a control directory at path.
func ClassifyDir(path string) ClassifiedDir {
	info, state, reason := resolve(path)
	if state != "" {
		return ClassifiedDir{State: state, Reason: reason}
	}
	if !info.IsDir() {
		return ClassifiedDir{State: StateWrongType, Reason: fmt.Sprintf("not a directory: %s", info.Mode().Type())}
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return ClassifiedDir{State: StateUnreadable, Reason: err.Error()}
	}
	if len(entries) == 0 {
		return ClassifiedDir{State: StateEmpty, Entries: entries}
	}
	return ClassifiedDir{State: StateParsed, Entries: entries}
}

// resolve stats path without opening it, returning the target's info for a caller to
// type-check, or the state that ends the classification. Lstat comes first because a read
// of a dangling symlink reports ENOENT, which would answer authoritative absence for a
// link whose target is gone: only a path with no link at all is absent, so once Lstat has
// succeeded a vanished target is unreadable. The link itself is then followed, because a
// control record linked elsewhere is still a control record.
func resolve(path string) (fs.FileInfo, FileState, string) {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, StateAbsent, ""
		}
		return nil, StateUnreadable, err.Error()
	}
	if info.Mode()&fs.ModeSymlink != 0 {
		target, err := os.Stat(path)
		if err != nil {
			return nil, StateUnreadable, err.Error()
		}
		info = target
	}
	return info, "", ""
}

// ClassifyNoFollow reports the state of a producer file at path without ever following
// a link or opening a non-regular object, reading at most ControlRecordLimit. It sits
// beside Classify rather than replacing it because the two answer one input
// differently on purpose: a control record behind a live link is still that record, but
// a producer file is authoritative input to generated output, so a link there redirects
// the generator at bytes outside the tree it is grading. Both link forms are therefore
// wrong-type, and the type check precedes every open so a FIFO cannot block the gate in
// open(2). The limit is not a parameter: one producer read under two bounds is the
// divergence ControlRecordLimit exists to forbid.
func ClassifyNoFollow(path string) Classified {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Classified{State: StateAbsent}
		}
		return Classified{State: StateUnreadable, Reason: err.Error()}
	}
	return gradeBytes(path, info, ControlRecordLimit)
}

// gradeBytes is what both public forms mean by "read this and say what it is". The two
// differ only in how they resolve the path — followed or refused — and in which limit
// they bind; everything from the regular-file check onward is one fact about how bounded
// bytes are graded, and splitting it would let the oversized or UTF-8 disposition drift
// between the forms one edit at a time. info must already be the caller's resolved
// stat: this helper does not decide whether a link is a control record.
func gradeBytes(path string, info fs.FileInfo, limit int64) Classified {
	if !info.Mode().IsRegular() {
		return Classified{State: StateWrongType, Reason: fmt.Sprintf("not a regular file: %s", info.Mode().Type())}
	}
	file, err := os.Open(path)
	if err != nil {
		return Classified{State: StateUnreadable, Reason: err.Error()}
	}
	defer file.Close()
	read := Read(file, limit)
	switch read.Status {
	case ReadOversized, ReadFailed:
		return Classified{State: StateUnreadable, Stream: read.Status, Reason: read.Err.Error()}
	}
	if !utf8.Valid(read.Data) {
		return Classified{State: StateMalformed, Stream: read.Status, Data: read.Data, Reason: "invalid UTF-8"}
	}
	if len(read.Data) == 0 {
		return Classified{State: StateEmpty, Stream: read.Status, Data: read.Data}
	}
	return Classified{State: StateParsed, Stream: read.Status, Data: read.Data}
}
