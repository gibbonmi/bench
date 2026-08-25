package worktree

import (
	"bytes"
	"context"
	"errors"
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
	root, creation, _, _, _, home := specLessLandingFixture(t, request)
	base, _ := addLocalCaptureIgnore(t, root, "")
	gitRun(t, creation.Path, "rebase", "main")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	writeLocalCapture(t, root)

	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", specLessLandArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 || !strings.Contains(stdout.String(), "worktree=released}") {
		t.Fatalf("land with local capture = (%d, %q, %q), want released", code, stdout.String(), stderr.String())
	}
	for _, rel := range localCapturePaths {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel))); err != nil {
			t.Fatalf("local capture %q was not preserved: %v", rel, err)
		}
	}
}

func TestResumeLandCommandAllowsLocalCaptureInDestination(t *testing.T) {
	t.Parallel()
	request := "local-capture-resume"
	root, creation, _, _, _, home := specLessLandingFixture(t, request)
	base, _ := addLocalCaptureIgnore(t, root, "")
	gitRun(t, creation.Path, "rebase", "main")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	writeLocalCapture(t, root)

	working := defaultJoins()
	broken := working
	broken.advanceLandingMarker = func(context.Context, string, string, string, string) error { return errors.New("interrupt") }
	var stdout, stderr bytes.Buffer
	if code := landWith(broken, root, home, "", specLessLandArgs(request, base, tip, creation.Path), &stdout, &stderr); code != 3 {
		t.Fatalf("interrupted land = (%d, %q, %q), want incomplete", code, stdout.String(), stderr.String())
	}
	published := gitOutput(t, root, "rev-parse", "main")
	stdout.Reset()
	stderr.Reset()
	args := []string{"--resume", published, "--request", request, "--base", base, "--source-tip", tip, creation.Path}
	if code := landWith(working, root, home, "", args, &stdout, &stderr); code != 0 || !strings.Contains(stdout.String(), "worktree=released}") {
		t.Fatalf("resume with local capture = (%d, %q, %q), want released", code, stdout.String(), stderr.String())
	}
}

func TestLandCommandIgnoredDiagnosticListsOnlyForeignResidue(t *testing.T) {
	t.Parallel()
	request := "local-capture-diagnostic"
	root, creation, _, _, _, home := specLessLandingFixture(t, request)
	base, _ := addLocalCaptureIgnore(t, root, "foreign.tmp")
	gitRun(t, creation.Path, "rebase", "main")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	writeLocalCapture(t, root)
	mustWrite(t, filepath.Join(root, "foreign.tmp"), []byte("foreign\n"), 0o600)

	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", specLessLandArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 1 || !strings.Contains(stdout.String(), "refusal_paths[1]{path}:\n  foreign.tmp\n") {
		t.Fatalf("mixed residue = (%d, %q, %q), want only foreign.tmp", code, stdout.String(), stderr.String())
	}
	for _, rel := range localCapturePaths {
		if strings.Contains(stdout.String(), rel) {
			t.Fatalf("diagnostic leaked allowed local capture %q: %q", rel, stdout.String())
		}
	}
}
