package preprelease

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/gate"
)

// seededFailure is the test name both isolation rows look for in a run's output. It is
// the whole observation: which tier reported it is what the split is.
const seededFailure = "TestSeededTierFailure"

// TestShipTierFailureIsolation is the story 1 and 7 acceptance. The seed is a failing
// test in a release-only package — the cheapest surface only the ship tier runs — so a
// single-tier implementation, which cannot produce two verdicts for one fault, fails
// this row by construction. Seeding the probe instead would activate only on a real
// tree, where the run costs the full 400 s+ path.
func TestShipTierFailureIsolation(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	root := newTierTree(t, "internal/preflight", "preflight")

	dev := gradeTier(t, root, registry.Dev)
	if strings.Contains(dev, seededFailure) {
		t.Fatalf("the dev tier ran the release-only package it stopped running:\n%s", dev)
	}
	ship := gradeTier(t, root, registry.Ship)
	if !strings.Contains(ship, seededFailure) {
		t.Fatalf("the ship tier did not run the release-only package it took over:\n%s", ship)
	}
}

// TestDevTierFailureRedsBothTiers is the row proving the split is a restaging and not a
// weakening: a fault in an ordinary package is reported by both tiers, so nothing the
// dev tier grades leaves the oracle when the ship tier takes the release checks.
func TestDevTierFailureRedsBothTiers(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	root := newTierTree(t, "internal/ordinary", "ordinary")

	for _, tier := range []registry.Tier{registry.Dev, registry.Ship} {
		if output := gradeTier(t, root, tier); !strings.Contains(output, seededFailure) {
			t.Fatalf("the %s tier did not report a fault in an ordinary package:\n%s", tier, output)
		}
	}
}

// newTierTree is a graded tree the tier split is observable in: a module carrying one
// package whose single test fails. Nothing else is present, so the only diagnostic that
// can name the seeded test is the `go test` whose package list the tier decides.
func newTierTree(t *testing.T, dir, pkg string) string {
	t.Helper()
	root := t.TempDir()
	contract.WriteFileAbs(t, filepath.Join(root, "go.mod"), "module benchtierfixture\n\ngo 1.25\n")
	contract.WriteFileAbs(t, filepath.Join(root, filepath.FromSlash(dir), "seed_test.go"),
		"package "+pkg+"\n\nimport \"testing\"\n\nfunc "+seededFailure+"(t *testing.T) { t.Fatal(\"seeded tier failure\") }\n")
	return root
}

// gradeTier runs the core test step over root at one tier and returns everything it
// reported. The tier-dependent package set lives behind `bench gate-go test`, so the dev
// call is byte-for-byte the gate's own `test` phase and the ship call is byte-for-byte
// prep-release's core-tests-ship step, which differs only by the tier variable. The argv
// comes from gate.GateGoArgv rather than a literal because that function is what both of
// those surfaces call: a row assembling its own argv could keep passing after the real
// invocation moved. The output is noisy by construction — the exit code answers for the
// whole enumeration — so the assertion is which tier names the seeded test.
func gradeTier(t *testing.T, root string, tier registry.Tier) string {
	t.Helper()
	f := contract.NewExecFixtureAt(t, root)
	env := map[string]string{}
	if tier == registry.Ship {
		env[registry.ConformanceTierEnv] = string(registry.Ship)
	}
	argv := gate.GateGoArgv(contract.SubjectRoot(t), "test", root)
	probe := contract.RunAt(t, f, root, env, argv[0], argv[1:]...)
	return probe.Stdout + probe.Stderr
}
