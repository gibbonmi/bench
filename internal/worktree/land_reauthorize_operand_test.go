package worktree

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

// Land and reauthorize consume the shared operand resolver, so an identifier reaches
// the registered assignment rather than being read as a path. The land case proves it
// through the refusal: a mismatched request against the resolved assignment addresses
// that assignment by id, which a path-misread operand never could. The reauthorize
// case completes end to end through an id prefix.
func TestLandAndReauthorizeResolveIdentifierOperands(t *testing.T) {
	root, creation, home := newOwnedAssignment(t, "operand-land")
	chdir(t, root)
	var stdout bytes.Buffer
	code := LandCommand(root, home, "", []string{"--request", "wrong-token", "--base", creation.Assignment.Start, "--source-tip", creation.Assignment.Start, "-m", "land", creation.Assignment.Label}, &stdout, io.Discard)
	if code == 0 {
		t.Fatal("a mismatched request landed")
	}
	if !strings.Contains(stdout.String(), creation.Assignment.ID) {
		t.Fatalf("land refusal does not address the resolved assignment:\n%s", stdout.String())
	}
	var stderr bytes.Buffer
	code = ReauthorizeCommand(root, home, []string{"--assignment", creation.Assignment.ID, "--request", "rotated-token", "--base", creation.Assignment.Start, "--source-tip", "HEAD", creation.Assignment.ID[:10]}, io.Discard, &stderr)
	if code != 0 {
		t.Fatalf("reauthorize via an id prefix exited %d: %s", code, stderr.String())
	}
}
