// Build window landing tests. While a chunked build is open, the landing refuses each other
// landing unless the build is in a chunk gap.
package worktree

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/intent"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
)

// openBuildFixture is the spec-less landing fixture outside a chunk gap. Its source binds the
// staged closureSpec, so the build window holds main for every other landing.
func openBuildFixture(t *testing.T, request string) landingFixture {
	t.Helper()
	return seededLandingFixture(t, request, "", "", filepath.Join(t.TempDir(), "bench-home"), false, nil, closureAdmission)
}

// windowLanding is an open build and the unrelated light-path landing that request names.
func windowLanding(t *testing.T, request string) (build, light landingFixture) {
	t.Helper()
	build = openBuildFixture(t, request+"-build")
	return build, lightLandingOn(t, build, request)
}

// lightLandArgs are the arguments of the light-path landing that closes the light folder.
func lightLandArgs(request string, light landingFixture) []string {
	return ticketsOnlyLandArgs(request, light.base, light.tip, lightLandingSlug, light.creation.Path)
}

// otherBuild is the outcome of the second build that startOtherBuild records.
const otherBuild = "other"

// startOtherBuild records a second open build in the ledger at root, as a start of an
// approved second outcome records it. The fixture policy approves one outcome, so the
// record copies the delivery binding. The assignment that owns worktree binds the same
// staged spec for outcome otherBuild.
func startOtherBuild(root, worktree string) error {
	return intent.Transact(root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
		owners := intent.AssignmentsOwning(ledger.Assignments, worktree)
		if len(owners) != 1 || ledger.Commitment == nil || len(ledger.Commitment.Bindings) == 0 {
			return ledger, false, errors.New("the ledger holds no assignment and delivery binding to start the other build")
		}
		state := *ledger.Commitment
		binding := state.Bindings[0]
		binding.Assignment, binding.Request, binding.Outcome = owners[0].ID, owners[0].Request, otherBuild
		state.Bindings = append(slices.Clone(state.Bindings), binding)
		state.Claims = append(slices.Clone(state.Claims), intent.OutcomeClaim{Milestone: binding.Milestone, Outcome: otherBuild})
		ledger.Commitment = &state
		return ledger, true, nil
	}, nil)
}

// requireWindowRefusal fails unless r is the window refusal of build outcome, and unless
// main stays at the base of f. The refusal names the build, its deliverable, and wait.
func requireWindowRefusal(t *testing.T, r verbResult, f landingFixture, outcome, wait string) {
	t.Helper()
	detail := "landing destination " + leasedMain + " is held by open build " + outcome + " (" + closureSpec + "); wait for " + wait
	next, printed := landingFaceNext(r.stdout, detail)
	if r.exit != 1 || !printed || !strings.HasPrefix(next, landingRepair(faceLandBuildWindow, nil)) {
		t.Fatalf("landing = (%d, %q, %q), want the %s refusal %q", r.exit, r.stdout, r.stderr, faceLandBuildWindow, detail)
	}
	if main := gitOutput(t, f.root, "rev-parse", "main"); main != f.base {
		t.Fatalf("refused landing moved main to %s, want %s", main, f.base)
	}
}

// requireLanded fails unless r is a landed record.
func requireLanded(t *testing.T, r verbResult, after string) {
	t.Helper()
	if _, landed := landedField(r.stdout, "worktree"); r.exit != 0 || !landed {
		t.Fatalf("landing %s = (%d, %q, %q), want a landed record", after, r.exit, r.stdout, r.stderr)
	}
}

// TestLandRefusesAnUnrelatedLandingInAnOpenBuild is BW46, BW47, and BW57. A spec-less
// landing and a light-path landing of an unbound assignment each refuse before the gate.
func TestLandRefusesAnUnrelatedLandingInAnOpenBuild(t *testing.T) {
	t.Parallel()
	request := "land-window-unrelated"
	_, light := windowLanding(t, request)
	for _, args := range [][]string{specLessLandArgs(request, light.base, light.tip, light.creation.Path), lightLandArgs(request, light)} {
		requireWindowRefusal(t, runVerb(t, verbLand, light.call(args...)), light, commitmenttest.DeliveryOutcome, "its first chunk checkpoint")
	}
	if tally, err := os.ReadFile(light.tally); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("gate tally = %q, %v; want no gate run", tally, err)
	}
}

// TestWindowAdmitsTheBuildsOwnLanding is BW50 and BW94. A sibling that `--from` binds to the
// build lands the bound spec through the window, and its spec-less landing refuses.
func TestWindowAdmitsTheBuildsOwnLanding(t *testing.T) {
	t.Parallel()
	request := "land-window-sibling"
	build := openBuildFixture(t, request+"-build")
	// The sibling starts at the tip of its --from source, so that source starts at main.
	from := mustCreate(t, build.root, build.home, request+"-from", "from")
	commitmenttest.Admit(t, from.Path, request+"-from", closureSpec)
	if created := runVerb(t, verbCreate, verbCall{root: build.root, home: build.home, args: []string{"--request", request, "--label", request, "--from", from.Assignment.ID}}); created.exit != 0 {
		t.Fatalf("sibling create = (%d, %q, %q)", created.exit, created.stdout, created.stderr)
	}
	sibling, found, err := intent.FindAssignmentForRequest(build.root, request)
	if err != nil || !found {
		t.Fatalf("sibling assignment = %t, %v", found, err)
	}
	f := build
	f.creation = Creation{Path: sibling.Worktree, Assignment: sibling}
	commitInWorktree(t, f.creation.Path, "owned.txt", reviewedSource, "sibling source")
	refreshLandingEvidence(t, f.creation.Path, f.base)
	f.tip = gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	requireWindowRefusal(t, runVerb(t, verbLand, f.call(specLessLandArgs(request, f.base, f.tip, f.creation.Path)...)), f, commitmenttest.DeliveryOutcome, "its first chunk checkpoint")
	requireLanded(t, runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...)), "of the bound spec from the sibling")
}

// TestWindowClosesAtCompletion is BW51. The build's completion landing closes the window,
// so the refused unrelated landing then publishes. That landing closes its folder, so it
// folds the moved main first.
func TestWindowClosesAtCompletion(t *testing.T) {
	t.Parallel()
	request := "land-window-completion"
	build, light := windowLanding(t, request)
	requireWindowRefusal(t, runVerb(t, verbLand, light.call(lightLandArgs(request, light)...)), light, commitmenttest.DeliveryOutcome, "its first chunk checkpoint")
	requireLanded(t, runVerb(t, verbLand, build.call(landArgs(request+"-build", build.base, build.tip, build.creation.Path)...)), "of the build's completion")
	gitRun(t, light.creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "merge", "-q", "--no-ff", "-m", "fold the destination", "main")
	light.tip = gitOutput(t, light.creation.Path, "rev-parse", "HEAD")
	requireLanded(t, runVerb(t, verbLand, light.call(lightLandArgs(request, light)...)), "after the build's completion")
}

// TestWindowClosesOnBlock is BW52. A block drops the claim, so the window closes although
// the binding stays.
func TestWindowClosesOnBlock(t *testing.T) {
	t.Parallel()
	request := "land-window-block"
	build, light := windowLanding(t, request)
	requireWindowRefusal(t, runVerb(t, verbLand, light.call(lightLandArgs(request, light)...)), light, commitmenttest.DeliveryOutcome, "its first chunk checkpoint")
	mustNoError(t, commitrepo.Store{Root: build.root}.Block(commitmenttest.DeliveryOutcome, "reviewer paused delivery"))
	requireLanded(t, runVerb(t, verbLand, light.call(lightLandArgs(request, light)...)), "after the block")
}

// TestWindowClosesOnAbandon is BW53. The release of the one bound assignment leaves no
// active owner, so the window closes although the claim stays. The assignment has not
// landed, so its release retains the worktree and ends the assignment as cleanup-pending.
func TestWindowClosesOnAbandon(t *testing.T) {
	t.Parallel()
	request := "land-window-abandon"
	build, light := windowLanding(t, request)
	requireWindowRefusal(t, runVerb(t, verbLand, light.call(lightLandArgs(request, light)...)), light, commitmenttest.DeliveryOutcome, "its first chunk checkpoint")
	released := runVerb(t, verbRelease, build.call("--request", request+"-build", build.creation.Path))
	if owner, found, err := intent.FindAssignmentForRequest(build.root, request+"-build"); err != nil || !found || owner.State == intent.StateActive {
		t.Fatalf("bound assignment after its release = %+v (found=%t, err=%v), release = (%d, %q, %q); want an assignment that is not active", owner, found, err, released.exit, released.stdout, released.stderr)
	}
	requireLanded(t, runVerb(t, verbLand, light.call(lightLandArgs(request, light)...)), "after the release")
}

// gapPlan is a completion plan of three chunks with one ticket each, for the gap decisions
// that a chunk checkpoint and a chunk dispatch make.
var gapPlan = rr.Plan{Chunks: []rr.PlannedChunk{{ID: "C1", Tickets: []string{"1.md"}}, {ID: "C2", Tickets: []string{"2.md"}}, {ID: "C3", Tickets: []string{"3.md"}}}}

// TestOpenGapAdmitsAnUnrelatedLanding is BW61. The gap that the build's chunk checkpoint
// opens lets the unbound light-path landing publish.
func TestOpenGapAdmitsAnUnrelatedLanding(t *testing.T) {
	t.Parallel()
	request := "land-window-gap"
	build, light := windowLanding(t, request)
	mustNoError(t, commitrepo.Store{Root: build.creation.Path}.OpenGap(gapPlan, "C1"))
	requireLanded(t, runVerb(t, verbLand, light.call(lightLandArgs(request, light)...)), "in the build's open gap")
}

// TestWindowRecheckAtPublish is BW56 and BW70. While the landing's gate runs, another build
// starts, or the first dispatch of its next chunk closes its gap. Either refuses the landing
// in the publish decision, so main stays at the base. The landing is not bound to that build.
func TestWindowRecheckAtPublish(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name, wait     string
		before, during func(root, other string) error
	}{
		{name: "build-start", wait: "its first chunk checkpoint", during: startOtherBuild},
		{name: "gap-close", wait: "the checkpoint of chunk C2", before: func(root, other string) error {
			if err := startOtherBuild(root, other); err != nil {
				return err
			}
			return commitrepo.Store{Root: other}.OpenGap(gapPlan, "C1")
		}, during: func(_, other string) error { return commitrepo.Store{Root: other}.CloseGap(gapPlan, "2.md") }},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			request := "land-window-publish-" + row.name
			f := barrierLandingFixture(t, request)
			other := commitmenttest.Assignment(t, f.root, request+"-other")
			if row.before != nil {
				mustNoError(t, row.before(f.root, other))
			}
			var duringErr error
			done := f.during(func() { duringErr = row.during(f.root, other) })
			r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
			select {
			case err := <-done:
				mustNoError(t, err)
			default:
				t.Fatalf("landing finished without reaching its gate barrier: (%d, %q, %q)", r.exit, r.stdout, r.stderr)
			}
			mustNoError(t, duringErr)
			requireWindowRefusal(t, r, f.landingFixture, otherBuild, row.wait)
		})
	}
}

// TestWindowWinsOverTheLease is BW59. The window decides before the lease, so a landing
// that waits for a gap never waits for a live lease holder first.
func TestWindowWinsOverTheLease(t *testing.T) {
	t.Parallel()
	request := "land-window-lease"
	_, light := windowLanding(t, request)
	holdLandingLease(t, light.root, otherLandingID)
	r := runVerb(t, verbLand, light.call(lightLandArgs(request, light)...))
	requireWindowRefusal(t, r, light, commitmenttest.DeliveryOutcome, "its first chunk checkpoint")
	if strings.Contains(r.stdout, "is leased by assignment") {
		t.Fatalf("landing = %q, want the window refusal alone", r.stdout)
	}
}
