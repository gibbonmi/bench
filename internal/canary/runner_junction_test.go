package canary

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
)

// TestSweepTierRunsPlantedBashGate is the one junction test that drives SweepTier with the
// real defaultRunner instead of a fake: resolvingRunner (canary_path_test.go) reproduces
// bash's cwd-relative resolution in-process, which proves the sweep's own call-shaping logic
// but never proves real bash agrees with that reproduction. This fixture is a hermetic gate
// planted under the fixture's own temp root, so the assertion is that the sweep's real exit
// verdict — from the actual subprocess, not a stand-in — reaches the sweep result.
func TestSweepTierRunsPlantedBashGate(t *testing.T) {
	root := t.TempDir()
	family := mappedFamily(t)
	fixture := canaryFixture(root, family, "gate-junction")
	mkdir(t, filepath.Join(fixture, "files"))
	write(t, filepath.Join(fixture, "EXPECT"), "targeted regression\n")
	write(t, filepath.Join(fixture, "files", "marker.txt"), "x\n")
	plantGate(t, root, "echo \"gate: running\"\nif [ -f marker.txt ]; then echo \"targeted regression\"; exit 1; fi\nexit 0\n")

	if err := SweepTier(root, registry.Dev, defaultRunner); err != nil {
		t.Fatalf("SweepTier err = %v, want nil: the planted gate's real exit verdict should clear a matching fixture", err)
	}
}

// TestSweepTierSurfacesMissingGateExit127 points the same real-runner sweep at a fixture
// root with no gate.sh planted at all. bash's own "no such file" failure is what the sweep
// result has to surface, not a fake standing in for it, so the gate is left entirely absent
// rather than swapped for a broken one.
func TestSweepTierSurfacesMissingGateExit127(t *testing.T) {
	root := t.TempDir()
	family := mappedFamily(t)
	fixture := canaryFixture(root, family, "missing-gate-junction")
	mkdir(t, filepath.Join(fixture, "files"))
	write(t, filepath.Join(fixture, "EXPECT"), "targeted regression\n")
	write(t, filepath.Join(fixture, "files", "marker.txt"), "x\n")

	err := SweepTier(root, registry.Dev, defaultRunner)
	if err == nil {
		t.Fatal("SweepTier err = nil, want the missing gate's exit-127 diagnostic to surface")
	}
	if !strings.Contains(err.Error(), "got exit 127") {
		t.Fatalf("SweepTier err = %v, want it to report the real bash exit-127 verdict", err)
	}
}

// plantGate writes root's inner gate as a real, hermetic bash script: body runs with the
// materialized fixture (or an empty baseline) as its cwd, so it is free to test for the
// fixture's own marker files by relative name.
func plantGate(t *testing.T, root, body string) {
	t.Helper()
	write(t, filepath.Join(root, ".bench", "gate.sh"), "#!/usr/bin/env bash\n"+body)
}
