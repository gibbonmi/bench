package worktree

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestReclaimable pins the four-way lease decision the black-box lease-hardening
// contract exercises but cannot cheaply enumerate: a recorded numeric pid gates on
// liveness, and unreadable/empty content reclaims only once aged past the threshold —
// so a fresh-empty writer mid-claim is never stolen while a legacy/crashed lease is.
func TestReclaimable(t *testing.T) {
	now := time.Date(2026, 7, 4, 12, 0, 0, 0, time.UTC)
	dead := func(int) bool { return false }
	live := func(int) bool { return true }

	cases := []struct {
		name    string
		content string
		age     time.Duration
		alive   func(int) bool
		want    bool
	}{
		{"live pid respected", "4242 2026-07-04T11:59:00Z", 30 * time.Second, live, false},
		{"dead pid reclaimed", "4242 2026-07-04T11:59:00Z", 30 * time.Second, dead, true},
		{"non-numeric legacy aged out reclaimed", "garbage content", 2 * time.Minute, dead, true},
		{"fresh empty writer mid-claim respected", "", 5 * time.Second, dead, false},
		{"aged-out empty crash reclaimed", "", 2 * time.Minute, dead, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mtime := now.Add(-tc.age)
			if got := reclaimable([]byte(tc.content), mtime, now, tc.alive); got != tc.want {
				t.Errorf("reclaimable(%q, age %s) = %v, want %v", tc.content, tc.age, got, tc.want)
			}
		})
	}
}

// leasedRepo builds a one-commit repo whose lease file carries the given content,
// returning the repo dir and the lease path. It is the fixture both Release contracts
// share: a tracked file to dirty, an untracked file to strand, and a lease to honor.
func leasedRepo(t *testing.T, leaseContent string) (dir, lease string) {
	t.Helper()
	dir = t.TempDir()
	for _, args := range [][]string{
		{"init", "-q"},
		{"config", "user.email", "t@t"},
		{"config", "user.name", "t"},
	} {
		if out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "tracked.txt"), []byte("clean\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"add", "."}, {"commit", "-qm", "init"}} {
		if out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
	}
	lease, err := LeaseFile(dir)
	if err != nil {
		t.Fatalf("LeaseFile: %v", err)
	}
	if err := os.WriteFile(lease, []byte(leaseContent), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir, lease
}

// dirty modifies the tracked file and adds an untracked one, so a release contract can
// assert both the reset and the clean happened.
func dirty(t *testing.T, dir string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "tracked.txt"), []byte("dirty\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "untracked.txt"), []byte("stray\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestReleaseOwnerRestoresCleanAndUnleases pins the owner path: after Release, the
// worktree is back to a reusable clean state and the lease is gone — the pool-entry
// invariant that an unleased entry is always claimably clean.
func TestReleaseOwnerRestoresCleanAndUnleases(t *testing.T) {
	dir, lease := leasedRepo(t, fmt.Sprintf("%d 2026-07-05T00:00:00Z\n", os.Getpid()))
	dirty(t, dir)

	Release(dir)

	if _, err := os.Stat(lease); !os.IsNotExist(err) {
		t.Errorf("lease still present after owner release: %v", err)
	}
	if !isClean(dir) {
		out, _ := exec.Command("git", "-C", dir, "status", "--porcelain").CombinedOutput()
		t.Errorf("worktree dirty after release:\n%s", out)
	}
	if got, err := os.ReadFile(filepath.Join(dir, "tracked.txt")); err != nil || string(got) != "clean\n" {
		t.Errorf("tracked.txt = %q, %v; want restored %q", got, err, "clean\n")
	}
}

// TestReleaseRespectsLiveForeignLease pins the non-owner path: a lease held by a
// different live process means the entry was stale-reclaimed and belongs to that
// owner, so Release must leave both the lease and the working state untouched.
// Pid 1 is live for any test runner (kill -0 yields nil or EPERM, both alive).
func TestReleaseRespectsLiveForeignLease(t *testing.T) {
	if os.Getpid() == 1 {
		t.Skip("running as pid 1")
	}
	dir, lease := leasedRepo(t, "1 2026-07-05T00:00:00Z\n")
	dirty(t, dir)

	Release(dir)

	if _, err := os.Stat(lease); err != nil {
		t.Errorf("foreign live lease removed: %v", err)
	}
	if got, err := os.ReadFile(filepath.Join(dir, "untracked.txt")); err != nil || string(got) != "stray\n" {
		t.Errorf("new owner's untracked file lost: %q, %v", got, err)
	}
	if got, err := os.ReadFile(filepath.Join(dir, "tracked.txt")); err != nil || string(got) != "dirty\n" {
		t.Errorf("new owner's edit reverted: %q, %v", got, err)
	}
}

// TestReleaseNeverClaimableMidCleanup pins the release ordering: at the moment the
// cleanup step runs, the owner's lease must still be held, so a concurrent Claim
// fails. This is the whole race contract — if the entry is never claimable before
// the single final lease removal, no claimant can win mid-cleanup (the dirty window)
// and no lease exists afterwards for a trailing remove to delete. The restoreClean
// seam is the interleave point; the pre-fix ordering (unlease, clean, unlease again)
// goes red here because the simulated claimant's create succeeds mid-cleanup.
func TestReleaseNeverClaimableMidCleanup(t *testing.T) {
	dir, lease := leasedRepo(t, fmt.Sprintf("%d 2026-07-05T00:00:00Z\n", os.Getpid()))
	dirty(t, dir)

	real := restoreClean
	t.Cleanup(func() { restoreClean = real })
	claimedMidCleanup := false
	restoreClean = func(wt string) {
		claimedMidCleanup = Claim(lease)
		real(wt)
	}

	Release(dir)

	if claimedMidCleanup {
		t.Error("entry was claimable mid-cleanup: a concurrent claimant can win while the worktree is dirty and lose its lease to the trailing remove")
	}
	if _, err := os.Stat(lease); !os.IsNotExist(err) {
		t.Errorf("lease still present after owner release: %v", err)
	}
}

// TestCandidateNameStaysInPool pins that a minted candidate never escapes the pool
// directory — a wrong name would mint outside the pool and silently break warm reuse.
func TestCandidateNameStaysInPool(t *testing.T) {
	pool := "/home/x/.bench/worktrees/bench-123"
	got := candidateName(pool, 1751630400, 4242, 2)
	if filepath.Dir(got) != pool {
		t.Errorf("candidateName parent = %q, want %q", filepath.Dir(got), pool)
	}
	if want := filepath.Join(pool, "1751630400-4242-2"); got != want {
		t.Errorf("candidateName = %q, want %q", got, want)
	}
	if strings.ContainsAny(filepath.Base(got), "/") {
		t.Errorf("candidate base %q escaped the pool", filepath.Base(got))
	}
}
