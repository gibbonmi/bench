package worktree

import (
	"github.com/gibbonmi/bench/internal/git"
	"os"
	"path/filepath"
	"testing"
)

// mustAcquire leases a fresh linked pool worktree from root. It matches the
// shift-charged worktree shape: a real `git worktree add`, a lease file at
// LeaseFile, reset and cleaned.
func mustAcquire(t *testing.T, root string) string {
	t.Helper()
	wt, err := Acquire(root, "", "")
	if err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	return wt
}

// TestRetainAndLockLocksDropsLeaseAndPreservesDirt pins the shift's preservation
// path. After RetainAndLock, `git worktree list --porcelain` reports the path
// locked with the given reason, and the pool lease file is gone. The dirty
// file on disk stays untouched: no reset, no clean, no release.
func TestRetainAndLockLocksDropsLeaseAndPreservesDirt(t *testing.T) {
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	worktreePath := mustAcquire(t, root)

	lease, err := LeaseFile(worktreePath)
	if err != nil {
		t.Fatalf("LeaseFile: %v", err)
	}
	if _, err := os.Stat(lease); err != nil {
		t.Fatalf("fixture lease missing before RetainAndLock: %v", err)
	}
	mustWrite(t, filepath.Join(worktreePath, "dirty.txt"), []byte("keep me\n"), 0o644)

	const reason = "bench shift recovery: gate red at the deadline"
	if err := RetainAndLock(worktreePath, reason); err != nil {
		t.Fatalf("RetainAndLock: %v", err)
	}

	worktrees, err := git.Worktrees(root)
	if err != nil {
		t.Fatalf("git.Worktrees: %v", err)
	}
	var found *git.Worktree
	for i := range worktrees {
		if samePath(worktrees[i].Path, worktreePath) {
			found = &worktrees[i]
			break
		}
	}
	if found == nil {
		t.Fatalf("worktree %s not found in worktree list: %#v", worktreePath, worktrees)
	}
	if !found.Locked {
		t.Errorf("worktree %s not locked after RetainAndLock", worktreePath)
	}
	if found.LockReason != reason {
		t.Errorf("lock reason = %q, want %q", found.LockReason, reason)
	}

	if _, err := os.Stat(lease); !os.IsNotExist(err) {
		t.Errorf("lease still present after RetainAndLock: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(worktreePath, "dirty.txt"))
	if err != nil {
		t.Fatalf("dirty file gone after RetainAndLock: %v", err)
	}
	if string(got) != "keep me\n" {
		t.Errorf("dirty file content = %q, want %q", got, "keep me\n")
	}
}
