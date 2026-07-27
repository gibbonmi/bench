package canary

import (
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
)

// TestSweepScopesFixtureRunsToTheirCheck grades the whole scope resolution at once:
// a conformance family's fixture takes the family binding, a CHECK file overrides
// it, contract and legacy flat fixtures take no scope at all, and an operator's
// ambient export of the scope variable reaches none of them. Each fixture is also
// enumerated so that an implementation merging same-check fixtures into one inner
// run grades fewer mutated trees than were selected.
func TestSweepScopesFixtureRunsToTheirCheck(t *testing.T) {
	t.Setenv(registry.ConformanceCheckEnv, "ambient-scope")
	override := devCheckNameOtherThan(t, mappedFamily)

	root := t.TempDir()
	fixture(t, canaryFixture(root, mappedFamily, "family-bound"), "")
	fixture(t, canaryFixture(root, mappedFamily, "check-bound"), override)
	fixture(t, canaryFixture(root, "behavior-owned", "contract-fx"), "")
	fixture(t, filepath.Join(root, "tests", "canary", "flat-fx"), "")

	calls := sweepCalls(t, root, registry.Dev)

	want := map[string]string{
		"family-bound": mappedFamily,
		"check-bound":  override,
		"contract-fx":  "",
		"flat-fx":      "",
	}
	for name, wantScope := range want {
		var got []RunCall
		for _, call := range calls {
			if call.FixtureDir != "" && filepath.Base(call.FixtureDir) == name {
				got = append(got, call)
			}
		}
		if len(got) != 1 {
			t.Errorf("fixture %s ran %d inner gates, want exactly 1", name, len(got))
			continue
		}
		scopes := scopeValues(got[0].Env)
		if wantScope == "" {
			if len(scopes) != 0 {
				t.Errorf("unscoped fixture %s carried scopes %v, want none", name, scopes)
			}
			continue
		}
		if !slices.Equal(scopes, []string{wantScope}) {
			t.Errorf("fixture %s carried scopes %v, want exactly [%s]", name, scopes, wantScope)
		}
	}
}

// TestSweepRunsUnmappedConformanceFamilyUnscoped pins the fallback every adopting
// repo depends on. The family table is knowledge about the kit's own fixture tree,
// while this sweep grades any repo that adopts Bench — including the seed canary
// `bench init` scaffolds, whose family no kit-owned table will ever carry. An unbound
// family therefore runs the full inner gate rather than erroring; the kit's own
// unbound family is a red the conformance layer raises.
func TestSweepRunsUnmappedConformanceFamilyUnscoped(t *testing.T) {
	root := t.TempDir()
	fixture(t, canaryFixture(root, "unmapped-family", "orphan"), "")

	calls := sweepCalls(t, root, registry.Dev)

	if len(calls) != 2 {
		t.Fatalf("sweep ran %d inner gates, want baseline + fixture", len(calls))
	}
	for _, call := range calls {
		if scopes := scopeValues(call.Env); len(scopes) != 0 {
			t.Errorf("call %q carried scopes %v, want an unscoped full run", call.FixtureDir, scopes)
		}
	}
}

// TestSweepRunsOneBaselinePerScopeGroup pins the shared-baseline cost model: one
// baseline per distinct scope, whatever the fixture count, and one full baseline
// shared by every unscoped fixture.
func TestSweepRunsOneBaselinePerScopeGroup(t *testing.T) {
	override := devCheckNameOtherThan(t, mappedFamily)

	root := t.TempDir()
	fixture(t, canaryFixture(root, mappedFamily, "family-a"), "")
	fixture(t, canaryFixture(root, mappedFamily, "family-b"), "")
	fixture(t, canaryFixture(root, mappedFamily, "check-bound"), override)
	fixture(t, canaryFixture(root, "behavior-owned", "contract-fx"), "")
	fixture(t, filepath.Join(root, "tests", "canary", "flat-fx"), "")

	got := baselineScopes(t, sweepCalls(t, root, registry.Dev))
	want := []string{"", mappedFamily, override}
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Fatalf("baseline scopes = %v, want %v", got, want)
	}
}

// TestAllUnscopedSweepRunsOneBaseline is the regression guard on the grouping key:
// grouping by anything finer than the resolved check name — the fixture's phase, say
// — splits contract fixtures from legacy flat ones and doubles today's cost.
func TestAllUnscopedSweepRunsOneBaseline(t *testing.T) {
	root := t.TempDir()
	fixture(t, canaryFixture(root, "behavior-owned", "contract-a"), "")
	fixture(t, canaryFixture(root, "behavior-owned", "contract-b"), "")
	fixture(t, filepath.Join(root, "tests", "canary", "flat-fx"), "")

	if got := baselineScopes(t, sweepCalls(t, root, registry.Dev)); !slices.Equal(got, []string{""}) {
		t.Fatalf("baseline scopes = %v, want one unscoped baseline", got)
	}
}

// TestSweepGradesVacuityAgainstItsOwnGroupBaseline proves the comparison is
// per-group in both directions: an EXPECT its own group's baseline already emits is
// vacuous, and the same text emitted only by another group's baseline is not.
func TestSweepGradesVacuityAgainstItsOwnGroupBaseline(t *testing.T) {
	root := t.TempDir()
	own := canaryFixture(root, mappedFamily, "own-group")
	fixture(t, own, "")
	write(t, filepath.Join(own, "EXPECT"), "shared noise\n")
	other := canaryFixture(root, secondMappedFamily, "other-group")
	fixture(t, other, "")
	write(t, filepath.Join(other, "EXPECT"), "shared noise\n")

	err := Sweep(root, func(call RunCall) RunResult {
		if call.FixtureDir == "" {
			if slices.Equal(scopeValues(call.Env), []string{mappedFamily}) {
				return RunResult{ExitCode: 1, Output: "shared noise\n"}
			}
			return RunResult{ExitCode: 1, Output: "unrelated noise\n"}
		}
		return RunResult{ExitCode: 1, Output: "shared noise\n"}
	})
	if err == nil {
		t.Fatal("Sweep err = nil, want the own-group fixture reported vacuous")
	}
	if !strings.Contains(err.Error(), "canary 'own-group' EXPECT is vacuous") {
		t.Fatalf("Sweep err = %v, want own-group graded vacuous against its own baseline", err)
	}
	if strings.Contains(err.Error(), "other-group") {
		t.Fatalf("Sweep err = %v, want other-group graded against its own baseline alone", err)
	}
}

// TestShipSweepScopesItsFixtures covers the release path: ship's fixtures all name
// the one ship-tier check, so the sweep is their scoped runs plus a single shared
// scoped baseline rather than a full inner gate apiece.
func TestShipSweepScopesItsFixtures(t *testing.T) {
	ship := shipCheckName(t)
	root := t.TempDir()
	fixture(t, canaryFixture(root, mappedFamily, "ship-a"), ship)
	fixture(t, canaryFixture(root, mappedFamily, "ship-b"), ship)

	calls := sweepCalls(t, root, registry.Ship)
	if len(calls) != 3 {
		t.Fatalf("ship sweep ran %d inner gates, want 2 fixtures + 1 shared baseline", len(calls))
	}
	if got := baselineScopes(t, calls); !slices.Equal(got, []string{ship}) {
		t.Fatalf("ship baseline scopes = %v, want one baseline scoped to %s", got, ship)
	}
	for _, call := range calls {
		if got := scopeValues(call.Env); !slices.Equal(got, []string{ship}) {
			t.Errorf("ship call %q carried scopes %v, want exactly [%s]", call.FixtureDir, got, ship)
		}
	}
}

// fixture writes a minimal biting fixture: a files/ tree, an EXPECT the sweep
// helpers below emit back, and a CHECK file only when check names one.
func fixture(t *testing.T, dir, check string) {
	t.Helper()
	mkdir(t, filepath.Join(dir, "files"))
	write(t, filepath.Join(dir, "EXPECT"), "target-"+filepath.Base(dir)+"\n")
	if check != "" {
		write(t, filepath.Join(dir, checkFileName), check+"\n")
	}
}

// sweepCalls runs a green sweep of tier and returns every RunCall it made.
func sweepCalls(t *testing.T, root string, tier registry.Tier) []RunCall {
	t.Helper()
	var mu sync.Mutex
	var calls []RunCall
	runner := func(call RunCall) RunResult {
		mu.Lock()
		calls = append(calls, call)
		mu.Unlock()
		if call.FixtureDir == "" {
			return RunResult{ExitCode: 1, Output: "baseline noise\n"}
		}
		return RunResult{ExitCode: 1, Output: "target-" + filepath.Base(call.FixtureDir) + "\n"}
	}
	if err := SweepTier(root, tier, runner); err != nil {
		t.Fatalf("%s sweep: %v", tier, err)
	}
	return calls
}

// baselineScopes lists the scope each empty-tree baseline ran under, sorted, so a
// duplicate baseline for one group shows up as a repeated entry.
func baselineScopes(t *testing.T, calls []RunCall) []string {
	t.Helper()
	var out []string
	for _, call := range calls {
		if call.FixtureDir != "" {
			continue
		}
		scopes := scopeValues(call.Env)
		switch len(scopes) {
		case 0:
			out = append(out, "")
		case 1:
			out = append(out, scopes[0])
		default:
			t.Fatalf("baseline call carried scopes %v, want at most one", scopes)
		}
	}
	slices.Sort(out)
	return out
}

func scopeValues(env []string) []string {
	var out []string
	for _, kv := range env {
		if value, ok := strings.CutPrefix(kv, registry.ConformanceCheckEnv+"="); ok {
			out = append(out, value)
		}
	}
	return out
}

// devCheckNameOtherThan is a registered dev-tier check that is not exclude, read
// rather than written down so a rename retires with the registry.
func devCheckNameOtherThan(t *testing.T, exclude string) string {
	t.Helper()
	for _, check := range registry.Checks {
		if check.Tier == registry.Dev && check.Name != exclude {
			return check.Name
		}
	}
	t.Fatalf("registry carries no dev-tier check other than %s", exclude)
	return ""
}
