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
	"github.com/gibbonmi/bench/internal/sanitize"
)

// missingTreeResetFixtures are the producing fixtures of the missing-tree face, one for each
// cause: the branch landed or not, and Git pruned the registration of the gone tree or not.
// The walk runs each printed route verbatim, and the route leaves the assignment inactive.
func missingTreeResetFixtures() []resetRefusalFixture {
	var fixtures []resetRefusalFixture
	for _, landed := range []bool{true, false} {
		for _, pruned := range []bool{false, true} {
			cause := map[bool]string{true: "landed", false: "unlanded"}[landed]
			if pruned {
				cause += " pruned"
			}
			fixtures = append(fixtures, resetRefusalFixture{
				face:  faceResetTreeMissing,
				cause: cause,
				build: func(t *testing.T) (ownedAssignment, []string) {
					f, args := ownedResetFixture(t, "missing")
					if landed {
						landAssignment(t, f.root, f.creation, "landed.txt")
					} else {
						commitInWorktree(t, f.creation.Path, "work.txt", "work\n", "unlanded work")
					}
					mustNoError(t, os.RemoveAll(f.creation.Path))
					if pruned {
						// Git prunes no locked registration, and the owned one is locked. A
						// foreign registration keeps the pool of the landed cause, so the
						// prune leaves the pool with no entry for the tree or no pool at all.
						if landed {
							gitRun(t, f.root, "worktree", "add", "-q", "--detach", filepath.Join(t.TempDir(), "foreign"))
						}
						gitRun(t, f.root, "worktree", "unlock", f.creation.Path)
						gitRun(t, f.root, "worktree", "prune")
						if registered, err := registeredAt(f.root, f.creation.Path); err != nil || registered {
							t.Fatalf("registration after the prune = %t, %v; want it pruned", registered, err)
						}
					}
					return f, args
				},
				// A landed route is the plan of the landed set, so the operator applies it.
				after: func(t *testing.T, wrapper string, f ownedAssignment, last verbResult) {
					if landed {
						runPrintedStep(t, wrapper, f.repoHome, "bench worktree clean --landed --apply "+last.mustFingerprint(t))
					}
					if assignmentActive(t, f.root, f.creation.Assignment.ID) {
						t.Fatalf("assignment %q is still active after its route", f.creation.Assignment.ID)
					}
				},
			})
		}
	}
	return fixtures
}

// holdLiveLease writes the lease of a live owner, this test process, for the assignment.
func holdLiveLease(t *testing.T, f ownedAssignment) {
	t.Helper()
	lease, err := LeaseFile(f.creation.Path)
	mustNoError(t, err)
	mustWrite(t, lease, []byte(strconv.Itoa(os.Getpid())+" 2026-10-07T00:00:00Z\n"), 0o600)
}

// requireOutOfMissingTree requires that the reset no longer refuses the assignment with
// the missing-tree face, because the route cleared the record.
func requireOutOfMissingTree(t *testing.T, f ownedAssignment) {
	t.Helper()
	r := runVerb(t, verbReset, f.call(resetToStart(f)...))
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
	holdLiveLease(t, f)
	mustNoError(t, os.RemoveAll(f.creation.Path))
	r := runVerb(t, verbRelease, f.call("--request", f.creation.Assignment.RequestToken, f.creation.Path))
	if r.exit == 0 || !strings.Contains(r.stderr, "("+string(ReasonLiveLease)+")") {
		t.Fatalf("release = (%d, %q, %q), want the live-lease retention", r.exit, r.stdout, r.stderr)
	}
	if registered, err := registeredAt(f.root, f.creation.Path); err != nil || !registered {
		t.Fatalf("registration after the refusal = %t, %v; want it kept", registered, err)
	}
}

// TestReleaseRetainsAMissingTreeUnderAnUnknownLease: a lease that no reader can parse
// leaves the lease state unknown, so the release of a gone tree retains the assignment with
// the uncertain reason and names the release that runs once the lease is repaired.
func TestReleaseRetainsAMissingTreeUnderAnUnknownLease(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "missing-unknown-lease")
	lease, err := LeaseFile(f.creation.Path)
	mustNoError(t, err)
	mustWrite(t, lease, []byte("not a lease\n"), 0o600)
	mustNoError(t, os.RemoveAll(f.creation.Path))
	r := runVerb(t, verbRelease, f.call("--request", f.creation.Assignment.RequestToken, f.creation.Path))
	detail := "worktree retained (" + string(ReasonUncertain) + "): " + unknownLeaseReason
	next := "next=bench worktree release --request <request> " + sanitize.ShellQuote(f.creation.Path)
	if r.exit == 0 || !strings.Contains(r.stderr, detail) || !strings.Contains(r.stderr, next) {
		t.Fatalf("release = (%d, %q, %q), want %q and %q", r.exit, r.stdout, r.stderr, detail, next)
	}
	if registered, err := registeredAt(f.root, f.creation.Path); err != nil || !registered {
		t.Fatalf("registration after the refusal = %t, %v; want it kept", registered, err)
	}
}

// TestMissingTreeRoutePrintsAPlaceholderInEachUnsafeSlot: a missing tree whose request or
// path is not line-safe prints the release with its command words, and the placeholder only
// in the slot of each unsafe value, so no control byte reaches the record.
func TestMissingTreeRoutePrintsAPlaceholderInEachUnsafeSlot(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	unsafePath := filepath.Join(t.TempDir(), "gone\nline")
	for _, tc := range []struct{ request, want string }{
		{"req-gone", "bench worktree release --request req-gone <checkout>"},
		{"req\x1bgone", "bench worktree release --request <request> <checkout>"},
	} {
		gone := intent.Assignment{ID: "gone", RequestToken: tc.request, Branch: "refs/heads/bench/gone", Worktree: unsafePath}
		var refused refusalError
		if err := missingTreeRefusal(root, gone); !errors.As(err, &refused) || refused.detail != "worktree tree is missing" || refused.next != tc.want {
			t.Fatalf("missing-tree refusal of request %q = %v (next %q), want the missing-tree sentence and %q", tc.request, err, refused.next, tc.want)
		}
	}
}

// TestMissingTreeRouteAgreesWithTheLandedSelector: the missing-tree route names the clean of
// the landed set exactly when that clean selects the assignment, and the `list` help row
// names the same route. A live lease keeps a landed assignment out of the landed set, so its
// route is the release, which names the lease.
func TestMissingTreeRouteAgreesWithTheLandedSelector(t *testing.T) {
	t.Parallel()
	for _, leased := range []bool{true, false} {
		request := map[bool]string{true: "missing-landed-leased", false: "missing-landed-free"}[leased]
		t.Run(request, func(t *testing.T) {
			t.Parallel()
			f := newOwnedAssignment(t, request)
			landAssignment(t, f.root, f.creation, "landed.txt")
			if leased {
				holdLiveLease(t, f)
			}
			mustNoError(t, os.RemoveAll(f.creation.Path))
			r := runVerb(t, verbReset, f.call(resetToStart(f)...))
			next, printed := recordField(r.stdout, "refused{detail="+refusalroute.Sentence(faceResetTreeMissing), refusalroute.NextField)
			plan := runVerb(t, verbClean, f.call("--landed"))
			selected := strings.Contains(plan.stdout, f.creation.Assignment.ID)
			if !printed || (next == "bench worktree clean --landed") != selected || selected == leased {
				t.Fatalf("missing-tree route = %q (printed %t) and clean --landed selects the assignment = %t, want the route to agree with the selector", next, printed, selected)
			}
			if list := runVerb(t, verbList, f.call()); !strings.Contains(list.stdout, "\n  "+next+",") {
				t.Fatalf("list = (%d, %q), want the help row %q that the missing-tree refusal prints", list.exit, list.stdout, next)
			}
		})
	}
}

// TestLandingRetiresAnAbsentLandedSibling: a landing that carries a sibling's work retires
// that sibling when its tree is gone, as it retires a present one.
func TestLandingRetiresAnAbsentLandedSibling(t *testing.T) {
	t.Parallel()
	request := "land-cleanup-absent-sibling"
	f := publicLandingFixture(t, request, "", "")
	folded := foldLandingSibling(t, f.root, f.home, request+"-sibling", f.creation)
	mustNoError(t, os.RemoveAll(folded.sibling.Path))
	j, _ := refreshJoins(nil)
	r := runVerb(t, verbLand, f.callWith(j, landArgs(request, f.base, folded.tip, f.creation.Path)...))
	if r.exit != 0 || !strings.Contains(r.stdout, wantEffects("not-applicable", "complete")) {
		t.Fatalf("absent-sibling landing = (%d, %q, %q), want a complete cleanup at exit 0", r.exit, r.stdout, r.stderr)
	}
	if assignmentActive(t, f.root, folded.sibling.Assignment.ID) {
		t.Fatalf("absent sibling assignment %q is still active after the landing", folded.sibling.Assignment.ID)
	}
}
