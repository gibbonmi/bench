package testreport

import (
	"errors"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/toon"
)

// invalidInventoryTitle is the error title of a canary inventory the fixtures face cannot read.
const invalidInventoryTitle = "canary inventory invalid"

// checkFixtures lists each fixture whose canary inventory owner is check, sorted by its
// repo-relative path. It reads the inventory only, so it selects no run binary and starts
// no Go child. A canary directory with no fixture lists nothing; that is an answer, not a
// refusal.
func checkFixtures(root, check string) (Outcome, string, int) {
	found, err := canary.Fixtures(filepath.Join(root, "tests", "canary"))
	if err != nil && !errors.Is(err, canary.ErrNoFixtures) {
		return refusedOutcome(toon.Errorf(invalidInventoryTitle, err.Error())+"\n", 1)
	}
	rows := [][]string{}
	for name, fixture := range found {
		if fixture.Check != check {
			continue
		}
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
