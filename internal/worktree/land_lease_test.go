// Landing lease tests: a landing holds one lease on its destination from composition
// through publish, and a second landing on that destination refuses before composition.
package worktree

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/testrepo"
)

// leasedMain is the destination ref of every landing fixture.
const leasedMain = "refs/heads/main"

// otherLandingID is the assignment of a landing outside the fixture that holds the lease.
const otherLandingID = "0123456789abcdef0123456789abcdef"

// holdLandingLease takes the destination's lease for another live landing, which this
// test process stands in for until the returned lease releases or the test ends.
func holdLandingLease(t *testing.T, root, assignment string) *commitrepo.Lease {
	t.Helper()
	lease, err := commitrepo.Store{Root: root}.AcquireLease(leasedMain, assignment)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(lease.Release)
	return lease
}

// barrierLanding is a landing fixture whose prospective gate signals started and then
// waits until a line reaches release.
type barrierLanding struct {
	landingFixture
	started, release string
}

func barrierLandingFixture(t *testing.T, request string) barrierLanding {
	t.Helper()
	barrier := t.TempDir()
	started, release := filepath.Join(barrier, "started"), filepath.Join(barrier, "release")
	for _, fifo := range []string{started, release} {
		mustNoError(t, syscall.Mkfifo(fifo, 0o600))
	}
	f := landingFixtureWithGateStep(t, request, "", "", filepath.Join(t.TempDir(), "bench-home"), true, func(*testrepo.GateFixture, string) string {
		return ": > " + sanitize.ShellQuote(started) + "\nIFS= read -r _ < " + sanitize.ShellQuote(release) + "\n"
	})
	return barrierLanding{landingFixture: f, started: started, release: release}
}

// during runs inside the gate once the gate reaches its barrier, and it releases the
// gate when it returns. The landing runs in the caller, so a landing that never reaches
// its gate leaves during unrun and the caller's receive empty.
func (f barrierLanding) during(inside func()) <-chan struct{} {
	done := make(chan struct{}, 1)
	go func() {
		if _, err := os.ReadFile(f.started); err == nil {
			inside()
		}
		done <- struct{}{}
		_ = os.WriteFile(f.release, []byte("go\n"), 0)
	}()
	return done
}

// TestLandHoldsTheDestinationLease is BW1 and BW6. The ledger names the landing's own
// assignment while its gate runs, and it names no lease once the landing published.
func TestLandHoldsTheDestinationLease(t *testing.T) {
	t.Parallel()
	request := "land-lease-holds"
	f := barrierLandingFixture(t, request)
	var atGate []intent.LandingLease
	var readErr error
	done := f.during(func() { atGate, readErr = commitmenttest.Leases(f.root) })
	r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
	select {
	case <-done:
	default:
		t.Fatalf("landing finished without reaching its gate barrier: (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	mustNoError(t, readErr)
	if r.exit != 0 {
		t.Fatalf("landing = (%d, %q, %q), want a publication", r.exit, r.stdout, r.stderr)
	}
	if want := (intent.LandingLease{Destination: leasedMain, Assignment: f.creation.Assignment.ID}); len(atGate) != 1 || atGate[0] != want {
		t.Fatalf("leases during the gate = %+v, want only %+v", atGate, want)
	}
	after, err := commitmenttest.Leases(f.root)
	mustNoError(t, err)
	if len(after) != 0 {
		t.Fatalf("leases after the publication = %+v, want none", after)
	}
}

// TestSecondLandingMeetsTheLease is BW2, BW3, BW4, and BW5. A second landing on the same
// destination, during the first landing's gate, refuses with the owner's assignment before
// it composes, so it meets neither the gate lock nor its diagnostic.
func TestSecondLandingMeetsTheLease(t *testing.T) {
	t.Parallel()
	request, secondRequest := "land-lease-first", "land-lease-second"
	f := barrierLandingFixture(t, request)
	second := mustCreate(t, f.root, f.home, secondRequest, "second landing")
	commitmenttest.Admit(t, second.Path, secondRequest, closureSpec)
	commitInWorktree(t, second.Path, "owned.txt", "second bytes\n", "second source")
	refreshLandingEvidence(t, second.Path, f.base)
	secondTip := gitOutput(t, second.Path, "rev-parse", "HEAD")
	composed := 0
	j := defaultJoins()
	j.landReviewed = func(ctx context.Context, request landing.ReviewedRequest, admission landing.Admission) (landing.ReviewedResult, error) {
		composed++
		return defaultJoins().landReviewed(ctx, request, admission)
	}
	var r2 verbResult
	done := f.during(func() {
		r2 = runVerb(t, verbLand, f.callWith(j, landArgs(secondRequest, f.base, secondTip, second.Path)...))
	})
	r1 := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
	select {
	case <-done:
	default:
		t.Fatalf("first landing finished without reaching its gate barrier: (%d, %q, %q)", r1.exit, r1.stdout, r1.stderr)
	}
	owner := f.creation.Assignment.ID
	if r2.exit != 1 || !strings.Contains(r2.stdout, "refused{detail=landing destination refs/heads/main is leased by assignment "+owner+",") {
		t.Fatalf("second landing = (%d, %q, %q), want the lease refusal naming %s", r2.exit, r2.stdout, r2.stderr, owner)
	}
	next, _ := landingFaceNext(r2.stdout, "")
	if rerun := landArgsRerun(secondRequest, f.base, secondTip, second.Path); !strings.HasPrefix(next, "wait until the landing of assignment "+owner+" ends") || !strings.HasSuffix(next, refusalroute.StepJoiner+rerun) {
		t.Fatalf("second landing next = %q, want the wait for %s and then the re-run %q", next, owner, rerun)
	}
	if composed != 0 {
		t.Fatalf("second landing composed %d times, want none", composed)
	}
	if tally, err := os.ReadFile(f.tally); err != nil || string(tally) != "g" {
		t.Fatalf("gate tally = %q, %v; want the first landing's one run", tally, err)
	}
	for _, stale := range []string{"bench doctor", "gate owner: pid"} {
		if strings.Contains(r2.stdout+r2.stderr, stale) {
			t.Fatalf("second landing = (%q, %q), want no %q", r2.stdout, r2.stderr, stale)
		}
	}
	if r1.exit != 0 {
		t.Fatalf("first landing = (%d, %q, %q), want a publication", r1.exit, r1.stdout, r1.stderr)
	}
}

// TestLandReleasesTheLeaseOnRefusal is BW8. A landing that refuses red at its gate leaves
// no lease, so the next landing on the destination takes the lease and publishes.
func TestLandReleasesTheLeaseOnRefusal(t *testing.T) {
	t.Parallel()
	request := "land-lease-refusal"
	red := filepath.Join(t.TempDir(), "red")
	mustWrite(t, red, []byte("red\n"), 0o600)
	f := landingFixtureWithGateStep(t, request, "", "", filepath.Join(t.TempDir(), "bench-home"), true, func(*testrepo.GateFixture, string) string {
		return "[ ! -e " + sanitize.ShellQuote(red) + " ]\n"
	})
	args := landArgs(request, f.base, f.tip, f.creation.Path)
	if r := runVerb(t, verbLand, f.call(args...)); r.exit != 1 || !strings.Contains(r.stdout, "refused{detail=prospective authorization refused: ") {
		t.Fatalf("red landing = (%d, %q, %q), want the red refusal", r.exit, r.stdout, r.stderr)
	}
	leases, err := commitmenttest.Leases(f.root)
	mustNoError(t, err)
	if len(leases) != 0 {
		t.Fatalf("leases after the red refusal = %+v, want none", leases)
	}
	mustRemove(t, red)
	if r := runVerb(t, verbLand, f.call(args...)); r.exit != 0 {
		t.Fatalf("next landing = (%d, %q, %q), want a publication", r.exit, r.stdout, r.stderr)
	}
}

// TestLandResumeTakesNoLease is BW11. A resume never composes or publishes, so it
// completes a published landing while another landing holds the destination's lease.
func TestLandResumeTakesNoLease(t *testing.T) {
	t.Parallel()
	request := "land-lease-resume"
	f := publicLandingFixture(t, request, "", "")
	interruptedLanding(t, f, request, f.tip)
	published := gitOutput(t, f.root, "rev-parse", "main")
	holdLandingLease(t, f.root, otherLandingID)
	r := runVerb(t, verbLand, f.call("--resume", published, "--request", request, "--base", f.base, "--source-tip", f.tip, "--spec", "x", f.creation.Path))
	if worktree, landed := landedField(r.stdout, "worktree"); r.exit != 0 || !landed || worktree != "released" {
		t.Fatalf("resume under another landing's lease = (%d, %q, %q), want a released landed record", r.exit, r.stdout, r.stderr)
	}
}

// TestLeaseLockUnavailable is BW15. A directory at the lease's lock path is a lock fault,
// so the landing refuses before composition rather than run without a lease.
func TestLeaseLockUnavailable(t *testing.T) {
	t.Parallel()
	request := "land-lease-unavailable"
	f := publicLandingFixture(t, request, "", "")
	lock, err := commitrepo.LeaseLockPath(f.root, leasedMain)
	mustNoError(t, err)
	mustMkdirAll(t, lock, 0o755)
	j, composed := forbidLandingComposition()
	r := runVerb(t, verbLand, f.callWith(j, landArgs(request, f.base, f.tip, f.creation.Path)...))
	r.mustViaJoins(t)
	if r.exit != 1 || !strings.Contains(r.stdout, "refused{detail=landing lease unavailable: ") {
		t.Fatalf("landing over a blocked lease lock = (%d, %q, %q), want the lease fault refusal", r.exit, r.stdout, r.stderr)
	}
	requireIdentityRefusalState(t, f.root, f.creation.Path, f.tally, *composed)
}
