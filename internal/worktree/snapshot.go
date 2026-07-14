package worktree

import (
	"errors"
	"fmt"
	"github.com/gibbonmi/bench/internal/git"
	"os"
	"os/exec"
	"strings"
)

// SnapshotDirty preserves a worktree's full dirty tree (staged, unstaged, and
// untracked, nested included) as a durable commit on ref, parented on parent.
// It is built entirely from plumbing under a temp index with an explicit
// synthetic identity, so a missing Git identity or a hostile commit hook
// cannot block preservation the way an ordinary commit would. exclude names
// worktree-relative paths to drop from the snapshot (the caller owns that
// policy — this primitive stays generic). Ref creation fails closed: if ref
// already exists, SnapshotDirty returns an error and leaves it untouched.
func SnapshotDirty(worktreePath, parent, ref string, exclude []string) (string, error) {
	admin, err := git.Output("-C", worktreePath, "rev-parse", "--path-format=absolute", "--git-dir")
	if err != nil {
		return "", fmt.Errorf("resolve snapshot admin dir: %w", err)
	}
	index, cleanup, err := temporaryIndex(admin)
	if err != nil {
		return "", fmt.Errorf("create snapshot index: %w", err)
	}
	defer cleanup()
	env := []string{"GIT_INDEX_FILE=" + index}
	if _, err := gitInput(worktreePath, env, nil, "read-tree", parent); err != nil {
		return "", fmt.Errorf("snapshot read-tree: %w", err)
	}
	if _, err := gitInput(worktreePath, env, nil, "add", "-A"); err != nil {
		return "", fmt.Errorf("snapshot add: %w", err)
	}
	for _, path := range exclude {
		args := []string{"rm", "--cached", "--ignore-unmatch", "--", ":(literal)" + path}
		if _, err := gitInput(worktreePath, env, nil, args...); err != nil {
			return "", fmt.Errorf("snapshot exclude %s: %w", path, err)
		}
	}
	tree, err := gitInput(worktreePath, env, nil, "write-tree")
	if err != nil {
		return "", fmt.Errorf("snapshot write-tree: %w", err)
	}
	commit, err := commitTree(worktreePath, tree, []string{parent}, "bench shift recovery snapshot\n")
	if err != nil {
		return "", fmt.Errorf("snapshot commit-tree: %w", err)
	}
	zero := strings.Repeat("0", len(commit))
	if out, err := exec.Command("git", "-C", worktreePath, "update-ref", ref, commit, zero).CombinedOutput(); err != nil {
		return "", fmt.Errorf("create recovery ref: %s", strings.TrimSpace(string(out)))
	}
	resolved, err := git.Output("-C", worktreePath, "rev-parse", "--verify", ref+"^{commit}")
	if err != nil || resolved != commit {
		return "", fmt.Errorf("recovery ref failed verification")
	}
	return commit, nil
}

// RetainAndLock is the fallback preservation path used only when SnapshotDirty
// itself fails: it locks the charged worktree in place with reason and drops
// its pool lease file WITHOUT restoring cleanliness, so the dirty work stays
// on disk untouched. No reset, no clean, no release.
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
