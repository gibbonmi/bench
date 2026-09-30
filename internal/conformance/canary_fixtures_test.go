package conformance

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
)

func canaryFixturePaths(t *testing.T, fixturesDir string) map[string]canary.Fixture {
	t.Helper()
	families, err := os.ReadDir(fixturesDir)
	if err != nil {
		t.Fatalf("read canary fixtures: %v", err)
	}
	for _, family := range families {
		if !family.IsDir() {
			continue
		}
		// A family is canonical when inventory resolution can identify it, through a
		// family binding or through its own fixture marker. Only a family with neither
		// is unattributable. A flat fixture carries its own binding, not a family.
		if canary.IsConformanceFamily(filepath.Join(fixturesDir, family.Name())) && !familyIsBound(family.Name()) {
			t.Errorf("canary family %q is not canonical", family.Name())
		}
	}
	// Inventory discovery enumerates fixtures and enforces base-name uniqueness. A second
	// walk here would disagree with the producer the check consumes.
	discovered, err := canary.Fixtures(fixturesDir)
	if err != nil {
		t.Fatalf("walk canary fixtures: %v", err)
	}
	return discovered
}
