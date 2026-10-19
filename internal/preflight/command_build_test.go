package preflight

import (
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// TestCommandBuildFresh is B1 (the fresh-build contract test). With no
// tickets/ directory at all, `build` counts a not-applicable verdict for each
// of rows-owned, rows-membership, and diff-nonempty, individually
// asserted. It runs the rest for real, and exits 0 when they are green. A
// fresh build runs before its first edit, and with no ticket no entry
// authorizes a change, so the feature branch returns to the base.
func TestCommandBuildFresh(t *testing.T) {
	_, slug := seedBuildFresh(t)
	preflighttest.RunGit(t, "reset", "-q", "--hard", "main")

	out, code := Command([]string{"build", slug})
	if code != 0 {
		t.Fatalf("Command exit = %d, want 0; output:\n%s", code, out)
	}
	verdicts := renderedVerdicts(t, out, modeBuild, slug)
	requireVerdict(t, verdicts, verdictGreen, "base-current", "paths-authorized")
	requireVerdict(t, verdicts, verdictNA, "rows-owned", "rows-membership", "diff-nonempty")
}

// TestCommandBuildResumedTicketsRunForReal is B2's present-tickets half (the
// resumed-build contract test). With a present tickets/ directory citing
// every declared row, `build` runs rows-owned and rows-membership for real
// — green, not not-applicable — while diff-nonempty stays not-applicable.
func TestCommandBuildResumedTicketsRunForReal(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t) // SeedConformant's tickets/one.md cites PF1 and PF2
	preflighttest.ActiveAssignment(t, root, root)

	out, code := Command([]string{"build", slug})
	if code != 0 {
		t.Fatalf("Command exit = %d, want 0; output:\n%s", code, out)
	}
	verdicts := renderedVerdicts(t, out, modeBuild, slug)
	requireVerdict(t, verdicts, verdictGreen, "rows-owned", "rows-membership")
	requireVerdict(t, verdicts, verdictNA, "diff-nonempty")
}

// TestCommandBuildCoversAcceptsEveryMapRowID pins one row-ID grammar across the
// map and the ticket. A zero-padded row the coverage map declares (PF02) is a
// valid Covers: token, so the ticket that cites it parses and owns the row.
func TestCommandBuildCoversAcceptsEveryMapRowID(t *testing.T) {
	slug := "example"
	root := preflighttest.StartRepo(t)
	spec := strings.Replace(preflighttest.SpecBody(slug), "| PF2 |", "| PF02 |", 1)
	preflighttest.MustWriteFile(t, "specs/"+slug+"/spec.md", spec)
	preflighttest.MustWriteFile(t, "specs/"+slug+"/tickets/one.md", preflighttest.TicketDoc("One", "PF1", "PF02"))
	commitmenttest.SeedAdmission(t, root, "specs/"+slug+"/spec.md")
	preflighttest.RunGit(t, "add", ".")
	preflighttest.RunGit(t, "commit", "-q", "-m", "c0")

	preflighttest.ActiveAssignment(t, root, root)
	out, _ := Command([]string{"build", slug})
	verdicts := renderedVerdicts(t, out, modeBuild, slug)
	requireVerdict(t, verdicts, verdictGreen, "tickets-parse", "rows-owned", "rows-membership")
}

// TestCommandBuildEmptyTicketsRed is B2's empty-tickets half (the empty-tickets
// contract test): a present-but-empty tickets/ directory is red — declared rows
// unowned — rather than not-applicable.
func TestCommandBuildEmptyTicketsRed(t *testing.T) {
	_, slug := seedBuildFresh(t)
	if err := os.MkdirAll(filepath.Join("specs", slug, "tickets"), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	out, code := Command([]string{"build", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.Contains(out, "rows-owned,red") {
		t.Errorf("output missing red rows-owned for a present-but-empty tickets/ dir:\n%s", out)
	}
}

// TestCommandBuildStaleBaseRedDespiteNA is B3 (the stale-base build contract test):
// base-current red in build mode exits 1 even while the ticket checks and
// diff-nonempty are not-applicable.
func TestCommandBuildStaleBaseRedDespiteNA(t *testing.T) {
	_, slug := seedBuildFresh(t)
	preflighttest.RunGit(t, "checkout", "-q", "main")
	preflighttest.MustWriteFile(t, "unrelated.txt", "advance main\n")
	preflighttest.RunGit(t, "add", "unrelated.txt")
	preflighttest.RunGit(t, "commit", "-q", "-m", "advance main")
	preflighttest.RunGit(t, "checkout", "-q", "feature")

	out, code := Command([]string{"build", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.Contains(out, "base-current,red") {
		t.Errorf("output missing red base-current row:\n%s", out)
	}
	verdicts := renderedVerdicts(t, out, modeBuild, slug)
	requireVerdict(t, verdicts, verdictNA, "rows-owned", "rows-membership", "diff-nonempty")
}

// TestCommandBuildOutOfFenceRed is B4 (the build out-of-fence contract test): a
// tracked change outside every fence entry makes paths-authorized red in build mode,
// exit 1.
func TestCommandBuildOutOfFenceRed(t *testing.T) {
	_, slug := seedBuildFresh(t)
	preflighttest.MustWriteFile(t, "unfenced/other.go", "package other\n")
	preflighttest.RunGit(t, "add", "unfenced/other.go")
	preflighttest.RunGit(t, "commit", "-q", "-m", "out of fence")

	out, code := Command([]string{"build", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.Contains(out, "paths-authorized,red") || !strings.Contains(out, "unfenced/other.go") {
		t.Errorf("output missing red paths-authorized naming unfenced/other.go:\n%s", out)
	}
}

// TestPreflightBuildRefusesRemovedForms covers CR07 and CR08. The grammar knows neither
// removed flag, so each call takes the usage refusal at exit 2 before any gather.
func TestPreflightBuildRefusesRemovedForms(t *testing.T) {
	for _, args := range [][]string{
		{"build", "s", "--propose-writes", "--ticket", "a.md", "--base", "B", "--source-tip", "T"},
		{"build", "s", "--breakdown"},
	} {
		out, code := Command(args)
		if want := "unknown argument: " + args[2]; code != 2 || !strings.Contains(out, want) {
			t.Errorf("%v = (%d), want exit 2 naming %q:\n%s", args, code, want, out)
		}
	}
}

// TestBuildVerdictReadsNoClosureInventory covers CR14. The committed canary inventory
// holds a fixture whose MUTATE.json does not parse, and the build verdict still renders,
// because no verdict row reads the fixture inventory.
func TestBuildVerdictReadsNoClosureInventory(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	preflighttest.CommitChargeSource(t, root, slug, map[string]string{
		"tests/canary/example-family/broken/EXPECT":      "planted diagnostic\n",
		"tests/canary/example-family/broken/MUTATE.json": "{",
	})
	out, code := Command([]string{"build", slug})
	if code != 0 || !strings.HasPrefix(out, "phase: build\n") {
		t.Fatalf("build over an unparsable canary inventory = (%d), want a green verdict:\n%s", code, out)
	}
}

// TestCommandBuildRefusesAnUncleanWritesEntry covers CR21 and CR22. The ticket grammar
// refuses a `Writes:` entry that escapes the repository or is not clean, so tickets-parse
// reds it. The fence names each entry, so no other row reds, and the outside file exists, so
// the tree probe alone would find it.
func TestCommandBuildRefusesAnUncleanWritesEntry(t *testing.T) {
	for _, entry := range []string{"internal/x/../../cmd/y.go", "../outside.txt"} {
		t.Run(entry, func(t *testing.T) {
			root, slug := preflighttest.SeedConformant(t)
			preflighttest.MustWriteFile(t, filepath.Join(filepath.Dir(root), "outside.txt"), "outside\n")
			preflighttest.CommitChargeSource(t, root, slug, map[string]string{
				"specs/" + slug + "/spec.md":        preflighttest.SpecBody(slug),
				"specs/" + slug + "/tickets/one.md": preflighttest.WritesTicketDoc("One", append(preflighttest.FenceWrites(preflighttest.ConformantFence), entry), "PF1", "PF2"),
			})
			out, code := Command([]string{"build", slug})
			if code != 1 || !strings.Contains(out, "  tickets-parse,red,") || !strings.Contains(out, entry) {
				t.Fatalf("build over the Writes: entry %s = (%d), want exit 1 with a tickets-parse red naming it:\n%s", entry, code, out)
			}
		})
	}
}

// ticketWithoutAcceptance is one ticket file whose Acceptance section is absent,
// so the rendered tickets-parse row has a fault of its own to name.
func ticketWithoutAcceptance(blocker, covers string, writes ...string) string {
	return "# A\n\nBlocked by: " + blocker + "\nWrites: " + strings.Join(writes, ", ") +
		"\nCovers: " + covers + "\n\n## What to build\n\nOne thing.\n"
}

// seedGrammarReds plants one tree in which each of the four grammar rows reds on its own
// fact: an absent section, a blocker cycle, a typo path, and a system-tagged test file. The
// ticket also writes a fixture-pinned path and a bound package without naming the fixture
// or the bound files, which a closure row would red.
func seedGrammarReds(t *testing.T) (root, slug string) {
	t.Helper()
	slug = "example"
	root = preflighttest.StartRepo(t)
	preflighttest.MustWriteFile(t, "specs/"+slug+"/spec.md", preflighttest.SpecBody(slug))
	preflighttest.MustWriteFile(t, "tests/canary/example-family/pinning-fixture/EXPECT", "planted diagnostic\n")
	preflighttest.MustWriteFile(t, "tests/canary/example-family/pinning-fixture/BASE", "internal/example/pinned.go\n")
	preflighttest.MustWriteFile(t, "internal/example/pinned.go", "package example\n")
	preflighttest.MustWriteFile(t, "internal/toon/toon_test.go", "package toon\n")
	preflighttest.MustWriteFile(t, "internal/example/sys_test.go", "//go:build system\n\npackage example\n")
	preflighttest.MustWriteFile(t, "specs/"+slug+"/tickets/a.md", ticketWithoutAcceptance(
		"b.md", "PF1",
		"internal/example/pinned.go", "internal/toon/toon_test.go",
		"internal/example/sys_test.go", "gone.go"))
	preflighttest.MustWriteFile(t, "specs/"+slug+"/tickets/b.md",
		"# B\n\nBlocked by: a.md\nWrites: specs\nCovers: PF2\n\n"+
			"## What to build\n\nBuild it.\n\n## Acceptance\n\n- [ ] It is built.\n")
	preflighttest.RunGit(t, "add", ".")
	preflighttest.RunGit(t, "commit", "-q", "-m", "c0")
	preflighttest.RunGit(t, "checkout", "-q", "-b", "feature")
	preflighttest.MustWriteFile(t, "internal/"+slug+"/foo.go", "package example\n")
	preflighttest.RunGit(t, "add", "internal/"+slug+"/foo.go")
	preflighttest.RunGit(t, "commit", "-q", "-m", "c1")
	return root, slug
}

// TestCommandBuildRendersFourGrammarRows covers TG26, CR02, CR16, and CR17. `bench
// preflight build` renders the four grammar rows in their fixed order, each red with its
// own detail, so one red never hides another. The red table holds no closure row for the
// pinned path or the bound package.
func TestCommandBuildRendersFourGrammarRows(t *testing.T) {
	_, slug := seedGrammarReds(t)

	out, code := Command([]string{"build", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	rows := []struct{ check, detail string }{
		{"tickets-parse", "ticket grammar fault(s): a.md: missing Acceptance section"},
		{"blockers-resolve", "blocker cycle(s): cycle edge b.md -> a.md"},
		{"writes-resolve", "Writes: entry names no tree path and carries no (new) marker: a.md: gone.go"},
		{"kit-pin", "ticket writes a system-tagged test file without stating BENCH_KIT: a.md: internal/example/sys_test.go"},
	}
	if strings.Contains(out, "-closure,") {
		t.Errorf("output renders a closure row:\n%s", out)
	}
	at := -1
	for i, row := range rows {
		want := "  " + row.check + ",red,\"" + row.detail + "\",\"\"\n"
		found := strings.Index(out, want)
		if found < 0 {
			t.Errorf("output missing the %s row line %q:\n%s", row.check, want, out)
			continue
		}
		if found <= at {
			t.Errorf("row %s (%d) renders out of the fixed grammar order:\n%s", row.check, i, out)
		}
		at = found
	}
}
