package worktree

import (
	"errors"
	"os"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

// assignmentLanded is the one landed verdict of an active assignment. The branch has to be
// in the default ref, no live lease may hold the tree, and the tree may not hold unlanded
// output. The list banner, the selector of the landed set, the resume sweep, and the
// automatic classifier all read it.
func assignmentLanded(root string, assignment intent.Assignment, plan CleanupPlan) bool {
	if assignment.State != intent.StateActive || !plan.landedTyped.ProvenLanded() || planHasLiveLease(plan) {
		return false
	}
	holds, err := assignmentHoldsUnlandedOutput(root, assignment)
	return err == nil && !holds
}

// assignmentHoldsUnlandedOutput reports an assignment with no commit of its own and a dirty
// tree. Its branch tip sits at its recorded start, which is a default-ref commit whatever the
// tree holds, so the ancestry proof alone calls the uncommitted output landed. A clean tree
// has no output to lose, and an assignment with an own commit has landed that commit. A
// missing tree holds nothing. The list cell reads this fact too, so no reader derives it twice.
func assignmentHoldsUnlandedOutput(root string, assignment intent.Assignment) (bool, error) {
	tip, err := git.ResolveCommit(root, assignment.Branch)
	if err != nil {
		return false, err
	}
	start, err := git.ResolveCommit(root, assignment.Start)
	if err != nil {
		return false, err
	}
	if tip != start {
		return false, nil
	}
	if _, statErr := os.Lstat(assignment.Worktree); errors.Is(statErr, os.ErrNotExist) {
		return false, nil
	}
	status, err := checkoutStatus(assignment.Worktree)
	return len(status) > 0, err
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
	// The match answers every state, and the scan keeps its own active filter: a retired
	// record names a tree whose phase is over, so its gone branch is the expected end
	// state rather than the recoverable loss this scan reports.
	for _, assignment := range intent.AssignmentsOwning(assignments, path) {
		if assignment.State != intent.StateActive {
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
