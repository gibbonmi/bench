package gocache

import (
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Bound is the footprint above which a report warns: 10 GiB in bytes. It is one
// constant, and there is no knob. Disk pressure is a machine fact, so an operator
// never tunes it per repository.
const Bound int64 = 10_737_418_240

// trimFile is the name Go writes the build cache's last trim time to, as unix seconds.
const trimFile = "trim.txt"

// trimReadLimit caps the trim file read. The file holds one integer, so a larger file is
// a corrupt file and not a value to parse.
const trimReadLimit = 64

// Footprint is the measured state of the Bench build cache: the apparent byte total, the
// regular-file count, and the last trim time.
type Footprint struct {
	Dir      string
	Bytes    int64
	Files    int64
	LastTrim string
}

// OverBound reports whether the measured bytes are above Bound.
func (f Footprint) OverBound() bool { return f.Bytes > Bound }

// Measure walks dir and returns its footprint. The walk reads directory entries and calls
// lstat; it opens no cache file and follows no symlink. So a FIFO or a dangling link in
// the tree is harmless. It recurses into every subdirectory, a `-d` executable directory
// included, because Go stores a cached executable as a directory. It counts regular files
// alone and sums their apparent sizes.
//
// An absent directory measures as zero bytes in zero files, which is the state of a
// machine that has not run a Go child yet. An unreadable entry is skipped rather than
// refused, because a partial count still answers the operator's question.
func Measure(dir string) Footprint {
	footprint := Footprint{Dir: dir, LastTrim: lastTrim(dir)}
	_ = filepath.WalkDir(dir, func(_ string, entry fs.DirEntry, err error) error {
		if err != nil {
			// A refused directory or a vanished entry ends that branch alone. SkipDir on
			// the root would end the walk, and fs.SkipDir is the only value that means
			// "this branch"; returning nil here is the equivalent for a failed entry.
			return nil
		}
		if !entry.Type().IsRegular() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return nil
		}
		footprint.Bytes += info.Size()
		footprint.Files++
		return nil
	})
	return footprint
}

// lastTrim reads dir's trim file and renders its unix seconds as UTC RFC 3339. It answers
// an empty string when the file is absent, is not a regular file, or does not parse. The
// lstat regular-file check comes before the read, so a symlink is not followed and a FIFO
// is not opened.
func lastTrim(dir string) string {
	path := filepath.Join(dir, trimFile)
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() {
		return ""
	}
	body, err := os.ReadFile(path)
	if err != nil || len(body) > trimReadLimit {
		return ""
	}
	seconds, err := strconv.ParseInt(strings.TrimSpace(string(body)), 10, 64)
	if err != nil {
		return ""
	}
	return time.Unix(seconds, 0).UTC().Format(time.RFC3339)
}
