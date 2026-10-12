// Reviewer override landing tests. An override lets one assignment's landing through every
// open build, and the landing that publishes consumes it.
package worktree

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/testrepo"
)

// landingOverrideReason is the reason of every override that these tests grant.
const landingOverrideReason = "back out the broken landing"

// grantOverride records the reviewer's override for the source assignment of f through the
// override owner, and it returns that record.
func grantOverride(t *testing.T, f landingFixture) intent.LandingOverride {
	t.Helper()
	mustNoError(t, commitrepo.Store{Root: f.root}.Override(f.creation.Assignment.ID, landingOverrideReason))
	return intent.LandingOverride{Assignment: f.creation.Assignment.ID, Reason: landingOverrideReason}
}

// requireOverrides fails unless the ledger at root holds exactly want.
func requireOverrides(t *testing.T, root string, want []intent.LandingOverride, after string) {
	t.Helper()
	got, err := commitmenttest.Overrides(root)
	mustNoError(t, err)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("overrides %s = %+v, want %+v", after, got, want)
	}
}

// TestOverrideAdmitsOneLandingThroughEveryWindow is BW72 and BW73. Two open builds are
// outside their gaps. The overridden landing publishes through both, and a landing of
// another assignment still refuses.
func TestOverrideAdmitsOneLandingThroughEveryWindow(t *testing.T) {
	t.Parallel()
	request := "land-override-every"
	build, light := windowLanding(t, request)
	mustNoError(t, startOtherBuild(build.root, commitmenttest.Assignment(t, build.root, request+"-other")))
	unrelated := lightLandingOn(t, build, request+"-unrelated")
	grantOverride(t, light)
	requireWindowRefusal(t, runVerb(t, verbLand, unrelated.call(lightLandArgs(request+"-unrelated", unrelated)...)), unrelated, commitmenttest.DeliveryOutcome, "its first chunk checkpoint")
	requireLanded(t, runVerb(t, verbLand, light.call(lightLandArgs(request, light)...)), "with the override")
}

// TestOverrideIsConsumedAtPublish is BW74. The publication drops the override of its
// assignment, so the override admits no later publication. The override of the build
// assignment stays, because only its own landing consumes it.
func TestOverrideIsConsumedAtPublish(t *testing.T) {
	t.Parallel()
	request := "land-override-consumed"
	build, light := windowLanding(t, request)
	kept := grantOverride(t, build)
	grantOverride(t, light)
	requireLanded(t, runVerb(t, verbLand, light.call(lightLandArgs(request, light)...)), "with the override")
	requireOverrides(t, light.root, []intent.LandingOverride{kept}, "after the publication")
}

// TestOverrideSurvivesARefusedLanding is BW75. The overridden landing refuses red at its
// gate and keeps the override, so its rerun publishes without a second grant.
func TestOverrideSurvivesARefusedLanding(t *testing.T) {
	t.Parallel()
	request := "land-override-red"
	red := filepath.Join(t.TempDir(), "red")
	mustWrite(t, red, []byte("red\n"), 0o600)
	build := seededLandingFixture(t, request+"-build", "", "", filepath.Join(t.TempDir(), "bench-home"), false, func(*testrepo.GateFixture, string) string {
		return "[ ! -e " + sanitize.ShellQuote(red) + " ]\n"
	}, closureAdmission)
	light := lightLandingOn(t, build, request)
	granted := grantOverride(t, light)
	args := lightLandArgs(request, light)
	if r := runVerb(t, verbLand, light.call(args...)); r.exit != 1 || !strings.Contains(r.stdout, "refused{detail=prospective authorization refused: ") {
		t.Fatalf("red landing = (%d, %q, %q), want the red refusal", r.exit, r.stdout, r.stderr)
	}
	requireOverrides(t, light.root, []intent.LandingOverride{granted}, "after the red refusal")
	mustRemove(t, red)
	requireLanded(t, runVerb(t, verbLand, light.call(args...)), "of the rerun")
}
