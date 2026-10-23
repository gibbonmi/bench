package worktree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

// A restore after later commits keeps the branch tip and merges the preserved staged,
// unstaged, and untracked work onto it.
func TestResetRestoreKeepsCommitsMadeAfterTheEnvelope(t *testing.T) {
	t.Parallel()
	f := restoreFixture(t)
	commitInWorktree(t, f.creation.Path, "later-one", "one\n", "later one")
	commitInWorktree(t, f.creation.Path, "later-two", "two\n", "later two")
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	plan := runVerb(t, verbReset, f.call("--restore", f.ref, f.creation.Assignment.ID))
	requireTest(t, isRestorePlanOf(plan, f.ref), "restore plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	result := runVerb(t, verbReset, f.call("--restore", f.ref, f.creation.Assignment.ID, "--apply", plan.mustFingerprint(t)))
	requireTest(t, result.exit == 0 && strings.Contains(result.stdout, "ref="+f.creation.Assignment.Branch) && !strings.Contains(result.stdout, "next="),
		"restore after later commits = %d %s %s", result.exit, result.stdout, result.stderr)
	requireTest(t, gitOutput(t, f.root, "rev-parse", f.creation.Assignment.Branch) == tip && gitOutput(t, f.creation.Path, "rev-parse", "HEAD") == tip &&
		gitOutput(t, f.creation.Path, "symbolic-ref", "HEAD") == f.creation.Assignment.Branch, "restore dropped the later commits or detached")
	requireTest(t, gitOutput(t, f.creation.Path, "status", "--porcelain=v1") == "MM README.md\n?? untracked", "restore lost staged or untracked status")
	requireTest(t, gitOutput(t, f.creation.Path, "show", ":README.md") == "staged", "restore lost staged bytes")
	body, err := os.ReadFile(filepath.Join(f.creation.Path, "README.md"))
	mustNoError(t, err)
	requireTest(t, string(body) == "working\n", "restore lost working bytes: %q", body)
}

// A staged edit that a later commit already carries restores as a clean index.
func TestResetRestoreDropsAStagedLayerThatALaterCommitCarries(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "restore-staged-carried")
	mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("carried\n"), 0o644)
	gitRun(t, f.creation.Path, "add", "README.md")
	mustWrite(t, filepath.Join(f.creation.Path, "untracked"), []byte("untracked\n"), 0o644)
	plan := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	result := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID, "--apply", plan.mustFingerprint(t)))
	requireTest(t, result.exit == 0, "fixture reset = %d %s %s", result.exit, result.stdout, result.stderr)
	commitInWorktree(t, f.creation.Path, "README.md", "carried\n", "carry the staged edit")
	ref := intent.ResetRefPrefix(f.creation.Assignment.OwnerID, f.creation.Assignment.ID) + "1"
	plan = runVerb(t, verbReset, f.call("--restore", ref, f.creation.Assignment.ID))
	requireTest(t, isRestorePlanOf(plan, ref), "restore plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	result = runVerb(t, verbReset, f.call("--restore", ref, f.creation.Assignment.ID, "--apply", plan.mustFingerprint(t)))
	requireTest(t, result.exit == 0 && gitOutput(t, f.creation.Path, "status", "--porcelain=v1") == "?? untracked",
		"carried staged restore = %d %s %s", result.exit, result.stdout, result.stderr)
}

// A preserved edit that conflicts with a later commit refuses the plan and moves nothing.
func TestResetRestoreRefusesAPayloadThatConflictsWithLaterCommits(t *testing.T) {
	t.Parallel()
	f := restoreFixture(t)
	commitInWorktree(t, f.creation.Path, "README.md", "later\n", "later readme")
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	result := runVerb(t, verbReset, f.call("--restore", f.ref, f.creation.Assignment.ID))
	requireTest(t, result.exit == 1 && strings.Contains(result.stdout, "envelope payload conflicts with the branch tip"),
		"conflicting restore = %d %s %s", result.exit, result.stdout, result.stderr)
	requireTest(t, gitOutput(t, f.root, "rev-parse", f.creation.Assignment.Branch) == tip && gitOutput(t, f.creation.Path, "rev-parse", "HEAD") == tip,
		"conflict refusal moved the checkout")
}

// A branch that diverged from the envelope tip refuses the restore, because the move to
// the envelope tip would drop the branch's own commits.
func TestResetRestoreRefusesABranchThatDivergedFromTheEnvelope(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "restore-diverged")
	commitInWorktree(t, f.creation.Path, "assignment", "assignment\n", "assignment tip")
	mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("dirty\n"), 0o644)
	plan := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	result := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID, "--apply", plan.mustFingerprint(t)))
	requireTest(t, result.exit == 0, "fixture reset = %d %s %s", result.exit, result.stdout, result.stderr)
	commitInWorktree(t, f.creation.Path, "diverged", "diverged\n", "diverged tip")
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	ref := intent.ResetRefPrefix(f.creation.Assignment.OwnerID, f.creation.Assignment.ID) + "1"
	result = runVerb(t, verbReset, f.call("--restore", ref, f.creation.Assignment.ID))
	requireTest(t, result.exit == 1 && strings.Contains(result.stdout, "restore would drop commits on the assignment branch"),
		"diverged restore = %d %s %s", result.exit, result.stdout, result.stderr)
	requireTest(t, gitOutput(t, f.root, "rev-parse", f.creation.Assignment.Branch) == tip, "diverged refusal moved the branch")
}
