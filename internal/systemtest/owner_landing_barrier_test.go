//go:build system

package systemtest

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/testrepo"
)

// TestLandingRaceBarrierEndsWithoutItsPartner shows that the barrier bound bites. Each case
// leaves one partner step undone, and the barrier must end inside the deadline that its
// window derives.
func TestLandingRaceBarrierEndsWithoutItsPartner(t *testing.T) {
	t.Parallel()
	for _, partner := range []struct {
		name       string
		readsReady bool
	}{{"no ready reader", false}, {"no release writer", true}} {
		t.Run(partner.name, func(t *testing.T) {
			t.Parallel()
			ready, release := landingBarrier(t, t.TempDir())
			if partner.readsReady {
				// A nonblocking read open of a FIFO succeeds with no writer, so the barrier's
				// ready signal completes and only the release stays undone.
				reader, err := os.OpenFile(ready, os.O_RDONLY|syscall.O_NONBLOCK, 0)
				if err != nil {
					t.Fatal(err)
				}
				defer reader.Close()
			}
			window := time.Second
			f := testrepo.NewGateFixture(t.TempDir())
			script := filepath.Join(t.TempDir(), "barrier.sh")
			if err := os.WriteFile(script, []byte(f.Script("set -eu\n"+landingRaceBarrier(f, window))), 0o755); err != nil {
				t.Fatal(err)
			}
			f.MustWrite(t, t.TempDir(), "", "")
			barrier, err := owner.startProcessGroup(t.TempDir(), []string{"LAND_RACE_READY=" + ready, "LAND_RACE_RELEASE=" + release}, script)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = barrier.kill(syscall.SIGKILL) })
			if code := exitWithin(t, barrier.cmd, window); code != 124 {
				t.Fatalf("barrier with %s exit = %d, want the timeout exit 124", partner.name, code)
			}
		})
	}
}
