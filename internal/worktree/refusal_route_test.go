// Refusal route tests for the landing: the authority each landing face prints, read from
// the refused record that the face's producing fixture makes the landing print.
package worktree

import (
	"context"
	"errors"
	"maps"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/refusalroute"
)

// reviewerRoute is the marker a reviewer route opens with. The spec pins it, so the proofs
// spell it rather than read it from the renderer they grade.
const reviewerRoute = "reviewer: "

// landingRoute renders a land face's route over the caller's re-run and the facts a
// fixture observed, so a proof that pins a whole route reads the shared registry.
func landingRoute(face, rerun string, values map[string]string) string {
	return landingRouteAfter(face, "", rerun, values)
}

// landingRouteAfter is landingRoute with the preface the landing states ahead of the
// route's first step, such as the skipped-proof sentence.
func landingRouteAfter(face, preface, rerun string, values map[string]string) string {
	facts := map[string]string{refusalroute.FactRerun: rerun}
	maps.Copy(facts, values)
	return refusalroute.New(face, refusalroute.Facts{Preface: preface, Values: facts}).Route
}

// landArgsRerun is the caller's own re-run of a landing that ran with landArgs, spelled
// from the same inputs.
func landArgsRerun(request, base, tip, path string) string {
	return landArgsRerunAt(request, base, "'"+tip+"'", path)
}

// landArgsRerunAt is landArgsRerun with its --source-tip argument as the route prints it.
func landArgsRerunAt(request, base, tipArg, path string) string {
	return "bench worktree land --request '" + request + "' --base '" + base + "' --source-tip " + tipArg +
		" --spec 'x' -m <message> '" + path + "'"
}

// repairedTipArg is the --source-tip argument of a re-run whose repair commits in the
// source. The spec pins the slot, so the proofs spell it.
const repairedTipArg = "<repaired-source-tip>"

// rerunMark stands in for the re-run where a proof reads only the repair ahead of it.
const rerunMark = "RERUN"

// landingRepair is a land face's route up to the caller's re-run.
func landingRepair(face string, values map[string]string) string {
	return strings.TrimSuffix(landingRoute(face, rerunMark, values), rerunMark)
}

// labelOf is the label fact of the assignment that owns the refusing tree.
func labelOf(creation Creation) map[string]string {
	return map[string]string{refusalroute.FactLabel: creation.Assignment.Label}
}

// producedFace is what a landing printed for one face's producing fixture: the fixture,
// the request it landed under, the run, and the next= value of the face's record.
type producedFace struct {
	f       landingFixture
	request string
	r       verbResult
	next    string
}

// produceLandingFace drives one producing fixture through the landing at the stage it
// declares and returns the face's next= value. It fails t unless the face printed a
// non-empty route.
func produceLandingFace(t *testing.T, fixture landingRefusalFixture) producedFace {
	t.Helper()
	request := "landing-face-" + fixture.face
	f := publicLandingFixture(t, request, "", "")
	if fixture.stage == stageIncomplete {
		r := interruptedLanding(t, f, request, f.tip)
		next, printed := landedNext(r.stdout)
		if !printed || next == "" {
			t.Fatalf("face %s = (%d, %q, %q), want a landed record with a non-empty next= field", fixture.face, r.exit, r.stdout, r.stderr)
		}
		return producedFace{f: f, request: request, r: r, next: next}
	}
	var r verbResult
	if fixture.stage == stageResume {
		r = landingFaceResume(t, fixture, f)
	} else {
		fixture.mutate(t, f.root, f.creation)
		// A mutation may add a source commit, so the pinned tip is read after it. An
		// unmoved source reads back the same commit the fixture created.
		tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
		if fixture.tip != nil {
			tip = fixture.tip(t, f.creation)
		}
		r = runVerb(t, verbLand, f.call(landArgs(request, f.base, tip, f.creation.Path)...))
	}
	next, printed := landingFaceNext(r.stdout, refusalroute.Sentence(fixture.face))
	if r.exit != 1 || !printed || next == "" {
		t.Fatalf("face %s = (%d, %q, %q), want exit 1 and a non-empty next= field", fixture.face, r.exit, r.stdout, r.stderr)
	}
	return producedFace{f: f, request: request, r: r, next: next}
}

// TestConflictRepairIsAReviewerRoute is RR16, RR17, and RR27. The hand merge of a
// composition conflict is a raw merge, which the destructive-git guard reserves to the
// reviewer, so each conflict route of the landing and of the merge must open with the
// reviewer marker. The pending arm's fixture holds a merge in progress, so its route names
// the continuation and not a second merge.
func TestConflictRepairIsAReviewerRoute(t *testing.T) {
	t.Parallel()
	landed := func(face string) func(t *testing.T) verbResult {
		return func(t *testing.T) verbResult { return produceLandingFace(t, landingFixtureFor(t, face)).r }
	}
	for _, tc := range []struct {
		face, wantStep string
		produce        func(t *testing.T) verbResult
	}{
		{face: faceCompositionConflict, wantStep: " merge '", produce: landed(faceCompositionConflict)},
		{face: faceCompositionConflictPending, wantStep: " merge --continue", produce: landed(faceCompositionConflictPending)},
		{face: faceMergeConflict, wantStep: " merge '", produce: func(t *testing.T) verbResult {
			return produceMergeFace(t, mergeFixtureFor(t, faceMergeConflict)).r
		}},
	} {
		t.Run(tc.face, func(t *testing.T) {
			t.Parallel()
			r := tc.produce(t)
			next, _ := landingFaceNext(r.stdout, refusalroute.Sentence(tc.face))
			if !strings.HasPrefix(next, reviewerRoute) || !strings.Contains(next, tc.wantStep) {
				t.Fatalf("%s next = %q in %q, want a value that opens with %q and names %q", tc.face, next, r.stdout, reviewerRoute, tc.wantStep)
			}
		})
	}
}

// TestLandingRedRouteNamesTheRepair is RR67. The authorization sentence names no action,
// so the land-red face carries the repair: the route commits in the source through
// `bench commit --in` and ends with the caller's re-run at the repaired source tip.
func TestLandingRedRouteNamesTheRepair(t *testing.T) {
	t.Parallel()
	produced := produceLandingFace(t, landingFixtureFor(t, faceLandRed))
	rerun := landArgsRerunAt(produced.request, produced.f.base, repairedTipArg, produced.f.creation.Path)
	if !strings.Contains(produced.r.stdout, "refused{detail=prospective authorization refused: ") ||
		!strings.Contains(produced.next, "bench commit --in ") || !strings.HasSuffix(produced.next, "; then "+rerun) {
		t.Fatalf("land-red next = %q in %q, want a commit through bench commit --in and the re-run %q", produced.next, produced.r.stdout, rerun)
	}
}

// A landing composition error that no face claims hands back to the reviewer under its
// own sentence, so no landing refusal prints without a route.
func TestUnclaimedLandingCompositionErrorHandsBack(t *testing.T) {
	t.Parallel()
	request := "landing-unclaimed-composition"
	f := publicLandingFixture(t, request, "", "")
	j := defaultJoins()
	j.landReviewed = func(context.Context, landing.ReviewedRequest, landing.Admission) (landing.ReviewedResult, error) {
		return landing.ReviewedResult{}, errors.New("an unclaimed composition fault")
	}
	r := runVerb(t, verbLand, f.callWith(j, landArgs(request, f.base, f.tip, f.creation.Path)...))
	next, printed := landingFaceNext(r.stdout, "an unclaimed composition fault")
	if r.exit != 1 || !printed || !strings.HasPrefix(next, reviewerRoute) {
		t.Fatalf("unclaimed composition error = (%d, %q, %q), want exit 1 and a next= value that opens with %q", r.exit, r.stdout, r.stderr, reviewerRoute)
	}
}

// TestReviewerLandFacesOpenWithTheMarker is RR18. The spec's face inventory gives each of
// these land faces to the reviewer, so the route each producing fixture prints must open
// with the reviewer marker. The list is the spec's, not the registry's, so a face the
// registry hands to the agent turns this test red.
func TestReviewerLandFacesOpenWithTheMarker(t *testing.T) {
	t.Parallel()
	for _, face := range []string{
		faceDestinationNotClean,
		faceDestinationCollision,
		faceCompositionConflict,
		faceCompositionConflictPending,
		faceResumeDestinationResidue,
		faceResumeMarker,
		faceLandHandback,
	} {
		t.Run(face, func(t *testing.T) {
			t.Parallel()
			produced := produceLandingFace(t, landingFixtureFor(t, face))
			if !strings.HasPrefix(produced.next, reviewerRoute) {
				t.Fatalf("%s next = %q in %q, want a value that opens with %q", face, produced.next, produced.r.stdout, reviewerRoute)
			}
		})
	}
}

// landingFaceNext reads the next= value out of the refused record whose detail names the
// face. next= is the last field of the refused record. The second result reports whether
// the face printed at all.
func landingFaceNext(stdout, detail string) (string, bool) {
	return recordField(stdout, "refused{detail="+detail, refusalroute.NextField)
}

// landedNext reads the next= value out of the landed record of an incomplete landing.
func landedNext(stdout string) (string, bool) {
	return landedField(stdout, refusalroute.NextField)
}

// landedField reads one field of the landed record. The census field is the record's last.
func landedField(stdout, field string) (string, bool) {
	return recordField(stdout, "landed{", field, "census")
}

// recordField is the one reader of the printed record layout. It reads field out of the
// first line that opens with opening. A value can hold a comma, so the value runs to the
// first of the later fields that the record prints after it, or to the closing brace. The
// second result reports whether the record printed at all.
func recordField(stdout, opening, field string, later ...string) (string, bool) {
	for _, line := range strings.Split(stdout, "\n") {
		if !strings.HasPrefix(line, opening) {
			continue
		}
		for _, separator := range []string{"{", ","} {
			if _, value, found := strings.Cut(line, separator+field+"="); found {
				for _, label := range later {
					value, _, _ = strings.Cut(value, ","+label+"=")
				}
				return strings.TrimSuffix(value, "}"), true
			}
		}
		return "", true
	}
	return "", false
}
