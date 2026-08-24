//go:build system

package systemtest

import (
	"errors"
	"syscall"
	"testing"
	"time"
)

func TestRedTimeoutAndDescendantTeardown(t *testing.T) {
	if red := owner.runSelected(owner.repos[1], "not-a-command"); red.code == 0 {
		t.Fatal("unknown command unexpectedly green")
	}
	owner.markTerminal("red")
	interrupted, err := owner.interruptProcessGroup()
	if err != nil {
		t.Fatal(err)
	}
	requireProcessGone(t, interrupted, "interrupt")
	owner.markTerminal("interrupt")
	descendant, err := owner.timeoutProcessGroup(100 * time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	requireProcessGone(t, descendant, "timeout")
	owner.markTerminal("timeout")
}

func requireProcessGone(t *testing.T, pid int, outcome string) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for {
		err := syscall.Kill(pid, 0)
		if errors.Is(err, syscall.ESRCH) {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("descendant %d remains after process-group %s", pid, outcome)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
