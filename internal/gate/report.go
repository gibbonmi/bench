package gate

// The engine's projection of a settled schedule: the red failure table, the green phase
// table, the verdict line, and the exit code. The execution engine lives in runner.go
// and knows nothing of the shape reported here.

import (
	"fmt"
	"io"
	"strings"

	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/testlines"
	"github.com/gibbonmi/bench/internal/toon"
)

// failureRowCap bounds what one phase contributes to the red table: at most twenty rows,
// and at most a twenty-line tail when a red Go phase's stream classified nothing. Both
// uses read the same number because a longer tail would only be cut again by the cap,
// and the more-row would then count lines the tail had already dropped. Twenty is the
// bound that lets a one-phase red fit the stop hook's thirty-line tail.
const failureRowCap = 20

// capabilityPhase files the skip reporter's red diagnoses in the table. It is not a
// scheduled phase, so its rows follow every phase's, and no phase may take this name.
const capabilityPhase = "capability"

// phaseTableCap is the largest phase count the green report still renders as a table.
// Seven phases cost eight lines with the header, and the skip-count line and the verdict
// take the run to ten. An eighth phase would pass that bound, so the table collapses to
// one row instead. The bound is on the report, not on the schedule: a longer table is a
// legal gate, and only its projection changes.
const phaseTableCap = 7

// aggregateAndReport is the one verdict tail every settled schedule reports through.
// It carries the failure table, any extra red-reporting checks (capability skips, the
// stripped-subject skip posture), and the `gate: red` / `gate: green` line. This is the
// operator's view of one command whichever schedule produced the results, so an edit to
// the reported shape lands everywhere at once.
//
// An extra check writes nothing itself. It answers its red diagnoses as rows and its own
// summary as one line, and this function decides where each belongs: the rows join the
// failure table, and the line prints on a green run only. A check that printed for itself
// would land ahead of the red table, which carries the whole of a red run's stdout.
//
// The verdict goes to stdout on both answers, so the rows and the word that grades them
// share one stream. The two operational refusals that print `gate: red` before any phase
// runs keep stderr: they are not this report.
//
// An interrupt is not a verdict, so a cancelled run publishes neither rows nor a gate
// line. Reporting one would grade phases that never got to answer. Naming the stragglers
// is not a verdict either; it only says what the run was doing.
func aggregateAndReport(results []phaseResult, cancelled bool, streams *phaseStreams, stdout, stderr io.Writer, redReports ...func() (rows []string, greenLine string, red bool)) int {
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
	// Every check answers before either shape prints, because a check can turn the run red
	// itself and the two shapes file its answers in different places.
	var capabilityRows, greenLines []string
	for _, report := range redReports {
		rows, greenLine, reportRed := report()
		capabilityRows = append(capabilityRows, rows...)
		if greenLine != "" {
			greenLines = append(greenLines, greenLine)
		}
		if reportRed {
			red = true
		}
	}
	if red {
		printFailures(results, capabilityRows, streams, stdout)
		fmt.Fprintln(stdout, "gate: red")
		return 1
	}
	// Green leads with the phase table: it is the run's whole answer, and the skip counts
	// read after it as a footnote about the host rather than as a preamble to it.
	printPhases(results, stdout)
	for _, line := range greenLines {
		fmt.Fprintln(stdout, line)
	}
	fmt.Fprintln(stdout, "gate: green")
	return 0
}

// printPhases emits the green run's one table, in phase-table order: what ran, how it
// settled, and how long it took. A green phase's own stream is not reported at all, so
// this table is the only account of a run that found nothing to say.
//
// Above phaseTableCap the table becomes one counted row. A reader of a green run wants
// the verdict and the shape of the run, and past that width the per-phase rows stop
// paying for the lines they cost.
func printPhases(results []phaseResult, stdout io.Writer) {
	if len(results) > phaseTableCap {
		fmt.Fprintf(stdout, "phases: %d/%d green\n", greenPhases(results), len(results))
		return
	}
	rows := make([][]any, 0, len(results))
	for _, result := range results {
		// The name passes the same control-byte filter a failure line does, so no cell
		// reaches the encoder unfiltered. The other two cells are the engine's own: a
		// fixed word and a count of milliseconds.
		rows = append(rows, []any{sanitize.Strip(result.Name), phaseVerdict(result), result.ElapsedMS})
	}
	block, err := toon.TableTyped("phases", []string{"phase", "verdict", "elapsed_ms"}, rows)
	if err != nil {
		// Unreachable through a filtered name and two engine-owned cells. It stays because
		// the encoder, not this file, owns which cells it refuses.
		return
	}
	fmt.Fprint(stdout, block)
}

// phaseVerdict answers one phase's cell. Both flavors of skip read alike: the cell says
// this phase produced no verdict, and whether an absent tool or an unsatisfied need
// caused that is not something a green run asks the reader to act on. Neither flavor is
// a red, so a run carrying either still settles green.
func phaseVerdict(result phaseResult) string {
	if result.Skipped || result.SkippedBy != "" {
		return "skipped"
	}
	return "green"
}

// greenPhases counts the phases that actually graded something. It is the numerator of
// the collapsed row, and it excludes a skip on purpose: the table this row replaces
// would have said `skipped` in that phase's cell, and folding skips into the green count
// would drop the one fact the table carried. The run is still green — the verdict line
// below the row says so — and `6/8 green` tells the reader that two phases never ran.
func greenPhases(results []phaseResult) int {
	count := 0
	for _, result := range results {
		if result.green() {
			count++
		}
	}
	return count
}

// printFailures emits the run's one failure table, in phase-table order. The skip
// reporter's rows follow, because `capability` is a filing name and not a phase the
// schedule ran.
//
// Only a phase's rows are capped. capabilityRows carries one row per diagnosis the
// reporter actually made, and a more-row pointing at the phase stream file would name a
// place those diagnoses were never written.
func printFailures(results []phaseResult, capabilityRows []string, streams *phaseStreams, stdout io.Writer) {
	rows := make([][]string, 0, len(results))
	for _, result := range results {
		rows = appendRows(rows, result.Name, boundRows(failureRows(result, streams), streams.path()))
	}
	rows = appendRows(rows, capabilityPhase, capabilityRows)
	block, err := toon.Table("failures", []string{"phase", "line"}, rows)
	if err != nil {
		// A cell spec-TOON cannot carry costs the table, never the verdict. appendRows
		// strips what the encoder refuses, so this answer is unreachable by a phase's own
		// output; it stays because the encoder, not this file, owns that contract.
		return
	}
	fmt.Fprint(stdout, block)
}

// appendRows files one phase's lines as table rows. Every cell passes the control-byte
// filter here, at the one point a line becomes a cell, so no row can reach the encoder
// unfiltered. The filter strips rather than escapes: the encoder escapes a backslash
// itself, and escaping twice would show a path with four.
func appendRows(rows [][]string, phase string, lines []string) [][]string {
	for _, line := range lines {
		rows = append(rows, []string{phase, sanitize.Strip(line)})
	}
	return rows
}

// boundRows caps one phase at failureRowCap rows and names what it dropped. A phase at
// exactly the cap dropped nothing and gets no extra row, so the reader is never told
// about an omission that did not happen.
func boundRows(lines []string, streamPath string) []string {
	if len(lines) <= failureRowCap {
		return lines
	}
	return append(append([]string(nil), lines[:failureRowCap]...), moreRow(len(lines)-failureRowCap, streamPath))
}

// moreRow says how many lines the cap dropped and where the rest is readable. An absent
// stream file gets its own wording rather than the same sentence with nothing after the
// colon: a reader sent to a place must be given one.
func moreRow(dropped int, streamPath string) string {
	if streamPath == "" {
		return fmt.Sprintf("+%d more lines (stream unavailable)", dropped)
	}
	return fmt.Sprintf("+%d more lines: %s", dropped, streamPath)
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
		// evidence and stdout carries none of it. A phase that never reached its program
		// says nothing because it never ran, and StartErr holds the reason it did not: a
		// bad working directory, an empty argv, an exec failure, or a deadlocked need. The
		// exit code names none of those, so the error's own text is the row.
		if result.StartErr != nil {
			return []string{result.StartErr.Error()}
		}
		return []string{fmt.Sprintf("exit %d with no output", result.Code)}
	}
	if goTestPhase(result.Argv) {
		if rows := testlines.FailureRows(lines); len(rows) > 0 {
			return rows
		}
		// A red Go phase whose stream classified nothing falls back to its own tail. An
		// empty table would report the failure as no failure.
		return lastLines(lines, failureRowCap)
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
