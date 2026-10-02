// Reauthorization and identity-expansion tests for the landing command: unknown requests and abbreviated identities.
package worktree

import (
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
	f := publicLandingFixture(t, request, "", "")
	working := defaultJoins()
	broken := working
	broken.advanceLandingMarker = func(context.Context, string, string, string, string) error {
		return errors.New("injected marker interruption")
	}
	if r := runVerb(t, verbLand, f.callWith(broken, landArgs(request, f.base, f.tip, f.creation.Path)...)); r.exit != 3 {
		t.Fatalf("interrupted landing = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	published := gitOutput(t, f.root, "rev-parse", "main")
	args := []string{"--resume", published, "--request", "unknown-request", "--base", f.base, "--source-tip", f.tip, "--spec", "x", f.creation.Path}
	wantNext := laterProofsSkipped + "; bench worktree reauthorize --assignment " + f.creation.Assignment.ID + " --request <new-request> --base '" + f.base + "' --source-tip '" + f.tip + "' '" + f.creation.Path + "'"
	want := "refused{detail=request token matches no assignment,observed=assignment:" + f.creation.Assignment.ID + ",next=" + wantNext + "}\n"
	if r := runVerb(t, verbLand, f.callWith(working, args...)); r.exit != 1 || r.stdout != want || len(r.stderr) != 0 {
		t.Fatalf("unknown-request resume = (%d, %q, %q), want exit 1 and %q", r.exit, r.stdout, r.stderr, want)
	}
}

// TestLandCommandUnknownRequestNamesReauthorizeRecovery is LR02: one active assignment
// owns the target, so the request refusal names the exact command that repairs it.
func TestLandCommandUnknownRequestNamesReauthorizeRecovery(t *testing.T) {
	t.Parallel()
	request := "land-reauthorize-recovery"
	f := publicLandingFixture(t, request, "", "")
	r := runVerb(t, verbLand, f.call(landArgs("unknown-request", f.base, f.tip, f.creation.Path)...))
	wantNext := laterProofsSkipped + "; bench worktree reauthorize --assignment " + f.creation.Assignment.ID + " --request <new-request> --base '" + f.base + "' --source-tip '" + f.tip + "' '" + f.creation.Path + "'"
	want := "refused{detail=request token matches no assignment,observed=assignment:" + f.creation.Assignment.ID + ",next=" + wantNext + "}\n"
	if r.exit != 1 || r.stdout != want {
		t.Fatalf("unknown-request land = (%d, %q, %q), want exit 1 and %q", r.exit, r.stdout, r.stderr, want)
	}
}

func TestLandCommandReauthorizeRecoveryExpandsAbbreviatedIdentityInputs(t *testing.T) {
	t.Parallel()
	request := "reauthorize-full-identities"
	f := publicLandingFixture(t, request, "", "")
	r := runVerb(t, verbLand, f.call(landArgs("unknown-request", f.base[:12], f.tip[:12], f.creation.Path)...))
	wantNext := "next=" + laterProofsSkipped + "; bench worktree reauthorize --assignment " + f.creation.Assignment.ID + " --request <new-request> --base '" + f.base + "' --source-tip '" + f.tip + "' '" + f.creation.Path + "'}\n"
	if r.exit != 1 || !strings.HasSuffix(r.stdout, wantNext) {
		t.Fatalf("abbreviated-identity recovery = (%d, %q, %q), want suffix %q with the expanded identities", r.exit, r.stdout, r.stderr, wantNext)
	}
}

func TestLandCommandReauthorizeRecoveryPointsThroughUnsafePath(t *testing.T) {
	t.Parallel()
	request := "reauthorize-unsafe-path"
	home := filepath.Join(t.TempDir(), "bench\n\x1bhome")
	f := publicLandingFixtureAtHome(t, request, "", "", home)
	r := runVerb(t, verbLand, repoHome{f.root, home}.call(landArgs("unknown-request", f.base, f.tip, f.creation.Path)...))
	wantNext := "next=" + laterProofsSkipped + "; bench worktree exec " + f.creation.Assignment.ID + " -- bench worktree reauthorize --assignment " + f.creation.Assignment.ID + " --request <new-request> --base '" + f.base + "' --source-tip '" + f.tip + "' .}\n"
	unsafe := strings.ContainsRune(r.stdout, '\x1b') || strings.Count(r.stdout, "\n") != 1
	if r.exit != 1 || unsafe || !strings.HasSuffix(r.stdout, wantNext) {
		t.Fatalf("unsafe-path recovery = (%d, %q, %q), want one safe record ending %q", r.exit, r.stdout, r.stderr, wantNext)
	}
	// LRS21: a landing-preflight face at the same unsafe path takes the same pointer form,
	// so no route quotes a path the operator cannot paste.
	mustWrite(t, filepath.Join(f.creation.Path, "scratch"), []byte("scratch\n"), 0o600)
	r = runVerb(t, verbLand, repoHome{f.root, home}.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
	wantSource := "; then bench worktree exec " + f.creation.Assignment.ID + " -- bench worktree land --request '" +
		request + "' --base '" + f.base + "' --source-tip '" + f.tip + "' --spec 'x' -m <message> .}\n"
	unsafe = strings.ContainsRune(r.stdout, '\x1b') || strings.Count(r.stdout, "\n") != 1
	if r.exit != 1 || unsafe || !strings.HasSuffix(r.stdout, wantSource) {
		t.Fatalf("unsafe-path source refusal = (%d, %q, %q), want one safe record ending %q", r.exit, r.stdout, r.stderr, wantSource)
	}
}

func TestLandCommandStoredRequestDigestCannotAuthenticate(t *testing.T) {
	t.Parallel()
	request := "stored-digest-is-not-a-token"
	f := publicLandingFixture(t, request, "", "")
	gitRun(t, f.root, "update-ref", "refs/bench/green/main", f.base)
	beforeDestination := gitOutput(t, f.root, "rev-parse", "refs/heads/main")
	beforeSource := gitOutput(t, f.root, "rev-parse", f.creation.Assignment.Branch)
	beforeMarker := gitOutput(t, f.root, "rev-parse", "refs/bench/green/main")
	r := runVerb(t, verbLand, f.call(landArgs(f.creation.Assignment.Request, f.base, f.tip, f.creation.Path)...))
	if r.exit != 1 || !strings.HasPrefix(r.stdout, "refused{detail=request token matches no assignment") {
		t.Fatalf("stored-digest land = (%d, %q, %q), want refusal", r.exit, r.stdout, r.stderr)
	}
	if got := gitOutput(t, f.root, "rev-parse", "refs/heads/main"); got != beforeDestination {
		t.Fatalf("stored digest published destination: got %s want %s", got, beforeDestination)
	}
	if got := gitOutput(t, f.root, "rev-parse", f.creation.Assignment.Branch); got != beforeSource {
		t.Fatalf("stored digest moved source branch: got %s want %s", got, beforeSource)
	}
	if got := gitOutput(t, f.root, "rev-parse", "refs/bench/green/main"); got != beforeMarker {
		t.Fatalf("stored digest moved project-green marker: got %s want %s", got, beforeMarker)
	}
}

func TestLandCommandAuthenticatesDigestShapedRequestToken(t *testing.T) {
	t.Parallel()
	request := strings.Repeat("a", 64)
	f := publicLandingFixture(t, request, "", "")
	j := stubLandJoins(f.base, f.tip)
	j.releaseLandingAssignment = func(joins, string, string, []string, io.Writer, io.Writer) int { return 0 }

	r := runVerb(t, verbLand, f.callWith(j, landArgs(request, f.base, f.tip, f.creation.Path)...))
	if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released,census=0}") {
		t.Fatalf("digest-shaped request land = (%d, %q, %q), want successful authentication", r.exit, r.stdout, r.stderr)
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
	r := runVerb(t, verbLand, repoHome{root, home}.call(landArgs("unknown-request", base, base, root)...))
	want := "refused{detail=request token matches no assignment,next=" + laterProofsSkipped + "; bench worktree list}\n"
	if r.exit != 1 || r.stdout != want || strings.Contains(r.stdout, "reauthorize") {
		t.Fatalf("assignment-free land = (%d, %q, %q), want exit 1 and %q", r.exit, r.stdout, r.stderr, want)
	}
}

// TestLandCommandUnknownRequestWithAmbiguousAssignmentsNamesTheListing is the other half
// of LR03: two active assignments at the target name no single id either.
func TestLandCommandUnknownRequestWithAmbiguousAssignmentsNamesTheListing(t *testing.T) {
	t.Parallel()
	request := "ambiguous-reauthorize-recovery"
	f := publicLandingFixture(t, request, "", "")
	second := f.creation.Assignment
	second.ID = strings.Repeat("f", 32)
	second.Request = intent.RequestDigest("second-request")
	second.Label = "second assignment"
	second.Branch = intent.AssignmentBranchRef(second.OwnerID, second.ID)
	if err := intent.PutAssignment(f.root, second); err != nil {
		t.Fatal(err)
	}
	r := runVerb(t, verbLand, f.call(landArgs("unknown-request", f.base, f.tip, f.creation.Path)...))
	want := "refused{detail=request token matches no assignment,next=" + laterProofsSkipped + "; bench worktree list}\n"
	if r.exit != 1 || r.stdout != want || strings.Contains(r.stdout, "reauthorize") {
		t.Fatalf("ambiguous-assignment land = (%d, %q, %q), want exit 1 and %q", r.exit, r.stdout, r.stderr, want)
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
		r := runVerb(t, verbLand, repoHome{root, home}.callWith(j, landArgs(request, base, abbreviated, creation.Path)...))
		if r.exit != 0 || !strings.Contains(r.stdout, "source_tip="+tip+",") || !strings.Contains(r.stdout, "worktree=released,census=0}") {
			t.Fatalf("abbreviated source tip %q = (%d, %q, %q), want released with the full tip", abbreviated, r.exit, r.stdout, r.stderr)
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
	r := runVerb(t, verbLand, repoHome{root, home}.callWith(j, landArgs(request, base[:12], tip, creation.Path)...))
	if r.exit != 0 || authorized != base || !strings.Contains(r.stdout, "worktree=released,census=0}") {
		t.Fatalf("abbreviated base = (%d, authorized=%q, %q, %q), want the full base authorized and released", r.exit, authorized, r.stdout, r.stderr)
	}
}

func TestResumeLandCommandExpandsAbbreviatedIdentities(t *testing.T) {
	t.Parallel()
	request := "resume-abbreviated-published"
	f := publicLandingFixture(t, request, "private/output", "dist/")
	if r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...)); r.exit != 3 || !strings.Contains(r.stdout, "worktree=incomplete:release") {
		t.Fatalf("interrupted landing = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	published := gitOutput(t, f.root, "rev-parse", "main")
	args := []string{"--resume", published[:12], "--request", request, "--base", f.base[:12], "--source-tip", f.tip[:12], "--spec", "x", f.creation.Path}
	r := runVerb(t, verbLandResume, f.call(args...))
	if r.exit != 3 || !strings.Contains(r.stdout, "source_base="+f.base+",source_tip="+f.tip+",destination_base="+f.base+",published_commit="+published+",") || !strings.Contains(r.stdout, "worktree=incomplete:release") {
		t.Fatalf("abbreviated resume = (%d, %q, %q), want the published landing resumed under its full identities", r.exit, r.stdout, r.stderr)
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
		r := runVerb(t, verbLand, repoHome{root, home}.call(landArgs(request, base, observed, creation.Path)...))
		// Drift is not abbreviation, so the refusal names both tips and routes the caller's
		// own command at the tip the worktree holds, whatever value the caller named.
		want := "refused{detail=worktree source tip mismatch,observed=" + observed + ",wanted=" + tip +
			",next=" + landingRerun(request, base, tip, "x", creation.Path, creation.Assignment.ID) + "}\n"
		if r.exit != 1 || r.stdout != want || len(r.stderr) != 0 || strings.Contains(r.stdout, "abbreviated") {
			t.Fatalf("source tip drift %q = (%d, %q, %q), want (1, %q, empty)", observed, r.exit, r.stdout, r.stderr, want)
		}
	}
}
