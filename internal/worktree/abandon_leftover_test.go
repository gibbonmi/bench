package worktree

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/intent"
)

// decayToHusk leaves the assignment directory in place, holding content, with its git
// metadata entry gone: the shape a partially removed checkout leaves behind.
func decayToHusk(t *testing.T, path string) {
	t.Helper()
	mustWrite(t, filepath.Join(path, "keep.txt"), []byte("husk bytes\n"), 0o644)
	mustRemove(t, filepath.Join(path, ".git"))
}

func TestApplyAbandonReleasesHuskWithoutDeletingBytes(t *testing.T) {
	const request = "landed-abandon-husk"
	root, creation := newOwnedAssignment(t, "abandon-husk")
	decayToHusk(t, creation.Path)

	plan, err := planAbandon(root, request, creation.Path)
	mustNoError(t, err)
	requireTest(t, plan.Action == actionReleaseLeftover && plan.leftover == creation.Path,
		"husk plan = %#v, want action %q naming leftover %s", plan, actionReleaseLeftover, creation.Path)
	requireTest(t, !plan.Action.removes(), "husk plan action %q still has a removal ahead of it", plan.Action)

	result, err := ApplyAbandon(root, request, creation.Path, plan.Fingerprint)
	requireTest(t, err == nil && result.Action == ActionRemoved, "ApplyAbandon over a husk = %#v, %v", result, err)
	registrations := gitOutput(t, root, "worktree", "list", "--porcelain")
	requireTest(t, !strings.Contains(registrations, "worktree "+creation.Path),
		"registration survived the husk abandon:\n%s", registrations)
	_, found, err := intent.FindAssignmentByRequest(root, requestDigest(request))
	requireTest(t, err == nil && !found, "intent entry survived the husk abandon: found=%t, %v", found, err)
	body, readErr := os.ReadFile(filepath.Join(creation.Path, "keep.txt"))
	requireTest(t, readErr == nil && string(body) == "husk bytes\n", "husk content after apply = %q, %v", body, readErr)
}

// abandonDeadline bounds an abandon that must decide by shape alone. A FIFO at an
// assignment path or in a discovered control record has no writer, so a read of one never
// returns; running the abandon off the test goroutine turns that into this test's own
// failure instead of a package-wide timeout.
const abandonDeadline = 15 * time.Second

func applyAbandonWithin(t *testing.T, root, request, path, fingerprint string) (CleanupPlan, error) {
	t.Helper()
	type outcome struct {
		plan CleanupPlan
		err  error
	}
	done := make(chan outcome, 1)
	go func() {
		plan, err := ApplyAbandon(root, request, path, fingerprint)
		done <- outcome{plan, err}
	}()
	select {
	case got := <-done:
		return got.plan, got.err
	case <-time.After(abandonDeadline):
		t.Fatalf("ApplyAbandon over %s blocked instead of deciding by shape", path)
		return CleanupPlan{}, nil
	}
}

// TestAbandonReleasesEveryPresentNonCheckoutShape drives every present-but-not-a-checkout
// shape through plan and apply. Each row's entry must survive apply exactly as it was, and
// the FIFO row has no writer on the other end, so an implementation that opens the path in
// either half never answers.
func TestAbandonReleasesEveryPresentNonCheckoutShape(t *testing.T) {
	shapes := []struct {
		name     string
		decay    func(*testing.T, string)
		survives func(*testing.T, string)
	}{
		{"husk", decayToHusk, func(t *testing.T, path string) {
			body, err := os.ReadFile(filepath.Join(path, "keep.txt"))
			requireTest(t, err == nil && string(body) == "husk bytes\n", "husk content after apply = %q, %v", body, err)
		}},
		{"dangling-symlink", func(t *testing.T, path string) {
			mustRemove(t, path)
			mustNoError(t, os.Symlink(filepath.Join(filepath.Dir(path), "gone"), path))
		}, func(t *testing.T, path string) {
			target, err := os.Readlink(path)
			requireTest(t, err == nil && target == filepath.Join(filepath.Dir(path), "gone"),
				"symlink after apply = %q, %v", target, err)
		}},
		{"fifo", func(t *testing.T, path string) {
			mustRemove(t, path)
			mustNoError(t, syscall.Mkfifo(path, 0o600))
		}, func(t *testing.T, path string) {
			info, err := os.Lstat(path)
			requireTest(t, err == nil && info.Mode()&os.ModeNamedPipe != 0, "fifo after apply = %v, %v", info, err)
		}},
		{"regular-file", func(t *testing.T, path string) {
			mustRemove(t, path)
			mustWrite(t, path, []byte("not a checkout\n"), 0o644)
		}, func(t *testing.T, path string) {
			body, err := os.ReadFile(path)
			requireTest(t, err == nil && string(body) == "not a checkout\n", "file content after apply = %q, %v", body, err)
		}},
	}
	for _, shape := range shapes {
		t.Run(shape.name, func(t *testing.T) {
			request := "landed-abandon-shape-" + shape.name
			root, creation := newOwnedAssignment(t, "abandon-shape-"+shape.name)
			shape.decay(t, creation.Path)

			plan, err := planAbandon(root, request, creation.Path)
			requireTest(t, err == nil && plan.Action == actionReleaseLeftover && plan.leftover == creation.Path,
				"%s plan = %#v, %v; want action %q naming leftover %s", shape.name, plan, err, actionReleaseLeftover, creation.Path)

			result, err := applyAbandonWithin(t, root, request, creation.Path, plan.Fingerprint)
			requireTest(t, err == nil && result.Action == ActionRemoved, "%s apply = %#v, %v", shape.name, result, err)
			registrations := gitOutput(t, root, "worktree", "list", "--porcelain")
			requireTest(t, !strings.Contains(registrations, "worktree "+creation.Path),
				"%s registration survived the abandon:\n%s", shape.name, registrations)
			_, found, err := intent.FindAssignmentByRequest(root, requestDigest(request))
			requireTest(t, err == nil && !found, "%s ledger entry survived the abandon: found=%t, %v", shape.name, found, err)
			shape.survives(t, creation.Path)
		})
	}
}

// TestHuskReleaseCompletesAfterRegistrationIsAlreadyReleased pins the retry of an abandon
// interrupted at its point of no return. The bytes stay put by design, so the target
// alone cannot tell the retry how far the first run got; only the released registration
// can, and mistaking the leftover for one never planned wedges the abandon.
func TestHuskReleaseCompletesAfterRegistrationIsAlreadyReleased(t *testing.T) {
	const request = "landed-abandon-husk-interrupted"
	root, creation := newOwnedAssignment(t, "abandon-husk-interrupted")
	decayToHusk(t, creation.Path)
	first, err := PlanAbandon(root, request, creation.Path)
	mustNoError(t, err)

	stop := errors.New("stop after registration release")
	old := cleanupTransactionBoundary
	cleanupTransactionBoundary = failLifecycleStep(StepRemoval, stop)
	_, err = ApplyAbandon(root, request, creation.Path, first)
	cleanupTransactionBoundary = old
	requireTest(t, errors.Is(err, stop), "interrupted husk release = %v, want %v", err, stop)
	interrupted := gitOutput(t, root, "worktree", "list", "--porcelain")
	requireTest(t, !strings.Contains(interrupted, "worktree "+creation.Path),
		"interruption landed before the registration release:\n%s", interrupted)
	pending, err := assignmentByID(root, creation.Assignment.ID)
	requireTest(t, err == nil && pending.State != intent.StateComplete,
		"interrupted assignment = %#v, %v; want a ledger entry still to release", pending, err)

	retry, err := PlanAbandon(root, request, creation.Path)
	requireTest(t, err == nil && retry == first, "retry fingerprint = %q, %v; want the first run's %q", retry, err, first)
	result, err := applyAbandonWithin(t, root, request, creation.Path, retry)
	requireTest(t, err == nil && result.Action == ActionRemoved, "husk release retry = %#v, %v", result, err)
	_, found, err := intent.FindAssignmentByRequest(root, requestDigest(request))
	requireTest(t, err == nil && !found, "ledger entry survived the retry: found=%t, %v", found, err)
	repo, target, err := cleanupIdentity(root, creation.Path)
	mustNoError(t, err)
	receipt, found, err := intent.CleanupReceiptFor(root, repo, cleanupOperation, target, first)
	requireTest(t, err == nil && found && receipt.State == intent.ReceiptComplete && receipt.Phase == intent.ReceiptPhaseTerminal,
		"retried cleanup receipt = %#v, found=%t error=%v", receipt, found, err)
	body, readErr := os.ReadFile(filepath.Join(creation.Path, "keep.txt"))
	requireTest(t, readErr == nil && string(body) == "husk bytes\n", "husk content after the retry = %q, %v", body, readErr)
}

// TestReleaseRegistrationSkipsUnrelatedSpecialControlRecords plants a FIFO where an
// unrelated private administration entry keeps its gitdir record. Releasing one
// registration sweeps the whole pool, so a stranger's record is reached too, and reading
// this one never returns.
//
// The release is driven directly rather than through ApplyAbandon: `git worktree list`
// reads the same pool and blocks on the same FIFO, so no abandon ever gets far enough to
// show what this sweep does with it. The record is planted after the assignment's own
// registration is known, which is also the only window a live repository has for one.
func TestReleaseRegistrationSkipsUnrelatedSpecialControlRecords(t *testing.T) {
	const request = "landed-abandon-husk-special-record"
	root, creation := newOwnedAssignment(t, "abandon-husk-special-record")
	decayToHusk(t, creation.Path)
	_, err := planAbandon(root, request, creation.Path)
	mustNoError(t, err)
	common := gitOutput(t, root, "rev-parse", "--path-format=absolute", "--git-common-dir")
	stranger := filepath.Join(filepath.Clean(common), "worktrees", "stranger")
	mustMkdirAll(t, stranger, 0o700)
	mustNoError(t, syscall.Mkfifo(filepath.Join(stranger, "gitdir"), 0o600))

	done := make(chan error, 1)
	go func() { done <- releaseRegistration(root, creation.Path) }()
	select {
	case releaseErr := <-done:
		mustNoError(t, releaseErr)
	case <-time.After(abandonDeadline):
		t.Fatal("releaseRegistration blocked on an unrelated special control record")
	}
	info, statErr := os.Lstat(filepath.Join(stranger, "gitdir"))
	requireTest(t, statErr == nil && info.Mode()&os.ModeNamedPipe != 0,
		"the release disturbed the unrelated administration entry: %v, %v", info, statErr)
	// Retiring the planted record before asking git anything: it blocks on this FIFO too.
	mustRemove(t, stranger)
	registrations := gitOutput(t, root, "worktree", "list", "--porcelain")
	requireTest(t, !strings.Contains(registrations, "worktree "+creation.Path),
		"registration survived the release:\n%s", registrations)
}

// TestClassifyPathShapeRefusesSpecialGitEntry plants a no-writer FIFO at the checkout's
// .git entry: the exact shape a fail-open classifier would hand to `git -C <path>
// rev-parse` without ever finishing, since git opens .git to follow a gitfile pointer.
// ClassifyPathShape must decide by shape alone, and planAbandon must refuse before
// invoking git at all, so both sides of this test run off the test goroutine and fail
// the moment they miss the deadline instead of wedging the suite.
func TestClassifyPathShapeRefusesSpecialGitEntry(t *testing.T) {
	root, creation := newOwnedAssignment(t, "abandon-special-git-fifo")
	mustRemove(t, filepath.Join(creation.Path, ".git"))
	mustNoError(t, syscall.Mkfifo(filepath.Join(creation.Path, ".git"), 0o600))

	type shapeOutcome struct {
		shape PathShape
		err   error
	}
	shapeDone := make(chan shapeOutcome, 1)
	go func() {
		shape, err := ClassifyPathShape(creation.Path)
		shapeDone <- shapeOutcome{shape, err}
	}()
	select {
	case got := <-shapeDone:
		requireTest(t, got.err == nil && got.shape == ShapeSpecialMetadata,
			"ClassifyPathShape over FIFO .git = %v, %v; want %v, <nil>", got.shape, got.err, ShapeSpecialMetadata)
	case <-time.After(abandonDeadline):
		t.Fatal("ClassifyPathShape blocked reading a no-writer FIFO .git entry")
	}

	request := "landed-abandon-special-git-fifo"
	type planOutcome struct {
		plan CleanupPlan
		err  error
	}
	planDone := make(chan planOutcome, 1)
	go func() {
		plan, err := planAbandon(root, request, creation.Path)
		planDone <- planOutcome{plan, err}
	}()
	select {
	case got := <-planDone:
		requireTest(t, got.err == errSpecialMetadata,
			"planAbandon over FIFO .git = %#v, %v; want %v", got.plan, got.err, errSpecialMetadata)
	case <-time.After(abandonDeadline):
		t.Fatal("planAbandon blocked instead of refusing a special .git entry before invoking git")
	}

	info, err := os.Lstat(filepath.Join(creation.Path, ".git"))
	requireTest(t, err == nil && info.Mode()&os.ModeNamedPipe != 0,
		"the FIFO .git entry was disturbed: %v, %v", info, err)
}

func TestPlanAbandonRefusesLeftoverRegistrationMismatch(t *testing.T) {
	mismatches := []struct {
		name    string
		request string
		mutate  func(*testing.T, string, Creation)
	}{
		{name: "registration-absent", mutate: func(t *testing.T, root string, creation Creation) {
			gitRun(t, root, "worktree", "unlock", creation.Path)
			gitRun(t, root, "worktree", "prune")
		}},
		{name: "branch-ref-disagrees", mutate: func(t *testing.T, root string, creation Creation) {
			branch := strings.TrimPrefix(creation.Assignment.Branch, "refs/heads/")
			gitRun(t, root, "branch", "-M", branch, branch+"-moved")
		}},
		{name: "path-disagrees", mutate: func(t *testing.T, root string, creation Creation) {
			assignment := creation.Assignment
			assignment.Worktree = filepath.Join(root, "elsewhere")
			mustNoError(t, intent.PutAssignment(root, assignment))
		}},
		{name: "request-disagrees", request: "landed-abandon-mismatch-stranger", mutate: func(*testing.T, string, Creation) {}},
	}
	for _, mismatch := range mismatches {
		t.Run(mismatch.name, func(t *testing.T) {
			request := "landed-abandon-mismatch-" + mismatch.name
			root, creation := newOwnedAssignment(t, "abandon-mismatch-"+mismatch.name)
			decayToHusk(t, creation.Path)
			mismatch.mutate(t, root, creation)
			if mismatch.request != "" {
				request = mismatch.request
			}

			plan, err := planAbandon(root, request, creation.Path)
			requireTest(t, err == errAbandonMismatch, "%s husk plan = %#v, %v; want %v", mismatch.name, plan, err, errAbandonMismatch)
			body, readErr := os.ReadFile(filepath.Join(creation.Path, "keep.txt"))
			requireTest(t, readErr == nil && string(body) == "husk bytes\n", "refusal disturbed the husk: %q, %v", body, readErr)
		})
	}
}

func TestApplyAbandonPreservesExistingRecoveryOverHusk(t *testing.T) {
	const request = "landed-abandon-husk-recovery"
	root, creation := newOwnedAssignment(t, "abandon-husk-recovery")
	head := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	ref := intent.RecoveryRefPrefix(creation.Assignment.OwnerID, creation.Assignment.ID) + "1"
	gitRun(t, root, "update-ref", ref, head)
	assignment := creation.Assignment
	assignment.State = intent.StateCleanupPending
	assignment.Recovery = []intent.Recovery{{Ref: ref, Root: head, Payloads: []string{head}}}
	mustNoError(t, intent.PutAssignment(root, assignment))
	decayToHusk(t, creation.Path)

	plan, err := planAbandon(root, request, creation.Path)
	mustNoError(t, err)
	requireTest(t, plan.Action == actionReleaseLeftover && plan.Recovery == ref,
		"husk plan with existing recovery = %#v", plan)
	requireTest(t, gitOutput(t, root, "rev-parse", "--verify", ref+"^{commit}") == head,
		"planning moved the recovery ref off %s", head)

	result, err := ApplyAbandon(root, request, creation.Path, plan.Fingerprint)
	requireTest(t, err == nil && result.Action == ActionRemoved, "ApplyAbandon over a husk with recovery = %#v, %v", result, err)
	requireTest(t, gitOutput(t, root, "rev-parse", "--verify", ref+"^{commit}") == head,
		"recovery ref no longer resolves to %s after apply", head)
	registrations := gitOutput(t, root, "worktree", "list", "--porcelain")
	requireTest(t, !strings.Contains(registrations, "worktree "+creation.Path),
		"registration survived the husk abandon:\n%s", registrations)
	body, readErr := os.ReadFile(filepath.Join(creation.Path, "keep.txt"))
	requireTest(t, readErr == nil && string(body) == "husk bytes\n", "husk content after apply = %q, %v", body, readErr)
}

func TestLeftoverAbandonFingerprintBindsPathAndAction(t *testing.T) {
	const request = "landed-abandon-husk-fingerprint"
	root, creation := newOwnedAssignment(t, "abandon-husk-fingerprint")
	decayToHusk(t, creation.Path)

	plan, err := planAbandon(root, request, creation.Path)
	mustNoError(t, err)
	same := leftoverFingerprint(plan.leftover, plan.registration, *plan.assignment, plan.Action, plan.Recovery)
	otherPath := leftoverFingerprint(plan.leftover+"-elsewhere", plan.registration, *plan.assignment, plan.Action, plan.Recovery)
	otherAction := leftoverFingerprint(plan.leftover, plan.registration, *plan.assignment, ActionRemove, plan.Recovery)
	requireTest(t, same == plan.Fingerprint, "leftover fingerprint = %q, want the plan's %q", same, plan.Fingerprint)
	requireTest(t, otherPath != plan.Fingerprint, "leftover path is not folded into the fingerprint: %q", otherPath)
	requireTest(t, otherAction != plan.Fingerprint, "action is not folded into the fingerprint: %q", otherAction)
}
