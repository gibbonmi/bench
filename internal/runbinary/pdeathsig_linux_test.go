//go:build linux

package runbinary

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
)

// pdeathsigOwnerSource names the parking-builder source root a re-executed owner builds
// from. Its presence turns the helper below into the owner process the test kills.
const pdeathsigOwnerSource = "BENCH_RUNBINARY_PDEATHSIG_SOURCE"

// TestPdeathsigOwnerHelperProcess is the owner half of the test after it. Without the
// variable it asserts nothing, so an ordinary package run passes straight over it.
func TestPdeathsigOwnerHelperProcess(t *testing.T) {
	source := os.Getenv(pdeathsigOwnerSource)
	if source == "" {
		return
	}
	factory := Factory{TempRoot: source, Verify: func(string, string) error { return nil }}
	if _, err := factory.Own(context.Background(), source); err != nil {
		t.Fatal(err)
	}
}

// TestBuilderChildDiesWithAnOwnerThatNeverDrains is LQ24. An owner killed mid-build never
// reaches drainBuilderGroup, so the parent-death signal is the only thing left that
// removes the builder child. The owner has to be its own process: the signal fires on the
// death of the process that started the child, and the observer must outlive it.
func TestBuilderChildDiesWithAnOwnerThatNeverDrains(t *testing.T) {
	builder := newParkingBuilder(t, true)
	owner := exec.Command(os.Args[0], "-test.run=^TestPdeathsigOwnerHelperProcess$")
	owner.Env = append(os.Environ(), pdeathsigOwnerSource+"="+builder.source)
	if err := owner.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = owner.Process.Kill()
		_ = owner.Wait()
	})
	child := awaitPID(t, builder.builderPID)
	t.Cleanup(func() { _ = syscall.Kill(-child, syscall.SIGKILL) })
	if err := owner.Process.Signal(syscall.SIGKILL); err != nil {
		t.Fatal(err)
	}
	requireProcessExit(t, child)
}

// awaitPID waits for path to hold the process id the parking builder records. The window
// derives from the grace a builder group gets, so a loaded machine cannot beat the wait.
func awaitPID(t *testing.T, path string) int {
	t.Helper()
	window := bounds.TestDeadline(BuilderCancelGrace)
	deadline := time.Now().Add(window)
	for time.Now().Before(deadline) {
		if data, err := os.ReadFile(path); err == nil && strings.TrimSpace(string(data)) != "" {
			return readPID(t, path)
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal(bounds.TestTimeoutVerdict("the builder child to record its process id at "+path, window))
	return 0
}
