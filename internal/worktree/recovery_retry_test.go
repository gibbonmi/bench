package worktree

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

func TestExplicitRetryFinalizesRecoveryAfterCleanDrift(t *testing.T) {
	t.Parallel()
	for _, afterRemoval := range []bool{false, true} {
		t.Run(fmt.Sprintf("after-removal=%t", afterRemoval), func(t *testing.T) {
			f := newOwnedAssignment(t, fmt.Sprintf("recovery-ref-clean-drift-%t", afterRemoval))
			dirty := filepath.Join(f.creation.Path, "dirty.txt")
			mustWrite(t, dirty, []byte("preserve once\n"), 0o644)
			first, err := PlanExplicit(f.root, f.creation.Path)
			mustNoError(t, err)
			repo, target, err := cleanupIdentity(f.root, f.creation.Path)
			mustNoError(t, err)

			stop := errors.New("stop before recovery ref creation")
			faulted := defaultJoins()
			faulted.cleanupBoundary = failLifecycleStep(StepRecoveryMetadata, stop)
			_, err = applyExplicitWith(faulted, f.ambient(), f.root, f.creation.Path, first.Fingerprint, CleanupOptions{})
			requireTest(t, errors.Is(err, stop), "first apply error = %v, want %v", err, stop)

			pending, err := assignmentByID(f.root, f.creation.Assignment.ID)
			mustNoError(t, err)
			requireTest(t, pending.State == intent.StateCleanupPending && len(pending.Recovery) == 1,
				"interrupted assignment = %#v", pending)
			recovery := pending.Recovery[0]
			requireTest(t, descendant(t, "git", "-C", f.root, "show-ref", "--verify", "--quiet", recovery.Ref).Run() != nil,
				"recovery ref exists before retry: %s", recovery.Ref)
			registration := gitOutput(t, f.root, "worktree", "list", "--porcelain")
			requireTest(t, strings.Contains(registration, "worktree "+f.creation.Path) && strings.Contains(registration, "locked "+lockReason(pending)),
				"interrupted checkout is not locked:\n%s", registration)

			mustRemove(t, dirty)
			retry, err := PlanExplicit(f.root, f.creation.Path)
			mustNoError(t, err)
			requireTest(t, retry.Action == ActionRecoverRemove && retry.Recovery == recovery.Ref,
				"clean-drift retry plan = %#v", retry)
			if afterRemoval {
				stop = errors.New("stop after worktree removal")
				faulted.cleanupBoundary = failLifecycleStep(StepRemoval, stop)
				_, err = applyExplicitWith(faulted, f.ambient(), f.root, f.creation.Path, retry.Fingerprint, CleanupOptions{})
				requireTest(t, errors.Is(err, stop), "removal interruption = %v, want %v", err, stop)
				interrupted, readErr := assignmentByID(f.root, f.creation.Assignment.ID)
				requireTest(t, readErr == nil && interrupted.State == intent.StateCleanupPending && len(interrupted.Recovery) == 1,
					"post-removal interrupted assignment = %#v, %v", interrupted, readErr)
			}
			result, err := ApplyExplicit(f.root, f.creation.Path, retry.Fingerprint)
			requireTest(t, err == nil && result.Action == ActionRemoved, "clean-drift retry = %#v, %v", result, err)

			_, statErr := os.Stat(f.creation.Path)
			requireTest(t, errors.Is(statErr, os.ErrNotExist), "retry left checkout: %v", statErr)
			resolved := gitOutput(t, f.root, "rev-parse", "--verify", recovery.Ref+"^{commit}")
			requireTest(t, resolved == recovery.Root, "recovery ref = %s, want %s", resolved, recovery.Root)
			refs := strings.Fields(gitOutput(t, f.root, "for-each-ref", "--format=%(refname)", intent.RecoveryRefPrefix(f.creation.Assignment.OwnerID, f.creation.Assignment.ID)))
			requireTest(t, len(refs) == 1 && refs[0] == recovery.Ref, "recovery refs = %#v, want only %s", refs, recovery.Ref)
			final, err := assignmentByID(f.root, f.creation.Assignment.ID)
			mustNoError(t, err)
			requireTest(t, final.State == intent.StateRecovered && len(final.Recovery) == 1 && final.Recovery[0].Ref == recovery.Ref && final.Recovery[0].Root == recovery.Root,
				"final recovered assignment = %#v", final)
			receipt, found, err := intent.CleanupReceiptFor(f.root, repo, cleanupOperation, target, retry.Fingerprint)
			requireTest(t, err == nil && found && receipt.State == intent.ReceiptComplete && receipt.Phase == intent.ReceiptPhaseTerminal && receipt.Recovery == recovery.Ref,
				"final cleanup receipt = %#v, found=%t error=%v", receipt, found, err)
		})
	}
	markProof(t, "lifecycle/journey/recovery")
}
