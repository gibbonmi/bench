package intent

import "github.com/gibbonmi/bench/internal/canonicalpath"

// AssignmentsOwning returns every assignment in the given slice whose worktree is the
// tree at path, in ledger order and in any state. The tree is matched by canonical path,
// because the ledger records a resolved path while a caller's spelling can arrive through
// a symlink. A row whose worktree no longer resolves owns nothing the caller can act on,
// so it is skipped, and an unresolvable path answers no owner at all.
//
// This is the one canonical-path match over assignments. The caller supplies the slice,
// because a caller that resolves a missing tree must still read the ledger from a root it
// can reach. The answer keeps every state, so a caller can refuse with the state named
// rather than report the tree as unassigned; the state filter belongs to the caller.
func AssignmentsOwning(assignments []Assignment, path string) []Assignment {
	canonical, err := canonicalpath.Resolve(path)
	if err != nil {
		return nil
	}
	var owning []Assignment
	for _, a := range assignments {
		if owned, ownedErr := canonicalpath.Resolve(a.Worktree); ownedErr == nil && owned == canonical {
			owning = append(owning, a)
		}
	}
	return owning
}
