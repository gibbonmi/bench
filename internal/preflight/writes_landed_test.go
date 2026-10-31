package preflight

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
	"github.com/gibbonmi/bench/internal/reviewrecord"
)

// seedLanded builds the conformant fixture and runs history on main. It then commits a
// ticket on main that also writes each plain entry, restarts the feature branch at main
// with the fenced change, and returns the slug.
func seedLanded(t *testing.T, history func(), entries ...string) string {
	t.Helper()
	root, slug := preflighttest.SeedConformant(t)
	preflighttest.RunGit(t, "checkout", "-q", "main")
	history()
	writes := append(preflighttest.FenceWrites(preflighttest.ConformantFence), entries...)
	preflighttest.CommitChargeSource(t, root, slug, map[string]string{
		"specs/" + slug + "/tickets/one.md": preflighttest.WritesTicketDoc("One", writes, "PF1", "PF2"),
	})
	return slug
}

// addPath commits path on the current branch and returns the commit.
func addPath(t *testing.T, path string) string {
	t.Helper()
	preflighttest.MustWriteFile(t, path, "landed work\n")
	return preflighttest.CommitAll(t, "add "+path)
}

// deletePath commits the deletion of path on the current branch and returns the commit.
func deletePath(t *testing.T, path string) string {
	t.Helper()
	preflighttest.RunGit(t, "rm", "-q", path)
	return preflighttest.CommitAll(t, "delete "+path)
}

// landedRow and resolveRow are the rendered red rows of writes-landed and writes-resolve
// over the named entries.
func landedRow(entries ...string) string {
	return `  writes-landed,red,"Writes: entry names a path that a landed commit deleted: ` + strings.Join(entries, ", ") + `",""`
}

func resolveRow(entries ...string) string {
	return `  writes-resolve,red,"Writes: entry names no tree path and carries no (new) marker: ` + strings.Join(entries, ", ") + `",""`
}

// requireRow fails unless out renders want as the row of check. An empty want requires
// that out renders no row of check, so the row is green or not applicable.
func requireRow(t *testing.T, out, check, want string) {
	t.Helper()
	if got, _ := rowOf(t, out, check); got != want {
		t.Errorf("%s row = %q, want %q:\n%s", check, got, want, out)
	}
}

// TestWritesLandedBareEntry checks that a bare build preflight with no recorded chunk names
// the plain entry that a default-branch commit deleted in writes-landed, with the full
// deleting commit, and keeps writes-resolve green.
func TestWritesLandedBareEntry(t *testing.T) {
	var deletion string
	slug := seedLanded(t, func() {
		addPath(t, "docs/gone.md")
		deletion = deletePath(t, "docs/gone.md")
	}, "docs/gone.md")

	out, code := Command([]string{"build", slug})
	if code != 1 {
		t.Fatalf("exit = %d, want 1:\n%s", code, out)
	}
	requireRow(t, out, "writes-landed", landedRow("one.md: docs/gone.md ("+deletion+")"))
	requireRow(t, out, "writes-resolve", "")
}

// TestWritesLandedShapes checks that a deletion that only a merge resolution made, a
// side-branch deletion that a merge brought in, and a rename each name the commit that the
// first-parent history of main answers. A never-tracked path, a path that the
// default-branch tip holds again, and a glob-shaped path beside a deleted sibling stay
// writes-resolve reds.
func TestWritesLandedShapes(t *testing.T) {
	var resolution, sideMerge, rename string
	slug := seedLanded(t, func() {
		addPath(t, "shapes/merge.md")
		preflighttest.RunGit(t, "checkout", "-q", "-b", "side")
		addPath(t, "shapes/side-work.md")
		preflighttest.RunGit(t, "checkout", "-q", "main")
		preflighttest.RunGit(t, "merge", "-q", "--no-ff", "--no-commit", "side")
		resolution = deletePath(t, "shapes/merge.md")

		addPath(t, "shapes/side.md")
		preflighttest.RunGit(t, "checkout", "-q", "-b", "side-deletion")
		deletePath(t, "shapes/side.md")
		preflighttest.RunGit(t, "checkout", "-q", "main")
		preflighttest.RunGit(t, "merge", "-q", "--no-ff", "-m", "merge the side deletion", "side-deletion")
		sideMerge = preflighttest.RunGit(t, "rev-parse", "HEAD")

		addPath(t, "shapes/old.md")
		preflighttest.RunGit(t, "mv", "shapes/old.md", "shapes/new.md")
		rename = preflighttest.CommitAll(t, "rename shapes/old.md")

		addPath(t, "shapes/restored.md")
		deletePath(t, "shapes/restored.md")
		addPath(t, "shapes/restored.md")

		addPath(t, "d/glXob x")
		deletePath(t, "d/glXob x")
	}, "shapes/merge.md", "shapes/side.md", "shapes/old.md", "shapes/never.md", "shapes/restored.md", "d/gl*ob x")
	// The build drops the restored path, so only the default-branch tip holds it.
	deletePath(t, "shapes/restored.md")

	out, _ := Command([]string{"build", slug})
	requireRow(t, out, "writes-landed", landedRow(
		"one.md: shapes/merge.md ("+resolution+")",
		"one.md: shapes/side.md ("+sideMerge+")",
		"one.md: shapes/old.md ("+rename+")",
	))
	requireRow(t, out, "writes-resolve", resolveRow("one.md: shapes/never.md", "one.md: shapes/restored.md", "one.md: d/gl*ob x"))
}

// TestWritesLandedPrecedence checks that each entry gets one answer: a mixed set reds both
// rows, a repeated landed entry is named once, an earlier-chunk deletion and a deletion
// inside the explicit --base range are green in both rows, and an unresolved default branch
// leaves a landed entry in writes-resolve.
func TestWritesLandedPrecedence(t *testing.T) {
	t.Run("mixed set", func(t *testing.T) {
		var deletion string
		slug := seedLanded(t, func() {
			addPath(t, "docs/gone.md")
			deletion = deletePath(t, "docs/gone.md")
		}, "docs/gone.md", "docs/typo.md")
		out, _ := Command([]string{"build", slug})
		requireRow(t, out, "writes-landed", landedRow("one.md: docs/gone.md ("+deletion+")"))
		requireRow(t, out, "writes-resolve", resolveRow("one.md: docs/typo.md"))
	})
	t.Run("repeated entry", func(t *testing.T) {
		var deletion string
		slug := seedLanded(t, func() {
			addPath(t, "docs/gone.md")
			deletion = deletePath(t, "docs/gone.md")
		}, "docs/gone.md", "docs/gone.md")
		out, _ := Command([]string{"build", slug})
		requireRow(t, out, "writes-landed", landedRow("one.md: docs/gone.md ("+deletion+")"))
		requireRow(t, out, "writes-resolve", "")
	})
	t.Run("earlier chunk deletion", func(t *testing.T) {
		slug := seedLanded(t, func() { addPath(t, "docs/chunk.md") }, "docs/chunk.md")
		base := preflighttest.RunGit(t, "rev-parse", "HEAD")
		deletion := deletePath(t, "docs/chunk.md")
		preflighttest.MustWriteFile(t, "reviews/"+slug+".md", recordFence(t, reviewrecord.Record{
			Version: 2, Spec: "specs/" + slug + "/spec.md", PlanDigest: "sha256:plan",
			Chunks: []reviewrecord.Chunk{{ID: "C1", Base: base, Tip: deletion, PlanDigest: "sha256:plan", SourceDigest: deletion}},
		}))
		preflighttest.CommitAll(t, "record the chunk")
		out, _ := Command([]string{"build", slug})
		requireRow(t, out, "writes-landed", "")
		requireRow(t, out, "writes-resolve", "")
	})
	t.Run("deletion inside the explicit range", func(t *testing.T) {
		var before string
		slug := seedLanded(t, func() {
			before = addPath(t, "docs/gone.md")
			deletePath(t, "docs/gone.md")
		}, "docs/gone.md")
		out, _ := Command([]string{"build", slug, "--base", before, "--source-tip", "HEAD"})
		requireRow(t, out, "writes-landed", "")
		requireRow(t, out, "writes-resolve", "")
	})
	t.Run("unresolved default branch", func(t *testing.T) {
		slug := seedLanded(t, func() {
			addPath(t, "docs/gone.md")
			deletePath(t, "docs/gone.md")
		}, "docs/gone.md")
		preflighttest.RunGit(t, "branch", "-m", "main", "trunk")
		out, _ := Command([]string{"build", slug})
		requireRow(t, out, "writes-landed", "")
		requireRow(t, out, "writes-resolve", resolveRow("one.md: docs/gone.md"))
	})
}

// TestWritesLandedReviewMode checks that the review preflight renders no writes-landed row,
// and it keeps a landed entry in writes-resolve, in the bare form and when the explicit
// --base comes after the deletion.
func TestWritesLandedReviewMode(t *testing.T) {
	slug := seedLanded(t, func() {
		addPath(t, "docs/gone.md")
		deletePath(t, "docs/gone.md")
	}, "docs/gone.md")
	for _, tc := range []struct {
		args []string
		rows int
	}{
		{[]string{"review", slug}, bareRowCount(modeReview)},
		{[]string{"review", slug, "--base", preflighttest.RunGit(t, "rev-parse", "main"), "--source-tip", "HEAD"}, bareRowCount(modeReview) + 1},
	} {
		out, _ := Command(tc.args)
		requireRow(t, out, "writes-resolve", resolveRow("one.md: docs/gone.md"))
		if strings.Contains(out, "writes-landed") || summaryCounts(t, out).total() != tc.rows {
			t.Errorf("%v renders a writes-landed row, or not %d rows:\n%s", tc.args, tc.rows, out)
		}
	}
}

// TestWritesLandedKeepsAuditManual checks that the audit form grades the writes-landed row
// of its explicit pair, and it requires the manual staleness pass while that row is red.
func TestWritesLandedKeepsAuditManual(t *testing.T) {
	var deletion string
	slug := seedLanded(t, func() {
		addPath(t, "docs/gone.md")
		deletion = deletePath(t, "docs/gone.md")
	}, "docs/gone.md")
	out, code := Command([]string{"build", slug, "--audit-spec", "--base", "main", "--source-tip", "HEAD"})
	if code != 1 {
		t.Fatalf("audit exit = %d, want 1:\n%s", code, out)
	}
	requireRow(t, out, "writes-landed", landedRow("one.md: docs/gone.md ("+deletion+")"))
	if manual := preflighttest.StalenessRow(t, out)["manual_required"]; manual != true {
		t.Errorf("manual_required = %v, want true while writes-landed is red:\n%s", manual, out)
	}
}
