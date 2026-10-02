package worktree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var localCapturePaths = []string{
	"capture/IDEAS.md",
	"capture/learnings.md",
	"capture/session-handoff.md",
}

func addLocalCaptureIgnore(t *testing.T, root string, foreign string) (string, string) {
	t.Helper()
	ignore := strings.Join(localCapturePaths, "\n") + "\n"
	if foreign != "" {
		ignore += foreign + "\n"
	}
	mustWrite(t, filepath.Join(root, ".gitignore"), []byte(ignore), 0o644)
	gitRun(t, root, "add", ".gitignore")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "ignore local capture")
	return gitOutput(t, root, "rev-parse", "HEAD"), ignore
}

func writeLocalCapture(t *testing.T, root string) {
	t.Helper()
	mustMkdirAll(t, filepath.Join(root, "capture"), 0o755)
	for _, rel := range localCapturePaths {
		mustWrite(t, filepath.Join(root, filepath.FromSlash(rel)), []byte(rel+"\n"), 0o600)
	}
}

func TestLandCommandAllowsLocalCaptureInDestinationAndReleases(t *testing.T) {
	t.Parallel()
	request := "local-capture-land"
	f := specLessLandingFixture(t, request)
	base, _ := addLocalCaptureIgnore(t, f.root, "")
	gitRun(t, f.creation.Path, "rebase", "main")
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	writeLocalCapture(t, f.root)

	r := runVerb(t, verbLand, f.call(specLessLandArgs(request, base, tip, f.creation.Path)...))
	if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released,census=0}") {
		t.Fatalf("land with local capture = (%d, %q, %q), want released", r.exit, r.stdout, r.stderr)
	}
	for _, rel := range localCapturePaths {
		if _, err := os.Stat(filepath.Join(f.root, filepath.FromSlash(rel))); err != nil {
			t.Fatalf("local capture %q was not preserved: %v", rel, err)
		}
	}
}

func TestResumeLandCommandAllowsLocalCaptureInDestination(t *testing.T) {
	t.Parallel()
	request := "local-capture-resume"
	f := markerLandingFixture(t, request, false)
	base, _ := addLocalCaptureIgnore(t, f.root, "")
	gitRun(t, f.creation.Path, "rebase", "main")
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	writeLocalCapture(t, f.root)

	_, published := interruptLandingAtMarker(t, f, specLessLandArgs(request, base, tip, f.creation.Path)...)
	args := []string{"--resume", published, "--request", request, "--base", base, "--source-tip", tip, f.creation.Path}
	r := runVerb(t, verbLand, f.call(args...))
	if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released,census=0}") {
		t.Fatalf("resume with local capture = (%d, %q, %q), want released", r.exit, r.stdout, r.stderr)
	}
}

// An ignored file that no declaration names and no landed path touches is the operator's
// own, so the landing publishes around it and leaves its bytes in place.
func TestLandCommandKeepsUndeclaredIgnoredFileInDestination(t *testing.T) {
	t.Parallel()
	request := "local-capture-foreign"
	f := specLessLandingFixture(t, request)
	base, _ := addLocalCaptureIgnore(t, f.root, "foreign.tmp")
	gitRun(t, f.creation.Path, "rebase", "main")
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	writeLocalCapture(t, f.root)
	mustWrite(t, filepath.Join(f.root, "foreign.tmp"), []byte("foreign\n"), 0o600)

	r := runVerb(t, verbLand, f.call(specLessLandArgs(request, base, tip, f.creation.Path)...))
	if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released,census=0}") {
		t.Fatalf("land with an undeclared ignored file = (%d, %q, %q), want released", r.exit, r.stdout, r.stderr)
	}
	if got, err := os.ReadFile(filepath.Join(f.root, "foreign.tmp")); err != nil || string(got) != "foreign\n" {
		t.Fatalf("undeclared ignored file after the landing = %q, %v, want its bytes kept", got, err)
	}
}
