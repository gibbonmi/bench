package surface

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
)

// setup_conflict_test.go covers FT76 story 4: conflicts are never merged or
// overwritten. A confirmed run writes every non-conflicting asset, preserves each
// conflicting asset byte-identical, and exits with the link lifecycle's partial
// status plus its machine-readable conflict list.

func TestSetupConflictContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "bench setup preserves conflicts and writes the rest in one run", testSetupPreservesConflictsAndWritesRest)
	contract.RunParallel(t, "bench setup treats a FIFO at AGENTS.md as a conflict, never hanging", testSetupFIFOAgentsNeverHangs)
	contract.RunParallel(t, "bench setup treats a FIFO at CLAUDE.md as a conflict, never hanging", testSetupFIFOClaudeNeverHangs)
	contract.RunParallel(t, "bench setup treats a directory at AGENTS.md as a conflict", testSetupDirectoryAtAgentsIsConflict)
	contract.RunParallel(t, "bench setup treats a directory at CLAUDE.md as a conflict", testSetupDirectoryAtClaudeIsConflict)
}

// testSetupPreservesConflictsAndWritesRest fixtures two simultaneous conflicts — a
// foreign .bench/gate.sh and a foreign managed command file — proving both are
// preserved byte-identical while the rest of the plan (AGENTS.md's managed block,
// CLAUDE.md's import lines) still converges in the same run.
func testSetupPreservesConflictsAndWritesRest(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("go.mod", "module example.com/fixture\n\ngo 1.22\n")
	const foreignGate = "#!/usr/bin/env bash\necho foreign gate\nexit 0\n"
	const foreignCommand = "project-owned command content\n"
	f.WriteExecutable(".bench/gate.sh", foreignGate)
	f.WriteFile(".agents/commands/bench-implement-spec.md", foreignCommand)

	probe := f.Bench("setup", "--yes")

	probe.RequireExit(3)
	probe.RequireContains(strings.ToLower(probe.Stdout+probe.Stderr), "conflict")
	probe.RequireContains(probe.Stdout, ".bench/gate.sh")
	probe.RequireContains(probe.Stdout, ".agents/commands/bench-implement-spec.md")
	// The verdict vocabulary, not just the word "conflict": a never-managed foreign
	// file reads "project-owned", the reason distinct from a tampered managed asset's
	// "modified-managed" (pinned by the rerun contract's hand-edit case).
	probe.RequireContains(probe.Stdout, "project-owned")
	if got := f.ReadFile(".bench/gate.sh"); got != foreignGate {
		t.Fatalf(".bench/gate.sh not preserved byte-identical:\n%s", got)
	}
	if got := f.ReadFile(".agents/commands/bench-implement-spec.md"); got != foreignCommand {
		t.Fatalf("conflicting command file not preserved byte-identical:\n%s", got)
	}
	requireLinkFile(t, f, "AGENTS.md")
	requireLiteralCount(t, f, "AGENTS.md", "<!-- bench:start -->", 1, "non-conflicting AGENTS.md was not written alongside the preserved conflicts")
	requireFixtureFileContains(t, f, "CLAUDE.md", "@AGENTS.md", "non-conflicting CLAUDE.md was not written alongside the preserved conflicts")
}

// testSetupFIFOAgentsNeverHangs plants a FIFO at the fixed AGENTS.md inspection path.
// A regression that opens the path for read (rather than routing straight to a
// conflict) would hang forever with no writer on the other end of the pipe — the
// bounded RunAtWithTimeout is what turns that hang into a fast, named test failure
// instead of a stuck CI job.
func testSetupFIFOAgentsNeverHangs(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("go.mod", "module example.com/fixture\n\ngo 1.22\n")
	agentsPath := filepath.Join(f.Root, "AGENTS.md")
	if err := syscall.Mkfifo(agentsPath, 0o644); err != nil {
		t.Fatalf("mkfifo AGENTS.md: %v", err)
	}

	bench := filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh")
	probe := contract.RunAtWithTimeout(t, f, f.Root, nil, 15*time.Second, "bash", bench, "setup", "--yes")

	probe.RequireExit(3)
	probe.RequireContains(strings.ToLower(probe.Stdout+probe.Stderr), "conflict")
	probe.RequireContains(probe.Stdout, "AGENTS.md")
	requireLinkFile(t, f, ".bench/gate.sh")
	requireFixtureFileContains(t, f, "CLAUDE.md", "@AGENTS.md", "non-conflicting CLAUDE.md was not written when AGENTS.md was a FIFO")
}

// testSetupFIFOClaudeNeverHangs is testSetupFIFOAgentsNeverHangs's CLAUDE.md mirror:
// os.ReadFile on a FIFO with no writer on the other end blocks forever, so a FIFO at
// CLAUDE.md must route straight to a conflict instead of ever being opened for read.
func testSetupFIFOClaudeNeverHangs(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("go.mod", "module example.com/fixture\n\ngo 1.22\n")
	claudePath := filepath.Join(f.Root, "CLAUDE.md")
	if err := syscall.Mkfifo(claudePath, 0o644); err != nil {
		t.Fatalf("mkfifo CLAUDE.md: %v", err)
	}

	bench := filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh")
	probe := contract.RunAtWithTimeout(t, f, f.Root, nil, 15*time.Second, "bash", bench, "setup", "--yes")

	probe.RequireExit(3)
	probe.RequireContains(strings.ToLower(probe.Stdout+probe.Stderr), "conflict")
	probe.RequireContains(probe.Stdout, "CLAUDE.md")
	requireLinkFile(t, f, ".bench/gate.sh")
	requireLiteralCount(t, f, "AGENTS.md", "<!-- bench:start -->", 1, "non-conflicting AGENTS.md was not written when CLAUDE.md was a FIFO")
}

// testSetupDirectoryAtAgentsIsConflict pins C2: a directory sitting where AGENTS.md
// belongs must be preserved as a conflict, the same posture as a FIFO, rather than
// falling through to the raw os.ReadFile EISDIR error the special-file guard used to
// leave uncaught (isSpecialFile previously excluded directories).
func testSetupDirectoryAtAgentsIsConflict(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("go.mod", "module example.com/fixture\n\ngo 1.22\n")
	contract.Mkdir(t, filepath.Join(f.Root, "AGENTS.md"))
	contract.WriteFileAbs(t, filepath.Join(f.Root, "AGENTS.md", "keep.txt"), "project directory content\n")

	probe := f.Bench("setup", "--yes")

	probe.RequireExit(3)
	probe.RequireContains(strings.ToLower(probe.Stdout+probe.Stderr), "conflict")
	probe.RequireContains(probe.Stdout, "AGENTS.md")
	if info, err := os.Stat(filepath.Join(f.Root, "AGENTS.md")); err != nil || !info.IsDir() {
		t.Fatalf("bench setup did not preserve a directory at AGENTS.md as a conflict: %v", err)
	}
	requireLinkFile(t, f, ".bench/gate.sh")
}

// testSetupDirectoryAtClaudeIsConflict is the CLAUDE.md mirror of the AGENTS.md case.
func testSetupDirectoryAtClaudeIsConflict(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("go.mod", "module example.com/fixture\n\ngo 1.22\n")
	contract.Mkdir(t, filepath.Join(f.Root, "CLAUDE.md"))
	contract.WriteFileAbs(t, filepath.Join(f.Root, "CLAUDE.md", "keep.txt"), "project directory content\n")

	probe := f.Bench("setup", "--yes")

	probe.RequireExit(3)
	probe.RequireContains(strings.ToLower(probe.Stdout+probe.Stderr), "conflict")
	probe.RequireContains(probe.Stdout, "CLAUDE.md")
	if info, err := os.Stat(filepath.Join(f.Root, "CLAUDE.md")); err != nil || !info.IsDir() {
		t.Fatalf("bench setup did not preserve a directory at CLAUDE.md as a conflict: %v", err)
	}
	requireLinkFile(t, f, ".bench/gate.sh")
}
