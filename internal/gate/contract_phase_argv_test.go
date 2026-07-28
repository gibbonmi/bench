package gate

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
)

// unnarrowedContractArgv is the whole contract subtree, which is what the phase grades in
// every mode — no run scopes it to one package any more.
var unnarrowedContractArgv = []string{"test", "-count=1", "./internal/contract/..."}

// TestContractPhaseGradesTheWholeSubtree pins the argv an inner gate runs to the one an
// outer gate runs. The two modes differ only in which phases survive the filter, never in
// what the contract phase grades.
func TestContractPhaseGradesTheWholeSubtree(t *testing.T) {
	kit := contractKit(t, "surface/artifact")
	root := t.TempDir()

	for name, mode := range map[string]phaseMode{"outer": outerMode, "inner": innerMode} {
		t.Run(name, func(t *testing.T) {
			assertContractArgv(t, root, kit, mode)
		})
	}
}

// TestRetiredContractPackageExportChangesNothing grades an ambient export of the variable
// the narrowing used to read. The name is spelled as a literal because no constant for it
// survives, which is the point: with no code reading it, the export is inert by
// construction rather than by a strip that has to be maintained.
func TestRetiredContractPackageExportChangesNothing(t *testing.T) {
	kit := contractKit(t, "surface/artifact")
	root := t.TempDir()
	t.Setenv("BENCH_CANARY_CONTRACT_PACKAGE", "surface/artifact")

	for name, mode := range map[string]phaseMode{"outer": outerMode, "inner": innerMode} {
		t.Run(name, func(t *testing.T) {
			assertContractArgv(t, root, kit, mode)
		})
	}
}

// contractKit is a kit checkout carrying pkg under internal/contract — the tree a narrowed
// argv would have resolved against, so a surviving narrowing rewrites rather than errors.
func contractKit(t *testing.T, pkg string) string {
	t.Helper()
	kit := t.TempDir()
	if err := os.MkdirAll(filepath.Join(kit, "internal", "contract", filepath.FromSlash(pkg)), 0o755); err != nil {
		t.Fatal(err)
	}
	return kit
}

// assertContractArgv grades the contract phase of the table the given mode actually runs,
// which is the only surface an argv rewrite could be observed through.
func assertContractArgv(t *testing.T, root, kit string, mode phaseMode) {
	t.Helper()
	phases, err := phaseTable(root, kit)
	if err != nil {
		t.Fatalf("phaseTable: %v", err)
	}
	phase, ok := phaseNamed(phasesForMode(phases, mode), canary.PhaseContract)
	if !ok {
		t.Fatalf("resolved table carries no %s phase: %v", canary.PhaseContract, phaseNames(phases))
	}
	want := append([]string{"go", "-C", kit}, unnarrowedContractArgv...)
	if !reflect.DeepEqual(phase.Argv, want) {
		t.Fatalf("contract argv = %#v, want %#v", phase.Argv, want)
	}
}
