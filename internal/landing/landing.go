// Package landing owns exact prospective Git-tree composition and publication.
package landing

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/gate/authorization"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/spec"
)

// Request is the complete, immutable input to one prospective landing.
type Request struct {
	Root, Destination, Expected, Message string
	Paths                                []string
	Spec                                 string
	Stdout, Stderr                       io.Writer
}

// Result identifies a successfully published commit and its authorized tree.
type Result struct{ Base, Commit, Tree string }

type composedSnapshot struct {
	tree            string
	specPath        string
	specPermissions os.FileMode
}

// Owner composes only Request.Paths from Request.Expected, authorizes that tree, then
// publishes it by expected-old update. The function fields are narrow operational seams
// for deterministic fault coverage; New supplies the real owners.
type Owner struct {
	authorize func(context.Context, string, string, io.Writer, io.Writer) authorization.Result
	updateRef func(string, string, string, string) error
	reconcile func(Request, []string, composedSnapshot) error
}

// New returns the production landing owner.
func New() Owner {
	return Owner{authorize: authorization.AuthorizeWithWriters, updateRef: updateRef, reconcile: reconcile}
}

// Land publishes a commit only when the exact composed tree receives green authorization.
func (o Owner) Land(ctx context.Context, r Request) (Result, error) {
	if err := validRequest(r); err != nil {
		return Result{}, err
	}
	paths, err := attributedPaths(r.Root, r.Expected, r.Paths)
	if err != nil {
		return Result{}, err
	}
	if r.Spec != "" {
		resolved, err := spec.CheckStaged(r.Root, r.Spec)
		if err != nil {
			return Result{}, err
		}
		rel, err := repositoryPath(r.Root, resolved)
		if err != nil {
			return Result{}, err
		}
		paths = append(paths, rel)
	}
	paths = unique(paths)
	snapshot, err := compose(r, paths)
	if err != nil {
		return Result{}, err
	}
	tree := snapshot.tree
	baseTree, err := output(r.Root, "rev-parse", r.Expected+"^{tree}")
	if err != nil {
		return Result{}, fmt.Errorf("read expected base tree: %w", err)
	}
	if tree == baseTree {
		return Result{}, errors.New("nothing to commit")
	}
	if got := o.authorize(ctx, r.Root, tree, r.Stdout, r.Stderr); got.Kind != authorization.Green {
		return Result{}, fmt.Errorf("prospective authorization refused: %s", got.Kind)
	}
	commit, err := output(r.Root, "commit-tree", tree, "-p", r.Expected, "-m", r.Message)
	if err != nil {
		return Result{}, fmt.Errorf("create landing commit: %w", err)
	}
	if err := o.updateRef(r.Root, r.Destination, commit, r.Expected); err != nil {
		return Result{}, fmt.Errorf("destination compare-and-swap refused: %w", err)
	}
	if err := o.reconcile(r, paths, snapshot); err != nil {
		return Result{Base: r.Expected, Commit: commit, Tree: tree}, fmt.Errorf("landed-but-checkout-incomplete: %w", err)
	}
	return Result{Base: r.Expected, Commit: commit, Tree: tree}, nil
}

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
	snapshot := composedSnapshot{}
	if r.Spec != "" {
		resolved, content, fileMode, err := transitionedSpec(r.Root, r.Spec)
		if err != nil {
			return composedSnapshot{}, err
		}
		rel, _ := repositoryPath(r.Root, resolved)
		blob, err := outputInput(r.Root, content, "hash-object", "-w", "--stdin")
		if err != nil {
			return composedSnapshot{}, err
		}
		if err := indexRun(r.Root, idx, "update-index", "--add", "--cacheinfo", gitRegularFileMode(fileMode)+","+blob+","+rel); err != nil {
			return composedSnapshot{}, err
		}
		snapshot.specPath = resolved
		snapshot.specPermissions = fileMode
	}
	snapshot.tree, err = indexOutput(r.Root, idx, "write-tree")
	return snapshot, err
}

func gitRegularFileMode(mode os.FileMode) string {
	if mode.Perm()&0o111 != 0 {
		return "100755"
	}
	return "100644"
}

func transitionedSpec(root, slug string) (string, []byte, os.FileMode, error) {
	resolved, err := spec.CheckStaged(root, slug)
	if err != nil {
		return "", nil, 0, err
	}
	content, err := os.ReadFile(resolved)
	if err != nil {
		return "", nil, 0, err
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", nil, 0, err
	}
	content, err = spec.Implemented(content)
	if err != nil {
		return "", nil, 0, err
	}
	return resolved, content, info.Mode().Perm(), nil
}

func reconcile(r Request, paths []string, snapshot composedSnapshot) error {
	for _, path := range paths {
		literal := ":(literal)" + path
		if err := run(r.Root, "restore", "--source="+snapshot.tree, "--staged", "--worktree", "--", literal); err != nil {
			if resetErr := run(r.Root, "reset", "-q", r.Expected, "--", literal); resetErr != nil {
				return err
			}
			if retryErr := run(r.Root, "restore", "--source="+snapshot.tree, "--staged", "--worktree", "--", literal); retryErr != nil {
				return retryErr
			}
		}
		if err := run(r.Root, "clean", "-f", "-d", "--", literal); err != nil {
			return err
		}
		if err := run(r.Root, "diff", "--quiet", "--cached", snapshot.tree, "--", literal); err != nil {
			return err
		}
		if err := run(r.Root, "diff", "--quiet", "--ignore-submodules=dirty", snapshot.tree, "--", literal); err != nil {
			return err
		}
		untracked, err := output(r.Root, "ls-files", "--others", "--exclude-standard", "--", literal)
		if err != nil {
			return err
		}
		if untracked != "" {
			return fmt.Errorf("named path %q still has untracked content", path)
		}
	}
	if snapshot.specPath != "" {
		return os.Chmod(snapshot.specPath, snapshot.specPermissions)
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
func updateRef(root, ref, new, old string) error { return run(root, "update-ref", ref, new, old) }
func output(root string, args ...string) (string, error) {
	return benchgit.Output(append([]string{"-C", root}, args...)...)
}
func outputInput(root string, input []byte, args ...string) (string, error) {
	c := exec.Command("git", append([]string{"-C", root}, args...)...)
	c.Stdin = strings.NewReader(string(input))
	b, err := c.Output()
	return strings.TrimSpace(string(b)), err
}
func run(root string, args ...string) error {
	return exec.Command("git", append([]string{"-C", root}, args...)...).Run()
}
func indexRun(root, idx string, args ...string) error {
	c := exec.Command("git", append([]string{"-C", root}, args...)...)
	c.Env = append(os.Environ(), "GIT_INDEX_FILE="+idx)
	return c.Run()
}
func indexOutput(root, idx string, args ...string) (string, error) {
	c := exec.Command("git", append([]string{"-C", root}, args...)...)
	c.Env = append(os.Environ(), "GIT_INDEX_FILE="+idx)
	b, err := c.Output()
	return strings.TrimSpace(string(b)), err
}
