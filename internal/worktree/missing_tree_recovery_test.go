// Missing-tree recovery tests: an assignment whose tree is gone refuses every verb that
// resolves it with the reset-tree-missing face, and the route that face prints clears the
// record. A landed assignment retires through the clean of the landed set, and an unlanded
// one through its own release.
package worktree

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/refusalroute"
)

// requireOutOfMissingTree requires that the reset no longer refuses the assignment with
// the missing-tree face, because the route cleared the record.
func requireOutOfMissingTree(t *testing.T, f ownedAssignment) {
	t.Helper()
	r := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID))
	if strings.Contains(r.stdout, "worktree tree is missing") {
		t.Fatalf("reset after the route = (%d, %q), want the missing-tree refusal cleared", r.exit, r.stdout)
	}
}

// TestCleanLandedRetiresAMissingTree is the landed arm: the plan of the landed set removes
// the record of a landed assignment whose tree is gone, and its apply retires it.
func TestCleanLandedRetiresAMissingTree(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "missing-landed")
	landAssignment(t, f.root, f.creation, "landed.txt")
	mustNoError(t, os.RemoveAll(f.creation.Path))
	plan := runVerb(t, verbClean, f.call("--landed"))
	if plan.exit != 0 || strings.Contains(plan.stdout, ",retain,") {
		t.Fatalf("clean --landed plan = (%d, %q, %q), want a plan that retires the record", plan.exit, plan.stdout, plan.stderr)
	}
	apply := runVerb(t, verbClean, f.call("--landed", "--apply", plan.mustFingerprint(t)))
	if apply.exit != 0 {
		t.Fatalf("clean --landed apply = (%d, %q, %q), want exit 0", apply.exit, apply.stdout, apply.stderr)
	}
	if registered, err := registeredAt(f.root, f.creation.Path); err != nil || registered {
		t.Fatalf("registration after the apply = %t, %v; want it released", registered, err)
	}
	// The landed branch goes, as the retirement of a present landed checkout deletes it.
	if git.OK("-C", f.root, "show-ref", "--verify", "--quiet", f.creation.Assignment.Branch) {
		t.Fatalf("landed branch %q remains after the apply", f.creation.Assignment.Branch)
	}
	requireOutOfMissingTree(t, f)
}

// TestReleaseReleasesAnUnlandedMissingTree is the unlanded arm: the release of an
// assignment whose tree is gone releases the record and never runs inside the tree.
func TestReleaseReleasesAnUnlandedMissingTree(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "missing-unlanded")
	commitInWorktree(t, f.creation.Path, "work.txt", "work\n", "unlanded work")
	mustNoError(t, os.RemoveAll(f.creation.Path))
	r := runVerb(t, verbRelease, f.call("--request", f.creation.Assignment.RequestToken, f.creation.Path))
	if r.exit != 0 {
		t.Fatalf("release = (%d, %q, %q), want exit 0", r.exit, r.stdout, r.stderr)
	}
	// The unlanded branch keeps the work that only it holds.
	if !git.OK("-C", f.root, "show-ref", "--verify", "--quiet", f.creation.Assignment.Branch) {
		t.Fatalf("unlanded branch %q is gone after the release", f.creation.Assignment.Branch)
	}
	requireOutOfMissingTree(t, f)
}

// TestReleaseRetainsAMissingTreeUnderALiveLease: a live lease still holds the assignment
// when its tree is gone, so the release retains it as it retains a present tree.
func TestReleaseRetainsAMissingTreeUnderALiveLease(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "missing-leased")
	lease, err := LeaseFile(f.creation.Path)
	mustNoError(t, err)
	mustWrite(t, lease, []byte(strconv.Itoa(os.Getpid())+" 2026-10-07T00:00:00Z\n"), 0o600)
	mustNoError(t, os.RemoveAll(f.creation.Path))
	r := runVerb(t, verbRelease, f.call("--request", f.creation.Assignment.RequestToken, f.creation.Path))
	if r.exit == 0 || !strings.Contains(r.stderr, "("+string(ReasonLiveLease)+")") {
		t.Fatalf("release = (%d, %q, %q), want the live-lease retention", r.exit, r.stdout, r.stderr)
	}
	if registered, err := registeredAt(f.root, f.creation.Path); err != nil || !registered {
		t.Fatalf("registration after the refusal = %t, %v; want it kept", registered, err)
	}
}

// TestMissingTreeRouteOfAnUnsafePathIsItsPlaceholder: a missing tree whose path is not
// line-safe prints the recovery slot's placeholder, so no control byte reaches the record.
func TestMissingTreeRouteOfAnUnsafePathIsItsPlaceholder(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	gone := intent.Assignment{ID: "gone", RequestToken: "req-gone", Branch: "refs/heads/bench/gone", Worktree: filepath.Join(t.TempDir(), "gone\nline")}
	var refused refusalError
	if err := missingTreeRefusal(root, gone); !errors.As(err, &refused) || refused.detail != "worktree tree is missing" || refused.next != "<"+refusalroute.FactRecovery+">" {
		t.Fatalf("missing-tree refusal = %v, want the missing-tree sentence and the recovery placeholder", err)
	}
}
