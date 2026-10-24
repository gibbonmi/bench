package runbinary

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// TestBuilderNormalDrain is PL34. The builder returns while its parked descendant waits on
// the selection, so Own must drain that group before it returns the executable.
func TestBuilderNormalDrain(t *testing.T) {
	t.Parallel()
	builder := newParkingBuilder(t, false)
	factory := Factory{TempRoot: t.TempDir(), Protection: binaryProtocol(), Verify: func(string, string) error { return nil }}
	selection, err := factory.Own(context.Background(), builder.source)
	if err != nil {
		t.Fatal(err)
	}
	child := readPID(t, builder.childPID)
	t.Cleanup(func() { _ = syscall.Kill(child, syscall.SIGKILL) })
	if state := protection.ObserveGroup(readPID(t, builder.builderPID)).State; state != protection.GroupAbsent {
		t.Fatalf("builder group at Own's return is %s, want absent", state)
	}
	store := selection.generation.Resource.Store
	if err := selection.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Dir(selection.Path)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("selection directory after Close = %v, want absent", err)
	}
	if _, err := os.Lstat(store); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("protection store after Close = %v, want absent", err)
	}
	if _, err := os.Stat(builder.premature); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("builder descendant observed premature cleanup: %v", err)
	}
}

// TestBuilderCancel is PL35. A descendant logs each TERM and survives it, so the log shows a
// TERM that opened the shutdown. Its absence after the grace shows that KILL followed.
func TestBuilderCancel(t *testing.T) {
	t.Parallel()
	source, marks := t.TempDir(), t.TempDir()
	signals, ready := filepath.Join(marks, "signals"), filepath.Join(marks, "descendant")
	if err := os.MkdirAll(filepath.Join(source, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	script := `#!/usr/bin/env bash
(
  exec >/dev/null 2>&1
  trap 'echo TERM >> "` + signals + `"' TERM
  echo $BASHPID > "` + ready + `.partial" && mv "` + ready + `.partial" "` + ready + `"
  while :; do sleep .01; done
) &
wait
`
	if err := os.WriteFile(filepath.Join(source, "scripts", "go-build.sh"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() {
		_, err := Factory{TempRoot: t.TempDir(), Protection: binaryProtocol()}.Own(ctx, source)
		done <- err
	}()
	descendant := processtest.WaitPID(t, ready)
	pgid, err := syscall.Getpgid(descendant)
	if err != nil {
		t.Fatal(err)
	}
	cancelled := time.Now()
	cancel()
	window := bounds.TestDeadline(BuilderCancelGrace + bounds.ProcessFinalWindow)
	select {
	case err = <-done:
	case <-time.After(window):
		t.Fatal(bounds.TestTimeoutVerdict("the cancelled build to return", window))
	}
	elapsed := time.Since(cancelled)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled Own = %v, want the cancellation cause", err)
	}
	if logged, readErr := os.ReadFile(signals); readErr != nil || !strings.HasPrefix(string(logged), "TERM\n") {
		t.Fatalf("descendant signal log = %q (%v), want TERM before KILL", logged, readErr)
	}
	if elapsed < BuilderCancelGrace {
		t.Fatalf("cancelled build returned after %v, want the full %v grace before KILL", elapsed, BuilderCancelGrace)
	}
	if state := protection.ObserveGroup(pgid).State; state != protection.GroupAbsent {
		t.Fatalf("builder group after the grace is %s, want absent", state)
	}
}
