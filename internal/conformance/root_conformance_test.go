package conformance

import (
	"os"
	"testing"
)

func TestRootConformance(t *testing.T) {
	root := os.Getenv("BENCH_CONFORMANCE_ROOT")
	if root == "" {
		t.Skip("BENCH_CONFORMANCE_ROOT not set")
	}
	h := NewHarness(t)
	for _, diag := range RunConformance(root, h.KitRoot) {
		t.Errorf("gate: %s", diag)
	}
}
