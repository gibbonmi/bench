package probe

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// lockHolderPathEnv, lockHolderReleaseEnv, and lockHolderReadyEnv carry the lock path, the
// release marker, and the ready marker to the lock-holding process. lockTrySubjectEnv carries
// a subject path to a process that tries that subject's lock once. A record lock belongs to its
// process, so a refusal under a held lock is only proven when another process holds that lock.
const (
	lockHolderPathEnv    = "BENCH_PROBE_LOCK_PATH"
	lockHolderReleaseEnv = "BENCH_PROBE_LOCK_RELEASE"
	lockHolderReadyEnv   = "BENCH_PROBE_LOCK_READY"
	lockTrySubjectEnv    = "BENCH_PROBE_LOCK_TRY"
)

// lockProcessArgv re-enters this test binary through TestProbeRefusesUnderALiveGateRun, which
// holds or tries a lock as its environment names.
func lockProcessArgv() []string {
	return []string{os.Args[0], "-test.run=^TestProbeRefusesUnderALiveGateRun$", "-test.timeout=" + interruptDeadline().String()}
}

// lockHolderEnv answers the environment of a holder of path, its release marker, and its ready
// marker. The holder writes the ready marker once it holds the lock, and removes it as it
// exits. The test's end writes the release and waits for that removal, so no holder outlives
// its directories.
func lockHolderEnv(t *testing.T, path string) (env []string, release, ready string) {
	t.Helper()
	dir := scratchDir(t)
	release, ready = filepath.Join(dir, "release"), filepath.Join(dir, "ready")
	t.Cleanup(func() {
		_ = os.WriteFile(release, nil, 0o644)
		for expiry := time.Now().Add(interruptDeadline()); time.Now().Before(expiry); time.Sleep(10 * time.Millisecond) {
			if _, err := os.Stat(ready); os.IsNotExist(err) {
				return
			}
		}
	})
	return []string{lockHolderPathEnv + "=" + path, lockHolderReleaseEnv + "=" + release, lockHolderReadyEnv + "=" + ready}, release, ready
}

// startLockHolder holds the record lock on path in its own process and answers the release.
func startLockHolder(t *testing.T, path string) func() {
	t.Helper()
	env, release, ready := lockHolderEnv(t, path)
	argv := lockProcessArgv()
	holder := exec.Command(argv[0], argv[1:]...)
	holder.Env = append(os.Environ(), env...)
	if err := holder.Start(); err != nil {
		t.Fatal(err)
	}
	awaitFile(t, ready)
	return func() {
		_ = os.WriteFile(release, nil, 0o644)
		_ = holder.Wait()
	}
}

// lockHolderHook answers the stub shell that starts a holder of path in the background and
// waits until it holds the lock. The holder's output goes nowhere, so the Go child's stream
// closes when the stub exits, and the test's end releases the holder.
func lockHolderHook(t *testing.T, path string) string {
	t.Helper()
	env, _, ready := lockHolderEnv(t, path)
	return shellCommand(slices.Concat([]string{"env"}, env, lockProcessArgv())) + " >/dev/null 2>&1 &\n" +
		"while [ ! -e " + sanitize.ShellQuote(ready) + " ]; do sleep 0.01; done\n"
}

// lockTryHook answers the stub shell that tries the lock of subject from another process and
// writes that process's answer to out.
func lockTryHook(subject, out string) string {
	return shellCommand(slices.Concat([]string{"env", lockTrySubjectEnv + "=" + subject}, lockProcessArgv())) + " > " + sanitize.ShellQuote(out) + " 2>&1\n"
}

// tryRecordLock takes and releases the subject lock of the file at subject once, and prints
// the refusal line that lockSubject answers, or "granted".
func tryRecordLock(subject string) {
	held, line := lockSubject(subjectFromPath(subject))
	if line != "" {
		fmt.Print(line)
		return
	}
	held.release()
	fmt.Println("granted")
}

// subjectFromPath names the fixture subject at path the way the probe's own rows name it.
func subjectFromPath(path string) subject {
	return subject{path: path, display: filepath.Base(path)}
}

// shellCommand joins words into one shell command, each word quoted.
func shellCommand(words []string) string {
	quoted := make([]string, len(words))
	for i, word := range words {
		quoted[i] = sanitize.ShellQuote(word)
	}
	return strings.Join(quoted, " ")
}

// startSubjectHolder holds the subject lock of the file at subject in another process.
func startSubjectHolder(t *testing.T, subject string) func() {
	t.Helper()
	return startLockHolder(t, subjectLockPath(subject))
}

// holdRecordLock takes the exclusive record lock on path, writes the ready marker, and parks
// until the release marker appears. So the refusal under test observes a genuinely held lock.
func holdRecordLock(t *testing.T, path, releaseMarker, readyMarker string) {
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
	if err := os.WriteFile(readyMarker, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(releaseMarker); err == nil {
			_ = os.Remove(readyMarker)
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// awaitFile waits for the file at path to appear and then settles. A process writes its marker
// just before the step that the waiter observes.
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
	t.Fatalf("no process created %s within %s", path, interruptDeadline())
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

// removeLockDir removes the empty lock directory, so the open that follows answers ENOENT.
func removeLockDir(t *testing.T, path string) {
	t.Helper()
	if err := os.Remove(filepath.Dir(path)); err != nil {
		t.Fatal(err)
	}
}

// permissions answers the permission bits of the file at path.
func permissions(t *testing.T, path string) os.FileMode {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	return info.Mode().Perm()
}
