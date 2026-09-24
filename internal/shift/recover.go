package shift

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/otelrecord"
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
		line, owner, present, ok := worktree.ReadLease(entry.Worktree)
		if present && ok && line == entry.Lease && !worktree.PIDAlive(owner) {
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
	line, owner, present, ok := worktree.ReadLease(entry.Worktree)
	switch {
	case !present:
		return abandon
	case !ok, worktree.PIDAlive(owner):
		return keep
	case line != entry.Lease:
		return abandon
	}
	return recoverShift
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
func finishRecovery(root string, entry intent.Entry, memory ...attribute.KeyValue) bool {
	kind, _ := splitRecovery(entry.Recovery)
	cleanup := otelrecord.CleanupReleased
	if kind == recoveryWorktreeKind {
		cleanup = otelrecord.CleanupRetained
		// A worktree that is already locked keeps its lock, and the pass drops its lease.
		if worktree.RetainAndLock(entry.Worktree, "bench shift recovery") != nil {
			if lease, err := worktree.LeaseFile(entry.Worktree); err == nil {
				_ = os.Remove(lease)
			}
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

// recordRecovery writes one shift.recovery span for a verdict on the entry key.
func recordRecovery(root, key, state, cleanup string, extra ...attribute.KeyValue) {
	attrs := append([]attribute.KeyValue{
		attribute.String(otelrecord.AttrIntentKey, key),
		attribute.String(otelrecord.AttrWorkState, state),
		attribute.String(otelrecord.AttrCleanup, cleanup)}, extra...)
	_, _, end := otelrecord.BeginIn(context.Background(), "", root, recoverySeam, recoverySeam, attrs...)
	end()
}
