package worktree

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

// The FT169 landing surface: one reviewed source must not pay six refusal
// round-trips. Each subtest drives the real land command through the public
// fixture and asserts the exact symptom the 2026-08-22 landing paid for.

func landSurface(t *testing.T, request string) (string, Creation, string, string) {
	t.Helper()
	root, creation, base, tip, _ := publicLandingFixture(t, request, "", "")
	return root, creation, base, tip
}

func landIn(t *testing.T, root string, args []string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, "", args, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

// seedCaptureBase commits capture files on the destination and rebases the
// source onto them, so both sides share the files in their merge base.
func seedCaptureBase(t *testing.T, root, source string, files map[string]string) string {
	t.Helper()
	for name, body := range files {
		mustMkdirAll(t, filepath.Dir(filepath.Join(root, filepath.FromSlash(name))), 0o755)
		mustWrite(t, filepath.Join(root, filepath.FromSlash(name)), []byte(body), 0o644)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "seed capture")
	gitRun(t, source, "rebase", "main")
	return gitOutput(t, root, "rev-parse", "HEAD")
}

func TestLandCommandAcceptsAbbreviatedIdentities(t *testing.T) {
	request := "land-surface-abbreviated"
	root, creation, base, tip := landSurface(t, request)
	code, stdout, stderr := landIn(t, root, landArgs(request, base[:12], tip[:12], creation.Path))
	if code != 0 || !strings.Contains(stdout, "worktree=released}") || !strings.Contains(stdout, "source_base="+base+",source_tip="+tip+",") {
		t.Fatalf("abbreviated landing = (%d, %q, %q), want released with full identities", code, stdout, stderr)
	}
	parents := strings.Fields(gitOutput(t, root, "rev-list", "--parents", "-n", "1", "main"))
	if len(parents) != 3 || parents[1] != base || parents[2] != tip {
		t.Fatalf("published parents = %q, want %s and %s", parents, base, tip)
	}
}

func TestLandCommandComposesCaptureOntoMovedDestination(t *testing.T) {
	request := "land-surface-capture-conflict"
	root, creation, _, _ := landSurface(t, request)
	base := seedCaptureBase(t, root, creation.Path, map[string]string{
		"capture/session-handoff.md": "handoff base\n",
		"capture/learnings.md":       "learnings base\n",
	})
	commitInWorktree(t, creation.Path, "capture/session-handoff.md", "handoff source\n", "source handoff")
	commitInWorktree(t, creation.Path, "capture/learnings.md", "learnings source\n", "source learnings")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	commitInWorktree(t, root, "capture/session-handoff.md", "handoff destination\n", "destination handoff")
	commitInWorktree(t, root, "capture/learnings.md", "learnings destination\n", "destination learnings")
	code, stdout, stderr := landIn(t, root, landArgs(request, base, tip, creation.Path))
	if code != 0 || !strings.Contains(stdout, "worktree=released}") {
		t.Fatalf("capture-conflict landing = (%d, %q, %q), want released", code, stdout, stderr)
	}
	if got := gitOutput(t, root, "show", "main:capture/session-handoff.md"); got != "handoff source" {
		t.Fatalf("published handoff = %q, want the source's", got)
	}
	if got := gitOutput(t, root, "show", "main:capture/learnings.md"); got != "learnings destination" {
		t.Fatalf("published learnings = %q, want the destination's", got)
	}
}

func TestLandCommandAuthorizesCaptureOutsideTheFence(t *testing.T) {
	request := "land-surface-capture-fence"
	root, creation, base, _ := landSurface(t, request)
	mustMkdirAll(t, filepath.Join(creation.Path, "capture"), 0o755)
	commitInWorktree(t, creation.Path, "capture/learnings.md", "learning\n", "phase-owned learning")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	code, stdout, stderr := landIn(t, root, landArgs(request, base, tip, creation.Path))
	if code != 0 || !strings.Contains(stdout, "worktree=released}") {
		t.Fatalf("capture landing = (%d, %q, %q), want released", code, stdout, stderr)
	}
	if got := gitOutput(t, root, "show", "main:capture/learnings.md"); got != "learning" {
		t.Fatalf("published learning = %q", got)
	}
}

func TestLandCommandReportsEveryRefusalInOnePreflight(t *testing.T) {
	request := "land-surface-one-preflight"
	root, creation, base, tip := landSurface(t, request)
	mustWrite(t, filepath.Join(root, "dirty"), []byte("dirty\n"), 0o600)
	mustWrite(t, filepath.Join(creation.Path, "scratch"), []byte("scratch\n"), 0o600)
	code, stdout, stderr := landIn(t, root, landArgs(request, base, tip, creation.Path))
	if code != 1 || !strings.Contains(stdout, "landing destination is not clean") || !strings.Contains(stdout, "reviewed source is not clean") {
		t.Fatalf("two-refusal preflight = (%d, %q, %q), want both refusals named", code, stdout, stderr)
	}
}

func TestLandCommandFenceRefusalNamesThePath(t *testing.T) {
	request := "land-surface-fence-path"
	root, creation, base, _ := landSurface(t, request)
	commitInWorktree(t, creation.Path, "stray.txt", "stray\n", "out of fence")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	code, stdout, stderr := landIn(t, root, landArgs(request, base, tip, creation.Path))
	if code != 1 || !strings.Contains(stdout, "stray.txt") {
		t.Fatalf("fence refusal = (%d, %q, %q), want the offending path named", code, stdout, stderr)
	}
}

func TestLandCommandConflictRefusalNamesThePath(t *testing.T) {
	request := "land-surface-conflict-path"
	root, creation, base, tip := landSurface(t, request)
	commitInWorktree(t, root, "owned.txt", "destination bytes\n", "destination conflict")
	code, stdout, stderr := landIn(t, root, landArgs(request, base, tip, creation.Path))
	if code != 1 || !strings.Contains(stdout, "composition conflict: textual") || !strings.Contains(stdout, "owned.txt") {
		t.Fatalf("conflict refusal = (%d, %q, %q), want the conflicted path named", code, stdout, stderr)
	}
}
