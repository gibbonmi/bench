package worktree

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
)

// mustAcquire leases a fresh linked pool worktree from root. It matches the
// shift-charged worktree shape: a real `git worktree add`, a lease file at
// LeaseFile, reset and cleaned.
func mustAcquire(t *testing.T, root, home string) string {
	t.Helper()
	wt, err := acquireAt(defaultJoins(), root, "", "", home, currentTime())
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
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	worktreePath := mustAcquire(t, root, home)

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

// recordedLeaseFile writes a lease that holds recorded and returns its path.
func recordedLeaseFile(t *testing.T, recorded string) string {
	t.Helper()
	lease := filepath.Join(t.TempDir(), "bench-lease")
	if err := os.WriteFile(lease, []byte(recorded+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return lease
}

// The identity claim takes a lease that holds the recorded line, and the lease then names
// this process.
func TestClaimRecordedLeaseTakesTheRecordedLine(t *testing.T) {
	t.Parallel()
	recorded := "4242 2026-07-05T00:00:00Z"
	lease := recordedLeaseFile(t, recorded)
	if !claimRecordedLease(defaultJoins(), lease, recorded) {
		t.Fatal("the claim of the recorded line conceded")
	}
	if got, _ := os.ReadFile(lease); !bytes.HasPrefix(got, []byte(fmt.Sprintf("%d ", os.Getpid()))) {
		t.Fatalf("lease = %q, want this process as the owner", got)
	}
}

// The identity claim refuses a lease that holds another line, and the lease stays.
func TestClaimRecordedLeaseRefusesAnotherLine(t *testing.T) {
	t.Parallel()
	lease := recordedLeaseFile(t, "4343 2026-07-05T00:00:01Z")
	if claimRecordedLease(defaultJoins(), lease, "4242 2026-07-05T00:00:00Z") {
		t.Fatal("the claim took a lease that holds another line")
	}
	if got, _ := os.ReadFile(lease); string(got) != "4343 2026-07-05T00:00:01Z\n" {
		t.Fatalf("lease = %q, want the other line unchanged", got)
	}
}

// LE97: a writer that replaces the lease in the takeover gap keeps it, and the identity
// claim concedes.
func TestClaimRecordedLeaseConcedesToAWriterInTheGap(t *testing.T) {
	t.Parallel()
	recorded := "4242 2026-07-05T00:00:00Z"
	lease := recordedLeaseFile(t, recorded)
	other := []byte("4343 2026-07-05T00:00:01Z\n")
	j := defaultJoins()
	j.claimTakeoverGap = func(path string) {
		if err := os.WriteFile(path, other, 0o600); err != nil {
			t.Error(err)
		}
	}
	if claimRecordedLease(j, lease, recorded) {
		t.Fatal("the claim won over a writer in the takeover gap")
	}
	if got, _ := os.ReadFile(lease); !bytes.Equal(got, other) {
		t.Fatalf("lease = %q, want the other writer's lease %q", got, other)
	}
}
