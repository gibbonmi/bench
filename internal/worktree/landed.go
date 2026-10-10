package worktree

import (
	"errors"
	"os"
	"strings"

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
// tree. Such a branch holds only default-branch commits, so the ancestry proof alone calls
// the uncommitted output landed. A clean tree has no output to lose, and an assignment with
// an own commit has landed that commit. A missing tree holds nothing. The list cell reads
// this fact too, so no reader derives it twice.
func assignmentHoldsUnlandedOutput(root string, assignment intent.Assignment) (bool, error) {
	own, err := assignmentOwnsCommit(root, assignment)
	if err != nil || own {
		return false, err
	}
	if _, statErr := os.Lstat(assignment.Worktree); errors.Is(statErr, os.ErrNotExist) {
		return false, nil
	}
	status, err := checkoutStatus(assignment.Worktree)
	return len(status) > 0, err
}

// assignmentOwnsCommit reports whether the first-parent chain from the branch tip back to
// the recorded start holds a commit that is not on the default branch's first-parent
// history. A landing publishes a two-parent commit on the default branch, with the
// assignment tip as its second parent, so a landed own commit is never a first-parent
// default-branch commit. A fast-forward merge of the default branch moves the tip past the
// start over commits that are all first-parent default-branch commits. The first-parent
// chain also leaves out the side commits of a landing merge that the fast-forward takes in.
// The scoped landed selector reads this fact too, as its no-contribution test.
func assignmentOwnsCommit(root string, assignment intent.Assignment) (bool, error) {
	def, ok := git.ResolvedDefault(root)
	if !ok {
		return false, errors.New("the default branch does not resolve")
	}
	tip, err := git.ResolveCommit(root, assignment.Branch)
	if err != nil {
		return false, err
	}
	start, err := git.ResolveCommit(root, assignment.Start)
	if err != nil {
		return false, err
	}
	chain, err := git.Output("-C", root, "rev-list", "--first-parent", start+".."+tip)
	if err != nil || chain == "" {
		return false, err
	}
	history, err := git.Output("-C", root, "rev-list", "--first-parent", def, "^"+start)
	if err != nil {
		return false, err
	}
	onDefault := map[string]bool{}
	for _, commit := range strings.Fields(history) {
		onDefault[commit] = true
	}
	for _, commit := range strings.Fields(chain) {
		if !onDefault[commit] {
			return true, nil
		}
	}
	return false, nil
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
