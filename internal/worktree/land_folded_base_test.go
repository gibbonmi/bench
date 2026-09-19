package worktree

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

// TestLandFenceRefusalNamesTheFoldedDefaultBase grades the fence face over a source that
// folded a later default-branch commit. The caller's --base predates the fold, so the
// default branch's own paths read as unfenced. The refusal names the fold's merge base as
// the wanted base, and its route re-points the caller's --base at that commit.
func TestLandFenceRefusalNamesTheFoldedDefaultBase(t *testing.T) {
	t.Parallel()
	request := "landing-folded-base"
	root, creation, base, _, _, home := publicLandingFixture(t, request, "", "")
	mustWrite(t, filepath.Join(root, "advance.txt"), []byte("destination advance\n"), 0o644)
	gitRun(t, root, "add", "advance.txt")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "destination advance")
	destination := gitOutput(t, root, "rev-parse", "HEAD")
	gitRun(t, creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local",
		"merge", "-q", "--no-ff", "-m", "fold the destination", destination)
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	next, printed := landingFaceNext(stdout.String(), landingRefusalFaceByName(faceSourceNotFenced).detail)
	if code != 1 || !printed {
		t.Fatalf("folded-base landing = (%d, %q, %q), want the fence refusal", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "wanted="+destination) {
		t.Fatalf("fence refusal = %q, want wanted=%s, the folded merge base", stdout.String(), destination)
	}
	if !strings.Contains(next, landingBaseFlag(destination)) || strings.Contains(next, landingBaseFlag(base)) {
		t.Fatalf("fence route = %q, want --base re-pointed from %s to %s", next, base, destination)
	}
}
