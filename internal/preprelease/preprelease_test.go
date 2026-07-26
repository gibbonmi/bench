package preprelease

import (
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

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
	if !slices.Contains(conformance.Env, ConformanceTierEnv+"=ship") {
		t.Errorf("the conformance step does not ask for the ship tier: %v", conformance.Env)
	}
	if !slices.Contains(stepNamed(t, steps, "artifacts").Argv, root) {
		t.Errorf("the artifact build does not receive the graded root")
	}
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
