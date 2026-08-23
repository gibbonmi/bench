package worktree

import (
	"fmt"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/intent"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestReclaimable pins the four-way lease decision the black-box lease-hardening
// contract exercises but cannot cheaply enumerate. A canonical lease gates on pid
// liveness, and unreadable/empty content reclaims only once aged past the threshold.
// So a fresh-empty writer mid-claim is never stolen while a legacy/crashed lease is.

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
		{"live pid respected", "4242 2026-07-04T11:59:00Z\n", bounds.LeaseStale / 2, live, false},
		{"dead pid reclaimed", "4242 2026-07-04T11:59:00Z\n", bounds.LeaseStale / 2, dead, true},
		{"non-numeric legacy aged out reclaimed", "garbage content", 2 * bounds.LeaseStale, dead, true},
		{"fresh empty writer mid-claim respected", "", bounds.LeaseStale / 2, dead, false},
		{"empty lease at exactly the stale window respected", "", bounds.LeaseStale, dead, false},
		{"aged-out empty crash reclaimed", "", 2 * bounds.LeaseStale, dead, true},
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
		if out, err := descendant(t, "git", append([]string{"-C", dir}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "tracked.txt"), []byte("clean\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"add", "."}, {"commit", "-qm", "init"}} {
		if out, err := descendant(t, "git", append([]string{"-C", dir}, args...)...).CombinedOutput(); err != nil {
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

// TestReleaseOwnerRestoresCleanAndUnleases pins the owner path. After Release, the
// worktree is back to a reusable clean state and the lease is gone. This is the
// pool-entry invariant: an unleased entry is always claimably clean.

func TestReleaseOwnerRestoresCleanAndUnleases(t *testing.T) {
	dir, lease := leasedRepo(t, fmt.Sprintf("%d 2026-07-05T00:00:00Z\n", os.Getpid()))
	dirty(t, dir)
	Release(dir)
	if _, err := os.Stat(lease); !os.IsNotExist(err) {
		t.Errorf("lease still present after owner release: %v", err)
	}
	if !isClean(dir) {
		out, _ := descendant(t, "git", "-C", dir, "status", "--porcelain").CombinedOutput()
		t.Errorf("worktree dirty after release:\n%s", out)
	}
	if got, err := os.ReadFile(filepath.Join(dir, "tracked.txt")); err != nil || string(got) != "clean\n" {
		t.Errorf("tracked.txt = %q, %v; want restored %q", got, err, "clean\n")
	}
	markProof(t, "lifecycle/journey/create-remove")
}

// TestReleaseRespectsLiveForeignLease pins the non-owner path. A lease held by a
// different live process means the entry was stale-reclaimed and belongs to that
// owner. Release must then leave both the lease and the working state untouched.
// Pid 1 is live for any test runner (kill -0 yields nil or EPERM, both alive).

func TestReleaseRespectsLiveForeignLease(t *testing.T) {
	if os.Getpid() == 1 {
		capability.Capability(t, capability.PID, "running as pid 1")
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
// fails. This is the whole race contract. If the entry is never claimable before
// the single final lease removal, no claimant can win mid-cleanup (the dirty window).
// And no lease exists afterwards for a trailing remove to delete. The restoreClean
// seam is the interleave point; the pre-fix ordering (unlease, clean, unlease again)
// goes red here because the simulated claimant's create succeeds mid-cleanup.

func TestReleaseNeverClaimableMidCleanup(t *testing.T) {
	dir, lease := leasedRepo(t, fmt.Sprintf("%d 2026-07-05T00:00:00Z\n", os.Getpid()))
	dirty(t, dir)
	real := restoreClean
	t.Cleanup(func() { restoreClean = real })
	claimedMidCleanup := false
	restoreClean = func(wt string) {
		claimedMidCleanup = claimAt(lease, time.Now())
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

// deadPidLine returns a lease line whose recorded pid is provably dead. A child
// process is spawned and reaped, so its pid is gone (kill -0 → ESRCH) yet was a
// real pid. This is exactly the crashed-owner case reclaimable takes over.
func deadPidLine(t *testing.T) string {
	t.Helper()
	cmd := descendant(t, "true")
	if err := cmd.Run(); err != nil {
		t.Fatalf("spawn reap victim: %v", err)
	}
	pid := cmd.Process.Pid
	if pidAlive(pid) {
		capability.Capability(t, capability.PID, "reaped pid reused before use")
	}
	return fmt.Sprintf("%d 2026-07-05T00:00:00Z\n", pid)
}

// TestClaimSecondReclaimerConcedes pins the takeover identity check: two reclaimers of
// the same dead-pid lease must not both win one worktree. The claimTakeoverGap seam
// drives the second (outer) reclaimer to pause after judging the lease reclaimable.
// It then runs a first (nested) reclaimer to a full takeover, installing a fresh live
// lease, before the outer's rename runs.
//
// A blind takeover renames (and so steals) that fresh lease and re-creates it,
// returning true. Both reclaimers then win and falsify the "cannot both win"
// guarantee. The fix sees the renamed bytes differ from the bytes it judged
// reclaimable, restores the fresh lease, and concedes.

func TestClaimSecondReclaimerConcedes(t *testing.T) {
	if os.Getpid() == 1 {
		capability.Capability(t, capability.PID, "running as pid 1")
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
		nestedWon = claimAt(leasePath, time.Now())
	}
	outerWon := claimAt(lease, time.Now())
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

// TestClaimStealDuringTakeoverKeepsFirstWriter pins the three-party takeover race. A
// concurrent reclaimer B wins the whole reclaim before the outer caller's rename runs.
// A fresh first-writer C then claims the slot the outer's rename vacates. The
// invariant is that C's lease survives: a first-writer that claims a vacated slot
// during a conceded takeover keeps its lease.
//
// A blind identity check clobbers C's lease with the stolen bytes it meant to
// restore. The fix restores only into a still-empty slot (no-clobber link),
// leaving C's lease alone.

func TestClaimStealDuringTakeoverKeepsFirstWriter(t *testing.T) {
	if os.Getpid() == 1 {
		capability.Capability(t, capability.PID, "running as pid 1")
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
		nestedWon = claimAt(lp, time.Now())
		inNestedTakeover = false
	}
	realSteal := claimStealGap
	t.Cleanup(func() { claimStealGap = realSteal })
	const sentinel = "999999 sentinel-first-writer\n"
	wrote := false
	claimStealGap = func(lp string) {
		// Skip B's own pass (still inside the nested takeover). The write must land only
		// in the slot the outer's rename vacates, not the one B's rename vacates.
		if inNestedTakeover || wrote {
			return
		}
		wrote = true
		if err := os.WriteFile(lp, []byte(sentinel), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	outerWon := claimAt(lease, time.Now())
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

func TestIgnoredInventoryEntryAndByteBoundaries(t *testing.T) {
	for _, count := range []int{0, 1, 20, 21, 1000, 1001} {
		t.Run(fmt.Sprintf("entries-%d", count), func(t *testing.T) {
			root := newWorktreeRepo(t)
			gitRun(t, root, "branch", "-M", "main")
			mustWrite(t, filepath.Join(root, ".git", "info", "exclude"), []byte("ignored-*\n"), 0o644)
			target := filepath.Join(filepath.Dir(root), fmt.Sprintf("ignored-%d", count))
			gitRun(t, root, "worktree", "add", "-q", "-b", fmt.Sprintf("ignored-%d", count), target, "HEAD")
			for i := 0; i < count; i++ {
				mustWrite(t, filepath.Join(target, fmt.Sprintf("ignored-%04d", i)), []byte("x"), 0o600)
			}
			plan, err := PlanExplicitWithOptions(root, target, CleanupOptions{DiscardIgnored: true})
			if err != nil {
				t.Fatal(err)
			}
			wantAction := ActionDiscardRemove
			if count == 0 {
				wantAction = ActionRemove
			} else if count > ignoredEntryLimit {
				wantAction = ActionRetain
			}
			if plan.Action != wantAction || plan.Ignored.Count != count {
				t.Fatalf("plan = %#v, want action %s count %d", plan, wantAction, count)
			}
			wantShown := count
			if wantShown > 20 {
				wantShown = 20
			}
			if plan.Ignored.Shown != wantShown || plan.Ignored.Truncated != (count > 20) {
				t.Fatalf("preview = %#v", plan.Ignored)
			}
			full, err := PlanExplicitWithOptions(root, target, CleanupOptions{DiscardIgnored: true, Full: true})
			if err != nil || full.Fingerprint != plan.Fingerprint {
				t.Fatalf("--full changed plan identity: %#v, %v", full, err)
			}
		})
	}
	for _, size := range []int64{ignoredByteLimit - 1, ignoredByteLimit + 1} {
		t.Run(fmt.Sprintf("bytes-%d", size), func(t *testing.T) {
			root := newWorktreeRepo(t)
			mustWrite(t, filepath.Join(root, ".git", "info", "exclude"), []byte("large.bin\n"), 0o644)
			target := filepath.Join(filepath.Dir(root), fmt.Sprintf("bytes-%d", size))
			gitRun(t, root, "worktree", "add", "-q", "-b", fmt.Sprintf("bytes-%d", size), target, "HEAD")
			large := filepath.Join(target, "large.bin")
			mustWrite(t, large, nil, 0o600)
			if err := os.Truncate(large, size); err != nil {
				t.Fatal(err)
			}
			plan, err := PlanExplicitWithOptions(root, target, CleanupOptions{DiscardIgnored: true})
			if err != nil {
				t.Fatal(err)
			}
			if (size < ignoredByteLimit) != (plan.Action == ActionDiscardRemove) {
				t.Fatalf("byte-boundary plan = %#v", plan)
			}
		})
	}
}
func TestReleaseReconcilesCompletedAutomaticCleanup(t *testing.T) {
	root, creation := newPendingAssignment(t, "release-crash-window")
	plan, err := ApplyAutomatic(root, creation.Path, nil)
	requireTest(t, err == nil && plan.Action == ActionRemoved, "automatic cleanup = %#v, %v", plan, err)
	args := []string{"--request", "landed-release-crash-window", creation.Path}
	var first, firstErr strings.Builder
	code := ReleaseCommand(root, args, &first, &firstErr)
	requireTest(t, code == 0 && firstErr.String() == "", "release reconciliation code=%d stderr=%q", code, firstErr.String())
	var replay, replayErr strings.Builder
	code = ReleaseCommand(root, args, &replay, &replayErr)
	requireTest(t, code == 0 && replay.String() == first.String() && replayErr.String() == "", "release replay code=%d stdout=%q stderr=%q", code, replay.String(), replayErr.String())
	repo, _, _ := cleanupIdentity(root, creation.Path)
	_, found, err := intent.CleanupReceiptFor(root, repo, releaseOperation, creation.Path, intent.RequestDigest("landed-release-crash-window"))
	requireTest(t, err == nil && found, "release receipt missing: %v", err)
}

// TestReleaseUnmergedAssignmentRetains pins CO4: ReleaseCommand decides purely through the
// automatic verdict (PlanAutomatic), not a second, independently-authored policy. An
// unlanded assignment branch is a retain PlanAutomatic itself decides, so release must
// honor it rather than removing. The assignment must land in cleanup-pending (the
// automatic transaction's own state, not a bespoke pre-transition refusal) so a retry
// replans through the same verdict.
func TestReleaseUnmergedAssignmentRetains(t *testing.T) {
	root, creation := newOwnedAssignment(t, "unmerged-release")
	commitInWorktree(t, creation.Path, "unique.txt", "preserve\n", "unique work")

	var stdout, stderr strings.Builder
	code := ReleaseCommand(root, []string{"--request", "landed-unmerged-release", creation.Path}, &stdout, &stderr)
	requireTest(t, code == 1, "unmerged release exit=%d stderr=%q", code, stderr.String())
	requireTest(t, strings.Contains(stderr.String(), "worktree retained (unmerged)") && strings.Contains(stderr.String(), "assignment branch has not landed"),
		"unmerged reason missing: %q", stderr.String())
	_, statErr := os.Stat(creation.Path)
	requireTest(t, statErr == nil, "unmerged worktree removed: %v", statErr)
	assignment, err := assignmentByID(root, creation.Assignment.ID)
	mustNoError(t, err)
	requireTest(t, assignment.State == intent.StateCleanupPending, "unmerged assignment state = %q, want cleanup-pending", assignment.State)
}

func mustRemove(t *testing.T, path string) {
	t.Helper()
	if err := os.RemoveAll(path); err != nil {
		t.Fatal(err)
	}
}
func mustChmod(t *testing.T, path string, mode os.FileMode) {
	t.Helper()
	if err := os.Chmod(path, mode); err != nil {
		t.Fatal(err)
	}
}
