// Refusal route tests for the landing: the authority each landing face prints, read from
// the refused record that the face's producing fixture makes the landing print.
package worktree

import (
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
			produced := produceLandingFace(t, landingFixtureFor(t, tc.face))
			if !strings.HasPrefix(produced.next, reviewerRoute) || !strings.Contains(produced.next, tc.wantStep) {
				t.Fatalf("%s next = %q in %q, want a value that opens with %q and names %q", tc.face, produced.next, produced.r.stdout, reviewerRoute, tc.wantStep)
			}
		})
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
// TestLandingRefusalRegistryHasAProducingFixture red. stage states how the walk drives
// the fixture: a resume fixture's mutation runs against the published destination of an
// interrupted landing, not against the first run.
type landingRefusalFixture struct {
	face   string
	stage  landingStage
	mutate func(t *testing.T, root string, creation Creation)
	// tip states the source tip the run names, read after the mutation. A face that
	// refuses on the caller's own tip needs a value other than the worktree's head, which
	// the walk names by default.
	tip func(t *testing.T, creation Creation) string
	// carry carries out a printed route step by the fixture's own means, keyed by the
	// step's index in the face's route: an instruction, or a step of a reviewer route.
	// The follow walk runs every other step verbatim.
	carry map[int]func(t *testing.T, f landingFixture)
	// names states what the printed route must name, read from the fixture.
	names func(f landingFixture) []string
	// repairsSource states that the face's repair commits in the source, so its re-run
	// names the repaired tip in place of the caller's.
	repairsSource bool
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
	// The reviewer discards the destination's uncommitted edit.
	discardDestination := func(t *testing.T, f landingFixture) { gitRun(t, f.root, "checkout", "--", "tracked.txt") }
	// The reviewer's hand merge of the destination into the source, its resolution, and the
	// commit that records it.
	handMerge := func(t *testing.T, f landingFixture) {
		merge := descendant(t, "git", "-C", f.creation.Path, "merge", "--no-commit", "main")
		if out, err := merge.CombinedOutput(); err == nil || !strings.Contains(string(out), "CONFLICT") {
			t.Fatalf("hand merge = %v, %s; want the conflict", err, out)
		}
		mustWrite(t, filepath.Join(f.creation.Path, "owned.txt"), []byte("destination bytes\nreviewed repair\n"), 0o644)
		gitRun(t, f.creation.Path, "add", "owned.txt")
	}
	commitResolution := func(t *testing.T, f landingFixture) {
		gitRun(t, f.creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "resolve the composition")
	}
	// The review of the repaired source refreshes its completion evidence.
	review := func(t *testing.T, f landingFixture) { refreshLandingEvidence(t, f.creation.Path, f.base) }
	return []landingRefusalFixture{
		{
			face: faceDestinationNotClean,
			mutate: func(t *testing.T, root string, _ Creation) {
				mustWrite(t, filepath.Join(root, "tracked.txt"), []byte("dirty\n"), 0o644)
			},
			carry: map[int]func(*testing.T, landingFixture){0: discardDestination},
		},
		{
			// The reviewed source adds owned.txt, so an untracked file at that path stands
			// where the landing writes.
			face: faceDestinationCollision,
			mutate: func(t *testing.T, root string, _ Creation) {
				mustWrite(t, filepath.Join(root, "owned.txt"), []byte("operator bytes\n"), 0o600)
			},
			carry: map[int]func(*testing.T, landingFixture){0: func(t *testing.T, f landingFixture) {
				mustRemove(t, filepath.Join(f.root, "owned.txt"))
			}},
		},
		{
			// The dirty path is one the ticket writes, so the printed commit lands it. The
			// project declares a commit lane in an ignored phase manifest, which leaves the
			// source's status alone, so the printed commit has a lane to pass.
			face:          faceSourceNotClean,
			repairsSource: true,
			mutate: func(t *testing.T, _ string, creation Creation) {
				mustWrite(t, filepath.Join(creation.Path, "owned.txt"), []byte("uncommitted repair\n"), 0o644)
				mustWrite(t, filepath.Join(creation.Path, ".bench", "phases.json"), []byte(`{"phases":[{"name":"build","argv":["true"]}],"lane":[{"name":"unit","argv":["true"]}]}`), 0o644)
				exclude := filepath.Join(gitOutput(t, creation.Path, "rev-parse", "--path-format=absolute", "--git-common-dir"), "info", "exclude")
				mustMkdirAll(t, filepath.Dir(exclude), 0o755)
				mustWrite(t, exclude, []byte(".bench/phases.json\n"), 0o644)
			},
			// The commit lane is the walk's scaffold, so it leaves the source before the
			// review. The release of the landed source then finds no ignored residue.
			carry: map[int]func(*testing.T, landingFixture){1: func(t *testing.T, f landingFixture) {
				mustRemove(t, filepath.Join(f.creation.Path, ".bench", "phases.json"))
				review(t, f)
			}},
			names: func(f landingFixture) []string { return []string{"bench commit --in ", f.creation.Assignment.Label} },
		},
		{
			face:          faceSourceNotFenced,
			repairsSource: true,
			mutate: func(t *testing.T, _ string, creation Creation) {
				commitInWorktree(t, creation.Path, "stray.txt", "stray\n", "out of fence")
			},
			carry: map[int]func(*testing.T, landingFixture){0: func(t *testing.T, f landingFixture) {
				gitRun(t, f.creation.Path, "rm", "-q", "stray.txt")
				gitRun(t, f.creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "take the stray path out of the range")
			}},
		},
		{
			// The source moves past the tip the caller names, which is the state an
			// operator reaches when a review repair adds a commit. The review of that repair
			// covers the moved source, and the unmoved destination's head is its review base.
			face: faceSourceTipMismatch,
			mutate: func(t *testing.T, root string, creation Creation) {
				commitInWorktree(t, creation.Path, "owned.txt", "moved\n", "source moved past the named tip")
				refreshLandingEvidence(t, creation.Path, gitOutput(t, root, "rev-parse", "HEAD"))
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
			carry:  map[int]func(*testing.T, landingFixture){0: handMerge, 1: commitResolution, 2: review},
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
			// The planted merge state records no merge Git can finish, so the reviewer
			// clears it and finishes a real one.
			carry: map[int]func(*testing.T, landingFixture){0: func(t *testing.T, f landingFixture) {
				mustRemove(t, filepath.Join(gitOutput(t, f.creation.Path, "rev-parse", "--absolute-git-dir"), "MERGE_HEAD"))
				handMerge(t, f)
				commitResolution(t, f)
			}, 1: review},
		},
		{
			face:  faceResumeDestinationResidue,
			stage: stageResume,
			mutate: func(t *testing.T, root string, _ Creation) {
				mustWrite(t, filepath.Join(root, "tracked.txt"), []byte("dirty\n"), 0o644)
			},
			carry: map[int]func(*testing.T, landingFixture){0: discardDestination},
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
			// The reviewer restores main to the published commit, the parent of the one
			// movement the mutation added, so the resume advances the marker itself.
			carry: map[int]func(*testing.T, landingFixture){0: func(t *testing.T, f landingFixture) {
				gitRun(t, f.root, "reset", "--hard", "HEAD~1")
			}},
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
			carry: map[int]func(*testing.T, landingFixture){0: func(t *testing.T, f landingFixture) {
				gitRun(t, f.root, "checkout", "-q", "main")
			}},
		},
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
