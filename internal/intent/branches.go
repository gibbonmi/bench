package intent

import "github.com/gibbonmi/bench/internal/git"

// PruneUnclaimedLandedBranches is Repo.PruneUnclaimedLandedBranches with a nil run.
func PruneUnclaimedLandedBranches(root string) (int, error) {
	return In(nil, root).PruneUnclaimedLandedBranches()
}

// PruneUnclaimedLandedBranches protects every recorded assignment branch while
// delegating landedness and exact ref deletion to the Git owner.
func (r Repo) PruneUnclaimedLandedBranches() (int, error) {
	assignments, err := r.Assignments()
	if err != nil {
		return 0, err
	}
	protected := make([]string, 0, len(assignments))
	for _, assignment := range assignments {
		protected = append(protected, assignment.Branch)
	}
	return git.PruneLandedBranches(r.Root, protected)
}
