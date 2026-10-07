// Assignment identity tests for the landing command: each identity dimension —
// request, review base, source tip, and source fingerprint — invalidates the landing
// before composition when it changes.
package worktree

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/testrepo"
)

// forbidLandingComposition returns a seam set that fails the test if the landing
// reaches composition, and a counter the assertions read. The identity seams must
// refuse first.
func forbidLandingComposition() (joins, *int) {
	composed := 0
	j := defaultJoins()
	j.landReviewed = func(context.Context, landing.ReviewedRequest, landing.Admission) (landing.ReviewedResult, error) {
		composed++
		return landing.ReviewedResult{}, errors.New("composition must not start")
	}
	return j, &composed
}

// TestLandCommandInvalidatesAChangedRequestBeforeComposition is SOL05.
func TestLandCommandInvalidatesAChangedRequestBeforeComposition(t *testing.T) {
	t.Parallel()
	request := "land-identity-request"
	f := publicLandingFixture(t, request, "", "")
	j, composed := forbidLandingComposition()

	r := runVerb(t, verbLand, f.callWith(j, landArgs("land-identity-request-changed", f.base, f.tip, f.creation.Path)...))
	r.mustViaJoins(t)
	if r.exit != 1 || !strings.HasPrefix(r.stdout, "refused{") {
		t.Fatalf("changed request = (%d, %q, %q), want a refusal", r.exit, r.stdout, r.stderr)
	}
	requireIdentityRefusalState(t, f.root, f.creation.Path, f.tally, *composed)
}

// TestLandCommandInvalidatesAChangedReviewBaseBeforeComposition is SOL06.
func TestLandCommandInvalidatesAChangedReviewBaseBeforeComposition(t *testing.T) {
	t.Parallel()
	request := "land-identity-base"
	f := publicLandingFixture(t, request, "", "")
	j, composed := forbidLandingComposition()

	r := runVerb(t, verbLand, f.callWith(j, landArgs(request, f.tip, f.tip, f.creation.Path)...))
	r.mustViaJoins(t)
	if r.exit != 1 || !strings.HasPrefix(r.stdout, "refused{") {
		t.Fatalf("changed review base = (%d, %q, %q), want a refusal", r.exit, r.stdout, r.stderr)
	}
	requireIdentityRefusalState(t, f.root, f.creation.Path, f.tally, *composed)
}

// TestLandCommandInvalidatesAChangedSourceTipBeforeComposition is SOL07.
func TestLandCommandInvalidatesAChangedSourceTipBeforeComposition(t *testing.T) {
	t.Parallel()
	request := "land-identity-tip"
	f := publicLandingFixture(t, request, "", "")
	commitInWorktree(t, f.creation.Path, "moved.txt", "moved\n", "tip moved after review")
	j, composed := forbidLandingComposition()

	r := runVerb(t, verbLand, f.callWith(j, landArgs(request, f.base, f.tip, f.creation.Path)...))
	r.mustViaJoins(t)
	if r.exit != 1 || !strings.Contains(r.stdout, "source tip mismatch") {
		t.Fatalf("changed source tip = (%d, %q, %q), want a tip-mismatch refusal", r.exit, r.stdout, r.stderr)
	}
	requireIdentityRefusalState(t, f.root, f.creation.Path, f.tally, *composed)
}

// TestLandCommandInvalidatesAChangedSourceFingerprintBeforeTheGate is SOL08. A
// commit-only comparison accepts dirty source content; the fingerprint proof refuses
// it before composition, so before the gate.
func TestLandCommandInvalidatesAChangedSourceFingerprintBeforeTheGate(t *testing.T) {
	t.Parallel()
	request := "land-identity-fingerprint"
	f := publicLandingFixture(t, request, "", "")
	mustWrite(t, filepath.Join(f.creation.Path, "dirty.txt"), []byte("uncommitted\n"), 0o600)
	j, composed := forbidLandingComposition()

	r := runVerb(t, verbLand, f.callWith(j, landArgs(request, f.base, f.tip, f.creation.Path)...))
	r.mustViaJoins(t)
	// The sentence and the repair read from the registry, so each keeps one source. The
	// hostile-source surface stays bounded: the refusal carries a route and no path table,
	// so no source-authored path name reaches the operator's terminal.
	face := landingRefusalFaceByName(faceSourceNotClean)
	if r.exit != 1 || !strings.Contains(r.stdout, face.detail) || !strings.Contains(r.stdout, "next="+face.route("")) ||
		strings.Contains(r.stdout, refusalPathsTable) || strings.Contains(r.stdout, "dirty.txt") {
		t.Fatalf("changed source fingerprint = (%d, %q, %q), want a routed not-clean refusal with no path table", r.exit, r.stdout, r.stderr)
	}
	requireIdentityRefusalState(t, f.root, f.creation.Path, f.tally, *composed)
}

// requireIdentityRefusalState pins what every identity refusal leaves behind: no
// composition, no gate run, an unmoved destination, and a retained source worktree.
func requireIdentityRefusalState(t *testing.T, root, path, tally string, composed int) {
	t.Helper()
	if composed != 0 {
		t.Fatalf("identity refusal reached composition %d times, want 0", composed)
	}
	if _, err := os.Stat(tally); !os.IsNotExist(err) {
		t.Fatalf("identity refusal ran the gate: %v", err)
	}
	if got := gitOutput(t, root, "status", "--porcelain=v1"); got != "" {
		t.Fatalf("identity refusal dirtied the destination: %q", got)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("identity refusal removed the source worktree: %v", err)
	}
}

// TestLandCommandRefusesAReviewBaseThatIsNotAnAncestorOfTheDestination covers the
// confusion the detail is for: the source folded the destination advance, and the
// operator names that fold commit as `--base`. The fold is a descendant of the recorded
// start, so the recorded-start binding accepts it, but it lives on the source branch and
// the destination never composed it. The detail names the two bases the flag can mean,
// because the fold commit a review reads and the destination tip the source folded are
// easy to confuse when only one prints as `--base`.
func TestLandCommandRefusesAReviewBaseThatIsNotAnAncestorOfTheDestination(t *testing.T) {
	t.Parallel()
	request := "land-identity-not-ancestor"
	f := foldedLandingFixture(t, request)
	if git.OK("-C", f.root, "merge-base", "--is-ancestor", f.fold, f.base) {
		t.Fatalf("premise failed: the fold commit %q must not already be on the destination %q", f.fold, f.base)
	}
	j, composed := forbidLandingComposition()

	r := runVerb(t, verbLand, f.callWith(j, specLessLandArgs(request, f.fold, f.tip, f.creation.Path)...))
	r.mustViaJoins(t)
	// The expectation is spelled out here rather than read from landingBaseNotAncestorDetail,
	// so a mutation of that constant turns this test red instead of passing silently.
	const wantDetail = "review base is not an ancestor of the landing destination: --base takes a default-branch commit, such as the tip merged before the first chunk, not the fold commit the review read"
	want := "detail=" + wantDetail + ",observed=" + f.fold + ",wanted=" + f.base
	if r.exit != 1 || !strings.Contains(r.stdout, want) {
		t.Fatalf("non-ancestor base = (%d, %q, %q), want a refusal carrying %q", r.exit, r.stdout, r.stderr, want)
	}
	requireIdentityRefusalState(t, f.root, f.creation.Path, f.tally, *composed)
}

// TestLandCommandRefusesAReviewBaseBehindTheRecordedStart is the SOL06 mutation guard
// the ancestry check alone cannot answer, and it closes C4. The base names an earlier
// ancestor of the recorded start, so it is a valid ancestor of the destination and the
// ancestry guard accepts it, and the spec-less landing names no ownership fence that
// would refuse it for another reason. Only the recorded-start binding refuses it,
// because that earlier base grades a range wider than the one the assignment
// authorized, and it refuses before composition.
func TestLandCommandRefusesAReviewBaseBehindTheRecordedStart(t *testing.T) {
	t.Parallel()
	request := "land-identity-recorded-start"
	f := specLessLandingFixture(t, request)
	earlier := gitOutput(t, f.root, "rev-parse", f.base+"~1")
	if earlier == f.base {
		t.Fatalf("fixture has no earlier ancestor than the recorded start %q", f.base)
	}
	if !git.OK("-C", f.root, "merge-base", "--is-ancestor", earlier, f.base) {
		t.Fatalf("premise failed: %q is not an ancestor of the destination", earlier)
	}
	j, composed := forbidLandingComposition()

	r := runVerb(t, verbLand, f.callWith(j, specLessLandArgs(request, earlier, f.tip, f.creation.Path)...))
	r.mustViaJoins(t)
	want := "detail=" + reviewedRangeDetail + ",observed=" + earlier + ",wanted=" + f.base
	if r.exit != 1 || !strings.Contains(r.stdout, want) {
		t.Fatalf("earlier ancestor base = (%d, %q, %q), want a refusal carrying %q", r.exit, r.stdout, r.stderr, want)
	}
	requireIdentityRefusalState(t, f.root, f.creation.Path, f.tally, *composed)
}

// switchActiveMilestone commits an approved switch at root that stops work on the fixture
// delivery outcome, and returns the new destination tip.
func switchActiveMilestone(root string) (string, error) {
	store := commitrepo.Store{Root: root}
	policy, _, err := store.Policy()
	if err != nil {
		return "", err
	}
	policy.Milestones = append(policy.Milestones, commitment.Milestone{ID: "next", Outcomes: []commitment.Outcome{{ID: "successor", Criteria: []commitment.Criterion{{ID: "successor-done", Text: "The successor is accepted."}}, Sources: []commitment.SourceBinding{}}}})
	policy.ActiveMilestone = "next"
	if err := commitmenttest.CommitPolicy(root, policy, "switch the active milestone"); err != nil {
		return "", err
	}
	tip, err := git.Output("-C", root, "rev-parse", "main")
	return strings.TrimSpace(tip), err
}

// An assignment that a later approved switch displaced cannot land its source,
// although its start-time binding admitted it.
func TestCommitmentStaleLanding(t *testing.T) {
	t.Parallel()
	request := "land-commitment-stale"
	f := publicLandingFixture(t, request, "", "")
	switched, err := switchActiveMilestone(f.root)
	if err != nil {
		t.Fatal(err)
	}
	r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
	if r.exit != 1 || !strings.Contains(r.stdout, "refused{detail=commitment: ") || !strings.Contains(r.stdout, "bench commitment start") {
		t.Fatalf("displaced landing = (%d, %q, %q), want a commitment refusal naming the start command", r.exit, r.stdout, r.stderr)
	}
	if got := gitOutput(t, f.root, "rev-parse", "main"); got != switched {
		t.Fatalf("displaced landing moved main to %s, want %s", got, switched)
	}
	requireIdentityRefusalState(t, f.root, f.creation.Path, f.tally, 0)
}

// A commitment change that lands while the gate runs refuses the now-disallowed
// publication. The gate waits at a FIFO barrier, so each change happens inside the gate.
func TestCommitmentGateRace(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name   string
		change func(root string) (string, error)
		want   []string
	}{
		{name: "blocker", change: func(root string) (string, error) {
			if err := (commitrepo.Store{Root: root}).Block(commitmenttest.DeliveryOutcome, "reviewer paused delivery"); err != nil {
				return "", err
			}
			tip, err := git.Output("-C", root, "rev-parse", "main")
			return strings.TrimSpace(tip), err
		}, want: []string{"refused{detail=commitment: ", "is blocked"}},
		{name: "default-branch-policy", change: switchActiveMilestone, want: []string{"landing destination checkout changed"}},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			request := "land-commitment-gate-race-" + row.name
			barrier := t.TempDir()
			started, release := filepath.Join(barrier, "started"), filepath.Join(barrier, "release")
			for _, fifo := range []string{started, release} {
				mustNoError(t, syscall.Mkfifo(fifo, 0o600))
			}
			f := landingFixtureWithGateStep(t, request, "", "", filepath.Join(t.TempDir(), "bench-home"), true, func(*testrepo.GateFixture, string) string {
				return ": > " + sanitize.ShellQuote(started) + "\nIFS= read -r _ < " + sanitize.ShellQuote(release) + "\n"
			})
			type changed struct {
				tip string
				err error
			}
			done := make(chan changed, 1)
			go func() {
				// The read returns once the gate has opened and closed its barrier. The
				// result is sent before the release, so the landing cannot finish first.
				if _, err := os.ReadFile(started); err != nil {
					done <- changed{err: err}
					return
				}
				tip, err := row.change(f.root)
				done <- changed{tip: tip, err: err}
				_ = os.WriteFile(release, []byte("go\n"), 0)
			}()
			r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
			var got changed
			select {
			case got = <-done:
			default:
				t.Fatalf("landing finished without reaching its gate barrier: (%d, %q, %q)", r.exit, r.stdout, r.stderr)
			}
			if got.err != nil {
				t.Fatal(got.err)
			}
			for _, want := range row.want {
				if r.exit != 1 || !strings.Contains(r.stdout, want) {
					t.Fatalf("gate-time change = (%d, %q, %q), want a refusal naming %q", r.exit, r.stdout, r.stderr, want)
				}
			}
			if main := gitOutput(t, f.root, "rev-parse", "main"); main != got.tip {
				t.Fatalf("gate-time change published main=%s, want %s", main, got.tip)
			}
			if tally, err := os.ReadFile(f.tally); err != nil || string(tally) != "g" {
				t.Fatalf("gate tally = %q, %v; want the one prospective run", tally, err)
			}
			if _, err := os.Stat(f.creation.Path); err != nil {
				t.Fatalf("refused landing removed the source worktree: %v", err)
			}
		})
	}
}

// publicationHoldWindow is how long the publication gap holds the intent lock while a
// competing blocker waits. It stays inside that blocker's own lock window, so the blocker
// still acquires the lock after the publication instead of timing out.
const publicationHoldWindow = bounds.IntentLockTimeout / 4

// A competing blocker that arrives between the final admission decision and the
// ref update waits for the publication. It then observes the published destination,
// whose policy records the outcome as delivered, so the blocker refuses.
func TestCommitmentPublishLock(t *testing.T) {
	t.Parallel()
	request := "land-commitment-publish-lock"
	f := publicLandingFixture(t, request, "", "")
	type blocked struct {
		main string
		err  error
	}
	competing := make(chan blocked, 1)
	early, entered := false, false
	j := defaultJoins()
	j.publicationGap = func(root string) {
		entered = true
		go func() {
			err := (commitrepo.Store{Root: root}).Block(commitmenttest.DeliveryOutcome, "competing blocker")
			main, readErr := git.Output("-C", root, "rev-parse", "main")
			competing <- blocked{main: strings.TrimSpace(main), err: errors.Join(err, readErr)}
		}()
		// A blocker that completes inside this window changed admission before the ref update.
		select {
		case got := <-competing:
			early = true
			competing <- got
		case <-time.After(publicationHoldWindow):
		}
	}
	r := runVerb(t, verbLand, f.callWith(j, landArgs(request, f.base, f.tip, f.creation.Path)...))
	r.mustViaJoins(t)
	if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released") {
		t.Fatalf("landing = (%d, %q, %q), want a published release", r.exit, r.stdout, r.stderr)
	}
	if !entered {
		t.Fatal("the landing published without entering its publication lock")
	}
	published := gitOutput(t, f.root, "rev-parse", "main")
	got := <-competing
	if early || got.err == nil || !strings.Contains(got.err.Error(), "is already delivered") || got.main != published {
		t.Fatalf("competing blocker = (early=%t, main=%s, err=%v), want it to wait for the publication %s", early, got.main, got.err, published)
	}
}
