// Resume tests for an interrupted landing: identity binding, marker completion, and destination-state refusals.
package worktree

import (
	"bytes"
	"context"
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
	root, creation, base, tip, _, home := publicLandingFixture(t, request, "private/output", "dist/")
	var stdout, stderr bytes.Buffer
	if code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr); code != 3 {
		t.Fatalf("first incomplete exit = %d, want 3; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	published := gitOutput(t, root, "rev-parse", "main")
	stdout.Reset()
	stderr.Reset()
	args := []string{"--resume", published, "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", creation.Path}
	code := LandCommand(root, home, "", args, &stdout, &stderr)
	if code != 3 || !strings.Contains(stdout.String(), "worktree=incomplete:release") {
		t.Fatalf("resume incomplete = (%d, %q, %q), want exit 3", code, stdout.String(), stderr.String())
	}
}

func TestLandCommandIncompleteNextUsesAssignmentPointerForUnsafePath(t *testing.T) {
	t.Parallel()
	request := "incomplete-unsafe-path"
	home := filepath.Join(t.TempDir(), "bench\n\x1bhome")
	root, creation, base, tip, _ := publicLandingFixtureAtHome(t, request, "private/output", "dist/", home)
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	wantNext := "next=bench worktree exec " + creation.Assignment.ID + " -- bench worktree land --resume '"
	unsafe := strings.ContainsRune(stdout.String(), '\x1b') || strings.Count(stdout.String(), "\n") != 1
	if code != 3 || unsafe || !strings.Contains(stdout.String(), wantNext) || !strings.Contains(stdout.String(), " --spec 'x' .,census=0}") {
		t.Fatalf("unsafe-path incomplete = (%d, %q, %q), want one safe pointer record containing %q", code, stdout.String(), stderr.String(), wantNext)
	}
}

func TestLandCommandPublicResumeCompletesPublishedReleaseWithoutRepublishing(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	request := "public-land-resume"
	root, creation, base, tip, tally, _ := publicLandingFixture(t, request, "private/output", "dist/")
	land := func(args ...string) (int, string, string) {
		var stdout, stderr bytes.Buffer
		cmd := descendant(t, binary, append([]string{"worktree", "land"}, args...)...)
		cmd.Dir, cmd.Stdout, cmd.Stderr = root, &stdout, &stderr
		return exitCode(cmd.Run()), stdout.String(), stderr.String()
	}
	code, stdout, stderr := land("--request", request, "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land reviewed source", creation.Path)
	if code != 3 || !strings.Contains(stdout, "source_base="+base) || !strings.Contains(stdout, "worktree=incomplete:release") || !strings.Contains(stderr, "worktree retained (ignored)") {
		t.Fatalf("interrupted landing = (%d, %q, %q)", code, stdout, stderr)
	}
	published := gitOutput(t, root, "rev-parse", "main")
	if err := os.Remove(filepath.Join(creation.Path, "private", "output")); err != nil {
		t.Fatal(err)
	}
	commitInWorktree(t, root, "destination-after-publication", "forward\n", "destination movement")
	// LF4: the destination declares Go build inputs only after publication, so the resume
	// faces a freshness proof it could not pass. This fixture's executable has no seal
	// these sources could match. It is committed because an untracked file would trip the
	// resume's own untracked-collision proof and hide the exemption behind another refusal.
	commitLandingBuildInputs(t, root, "build_script=scripts/go-build.sh\n")
	destination := gitOutput(t, root, "rev-parse", "main")
	code, stdout, stderr = land("--resume", published, "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", creation.Path)
	if code != 0 || !strings.Contains(stdout, "source_base="+base) || !strings.Contains(stdout, "worktree=released,census=0}") || stderr != "" {
		t.Fatalf("resume = (%d, %q, %q)", code, stdout, stderr)
	}
	if got := gitOutput(t, root, "rev-parse", "main"); got != destination {
		t.Fatalf("resume moved destination backward: main=%s destination=%s", got, destination)
	}
	if got, err := os.ReadFile(tally); err != nil || string(got) != "g" {
		t.Fatalf("resume reran gate: tally=%q error=%v", got, err)
	}
	code, stdout, stderr = land("--resume", published, "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", creation.Path)
	if code != 0 || !strings.Contains(stdout, "source_base="+base) || !strings.Contains(stdout, "worktree=already-complete,census=0}") || stderr != "" {
		t.Fatalf("completed resume = (%d, %q, %q)", code, stdout, stderr)
	}
}

func TestResumeLandCommandPublicBindsPublishedLandingIdentity(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	request := "public-resume-identity"
	root, creation, base, tip, _, _ := publicLandingFixture(t, request, "private/output", "dist/")
	land := func(args ...string) (int, string, string) {
		var stdout, stderr bytes.Buffer
		cmd := descendant(t, binary, append([]string{"worktree", "land"}, args...)...)
		cmd.Dir, cmd.Stdout, cmd.Stderr = root, &stdout, &stderr
		return exitCode(cmd.Run()), stdout.String(), stderr.String()
	}

	if code, stdout, stderr := land("--request", request, "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land reviewed source", creation.Path); code != 3 || !strings.Contains(stdout, "worktree=incomplete:release") || stderr == "" {
		t.Fatalf("interrupted landing = (%d, %q, %q)", code, stdout, stderr)
	}
	published := gitOutput(t, root, "rev-parse", "main")
	commitInWorktree(t, root, "destination-after-publication", "forward\n", "destination movement")
	destination := gitOutput(t, root, "rev-parse", "main")
	for _, tc := range []struct {
		name, published, source string
	}{
		{name: "wrong-published", published: destination, source: tip},
		{name: "wrong-source", published: published, source: base},
	} {
		t.Run(tc.name, func(t *testing.T) {
			code, stdout, stderr := land("--resume", tc.published, "--request", request, "--base", base, "--source-tip", tc.source, "--spec", "x", creation.Path)
			if code != 1 || !strings.HasPrefix(stdout, "refused{detail=") || stderr != "" {
				t.Fatalf("resume refusal = (%d, %q, %q)", code, stdout, stderr)
			}
			if got := gitOutput(t, root, "rev-parse", "main"); got != destination {
				t.Fatalf("resume moved destination: got %s want %s", got, destination)
			}
			if got := gitOutput(t, root, "rev-parse", "refs/bench/green/main"); got != published {
				t.Fatalf("resume moved project-green: got %s want %s", got, published)
			}
		})
	}
	if err := os.Remove(filepath.Join(creation.Path, "private", "output")); err != nil {
		t.Fatal(err)
	}
	if code, stdout, stderr := land("--resume", published, "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", creation.Path); code != 0 || !strings.Contains(stdout, "worktree=released,census=0}") || stderr != "" {
		t.Fatalf("active resume = (%d, %q, %q)", code, stdout, stderr)
	}
	repo, _, err := cleanupIdentity(root, creation.Path)
	if err != nil {
		t.Fatal(err)
	}
	receipt, found, err := intent.CleanupReceiptFor(root, repo, releaseOperation, creation.Path, intent.RequestDigest(request))
	if err != nil || !found || receipt.Branch != creation.Assignment.Branch || receipt.BranchOID != tip {
		t.Fatalf("terminal receipt = %#v, found=%t error=%v", receipt, found, err)
	}
	checkout := gitOutput(t, root, "status", "--porcelain=v1", "--untracked-files=all")
	for _, tc := range []struct {
		name, want string
		mutate     func(*intent.CleanupReceipt)
	}{
		{name: "wrong-branch", want: "refused{detail=missing-terminal-receipt}\n", mutate: func(receipt *intent.CleanupReceipt) {
			receipt.Branch = intent.AssignmentBranchRef(strings.Repeat("a", 32), strings.Repeat("b", 32))
		}},
		{name: "wrong-source", want: "refused{detail=terminal receipt source tip mismatch,observed=" + tip + ",wanted=" + base + "}\n", mutate: func(receipt *intent.CleanupReceipt) {
			receipt.BranchOID = base
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			forged := receipt
			tc.mutate(&forged)
			if err := intent.PutCleanupReceipt(root, forged); err != nil {
				t.Fatal(err)
			}
			code, stdout, stderr := land("--resume", published, "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", creation.Path)
			if got := gitOutput(t, root, "rev-parse", "main"); got != destination {
				t.Fatalf("mismatched receipt moved destination: got %s want %s", got, destination)
			}
			if got := gitOutput(t, root, "rev-parse", "refs/bench/green/main"); got != published {
				t.Fatalf("mismatched receipt moved project-green: got %s want %s", got, published)
			}
			if got := gitOutput(t, root, "status", "--porcelain=v1", "--untracked-files=all"); got != checkout {
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
	root, creation, base, tip, tally, home := publicLandingFixture(t, request, "", "")
	working := defaultJoins()
	broken := working
	broken.reconcileLanding = func(joins, string, string, string, string) error {
		return errors.New("injected reconciliation interruption")
	}
	var stdout, stderr bytes.Buffer
	if code := landWith(broken, root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr); code != 3 || !strings.Contains(stdout.String(), "worktree=incomplete:reconcile") {
		t.Fatalf("interrupted landing = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
	published := gitOutput(t, root, "rev-parse", "main")
	stdout.Reset()
	stderr.Reset()
	args := []string{"--resume", published, "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", creation.Path}
	if code := landWith(working, root, home, "", args, &stdout, &stderr); code != 0 || !strings.Contains(stdout.String(), "worktree=released,census=0}") || stderr.Len() != 0 {
		t.Fatalf("resume reconciliation = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
	if got := gitOutput(t, root, "rev-parse", "HEAD"); got != published {
		t.Fatalf("destination checkout = %s, want %s", got, published)
	}
	if got, err := os.ReadFile(tally); err != nil || string(got) != "g" {
		t.Fatalf("resume reran gate: tally=%q error=%v", got, err)
	}
}

func TestResumeLandCommandAcceptsSpecSlugAndPath(t *testing.T) {
	for _, specArg := range []string{"x", "./specs/x/spec.md"} {
		t.Run(specArg, func(t *testing.T) {
			request := "resume-spec-form-" + strings.ReplaceAll(specArg, "/", "-")
			root, creation, base, tip, tally, home := publicLandingFixture(t, request, "", "")
			chdir(t, root)
			working := defaultJoins()
			broken := working
			broken.advanceLandingMarker = func(context.Context, string, string, string, string) error {
				return errors.New("injected marker interruption")
			}

			var stdout, stderr bytes.Buffer
			if code := landWith(broken, root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr); code != 3 || !strings.Contains(stdout.String(), "worktree=incomplete:marker") {
				t.Fatalf("interrupted landing = (%d, %q, %q)", code, stdout.String(), stderr.String())
			}
			published := gitOutput(t, root, "rev-parse", "main")
			stdout.Reset()
			stderr.Reset()
			args := []string{"--resume", published, "--request", request, "--base", base, "--source-tip", tip, "--spec", specArg, creation.Path}
			if code := landWith(working, root, home, "", args, &stdout, &stderr); code != 0 || !strings.Contains(stdout.String(), "worktree=released,census=0}") || stderr.Len() != 0 {
				t.Fatalf("resume with spec %q = (%d, %q, %q)", specArg, code, stdout.String(), stderr.String())
			}
			if got, err := os.ReadFile(tally); err != nil || string(got) != "g" {
				t.Fatalf("resume reran gate: tally=%q error=%v", got, err)
			}
		})
	}
}

func TestResumeLandCommandCompletesAnInterruptedMarker(t *testing.T) {
	t.Parallel()
	request := "resume-marker"
	root, creation, base, tip, tally, home := publicLandingFixture(t, request, "", "")
	working := defaultJoins()
	broken := working
	broken.advanceLandingMarker = func(context.Context, string, string, string, string) error {
		return errors.New("injected marker interruption")
	}
	var stdout, stderr bytes.Buffer
	if code := landWith(broken, root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr); code != 3 || !strings.Contains(stdout.String(), "worktree=incomplete:marker") {
		t.Fatalf("interrupted landing = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
	published := gitOutput(t, root, "rev-parse", "main")
	stdout.Reset()
	stderr.Reset()
	args := []string{"--resume", published, "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", creation.Path}
	if code := landWith(working, root, home, "", args, &stdout, &stderr); code != 0 || !strings.Contains(stdout.String(), "worktree=released,census=0}") || stderr.Len() != 0 {
		t.Fatalf("resume marker = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
	if got := gitOutput(t, root, "rev-parse", "refs/bench/green/main"); got != published {
		t.Fatalf("project-green = %s, want %s", got, published)
	}
	if got, err := os.ReadFile(tally); err != nil || string(got) != "g" {
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
			root, creation, base, tip, _, home := publicLandingFixture(t, request, "", "")
			var stdout, stderr bytes.Buffer
			if code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr); code != 0 {
				t.Fatalf("landing = (%d, %q, %q)", code, stdout.String(), stderr.String())
			}
			published := gitOutput(t, root, "rev-parse", "main")
			gitRun(t, root, "update-ref", "-d", "refs/bench/green/main")
			stdout.Reset()
			stderr.Reset()
			if code := LandCommand(root, home, "", tc.args(published, base, tip, request, creation.Path), &stdout, &stderr); code != 1 || !strings.Contains(stdout.String(), "missing-terminal-receipt") || stderr.Len() != 0 {
				t.Fatalf("preauthentication refusal = (%d, %q, %q)", code, stdout.String(), stderr.String())
			}
			if descendant(t, "git", "-C", root, "show-ref", "--verify", "--quiet", "refs/bench/green/main").Run() == nil {
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
	root, creation, base, tip, _, home := publicLandingFixture(t, request, "private/output", "dist/")
	recordRawCalls(t, home, root, creation.Path, 2)
	var stdout, stderr bytes.Buffer
	if code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr); code != 3 || !strings.HasSuffix(stdout.String(), ",census=2}\n") {
		t.Fatalf("incomplete landing = (%d, %q, %q), want exit 3 and census=2", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(censusRecordPath(home, root, creation.Assignment.ID)); err != nil {
		t.Fatalf("the incomplete landing dropped the census records: %v", err)
	}
	published := gitOutput(t, root, "rev-parse", "main")
	stdout.Reset()
	stderr.Reset()
	args := []string{"--resume", published, "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", creation.Path}
	if code := LandCommand(root, home, "", args, &stdout, &stderr); code != 3 || !strings.HasSuffix(stdout.String(), ",census=2}\n") {
		t.Fatalf("resumed landing = (%d, %q, %q), want exit 3 and census=2 again", code, stdout.String(), stderr.String())
	}
}
