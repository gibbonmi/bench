package intent

import "github.com/gibbonmi/bench/internal/git"

// PruneUnclaimedLandedBranches protects every recorded assignment branch while
// delegating landedness and exact ref deletion to the Git owner.
func PruneUnclaimedLandedBranches(root string) (int, error) {
	assignments, err := Assignments(root)
	if err != nil {
		return 0, err
	}
	protected := make([]string, 0, len(assignments))
	for _, assignment := range assignments {
		protected = append(protected, assignment.Branch)
	}
	return git.PruneLandedBranches(root, protected)
}
