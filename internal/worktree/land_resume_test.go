// Resume tests for an interrupted landing: identity binding, marker completion, and destination-state refusals.
package worktree

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

func TestResumeLandCommandFollowupFailureExitsIncomplete(t *testing.T) {
	t.Parallel()
	request := "resume-release-incomplete"
	f := publicLandingFixture(t, request, "private/output", "dist/")
	if r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...)); r.exit != 3 {
		t.Fatalf("first incomplete exit = %d, want 3; stdout=%q stderr=%q", r.exit, r.stdout, r.stderr)
	}
	published := gitOutput(t, f.root, "rev-parse", "main")
	args := []string{"--resume", published, "--request", request, "--base", f.base, "--source-tip", f.tip, "--spec", "x", f.creation.Path}
	r := runVerb(t, verbLand, f.call(args...))
	if r.exit != 3 || !strings.Contains(r.stdout, "worktree=incomplete:release") {
		t.Fatalf("resume incomplete = (%d, %q, %q), want exit 3", r.exit, r.stdout, r.stderr)
	}
}

func TestLandCommandIncompleteNextUsesAssignmentPointerForUnsafePath(t *testing.T) {
	t.Parallel()
	request := "incomplete-unsafe-path"
	home := filepath.Join(t.TempDir(), "bench\n\x1bhome")
	f := publicLandingFixtureAtHome(t, request, "private/output", "dist/", home)
	r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
	wantNext := "next=bench worktree exec " + f.creation.Assignment.ID + " -- bench worktree land --resume '"
	unsafe := strings.ContainsRune(r.stdout, '\x1b') || strings.Count(r.stdout, "\n") != 1
	if r.exit != 3 || unsafe || !strings.Contains(r.stdout, wantNext) || !strings.Contains(r.stdout, " --spec 'x' .,census=0}") {
		t.Fatalf("unsafe-path incomplete = (%d, %q, %q), want one safe pointer record containing %q", r.exit, r.stdout, r.stderr, wantNext)
	}
}

func TestLandCommandPublicResumeCompletesPublishedReleaseWithoutRepublishing(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	request := "public-land-resume"
	f := publicLandingFixture(t, request, "private/output", "dist/")
	land := func(args ...string) (int, string, string) {
		var stdout, stderr bytes.Buffer
		cmd := descendant(t, binary, append([]string{"worktree", "land"}, args...)...)
		cmd.Dir, cmd.Stdout, cmd.Stderr = f.root, &stdout, &stderr
		return exitCode(cmd.Run()), stdout.String(), stderr.String()
	}
	code, stdout, stderr := land("--request", request, "--base", f.base, "--source-tip", f.tip, "--spec", "x", "-m", "land reviewed source", f.creation.Path)
	if code != 3 || !strings.Contains(stdout, "source_base="+f.base) || !strings.Contains(stdout, "worktree=incomplete:release") || !strings.Contains(stderr, "worktree retained (ignored)") {
		t.Fatalf("interrupted landing = (%d, %q, %q)", code, stdout, stderr)
	}
	published := gitOutput(t, f.root, "rev-parse", "main")
	if err := os.Remove(filepath.Join(f.creation.Path, "private", "output")); err != nil {
		t.Fatal(err)
	}
	commitInWorktree(t, f.root, "destination-after-publication", "forward\n", "destination movement")
	// LF4: the destination declares Go build inputs only after publication, so the resume
	// owns a broker refresh it cannot finish. This fixture carries no build entry point,
	// and no seal these sources could match, so the effect reports failed at exit 3 while
	// the published commit stands. The manifest is committed because an untracked file
	// would trip the resume's own untracked-collision proof and hide the effect behind
	// another refusal.
	commitLandingBuildInputs(t, f.root)
	destination := gitOutput(t, f.root, "rev-parse", "main")
	code, stdout, stderr = land("--resume", published, "--request", request, "--base", f.base, "--source-tip", f.tip, "--spec", "x", f.creation.Path)
	if code != 3 || !strings.Contains(stdout, "source_base="+f.base) || !strings.Contains(stdout, "worktree=incomplete:refresh,next=") || !strings.Contains(stderr, "landing refresh failed") {
		t.Fatalf("resume = (%d, %q, %q)", code, stdout, stderr)
	}
	if got := gitOutput(t, f.root, "rev-parse", "main"); got != destination {
		t.Fatalf("resume moved destination backward: main=%s destination=%s", got, destination)
	}
	if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
		t.Fatalf("resume reran gate: tally=%q error=%v", got, err)
	}
	// The release settled on the resume above, so this call reaches the terminal path.
	// That path owns the effects too, and the refresh still cannot finish.
	code, stdout, stderr = land("--resume", published, "--request", request, "--base", f.base, "--source-tip", f.tip, "--spec", "x", f.creation.Path)
	if code != 3 || !strings.Contains(stdout, "source_base="+f.base) || !strings.Contains(stdout, "worktree=incomplete:refresh,next=") {
		t.Fatalf("completed resume = (%d, %q, %q)", code, stdout, stderr)
	}
}

func TestResumeLandCommandPublicBindsPublishedLandingIdentity(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	request := "public-resume-identity"
	f := publicLandingFixture(t, request, "private/output", "dist/")
	land := func(args ...string) (int, string, string) {
		var stdout, stderr bytes.Buffer
		cmd := descendant(t, binary, append([]string{"worktree", "land"}, args...)...)
		cmd.Dir, cmd.Stdout, cmd.Stderr = f.root, &stdout, &stderr
		return exitCode(cmd.Run()), stdout.String(), stderr.String()
	}

	if code, stdout, stderr := land("--request", request, "--base", f.base, "--source-tip", f.tip, "--spec", "x", "-m", "land reviewed source", f.creation.Path); code != 3 || !strings.Contains(stdout, "worktree=incomplete:release") || stderr == "" {
		t.Fatalf("interrupted landing = (%d, %q, %q)", code, stdout, stderr)
	}
	published := gitOutput(t, f.root, "rev-parse", "main")
	commitInWorktree(t, f.root, "destination-after-publication", "forward\n", "destination movement")
	destination := gitOutput(t, f.root, "rev-parse", "main")
	for _, tc := range []struct {
		name, published, source string
	}{
		{name: "wrong-published", published: destination, source: f.tip},
		{name: "wrong-source", published: published, source: f.base},
	} {
		t.Run(tc.name, func(t *testing.T) {
			code, stdout, stderr := land("--resume", tc.published, "--request", request, "--base", f.base, "--source-tip", tc.source, "--spec", "x", f.creation.Path)
			if code != 1 || !strings.HasPrefix(stdout, "refused{detail=") || stderr != "" {
				t.Fatalf("resume refusal = (%d, %q, %q)", code, stdout, stderr)
			}
			if got := gitOutput(t, f.root, "rev-parse", "main"); got != destination {
				t.Fatalf("resume moved destination: got %s want %s", got, destination)
			}
			if got := gitOutput(t, f.root, "rev-parse", "refs/bench/green/main"); got != published {
				t.Fatalf("resume moved project-green: got %s want %s", got, published)
			}
		})
	}
	if err := os.Remove(filepath.Join(f.creation.Path, "private", "output")); err != nil {
		t.Fatal(err)
	}
	if code, stdout, stderr := land("--resume", published, "--request", request, "--base", f.base, "--source-tip", f.tip, "--spec", "x", f.creation.Path); code != 0 || !strings.Contains(stdout, "worktree=released,census=0}") || stderr != "" {
		t.Fatalf("active resume = (%d, %q, %q)", code, stdout, stderr)
	}
	repo, _, err := cleanupIdentity(f.root, f.creation.Path)
	if err != nil {
		t.Fatal(err)
	}
	receipt, found, err := intent.CleanupReceiptFor(f.root, repo, releaseOperation, f.creation.Path, intent.RequestDigest(request))
	if err != nil || !found || receipt.Branch != f.creation.Assignment.Branch || receipt.BranchOID != f.tip {
		t.Fatalf("terminal receipt = %#v, found=%t error=%v", receipt, found, err)
	}
	checkout := gitOutput(t, f.root, "status", "--porcelain=v1", "--untracked-files=all")
	for _, tc := range []struct {
		name, want string
		mutate     func(*intent.CleanupReceipt)
	}{
		{name: "wrong-branch", want: "refused{detail=missing-terminal-receipt}\n", mutate: func(receipt *intent.CleanupReceipt) {
			receipt.Branch = intent.AssignmentBranchRef(strings.Repeat("a", 32), strings.Repeat("b", 32))
		}},
		{name: "wrong-source", want: "refused{detail=terminal receipt source tip mismatch,observed=" + f.tip + ",wanted=" + f.base + "}\n", mutate: func(receipt *intent.CleanupReceipt) {
			receipt.BranchOID = f.base
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			forged := receipt
			tc.mutate(&forged)
			if err := intent.PutCleanupReceipt(f.root, forged); err != nil {
				t.Fatal(err)
			}
			code, stdout, stderr := land("--resume", published, "--request", request, "--base", f.base, "--source-tip", f.tip, "--spec", "x", f.creation.Path)
			if got := gitOutput(t, f.root, "rev-parse", "main"); got != destination {
				t.Fatalf("mismatched receipt moved destination: got %s want %s", got, destination)
			}
			if got := gitOutput(t, f.root, "rev-parse", "refs/bench/green/main"); got != published {
				t.Fatalf("mismatched receipt moved project-green: got %s want %s", got, published)
			}
			if got := gitOutput(t, f.root, "status", "--porcelain=v1", "--untracked-files=all"); got != checkout {
				t.Fatalf("mismatched receipt changed checkout: got %q want %q", got, checkout)
			}
			if code != 1 || stdout != tc.want || stderr != "" {
				t.Fatalf("mismatched receipt resume = (%d, %q, %q)", code, stdout, stderr)
			}
		})
	}
}

func TestResumeLandCommandReconcilesAnUnreconciledPublishedCheckout(t *testing.T) {
	t.Parallel()
	request := "resume-reconcile"
	f := publicLandingFixture(t, request, "", "")
	working := defaultJoins()
	broken := working
	broken.reconcileLanding = func(joins, string, string, string, string) error {
		return errors.New("injected reconciliation interruption")
	}
	if r := runVerb(t, verbLand, f.callWith(broken, landArgs(request, f.base, f.tip, f.creation.Path)...)); r.exit != 3 || !strings.Contains(r.stdout, "worktree=incomplete:reconcile") {
		t.Fatalf("interrupted landing = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	published := gitOutput(t, f.root, "rev-parse", "main")
	args := []string{"--resume", published, "--request", request, "--base", f.base, "--source-tip", f.tip, "--spec", "x", f.creation.Path}
	if r := runVerb(t, verbLand, f.callWith(working, args...)); r.exit != 0 || !strings.Contains(r.stdout, "worktree=released,census=0}") || len(r.stderr) != 0 {
		t.Fatalf("resume reconciliation = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	if got := gitOutput(t, f.root, "rev-parse", "HEAD"); got != published {
		t.Fatalf("destination checkout = %s, want %s", got, published)
	}
	if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
		t.Fatalf("resume reran gate: tally=%q error=%v", got, err)
	}
}

func TestResumeLandCommandAcceptsSpecSlugAndPath(t *testing.T) {
	for _, specArg := range []string{"x", "./specs/x/spec.md"} {
		t.Run(specArg, func(t *testing.T) {
			request := "resume-spec-form-" + strings.ReplaceAll(specArg, "/", "-")
			f := markerLandingFixture(t, request, true)
			chdir(t, f.root)
			_, published := interruptLandingAtMarker(t, f, landArgs(request, f.base, f.tip, f.creation.Path)...)
			args := []string{"--resume", published, "--request", request, "--base", f.base, "--source-tip", f.tip, "--spec", specArg, f.creation.Path}
			if r := runVerb(t, verbLand, f.call(args...)); r.exit != 0 || !strings.Contains(r.stdout, "worktree=released,census=0}") || len(r.stderr) != 0 {
				t.Fatalf("resume with spec %q = (%d, %q, %q)", specArg, r.exit, r.stdout, r.stderr)
			}
			if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
				t.Fatalf("resume reran gate: tally=%q error=%v", got, err)
			}
		})
	}
}

func TestResumeLandCommandCompletesAnInterruptedMarker(t *testing.T) {
	t.Parallel()
	request := "resume-marker"
	f := markerLandingFixture(t, request, true)
	_, published := interruptLandingAtMarker(t, f, landArgs(request, f.base, f.tip, f.creation.Path)...)
	args := []string{"--resume", published, "--request", request, "--base", f.base, "--source-tip", f.tip, "--spec", "x", f.creation.Path}
	if r := runVerb(t, verbLand, f.call(args...)); r.exit != 0 || !strings.Contains(r.stdout, "worktree=released,census=0}") || len(r.stderr) != 0 {
		t.Fatalf("resume marker = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	if got := gitOutput(t, f.root, "rev-parse", "refs/bench/green/main"); got != published {
		t.Fatalf("project-green = %s, want %s", got, published)
	}
	if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
		t.Fatalf("resume reran gate: tally=%q error=%v", got, err)
	}
	markProof(t, "landing/journey/interrupted-resume")
}

func TestResumeLandCommandPreauthenticatesCompletedRequestAndPath(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		args func(string, string, string, string, string) []string
	}{
		{name: "wrong-request", args: func(published, base, tip, request, path string) []string {
			return []string{"--resume", published, "--request", request + "-forged", "--base", base, "--source-tip", tip, "--spec", "x", path}
		}},
		{name: "wrong-path", args: func(published, base, tip, request, path string) []string {
			return []string{"--resume", published, "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", filepath.Join(filepath.Dir(path), "forged")}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := "resume-preauth-" + tc.name
			f := publicLandingFixture(t, request, "", "")
			if r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...)); r.exit != 0 {
				t.Fatalf("landing = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
			}
			published := gitOutput(t, f.root, "rev-parse", "main")
			gitRun(t, f.root, "update-ref", "-d", "refs/bench/green/main")
			if r := runVerb(t, verbLand, f.call(tc.args(published, f.base, f.tip, request, f.creation.Path)...)); r.exit != 1 || !strings.Contains(r.stdout, "missing-terminal-receipt") || len(r.stderr) != 0 {
				t.Fatalf("preauthentication refusal = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
			}
			if descendant(t, "git", "-C", f.root, "show-ref", "--verify", "--quiet", "refs/bench/green/main").Run() == nil {
				t.Fatal("forged completed resume recreated project-green marker")
			}
		})
	}
}

// TestResumeLandCommandRepeatsTheCensusCount is EC22. The incomplete landing keeps
// the record file, because its release step never ran, so the resume reads the same
// count and states it again. A drop before the landed record would lose the evidence.
func TestResumeLandCommandRepeatsTheCensusCount(t *testing.T) {
	t.Parallel()
	request := "census-resume-count"
	f := publicLandingFixture(t, request, "private/output", "dist/")
	recordRawCalls(t, f.home, f.root, f.creation.Path, 2)
	r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
	if r.exit != 3 || !strings.HasSuffix(r.stdout, ",census=2}\n") {
		t.Fatalf("incomplete landing = (%d, %q, %q), want exit 3 and census=2", r.exit, r.stdout, r.stderr)
	}
	if !strings.Contains(r.stderr, "census heads{sed=2}\n") {
		t.Fatalf("incomplete landing stderr = %q, want the head breakdown", r.stderr)
	}
	if _, err := os.Stat(censusRecordPath(f.home, f.root, f.creation.Assignment.ID)); err != nil {
		t.Fatalf("the incomplete landing dropped the census records: %v", err)
	}
	published := gitOutput(t, f.root, "rev-parse", "main")
	args := []string{"--resume", published, "--request", request, "--base", f.base, "--source-tip", f.tip, "--spec", "x", f.creation.Path}
	r = runVerb(t, verbLand, f.call(args...))
	if r.exit != 3 || !strings.HasSuffix(r.stdout, ",census=2}\n") {
		t.Fatalf("resumed landing = (%d, %q, %q), want exit 3 and census=2 again", r.exit, r.stdout, r.stderr)
	}
	if !strings.Contains(r.stderr, "census heads{sed=2}\n") {
		t.Fatalf("resumed landing stderr = %q, want the head breakdown again", r.stderr)
	}
}
