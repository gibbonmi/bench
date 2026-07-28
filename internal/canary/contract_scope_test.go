package canary

import (
	"maps"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
)

// TestSweepDiscoversContractFixturesUnderTheirPackage grades the walk's whole answer for
// a nested package: the fixture is the directory holding EXPECT, it is graded by exactly
// one run, and the package it carries is every segment between the family and it. A walk
// keeping only the last segment binds the fixture to a package that does not exist, which
// is why the multi-segment path is the case under test.
func TestSweepDiscoversContractFixturesUnderTheirPackage(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "surface/artifact", "nested-fx")
	contractFixture(t, root, "axi", "flat-pkg-fx")

	calls := sweepCalls(t, root, registry.Dev)

	want := map[string]string{
		"nested-fx":   "surface/artifact",
		"flat-pkg-fx": "axi",
	}
	for name, pkg := range want {
		got := fixtureCalls(calls, name)
		if len(got) != 1 {
			t.Errorf("fixture %s ran %d graded runs, want exactly 1", name, len(got))
			continue
		}
		if got[0].Kind != RunBite || got[0].Package != pkg {
			t.Errorf("fixture %s ran kind %v for package %q, want a bite of %s", name, got[0].Kind, got[0].Package, pkg)
		}
	}
}

// TestSweepRejectsDuplicateNamesAcrossContractPackages extends global base-name
// uniqueness across the nesting: two packages are two directories, so a walk checking
// uniqueness per directory reads them as distinct and lets one fixture's work dir and
// diagnostics collide with the other's.
func TestSweepRejectsDuplicateNamesAcrossContractPackages(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "axi", "duplicate")
	contractFixture(t, root, "surface/artifact", "duplicate")

	err := Sweep(root, (&recordingRunner{}).Run)
	want := `canary fixture name "duplicate" appears in multiple families; base names must be globally unique`
	if err == nil || err.Error() != want {
		t.Fatalf("Sweep err = %v, want %s", err, want)
	}
}

// TestSweepRejectsStructurallyUnscopableContractFixtures grades the two structural defects
// that each leave a behavior-owned fixture with no contract package to be graded by at
// all. Every diagnostic names the fixture, because the sweep's tree is the only place the
// defect can be repaired.
func TestSweepRejectsStructurallyUnscopableContractFixtures(t *testing.T) {
	cases := []struct {
		name  string
		build func(t *testing.T, root string)
		want  []string
	}{
		{
			name: "no package directory",
			build: func(t *testing.T, root string) {
				fixture(t, canaryFixture(root, "behavior-owned", "unbound-fx"), "")
			},
			want: []string{"unbound-fx", "behavior-owned"},
		},
		{
			name: "package absent from the kit tree",
			build: func(t *testing.T, root string) {
				fixture(t, filepath.Join(root, "tests", "canary", "behavior-owned", "ghost", "rotted-fx"), "")
			},
			want: []string{"rotted-fx", "ghost"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			tc.build(t, root)

			calls, err := countedSweep(t, root)
			if err == nil {
				t.Fatalf("Sweep err = nil, want the fixture refused")
			}
			for _, want := range tc.want {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("Sweep err = %v, want a diagnostic naming %q", err, want)
				}
			}
			if calls != 0 {
				t.Errorf("sweep ran %d inner gates before refusing, want none", calls)
			}
		})
	}
}

// TestSweepRejectsExpectAbovePackagedFixtures grades the walk's stopping rule. An EXPECT
// left at package depth makes that directory look like the fixture, and every real
// fixture below it disappears from the sweep — the harness silently stops grading them
// while staying green. A fixture is a leaf: files/ and its own marker files, nothing
// else, so a directory holding both an EXPECT and further fixture directories is the
// defect and is named.
func TestSweepRejectsExpectAbovePackagedFixtures(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "surface/artifact", "buried-a")
	contractFixture(t, root, "surface/artifact", "buried-b")
	pkgDir := filepath.Join(root, "tests", "canary", "behavior-owned", "surface", "artifact")
	write(t, filepath.Join(pkgDir, "EXPECT"), "stray expectation\n")

	calls, err := countedSweep(t, root)
	if err == nil {
		t.Fatal("Sweep err = nil, want the stray EXPECT refused")
	}
	for _, want := range []string{"artifact", "EXPECT"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("Sweep err = %v, want a diagnostic naming %q", err, want)
		}
	}
	for _, buried := range []string{"buried-a", "buried-b"} {
		if strings.Contains(err.Error(), buried) {
			t.Errorf("Sweep err = %v, want the containing directory named rather than %s", err, buried)
		}
	}
	if calls != 0 {
		t.Errorf("sweep ran %d inner gates before refusing, want none", calls)
	}
}

// TestContractFixtureBiteCarriesOnlyItsSubjectRoot pins the whole control set a compiled
// bite runs under: the tree it grades, and the width pin the sweep's worker budget is
// derived from. Every gate-era variable is asserted absent rather than ignored — there is
// no gate in this run to read one, and a binary carrying the inner-gate marker or a phase
// pin claims to be a nested gate to whatever reads them next.
func TestContractFixtureBiteCarriesOnlyItsSubjectRoot(t *testing.T) {
	t.Setenv(PhaseEnv, "ambient-phase")
	t.Setenv(registry.ConformanceCheckEnv, "ambient-check")

	root := t.TempDir()
	contractFixture(t, root, "surface/artifact", "scoped-fx")

	got := fixtureCalls(sweepCalls(t, root, registry.Dev), "scoped-fx")
	if len(got) != 1 {
		t.Fatalf("fixture ran %d graded runs, want exactly 1", len(got))
	}
	for _, key := range []string{PhaseEnv, registry.ConformanceCheckEnv, registry.ConformanceTierEnv, "BENCH_CANARY_INNER"} {
		if values := envValues(got[0].Env, key); len(values) != 0 {
			t.Errorf("bite carried %s=%v, want the variable absent", key, values)
		}
	}
	if values := envValues(got[0].Env, "GOMAXPROCS"); len(values) != 1 {
		t.Errorf("bite carried GOMAXPROCS=%v, want exactly the sweep's width pin", values)
	}
}

// TestSweepBaselinesContractGroupsPerPackage grades the vacuity groups two packages
// produce: one baseline each, run in the same shape as the group's own fixtures and with
// the group's own binary. A single shared baseline — the degenerate grouping — fails on
// the count, and a baseline that spawned a gate instead fails on the kind.
func TestSweepBaselinesContractGroupsPerPackage(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "axi", "axi-a")
	contractFixture(t, root, "axi", "axi-b")
	contractFixture(t, root, "surface/artifact", "artifact-fx")

	calls := sweepCalls(t, root, registry.Dev)
	compiled := compileOutputs(t, calls)

	baselines := map[string]int{}
	for _, call := range baselineCalls(calls) {
		if call.Kind != RunBite {
			t.Fatalf("contract baseline ran kind %v, want the same bite its fixtures run", call.Kind)
		}
		if want := compiled[call.Package]; call.Binary != want {
			t.Errorf("baseline for %s ran binary %q, want its group's compile output %q", call.Package, call.Binary, want)
		}
		baselines[call.Package]++
	}
	want := map[string]int{"axi": 1, "surface/artifact": 1}
	if !maps.Equal(baselines, want) {
		t.Fatalf("contract baselines = %v, want one per package", baselines)
	}
}

// TestContractFixtureGradesVacuityAgainstItsOwnPackageBaseline proves the group key the
// baseline is written under is the one the fixture reads back. A desync yields an absent
// baseline, and an EXPECT compared against nothing is never vacuous, so vacuity would be
// silently un-guarded for every behavior-owned fixture.
func TestContractFixtureGradesVacuityAgainstItsOwnPackageBaseline(t *testing.T) {
	root := t.TempDir()
	own := contractFixture(t, root, "axi", "own-package")
	write(t, filepath.Join(own, "EXPECT"), "shared noise\n")
	other := contractFixture(t, root, "surface/artifact", "other-package")
	write(t, filepath.Join(other, "EXPECT"), "shared noise\n")

	err := Sweep(root, func(call RunCall) RunResult {
		if result, done := stubCompile(call); done {
			return result
		}
		if call.FixtureDir == "" {
			if call.Package == "axi" {
				return RunResult{ExitCode: 1, Output: "shared noise\n"}
			}
			return RunResult{ExitCode: 1, Output: "unrelated noise\n"}
		}
		return RunResult{ExitCode: 1, Output: "shared noise\n"}
	})
	if err == nil {
		t.Fatal("Sweep err = nil, want the axi fixture reported vacuous")
	}
	if !strings.Contains(err.Error(), "canary 'own-package' EXPECT is vacuous") {
		t.Fatalf("Sweep err = %v, want own-package graded vacuous against its own baseline", err)
	}
	if strings.Contains(err.Error(), "other-package") {
		t.Fatalf("Sweep err = %v, want other-package graded against its own baseline alone", err)
	}
}

// contractFixture writes a biting behavior-owned fixture under pkg, together with the
// contract package directory the sweep resolves the binding against, and returns the
// fixture directory.
func contractFixture(t *testing.T, root, pkg, name string) string {
	t.Helper()
	dir := filepath.Join(root, "tests", "canary", "behavior-owned", filepath.FromSlash(pkg), name)
	fixture(t, dir, "")
	mkdir(t, filepath.Join(root, "internal", "contract", filepath.FromSlash(pkg)))
	return dir
}

// fixtureCalls lists the graded runs of one fixture by base name, so a count is what a
// test states rather than the first matching call.
func fixtureCalls(calls []RunCall, name string) []RunCall {
	var out []RunCall
	for _, call := range calls {
		if call.FixtureDir != "" && filepath.Base(call.FixtureDir) == name {
			out = append(out, call)
		}
	}
	return out
}

// baselineCalls lists the empty-tree baselines. A compile shares their empty FixtureDir
// and grades no tree, so it is excluded by kind.
func baselineCalls(calls []RunCall) []RunCall {
	var out []RunCall
	for _, call := range calls {
		if call.Kind != RunCompile && call.FixtureDir == "" {
			out = append(out, call)
		}
	}
	return out
}

// compileOutputs maps each compiled package to the binary path its compile wrote, and
// fails on a second compile of one package — the per-fixture compile this slice removes.
func compileOutputs(t *testing.T, calls []RunCall) map[string]string {
	t.Helper()
	out := map[string]string{}
	for _, call := range calls {
		if call.Kind != RunCompile {
			continue
		}
		if _, seen := out[call.Package]; seen {
			t.Errorf("package %s was compiled more than once, want one compile per group", call.Package)
		}
		out[call.Package] = call.Binary
	}
	return out
}

func envValues(env []string, key string) []string {
	var out []string
	for _, kv := range env {
		if value, ok := strings.CutPrefix(kv, key+"="); ok {
			out = append(out, value)
		}
	}
	return out
}
