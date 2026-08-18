package gate

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
)

func TestBenchkitPhasesUseBranchNativeDrivers(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	phases := BenchkitPhases(root, root)
	gotNames := make([]string, 0, len(phases))
	byName := map[string]Phase{}
	for _, phase := range phases {
		gotNames = append(gotNames, phase.Name)
		byName[phase.Name] = phase
	}
	wantNames := []string{"gofmt", "vet", "test", "race", "system", "shellcheck"}
	if !reflect.DeepEqual(gotNames, wantNames) {
		t.Fatalf("phase names = %v, want %v", gotNames, wantNames)
	}
	if want := []string{"go", "test", "-count=1", "./..."}; !reflect.DeepEqual(byName["test"].Argv, want) {
		t.Fatalf("ordinary test argv = %#v, want %#v", byName["test"].Argv, want)
	}
	if want := []string{"go", "test", "-count=1", "-tags=system", "./internal/systemtest"}; !reflect.DeepEqual(byName["system"].Argv, want) {
		t.Fatalf("system test argv = %#v, want %#v", byName["system"].Argv, want)
	}
	if want := raceDriverArgv(); !reflect.DeepEqual(byName["race"].Argv, want) {
		t.Fatalf("race test argv = %#v, want %#v", byName["race"].Argv, want)
	}
	wantEnv := []string{registry.ConformanceRootEnv + "=" + root, registry.ConformanceTierEnv + "=" + string(registry.Dev)}
	if !reflect.DeepEqual(byName["test"].Env, wantEnv) {
		t.Fatalf("ordinary test env = %#v, want %#v", byName["test"].Env, wantEnv)
	}
}

// TestLinkedRootTestPhaseCarriesNoConformanceEnv holds the materialization rule the
// race and system phases already follow: a graded root that does not declare the
// conformance entry test gets no variable its test binaries cannot honor. Dropping the
// probe and setting the env unconditionally turns this red.
func TestLinkedRootTestPhaseCarriesNoConformanceEnv(t *testing.T) {
	kit, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module linked\n\ngo 1.24\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, phase := range BenchkitPhases(root, kit) {
		if phase.Name != "test" {
			continue
		}
		if phase.Env != nil {
			t.Fatalf("linked-root test env = %#v, want none", phase.Env)
		}
		return
	}
	t.Fatal("linked root produced no test phase")
}

func TestPhaseTableGateRequiresSelectedBinaryRoute(t *testing.T) {
	root := t.TempDir()
	gateDir := filepath.Join(root, ".bench")
	if err := os.MkdirAll(gateDir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(gateDir, "gate.sh")
	if err := os.WriteFile(path, []byte("#!/bin/sh\nexec true gate-phases\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if phaseTableGate(root, Resolution{Kind: GateSh}) {
		t.Fatal("unrelated gate mentioning gate-phases claimed the selected-binary route")
	}
	if err := os.WriteFile(path, []byte("#!/bin/sh\nexec env BENCH_RUN_BINARY=\"$bench\" \"$bench\" gate-phases \"$root\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if !phaseTableGate(root, Resolution{Kind: ProspectiveGateSh}) {
		t.Fatal("selected-binary gate-phases route was not recognized")
	}
}
