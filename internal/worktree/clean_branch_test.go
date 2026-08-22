package worktree

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

// unprovableLandedAssignment returns one owned, cleanup-pending assignment whose branch
// commits were composed into a single commit on the default branch. One commit is a mode
// change on a file that survives the landing. Every byte of the branch is in the default
// branch, and still no derived proof can say so. Ancestry fails, no merge exists, and the
// composed commit shares no patch-id with either original. Reverse-apply also refuses the
// mode change on a surviving entry, because git apply reports a preimage mode mismatch as
// a warning rather than a failure. This is the exact case the operator override exists for.
func unprovableLandedAssignment(t *testing.T, request string) (string, Creation) {
	t.Helper()
	root, creation := newOwnedAssignment(t, request)
	commitInWorktree(t, creation.Path, "one.txt", "one\n", "one")
	if err := os.Chmod(filepath.Join(creation.Path, "tracked.txt"), 0o755); err != nil {
		t.Fatal(err)
	}
	gitRun(t, creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-q", "-a", "-m", "make tracked executable")
	short := strings.TrimPrefix(creation.Assignment.Branch, "refs/heads/")
	gitRun(t, root, "cherry-pick", "--no-commit", creation.Assignment.Start+".."+short)
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "squashed")
	markPending(t, root, creation.Assignment)
	return root, creation
}

func branchExists(root, ref string) bool {
	return exec.Command("git", "-C", root, "show-ref", "--verify", "--quiet", ref).Run() == nil
}

func recoveryRefs(t *testing.T, root string, assignment intent.Assignment) string {
	t.Helper()
	return gitOutput(t, root, "for-each-ref", "--format=%(refname)", intent.RecoveryRefPrefix(assignment.OwnerID, assignment.ID))
}

func renderedPlan(t *testing.T, plan CleanupPlan) string {
	t.Helper()
	var out bytes.Buffer
	mustNoError(t, renderCleanup(&out, plan))
	return out.String()
}

// [RW1]
func TestDiscardBranchRetiresTheBranchAndLeavesNoRecoveryRef(t *testing.T) {
	root, creation := unprovableLandedAssignment(t, "rw1")
	options := CleanupOptions{DiscardBranch: true}
	plan, err := PlanExplicitWithOptions(root, creation.Path, options)
	mustNoError(t, err)
	requireTest(t, plan.Action == ActionRemove, "plan action = %q, want remove", plan.Action)
	requireTest(t, plan.deleteBranch, "plan did not authorize branch deletion under the override")
	requireTest(t, plan.branchRef == creation.Assignment.Branch, "plan branch ref = %q, want %q", plan.branchRef, creation.Assignment.Branch)
	requireTest(t, plan.Recovery == "none", "plan recovery = %q, want none", plan.Recovery)

	applied, err := ApplyExplicitWithOptions(root, creation.Path, plan.Fingerprint, options)
	mustNoError(t, err)
	requireTest(t, applied.Action == ActionRemoved, "applied action = %q, want removed", applied.Action)
	_, statErr := os.Lstat(creation.Path)
	requireTest(t, os.IsNotExist(statErr), "checkout survived the apply: %v", statErr)
	requireTest(t, !branchExists(root, creation.Assignment.Branch), "assignment branch %s survived the apply", creation.Assignment.Branch)
	refs := recoveryRefs(t, root, creation.Assignment)
	requireTest(t, refs == "", "apply left recovery refs behind: %q", refs)
}

// [RW2]
func TestDiscardBranchPlanNamesTheBranchBeforeAnyRemoval(t *testing.T) {
	root, creation := unprovableLandedAssignment(t, "rw2")
	plan, err := PlanExplicitWithOptions(root, creation.Path, CleanupOptions{DiscardBranch: true})
	mustNoError(t, err)
	output := renderedPlan(t, plan)
	requireTest(t, strings.Contains(output, "discards branch "+creation.Assignment.Branch), "plan output does not name the branch: %q", output)
	requireTest(t, strings.Contains(output, ",none,"), "plan output does not name the recovery ref: %q", output)

	_, statErr := os.Lstat(creation.Path)
	requireTest(t, statErr == nil, "planning removed the checkout: %v", statErr)
	requireTest(t, branchExists(root, creation.Assignment.Branch), "planning deleted the branch")
}

// [RW3] The override is an argument to the explicit path only. This is the regression
// guard on that boundary: it passed before the override existed and must keep passing.
func TestDiscardBranchLeavesTheDerivedClassificationUnchanged(t *testing.T) {
	t.Run("automatic path retains an unproven branch and authorizes no deletion", func(t *testing.T) {
		root, creation := unprovableLandedAssignment(t, "rw3-automatic")
		plan, err := PlanAutomatic(root, creation.Path)
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRetain, "automatic action = %q, want retain", plan.Action)
		requireTest(t, plan.ReasonCode == ReasonUnmerged, "automatic reason = %q, want unmerged", plan.ReasonCode)
		requireTest(t, !plan.deleteBranch, "automatic plan authorized branch deletion")
		requireTest(t, plan.branchRef == "", "automatic plan named a deletable branch: %q", plan.branchRef)
	})
	t.Run("explicit path without the override removes only the checkout", func(t *testing.T) {
		root, creation := unprovableLandedAssignment(t, "rw3-explicit")
		plan, err := PlanExplicit(root, creation.Path)
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRemove, "explicit action = %q, want remove", plan.Action)
		requireTest(t, !plan.deleteBranch && plan.branchRef == "", "explicit plan authorized branch deletion: %t %q", plan.deleteBranch, plan.branchRef)
		output := renderedPlan(t, plan)
		requireTest(t, strings.HasSuffix(output, ",apply with exact fingerprint\n"), "explicit plan detail changed: %q", output)

		applied, err := ApplyExplicit(root, creation.Path, plan.Fingerprint)
		mustNoError(t, err)
		requireTest(t, applied.Action == ActionRemoved, "explicit applied action = %q, want removed", applied.Action)
		requireTest(t, branchExists(root, creation.Assignment.Branch), "cleanup without the override deleted branch %s", creation.Assignment.Branch)
	})
}

// [RW4]
func TestDiscardBranchNeverBypassesARefusal(t *testing.T) {
	options := CleanupOptions{DiscardBranch: true}
	t.Run("primary checkout", func(t *testing.T) {
		root, _ := unprovableLandedAssignment(t, "rw4-primary")
		plan, err := PlanExplicitWithOptions(root, root, options)
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRetain, "primary action = %q, want retain", plan.Action)
		requireTest(t, plan.Reason == "primary checkout is never removable", "primary reason = %q", plan.Reason)
		applied, err := ApplyExplicitWithOptions(root, root, plan.Fingerprint, options)
		mustNoError(t, err)
		requireTest(t, applied.Action == ActionRetain, "apply acted on the primary checkout: %q", applied.Action)
		requireTest(t, branchExists(root, "refs/heads/main"), "apply deleted the default branch")
	})
	t.Run("path outside any registration", func(t *testing.T) {
		root := newWorktreeRepo(t)
		outside := filepath.Join(t.TempDir(), "unregistered")
		mustMkdirAll(t, outside, 0o700)
		plan, err := PlanExplicitWithOptions(root, outside, options)
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRetain, "unregistered action = %q, want retain", plan.Action)
		requireTest(t, plan.ReasonCode == ReasonForeign && plan.Reason == "target is not registered", "unregistered reason = %q/%q", plan.ReasonCode, plan.Reason)
		applied, err := ApplyExplicitWithOptions(root, outside, plan.Fingerprint, options)
		mustNoError(t, err)
		requireTest(t, applied.Action == ActionRetain, "apply acted on an unregistered path: %q", applied.Action)
		_, statErr := os.Lstat(outside)
		requireTest(t, statErr == nil, "apply removed an unregistered path: %v", statErr)
	})
	t.Run("foreign lock", func(t *testing.T) {
		root := newWorktreeRepo(t)
		target := filepath.Join(t.TempDir(), "foreign locked")
		gitRun(t, root, "worktree", "add", "-q", "-b", "foreign-locked", target, "HEAD")
		gitRun(t, root, "worktree", "lock", "--reason", "foreign", target)
		plan, err := PlanExplicitWithOptions(root, target, options)
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRetain, "foreign action = %q, want retain", plan.Action)
		requireTest(t, plan.ReasonCode == ReasonUnexpectedLock, "foreign reason = %q", plan.ReasonCode)
		applied, err := ApplyExplicitWithOptions(root, target, plan.Fingerprint, options)
		mustNoError(t, err)
		requireTest(t, applied.Action == ActionRetain, "apply acted on a foreign worktree: %q", applied.Action)
		requireTest(t, branchExists(root, "refs/heads/foreign-locked"), "apply deleted a foreign branch")
		_, statErr := os.Lstat(target)
		requireTest(t, statErr == nil, "apply removed a foreign worktree: %v", statErr)
	})
	t.Run("assignment identity mismatch", func(t *testing.T) {
		root, creation := unprovableLandedAssignment(t, "rw4-identity")
		gitRun(t, creation.Path, "switch", "-q", "-c", "drifted-identity")
		plan, err := PlanExplicitWithOptions(root, creation.Path, options)
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRetain, "mismatch action = %q, want retain", plan.Action)
		requireTest(t, plan.Reason == "assignment does not match current branch", "mismatch reason = %q", plan.Reason)
		applied, err := ApplyExplicitWithOptions(root, creation.Path, plan.Fingerprint, options)
		mustNoError(t, err)
		requireTest(t, applied.Action == ActionRetain, "apply acted on a mismatched assignment: %q", applied.Action)
		requireTest(t, branchExists(root, creation.Assignment.Branch), "apply deleted the assignment branch behind an identity refusal")
		requireTest(t, branchExists(root, "refs/heads/drifted-identity"), "apply deleted the checked-out branch behind an identity refusal")
		_, statErr := os.Lstat(creation.Path)
		requireTest(t, statErr == nil, "apply removed a mismatched checkout: %v", statErr)
	})
}

// A detached HEAD has no branch for the DiscardBranch override to name. A registered,
// detached-HEAD worktree planned with DiscardBranch leaves deleteBranch and branchRef
// at their zero values. This contrasts with the attached-branch cases above, which do
// carry branch-deletion authority.
func TestDiscardBranchLeavesADetachedHeadUnaffected(t *testing.T) {
	root := newWorktreeRepo(t)
	target := filepath.Join(filepath.Dir(root), "detached discard target")
	gitRun(t, root, "worktree", "add", "-q", "--detach", target, "HEAD")
	plan, err := PlanExplicitWithOptions(root, target, CleanupOptions{DiscardBranch: true})
	mustNoError(t, err)
	requireTest(t, !plan.deleteBranch, "detached-HEAD plan authorized branch deletion under the override")
	requireTest(t, plan.branchRef == "", "detached-HEAD plan named a deletable branch: %q", plan.branchRef)
	requireTest(t, plan.Action == ActionRecoverRemove, "detached-HEAD action = %q, want recover-remove — the two checks above only bind while the fixture still reaches a removing verdict", plan.Action)
}
