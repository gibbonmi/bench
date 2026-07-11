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
//
// It lists every local head and applies the scratch-prefix filter in Go rather than as a
// git-side `worktree-*` glob: git's for-each-ref wildcard stops at a `/`, so a slashed scratch
// name (`worktree-agent-<hash>/x`) would escape the glob and never be swept.
func OrphanedDelegateBranches(root string) []string {
	branchesOut, err := git.Output("-C", root, "for-each-ref", "--format=%(refname:short)", "refs/heads/")
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

// ResolvedDefault returns the repo's default branch and whether it resolves to a
// commit. The sweep and the status board both refuse to classify orphans against a
// default that does not resolve — mergedness against a missing ref would read every
// branch as un-landed, or worse, let a silent all-clean sweep stand.
func ResolvedDefault(root string) (string, bool) {
	return git.ResolvedDefault(root)
}

// LandedInDefault reports whether branch's work is already present in def: by
// ancestry (merge-base) or by patch containment — `git cherry` reporting every
// commit on the branch as already applied (all `-`). byContent distinguishes the
// cherry proof so the sweep can report which proof landed the branch. Any git
// failure classifies as not landed: keeping a branch is recoverable, deleting
// it is not.
//
// The cherry proof only speaks for non-merge commits: `git cherry` never lists a
// merge, so a merge commit's own content (a conflict resolution) is invisible to
// the patch comparison. A branch carrying any merge commit def lacks is therefore
// unprovable by content and stays kept.
func LandedInDefault(root, branch, def string) (landed, byContent bool) {
	return git.LandedInDefault(root, branch, def)
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
