package preprelease

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/contract"
)

// interruptDuringPreflight runs prep-release with the preflight script held open between
// staging and promotion, signals the run once it is demonstrably in that window, and
// returns the exit code. Signalling on a timer instead would race the toolchain: the
// steps before the preflight take as long as the host's Go cache says they do.
func interruptDuringPreflight(t *testing.T, r shipRepo) int {
	t.Helper()
	cmd := exec.Command("bash", r.benchScript(), "prep-release")
	cmd.Dir = r.Root
	cmd.Env = contract.ProcessEnv(r.Env, r.runEnv(map[string]string{stubPreflightSleepEnv: "60"}))
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Start(); err != nil {
		t.Fatalf("start prep-release: %v", err)
	}
	t.Cleanup(func() { _ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) })

	waitForStagedPreflight(t, r)
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGINT); err != nil {
		t.Fatalf("interrupt prep-release: %v", err)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		if err != nil {
			t.Fatalf("wait for prep-release: %v\n%s", err, output.String())
		}
		return 0
	case <-time.After(bounds.TestDeadline(0)):
		t.Fatalf("prep-release survived SIGINT past the test deadline\n%s", output.String())
		return -1
	}
}

// waitForStagedPreflight blocks until the preflight script has written its staging
// directory, which is the only moment at which a bypassed atomic promotion would be
// observable as a partial index.
func waitForStagedPreflight(t *testing.T, r shipRepo) {
	t.Helper()
	deadline := time.Now().Add(bounds.TestDeadline(0))
	for time.Now().Before(deadline) {
		staged, err := filepath.Glob(filepath.Join(r.Root, "dist", ".preflight.*", "release-index.json"))
		if err == nil && len(staged) > 0 {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	listing, _ := os.ReadDir(filepath.Join(r.Root, "dist"))
	t.Fatalf("preflight never staged its evidence; dist/ holds %v", listing)
}
