package testreport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// The helper half of TestCommandDrainsBuilderGroupOnCancelSignal runs as its own
// process, because the fact under test is what a signal does to the owner — and a
// signal sent to the test binary itself would take the assertions down with it.
const (
	cancelHelperEnv   = "BENCH_TEST_CANCEL_HELPER"
	cancelTempRootEnv = "BENCH_TEST_CANCEL_TEMPROOT"
	cancelSourceEnv   = "BENCH_TEST_CANCEL_SOURCE"
)

// TestCommandDrainsBuilderGroupOnCancelSignal grades what a cancel signal does to
// the builder group `bench test` detaches: a signal the owner does not trap leaves
// that group with no cleanup path at all. The stub builder parks with a descendant
// alive, which is the state the real `go build` holds for the length of a build.
//
// Each signal is its own subtest, so a set missing one signal reds only that leg.
func TestCommandDrainsBuilderGroupOnCancelSignal(t *testing.T) {
	if helper := os.Getenv(cancelHelperEnv); helper != "" {
		runCancelHelper(t)
		return
	}
	for _, signalCase := range []struct {
		name   string
		signal syscall.Signal
	}{
		{"SIGINT", syscall.SIGINT},
		{"SIGTERM", syscall.SIGTERM},
		{"SIGHUP", syscall.SIGHUP},
	} {
		t.Run(signalCase.name, func(t *testing.T) {
			tempRoot, descendant := signalCancelHelper(t, signalCase.signal)
			requireProcessGone(t, descendant)
			entries, err := filepath.Glob(filepath.Join(tempRoot, "bench-run-*"))
			if err != nil {
				t.Fatal(err)
			}
			if len(entries) != 0 {
				t.Errorf("cancelled selection left %v behind, want the private directory removed", entries)
			}
		})
	}
}

// signalCancelHelper starts the owner, waits for its builder to report a live
// descendant, signals it, and returns once the owner has exited. The returned pid
// is the descendant the drain has to have killed.
func signalCancelHelper(t *testing.T, signal syscall.Signal) (string, int) {
	t.Helper()
	tempRoot := t.TempDir()
	pidFile := filepath.Join(t.TempDir(), "builder-descendant")
	helper := exec.Command(os.Args[0], "-test.run=^TestCommandDrainsBuilderGroupOnCancelSignal$", "-test.timeout="+deadline().String())
	helper.Env = append(os.Environ(),
		cancelHelperEnv+"="+signal.String(),
		cancelTempRootEnv+"="+tempRoot,
		cancelSourceEnv+"="+parkingBuilderSource(t, pidFile),
	)
	output, err := helper.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	helper.Stderr = helper.Stdout
	if err := helper.Start(); err != nil {
		t.Fatal(err)
	}
	waited := make(chan error, 1)
	go func() { _, _ = io.ReadAll(output); waited <- helper.Wait() }()
	t.Cleanup(func() { _ = helper.Process.Kill() })

	descendant := waitForPID(t, pidFile)
	t.Cleanup(func() { _ = syscall.Kill(descendant, syscall.SIGKILL) })
	if err := helper.Process.Signal(signal); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-waited:
		if err != nil && !errors.As(err, new(*exec.ExitError)) {
			t.Fatalf("cancel helper: %v", err)
		}
	case <-time.After(deadline()):
		t.Fatalf("cancel helper did not exit within %s of %s", deadline(), signal)
	}
	return tempRoot, descendant
}

// runCancelHelper is the owner process. It drives the real Command through a
// factory whose builder is the parking script, so the cancellation path under test
// is production's: Command's signal context cancels canonicalBuild, which drains
// the group it detached.
func runCancelHelper(t *testing.T) {
	factory := runbinary.Factory{
		TempRoot: os.Getenv(cancelTempRootEnv),
		Verify:   func(string, string) error { return nil },
	}
	source := os.Getenv(cancelSourceEnv)
	previous := selectRunBinary
	selectRunBinary = func(ctx context.Context, _ string) (*runbinary.Selection, error) {
		return factory.Own(ctx, source)
	}
	t.Cleanup(func() { selectRunBinary = previous })
	out, code := Command(t.TempDir(), []string{"./..."})
	fmt.Print(out, code)
}

// parkingBuilderSource writes a source root whose builder reports one live
// descendant and then parks, holding the cancellation window open. The descendant
// is an ordinary background child, so it belongs to the detached group and only a
// group-wide signal reaches it.
func parkingBuilderSource(t *testing.T, pidFile string) string {
	t.Helper()
	source := t.TempDir()
	if err := os.MkdirAll(filepath.Join(source, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	script := `#!/usr/bin/env bash
sleep ` + strconv.Itoa(int(parkSeconds())) + ` &
printf %s "$!" > "` + pidFile + `.partial"
mv "` + pidFile + `.partial" "` + pidFile + `"
wait
`
	if err := os.WriteFile(filepath.Join(source, "scripts", "go-build.sh"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return source
}

func waitForPID(t *testing.T, pidFile string) int {
	t.Helper()
	expiry := time.Now().Add(deadline())
	for time.Now().Before(expiry) {
		raw, err := os.ReadFile(pidFile)
		if err == nil {
			pid, convErr := strconv.Atoi(strings.TrimSpace(string(raw)))
			if convErr == nil && pid > 0 {
				return pid
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("builder never reported a descendant within %s", deadline())
	return 0
}

// requireProcessGone polls rather than sampling once: the drain kills the group and
// the kernel reaps asynchronously, so a single probe grades scheduler lag.
func requireProcessGone(t *testing.T, pid int) {
	t.Helper()
	expiry := time.Now().Add(deadline())
	for time.Now().Before(expiry) {
		if err := syscall.Kill(pid, 0); errors.Is(err, syscall.ESRCH) {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Errorf("builder descendant %d survived the cancelled run", pid)
}

// deadline outlasts the window it waits on — the grace canonicalBuild gives a
// cancelled group before escalating to SIGKILL.
func deadline() time.Duration { return bounds.TestDeadline(runbinary.BuilderCancelGrace) }

// parkSeconds keeps the builder's descendant alive well past the deadline, so a
// surviving process is always a drain failure and never the child exiting on time.
func parkSeconds() int64 { return int64(deadline()/time.Second) * 4 }
