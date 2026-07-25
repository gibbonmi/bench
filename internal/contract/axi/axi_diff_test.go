package axi

import (
	"github.com/gibbonmi/bench/internal/contract"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAXIDiffContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "AXI diff recorded-base contract", testAXIDiffRecordedBase)
	contract.RunParallel(t, "AXI diff fallback/shape contract", testAXIDiffFallbackShape)
	contract.RunParallel(t, "AXI diff error-posture contract", testAXIDiffErrorPosture)
	contract.RunParallel(t, "AXI diff --full contract", testAXIDiffFullContract)
	contract.RunParallel(t, "AXI diff git-failure propagation contract", testAXIDiffGitFailurePropagation)
	contract.RunParallel(t, "AXI diff control-byte posture contract", testAXIDiffControlBytePosture)
}

func testAXIDiffRecordedBase(t *testing.T) {
	contract.NoteContractFailure(t, "AXI diff recorded-base contract failed")
	f := contract.NewFixture(t)
	f.WriteFile("README.md", "r\n")
	f.CommitAll("c1")
	f.Git("branch", "-m", "main")
	f.Git("switch", "-qc", "feature")
	f.WriteFile("f.txt", "f\n")
	f.CommitAll("c2")
	c2 := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
	f.Git("switch", "-qc", "bench/shift-stacked")
	f.Git("config", "branch.bench/shift-stacked.benchBase", c2)
	f.WriteFile("work.txt", "w\n")
	f.CommitAll("c3")

	out := f.Bench("diff")

	out.RequireExit(0)
	requireOutputLine(t, out, "base: "+c2)
	requireOutputLine(t, out, "method: recorded")
	requireOutputLine(t, out, "files[1]{status,path}:")
	requireOutputLine(t, out, "  A,work.txt")

	f.Git("switch", "-q", "main")
	f.Git("switch", "-qc", "other")
	f.WriteFile("o.txt", "o\n")
	f.CommitAll("o1")
	o1 := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
	f.Git("switch", "-q", "bench/shift-stacked")
	f.Git("config", "branch.bench/shift-stacked.benchBase", o1)
	mergeBase := strings.TrimSpace(f.Git("merge-base", "main", "HEAD").Stdout)

	out = f.Bench("diff")

	out.RequireExit(0)
	requireOutputLine(t, out, "method: merge-base (recorded sha not an ancestor)")
	requireOutputLine(t, out, "base: "+mergeBase)
}

func testAXIDiffFallbackShape(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("README.md", "r\n")
	f.CommitAll("c1")
	f.Git("branch", "-m", "main")
	c1 := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
	f.Git("switch", "-qc", "feature")
	f.WriteFile("sub/deeper/.keep", "")
	f.WriteFile("a b.txt", "f\n")
	f.WriteFile("café.txt", "u\n")
	f.WriteFile("a\"q.txt", "q\n")
	f.CommitAll("c2")

	out := f.Run("bash", "-c", "cd sub/deeper && "+shellQuote(contract.SubjectRoot(t)+"/bin/bench.sh")+" diff")

	out.RequireExit(0)
	requireOutputLine(t, out, "base: "+c1)
	requireOutputLine(t, out, "method: merge-base")
	requireOutputLine(t, out, "  A,a b.txt")
	requireOutputLine(t, out, "  A,café.txt")
	requireOutputLine(t, out, `  A,"a\"q.txt"`)

	f.Git("config", "branch.feature.benchBase", "0123456789abcdef0123456789abcdef01234567")
	out = f.Bench("diff")

	out.RequireExit(0)
	requireOutputLine(t, out, "base: "+c1)
	requireOutputLine(t, out, "method: merge-base (recorded sha unreachable)")

	f.Git("switch", "-q", "main")
	f.Git("switch", "-qc", "idle")
	out = f.Bench("diff")

	out.RequireExit(0)
	requireOutputLine(t, out, "files[0]{status,path}:")
}

func testAXIDiffErrorPosture(t *testing.T) {
	f := contract.NewFixture(t, contract.WithNoRepo())

	out := f.Bench("diff")

	out.RequireExit(1)
	requireOutputPrefix(t, out, "error: ")

	f.Git("init", "-q")
	f.WriteFile("README.md", "r\n")
	f.CommitAll("c1")
	f.Git("branch", "-m", "trunk")
	// Two local branches with no origin/HEAD is what leaves the default genuinely
	// unresolvable: a lone local branch is itself the resolved default.
	f.Git("branch", "spare")

	out = f.Bench("diff")

	out.RequireExit(1)
	out.RequireContains(out.Stdout, "error: cannot resolve a review base")

	out = f.Bench("diff", "bogusarg")

	out.RequireExit(2)
	out.RequireContains(strings.ToLower(out.Stdout), "usage")

	// row 4b: --full does not change the error postures — base-unresolvable stays
	// exit 1, and an unrecognized arg (with or without --full) stays exit 2.
	out = f.Bench("diff", "--full")

	out.RequireExit(1)
	out.RequireContains(out.Stdout, "error: cannot resolve a review base")

	out = f.Bench("diff", "--full", "bogusarg")

	out.RequireExit(2)
	out.RequireContains(strings.ToLower(out.Stdout), "usage")
	// row 5: attribution names the actual offender ("bogusarg"), not the
	// recognized leading flag ("--full") that parsed fine.
	out.RequireContains(out.Stdout, "unknown argument: bogusarg")
}

// testAXIDiffFullContract drives `bench diff --full` end to end: an empty-since-base
// fixture (row 4a), a commit since base rendering the log table and raw diff body
// (rows 1/3a), a comma-and-quote commit subject escaped exactly once in the log table
// (row 3b), and a bare-`diff` regression guard that the new sections never leak into
// the unflagged path (row 2 — already covered, not TDD-able, guard only).
func testAXIDiffFullContract(t *testing.T) {
	contract.NoteContractFailure(t, "AXI diff --full contract failed")
	f := contract.NewFixture(t)
	f.WriteFile("README.md", "r\n")
	f.CommitAll("c1")
	f.Git("branch", "-m", "main")
	c1 := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
	f.Git("switch", "-qc", "feature")

	// row 2: bare `bench diff` never carries the --full sections. Not TDD-able
	// (bare diff never emitted these), so this is a regression guard, not a
	// red-first assertion.
	bare := f.Bench("diff")
	bare.RequireExit(0)
	bare.RequireNotContains(bare.Stdout, "log[")
	bare.RequireNotContains(bare.Stdout, "@@ ")

	// row 4a: base == HEAD (feature just branched off c1, no commits since).
	empty := f.Bench("diff", "--full")
	empty.RequireExit(0)
	requireOutputLine(t, empty, "base: "+c1)
	requireOutputLine(t, empty, "log[0]{sha,subject}:")
	requireOutputLine(t, empty, "diff_body:")
	empty.RequireNotContains(empty.Stdout, "@@ ")

	// A tracked worktree edit belongs to the branch-relative review picture even
	// when HEAD has no commits since base. The log remains commit-only.
	f.WriteFile("README.md", "dirty\n")
	dirty := f.Bench("diff", "--full")
	dirty.RequireExit(0)
	requireOutputLine(t, dirty, "files[1]{status,path}:")
	requireOutputLine(t, dirty, "  M,README.md")
	requireOutputLine(t, dirty, "log[0]{sha,subject}:")
	dirty.RequireContains(dirty.Stdout, "+dirty")
	f.WriteFile("README.md", "r\n")

	// A staged addition is part of the same branch-relative files/body picture.
	f.WriteFile("staged.txt", "staged\n")
	f.Git("add", "staged.txt")
	staged := f.Bench("diff", "--full")
	staged.RequireExit(0)
	requireOutputLine(t, staged, "files[1]{status,path}:")
	requireOutputLine(t, staged, "  A,staged.txt")
	requireOutputLine(t, staged, "log[0]{sha,subject}:")
	staged.RequireContains(staged.Stdout, "+staged")
	f.Git("reset", "--hard", "-q", "HEAD")

	// rows 1/3a: a commit since base emits the log table and the raw diff body,
	// appended last, with an unmangled line-start `@@ ` hunk marker.
	f.WriteFile("f.txt", "line one\nline two\n")
	f.CommitAll("c2")
	c2 := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
	shortC2 := strings.TrimSpace(f.Git("rev-parse", "--short", c2).Stdout)

	out := f.Bench("diff", "--full")
	out.RequireExit(0)
	requireOutputLine(t, out, "log[1]{sha,subject}:")
	requireLogRow(t, out, shortC2, "c2")
	requireOutputLine(t, out, "diff_body:")
	requireOutputPrefix(t, out, "@@ ")

	// row 3b: a commit subject with a comma and a quote is TOON-escaped exactly
	// once in the log table (mirrors testAXITOONFieldEscaping). The sha itself is
	// asserted via requireLogRow, which tolerates the TOON encoder quoting a
	// leading-digit sha (e.g. "05af9b0" reads as numeric-ish) — that quoting is
	// unrelated to what this row tests.
	f.WriteFile("g.txt", "g\n")
	f.CommitAll(`a, "b"`)
	c3 := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
	shortC3 := strings.TrimSpace(f.Git("rev-parse", "--short", c3).Stdout)

	out = f.Bench("diff", "--full")
	out.RequireExit(0)
	requireOutputLine(t, out, "log[2]{sha,subject}:")
	requireLogRow(t, out, shortC3, `"a, \"b\""`)
}

// testAXIDiffGitFailurePropagation drives `bench diff --full` with a PATH-shimmed
// git that fails one post-resolution call at a time (row 6): the changed-files
// listing, the commit log, and the raw diff body. Base resolution runs against the
// same shim and must still succeed — only the named call fails — so an empty
// section at exit 0 (the pre-fix behavior) never masks a broken subprocess.
func testAXIDiffGitFailurePropagation(t *testing.T) {
	contract.NoteContractFailure(t, "AXI diff git-failure propagation contract failed")
	cases := []struct {
		name string
		call string
	}{
		{"changed-files listing", "files"},
		{"commit log", "log"},
		{"diff body", "body"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			f := contract.NewFixture(t)
			f.WriteFile("README.md", "r\n")
			f.CommitAll("c1")
			f.Git("branch", "-m", "main")

			stubBin := filepath.Join(f.Root, "stubbin")
			if err := os.MkdirAll(stubBin, 0o755); err != nil {
				t.Fatalf("create stubbin: %v", err)
			}
			writeGitFailureShim(t, stubBin)

			env := map[string]string{
				"PATH":          stubBin + string(os.PathListSeparator) + os.Getenv("PATH"),
				"FAIL_GIT_CALL": c.call,
			}
			out := f.BenchEnv(env, "diff", "--full")
			out.RequireExit(1)
			requireOutputPrefix(t, out, "error: ")
		})
	}
}

// testAXIDiffControlBytePosture pins row 7: a commit subject carrying a control
// byte (ESC) makes `--full` exit 1 with the unrepresentable-TOON-cell error instead
// of a mangled log row — the refusal already exists (toon.Table's own guard), so
// this probe passes on day one and is the regression guard from here on. The help
// assertion is the TDD-able half: it was red before the diff --full help sentence
// named the refusal.
func testAXIDiffControlBytePosture(t *testing.T) {
	contract.NoteContractFailure(t, "AXI diff control-byte posture contract failed")
	f := contract.NewFixture(t)
	f.WriteFile("README.md", "r\n")
	f.CommitAll("c1")
	f.Git("branch", "-m", "main")
	f.Git("switch", "-qc", "feature")
	f.WriteFile("f.txt", "f\n")
	f.CommitAll("subject with \x1b esc")

	out := f.Bench("diff", "--full")
	out.RequireExit(1)
	requireOutputPrefix(t, out, "error: unrepresentable TOON cell")

	help := f.Bench("diff", "-h")
	help.RequireExit(0)
	requireContainsFold(t, help.Stdout, "control byte")
}

// TestAXIDiffUnresolvableDefault pins the fail-closed posture: with no default branch to
// compute a review base from, diff names that reality and the escape hatch rather than
// diffing against a branch the repository does not have.
func TestAXIDiffUnresolvableDefault(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)
	f.WriteFile("README.md", "r\n")
	f.CommitAll("c1")
	f.Git("branch", "-M", "master")
	f.Git("switch", "-qc", "feature")
	f.WriteFile("work.txt", "w\n")
	f.CommitAll("c2")

	out := f.Bench("diff")

	out.RequireExit(1)
	out.RequireContains(out.Stdout, "no resolvable default branch")
	out.RequireContains(out.Stdout, "git config branch.<name>.benchBase <sha>")
	requireNoAXILineMatching(t, out.Stdout, `'main'`)
}
