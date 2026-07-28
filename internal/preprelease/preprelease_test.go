package preprelease

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/gate"
)

// TestStepsNameScriptsThisRepoShips grades the step plan against the tree it will run
// in. The surface contracts drive throwaway scripts at the same interface, which proves
// the orchestration but not that the interface still matches the kit's own scripts; a
// renamed or moved script is only observable here.
func TestStepsNameScriptsThisRepoShips(t *testing.T) {
	root := repoRoot(t)
	for _, step := range Steps(root, root) {
		if len(step.Argv) < 2 || step.Argv[0] != "bash" {
			continue
		}
		if _, err := os.Stat(step.Argv[1]); err != nil {
			t.Errorf("step %s runs %s, which this repo does not ship: %v", step.Name, step.Argv[1], err)
		}
	}
}

// TestStepsGradeTheRootTheyWereGiven pins the two arguments every scripted step is wrong
// without: the graded root, and — for the conformance suite alone — the kit whose module
// the suite compiles in. A linked repo is the case where those differ.
func TestStepsGradeTheRootTheyWereGiven(t *testing.T) {
	root, kit := filepath.Join("elsewhere", "graded root"), filepath.Join("other", "kit")
	steps := Steps(root, kit)

	for _, step := range steps {
		if step.Argv == nil && step.Run == nil {
			t.Errorf("step %s runs nothing at all", step.Name)
		}
		if step.Argv != nil && step.Run != nil {
			t.Errorf("step %s is both a subprocess and a library call", step.Name)
		}
	}

	conformance := stepNamed(t, steps, "conformance-ship")
	if !slices.Contains(conformance.Argv, kit) {
		t.Errorf("the conformance step compiles somewhere other than the kit: %v", conformance.Argv)
	}
	if !slices.Contains(conformance.Env, "BENCH_CONFORMANCE_ROOT="+root) {
		t.Errorf("the conformance step grades a root other than the one asked for: %v", conformance.Env)
	}
	if !slices.Contains(conformance.Env, registry.ConformanceTierEnv+"=ship") {
		t.Errorf("the conformance step does not ask for the ship tier: %v", conformance.Env)
	}
	if !slices.Contains(stepNamed(t, steps, "artifacts").Argv, root) {
		t.Errorf("the artifact build does not receive the graded root")
	}
}

// TestStepsRunReleaseOnlyPackages pins the step that keeps internal/preflight,
// internal/releaseevidence, and internal/publication in ship green. No dev-tier run
// executes those suites, so nothing else observes this step going missing, and an
// end-to-end prep-release run cannot stand in here because the command is blocked on an
// unrelated guard-enumeration leak.
func TestStepsRunReleaseOnlyPackages(t *testing.T) {
	root, kit := filepath.Join("elsewhere", "graded root"), filepath.Join("other", "kit")
	step := stepNamed(t, Steps(root, kit), "core-tests-ship")

	if want := gate.GateGoArgv(kit, "test", root); !slices.Equal(step.Argv, want) {
		t.Errorf("the release-only suites run as %v, want the shared gate-go argv %v", step.Argv, want)
	}
	if !slices.Contains(step.Env, registry.ConformanceTierEnv+"="+string(registry.Ship)) {
		t.Errorf("the release-only suites step does not ask for the ship tier: %v", step.Env)
	}
}

// TestArtifactStepIsHermetic keeps the release build off the dev tier's shared build
// cache. No step sets the opt-in, so this is a pin rather than a description: it holds on
// an unmodified tree and reds the moment an edit hands the token to the artifacts step.
// The ship tier catches that edit too, but only once per release and only after a full
// hermetic build. A token inherited from the surrounding shell is out of reach here by
// construction, since each step runs with append(os.Environ(), step.Env...).
func TestArtifactStepIsHermetic(t *testing.T) {
	step := stepNamed(t, Steps("root", "kit"), "artifacts")
	for _, env := range step.Env {
		if strings.Contains(env, contract.SharedBuildCacheEnv) {
			t.Errorf("the artifact build opts into the shared build cache via %q, so the release bytes are not hermetic", env)
		}
	}
}

// TestConformanceShipRunsTheStressAssertions pins the `-run` filter as the thing that
// makes `-tags stress` more than a compile: a filter naming only the entry point builds
// the stress-tagged files and executes nothing in them, so an assertion that exists on
// no other surface never runs anywhere. That each name is declared in the conformance
// package is graded there, where the declarations are.
func TestConformanceShipRunsTheStressAssertions(t *testing.T) {
	step := stepNamed(t, Steps("root", "kit"), "conformance-ship")
	filter := regexp.MustCompile(runFilter(t, step.Argv))
	if len(ShipConformanceTests) < 2 {
		t.Fatalf("the ship run names only %v, so the stress assertions compile and never run", ShipConformanceTests)
	}
	for _, name := range ShipConformanceTests {
		if !filter.MatchString(name) {
			t.Errorf("the ship conformance filter %q does not select %q", filter, name)
		}
	}
}

func runFilter(t *testing.T, argv []string) string {
	t.Helper()
	for i, arg := range argv {
		if arg == "-run" && i+1 < len(argv) {
			return argv[i+1]
		}
	}
	t.Fatalf("argv carries no -run filter: %v", argv)
	return ""
}

// TestRefusalNamesTheCauseAndTheRemedy covers the state the gate reports with no reason
// attached — an absent cache — where a message built only from Inspection.Reason would
// name nothing at all.
func TestRefusalNamesTheCauseAndTheRemedy(t *testing.T) {
	for _, test := range []struct {
		inspection gate.Inspection
		want       string
	}{
		{gate.Inspection{State: gate.Absent}, "absent"},
		{gate.Inspection{State: gate.Ready, Reason: "recorded red"}, "recorded red"},
	} {
		got := Refusal(test.inspection)
		if !strings.Contains(got, test.want) || !strings.Contains(got, "bench gate") {
			t.Errorf("Refusal(%+v) = %q, want the cause %q and the remedy `bench gate`", test.inspection, got, test.want)
		}
	}
}

// TestPreflightAttributionsNameEveryRedPhase covers the evidence set a real red run
// leaves: a green phase, a red one carrying its own message, a red one carrying none,
// and a phase whose record never landed. Reporting only the first red would hide the
// second cause, and a missing record is the interrupted run, not a failure to report.
func TestPreflightAttributionsNameEveryRedPhase(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, filepath.FromSlash(EvidenceDir)), 0o755); err != nil {
		t.Fatalf("create evidence dir: %v", err)
	}
	for name, body := range map[string]string{
		"gate":          `{"phase":"gate","status":"green","error":null}`,
		"vulnerability": `{"phase":"vulnerability","status":"red","error":{"kind":"tool","message":"govulncheck is missing"}}`,
		"race":          `{"phase":"race","status":"red","error":null}`,
	} {
		path := filepath.Join(root, filepath.FromSlash(EvidenceDir), name+".json")
		if err := os.WriteFile(path, []byte(body+"\n"), 0o644); err != nil {
			t.Fatalf("write %s evidence: %v", name, err)
		}
	}

	got := strings.Join(PreflightAttributions(root), "\n")
	for _, want := range []string{"phase race is red", "phase vulnerability is red: govulncheck is missing"} {
		if !strings.Contains(got, want) {
			t.Errorf("attributions = %q, want it to contain %q", got, want)
		}
	}
	if strings.Contains(got, "gate") {
		t.Errorf("attributions = %q, want no line for a green phase", got)
	}
	if strings.Contains(got, "smoke") {
		t.Errorf("attributions = %q, want no line for a phase with no record", got)
	}
}

func stepNamed(t *testing.T, steps []Step, name string) Step {
	t.Helper()
	for _, step := range steps {
		if step.Name == name {
			return step
		}
	}
	t.Fatalf("no step named %s", name)
	return Step{}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve repo root: runtime caller unavailable")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}
