package worktree

import (
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
)

const delegateBranchPrefix = "worktree-"

// OrphanedDelegateBranches returns local harness scratch branches (`worktree-*`) that
// are not attached to any registered git worktree. Bench shift review branches use a
// different namespace (`bench/shift-*`) and are deliberately preserved.
func OrphanedDelegateBranches(root string) []string {
	branchesOut, err := git.Output("-C", root, "for-each-ref", "--format=%(refname:short)", "refs/heads/"+delegateBranchPrefix+"*")
	if err != nil || branchesOut == "" {
		return nil
	}
	active := activeWorktreeBranches(root)
	var orphans []string
	for _, branch := range strings.Split(branchesOut, "\n") {
		if branch == "" || !strings.HasPrefix(branch, delegateBranchPrefix) || active[branch] {
			continue
		}
		orphans = append(orphans, branch)
	}
	sort.Strings(orphans)
	return orphans
}

func activeWorktreeBranches(root string) map[string]bool {
	out, err := git.Output("-C", root, "worktree", "list", "--porcelain")
	if err != nil || out == "" {
		return nil
	}
	active := map[string]bool{}
	for _, line := range strings.Split(out, "\n") {
		const prefix = "branch refs/heads/"
		if strings.HasPrefix(line, prefix) {
			active[line[len(prefix):]] = true
		}
	}
	return active
}
