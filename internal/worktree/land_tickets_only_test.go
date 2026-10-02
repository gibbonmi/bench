// Tickets-only spec landing tests: folder closure, an already-removed folder, and an interrupted close.
package worktree

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// WL8: a --spec naming a tickets-only folder closes that folder on the landing rather
// than refusing it as an unreadable staged spec.
func TestLandCommandTicketsOnlySpecClosesTheFolder(t *testing.T) {
	t.Parallel()
	request := "tickets-only-close"
	f := ticketsOnlyLandingFixture(t, request)
	r := runVerb(t, verbLand, f.call(ticketsOnlyLandArgs(request, f.base, f.tip, "t", f.creation.Path)...))
	published := gitOutput(t, f.root, "rev-parse", "main")
	if r.exit != 0 || !strings.Contains(r.stdout, "published_commit="+published+",") || !strings.HasSuffix(r.stdout, "worktree=released,census=0}\n") {
		t.Fatalf("tickets-only close = (%d, %q, %q), want exit 0 and a released landing", r.exit, r.stdout, r.stderr)
	}
	if descendant(t, "git", "-C", f.root, "cat-file", "-e", published+":specs/t").Run() == nil {
		t.Fatalf("published tree still carries specs/t")
	}
	if _, err := os.Stat(filepath.Join(f.root, "specs", "t")); !os.IsNotExist(err) {
		t.Fatalf("destination checkout still carries specs/t: %v", err)
	}
	if got := gitOutput(t, f.root, "show", published+":specs/x/spec.md"); strings.Contains(got, "Status: implemented") {
		t.Fatalf("tickets-only close transitioned a spec: %q", got)
	}
	if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
		t.Fatalf("gate tally = %q, %v", got, err)
	}
}

// Edge under WL8: the destination already removed the folder, so the close composes as
// a no-op and the landing still publishes and releases.
func TestLandCommandTicketsOnlySpecLandsWhenTheDestinationAlreadyRemovedTheFolder(t *testing.T) {
	t.Parallel()
	request := "tickets-only-already-removed"
	f := ticketsOnlyLandingFixture(t, request)
	gitRun(t, f.root, "rm", "-r", "-q", "specs/t")
	gitRun(t, f.root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "destination closed the folder")
	gitRun(t, f.root, "update-ref", "refs/bench/green/main", gitOutput(t, f.root, "rev-parse", "HEAD"))
	r := runVerb(t, verbLand, f.call(ticketsOnlyLandArgs(request, f.base, f.tip, "t", f.creation.Path)...))
	published := gitOutput(t, f.root, "rev-parse", "main")
	if r.exit != 0 || !strings.HasSuffix(r.stdout, "worktree=released,census=0}\n") {
		t.Fatalf("already-removed close = (%d, %q, %q), want exit 0 and a released landing", r.exit, r.stdout, r.stderr)
	}
	if descendant(t, "git", "-C", f.root, "cat-file", "-e", published+":specs/t").Run() == nil {
		t.Fatalf("published tree still carries specs/t")
	}
}

// Edge under WL8: a --spec naming a folder absent from the source is neither a staged
// spec.md nor a tickets-only folder, so it keeps the refusal it has today, which names
// the unreadable spec through the fence resolve.
func TestLandCommandAbsentSpecFolderKeepsTheUnreadableRefusal(t *testing.T) {
	t.Parallel()
	request := "tickets-only-absent"
	f := ticketsOnlyLandingFixture(t, request)
	r := runVerb(t, verbLand, f.call(ticketsOnlyLandArgs(request, f.base, f.tip, "absent", f.creation.Path)...))
	if r.exit != 1 || !strings.Contains(r.stdout, "reviewed source range or ownership fence is invalid: spec not found: no spec resolved for absent") {
		t.Fatalf("absent spec folder = (%d, %q, %q), want the unreadable staged-spec refusal", r.exit, r.stdout, r.stderr)
	}
	if _, err := os.Stat(f.tally); !os.IsNotExist(err) {
		t.Fatalf("unreadable-spec refusal ran the gate: %v", err)
	}
	// LRS3: the fence face routes its wrapped cause too, so the operator reads the
	// face's own repair with the caller's own re-run behind it.
	rerun := "bench worktree land --request '" + request + "' --base '" + f.base +
		"' --source-tip '" + f.tip + "' --spec 'absent' -m <message> '" + f.creation.Path + "'"
	next, printed := landingFaceNext(r.stdout, landingRefusalFaceByName(faceSourceNotFenced).detail)
	repair := landingRefusalFaceByName(faceSourceNotFenced).route(rerun)
	if !printed || next != repair {
		t.Fatalf("absent spec folder next = %q (printed=%t) in %q, want %q", next, printed, r.stdout, repair)
	}
}

// Edge under WL8: a --resume carrying the tickets-only slug authenticates the folder's
// absence from the published commit, never a spec.md transition the first run never made.
func TestResumeLandCommandTicketsOnlySpecCompletesAnInterruptedClose(t *testing.T) {
	t.Parallel()
	request := "tickets-only-resume"
	f := ticketsOnlyLandingFixture(t, request)
	working := defaultJoins()
	broken := working
	broken.advanceLandingMarker = func(context.Context, string, string, string, string) error {
		return errors.New("injected marker interruption")
	}
	r := runVerb(t, verbLand, f.callWith(broken, ticketsOnlyLandArgs(request, f.base, f.tip, "t", f.creation.Path)...))
	if r.exit != 3 || !strings.Contains(r.stdout, "worktree=incomplete:marker") {
		t.Fatalf("interrupted tickets-only landing = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	published := gitOutput(t, f.root, "rev-parse", "main")
	args := []string{"--resume", published, "--request", request, "--base", f.base, "--source-tip", f.tip, "--spec", "t", f.creation.Path}
	r = runVerb(t, verbLand, f.callWith(working, args...))
	if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released,census=0}") || len(r.stderr) != 0 {
		t.Fatalf("tickets-only resume = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	if got := gitOutput(t, f.root, "rev-parse", "refs/bench/green/main"); got != published {
		t.Fatalf("project-green = %s, want %s", got, published)
	}
	if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
		t.Fatalf("tickets-only resume reran the gate: tally=%q error=%v", got, err)
	}
}

// SR61: the interruption lands after the reconcile, so the destination checkout no longer
// carries specs/t when the resume classifies the landing. The classification must read the
// source commit's objects; a reader over the destination working tree finds no folder and
// authenticates a spec transition the first run never published.
func TestResumeLandCommandTicketsOnlyCloseSurvivesTheConsumedCheckout(t *testing.T) {
	t.Parallel()
	request := "tickets-only-resume-released"
	f := ticketsOnlyLandingFixture(t, request)
	working := defaultJoins()
	broken := working
	broken.releaseLandingAssignment = func(joins, string, string, []string, io.Writer, io.Writer) int { return 1 }
	r := runVerb(t, verbLand, f.callWith(broken, ticketsOnlyLandArgs(request, f.base, f.tip, "t", f.creation.Path)...))
	if r.exit != 3 || !strings.Contains(r.stdout, "worktree=incomplete:release") {
		t.Fatalf("interrupted release = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	if _, err := os.Stat(filepath.Join(f.root, "specs", "t")); !os.IsNotExist(err) {
		t.Fatalf("destination checkout still carries specs/t after the reconcile: %v", err)
	}
	published := gitOutput(t, f.root, "rev-parse", "main")
	args := []string{"--resume", published, "--request", request, "--base", f.base, "--source-tip", f.tip, "--spec", "t", f.creation.Path}
	r = runVerb(t, verbLand, f.callWith(working, args...))
	if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released,census=0}") || len(r.stderr) != 0 {
		t.Fatalf("resume after a consumed checkout = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
		t.Fatalf("resume reran the gate: tally=%q error=%v", got, err)
	}
}
