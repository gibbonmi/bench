package worktree

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
)

const delegateBranchPrefix = "worktree-"

// OrphanedDelegateBranches returns local harness scratch branches (`worktree-*`) that
// are not attached to any registered git worktree. Bench shift review branches use a
// different namespace (`bench/shift-*`) and are deliberately preserved.
//
// It lists every local head and applies the scratch-prefix filter in Go rather than as a
// git-side `worktree-*` glob: git's for-each-ref wildcard stops at a `/`, so a slashed scratch
// name (`worktree-agent-<hash>/x`) would escape the glob and never be swept.
func OrphanedDelegateBranches(root string) ([]string, error) {
	branches, err := git.LocalBranches(root)
	if err != nil {
		return nil, fmt.Errorf("list local branches: %w", err)
	}
	registered, err := git.Worktrees(root)
	if err != nil {
		return nil, fmt.Errorf("list registered worktrees: %w", err)
	}
	active := map[string]bool{}
	for _, worktree := range registered {
		if worktree.Branch != "" {
			active[worktree.Branch] = true
		}
	}
	var orphans []string
	for _, branch := range branches {
		if branch == "" || !strings.HasPrefix(branch, delegateBranchPrefix) || active[branch] {
			continue
		}
		orphans = append(orphans, branch)
	}
	sort.Strings(orphans)
	return orphans, nil
}
