package harnesstranscript

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"

	"github.com/gibbonmi/bench/internal/bounds"
)

// This file owns one question: how a record file becomes a stream of lines and a digest. The
// question is the same for every source, so it sits beside the source mappings rather than
// inside one of them. A mapping decides only what a line means.

// maxRecordLine bounds one line of the record. A rollout line holds a whole tool result, so
// the bound sits far above an ordinary line and only a line no reader could use reaches it.
// A line past the bound is one skipped event: the reader discards that line, marks the counts
// incomplete, and keeps reading, so a later turn, compaction, or timestamp is not lost with
// it. It is a variable so the package test can lower it instead of building a fixture of the
// production size.
var maxRecordLine = 16 << 20

// lineReader accumulates one record's lines into a source-specific tally.
type lineReader interface {
	line(raw []byte)
	skipped()
}

// readRecord reads the record at path once, feeds every line it can read to tally, and
// returns the digest of the whole file.
//
// The path is typed before it is opened. A FIFO would block the reader inside open(2), and a
// link would supply bytes from somewhere other than the path the agent named, so both are
// refused on the stat that precedes every open. The record is then streamed a line at a time,
// because a session record outgrows any whole-file bound, and the digest is taken from the
// same pass, so the identity names exactly the bytes the counts came from.
//
// Record contents are data throughout. The reader decodes and counts. It never executes a
// field, resolves a path out of one, or copies record text into output.
func readRecord(path string, tally lineReader) (string, Failure) {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", Failure{State: bounds.StateAbsent, Reason: "no record at the named path"}
		}
		return "", Failure{State: bounds.StateUnreadable, Reason: err.Error()}
	}
	if !info.Mode().IsRegular() {
		return "", Failure{State: bounds.StateWrongType, Reason: "not a regular file: " + info.Mode().Type().String()}
	}
	file, err := os.Open(path)
	if err != nil {
		return "", Failure{State: bounds.StateUnreadable, Reason: err.Error()}
	}
	defer file.Close()

	digest := sha256.New()
	reader := bufio.NewReaderSize(io.TeeReader(file, digest), 64<<10)
	for {
		line, oversized, err := readRecordLine(reader)
		if oversized {
			tally.skipped()
		} else if len(line) > 0 {
			tally.line(line)
		}
		if err != nil {
			if err != io.EOF {
				return "", Failure{State: bounds.StateUnreadable, Reason: err.Error()}
			}
			break
		}
	}
	return "sha256:" + hex.EncodeToString(digest.Sum(nil)), Failure{}
}

// readRecordLine returns the next line of the record. A line past maxRecordLine is read to its
// end and discarded, and the second result reports it, so the reader keeps its place and the
// later lines still count. Every byte still passes the digest, whether or not the line it
// belongs to was kept.
func readRecordLine(reader *bufio.Reader) ([]byte, bool, error) {
	var line []byte
	oversized := false
	for {
		chunk, isPrefix, err := reader.ReadLine()
		if !oversized && len(line)+len(chunk) > maxRecordLine {
			oversized, line = true, nil
		}
		if !oversized {
			line = append(line, chunk...)
		}
		if err != nil {
			return line, oversized, err
		}
		if !isPrefix {
			return line, oversized, nil
		}
	}
}
