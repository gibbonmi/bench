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
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(t.TempDir(), "bench-home")
	creation := mustCreate(t, root, home, "refusal-destination", "refusal")
	stageLandSpec(t, root, creation.Path)
	base := gitOutput(t, root, "rev-parse", "HEAD")
	commitInWorktree(t, creation.Path, "owned.txt", "owned\n", "owned")
	mustWrite(t, filepath.Join(root, "dirty"), []byte("dirty\n"), 0o600)
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs("refusal-destination", base, gitOutput(t, creation.Path, "rev-parse", "HEAD"), creation.Path), &stdout, &stderr)
	// The route reads from the registry, so the face's repair keeps one source. The
	// caller's own re-run rides behind it and this row does not pin it.
	wantNext := "next=" + landingRefusalFaceByName(faceDestinationNotClean).route("")
	if code != 1 || !strings.Contains(stdout.String(), wantNext) || !strings.Contains(stdout.String(), "paths_total=1") || !strings.Contains(stdout.String(), "refusal_paths[1]{path}:") || !strings.Contains(stdout.String(), "dirty") || stderr.Len() != 0 {
		t.Fatalf("destination refusal = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
}

func TestLandCommandRefusalListsIgnoredPaths(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(t.TempDir(), "bench-home")
	creation := mustCreate(t, root, home, "refusal-ignored", "refusal")
	stageLandSpec(t, root, creation.Path)
	base := gitOutput(t, root, "rev-parse", "HEAD")
	commitInWorktree(t, creation.Path, "owned.txt", "owned\n", "owned")
	mustWrite(t, filepath.Join(root, ".git", "info", "exclude"), []byte("ignored/\n"), 0o644)
	mustMkdirAll(t, filepath.Join(root, "ignored"), 0o755)
	mustWrite(t, filepath.Join(root, "ignored", "residue"), []byte("residue\n"), 0o600)
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs("refusal-ignored", base, gitOutput(t, creation.Path, "rev-parse", "HEAD"), creation.Path), &stdout, &stderr)
	if code != 1 || !strings.Contains(stdout.String(), "paths_total=1") || !strings.Contains(stdout.String(), "refusal_paths[1]{path}:") || !strings.Contains(stdout.String(), "ignored/residue") || stderr.Len() != 0 {
		t.Fatalf("ignored refusal = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
	// LRS10 and LRS11. The operator chooses between the two routes out of undeclared
	// residue, so the field names the declaration file and the removal of the exact path.
	next, printed := landingFaceNext(stdout.String(), landingRefusalFaceByName(faceDestinationResidue).detail)
	if !printed || !strings.Contains(next, ".bench/build-outputs.json") || !strings.Contains(next, "rm -rf 'ignored/residue'") {
		t.Fatalf("ignored refusal next = (%v, %q), want the declaration file and the removal command", printed, next)
	}
}

// LRS10 edge: the residue face names no path when its own read fails, so the removal
// command it states takes the placeholder the operator resolves by hand. The ignored
// inventory refuses a control-bearing path before the route reads it, so this empty
// list is the one state the placeholder covers.
func TestLandCommandResidueRouteHoldsThePlaceholderWithoutPaths(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name  string
		setup func(*testing.T, string)
	}{
		{name: "unreadable inventory", setup: func(t *testing.T, root string) {
			mustWrite(t, filepath.Join(root, ".git", "info", "exclude"), []byte("ignored/\n"), 0o644)
			mustMkdirAll(t, filepath.Join(root, "ignored"), 0o755)
			mustWrite(t, filepath.Join(root, "ignored", "res\x1bidue"), []byte("residue\n"), 0o600)
		}},
		// A directory in the declaration's place is not a regular file, so the
		// declaration read fails.
		{name: "unreadable declaration", setup: func(t *testing.T, root string) {
			mustMkdirAll(t, filepath.Join(root, ".bench", "build-outputs.json"), 0o755)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			request := "refusal-residue-" + strings.ReplaceAll(tc.name, " ", "-")
			root := newWorktreeRepo(t)
			home := filepath.Join(t.TempDir(), "bench-home")
			creation := mustCreate(t, root, home, request, "refusal")
			stageLandSpec(t, root, creation.Path)
			base := gitOutput(t, root, "rev-parse", "HEAD")
			commitInWorktree(t, creation.Path, "owned.txt", "owned\n", "owned")
			tc.setup(t, root)
			var stdout, stderr bytes.Buffer
			code := LandCommand(root, home, "", landArgs(request, base, gitOutput(t, creation.Path, "rev-parse", "HEAD"), creation.Path), &stdout, &stderr)
			next, printed := landingFaceNext(stdout.String(), landingRefusalFaceByName(faceDestinationResidue).detail)
			if code != 1 || !printed || !strings.Contains(next, "rm -rf <refusal_paths entries>") || strings.Contains(stdout.String(), "paths_total=") {
				t.Fatalf("%s residue route = (%d, %q, %q), next %q, want the placeholder removal and no path table", tc.name, code, stdout.String(), stderr.String(), next)
			}
		})
	}
}

func TestLandCommandRefusalKeepsControlBearingPathInOneTableRow(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(t.TempDir(), "bench-home")
	creation := mustCreate(t, root, home, "refusal-controls", "refusal")
	stageLandSpec(t, root, creation.Path)
	base := gitOutput(t, root, "rev-parse", "HEAD")
	commitInWorktree(t, creation.Path, "owned.txt", "owned\n", "owned")
	path := "bad\n\x1b,comma"
	mustWrite(t, filepath.Join(root, path), []byte("residue\n"), 0o600)
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
			tip, branch, marker, fingerprint, err := landingDestination(defaultJoins(), root)
			if err != nil || tip == "" || branch != "main" || marker != "" || fingerprint == "" {
				t.Fatalf("destination %s = (%q, %q, %q, %q, %v)", tc.name, tip, branch, marker, fingerprint, err)
			}
		})
	}
}
