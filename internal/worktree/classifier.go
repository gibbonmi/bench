package worktree

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
)

type Class string

const (
	ClassRoot      Class = "root"
	ClassPoolWarm  Class = "pool-warm"
	ClassPoolLease Class = "pool-leased"
	ClassOutOfPool Class = "out-of-pool"
)

type Registered struct {
	Path     string
	Class    Class
	Branch   string
	Detached bool
	Locked   bool
}

// ClassifyRegisteredWorktrees returns every worktree `git worktree list` knows about,
// classified by pool membership. It returns the git query's error rather than
// swallowing it into an empty slice: every caller must distinguish "no worktrees" from
// "the classify query itself failed" so a git failure can never read as a silent
// all-clear.
func ClassifyRegisteredWorktrees(root string) ([]Registered, error) {
	out, err := git.Output("-C", root, "worktree", "list", "--porcelain")
	if err != nil {
		return nil, err
	}
	mainRoot := canonicalRoot(root)
	return classifyRegistered(mainRoot, out), nil
}

func canonicalRoot(root string) string {
	common, err := git.Output("-C", root, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil || filepath.Base(common) != ".git" {
		return root
	}
	return filepath.Dir(common)
}

func classifyRegistered(mainRoot, porcelain string) []Registered {
	var out []Registered
	var current *Registered
	for _, line := range strings.Split(porcelain, "\n") {
		switch {
		case strings.HasPrefix(line, "worktree "):
			out = append(out, Registered{Path: strings.TrimPrefix(line, "worktree ")})
			current = &out[len(out)-1]
		case current != nil && strings.HasPrefix(line, "branch refs/heads/"):
			current.Branch = strings.TrimPrefix(line, "branch refs/heads/")
		case current != nil && line == "detached":
			current.Detached = true
		case current != nil && (line == "locked" || strings.HasPrefix(line, "locked ")):
			current.Locked = true
		}
	}
	pool := Pool(mainRoot)
	for i := range out {
		out[i].Class = classifyPath(mainRoot, pool, out[i].Path)
	}
	return out
}

func classifyPath(root, pool, path string) Class {
	if samePath(path, root) {
		return ClassRoot
	}
	if insidePool(pool, path) {
		lease, _ := LeaseFile(path)
		if isRegularFile(lease) {
			return ClassPoolLease
		}
		return ClassPoolWarm
	}
	return ClassOutOfPool
}

func insidePool(pool, path string) bool {
	rel, err := filepath.Rel(pool, path)
	return err == nil && rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func samePath(a, b string) bool {
	ar, aerr := filepath.EvalSymlinks(a)
	br, berr := filepath.EvalSymlinks(b)
	if aerr == nil {
		a = ar
	}
	if berr == nil {
		b = br
	}
	return filepath.Clean(a) == filepath.Clean(b)
}

func isRegularFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}
