package landing

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

func validRequest(r Request) error {
	if r.Root == "" || r.Destination == "" || r.Expected == "" || strings.TrimSpace(r.Message) == "" {
		return errors.New("landing request is incomplete")
	}
	if _, err := benchgit.ResolveCommit(r.Root, r.Expected); err != nil {
		return errors.New("expected base is not a commit")
	}
	return nil
}

// NamedPathError is a refusal that the caller clears by correcting a path it named, or the
// file at that path. It carries the cause, so a caller routes the refusal by the type and
// never by the sentence.
type NamedPathError struct{ Err error }

func (e NamedPathError) Error() string { return e.Err.Error() }
func (e NamedPathError) Unwrap() error { return e.Err }

// attributedPaths returns the attributed names of raw and the merge heads of the pending
// merge whose changed paths those names cover.
func attributedPaths(root, expected string, raw []string) ([]string, []string, error) {
	if len(raw) == 0 {
		return nil, nil, NamedPathError{errors.New("at least one path is required")}
	}
	paths := make([]string, 0, len(raw))
	for _, p := range raw {
		rel, err := repositoryPath(root, p)
		if err != nil {
			return nil, nil, NamedPathError{err}
		}
		if rel == "." || rel == "" {
			return nil, nil, NamedPathError{errors.New("repository root is not an attributed path")}
		}
		if err := safePath(root, expected, rel); err != nil {
			return nil, nil, NamedPathError{err}
		}
		paths = append(paths, rel)
	}
	named := unique(paths)
	heads, err := coverPendingMerge(root, expected, named)
	if err != nil {
		return nil, nil, err
	}
	return named, heads, nil
}

// coverPendingMerge refuses names that omit a path the pending merge changed: a merge commit
// without that path makes the later landing revert it. An unreadable merge state refuses as
// a plain error, because the caller cannot clear it by naming paths. It returns the merge
// heads that it read.
func coverPendingMerge(root, expected string, names []string) ([]string, error) {
	heads, err := benchgit.PendingMerge(root)
	if err != nil || len(heads) == 0 {
		return nil, err
	}
	merged, err := stagedPaths(root, expected, nil)
	if err != nil {
		return nil, err
	}
	covered, err := stagedPaths(root, expected, names)
	if err != nil {
		return nil, err
	}
	covers := make(map[string]bool, len(covered))
	for _, path := range covered {
		covers[path] = true
	}
	var omitted []string
	for _, path := range merged {
		if !covers[path] {
			omitted = append(omitted, strconv.Quote(path))
		}
	}
	if len(omitted) > 0 {
		return nil, NamedPathError{errors.New("pending merge changed paths that the commit does not name: " + strings.Join(omitted, " "))}
	}
	return heads, nil
}

// stagedPaths lists each index path that differs from base, limited to what names cover
// when names is not empty. Git's pathspec decides what a name covers, as it does for compose
// and reconcile, and the NUL framing keeps each path byte for byte.
func stagedPaths(root, base string, names []string) ([]string, error) {
	args := []string{"diff-index", "--cached", "--name-only", "--no-renames", "-z", base, "--"}
	for _, name := range names {
		args = append(args, ":(literal)"+name)
	}
	out, err := outputRaw(root, args...)
	if err != nil {
		return nil, fmt.Errorf("list the paths of the pending merge: %w", err)
	}
	// A NUL ends each path, so the last field is always empty.
	paths := strings.Split(string(out), "\x00")
	return paths[:len(paths)-1], nil
}

// ResolveAttributedPaths returns the safe repository-relative ownership fence for raw.
func ResolveAttributedPaths(root, expected string, raw []string) ([]string, error) {
	named, _, err := attributedPaths(root, expected, raw)
	return named, err
}

// CandidateTree composes the attributed working changes without modifying the real index.
func CandidateTree(root, expected string, paths []string) (string, error) {
	named, _, err := attributedPaths(root, expected, paths)
	if err != nil {
		return "", err
	}
	snapshot, err := compose(Request{Root: root, Expected: expected}, named)
	return snapshot.tree, err
}

func repositoryPath(root, path string) (string, error) {
	abs := path
	if !filepath.IsAbs(abs) {
		abs = filepath.Join(root, path)
	}
	abs, err := filepath.Abs(abs)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(root, abs)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path %q escapes repository", path)
	}
	return filepath.ToSlash(rel), nil
}

func safePath(root, expected, rel string) error {
	for current := filepath.Join(root, filepath.FromSlash(rel)); ; current = filepath.Dir(current) {
		info, err := os.Lstat(current)
		if err == nil && info.Mode().Perm()&0o444 == 0 {
			return fmt.Errorf("unreadable path %q is not attributable", rel)
		}
		if err == nil && special(info.Mode()) {
			return fmt.Errorf("special file %q is not attributable", rel)
		}
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("inspect %q: %w", rel, err)
		}
		if current == root {
			break
		}
	}
	path := filepath.Join(root, filepath.FromSlash(rel))
	info, err := os.Lstat(path)
	if err == nil && info.IsDir() && !gitlinkAt(root, expected, rel) {
		err = filepath.WalkDir(path, func(current string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.Type()&os.ModeSymlink != 0 {
				return nil
			}
			entryInfo, statErr := entry.Info()
			if statErr != nil {
				return statErr
			}
			if entryInfo.Mode().Perm()&0o444 == 0 || special(entryInfo.Mode()) {
				return errors.New("special or unreadable descendant")
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("inspect attributed directory %q: %w", rel, err)
		}
	}
	return nil
}

func special(mode os.FileMode) bool {
	return mode&os.ModeType != 0 && mode&os.ModeSymlink == 0 && !mode.IsDir()
}

func gitlinkAt(root, expected, path string) bool {
	for _, args := range [][]string{
		{"ls-tree", expected, "--", path},
		{"ls-files", "--stage", "--", path},
	} {
		line, err := output(root, args...)
		if err == nil && strings.HasPrefix(line, "160000 ") {
			return true
		}
	}
	return false
}

func compose(r Request, paths []string) (composedSnapshot, error) {
	dir, err := os.MkdirTemp("", "bench-landing-index-")
	if err != nil {
		return composedSnapshot{}, err
	}
	defer os.RemoveAll(dir)
	idx := filepath.Join(dir, "index")
	if err := indexRun(r.Root, idx, "read-tree", r.Expected); err != nil {
		return composedSnapshot{}, err
	}
	for _, path := range paths {
		if err := indexRun(r.Root, idx, "add", "-A", "--", ":(literal)"+path); err != nil {
			if !trackedAt(r.Root, r.Expected, path) {
				return composedSnapshot{}, NamedPathError{fmt.Errorf("named path %q not found in worktree, index, or expected base", path)}
			}
			return composedSnapshot{}, fmt.Errorf("compose attributed path %q: %w", path, err)
		}
	}
	var snapshot composedSnapshot
	snapshot.tree, err = indexOutput(r.Root, idx, "write-tree")
	return snapshot, err
}

func gitRegularFileMode(mode os.FileMode) string {
	if mode.Perm()&0o111 != 0 {
		return "100755"
	}
	return "100644"
}

func reconcile(r Request, paths []string, snapshot composedSnapshot) error {
	for _, path := range paths {
		literal := ":(literal)" + path
		failed := func(err error) error { return &ReconcileError{Path: path, Err: err} }
		if err := run(r.Root, "restore", "--source="+snapshot.tree, "--staged", "--worktree", "--", literal); err != nil {
			if resetErr := run(r.Root, "reset", "-q", r.Expected, "--", literal); resetErr != nil {
				return failed(err)
			}
			if retryErr := run(r.Root, "restore", "--source="+snapshot.tree, "--staged", "--worktree", "--", literal); retryErr != nil {
				return failed(retryErr)
			}
		}
		if err := run(r.Root, "clean", "-f", "-d", "--", literal); err != nil {
			return failed(err)
		}
		if err := run(r.Root, "diff", "--quiet", "--cached", snapshot.tree, "--", literal); err != nil {
			return failed(err)
		}
		if err := run(r.Root, "diff", "--quiet", "--ignore-submodules=dirty", snapshot.tree, "--", literal); err != nil {
			return failed(err)
		}
		untracked, err := output(r.Root, "ls-files", "--others", "--exclude-standard", "--", literal)
		if err != nil {
			return failed(err)
		}
		if untracked != "" {
			return failed(errors.New("named path still has untracked content"))
		}
	}
	return nil
}

func trackedAt(root, base, path string) bool {
	return run(root, "cat-file", "-e", base+":"+path) == nil
}
func unique(values []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, v := range values {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	sort.Strings(out)
	return out
}
