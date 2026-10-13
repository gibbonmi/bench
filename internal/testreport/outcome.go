package testreport

import (
	"errors"
	"strings"

	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/toon"
)

// OutcomeKind names what one focused run did. The six kinds partition every path
// `Execute` can take: a run that reached a verdict gives passed, failed,
// build-failed, or no-test-run, and a run that never reached one gives refused,
// interrupted, or build-failed for a failed run-binary build script.
type OutcomeKind string

const (
	// OutcomePassed is a run that ran at least one test and observed no failure.
	OutcomePassed OutcomeKind = "passed"
	// OutcomeFailed is a run with at least one failed test.
	OutcomeFailed OutcomeKind = "failed"
	// OutcomeBuildFailed is a package the run reported as failed with no failing test, or a
	// run-binary build script that failed.
	OutcomeBuildFailed OutcomeKind = "build-failed"
	// OutcomeNoTestRun is a run that ran no test, including a run pattern that matched none.
	OutcomeNoTestRun OutcomeKind = "no-test-run"
	// OutcomeRefused is every refusal before or after the Go child that is not an interrupt
	// or a failed run-binary build script.
	OutcomeRefused OutcomeKind = "refused"
	// OutcomeInterrupted is a run a cancel signal ended.
	OutcomeInterrupted OutcomeKind = "interrupted"
)

// Outcome is the typed verdict beside the rendered report. FailedTests counts the
// distinct failed tests, not the printed rows, so a caller cites the count without
// parsing the report. FailedNames names those Go tests in the order that the failures
// table prints them; a prose run counts its findings in FailedTests and names no test.
// Ran counts the distinct tests that emitted a run event, and it reads 0 when the run did
// not start, so a caller can never read no execution as evidence. BuildLine names the
// build-failed cause: the first diagnostic line of a failed package, or the first builder
// output line of a failed build script.
type Outcome struct {
	Kind        OutcomeKind
	FailedTests int
	FailedNames []string
	Ran         int
	BuildLine   string
}

// Request is one parsed focused-run selection. Its fields stay unexported, because
// the grammar owns the selection and a caller only carries the request from Prepare
// to Execute.
type Request struct {
	focused focusedRequest
}

// Command runs Go from root and renders one stable row for each observed result.
func Command(root string, args []string) (string, int) {
	request, line, code := Prepare(root, args)
	if line != "" {
		return line, code
	}
	_, out, code := Execute(root, request)
	return out, code
}

// Prepare parses the selection with `bench test`'s grammar. A refusal returns the
// usage line and its code; an accepted selection returns the request with an empty
// line, which is the caller's signal to Execute it.
func Prepare(root string, args []string) (Request, string, int) {
	focused, line, code := parseFocusedRequest(root, args)
	if line != "" {
		return Request{}, line + "\n", code
	}
	return Request{focused: focused}, "", 0
}

// Execute runs one prepared selection and returns the typed outcome beside the bytes
// and the exit `bench test` prints.
func Execute(root string, request Request) (Outcome, string, int) {
	return runFocusedRequest(root, request.focused)
}

func refusedOutcome(line string, code int) (Outcome, string, int) {
	return Outcome{Kind: OutcomeRefused}, line, code
}

func interruptedOutcome(line string, code int) (Outcome, string, int) {
	return Outcome{Kind: OutcomeInterrupted}, line, code
}

// selectionFailure answers a run whose run-binary selection failed. A failed build script
// reads build-failed, and every other selection failure is a refusal; both print the same
// refusal line and exit.
func selectionFailure(err error) (Outcome, string, int) {
	line := toon.Errorf("Bench executable selection failed", err.Error()) + "\n"
	var built runbinary.BuildError
	if !errors.As(err, &built) {
		return refusedOutcome(line, 1)
	}
	return Outcome{Kind: OutcomeBuildFailed, BuildLine: builderLine(built.Output)}, line, 1
}

// builderLine answers the first builder output line that names a Go source position, or
// else the first non-empty line, because a build script can fail before the compiler runs
// and a compiler heading precedes the error it names.
func builderLine(output string) string {
	first := ""
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, ".go:") {
			return line
		}
		if first == "" {
			first = line
		}
	}
	return first
}

// outcome derives the kind from the same report the renderer consumed, so the verdict
// and the printed tables can never disagree. The failed tests are the failures that the
// renderer splits into rows, counted once each; a failure with no test name is the
// package's own diagnostic, which is the build failure the kinds separate from a test
// failure. The build line is the first diagnostic line of those package failures in table
// order.
func (r *report) outcome() Outcome {
	ran := 0
	for _, tests := range r.ranTests {
		ran += len(tests)
	}
	var names []string
	buildLine := ""
	for _, f := range r.failed() {
		if f.test != "" {
			names = append(names, f.test)
		} else if buildLine == "" && len(f.lines) != 0 {
			buildLine = f.lines[0]
		}
	}
	if len(names) != 0 {
		return Outcome{Kind: OutcomeFailed, FailedTests: len(names), FailedNames: names, Ran: ran}
	}
	for _, status := range r.statuses {
		if status == "fail" {
			return Outcome{Kind: OutcomeBuildFailed, Ran: ran, BuildLine: buildLine}
		}
	}
	if ran == 0 {
		return Outcome{Kind: OutcomeNoTestRun}
	}
	return Outcome{Kind: OutcomePassed, Ran: ran}
}
