package testreport

import (
	"slices"
	"strconv"
	"strings"
	"testing"
)

// The header and the kind cells here are authored apart from the face, and the name
// list comes from the help output rather than from its producer, so a lost check, a
// fixed kind, or a count of fixtures in place of families reds these rows.

const checksHeader = "{name,kind,families}:"

// checksFace runs `--checks` through inventoryFace and answers its rows by name. It
// fails the test unless the face exits 0 with the authored header and one row per name.
func checksFace(t *testing.T, root string) (names []string, rows map[string][]string) {
	t.Helper()
	output, code := inventoryFace(t, root, "--checks")
	lines := strings.Split(strings.TrimSuffix(output, "\n"), "\n")
	if code != 0 || !strings.HasPrefix(output, "checks[") || !strings.HasSuffix(lines[0], checksHeader) {
		t.Fatalf("--checks = %d, %q; want 0 and the checks table", code, output)
	}
	rows = map[string][]string{}
	for _, line := range lines[1:] {
		cells := strings.Split(strings.TrimPrefix(line, "  "), ",")
		names = append(names, cells[0])
		rows[cells[0]] = cells
	}
	if lines[0] != "checks["+strconv.Itoa(len(names))+"]"+checksHeader {
		t.Fatalf("--checks header %q counts other than its %d rows", lines[0], len(names))
	}
	return names, rows
}

// familiesCell answers the families cell of check in a --checks run over root.
func familiesCell(t *testing.T, root, check string) string {
	t.Helper()
	_, rows := checksFace(t, root)
	row, found := rows[check]
	if !found {
		t.Fatalf("--checks has no row for %q", check)
	}
	return row[2]
}

func TestInventoryFacesStartNoChild(t *testing.T) {
	root := t.TempDir()
	plantFixture(t, root, "package-core-guard", "a", "")
	for _, args := range [][]string{{"--checks"}, {"--check", "package-core-guard", "--fixtures"}} {
		if output, code := inventoryFace(t, root, args...); code != 0 {
			t.Errorf("Command(%q) = %d, %q; want 0", args, code, output)
		}
	}
}

func TestInventoryFacesAbsentCanaryDirectory(t *testing.T) {
	root := t.TempDir()
	output, code := inventoryFace(t, root, "--check", "package-core-guard", "--fixtures")
	if code != 0 || output != emptyFixturesTable {
		t.Errorf("absent tests/canary fixtures = %d, %q; want 0 and %q", code, output, emptyFixturesTable)
	}
	names, rows := checksFace(t, root)
	for _, name := range names {
		if rows[name][2] != "0" {
			t.Errorf("absent tests/canary row %q; want families 0", rows[name])
		}
	}
}

func TestChecksFaceEqualsHelpList(t *testing.T) {
	help, code := Command(t.TempDir(), []string{"--help"})
	_, list, found := strings.Cut(help, "\nchecks:\n")
	if code != 0 || !found {
		t.Fatalf("--help = %d, %q; want 0 and the check list", code, help)
	}
	want := strings.Fields(list)
	names, _ := checksFace(t, t.TempDir())
	if !slices.Equal(names, want) {
		t.Fatalf("--checks names = %q; want the help list %q", names, want)
	}
}

func TestChecksFaceKinds(t *testing.T) {
	_, rows := checksFace(t, t.TempDir())
	for check, kind := range map[string]string{"line-routing": "conformance", "system": "system", "prose": "prose"} {
		if row := rows[check]; len(row) != 3 || row[1] != kind {
			t.Errorf("--checks row %q; want kind %q", row, kind)
		}
	}
}

func TestChecksFaceCountsFamilies(t *testing.T) {
	root := t.TempDir()
	plantFixture(t, root, "package-core-guard", "a", "")
	plantFixture(t, root, "package-core-guard", "b", "")
	plantFixture(t, root, "guard-classifier-table", "c", "package-core-guard")
	if got := familiesCell(t, root, "package-core-guard"); got != "2" {
		t.Fatalf("package-core-guard families = %q; want 2", got)
	}
}

func TestChecksFaceCountsMarkedFixtureFamily(t *testing.T) {
	if got := familiesCell(t, reassignedFixtureTree(t), "default-branch-single-source"); got != "1" {
		t.Fatalf("default-branch-single-source families = %q; want 1", got)
	}
}

func TestChecksFaceIgnoresEmptyFamily(t *testing.T) {
	root := t.TempDir()
	plantFixture(t, root, "", "loose", "package-core-guard")
	output, code := inventoryFace(t, root, "--check", "package-core-guard", "--fixtures")
	want := "fixtures[1]{family,fixture,path}:\n  \"\",loose,tests/canary/loose\n"
	if code != 0 || output != want {
		t.Errorf("loose fixture = %d, %q; want 0 and %q", code, output, want)
	}
	if got := familiesCell(t, root, "package-core-guard"); got != "0" {
		t.Errorf("package-core-guard families with a loose fixture = %q; want 0", got)
	}
}

func TestChecksFaceRefusesInvalidInventory(t *testing.T) {
	root := t.TempDir()
	plantFixture(t, root, "package-core-guard", "a", "no-such-check")
	output, code := inventoryFace(t, root, "--checks")
	if code != 1 || !strings.Contains(output, "names unknown check") {
		t.Fatalf("--checks over an invalid inventory = %d, %q; want 1 and names unknown check", code, output)
	}
}

func TestHelpStartsWithGrammar(t *testing.T) {
	const grammarText = "bench test [--full] [--package <expr> | <legacy-package> | --changed] [--base <commit> [--source-tip <commit>]] [--run <go-regex>] | bench test [--full] --check <name> | bench test [--full] --check system --run <go-regex> | bench test --check <name> --fixtures | bench test --checks"
	if help, code := Command(t.TempDir(), []string{"--help"}); code != 0 || !strings.HasPrefix(help, "usage: "+grammarText+"\nnotes:\n") {
		t.Fatalf("--help = %d, %q; want 0 and the grammar text first", code, help)
	}
}
