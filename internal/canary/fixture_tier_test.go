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
	union := append(fixtureDirs(dev), fixtureDirs(ship)...)
	slices.Sort(union)
	if !slices.Equal(union, all) {
		t.Fatalf("tier selections do not reconstruct the fixture list")
	}

	for _, fx := range dev {
		tier, _, err := fixtureCheck(fx.dir)
		if err != nil {
			t.Fatalf("%s: %v", fx.dir, err)
		}
		if tier != registry.Dev {
			t.Errorf("dev sweep selected %s, whose tier is %q", filepath.Base(fx.dir), tier)
		}
	}

	var shipNames []string
	for _, fx := range ship {
		shipNames = append(shipNames, filepath.Base(fx.dir))
	}
	slices.Sort(shipNames)
	want := slices.Clone(shipFixtures)
	slices.Sort(want)
	if !slices.Equal(shipNames, want) {
		t.Errorf("ship-tier fixtures: got %v, want %v", shipNames, want)
	}
}

// TestFixtureTierResolution covers the synthetic cases the real tree cannot show at
// once: a fixture with no CHECK file, one naming a ship-tier check, one naming a check
// the registry does not carry, and the two blank files that name no check at all —
// each of which has to reach a diagnostic describing the condition it actually is.
func TestFixtureTierResolution(t *testing.T) {
	cases := []struct {
		name     string
		absent   bool
		check    string
		wantTier registry.Tier
		wantErr  string
	}{
		{name: "plain", absent: true, wantTier: registry.Dev},
		{name: "shipped", check: shipCheckName(t) + "\n", wantTier: registry.Ship},
		{name: "bogus", check: "no-such-check\n", wantErr: "no-such-check"},
		{name: "empty", check: "", wantErr: "empty " + checkFileName + " file"},
		{name: "blank", check: " \t\n", wantErr: "empty " + checkFileName + " file"},
	}

	root := t.TempDir()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fx := canaryFixture(root, mappedFamily(t), tc.name)
			mkdir(t, filepath.Join(fx, "files"))
			write(t, filepath.Join(fx, "EXPECT"), "target-"+tc.name+"\n")
			if !tc.absent {
				write(t, filepath.Join(fx, checkFileName), tc.check)
			}

			tier, _, err := fixtureCheck(fx)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Errorf("got (%q, %v), want a diagnostic naming %q", tier, err, tc.wantErr)
				}
				return
			}
			if err != nil || tier != tc.wantTier {
				t.Errorf("got (%q, %v), want %q", tier, err, tc.wantTier)
			}
		})
	}
}

// TestSweepTierRunsOnlyItsOwnTier proves the selection reaches the sweep rather than
// stopping at the helper: each tier's runner sees its own fixtures and no others.
func TestSweepTierRunsOnlyItsOwnTier(t *testing.T) {
	root := t.TempDir()
	plain := canaryFixture(root, mappedFamily(t), "plain")
	mkdir(t, filepath.Join(plain, "files"))
	write(t, filepath.Join(plain, "EXPECT"), "target-plain\n")
	shipped := canaryFixture(root, mappedFamily(t), "shipped")
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

// TestSweepTierPinsInnerTier proves the inner gate grades the tier whose fixtures are
// being swept, and that an ambient export of the tier variable cannot select it — a
// ship fixture graded by a dev inner gate reports "did not bite" forever.
func TestSweepTierPinsInnerTier(t *testing.T) {
	t.Setenv(registry.ConformanceTierEnv, "ambient")
	root := t.TempDir()
	plain := canaryFixture(root, mappedFamily(t), "plain")
	mkdir(t, filepath.Join(plain, "files"))
	write(t, filepath.Join(plain, "EXPECT"), "target-plain\n")
	shipped := canaryFixture(root, mappedFamily(t), "shipped")
	mkdir(t, filepath.Join(shipped, "files"))
	write(t, filepath.Join(shipped, "EXPECT"), "target-shipped\n")
	write(t, filepath.Join(shipped, "CHECK"), shipCheckName(t)+"\n")

	tiersSeen := func(t *testing.T, tier registry.Tier) []string {
		t.Helper()
		var mu sync.Mutex
		var seen []string
		runner := func(call RunCall) RunResult {
			mu.Lock()
			for _, kv := range call.Env {
				if value, ok := strings.CutPrefix(kv, registry.ConformanceTierEnv+"="); ok {
					seen = append(seen, value)
				}
			}
			mu.Unlock()
			if call.FixtureDir == "" {
				return RunResult{ExitCode: 1, Output: "baseline\n"}
			}
			return RunResult{ExitCode: 1, Output: "target-" + filepath.Base(call.FixtureDir) + "\n"}
		}
		if err := SweepTier(root, tier, runner); err != nil {
			t.Fatalf("%s sweep: %v", tier, err)
		}
		if len(seen) == 0 {
			t.Fatalf("%s sweep handed the runner no %s", tier, registry.ConformanceTierEnv)
		}
		return seen
	}

	for _, tier := range []registry.Tier{registry.Dev, registry.Ship} {
		for _, got := range tiersSeen(t, tier) {
			if got != string(tier) {
				t.Errorf("%s sweep ran an inner gate at tier %q", tier, got)
			}
		}
	}
}

// fixtureDirs is the selection's fixture paths, for comparison against the
// directory listing the selection was drawn from.
func fixtureDirs(selection []selected) []string {
	out := make([]string, 0, len(selection))
	for _, fx := range selection {
		out = append(out, fx.dir)
	}
	return out
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
