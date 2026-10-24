package protection

import (
	"errors"
	"os/exec"
	"syscall"
	"testing"
)

// startGroup starts a real leader in its own process group and reaps it at cleanup.
func startGroup(t *testing.T) *exec.Cmd {
	t.Helper()
	cmd := exec.Command("sleep", "30")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		_ = cmd.Wait()
	})
	return cmd
}

func TestObserveGroupAbsent(t *testing.T) {
	cmd := startGroup(t)
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	_ = cmd.Wait()
	if got := ObserveGroup(cmd.Process.Pid); got.State != GroupAbsent {
		t.Fatalf("reaped group observed %+v, want absent", got)
	}
}

func TestObserveGroupPresent(t *testing.T) {
	cmd := startGroup(t)
	if got := ObserveGroup(cmd.Process.Pid); got.State != GroupPresent {
		t.Fatalf("live group observed %+v, want present", got)
	}
}

func TestObserveGroupPermission(t *testing.T) {
	requireUncertain(t, syscall.EPERM)
}

func TestObserveGroupUnknownError(t *testing.T) {
	requireUncertain(t, syscall.EINVAL)
}

// requireUncertain feeds one probe failure that a real group cannot produce on demand.
func requireUncertain(t *testing.T, failure error) {
	t.Helper()
	got := observeGroup(func(int, syscall.Signal) error { return failure }, 4242)
	if got.State != GroupUncertain || !errors.Is(got.Err, failure) {
		t.Fatalf("probe failure %v observed %+v, want uncertain with the failure", failure, got)
	}
}
