package worktree

import (
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
	f := publicLandingFixture(t, request, "", "")
	mustWrite(t, filepath.Join(f.root, "advance.txt"), []byte("destination advance\n"), 0o644)
	gitRun(t, f.root, "add", "advance.txt")
	gitRun(t, f.root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "destination advance")
	destination := gitOutput(t, f.root, "rev-parse", "HEAD")
	gitRun(t, f.creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local",
		"merge", "-q", "--no-ff", "-m", "fold the destination", destination)
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	r := runVerb(t, verbLand, f.call(landArgs(request, f.base, tip, f.creation.Path)...))
	next, printed := landingFaceNext(r.stdout, landingRefusalFaceByName(faceSourceNotFenced).detail)
	if r.exit != 1 || !printed {
		t.Fatalf("folded-base landing = (%d, %q, %q), want the fence refusal", r.exit, r.stdout, r.stderr)
	}
	if !strings.Contains(r.stdout, "wanted="+destination) {
		t.Fatalf("fence refusal = %q, want wanted=%s, the folded merge base", r.stdout, destination)
	}
	if !strings.Contains(next, landingBaseFlag(destination)) || strings.Contains(next, landingBaseFlag(f.base)) {
		t.Fatalf("fence route = %q, want --base re-pointed from %s to %s", next, f.base, destination)
	}
}
