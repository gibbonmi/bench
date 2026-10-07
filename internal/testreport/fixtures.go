package testreport

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/toon"
)

// invalidInventoryTitle is the error title of a canary inventory the inventory faces cannot read.
const invalidInventoryTitle = "canary inventory invalid"

// ownedFixtures answers, by fixture name, each fixture whose canary inventory owner is check.
func ownedFixtures(found map[string]canary.Fixture, check string) map[string]canary.Fixture {
	owned := map[string]canary.Fixture{}
	for name, fixture := range found {
		if fixture.Check == check {
			owned[name] = fixture
		}
	}
	return owned
}

// checkFixtures lists each fixture that check owns, sorted by its repo-relative path. It
// reads the inventory only, so it selects no run binary and starts no Go child.
func checkFixtures(root, check string) (Outcome, string, int) {
	found, err := canary.RootFixtures(root)
	if err != nil {
		return refusedOutcome(toon.Errorf(invalidInventoryTitle, err.Error())+"\n", 1)
	}
	rows := [][]string{}
	for name, fixture := range ownedFixtures(found, check) {
		path, err := filepath.Rel(root, fixture.Dir)
		if err != nil {
			return refusedOutcome(toon.Errorf(invalidInventoryTitle, err.Error())+"\n", 1)
		}
		rows = append(rows, []string{fixture.Family, name, filepath.ToSlash(path)})
	}
	slices.SortFunc(rows, func(a, b []string) int { return strings.Compare(a[2], b[2]) })
	table, err := toon.Table("fixtures", []string{"family", "fixture", "path"}, rows)
	if err != nil {
		return refusedOutcome(toon.RenderError(err)+"\n", 1)
	}
	return Outcome{Kind: OutcomeNoTestRun}, table, 0
}

// checksInventory lists each named check in the order of the help list, with its kind and
// the count of distinct non-empty family names in which it owns a fixture. It reads the
// inventory only, so it selects no run binary and starts no Go child.
func checksInventory(root string) (Outcome, string, int) {
	found, err := canary.RootFixtures(root)
	if err != nil {
		return refusedOutcome(toon.Errorf(invalidInventoryTitle, err.Error())+"\n", 1)
	}
	rows := [][]any{}
	for _, check := range namedChecks() {
		owned := ownedFixtures(found, check)
		families := map[string]bool{}
		for _, fixture := range owned {
			if fixture.Family != "" {
				families[fixture.Family] = true
			}
		}
		rows = append(rows, []any{check, namedCheckKind(check), len(families)})
	}
	table, err := toon.TableTyped("checks", []string{"name", "kind", "families"}, rows)
	if err != nil {
		return refusedOutcome(toon.RenderError(err)+"\n", 1)
	}
	return Outcome{Kind: OutcomeNoTestRun}, table, 0
}
