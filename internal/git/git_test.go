package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/gittest"
)

// runGit runs `git -C root <args>` and fails the test on a nonzero exit.
func runGit(t *testing.T, root string, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", append([]string{"-C", root}, args...)...).Output()
	if err != nil {
		t.Fatalf("git %v: %v", args, err)
	}
	return string(out)
}

func worktreesWithin(t *testing.T, root string) ([]Worktree, error) {
	t.Helper()
	type result struct {
		worktrees []Worktree
		err       error
	}
	done := make(chan result, 1)
	go func() {
		worktrees, err := Worktrees(root)
		done <- result{worktrees: worktrees, err: err}
	}()
	select {
	case result := <-done:
		return result.worktrees, result.err
	case <-time.After(bounds.TestDeadline(0)):
		t.Fatal("git.Worktrees blocked on worktree admin entry")
		return nil, nil
	}
}

func requireAdminRefusal(t *testing.T, err error, path, shape string) {
	t.Helper()
	var got *WorktreeAdminError
	if !errors.As(err, &got) || got.Path != path || got.Shape != shape || got.Action != "inspect and remove it" || !strings.Contains(err.Error(), path) {
		t.Fatalf("admin refusal = %v, want path=%q shape=%q", err, path, shape)
	}
}

func TestCommonDirReturnsUnvalidatedOutput(t *testing.T) {
	root := newRepo(t)
	gittest.StubGit(t, root, "bad-rev-parse", filepath.Join(t.TempDir(), "argv"))
	want := filepath.Join(root, "missing-common")
	got, err := CommonDir(root)
	if err != nil || got != want {
		t.Fatalf("CommonDir = %q, %v, want %q, nil", got, err, want)
	}
}

func TestCommonDirKeepsPlainOutputFailure(t *testing.T) {
	root := newRepo(t)
	gittest.StubGit(t, root, "fail-rev-parse", filepath.Join(t.TempDir(), "argv"))
	_, err := CommonDir(root)
	var exitErr *exec.ExitError
	var resolution *ResolutionError
	if !errors.As(err, &exitErr) || errors.As(err, &resolution) {
		t.Fatalf("CommonDir failure = %T %v, want plain git.Output exit error", err, err)
	}
}

func TestScanWorktreeAdminAcceptanceShapes(t *testing.T) {
	root := newRepo(t)
	common, err := CommonDir(root)
	if err != nil {
		t.Fatal(err)
	}
	base := filepath.Join(common, "worktrees")
	admin := filepath.Join(base, "x y*")
	if err := os.MkdirAll(admin, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(filepath.Join(admin, "gitdir"), 0o600); err != nil {
		t.Fatal(err)
	}
	var got *WorktreeAdminError
	if err := ScanWorktreeAdmin(common); !errors.As(err, &got) || got.Shape != "fifo" || !strings.Contains(got.Error(), "worktrees/x y*/gitdir") || !strings.Contains(got.Error(), "inspect and remove it") {
		t.Fatalf("fifo refusal = %v", err)
	}
	_ = os.Remove(filepath.Join(admin, "gitdir"))
	if err := os.Symlink("target", filepath.Join(admin, "gitdir")); err != nil {
		t.Fatal(err)
	}
	if err := ScanWorktreeAdmin(common); !errors.As(err, &got) || got.Shape != "symlink" {
		t.Fatalf("symlink refusal = %v", err)
	}
	_ = os.Remove(filepath.Join(admin, "gitdir"))
	if err := os.Symlink("lease-target", filepath.Join(admin, "bench-lease")); err != nil {
		t.Fatal(err)
	}
	if err := ScanWorktreeAdmin(common); !errors.As(err, &got) || got.Shape != "symlink" || !strings.Contains(got.Error(), "bench-lease") {
		t.Fatalf("Bench lease refusal = %v", err)
	}
	_ = os.Remove(filepath.Join(admin, "bench-lease"))
	if err := os.WriteFile(filepath.Join(admin, "gitdir"), []byte("ok"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(admin, "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(admin, "nested", "deep"), []byte("ok"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := ScanWorktreeAdmin(common); err != nil {
		t.Fatalf("regular/deep state refused: %v", err)
	}
	if err := os.RemoveAll(base); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(base, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := ScanWorktreeAdmin(common); err != nil {
		t.Fatalf("fifo worktrees should be absent: %v", err)
	}
}

func TestWorktreesRefusesMalformedAdminShapes(t *testing.T) {
	t.Run("FIFO gitdir", func(t *testing.T) {
		root := newRepo(t)
		common, err := CommonDir(root)
		if err != nil {
			t.Fatal(err)
		}
		id := filepath.Join(common, "worktrees", "fifo")
		if err := os.MkdirAll(id, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := syscall.Mkfifo(filepath.Join(id, "gitdir"), 0o600); err != nil {
			capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
		}
		_, err = worktreesWithin(t, root)
		requireAdminRefusal(t, err, "worktrees/fifo/gitdir", "fifo")
	})
	t.Run("first-level symlink", func(t *testing.T) {
		root := newRepo(t)
		common, err := CommonDir(root)
		if err != nil {
			t.Fatal(err)
		}
		target := filepath.Join(t.TempDir(), "target")
		if err := os.Mkdir(target, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(common, "worktrees"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, filepath.Join(common, "worktrees", "linked")); err != nil {
			capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable: %v", err))
		}
		_, err = worktreesWithin(t, root)
		requireAdminRefusal(t, err, "worktrees/linked", "symlink")
	})
	t.Run("first-level FIFO", func(t *testing.T) {
		root := newRepo(t)
		common, err := CommonDir(root)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(common, "worktrees"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := syscall.Mkfifo(filepath.Join(common, "worktrees", "fifo"), 0o600); err != nil {
			capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
		}
		_, err = worktreesWithin(t, root)
		requireAdminRefusal(t, err, "worktrees/fifo", "fifo")
	})
	t.Run("symlinked gitdir", func(t *testing.T) {
		root := newRepo(t)
		common, err := CommonDir(root)
		if err != nil {
			t.Fatal(err)
		}
		id := filepath.Join(common, "worktrees", "linked-gitdir")
		if err := os.MkdirAll(id, 0o755); err != nil {
			t.Fatal(err)
		}
		target := filepath.Join(t.TempDir(), "regular")
		if err := os.WriteFile(target, []byte("regular\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, filepath.Join(id, "gitdir")); err != nil {
			capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable: %v", err))
		}
		_, err = worktreesWithin(t, root)
		requireAdminRefusal(t, err, "worktrees/linked-gitdir/gitdir", "symlink")
	})
	t.Run("stray FIFO preserves hostile id", func(t *testing.T) {
		root := newRepo(t)
		common, err := CommonDir(root)
		if err != nil {
			t.Fatal(err)
		}
		id := filepath.Join(common, "worktrees", "x y*")
		if err := os.MkdirAll(id, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := syscall.Mkfifo(filepath.Join(id, "stray"), 0o600); err != nil {
			capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
		}
		_, err = worktreesWithin(t, root)
		requireAdminRefusal(t, err, "worktrees/x y*/stray", "fifo")
	})
}

func TestScanWorktreeAdminRefusesUninspectableRoot(t *testing.T) {
	common := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(common, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	err := ScanWorktreeAdmin(common)
	var got *WorktreeScanError
	if !errors.As(err, &got) || got.Path != "worktrees" || got.Action != "investigate the git failure" {
		t.Fatalf("scan error = %v", err)
	}
}

func TestWorktreesPropagatesScanTraversalFailureBeforePorcelain(t *testing.T) {
	root := newRepo(t)
	logPath := filepath.Join(t.TempDir(), "argv")
	common := gittest.StubGit(t, root, "clean", logPath)
	id := filepath.Join(common, "worktrees", "unreadable")
	if err := os.MkdirAll(id, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(id, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(id, 0o755) })
	if _, err := os.ReadDir(id); err == nil {
		capability.Capability(t, capability.Privilege, "host privileges bypass unreadable-directory traversal")
	}

	_, err := Worktrees(root)
	var got *WorktreeScanError
	if !errors.As(err, &got) || got.Path != "worktrees/unreadable" || got.Action != "investigate the git failure" {
		t.Fatalf("scan traversal failure = %v", err)
	}
	data, readErr := os.ReadFile(logPath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if strings.Contains(string(data), "worktree list") {
		t.Fatalf("porcelain invoked after scan failure: %s", data)
	}
}

func TestWorktreesAllowsBenignAdminShapes(t *testing.T) {
	t.Run("absent and FIFO worktrees root", func(t *testing.T) {
		root := newRepo(t)
		if _, err := worktreesWithin(t, root); err != nil {
			t.Fatalf("absent worktrees root: %v", err)
		}
		common, err := CommonDir(root)
		if err != nil {
			t.Fatal(err)
		}
		base := filepath.Join(common, "worktrees")
		if err := syscall.Mkfifo(base, 0o600); err != nil {
			capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
		}
		if _, err := worktreesWithin(t, root); err != nil {
			t.Fatalf("FIFO worktrees root: %v", err)
		}
	})
	for _, state := range []string{"prunable", "gitdir-less", "empty", "deep FIFO"} {
		t.Run(state, func(t *testing.T) {
			assertBenignAdminShape(t, state)
		})
	}
}

func TestWorktreesTreatsRegularFileAdminRootAsAbsent(t *testing.T) {
	root := newRepo(t)
	common, err := CommonDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(common, "worktrees"), []byte("ignored\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := worktreesWithin(t, root); err != nil {
		t.Fatalf("regular-file worktrees root: %v", err)
	}
}

func TestPruneLandedBranchesUsesNeutralDiscoveryFailure(t *testing.T) {
	root := newRepo(t)
	common, err := CommonDir(root)
	if err != nil {
		t.Fatal(err)
	}
	id := filepath.Join(common, "worktrees", "fifo")
	if err := os.MkdirAll(id, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(filepath.Join(id, "gitdir"), 0o600); err != nil {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
	}

	_, err = PruneLandedBranches(root, nil)
	if err == nil || !strings.Contains(err.Error(), "worktree discovery failed") || strings.Contains(err.Error(), "git worktree list") {
		t.Fatalf("worktree discovery refusal = %v", err)
	}
}

func TestWorktreesAcceptsRegularFirstLevelAdminEntry(t *testing.T) {
	root := newRepo(t)
	common, err := CommonDir(root)
	if err != nil {
		t.Fatal(err)
	}
	base := filepath.Join(common, "worktrees")
	if err := os.Mkdir(base, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "regular-id"), []byte("ignored\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := worktreesWithin(t, root); err != nil {
		t.Fatalf("regular first-level admin entry: %v", err)
	}
}

func assertBenignAdminShape(t *testing.T, state string) {
	t.Helper()
	root := newRepo(t)
	common, err := CommonDir(root)
	if err != nil {
		t.Fatal(err)
	}
	id := filepath.Join(common, "worktrees", "benign")
	switch state {
	case "prunable":
		linked := filepath.Join(t.TempDir(), "linked")
		runGit(t, root, "worktree", "add", "-q", "--detach", linked, "HEAD")
		if err := os.RemoveAll(linked); err != nil {
			t.Fatal(err)
		}
	case "gitdir-less", "empty":
		if err := os.MkdirAll(id, 0o755); err != nil {
			t.Fatal(err)
		}
	case "deep FIFO":
		if err := os.MkdirAll(filepath.Join(id, "logs"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := syscall.Mkfifo(filepath.Join(id, "logs", "HEAD"), 0o600); err != nil {
			capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
		}
	}
	if _, err := worktreesWithin(t, root); err != nil {
		t.Fatalf("%s state refused: %v", state, err)
	}
}

func TestWorktreesRefusesSymlinkedAdminRootBeforePorcelain(t *testing.T) {
	root := newRepo(t)
	common, err := CommonDir(root)
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(t.TempDir(), "worktrees")
	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, filepath.Join(common, "worktrees")); err != nil {
		t.Fatal(err)
	}
	var got *WorktreeAdminError
	if _, err := Worktrees(root); !errors.As(err, &got) || got.Shape != "symlink" {
		t.Fatalf("symlink root refusal = %v", err)
	}
}

func TestWorktreesRefusesSharedAdminFromHostileLinkedRoot(t *testing.T) {
	root := newRepo(t)
	linked := filepath.Join(t.TempDir(), "linked root [*];$(nope)")
	runGit(t, root, "worktree", "add", "-q", "--detach", linked, "HEAD")
	common, err := CommonDir(root)
	if err != nil {
		t.Fatal(err)
	}
	id := filepath.Join(common, "worktrees", "linked-fixture")
	if err := os.MkdirAll(id, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(filepath.Join(id, "gitdir"), 0o600); err != nil {
		t.Fatal(err)
	}
	var got *WorktreeAdminError
	if _, err := Worktrees(linked); !errors.As(err, &got) || got.Shape != "fifo" {
		t.Fatalf("linked refusal = %v", err)
	}
}

func TestWorktreesRefusesMalformedAdminBeforePorcelain(t *testing.T) {
	root := newRepo(t)
	logPath := filepath.Join(t.TempDir(), "argv")
	gittest.StubGit(t, root, "block-worktree", logPath)
	id := filepath.Join(root, ".git", "worktrees", "before-porcelain")
	if err := os.MkdirAll(id, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(filepath.Join(id, "gitdir"), 0o600); err != nil {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
	}
	done := make(chan error, 1)
	go func() {
		_, err := Worktrees(root)
		done <- err
	}()
	select {
	case err := <-done:
		var got *WorktreeAdminError
		if !errors.As(err, &got) || !strings.Contains(err.Error(), "worktrees/before-porcelain/gitdir") || got.Shape != "fifo" || got.Action != "inspect and remove it" {
			t.Fatalf("admin refusal = %v", err)
		}
	case <-time.After(bounds.TestDeadline(0)):
		t.Fatal("malformed worktree admin entry reached porcelain")
	}
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "worktree") {
		t.Fatalf("worktree invoked: %s", data)
	}
}

func TestWorktreesRejectsBadCommonDirBeforePorcelain(t *testing.T) {
	for _, tc := range []struct {
		name, mode string
		want       []string
	}{
		{"missing", "bad-rev-parse", []string{"missing-common", "missing path"}},
		{"empty", "empty-rev-parse", []string{"empty path"}},
		{"symlink to directory", "symlink-rev-parse", []string{"symlink-common", "symlink"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assertBadCommonDirRefusal(t, tc.mode, tc.want)
		})
	}
}

func assertBadCommonDirRefusal(t *testing.T, mode string, want []string) {
	t.Helper()
	root := newRepo(t)
	logPath := filepath.Join(t.TempDir(), "argv")
	commonDir := gittest.StubGit(t, root, mode, logPath)
	if mode == "symlink-rev-parse" {
		if err := os.Symlink(filepath.Join(root, ".git"), commonDir); err != nil {
			capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable: %v", err))
		}
	}
	_, err := Worktrees(root)
	var resolution *ResolutionError
	if !errors.As(err, &resolution) || resolution.Action != "investigate the git failure" {
		t.Fatalf("resolution refusal = %v", err)
	}
	for _, fragment := range want {
		if !strings.Contains(err.Error(), fragment) {
			t.Fatalf("resolution refusal = %v, want %q", err, fragment)
		}
	}
	data, _ := os.ReadFile(logPath)
	if strings.Contains(string(data), "worktree") {
		t.Fatalf("worktree invoked: %s", data)
	}
}

func TestWorktreesPropagatesRevParseFailureBeforePorcelain(t *testing.T) {
	root := newRepo(t)
	logPath := filepath.Join(t.TempDir(), "argv")
	gittest.StubGit(t, root, "fail-rev-parse-noisy", logPath)
	var err error
	_, err = Worktrees(root)
	var resolution *ResolutionError
	if !errors.As(err, &resolution) || resolution.Err == nil || !strings.Contains(resolution.Err.Error(), "rev-parse") || !strings.Contains(resolution.Err.Error(), "fatal: common directory unavailable") || resolution.Action != "investigate the git failure" {
		t.Fatalf("err=%v", err)
	}
	data, _ := os.ReadFile(logPath)
	if strings.Contains(string(data), "worktree") {
		t.Fatalf("worktree invoked: %s", data)
	}
}

func TestWorktreeListTimeoutDefaultUsesPolicy(t *testing.T) {
	if worktreeListTimeout != bounds.WorktreeListTimeout {
		t.Fatalf("worktreeListTimeout=%s, want %s", worktreeListTimeout, bounds.WorktreeListTimeout)
	}
}

func TestWorktreesBoundsEachChildAndPreservesStdout(t *testing.T) {
	for _, tc := range []struct{ mode, invocation string }{
		{"block-worktree", "worktree list"},
		{"block-rev-parse", "rev-parse"},
	} {
		t.Run(tc.mode, func(t *testing.T) {
			restore := SetWorktreeListTimeoutForTest(100 * time.Millisecond)
			t.Cleanup(restore)
			root := newRepo(t)
			gittest.StubGit(t, root, tc.mode, filepath.Join(t.TempDir(), "argv"))
			done := make(chan error, 1)
			go func() { _, err := Worktrees(root); done <- err }()
			select {
			case err := <-done:
				var typed *ResolutionError
				if !errors.As(err, &typed) || !strings.Contains(err.Error(), tc.invocation) || !strings.Contains(err.Error(), "100ms") || !strings.Contains(err.Error(), "investigate the git failure") {
					t.Fatalf("timeout refusal = %v", err)
				}
			case <-time.After(time.Second):
				t.Fatalf("%s did not return within the overridden bound", tc.mode)
			}
		})
	}

	t.Run("noisy list", func(t *testing.T) {
		root := newRepo(t)
		gittest.StubGit(t, root, "noisy-list", filepath.Join(t.TempDir(), "argv"))
		worktrees, err := Worktrees(root)
		if err != nil || len(worktrees) != 1 || worktrees[0].Path != root {
			t.Fatalf("noisy list = %#v, %v", worktrees, err)
		}
	})
}

func TestWorktreesTypesStartFailures(t *testing.T) {
	for _, tc := range []struct {
		name, mode, invocation, failure string
		setup                           func(*testing.T)
	}{
		{"porcelain", "vanish-after-rev-parse", "worktree list", "executable file not found", nil},
		{"rev parse", "", "rev-parse", "executable file not found", func(t *testing.T) { t.Setenv("PATH", t.TempDir()) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newRepo(t)
			if tc.mode != "" {
				gittest.StubGit(t, root, tc.mode, filepath.Join(t.TempDir(), "argv"))
			}
			if tc.setup != nil {
				tc.setup(t)
			}
			_, err := Worktrees(root)
			var typed *ResolutionError
			if !errors.As(err, &typed) || !strings.Contains(err.Error(), tc.invocation) || !strings.Contains(err.Error(), tc.failure) || typed.Action != "investigate the git failure" {
				t.Fatalf("start failure = %v", err)
			}
		})
	}
}

// newRepo initialises a repo with one commit and returns its root.
func newRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	runGit(t, root, "init")
	if err := os.WriteFile(filepath.Join(root, "a.txt"), []byte("hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "add", "-A")
	runGit(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-m", "init")
	return root
}

// TestTreeHashCleanTreeMatchesHEAD is the same-tree property the gate relies on. A
// clean working tree hashes to HEAD's tree object.
func TestTreeHashCleanTreeMatchesHEAD(t *testing.T) {
	root := newRepo(t)
	want := runGit(t, root, "rev-parse", "HEAD^{tree}")
	want = want[:len(want)-1] // strip trailing newline
	got := TreeHash(root)
	if got != want {
		t.Fatalf("clean tree: got %q, want %q", got, want)
	}
	if got == "none" {
		t.Fatal("clean repo yielded none")
	}
}

// TestTreeHashDirtyTreeDiffers checks that changing on-disk content changes the hash.
func TestTreeHashDirtyTreeDiffers(t *testing.T) {
	root := newRepo(t)
	clean := TreeHash(root)

	// Modify an existing tracked file.
	if err := os.WriteFile(filepath.Join(root, "a.txt"), []byte("changed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if modified := TreeHash(root); modified == clean {
		t.Fatalf("modifying a file did not change the hash (%q)", modified)
	}

	// Add a new untracked file to a fresh clean repo.
	root2 := newRepo(t)
	base := TreeHash(root2)
	if err := os.WriteFile(filepath.Join(root2, "b.txt"), []byte("new\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if added := TreeHash(root2); added == base {
		t.Fatalf("adding a file did not change the hash (%q)", added)
	}
}

func TestAllFilesFactsExpandsWithoutChangingFacts(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init", "-q", "-b", "main")
	runGit(t, root, "config", "user.email", "t@example.com")
	runGit(t, root, "config", "user.name", "t")
	if err := os.MkdirAll(filepath.Join(root, "nested", "new"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "nested", "new", "file.txt"), []byte("new\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	legacy, err := Facts(root)
	if err != nil {
		t.Fatal(err)
	}
	all, err := AllFilesFacts(root)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := legacy.Changes[0].Path, "nested/"; got != want {
		t.Fatalf("Facts changed its legacy collapsed-directory result: got %q, want %q", got, want)
	}
	if got, want := all.Changes[0].Path, "nested/new/file.txt"; got != want {
		t.Fatalf("AllFilesFacts path = %q, want %q", got, want)
	}
}

func TestTreeHashExcludesIgnoredContent(t *testing.T) {
	root := newRepo(t)
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("ignored/\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "add", ".gitignore")
	runGit(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-qm", "ignore fixture")
	clean := TreeHash(root)
	if err := os.MkdirAll(filepath.Join(root, "ignored"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "ignored", "smuggled.txt"), []byte("ignored\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := TreeHash(root); got != clean {
		t.Fatalf("ignored content changed tree from %q to %q", clean, got)
	}
}

// TestTreeHashNonRepoReturnsNone checks the failure posture for a non-repo dir.
func TestTreeHashNonRepoReturnsNone(t *testing.T) {
	dir := t.TempDir()
	if got := TreeHash(dir); got != "none" {
		t.Fatalf("non-repo dir: got %q, want \"none\"", got)
	}
}

// TestTreeHashLeavesNoStrayIndex confirms the throwaway index is external. The call
// must not drop an index file into the working tree.
func TestTreeHashLeavesNoStrayIndex(t *testing.T) {
	root := newRepo(t)
	before := listDir(t, root)
	index := runGit(t, root, "write-tree")
	_ = TreeHash(root)
	after := listDir(t, root)

	if _, err := os.Stat(filepath.Join(root, "index")); err == nil {
		t.Fatal("TreeHash created root/index")
	}
	if len(after) != len(before) {
		t.Fatalf("TreeHash changed the working tree: before %v, after %v", before, after)
	}
	if got := runGit(t, root, "write-tree"); got != index {
		t.Fatalf("TreeHash changed the real index from %q to %q", index, got)
	}
}

func TestLandedStateCountsDirtyPathsOnlyInNamedCheckout(t *testing.T) {
	root := newRepo(t)
	runGit(t, root, "branch", "-M", "main")
	runGit(t, root, "remote", "add", "origin", root)
	runGit(t, root, "update-ref", "refs/remotes/origin/main", "HEAD")
	linked := filepath.Join(t.TempDir(), "feature worktree")
	runGit(t, root, "worktree", "add", "-q", "-b", "feature", linked, "HEAD")
	if err := os.WriteFile(filepath.Join(linked, "a.txt"), []byte("dirty\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "checkout", "-q", "-b", "ahead")
	if err := os.WriteFile(filepath.Join(root, "b.txt"), []byte("ahead\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "add", "b.txt")
	runGit(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-qm", "ahead")
	runGit(t, root, "branch", "ahead-copy")
	for _, branch := range []string{"ahead", "ahead-copy"} {
		runGit(t, root, "config", "branch."+branch+".remote", "origin")
		runGit(t, root, "config", "branch."+branch+".merge", "refs/heads/main")
	}
	if err := os.WriteFile(filepath.Join(root, "a.txt"), []byte("root dirty\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(linked, "linked.txt"), []byte("linked dirty\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name      string
		root      string
		wantDirty int
	}{
		{name: "primary", root: root, wantDirty: 1},
		{name: "linked", root: linked, wantDirty: 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			state, err := LandedState(tc.root)
			if err != nil {
				t.Fatal(err)
			}
			if state.DirtyPaths != tc.wantDirty || state.UnpushedCommits != 1 || state.UniqueBranches != 2 {
				t.Fatalf("LandedState = %#v, want dirty=%d ahead=1 unique=2", state, tc.wantDirty)
			}
		})
	}
}

func TestLandedStateCountsDefaultBranchUpstreamAhead(t *testing.T) {
	root := newRepo(t)
	runGit(t, root, "branch", "-M", "main")
	runGit(t, root, "remote", "add", "origin", root)
	runGit(t, root, "update-ref", "refs/remotes/origin/main", "HEAD")
	runGit(t, root, "config", "branch.main.remote", "origin")
	runGit(t, root, "config", "branch.main.merge", "refs/heads/main")
	if err := os.WriteFile(filepath.Join(root, "ahead.txt"), []byte("ahead\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "add", "ahead.txt")
	runGit(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-qm", "ahead")

	state, err := LandedState(root)
	if err != nil {
		t.Fatal(err)
	}
	if state.UnpushedCommits != 1 {
		t.Fatalf("default branch ahead count = %d, want 1", state.UnpushedCommits)
	}
}

func TestLandedStateIgnoresDirtySiblingWithNewlinePath(t *testing.T) {
	root := newRepo(t)
	runGit(t, root, "branch", "-M", "main")
	linked := filepath.Join(t.TempDir(), "linked\nworktree")
	runGit(t, root, "worktree", "add", "-q", "--detach", linked, "HEAD")
	if err := os.WriteFile(filepath.Join(linked, "a.txt"), []byte("dirty\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	state, err := LandedState(root)
	if err != nil {
		t.Fatalf("LandedState with newline path: %v", err)
	}
	if state.DirtyPaths != 0 {
		t.Fatalf("newline sibling worktree dirty count = %d, want 0", state.DirtyPaths)
	}
}

func TestLandedStateIgnoresUnreadableSibling(t *testing.T) {
	root := newRepo(t)
	runGit(t, root, "branch", "-M", "main")
	linked := filepath.Join(t.TempDir(), "missing-worktree")
	runGit(t, root, "worktree", "add", "-q", "--detach", linked, "HEAD")
	if err := os.RemoveAll(linked); err != nil {
		t.Fatal(err)
	}

	state, err := LandedState(root)
	if err != nil {
		t.Fatalf("LandedState with unreadable sibling: %v", err)
	}
	if state.DirtyPaths != 0 {
		t.Fatalf("unreadable sibling dirty count = %d, want 0", state.DirtyPaths)
	}
}

func TestLandedStateGitFailureIsUnknown(t *testing.T) {
	root := newRepo(t)
	t.Setenv("PATH", "")
	if _, err := LandedState(root); err == nil {
		t.Fatal("LandedState swallowed missing git")
	}
}

// TestRefResolvesAndBranchExists exercises the two guard probes and their fail-safe
// posture. They run in the process cwd (the agent's working dir), so the test chdirs
// into a fixture repo and restores.
func TestRefResolvesAndBranchExists(t *testing.T) {
	root := newRepo(t)
	runGit(t, root, "branch", "known-branch")
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(old)

	if !RefResolves("HEAD") {
		t.Error("RefResolves(HEAD) = false, want true")
	}
	if RefResolves("definitely-not-a-ref-xyz") {
		t.Error("RefResolves(bogus) = true, want false")
	}
	if !BranchExists("known-branch") {
		t.Error("BranchExists(known-branch) = false, want true")
	}
	if BranchExists("no-such-branch-xyz") {
		t.Error("BranchExists(absent) = true, want false")
	}
}

func listDir(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return names
}

// newTwoBranchRepo initialises a repo with one commit and exactly two local branches.
// Neither branch is a resolvable default candidate: there is no origin/HEAD, and no
// branch named "main" for the candidate probe to verify.
func newTwoBranchRepo(t *testing.T) string {
	t.Helper()
	root := newRepo(t)
	runGit(t, root, "branch", "-M", "master")
	runGit(t, root, "branch", "feature")
	return root
}

// TestFactsUnresolvableDefault pins what Facts reports when no default branch resolves:
// the state, not a guess. Ahead/Behind stay zero because there is nothing to measure
// against, and the rest of the snapshot still has to arrive.
func TestFactsUnresolvableDefault(t *testing.T) {
	root := newTwoBranchRepo(t)
	if err := os.WriteFile(filepath.Join(root, "b.txt"), []byte("dirty\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	facts, err := Facts(root)

	if err != nil {
		t.Fatalf("Facts on an unresolvable default: %v", err)
	}
	if facts.DefaultResolved {
		t.Error("DefaultResolved = true, want false")
	}
	if facts.DefaultBranch != "" {
		t.Errorf("DefaultBranch = %q, want empty", facts.DefaultBranch)
	}
	if facts.Ahead != 0 || facts.Behind != 0 {
		t.Errorf("Ahead/Behind = %d/%d, want 0/0", facts.Ahead, facts.Behind)
	}
	if facts.Branch != "master" || !facts.Dirty {
		t.Errorf("rest of the snapshot lost: branch %q, dirty %v", facts.Branch, facts.Dirty)
	}
}

// TestResolvedDefaultSoleMaster is the sole-local-branch fallback. A master-only
// repository has no origin/HEAD and no "main" to verify. The lone local branch is
// the only evidence of its default.
func TestResolvedDefaultSoleMaster(t *testing.T) {
	root := newRepo(t)
	runGit(t, root, "branch", "-M", "master")

	def, ok := ResolvedDefault(root)

	if !ok || def != "master" {
		t.Fatalf("ResolvedDefault = (%q, %v), want (\"master\", true)", def, ok)
	}
}

// TestResolvedDefaultNoLocalBranches covers the empty end of the sole-local-branch
// fallback. A repository with no commits has no branch to fall back to. If the code
// indexes the list before counting it, it panics instead of reporting the unresolved
// state.
func TestResolvedDefaultNoLocalBranches(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")

	def, ok := ResolvedDefault(root)

	if ok || def != "" {
		t.Fatalf("ResolvedDefault = (%q, %v), want (\"\", false)", def, ok)
	}
}

// TestResolvedDefaultUnresolvableNamesNothing pins the ok=false return as an empty name.
// No caller can put a branch this repository does not have into a message or a ref.
func TestResolvedDefaultUnresolvableNamesNothing(t *testing.T) {
	def, ok := ResolvedDefault(newTwoBranchRepo(t))

	if ok || def != "" {
		t.Fatalf("ResolvedDefault = (%q, %v), want (\"\", false)", def, ok)
	}
}
