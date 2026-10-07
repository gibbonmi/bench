// Resume refusal tests: destructive destination state, a non-ancestor review base, a stale
// marker, an evicted receipt, and the handback of a cause no landing route repairs. The
// interrupted landing a resume starts from is built here too.
package worktree

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/refusalroute"
)

func TestResumeLandCommandPublicRefusesDestructiveDestinationState(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	for _, journey := range []struct {
		name  string
		later bool
	}{
		{name: "PL25"},
		{name: "PL30", later: true},
	} {
		for _, tc := range []struct {
			name    string
			allowed bool
			detail  string
			setup   func(*testing.T, string)
		}{
			{name: "clean"},
			{name: "declared ignored output", allowed: true, setup: func(t *testing.T, root string) {
				mustWrite(t, filepath.Join(root, ".git", "info", "exclude"), []byte("dist/\n"), 0o644)
				mustMkdirAll(t, filepath.Join(root, "dist"), 0o755)
				mustWrite(t, filepath.Join(root, "dist", "out"), []byte("build output\n"), 0o600)
			}},
			{name: "staged changes", detail: "landing destination has staged changes", setup: func(t *testing.T, root string) {
				mustWrite(t, filepath.Join(root, "staged.txt"), []byte("staged\n"), 0o600)
				gitRun(t, root, "add", "staged.txt")
			}},
			{name: "tracked-worktree changes", detail: "landing destination has tracked-worktree changes", setup: func(t *testing.T, root string) {
				mustWrite(t, filepath.Join(root, "owned.txt"), []byte("caller bytes\n"), 0o600)
			}},
			// An untracked or ignored file outside the published tree is the operator's own,
			// so the resume completes around it.
			{name: "untracked file", allowed: true, setup: func(t *testing.T, root string) {
				mustWrite(t, filepath.Join(root, "untracked-file.txt"), []byte("caller bytes\n"), 0o600)
			}},
			{name: "undeclared ignored file", allowed: true, setup: func(t *testing.T, root string) {
				mustWrite(t, filepath.Join(root, ".git", "info", "exclude"), []byte(".env\n"), 0o644)
				mustWrite(t, filepath.Join(root, ".env"), []byte("caller bytes\n"), 0o600)
			}},
			{name: "nested repository", detail: "landing destination has nested repositories", setup: func(t *testing.T, root string) {
				plantNestedRepository(t, root)
			}},
		} {
			t.Run(journey.name+"/"+tc.name, func(t *testing.T) {
				request := "resume-destination-state-" + strings.ReplaceAll(journey.name+"-"+tc.name, " ", "-")
				f := publicLandingFixture(t, request, "private/output", "dist/")
				land := func(args ...string) (int, string, string) {
					var stdout, stderr bytes.Buffer
					cmd := descendant(t, binary, append([]string{"worktree", "land"}, args...)...)
					cmd.Dir, cmd.Stdout, cmd.Stderr = f.root, &stdout, &stderr
					return exitCode(cmd.Run()), stdout.String(), stderr.String()
				}
				if code, stdout, stderr := land("--request", request, "--base", f.base, "--source-tip", f.tip, "--spec", "x", "-m", "land reviewed source", f.creation.Path); code != 3 || !strings.Contains(stdout, "worktree=incomplete:release") || !strings.Contains(stderr, "worktree retained (ignored)") {
					t.Fatalf("interrupted landing = (%d, %q, %q)", code, stdout, stderr)
				}
				published := gitOutput(t, f.root, "rev-parse", "main")
				mustRemove(t, filepath.Join(f.creation.Path, "private", "output"))
				if journey.later {
					commitInWorktree(t, f.root, "destination-after-publication", "forward\n", "destination movement")
				}
				if tc.setup != nil {
					tc.setup(t, f.root)
				}
				before := resumeDestinationState(t, f.root)
				code, stdout, stderr := land("--resume", published, "--request", request, "--base", f.base, "--source-tip", f.tip, "--spec", "x", f.creation.Path)
				if tc.setup == nil || tc.allowed {
					if code != 0 || !strings.Contains(stdout, "worktree=released,census=0}") || stderr != "" {
						t.Fatalf("resume = (%d, %q, %q)", code, stdout, stderr)
					}
					return
				}
				if code != 1 || !strings.HasPrefix(stdout, "refused{detail="+tc.detail+",next=") || stderr != "" {
					t.Fatalf("destructive-state refusal = (%d, %q, %q), want detail %q", code, stdout, stderr, tc.detail)
				}
				// LRS5: the resume refusal carries the face's own repair with the
				// caller's resume continuation behind it, so the route survives the
				// interruption.
				next, printed := landingFaceNext(stdout, tc.detail)
				repair := landingRepair(faceResumeDestinationResidue, nil)
				if !printed || !strings.HasPrefix(next, repair) || !strings.Contains(next, "bench worktree land --resume ") {
					t.Fatalf("destructive-state route = (%q, %v), want %q ahead of the resume continuation", next, printed, repair)
				}
				if after := resumeDestinationState(t, f.root); after != before {
					t.Fatalf("refusal changed destination state:\nbefore:\n%safter:\n%s", before, after)
				}
				if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
					t.Fatalf("resume reran gate: tally=%q error=%v", got, err)
				}
			})
		}
	}
}

// LRS5: a source-side face refuses on the resume too. The resume rebuilds the caller's
// own continuation from the flags it passed, so the route ends with that resume rather
// than with a first-run landing the operator must not repeat.
func TestResumeLandCommandSourceRefusalNamesTheCallersResume(t *testing.T) {
	t.Parallel()
	request := "resume-source-not-clean"
	f := publicLandingFixture(t, request, "", "")
	interruptedLanding(t, f, request, f.tip)
	published := gitOutput(t, f.root, "rev-parse", "main")
	landingFixtureFor(t, faceSourceNotClean).mutate(t, f.root, f.creation)
	args := []string{"--resume", published, "--request", request, "--base", f.base, "--source-tip", f.tip, "--spec", "x", f.creation.Path}
	r := runVerb(t, verbLand, f.callWith(defaultJoins(), args...))
	resume := "bench worktree land --resume '" + published + "' --request <request> --base '" + f.base +
		"' --source-tip '" + f.tip + "' --spec 'x' '" + f.creation.Path + "'"
	want := landingRoute(faceSourceNotClean, resume, labelOf(f.creation))
	next, printed := landingFaceNext(r.stdout, refusalroute.Sentence(faceSourceNotClean))
	if r.exit != 1 || !printed || next != want {
		t.Fatalf("resume source refusal = (%d, %q, %q), want next %q", r.exit, r.stdout, r.stderr, want)
	}
}

func TestResumeLandCommandRefusesNonAncestorReviewBaseWithoutMutation(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	request := "resume-nonancestor-base"
	f := publicLandingFixture(t, request, "", "")
	run := func(args ...string) (int, string, string) {
		var stdout, stderr bytes.Buffer
		cmd := descendant(t, binary, append([]string{"worktree", "land"}, args...)...)
		cmd.Dir, cmd.Stdout, cmd.Stderr = f.root, &stdout, &stderr
		return exitCode(cmd.Run()), stdout.String(), stderr.String()
	}
	if code, stdout, stderr := run("--request", request, "--base", f.base, "--source-tip", f.tip, "--spec", "x", "-m", "land reviewed source", f.creation.Path); code != 0 || !strings.Contains(stdout, "worktree=released,census=0}") || stderr == "" {
		t.Fatalf("landing = (%d, %q, %q)", code, stdout, stderr)
	}
	published := gitOutput(t, f.root, "rev-parse", "main")
	commitInWorktree(t, f.root, "destination-after-publication", "forward\n", "destination movement")
	destination := gitOutput(t, f.root, "rev-parse", "main")
	code, stdout, stderr := run("--resume", published, "--request", request, "--base", destination, "--source-tip", f.tip, "--spec", "x", f.creation.Path)
	if code != 1 || !strings.Contains(stdout, "review base does not authenticate the published source") || stderr != "" {
		t.Fatalf("nonancestor resume = (%d, %q, %q)", code, stdout, stderr)
	}
	if got := gitOutput(t, f.root, "rev-parse", "main"); got != destination {
		t.Fatalf("nonancestor resume moved destination: got %s want %s", got, destination)
	}
	if got := gitOutput(t, f.root, "rev-parse", "refs/bench/green/main"); got != published {
		t.Fatalf("nonancestor resume moved project-green: got %s want %s", got, published)
	}
	if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
		t.Fatalf("nonancestor resume reran gate: tally=%q error=%v", got, err)
	}
}

func TestResumeLandCommandRefusesAbsentOrBehindMarkerAfterDestinationMoves(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		marker func(*testing.T, string, string, string)
	}{
		{name: "absent", marker: func(t *testing.T, root, _, _ string) { gitRun(t, root, "update-ref", "-d", "refs/bench/green/main") }},
		{name: "behind", marker: func(t *testing.T, root, base, _ string) { gitRun(t, root, "update-ref", "refs/bench/green/main", base) }},
		{name: "divergent", marker: func(t *testing.T, root, _, tip string) { gitRun(t, root, "update-ref", "refs/bench/green/main", tip) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := "resume-marker-" + tc.name
			f := publicLandingFixture(t, request, "", "")
			working := defaultJoins()
			broken := working
			broken.releaseLandingAssignment = func(joins, ambient, string, []string, io.Writer, io.Writer) int { return 1 }
			if r := runVerb(t, verbLand, f.callWith(broken, landArgs(request, f.base, f.tip, f.creation.Path)...)); r.exit != 3 || !strings.Contains(r.stdout, "worktree=incomplete:release") {
				t.Fatalf("interrupted landing = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
			}
			commitInWorktree(t, f.root, "destination-after-publication", "forward\n", "destination movement")
			destination := gitOutput(t, f.root, "rev-parse", "main")
			tc.marker(t, f.root, f.base, f.tip)
			args := []string{"--resume", gitOutput(t, f.root, "rev-parse", "main~1"), "--request", request, "--base", f.base, "--source-tip", f.tip, "--spec", "x", f.creation.Path}
			if r := runVerb(t, verbLand, f.callWith(working, args...)); r.exit != 1 || !strings.HasPrefix(r.stdout, "refused{detail=project-green marker") || len(r.stderr) != 0 {
				t.Fatalf("marker refusal = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
			}
			if got := gitOutput(t, f.root, "rev-parse", "main"); got != destination {
				t.Fatalf("marker refusal moved destination: got %s want %s", got, destination)
			}
			if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
				t.Fatalf("resume reran gate: tally=%q error=%v", got, err)
			}
		})
	}
}

func TestResumeLandCommandRefusesWhenTerminalReceiptWasEvicted(t *testing.T) {
	t.Parallel()
	request := "resume-evicted-receipt"
	f := publicLandingFixture(t, request, "", "")
	if r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...)); r.exit != 0 || !strings.Contains(r.stdout, "worktree=released,census=0}") {
		t.Fatalf("landing = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	published := gitOutput(t, f.root, "rev-parse", "main")
	for i := 0; i < intent.MaxCleanupReceipts; i++ {
		receipt := intent.CleanupReceipt{
			Schema: intent.CleanupReceiptSchema, Repo: f.root, Operation: "eviction-test", Target: filepath.Join(f.root, fmt.Sprintf("target-%d", i)),
			Fingerprint: intent.RequestDigest(fmt.Sprintf("evict-%d", i)), State: intent.ReceiptComplete, Phase: intent.ReceiptPhaseTerminal, Action: "removed",
		}
		if err := intent.PutCleanupReceipt(f.root, receipt); err != nil {
			t.Fatal(err)
		}
	}
	gitRun(t, f.root, "update-ref", "-d", "refs/bench/green/main")
	gitRun(t, f.root, "read-tree", "main^")
	gitRun(t, f.root, "checkout-index", "-a", "-f")
	staged := gitOutput(t, f.root, "diff", "--cached", "--name-only")
	args := []string{"--resume", published, "--request", request, "--base", f.base, "--source-tip", f.tip, "--spec", "x", f.creation.Path}
	if r := runVerb(t, verbLand, f.call(args...)); r.exit != 1 || !strings.Contains(r.stdout, "missing-terminal-receipt") || len(r.stderr) != 0 {
		t.Fatalf("evicted resume = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	if got := gitOutput(t, f.root, "rev-parse", "main"); got != published {
		t.Fatalf("evicted resume moved destination: got %s want %s", got, published)
	}
	if descendant(t, "git", "-C", f.root, "show-ref", "--verify", "--quiet", "refs/bench/green/main").Run() == nil {
		t.Fatal("evicted resume recreated project-green marker")
	}
	if got := gitOutput(t, f.root, "diff", "--cached", "--name-only"); got != staged {
		t.Fatalf("evicted resume reconciled checkout: got %q want %q", got, staged)
	}
	if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
		t.Fatalf("evicted resume reran gate: tally=%q error=%v", got, err)
	}
}

func resumeDestinationState(t *testing.T, root string) string {
	t.Helper()
	index, err := os.ReadFile(filepath.Join(root, ".git", "index"))
	if err != nil {
		t.Fatal(err)
	}
	assignment, err := os.ReadFile(filepath.Join(root, ".git", intent.Filename))
	if err != nil {
		t.Fatal(err)
	}
	marker := gitOutput(t, root, "rev-parse", "refs/bench/green/main")
	return fmt.Sprintf("refs=%smarker=%sindex=%x\nassignment=%x\nworktree=%s", gitOutput(t, root, "show-ref", "--head"), marker, index, assignment, resumeWorktreeBytes(t, root))
}

func resumeWorktreeBytes(t *testing.T, root string) string {
	t.Helper()
	var snapshot strings.Builder
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == filepath.Join(root, ".git") && entry.IsDir() {
			return filepath.SkipDir
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		fmt.Fprintf(&snapshot, "%s %o ", filepath.ToSlash(rel), info.Mode())
		if info.Mode().IsRegular() {
			body, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			fmt.Fprintf(&snapshot, "%x", body)
		} else if info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(path)
			if err != nil {
				return err
			}
			snapshot.WriteString(target)
		}
		snapshot.WriteByte('\n')
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return snapshot.String()
}

// resumeRerunOf is the caller's own resume of a landing that published published, as a
// resume that names the spec x prints it before its assignment resolves.
func resumeRerunOf(published, base, tip, path string) string {
	return "bench worktree land --resume '" + published + "' --request <request> --base '" + base +
		"' --source-tip '" + tip + "' --spec 'x' '" + path + "'"
}

// interruptedLanding runs the first landing with a release step that fails, so the
// landing publishes and exits incomplete.
func interruptedLanding(t *testing.T, f landingFixture, request, tip string) verbResult {
	t.Helper()
	broken := defaultJoins()
	broken.releaseLandingAssignment = func(joins, ambient, string, []string, io.Writer, io.Writer) int { return 1 }
	r := runVerb(t, verbLand, f.callWith(broken, landArgs(request, f.base, tip, f.creation.Path)...))
	if r.exit != 3 {
		t.Fatalf("interrupted landing = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	return r
}

// landingFaceResume drives a resume fixture. It interrupts a landing at the release step,
// applies the fixture's mutation to the published destination, and resumes.
func landingFaceResume(t *testing.T, fixture landingRefusalFixture, f landingFixture) verbResult {
	t.Helper()
	request := "landing-face-" + fixture.face
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	interruptedLanding(t, f, request, tip)
	published := gitOutput(t, f.root, "rev-parse", "main")
	fixture.mutate(t, f.root, f.creation)
	args := []string{"--resume", published, "--request", request, "--base", f.base, "--source-tip", tip, "--spec", "x", f.creation.Path}
	return runVerb(t, verbLand, f.callWith(defaultJoins(), args...))
}

// requireHandback fails t unless the refused record whose detail opens with sentence
// prints a route that opens with the reviewer marker.
func requireHandback(t *testing.T, r verbResult, sentence string) {
	t.Helper()
	next, printed := landingFaceNext(r.stdout, sentence)
	if r.exit != 1 || !printed || !strings.HasPrefix(next, reviewerRoute) {
		t.Fatalf("%q refusal = (%d, %q, %q), want exit 1 and a next= that opens with %q", sentence, r.exit, r.stdout, r.stderr, reviewerRoute)
	}
}

// TestLandingHandsBackAPathWithNoCanonicalForm is a land-handback site. A relative operand
// resolves against the working directory, and a removed working directory leaves it no
// canonical form. No Bench verb restores the directory, so the first run and the resume
// both hand the refusal back.
func TestLandingHandsBackAPathWithNoCanonicalForm(t *testing.T) {
	gone := filepath.Join(t.TempDir(), "gone")
	mustMkdirAll(t, gone, 0o755)
	chdir(t, gone)
	mustRemove(t, gone)
	root, home, tip := t.TempDir(), t.TempDir(), strings.Repeat("a", 40)
	for name, args := range map[string][]string{
		"first run": landArgs("no-canonical-path", tip, tip, "worktree"),
		"resume":    {"--resume", tip, "--request", "no-canonical-path", "--base", tip, "--source-tip", tip, "--spec", "x", "worktree"},
	} {
		t.Run(name, func(t *testing.T) {
			requireHandback(t, runVerb(t, verbLand, verbCall{root: root, home: home, args: args}), "worktree path is not canonical")
		})
	}
}

// TestResumeHandsBackADetachedDestination is a land-handback site on the resume path. The
// resume reads the destination's identity before it reads the publication, and a detached
// landing checkout is the primary checkout's own state, which only the reviewer clears.
func TestResumeHandsBackADetachedDestination(t *testing.T) {
	t.Parallel()
	request := "resume-detached-destination"
	f := publicLandingFixture(t, request, "", "")
	interruptedLanding(t, f, request, f.tip)
	published := gitOutput(t, f.root, "rev-parse", "main")
	gitRun(t, f.root, "checkout", "-q", "--detach")
	args := []string{"--resume", published, "--request", request, "--base", f.base, "--source-tip", f.tip, "--spec", "x", f.creation.Path}
	requireHandback(t, runVerb(t, verbLand, f.call(args...)), "landing checkout is not attached to the default branch")
}
