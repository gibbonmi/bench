package shift

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/otelrecord"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
	"github.com/gibbonmi/bench/internal/worktree"
	"go.opentelemetry.io/otel/attribute"
)

// recoverySeam is the seam of a recovery verdict's span.
const recoverySeam = "shift.recovery"

// verdict is what the recovery pass does with one shift entry.
type verdict int

const (
	keep verdict = iota
	abandon
	recoverShift
	resumeRecovery
)

// Recover is the recovery pass over root's shift entries. An open entry of a gone owner
// is abandoned or recovered by its lease identity, and a recovery that stopped before its
// act is finished. Every other entry stays as it is. The pass prints one line only when it
// acted, so a quiet session start stays quiet.
func Recover(root string, out io.Writer) {
	current, err := intent.Read(root)
	if err != nil {
		return
	}
	recovered, abandoned := 0, 0
	for _, entry := range current.Entries {
		if entry.Kind != intent.KindShift {
			continue
		}
		switch judge(entry) {
		case abandon:
			if abandonEntry(root, entry) {
				abandoned++
			}
		case recoverShift:
			if recoverEntry(root, entry) {
				recovered++
			}
		case resumeRecovery:
			if finishRecovery(root, entry) {
				recovered++
			}
		}
	}
	if recovered+abandoned > 0 {
		fmt.Fprintf(out, "bench shift recovery: recovered %d, abandoned %d\n", recovered, abandoned)
	}
}

// judge gives one shift entry its verdict. The lease is graded before the comparison, so
// a malformed or special lease never counts as another owner.
func judge(entry intent.Entry) verdict {
	switch entry.Outcome {
	case "":
	case otelrecord.WorkRecovered:
		// A recovery that wrote its entry and stopped before its act still holds its lease.
		if _, _, held, gone := leaseProof(entry); held && gone {
			return resumeRecovery
		}
		return keep
	default:
		return keep
	}
	if entry.Lease == "" {
		if pid, ok := intent.KeyOwner(intent.KindShift, entry.Key); ok && !worktree.PIDAlive(pid) {
			return abandon
		}
		return keep
	}
	present, ok, held, gone := leaseProof(entry)
	switch {
	case !present:
		return abandon
	case !ok, !gone:
		return keep
	case !held:
		return abandon
	}
	return recoverShift
}

// leaseProof grades the lease of entry's worktree. A special, unreadable, or malformed
// lease is present and not ok. held reports that it holds the entry's recorded line, and
// gone reports that its owner process is gone.
func leaseProof(entry intent.Entry) (present, ok, held, gone bool) {
	line, owner, present, ok := worktree.ReadLease(entry.Worktree)
	return present, ok, ok && line == entry.Lease, ok && !worktree.PIDAlive(owner)
}

// abandonEntry closes an entry whose owner is gone. It changes no worktree file, lease,
// or lock.
func abandonEntry(root string, entry intent.Entry) bool {
	entry.Outcome = otelrecord.WorkAbandoned
	if intent.Upsert(root, entry) != nil {
		return false
	}
	recordRecovery(root, entry.Key, otelrecord.WorkAbandoned, otelrecord.CleanupNone)
	return true
}

// recoverEntry recovers a crashed shift. It keeps the notes first, because a concurrent
// acquire of a clean tree removes them, and then takes the lease through the pool's
// takeover protocol. A lost claim changes nothing more.
func recoverEntry(root string, entry intent.Entry) bool {
	memory := keepNotes(root, entry.Worktree, entry.Key)
	if hitShift(shiftFault, stepRecoveryClaim) != nil || !worktree.ClaimRecordedLease(entry.Worktree, entry.Lease) {
		return false
	}
	entry.Lease, _, _, _ = worktree.ReadLease(entry.Worktree)
	entry.Outcome, entry.Recovery = otelrecord.WorkRecovered, RecoveryNone
	if len(dirtyPaths(entry.Worktree)) > 0 {
		entry.Recovery = recoveryWorktree(entry.Worktree)
	}
	if intent.Upsert(root, entry) != nil || hitShift(shiftFault, stepRecoveryAct) != nil {
		return false
	}
	return finishRecovery(root, entry, memory...)
}

// finishRecovery acts on a claimed worktree as preserveAndRecover does, writes the entry
// again so an abandon that landed between the writes loses, and records the recovery.
// Under the claimed lease, it rechecks the users of a worktree that the plan would
// release. A user that remains unresolved retains the worktree and its scratch, even when
// the dirty set is empty.
func finishRecovery(root string, entry intent.Entry, memory ...attribute.KeyValue) bool {
	kind, _ := splitRecovery(entry.Recovery)
	if kind != recoveryWorktreeKind && worktree.CheckoutUnprotected(entry.Worktree) != nil {
		entry.Recovery, kind = recoveryWorktree(entry.Worktree), recoveryWorktreeKind
	}
	cleanup := otelrecord.CleanupReleased
	if kind == recoveryWorktreeKind {
		cleanup = otelrecord.CleanupRetained
		if lockedTree(root, entry.Worktree) {
			// A worktree that is already locked keeps its lock, and the pass drops its lease.
			if lease, err := worktree.LeaseFile(entry.Worktree); err == nil {
				_ = os.Remove(lease)
			}
		} else {
			// A failed lock keeps the lease, so a later pass resumes the act.
			_ = worktree.RetainAndLock(entry.Worktree, "bench shift recovery")
		}
	} else {
		cleanupScratch(entry.Worktree)
		worktree.Release(entry.Worktree)
	}
	if intent.Upsert(root, entry) != nil {
		return false
	}
	recordRecovery(root, entry.Key, otelrecord.WorkRecovered, cleanup, memory...)
	return true
}

// lockedTree reports whether git lists wt as a locked worktree of root.
func lockedTree(root, wt string) bool {
	listed, _ := git.Worktrees(root)
	for _, tree := range listed {
		if tree.Path == wt {
			return tree.Locked
		}
	}
	return false
}

// recordRecovery writes one shift.recovery span for a verdict on the entry key.
func recordRecovery(root, key, state, cleanup string, extra ...attribute.KeyValue) {
	attrs := append([]attribute.KeyValue{
		attribute.String(otelrecord.AttrIntentKey, key),
		attribute.String(otelrecord.AttrWorkState, state),
		attribute.String(otelrecord.AttrCleanup, cleanup)}, extra...)
	_, _, end := otelrecord.BeginIn(context.Background(), "", root, recoverySeam, recoverySeam, attrs...)
	end()
}

// teardown removes the shift scratch and releases the pool lease, exactly once whether
// reached by the normal path, the deferred cleanup, or preserveAndRecover. It returns
// the first error hit: an unresolved user of the worktree, which keeps the scratch and
// the lease, or the injectable teardown fault, since worktree.Release does not itself
// report git failures. A real teardown problem is reported rather than silently
// swallowed. A second call, the deferred safety net once preserve is already set, is a
// no-op via sync.Once and returns nil.
func (s *session) teardown() error {
	var err error
	s.teardownOne.Do(func() {
		if err = worktree.CheckoutUnprotected(s.root); err != nil {
			return
		}
		cleanupScratch(s.root)
		if ferr := hitShift(shiftFault, stepTeardown); ferr != nil {
			err = ferr
			return
		}
		worktree.Release(s.root)
		s.released = true
	})
	return err
}

// retainNotes keeps the notes of a shift that reached its first iteration. Every exit
// after the first iteration runs preserveAndRecover, which calls it before any scratch
// cleanup and on the retain path. An earlier exit has no notes to keep.
func (s *session) retainNotes() {
	if s.iterationsUsed > 0 {
		s.record.retainNotes(s.root)
	}
}

// preserveAndRecover is the one place a post-mutation failure preserves work before this
// session's charged worktree leaves the process's hands. When nothing beyond scratch is
// dirty and no user of the worktree remains unresolved, it releases through teardown and
// reports RecoveryNone. Otherwise it retains and locks the worktree, leaving the tree and
// its scratch exactly where they are and never running Release. An empty dirty set grants
// no removal while an adapter group can still change the worktree. The work stays visible
// at a path the operator can read, rather than in a ref no command hands back. It always
// prints the resulting location and marks the session so the deferred cleanup never
// re-finalizes this worktree.
func (s *session) preserveAndRecover(reason string) (recovery string, teardownErr error) {
	defer s.preserve.Store(true)
	s.retainNotes()
	if len(dirtyPaths(s.root)) == 0 {
		err := s.teardown()
		if !errors.Is(err, protection.ErrRetained) {
			return RecoveryNone, err
		}
		fmt.Fprintln(s.stderr, err)
	}
	if lockErr := worktree.RetainAndLock(s.root, "bench shift recovery: "+reason); lockErr != nil {
		fmt.Fprintf(s.stderr, "could not retain worktree %s: %v\n", s.root, lockErr)
	}
	fmt.Fprintf(s.stdout, "  recovery: %s\n", recoveryWorktree(s.root))
	return recoveryWorktree(s.root), nil
}

// teardownFailureResult is the one Result a teardown error resolves to, regardless of
// what post-mutation path reached it. Its outcome is failed/1, with a detail naming
// what is already safe, the branch, and the recovery pointer when one exists. The
// teardown failure is real even when the work is not lost.
func teardownFailureResult(s *session, recovery string, err error) Result {
	detail := fmt.Sprintf("teardown failed: %v; branch %s is safe", err, s.branch)
	if intent.HoldsRecovery(recovery) {
		detail += "; recovery " + recovery
	}
	return Result{Outcome: OutcomeFailed, Branch: s.branch, Committed: s.committed, IterationsUsed: s.iterationsUsed, Recovery: recovery, Detail: detail}
}

// rollback discards a red iteration's work while preserving the shift scratch, so the
// next iteration still reads what the last one learned. An unresolved user of the
// worktree refuses it before any reset or clean, because that user can still change the
// evidence that the rollback would discard.
func rollback(root string) error {
	if err := worktree.CheckoutUnprotected(root); err != nil {
		return err
	}
	_ = exec.Command("git", "-C", root, "reset", "-q", "--hard").Run()
	_ = exec.Command("git", "-C", root, "clean", "-qfdx", "-e", ".bench-objective", "-e", ".bench-notes.md").Run()
	return nil
}
