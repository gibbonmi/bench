package worktree

import (
	"os"
	"syscall"
	"testing"
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
