package surface

import (
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

// setup_convergence_matrix_test.go covers FT76 story 5: instruction files converge
// across the file-state matrix (absent/empty/present-without-block or -imports,
// missing trailing newline, unclosed fence) with project content preserved in every
// non-conflict cell.

func TestSetupConvergenceMatrixContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "AGENTS.md absent converges", testSetupAgentsAbsent)
	contract.RunParallel(t, "AGENTS.md empty converges", testSetupAgentsEmpty)
	contract.RunParallel(t, "AGENTS.md present without block converges, project text preserved", testSetupAgentsNoBlock)
	contract.RunParallel(t, "AGENTS.md missing trailing newline converges, project text preserved", testSetupAgentsNoTrailingNewline)
	contract.RunParallel(t, "AGENTS.md unclosed fence around markers is a conflict", testSetupAgentsUnclosedFence)
	contract.RunParallel(t, "CLAUDE.md absent converges", testSetupClaudeAbsent)
	contract.RunParallel(t, "CLAUDE.md empty is left project-owned, untouched", testSetupClaudeEmpty)
	contract.RunParallel(t, "CLAUDE.md present without imports is left project-owned", testSetupClaudeNoImports)
}

func setupMatrixFixture(t *testing.T) contract.Fixture {
	t.Helper()
	f := contract.NewFixture(t)
	f.WriteFile("go.mod", "module example.com/fixture\n\ngo 1.22\n")
	return f
}

func testSetupAgentsAbsent(t *testing.T) {
	f := setupMatrixFixture(t)
	probe := f.Bench("setup", "--yes")
	probe.RequireExit(0)
	requireLiteralCount(t, f, "AGENTS.md", "<!-- bench:start -->", 1, "absent AGENTS.md did not converge")
}

func testSetupAgentsEmpty(t *testing.T) {
	f := setupMatrixFixture(t)
	f.WriteFile("AGENTS.md", "")
	probe := f.Bench("setup", "--yes")
	probe.RequireExit(0)
	requireLiteralCount(t, f, "AGENTS.md", "<!-- bench:start -->", 1, "empty AGENTS.md did not converge")
}

func testSetupAgentsNoBlock(t *testing.T) {
	f := setupMatrixFixture(t)
	f.WriteFile("AGENTS.md", "# Project rules\n\nKEEP-ME project text.\n")
	probe := f.Bench("setup", "--yes")
	probe.RequireExit(0)
	requireLiteralCount(t, f, "AGENTS.md", "<!-- bench:start -->", 1, "AGENTS.md without a block did not converge")
	requireFixtureFileContains(t, f, "AGENTS.md", "KEEP-ME project text.", "project text lost converging AGENTS.md")
}

func testSetupAgentsNoTrailingNewline(t *testing.T) {
	f := setupMatrixFixture(t)
	f.WriteFile("AGENTS.md", "# Project rules\n\nKEEP-ME no trailing newline")
	probe := f.Bench("setup", "--yes")
	probe.RequireExit(0)
	requireLiteralCount(t, f, "AGENTS.md", "<!-- bench:start -->", 1, "AGENTS.md lacking a trailing newline did not converge")
	requireFixtureFileContains(t, f, "AGENTS.md", "KEEP-ME no trailing newline", "project text lost converging a no-trailing-newline AGENTS.md")
}

func testSetupAgentsUnclosedFence(t *testing.T) {
	f := setupMatrixFixture(t)
	f.WriteFile("AGENTS.md", "# Project rules\n\nBroken docs with an unclosed fence:\n\n```\n<!-- bench:start -->\n<!-- bench:end -->\n\nKEEP-ME text after the unclosed fence.\n")
	probe := f.Bench("setup", "--yes")
	if probe.ExitCode == 0 {
		t.Fatalf("setup succeeded despite an unclosed fence around Bench markers\nstdout:\n%s\nstderr:\n%s", probe.Stdout, probe.Stderr)
	}
	requireFixtureFileContains(t, f, "AGENTS.md", "KEEP-ME text after the unclosed fence.", "unclosed-fence failure rewrote project text")
	requireLiteralCount(t, f, "AGENTS.md", "## Bench", 0, "unclosed-fence run still installed a managed block")
}

func testSetupClaudeAbsent(t *testing.T) {
	f := setupMatrixFixture(t)
	probe := f.Bench("setup", "--yes")
	probe.RequireExit(0)
	requireFixtureFileContains(t, f, "CLAUDE.md", "@AGENTS.md", "absent CLAUDE.md did not converge")
	requireFixtureFileContains(t, f, "CLAUDE.md", "@.bench/BENCH.md", "absent CLAUDE.md did not converge")
}

// testSetupClaudeEmpty pins the same rule bench link/unlink already enforce (see
// unlink_test.go's pre-existing-CLAUDE.md matrix): a pre-existing CLAUDE.md, even an
// empty one, is a path the user already claimed. Setup never writes into it - an
// empty file is not "no project content", it is content the user chose to leave
// blank, and only an absent CLAUDE.md is safe to create from scratch. Story 10's exit
// contract then makes the run partial: doctor's CLAUDE.md row is the honest red
// signal that the import lines still need adding by hand.
func testSetupClaudeEmpty(t *testing.T) {
	f := setupMatrixFixture(t)
	f.WriteFile("CLAUDE.md", "")
	probe := f.Bench("setup", "--yes")
	probe.RequireExit(3)
	probe.RequireContains(probe.Stdout, "red: CLAUDE.md")
	if got := f.ReadFile("CLAUDE.md"); got != "" {
		t.Fatalf("setup wrote into a pre-existing empty CLAUDE.md: %q", got)
	}
}

// testSetupClaudeNoImports pins the existing, deliberate link-lifecycle rule this
// slice reuses verbatim: a CLAUDE.md carrying unrelated project content is never
// rewritten to inject the import lines (see the "relink injected an import into a
// project-owned CLAUDE.md" regression guard in link_test.go) — project content stays
// exactly as the user left it, and doctor's CLAUDE.md row is the honest signal that
// it still needs the imports added by hand.
func testSetupClaudeNoImports(t *testing.T) {
	f := setupMatrixFixture(t)
	f.WriteFile("CLAUDE.md", "# Custom\n\nproject-owned claude config\n")
	probe := f.Bench("setup", "--yes")
	probe.RequireExit(3)
	probe.RequireContains(probe.Stdout, "red: CLAUDE.md")
	requireFixtureFileContains(t, f, "CLAUDE.md", "project-owned claude config", "setup rewrote a project-owned CLAUDE.md")
	requireFixtureFileNotContains(t, f, "CLAUDE.md", "@.bench/BENCH.md", "setup injected an import into a project-owned CLAUDE.md")
}
