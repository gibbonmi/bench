package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// IsPrimaryCheckout reports whether root owns the repository's common Git directory.
// Linked worktrees have a private Git directory beneath that common directory.
func IsPrimaryCheckout(root string) (bool, error) {
	gitDir, err := AdminDir(root)
	if err != nil {
		return false, err
	}
	common, err := CommonDir(root)
	if err != nil {
		return false, err
	}
	return filepath.Clean(gitDir) == filepath.Clean(common), nil
}

// PendingMerge answers the merge heads that MERGE_HEAD records for the checkout at root,
// in file order. No error and no merge head mean no pending merge, and an error means an
// unreadable merge state.
func PendingMerge(root string) ([]string, error) {
	heads, err := mergeHeads(root)
	if err != nil {
		return nil, fmt.Errorf("read pending merge: %w", err)
	}
	return heads, nil
}

// mergeHeads reads only a regular MERGE_HEAD: a read of a FIFO blocks, and a followed
// symbolic link reads bytes outside the checkout administration directory.
func mergeHeads(root string) ([]string, error) {
	path, err := AdminPath(root, "MERGE_HEAD")
	if err != nil {
		return nil, err
	}
	info, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%q is not a regular file", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := strings.TrimSuffix(string(data), "\n")
	if content == "" {
		return nil, fmt.Errorf("%q is empty", path)
	}
	var heads []string
	for _, line := range strings.Split(content, "\n") {
		if !fullCommitID(root, line) {
			return nil, fmt.Errorf("%q holds the line %q, which is not a full commit ID", path, line)
		}
		heads = append(heads, line)
	}
	return heads, nil
}

// fullCommitID reports whether line is a commit ID exactly as rev-parse answers it, the
// only form in which Git writes a merge head.
func fullCommitID(root, line string) bool {
	if line == "" || strings.Trim(line, "0123456789abcdef") != "" {
		return false
	}
	commit, err := ResolveCommit(root, line)
	return err == nil && commit == line
}
