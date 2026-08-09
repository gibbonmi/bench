package surface

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

// setup_route_test.go covers FT76 stories 1 (wrapper dispatch + convergence), 6
// (gitless stop), and 7 (the non-TTY invocation matrix). Story 2's preview-text
// assertions live in setup_plan_test.go; story 9's gate-table cases live in
// setup_gate_table_test.go.

func TestSetupRouteContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "bench setup converges a fresh git repo under --yes", testSetupConvergesFreshRepo)
	contract.RunParallel(t, "bench setup anchors at the repo root from a subdirectory", testSetupFromSubdirectory)
	contract.RunParallel(t, "bench setup in a gitless directory stops and writes nothing", testSetupGitlessStops)
	contract.RunParallel(t, "bench setup non-TTY matrix", testSetupNonTTYMatrix)
	contract.RunParallel(t, "bench setup seeds the profile through the link transaction, not after it", testSetupProfileWriteIsTransactional)
	contract.RunParallel(t, "bench setup flag edge cases", testSetupFlagEdgeCases)
}

func testSetupConvergesFreshRepo(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("go.mod", "module example.com/fixture\n\ngo 1.22\n")

	probe := f.BenchWrapper("setup", "--yes")

	probe.RequireExit(0)
	requireLinkFile(t, f, "AGENTS.md")
	requireLiteralCount(t, f, "AGENTS.md", "<!-- bench:start -->", 1, "setup did not create exactly one managed AGENTS.md start marker")
	requireLiteralCount(t, f, "AGENTS.md", "<!-- bench:end -->", 1, "setup did not create exactly one managed AGENTS.md end marker")
	requireFixtureFileContains(t, f, "CLAUDE.md", "@AGENTS.md", "setup CLAUDE.md is missing the AGENTS.md import line")
	requireFixtureFileContains(t, f, "CLAUDE.md", "@.bench/BENCH.md", "setup CLAUDE.md is missing the BENCH.md import line")
	requireLinkFile(t, f, ".bench/gate.sh")
	requireExecutable(t, filepath.Join(f.Root, ".bench", "gate.sh"), "setup did not write an executable .bench/gate.sh")
	requireLinkFile(t, f, ".bench/link-manifest.tsv")
	requireExecutable(t, filepath.Join(f.Root, ".bench", "bin", "bench.sh"), "setup did not install the repo-local launcher .bench/bin/bench.sh")
	if !hasProfileFile(f) {
		t.Fatal("setup did not write a projects/<name>.md profile")
	}
	// The profile is reviewer-owned judgment content, not a managed/converged asset:
	// it is seeded through the same FT84 transaction as everything else, but a manifest
	// row would make a later reviewer hand-edit read back as a modified-managed
	// conflict on the next run, so it must never appear in the manifest.
	manifest := f.ReadFile(".bench/link-manifest.tsv")
	for _, line := range strings.Split(manifest, "\n") {
		if strings.HasPrefix(line, "projects/") {
			t.Fatalf("link-manifest.tsv unexpectedly tracks the reviewer-owned profile: %q", line)
		}
	}
}

// testSetupProfileWriteIsTransactional proves the profile write is part of the single
// FT84 transaction, not a second write that lands after it: a regular file sitting
// where the profile's parent directory belongs makes the profile write fail, and that
// failure must roll back everything else with it rather than leaving AGENTS.md,
// CLAUDE.md, and .bench/gate.sh converged behind a non-zero exit.
func testSetupProfileWriteIsTransactional(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("go.mod", "module example.com/fixture\n\ngo 1.22\n")
	f.WriteFile("projects", "not a directory\n")

	probe := f.BenchWrapper("setup", "--yes")

	if probe.ExitCode == 0 {
		t.Fatalf("bench setup succeeded despite an unwritable profile path\nstdout:\n%s\nstderr:\n%s", probe.Stdout, probe.Stderr)
	}
	if f.Exists("AGENTS.md") || f.Exists(".bench") {
		t.Fatal("bench setup left a partially-converged tree when the profile write failed - the profile write is not transactional with the rest of the plan")
	}
}

func testSetupFlagEdgeCases(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("go.mod", "module example.com/fixture\n\ngo 1.22\n")

	both := f.BenchWrapper("setup", "--plan", "--yes")
	if both.ExitCode != 2 {
		t.Fatalf("bench setup --plan --yes exited %d, want 2\nstdout:\n%s\nstderr:\n%s", both.ExitCode, both.Stdout, both.Stderr)
	}
	both.RequireContains(both.Stderr, "--plan")
	both.RequireContains(both.Stderr, "--yes")

	unknown := f.BenchWrapper("setup", "--bogus")
	if unknown.ExitCode != 2 {
		t.Fatalf("bench setup --bogus exited %d, want 2\nstdout:\n%s\nstderr:\n%s", unknown.ExitCode, unknown.Stdout, unknown.Stderr)
	}
	unknown.RequireContains(unknown.Stderr, "usage")
}

func testSetupFromSubdirectory(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("go.mod", "module example.com/fixture\n\ngo 1.22\n")
	sub := filepath.Join(f.Root, "pkg", "nested")
	contract.Mkdir(t, sub)

	bench := filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh")
	probe := contract.RunAt(t, f, sub, nil, "bash", bench, "setup", "--yes")

	probe.RequireExit(0)
	requireLinkFile(t, f, ".bench/gate.sh")
	requireLinkFile(t, f, "AGENTS.md")
}

func testSetupGitlessStops(t *testing.T) {
	f := contract.NewFixture(t, contract.WithNoRepo())

	probe := f.BenchWrapper("setup", "--yes")

	if probe.ExitCode == 0 {
		t.Fatalf("bench setup succeeded outside a git repository\nstdout:\n%s\nstderr:\n%s", probe.Stdout, probe.Stderr)
	}
	probe.RequireContains(probe.Stderr, "git init")
	if f.Exists("AGENTS.md") || f.Exists(".bench") {
		t.Fatal("bench setup wrote files outside a git repository")
	}
}

func testSetupNonTTYMatrix(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("go.mod", "module example.com/fixture\n\ngo 1.22\n")

	// A default nil cmd.Stdin connects the child to /dev/null, a character device that
	// terminal.IsTerminal's ModeCharDevice heuristic misreads as a terminal (see
	// runtime_gate_test.go's identical gate-pin non-TTY contract) - supply real stdin
	// content so the child gets a pipe instead, the genuinely non-TTY case this row means
	// to exercise.
	bench := filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh")
	bare := contract.RunAtWithInput(t, f, f.Root, nil, "\n", "bash", bench, "setup")
	if bare.ExitCode == 0 {
		t.Fatalf("bare non-TTY bench setup succeeded\nstdout:\n%s\nstderr:\n%s", bare.Stdout, bare.Stderr)
	}
	bare.RequireContains(bare.Stderr, "--yes")
	bare.RequireContains(bare.Stderr, "--plan")
	if f.Exists("AGENTS.md") || f.Exists(".bench") {
		t.Fatal("bare non-TTY bench setup wrote files")
	}

	plan := f.BenchWrapper("setup", "--plan")
	plan.RequireExit(0)
	if f.Exists("AGENTS.md") || f.Exists(".bench") {
		t.Fatal("bench setup --plan wrote files")
	}

	ambiguous := contract.NewFixture(t)
	ambiguous.WriteFile("go.mod", "module example.com/fixture\n\ngo 1.22\n")
	ambiguous.WriteFile("package.json", `{"name":"fixture","scripts":{"test":"echo ok"}}`+"\n")
	ambiguousProbe := ambiguous.BenchWrapper("setup", "--yes")
	if ambiguousProbe.ExitCode == 0 {
		t.Fatalf("--yes over an ambiguous gate plan succeeded\nstdout:\n%s\nstderr:\n%s", ambiguousProbe.Stdout, ambiguousProbe.Stderr)
	}
	ambiguousProbe.RequireContains(strings.ToLower(ambiguousProbe.Stderr), "open question")
	if ambiguous.Exists("AGENTS.md") || ambiguous.Exists(".bench") {
		t.Fatal("--yes over an ambiguous gate plan wrote files")
	}

	conflicted := contract.NewFixture(t)
	conflicted.WriteFile("go.mod", "module example.com/fixture\n\ngo 1.22\n")
	conflicted.WriteFile(".agents/commands/bench-implement-spec.md", "project command\n")
	conflictedProbe := conflicted.BenchWrapper("setup", "--yes")
	conflictedProbe.RequireExit(3)
	requireFixtureFileContains(t, conflicted, ".agents/commands/bench-implement-spec.md", "project command", "--yes over a conflicted but ambiguity-free plan overwrote a project-owned asset")
	if !conflicted.Exists(".bench/gate.sh") {
		t.Fatal("--yes over a conflicted but ambiguity-free plan did not write the non-conflicting gate.sh")
	}
}

func hasProfileFile(f contract.Fixture) bool {
	matches, _ := filepath.Glob(filepath.Join(f.Root, "projects", "*.md"))
	return len(matches) > 0
}
