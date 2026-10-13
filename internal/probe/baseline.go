package probe

import (
	"github.com/gibbonmi/bench/internal/testreport"
	"github.com/gibbonmi/bench/internal/toon"
)

// untouchedCell is the restored cell of a probe that ended before it wrote anything. A
// baseline refusal preserved nothing and mutated nothing, so neither `yes` nor `no`
// describes what happened to the subject.
const untouchedCell = "untouched"

// untouchedCells answers the cells of a probe that ended with cause before any write.
func untouchedCells(cause string) verdictCells {
	return verdictCells{verdict: "invalid", cause: cause, restored: untouchedCell}
}

// baselineCause prefixes the baseline's own outcome kind, so a cause names which of the two
// runs reached the kind it carries.
const baselineCause = "baseline-"

// gradeBaseline runs the prepared selection once over the unmutated tree and answers the
// whole probe output when that run did not pass. A tree that is already red, that does not
// build, or that runs no test cannot support a verdict about the mutation, so the probe
// ends here: before the preserved copy, before the mutation, and before any write.
//
// The evidence after the rows is the baseline's own report, so its failures table names the
// tests that were red before the probe ran. The `ran` cell reads 0 and the focused exit is
// null, because the mutated run never started, and the `baseline` cell names the kind the
// caller has to fix first.
//
// The third result is the kind this run observed. The probe that follows prints it in its
// own `baseline` cell, so that cell reports one run's answer rather than a constant.
func gradeBaseline(root string, subject subject, mutation string, request testreport.Request) (string, int, testreport.OutcomeKind) {
	outcome, report, _ := testreport.Execute(root, request)
	if outcome.Kind == testreport.OutcomePassed {
		return "", 0, outcome.Kind
	}
	e := selected(subject, mutation, baselineCells(outcome), request)
	e.command = 1
	out, err := e.answer(report)
	if err != nil {
		return toon.RenderError(err) + "\n", 1, outcome.Kind
	}
	return out, 1, outcome.Kind
}

// baselineCells answers the cells of a probe whose baseline did not pass. The row is invalid,
// its cause and failures are the baseline's own, and the subject is untouched.
func baselineCells(outcome testreport.Outcome) verdictCells {
	cells := untouchedCells(baselineCause + string(outcome.Kind))
	cells.failed, cells.first, cells.baseline = outcome.FailedTests, firstFailure(outcome), string(outcome.Kind)
	return cells
}
