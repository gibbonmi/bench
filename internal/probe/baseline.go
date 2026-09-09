package probe

import (
	"github.com/gibbonmi/bench/internal/testreport"
	"github.com/gibbonmi/bench/internal/toon"
)

// untouchedCell is the restored cell of a probe that ended before it wrote anything. A
// baseline refusal preserved nothing and mutated nothing, so neither `yes` nor `no`
// describes what happened to the subject.
const untouchedCell = "untouched"

// baselineCause prefixes the baseline's own outcome kind, so a cause names which of the two
// runs reached the kind it carries.
const baselineCause = "baseline-"

// gradeBaseline runs the prepared selection once over the unmutated tree and answers the
// whole probe output when that run did not pass. A tree that is already red, that does not
// build, or that runs no test cannot support a verdict about the mutation, so the probe
// ends here: before the preserved copy, before the mutation, and before any write.
//
// The evidence after the rows is the baseline's own report, so its failures table names the
// tests that were red before the probe ran. The `ran` cell reads 0, because the mutated run
// never started, and the `baseline` cell names the kind the caller has to fix first.
func gradeBaseline(root string, subject subject, mutation string, request testreport.Request) (string, int) {
	outcome, report, _ := testreport.Execute(root, request)
	if outcome.Kind == testreport.OutcomePassed {
		return "", 0
	}
	cells := verdictCells{
		verdict:  "invalid",
		cause:    baselineCause + string(outcome.Kind),
		failed:   outcome.FailedTests,
		restored: untouchedCell,
		baseline: string(outcome.Kind),
	}
	out, err := rows(subject, mutation, cells, request)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return out + report, 1
}
