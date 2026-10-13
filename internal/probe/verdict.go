package probe

import (
	"strings"

	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/testreport"
	"github.com/gibbonmi/bench/internal/toon"
)

// expectMissedCause is the cause of a failed mutated run whose failed tests do not include
// the --expect witness. The bite came from another test, so it proves nothing about the
// test the caller named.
const expectMissedCause = "expect-missed"

// mutantFields is the schema of the table that follows probe_execution when the mutated run
// failed to build.
var mutantFields = []string{"first_build_line"}

// firstFailure answers the first_failure cell of one run: the first failed test name in
// the failures-table order, or nil when the run named no failed test, so the cell prints
// null rather than a name that no test carries.
func firstFailure(outcome testreport.Outcome) *string {
	if len(outcome.FailedNames) == 0 {
		return nil
	}
	return &outcome.FailedNames[0]
}

// mutatedCells answers the cells and the exit of the mutated run's answer. expect is the
// --expect value, or empty when the caller named no witness. The baseline cell carries the
// kind the baseline run observed, so the row joins the two runs it reports rather than
// asserting the kind the refusal order implies.
func mutatedCells(outcome testreport.Outcome, expect string, baseline testreport.OutcomeKind) (verdictCells, int) {
	cells := verdictCells{
		cause:    string(outcome.Kind),
		failed:   outcome.FailedTests,
		first:    firstFailure(outcome),
		restored: restoredCell,
		baseline: string(baseline),
		ran:      outcome.Ran,
	}
	var code int
	cells.verdict, code = verdictFor(outcome.Kind)
	if outcome.Kind == testreport.OutcomeFailed && expect != "" && !witnessed(outcome.FailedNames, expect) {
		cells.verdict, cells.cause, code = "invalid", expectMissedCause, 1
	}
	if outcome.Kind == testreport.OutcomeBuildFailed {
		line := sanitize.Preview(outcome.BuildLine)
		cells.buildLine = &line
	}
	return cells, code
}

// mutantTable answers the mutant table of a mutated run that failed to build, or nothing
// for every other answer.
func (c verdictCells) mutantTable() (string, error) {
	if c.buildLine == nil {
		return "", nil
	}
	return toon.TableTyped("mutant", mutantFields, [][]any{{*c.buildLine}})
}

// witnessed reports whether a failed test is the expected test or one of its subtests. A
// subtest failure fails its parent, so the parent's name matches it, but a longer name
// that only starts with the value is another test.
func witnessed(failed []string, expect string) bool {
	for _, name := range failed {
		if name == expect || strings.HasPrefix(name, expect+"/") {
			return true
		}
	}
	return false
}
