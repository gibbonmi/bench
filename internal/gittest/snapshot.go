package gittest

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/visibletree"
)

// CommitWorkingTree copies tracked and visible untracked files into destination. It
// commits that exact snapshot, so tests can authenticate the tree they are grading.
func CommitWorkingTree(source, destination string) error {
	visible, err := visibletree.Files(source)
	if err != nil {
		return err
	}
	for _, relative := range visible {
		relative = filepath.FromSlash(relative)
		from, to := filepath.Join(source, relative), filepath.Join(destination, relative)
		info, err := os.Lstat(from)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(to), 0o755); err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(from)
			if err == nil {
				err = os.Symlink(target, to)
			}
			if err != nil {
				return err
			}
			continue
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("source entry is not regular: %s", relative)
		}
		data, err := os.ReadFile(from)
		if err == nil {
			err = os.WriteFile(to, data, info.Mode().Perm())
		}
		if err != nil {
			return err
		}
	}
	for _, args := range append([][]string{{"init", "-q"}}, IdentityConfig()...) {
		if _, err := Run(destination, args...); err != nil {
			return err
		}
	}
	return CommitAll(destination, "authenticated test snapshot")
}

// CommitAll authenticates the current visible working tree in an initialized test repo.
// It force-adds every path, unlike Commit, which commits the index as the caller left it.
func CommitAll(root, message string) error {
	for _, args := range [][]string{{"add", "-f", "."}, {"commit", "-qm", message}} {
		if _, err := Run(root, args...); err != nil {
			return err
		}
	}
	return nil
}
