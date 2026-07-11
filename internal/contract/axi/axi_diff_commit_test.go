package axi

import (
	"github.com/gibbonmi/bench/internal/contract"
	"strings"
	"testing"
)

func TestAXIDiffCommitContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "AXI diff --commit range contract", testAXIDiffCommitRange)
	contract.RunParallel(t, "AXI diff --commit merge-commit contract", testAXIDiffCommitMerge)
	contract.RunParallel(t, "AXI diff --commit error-posture contract", testAXIDiffCommitErrorPosture)
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
	// and the raw diff body with an unmangled @@ hunk marker. Unrelated worktree
	// state stays outside this exact landed-commit view.
	f.WriteFile("README.md", "dirty\n")
	shortC2 := strings.TrimSpace(f.Git("rev-parse", "--short", c2).Stdout)
	full := f.Bench("diff", "--full", "--commit", c2)
	full.RequireExit(0)
	requireOutputLine(t, full, "log[1]{sha,subject}:")
	requireLogRow(t, full, shortC2, "c2")
	requireOutputLine(t, full, "diff_body:")
	full.RequireContains(full.Stdout, "@@ ")
	full.RequireNotContains(full.Stdout, "+dirty")
	f.WriteFile("README.md", "r\n")

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
