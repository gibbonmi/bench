package testreport

import (
	"context"
	"os"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/prose"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/toon"
)

const proseCheckName = "prose"

// ranNothingTitle is the error title of a named check that reached no evidence that it ran.
const ranNothingTitle = "named check ran nothing"

func unknownCheck(check string) string {
	return "unknown check: " + check + "\n" + namedCheckInventory()
}

func namedCheckInventory() string {
	checks := namedChecks()
	return "checks:\n  " + strings.Join(checks, "\n  ")
}

func namedChecks() []string {
	return append(registry.Names(registry.Dev), gate.SystemPhaseName, proseCheckName)
}

func isNamedCheck(check string) bool {
	for _, name := range namedChecks() {
		if name == check {
			return true
		}
	}
	return false
}

// namedCheckKind is the one owner of the kind cell: the system suite, the prose
// grader, or a conformance scope of the root conformance test.
func namedCheckKind(check string) string {
	switch check {
	case gate.SystemPhaseName:
		return "system"
	case proseCheckName:
		return "prose"
	}
	return "conformance"
}

// checkRow renders the `check` block that leads each named-check result with a verdict.
func checkRow(check string, testsRun, subjects int) (string, error) {
	return toon.TableTyped("check", []string{"name", "kind", "tests_run", "subjects"}, [][]any{{check, namedCheckKind(check), testsRun, subjects}})
}

// runNamedCheck leads a Go-backed verdict with the check row. A refusal, an interrupt,
// and the `--run` refusal reached no check verdict, so they keep their bytes. A run with
// no failure and no run event is no evidence that the check ran, so it exits 1.
func runNamedCheck(ctx context.Context, root string, request focusedRequest, selection *runbinary.Selection) (Outcome, string, int) {
	outcome, out, code := runGoBackedCheck(ctx, root, request, selection)
	if outcome.Kind == OutcomeRefused || outcome.Kind == OutcomeInterrupted || request.run != "" && outcome.Kind == OutcomeNoTestRun {
		return outcome, out, code
	}
	row, err := checkRow(request.check, outcome.Ran, 0)
	if err != nil {
		return refusedOutcome(toon.RenderError(err)+"\n", 1)
	}
	if outcome.Kind == OutcomeNoTestRun {
		return outcome, row + toon.Errorf(ranNothingTitle, "no test emitted a run event") + "\n", 1
	}
	return outcome, row + out, code
}

func runGoBackedCheck(ctx context.Context, root string, request focusedRequest, selection *runbinary.Selection) (Outcome, string, int) {
	if request.check == gate.SystemPhaseName {
		return runSystemCheck(ctx, root, request, selection)
	}
	argv := focusedTestArgv("./internal/conformance", "-run", namedCheckRunPattern())
	env, err := conformanceEnvironment(os.Environ(), root, request.check, selection)
	if err != nil {
		return refusedOutcome(toon.Errorf("go test failed to start", err.Error())+"\n", 1)
	}
	return runGoTest(ctx, selection.SourceRoot, request, argv, env)
}

// runProseCheck grades sentences rather than a Go test, so its check row counts graded
// subjects and runs no test. Each finding or grader refusal line follows the row as a
// failure, and a grade of no subject is no evidence that the check ran.
func runProseCheck(root string, full bool) (Outcome, string, int) {
	grade := prose.GradeTree(root)
	row, err := checkRow(proseCheckName, 0, len(grade.Subjects))
	if err != nil {
		return refusedOutcome(toon.RenderError(err)+"\n", 1)
	}
	if len(grade.Findings) > 0 {
		return Outcome{Kind: OutcomeFailed, FailedTests: len(grade.Findings)}, row + strings.Join(grade.Findings, "\n") + "\n", 1
	}
	if len(grade.Subjects) == 0 {
		return Outcome{Kind: OutcomeNoTestRun}, row + toon.Errorf(ranNothingTitle, "the prose grader graded no subject") + "\n", 1
	}
	if !full {
		return Outcome{Kind: OutcomePassed}, row, 0
	}
	subjects := make([][]string, 0, len(grade.Subjects))
	for _, path := range slices.Sorted(slices.Values(grade.Subjects)) {
		subjects = append(subjects, []string{path})
	}
	table, err := toon.Table("subjects", []string{"path"}, subjects)
	if err != nil {
		return refusedOutcome(toon.RenderError(err)+"\n", 1)
	}
	return Outcome{Kind: OutcomePassed}, row + table, 0
}

// runSystemCheck runs the gate's system phase as a focused run. It reads the phase's
// operands and environment from the gate's producer, and it sets no conformance
// variable, because the system suite is a build-tagged package rather than a
// conformance scope.
func runSystemCheck(ctx context.Context, root string, request focusedRequest, selection *runbinary.Selection) (Outcome, string, int) {
	operands, suiteEnv := gate.SystemSuite(root)
	env, err := selectedRunEnvironment(os.Environ(), selection)
	if err != nil {
		return refusedOutcome(toon.Errorf("go test failed to start", err.Error())+"\n", 1)
	}
	return runGoTest(ctx, root, request, focusedTestArgv(operands...), append(env, suiteEnv...))
}
