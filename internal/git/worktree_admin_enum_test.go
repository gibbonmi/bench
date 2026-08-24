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
