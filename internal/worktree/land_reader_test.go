// The reader landing: a publication on a green reader set keeps the green marker.
package worktree

import (
	"context"
	"testing"

	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/landing"
)

// TestReaderLandingLeavesTheGreenMarker is NG58. The landing hands the destination's green
// marker to the reviewed landing, and a publication on `reader green` leaves the marker at
// that prior commit, because the reader set graded the readers alone.
func TestReaderLandingLeavesTheGreenMarker(t *testing.T) {
	t.Parallel()
	request := "reader-landing-marker"
	f := specLessLandingFixture(t, request)
	gitRun(t, f.root, "update-ref", markerRef, f.base)
	var handed string
	j := defaultJoins()
	j.landReviewed = func(ctx context.Context, request landing.ReviewedRequest, admission landing.Admission) (landing.ReviewedResult, error) {
		handed = request.GreenMarker
		result, err := defaultJoins().landReviewed(ctx, request, admission)
		result.Kind = authorization.ReaderGreen
		return result, err
	}
	r := runVerb(t, verbLand, f.callWith(j, specLessLandArgs(request, f.base, f.tip, f.creation.Path)...))
	r.mustViaJoins(t)
	published := gitOutput(t, f.root, "rev-parse", "main")
	if r.exit != 0 || published == f.base {
		t.Fatalf("reader landing = (%d, %q, %q), want a publication", r.exit, r.stdout, r.stderr)
	}
	if handed != f.base {
		t.Fatalf("reviewed request green marker = %q, want the prior marker %s", handed, f.base)
	}
	if got := gitOutput(t, f.root, "rev-parse", markerRef); got != f.base {
		t.Fatalf("%s = %s after a reader green publication, want the prior marker %s", markerRef, got, f.base)
	}
}
