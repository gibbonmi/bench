// Release refusal tests for the landing command: destination and ignored path tables in the refusal output.
package worktree

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"unicode"

	"github.com/gibbonmi/bench/internal/diff"
)

func TestLandCommandRefusalListsDestinationPaths(t *testing.T) {
	root := newWorktreeRepo(t)
	bindEnv(t, "BENCH_HOME", filepath.Join(t.TempDir(), "bench-home"))
	creation := mustCreate(t, root, Home(), "refusal-destination", "refusal")
	stageLandSpec(t, root, creation.Path)
	base := gitOutput(t, root, "rev-parse", "HEAD")
	commitInWorktree(t, creation.Path, "owned.txt", "owned\n", "owned")
	mustWrite(t, filepath.Join(root, "dirty"), []byte("dirty\n"), 0o600)
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, Home(), "", landArgs("refusal-destination", base, gitOutput(t, creation.Path, "rev-parse", "HEAD"), creation.Path), &stdout, &stderr)
	if code != 1 || !strings.Contains(stdout.String(), "paths_total=1") || !strings.Contains(stdout.String(), "refusal_paths[1]{path}:") || !strings.Contains(stdout.String(), "dirty") || stderr.Len() != 0 {
		t.Fatalf("destination refusal = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
}

func TestLandCommandRefusalListsIgnoredPaths(t *testing.T) {
	root := newWorktreeRepo(t)
	bindEnv(t, "BENCH_HOME", filepath.Join(t.TempDir(), "bench-home"))
	creation := mustCreate(t, root, Home(), "refusal-ignored", "refusal")
	stageLandSpec(t, root, creation.Path)
	base := gitOutput(t, root, "rev-parse", "HEAD")
	commitInWorktree(t, creation.Path, "owned.txt", "owned\n", "owned")
	mustWrite(t, filepath.Join(root, ".git", "info", "exclude"), []byte("ignored/\n"), 0o644)
	mustMkdirAll(t, filepath.Join(root, "ignored"), 0o755)
	mustWrite(t, filepath.Join(root, "ignored", "residue"), []byte("residue\n"), 0o600)
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, Home(), "", landArgs("refusal-ignored", base, gitOutput(t, creation.Path, "rev-parse", "HEAD"), creation.Path), &stdout, &stderr)
	if code != 1 || !strings.Contains(stdout.String(), "paths_total=1") || !strings.Contains(stdout.String(), "refusal_paths[1]{path}:") || !strings.Contains(stdout.String(), "ignored/residue") || stderr.Len() != 0 {
		t.Fatalf("ignored refusal = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
}

func TestLandCommandRefusalKeepsControlBearingPathInOneTableRow(t *testing.T) {
	root := newWorktreeRepo(t)
	bindEnv(t, "BENCH_HOME", filepath.Join(t.TempDir(), "bench-home"))
	creation := mustCreate(t, root, Home(), "refusal-controls", "refusal")
	stageLandSpec(t, root, creation.Path)
	base := gitOutput(t, root, "rev-parse", "HEAD")
	commitInWorktree(t, creation.Path, "owned.txt", "owned\n", "owned")
	path := "bad\n\x1b,comma"
	mustWrite(t, filepath.Join(root, path), []byte("residue\n"), 0o600)
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	// The minimal fixture spec has no valid coverage map, so the source proof would add
	// a second record; this test pins the destination record's shape alone.
	oldAuthorize := authorizeLandingSource
	authorizeLandingSource = func(string, string, string) (diff.SourceRange, error) {
		return diff.SourceRange{Base: base, Tip: tip}, nil
	}
	t.Cleanup(func() { authorizeLandingSource = oldAuthorize })
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, Home(), "", landArgs("refusal-controls", base, tip, creation.Path), &stdout, &stderr)
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
	for _, tc := range []struct {
		name, request string
	}{
		{name: "line-safe request", request: "release request[*]"},
		{name: "control-bearing request", request: "release\n\x1brequest"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newWorktreeRepo(t)
			bindEnv(t, "BENCH_HOME", filepath.Join(root, "home\n\x1bunsafe"))
			creation := mustCreate(t, root, Home(), tc.request, "unsafe release pointer")
			wantNext := "bench worktree exec " + creation.Assignment.ID + " -- bench worktree release --request <request> ."

			var stdout, stderr bytes.Buffer
			code := ReleaseCommand(root, Home(), []string{"--request", tc.request, creation.Path}, &stdout, &stderr)
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
	request := "release\n\x1brequest"
	root := newWorktreeRepo(t)
	bindEnv(t, "BENCH_HOME", filepath.Join(root, ".bench-home"))
	creation := mustCreate(t, root, Home(), request, "safe release pointer")
	mustWrite(t, filepath.Join(root, ".git", "info", "exclude"), []byte("residue\n"), 0o644)
	mustWrite(t, filepath.Join(creation.Path, "residue"), []byte("retained\n"), 0o600)
	var stdout, stderr bytes.Buffer
	code := ReleaseCommand(root, Home(), []string{"--request", request, creation.Path}, &stdout, &stderr)
	out := stderr.String()
	wantNext := "next=bench worktree release --request <request> '" + creation.Path + "'"
	unsafe := strings.ContainsFunc(out, func(r rune) bool { return r != '\n' && unicode.IsControl(r) })
	if code != 1 || stdout.Len() != 0 || unsafe || !strings.Contains(out, wantNext) || strings.Contains(out, request) {
		t.Fatalf("safe-path release refusal = (%d, %q, %q), want stderr recovery %q without caller token or controls", code, stdout.String(), out, wantNext)
	}
}

func TestLandingDestinationAllowsDeclaredAndRuntimeIgnoredOutput(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, ignore, output, declaration string
	}{
		{name: "declared", ignore: "dist/", output: "dist/bench", declaration: "{\"schema\":1,\"paths\":[\"dist/\"]}\n"},
		{name: "runtime", ignore: ".logs/", output: ".logs/gate.jsonl"},
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
