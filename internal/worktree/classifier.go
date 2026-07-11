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
	facts, err := git.Worktrees(root)
	if err != nil {
		return nil, err
	}
	mainRoot := canonicalRoot(root)
	return classifyRegistered(mainRoot, facts), nil
}

func canonicalRoot(root string) string {
	common, err := git.Output("-C", root, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil || filepath.Base(common) != ".git" {
		return root
	}
	return filepath.Dir(common)
}

func classifyRegistered(mainRoot string, facts []git.Worktree) []Registered {
	out := make([]Registered, 0, len(facts))
	for _, fact := range facts {
		out = append(out, Registered{Path: fact.Path, Branch: fact.Branch, Detached: fact.Detached, Locked: fact.Locked})
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
