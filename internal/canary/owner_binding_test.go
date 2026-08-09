package canary

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// TestSweepRefusesEveryTestBindingDefectTogether grades the four shapes a binding fails in,
// planted in one tree because each is a distinct lie and the author fixing them reads one
// sweep, not four. Absent and blank are separate answers everywhere else in this harness —
// deleting a marker is how a fixture asks for the default, so a blank one cannot be read as
// a deletion — and a binding outside the behavior-owned family or above its fixtures is read
// by nothing at all, which is why each diagnostic names the path that holds it. The refusal
// is what removes the silent fallback: without it every one of these fixtures runs its whole
// package while the tree claims per-test scoping.
func TestSweepRefusesEveryTestBindingDefectTogether(t *testing.T) {
	root := t.TempDir()
	canaryDir := filepath.Join(root, "tests", "canary")

	absent := contractFixture(t, root, "axi", "absent-owner")
	if err := os.Remove(filepath.Join(absent, testFileName)); err != nil {
		t.Fatal(err)
	}
	blank := contractFixture(t, root, "axi", "blank-owner")
	write(t, filepath.Join(blank, testFileName), "\n")

	// A conformance fixture is graded by a whole inner gate, which no owner can narrow, so
	// its binding is unread wherever the file sits inside it.
	strayFamily := canaryFixture(root, mappedFamily(t), "stray-fx")
	fixture(t, strayFamily, "")
	write(t, filepath.Join(strayFamily, testFileName), defaultOwner+"\n")

	// Both levels above a behavior-owned fixture: the family root and a package directory.
	// Nothing resolves an owner at either, so a binding left there names a test no run applies.
	familyLevel := filepath.Join(canaryDir, "behavior-owned", testFileName)
	write(t, familyLevel, defaultOwner+"\n")
	packageLevel := filepath.Join(canaryDir, "behavior-owned", "axi", testFileName)
	write(t, packageLevel, defaultOwner+"\n")

	calls, err := countedSweep(t, root)
	if err == nil {
		t.Fatal("Sweep err = nil, want every binding defect refused")
	}
	if calls != 0 {
		t.Errorf("sweep ran %d calls before refusing, want none", calls)
	}

	for _, want := range []struct {
		defect string
		lines  []string
	}{
		{defect: "absent binding", lines: []string{"absent-owner", filepath.Join(absent, testFileName), "no TEST file"}},
		{defect: "blank binding", lines: []string{"blank-owner", filepath.Join(blank, testFileName), "empty TEST file"}},
		{defect: "binding outside the family", lines: []string{filepath.Join(strayFamily, testFileName), "outside the behavior-owned family"}},
		{defect: "binding at the family level", lines: []string{familyLevel, "above the fixtures"}},
		{defect: "binding at the package level", lines: []string{packageLevel, "above the fixtures"}},
	} {
		for _, line := range want.lines {
			if !strings.Contains(err.Error(), line) {
				t.Errorf("Sweep err reports no %s: want a diagnostic naming %q, got:\n%s", want.defect, line, err)
			}
		}
	}
}

// TestSweepRefusesOwnersTheCompiledBinaryDoesNotCarry grades the names against the binary
// that would run them, before anything is graded against anything. A renamed or mistyped
// owner narrows a bite to a test that does not exist, which surfaces as a did-not-bite the
// author has to trace back to the marker one fixture at a time; the refusal names the fixture
// and the name it declared, and reports every such fixture in one sweep. The graded-run count
// is asserted because a refusal arriving after the bites has already paid for them.
func TestSweepRefusesOwnersTheCompiledBinaryDoesNotCarry(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "axi", "sound-owner")
	renamed := contractFixture(t, root, "axi", "renamed-owner")
	write(t, filepath.Join(renamed, testFileName), "TestOwnerRenamedAway\n")
	mistyped := contractFixture(t, root, "surface/artifact", "mistyped-owner")
	write(t, filepath.Join(mistyped, testFileName), "TestOnwer\n")

	var mu sync.Mutex
	var graded []RunCall
	err := Sweep(root, func(call RunCall) RunResult {
		if result, done := stubToolchain(call); done {
			return result
		}
		mu.Lock()
		graded = append(graded, call)
		mu.Unlock()
		return RunResult{ExitCode: 1, Output: "target-" + filepath.Base(call.FixtureDir) + "\n"}
	})

	if err == nil {
		t.Fatal("Sweep err = nil, want the unknown owners refused")
	}
	if len(graded) != 0 {
		t.Errorf("sweep made %d graded runs before refusing, want none", len(graded))
	}
	if got := strings.Split(err.Error(), "\n"); len(got) != 2 {
		t.Errorf("Sweep err carries %d lines, want one per defective fixture:\n%s", len(got), err)
	}
	for _, want := range []string{"renamed-owner", "TestOwnerRenamedAway", "mistyped-owner", "TestOnwer"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("Sweep err = %v, want a diagnostic naming %q", err, want)
		}
	}
	if strings.Contains(err.Error(), "sound-owner") {
		t.Errorf("Sweep err = %v, want the fixture whose owner the binary carries left unreported", err)
	}
}

func TestSweepAcceptsNamedSubtestWhenTheBinaryCarriesItsTopLevelOwner(t *testing.T) {
	root := t.TempDir()
	dir := contractFixture(t, root, "axi", "subtest-owner")
	write(t, filepath.Join(dir, testFileName), "TestCanaryOwner/owned_case\n")

	err := Sweep(root, func(call RunCall) RunResult {
		if result, done := stubToolchain(call); done {
			return result
		}
		if isBaseline(call) {
			return RunResult{ExitCode: 1, Output: "baseline\n"}
		}
		return RunResult{ExitCode: 1, Output: "target-subtest-owner\n"}
	})
	if err != nil {
		t.Fatalf("Sweep named subtest: %v", err)
	}
}

// TestSweepRedsWhenAPackageWillNotListItsTests covers the answer the validation cannot work
// without. A binary that refuses to say what it carries makes every owner in its group
// unresolvable, and a list failure read as an empty membership would refuse each of those
// fixtures for naming a test the package holds perfectly well — so the package is named and
// its fixtures are not.
func TestSweepRedsWhenAPackageWillNotListItsTests(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "axi", "axi-fx")

	var mu sync.Mutex
	var graded []RunCall
	err := Sweep(root, func(call RunCall) RunResult {
		if call.Kind == RunList {
			return RunResult{ExitCode: 2, Output: "exec format error\n"}
		}
		if result, done := stubToolchain(call); done {
			return result
		}
		mu.Lock()
		graded = append(graded, call)
		mu.Unlock()
		return RunResult{ExitCode: 1, Output: "target-" + filepath.Base(call.FixtureDir) + "\n"}
	})

	if err == nil {
		t.Fatal("Sweep err = nil, want the failed list reported")
	}
	for _, want := range []string{`"axi"`, "exec format error"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("Sweep err = %v, want a diagnostic carrying %q", err, want)
		}
	}
	if strings.Contains(err.Error(), "axi-fx") {
		t.Errorf("Sweep err = %v, want the package named rather than a fixture whose owner went ungraded", err)
	}
	if len(graded) != 0 {
		t.Errorf("sweep made %d graded runs after a failed list, want none", len(graded))
	}
}

// TestListedTestsAcceptsOnlyTestNames grades the membership the owner check is made against.
// The flag lists every runnable name a binary holds, and a run narrowed to a benchmark, a
// fuzz target, or an example matches no test at all — so admitting one turns a refusal the
// author can read into the did-not-bite hunt this validation exists to end.
func TestListedTestsAcceptsOnlyTestNames(t *testing.T) {
	listed := listedTests("TestOwner\nBenchmarkOwner\nFuzzOwner\nExampleOwner\nok  \tpkg\t0.01s\n")
	if !listed["TestOwner"] {
		t.Errorf("membership %v omits the test name the binary listed", listed)
	}
	for _, other := range []string{"BenchmarkOwner", "FuzzOwner", "ExampleOwner"} {
		if listed[other] {
			t.Errorf("membership admits %q, which no -test.run filter can reach", other)
		}
	}
}
