// Package processtest supplies the real-process fixtures that lifetime tests share. A test
// binary that imports it becomes an escaped holder when holdEnv names a PID file.
package processtest

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
)

const holdEnv = "BENCH_PROCESSTEST_HOLD"

// init runs before the test main, so a holder never runs a test. It leaves its parent's
// group, publishes its PID, and becomes a sleep that keeps every inherited descriptor.
func init() {
	path := os.Getenv(holdEnv)
	if path == "" {
		return
	}
	if syscall.Setpgid(0, 0) != nil || os.WriteFile(path+".partial", []byte(strconv.Itoa(os.Getpid())), 0o600) != nil || os.Rename(path+".partial", path) != nil {
		os.Exit(3)
	}
	if sleep, err := exec.LookPath("sleep"); err == nil {
		_ = syscall.Exec(sleep, []string{"sleep", "300"}, os.Environ())
	}
	os.Exit(3)
}

// Holder is one escaped descendant. A child passes Env to it and starts Command in the
// background; the test kills the holder at cleanup.
type Holder struct {
	Env, Command, path string
}

// NewHolder prepares a holder for one test.
func NewHolder(t *testing.T) *Holder {
	t.Helper()
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "holder.pid")
	h := &Holder{Env: holdEnv + "=" + path, Command: executable, path: path}
	t.Cleanup(func() {
		if pid := readPID(path); pid > 0 {
			_ = syscall.Kill(pid, syscall.SIGKILL)
		}
	})
	return h
}

// PID waits until the holder has left its parent's group and returns its PID.
func (h *Holder) PID(t *testing.T) int {
	t.Helper()
	return WaitPID(t, h.path)
}

// WaitPID waits for a child to publish its PID in path, and kills that PID at cleanup.
func WaitPID(t *testing.T, path string) int {
	t.Helper()
	window := bounds.TestDeadline(0)
	for deadline := time.Now().Add(bounds.FixedWindow(window)); time.Now().Before(deadline); {
		if pid := readPID(path); pid > 0 {
			t.Cleanup(func() { _ = syscall.Kill(pid, syscall.SIGKILL) })
			return pid
		}
		time.Sleep(bounds.FixedWindow(10 * time.Millisecond))
	}
	t.Fatal(bounds.TestTimeoutVerdict("a child to publish "+filepath.Base(path), window))
	return 0
}

func readPID(path string) int {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	pid, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	return pid
}
