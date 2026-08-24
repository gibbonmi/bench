package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/gittest"
)

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
