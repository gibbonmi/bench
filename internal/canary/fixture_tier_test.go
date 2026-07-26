package canary

import (
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
)

// shipFixtures is the spec's own statement of which fixtures follow the
// release-evidence probe to the ship tier. It is written out rather than derived so
// that an implementation which ships nothing — or ships everything — is red here.
var shipFixtures = []string{
	"release-digest-binding-omitted",
	"release-package-evidence-omitted",
}

// TestFixtureTierMatchesCheckTier grades the real fixture tree: every fixture resolves
// to the tier of the registry check it names, the two tiers partition the harness so no
// fixture can quietly stop being swept, and the dev sweep runs dev-tier fixtures alone.
func TestFixtureTierMatchesCheckTier(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	all, err := fixtures(filepath.Join(root, "tests", "canary"))
	if err != nil {
		t.Fatal(err)
	}

	dev, err := selectTier(all, registry.Dev)
	if err != nil {
		t.Fatalf("dev selection: %v", err)
	}
	ship, err := selectTier(all, registry.Ship)
	if err != nil {
		t.Fatalf("ship selection: %v", err)
	}

	if len(dev)+len(ship) != len(all) {
		t.Fatalf("tiers do not partition the harness: %d dev + %d ship != %d fixtures", len(dev), len(ship), len(all))
	}
	union := append(append([]string(nil), dev...), ship...)
	slices.Sort(union)
	if !slices.Equal(union, all) {
		t.Fatalf("tier selections do not reconstruct the fixture list")
	}

	for _, fx := range dev {
		tier, err := fixtureTier(fx)
		if err != nil {
			t.Fatalf("%s: %v", fx, err)
		}
		if tier != registry.Dev {
			t.Errorf("dev sweep selected %s, whose tier is %q", filepath.Base(fx), tier)
		}
	}

	var shipNames []string
	for _, fx := range ship {
		shipNames = append(shipNames, filepath.Base(fx))
	}
	slices.Sort(shipNames)
	want := slices.Clone(shipFixtures)
	slices.Sort(want)
	if !slices.Equal(shipNames, want) {
		t.Errorf("ship-tier fixtures: got %v, want %v", shipNames, want)
	}
}

// TestFixtureTierResolution covers the three synthetic cases the real tree cannot show
// at once: a fixture with no CHECK file, one naming a ship-tier check, and one naming a
// check the registry does not carry.
func TestFixtureTierResolution(t *testing.T) {
	root := t.TempDir()

	plain := canaryFixture(root, "test-family", "plain")
	mkdir(t, filepath.Join(plain, "files"))
	write(t, filepath.Join(plain, "EXPECT"), "target-plain\n")

	shipped := canaryFixture(root, "test-family", "shipped")
	mkdir(t, filepath.Join(shipped, "files"))
	write(t, filepath.Join(shipped, "EXPECT"), "target-shipped\n")
	write(t, filepath.Join(shipped, "CHECK"), shipCheckName(t)+"\n")

	bogus := canaryFixture(root, "test-family", "bogus")
	mkdir(t, filepath.Join(bogus, "files"))
	write(t, filepath.Join(bogus, "EXPECT"), "target-bogus\n")
	write(t, filepath.Join(bogus, "CHECK"), "no-such-check\n")

	if tier, err := fixtureTier(plain); err != nil || tier != registry.Dev {
		t.Errorf("fixture without CHECK: got (%q, %v), want dev", tier, err)
	}
	if tier, err := fixtureTier(shipped); err != nil || tier != registry.Ship {
		t.Errorf("fixture naming a ship check: got (%q, %v), want ship", tier, err)
	}
	_, err := fixtureTier(bogus)
	if err == nil || !strings.Contains(err.Error(), "no-such-check") {
		t.Errorf("fixture naming an unknown check: got %v, want a diagnostic naming it", err)
	}
}

// TestSweepTierRunsOnlyItsOwnTier proves the selection reaches the sweep rather than
// stopping at the helper: each tier's runner sees its own fixtures and no others.
func TestSweepTierRunsOnlyItsOwnTier(t *testing.T) {
	root := t.TempDir()
	plain := canaryFixture(root, "test-family", "plain")
	mkdir(t, filepath.Join(plain, "files"))
	write(t, filepath.Join(plain, "EXPECT"), "target-plain\n")
	shipped := canaryFixture(root, "test-family", "shipped")
	mkdir(t, filepath.Join(shipped, "files"))
	write(t, filepath.Join(shipped, "EXPECT"), "target-shipped\n")
	write(t, filepath.Join(shipped, "CHECK"), shipCheckName(t)+"\n")

	swept := func(t *testing.T, tier registry.Tier) []string {
		t.Helper()
		var mu sync.Mutex
		var seen []string
		runner := func(call RunCall) RunResult {
			if call.FixtureDir == "" {
				return RunResult{ExitCode: 1, Output: "baseline\n"}
			}
			mu.Lock()
			seen = append(seen, filepath.Base(call.FixtureDir))
			mu.Unlock()
			return RunResult{ExitCode: 1, Output: "target-" + filepath.Base(call.FixtureDir) + "\n"}
		}
		if err := SweepTier(root, tier, runner); err != nil {
			t.Fatalf("%s sweep: %v", tier, err)
		}
		slices.Sort(seen)
		return seen
	}

	if got := swept(t, registry.Dev); !slices.Equal(got, []string{"plain"}) {
		t.Errorf("dev sweep ran %v, want [plain]", got)
	}
	if got := swept(t, registry.Ship); !slices.Equal(got, []string{"shipped"}) {
		t.Errorf("ship sweep ran %v, want [shipped]", got)
	}
}

// shipCheckName is the registry's ship-tier check, read rather than written down so the
// synthetic fixtures follow a retiering instead of pinning one name.
func shipCheckName(t *testing.T) string {
	t.Helper()
	for _, check := range registry.Checks {
		if check.Tier == registry.Ship {
			return check.Name
		}
	}
	t.Fatal("registry carries no ship-tier check")
	return ""
}
