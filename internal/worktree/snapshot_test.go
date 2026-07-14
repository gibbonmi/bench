package worktree

import (
	"github.com/gibbonmi/bench/internal/git"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// unsetEnv removes each key for the test's duration, restoring any prior value
// on cleanup. Used to prove the snapshot primitive needs no ambient Git
// identity — a missing GIT_AUTHOR_NAME must not block preservation.
func unsetEnv(t *testing.T, keys ...string) {
	t.Helper()
	for _, key := range keys {
		old, had := os.LookupEnv(key)
		if err := os.Unsetenv(key); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if had {
				_ = os.Setenv(key, old)
			}
		})
	}
}

// mustAcquire leases a fresh linked pool worktree from root, matching the
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

// dirtySnapshotFixture builds a linked worktree (via the real pool Acquire
// path, so it matches production shape) containing: a tracked file modified after
// checkout, a hostile-named untracked file, a nested untracked path, and two
// scratch files that the caller (internal/shift) would exclude by name. It
// returns the worktree path and the commit its branch currently points at,
// the parent SnapshotDirty should be built on.
func dirtySnapshotFixture(t *testing.T) (worktreePath, parent string) {
	t.Helper()
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	worktreePath = mustAcquire(t, root)
	parent = gitOutput(t, worktreePath, "rev-parse", "HEAD")

	mustWrite(t, filepath.Join(worktreePath, "tracked.txt"), []byte("modified\n"), 0o644)
	mustWrite(t, filepath.Join(worktreePath, "step 1 [a].txt"), []byte("hostile\n"), 0o644)
	mustMkdirAll(t, filepath.Join(worktreePath, "sub", "dir"), 0o755)
	mustWrite(t, filepath.Join(worktreePath, "sub", "dir", "new.txt"), []byte("nested\n"), 0o644)
	mustWrite(t, filepath.Join(worktreePath, ".bench-objective"), []byte("do the thing\n"), 0o644)
	mustWrite(t, filepath.Join(worktreePath, ".bench-notes.md"), []byte("notes\n"), 0o644)
	return worktreePath, parent
}

// TestSnapshotDirtyCapturesUntrackedNestedAndExcludesScratch pins row 1,18: the
// snapshot's tree contains the hostile-named and nested untracked paths (proving
// `add -A` under the temp index captures untracked/nested files in a linked
// worktree, not just tracked changes), excludes the caller's scratch names, is
// parented on the given parent, and is authored under the synthetic bench/
// bench@local identity — all without any ambient Git identity present.
func TestSnapshotDirtyCapturesUntrackedNestedAndExcludesScratch(t *testing.T) {
	unsetEnv(t, "GIT_AUTHOR_NAME", "GIT_AUTHOR_EMAIL", "GIT_COMMITTER_NAME", "GIT_COMMITTER_EMAIL")
	worktreePath, parent := dirtySnapshotFixture(t)
	ref := "refs/bench/recovery/snapshot-fixture"
	exclude := []string{".bench-objective", ".bench-notes.md"}

	commit, err := SnapshotDirty(worktreePath, parent, ref, exclude)
	if err != nil {
		t.Fatalf("SnapshotDirty: %v", err)
	}
	if commit == "" {
		t.Fatal("SnapshotDirty returned an empty commit oid")
	}

	tree := gitOutput(t, worktreePath, "ls-tree", "-r", "--name-only", commit)
	paths := strings.Split(tree, "\n")
	want := map[string]bool{"step 1 [a].txt": false, "sub/dir/new.txt": false}
	unwanted := map[string]bool{".bench-objective": true, ".bench-notes.md": true}
	for _, p := range paths {
		if _, ok := want[p]; ok {
			want[p] = true
		}
		if unwanted[p] {
			t.Errorf("snapshot tree contains excluded scratch path %q", p)
		}
	}
	for p, found := range want {
		if !found {
			t.Errorf("snapshot tree missing %q; got tree:\n%s", p, tree)
		}
	}

	gotParent := gitOutput(t, worktreePath, "rev-parse", commit+"^")
	if gotParent != parent {
		t.Errorf("snapshot parent = %q, want %q", gotParent, parent)
	}

	body := gitOutput(t, worktreePath, "show", "-s", "--format=%an <%ae>%n%cn <%ce>", commit)
	lines := strings.Split(body, "\n")
	if len(lines) != 2 {
		t.Fatalf("unexpected commit header %q", body)
	}
	if lines[0] != "bench <bench@local>" {
		t.Errorf("author = %q, want %q", lines[0], "bench <bench@local>")
	}
	if lines[1] != "bench <bench@local>" {
		t.Errorf("committer = %q, want %q", lines[1], "bench <bench@local>")
	}

	resolved := gitOutput(t, worktreePath, "rev-parse", "--verify", ref+"^{commit}")
	if resolved != commit {
		t.Errorf("ref %s resolves to %q, want %q", ref, resolved, commit)
	}
}

// TestSnapshotDirtyFailsClosedOnExistingRef pins row 3: a pre-existing ref at
// the target name must never be overwritten — SnapshotDirty returns an error
// and the ref keeps resolving to its original commit, so a timestamp-collision
// branch name cannot clobber another shift's evidence.
func TestSnapshotDirtyFailsClosedOnExistingRef(t *testing.T) {
	worktreePath, parent := dirtySnapshotFixture(t)
	ref := "refs/bench/recovery/snapshot-fixture"

	gitRun(t, worktreePath, "update-ref", ref, parent)
	original := gitOutput(t, worktreePath, "rev-parse", "--verify", ref+"^{commit}")

	_, err := SnapshotDirty(worktreePath, parent, ref, nil)
	if err == nil {
		t.Fatal("SnapshotDirty succeeded against a pre-existing ref, want an error")
	}

	resolved := gitOutput(t, worktreePath, "rev-parse", "--verify", ref+"^{commit}")
	if resolved != original {
		t.Errorf("pre-existing ref moved from %q to %q after a failed snapshot", original, resolved)
	}
}

// TestRetainAndLockLocksDropsLeaseAndPreservesDirt pins the fallback path: after
// RetainAndLock, `git worktree list --porcelain` reports the path locked with
// the given reason, the pool lease file is gone, and the dirty file on disk is
// untouched — no reset, no clean, no release.
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

	const reason = "bench shift recovery: snapshot failed"
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
