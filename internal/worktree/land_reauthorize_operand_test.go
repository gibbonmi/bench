package worktree

import (
	"strings"
	"testing"
)

// Land and reauthorize consume the shared operand resolver, so an identifier reaches
// the registered assignment rather than being read as a path. The land case proves it
// through the refusal: a mismatched request against the resolved assignment addresses
// that assignment by id, which a path-misread operand never could. The reauthorize
// case completes end to end through an id prefix.
func TestLandAndReauthorizeResolveIdentifierOperands(t *testing.T) {
	f := newOwnedAssignment(t, "operand-land")
	chdir(t, f.root)
	landed := runVerb(t, verbLand, f.call("--request", "wrong-token", "--base", f.creation.Assignment.Start, "--source-tip", f.creation.Assignment.Start, "-m", "land", f.creation.Assignment.Label))
	if landed.exit == 0 {
		t.Fatal("a mismatched request landed")
	}
	if !strings.Contains(landed.stdout, f.creation.Assignment.ID) {
		t.Fatalf("land refusal does not address the resolved assignment:\n%s", landed.stdout)
	}
	reauthorized := runVerb(t, verbReauthorize, f.call("--assignment", f.creation.Assignment.ID, "--request", "rotated-token", "--base", f.creation.Assignment.Start, "--source-tip", "HEAD", f.creation.Assignment.ID[:10]))
	if reauthorized.exit != 0 {
		t.Fatalf("reauthorize via an id prefix exited %d: %s", reauthorized.exit, reauthorized.stderr)
	}
}
