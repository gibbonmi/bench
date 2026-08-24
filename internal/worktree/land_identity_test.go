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

	"github.com/gibbonmi/bench/internal/landing"
)

// forbidLandingComposition fails the test if the landing reaches composition, and
// returns a counter the assertions read. The identity seams must refuse first.
func forbidLandingComposition(t *testing.T) *int {
	t.Helper()
	composed := 0
	old := landReviewed
	landReviewed = func(context.Context, landing.ReviewedRequest) (landing.ReviewedResult, error) {
		composed++
		return landing.ReviewedResult{}, errors.New("composition must not start")
	}
	t.Cleanup(func() { landReviewed = old })
	return &composed
}

// TestLandCommandInvalidatesAChangedRequestBeforeComposition is SOL05.
func TestLandCommandInvalidatesAChangedRequestBeforeComposition(t *testing.T) {
	request := "land-identity-request"
	root, creation, base, tip, tally := publicLandingFixture(t, request, "", "")
	composed := forbidLandingComposition(t)

	var stdout, stderr bytes.Buffer
	code := LandCommand(root, "", landArgs("land-identity-request-changed", base, tip, creation.Path), &stdout, &stderr)
	if code != 1 || !strings.HasPrefix(stdout.String(), "refused{") {
		t.Fatalf("changed request = (%d, %q, %q), want a refusal", code, stdout.String(), stderr.String())
	}
	requireIdentityRefusalState(t, root, creation.Path, tally, *composed)
}

// TestLandCommandInvalidatesAChangedReviewBaseBeforeComposition is SOL06.
func TestLandCommandInvalidatesAChangedReviewBaseBeforeComposition(t *testing.T) {
	request := "land-identity-base"
	root, creation, _, tip, tally := publicLandingFixture(t, request, "", "")
	composed := forbidLandingComposition(t)

	var stdout, stderr bytes.Buffer
	code := LandCommand(root, "", landArgs(request, tip, tip, creation.Path), &stdout, &stderr)
	if code != 1 || !strings.HasPrefix(stdout.String(), "refused{") {
		t.Fatalf("changed review base = (%d, %q, %q), want a refusal", code, stdout.String(), stderr.String())
	}
	requireIdentityRefusalState(t, root, creation.Path, tally, *composed)
}

// TestLandCommandInvalidatesAChangedSourceTipBeforeComposition is SOL07.
func TestLandCommandInvalidatesAChangedSourceTipBeforeComposition(t *testing.T) {
	request := "land-identity-tip"
	root, creation, base, tip, tally := publicLandingFixture(t, request, "", "")
	commitInWorktree(t, creation.Path, "moved.txt", "moved\n", "tip moved after review")
	composed := forbidLandingComposition(t)

	var stdout, stderr bytes.Buffer
	code := LandCommand(root, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 1 || !strings.Contains(stdout.String(), "source tip mismatch") {
		t.Fatalf("changed source tip = (%d, %q, %q), want a tip-mismatch refusal", code, stdout.String(), stderr.String())
	}
	requireIdentityRefusalState(t, root, creation.Path, tally, *composed)
}

// TestLandCommandInvalidatesAChangedSourceFingerprintBeforeTheGate is SOL08. A
// commit-only comparison accepts dirty source content; the fingerprint proof refuses
// it before composition, so before the gate.
func TestLandCommandInvalidatesAChangedSourceFingerprintBeforeTheGate(t *testing.T) {
	request := "land-identity-fingerprint"
	root, creation, base, tip, tally := publicLandingFixture(t, request, "", "")
	mustWrite(t, filepath.Join(creation.Path, "dirty.txt"), []byte("uncommitted\n"), 0o600)
	composed := forbidLandingComposition(t)

	var stdout, stderr bytes.Buffer
	code := LandCommand(root, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 1 || !strings.Contains(stdout.String(), "reviewed source is not clean") {
		t.Fatalf("changed source fingerprint = (%d, %q, %q), want a not-clean refusal", code, stdout.String(), stderr.String())
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
