package worktree

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
)

// RetainAndLock is the preservation path a shift takes when its charged worktree is
// dirty: it locks the worktree in place with reason and drops its pool lease file
// WITHOUT restoring cleanliness, so the dirty work stays on disk untouched. No reset,
// no clean, no release.
func RetainAndLock(worktreePath, reason string) error {
	if out, err := exec.Command("git", "-C", worktreePath, "worktree", "lock", "--reason", reason, worktreePath).CombinedOutput(); err != nil {
		return fmt.Errorf("lock worktree for retention: %s", strings.TrimSpace(string(out)))
	}
	lease, err := LeaseFile(worktreePath)
	if err != nil {
		return fmt.Errorf("resolve lease for retention: %w", err)
	}
	if err := os.Remove(lease); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("drop lease for retention: %w", err)
	}
	return nil
}

// leaseEnd ends every lease line that leaseLine writes.
const leaseEnd = "\n"

// ClaimRecordedLease takes wt's lease through the pool's takeover protocol when the
// lease file holds exactly the recorded line. It returns false when the claim concedes.
func ClaimRecordedLease(wt, recorded string) bool {
	lease, err := LeaseFile(wt)
	if err != nil {
		return false
	}
	return claimRecordedLease(defaultJoins(), lease, recorded)
}

// claimRecordedLease is ClaimRecordedLease on a lease path with the seam set given.
// The recorded line is the one accepted judgment, so a lease another writer holds, or
// one that changes in the takeover gap, concedes.
func claimRecordedLease(j joins, lease, recorded string) bool {
	return claimAt(j, lease, currentTime(), func(content []byte, _, _ time.Time) bool {
		return string(content) == recorded+leaseEnd
	})
}

// ReadLease grades wt's lease file without a follow. It returns the lease line without
// its final leaseEnd and the line's owner. present is false when the file is absent, and
// ok is false for a special, unreadable, or malformed file. The policy parse accepts
// only one line that ends in leaseEnd.
func ReadLease(wt string) (line string, owner int, present, ok bool) {
	lease, err := LeaseFile(wt)
	if err != nil {
		return "", 0, false, false
	}
	read := bounds.ClassifyNoFollow(lease)
	if read.State == bounds.StateAbsent {
		return "", 0, false, false
	}
	if read.State != bounds.StateParsed {
		return "", 0, true, false
	}
	owner, ok = leaseOwnerPID(read.Data)
	return strings.TrimSuffix(string(read.Data), leaseEnd), owner, true, ok
}
