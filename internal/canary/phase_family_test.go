package canary

import (
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
)

// TestSweepRoutesPhaseNamedFamilies grades the routing a phase-named family buys: the
// fixture's inner gate is told the one phase that owns its failure, not the conformance
// phase every other family resolves to. The phase name itself is asserted, since a
// fixture routed to the wrong phase still reds — for the wrong reason, and forever.
func TestSweepRoutesPhaseNamedFamilies(t *testing.T) {
	families := []string{"build", "gofmt", "vet", "test", "race", "conformance-suite"}
	root := t.TempDir()
	for _, family := range families {
		fixture(t, canaryFixture(root, family, family+"-fx"), "")
	}

	calls := sweepCalls(t, root, registry.Dev)

	for _, family := range families {
		var got []string
		for _, call := range calls {
			if call.FixtureDir == "" || filepath.Base(filepath.Dir(call.FixtureDir)) != family {
				continue
			}
			got = append(got, phaseValues(call.Env)...)
		}
		sort.Strings(got)
		if len(got) != 1 || got[0] != family {
			t.Errorf("family %s routed to phases %v, want exactly [%s]", family, got, family)
		}
	}
}

// TestPhaseNamedFamiliesJoinTheUnscopedBaseline pins the cost model the routing must not
// disturb: a phase-named family resolves to no conformance check, so its fixtures share
// the single full baseline every unscoped fixture already runs against.
func TestPhaseNamedFamiliesJoinTheUnscopedBaseline(t *testing.T) {
	root := t.TempDir()
	fixture(t, canaryFixture(root, "gofmt", "gofmt-fx"), "")
	fixture(t, canaryFixture(root, "behavior-owned", "contract-fx"), "")

	calls := sweepCalls(t, root, registry.Dev)

	var baselines int
	for _, call := range calls {
		if call.FixtureDir != "" {
			continue
		}
		baselines++
		if scopes := scopeValues(call.Env); len(scopes) != 0 {
			t.Errorf("baseline carried scopes %v, want the unscoped group", scopes)
		}
	}
	if baselines != 1 {
		t.Errorf("sweep ran %d baselines, want the one unscoped baseline shared with contract fixtures", baselines)
	}
}

func phaseValues(env []string) []string {
	var out []string
	for _, kv := range env {
		if value, ok := strings.CutPrefix(kv, PhaseEnv+"="); ok {
			out = append(out, value)
		}
	}
	return out
}
