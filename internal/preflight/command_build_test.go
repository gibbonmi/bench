package preflight

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCommandBuildFresh is B1 (the fresh-build contract test). With no
// tickets/ directory at all, `build` prints a not-applicable row for each
// of rows-owned, rows-membership, and diff-nonempty, individually
// asserted. It runs the rest for real, and exits 0 when they are green.
func TestCommandBuildFresh(t *testing.T) {
	_, slug := seedBuildFresh(t)

	out, code := Command([]string{"build", slug})
	if code != 0 {
		t.Fatalf("Command exit = %d, want 0; output:\n%s", code, out)
	}
	for _, name := range []string{"base-current,green", "paths-authorized,green"} {
		if !strings.Contains(out, name) {
			t.Errorf("output missing %s row:\n%s", name, out)
		}
	}
	for _, name := range []string{"rows-owned,not-applicable", "rows-membership,not-applicable", "diff-nonempty,not-applicable"} {
		if !strings.Contains(out, name) {
			t.Errorf("output missing %s row:\n%s", name, out)
		}
	}
}

// TestCommandBuildResumedTicketsRunForReal is B2's present-tickets half (the
// resumed-build contract test). With a present tickets/ directory citing
// every declared row, `build` runs rows-owned and rows-membership for real
// — green, not not-applicable — while diff-nonempty stays not-applicable.
func TestCommandBuildResumedTicketsRunForReal(t *testing.T) {
	_, slug := seedConformant(t) // seedConformant's tickets/one.md cites PF1 and PF2

	out, code := Command([]string{"build", slug})
	if code != 0 {
		t.Fatalf("Command exit = %d, want 0; output:\n%s", code, out)
	}
	for _, name := range []string{"rows-owned,green", "rows-membership,green"} {
		if !strings.Contains(out, name) {
			t.Errorf("output missing %s row (present tickets/ must run the check for real):\n%s", name, out)
		}
	}
	if !strings.Contains(out, "diff-nonempty,not-applicable") {
		t.Errorf("output missing diff-nonempty,not-applicable row:\n%s", out)
	}
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
	runGit(t, "checkout", "-q", "main")
	mustWriteFile(t, "unrelated.txt", "advance main\n")
	runGit(t, "add", "unrelated.txt")
	runGit(t, "commit", "-q", "-m", "advance main")
	runGit(t, "checkout", "-q", "feature")

	out, code := Command([]string{"build", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.Contains(out, "base-current,red") {
		t.Errorf("output missing red base-current row:\n%s", out)
	}
	for _, name := range []string{"rows-owned,not-applicable", "rows-membership,not-applicable", "diff-nonempty,not-applicable"} {
		if !strings.Contains(out, name) {
			t.Errorf("output missing %s row:\n%s", name, out)
		}
	}
}

// TestCommandBuildOutOfFenceRed is B4 (the build out-of-fence contract test): a
// tracked change outside every fence entry makes paths-authorized red in build mode,
// exit 1.
func TestCommandBuildOutOfFenceRed(t *testing.T) {
	_, slug := seedBuildFresh(t)
	mustWriteFile(t, "unfenced/other.go", "package other\n")
	runGit(t, "add", "unfenced/other.go")
	runGit(t, "commit", "-q", "-m", "out of fence")

	out, code := Command([]string{"build", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.Contains(out, "paths-authorized,red") || !strings.Contains(out, "unfenced/other.go") {
		t.Errorf("output missing red paths-authorized naming unfenced/other.go:\n%s", out)
	}
}

// ticketWithoutAcceptance is one ticket file whose Acceptance section is absent,
// so the rendered tickets-parse row has a fault of its own to name.
func ticketWithoutAcceptance(blocker, covers string, writes ...string) string {
	return "# A\n\nBlocked by: " + blocker + "\nWrites: " + strings.Join(writes, ", ") +
		"\nCovers: " + covers + "\n\n## What to build\n\nOne thing.\n"
}

// seedSixGrammarReds plants one tree in which each of the six grammar rows reds
// on its own fact: an absent section, a blocker cycle, a typo path, a
// fixture-pinned path, a bound package, and a system-tagged test file.
func seedSixGrammarReds(t *testing.T) (root, slug string) {
	t.Helper()
	slug = "example"
	root = initRepo(t)
	mustWriteFile(t, "specs/"+slug+"/spec.md", specBody(slug))
	mustWriteFile(t, "tests/canary/example-family/pinning-fixture/EXPECT", "planted diagnostic\n")
	mustWriteFile(t, "tests/canary/example-family/pinning-fixture/BASE", "internal/example/pinned.go\n")
	mustWriteFile(t, "internal/example/pinned.go", "package example\n")
	mustWriteFile(t, "internal/toon/toon_test.go", "package toon\n")
	mustWriteFile(t, "internal/example/sys_test.go", "//go:build system\n\npackage example\n")
	mustWriteFile(t, "specs/"+slug+"/tickets/a.md", ticketWithoutAcceptance(
		"b.md", "PF1",
		"internal/example/pinned.go", "internal/toon/toon_test.go",
		"internal/example/sys_test.go", "gone.go"))
	mustWriteFile(t, "specs/"+slug+"/tickets/b.md",
		"# B\n\nBlocked by: a.md\nWrites: specs\nCovers: PF2\n\n"+
			"## What to build\n\nBuild it.\n\n## Acceptance\n\n- [ ] It is built.\n")
	runGit(t, "add", ".")
	runGit(t, "commit", "-q", "-m", "c0")
	runGit(t, "checkout", "-q", "-b", "feature")
	mustWriteFile(t, "internal/"+slug+"/foo.go", "package example\n")
	runGit(t, "add", "internal/"+slug+"/foo.go")
	runGit(t, "commit", "-q", "-m", "c1")
	return root, slug
}

// TestCommandBuildRendersSixGrammarRows covers TG26. `bench preflight build`
// renders the six grammar rows in their fixed order, each red with its own
// detail, so one red never hides another.
func TestCommandBuildRendersSixGrammarRows(t *testing.T) {
	_, slug := seedSixGrammarReds(t)

	out, code := Command([]string{"build", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	rows := []struct{ check, detail string }{
		{"tickets-parse", "ticket grammar fault(s): a.md: missing Acceptance section"},
		{"blockers-resolve", "blocker cycle(s): cycle edge b.md -> a.md"},
		{"writes-resolve", "Writes: entry names no tree path and carries no (new) marker: a.md: gone.go"},
		{"fixture-closure", "Writes: entry names a fixture-pinned path without naming the fixture: a.md: internal/example/pinned.go is pinned by tests/canary/example-family/pinning-fixture"},
		{"registry-closure", "Writes: entry names a bound package without naming every bound file: a.md: internal/toon/toon_test.go requires internal/conformance/data_handling_test.go"},
		{"kit-pin", "ticket writes a system-tagged test file without stating BENCH_KIT: a.md: internal/example/sys_test.go"},
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
