package axi

import (
	"github.com/gibbonmi/bench/internal/contract"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAXIWave2Contracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "AXI diff recorded-base contract", testAXIDiffRecordedBase)
	contract.RunParallel(t, "AXI diff fallback/shape contract", testAXIDiffFallbackShape)
	contract.RunParallel(t, "AXI diff error-posture contract", testAXIDiffErrorPosture)
	contract.RunParallel(t, "AXI diff --full contract", testAXIDiffFullContract)
	contract.RunParallel(t, "AXI diff git-failure propagation contract", testAXIDiffGitFailurePropagation)
	contract.RunParallel(t, "AXI diff control-byte posture contract", testAXIDiffControlBytePosture)
	contract.RunParallel(t, "AXI diff --commit range contract", testAXIDiffCommitRange)
	contract.RunParallel(t, "AXI diff --commit merge-commit contract", testAXIDiffCommitMerge)
	contract.RunParallel(t, "AXI diff --commit error-posture contract", testAXIDiffCommitErrorPosture)
	contract.RunParallel(t, "AXI coverage extraction contract", testAXICoverageExtraction)
	contract.RunParallel(t, "AXI coverage state/error contract", testAXICoverageStateError)
	contract.RunParallel(t, "AXI coverage --check validation contract", testAXICoverageCheckValidation)
	contract.RunParallel(t, "AXI maps/guards help contract", testAXIMapsGuardsHelp)
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

// testAXIDiffCommitRange pins the review-after-merge story-1/2 rows: `--commit
// HEAD` bounds base to the commit's own first parent and the files table to that
// commit's own change (not the whole branch history); `--full` composes with
// `--commit` in either argument order to append the log table and raw diff body;
// an --allow-empty commit renders files[0]/log[0] at exit 0 through the shared
// renderers; and a detached-HEAD checkout still renders `branch: (detached)`
// through the existing preamble path, since landed-commit review often runs there.
func testAXIDiffCommitRange(t *testing.T) {
	contract.NoteContractFailure(t, "AXI diff --commit range contract failed")
	f := contract.NewFixture(t)
	f.WriteFile("README.md", "r\n")
	f.CommitAll("c1")
	f.Git("branch", "-m", "main")
	c1 := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)

	f.WriteFile("f.txt", "line one\nline two\n")
	f.CommitAll("c2")
	c2 := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)

	// row 1: base is c2's first parent (c1), method names the resolved commit,
	// and the files table is c2's own change alone.
	out := f.Bench("diff", "--commit", "HEAD")
	out.RequireExit(0)
	requireOutputLine(t, out, "base: "+c1)
	requireOutputLine(t, out, "method: commit "+c2)
	requireOutputLine(t, out, "files[1]{status,path}:")
	requireOutputLine(t, out, "  A,f.txt")

	// row 2: --full --commit <sha> appends the one brought-in commit's log row
	// and the raw diff body with an unmangled @@ hunk marker.
	shortC2 := strings.TrimSpace(f.Git("rev-parse", "--short", c2).Stdout)
	full := f.Bench("diff", "--full", "--commit", c2)
	full.RequireExit(0)
	requireOutputLine(t, full, "log[1]{sha,subject}:")
	requireLogRow(t, full, shortC2, "c2")
	requireOutputLine(t, full, "diff_body:")
	full.RequireContains(full.Stdout, "@@ ")

	// the two flags compose in the other order too.
	full2 := f.Bench("diff", "--commit", c2, "--full")
	full2.RequireExit(0)
	requireOutputLine(t, full2, "log[1]{sha,subject}:")

	// edge inventory: an --allow-empty commit renders files[0] (no file change)
	// at exit 0 through the shared renderer already pinned for empty-since-base
	// — the log table still carries the empty commit's own row (log a..b always
	// includes b itself), and the diff body stays empty (no @@ marker).
	f.Git("-c", "user.email=bench@local", "-c", "user.name=bench",
		"commit", "--allow-empty", "-q", "-m", "empty")
	emptySha := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
	shortEmpty := strings.TrimSpace(f.Git("rev-parse", "--short", emptySha).Stdout)
	empty := f.Bench("diff", "--full", "--commit", emptySha)
	empty.RequireExit(0)
	requireOutputLine(t, empty, "files[0]{status,path}:")
	requireOutputLine(t, empty, "log[1]{sha,subject}:")
	requireLogRow(t, empty, shortEmpty, "empty")
	requireOutputLine(t, empty, "diff_body:")
	empty.RequireNotContains(empty.Stdout, "@@ ")

	// edge inventory: detached HEAD still renders `branch: (detached)` via the
	// existing preamble path.
	f.Git("switch", "-q", "--detach", c2)
	detached := f.Bench("diff", "--commit", c2)
	detached.RequireExit(0)
	requireOutputLine(t, detached, "branch: (detached)")
	requireOutputLine(t, detached, "base: "+c1)
}

// testAXIDiffCommitMerge pins the story-3 merge-commit row: a merge commit's
// `--commit` bundle shows the full merged-in change — the feature branch's file,
// the two commits (the feature commit and the merge commit itself) `A..B` (first
// parent .. merge) walks, and the raw body — proving first-parent semantics
// rather than a merge-base or second-parent view.
func testAXIDiffCommitMerge(t *testing.T) {
	contract.NoteContractFailure(t, "AXI diff --commit merge-commit contract failed")
	f := contract.NewFixture(t)
	f.WriteFile("README.md", "r\n")
	f.CommitAll("c1")
	f.Git("branch", "-m", "main")

	f.Git("switch", "-qc", "feature")
	f.WriteFile("feature.txt", "feature work\n")
	f.CommitAll("feature work")
	featureTip := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
	shortFeature := strings.TrimSpace(f.Git("rev-parse", "--short", featureTip).Stdout)

	f.Git("switch", "-q", "main")
	f.WriteFile("mainline.txt", "mainline work\n")
	f.CommitAll("mainline work")

	f.Git("-c", "user.email=bench@local", "-c", "user.name=bench",
		"merge", "--no-ff", "-q", "-m", "merge feature", "feature")
	mergeSha := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
	shortMerge := strings.TrimSpace(f.Git("rev-parse", "--short", mergeSha).Stdout)
	mergeParent1 := strings.TrimSpace(f.Git("rev-parse", mergeSha+"^").Stdout)

	out := f.Bench("diff", "--full", "--commit", mergeSha)
	out.RequireExit(0)
	requireOutputLine(t, out, "base: "+mergeParent1)
	requireOutputLine(t, out, "method: commit "+mergeSha)
	// first-parent pinned: the merge's bundle shows the feature branch's file —
	// the full merged-in change — not an empty (merge-commit-only) diff.
	requireOutputLine(t, out, "files[1]{status,path}:")
	requireOutputLine(t, out, "  A,feature.txt")
	// log = commits reachable from the merge but not from its first parent: the
	// feature commit and the merge commit itself.
	requireOutputLine(t, out, "log[2]{sha,subject}:")
	requireLogRow(t, out, shortFeature, "feature work")
	requireLogRow(t, out, shortMerge, "merge feature")
	out.RequireContains(out.Stdout, "@@ ")
}

// testAXIDiffCommitErrorPosture pins the story-3 edge postures: a root commit's
// missing parent and an unresolvable sha each exit 1 with a structured error (no
// leaked git stderr), and a bare or repeated --commit exits 2 usage rather than
// silently falling back to branch-relative resolution.
func testAXIDiffCommitErrorPosture(t *testing.T) {
	contract.NoteContractFailure(t, "AXI diff --commit error-posture contract failed")
	f := contract.NewFixture(t)
	f.WriteFile("README.md", "r\n")
	f.CommitAll("root commit")
	f.Git("branch", "-m", "main")
	root := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)

	f.WriteFile("f.txt", "f\n")
	f.CommitAll("c2")

	// a root commit has no parent to diff against — a structured, named error.
	out := f.Bench("diff", "--commit", root)
	out.RequireExit(1)
	requireOutputPrefix(t, out, "error: ")
	out.RequireContains(out.Stdout, "no parent")

	// an unresolvable sha gets the same exit-1 structured-error posture.
	bogus := f.Bench("diff", "--commit", "not-a-real-sha")
	bogus.RequireExit(1)
	requireOutputPrefix(t, bogus, "error: ")

	// --commit missing its value — exit 2 usage.
	missing := f.Bench("diff", "--commit")
	missing.RequireExit(2)
	missing.RequireContains(strings.ToLower(missing.Stdout), "usage")

	// --commit given twice — exit 2 usage.
	repeated := f.Bench("diff", "--commit", "HEAD", "--commit", "HEAD")
	repeated.RequireExit(2)
	repeated.RequireContains(strings.ToLower(repeated.Stdout), "usage")
}

// testAXIMapsGuardsHelp pins row 8: the real, agent-invocable flags are documented
// in their own command's help, not findable only in Go source.
func testAXIMapsGuardsHelp(t *testing.T) {
	contract.NoteContractFailure(t, "AXI maps/guards help contract failed")
	f := contract.NewFixture(t)

	maps := f.Bench("maps", "-h")
	maps.RequireExit(0)
	requireContainsFold(t, maps.Stdout, "--count")

	guards := f.Bench("guards", "-h")
	guards.RequireExit(0)
	requireContainsFold(t, guards.Stdout, "--brief")
}

// requireLogRow asserts a `log` table row for sha/renderedSubject, tolerating the
// TOON encoder's own choice to quote the sha cell (a short sha can start with a
// digit and read as numeric-ish, e.g. "05af9b0") — that quoting is independent of
// anything this seam controls, so the row match accepts either form.
func requireLogRow(t *testing.T, probe contract.Probe, sha, renderedSubject string) {
	t.Helper()
	plain := "  " + sha + "," + renderedSubject
	quoted := "  \"" + sha + "\"," + renderedSubject
	for _, got := range strings.Split(strings.TrimSuffix(probe.Stdout, "\n"), "\n") {
		if got == plain || got == quoted {
			return
		}
	}
	t.Fatalf("missing log row for sha %q\nstdout:\n%s\nstderr:\n%s", sha, probe.Stdout, probe.Stderr)
}

func testAXICoverageExtraction(t *testing.T) {
	contract.NoteContractFailure(t, "AXI coverage extraction contract failed")
	f := contract.NewFixture(t)
	spec := `# t

## User stories
1. As a, I want b, so c.
2. As d, I want e, so f.
3. As g, I want h, so i.

## Testing decisions

### Acceptance coverage map
| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 2–3 | does x \| y | cli seam | cmd fails, loudly | catches z |
| edge of 1 | edge case | gate | already covered | catches w |

## Out of scope
`
	f.WriteFile("specs/t.md", spec)

	out := f.Bench("coverage", "specs/t.md")

	out.RequireExit(0)
	requireOutputLine(t, out, "spec: specs/t.md")
	requireOutputLine(t, out, "state: mapped")
	requireOutputLine(t, out, "rows[2]{story,seam,red_signal}:")
	requireOutputLine(t, out, `  2–3,cli seam,"cmd fails, loudly"`)
	requireOutputLine(t, out, "  edge of 1,gate,already covered")

	f.WriteFile("specs/t.md", strings.ReplaceAll(spec, "\n", "\r\n"))
	out = f.Bench("coverage", "specs/t.md")

	out.RequireExit(0)
	requireOutputLine(t, out, "state: mapped")
	requireOutputLine(t, out, "rows[2]{story,seam,red_signal}:")
}

func testAXICoverageStateError(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("specs/h.md", "# h\n\n<!-- coverage-map: historical -->\n\n### Acceptance coverage map\n| a |\n")

	out := f.Bench("coverage", "specs/h.md")

	out.RequireExit(0)
	requireOutputLine(t, out, "state: historical")
	requireOutputLine(t, out, "rows[0]{story,seam,red_signal}:")

	f.WriteFile("specs/n.md", "# n\nno map here\n")
	out = f.Bench("coverage", "specs/n.md")

	out.RequireExit(0)
	requireOutputLine(t, out, "state: no-map")

	out = f.Bench("coverage")

	out.RequireExit(2)
	out.RequireContains(strings.ToLower(out.Stdout), "usage")

	out = f.Bench("coverage", "--bogus", "specs/n.md")

	out.RequireExit(2)

	out = f.Bench("coverage", "specs/absent.md")

	out.RequireExit(1)
	requireOutputPrefix(t, out, "error: ")
}

func testAXICoverageCheckValidation(t *testing.T) {
	f := contract.NewFixture(t)
	stories := "## User stories\n1. As a, I want b, so c.\n2. As d, I want e, so f.\n3. As g, I want h, so i.\n"
	header := "| story | behavior | seam | red signal | why it catches the failure |\n|---|---|---|---|---|"
	f.WriteFile("specs/v.md", "# v\n\n"+stories+"\n### Acceptance coverage map\n"+header+"\n| 1, 2–3 | b | s | r | w |\n")

	out := f.Bench("coverage", "--check", "specs/v.md")

	// row 2: a valid map's --check gets a definitive pass line, not silence —
	// stdout was empty before this story, indistinguishable from a check that
	// silently produced nothing.
	out.RequireExit(0)
	requireOutputLine(t, out, "ok: coverage map valid — 1 row(s)")

	f.WriteFile("specs/h.md", "# h\n<!-- coverage-map: historical -->\n### Acceptance coverage map\n|bad|\n")
	out = f.Bench("coverage", "--check", "specs/h.md")
	out.RequireExit(0)
	requireNoOutput(t, out)

	f.WriteFile("specs/n.md", "# n\nprose only\n")
	out = f.Bench("coverage", "--check", "specs/n.md")
	out.RequireExit(0)
	requireNoOutput(t, out)

	cases := []struct {
		name string
		path string
		body string
		want string
	}{
		{
			name: "missing canonical header",
			path: "specs/b1.md",
			body: "# b\n\n" + stories + "\n### Acceptance coverage map\n| a | b |\n|---|---|\n| 1 | x |\n",
			want: "coverage map missing the canonical header",
		},
		{
			name: "no data rows",
			path: "specs/b2.md",
			body: "# b\n\n" + stories + "\n### Acceptance coverage map\n" + header + "\n",
			want: "coverage map has no data rows",
		},
		{
			name: "wrong cell count",
			path: "specs/b3.md",
			body: "# b\n\n" + stories + "\n### Acceptance coverage map\n" + header + "\n| 1 | b | s | r |\n",
			want: "coverage map row 1 has 4 cells (want 5)",
		},
		{
			name: "empty cell",
			path: "specs/b4.md",
			body: "# b\n\n" + stories + "\n### Acceptance coverage map\n" + header + "\n| 1 | b |  | r | w |\n",
			want: "coverage map row 1 has an empty 'seam' cell",
		},
		{
			name: "story reference beyond numbering",
			path: "specs/b5.md",
			body: "# b\n\n" + stories + "\n### Acceptance coverage map\n" + header + "\n| 9 | b | s | r | w |\n",
			want: "references story 9 but the spec numbers only 3",
		},
		{
			name: "unrecognized story reference",
			path: "specs/b6.md",
			body: "# b\n\n" + stories + "\n### Acceptance coverage map\n" + header + "\n| x | b | s | r | w |\n",
			want: "has an unrecognized story reference 'x'",
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			f := contract.NewFixture(t)
			f.WriteFile(c.path, c.body)

			out := f.Bench("coverage", "--check", c.path)

			out.RequireExit(1)
			out.RequireContains(out.Stdout, c.want)
		})
	}

	// row 3: a violation line renders through the one canonical `toon.Errorf`
	// shape (`error: <kind> — <hint>`), not a hand-rolled line of its own —
	// pinned byte-for-byte rather than by substring.
	f.WriteFile("specs/b1.md", "# b\n\n"+stories+"\n### Acceptance coverage map\n| a | b |\n|---|---|\n| 1 | x |\n")
	out = f.Bench("coverage", "--check", "specs/b1.md")
	out.RequireExit(1)
	requireOutputLine(t, out, "error: specs/b1.md coverage map missing the canonical header — fix the map or mark it <!-- coverage-map: historical -->")
}

func requireOutputLine(t *testing.T, probe contract.Probe, line string) {
	t.Helper()
	for _, got := range strings.Split(strings.TrimSuffix(probe.Stdout, "\n"), "\n") {
		if got == line {
			return
		}
	}
	t.Fatalf("missing output line %q\nstdout:\n%s\nstderr:\n%s", line, probe.Stdout, probe.Stderr)
}

func requireOutputPrefix(t *testing.T, probe contract.Probe, prefix string) {
	t.Helper()
	for _, got := range strings.Split(strings.TrimSuffix(probe.Stdout, "\n"), "\n") {
		if strings.HasPrefix(got, prefix) {
			return
		}
	}
	t.Fatalf("missing output line with prefix %q\nstdout:\n%s\nstderr:\n%s", prefix, probe.Stdout, probe.Stderr)
}

func requireNoOutput(t *testing.T, probe contract.Probe) {
	t.Helper()
	if probe.Stdout != "" || probe.Stderr != "" {
		t.Fatalf("expected no output\nstdout:\n%s\nstderr:\n%s", probe.Stdout, probe.Stderr)
	}
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}
