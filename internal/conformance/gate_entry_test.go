package conformance

import (
	"github.com/gibbonmi/bench/internal/capability"
	"os"
	"strings"
	"testing"
)

func TestRootConformance(t *testing.T) {
	root := os.Getenv("BENCH_CONFORMANCE_ROOT")
	if root == "" {
		capability.Environment(t, "BENCH_CONFORMANCE_ROOT not set")
	}
	h := NewHarness(t)
	for _, diag := range RunConformance(root, h.KitRoot) {
		t.Errorf("gate: %s", diag)
	}
}

func TestGateEntryRunsGoConformanceAndBehaviorContracts(t *testing.T) {
	h := NewHarness(t)
	gate := h.ReadRootFile(".bench", "gate.sh")

	needle := `exec "$kit/bin/bench.sh" gate-phases "$root"`
	if !strings.Contains(gate, needle) {
		t.Fatalf(".bench/gate.sh missing %q", needle)
	}

	for _, retired := range []string{
		`BENCH_CONFORMANCE_ROOT="$root" go test -count=1 ./internal/conformance -run '^TestRootConformance$'`,
		`BENCH_CONTRACT_ROOT="$root" go test -count=1 ./internal/contract/...`,
		`bin/bench.sh" canary "$root"`,
		"gate-docs-contracts.sh",
		"gate-line-contracts.sh",
		"gate-package-contracts.sh",
		"gate-runtime-shift-contracts.sh",
		"gate-contract-runner.sh",
		"gate-go-contracts.sh",
		"gate-link-contracts.sh",
		"gate-runtime-contracts.sh",
		"gate-runtime-git-contracts.sh",
		"gate-doctor-contracts.sh",
		"gate-axi-contracts.sh",
		"gate-axi-guards-contracts.sh",
		"gate-axi-wave2-contracts.sh",
	} {
		if strings.Contains(gate, retired) {
			t.Fatalf(".bench/gate.sh still references retired conformance fragment %s", retired)
		}
	}
}
