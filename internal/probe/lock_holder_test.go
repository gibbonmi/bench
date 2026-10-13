package probe

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/gocache"
)

// lockHolderPathEnv and lockHolderReleaseEnv carry the lock path and the release marker to
// the lock-holding process. A record lock belongs to its process, so a refusal under a held
// lock is only proven when another process holds that lock.
const (
	lockHolderPathEnv    = "BENCH_PROBE_LOCK_PATH"
	lockHolderReleaseEnv = "BENCH_PROBE_LOCK_RELEASE"
)

// startLockHolder holds the record lock on path in its own process and answers the release.
// The process re-enters this test binary through TestProbeRefusesUnderALiveGateRun.
func startLockHolder(t *testing.T, path string) func() {
	t.Helper()
	releaseMarker := filepath.Join(scratchDir(t), "release")
	holder := exec.Command(os.Args[0], "-test.run=^TestProbeRefusesUnderALiveGateRun$", "-test.timeout="+interruptDeadline().String())
	holder.Env = append(os.Environ(), lockHolderPathEnv+"="+path, lockHolderReleaseEnv+"="+releaseMarker)
	if err := holder.Start(); err != nil {
		t.Fatal(err)
	}
	awaitFile(t, path)
	return func() {
		_ = os.WriteFile(releaseMarker, nil, 0o644)
		_ = holder.Wait()
	}
}

// startSubjectHolder holds the subject lock of the file at subject in another process.
func startSubjectHolder(t *testing.T, subject string) func() {
	t.Helper()
	return startLockHolder(t, subjectLockPath(subject))
}

// holdRecordLock takes the exclusive record lock on path and parks until the release marker
// appears, so the refusal under test observes a genuinely held lock.
func holdRecordLock(t *testing.T, path, releaseMarker string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	lock, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Close()
	held := gocache.RecordLock(syscall.F_WRLCK)
	if err := syscall.FcntlFlock(lock.Fd(), syscall.F_SETLK, &held); err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(releaseMarker); err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// awaitFile waits for the holder's lock file to appear and then settles, because the
// create and the lock are two steps in the holder's own process.
func awaitFile(t *testing.T, path string) {
	t.Helper()
	expiry := time.Now().Add(interruptDeadline())
	for time.Now().Before(expiry) {
		if info, err := os.Stat(path); err == nil && info.Mode().IsRegular() {
			time.Sleep(50 * time.Millisecond)
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("the lock holder did not create %s within %s", path, interruptDeadline())
}

// setHook installs hook in slot for the test's span and then puts the earlier value back.
func setHook[T any](t *testing.T, slot *T, hook T) {
	t.Helper()
	earlier := *slot
	*slot = hook
	t.Cleanup(func() { *slot = earlier })
}

// replaceLockFile puts a new file at the lock path, so a descriptor opened before the call
// holds a file that the path no longer names.
func replaceLockFile(t *testing.T, path string) {
	t.Helper()
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	writeFixtureFile(t, path, "", 0o600)
}
