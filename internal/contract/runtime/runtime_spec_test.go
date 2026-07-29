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
	contract.RunParallel(t, "live flat forms refuse with migration guidance", testSpecLiveFlatRefusals)
	contract.RunParallel(t, "usage errors exit 2", testSpecUsageExitTwo)
}

func testSpecLiveFlatRefusals(t *testing.T) {
	const instruction = "flat spec layout: move to specs/<slug>/spec.md"
	for _, tc := range []struct {
		name, slug string
		write      func(contract.Fixture)
		want       []string
	}{
		{"flat only", "flat", func(f contract.Fixture) { f.WriteFile("specs/flat.md", "Status: staged\n") }, []string{instruction, "specs/flat.md", "specs/flat/spec.md"}},
		{"both forms", "collision", func(f contract.Fixture) {
			f.WriteFile("specs/collision.md", "Status: staged\n")
			f.WriteFile("specs/collision/spec.md", "Status: staged\n")
		}, []string{instruction, "specs/collision.md", "specs/collision/spec.md"}},
		{"incomplete folder", "incomplete", func(f contract.Fixture) { contract.Mkdir(t, filepath.Join(f.Root, "specs", "incomplete")) }, []string{"spec folder is missing", "specs/incomplete/spec.md"}},
		{"flat precedes incomplete folder", "partial", func(f contract.Fixture) {
			f.WriteFile("specs/partial.md", "Status: staged\n")
			contract.Mkdir(t, filepath.Join(f.Root, "specs", "partial"))
		}, []string{instruction, "specs/partial.md", "specs/partial/spec.md"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := contract.NewFixture(t)
			tc.write(f)
			p := f.Bench("spec", "implemented", tc.slug)
			p.RequireExit(1)
			for _, want := range tc.want {
				p.RequireContains(p.Stdout+p.Stderr, want)
			}
			if tc.name == "flat precedes incomplete folder" && strings.Contains(p.Stdout+p.Stderr, "spec folder is missing") {
				t.Fatalf("combined refusal used incomplete-folder diagnostic: %s", p.Stdout+p.Stderr)
			}
		})
	}
}

func testSpecFlipForm(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("specs/feature/spec.md", "# feature\n\nStatus: staged\n\n## body\n")

	p := f.Bench("spec", "implemented", "feature")
	p.RequireExit(0)
	p.RequireContains(p.Stdout, "spec implemented")
	if got := f.ReadFile("specs/feature/spec.md"); got != "# feature\n\nStatus: implemented\n\n## body\n" {
		t.Fatalf("flipped content = %q", got)
	}
}

func testSpecSlugAnchored(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("specs/feature/spec.md", "Status: staged\n")
	contract.Mkdir(t, filepath.Join(f.Root, "sub"))

	// A bare slug resolves specs/<slug>/spec.md at the repo root even when run from a subdir.
	contract.RunAt(t, f, filepath.Join(f.Root, "sub"), nil, "bash", benchPath(t), "spec", "implemented", "feature").RequireExit(0)
	contract.RequireContains(t, f.ReadFile("specs/feature/spec.md"), "Status: implemented")
}

func testSpecByteIdentity(t *testing.T) {
	f := contract.NewFixture(t)
	// No trailing newline, tab separator, and a later line that also contains "staged".
	f.WriteFile("specs/edge/spec.md", "# edge\nStatus:\tstaged\nwork was staged elsewhere")

	f.Bench("spec", "implemented", "edge").RequireExit(0)
	if got := f.ReadFile("specs/edge/spec.md"); got != "# edge\nStatus:\timplemented\nwork was staged elsewhere" {
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
				f.WriteFile("specs/s/spec.md", tc.content)
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
	contract.RunParallel(t, "retire deletes pickup and complete spec folder", testSpecRetireMergedDeletesBoth)
	contract.RunParallel(t, "retire with no pickup retires the spec alone", testSpecRetireNoPickup)
	contract.RunParallel(t, "retire refuses a staged spec without deleting", testSpecRetireStagedRefuses)
	contract.RunParallel(t, "retire refuses an implemented spec not yet at HEAD", testSpecRetireNotAtHEAD)
	contract.RunParallel(t, "retire on unknown slug and bad args", testSpecRetireBadArgs)
	contract.RunParallel(t, "retire refuses an orphaned pickup", testSpecRetireOrphanedPickup)
	contract.RunParallel(t, "retire recovers interrupted folder deletion", testSpecRetireRecoversInterruptedFolder)
	contract.RunParallel(t, "retire refuses terminal folder residue", testSpecRetireRefusesFolderResidue)
	contract.RunParallel(t, "retire keeps flat refusal ahead of incomplete-folder residue", testSpecRetireFlatPrecedence)
}

func testSpecRetireMergedDeletesBoth(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("specs/x/spec.md", "# x\n\nStatus: implemented\n")
	f.WriteFile("specs/x/tickets/work.md", "# work\n")
	f.WriteFile("reviews/x.md", "# review of x\n")
	f.CommitAll("merge x")

	p := f.Bench("spec", "retire", "x")
	p.RequireExit(0)
	p.RequireContains(p.Stdout, "retired: reviews/x.md")
	p.RequireContains(p.Stdout, "retired: specs/x")
	p.RequireContains(p.Stdout, "spec-retire: x")
	if f.Exists("specs/x") {
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
	g.WriteFile("specs/y/spec.md", "# y\n\nStatus: implemented\n")
	g.CommitAll("merge y")
	pp := g.Bench("spec", "retire", "specs/y/spec.md")
	pp.RequireExit(0)
	pp.RequireContains(pp.Stdout, "retired: specs/y")
	if g.Exists("specs/y") {
		t.Fatal("path-arg retire did not delete the spec")
	}
}

func testSpecRetireNoPickup(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("specs/x/spec.md", "# x\n\nStatus: implemented\n")
	f.CommitAll("merge x")

	p := f.Bench("spec", "retire", "x")
	p.RequireExit(0)
	p.RequireContains(p.Stdout, "retired: specs/x")
	p.RequireNotContains(p.Stdout, "reviews/") // no phantom pickup line when none exists
	if f.Exists("specs/x") {
		t.Fatal("spec not deleted")
	}
}

func testSpecRetireStagedRefuses(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("specs/x/spec.md", "# x\n\nStatus: staged\n")
	f.CommitAll("stage x")

	p := f.Bench("spec", "retire", "x")
	p.RequireExit(1)
	p.RequireContains(p.Stdout+p.Stderr, "not merged-implemented")
	if !f.Exists("specs/x/spec.md") {
		t.Fatal("staged spec was deleted — retire must refuse without deleting")
	}
}

func testSpecRetireNotAtHEAD(t *testing.T) {
	f := contract.NewFixture(t)
	// Committed as staged, then flipped to implemented only in the working tree.
	f.WriteFile("specs/x/spec.md", "# x\n\nStatus: staged\n")
	f.CommitAll("stage x")
	f.WriteFile("specs/x/spec.md", "# x\n\nStatus: implemented\n")

	p := f.Bench("spec", "retire", "x")
	p.RequireExit(1)
	p.RequireContains(p.Stdout+p.Stderr, "not at HEAD")
	if !f.Exists("specs/x/spec.md") {
		t.Fatal("uncommitted implemented spec was deleted — retire must refuse")
	}
}

func testSpecRetireRecoversInterruptedFolder(t *testing.T) {
	for _, tc := range []struct{ name, removed string }{
		{"pickup already removed", "reviews/x.md"},
		{"tickets partly removed", "specs/x/tickets"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := contract.NewFixture(t)
			f.WriteFile("specs/x/spec.md", "# x\n\nStatus: implemented\n")
			f.WriteFile("specs/x/tickets/work.md", "# work\n")
			f.WriteFile("reviews/x.md", "# pickup\n")
			f.CommitAll("merge x")
			contract.Remove(t, filepath.Join(f.Root, tc.removed))

			f.Bench("spec", "retire", "x").RequireExit(0)
			if f.Exists("specs/x") {
				t.Fatal("recoverable folder state was not retired")
			}
		})
	}
}

func testSpecRetireRefusesFolderResidue(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("specs/x/tickets/work.md", "# residue\n")
	f.CommitAll("interrupted retire x")

	p := f.Bench("spec", "retire", "x")
	p.RequireExit(1)
	p.RequireContains(p.Stdout+p.Stderr, "incomplete retired spec folder: specs/x")
	p.RequireContains(p.Stdout+p.Stderr, "remove specs/x by hand")
	if !f.Exists("specs/x/tickets/work.md") {
		t.Fatal("terminal residue was auto-cleaned")
	}
}

func testSpecRetireFlatPrecedence(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("specs/partial.md", "Status: staged\n")
	f.WriteFile("specs/partial/tickets/work.md", "# residue\n")

	p := f.Bench("spec", "retire", "partial")
	p.RequireExit(1)
	p.RequireContains(p.Stdout+p.Stderr, "flat spec layout: move to specs/<slug>/spec.md")
	p.RequireContains(p.Stdout+p.Stderr, "specs/partial.md")
	p.RequireContains(p.Stdout+p.Stderr, "specs/partial/spec.md")
	p.RequireNotContains(p.Stdout+p.Stderr, "incomplete retired spec folder")
}

func testSpecRetireBadArgs(t *testing.T) {
	f := contract.NewFixture(t)
	// Unknown slug: exit 1 naming the tried paths.
	p := f.Bench("spec", "retire", "nope")
	p.RequireExit(1)
	p.RequireContains(p.Stdout+p.Stderr, "not found")
	p.RequireContains(p.Stdout+p.Stderr, "specs/nope/spec.md")
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
