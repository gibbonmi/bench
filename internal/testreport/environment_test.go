package testreport

import (
	"os"
	"os/exec"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// A focused run of the kit's tests starts git with no auto-maintenance. A background
// prune from that maintenance can remove a worktree admin entry a test fixture planted.
func TestKitRunEnvironmentGitStartsNoAutoMaintenance(t *testing.T) {
	kit := t.TempDir()
	child, err := selectedRunEnvironment(os.Environ(), kit, &runbinary.Selection{Path: "/selected/bench", SourceRoot: kit})
	if err != nil {
		t.Fatal(err)
	}
	probe := exec.Command(gittest.MaintenanceProbe(t))
	probe.Env = child
	if out, err := probe.CombinedOutput(); err != nil {
		t.Fatalf("maintenance probe = %v: %s", err, out)
	}
}
