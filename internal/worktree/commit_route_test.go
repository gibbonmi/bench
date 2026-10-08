package worktree

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/refusalroute"
)

// TestCommitExitThreeRouteReconcilesTheCheckout is RR35. A commit in an assignment worktree
// publishes, and its lane's stale index lock keeps the checkout from following, so the
// commit exits 3. The operator runs the printed reset plan and then the apply that the plan
// prints, and the checkout is clean at the published commit.
func TestCommitExitThreeRouteReconcilesTheCheckout(t *testing.T) {
	t.Parallel()
	f := admittedMergeFixture(t, "integration")
	target := f.created[0]
	lock := commitIndexLockLane(t, target.Path)
	mustWrite(t, filepath.Join(target.Path, "owned.txt"), []byte("owned\n"), 0o644)
	wrapper := installedWrapper(t, testRunBinary(t))
	var stdout, stderr bytes.Buffer
	commit := descendant(t, wrapper, "commit", "-m", "publish past the reconcile", "--", "owned.txt")
	commit.Dir, commit.Env, commit.Stdout, commit.Stderr = target.Path, journeyChildEnv(t, f.home), &stdout, &stderr
	code := exitCode(commit.Run())
	next, printed := recordField(stdout.String(), "committed{", refusalroute.NextField)
	if code != 3 || !printed || next == "" {
		t.Fatalf("commit = (%d, %q, %q), want exit 3 and a committed record with a next= route", code, stdout.String(), stderr.String())
	}
	published := gitOutput(t, f.root, "rev-parse", target.Assignment.Branch)
	// The lane's stale lock is the fixture's fault, and the reset needs the index lock too.
	mustRemove(t, lock)
	applyResetPlan(t, wrapper, f.repoHome, runPrintedStep(t, wrapper, f.repoHome, next))
	if head, status := gitOutput(t, target.Path, "rev-parse", "HEAD"), gitOutput(t, target.Path, "status", "--porcelain"); head != published || status != "" {
		t.Fatalf("after the printed route the checkout is at %s with status %q, want a clean checkout at the published %s", head, status, published)
	}
}
