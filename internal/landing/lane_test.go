package landing

import (
	"context"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gate/authorization"
)

// SR15: the mapping decides which tree the lane grades and which checks run on it. Every
// field the authority carries is pinned here, so a dropped field is a red rather than a
// silently wider or narrower grade.
func TestLaneAuthorityCarriesTheLaneAndTheBase(t *testing.T) {
	lane := &gate.Lane{
		Checks:    []gate.Phase{{Name: "build", Argv: []string{"go", "build", "./..."}}},
		Kit:       "/kit/root",
		Selective: true,
	}
	want := authorization.LaneAuthority{
		Checks:    lane.Checks,
		Kit:       "/kit/root",
		Selective: true,
		Base:      "abc123",
	}
	if got := laneAuthority(lane, "abc123"); !reflect.DeepEqual(got, want) {
		t.Fatalf("laneAuthority = %+v, want %+v", got, want)
	}
}

// A root that declares no lane answers the whole-project gate owner, not a lane owner
// over a nil lane. The gate owner is identified by the authority it grades under.
func TestNoLaneAnswersTheGateOwner(t *testing.T) {
	gateOwner := New()
	got := NewForLane(nil, "abc123")
	if reflect.ValueOf(got.authorize).Pointer() != reflect.ValueOf(gateOwner.authorize).Pointer() {
		t.Fatal("a nil lane answers an owner other than the gate owner")
	}
	lane := &gate.Lane{Checks: []gate.Phase{{Name: "build"}}}
	if reflect.ValueOf(NewForLane(lane, "abc123").authorize).Pointer() == reflect.ValueOf(gateOwner.authorize).Pointer() {
		t.Fatal("a declared lane answers the gate owner")
	}
}

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
