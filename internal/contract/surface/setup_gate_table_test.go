package surface

import (
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

// setup_gate_table_test.go covers FT76 story 9: the deterministic gate-inference
// table, its one-question ambiguity posture, and the fail-closed zero-signal stub.

func TestSetupGateTableContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "bench setup proposes go test ./... for go.mod", testSetupGateGoMod)
	contract.RunParallel(t, "bench setup proposes npm test for a package.json test script", testSetupGateNpm)
	contract.RunParallel(t, "bench setup proposes cargo test for Cargo.toml", testSetupGateCargo)
	contract.RunParallel(t, "bench setup proposes make test for a Makefile test target", testSetupGateMake)
	contract.RunParallel(t, "bench setup zero-signal repo writes a fail-closed stub and exits partial", testSetupGateZeroSignal)
}

func testSetupGateGoMod(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("go.mod", "module example.com/fixture\n\ngo 1.22\n")

	probe := f.BenchWrapper("setup", "--yes")

	probe.RequireExit(0)
	requireFixtureFileContains(t, f, ".bench/gate.sh", "go test ./...", "gate.sh does not embed the proposed go.mod command")
}

func testSetupGateNpm(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("package.json", `{"name":"fixture","scripts":{"test":"echo ok"}}`+"\n")

	probe := f.BenchWrapper("setup", "--yes")

	probe.RequireExit(0)
	requireFixtureFileContains(t, f, ".bench/gate.sh", "npm test", "gate.sh does not embed the proposed package.json command")
}

func testSetupGateCargo(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("Cargo.toml", "[package]\nname = \"fixture\"\nversion = \"0.1.0\"\n")

	probe := f.BenchWrapper("setup", "--yes")

	probe.RequireExit(0)
	requireFixtureFileContains(t, f, ".bench/gate.sh", "cargo test", "gate.sh does not embed the proposed Cargo.toml command")
}

func testSetupGateMake(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("Makefile", "test:\n\techo ok\n")

	probe := f.BenchWrapper("setup", "--yes")

	probe.RequireExit(0)
	requireFixtureFileContains(t, f, ".bench/gate.sh", "make test", "gate.sh does not embed the proposed Makefile command")
}

func testSetupGateZeroSignal(t *testing.T) {
	f := contract.NewFixture(t)

	probe := f.BenchWrapper("setup", "--yes")

	probe.RequireExit(3)
	requireLinkFile(t, f, ".bench/gate.sh")
	requireExecutable(t, filepath.Join(f.Root, ".bench", "gate.sh"), "zero-signal gate.sh is not executable")
	stub := f.RunEnv(map[string]string{"BENCH_CANARY_INNER": "1"}, "bash", filepath.Join(f.Root, ".bench", "gate.sh"))
	if stub.ExitCode == 0 {
		t.Fatal("zero-signal gate stub exits zero; it must fail closed until configured")
	}
	requireFixtureFileContains(t, f, ".bench/gate.sh", "configure", "zero-signal gate stub does not name its configuration step")
	// Deferred row-9 cell: a fabricated green oracle must never hide behind a green
	// doctor row either - the unconfigured stub is red on its own doctor row, not just
	// non-zero when actually run.
	probe.RequireContains(probe.Stdout, "red: .bench/gate.sh")
}
