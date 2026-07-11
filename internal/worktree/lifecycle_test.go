package worktree

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestAcquireCreatesPrivatePoolAndLease(t *testing.T) {
	oldUmask := syscall.Umask(0)
	t.Cleanup(func() { syscall.Umask(oldUmask) })

	home := t.TempDir()
	t.Setenv("BENCH_HOME", home)
	root := newWorktreeRepo(t)
	wt, err := Acquire(root, "", "")
	if err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	t.Cleanup(func() { Release(wt) })

	assertMode := func(path string, want os.FileMode) {
		t.Helper()
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("stat %s: %v", path, err)
		}
		if got := info.Mode().Perm(); got != want {
			t.Errorf("mode %s = %04o, want %04o", path, got, want)
		}
	}
	assertMode(Pool(root), 0o700)
	lease, err := LeaseFile(wt)
	if err != nil {
		t.Fatalf("LeaseFile: %v", err)
	}
	assertMode(lease, 0o600)
}

func TestAcquireTightensExistingPool(t *testing.T) {
	home := t.TempDir()
	t.Setenv("BENCH_HOME", home)
	root := newWorktreeRepo(t)
	pool := Pool(root)
	if err := os.MkdirAll(pool, 0o777); err != nil {
		t.Fatalf("mkdir loose pool: %v", err)
	}
	if err := os.Chmod(pool, 0o777); err != nil {
		t.Fatalf("chmod loose pool: %v", err)
	}

	wt, err := Acquire(root, "", "")
	if err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	t.Cleanup(func() { Release(wt) })
	info, err := os.Stat(pool)
	if err != nil {
		t.Fatalf("stat pool: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o700 {
		t.Fatalf("pool mode after Acquire = %04o, want 0700", got)
	}
}

func TestAcquireContinuesWhenPoolTightenFails(t *testing.T) {
	home := t.TempDir()
	t.Setenv("BENCH_HOME", home)
	root := newWorktreeRepo(t)
	pool := Pool(root)
	old := chmodPool
	called := false
	chmodPool = func(path string, mode os.FileMode) error {
		if path == pool {
			called = true
			return os.ErrPermission
		}
		return os.Chmod(path, mode)
	}
	t.Cleanup(func() { chmodPool = old })

	wt, err := Acquire(root, "", "")
	if err != nil {
		t.Fatalf("Acquire after pool chmod failure: %v", err)
	}
	t.Cleanup(func() { Release(wt) })
	if !called {
		t.Fatal("Acquire did not attempt to tighten the pool")
	}
}

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

// deadPidLine returns a lease line whose recorded pid is provably dead: a child
// process is spawned and reaped, so its pid is gone (kill -0 → ESRCH) yet was a real
// pid, exactly the crashed-owner case reclaimable takes over.
func deadPidLine(t *testing.T) string {
	t.Helper()
	cmd := exec.Command("true")
	if err := cmd.Run(); err != nil {
		t.Fatalf("spawn reap victim: %v", err)
	}
	pid := cmd.Process.Pid
	if pidAlive(pid) {
		t.Skip("reaped pid reused before use")
	}
	return fmt.Sprintf("%d 2026-07-05T00:00:00Z\n", pid)
}

// TestClaimSecondReclaimerConcedes pins the takeover identity check: two reclaimers of
// the same dead-pid lease must not both win one worktree. The claimTakeoverGap seam
// drives the second (outer) reclaimer to pause after judging the lease reclaimable, and
// runs a first (nested) reclaimer to a full takeover — installing a fresh live lease —
// before the outer's rename runs. The pre-fix takeover renames (and so steals) that
// fresh lease and re-creates it, returning true: both reclaimers win, falsifying the
// "cannot both win" guarantee. The fix sees the renamed bytes differ from the bytes it
// judged reclaimable, restores the fresh lease, and concedes.
func TestClaimSecondReclaimerConcedes(t *testing.T) {
	if os.Getpid() == 1 {
		t.Skip("running as pid 1")
	}
	lease := filepath.Join(t.TempDir(), "bench-lease")
	if err := os.WriteFile(lease, []byte(deadPidLine(t)), 0o644); err != nil {
		t.Fatal(err)
	}

	real := claimTakeoverGap
	t.Cleanup(func() { claimTakeoverGap = real })
	var nestedWon, reentered bool
	claimTakeoverGap = func(leasePath string) {
		if reentered {
			return // the nested reclaimer's own gap is a no-op
		}
		reentered = true
		nestedWon = Claim(leasePath)
	}

	outerWon := Claim(lease)

	if !nestedWon {
		t.Error("first (nested) reclaimer did not win the dead-pid lease")
	}
	if outerWon {
		t.Error("second reclaimer also won: two reclaimers both took over one worktree")
	}
	// The surviving lease is the nested winner's fresh lease: present and owned by us.
	got, err := os.ReadFile(lease)
	if err != nil {
		t.Fatalf("lease missing after takeover: %v", err)
	}
	if want := fmt.Sprintf("%d ", os.Getpid()); !strings.HasPrefix(string(got), want) {
		t.Errorf("lease = %q, want fresh winner prefix %q", got, want)
	}
	// No .stale.* leftover beside the lease.
	entries, err := os.ReadDir(filepath.Dir(lease))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.Contains(e.Name(), ".stale.") {
			t.Errorf("stale leftover remains beside lease: %s", e.Name())
		}
	}
}

// TestClaimStealDuringTakeoverKeepsFirstWriter pins the three-party takeover race: a
// concurrent reclaimer B wins the whole reclaim before the outer caller's rename runs,
// and a fresh first-writer C then claims the slot the outer's rename vacates. The
// invariant is that C's lease survives: a first-writer that claims a vacated slot
// during a conceded takeover keeps its lease. Pre-fix, the identity check's blind
// rename-back clobbers C's lease with the stolen bytes it meant to restore; the fix
// restores only into a still-empty slot (no-clobber link), leaving C's lease alone.
func TestClaimStealDuringTakeoverKeepsFirstWriter(t *testing.T) {
	if os.Getpid() == 1 {
		t.Skip("running as pid 1")
	}
	lease := filepath.Join(t.TempDir(), "bench-lease")
	if err := os.WriteFile(lease, []byte(deadPidLine(t)), 0o644); err != nil {
		t.Fatal(err)
	}

	realTakeover := claimTakeoverGap
	t.Cleanup(func() { claimTakeoverGap = realTakeover })
	var nestedWon bool
	inNestedTakeover := false
	claimTakeoverGap = func(lp string) {
		if inNestedTakeover {
			return // B's own pass through this gap must not recurse again
		}
		inNestedTakeover = true
		nestedWon = Claim(lp)
		inNestedTakeover = false
	}

	realSteal := claimStealGap
	t.Cleanup(func() { claimStealGap = realSteal })
	const sentinel = "999999 sentinel-first-writer\n"
	wrote := false
	claimStealGap = func(lp string) {
		// Skip B's own pass (still inside the nested takeover) so the write lands only
		// in the slot the outer's rename vacates, not the one B's rename vacates.
		if inNestedTakeover || wrote {
			return
		}
		wrote = true
		if err := os.WriteFile(lp, []byte(sentinel), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	outerWon := Claim(lease)

	if outerWon {
		t.Error("outer Claim = true, want false (it must concede to the first-writer)")
	}
	if !nestedWon {
		t.Error("nested reclaimer B did not win its takeover")
	}
	got, err := os.ReadFile(lease)
	if err != nil {
		t.Fatalf("lease missing after conceded takeover: %v", err)
	}
	if string(got) != sentinel {
		t.Errorf("lease = %q, want first-writer sentinel %q — the restore clobbered it", got, sentinel)
	}
	entries, err := os.ReadDir(filepath.Dir(lease))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.Contains(e.Name(), ".stale.") {
			t.Errorf("stale leftover remains beside lease: %s", e.Name())
		}
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
