package releasepreflight

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
)

// An external phase's git starts no auto-maintenance. The race phase runs the kit's
// tests, and a background prune can remove a worktree admin entry a test fixture planted.
func TestExternalPhaseGitStartsNoAutoMaintenance(t *testing.T) {
	t.Setenv("BENCH_PREFLIGHT_RACE", gittest.MaintenanceProbe(t))
	version, err := exec.Command("go", "env", "GOVERSION").Output()
	if err != nil {
		t.Fatal(err)
	}
	toolchain := strings.TrimPrefix(strings.TrimSpace(string(version)), "go")
	var stderr bytes.Buffer
	r := &runner{root: t.TempDir(), stderr: &stderr, identity: Identity{Toolchain: &toolchain}}

	code, err := r.runExternal(context.Background(), "race")

	if err != nil || code != 0 {
		t.Fatalf("maintenance probe = (%d, %v), want (0, nil); stderr=%q", code, err, stderr.String())
	}
}
