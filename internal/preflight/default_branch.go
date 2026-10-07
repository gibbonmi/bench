package preflight

import "github.com/gibbonmi/bench/internal/git"

// baseCurrentFacts backs the base-current check: it resolves the default
// branch and reports whether its tip is an ancestor of HEAD:
// merge-base(default, HEAD) equal to rev-parse(default).
//
// The resolved name comes back with the two predicates, because the stale-base
// remedy names that branch and git.ResolvedDefault is its one source.
//
// An unresolved default branch answers no name and both predicates false. The
// check itself renders that as red without a separate bootstrap failure, since
// map #7 names this a per-check red rather than a bootstrap precondition.
func baseCurrentFacts(root string) (branch string, resolved, current bool) {
	def, ok := git.ResolvedDefault(root)
	if !ok {
		return "", false, false
	}
	mergeBase, err1 := git.Output("merge-base", def, "HEAD")
	tip, err2 := git.Output("rev-parse", def)
	if err1 != nil || err2 != nil {
		return def, true, false
	}
	return def, true, mergeBase == tip
}

// defaultTipEqualPaths is the changed subset whose bytes and mode at sourceTip equal those
// at the default-branch tip. git.ResolvedDefault owns the default branch, and the raw tree
// diff owns the comparison, so a gitlink compares by its commit too. A path absent at both
// tips reads as equal, and an untracked path is absent at both. So a dirty source, a
// missing tip, or an unresolved default branch answers none, and the fence rule grades
// every path.
func defaultTipEqualPaths(root, sourceTip string, changed []string) ([]string, *BootstrapFailure) {
	branch, resolved := git.ResolvedDefault(root)
	if !resolved || sourceTip == "" || len(changed) == 0 {
		return nil, nil
	}
	dirty, err := git.WorktreeDirty(root)
	if err != nil {
		return nil, &BootstrapFailure{"source status unreadable", err.Error()}
	}
	if dirty {
		return nil, nil
	}
	// The commit identity keeps a file named like the branch from making the diff ambiguous.
	defaultTip, err := git.ResolveCommit(root, branch)
	if err != nil {
		return nil, &BootstrapFailure{"default-branch tip not readable", err.Error()}
	}
	differences, err := git.TreeChangesIncludingSubmodules(root, defaultTip, sourceTip)
	if err != nil {
		return nil, &BootstrapFailure{"default-branch tree not readable", err.Error()}
	}
	differs := make(map[string]bool, len(differences))
	for _, change := range differences {
		differs[change.Path] = true
	}
	var equal []string
	for _, path := range changed {
		if !differs[path] {
			equal = append(equal, path)
		}
	}
	return equal, nil
}
