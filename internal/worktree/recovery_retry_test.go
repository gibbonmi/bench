package worktree

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

func TestExplicitRetryFinalizesRecoveryAfterCleanDrift(t *testing.T) {
	for _, afterRemoval := range []bool{false, true} {
		t.Run(fmt.Sprintf("after-removal=%t", afterRemoval), func(t *testing.T) {
			root, creation := newOwnedAssignment(t, fmt.Sprintf("recovery-ref-clean-drift-%t", afterRemoval))
			dirty := filepath.Join(creation.Path, "dirty.txt")
			mustWrite(t, dirty, []byte("preserve once\n"), 0o644)
			first, err := PlanExplicit(root, creation.Path)
			mustNoError(t, err)
			repo, target, err := cleanupIdentity(root, creation.Path)
			mustNoError(t, err)

			stop := errors.New("stop before recovery ref creation")
			old := cleanupTransactionBoundary
			cleanupTransactionBoundary = failLifecycleStep(StepRecoveryMetadata, stop)
			_, err = ApplyExplicit(root, creation.Path, first.Fingerprint)
			cleanupTransactionBoundary = old
			requireTest(t, errors.Is(err, stop), "first apply error = %v, want %v", err, stop)

			pending, err := assignmentByID(root, creation.Assignment.ID)
			mustNoError(t, err)
			requireTest(t, pending.State == intent.StateCleanupPending && len(pending.Recovery) == 1,
				"interrupted assignment = %#v", pending)
			recovery := pending.Recovery[0]
			requireTest(t, exec.Command("git", "-C", root, "show-ref", "--verify", "--quiet", recovery.Ref).Run() != nil,
				"recovery ref exists before retry: %s", recovery.Ref)
			registration := gitOutput(t, root, "worktree", "list", "--porcelain")
			requireTest(t, strings.Contains(registration, "worktree "+creation.Path) && strings.Contains(registration, "locked "+lockReason(pending)),
				"interrupted checkout is not locked:\n%s", registration)

			mustRemove(t, dirty)
			retry, err := PlanExplicit(root, creation.Path)
			mustNoError(t, err)
			requireTest(t, retry.Action == ActionRecoverRemove && retry.Recovery == recovery.Ref,
				"clean-drift retry plan = %#v", retry)
			if afterRemoval {
				stop = errors.New("stop after worktree removal")
				cleanupTransactionBoundary = failLifecycleStep(StepRemoval, stop)
				_, err = ApplyExplicit(root, creation.Path, retry.Fingerprint)
				cleanupTransactionBoundary = old
				requireTest(t, errors.Is(err, stop), "removal interruption = %v, want %v", err, stop)
				interrupted, readErr := assignmentByID(root, creation.Assignment.ID)
				requireTest(t, readErr == nil && interrupted.State == intent.StateCleanupPending && len(interrupted.Recovery) == 1,
					"post-removal interrupted assignment = %#v, %v", interrupted, readErr)
			}
			result, err := ApplyExplicit(root, creation.Path, retry.Fingerprint)
			requireTest(t, err == nil && result.Action == ActionRemoved, "clean-drift retry = %#v, %v", result, err)

			_, statErr := os.Stat(creation.Path)
			requireTest(t, errors.Is(statErr, os.ErrNotExist), "retry left checkout: %v", statErr)
			resolved := gitOutput(t, root, "rev-parse", "--verify", recovery.Ref+"^{commit}")
			requireTest(t, resolved == recovery.Root, "recovery ref = %s, want %s", resolved, recovery.Root)
			refs := strings.Fields(gitOutput(t, root, "for-each-ref", "--format=%(refname)", intent.RecoveryRefPrefix(creation.Assignment.OwnerID, creation.Assignment.ID)))
			requireTest(t, len(refs) == 1 && refs[0] == recovery.Ref, "recovery refs = %#v, want only %s", refs, recovery.Ref)
			final, err := assignmentByID(root, creation.Assignment.ID)
			mustNoError(t, err)
			requireTest(t, final.State == intent.StateRecovered && len(final.Recovery) == 1 && final.Recovery[0].Ref == recovery.Ref && final.Recovery[0].Root == recovery.Root,
				"final recovered assignment = %#v", final)
			receipt, found, err := intent.CleanupReceiptFor(root, repo, cleanupOperation, target, retry.Fingerprint)
			requireTest(t, err == nil && found && receipt.State == intent.ReceiptComplete && receipt.Phase == intent.ReceiptPhaseTerminal && receipt.Recovery == recovery.Ref,
				"final cleanup receipt = %#v, found=%t error=%v", receipt, found, err)
		})
	}
}

func TestPlanAbandonAcceptsRemovedDirectory(t *testing.T) {
	const request = "landed-abandon-removed-directory"
	root, creation := newOwnedAssignment(t, "abandon-removed-directory")
	mustRemove(t, creation.Path)
	repo, target, err := cleanupIdentity(root, creation.Path)
	mustNoError(t, err)
	_, found, err := intent.CleanupReceiptForRequest(root, repo, cleanupOperation, target, requestDigest(request))
	requireTest(t, err == nil && !found, "removed directory already has a cleanup receipt: found=%t error=%v", found, err)

	fingerprint, err := PlanAbandon(root, request, creation.Path)
	requireTest(t, err == nil && fingerprint != "", "PlanAbandon over a removed directory = %q, %v", fingerprint, err)
}

func TestApplyAbandonCompletesForRemovedDirectory(t *testing.T) {
	const request = "landed-abandon-removed-directory-apply"
	root, creation := newOwnedAssignment(t, "abandon-removed-directory-apply")
	mustRemove(t, creation.Path)
	fingerprint, err := PlanAbandon(root, request, creation.Path)
	mustNoError(t, err)

	result, err := ApplyAbandon(root, request, creation.Path, fingerprint)
	requireTest(t, err == nil && result.Action == ActionRemoved, "ApplyAbandon over a removed directory = %#v, %v", result, err)
	registrations := gitOutput(t, root, "worktree", "list", "--porcelain")
	requireTest(t, !strings.Contains(registrations, "worktree "+creation.Path),
		"registration survived the abandon:\n%s", registrations)
}

// TestApplyAbandonRecoverRemovesDirectoryWithExistingRecovery covers planRemovedCheckout's
// other branch: a removed directory whose assignment already carries a recovery ref plans
// ActionRecoverRemove against that ref rather than a plain removal, and the ref is what
// ApplyAbandon must still be holding when the assignment lands terminal.
func TestApplyAbandonRecoverRemovesDirectoryWithExistingRecovery(t *testing.T) {
	const request = "landed-recovery-remove-existing"
	root, creation := newOwnedAssignment(t, "recovery-remove-existing")
	head := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	ref := intent.RecoveryRefPrefix(creation.Assignment.OwnerID, creation.Assignment.ID) + "1"
	gitRun(t, root, "update-ref", ref, head)
	assignment := creation.Assignment
	assignment.State = intent.StateCleanupPending
	assignment.Recovery = []intent.Recovery{{Ref: ref, Root: head, Payloads: []string{head}}}
	mustNoError(t, intent.PutAssignment(root, assignment))
	mustRemove(t, creation.Path)

	plan, err := planAbandon(root, request, creation.Path)
	mustNoError(t, err)
	requireTest(t, plan.Action == ActionRecoverRemove && plan.Recovery == ref,
		"removed-directory plan with existing recovery = %#v", plan)

	result, err := ApplyAbandon(root, request, creation.Path, plan.Fingerprint)
	requireTest(t, err == nil && result.Action == ActionRemoved, "ApplyAbandon over an existing recovery ref = %#v, %v", result, err)
	resolved := gitOutput(t, root, "rev-parse", "--verify", ref+"^{commit}")
	requireTest(t, resolved == head, "recovery ref = %s after abandon, want %s", resolved, head)
	registrations := gitOutput(t, root, "worktree", "list", "--porcelain")
	requireTest(t, !strings.Contains(registrations, "worktree "+creation.Path),
		"registration survived the abandon:\n%s", registrations)
}

func TestPlanAbandonRefusesForeignCheckout(t *testing.T) {
	const request = "landed-abandon-foreign-checkout"
	root, creation := newOwnedAssignment(t, "abandon-foreign-checkout")
	mustRemove(t, creation.Path)
	gitRun(t, filepath.Dir(creation.Path), "init", "-q", "-b", "main", creation.Path)
	gitRun(t, creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-q", "--allow-empty", "-m", "foreign")

	fingerprint, err := PlanAbandon(root, request, creation.Path)
	requireTest(t, err != nil && fingerprint == "", "PlanAbandon over a foreign checkout = %q, %v", fingerprint, err)
	requireTest(t, err.Error() == "abandon request, assignment, or path mismatch; checkout retained",
		"foreign checkout refusal = %q", err)
	_, statErr := os.Stat(filepath.Join(creation.Path, ".git"))
	requireTest(t, statErr == nil, "refusal disturbed the foreign checkout: %v", statErr)
}

func TestInterruptedReleaseStillResumesThroughReceipt(t *testing.T) {
	const request = "landed-abandon-receipt-resume"
	root, creation := newOwnedAssignment(t, "abandon-receipt-resume")
	mustWrite(t, filepath.Join(creation.Path, "dirty.txt"), []byte("preserve once\n"), 0o644)
	first, err := PlanAbandon(root, request, creation.Path)
	mustNoError(t, err)
	stop := errors.New("stop before worktree removal")
	old := cleanupTransactionBoundary
	cleanupTransactionBoundary = failLifecycleStep(StepRecoveryRef, stop)
	_, err = ApplyAbandon(root, request, creation.Path, first)
	cleanupTransactionBoundary = old
	requireTest(t, errors.Is(err, stop), "interrupted abandon error = %v, want %v", err, stop)
	mustRemove(t, creation.Path)

	repo, target, err := cleanupIdentity(root, creation.Path)
	mustNoError(t, err)
	receipt, found, err := intent.CleanupReceiptForRequest(root, repo, cleanupOperation, target, requestDigest(request))
	requireTest(t, err == nil && found && receipt.State == intent.ReceiptInFlight,
		"interrupted cleanup receipt = %#v, found=%t error=%v", receipt, found, err)
	plan, err := planAbandon(root, request, creation.Path)
	requireTest(t, err == nil && plan.Fingerprint == first, "resumed abandon plan = %#v, %v; want fingerprint %q", plan, err, first)
	requireTest(t, reflect.DeepEqual(plan, planFromReceipt(receipt)), "resumed abandon plan = %#v, want the receipt-derived %#v", plan, planFromReceipt(receipt))
}

func TestRemovedDirectoryWithHostilePathPlansAndApplies(t *testing.T) {
	const request = "landed-abandon-removed-hostile-path"
	root := newWorktreeRepo(t)
	pool := t.TempDir()
	t.Setenv("BENCH_HOME", filepath.Join(pool, "a b*c"))
	decoy := filepath.Join(pool, "a bzc", "keep.txt")
	mustNoError(t, os.MkdirAll(filepath.Dir(decoy), 0o700))
	mustWrite(t, decoy, []byte("keep\n"), 0o644)
	creation := mustCreate(t, root, request, "hostile path")
	requireTest(t, strings.Contains(creation.Path, "a b*c"), "assignment path is not hostile: %s", creation.Path)
	mustRemove(t, creation.Path)

	fingerprint, err := PlanAbandon(root, request, creation.Path)
	mustNoError(t, err)
	result, err := ApplyAbandon(root, request, creation.Path, fingerprint)
	requireTest(t, err == nil && result.Action == ActionRemoved, "hostile-path abandon = %#v, %v", result, err)
	registrations := gitOutput(t, root, "worktree", "list", "--porcelain")
	requireTest(t, !strings.Contains(registrations, "worktree "+creation.Path),
		"registration survived the abandon:\n%s", registrations)
	body, readErr := os.ReadFile(decoy)
	requireTest(t, readErr == nil && string(body) == "keep\n", "abandon removed a glob sibling: %q, %v", body, readErr)
}

func TestAbandonRetryUsesInFlightReceipt(t *testing.T) {
	const request = "landed-abandon-retry-in-flight"
	root, creation := newOwnedAssignment(t, "abandon-retry-in-flight")
	mustWrite(t, filepath.Join(creation.Path, "dirty.txt"), []byte("preserve once\n"), 0o644)
	first, err := PlanAbandon(root, request, creation.Path)
	mustNoError(t, err)
	stop := errors.New("stop before worktree removal")
	old := cleanupTransactionBoundary
	cleanupTransactionBoundary = failLifecycleStep(StepRemoval, stop)
	_, err = ApplyAbandon(root, request, creation.Path, first)
	cleanupTransactionBoundary = old
	requireTest(t, errors.Is(err, stop), "first abandon apply error = %v, want %v", err, stop)

	pending, err := assignmentByID(root, creation.Assignment.ID)
	mustNoError(t, err)
	requireTest(t, pending.State == intent.StateCleanupPending && len(pending.Recovery) == 1,
		"interrupted assignment = %#v", pending)
	retry, err := PlanAbandon(root, request, creation.Path)
	requireTest(t, err == nil && retry == first, "abandon retry fingerprint = %q, %v; want %q", retry, err, first)
	result, err := ApplyAbandon(root, request, creation.Path, retry)
	requireTest(t, err == nil && result.Action == ActionRemoved, "abandon retry = %#v, %v", result, err)

	if _, err := os.Stat(creation.Path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("retry left checkout: %v", err)
	}
	recovery := pending.Recovery[0]
	refs := strings.Fields(gitOutput(t, root, "for-each-ref", "--format=%(refname)", intent.RecoveryRefPrefix(creation.Assignment.OwnerID, creation.Assignment.ID)))
	requireTest(t, len(refs) == 1 && refs[0] == recovery.Ref, "recovery refs = %#v, want only %s", refs, recovery.Ref)
	repo, target, err := cleanupIdentity(root, creation.Path)
	mustNoError(t, err)
	receipt, found, err := intent.CleanupReceiptFor(root, repo, cleanupOperation, target, first)
	requireTest(t, err == nil && found && receipt.State == intent.ReceiptComplete && receipt.Recovery == recovery.Ref,
		"final cleanup receipt = %#v, found=%t error=%v", receipt, found, err)
}
