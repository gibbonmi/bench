package transaction

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
)

func TestReleasedLeaseLeavesNoLockFile(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "locks")
	root := t.TempDir()
	destinations := make([]string, 8)
	last, lastKey := 0, ""
	for i := range destinations {
		destinations[i] = filepath.Join(root, "managed", fmt.Sprint(i))
		key, err := destinationKey(destinations[i])
		if err != nil {
			t.Fatal(err)
		}
		if key > lastKey {
			last, lastKey = i, key
		}
	}
	// The held key sorts last, so the refused writer acquires every other key before it fails.
	holder, err := lockIn(dir, destinations[last:last+1])
	if err != nil {
		t.Fatal(err)
	}
	if _, err := lockIn(dir, destinations); !errors.Is(err, syscall.EWOULDBLOCK) {
		t.Fatalf("contended lease = %v, want the busy refusal", err)
	}
	assertLockFiles(t, dir, 1)
	if err := holder.Close(); err != nil {
		t.Fatal(err)
	}
	lease, err := lockIn(dir, destinations)
	if err != nil {
		t.Fatal(err)
	}
	if err := lease.Close(); err != nil {
		t.Fatal(err)
	}
	assertLockFiles(t, dir, 0)
}

func TestReleasedLockFileKeepsExclusion(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "locks")
	destination := []string{filepath.Join(t.TempDir(), "managed")}
	var holders, overlaps, acquired atomic.Int64
	failures := make(chan error, 8)
	var writers sync.WaitGroup
	for range 8 {
		writers.Add(1)
		go func() {
			defer writers.Done()
			for range 2000 {
				lease, err := lockIn(dir, destination)
				if errors.Is(err, syscall.EWOULDBLOCK) {
					continue
				}
				if err != nil {
					failures <- err
					return
				}
				acquired.Add(1)
				if holders.Add(1) > 1 {
					overlaps.Add(1)
				}
				runtime.Gosched()
				holders.Add(-1)
				if err := lease.Close(); err != nil {
					failures <- err
					return
				}
			}
		}()
	}
	writers.Wait()
	close(failures)
	var failed []error
	for err := range failures {
		failed = append(failed, err)
	}
	if overlaps.Load() != 0 || len(failed) != 0 || acquired.Load() == 0 {
		t.Fatalf("concurrent writers: %d overlapping holds, %d acquisitions, failures: %v", overlaps.Load(), acquired.Load(), errors.Join(failed...))
	}
}

func TestCloseRemovesOnlyTheHeldLockFile(t *testing.T) {
	for _, tc := range []struct {
		name    string
		replace bool
	}{
		{name: "removed"},
		{name: "replaced", replace: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			lease, err := lockIn(filepath.Join(t.TempDir(), "locks"), []string{filepath.Join(t.TempDir(), "managed")})
			if err != nil {
				t.Fatal(err)
			}
			path := lease.files[0].Name()
			// An outside actor removes the held lock file, and another writer can then create a new one.
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			if tc.replace {
				if err := os.WriteFile(path, nil, 0o600); err != nil {
					t.Fatal(err)
				}
			}
			if err := lease.Close(); err != nil {
				t.Fatalf("close after the lock file was %s = %v, want nil", tc.name, err)
			}
			if _, err := os.Lstat(path); tc.replace && err != nil {
				t.Fatalf("close removed the replacement lock file: %v", err)
			}
		})
	}
}

func TestCloseReportsAnUnreadableLockPath(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "locks")
	lease, err := lockIn(dir, []string{filepath.Join(t.TempDir(), "managed")})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })
	if err := os.Chmod(dir, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(lease.files[0].Name()); err == nil {
		capability.Capability(t, capability.Privilege, "lock directory without search permission remains searchable by this user")
	}
	if err := lease.Close(); !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("close on an unsearchable lock directory = %v, want the permission error", err)
	}
}

func assertLockFiles(t *testing.T, dir string, want int) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != want {
		t.Fatalf("lock directory holds %d files, want %d", len(entries), want)
	}
}
