package git

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
)

// IndexEntry is one record of `git ls-files --stage -z`: the entry's index mode and its
// repository-relative path. The mode is the six-digit octal Git stores, so a caller tells
// a regular file (100644 or 100755) from a symbolic link (120000) or a gitlink (160000)
// without a second read.
type IndexEntry struct {
	Mode string
	Path string
}

// regularFileModes are the index modes of a regular file. IsRegularFile is the one reader
// of the set, so the mode vocabulary has one source and no caller can widen it.
var regularFileModes = map[string]bool{"100644": true, "100755": true}

// IsRegularFile reports whether the entry holds regular-file content.
func (e IndexEntry) IsRegularFile() bool { return regularFileModes[e.Mode] }

// StagedIndex is one read of a repository's index. Entries is every index entry with its
// mode, which is what a policy validates its targets against. Staged is the subset whose
// content differs from HEAD, which is what a staged grade takes as its subjects. Both
// come from one read, so the two lists cannot disagree about a mode.
type StagedIndex struct {
	Entries []IndexEntry
	Staged  []IndexEntry
}

// ReadStagedIndex lists root's index entries and the subset that differs from HEAD.
//
// The changed names come from `diff --cached --name-only -z --diff-filter=ACMRT
// --no-renames`, so the added, copied, modified, renamed, and type-changed kinds count and
// a deletion does not. The modes come from `ls-files --stage -z`. Both streams are
// NUL-framed, so a path holding a space, a double quote, a tab, or a newline survives
// whole; Git C-quotes those paths in every non-`-z` form.
//
// On an unborn branch there is no HEAD to diff against, so every index entry is staged.
//
// The function refuses a root that is not the top of a working tree, because a staged
// grade names its subjects by a path relative to that top.
func ReadStagedIndex(root string) (StagedIndex, error) {
	if !IsWorkTreeTop(root) {
		return StagedIndex{}, fmt.Errorf("%q is not the top of a git working tree", root)
	}
	raw, err := Raw("-C", root, "ls-files", "--stage", "-z")
	if err != nil {
		return StagedIndex{}, fmt.Errorf("git ls-files %s: %w", root, err)
	}
	entries, err := parseIndexEntriesZ(raw)
	if err != nil {
		return StagedIndex{}, fmt.Errorf("git ls-files %s: %w", root, err)
	}
	index := StagedIndex{Entries: entries}
	if !OK("-C", root, "rev-parse", "--verify", "--quiet", "HEAD") {
		index.Staged = entries
		return index, nil
	}
	diff, err := Raw("-C", root, "diff", "--cached", "--name-only", "-z", "--diff-filter=ACMRT", "--no-renames")
	if err != nil {
		return StagedIndex{}, fmt.Errorf("git diff --cached %s: %w", root, err)
	}
	changed := make(map[string]bool)
	for _, name := range strings.Split(string(diff), "\x00") {
		if name != "" {
			changed[name] = true
		}
	}
	for _, entry := range entries {
		if changed[entry.Path] {
			index.Staged = append(index.Staged, entry)
		}
	}
	return index, nil
}

// IndexBlobReader starts `git show :<path>` in root and returns its stdout stream. The
// bytes are verbatim: a staged grade reads what a commit would carry, which is never the
// working file. The caller reads under its own bound and closes the stream, so a blob
// larger than that bound never reaches memory whole.
//
// Close closes the pipe and reaps the child, so a bounded read leaves no running git. A
// caller that stops short at its bound leaves the child writing into a closed pipe, and
// that exit is the bound rather than a fault; Close therefore reports an error only when
// the child failed before it delivered a byte, which is the missing-path case.
func IndexBlobReader(root, path string) (io.ReadCloser, error) {
	cmd := exec.Command("git", "-C", root, "show", ":"+path)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("git show :%s in %s: %w", path, root, err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("git show :%s in %s: %w", path, root, err)
	}
	return &indexBlobStream{cmd: cmd, stdout: stdout, root: root, path: path}, nil
}

// indexBlobStream is one running `git show` beside its stdout pipe. It counts the bytes it
// delivers, because that count is what tells a child that failed from a child the caller
// cut off at its own bound.
type indexBlobStream struct {
	cmd       *exec.Cmd
	stdout    io.ReadCloser
	root      string
	path      string
	delivered int64
}

func (s *indexBlobStream) Read(p []byte) (int, error) {
	n, err := s.stdout.Read(p)
	s.delivered += int64(n)
	return n, err
}

func (s *indexBlobStream) Close() error {
	// The pipe closes before the reap, so a child still writing sees the closed pipe and
	// exits rather than blocking the wait below forever.
	pipeErr := s.stdout.Close()
	if err := errors.Join(s.cmd.Wait(), pipeErr); err != nil && s.delivered == 0 {
		return fmt.Errorf("git show :%s in %s: %w", s.path, s.root, err)
	}
	return nil
}

// IsWorkTreeTop reports whether dir is the top of a Git working tree. A directory inside a
// repository is not the top, and a directory outside every repository is not either.
func IsWorkTreeTop(dir string) bool {
	top, err := RootAt(dir)
	if err != nil || top == "" {
		return false
	}
	return sameDirectory(top, dir)
}

// sameDirectory reports whether two directory operands name one directory. Both sides are
// made absolute first, because Git answers an absolute top while a caller may spell its
// root relative to the process directory: `.` names the top as often as the top's own
// path does. The resolved comparison follows, because a temporary directory is often
// reached through a symbolic link and Git answers the resolved path.
func sameDirectory(a, b string) bool {
	absA, absB := absolute(a), absolute(b)
	if absA == absB {
		return true
	}
	resolvedA, errA := filepath.EvalSymlinks(absA)
	resolvedB, errB := filepath.EvalSymlinks(absB)
	if errA != nil || errB != nil {
		return false
	}
	return resolvedA == resolvedB
}

// absolute makes path absolute against the process directory, and cleans it in place when
// the process directory is unavailable. A cleaned relative path never equals an absolute
// one, so the caller's comparison stays false rather than becoming accidentally true.
func absolute(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return filepath.Clean(path)
	}
	return abs
}

// parseIndexEntriesZ splits `ls-files --stage -z` output. Each record is
// `<mode> <object> <stage>\t<path>` and ends at a NUL, so the path is everything past the
// first tab and no byte of it is escaped. This function is the one source of that framing
// knowledge.
func parseIndexEntriesZ(raw []byte) ([]IndexEntry, error) {
	var entries []IndexEntry
	for offset := 0; offset < len(raw); {
		end := bytes.IndexByte(raw[offset:], 0)
		if end < 0 {
			return nil, errors.New("malformed index listing")
		}
		record := raw[offset : offset+end]
		offset += end + 1
		tab := bytes.IndexByte(record, '\t')
		if tab < 0 {
			return nil, errors.New("malformed index listing")
		}
		fields := strings.Fields(string(record[:tab]))
		if len(fields) != 3 {
			return nil, errors.New("malformed index listing")
		}
		entries = append(entries, IndexEntry{Mode: fields[0], Path: string(record[tab+1:])})
	}
	return entries, nil
}
