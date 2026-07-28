package canary

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
)

// TestSweepDiscoversContractFixturesUnderTheirPackage grades the walk's whole answer for
// a nested package: the fixture is the directory holding EXPECT, it is enumerated once
// under its own base name, and the package it carries is every segment between the family
// and it. A walk keeping only the last segment binds the fixture to a package that does
// not exist, which is why the multi-segment path is the case under test.
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
		if values := envValues(got[0].Env, ContractPackageEnv); !slices.Equal(values, []string{pkg}) {
			t.Errorf("fixture %s carried packages %v, want exactly [%s]", name, values, pkg)
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

// TestSweepRejectsStructurallyUnscopableContractFixtures grades the three structural
// defects that would each leave a behavior-owned fixture paying the full contract suite
// in silence — the cost this scoping exists to remove. Every diagnostic names the
// fixture, because the sweep's tree is the only place the defect can be repaired.
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
		{
			name: "fixture declares a phase manifest",
			build: func(t *testing.T, root string) {
				fx := contractFixture(t, root, "axi", "manifest-fx")
				write(t, filepath.Join(fx, "files", "dot-bench", "phases.json"), "{}\n")
			},
			want: []string{"manifest-fx", "phases.json", "flat"},
		},
		{
			// The dot- prefix is a storage convention, not the only way to write the
			// path: a literal .bench directory materializes the same manifest-declaring
			// root, so a check reading one spelling passes the other through unscoped.
			name: "fixture declares a phase manifest in a literal dot directory",
			build: func(t *testing.T, root string) {
				fx := contractFixture(t, root, "axi", "literal-manifest-fx")
				write(t, filepath.Join(fx, "files", ".bench", "phases.json"), "{}\n")
			},
			want: []string{"literal-manifest-fx", "phases.json", "flat"},
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

// TestContractFixtureCarriesExactlyItsPhaseAndPackage pins the whole control set a
// behavior-owned fixture's inner gate runs under. The conformance check variable is
// asserted absent rather than ignored: it is what a scope resolved from the fixture's
// parent directory would leak in, and an inner gate carrying one grades a conformance
// check instead of the contract package that owns the EXPECT.
func TestContractFixtureCarriesExactlyItsPhaseAndPackage(t *testing.T) {
	t.Setenv(PhaseEnv, "ambient-phase")
	t.Setenv(ContractPackageEnv, "ambient/package")
	t.Setenv(registry.ConformanceCheckEnv, "ambient-check")

	root := t.TempDir()
	contractFixture(t, root, "surface/artifact", "scoped-fx")

	var fixtureCall RunCall
	for _, call := range sweepCalls(t, root, registry.Dev) {
		if call.FixtureDir != "" {
			fixtureCall = call
		}
	}
	want := []string{ContractPackageEnv + "=surface/artifact", PhaseEnv + "=contract"}
	if got := controlValues(fixtureCall.Env); !slices.Equal(got, want) {
		t.Fatalf("fixture control env = %v, want exactly %v", got, want)
	}
}

// TestSweepBaselinesContractGroupsPerPackage grades the vacuity groups two packages
// produce: one baseline each, and each baseline running under the same phase and package
// its own fixtures carry. A single shared baseline — the degenerate grouping — fails on
// the call count and on every env comparison at once.
func TestSweepBaselinesContractGroupsPerPackage(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "axi", "axi-a")
	contractFixture(t, root, "axi", "axi-b")
	contractFixture(t, root, "surface/artifact", "artifact-fx")

	var baselines [][]string
	for _, call := range sweepCalls(t, root, registry.Dev) {
		if call.FixtureDir == "" {
			baselines = append(baselines, controlValues(call.Env))
		}
	}
	slices.SortFunc(baselines, func(a, b []string) int { return slices.Compare(a, b) })
	want := [][]string{
		{ContractPackageEnv + "=axi", PhaseEnv + "=contract"},
		{ContractPackageEnv + "=surface/artifact", PhaseEnv + "=contract"},
	}
	if len(baselines) != len(want) {
		t.Fatalf("sweep ran %d baselines (%v), want one per contract package", len(baselines), baselines)
	}
	for idx := range want {
		if !slices.Equal(baselines[idx], want[idx]) {
			t.Errorf("baseline control env = %v, want %v", baselines[idx], want[idx])
		}
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
		if call.FixtureDir == "" {
			if slices.Contains(call.Env, ContractPackageEnv+"=axi") {
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

// controlValues lists the sorted scope variables one inner gate runs under, so an
// assertion states the whole set a run carries rather than the presence of one entry.
func controlValues(env []string) []string {
	var out []string
	for _, key := range []string{PhaseEnv, ContractPackageEnv, registry.ConformanceCheckEnv} {
		for _, value := range envValues(env, key) {
			out = append(out, key+"="+value)
		}
	}
	slices.Sort(out)
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
