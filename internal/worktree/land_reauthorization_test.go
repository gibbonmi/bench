// Reauthorization and identity-expansion tests for the landing command: unknown requests and abbreviated identities.
package worktree

import (
	"bytes"
	"context"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/intent"
)

func TestResumeLandCommandUnknownRequestNamesReauthorizeRecovery(t *testing.T) {
	t.Parallel()
	request := "resume-reauthorize-recovery"
	root, creation, base, tip, _, home := publicLandingFixture(t, request, "", "")
	working := defaultJoins()
	broken := working
	broken.advanceLandingMarker = func(context.Context, string, string, string, string) error {
		return errors.New("injected marker interruption")
	}
	var stdout, stderr bytes.Buffer
	if code := landWith(broken, root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr); code != 3 {
		t.Fatalf("interrupted landing = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
	published := gitOutput(t, root, "rev-parse", "main")
	stdout.Reset()
	stderr.Reset()
	args := []string{"--resume", published, "--request", "unknown-request", "--base", base, "--source-tip", tip, "--spec", "x", creation.Path}
	wantNext := "bench worktree reauthorize --assignment " + creation.Assignment.ID + " --request <new-request> --base '" + base + "' --source-tip '" + tip + "' '" + creation.Path + "'"
	want := "refused{detail=request token matches no assignment,observed=assignment:" + creation.Assignment.ID + ",next=" + wantNext + "}\n"
	if code := landWith(working, root, home, "", args, &stdout, &stderr); code != 1 || stdout.String() != want || stderr.Len() != 0 {
		t.Fatalf("unknown-request resume = (%d, %q, %q), want exit 1 and %q", code, stdout.String(), stderr.String(), want)
	}
}

// TestLandCommandUnknownRequestNamesReauthorizeRecovery is LR02: one active assignment
// owns the target, so the request refusal names the exact command that repairs it.
func TestLandCommandUnknownRequestNamesReauthorizeRecovery(t *testing.T) {
	t.Parallel()
	request := "land-reauthorize-recovery"
	root, creation, base, tip, _, home := publicLandingFixture(t, request, "", "")
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs("unknown-request", base, tip, creation.Path), &stdout, &stderr)
	wantNext := "bench worktree reauthorize --assignment " + creation.Assignment.ID + " --request <new-request> --base '" + base + "' --source-tip '" + tip + "' '" + creation.Path + "'"
	want := "refused{detail=request token matches no assignment,observed=assignment:" + creation.Assignment.ID + ",next=" + wantNext + "}\n"
	if code != 1 || stdout.String() != want {
		t.Fatalf("unknown-request land = (%d, %q, %q), want exit 1 and %q", code, stdout.String(), stderr.String(), want)
	}
}

func TestLandCommandReauthorizeRecoveryExpandsAbbreviatedIdentityInputs(t *testing.T) {
	t.Parallel()
	request := "reauthorize-full-identities"
	root, creation, base, tip, _, home := publicLandingFixture(t, request, "", "")
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs("unknown-request", base[:12], tip[:12], creation.Path), &stdout, &stderr)
	wantNext := "next=bench worktree reauthorize --assignment " + creation.Assignment.ID + " --request <new-request> --base '" + base + "' --source-tip '" + tip + "' '" + creation.Path + "'}\n"
	if code != 1 || !strings.HasSuffix(stdout.String(), wantNext) {
		t.Fatalf("abbreviated-identity recovery = (%d, %q, %q), want suffix %q with the expanded identities", code, stdout.String(), stderr.String(), wantNext)
	}
}

func TestLandCommandReauthorizeRecoveryPointsThroughUnsafePath(t *testing.T) {
	t.Parallel()
	request := "reauthorize-unsafe-path"
	home := filepath.Join(t.TempDir(), "bench\n\x1bhome")
	root, creation, base, tip, _ := publicLandingFixtureAtHome(t, request, "", "", home)
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs("unknown-request", base, tip, creation.Path), &stdout, &stderr)
	wantNext := "next=bench worktree exec " + creation.Assignment.ID + " -- bench worktree reauthorize --assignment " + creation.Assignment.ID + " --request <new-request> --base '" + base + "' --source-tip '" + tip + "' .}\n"
	unsafe := strings.ContainsRune(stdout.String(), '\x1b') || strings.Count(stdout.String(), "\n") != 1
	if code != 1 || unsafe || !strings.HasSuffix(stdout.String(), wantNext) {
		t.Fatalf("unsafe-path recovery = (%d, %q, %q), want one safe record ending %q", code, stdout.String(), stderr.String(), wantNext)
	}
}

func TestLandCommandStoredRequestDigestCannotAuthenticate(t *testing.T) {
	t.Parallel()
	request := "stored-digest-is-not-a-token"
	root, creation, base, tip, _, home := publicLandingFixture(t, request, "", "")
	gitRun(t, root, "update-ref", "refs/bench/green/main", base)
	beforeDestination := gitOutput(t, root, "rev-parse", "refs/heads/main")
	beforeSource := gitOutput(t, root, "rev-parse", creation.Assignment.Branch)
	beforeMarker := gitOutput(t, root, "rev-parse", "refs/bench/green/main")
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs(creation.Assignment.Request, base, tip, creation.Path), &stdout, &stderr)
	if code != 1 || !strings.HasPrefix(stdout.String(), "refused{detail=request token matches no assignment") {
		t.Fatalf("stored-digest land = (%d, %q, %q), want refusal", code, stdout.String(), stderr.String())
	}
	if got := gitOutput(t, root, "rev-parse", "refs/heads/main"); got != beforeDestination {
		t.Fatalf("stored digest published destination: got %s want %s", got, beforeDestination)
	}
	if got := gitOutput(t, root, "rev-parse", creation.Assignment.Branch); got != beforeSource {
		t.Fatalf("stored digest moved source branch: got %s want %s", got, beforeSource)
	}
	if got := gitOutput(t, root, "rev-parse", "refs/bench/green/main"); got != beforeMarker {
		t.Fatalf("stored digest moved project-green marker: got %s want %s", got, beforeMarker)
	}
}

func TestLandCommandAuthenticatesDigestShapedRequestToken(t *testing.T) {
	t.Parallel()
	request := strings.Repeat("a", 64)
	root, creation, base, tip, _, home := publicLandingFixture(t, request, "", "")
	j := stubLandJoins(base, tip)
	j.releaseLandingAssignment = func(joins, string, string, []string, io.Writer, io.Writer) int { return 0 }

	var stdout, stderr bytes.Buffer
	code := landWith(j, root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 || !strings.Contains(stdout.String(), "worktree=released,census=0}") {
		t.Fatalf("digest-shaped request land = (%d, %q, %q), want successful authentication", code, stdout.String(), stderr.String())
	}
}

// TestLandCommandUnknownRequestWithoutAssignmentNamesTheListing is half of LR03: with no
// active assignment at the target there is no id to reauthorize, so the route is the
// listing that names which assignment owns which tree.
func TestLandCommandUnknownRequestWithoutAssignmentNamesTheListing(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(t.TempDir(), "bench-home")
	base := gitOutput(t, root, "rev-parse", "HEAD")
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs("unknown-request", base, base, root), &stdout, &stderr)
	want := "refused{detail=request token matches no assignment,next=bench worktree list}\n"
	if code != 1 || stdout.String() != want || strings.Contains(stdout.String(), "reauthorize") {
		t.Fatalf("assignment-free land = (%d, %q, %q), want exit 1 and %q", code, stdout.String(), stderr.String(), want)
	}
}

// TestLandCommandUnknownRequestWithAmbiguousAssignmentsNamesTheListing is the other half
// of LR03: two active assignments at the target name no single id either.
func TestLandCommandUnknownRequestWithAmbiguousAssignmentsNamesTheListing(t *testing.T) {
	t.Parallel()
	request := "ambiguous-reauthorize-recovery"
	root, creation, base, tip, _, home := publicLandingFixture(t, request, "", "")
	second := creation.Assignment
	second.ID = strings.Repeat("f", 32)
	second.Request = intent.RequestDigest("second-request")
	second.Label = "second assignment"
	second.Branch = intent.AssignmentBranchRef(second.OwnerID, second.ID)
	if err := intent.PutAssignment(root, second); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs("unknown-request", base, tip, creation.Path), &stdout, &stderr)
	want := "refused{detail=request token matches no assignment,next=bench worktree list}\n"
	if code != 1 || stdout.String() != want || strings.Contains(stdout.String(), "reauthorize") {
		t.Fatalf("ambiguous-assignment land = (%d, %q, %q), want exit 1 and %q", code, stdout.String(), stderr.String(), want)
	}
}

// An abbreviated identity expands to the exact commit before any identity proof runs,
// so the proof that compares it sees the full value and the landing pins it.
func TestLandCommandExpandsAbbreviatedSourceTip(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(t.TempDir(), "bench-home")
	request := "landed-abbreviated-source-tip"
	creation := mustCreate(t, root, home, request, "abbreviated source tip")
	stageLandSpec(t, root, creation.Path)
	base := gitOutput(t, root, "rev-parse", "HEAD")
	commitInWorktree(t, creation.Path, "owned.txt", "owned\n", "owned")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	j := stubLandJoins(base, tip)
	j.releaseLandingAssignment = func(joins, string, string, []string, io.Writer, io.Writer) int { return 0 }
	for _, abbreviated := range []string{tip[:4], tip[:12], tip[:39], strings.ToUpper(tip[:12])} {
		var stdout, stderr bytes.Buffer
		code := landWith(j, root, home, "", landArgs(request, base, abbreviated, creation.Path), &stdout, &stderr)
		if code != 0 || !strings.Contains(stdout.String(), "source_tip="+tip+",") || !strings.Contains(stdout.String(), "worktree=released,census=0}") {
			t.Fatalf("abbreviated source tip %q = (%d, %q, %q), want released with the full tip", abbreviated, code, stdout.String(), stderr.String())
		}
	}
}

func TestLandCommandExpandsAbbreviatedBase(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(t.TempDir(), "bench-home")
	request := "landed-abbreviated-base"
	creation := mustCreate(t, root, home, request, "abbreviated base")
	stageLandSpec(t, root, creation.Path)
	base := gitOutput(t, root, "rev-parse", "HEAD")
	commitInWorktree(t, creation.Path, "owned.txt", "owned\n", "owned")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	j := stubLandJoins(base, tip)
	j.releaseLandingAssignment = func(joins, string, string, []string, io.Writer, io.Writer) int { return 0 }
	authorized := ""
	j.authorizeLandingSource = func(_, _ string, reviewBase string) (diff.SourceRange, error) {
		authorized = reviewBase
		return diff.SourceRange{Base: base, Tip: tip}, nil
	}
	var stdout, stderr bytes.Buffer
	code := landWith(j, root, home, "", landArgs(request, base[:12], tip, creation.Path), &stdout, &stderr)
	if code != 0 || authorized != base || !strings.Contains(stdout.String(), "worktree=released,census=0}") {
		t.Fatalf("abbreviated base = (%d, authorized=%q, %q, %q), want the full base authorized and released", code, authorized, stdout.String(), stderr.String())
	}
}

func TestResumeLandCommandExpandsAbbreviatedIdentities(t *testing.T) {
	t.Parallel()
	request := "resume-abbreviated-published"
	root, creation, base, tip, _, home := publicLandingFixture(t, request, "private/output", "dist/")
	var stdout, stderr bytes.Buffer
	if code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr); code != 3 || !strings.Contains(stdout.String(), "worktree=incomplete:release") {
		t.Fatalf("interrupted landing = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
	published := gitOutput(t, root, "rev-parse", "main")
	stdout.Reset()
	stderr.Reset()
	args := []string{"--resume", published[:12], "--request", request, "--base", base[:12], "--source-tip", tip[:12], "--spec", "x", creation.Path}
	code := ResumeLandCommand(root, home, args, &stdout, &stderr)
	if code != 3 || !strings.Contains(stdout.String(), "source_base="+base+",source_tip="+tip+",destination_base="+base+",published_commit="+published+",") || !strings.Contains(stdout.String(), "worktree=incomplete:release") {
		t.Fatalf("abbreviated resume = (%d, %q, %q), want the published landing resumed under its full identities", code, stdout.String(), stderr.String())
	}
}

func TestLandCommandDistinguishesSourceTipDriftFromAbbreviation(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(t.TempDir(), "bench-home")
	request := "land-source-tip-drift"
	creation := mustCreate(t, root, home, request, "source tip drift")
	stageLandSpec(t, root, creation.Path)
	base := gitOutput(t, root, "rev-parse", "HEAD")
	commitInWorktree(t, creation.Path, "owned.txt", "owned\n", "owned")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	for _, observed := range []string{base, tip[:3], "not-a-commit"} {
		var stdout, stderr bytes.Buffer
		code := LandCommand(root, home, "", landArgs(request, base, observed, creation.Path), &stdout, &stderr)
		want := "refused{detail=worktree source tip mismatch,observed=" + observed + ",wanted=" + tip + "}\n"
		if code != 1 || stdout.String() != want || stderr.Len() != 0 || strings.Contains(stdout.String(), "abbreviated") {
			t.Fatalf("source tip drift %q = (%d, %q, %q), want (1, %q, empty)", observed, code, stdout.String(), stderr.String(), want)
		}
	}
}
