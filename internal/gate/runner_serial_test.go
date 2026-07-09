package gate

// The serial-phase contract: a Serial phase completes before any concurrent phase
// starts regardless of table position, its red fails the run fast, and an interrupt
// during it is a 130, never a graded red — in both outer and inner modes.

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestRunnerSerialPhaseCompletesBeforeConcurrentPhasesStart(t *testing.T) {
	root := t.TempDir()
	marker := filepath.Join(root, "built")
	phases := []Phase{
		{
			Name:   "build",
			Argv:   []string{"bash", "-c", `sleep 0.1; touch "$1"`, "bash", marker},
			Serial: true,
		},
		// Each reader phase goes red unless the serial phase already finished:
		// a runner that launches everything concurrently starts these ~100ms
		// before the marker exists.
		{Name: "alpha", Argv: []string{"bash", "-c", `test -f "$1"`, "bash", marker}},
		{Name: "bravo", Argv: []string{"bash", "-c", `test -f "$1"`, "bash", marker}},
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phases, outerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("runPhases rc = %d; a reader phase started before the serial phase finished\nstdout=%q\nstderr=%q", rc, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "phase build: green") {
		t.Fatalf("serial phase missing from summaries:\n%s", stdout.String())
	}
}

func TestRunnerSerialRedFailsFast(t *testing.T) {
	root := t.TempDir()
	leak := filepath.Join(root, "leak")
	phases := []Phase{
		{Name: "build", Argv: []string{"bash", "-c", "exit 3"}, Serial: true},
		{Name: "alpha", Argv: []string{"bash", "-c", `touch "$1"`, "bash", leak}},
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phases, outerMode, &stdout, &stderr)
	if rc != 1 {
		t.Fatalf("runPhases rc = %d, want 1", rc)
	}
	out := stdout.String() + stderr.String()
	if !strings.Contains(out, "phase build: red (exit 3)") || !strings.Contains(out, "gate: red") {
		t.Fatalf("serial red not attributed:\n%s", out)
	}
	if _, err := os.Stat(leak); err == nil {
		t.Fatalf("phase after failed serial phase still ran (it would grade a stale or absent binary)")
	}
}

func TestRunnerSerialRedFailsFastInnerMode(t *testing.T) {
	root := t.TempDir()
	leak := filepath.Join(root, "leak")
	phases := []Phase{
		{Name: "build", Argv: []string{"bash", "-c", "exit 3"}, Serial: true},
		{Name: "alpha", Argv: []string{"bash", "-c", `touch "$1"`, "bash", leak}},
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phases, innerMode, &stdout, &stderr)
	if rc != 1 {
		t.Fatalf("runPhases rc = %d, want 1", rc)
	}
	if !strings.HasSuffix(stderr.String(), "gate: red\n") {
		t.Fatalf("stderr final line = %q, want gate: red", stderr.String())
	}
	if _, err := os.Stat(leak); err == nil {
		t.Fatalf("inner-mode phase after failed serial phase still ran")
	}
}

func TestRunnerCancelDuringSerialPhaseReturns130(t *testing.T) {
	root := t.TempDir()
	pidfile := filepath.Join(root, "sleep.pid")
	leak := filepath.Join(root, "leak")
	phases := []Phase{
		{
			Name:   "build",
			Argv:   []string{"bash", "-c", `sleep 30 & echo $! > "$1"; wait`, "bash", pidfile},
			Serial: true,
		},
		{Name: "alpha", Argv: []string{"bash", "-c", `touch "$1"`, "bash", leak}},
	}
	ctx, cancel := context.WithCancel(context.Background())
	var stdout, stderr bytes.Buffer
	done := make(chan int, 1)
	go func() {
		done <- runPhases(ctx, root, phases, outerMode, &stdout, &stderr)
	}()

	pid := waitForPIDFile(t, pidfile)
	t.Cleanup(func() { _ = syscall.Kill(pid, syscall.SIGKILL) })
	cancel()

	select {
	case rc := <-done:
		if rc != 130 {
			t.Fatalf("runPhases rc = %d, want 130; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
		}
	case <-time.After(5 * time.Second):
		t.Fatal("runPhases did not return after cancellation during the serial phase")
	}
	waitForProcessExit(t, pid)
	// An interrupt is not a verdict: no phase behind the interrupted build may run,
	// and the run must not present itself as a graded red.
	if _, err := os.Stat(leak); err == nil {
		t.Fatalf("concurrent phase ran after the serial phase was interrupted")
	}
	if strings.Contains(stdout.String()+stderr.String(), "gate: red") {
		t.Fatalf("interrupted serial phase was graded as red:\nstdout=%q\nstderr=%q", stdout.String(), stderr.String())
	}
}

func TestRunnerSerialPhaseNotFirstStillRunsFirst(t *testing.T) {
	root := t.TempDir()
	marker := filepath.Join(root, "built")
	phases := []Phase{
		// Readers listed ahead of the serial phase: only the runner's own
		// reordering keeps them from starting before the marker exists.
		{Name: "alpha", Argv: []string{"bash", "-c", `test -f "$1"`, "bash", marker}},
		{Name: "bravo", Argv: []string{"bash", "-c", `test -f "$1"`, "bash", marker}},
		{
			Name:   "build",
			Argv:   []string{"bash", "-c", `sleep 0.1; touch "$1"`, "bash", marker},
			Serial: true,
		},
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phases, outerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("runPhases rc = %d; a reader phase ran before the non-first serial phase\nstdout=%q\nstderr=%q", rc, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "phase build: green") {
		t.Fatalf("serial phase missing from summaries:\n%s", stdout.String())
	}
}

func TestRunnerSerialPhaseNotFirstInnerMode(t *testing.T) {
	t.Run("runs first", func(t *testing.T) {
		root := t.TempDir()
		marker := filepath.Join(root, "built")
		phases := []Phase{
			// A reader listed ahead of the serial phase: sequential execution in
			// table order would run it before the marker exists.
			{Name: "alpha", Argv: []string{"bash", "-c", `test -f "$1"`, "bash", marker}},
			{
				Name:   "build",
				Argv:   []string{"bash", "-c", `touch "$1"`, "bash", marker},
				Serial: true,
			},
		}

		var stdout, stderr bytes.Buffer
		rc := runPhases(context.Background(), root, phases, innerMode, &stdout, &stderr)
		if rc != 0 {
			t.Fatalf("runPhases rc = %d; the reader ran before the non-first serial phase\nstdout=%q\nstderr=%q", rc, stdout.String(), stderr.String())
		}
	})

	t.Run("red fails fast", func(t *testing.T) {
		root := t.TempDir()
		leak := filepath.Join(root, "leak")
		phases := []Phase{
			{Name: "alpha", Argv: []string{"bash", "-c", `touch "$1"`, "bash", leak}},
			{Name: "build", Argv: []string{"bash", "-c", "exit 3"}, Serial: true},
		}

		var stdout, stderr bytes.Buffer
		rc := runPhases(context.Background(), root, phases, innerMode, &stdout, &stderr)
		if rc != 1 {
			t.Fatalf("runPhases rc = %d, want 1", rc)
		}
		if _, err := os.Stat(leak); err == nil {
			t.Fatalf("phase listed ahead of the failed serial phase still ran")
		}
	})
}
