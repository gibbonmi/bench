package testreport

// OutcomeKind names what one focused run did. The six kinds partition every path
// `Execute` can take: a run that reached a verdict gives passed, failed,
// build-failed, or no-test-run, and a run that never reached one gives refused or
// interrupted.
type OutcomeKind string

const (
	// OutcomePassed is a run that ran at least one test and observed no failure.
	OutcomePassed OutcomeKind = "passed"
	// OutcomeFailed is a run with at least one failing test row.
	OutcomeFailed OutcomeKind = "failed"
	// OutcomeBuildFailed is a package the run reported as failed with no failing test.
	OutcomeBuildFailed OutcomeKind = "build-failed"
	// OutcomeNoTestRun is a run that ran no test, including a run pattern that matched none.
	OutcomeNoTestRun OutcomeKind = "no-test-run"
	// OutcomeRefused is every refusal before or after the Go child that is not an interrupt.
	OutcomeRefused OutcomeKind = "refused"
	// OutcomeInterrupted is a run a cancel signal ended.
	OutcomeInterrupted OutcomeKind = "interrupted"
)

// Outcome is the typed verdict beside the rendered report. FailedTests counts the
// failing test rows, so a caller cites the count without parsing the report.
type Outcome struct {
	Kind        OutcomeKind
	FailedTests int
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

// outcome derives the kind from the same report the renderer consumed, so the verdict
// and the printed tables can never disagree. The failing rows are the renderer's own
// failure rows; a row with no test name is the package's own diagnostic, which is the
// build failure the kinds separate from a test failure.
func (r *report) outcome(full bool) Outcome {
	failed := 0
	for _, row := range r.failures(full) {
		if row[1] != "" {
			failed++
		}
	}
	if failed != 0 {
		return Outcome{Kind: OutcomeFailed, FailedTests: failed}
	}
	for _, status := range r.statuses {
		if status == "fail" {
			return Outcome{Kind: OutcomeBuildFailed}
		}
	}
	if !r.ranTest {
		return Outcome{Kind: OutcomeNoTestRun}
	}
	return Outcome{Kind: OutcomePassed}
}
