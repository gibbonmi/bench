package landing

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate/authorization"
)

// OG15: the reviewed landing publishes onto main, so it accepts a graded green alone. A
// lane pass grades a declared check list in a worktree; accepting it here would put an
// ungraded tree on main.
func TestReviewedLandingRefusesALanePass(t *testing.T) {
	o := New()
	o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
		return authorization.Result{Kind: authorization.LanePass}
	}
	if o.reviewedPublishes.permits(authorization.LanePass) {
		t.Fatal("the reviewed landing accepts a lane pass")
	}
	if !o.reviewedPublishes.permits(authorization.Green) {
		t.Fatal("the reviewed landing refuses a graded green")
	}
	// The worktree commit's landing is the one that accepts a lane pass.
	if !o.publishes.permits(authorization.LanePass) || !o.publishes.permits(authorization.Green) {
		t.Fatal("the commit landing refuses a kind it must publish on")
	}
}

// The refusal reaches a caller as the reviewed landing's own error, naming the kind.
func TestReviewedLandingRefusalNamesTheLanePassKind(t *testing.T) {
	o := New()
	o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
		return authorization.Result{Kind: authorization.LanePass}
	}
	f := newReviewedLanding(t, reviewedBaseSpec, "", reviewedAmendedSpec, 0)
	_, err := o.LandReviewed(context.Background(), f.request(t, reviewedAmendedSpec, "must refuse"))
	if err == nil || !strings.Contains(err.Error(), "prospective authorization refused: lane pass") {
		t.Fatalf("error = %v, want the refusal naming the lane pass kind", err)
	}
}
