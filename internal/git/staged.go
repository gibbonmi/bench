package git

import (
	"bytes"
	"errors"
	"fmt"
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

// RegularFileModes are the index modes of a regular file. A caller that grades file
// content keys on this set, so the mode vocabulary has one source.
var RegularFileModes = map[string]bool{"100644": true, "100755": true}

// IsRegularFile reports whether the entry holds regular-file content.
func (e IndexEntry) IsRegularFile() bool { return RegularFileModes[e.Mode] }

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

// IndexBlob returns the bytes the index holds for the repository-relative path, through
// `git show :<path>`. The bytes are verbatim: a staged grade reads what a commit would
// carry, which is never the working file.
func IndexBlob(root, path string) ([]byte, error) {
	out, err := Raw("-C", root, "show", ":"+path)
	if err != nil {
		return nil, fmt.Errorf("git show :%s in %s: %w", path, root, err)
	}
	return out, nil
}

// IsWorkTreeTop reports whether dir is the top of a Git working tree. A directory inside a
// repository is not the top, and a directory outside every repository is not either. The
// comparison resolves both sides, because a temporary directory is often reached through a
// symbolic link and Git answers the resolved path.
func IsWorkTreeTop(dir string) bool {
	top, err := RootAt(dir)
	if err != nil || top == "" {
		return false
	}
	return sameDirectory(top, dir)
}

func sameDirectory(a, b string) bool {
	if filepath.Clean(a) == filepath.Clean(b) {
		return true
	}
	resolvedA, errA := filepath.EvalSymlinks(a)
	resolvedB, errB := filepath.EvalSymlinks(b)
	if errA != nil || errB != nil {
		return false
	}
	return filepath.Clean(resolvedA) == filepath.Clean(resolvedB)
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
