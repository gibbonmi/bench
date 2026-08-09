package surface

import (
	"github.com/gibbonmi/bench/internal/contract"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

// TestDoctorRowContracts pins story 11 of specs/repo-aware-setup.md: doctor asserts
// per-harness rows independently, so a broken row is always visible even when every
// other signal (including shim health) is green. Each row gets its own red fixture —
// never one aggregate — matching the acceptance map's row-11 red signal.
func TestDoctorRowContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "bench doctor AGENTS.md row contract", testDoctorAgentsRow)
	contract.RunParallel(t, "bench doctor CLAUDE.md row contract", testDoctorClaudeRow)
	contract.RunParallel(t, "bench doctor gate row contract", testDoctorGateRow)
	contract.RunParallel(t, "bench doctor profile row contract", testDoctorProfileRow)
	contract.RunParallel(t, "bench doctor repo-local bench row contract", testDoctorRepoLocalBenchRow)
	contract.RunParallel(t, "bench doctor setup pointers row contract", testDoctorSetupPointersRow)
	contract.RunParallel(t, "bench doctor all-rows-green contract", testDoctorAllRowsGreen)
	contract.RunParallel(t, "bench doctor AGENTS.md row against a FIFO never hangs", testDoctorAgentsRowFIFO)
	contract.RunParallel(t, "bench doctor CLAUDE.md row against a FIFO never hangs", testDoctorClaudeRowFIFO)
}

// testDoctorAgentsRowFIFO drives evalAgentsRow against a FIFO at AGENTS.md's fixed
// inspection path: os.ReadFile on a FIFO with no writer on the other end blocks
// forever, so the row must route to isSpecialFile's guard and report red without ever
// attempting the read. The bounded RunAtWithTimeout turns a regression into a fast,
// named failure instead of a stuck test run.
func testDoctorAgentsRowFIFO(t *testing.T) {
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	writeDoctorGreenTree(t, f)
	f.BenchWrapperEnv(sb.env, "doctor", "--fix").RequireExit(0)

	contract.Remove(t, f.Root+"/AGENTS.md")
	if err := syscall.Mkfifo(f.Root+"/AGENTS.md", 0o644); err != nil {
		t.Fatalf("mkfifo AGENTS.md: %v", err)
	}

	bench := filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh")
	probe := contract.RunAtWithTimeout(t, f, f.Root, sb.env, 15*time.Second, "bash", bench, "doctor")
	probe.RequireExit(1)
	probe.RequireContains(probe.Stdout, "red: AGENTS.md")
}

// testDoctorClaudeRowFIFO is the CLAUDE.md mirror of testDoctorAgentsRowFIFO.
func testDoctorClaudeRowFIFO(t *testing.T) {
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	writeDoctorGreenTree(t, f)
	f.BenchWrapperEnv(sb.env, "doctor", "--fix").RequireExit(0)

	contract.Remove(t, f.Root+"/CLAUDE.md")
	if err := syscall.Mkfifo(f.Root+"/CLAUDE.md", 0o644); err != nil {
		t.Fatalf("mkfifo CLAUDE.md: %v", err)
	}

	bench := filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh")
	probe := contract.RunAtWithTimeout(t, f, f.Root, sb.env, 15*time.Second, "bash", bench, "doctor")
	probe.RequireExit(1)
	probe.RequireContains(probe.Stdout, "red: CLAUDE.md")
}

// writeDoctorGreenTree writes every asset a converged repo carries so a single row can
// be flipped red while every other row (and the shim/pre-push rows) stays green — the
// only way to prove a row is independent rather than folded into an aggregate.
func writeDoctorGreenTree(t *testing.T, f contract.Fixture) {
	t.Helper()
	f.WriteFile("AGENTS.md", "# Working agreement\n\n<!-- bench:start -->\nBench block\n<!-- bench:end -->\n")
	f.WriteFile("CLAUDE.md", "# Bench\n\nCanonical agreement in AGENTS.md; platform rules in .bench/BENCH.md.\n\n@AGENTS.md\n@.bench/BENCH.md\n")
	f.WriteFile("projects/example.md", "# Project: example\n")
	// os.WriteFile (behind WriteFile/WriteExecutable) only applies its mode when creating
	// a file, not when overwriting one that already exists — an explicit chmod after the
	// write is what actually restores exec bits a prior fixture step stripped.
	writeDoctorExecutable(t, f, ".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	writeDoctorExecutable(t, f, ".bench/bin/bench.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.WriteFile(".agents/commands/bench-setup-repo.md", "# /bench-setup-repo\n")
	mustWriteFile(t, f.Root+"/.git/hooks/pre-push", managedPrePushBody, 0o755)
}

func writeDoctorExecutable(t *testing.T, f contract.Fixture, path, contents string) {
	t.Helper()
	f.WriteExecutable(path, contents)
	if err := os.Chmod(f.Root+"/"+path, 0o755); err != nil {
		t.Fatalf("chmod %s: %v", path, err)
	}
}

// testDoctorSetupPointersRow pins the deferred row-11 sub-row: doctor validates that
// the pointer setup's printed next action relies on - the /bench-setup-repo command
// file - actually exists, so a false "continue the /bench-setup-repo conversation"
// print is never left unbacked.
func testDoctorSetupPointersRow(t *testing.T) {
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	writeDoctorGreenTree(t, f)
	f.BenchWrapperEnv(sb.env, "doctor", "--fix").RequireExit(0)

	contract.Remove(t, f.Root+"/.agents/commands/bench-setup-repo.md")
	probe := f.BenchWrapperEnv(sb.env, "doctor")
	probe.RequireExit(1)
	probe.RequireContains(probe.Stdout, "red:")
	probe.RequireContains(probe.Stdout, "bench-setup-repo.md")

	writeDoctorGreenTree(t, f)
	probe = f.BenchWrapperEnv(sb.env, "doctor")
	probe.RequireContains(probe.Stdout, "ok:")
	probe.RequireContains(probe.Stdout, "bench-setup-repo.md")
}

func testDoctorAllRowsGreen(t *testing.T) {
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	writeDoctorGreenTree(t, f)
	f.BenchWrapperEnv(sb.env, "doctor", "--fix").RequireExit(0)

	probe := f.BenchWrapperEnv(sb.env, "doctor")
	probe.RequireExit(0)
	for _, want := range []string{"AGENTS.md", "CLAUDE.md", "gate", "profile", "repo-local bench"} {
		probe.RequireContains(probe.Stdout, want)
	}
}

func testDoctorAgentsRow(t *testing.T) {
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	writeDoctorGreenTree(t, f)
	f.BenchWrapperEnv(sb.env, "doctor", "--fix").RequireExit(0)

	// Absent — no AGENTS.md at all.
	contract.Remove(t, f.Root+"/AGENTS.md")
	probe := f.BenchWrapperEnv(sb.env, "doctor")
	probe.RequireExit(1)
	probe.RequireContains(probe.Stdout, "red: AGENTS.md")

	// Present but no managed block.
	f.WriteFile("AGENTS.md", "# Working agreement\n\nno bench block here\n")
	probe = f.BenchWrapperEnv(sb.env, "doctor")
	probe.RequireExit(1)
	probe.RequireContains(probe.Stdout, "red: AGENTS.md")

	// Malformed markers — a lone start with no end.
	f.WriteFile("AGENTS.md", "# Working agreement\n\n<!-- bench:start -->\nunterminated\n")
	probe = f.BenchWrapperEnv(sb.env, "doctor")
	probe.RequireExit(1)
	probe.RequireContains(probe.Stdout, "red: AGENTS.md")

	// Restored to green.
	writeDoctorGreenTree(t, f)
	probe = f.BenchWrapperEnv(sb.env, "doctor")
	probe.RequireContains(probe.Stdout, "ok: AGENTS.md")
}

// testDoctorClaudeRow pins the spec's headline cell: a preserved CLAUDE.md that exists
// but lacks the marker-owned import lines is red on its own, never hidden inside an
// aggregate green even though every other row (including the shim) stays healthy.
func testDoctorClaudeRow(t *testing.T) {
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	writeDoctorGreenTree(t, f)
	f.BenchWrapperEnv(sb.env, "doctor", "--fix").RequireExit(0)

	// Preserved project content, imports stripped — the spec's named red cell.
	f.WriteFile("CLAUDE.md", "# My project\n\nProject-owned notes, no Bench imports.\n")
	probe := f.BenchWrapperEnv(sb.env, "doctor")
	probe.RequireExit(1)
	probe.RequireContains(probe.Stdout, "red: CLAUDE.md")
	probe.RequireContains(probe.Stdout, "ok: AGENTS.md")
	probe.RequireContains(probe.Stdout, "ok: gate")

	// Absent entirely.
	contract.Remove(t, f.Root+"/CLAUDE.md")
	probe = f.BenchWrapperEnv(sb.env, "doctor")
	probe.RequireExit(1)
	probe.RequireContains(probe.Stdout, "red: CLAUDE.md")

	// Legacy form (import line only, no .bench/BENCH.md line) counts as green.
	f.WriteFile("CLAUDE.md", "# Bench\n\nCanonical agreement in AGENTS.md.\n\n@AGENTS.md\n")
	probe = f.BenchWrapperEnv(sb.env, "doctor")
	probe.RequireContains(probe.Stdout, "ok: CLAUDE.md")

	// Canonical form restored.
	writeDoctorGreenTree(t, f)
	probe = f.BenchWrapperEnv(sb.env, "doctor")
	probe.RequireContains(probe.Stdout, "ok: CLAUDE.md")
}

func testDoctorGateRow(t *testing.T) {
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	writeDoctorGreenTree(t, f)
	f.BenchWrapperEnv(sb.env, "doctor", "--fix").RequireExit(0)

	// Absent.
	contract.Remove(t, f.Root+"/.bench/gate.sh")
	probe := f.BenchWrapperEnv(sb.env, "doctor")
	probe.RequireExit(1)
	probe.RequireContains(probe.Stdout, "red: .bench/gate.sh")

	// Present but not executable.
	f.WriteFile(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	probe = f.BenchWrapperEnv(sb.env, "doctor")
	probe.RequireExit(1)
	probe.RequireContains(probe.Stdout, "not executable")

	// Restored to green.
	writeDoctorGreenTree(t, f)
	probe = f.BenchWrapperEnv(sb.env, "doctor")
	probe.RequireContains(probe.Stdout, "ok: gate")
}

func testDoctorProfileRow(t *testing.T) {
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	writeDoctorGreenTree(t, f)
	f.BenchWrapperEnv(sb.env, "doctor", "--fix").RequireExit(0)

	// No projects/ directory at all.
	contract.Remove(t, f.Root+"/projects")
	probe := f.BenchWrapperEnv(sb.env, "doctor")
	probe.RequireExit(1)
	probe.RequireContains(probe.Stdout, "red: no projects/")

	// projects/ exists but carries no .md profile.
	f.WriteFile("projects/.keep", "")
	probe = f.BenchWrapperEnv(sb.env, "doctor")
	probe.RequireExit(1)
	probe.RequireContains(probe.Stdout, "red: no projects/")

	// Restored to green.
	writeDoctorGreenTree(t, f)
	probe = f.BenchWrapperEnv(sb.env, "doctor")
	probe.RequireContains(probe.Stdout, "ok: profile")
}

func testDoctorRepoLocalBenchRow(t *testing.T) {
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	writeDoctorGreenTree(t, f)
	f.BenchWrapperEnv(sb.env, "doctor", "--fix").RequireExit(0)

	// Absent.
	contract.Remove(t, f.Root+"/.bench/bin")
	probe := f.BenchWrapperEnv(sb.env, "doctor")
	probe.RequireExit(1)
	probe.RequireContains(probe.Stdout, "red: .bench/bin/bench.sh")

	// Present but not executable.
	f.WriteFile(".bench/bin/bench.sh", "#!/usr/bin/env bash\nexit 0\n")
	probe = f.BenchWrapperEnv(sb.env, "doctor")
	probe.RequireExit(1)
	probe.RequireContains(probe.Stdout, "not executable")

	// Restored to green.
	writeDoctorGreenTree(t, f)
	probe = f.BenchWrapperEnv(sb.env, "doctor")
	probe.RequireContains(probe.Stdout, "ok: repo-local bench")
}
