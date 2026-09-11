// Landed-sibling cleanup tests for the landing command: the narrowed selector scope, the
// cleanup effect's own result, and the resume that finishes an unfinished effect.
package worktree

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

const siblingReviewPath = "reviews/sibling.md"

// foldLandingSibling mints one more assignment, commits its own reviewed bytes, and folds
// that branch into the landing source. The landing then carries the sibling's commits, so
// the landed proof holds for the sibling against the published commit and fails against
// the destination base. It answers the sibling and the source's new tip.
//
// The sibling writes its own review file, which the landing fixture's own spec declares in its
// ownership fence, so the folded range still authorizes.
func foldLandingSibling(t *testing.T, root, home, request string, source Creation) (Creation, string) {
	t.Helper()
	sibling := mustCreate(t, root, home, request, "folded sibling")
	mustMkdirAll(t, filepath.Join(sibling.Path, "reviews"), 0o755)
	commitInWorktree(t, sibling.Path, siblingReviewPath, "sibling review\n", "sibling review")
	gitRun(t, source.Path, "-c", "user.name=bench", "-c", "user.email=bench@local",
		"merge", "-q", "--no-ff", "-m", "fold the sibling", strings.TrimPrefix(sibling.Assignment.Branch, "refs/heads/"))
	refreshLandingEvidence(t, source.Path, gitOutput(t, root, "merge-base", "main", source.Assignment.Branch))
	return sibling, gitOutput(t, source.Path, "rev-parse", "HEAD")
}

// assignmentActive reports whether root still records assignment as an active one. The
// cleanup retires the record beside the checkout, so a removal is only proved when both
// are gone.
func assignmentActive(t *testing.T, root, id string) bool {
	t.Helper()
	assignments, err := intent.Assignments(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, assignment := range assignments {
		if assignment.ID == id && assignment.State == intent.StateActive {
			return true
		}
	}
	return false
}

// requireAbsent fails unless path is gone. A cleanup that reports complete over a
// surviving checkout is the failure this reads.
func requireAbsent(t *testing.T, path, what string) {
	t.Helper()
	if _, err := os.Lstat(path); !os.IsNotExist(err) {
		t.Fatalf("%s %q survived the cleanup: %v", what, path, err)
	}
}

// requirePresent fails unless path is still there.
func requirePresent(t *testing.T, path, what string) {
	t.Helper()
	if _, err := os.Lstat(path); err != nil {
		t.Fatalf("%s %q was removed: %v", what, path, err)
	}
}

// TestLandCleansTheFoldedSibling is LC9 and LC42. A landing that folds a sibling
// assignment retires that sibling's checkout and its record inside the same process, so
// the effects row reports a complete cleanup. The landing plans and applies in one call,
// so its stdout carries no fingerprint and no apply action.
func TestLandCleansTheFoldedSibling(t *testing.T) {
	t.Parallel()
	request := "land-cleanup-folded-sibling"
	root, creation, base, _, _, home := publicLandingFixture(t, request, "", "")
	sibling, tip := foldLandingSibling(t, root, home, request+"-sibling", creation)
	j, _ := refreshJoins(nil)

	var stdout, stderr bytes.Buffer
	code := landWith(j, root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 || !strings.Contains(stdout.String(), wantEffects("not-applicable", "complete")) {
		t.Fatalf("folded-sibling landing = (%d, %q, %q), want a complete cleanup at exit 0", code, stdout.String(), stderr.String())
	}
	requireAbsent(t, sibling.Path, "sibling worktree")
	if assignmentActive(t, root, sibling.Assignment.ID) {
		t.Fatalf("sibling assignment %q is still active", sibling.Assignment.ID)
	}
	if strings.Contains(stdout.String(), "--apply") || regexp.MustCompile(`[0-9a-f]{64}`).MatchString(stdout.String()) {
		t.Fatalf("landing stdout = %q, want no cleanup fingerprint and no apply action", stdout.String())
	}
	if !strings.Contains(stderr.String(), "landing cleanup{assignment="+sibling.Assignment.ID+",action=remove,") {
		t.Fatalf("landing stderr = %q, want the sibling's plan row", stderr.String())
	}
}

// TestLandLeavesAPriorLandedAssignment is LC10. An assignment whose work the destination
// base already carries belongs to an earlier landing, not to this one. The narrowed scope
// keeps it out of the set, so one landing does no repository-wide maintenance.
func TestLandLeavesAPriorLandedAssignment(t *testing.T) {
	t.Parallel()
	request := "land-cleanup-prior-landed"
	root, creation, _, _, _, home := publicLandingFixture(t, request, "", "")
	prior := mustCreate(t, root, home, request+"-prior", "prior landing")
	landAssignment(t, root, prior, "prior.txt")
	gitRun(t, creation.Path, "rebase", "main")
	base := gitOutput(t, root, "rev-parse", "HEAD")
	sibling, tip := foldLandingSibling(t, root, home, request+"-sibling", creation)
	j, _ := refreshJoins(nil)

	var stdout, stderr bytes.Buffer
	code := landWith(j, root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 || !strings.Contains(stdout.String(), wantEffects("not-applicable", "complete")) {
		t.Fatalf("narrowed landing = (%d, %q, %q), want a complete cleanup at exit 0", code, stdout.String(), stderr.String())
	}
	requireAbsent(t, sibling.Path, "sibling worktree")
	requirePresent(t, prior.Path, "prior-landed worktree")
	if !assignmentActive(t, root, prior.Assignment.ID) {
		t.Fatalf("prior-landed assignment %q was retired", prior.Assignment.ID)
	}
	if strings.Contains(stderr.String(), prior.Assignment.ID) {
		t.Fatalf("landing stderr = %q, want no plan row for the prior-landed assignment", stderr.String())
	}
}

// TestLandRetainsAnUnprovenSibling is LC11, LC12, and LC43. A folded sibling that holds
// uncommitted tracked work keeps its checkout, and its plan row names the action and the
// reason code the preservation verdict returned. A retained row whose path cannot be
// carried on one line names its assignment pointer instead.
func TestLandRetainsAnUnprovenSibling(t *testing.T) {
	t.Parallel()
	request := "land-cleanup-unproven-sibling"
	root, creation, base, _, _, home := publicLandingFixture(t, request, "", "")
	sibling, tip := foldLandingSibling(t, root, home, request+"-sibling", creation)
	mustWrite(t, filepath.Join(sibling.Path, siblingReviewPath), []byte("uncommitted review\n"), 0o644)
	j, _ := refreshJoins(nil)

	var stdout, stderr bytes.Buffer
	code := landWith(j, root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 || !strings.Contains(stdout.String(), wantEffects("not-applicable", "complete")) {
		t.Fatalf("retaining landing = (%d, %q, %q), want a complete cleanup at exit 0", code, stdout.String(), stderr.String())
	}
	requirePresent(t, sibling.Path, "retained sibling worktree")
	want := "landing cleanup{assignment=" + sibling.Assignment.ID + ",action=retain,reason=dirty,target=" + sibling.Path + "}\n"
	if !strings.Contains(stderr.String(), want) {
		t.Fatalf("landing stderr = %q, want the retained row %q", stderr.String(), want)
	}

	t.Run("hostile path", func(t *testing.T) {
		hostile := sibling
		hostile.Assignment.Worktree = sibling.Path + "\nforged"
		var rows bytes.Buffer
		printLandedCleanupRows(&rows, landedCleanupSet{rows: []landedCleanupRow{{
			assignment: hostile.Assignment,
			plan:       CleanupPlan{Action: ActionRetain, ReasonCode: ReasonDirty},
		}}})
		pointer := "landing cleanup{assignment=" + sibling.Assignment.ID + ",action=retain,reason=dirty,target=assignment/" + sibling.Assignment.ID + "}\n"
		if rows.String() != pointer {
			t.Fatalf("hostile-path row = %q, want the assignment pointer %q", rows.String(), pointer)
		}
	})
}

// TestLandReportsAFailedCleanup is LC14. A fault inside the apply leaves the landing
// unfinished, so the effect reports failed, the record names the cleanup step, and the
// exit carries the resume every incomplete step carries.
func TestLandReportsAFailedCleanup(t *testing.T) {
	t.Parallel()
	request := "land-cleanup-apply-fault"
	root, creation, base, _, _, home := publicLandingFixture(t, request, "", "")
	sibling, tip := foldLandingSibling(t, root, home, request+"-sibling", creation)
	j, _ := refreshJoins(nil)
	// The release step runs its own cleanup transaction first, so the fault is bound to
	// the lock the cleanup effect's own row takes.
	locks := 0
	j.cleanupBoundary = func(step LifecycleStep) error {
		if step != StepApplyLocked {
			return nil
		}
		locks++
		if locks > 1 {
			return errStaleFingerprint
		}
		return nil
	}

	var stdout, stderr bytes.Buffer
	code := landWith(j, root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 3 || !strings.Contains(stdout.String(), wantEffects("not-applicable", "failed")) {
		t.Fatalf("faulted cleanup = (%d, %q, %q), want exit 3 with a failed cleanup", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "worktree=incomplete:cleanup,next=bench worktree land --resume ") {
		t.Fatalf("faulted cleanup record = %q, want an incomplete:cleanup cell with a resume", stdout.String())
	}
	requirePresent(t, sibling.Path, "sibling worktree")
}

// resumeLandArgs is the resume grammar's own argument list for a published landing.
func resumeLandArgs(published, request, base, tip, path string) []string {
	return []string{"--resume", published, "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", path}
}

// TestResumeLandCompletesTheUnfinishedEffects is LC16 and LC17. A resume runs the refresh
// and then the cleanup, it starts no cleanup while the refresh fails, and it repeats
// neither effect once each is complete. It publishes nothing a second time.
func TestResumeLandCompletesTheUnfinishedEffects(t *testing.T) {
	t.Parallel()
	request := "land-cleanup-resume-unfinished"
	root, creation, base, _, home := brokerDestinationFixture(t, request)
	sibling, tip := foldLandingSibling(t, root, home, request+"-sibling", creation)
	failing, failingCalls := refreshJoins(nil)
	working, workingCalls := refreshJoins(func(root, executable string) error {
		return publishVerifyingBroker(t, root, executable)
	})
	interrupted := failing
	interrupted.releaseLandingAssignment = func(joins, string, string, []string, io.Writer, io.Writer) int { return 1 }

	var stdout, stderr bytes.Buffer
	if code := landWith(interrupted, root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr); code != 3 || !strings.Contains(stdout.String(), "worktree=incomplete:release") {
		t.Fatalf("interrupted landing = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
	published := gitOutput(t, root, "rev-parse", "main")
	args := resumeLandArgs(published, request, base, tip, creation.Path)

	stdout.Reset()
	stderr.Reset()
	if code := resumeLandWith(failing, root, home, args, &stdout, &stderr); code != 3 || !strings.Contains(stdout.String(), wantEffects("failed")) {
		t.Fatalf("failed-refresh resume = (%d, %q, %q), want a pending cleanup at exit 3", code, stdout.String(), stderr.String())
	}
	requirePresent(t, sibling.Path, "sibling worktree")
	if *failingCalls != 1 {
		t.Fatalf("failed-refresh resume build calls = %d, want exactly one", *failingCalls)
	}

	stdout.Reset()
	stderr.Reset()
	if code := resumeLandWith(working, root, home, args, &stdout, &stderr); code != 0 || !strings.Contains(stdout.String(), wantEffects("complete", "complete")) {
		t.Fatalf("repairing resume = (%d, %q, %q), want both effects complete", code, stdout.String(), stderr.String())
	}
	requireAbsent(t, sibling.Path, "sibling worktree")
	if *workingCalls != 1 {
		t.Fatalf("repairing resume build calls = %d, want exactly one", *workingCalls)
	}

	stdout.Reset()
	stderr.Reset()
	if code := resumeLandWith(working, root, home, args, &stdout, &stderr); code != 0 || !strings.Contains(stdout.String(), wantEffects("complete", "complete")) {
		t.Fatalf("second resume = (%d, %q, %q), want both effects left alone", code, stdout.String(), stderr.String())
	}
	if *workingCalls != 1 {
		t.Fatalf("second resume build calls = %d, want no second build", *workingCalls)
	}
	if strings.Contains(stderr.String(), "landing cleanup{") {
		t.Fatalf("second resume stderr = %q, want an empty narrowed set", stderr.String())
	}
	if got := gitOutput(t, root, "rev-parse", "main"); got != published {
		t.Fatalf("resume republished: main = %s, want %s", got, published)
	}
}

// TestResumeLandRunsTheEffectsAfterRelease is LC18. A landing whose refresh failed after
// its release leaves no active assignment, and the resume's terminal path still owns both
// effects. It runs them rather than reporting a complete landing over an unfinished tree.
func TestResumeLandRunsTheEffectsAfterRelease(t *testing.T) {
	t.Parallel()
	request := "land-cleanup-resume-released"
	root, creation, base, _, home := brokerDestinationFixture(t, request)
	sibling, tip := foldLandingSibling(t, root, home, request+"-sibling", creation)
	failing, _ := refreshJoins(nil)
	working, workingCalls := refreshJoins(func(root, executable string) error {
		return publishVerifyingBroker(t, root, executable)
	})

	var stdout, stderr bytes.Buffer
	if code := landWith(failing, root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr); code != 3 || !strings.Contains(stdout.String(), wantEffects("failed")) {
		t.Fatalf("failed-refresh landing = (%d, %q, %q), want a pending cleanup at exit 3", code, stdout.String(), stderr.String())
	}
	requirePresent(t, sibling.Path, "sibling worktree")
	published := gitOutput(t, root, "rev-parse", "main")

	stdout.Reset()
	stderr.Reset()
	code := resumeLandWith(working, root, home, resumeLandArgs(published, request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 || !strings.Contains(stdout.String(), wantEffects("complete", "complete")) {
		t.Fatalf("released resume = (%d, %q, %q), want both effects complete", code, stdout.String(), stderr.String())
	}
	requireAbsent(t, sibling.Path, "sibling worktree")
	if *workingCalls != 1 {
		t.Fatalf("released resume build calls = %d, want exactly one", *workingCalls)
	}
}
