package worktree

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
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
