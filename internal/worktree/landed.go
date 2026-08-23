package worktree

import (
	"os"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

func assignmentLanded(assignment intent.Assignment, plan CleanupPlan) bool {
	return assignment.State == intent.StateActive &&
		plan.landedTyped.ProvenLanded() &&
		!planHasLiveLease(plan)
}

func planHasLiveLease(plan CleanupPlan) bool {
	if plan.ReasonCode == ReasonLiveLease {
		return true
	}
	lease, err := LeaseFile(plan.Target)
	if err != nil {
		return false
	}
	info, err := os.Lstat(lease)
	return err == nil && info.Mode().IsRegular() && ProbeLease(lease) == LeaseLive
}

func activeAssignmentWithMissingBranch(root, path string) (intent.Assignment, bool) {
	assignments, err := intent.Assignments(root)
	if err != nil {
		return intent.Assignment{}, false
	}
	registrations, err := git.Worktrees(root)
	if err != nil {
		return intent.Assignment{}, false
	}
	for _, assignment := range assignments {
		if assignment.State != intent.StateActive || !samePath(assignment.Worktree, path) {
			continue
		}
		registered := false
		for _, registration := range registrations {
			if samePath(registration.Path, path) && registration.BranchRef == assignment.Branch {
				registered = true
				break
			}
		}
		if registered && !git.OK("-C", root, "rev-parse", "--verify", "--quiet", assignment.Branch+"^{commit}") {
			return assignment, true
		}
	}
	return intent.Assignment{}, false
}
