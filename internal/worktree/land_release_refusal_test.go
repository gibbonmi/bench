// Release refusal tests for the landing command: destination and collision path tables in
// the refusal output, and the checkout lookup a route takes for a path that is not
// line-safe.
package worktree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode"

	"github.com/gibbonmi/bench/internal/refusalroute"
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
	r := runVerb(t, verbLand, repoHome{root, home}.call(landArgs("refusal-destination", base, gitOutput(t, creation.Path, "rev-parse", "HEAD"), creation.Path)...))
	// The route reads from the registry, so the face's repair keeps one source. The
	// caller's own re-run rides behind it and this row does not pin it.
	wantNext := "next=" + landingRepair(faceDestinationNotClean, nil)
	if r.exit != 1 || !strings.Contains(r.stdout, wantNext) || !strings.Contains(r.stdout, "paths_total=1") || !strings.Contains(r.stdout, refusalPathsTable+"[1]{path}:") || !strings.Contains(r.stdout, "tracked.txt") || len(r.stderr) != 0 {
		t.Fatalf("destination refusal = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
}

// The destination refuses only the untracked and ignored paths where the landing writes.
// The reviewed source adds owned.txt, so an ignored owned.txt collides, while an ignored
// .env and an untracked notes.txt stay the operator's own and are not named.
func TestLandCommandRefusalListsCollidingPaths(t *testing.T) {
	t.Parallel()
	request := "refusal-collision"
	f := publicLandingFixture(t, request, "", "")
	root := f.root
	mustWrite(t, filepath.Join(root, ".git", "info", "exclude"), []byte("owned.txt\n.env\n"), 0o644)
	mustWrite(t, filepath.Join(root, "owned.txt"), []byte("operator bytes\n"), 0o600)
	mustWrite(t, filepath.Join(root, ".env"), []byte("SECRET=1\n"), 0o600)
	mustWrite(t, filepath.Join(root, "notes.txt"), []byte("notes\n"), 0o600)
	r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
	next, printed := landingFaceNext(r.stdout, refusalroute.Sentence(faceDestinationCollision))
	if r.exit != 1 || !printed || !strings.HasPrefix(next, landingRepair(faceDestinationCollision, nil)) ||
		!strings.Contains(r.stdout, refusalPathsTable+"[1]{path}:\n  owned.txt\n") || strings.Contains(r.stdout, ".env") ||
		strings.Contains(r.stdout, "notes.txt") || len(r.stderr) != 0 {
		t.Fatalf("collision refusal = (%d, %q, %q), want owned.txt alone", r.exit, r.stdout, r.stderr)
	}
	if got, err := os.ReadFile(filepath.Join(root, "owned.txt")); err != nil || string(got) != "operator bytes\n" {
		t.Fatalf("refused landing changed the operator's file: %q, %v", got, err)
	}
}

func TestLandCommandRefusalKeepsControlBearingPathInOneTableRow(t *testing.T) {
	t.Parallel()
	request := "refusal-controls"
	f := publicLandingFixture(t, request, "", "")
	path := "bad\n\x1b,comma"
	commitInWorktree(t, f.root, path, "committed\n", "control-bearing path")
	mustWrite(t, filepath.Join(f.root, path), []byte("residue\n"), 0o644)
	r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
	unsafe := strings.ContainsFunc(r.stdout, func(r rune) bool { return r != '\n' && unicode.IsControl(r) })
	wantPathRow := `  "bad\\n\\u001b,comma"` + "\n"
	if r.exit != 1 || unsafe || !strings.Contains(r.stdout, "refused{") || !strings.Contains(r.stdout, refusalPathsTable+"[1]{path}:\n"+wantPathRow) || strings.Count(r.stdout, "\n") != 4 || len(r.stderr) != 0 {
		t.Fatalf("control refusal = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
}

func TestReleaseCommandRefusalListsBoundedIgnoredPathsWithTrueTotal(t *testing.T) {
	t.Parallel()
	request := "landed-release-refusal-paths"
	f := newOwnedAssignment(t, "release-refusal-paths")
	mustWrite(t, filepath.Join(f.root, ".git", "info", "exclude"), []byte("residue-*\n"), 0o644)
	for i := 0; i < 1003; i++ {
		name := fmt.Sprintf("residue-%04d", i)
		if i == 0 {
			name += " space[*]"
		}
		mustWrite(t, filepath.Join(f.creation.Path, name), []byte("residue\n"), 0o600)
	}

	r := runVerb(t, verbRelease, f.call("--request", request, f.creation.Path))
	out := r.stderr
	wantNext := "next=bench worktree release --request <request> '" + f.creation.Path + "'"
	if r.exit != 1 || len(r.stdout) != 0 || !strings.HasPrefix(out, "bench worktree release: worktree retained (ignored):") ||
		!strings.Contains(out, "paths_total=1003\n") || !strings.Contains(out, refusalPathsTable+"[1000]{path}:") ||
		!strings.Contains(out, "residue-0000 space[*]") || strings.Contains(out, "residue-1000") || !strings.Contains(out, wantNext) || strings.Contains(out, request) {
		t.Fatalf("release refusal: code=%d stdout=%q prefix=%t total=%t table=%t hostile=%t bounded=%t next=%t", r.exit, r.stdout,
			strings.HasPrefix(out, "bench worktree release: worktree retained (ignored):"), strings.Contains(out, "paths_total=1003\n"),
			strings.Contains(out, refusalPathsTable+"[1000]{path}:"), strings.Contains(out, "residue-0000 space[*]"), !strings.Contains(out, "residue-1000"), strings.Contains(out, wantNext))
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

			r := runVerb(t, verbRelease, repoHome{root, home}.call("--request", tc.request, creation.Path))
			out := r.stderr
			unsafe := strings.ContainsFunc(out, func(r rune) bool { return r != '\n' && unicode.IsControl(r) })
			if r.exit != 1 || len(r.stdout) != 0 || unsafe || strings.Count(out, "\n") != 1 || !strings.Contains(out, "; next="+wantNext+"\n") || strings.Contains(out, tc.request) {
				t.Fatalf("release pointer: code=%d stdout=%q safe=%t one-line=%t next=%t stderr=%q", r.exit, r.stdout, !unsafe,
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
	r := runVerb(t, verbRelease, repoHome{root, home}.call("--request", request, creation.Path))
	out := r.stderr
	wantNext := "next=bench worktree release --request <request> '" + creation.Path + "'"
	unsafe := strings.ContainsFunc(out, func(r rune) bool { return r != '\n' && unicode.IsControl(r) })
	if r.exit != 1 || len(r.stdout) != 0 || unsafe || !strings.Contains(out, wantNext) || strings.Contains(out, request) {
		t.Fatalf("safe-path release refusal = (%d, %q, %q), want stderr recovery %q without caller token or controls", r.exit, r.stdout, out, wantNext)
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

// TestUnsafePathRouteUsesThePlaceholder is RR58 and RR60. A source path that is not
// line-safe cannot be pasted, and `bench worktree exec` refuses a Bench child, so the
// preflight re-run and the incomplete landing's resume both look the path up by the
// assignment id and print the checkout placeholder.
func TestUnsafePathRouteUsesThePlaceholder(t *testing.T) {
	t.Parallel()
	unsafeFixture := func(t *testing.T, request string) landingFixture {
		return publicLandingFixtureAtHome(t, request, "", "", filepath.Join(t.TempDir(), "bench\n\x1bhome"))
	}
	// requireCheckoutLookup fails t unless next looks the checkout up by the assignment
	// id and names the checkout placeholder, with no exec wrapper around a Bench child.
	requireCheckoutLookup := func(t *testing.T, kind, next string, printed bool, f landingFixture, stdout string) {
		t.Helper()
		lookup := "bench worktree path '" + f.creation.Assignment.ID + "'"
		if !printed || !strings.Contains(next, "<checkout>") || !strings.Contains(next, lookup) || strings.Contains(next, "bench worktree exec") {
			t.Fatalf("unsafe-path %s next = %q (printed=%t) in %q, want %q, <checkout>, and no bench worktree exec",
				kind, next, printed, stdout, lookup)
		}
	}
	t.Run("preflight re-run", func(t *testing.T) {
		t.Parallel()
		request := "unsafe-path-preflight"
		f := unsafeFixture(t, request)
		landingFixtureFor(t, faceSourceNotClean).mutate(t, f.root, f.creation)
		r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
		next, printed := landingFaceNext(r.stdout, refusalroute.Sentence(faceSourceNotClean))
		requireCheckoutLookup(t, "preflight", next, printed && r.exit == 1, f, r.stdout)
	})
	t.Run("incomplete resume", func(t *testing.T) {
		t.Parallel()
		request := "unsafe-path-incomplete"
		f := unsafeFixture(t, request)
		r := interruptedLanding(t, f, request, f.tip)
		next, printed := landedNext(r.stdout)
		requireCheckoutLookup(t, "resume", next, printed, f, r.stdout)
	})
}
