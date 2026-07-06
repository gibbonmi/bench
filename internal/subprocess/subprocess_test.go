package subprocess

import (
	"os/exec"
	"strings"
	"testing"
)

func TestCaptureSeparatesStreams(t *testing.T) {
	r := Capture(exec.Command("sh", "-c", "echo out; echo err 1>&2"))
	if r.ExitCode != 0 || r.Err != nil {
		t.Fatalf("exit=%d err=%v, want 0/nil", r.ExitCode, r.Err)
	}
	if strings.TrimSpace(r.Stdout) != "out" {
		t.Errorf("stdout = %q, want %q", r.Stdout, "out")
	}
	if strings.TrimSpace(r.Stderr) != "err" {
		t.Errorf("stderr = %q, want %q", r.Stderr, "err")
	}
}

func TestCapturePropagatesExitCode(t *testing.T) {
	r := Capture(exec.Command("sh", "-c", "exit 3"))
	if r.ExitCode != 3 {
		t.Errorf("exit = %d, want 3", r.ExitCode)
	}
	if r.Err == nil {
		t.Error("Err = nil, want non-nil for a failed process")
	}
}

func TestCaptureSpawnFailureIsOne(t *testing.T) {
	r := Capture(exec.Command("bench-no-such-binary"))
	if r.ExitCode != 1 {
		t.Errorf("exit = %d, want 1 for a spawn failure", r.ExitCode)
	}
	if r.Err == nil {
		t.Error("Err = nil, want the spawn error")
	}
}

func TestCaptureMergedInterleavesIntoStdout(t *testing.T) {
	r := CaptureMerged(exec.Command("sh", "-c", "echo out; echo err 1>&2"))
	if r.ExitCode != 0 {
		t.Fatalf("exit = %d, want 0", r.ExitCode)
	}
	if !strings.Contains(r.Stdout, "out") || !strings.Contains(r.Stdout, "err") {
		t.Errorf("merged stdout = %q, want both streams", r.Stdout)
	}
	if r.Stderr != "" {
		t.Errorf("merged Stderr = %q, want empty", r.Stderr)
	}
}
