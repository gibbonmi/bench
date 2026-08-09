package canary

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
)

// TestBehaviorOwnedFixturesSpawnNoGate is the shape assertion the whole slice exists for.
// The cheapest wrong implementation adds compile-and-invoke beside the gate spawn it was
// meant to replace, reports the same green, and costs more rather than less; asserting the
// absence of the spawn kind is what that implementation fails.
func TestBehaviorOwnedFixturesSpawnNoGate(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "axi", "axi-fx")
	contractFixture(t, root, "surface/artifact", "artifact-fx")

	for _, call := range sweepCalls(t, root, registry.Dev) {
		if call.Kind == RunGate {
			t.Errorf("sweep spawned a gate for %q, want the behavior-owned family graded by compiled bites alone", call.FixtureDir)
		}
	}
}

// TestDefaultRunnerDispatchesOnCallKind grades the real runner rather than the recorded
// call metadata. It is the one assertion in this family that a relabelling cannot satisfy:
// every other test injects a fake, so a sweep that emitted correct kinds while still
// building `bash <gate>` for all of them would pass those and buy nothing.
func TestDefaultRunnerDispatchesOnCallKind(t *testing.T) {
	cases := []struct {
		name string
		call RunCall
		args []string
		dir  string
	}{
		{
			name: "gate spawn runs the gate script over the subject tree",
			call: RunCall{Kind: RunGate, Cwd: "/work", Gate: "/root/.bench/gate.sh"},
			args: []string{"bash", "/root/.bench/gate.sh"},
			dir:  "/work",
		},
		{
			name: "compile runs the go toolchain against the swept root",
			call: RunCall{Kind: RunCompile, Cwd: "/root", Package: "surface/artifact", Binary: "/bin/surface_sartifact.test"},
			args: []string{"go", "-C", "/root", "test", "-c", "-o", "/bin/surface_sartifact.test", "./internal/contract/surface/artifact"},
		},
		{
			name: "bite runs the compiled binary in its package source directory",
			call: RunCall{Kind: RunBite, Cwd: "/root/internal/contract/axi", Binary: "/bin/axi.test"},
			args: []string{"/bin/axi.test"},
			dir:  "/root/internal/contract/axi",
		},
		{
			name: "bite naming an owning test runs that test alone",
			call: RunCall{Kind: RunBite, Cwd: "/root/internal/contract/axi", Binary: "/bin/axi.test", Test: "TestOwner"},
			args: []string{"/bin/axi.test", "-test.run", "^TestOwner$"},
			dir:  "/root/internal/contract/axi",
		},
		{
			name: "bite naming an owning subtest runs only that path",
			call: RunCall{Kind: RunBite, Cwd: "/root/internal/contract/axi", Binary: "/bin/axi.test", Test: "TestOwner/owned_case"},
			args: []string{"/bin/axi.test", "-test.run", "^TestOwner$/^owned_case$"},
			dir:  "/root/internal/contract/axi",
		},
		{
			name: "list asks the compiled binary which tests it carries",
			call: RunCall{Kind: RunList, Cwd: "/root/internal/contract/axi", Binary: "/bin/axi.test"},
			args: []string{"/bin/axi.test", "-test.list", ".*"},
			dir:  "/root/internal/contract/axi",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := runnerCommand(tc.call)
			if !slices.Equal(cmd.Args, tc.args) {
				t.Errorf("command = %v, want %v", cmd.Args, tc.args)
			}
			if cmd.Dir != tc.dir {
				t.Errorf("command dir = %q, want %q", cmd.Dir, tc.dir)
			}
		})
	}
}

// TestPhaseNamedFamilyStillSpawnsItsGate is the other side of the migration: only the
// behavior-owned family loses its nested gate. The existing phase-family test asserts the
// phase pin's value alone, which a wholesale conversion carrying that pin on a compiled
// call would still satisfy, so the spawn kind is asserted directly here.
func TestPhaseNamedFamilyStillSpawnsItsGate(t *testing.T) {
	root := t.TempDir()
	fixture(t, canaryFixture(root, PhaseVet, "vet-fx"), "")

	got := fixtureCalls(sweepCalls(t, root, registry.Dev), "vet-fx")
	if len(got) != 1 {
		t.Fatalf("fixture ran %d graded runs, want exactly 1", len(got))
	}
	if got[0].Kind != RunGate || got[0].Gate == "" {
		t.Fatalf("phase-named fixture ran kind %v with gate %q, want a gate spawn", got[0].Kind, got[0].Gate)
	}
}

// TestSweepCompilesOncePerPackageGroup pins the saving the slice is for: the compile count
// follows the packages, not the fixtures. Compiling per fixture passes every bite
// assertion while buying nothing, which is precisely the failure being removed.
func TestSweepCompilesOncePerPackageGroup(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "axi", "axi-a")
	contractFixture(t, root, "axi", "axi-b")
	contractFixture(t, root, "surface/artifact", "artifact-fx")

	calls := sweepCalls(t, root, registry.Dev)
	compiled := compileOutputs(t, calls)
	if got := slices.Sorted(maps.Keys(compiled)); !slices.Equal(got, []string{"axi", "surface/artifact"}) {
		t.Fatalf("compiled packages = %v, want one compile for each of the two groups", got)
	}

	for _, name := range []string{"axi-a", "axi-b", "artifact-fx"} {
		got := fixtureCalls(calls, name)
		if len(got) != 1 {
			t.Fatalf("fixture %s ran %d graded runs, want exactly 1", name, len(got))
		}
		if want := compiled[got[0].Package]; got[0].Binary != want {
			t.Errorf("fixture %s ran binary %q, want its group's compile output %q", name, got[0].Binary, want)
		}
	}
}

// TestBiteCarriesItsOwnFixtureTreeAsSubjectRoot grades the isolation the whole family
// depends on. A shared root would grade every fixture against one tree, so all but one
// would report did-not-bite — or worse, all of them would bite for one fixture's reason.
func TestBiteCarriesItsOwnFixtureTreeAsSubjectRoot(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "axi", "axi-a")
	contractFixture(t, root, "axi", "axi-b")

	calls := sweepCalls(t, root, registry.Dev)
	seen := map[string]string{}
	for _, name := range []string{"axi-a", "axi-b"} {
		got := fixtureCalls(calls, name)
		if len(got) != 1 {
			t.Fatalf("fixture %s ran %d graded runs, want exactly 1", name, len(got))
		}
		roots := envValues(got[0].Env, SubjectRootEnv)
		if len(roots) != 1 {
			t.Fatalf("fixture %s carried subject roots %v, want exactly one", name, roots)
		}
		if prefix := "bench-canary-" + name + "-"; !strings.HasPrefix(filepath.Base(roots[0]), prefix) {
			t.Errorf("fixture %s graded %q, want its own materialized work directory", name, roots[0])
		}
		if other, clash := seen[roots[0]]; clash {
			t.Errorf("fixtures %s and %s graded the same tree %q", name, other, roots[0])
		}
		seen[roots[0]] = name
	}
}

// TestEachFixtureBitesTheTestItsMarkerNames grades the owner resolution per fixture rather
// than per package. Resolving one owner for a whole group — or letting the first fixture
// read win — satisfies any single-fixture assertion while grading every other fixture in
// the group against a test that is not its own.
func TestEachFixtureBitesTheTestItsMarkerNames(t *testing.T) {
	root := t.TempDir()
	owners := map[string]string{"axi-a": "TestOwnerA", "axi-b": "TestOwnerB"}
	for name, owner := range owners {
		bindOwner(t, root, "axi", contractFixture(t, root, "axi", name), owner)
	}

	calls := sweepCalls(t, root, registry.Dev)
	seen := map[string]string{}
	for name, owner := range owners {
		got := fixtureCalls(calls, name)
		if len(got) != 1 {
			t.Fatalf("fixture %s ran %d graded runs, want exactly 1", name, len(got))
		}
		if got[0].Test != owner {
			t.Errorf("fixture %s bit test %q, want %q", name, got[0].Test, owner)
		}
		filter := runFilter(t, got[0])
		if want := "^" + owner + "$"; filter != want {
			t.Errorf("fixture %s ran filter %q, want %q", name, filter, want)
		}
		if other, clash := seen[filter]; clash {
			t.Errorf("fixtures %s and %s both ran filter %q", name, other, filter)
		}
		seen[filter] = name
	}
}

// TestBiteFilterMatchesOnlyTheOwnerItNames grades the quoting rather than the filter's
// existence. An owner interpolated raw leaves its metacharacters reading as pattern syntax,
// so a name carrying a dot matches a superset of itself and the fixture is graded by
// whichever test that superset catches; an unanchored one matches every test the name is a
// substring of.
func TestBiteFilterMatchesOnlyTheOwnerItNames(t *testing.T) {
	const owner = "TestOwner.Case"
	filter := runFilter(t, RunCall{Kind: RunBite, Binary: "/bin/axi.test", Test: owner})
	pattern, err := regexp.Compile(filter)
	if err != nil {
		t.Fatalf("filter %q does not compile: %v", filter, err)
	}
	if !pattern.MatchString(owner) {
		t.Fatalf("filter %q does not match the owner %q it was built from", filter, owner)
	}
	for _, other := range []string{"TestOwnerXCase", "TestOwner.CaseTwo", "OuterTestOwner.Case"} {
		if pattern.MatchString(other) {
			t.Errorf("filter %q also matches %q", filter, other)
		}
	}
}

// TestContractBaselineRunsItsPackageWide keeps the vacuity screen wider than the runs it
// grades. A baseline narrowed to one fixture's owner prints a fraction of what its group's
// bites can, so an EXPECT the wide run already emits clears the screen in silence and every
// fixture in the group goes ungraded for vacuity.
func TestContractBaselineRunsItsPackageWide(t *testing.T) {
	root := t.TempDir()
	bindOwner(t, root, "axi", contractFixture(t, root, "axi", "axi-fx"), "TestOwnerA")

	var baselines int
	for _, call := range baselineCalls(sweepCalls(t, root, registry.Dev)) {
		if call.Kind != RunBite {
			continue
		}
		baselines++
		if call.Test != "" {
			t.Errorf("contract baseline named owner %q, want none", call.Test)
		}
		if filter := runFilter(t, call); filter != "" {
			t.Errorf("contract baseline ran filter %q, want its package whole", filter)
		}
	}
	if baselines != 1 {
		t.Fatalf("contract baselines = %d, want exactly 1", baselines)
	}
}

// runFilter is the -test.run value a call actually executes, empty for a call that runs its
// binary unfiltered. It reads the built argv rather than the call's own field, so an owner
// is graded where it takes effect instead of where it is recorded.
func runFilter(t *testing.T, call RunCall) string {
	t.Helper()
	args := runnerCommand(call).Args
	for idx, arg := range args {
		if arg != "-test.run" {
			continue
		}
		if idx+1 == len(args) {
			t.Fatalf("argv %v ends at -test.run, naming no filter", args)
		}
		return args[idx+1]
	}
	return ""
}

// TestAmbientSubjectRootDoesNotReachABite grades the strip half of the strip-then-set
// discipline. Go's exec environment has no duplicate-key precedence, so a set without its
// matching strip leaves two entries and hands an operator's ambient export control of what
// every fixture grades — which is why the count, not just the value, is asserted.
func TestAmbientSubjectRootDoesNotReachABite(t *testing.T) {
	t.Setenv(SubjectRootEnv, "/ambient/tree")

	root := t.TempDir()
	contractFixture(t, root, "axi", "axi-fx")

	for _, call := range sweepCalls(t, root, registry.Dev) {
		if call.Kind != RunBite {
			continue
		}
		roots := envValues(call.Env, SubjectRootEnv)
		if len(roots) != 1 {
			t.Fatalf("bite carried subject roots %v, want exactly the sweep's own", roots)
		}
		if roots[0] == "/ambient/tree" {
			t.Fatalf("bite graded the ambient tree %q", roots[0])
		}
	}
}

// TestContractGroupsDoNotFoldIntoTheUnscopedBaseline keeps the two baseline populations
// apart. The existing unscoped-baseline test builds no behavior-owned group at all, so a
// contract group folded into the flat fixtures' shared baseline would be invisible to it —
// and a baseline of the wrong shape both misses a genuinely vacuous EXPECT and flags a
// sound one.
func TestContractGroupsDoNotFoldIntoTheUnscopedBaseline(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "axi", "axi-fx")
	fixture(t, filepath.Join(root, "tests", "canary", "flat-fx"), "")

	var unscoped, contractual int
	for _, call := range baselineCalls(sweepCalls(t, root, registry.Dev)) {
		if call.Kind == RunBite {
			contractual++
			continue
		}
		unscoped++
	}
	if unscoped != 1 || contractual != 1 {
		t.Fatalf("baselines = %d unscoped + %d contract, want one of each", unscoped, contractual)
	}
}

// TestContractGroupBaselineThatPrintedNothingIsRefused reaches the empty-baseline refusal
// through the compiled path. The existing test drives conformance groups through the gate
// spawn, so a compiled path routing around the refusal would still pass it — and an empty
// baseline contains no EXPECT, so every fixture in the group would clear vacuity unguarded.
func TestContractGroupBaselineThatPrintedNothingIsRefused(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "axi", "axi-fx")

	var mu sync.Mutex
	var graded []string
	err := Sweep(root, func(call RunCall) RunResult {
		if result, done := stubToolchain(call); done {
			return result
		}
		if call.FixtureDir == "" {
			return RunResult{ExitCode: 1, Output: ""}
		}
		mu.Lock()
		graded = append(graded, filepath.Base(call.FixtureDir))
		mu.Unlock()
		return RunResult{ExitCode: 1, Output: "target-axi-fx\n"}
	})

	if err == nil {
		t.Fatal("Sweep err = nil, want the empty contract baseline reported")
	}
	if want := fmt.Sprintf("scope group %q", contractGroupPrefix+"axi"); !strings.Contains(err.Error(), want) {
		t.Fatalf("Sweep err = %v, want a diagnostic naming %s", err, want)
	}
	if len(graded) != 0 {
		t.Fatalf("fixtures %v ran against an ungradeable baseline, want none", graded)
	}
}

// TestFailedCompileRedsNamingItsPackage covers the two ways a compile leaves a group with
// no binary. A swallowed compile failure turns a broken package into a silently unswept
// family; the zero-exit case is not hypothetical, because `go test -c` on a package holding
// no test files succeeds and writes nothing, and an exit-code-only check would go on to
// invoke a path that does not exist.
func TestFailedCompileRedsNamingItsPackage(t *testing.T) {
	cases := []struct {
		name    string
		compile func(call RunCall) RunResult
	}{
		{
			// The binary is written before the failure is reported, so only the exit
			// code distinguishes this run: a sweep grading the binary's existence alone
			// would take a broken package's stale artifact for a good compile.
			name: "compile exits nonzero",
			compile: func(call RunCall) RunResult {
				stubToolchain(call)
				return RunResult{ExitCode: 2, Output: "undefined: Foo\n"}
			},
		},
		{
			name:    "compile exits zero and writes no binary",
			compile: func(RunCall) RunResult { return RunResult{} },
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			contractFixture(t, root, "axi", "axi-fx")

			var mu sync.Mutex
			var ran []RunCall
			err := Sweep(root, func(call RunCall) RunResult {
				if call.Kind == RunCompile {
					return tc.compile(call)
				}
				mu.Lock()
				ran = append(ran, call)
				mu.Unlock()
				return RunResult{ExitCode: 1, Output: "target-axi-fx\n"}
			})

			if err == nil {
				t.Fatal("Sweep err = nil, want the compile failure reported")
			}
			if !strings.Contains(err.Error(), `"axi"`) {
				t.Errorf("Sweep err = %v, want a diagnostic naming the package", err)
			}
			if len(ran) != 0 {
				t.Errorf("sweep ran %d graded runs for a group with no binary, want none", len(ran))
			}
		})
	}
}

// TestCompiledBinariesShareOneSweepOwnedDirectory keeps the artifacts out of the tree the
// gate grades: a binary written beside its source turns a sweep into a git-status change.
func TestCompiledBinariesShareOneSweepOwnedDirectory(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "axi", "axi-fx")
	contractFixture(t, root, "surface/artifact", "artifact-fx")

	var parents []string
	for _, binary := range compileOutputs(t, sweepCalls(t, root, registry.Dev)) {
		parents = append(parents, filepath.Dir(binary))
	}
	slices.Sort(parents)
	parents = slices.Compact(parents)
	if len(parents) != 1 {
		t.Fatalf("compiled binaries landed under %v, want one sweep-owned directory", parents)
	}
	if strings.HasPrefix(parents[0], root+string(filepath.Separator)) {
		t.Fatalf("compiled binaries landed under the swept tree at %q", parents[0])
	}
}

// TestBinaryDirectoryIsRemovedOnEveryErrorPath grades the exit paths the obvious
// implementation forgets: cleanup deferred after the happy path, or placed below the
// compile's own error return. Both leak a directory per failing sweep, and the leak is
// invisible until a disk fills — so each error return that can follow the first compile is
// driven separately.
func TestBinaryDirectoryIsRemovedOnEveryErrorPath(t *testing.T) {
	cases := []struct {
		name    string
		compile func(call RunCall) RunResult
		graded  RunResult
	}{
		{
			// The compile's own failure return is the earliest exit past the directory's
			// creation, and the one a defer written below it never reaches.
			name: "the compile itself fails",
			compile: func(call RunCall) RunResult {
				stubToolchain(call)
				return RunResult{ExitCode: 2, Output: "undefined: Foo\n"}
			},
		},
		{
			// A fixture that reds without its EXPECT is a did-not-bite, which is the
			// cheapest way to reach the sweep's last error return.
			name:   "a fixture does not bite",
			graded: RunResult{ExitCode: 1, Output: "unrelated failure\n"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			contractFixture(t, root, "axi", "axi-fx")

			var mu sync.Mutex
			var binDir string
			err := Sweep(root, func(call RunCall) RunResult {
				if call.Kind == RunCompile {
					mu.Lock()
					binDir = filepath.Dir(call.Binary)
					mu.Unlock()
					if tc.compile != nil {
						return tc.compile(call)
					}
					result, _ := stubToolchain(call)
					return result
				}
				if call.FixtureDir == "" {
					return RunResult{ExitCode: 1, Output: "baseline noise\n"}
				}
				return tc.graded
			})

			if err == nil {
				t.Fatal("Sweep err = nil, want the failing sweep reported")
			}
			if binDir == "" {
				t.Fatal("sweep compiled nothing, so the cleanup path was never reached")
			}
			if _, statErr := os.Stat(binDir); !os.IsNotExist(statErr) {
				t.Fatalf("binary directory %q survived a failing sweep (stat err = %v)", binDir, statErr)
			}
		})
	}
}

// TestPackagePathsCannotShareABinaryName grades the encoding on the pair that a naive
// separator-flattening collides: one group's binary would overwrite the other's mid-sweep,
// and the loser would grade the wrong package's tests.
func TestPackagePathsCannotShareABinaryName(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "a/bc", "left-fx")
	contractFixture(t, root, "ab/c", "right-fx")

	compiled := compileOutputs(t, sweepCalls(t, root, registry.Dev))
	if compiled["a/bc"] == compiled["ab/c"] {
		t.Fatalf("packages a/bc and ab/c both compiled to %q", compiled["a/bc"])
	}
}

// TestContractBinaryNameIsInjective grades the encoding directly on the pairs a
// fixture-tree test cannot reach, including a path already containing the escape
// character — the case that makes encoding the separator alone insufficient.
func TestContractBinaryNameIsInjective(t *testing.T) {
	paths := []string{"a/bc", "ab/c", "a_sb", "a/b", "a_ub", "a_b"}
	seen := map[string]string{}
	for _, pkg := range paths {
		name := contractBinaryName(pkg)
		if other, clash := seen[name]; clash {
			t.Errorf("packages %q and %q both name %q", other, pkg, name)
		}
		seen[name] = pkg
	}
}

// TestContractFixtureDeclaringItsOwnPhaseTableIsSwept grades the refusal's removal. A
// declared phase table replaced the built-in one a nested gate would have read; with no
// gate spawned for this family the file is an inert artifact in a subject tree, and the
// fixture is graded like any other.
func TestContractFixtureDeclaringItsOwnPhaseTableIsSwept(t *testing.T) {
	root := t.TempDir()
	fx := contractFixture(t, root, "axi", "manifest-fx")
	write(t, filepath.Join(fx, "files", "dot-bench", "phases.json"), "{}\n")

	got := fixtureCalls(sweepCalls(t, root, registry.Dev), "manifest-fx")
	if len(got) != 1 || got[0].Kind != RunBite {
		t.Fatalf("fixture ran %v, want exactly one bite", got)
	}
}
