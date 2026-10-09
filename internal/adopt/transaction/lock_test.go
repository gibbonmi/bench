package transaction

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
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
