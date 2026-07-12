package worktree

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
