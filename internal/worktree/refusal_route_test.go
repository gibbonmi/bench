// Refusal route tests for the landing: the authority each landing face prints, read from
// the refused record that the face's producing fixture makes the landing print.
package worktree

import (
	"io"
	"maps"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/refusalroute"
)

// reviewerRoute is the marker a reviewer route opens with. The spec pins it, so the proofs
// spell it rather than read it from the renderer they grade.
const reviewerRoute = "reviewer: "

// landingRoute renders a land face's route over the caller's re-run and the facts a
// fixture observed, so a proof that pins a whole route reads the shared registry.
func landingRoute(face, rerun string, values map[string]string) string {
	facts := map[string]string{refusalroute.FactRerun: rerun}
	maps.Copy(facts, values)
	return refusalroute.New(face, refusalroute.Facts{Values: facts}).Route
}

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

// landedFaceNext runs a first landing over the face's producing fixture and returns the
// next= value of the refused record whose sentence opens with detail.
func landedFaceNext(t *testing.T, face, detail string) (landingFixture, verbResult, string) {
	t.Helper()
	request := "face-route-" + face
	f := publicLandingFixture(t, request, "", "")
	landingFixtureFor(t, face).mutate(t, f.root, f.creation)
	r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
	next, printed := landingFaceNext(r.stdout, detail)
	if r.exit != 1 || !printed {
		t.Fatalf("%s landing = (%d, %q, %q), want exit 1 and the face's refused record", face, r.exit, r.stdout, r.stderr)
	}
	return f, r, next
}

// TestConflictRepairIsAReviewerRoute is RR16 and RR17. The hand merge of a composition
// conflict is a raw merge, which the destructive-git guard reserves to the reviewer, so
// both conflict routes must open with the reviewer marker. The pending arm's fixture holds
// a merge in progress, so its route names the continuation and not a second merge.
func TestConflictRepairIsAReviewerRoute(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		face, wantStep string
	}{
		{face: faceCompositionConflict, wantStep: " merge '"},
		{face: faceCompositionConflictPending, wantStep: " merge --continue"},
	} {
		t.Run(tc.face, func(t *testing.T) {
			t.Parallel()
			_, r, next := landedFaceNext(t, tc.face, "composition conflict")
			if !strings.HasPrefix(next, reviewerRoute) || !strings.Contains(next, tc.wantStep) {
				t.Fatalf("%s next = %q in %q, want a value that opens with %q and names %q", tc.face, next, r.stdout, reviewerRoute, tc.wantStep)
			}
		})
	}
}

// TestLandingCleanlinessFacesRouteByAuthority is RR18 and RR19. A dirty destination is the
// reviewer's own uncommitted work in the primary checkout, so its route hands back. A dirty
// source is the agent's own work, so its route commits it through Bench at the source's
// label rather than through a raw commit the guard denies.
func TestLandingCleanlinessFacesRouteByAuthority(t *testing.T) {
	t.Parallel()
	t.Run(faceDestinationNotClean, func(t *testing.T) {
		t.Parallel()
		_, r, next := landedFaceNext(t, faceDestinationNotClean, refusalroute.Sentence(faceDestinationNotClean))
		if !strings.HasPrefix(next, reviewerRoute) {
			t.Fatalf("dirty destination next = %q in %q, want a value that opens with %q", next, r.stdout, reviewerRoute)
		}
	})
	t.Run(faceSourceNotClean, func(t *testing.T) {
		t.Parallel()
		f, r, next := landedFaceNext(t, faceSourceNotClean, refusalroute.Sentence(faceSourceNotClean))
		label := f.creation.Assignment.Label
		if strings.HasPrefix(next, reviewerRoute) || !strings.Contains(next, "bench commit --in ") || !strings.Contains(next, label) {
			t.Fatalf("dirty source next = %q in %q, want an agent route that names bench commit --in and the label %q", next, r.stdout, label)
		}
	})
}

// landingStage names when in a landing a fixture's face prints, which decides how the walk
// drives the fixture and the command the face's route ends with.
type landingStage int

const (
	// stagePreflight is the first run's preflight. Its route ends with the caller's own
	// re-run of the landing.
	stagePreflight landingStage = iota
	// stageComposition is the first run past its preflight. Its route ends with a re-run
	// the operator points at a repaired source tip.
	stageComposition
	// stageResume is the resume path. Its route ends with the caller's own resume.
	stageResume
	// stageIncomplete is a published landing whose later step did not finish. Its landed
	// record ends with the caller's resume.
	stageIncomplete
)

// landingRefusalFixture produces exactly one landing face. mutate breaks the landing
// fixture so that face is the one the landing prints. The registry walk requires one
// fixture per land face of the shared registry, so a face added with no fixture turns
// TestLandingRefusalRegistryHasAProducingFixture red. stage states how the walk drives the fixture: a resume fixture's mutation
// runs against the published destination of an interrupted landing, not against the
// first run.
type landingRefusalFixture struct {
	face   string
	stage  landingStage
	mutate func(t *testing.T, root string, creation Creation)
	// tip states the source tip the run names, read after the mutation. A face that
	// refuses on the caller's own tip needs a value other than the worktree's head, which
	// the walk names by default.
	tip func(t *testing.T, creation Creation) string
}

// landingFixtureFor returns the one fixture that produces the named face, so a proof that
// raises one face composes its mutation rather than copying it.
func landingFixtureFor(t *testing.T, face string) landingRefusalFixture {
	t.Helper()
	for _, fixture := range landingRefusalFixtures() {
		if fixture.face == face {
			return fixture
		}
	}
	t.Fatalf("no producing fixture for landing face %q", face)
	return landingRefusalFixture{}
}

func landingRefusalFixtures() []landingRefusalFixture {
	destinationConflict := func(t *testing.T, root string, _ Creation) {
		commitInWorktree(t, root, "owned.txt", "destination bytes\n", "destination conflict")
	}
	return []landingRefusalFixture{
		{
			face: faceDestinationNotClean,
			mutate: func(t *testing.T, root string, _ Creation) {
				mustWrite(t, filepath.Join(root, "tracked.txt"), []byte("dirty\n"), 0o644)
			},
		},
		{
			// The reviewed source adds owned.txt, so an untracked file at that path stands
			// where the landing writes.
			face: faceDestinationCollision,
			mutate: func(t *testing.T, root string, _ Creation) {
				mustWrite(t, filepath.Join(root, "owned.txt"), []byte("operator bytes\n"), 0o600)
			},
		},
		{
			face: faceSourceNotClean,
			mutate: func(t *testing.T, _ string, creation Creation) {
				mustWrite(t, filepath.Join(creation.Path, "scratch"), []byte("scratch\n"), 0o600)
			},
		},
		{
			face: faceSourceNotFenced,
			mutate: func(t *testing.T, _ string, creation Creation) {
				commitInWorktree(t, creation.Path, "stray.txt", "stray\n", "out of fence")
			},
		},
		{
			// The source moves past the tip the caller names, which is the state an
			// operator reaches when a review repair adds a commit.
			face: faceSourceTipMismatch,
			mutate: func(t *testing.T, _ string, creation Creation) {
				commitInWorktree(t, creation.Path, "owned.txt", "moved\n", "source moved past the named tip")
			},
			tip: func(t *testing.T, creation Creation) string {
				return gitOutput(t, creation.Path, "rev-parse", "HEAD~1")
			},
		},
		{
			// The destination moves onto a commit the reviewed source also changed, which
			// is the state a composition conflict outside the rule table reaches.
			face:   faceCompositionConflict,
			stage:  stageComposition,
			mutate: destinationConflict,
		},
		{
			// The same conflict over a source that already holds a merge in progress.
			face:  faceCompositionConflictPending,
			stage: stageComposition,
			mutate: func(t *testing.T, root string, creation Creation) {
				destinationConflict(t, root, creation)
				admin := gitOutput(t, creation.Path, "rev-parse", "--absolute-git-dir")
				mustWrite(t, filepath.Join(admin, "MERGE_HEAD"), []byte(gitOutput(t, root, "rev-parse", "HEAD")+"\n"), 0o644)
			},
		},
		{
			face:  faceResumeDestinationResidue,
			stage: stageResume,
			mutate: func(t *testing.T, root string, _ Creation) {
				mustWrite(t, filepath.Join(root, "tracked.txt"), []byte("dirty\n"), 0o644)
			},
		},
		{
			face:  faceResumeMarker,
			stage: stageResume,
			mutate: func(t *testing.T, root string, _ Creation) {
				// The marker refuses only once the destination has moved past the
				// published landing, so the mutation moves it and then drops the marker.
				commitInWorktree(t, root, "destination-after-publication", "forward\n", "destination movement")
				gitRun(t, root, "update-ref", "-d", "refs/bench/green/main")
			},
		},
		{
			// The walk interrupts the landing at its release step, so the fixture adds no
			// mutation of its own.
			face:   faceLandIncomplete,
			stage:  stageIncomplete,
			mutate: func(*testing.T, string, Creation) {},
		},
		{
			// A detached landing checkout is the primary checkout's own state, which only
			// the reviewer clears.
			face: faceLandHandback,
			mutate: func(t *testing.T, root string, _ Creation) {
				gitRun(t, root, "checkout", "-q", "--detach")
			},
		},
	}
}

// interruptedLanding runs the first landing with a release step that fails, so the
// landing publishes and exits incomplete.
func interruptedLanding(t *testing.T, f landingFixture, request, tip string) verbResult {
	t.Helper()
	broken := defaultJoins()
	broken.releaseLandingAssignment = func(joins, ambient, string, []string, io.Writer, io.Writer) int { return 1 }
	r := runVerb(t, verbLand, f.callWith(broken, landArgs(request, f.base, tip, f.creation.Path)...))
	if r.exit != 3 {
		t.Fatalf("interrupted landing = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	return r
}

// landingFaceResume drives a resume fixture. It interrupts a landing at the release step,
// applies the fixture's mutation to the published destination, and resumes.
func landingFaceResume(t *testing.T, fixture landingRefusalFixture, f landingFixture) verbResult {
	t.Helper()
	request := "landing-face-" + fixture.face
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	interruptedLanding(t, f, request, tip)
	published := gitOutput(t, f.root, "rev-parse", "main")
	fixture.mutate(t, f.root, f.creation)
	args := []string{"--resume", published, "--request", request, "--base", f.base, "--source-tip", tip, "--spec", "x", f.creation.Path}
	return runVerb(t, verbLand, f.callWith(defaultJoins(), args...))
}

// landingFaceNext reads the next= value out of the refused record whose detail names the
// face. next= is the last field the formatter writes, so the value runs to the closing
// brace. The second result reports whether the face printed at all.
func landingFaceNext(stdout, detail string) (string, bool) {
	for _, line := range strings.Split(stdout, "\n") {
		if !strings.HasPrefix(line, "refused{detail="+detail) {
			continue
		}
		_, next, found := strings.Cut(line, ",next=")
		if !found {
			return "", true
		}
		return strings.TrimSuffix(next, "}"), true
	}
	return "", false
}

// landedNext reads the next= value out of the landed record of an incomplete landing. The
// census field follows next=, so the value runs to that field.
func landedNext(stdout string) (string, bool) {
	_, rest, found := strings.Cut(stdout, "landed{")
	if !found {
		return "", false
	}
	_, next, found := strings.Cut(rest, ",next=")
	if !found {
		return "", true
	}
	next, _, _ = strings.Cut(next, ",census=")
	return next, true
}
