// The landing face producing fixtures: one fixture per land face of the shared registry,
// and the fault scaffolding the fixtures plant and clear.
package worktree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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
	// clear removes the fixture's own fault scaffolding once the face printed, which no
	// operator's tree holds, so the route repairs only the cause the face names.
	clear func(t *testing.T, f landingFixture)
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
		stageHandMerge(t, f.creation.Path, "owned.txt", landingResolution)
	}
	commitResolution := func(t *testing.T, f landingFixture) { commitHandMerge(t, f.creation.Path) }
	// The review of the repaired source refreshes its completion evidence.
	review := func(t *testing.T, f landingFixture) { refreshLandingEvidence(t, f.creation.Path, f.base) }
	// The commit lane is the walk's scaffold, so it leaves the source before the review. The
	// release of the landed source then finds no ignored residue.
	reviewWithoutLane := func(t *testing.T, f landingFixture) {
		mustRemove(t, filepath.Join(f.creation.Path, ".bench", "phases.json"))
		review(t, f)
	}
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
			// The dirty path is one the ticket writes, so the printed commit lands it.
			face:          faceSourceNotClean,
			repairsSource: true,
			mutate: func(t *testing.T, _ string, creation Creation) {
				mustWrite(t, filepath.Join(creation.Path, "owned.txt"), []byte("uncommitted repair\n"), 0o644)
				plantSourceCommitLane(t, creation)
			},
			carry: map[int]func(*testing.T, landingFixture){1: reviewWithoutLane},
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
			// The destination's gate now asks for a repair that only the source can carry,
			// and the reviewed source folds that destination, so the composed tree grades red
			// until the source commits the repair.
			face:          faceLandRed,
			stage:         stageComposition,
			repairsSource: true,
			mutate: func(t *testing.T, root string, creation Creation) {
				base := gitOutput(t, root, "rev-parse", "HEAD")
				commitGateLine(t, root, `IFS= read -r owned < owned.txt
[ "$owned" = "repaired" ]`)
				gitRun(t, creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "merge", "-q", "--no-ff", "-m", "fold the destination", "main")
				refreshLandingEvidence(t, creation.Path, base)
			},
			carry: map[int]func(*testing.T, landingFixture){
				0: func(t *testing.T, f landingFixture) {
					mustWrite(t, filepath.Join(f.creation.Path, "owned.txt"), []byte("repaired\n"), 0o644)
					plantSourceCommitLane(t, f.creation)
				},
				2: reviewWithoutLane,
			},
			names: func(f landingFixture) []string { return []string{"bench commit --in ", f.creation.Assignment.Label} },
		},
		{
			// The gate cannot open its lock, so it grades no tree. The fixture frees the lock
			// once the face printed, so the re-run after the diagnosis grades the composed
			// tree.
			face:   faceLandInfrastructure,
			stage:  stageComposition,
			mutate: func(t *testing.T, root string, _ Creation) { blockGateLock(t, root) },
			clear:  func(t *testing.T, f landingFixture) { mustRemove(t, gateLock(t, f.root)) },
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

// landingResolution resolves the landing fixture's conflict over owned.txt, which the
// destination and the reviewed source both change.
const landingResolution = "destination bytes\nreviewed repair\n"

// stageHandMerge is the reviewer's hand merge of main into checkout: the merge must
// conflict, and the resolution body at path is staged. A printed route names the commit
// that records the resolution as its own later step, so commitHandMerge is a separate call.
func stageHandMerge(t *testing.T, checkout, path, body string) {
	t.Helper()
	merge := descendant(t, "git", "-C", checkout, "merge", "--no-commit", "main")
	if out, err := merge.CombinedOutput(); err == nil || !strings.Contains(string(out), "CONFLICT") {
		t.Fatalf("hand merge = %v, %s; want the conflict", err, out)
	}
	mustWrite(t, filepath.Join(checkout, path), []byte(body), 0o644)
	gitRun(t, checkout, "add", path)
}

// commitHandMerge records the staged resolution of a hand merge in checkout.
func commitHandMerge(t *testing.T, checkout string) {
	t.Helper()
	gitRun(t, checkout, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "resolve the composition")
}

// plantSourceCommitLane declares a commit lane for the source in an ignored phase
// manifest, which leaves the source's status alone, so a printed `bench commit` has a lane
// to pass. The project's own gate grades the staged spec, which no commit in the source
// passes.
func plantSourceCommitLane(t *testing.T, creation Creation) {
	t.Helper()
	mustWrite(t, filepath.Join(creation.Path, ".bench", "phases.json"), []byte(`{"phases":[{"name":"build","argv":["true"]}],"lane":[{"name":"unit","argv":["true"]}]}`), 0o644)
	exclude := filepath.Join(gitOutput(t, creation.Path, "rev-parse", "--path-format=absolute", "--git-common-dir"), "info", "exclude")
	mustMkdirAll(t, filepath.Dir(exclude), 0o755)
	mustWrite(t, exclude, []byte(".bench/phases.json\n"), 0o644)
}

// commitGateLine appends one line to the destination's prospective gate and commits it,
// so the landing grades every composed tree with that line too.
func commitGateLine(t *testing.T, root, line string) {
	t.Helper()
	script := filepath.Join(root, ".bench", "gate-prospective.sh")
	body, err := os.ReadFile(script)
	mustNoError(t, err)
	mustWrite(t, script, append(body, []byte(line+"\n")...), 0o755)
	gitRun(t, root, "add", ".bench/gate-prospective.sh")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "destination gate line")
}

// gateLock is the gate's execution lock for the gate runs that checkout stores.
func gateLock(t *testing.T, checkout string) string {
	t.Helper()
	return mustAdminPath(t, checkout, "bench-gate.lock")
}

// blockGateLock puts a directory where the gate opens its execution lock, so every gate
// run that checkout stores stops on an infrastructure outcome until the directory goes.
func blockGateLock(t *testing.T, checkout string) {
	t.Helper()
	mustMkdirAll(t, gateLock(t, checkout), 0o755)
}
