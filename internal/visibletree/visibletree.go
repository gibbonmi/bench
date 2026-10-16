// Package visibletree owns the one definition of "the tree" for a check that grades files
// in place: every tracked file still present, plus every untracked file that git's
// standard excludes do not ignore. An ignored path is local state outside the gate's
// subject tree, so a check that read it would grade content no commit can carry.
//
// It is a leaf that imports only the standard library, so production checks, internal/git,
// and the test-repository scaffolds can all share it without an import cycle.
package visibletree

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Files lists root's visible files as sorted slash paths relative to root. Outside a
// repository nothing is ignored, so every file below root except a .git entry is listed.
func Files(root string) ([]string, error) {
	if !insideWorkTree(root) {
		return allFiles(root)
	}
	listed, err := runGit(root, nil, "ls-files", "-z", "--cached", "--others", "--exclude-standard")
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	var files []string
	for _, record := range bytes.Split(listed, []byte{0}) {
		rel := string(record)
		// An untracked nested repository lists as its directory with a trailing slash;
		// its content belongs to another repository's tree, not this one. A conflicted
		// path lists once per stage, and a tracked file deleted from the checkout lists
		// although no bytes remain to grade.
		if rel == "" || strings.HasSuffix(rel, "/") || seen[rel] {
			continue
		}
		seen[rel] = true
		if _, err := os.Lstat(filepath.Join(root, filepath.FromSlash(rel))); errors.Is(err, fs.ErrNotExist) {
			continue
		} else if err != nil {
			return nil, err
		}
		files = append(files, rel)
	}
	sort.Strings(files)
	return files, nil
}

// Ignored reports which of rels, slash paths relative to root, git's standard excludes
// ignore. A tracked path is never ignored. Outside a repository nothing is ignored.
func Ignored(root string, rels ...string) (map[string]bool, error) {
	ignored := map[string]bool{}
	if len(rels) == 0 || !insideWorkTree(root) {
		return ignored, nil
	}
	listed, err := runGit(root, strings.NewReader(strings.Join(rels, "\x00")+"\x00"), "check-ignore", "-z", "--stdin")
	// check-ignore exits 1 when it matched none of the paths; only a higher exit is a fault.
	var exitErr *exec.ExitError
	if err != nil && !(errors.As(err, &exitErr) && exitErr.ExitCode() == 1) {
		return nil, err
	}
	for _, record := range bytes.Split(listed, []byte{0}) {
		if rel := strings.TrimSuffix(string(record), "/"); rel != "" {
			ignored[rel] = true
		}
	}
	return ignored, nil
}

func runGit(root string, stdin *strings.Reader, args ...string) ([]byte, error) {
	var out bytes.Buffer
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	if stdin != nil {
		cmd.Stdin = stdin
	}
	cmd.Stdout = &out
	err := cmd.Run()
	return out.Bytes(), err
}

func insideWorkTree(root string) bool {
	inside, err := runGit(root, nil, "rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(string(inside)) == "true"
}

func allFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Name() == ".git" && path != root {
			if entry.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files = append(files, filepath.ToSlash(rel))
		return nil
	})
	sort.Strings(files)
	return files, err
}
