// Assignment identity tests for the landing command: each identity dimension —
// request, review base, source tip, and source fingerprint — invalidates the landing
// before composition when it changes.
package worktree

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/landing"
)

// forbidLandingComposition returns a seam set that fails the test if the landing
// reaches composition, and a counter the assertions read. The identity seams must
// refuse first.
func forbidLandingComposition() (joins, *int) {
	composed := 0
	j := defaultJoins()
	j.landReviewed = func(context.Context, landing.ReviewedRequest) (landing.ReviewedResult, error) {
		composed++
		return landing.ReviewedResult{}, errors.New("composition must not start")
	}
	return j, &composed
}

// TestLandCommandInvalidatesAChangedRequestBeforeComposition is SOL05.
func TestLandCommandInvalidatesAChangedRequestBeforeComposition(t *testing.T) {
	t.Parallel()
	request := "land-identity-request"
	f := publicLandingFixture(t, request, "", "")
	j, composed := forbidLandingComposition()

	r := runVerb(t, verbLand, f.callWith(j, landArgs("land-identity-request-changed", f.base, f.tip, f.creation.Path)...))
	if r.exit != 1 || !strings.HasPrefix(r.stdout, "refused{") {
		t.Fatalf("changed request = (%d, %q, %q), want a refusal", r.exit, r.stdout, r.stderr)
	}
	requireIdentityRefusalState(t, f.root, f.creation.Path, f.tally, *composed)
}

// TestLandCommandInvalidatesAChangedReviewBaseBeforeComposition is SOL06.
func TestLandCommandInvalidatesAChangedReviewBaseBeforeComposition(t *testing.T) {
	t.Parallel()
	request := "land-identity-base"
	f := publicLandingFixture(t, request, "", "")
	j, composed := forbidLandingComposition()

	r := runVerb(t, verbLand, f.callWith(j, landArgs(request, f.tip, f.tip, f.creation.Path)...))
	if r.exit != 1 || !strings.HasPrefix(r.stdout, "refused{") {
		t.Fatalf("changed review base = (%d, %q, %q), want a refusal", r.exit, r.stdout, r.stderr)
	}
	requireIdentityRefusalState(t, f.root, f.creation.Path, f.tally, *composed)
}

// TestLandCommandInvalidatesAChangedSourceTipBeforeComposition is SOL07.
func TestLandCommandInvalidatesAChangedSourceTipBeforeComposition(t *testing.T) {
	t.Parallel()
	request := "land-identity-tip"
	f := publicLandingFixture(t, request, "", "")
	commitInWorktree(t, f.creation.Path, "moved.txt", "moved\n", "tip moved after review")
	j, composed := forbidLandingComposition()

	r := runVerb(t, verbLand, f.callWith(j, landArgs(request, f.base, f.tip, f.creation.Path)...))
	if r.exit != 1 || !strings.Contains(r.stdout, "source tip mismatch") {
		t.Fatalf("changed source tip = (%d, %q, %q), want a tip-mismatch refusal", r.exit, r.stdout, r.stderr)
	}
	requireIdentityRefusalState(t, f.root, f.creation.Path, f.tally, *composed)
}

// TestLandCommandInvalidatesAChangedSourceFingerprintBeforeTheGate is SOL08. A
// commit-only comparison accepts dirty source content; the fingerprint proof refuses
// it before composition, so before the gate.
func TestLandCommandInvalidatesAChangedSourceFingerprintBeforeTheGate(t *testing.T) {
	t.Parallel()
	request := "land-identity-fingerprint"
	f := publicLandingFixture(t, request, "", "")
	mustWrite(t, filepath.Join(f.creation.Path, "dirty.txt"), []byte("uncommitted\n"), 0o600)
	j, composed := forbidLandingComposition()

	r := runVerb(t, verbLand, f.callWith(j, landArgs(request, f.base, f.tip, f.creation.Path)...))
	// The sentence and the repair read from the registry, so each keeps one source. The
	// hostile-source surface stays bounded: the refusal carries a route and no path table,
	// so no source-authored path name reaches the operator's terminal.
	face := landingRefusalFaceByName(faceSourceNotClean)
	if r.exit != 1 || !strings.Contains(r.stdout, face.detail) || !strings.Contains(r.stdout, "next="+face.route("")) ||
		strings.Contains(r.stdout, "refusal_paths") || strings.Contains(r.stdout, "dirty.txt") {
		t.Fatalf("changed source fingerprint = (%d, %q, %q), want a routed not-clean refusal with no path table", r.exit, r.stdout, r.stderr)
	}
	requireIdentityRefusalState(t, f.root, f.creation.Path, f.tally, *composed)
}

// requireIdentityRefusalState pins what every identity refusal leaves behind: no
// composition, no gate run, an unmoved destination, and a retained source worktree.
func requireIdentityRefusalState(t *testing.T, root, path, tally string, composed int) {
	t.Helper()
	if composed != 0 {
		t.Fatalf("identity refusal reached composition %d times, want 0", composed)
	}
	if _, err := os.Stat(tally); !os.IsNotExist(err) {
		t.Fatalf("identity refusal ran the gate: %v", err)
	}
	if got := gitOutput(t, root, "status", "--porcelain=v1"); got != "" {
		t.Fatalf("identity refusal dirtied the destination: %q", got)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("identity refusal removed the source worktree: %v", err)
	}
}

// TestLandCommandRefusesAReviewBaseThatIsNotAnAncestorOfTheDestination covers the
// confusion the detail is for: the source folded the destination advance, and the
// operator names that fold commit as `--base`. The fold is a descendant of the recorded
// start, so the recorded-start binding accepts it, but it lives on the source branch and
// the destination never composed it. The detail names the two bases the flag can mean,
// because the fold commit a review reads and the destination tip the source folded are
// easy to confuse when only one prints as `--base`.
func TestLandCommandRefusesAReviewBaseThatIsNotAnAncestorOfTheDestination(t *testing.T) {
	t.Parallel()
	request := "land-identity-not-ancestor"
	f := foldedLandingFixture(t, request)
	if git.OK("-C", f.root, "merge-base", "--is-ancestor", f.fold, f.base) {
		t.Fatalf("premise failed: the fold commit %q must not already be on the destination %q", f.fold, f.base)
	}
	j, composed := forbidLandingComposition()

	r := runVerb(t, verbLand, f.callWith(j, specLessLandArgs(request, f.fold, f.tip, f.creation.Path)...))
	// The expectation is spelled out here rather than read from landingBaseNotAncestorDetail,
	// so a mutation of that constant turns this test red instead of passing silently.
	const wantDetail = "review base is not an ancestor of the landing destination: --base takes the landing base, the default-branch tip the source folded, not the fold commit the review read"
	want := "detail=" + wantDetail + ",observed=" + f.fold + ",wanted=" + f.base
	if r.exit != 1 || !strings.Contains(r.stdout, want) {
		t.Fatalf("non-ancestor base = (%d, %q, %q), want a refusal carrying %q", r.exit, r.stdout, r.stderr, want)
	}
	requireIdentityRefusalState(t, f.root, f.creation.Path, f.tally, *composed)
}

// TestLandCommandRefusesAReviewBaseBehindTheRecordedStart is the SOL06 mutation guard
// the ancestry check alone cannot answer, and it closes C4. The base names an earlier
// ancestor of the recorded start, so it is a valid ancestor of the destination and the
// ancestry guard accepts it, and the spec-less landing names no ownership fence that
// would refuse it for another reason. Only the recorded-start binding refuses it,
// because that earlier base grades a range wider than the one the assignment
// authorized, and it refuses before composition.
func TestLandCommandRefusesAReviewBaseBehindTheRecordedStart(t *testing.T) {
	t.Parallel()
	request := "land-identity-recorded-start"
	f := specLessLandingFixture(t, request)
	earlier := gitOutput(t, f.root, "rev-parse", f.base+"~1")
	if earlier == f.base {
		t.Fatalf("fixture has no earlier ancestor than the recorded start %q", f.base)
	}
	if !git.OK("-C", f.root, "merge-base", "--is-ancestor", earlier, f.base) {
		t.Fatalf("premise failed: %q is not an ancestor of the destination", earlier)
	}
	j, composed := forbidLandingComposition()

	r := runVerb(t, verbLand, f.callWith(j, specLessLandArgs(request, earlier, f.tip, f.creation.Path)...))
	want := "detail=" + reviewedRangeDetail + ",observed=" + earlier + ",wanted=" + f.base
	if r.exit != 1 || !strings.Contains(r.stdout, want) {
		t.Fatalf("earlier ancestor base = (%d, %q, %q), want a refusal carrying %q", r.exit, r.stdout, r.stderr, want)
	}
	requireIdentityRefusalState(t, f.root, f.creation.Path, f.tally, *composed)
}
