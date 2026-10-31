package preflight

import (
	"errors"

	"github.com/gibbonmi/bench/internal/git"
)

// baseCurrentFacts backs the base-current check: it resolves the default
// branch and reports whether its tip is an ancestor of HEAD:
// merge-base(default, HEAD) equal to rev-parse(default).
//
// The resolved name comes back with the two predicates, because the stale-base
// remedy names that branch and git.DefaultTip is its one source.
//
// An unresolved default branch answers no name and both predicates false. The
// check itself renders that as red without a separate bootstrap failure, since
// map #7 names this a per-check red rather than a bootstrap precondition.
func baseCurrentFacts(root string) (branch string, resolved, current bool) {
	def, tip, err := git.DefaultTip(root)
	if errors.Is(err, git.ErrDefaultUnresolved) {
		return "", false, false
	}
	if err != nil {
		return def, true, false
	}
	mergeBase, err := git.Output("-C", root, "merge-base", tip, "HEAD")
	if err != nil {
		return def, true, false
	}
	return def, true, mergeBase == tip
}

// defaultTipEqualPaths is the changed subset whose bytes and mode at sourceTip equal those
// at the default-branch tip. git.DefaultTip owns the default-branch tip, and
// unchangedPaths owns the comparison and its dirty-source rule. A missing tip or an
// unresolved default branch answers none, and the fence rule grades every path.
func defaultTipEqualPaths(root, sourceTip string, changed []string) ([]string, *BootstrapFailure) {
	if sourceTip == "" || len(changed) == 0 {
		return nil, nil
	}
	// The commit identity keeps a file named like the branch from making the diff ambiguous.
	_, defaultTip, err := git.DefaultTip(root)
	if errors.Is(err, git.ErrDefaultUnresolved) {
		return nil, nil
	}
	if err != nil {
		return nil, &BootstrapFailure{"default-branch tip not readable", err.Error()}
	}
	return unchangedPaths(root, defaultTip, sourceTip, changed)
}
