package conformance

import (
	"os"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/conformance/registry"
)

func TestRootConformance(t *testing.T) {
	root := os.Getenv("BENCH_CONFORMANCE_ROOT")
	if root == "" {
		capability.Environment(t, "BENCH_CONFORMANCE_ROOT not set")
	}
	h := NewHarness(t)
	for _, diag := range RunConformance(root, h.KitRoot, entryTier(os.Getenv(registry.ConformanceTierEnv))) {
		t.Errorf("gate: %s", diag)
	}
}

// entryTier reads the tier this entry point grades. The env var and the token that
// selects the ship surface both come from the registry, so nothing here restates that
// contract. Any other value is the dev tier: a typo or a stray export must never widen
// what the gate runs.
func entryTier(value string) registry.Tier {
	if value == string(registry.Ship) {
		return registry.Ship
	}
	return registry.Dev
}

func TestEntryTierDefaultsToDev(t *testing.T) {
	for _, test := range []struct {
		value string
		want  registry.Tier
	}{
		{"", registry.Dev},
		{"Ship", registry.Dev},
		{"dev", registry.Dev},
		{"anything", registry.Dev},
		{string(registry.Ship), registry.Ship},
	} {
		if got := entryTier(test.value); got != test.want {
			t.Errorf("entryTier(%q) = %q, want %q", test.value, got, test.want)
		}
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
