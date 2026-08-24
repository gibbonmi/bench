package landing

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func validRequest(r Request) error {
	if r.Root == "" || r.Destination == "" || r.Expected == "" || strings.TrimSpace(r.Message) == "" {
		return errors.New("landing request is incomplete")
	}
	if _, err := output(r.Root, "rev-parse", "--verify", r.Expected+"^{commit}"); err != nil {
		return errors.New("expected base is not a commit")
	}
	return nil
}

func attributedPaths(root, expected string, raw []string) ([]string, error) {
	if len(raw) == 0 {
		return nil, errors.New("at least one path is required")
	}
	paths := make([]string, 0, len(raw))
	for _, p := range raw {
		rel, err := repositoryPath(root, p)
		if err != nil {
			return nil, err
		}
		if rel == "." || rel == "" {
			return nil, errors.New("repository root is not an attributed path")
		}
		if err := safePath(root, expected, rel); err != nil {
			return nil, err
		}
		paths = append(paths, rel)
	}
	return unique(paths), nil
}

// ResolveAttributedPaths returns the safe repository-relative ownership fence for raw.
func ResolveAttributedPaths(root, expected string, raw []string) ([]string, error) {
	return attributedPaths(root, expected, raw)
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
				return composedSnapshot{}, fmt.Errorf("named path %q not found in worktree, index, or expected base", path)
			}
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

// removeIndexTree drops every entry beneath rel from the prospective index, so the
// published tree carries the deletion rather than the checkout carrying it afterwards.
// The pathspec is literal and the removals name exact index paths, so a folder name
// holding a space or a glob character resolves to itself.
func removeIndexTree(root, idx, rel string) error {
	listed, err := indexOutputRaw(root, idx, "ls-files", "-z", "--cached", "--", ":(literal)"+rel)
	if err != nil {
		return fmt.Errorf("list tracked entries under %q: %w", rel, err)
	}
	for _, path := range strings.Split(string(listed), "\x00") {
		if path == "" {
			continue
		}
		if err := indexRun(root, idx, "update-index", "--force-remove", "--", path); err != nil {
			return fmt.Errorf("remove %q from prospective index: %w", path, err)
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
