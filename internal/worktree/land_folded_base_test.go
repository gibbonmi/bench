package worktree

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestLandFenceRefusalOverAFoldedSourceKeepsTheCallerBase grades the fence face over a
// source that folded a later default-branch commit and then wrote one unfenced path. The
// folded path arrives unchanged from the default-branch tip, so the fence authorizes it,
// and the refusal names only the build's own path. The route keeps the caller's --base,
// because a later base would narrow the range the landing authorizes.
func TestLandFenceRefusalOverAFoldedSourceKeepsTheCallerBase(t *testing.T) {
	t.Parallel()
	request := "landing-folded-base"
	f := publicLandingFixture(t, request, "", "")
	mustWrite(t, filepath.Join(f.root, "advance.txt"), []byte("destination advance\n"), 0o644)
	gitRun(t, f.root, "add", "advance.txt")
	gitRun(t, f.root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "destination advance")
	destination := gitOutput(t, f.root, "rev-parse", "HEAD")
	gitRun(t, f.creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local",
		"merge", "-q", "--no-ff", "-m", "fold the destination", destination)
	commitInWorktree(t, f.creation.Path, "stray.txt", "stray\n", "out of fence")
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	r := runVerb(t, verbLand, f.call(landArgs(request, f.base, tip, f.creation.Path)...))
	face := landingRefusalFaceByName(faceSourceNotFenced)
	next, printed := landingFaceNext(r.stdout, face.detail)
	if r.exit != 1 || !printed {
		t.Fatalf("folded-source landing = (%d, %q, %q), want the fence refusal", r.exit, r.stdout, r.stderr)
	}
	if !strings.Contains(r.stdout, refusalPathsTable+"[1]{path}:\n  stray.txt\n") {
		t.Fatalf("fence refusal = %q, want only the build's stray.txt in refusal_paths", r.stdout)
	}
	if !strings.HasPrefix(next, face.route("")) || !strings.Contains(next, landingBaseFlag(f.base)) || strings.Contains(next, landingBaseFlag(destination)) {
		t.Fatalf("fence route = %q, want the fence repair at the caller's --base %s", next, f.base)
	}
}
