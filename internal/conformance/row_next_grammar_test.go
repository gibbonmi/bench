package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/roadmap"
)

const (
	// rowNextGrammarCommand is the phase file that carries the human-facing token table.
	rowNextGrammarCommand = ".agents/commands/bench-drain.md"
	// rowNextGrammarTableHead names the table's first header cell, which is how the table
	// is found among the command's other tables.
	rowNextGrammarTableHead = "token"
	rowNextGrammarFamily    = "row-next-grammar"
)

// checkRowNextGrammar binds the drain command's token table to the parser's token set, so
// the grammar and its documentation cannot drift apart. The table is the human-facing
// form. `roadmap.RowNextTokens` is the one source.
func checkRowNextGrammar(root string) []string {
	body := readIfExists(filepath.Join(root, filepath.FromSlash(rowNextGrammarCommand)))
	if body == "" {
		return []string{"row-next-grammar drift: " + rowNextGrammarCommand + " is absent or unreadable, so the token table cannot be compared"}
	}
	documented, found := rowNextGrammarTable(body)
	if !found {
		return []string{"row-next-grammar drift: " + rowNextGrammarCommand + " renders no '| token | phase |' table"}
	}
	var diags []string
	seen := map[string]bool{}
	for _, token := range documented {
		if seen[token] {
			diags = append(diags, "row-next-grammar drift: "+rowNextGrammarCommand+" token table lists token '"+token+"' more than once")
		}
		seen[token] = true
	}
	expected := roadmap.RowNextTokens()
	for _, token := range expected {
		if !seen[token] {
			diags = append(diags, "row-next-grammar drift: "+rowNextGrammarCommand+" token table lacks token '"+token+"'")
		}
	}
	for _, token := range documented {
		if !slices.Contains(expected, token) {
			diags = append(diags, "row-next-grammar drift: "+rowNextGrammarCommand+" token table carries token '"+token+"', which the parser does not accept")
		}
	}
	return diags
}

// rowNextGrammarTable returns the first-column tokens of the command's token table, in
// document order, and whether the table exists. The table ends at the first line that is
// not a row, so prose after it is never read as a token.
func rowNextGrammarTable(body string) ([]string, bool) {
	var tokens []string
	header := false
	for _, line := range strings.Split(body, "\n") {
		row, isRow := markdownRow(line)
		if !isRow {
			if header {
				break
			}
			continue
		}
		if !header {
			header = len(row) >= 2 && strings.ToLower(row[0]) == rowNextGrammarTableHead
			continue
		}
		if isRuleRow(row) || row[0] == "" {
			continue
		}
		tokens = append(tokens, row[0])
	}
	return tokens, header
}

// rowNextGrammarFixtureClasses names the diagnostic class each fixture in the family
// plants. The independently authored inventory makes a deleted fixture red.
var rowNextGrammarFixtureClasses = map[string]string{
	"token-table-lacks-kit-edit": "token table lacks token 'kit-edit'",
}

func validateRowNextGrammarFixtureInventory(fixtures map[string]canary.Fixture) error {
	expected := make([]string, 0, len(rowNextGrammarFixtureClasses))
	for name := range rowNextGrammarFixtureClasses {
		expected = append(expected, name)
	}
	if len(expected) != 1 {
		return fmt.Errorf("row-next-grammar fixture inventory has %d entries, want 1", len(expected))
	}
	sort.Strings(expected)
	actual := make([]string, 0, len(expected))
	for name, fixture := range fixtures {
		if fixture.Family != rowNextGrammarFamily {
			continue
		}
		actual = append(actual, name)
		class, named := rowNextGrammarFixtureClasses[name]
		if !named {
			continue
		}
		data, err := os.ReadFile(filepath.Join(fixture.Dir, "EXPECT"))
		if err != nil {
			return fmt.Errorf("read %s EXPECT: %w", name, err)
		}
		if expect := strings.TrimSpace(string(data)); !strings.Contains(expect, class) {
			return fmt.Errorf("fixture %q EXPECT %q does not plant its %q class", name, expect, class)
		}
	}
	sort.Strings(actual)
	if !slices.Equal(actual, expected) {
		return fmt.Errorf("row-next-grammar fixture inventory = %v, want %v", actual, expected)
	}
	return nil
}

func TestRowNextGrammarFixturesCoverEveryDiagnosticClass(t *testing.T) {
	h := NewHarness(t)
	fixtures, err := canary.Fixtures(filepath.Join(h.KitRoot, "tests", "canary"))
	if err != nil {
		t.Fatal(err)
	}
	if err := validateRowNextGrammarFixtureInventory(fixtures); err != nil {
		t.Fatal(err)
	}
}

func TestRowNextGrammarFixtureInventoryRejectsDeletion(t *testing.T) {
	h := NewHarness(t)
	fixtures, err := canary.Fixtures(filepath.Join(h.KitRoot, "tests", "canary"))
	if err != nil {
		t.Fatal(err)
	}
	// The family holds one fixture, so removing its directory would empty the family.
	// Dropping the entry from the loaded inventory presents the validator with the same
	// absence.
	delete(fixtures, "token-table-lacks-kit-edit")
	err = validateRowNextGrammarFixtureInventory(fixtures)
	if err == nil || !strings.Contains(err.Error(), "token-table-lacks-kit-edit") {
		t.Fatalf("deleted fixture inventory error = %v, want token-table-lacks-kit-edit omission", err)
	}
}

// TestRowNextGrammarBindsTableToParserTokens pins RF16: the live command's table agrees
// with the parser. A table that lacks `kit-edit` or carries a stranger yields a drift
// diagnostic naming that token.
func TestRowNextGrammarBindsTableToParserTokens(t *testing.T) {
	h := NewHarness(t)
	if diags := checkRowNextGrammar(h.KitRoot); len(diags) != 0 {
		t.Fatalf("live tree drifted: %v", diags)
	}
	write := func(t *testing.T, table string) string {
		t.Helper()
		root := t.TempDir()
		dir := filepath.Join(root, ".agents", "commands")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "bench-drain.md"), []byte("# /bench-drain\n\n"+table), 0o644); err != nil {
			t.Fatal(err)
		}
		return root
	}
	const complete = "| token | phase |\n|---|---|\n| `shape` | a |\n| `spec` | b |\n| `ticket` | c |\n| `decide` | d |\n| `ready-for-agent` | e |\n| `kit-edit` | f |\n"
	if diags := checkRowNextGrammar(write(t, complete)); len(diags) != 0 {
		t.Fatalf("complete table = %v, want no diagnostic", diags)
	}
	lacking := strings.Replace(complete, "| `kit-edit` | f |\n", "", 1)
	if diags := checkRowNextGrammar(write(t, lacking)); !containsDiagnostic(diags, "token table lacks token 'kit-edit'") {
		t.Fatalf("table without kit-edit = %v, want drift naming kit-edit", diags)
	}
	extra := complete + "| `refactor` | f |\n"
	if diags := checkRowNextGrammar(write(t, extra)); !containsDiagnostic(diags, "carries token 'refactor'") {
		t.Fatalf("table with refactor = %v, want drift naming refactor", diags)
	}
	if diags := checkRowNextGrammar(write(t, "no table here\n")); !containsDiagnostic(diags, "renders no '| token | phase |' table") {
		t.Fatalf("command without table = %v, want missing-table diagnostic", diags)
	}
}
