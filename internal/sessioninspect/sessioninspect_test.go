package sessioninspect

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
)

func TestInspectDeadlineWarnsAndReturnsZero(t *testing.T) {
	original := phases
	t.Cleanup(func() { phases = original })
	phases = []phase{func(ctx context.Context, _, _ io.Writer, _ string) int {
		<-ctx.Done()
		return 1
	}}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	var out bytes.Buffer
	if code := Inspect(ctx, &out, t.TempDir()); code != 0 {
		t.Fatalf("Inspect exit = %d, want 0", code)
	}
	const warning = "warning: bench session-inspect: deadline exceeded; session inspection stopped\n"
	if got := out.String(); got != warning {
		t.Fatalf("Inspect warning = %q, want %q", got, warning)
	}
}

func TestEnvironmentPhaseTE15StopsAtDiscoveryBound(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	if err := os.WriteFile(filepath.Join(bin, "bash"), []byte("#!/bin/sh\n/bin/sleep 30\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin)
	started := time.Now()
	if code := environmentPhase(context.Background(), io.Discard, io.Discard, root); code != 0 {
		t.Fatalf("environmentPhase exit = %d, want 0", code)
	}
	if elapsed := time.Since(started); elapsed >= bounds.EnvironmentDiscoveryTimeout+time.Second {
		t.Fatalf("environmentPhase elapsed = %s, want discovery stopped near %s", elapsed, bounds.EnvironmentDiscoveryTimeout)
	}
}

func TestPhaseFinishedHonorsCancellationAfterResult(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan int, 1)
	done <- 1
	cancel()

	if phaseFinished(ctx, done) {
		t.Fatal("phaseFinished = true after cancellation and result, want false")
	}
}

func TestCommandInstallsTenSecondDeadline(t *testing.T) {
	original := runInspect
	t.Cleanup(func() { runInspect = original })
	runInspect = func(ctx context.Context, _ io.Writer, _ string) int {
		deadline, ok := ctx.Deadline()
		if !ok {
			t.Fatal("Command handed Inspect a context without a deadline")
		}
		remaining := time.Until(deadline)
		if remaining < 9*time.Second || remaining > 10*time.Second {
			t.Fatalf("deadline remaining = %v, want concrete 10s timeout", remaining)
		}
		return 0
	}
	if code := Command(nil, io.Discard, io.Discard); code != 0 {
		t.Fatalf("Command exit = %d, want 0", code)
	}
}

func TestResumePhaseForwardsUnderlyingFailure(t *testing.T) {
	t.Chdir(t.TempDir())
	var stdout, stderr bytes.Buffer
	if code := resumePhase(context.Background(), &stdout, &stderr, t.TempDir()); code == 0 {
		t.Fatal("resumePhase exit = 0 outside a repository")
	}
	for _, want := range []string{
		"error: not in a git repository — run inside a Bench-linked repo",
		"warning: bench session-start: resume-clean failed; inspect retained worktree state",
	} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("resumePhase stderr missing %q:\n%s", want, stderr.String())
		}
	}
}
