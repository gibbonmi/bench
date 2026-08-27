package gate

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
)

func TestPhaseTableRequiresGoForBuiltInGoModule(t *testing.T) {
	for _, test := range []struct {
		name  string
		goMod string
	}{
		{name: "module declaration", goMod: "module example.com/graded\n"},
		{name: "empty regular file"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte(test.goMod), 0o644); err != nil {
				t.Fatal(err)
			}
			t.Setenv("PATH", t.TempDir())

			phases, err := phaseTable(root, root)
			if err == nil {
				t.Fatalf("phase table = %#v, want missing Go refusal", phases)
			}
			for _, want := range []string{"go", "PATH", root} {
				if !strings.Contains(err.Error(), want) {
					t.Fatalf("phase-table error = %q, want %q", err, want)
				}
			}
		})
	}
}

func TestPhaseTableWithoutGoModuleKeepsBuiltInNonGoPhases(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PATH", t.TempDir())

	phases, err := phaseTable(root, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"shellcheck"}; !reflect.DeepEqual(phaseNames(phases), want) {
		t.Fatalf("phase names = %v, want %v", phaseNames(phases), want)
	}
}

func TestPhaseTableManifestRemainsAuthoritativeWithoutGo(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/graded\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := []byte(`{"phases":[{"name":"declared","argv":["true"]}]}`)
	if err := os.WriteFile(filepath.Join(root, ".bench", "phases.json"), manifest, 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", t.TempDir())

	phases, err := phaseTable(root, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	want := []Phase{{Name: "declared", Argv: []string{"true"}, Dir: root}}
	if !reflect.DeepEqual(phases, want) {
		t.Fatalf("phases = %#v, want %#v", phases, want)
	}
}

func phaseNames(phases []Phase) []string {
	names := make([]string, 0, len(phases))
	for _, phase := range phases {
		names = append(names, phase.Name)
	}
	return names
}

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
	if want := []string{"go", "test", "-trimpath", "-count=1", "./..."}; !reflect.DeepEqual(byName["test"].Argv, want) {
		t.Fatalf("ordinary test argv = %#v, want %#v", byName["test"].Argv, want)
	}
	if want := []string{"go", "test", "-trimpath", "-count=1", "-tags=system", "./internal/systemtest"}; !reflect.DeepEqual(byName["system"].Argv, want) {
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
// race and system phases already follow. A graded root that does not declare the
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

// T02: the vet phase argv is verbatim `go -C <root> vet -trimpath ./...`. Vet compiles
// and writes facts entries of its own, so it keys on the package directory exactly as a
// test build does.
func TestVetPhaseArgvCarriesTrimPath(t *testing.T) {
	root := t.TempDir()
	writeLaneFile(t, filepath.Join(root, "go.mod"), "module example.com/x\n\ngo 1.24\n")
	for _, phase := range BenchkitPhases(root, root) {
		if phase.Name != "vet" {
			continue
		}
		if want := []string{"go", "-C", root, "vet", "-trimpath", "./..."}; !reflect.DeepEqual(phase.Argv, want) {
			t.Fatalf("vet phase argv = %#v, want %#v", phase.Argv, want)
		}
		return
	}
	t.Fatal("the kit phase table declares no vet phase")
}

// T05: the release core test step runs the enumerated packages under -trimpath and
// -count=1. Nothing else asserts this argv, so a release run that keys on its own
// checkout path would otherwise go unobserved.
func TestCoreTestArgvCarriesTrimPathAndCountOne(t *testing.T) {
	want := []string{"go", "test", "-trimpath", "-count=1", "./internal/gate", "./internal/toon"}
	if got := coreTestArgv([]string{"./internal/gate", "./internal/toon"}); !reflect.DeepEqual(got, want) {
		t.Fatalf("core test argv = %#v, want %#v", got, want)
	}
}
