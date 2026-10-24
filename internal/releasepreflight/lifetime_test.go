package releasepreflight

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/env"
	"github.com/gibbonmi/bench/internal/releaseevidence"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// TestMain installs the kit-test and release-evidence rows of the composition root's table, so
// a phase protects its private run and the release evidence as the CLI does. The test owns the
// only vector, so an outer run's resources stay out of its launches.
func TestMain(m *testing.M) {
	if err := os.Unsetenv(protection.Env); err != nil {
		os.Exit(3)
	}
	protection.Install(protection.Protocol{Validators: protection.Validators{
		protection.KindKitTest:         env.ValidateProtection,
		protection.KindReleaseEvidence: releaseevidence.ValidateProtection,
	}, ReadLimit: bounds.ControlRecordLimit})
	os.Exit(m.Run())
}

// TestExternalIncomplete is PL62. An external phase exits zero and leaves a live member of its
// group. Its raw exit cannot certify the phase, and promotion of its record refuses while that
// member keeps the release evidence protected.
func TestExternalIncomplete(t *testing.T) {
	t.Setenv("TMPDIR", t.TempDir())
	root := t.TempDir()
	for _, input := range []string{"bin/bench.sh", ".bench/gate.sh"} {
		path := filepath.Join(root, filepath.FromSlash(input))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	pidFile := filepath.Join(t.TempDir(), "member")
	t.Setenv("BENCH_PREFLIGHT_GATE", phaseScript(t, lingeringMember(pidFile)))
	var stderr bytes.Buffer
	r := &runner{root: root, mode: ModeVerify, stderr: &stderr, evidence: protectEvidence(t, root)}

	result := r.runPhase(context.Background(), "gate")
	processtest.WaitPID(t, pidFile)
	if result.Status != StatusRed || result.Failure == nil || !strings.Contains(result.Failure.Message, "process protection "+string(protection.Protected)) {
		t.Fatalf("phase result = %+v, want red on its protected group; stderr %q", result, stderr.String())
	}
	err := PromoteEvidence(root, ModeVerify, []Result{result}, Manifest{SchemaVersion: 1, Mode: ModeVerify})
	if !errors.Is(err, protection.ErrRetained) {
		t.Fatalf("promotion = %v, want the retained-evidence refusal", err)
	}
	if _, err := os.Stat(filepath.Join(root, "dist", "preflight")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("promoted evidence = %v, want none while the phase member lives", err)
	}
}

// TestScannerCancel is PL63. A scanner descendant logs each TERM and survives it, so the log
// shows the TERM that opened the shutdown. Its absence when the phase returns shows that KILL
// followed the grace.
func TestScannerCancel(t *testing.T) {
	marks := t.TempDir()
	signals, ready := filepath.Join(marks, "signals"), filepath.Join(marks, "descendant")
	t.Setenv("BENCH_PREFLIGHT_VULNERABILITY", phaseScript(t, `(
  exec >/dev/null 2>&1
  trap 'echo TERM >> `+sanitize.ShellQuote(signals)+`' TERM
  echo $BASHPID > `+sanitize.ShellQuote(ready+".partial")+` && mv `+sanitize.ShellQuote(ready+".partial")+` `+sanitize.ShellQuote(ready)+`
  while :; do sleep .01; done
) &
wait
`))
	var stderr bytes.Buffer
	r := &runner{root: t.TempDir(), mode: ModeVerify, stderr: &stderr}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan Result, 1)
	go func() { done <- r.runPhase(ctx, "vulnerability") }()
	descendant := processtest.WaitPID(t, ready)
	pgid, err := syscall.Getpgid(descendant)
	if err != nil {
		t.Fatal(err)
	}
	cancelled := time.Now()
	cancel()
	window := bounds.TestDeadline(bounds.ProcessShutdownGrace + bounds.ProcessFinalWindow)
	var result Result
	select {
	case result = <-done:
	case <-time.After(window):
		t.Fatal(bounds.TestTimeoutVerdict("the cancelled scanner to return", window))
	}
	elapsed := time.Since(cancelled)
	if result.Status != StatusInterrupted {
		t.Fatalf("cancelled scanner result = %+v, want interrupted", result)
	}
	if logged, readErr := os.ReadFile(signals); readErr != nil || !strings.HasPrefix(string(logged), "TERM\n") {
		t.Fatalf("scanner signal log = %q (%v), want TERM before KILL", logged, readErr)
	}
	if elapsed < bounds.ProcessShutdownGrace {
		t.Fatalf("cancelled scanner returned after %v, want the full %v grace before KILL", elapsed, bounds.ProcessShutdownGrace)
	}
	if state := protection.ObserveGroup(pgid).State; state != protection.GroupAbsent {
		t.Fatalf("scanner group after the phase returned is %s, want absent", state)
	}
}

// TestScannerZeroUncertified pins that a scanner's exit zero and complete output cannot certify
// its findings while a member of its group outlives it. The phase fails on its protection.
func TestScannerZeroUncertified(t *testing.T) {
	pidFile := filepath.Join(t.TempDir(), "member")
	t.Setenv("BENCH_PREFLIGHT_VULNERABILITY", phaseScript(t, "printf '{}\\n'\n"+lingeringMember(pidFile)))
	root := t.TempDir()
	var stderr bytes.Buffer
	r := &runner{root: root, mode: ModeVerify, stderr: &stderr, evidence: protectEvidence(t, root)}

	result := r.runPhase(context.Background(), "vulnerability")
	processtest.WaitPID(t, pidFile)
	if result.Status != StatusRed || result.Failure == nil || !strings.Contains(result.Failure.Message, "process protection "+string(protection.Protected)) {
		t.Fatalf("scanner result = %+v, want red on its protected group", result)
	}
}

// phaseScript writes body as an executable bash script and returns its path.
func phaseScript(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "phase")
	if err := os.WriteFile(path, []byte("#!/usr/bin/env bash\n"+body), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

// lingeringMember is a script tail that leaves a live member in the script's group and
// publishes that member's PID in pidFile.
func lingeringMember(pidFile string) string {
	partial := sanitize.ShellQuote(pidFile + ".partial")
	return "sleep 300 </dev/null >/dev/null 2>&1 &\necho $! >" + partial + " && mv " + partial + " " + sanitize.ShellQuote(pidFile) + "\n"
}

// protectEvidence opens the release-evidence generation of root as Command does, and it removes
// that generation's store, which outlives the temporary checkout.
func protectEvidence(t *testing.T, root string) protection.Generation {
	t.Helper()
	evidence, err := releaseevidence.Protect(root)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := evidence.Discard(); err != nil {
			t.Error(err)
		}
	})
	return evidence
}
