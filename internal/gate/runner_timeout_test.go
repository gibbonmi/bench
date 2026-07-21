package gate

import (
	"bytes"
	"context"
	"io"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestExecuteDeadlineRecordsDistinctTimeout(t *testing.T) {
	root := gateTestRepo(t, "#!/bin/sh\nsleep 30 &\necho $! > .git/timeout-child\nwait\n", "")
	old := gateTimeout
	gateTimeout = 50 * time.Millisecond
	t.Cleanup(func() { gateTimeout = old })
	parent, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	var stderr bytes.Buffer
	got := Execute(parent, root, io.Discard, &stderr)
	inspection := Inspect(root)
	if got.ActionExit != 124 || got.GateExit != 124 || inspection.State != Ready || inspection.Status != "timeout" || inspection.ReusableGreen {
		t.Fatalf("result=%+v inspection=%+v stderr=%q", got, inspection, stderr.String())
	}
	if !strings.Contains(stderr.String(), "gate: timeout") {
		t.Fatalf("timeout message missing: %q", stderr.String())
	}
	child := waitForPIDFile(t, filepath.Join(root, ".git", "timeout-child"))
	t.Cleanup(func() { _ = syscall.Kill(child, syscall.SIGKILL) })
	waitForProcessExit(t, child)
}

func TestExecuteHealthyGateJustBelowDeadlineRecordsOrdinaryGreen(t *testing.T) {
	root := gateTestRepo(t, "#!/bin/sh\nsleep 0.6\nexit 0\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	old := gateTimeout
	gateTimeout = time.Second
	t.Cleanup(func() { gateTimeout = old })
	var stderr bytes.Buffer
	got := Execute(context.Background(), root, io.Discard, &stderr)
	inspection := Inspect(root)
	if got.ActionExit != 0 || got.GateExit != 0 || inspection.State != Ready || inspection.Status != "green" || !inspection.ReusableGreen {
		t.Fatalf("near-boundary result=%+v inspection=%+v stderr=%q", got, inspection, stderr.String())
	}
}
