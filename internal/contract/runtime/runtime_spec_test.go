package runtime

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestRuntimeSpecImplementedContracts(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	contract.RunParallel(t, "flip writes the retirement form", testSpecFlipForm)
	contract.RunParallel(t, "bare slug anchors at repo root from a subdir", testSpecSlugAnchored)
	contract.RunParallel(t, "flip preserves every other byte", testSpecByteIdentity)
	contract.RunParallel(t, "error branches exit 1", testSpecErrorBranches)
	contract.RunParallel(t, "usage errors exit 2", testSpecUsageExitTwo)
}

func testSpecFlipForm(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("specs/feature.md", "# feature\n\nStatus: staged\n\n## body\n")

	p := f.Bench("spec", "implemented", "feature")
	p.RequireExit(0)
	p.RequireContains(p.Stdout, "spec implemented")
	if got := f.ReadFile("specs/feature.md"); got != "# feature\n\nStatus: implemented\n\n## body\n" {
		t.Fatalf("flipped content = %q", got)
	}
}

func testSpecSlugAnchored(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("specs/feature.md", "Status: staged\n")
	contract.Mkdir(t, filepath.Join(f.Root, "sub"))

	// A bare slug resolves specs/<slug>.md at the repo root even when run from a subdir.
	contract.RunAt(t, f, filepath.Join(f.Root, "sub"), nil, "bash", benchPath(t), "spec", "implemented", "feature").RequireExit(0)
	contract.RequireContains(t, f.ReadFile("specs/feature.md"), "Status: implemented")
}

func testSpecByteIdentity(t *testing.T) {
	f := contract.NewFixture(t)
	// No trailing newline, tab separator, and a later line that also contains "staged".
	f.WriteFile("specs/edge.md", "# edge\nStatus:\tstaged\nwork was staged elsewhere")

	f.Bench("spec", "implemented", "edge").RequireExit(0)
	if got := f.ReadFile("specs/edge.md"); got != "# edge\nStatus:\timplemented\nwork was staged elsewhere" {
		t.Fatalf("byte-identity violated: %q", got)
	}
}

func testSpecErrorBranches(t *testing.T) {
	for _, tc := range []struct {
		name    string
		content string // "" => don't create the file
		wantSub string
	}{
		{"not found", "", "not found"},
		{"no staged line", "# s\nStatus: draft\n", "no `Status: staged`"},
		{"already implemented", "# s\nStatus: implemented\n", "no `Status: staged`"},
		{"more than one", "Status: staged\nStatus: staged\n", "expected exactly one"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := contract.NewFixture(t)
			if tc.content != "" {
				f.WriteFile("specs/s.md", tc.content)
			}
			p := f.Bench("spec", "implemented", "s")
			p.RequireExit(1)
			p.RequireContains(p.Stdout+p.Stderr, tc.wantSub)
		})
	}
}

func testSpecUsageExitTwo(t *testing.T) {
	f := contract.NewFixture(t)
	for _, args := range [][]string{
		{"spec"},                       // no subcommand
		{"spec", "retire"},             // no argument
		{"spec", "unknowncmd"},         // unknown subcommand
		{"spec", "implemented"},        // no argument
		{"spec", "implemented", "--x"}, // unknown flag
	} {
		f.Bench(args...).RequireExit(2)
	}
}

func TestRuntimeSpecRetireContracts(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	contract.RunParallel(t, "retire deletes pickup then spec and prints duties", testSpecRetireMergedDeletesBoth)
	contract.RunParallel(t, "retire with no pickup retires the spec alone", testSpecRetireNoPickup)
	contract.RunParallel(t, "retire refuses a staged spec without deleting", testSpecRetireStagedRefuses)
	contract.RunParallel(t, "retire refuses an implemented spec not yet at HEAD", testSpecRetireNotAtHEAD)
	contract.RunParallel(t, "retire on unknown slug and bad args", testSpecRetireBadArgs)
	contract.RunParallel(t, "retire refuses an orphaned pickup", testSpecRetireOrphanedPickup)
}

func testSpecRetireMergedDeletesBoth(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("specs/x.md", "# x\n\nStatus: implemented\n")
	f.WriteFile("reviews/x.md", "# review of x\n")
	f.CommitAll("merge x")

	p := f.Bench("spec", "retire", "x")
	p.RequireExit(0)
	p.RequireContains(p.Stdout, "retired: reviews/x.md")
	p.RequireContains(p.Stdout, "retired: specs/x.md")
	p.RequireContains(p.Stdout, "spec-retire: x")
	if f.Exists("specs/x.md") {
		t.Fatalf("spec not deleted:\n%s", p.Stdout)
	}
	if f.Exists("reviews/x.md") {
		t.Fatalf("review pickup not deleted:\n%s", p.Stdout)
	}

	// Re-run idempotency: the spec is gone, so a second retire hits the unknown-slug path.
	p2 := f.Bench("spec", "retire", "x")
	p2.RequireExit(1)
	p2.RequireContains(p2.Stdout+p2.Stderr, "not found")

	// Path-argument resolution (asserted once): a path arg resolves the same spec.
	g := contract.NewFixture(t)
	g.WriteFile("specs/y.md", "# y\n\nStatus: implemented\n")
	g.CommitAll("merge y")
	pp := g.Bench("spec", "retire", "specs/y.md")
	pp.RequireExit(0)
	pp.RequireContains(pp.Stdout, "retired: specs/y.md")
	if g.Exists("specs/y.md") {
		t.Fatal("path-arg retire did not delete the spec")
	}
}

func testSpecRetireNoPickup(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("specs/x.md", "# x\n\nStatus: implemented\n")
	f.CommitAll("merge x")

	p := f.Bench("spec", "retire", "x")
	p.RequireExit(0)
	p.RequireContains(p.Stdout, "retired: specs/x.md")
	p.RequireNotContains(p.Stdout, "reviews/") // no phantom pickup line when none exists
	if f.Exists("specs/x.md") {
		t.Fatal("spec not deleted")
	}
}

func testSpecRetireStagedRefuses(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("specs/x.md", "# x\n\nStatus: staged\n")
	f.CommitAll("stage x")

	p := f.Bench("spec", "retire", "x")
	p.RequireExit(1)
	p.RequireContains(p.Stdout+p.Stderr, "not merged-implemented")
	if !f.Exists("specs/x.md") {
		t.Fatal("staged spec was deleted — retire must refuse without deleting")
	}
}

func testSpecRetireNotAtHEAD(t *testing.T) {
	f := contract.NewFixture(t)
	// Committed as staged, then flipped to implemented only in the working tree.
	f.WriteFile("specs/x.md", "# x\n\nStatus: staged\n")
	f.CommitAll("stage x")
	f.WriteFile("specs/x.md", "# x\n\nStatus: implemented\n")

	p := f.Bench("spec", "retire", "x")
	p.RequireExit(1)
	p.RequireContains(p.Stdout+p.Stderr, "not at HEAD")
	if !f.Exists("specs/x.md") {
		t.Fatal("uncommitted implemented spec was deleted — retire must refuse")
	}
}

func testSpecRetireBadArgs(t *testing.T) {
	f := contract.NewFixture(t)
	// Unknown slug: exit 1 naming the tried paths.
	p := f.Bench("spec", "retire", "nope")
	p.RequireExit(1)
	p.RequireContains(p.Stdout+p.Stderr, "not found")
	p.RequireContains(p.Stdout+p.Stderr, "specs/nope.md")
	// Extra positional, unknown flag, and missing slug: usage errors exit 2.
	f.Bench("spec", "retire", "x", "y").RequireExit(2)
	f.Bench("spec", "retire", "--x").RequireExit(2)
	f.Bench("spec", "retire").RequireExit(2)
}

func testSpecRetireOrphanedPickup(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("reviews/x.md", "# orphan review\n")
	f.CommitAll("orphan pickup")

	p := f.Bench("spec", "retire", "x")
	p.RequireExit(1)
	p.RequireContains(p.Stdout+p.Stderr, "orphaned review pickup")
	p.RequireContains(p.Stdout+p.Stderr, "reviews/x.md")
	if !f.Exists("reviews/x.md") {
		t.Fatal("orphaned pickup was auto-cleaned — retire must refuse")
	}
}

func TestRuntimeSpecHistoryContracts(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	contract.RunParallel(t, "history merges retire and delete commits, deduped, newest first", testSpecHistoryMergedDeduped)
	contract.RunParallel(t, "history accepts a bare slug or a specs/<slug>.md path", testSpecHistorySlugForms)
	contract.RunParallel(t, "history renders the definitive empty state for an unknown slug", testSpecHistoryEmptyState)
	contract.RunParallel(t, "history usage errors exit 2", testSpecHistoryUsageErrors)
	contract.RunParallel(t, "history refuses a control-byte subject at exit 1", testSpecHistoryControlBytePosture)
	contract.RunParallel(t, "history slugs with spaces or glob characters find only their own commit", testSpecHistoryGlobSlug)
	contract.RunParallel(t, "history invoked from a subdirectory sees the full-repo history", testSpecHistorySubdirInvocation)
	contract.RunParallel(t, "history outside a git repo exits 1 with the not-in-repo error", testSpecHistoryNotInRepo)
}

// testSpecHistoryMergedDeduped pins story 3 and the story-1 basic listing: a
// `bench spec retire`-shaped commit (deletes the file AND carries the `spec-retire:`
// message) is merged and deduped to one `retire`-tagged row, a delete-only commit (no
// retire message) is tagged `delete`, and the merged table is newest-first.
func testSpecHistoryMergedDeduped(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("README.md", "r\n")
	f.CommitAll("c1")
	f.Git("branch", "-m", "main")

	// foo: a spec-retire-shaped commit — deletes the file and carries the message —
	// must appear exactly once, tagged retire.
	f.WriteFile("specs/foo.md", "# foo\n")
	f.CommitAll("add foo")
	f.Git("rm", "-q", "specs/foo.md")
	f.CommitAll("spec-retire: foo")
	fooSha := strings.TrimSpace(f.Git("rev-parse", "--short", "HEAD").Stdout)

	// bar: a delete-only commit (older) — no spec-retire message — tagged delete.
	f.WriteFile("specs/bar.md", "# bar\n")
	f.CommitAll("add bar")
	f.Git("rm", "-q", "specs/bar.md")
	f.CommitAll("remove bar without retiring")
	barSha := strings.TrimSpace(f.Git("rev-parse", "--short", "HEAD").Stdout)

	fooOut := f.Bench("spec", "history", "foo")
	fooOut.RequireExit(0)
	requireOutputLine(t, fooOut, "history[1]{hash,date,kind,subject}:")
	requireHistoryRow(t, fooOut, fooSha, "retire", "spec-retire: foo")

	barOut := f.Bench("spec", "history", "bar")
	barOut.RequireExit(0)
	requireOutputLine(t, barOut, "history[1]{hash,date,kind,subject}:")
	requireHistoryRow(t, barOut, barSha, "delete", "remove bar without retiring")
}

// testSpecHistorySlugForms pins story 2: a bare slug and its specs/<slug>.md path
// resolve identically.
func testSpecHistorySlugForms(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("README.md", "r\n")
	f.CommitAll("c1")
	f.WriteFile("specs/x.md", "# x\n")
	f.CommitAll("add x")
	f.Git("rm", "-q", "specs/x.md")
	f.CommitAll("spec-retire: x")

	bySlug := f.Bench("spec", "history", "x")
	byPath := f.Bench("spec", "history", "specs/x.md")
	bySlug.RequireExit(0)
	byPath.RequireExit(0)
	if bySlug.Stdout != byPath.Stdout {
		t.Fatalf("slug and path forms diverged:\nslug: %s\npath: %s", bySlug.Stdout, byPath.Stdout)
	}
}

// testSpecHistoryEmptyState pins story 4's definitive empty state: a slug with no
// matching commits still exits 0 with an explicit history[0] table, not silence or an
// error.
func testSpecHistoryEmptyState(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("README.md", "r\n")
	f.CommitAll("c1")

	out := f.Bench("spec", "history", "never-existed")
	out.RequireExit(0)
	requireOutputLine(t, out, "history[0]{hash,date,kind,subject}:")
}

func testSpecHistoryUsageErrors(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("README.md", "r\n")
	f.CommitAll("c1")
	for _, args := range [][]string{
		{"spec", "history"},           // no argument
		{"spec", "history", "a", "b"}, // extra positional
		{"spec", "history", "--x"},    // unknown flag
	} {
		f.Bench(args...).RequireExit(2)
	}
}

// testSpecHistoryControlBytePosture pins story 5: a commit subject carrying a control
// byte (ESC) makes the command exit 1 with the unrepresentable-TOON-cell error — the
// same refusal bench diff --full already has for the same failure — instead of a
// mangled row.
func testSpecHistoryControlBytePosture(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("README.md", "r\n")
	f.CommitAll("c1")
	f.WriteFile("specs/ctrl.md", "# ctrl\n")
	f.CommitAll("add ctrl")
	f.Git("rm", "-q", "specs/ctrl.md")
	f.CommitAll("spec-retire: ctrl \x1b esc")

	out := f.Bench("spec", "history", "ctrl")
	out.RequireExit(1)
	requireOutputPrefix(t, out, "error: unrepresentable TOON cell")
}

// testSpecHistoryGlobSlug pins story 5's hostile-input-checklist row: a slug
// containing a space, and a slug containing a glob character (*), each find only
// their own commit — a non-literal pathspec or grep pattern would glob-match the
// wrong file or misinterpret the metacharacter.
func testSpecHistoryGlobSlug(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("README.md", "r\n")
	f.CommitAll("c1")

	f.WriteFile("specs/weird name.md", "# weird name\n")
	f.CommitAll("add weird name")
	f.Git("rm", "-q", "specs/weird name.md")
	f.CommitAll("spec-retire: weird name")

	f.WriteFile("specs/weird*name.md", "# weird star name\n")
	f.CommitAll("add weird star name")
	f.Git("rm", "-q", "specs/weird*name.md")
	f.CommitAll("spec-retire: weird*name")

	spaceOut := f.Bench("spec", "history", "weird name")
	spaceOut.RequireExit(0)
	requireOutputLine(t, spaceOut, "history[1]{hash,date,kind,subject}:")
	spaceOut.RequireContains(spaceOut.Stdout, "spec-retire: weird name")
	spaceOut.RequireNotContains(spaceOut.Stdout, "weird*name")

	starOut := f.Bench("spec", "history", "weird*name")
	starOut.RequireExit(0)
	requireOutputLine(t, starOut, "history[1]{hash,date,kind,subject}:")
	starOut.RequireContains(starOut.Stdout, "weird*name")
}

// testSpecHistorySubdirInvocation pins story 5's cwd-deeper-than-root row: the
// `:(literal,top)` pathspec magic anchors the delete query at the repo root, so
// invocation from a subdirectory still sees the same history as from the root.
func testSpecHistorySubdirInvocation(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("README.md", "r\n")
	f.CommitAll("c1")
	f.WriteFile("specs/sub.md", "# sub\n")
	f.CommitAll("add sub")
	f.Git("rm", "-q", "specs/sub.md")
	f.CommitAll("spec-retire: sub")
	contract.Mkdir(t, filepath.Join(f.Root, "deep", "sub"))

	fromRoot := f.Bench("spec", "history", "sub")
	fromRoot.RequireExit(0)

	fromSub := contract.RunAt(t, f, filepath.Join(f.Root, "deep", "sub"), nil, "bash", benchPath(t), "spec", "history", "sub")
	fromSub.RequireExit(0)
	if fromRoot.Stdout != fromSub.Stdout {
		t.Fatalf("subdirectory invocation diverged from repo root:\nroot: %s\nsub:  %s", fromRoot.Stdout, fromSub.Stdout)
	}
}

// testSpecHistoryNotInRepo pins the not-in-repo posture every AXI query command
// shares: outside a git repository, history exits 1 with the structured not-in-repo
// error rather than a leaked git failure.
func testSpecHistoryNotInRepo(t *testing.T) {
	f := contract.NewFixture(t, contract.WithNoRepo())
	out := f.Bench("spec", "history", "anything")
	out.RequireExit(1)
	requireOutputPrefix(t, out, "error: ")
}

// requireOutputLine asserts one exact stdout line is present — the same shape helper
// internal/contract/axi's tests use for bench diff, ported here for bench spec
// history since the two packages don't share test-only helpers.
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

// requireHistoryRow asserts one history table row, tolerating the TOON encoder's own
// choice to quote a leading-digit-looking short hash or a subject containing a colon
// — the same tolerance requireLogRow applies to bench diff's log table.
func requireHistoryRow(t *testing.T, probe contract.Probe, hash, kind, subject string) {
	t.Helper()
	for _, got := range strings.Split(strings.TrimSuffix(probe.Stdout, "\n"), "\n") {
		if !strings.HasPrefix(got, "  ") {
			continue
		}
		cells := strings.SplitN(strings.TrimPrefix(got, "  "), ",", 4)
		if len(cells) != 4 {
			continue
		}
		gotHash := strings.Trim(cells[0], `"`)
		gotSubject := strings.Trim(cells[3], `"`)
		if gotHash == hash && cells[2] == kind && gotSubject == subject {
			return
		}
	}
	t.Fatalf("missing history row hash=%q kind=%q subject=%q\nstdout:\n%s", hash, kind, subject, probe.Stdout)
}
