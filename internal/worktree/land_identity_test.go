// Assignment identity tests for the landing command: each identity dimension —
// request, review base, source tip, and source fingerprint — invalidates the landing
// before composition when it changes.
package worktree

import (
	"bytes"
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
	root, creation, base, tip, tally, home := publicLandingFixture(t, request, "", "")
	j, composed := forbidLandingComposition()

	var stdout, stderr bytes.Buffer
	code := landWith(j, root, home, "", landArgs("land-identity-request-changed", base, tip, creation.Path), &stdout, &stderr)
	if code != 1 || !strings.HasPrefix(stdout.String(), "refused{") {
		t.Fatalf("changed request = (%d, %q, %q), want a refusal", code, stdout.String(), stderr.String())
	}
	requireIdentityRefusalState(t, root, creation.Path, tally, *composed)
}

// TestLandCommandInvalidatesAChangedReviewBaseBeforeComposition is SOL06.
func TestLandCommandInvalidatesAChangedReviewBaseBeforeComposition(t *testing.T) {
	t.Parallel()
	request := "land-identity-base"
	root, creation, _, tip, tally, home := publicLandingFixture(t, request, "", "")
	j, composed := forbidLandingComposition()

	var stdout, stderr bytes.Buffer
	code := landWith(j, root, home, "", landArgs(request, tip, tip, creation.Path), &stdout, &stderr)
	if code != 1 || !strings.HasPrefix(stdout.String(), "refused{") {
		t.Fatalf("changed review base = (%d, %q, %q), want a refusal", code, stdout.String(), stderr.String())
	}
	requireIdentityRefusalState(t, root, creation.Path, tally, *composed)
}

// TestLandCommandInvalidatesAChangedSourceTipBeforeComposition is SOL07.
func TestLandCommandInvalidatesAChangedSourceTipBeforeComposition(t *testing.T) {
	t.Parallel()
	request := "land-identity-tip"
	root, creation, base, tip, tally, home := publicLandingFixture(t, request, "", "")
	commitInWorktree(t, creation.Path, "moved.txt", "moved\n", "tip moved after review")
	j, composed := forbidLandingComposition()

	var stdout, stderr bytes.Buffer
	code := landWith(j, root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 1 || !strings.Contains(stdout.String(), "source tip mismatch") {
		t.Fatalf("changed source tip = (%d, %q, %q), want a tip-mismatch refusal", code, stdout.String(), stderr.String())
	}
	requireIdentityRefusalState(t, root, creation.Path, tally, *composed)
}

// TestLandCommandInvalidatesAChangedSourceFingerprintBeforeTheGate is SOL08. A
// commit-only comparison accepts dirty source content; the fingerprint proof refuses
// it before composition, so before the gate.
func TestLandCommandInvalidatesAChangedSourceFingerprintBeforeTheGate(t *testing.T) {
	t.Parallel()
	request := "land-identity-fingerprint"
	root, creation, base, tip, tally, home := publicLandingFixture(t, request, "", "")
	mustWrite(t, filepath.Join(creation.Path, "dirty.txt"), []byte("uncommitted\n"), 0o600)
	j, composed := forbidLandingComposition()

	var stdout, stderr bytes.Buffer
	code := landWith(j, root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	// The sentence and the repair read from the registry, so each keeps one source. The
	// hostile-source surface stays bounded: the refusal carries a route and no path table,
	// so no source-authored path name reaches the operator's terminal.
	face := landingRefusalFaceByName(faceSourceNotClean)
	if code != 1 || !strings.Contains(stdout.String(), face.detail) || !strings.Contains(stdout.String(), "next="+face.route("")) ||
		strings.Contains(stdout.String(), "refusal_paths") || strings.Contains(stdout.String(), "dirty.txt") {
		t.Fatalf("changed source fingerprint = (%d, %q, %q), want a routed not-clean refusal with no path table", code, stdout.String(), stderr.String())
	}
	requireIdentityRefusalState(t, root, creation.Path, tally, *composed)
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
	root, creation, base, tip, tally, home := specLessLandingFixture(t, request)
	earlier := gitOutput(t, root, "rev-parse", base+"~1")
	if earlier == base {
		t.Fatalf("fixture has no earlier ancestor than the recorded start %q", base)
	}
	if !git.OK("-C", root, "merge-base", "--is-ancestor", earlier, base) {
		t.Fatalf("premise failed: %q is not an ancestor of the destination", earlier)
	}
	j, composed := forbidLandingComposition()

	var stdout, stderr bytes.Buffer
	code := landWith(j, root, home, "", specLessLandArgs(request, earlier, tip, creation.Path), &stdout, &stderr)
	want := "detail=" + reviewedRangeDetail + ",observed=" + earlier + ",wanted=" + base
	if code != 1 || !strings.Contains(stdout.String(), want) {
		t.Fatalf("earlier ancestor base = (%d, %q, %q), want a refusal carrying %q", code, stdout.String(), stderr.String(), want)
	}
	requireIdentityRefusalState(t, root, creation.Path, tally, *composed)
}
