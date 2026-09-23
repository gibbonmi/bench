// Release refusal tests for the landing command: destination and collision path tables in the refusal output.
package worktree

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode"

	"github.com/gibbonmi/bench/internal/diff"
)

func TestLandCommandRefusalListsDestinationPaths(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(t.TempDir(), "bench-home")
	creation := mustCreate(t, root, home, "refusal-destination", "refusal")
	stageLandSpec(t, root, creation.Path)
	base := gitOutput(t, root, "rev-parse", "HEAD")
	commitInWorktree(t, creation.Path, "owned.txt", "owned\n", "owned")
	mustWrite(t, filepath.Join(root, "tracked.txt"), []byte("dirty\n"), 0o644)
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs("refusal-destination", base, gitOutput(t, creation.Path, "rev-parse", "HEAD"), creation.Path), &stdout, &stderr)
	// The route reads from the registry, so the face's repair keeps one source. The
	// caller's own re-run rides behind it and this row does not pin it.
	wantNext := "next=" + landingRefusalFaceByName(faceDestinationNotClean).route("")
	if code != 1 || !strings.Contains(stdout.String(), wantNext) || !strings.Contains(stdout.String(), "paths_total=1") || !strings.Contains(stdout.String(), "refusal_paths[1]{path}:") || !strings.Contains(stdout.String(), "tracked.txt") || stderr.Len() != 0 {
		t.Fatalf("destination refusal = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
}

// The destination refuses only the untracked and ignored paths where the landing writes.
// The reviewed source adds owned.txt, so an ignored owned.txt collides, while an ignored
// .env and an untracked notes.txt stay the operator's own and are not named.
func TestLandCommandRefusalListsCollidingPaths(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(t.TempDir(), "bench-home")
	creation := mustCreate(t, root, home, "refusal-collision", "refusal")
	stageLandSpec(t, root, creation.Path)
	base := gitOutput(t, root, "rev-parse", "HEAD")
	commitInWorktree(t, creation.Path, "owned.txt", "owned\n", "owned")
	mustWrite(t, filepath.Join(root, ".git", "info", "exclude"), []byte("owned.txt\n.env\n"), 0o644)
	mustWrite(t, filepath.Join(root, "owned.txt"), []byte("operator bytes\n"), 0o600)
	mustWrite(t, filepath.Join(root, ".env"), []byte("SECRET=1\n"), 0o600)
	mustWrite(t, filepath.Join(root, "notes.txt"), []byte("notes\n"), 0o600)
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	// The collision proof runs after the source proofs, and the minimal fixture spec has
	// no valid coverage map, so the source range is stubbed to pass.
	j := defaultJoins()
	j.authorizeLandingSource = func(string, string, string) (diff.SourceRange, error) {
		return diff.SourceRange{Base: base, Tip: tip}, nil
	}
	var stdout, stderr bytes.Buffer
	code := landWith(j, root, home, "", landArgs("refusal-collision", base, tip, creation.Path), &stdout, &stderr)
	next, printed := landingFaceNext(stdout.String(), landingRefusalFaceByName(faceDestinationCollision).detail)
	if code != 1 || !printed || !strings.HasPrefix(next, landingRefusalFaceByName(faceDestinationCollision).route("")) ||
		!strings.Contains(stdout.String(), "refusal_paths[1]{path}:\n  owned.txt\n") || strings.Contains(stdout.String(), ".env") ||
		strings.Contains(stdout.String(), "notes.txt") || stderr.Len() != 0 {
		t.Fatalf("collision refusal = (%d, %q, %q), want owned.txt alone", code, stdout.String(), stderr.String())
	}
	if got, err := os.ReadFile(filepath.Join(root, "owned.txt")); err != nil || string(got) != "operator bytes\n" {
		t.Fatalf("refused landing changed the operator's file: %q, %v", got, err)
	}
}

func TestLandCommandRefusalKeepsControlBearingPathInOneTableRow(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(t.TempDir(), "bench-home")
	path := "bad\n\x1b,comma"
	mustWrite(t, filepath.Join(root, path), []byte("committed\n"), 0o644)
	gitRun(t, root, "add", "--", path)
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "control-bearing path")
	creation := mustCreate(t, root, home, "refusal-controls", "refusal")
	stageLandSpec(t, root, creation.Path)
	base := gitOutput(t, root, "rev-parse", "HEAD")
	commitInWorktree(t, creation.Path, "owned.txt", "owned\n", "owned")
	mustWrite(t, filepath.Join(root, path), []byte("residue\n"), 0o644)
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	// The minimal fixture spec has no valid coverage map, so the source proof would add
	// a second record; this test pins the destination record's shape alone.
	j := defaultJoins()
	j.authorizeLandingSource = func(string, string, string) (diff.SourceRange, error) {
		return diff.SourceRange{Base: base, Tip: tip}, nil
	}
	var stdout, stderr bytes.Buffer
	code := landWith(j, root, home, "", landArgs("refusal-controls", base, tip, creation.Path), &stdout, &stderr)
	unsafe := strings.ContainsFunc(stdout.String(), func(r rune) bool { return r != '\n' && unicode.IsControl(r) })
	wantPathRow := `  "bad\\n\\u001b,comma"` + "\n"
	if code != 1 || unsafe || !strings.Contains(stdout.String(), "refused{") || !strings.Contains(stdout.String(), "refusal_paths[1]{path}:\n"+wantPathRow) || strings.Count(stdout.String(), "\n") != 4 || stderr.Len() != 0 {
		t.Fatalf("control refusal = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
}

func TestReleaseCommandRefusalListsBoundedIgnoredPathsWithTrueTotal(t *testing.T) {
	t.Parallel()
	request := "landed-release-refusal-paths"
	root, creation, home := newOwnedAssignment(t, "release-refusal-paths")
	mustWrite(t, filepath.Join(root, ".git", "info", "exclude"), []byte("residue-*\n"), 0o644)
	for i := 0; i < 1003; i++ {
		name := fmt.Sprintf("residue-%04d", i)
		if i == 0 {
			name += " space[*]"
		}
		mustWrite(t, filepath.Join(creation.Path, name), []byte("residue\n"), 0o600)
	}

	var stdout, stderr bytes.Buffer
	code := ReleaseCommand(root, home, []string{"--request", request, creation.Path}, &stdout, &stderr)
	out := stderr.String()
	wantNext := "next=bench worktree release --request <request> '" + creation.Path + "'"
	if code != 1 || stdout.Len() != 0 || !strings.HasPrefix(out, "bench worktree release: worktree retained (ignored):") ||
		!strings.Contains(out, "paths_total=1003\n") || !strings.Contains(out, "refusal_paths[1000]{path}:") ||
		!strings.Contains(out, "residue-0000 space[*]") || strings.Contains(out, "residue-1000") || !strings.Contains(out, wantNext) || strings.Contains(out, request) {
		t.Fatalf("release refusal: code=%d stdout=%q prefix=%t total=%t table=%t hostile=%t bounded=%t next=%t", code, stdout.String(),
			strings.HasPrefix(out, "bench worktree release: worktree retained (ignored):"), strings.Contains(out, "paths_total=1003\n"),
			strings.Contains(out, "refusal_paths[1000]{path}:"), strings.Contains(out, "residue-0000 space[*]"), !strings.Contains(out, "residue-1000"), strings.Contains(out, wantNext))
	}
}

func TestReleaseCommandRefusalPointsThroughAssignmentForControlBearingPath(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, request string
	}{
		{name: "line-safe request", request: "release request[*]"},
		{name: "control-bearing request", request: "release\n\x1brequest"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newWorktreeRepo(t)
			home := filepath.Join(root, "home\n\x1bunsafe")
			creation := mustCreate(t, root, home, tc.request, "unsafe release pointer")
			wantNext := "bench worktree exec " + creation.Assignment.ID + " -- bench worktree release --request <request> ."

			var stdout, stderr bytes.Buffer
			code := ReleaseCommand(root, home, []string{"--request", tc.request, creation.Path}, &stdout, &stderr)
			out := stderr.String()
			unsafe := strings.ContainsFunc(out, func(r rune) bool { return r != '\n' && unicode.IsControl(r) })
			if code != 1 || stdout.Len() != 0 || unsafe || strings.Count(out, "\n") != 1 || !strings.Contains(out, "; next="+wantNext+"\n") || strings.Contains(out, tc.request) {
				t.Fatalf("release pointer: code=%d stdout=%q safe=%t one-line=%t next=%t stderr=%q", code, stdout.String(), !unsafe,
					strings.Count(out, "\n") == 1, strings.Contains(out, "; next="+wantNext+"\n"), out)
			}
		})
	}
}

func TestReleaseCommandRefusalHidesControlBearingRequestForSafePath(t *testing.T) {
	t.Parallel()
	request := "release\n\x1brequest"
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	creation := mustCreate(t, root, home, request, "safe release pointer")
	mustWrite(t, filepath.Join(root, ".git", "info", "exclude"), []byte("residue\n"), 0o644)
	mustWrite(t, filepath.Join(creation.Path, "residue"), []byte("retained\n"), 0o600)
	var stdout, stderr bytes.Buffer
	code := ReleaseCommand(root, home, []string{"--request", request, creation.Path}, &stdout, &stderr)
	out := stderr.String()
	wantNext := "next=bench worktree release --request <request> '" + creation.Path + "'"
	unsafe := strings.ContainsFunc(out, func(r rune) bool { return r != '\n' && unicode.IsControl(r) })
	if code != 1 || stdout.Len() != 0 || unsafe || !strings.Contains(out, wantNext) || strings.Contains(out, request) {
		t.Fatalf("safe-path release refusal = (%d, %q, %q), want stderr recovery %q without caller token or controls", code, stdout.String(), out, wantNext)
	}
}

// The destination proof reads tracked changes only, so an ignored file needs no
// declaration and an untracked file does not refuse.
func TestLandingDestinationAllowsUntrackedAndIgnoredFiles(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, ignore, output, declaration string
	}{
		{name: "declared", ignore: "dist/", output: "dist/bench", declaration: "{\"schema\":1,\"paths\":[\"dist/\"]}\n"},
		{name: "runtime", ignore: ".logs/", output: ".logs/gate.jsonl"},
		{name: "undeclared ignored", ignore: ".env", output: ".env"},
		{name: "untracked", ignore: ".env", output: "notes.txt"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newWorktreeRepo(t)
			mustMkdirAll(t, filepath.Join(root, ".bench"), 0o755)
			mustWrite(t, filepath.Join(root, ".gitignore"), []byte(tc.ignore+"\n"), 0o644)
			if tc.declaration != "" {
				mustWrite(t, filepath.Join(root, ".bench", "build-outputs.json"), []byte(tc.declaration), 0o644)
			}
			gitRun(t, root, "add", "-A")
			gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "declare output")
			mustMkdirAll(t, filepath.Dir(filepath.Join(root, tc.output)), 0o755)
			mustWrite(t, filepath.Join(root, tc.output), []byte("output\n"), 0o755)
			tip, branch, marker, fingerprint, err := landingDestination(root)
			if err != nil || tip == "" || branch != "main" || marker != "" || fingerprint == "" {
				t.Fatalf("destination %s = (%q, %q, %q, %q, %v)", tc.name, tip, branch, marker, fingerprint, err)
			}
		})
	}
}
