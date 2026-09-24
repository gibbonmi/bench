package sessioninspect

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/intent"
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

// phaseIndex returns the position of want in the inspection sequence, or -1.
func phaseIndex(want phase) int {
	for index, run := range phases {
		if reflect.ValueOf(run).Pointer() == reflect.ValueOf(want).Pointer() {
			return index
		}
	}
	return -1
}

// LE76: the sequence runs the recovery pass right after the resume phase, and the pass
// abandons a lease-less entry whose owner process is gone.
func TestInspectRecoversAfterTheResumePhase(t *testing.T) {
	if resume, recovery := phaseIndex(resumePhase), phaseIndex(recoveryPhase); resume < 0 || recovery != resume+1 {
		t.Fatalf("resume phase %d, recovery phase %d, want recovery right after resume", resume, recovery)
	}
	root := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", root).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	t.Setenv("BENCH_HOME", t.TempDir())
	// A reaped child's process id names no live process.
	child := exec.Command("true")
	if err := child.Run(); err != nil {
		t.Fatal(err)
	}
	entry := intent.EntryOwnedBy(intent.KindShift, child.Process.Pid)
	if err := intent.Upsert(root, entry); err != nil {
		t.Fatal(err)
	}
	original := phases
	t.Cleanup(func() { phases = original })
	phases = []phase{recoveryPhase}
	var out bytes.Buffer
	Inspect(context.Background(), &out, root)
	current, err := intent.Read(root)
	if err != nil || len(current.Entries) != 1 || current.Entries[0].Outcome == "" {
		t.Fatalf("ledger after Inspect = %+v (%v), want the entry abandoned", current.Entries, err)
	}
}
