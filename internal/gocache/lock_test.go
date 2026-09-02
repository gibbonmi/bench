package gocache

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
)

// The helper-process protocol. A lock row needs two processes, because a POSIX record lock
// is owned per process and one process can never contend with itself. The second process
// is this same test binary re-executed with one role in its environment, so no row depends
// on a compiled fixture or on `go run`.
const (
	helperRoleEnv  = "BENCH_GOCACHE_HELPER_ROLE"
	helperHomeEnv  = "BENCH_GOCACHE_HELPER_HOME"
	helperReadyEnv = "BENCH_GOCACHE_HELPER_READY"
	helperStopEnv  = "BENCH_GOCACHE_HELPER_STOP"
	helperClean    = "clean"
	helperHold     = "hold"
)

// helperWait bounds every wait a lock row makes on the second process. A row that hangs
// reports nothing, so the deadline turns a lost signal into a named failure. A sentinel
// handshake between two processes holds no window of its own, so the wait takes the floor
// the derivation gives a zero bound.
var helperWait = bounds.TestDeadline(0)

// TestGocacheHelperProcess is the second process every lock row drives, inert unless the
// parent selects it through the role environment variable. It is a test function so the
// row can re-execute this binary.
func TestGocacheHelperProcess(t *testing.T) {
	home := os.Getenv(helperHomeEnv)
	switch os.Getenv(helperRoleEnv) {
	case helperClean:
		out, code := clean([]string{"HOME=" + home, "PATH=" + os.Getenv("PATH")})
		fmt.Print(out)
		os.Exit(code)
	case helperHold:
		holder, err := HoldDir(cacheDir(home))
		if err != nil {
			fmt.Print("hold failed: " + err.Error())
			os.Exit(1)
		}
		if err := os.WriteFile(os.Getenv(helperReadyEnv), []byte("held\n"), 0o600); err != nil {
			os.Exit(1)
		}
		if verdict := waitForFile(os.Getenv(helperStopEnv), helperWait); verdict != "" {
			fmt.Fprintln(os.Stderr, verdict)
		}
		_ = holder.Release()
		os.Exit(0)
	default:
		return
	}
}

// cacheDir is the directory the fixture home derives, spelled through the production
// derivation so a test never states the path a second time.
func cacheDir(home string) string {
	dir, err := Dir([]string{"HOME=" + home})
	if err != nil {
		panic(err)
	}
	return dir
}

// waitForFile polls until path exists or the window runs out. It is the helper's half of
// the sentinel handshake; a lost signal ends the helper instead of leaking it. The helper
// process holds no *testing.T, so the wait answers the timeout verdict on expiry and the
// empty string on arrival, and its caller reports it. The poll interval stays a literal:
// it paces the loop and bounds nothing.
func waitForFile(path string, window time.Duration) string {
	deadline := time.Now().Add(window)
	for {
		if _, err := os.Stat(path); err == nil {
			return ""
		}
		if !time.Now().Before(deadline) {
			return bounds.TestTimeoutVerdict("the sentinel file "+path, window)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// runCleanProcess runs `bench cache clean` in a second process against the fixture home and
// returns its stdout and exit code.
func runCleanProcess(t *testing.T, home string) (string, int) {
	t.Helper()
	command := exec.Command(os.Args[0], "-test.run=^TestGocacheHelperProcess$")
	command.Env = []string{
		helperRoleEnv + "=" + helperClean,
		helperHomeEnv + "=" + home,
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + home,
	}
	output, err := command.CombinedOutput()
	if err == nil {
		return string(output), 0
	}
	exit, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("clean process failed to run: %v; output=%q", err, output)
	}
	return string(output), exit.ExitCode()
}

// startHolderProcess starts a second process that holds the shared lock and returns a stop
// function. The process signals through a ready file, so the row never races the hold.
func startHolderProcess(t *testing.T, home string) func() {
	t.Helper()
	ready := filepath.Join(t.TempDir(), "ready")
	stop := filepath.Join(t.TempDir(), "stop")
	command := exec.Command(os.Args[0], "-test.run=^TestGocacheHelperProcess$")
	command.Env = []string{
		helperRoleEnv + "=" + helperHold,
		helperHomeEnv + "=" + home,
		helperReadyEnv + "=" + ready,
		helperStopEnv + "=" + stop,
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + home,
	}
	if err := command.Start(); err != nil {
		t.Fatalf("holder process did not start: %v", err)
	}
	released := false
	release := func() {
		if released {
			return
		}
		released = true
		_ = os.WriteFile(stop, []byte("stop\n"), 0o600)
		_ = command.Wait()
	}
	t.Cleanup(release)
	if verdict := waitForFile(ready, helperWait); verdict != "" {
		t.Fatalf("holder process did not signal a hold: %s", verdict)
	}
	return release
}

// L05: two holders take the shared lock at the same time, so the bound never serializes two
// concurrent gates. The rows are two processes, because one process's second request
// replaces its own lock rather than contending with it. An exclusive-mode holder makes the
// second acquisition fail here.
func TestTwoHoldersAcquireTheSharedLockAtTheSameTime(t *testing.T) {
	home := t.TempDir()
	startHolderProcess(t, home)
	acquired := make(chan error, 1)
	go func() {
		holder, err := HoldDir(cacheDir(home))
		if err == nil {
			_ = holder.Release()
		}
		acquired <- err
	}()
	select {
	case err := <-acquired:
		if err != nil {
			t.Fatalf("second holder = %v, want a shared lock beside the first", err)
		}
	case <-time.After(helperWait):
		t.Fatal(bounds.TestTimeoutVerdict("the second holder to acquire the shared lock", helperWait))
	}
}

// A migrated wait reports a timeout instead of returning silently. A zero window expires
// on the first pass, so the row observes the verdict with no wall-clock wait, and it
// grades the two halves the reader of a red run needs: the wait's name and its window.
func TestWaitForFileAnswersTheTimeoutVerdictOnAnExpiredWindow(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "never-written")
	verdict := waitForFile(missing, 0)
	if !strings.Contains(verdict, missing) || !strings.Contains(verdict, "0s") {
		t.Fatalf("expired waitForFile = %q, want the wait name %q and the window", verdict, missing)
	}
	present := filepath.Join(t.TempDir(), "written")
	if err := os.WriteFile(present, []byte("ready\n"), 0o600); err != nil {
		t.Fatalf("write sentinel: %v", err)
	}
	if arrived := waitForFile(present, helperWait); arrived != "" {
		t.Fatalf("waitForFile over an existing sentinel = %q, want no verdict", arrived)
	}
}

// A holder creates the directory and the lock file it needs, so a first run on a fresh
// machine locks rather than refusing.
func TestHolderCreatesTheDirectoryAndTheLockFile(t *testing.T) {
	home := t.TempDir()
	holder, err := Hold([]string{"HOME=" + home})
	if err != nil {
		t.Fatalf("Hold = %v, want a holder", err)
	}
	defer holder.Release()
	if _, err := os.Stat(filepath.Join(cacheDir(home), LockFile)); err != nil {
		t.Fatalf("lock file after Hold = %v, want it created", err)
	}
}

// Two holds inside one process share one descriptor, so releasing the first leaves the
// second's lock in place. Closing any descriptor for a file drops every record lock the
// process holds on it, which is the failure this asserts against.
func TestReleasingOneInProcessHoldKeepsTheOther(t *testing.T) {
	home := t.TempDir()
	dir := cacheDir(home)
	first, err := HoldDir(dir)
	if err != nil {
		t.Fatalf("first hold = %v", err)
	}
	second, err := HoldDir(dir)
	if err != nil {
		t.Fatalf("second hold = %v", err)
	}
	if err := first.Release(); err != nil {
		t.Fatalf("first release = %v", err)
	}
	out, code := runCleanProcess(t, home)
	if code != 1 || !strings.Contains(out, "cache in use") {
		t.Fatalf("clean beside the surviving hold = %q exit %d, want a refusal", out, code)
	}
	if err := second.Release(); err != nil {
		t.Fatalf("second release = %v", err)
	}
}
