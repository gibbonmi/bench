package probe

import (
	"slices"
	"testing"

	"github.com/gibbonmi/bench/internal/probe/probetest"
	"github.com/gibbonmi/bench/internal/toon"
)

// These rows run the real Go toolchain over the clamp repository, which holds
// TestClampNegative and TestClampPositive. The guard swap fails both tests and the return swap
// fails only TestClampPositive, so the ran and failed_tests cells show which tests a run
// selected.

// unpatternedSelection is the clamp package selection with no run pattern.
var unpatternedSelection = []string{"--package", "./"}

// guardArgs is the guard swap of clamp.go under the clamp package with no run pattern, then
// extra.
func guardArgs(extra ...string) []string {
	return slices.Concat([]string{probetest.ClampPath}, guardSwap, unpatternedSelection, extra)
}

// A single probe's --expect narrows the baseline and the mutated run to the named test, so
// the guard swap fails one test of the one that ran. Without --expect, both tests run.
func TestProbeExpectNarrowsTheSelection(t *testing.T) {
	for _, tc := range []struct {
		name   string
		extra  []string
		failed int
		run    string
		ran    int
	}{
		{"an expect test", []string{"--expect", "TestClampNegative"}, 1, "^TestClampNegative$", 1},
		{"no expect test", nil, 2, "all", 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := clampFixture(t)
			out, code := runProbe(t, guardArgs(tc.extra...)...)
			want := probeRow(t, "bit", "clamp.go", "swap", "failed", tc.failed, "TestClampNegative", "yes") +
				selectionRow(t, "package", "./", tc.run, "passed", tc.ran)
			requireRow(t, out, code, 0, want)
			requireSubjectBytes(t, f, probetest.ClampSource)
		})
	}
}

// An explicit --run governs the selection, so a witness outside it reads expect-missed.
func TestProbeRunOverridesTheExpectNarrowing(t *testing.T) {
	clampFixture(t)
	out, code := runProbe(t, guardArgs("--expect", "TestClampPositive", "--run", "^TestClampNegative$")...)
	want := probeRow(t, "invalid", "clamp.go", "swap", "expect-missed", 1, "TestClampNegative", "yes") +
		selectionRow(t, "package", "./", "^TestClampNegative$", "passed", 1)
	requireRow(t, out, code, 1, want)
}

// A subtest witness selects its parent and only that subtest. A parent that fails outside the
// named subtest still reads expect-missed.
func TestProbeSubtestExpectMissesAParentFailure(t *testing.T) {
	clampFixture(t)
	out, code := runProbe(t, guardArgs("--expect", "TestClampNegative/zero")...)
	want := probeRow(t, "invalid", "clamp.go", "swap", "expect-missed", 1, "TestClampNegative", "yes") +
		selectionRow(t, "package", "./", "^TestClampNegative$/^zero$", "passed", 1)
	requireRow(t, out, code, 1, want)
}

// The narrowed pattern matches the whole name literally, so a name prefix and a name that
// holds a regular-expression metacharacter select no test, and the baseline reads no-test-run.
func TestProbeExpectMatchesTheWholeNameLiterally(t *testing.T) {
	for _, tc := range []struct{ expect, run string }{
		{"TestClampNeg", "^TestClampNeg$"},
		{"TestClamp.egative", `^TestClamp\.egative$`},
	} {
		t.Run(tc.expect, func(t *testing.T) {
			f := clampFixture(t)
			out, code := runProbe(t, guardArgs("--expect", tc.expect)...)
			want := probeRow(t, "invalid", "clamp.go", "swap", "baseline-no-test-run", 0, nil, "untouched") +
				selectionRow(t, "package", "./", tc.run, "no-test-run", 0)
			requireRow(t, out, code, 1, want)
			requireSubjectBytes(t, f, probetest.ClampSource)
		})
	}
}

// Each set entry with an --expect test runs only that test over the one held run binary, and
// the baseline runs the union of the entry tests. An entry with no --expect test makes the
// baseline and its own run take the whole selection.
func TestMutantSetExpectNarrowsEachEntry(t *testing.T) {
	positiveBite := func(n, ran int) []any {
		return []any{n, "bit", "swap", "failed", 1, "TestClampPositive", ran, 1, 0, nil, "yes"}
	}
	expect := func(entry []string, test string) []string { return append(slices.Clone(entry), "--expect", test) }
	for _, tc := range []struct {
		name    string
		entries [][]string
		rows    [][]any
		run     string
	}{
		{"every entry names a test", [][]string{expect(guardSwap, "TestClampNegative"), expect(returnSwap, "TestClampPositive")},
			[][]any{biteRow(1, "swap"), positiveBite(2, 1)}, "^TestClampNegative$|^TestClampPositive$"},
		{"one entry names no test", [][]string{expect(guardSwap, "TestClampNegative"), returnSwap},
			[][]any{biteRow(1, "swap"), positiveBite(2, 2)}, "all"},
	} {
		for _, binary := range runBinaryCases {
			t.Run(tc.name+"/"+binary.name, func(t *testing.T) {
				f := clampFixture(t)
				if binary.owned {
					buildEachRun(t, f, "clamp.go")
				}
				out, code := runProbe(t, setArgs(mutantsFile(t, tc.entries...), unpatternedSelection...)...)
				selection, err := toon.TableTyped("selection", setSelectionFields, [][]any{{"package", "./", tc.run, "passed", binary.builds}})
				if err != nil {
					t.Fatal(err)
				}
				requireRow(t, out, code, 0, setRows(t, tc.rows...)+selection)
				requireSubjectBytes(t, f, probetest.ClampSource)
			})
		}
	}
}
