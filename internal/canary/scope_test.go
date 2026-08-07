package canary

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/conformance/registry"
)

func TestFixturesExposeResolvedChecks(t *testing.T) {
	const (
		fixtureName = "default-branch-refabricated"
		wantCheck   = "default-branch-single-source"
	)
	fixtures, err := Fixtures(filepath.Join(kitRoot(t), "tests", "canary"))
	if err != nil {
		t.Fatal(err)
	}
	got, found := fixtures[fixtureName]
	if !found {
		t.Fatalf("fixtures did not include %q", fixtureName)
	}
	if got.Check != wantCheck {
		t.Errorf("fixture %q check = %q, want %q", fixtureName, got.Check, wantCheck)
	}

	family := mappedFamily(t)
	wantFamilyCheck := boundCheck(t, family)
	root := t.TempDir()
	fixture(t, canaryFixture(root, family, "family-fallback"), "")
	fixtures, err = Fixtures(filepath.Join(root, "tests", "canary"))
	if err != nil {
		t.Fatal(err)
	}
	if got := fixtures["family-fallback"].Check; got != wantFamilyCheck {
		t.Errorf("fixture without CHECK resolved check = %q, want family check %q", got, wantFamilyCheck)
	}
}

// TestSweepScopesFixtureRunsToTheirCheck grades the whole scope resolution at once:
// a conformance family's fixture takes the family binding, a CHECK file overrides
// it, contract and legacy flat fixtures take no scope at all, and an operator's
// ambient export of any conformance selection control reaches none of them. Each
// fixture is also enumerated so that an implementation merging same-check fixtures
// into one inner run grades fewer mutated trees than were selected.
func TestSweepScopesFixtureRunsToTheirCheck(t *testing.T) {
	const ambient = "ambient-scope"
	t.Setenv(registry.ConformanceCheckEnv, ambient)
	t.Setenv(registry.ConformanceChecksEnv, ambient)
	t.Setenv(registry.ConformanceInheritedEnv, ambient)
	family := mappedFamily(t)
	bound := boundCheck(t, family)
	override := devCheckNameOtherThan(t, bound)

	root := t.TempDir()
	fixture(t, canaryFixture(root, family, "family-bound"), "")
	fixture(t, canaryFixture(root, family, "check-bound"), override)
	contractFixture(t, root, "axi", "contract-fx")
	fixture(t, filepath.Join(root, "tests", "canary", "flat-fx"), "")

	calls := sweepCalls(t, root, registry.Dev)

	want := map[string]string{
		"family-bound": bound,
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

	// Baselines are the half the loop above skips, having no FixtureDir to key on, and
	// they are where an inherited scope does the most damage: an empty tree graded under
	// the operator's check produces output belonging to no group, and every fixture in
	// the group is then graded vacuous or not against the wrong run. The unscoped group
	// carries no scope variable at all, so an ambient value surviving the strip shows up
	// here as a scope where none belongs.
	wantGroups := sortedGroups("", bound, override, contractGroupPrefix+"axi")
	if got := baselineGroups(t, calls); !slices.Equal(got, wantGroups) {
		t.Errorf("baseline groups = %v, want the four groups' own keys and no ambient value", got)
	}
	for _, call := range calls {
		if slices.Contains(scopeValues(call.Env), ambient) {
			t.Errorf("call %q carried the ambient scope %s, want it stripped", call.FixtureDir, ambient)
		}
		for _, key := range []string{registry.ConformanceChecksEnv, registry.ConformanceInheritedEnv} {
			if values := envValues(call.Env, key); len(values) != 0 {
				t.Errorf("call %q carried %s=%v, want the outer selection absent", call.FixtureDir, key, values)
			}
		}
	}
}

// sortedGroups is the sorted form baselineGroups returns, so a wanted set is written in
// the order it reads best and compared in the order the helper produces.
func sortedGroups(groups ...string) []string {
	out := append([]string(nil), groups...)
	slices.Sort(out)
	return out
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
// baseline per distinct group, whatever the fixture count, and one full baseline
// shared by every fixture needing the full inner gate.
func TestSweepRunsOneBaselinePerScopeGroup(t *testing.T) {
	family := mappedFamily(t)
	bound := boundCheck(t, family)
	override := devCheckNameOtherThan(t, bound)

	root := t.TempDir()
	fixture(t, canaryFixture(root, family, "family-a"), "")
	fixture(t, canaryFixture(root, family, "family-b"), "")
	fixture(t, canaryFixture(root, family, "check-bound"), override)
	contractFixture(t, root, "axi", "contract-a")
	contractFixture(t, root, "axi", "contract-b")
	fixture(t, filepath.Join(root, "tests", "canary", "flat-fx"), "")

	got := baselineGroups(t, sweepCalls(t, root, registry.Dev))
	want := sortedGroups("", bound, override, contractGroupPrefix+"axi")
	if !slices.Equal(got, want) {
		t.Fatalf("baseline groups = %v, want %v", got, want)
	}
}

// TestAllUnscopedSweepRunsOneBaseline is the regression guard on the grouping key:
// grouping by anything finer than the group a fixture resolves to — the fixture's
// family, say — splits legacy flat fixtures from phase-named ones and multiplies the
// cost of the one baseline they all share.
func TestAllUnscopedSweepRunsOneBaseline(t *testing.T) {
	root := t.TempDir()
	fixture(t, canaryFixture(root, PhaseGofmt, "gofmt-fx"), "")
	fixture(t, canaryFixture(root, PhaseVet, "vet-fx"), "")
	fixture(t, filepath.Join(root, "tests", "canary", "flat-fx"), "")

	if got := baselineGroups(t, sweepCalls(t, root, registry.Dev)); !slices.Equal(got, []string{""}) {
		t.Fatalf("baseline groups = %v, want one unscoped baseline", got)
	}
}

// TestSweepGradesVacuityAgainstItsOwnGroupBaseline proves the comparison is
// per-group in both directions: an EXPECT its own group's baseline already emits is
// vacuous, and the same text emitted only by another group's baseline is not.
func TestSweepGradesVacuityAgainstItsOwnGroupBaseline(t *testing.T) {
	root := t.TempDir()
	ownScope := boundCheck(t, mappedFamily(t))
	own := canaryFixture(root, mappedFamily(t), "own-group")
	fixture(t, own, "")
	write(t, filepath.Join(own, "EXPECT"), "shared noise\n")
	other := canaryFixture(root, secondMappedFamily(t), "other-group")
	fixture(t, other, "")
	write(t, filepath.Join(other, "EXPECT"), "shared noise\n")

	err := Sweep(root, func(call RunCall) RunResult {
		if call.FixtureDir == "" {
			if slices.Equal(scopeValues(call.Env), []string{ownScope}) {
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

// TestSweepRejectsGroupBaselineThatPrintedNothing grades the group the vacuity check
// cannot defend on its own: a baseline that dies before printing contains no EXPECT,
// so every fixture in its group would clear the vacuity test unguarded while the other
// groups stayed graded. The sweep stops instead, and names the group so the failure is
// attributable.
func TestSweepRejectsGroupBaselineThatPrintedNothing(t *testing.T) {
	root := t.TempDir()
	starved := boundCheck(t, secondMappedFamily(t))
	fixture(t, canaryFixture(root, mappedFamily(t), "graded"), "")
	fixture(t, canaryFixture(root, secondMappedFamily(t), "starved"), "")

	var mu sync.Mutex
	var ran []string
	err := Sweep(root, func(call RunCall) RunResult {
		if call.FixtureDir != "" {
			mu.Lock()
			ran = append(ran, filepath.Base(call.FixtureDir))
			mu.Unlock()
			return RunResult{ExitCode: 1, Output: "target-" + filepath.Base(call.FixtureDir) + "\n"}
		}
		if slices.Equal(scopeValues(call.Env), []string{starved}) {
			return RunResult{ExitCode: 1, Output: ""}
		}
		return RunResult{ExitCode: 1, Output: "baseline noise\n"}
	})

	if err == nil {
		t.Fatal("Sweep err = nil, want the empty baseline reported")
	}
	if want := fmt.Sprintf("scope group %q", starved); !strings.Contains(err.Error(), want) {
		t.Fatalf("Sweep err = %v, want a diagnostic naming %s", err, want)
	}
	if len(ran) != 0 {
		t.Fatalf("fixtures %v ran against an ungradeable baseline, want none", ran)
	}
}

// TestShipSweepScopesItsFixtures covers the release path: ship's fixtures all name
// the one ship-tier check, so the sweep is their scoped runs plus a single shared
// scoped baseline rather than a full inner gate apiece.
func TestShipSweepScopesItsFixtures(t *testing.T) {
	ship := shipCheckName(t)
	root := t.TempDir()
	fixture(t, canaryFixture(root, mappedFamily(t), "ship-a"), ship)
	fixture(t, canaryFixture(root, mappedFamily(t), "ship-b"), ship)

	calls := sweepCalls(t, root, registry.Ship)
	if len(calls) != 3 {
		t.Fatalf("ship sweep ran %d inner gates, want 2 fixtures + 1 shared baseline", len(calls))
	}
	if got := baselineGroups(t, calls); !slices.Equal(got, []string{ship}) {
		t.Fatalf("ship baseline groups = %v, want one baseline scoped to %s", got, ship)
	}
	for _, call := range calls {
		if got := scopeValues(call.Env); !slices.Equal(got, []string{ship}) {
			t.Errorf("ship call %q carried scopes %v, want exactly [%s]", call.FixtureDir, got, ship)
		}
	}
}

func TestGateOwnedFamilySelectionNarrowsDevButNotShip(t *testing.T) {
	selected := mappedFamily(t)
	other := secondMappedFamily(t)
	root := t.TempDir()
	fixture(t, canaryFixture(root, selected, "selected"), "")
	fixture(t, canaryFixture(root, other, "other"), "")

	t.Setenv(FamilySelectionEnv, selected)
	if got := fixtureCallNames(sweepCalls(t, root, registry.Dev)); !slices.Equal(got, []string{"other", "selected"}) {
		t.Fatalf("ambient family selection ran %v, want both families", got)
	}
	t.Setenv(FamilySelectionOwnerEnv, "gate")
	if got := fixtureCallNames(sweepCalls(t, root, registry.Dev)); !slices.Equal(got, []string{"other", "selected"}) {
		t.Fatalf("ambient gate-owner forgery ran %v, want both families", got)
	}
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.WriteString(selected); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	t.Setenv(FamilySelectionAuthorityEnv, fmt.Sprint(reader.Fd()))
	if got := fixtureCallNames(sweepCalls(t, root, registry.Dev)); !slices.Equal(got, []string{"selected"}) {
		t.Fatalf("gate-owned family selection ran %v, want selected family", got)
	}
	write(t, filepath.Join(canaryFixture(root, selected, "selected"), checkFileName), "release-evidence-probe\n")
	write(t, filepath.Join(canaryFixture(root, other, "other"), checkFileName), "release-evidence-probe\n")
	if got := fixtureCallNames(sweepCalls(t, root, registry.Ship)); !slices.Equal(got, []string{"other", "selected"}) {
		t.Fatalf("ship family selection ran %v, want both families", got)
	}
}

func TestGateFamilySelectionAuthorityReadDoesNotBlockOnHeldOpenDescriptor(t *testing.T) {
	selected := mappedFamily(t)
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		reader.Close()
		writer.Close()
	})
	if _, err := writer.WriteString(selected); err != nil {
		t.Fatal(err)
	}
	t.Setenv(FamilySelectionEnv, selected)
	t.Setenv(FamilySelectionOwnerEnv, "gate")
	t.Setenv(FamilySelectionAuthorityEnv, fmt.Sprint(reader.Fd()))

	done := make(chan error, 1)
	go func() {
		_, scoped, err := gateSelectedFamilies(registry.Dev)
		if err == nil || scoped {
			done <- fmt.Errorf("held-open authority resolved scoped=%t err=%v, want an invalid non-scoped selection", scoped, err)
			return
		}
		done <- nil
	}()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(bounds.TestDeadline(0)):
		t.Fatal("held-open authority descriptor blocked canary selection")
	}
}

func fixtureCallNames(calls []RunCall) []string {
	var names []string
	for _, call := range calls {
		if call.FixtureDir != "" {
			names = append(names, filepath.Base(call.FixtureDir))
		}
	}
	slices.Sort(names)
	return names
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

// stubToolchain answers the two calls a contract package's binary is the subject of, and
// reports whether it handled the call. Every fake runner routes both through it: the sweep
// checks the binary exists before invoking it, so a runner that only returned exit zero for
// a compile would red every contract group in every test, and a list answered as if it were
// a graded run would fail every fixture's owner against output that names no test.
//
// A compile answers the way a successful `go test -c` does — exit zero and a binary on disk.
// A list answers from the roster the fixture helpers write beside the package source, which
// is what keeps a fixture's declared owner and the membership it is graded against on one
// source in the synthetic tree, the way the kit's own fixtures and binaries are.
func stubToolchain(call RunCall) (RunResult, bool) {
	switch call.Kind {
	case RunCompile:
		if err := os.WriteFile(call.Binary, nil, 0o755); err != nil {
			return RunResult{ExitCode: 1, Output: err.Error()}, true
		}
		return RunResult{}, true
	case RunList:
		roster, err := os.ReadFile(filepath.Join(call.Cwd, testRosterName))
		if err != nil {
			return RunResult{ExitCode: 1, Output: err.Error()}, true
		}
		return RunResult{Output: string(roster)}, true
	}
	return RunResult{}, false
}

// isBaseline reports whether a call is one of the empty-tree vacuity baselines. The compile
// and list calls a contract group needs share the baselines' empty FixtureDir and grade no
// tree, so membership is stated by kind rather than read from the absent fixture.
func isBaseline(call RunCall) bool {
	return call.FixtureDir == "" && (call.Kind == RunGate || call.Kind == RunBite)
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
		if result, done := stubToolchain(call); done {
			return result
		}
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

// baselineGroups lists the group each empty-tree baseline ran for, sorted, so a
// duplicate baseline for one group shows up as a repeated entry.
func baselineGroups(t *testing.T, calls []RunCall) []string {
	t.Helper()
	var out []string
	for _, call := range calls {
		if !isBaseline(call) {
			continue
		}
		out = append(out, callGroup(t, call))
	}
	slices.Sort(out)
	return out
}

func scopeValues(env []string) []string {
	return envValues(env, registry.ConformanceCheckEnv)
}

// boundCheck is the check the registry binds family to, which is what the sweep resolves
// a family's fixtures to. A family name reads as its own scope only for the families
// named after the check they grade, so the two are never used interchangeably here.
func boundCheck(t *testing.T, family string) string {
	t.Helper()
	check, bound := registry.FamilyCheck(family)
	if !bound {
		t.Fatalf("registry binds no check to family %s", family)
	}
	return check
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
