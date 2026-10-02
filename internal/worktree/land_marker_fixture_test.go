// The marker-interrupt landing fixture: a real landing whose gate deletes the destination's green marker.
package worktree

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/testrepo"
)

// markerRef is the green marker that the fixture records and that its gate deletes.
const markerRef = "refs/bench/green/main"

// markerLandingFixture is the landing fixture whose prospective gate deletes the
// destination's green marker. The fixture records the marker at the review base, so the
// landing reads that marker as its expected prior tip. After the publication, the real
// marker swap then finds no marker and fails. gradeSpec selects the spec-backed gate.
func markerLandingFixture(t *testing.T, request string, gradeSpec bool) landingFixture {
	t.Helper()
	f := landingFixtureWithGateStep(t, request, "", "", filepath.Join(t.TempDir(), "bench-home"), gradeSpec, deleteGreenMarker)
	gitRun(t, f.root, "update-ref", markerRef, f.base)
	return f
}

// deleteGreenMarker is the gate line that deletes the green marker in the destination
// repository at root. The gate script runs with a private command path, so the line
// declares git through the gate fixture.
func deleteGreenMarker(gate *testrepo.GateFixture, root string) string {
	return gate.Command("git") + " -C " + sanitize.ShellQuote(root) + " update-ref -d " + markerRef + "\n"
}

// interruptLandingAtMarker runs the landing on a marker landing fixture with the default
// joins. It requires the marker interrupt, and it returns the run and the published
// commit, which is the state a resume exists to finish.
func interruptLandingAtMarker(t *testing.T, f landingFixture, args ...string) (verbResult, string) {
	t.Helper()
	r := runVerb(t, verbLand, f.call(args...))
	if r.exit != 3 || !strings.Contains(r.stdout, "worktree=incomplete:marker") {
		t.Fatalf("interrupted landing = (%d, %q, %q), want exit 3 at the marker", r.exit, r.stdout, r.stderr)
	}
	return r, gitOutput(t, f.root, "rev-parse", "main")
}
