// Moved-destination tests: a landing over a moved destination names the fold and the re-run.
// A cause that no fold repairs hands back.
package worktree

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/refusalroute"
)

// foldOfMain opens the printed fold of main into the source.
const foldOfMain = "bench worktree merge --from 'main' "

// movedFacts are the facts the moved-destination route reads: the source's label, and the
// destination branch that the fold merges from.
func movedFacts(f landingFixture) map[string]string {
	return map[string]string{refusalroute.FactLabel: f.creation.Assignment.Label, refusalroute.FactFrom: "main"}
}

// requireMovedRoute requires the moved-destination face over cause, and its route with the
// re-run at the repaired source tip.
func requireMovedRoute(t *testing.T, r verbResult, f landingFixture, request, cause string) {
	t.Helper()
	next, printed := landingFaceNext(r.stdout, refusalroute.Sentence(faceLandDestinationMoved)+causeJoiner+cause)
	want := landingRoute(faceLandDestinationMoved, landArgsRerunAt(request, f.base, repairedTipArg, f.creation.Path), movedFacts(f))
	if r.exit != 1 || !printed || next != want {
		t.Fatalf("moved destination = (%d, %q, %q), want the %s face over %q with next=%q", r.exit, r.stdout, r.stderr, faceLandDestinationMoved, cause, want)
	}
}

// TestCompletionOverAMovedDestinationNamesTheFold is BW17, BW18, and BW19. The destination
// gains one commit, so the completion proof refuses. The route folds main and ends with the
// re-run, and the output names no bench doctor.
func TestCompletionOverAMovedDestinationNamesTheFold(t *testing.T) {
	t.Parallel()
	produced := produceLandingFace(t, landingFixtureFor(t, faceLandDestinationMoved))
	requireMovedRoute(t, produced.r, produced.f, produced.request, "prospective authorization refused: ")
	rerun := landArgsRerunAt(produced.request, produced.f.base, repairedTipArg, produced.f.creation.Path)
	if !strings.HasPrefix(produced.next, foldOfMain) || !strings.HasSuffix(produced.next, refusalroute.StepJoiner+rerun) {
		t.Fatalf("next = %q, want the fold of main first and the re-run %q last", produced.next, rerun)
	}
	if output := produced.r.stdout + produced.r.stderr; strings.Contains(output, "bench doctor") {
		t.Fatalf("output = %q, want no bench doctor", output)
	}
}

// TestCompareAndSwapRefusalNamesTheFold is BW21. The destination moves inside the
// publication gap, so the compare-and-swap refuses the ref update.
func TestCompareAndSwapRefusalNamesTheFold(t *testing.T) {
	t.Parallel()
	request := "land-moved-compare-and-swap"
	f := publicLandingFixture(t, request, "", "")
	j := defaultJoins()
	j.publicationGap = func(root string) {
		commitInWorktree(t, root, "moved.txt", "destination bytes\n", "destination moves in the publication gap")
	}
	r := runVerb(t, verbLand, f.callWith(j, landArgs(request, f.base, f.tip, f.creation.Path)...))
	r.mustViaJoins(t)
	requireMovedRoute(t, r, f, request, "destination compare-and-swap refused")
}

// TestCheckoutChangeWithoutARefMoveKeepsTheHandback is BW22. An ignored file appears in the
// destination checkout during the gate, and the destination ref stays, so no fold helps.
func TestCheckoutChangeWithoutARefMoveKeepsTheHandback(t *testing.T) {
	t.Parallel()
	request := "land-moved-checkout-change"
	f := barrierLandingFixture(t, request)
	exclude := mustAdminPath(t, f.root, "info/exclude")
	mustMkdirAll(t, filepath.Dir(exclude), 0o755)
	mustWrite(t, exclude, []byte("late-ignored\n"), 0o644)
	var writeErr error
	done := f.during(func() { writeErr = os.WriteFile(filepath.Join(f.root, "late-ignored"), []byte("late\n"), 0o644) })
	r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
	mustNoError(t, <-done)
	mustNoError(t, writeErr)
	requireHandback(t, r, "landing destination checkout changed")
	if main := gitOutput(t, f.root, "rev-parse", "main"); main != f.base {
		t.Fatalf("main = %s, want the unmoved %s", main, f.base)
	}
}

// TestProofFaultOverACurrentSourceHandsBack is BW23. The reviewed source holds the
// destination base, so a completion proof fault is no moved destination.
func TestProofFaultOverACurrentSourceHandsBack(t *testing.T) {
	t.Parallel()
	request := "land-moved-current-source"
	f := publicLandingFixture(t, request, "", "")
	j := defaultJoins()
	j.landReviewed = func(context.Context, landing.ReviewedRequest, landing.Admission) (landing.ReviewedResult, error) {
		return landing.ReviewedResult{}, landing.AuthorizationRefusal{Result: authorization.Result{Kind: authorization.Infrastructure, Reason: "completion proof fault", Cause: authorization.CompletionProof}}
	}
	r := runVerb(t, verbLand, f.callWith(j, landArgs(request, f.base, f.tip, f.creation.Path)...))
	r.mustViaJoins(t)
	requireHandback(t, r, "prospective authorization refused: infrastructure (completion proof fault)")
}
