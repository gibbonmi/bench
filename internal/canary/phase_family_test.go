package canary

import (
	"path/filepath"
	"sort"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
)

// TestSweepRoutesPhaseNamedFamilies grades the routing a phase-named family buys: the
// fixture's inner gate is told the one phase that owns its failure, not the conformance
// phase every other family resolves to. The phase name itself is asserted, since a
// fixture routed to the wrong phase still reds — for the wrong reason, and forever.
func TestSweepRoutesPhaseNamedFamilies(t *testing.T) {
	families := []string{PhaseGofmt, PhaseVet, PhaseTest, PhaseRace, PhaseConformanceSuite}
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
// disturb: a phase-named family resolves to no conformance check and to no contract
// package, so its fixtures share the single full baseline every legacy flat fixture
// already runs against.
func TestPhaseNamedFamiliesJoinTheUnscopedBaseline(t *testing.T) {
	root := t.TempDir()
	fixture(t, canaryFixture(root, "gofmt", "gofmt-fx"), "")
	fixture(t, filepath.Join(root, "tests", "canary", "flat-fx"), "")

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
		t.Errorf("sweep ran %d baselines, want the one unscoped baseline shared with legacy flat fixtures", baselines)
	}
}

// TestBaselinesRunUnpinnedByPhase pins the shape every vacuity baseline keeps: the whole
// inner gate, never one phase of it. A baseline narrowed to a phase prints a fraction of
// what the empty tree can produce, so an EXPECT the full run already emits stops being
// caught as vacuous and every fixture in the group clears the check unguarded.
//
// A group is represented by whichever of its fixtures sorts first, and the only pin a
// baseline could inherit is that representative's family phase. Every fixture here belongs
// to a family that names one for exactly that reason: a legacy flat fixture sorting first
// would name no phase, and the assertion would then hold over a baseline never at risk.
func TestBaselinesRunUnpinnedByPhase(t *testing.T) {
	root := t.TempDir()
	fixture(t, canaryFixture(root, PhaseVet, "vet-fx"), "")
	fixture(t, canaryFixture(root, PhaseGofmt, "gofmt-fx"), "")
	conformance := mappedFamily(t)
	fixture(t, canaryFixture(root, conformance, "conformance-fx"), "")

	baselines := baselineCalls(sweepCalls(t, root, registry.Dev))

	if len(baselines) != 2 {
		t.Fatalf("sweep ran %d baselines, want the unscoped group's and family %s's", len(baselines), conformance)
	}
	for _, call := range baselines {
		if phases := phaseValues(call.Env); len(phases) != 0 {
			t.Errorf("baseline for %s carried phases %v, want the full inner gate", groupLabel(callGroup(t, call)), phases)
		}
	}
}

func phaseValues(env []string) []string {
	return envValues(env, PhaseEnv)
}
