package gate

// The engine's projection of a settled schedule: the failure table, the verdict line,
// and the exit code. The execution engine lives in runner.go and knows nothing of the
// shape reported here.

import (
	"fmt"
	"io"
	"strings"

	"github.com/gibbonmi/bench/internal/testlines"
	"github.com/gibbonmi/bench/internal/toon"
)

// failureTailLines bounds the tail a red Go phase falls back to when its stream
// classified no failure row. Ticket 03 adds the general per-phase row cap.
const failureTailLines = 20

// aggregateAndReport is the one verdict tail every settled schedule reports through.
// It carries the failure table, any extra red-reporting checks (capability skips, the
// stripped-subject skip posture), and the `gate: red` / `gate: green` line. This is the
// operator's view of one command whichever schedule produced the results, so an edit to
// the reported shape lands everywhere at once.
//
// The verdict goes to stdout on both answers, so the rows and the word that grades them
// share one stream. The two operational refusals that print `gate: red` before any phase
// runs keep stderr: they are not this report.
//
// An interrupt is not a verdict, so a cancelled run publishes neither rows nor a gate
// line. Reporting one would grade phases that never got to answer. Naming the stragglers
// is not a verdict either; it only says what the run was doing.
func aggregateAndReport(results []phaseResult, cancelled bool, streams *phaseStreams, stdout, stderr io.Writer, redReports ...func() bool) int {
	if cancelled {
		reportStragglers(results, stderr)
		return 130
	}
	red := false
	for _, result := range results {
		if result.Code != 0 {
			red = true
		}
	}
	// A red report prints as it diagnoses, so its lines precede the table. Ticket 03
	// makes each such diagnosis a row under phase `capability` instead.
	for _, report := range redReports {
		if report() {
			red = true
		}
	}
	if red {
		printFailures(results, streams, stdout)
		fmt.Fprintln(stdout, "gate: red")
		return 1
	}
	// Ticket 04 adds the phase table and the capability-skips line above this verdict.
	fmt.Fprintln(stdout, "gate: green")
	return 0
}

// printFailures emits the run's one failure table, in phase-table order.
func printFailures(results []phaseResult, streams *phaseStreams, stdout io.Writer) {
	rows := make([][]string, 0, len(results))
	for _, result := range results {
		for _, line := range failureRows(result, streams) {
			rows = append(rows, []string{result.Name, line})
		}
	}
	block, err := toon.Table("failures", []string{"phase", "line"}, rows)
	if err != nil {
		// A cell spec-TOON cannot carry costs the table, never the verdict. Ticket 03
		// adds the control-byte filter that keeps every cell representable.
		return
	}
	fmt.Fprint(stdout, block)
}

// failureRows answers one phase's rows. A green phase contributes none, and neither does
// a phase a red need kept from launching: that phase is a consequence, not a failure, and
// naming it would read as a second thing to fix.
func failureRows(result phaseResult, streams *phaseStreams) []string {
	if result.Code == 0 || result.SkippedBy != "" {
		return nil
	}
	lines := nonEmptyLines(streams.lines(result.Name))
	if len(lines) == 0 {
		// A red that said nothing still has to show, or its exit code is the run's only
		// evidence and stdout carries none of it.
		return []string{fmt.Sprintf("exit %d with no output", result.Code)}
	}
	if goTestPhase(result.Argv) {
		if rows := testlines.FailureRows(lines); len(rows) > 0 {
			return rows
		}
		// A red Go phase whose stream classified nothing — a race report, say — falls
		// back to its own tail. An empty table would report the failure as no failure.
		return lastLines(lines, failureTailLines)
	}
	return lines
}

// goTestPhase reports whether a phase runs `go test`, the one phase kind whose red stream
// has a classifier. Every other phase's red stream is failure rows line by line.
func goTestPhase(argv []string) bool {
	return len(argv) >= 2 && argv[0] == "go" && argv[1] == "test"
}

// nonEmptyLines drops the blank lines a tool pads its output with. A blank row carries
// nothing and still costs the reader a line.
func nonEmptyLines(lines []string) []string {
	kept := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		kept = append(kept, line)
	}
	return kept
}

func lastLines(lines []string, n int) []string {
	if len(lines) <= n {
		return lines
	}
	return lines[len(lines)-n:]
}
