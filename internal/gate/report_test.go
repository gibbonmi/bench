package gate

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strings"
	"testing"
)

// writePhaseStream fills one phase's buffer through the writers the engine hands the
// scheduler, so every report test exercises the same path a real phase writes through.
func writePhaseStream(t *testing.T, streams *phaseStreams, phase, stdout, stderr string) {
	t.Helper()
	out, errOut, closeWriters := streams.open(Phase{Name: phase})
	if stdout != "" {
		if _, err := io.WriteString(out, stdout); err != nil {
			t.Fatalf("phase %s stdout write: %v", phase, err)
		}
	}
	if stderr != "" {
		if _, err := io.WriteString(errOut, stderr); err != nil {
			t.Fatalf("phase %s stderr write: %v", phase, err)
		}
	}
	closeWriters()
}

// tableRows answers the `line` cell of every row of a rendered failures block, so a test
// asserts the rows it cares about without restating the TOON escaping rules.
func tableRows(t *testing.T, stdout string) []string {
	t.Helper()
	var rows []string
	for _, line := range strings.Split(stdout, "\n") {
		if !strings.HasPrefix(line, "  ") {
			continue
		}
		cells := strings.SplitN(strings.TrimPrefix(line, "  "), ",", 2)
		if len(cells) != 2 {
			t.Fatalf("row %q is not a two-cell row", line)
		}
		rows = append(rows, unquoteCell(cells[1]))
	}
	return rows
}

// unquoteCell undoes the spec-TOON quoting the emitter applies to a cell that opens with
// a dash or carries a colon. The report owns which lines become rows; the encoder owns
// how a cell is spelled, and no test here asserts the second.
func unquoteCell(cell string) string {
	if len(cell) < 2 || !strings.HasPrefix(cell, `"`) || !strings.HasSuffix(cell, `"`) {
		return cell
	}
	return strings.NewReplacer(`\"`, `"`, `\\`, `\`).Replace(cell[1 : len(cell)-1])
}

// BG01, BG11: a red run's stdout is the failure table and the verdict, and nothing else.
// The verdict shares the rows' stream, so stderr carries no part of the answer.
func TestRedRunPrintsOnlyTheFailureTableAndTheVerdict(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	writePhaseStream(t, streams, "vet", "vet found a shadowed err\n", "")
	var stdout, stderr bytes.Buffer
	code := aggregateAndReport([]phaseResult{{Name: "vet", Argv: []string{"go", "vet", "./..."}, Code: 1}}, false, streams, &stdout, &stderr)

	want := "failures[1]{phase,line}:\n  vet,vet found a shadowed err\ngate: red\n"
	if stdout.String() != want {
		t.Errorf("red stdout = %q, want %q", stdout.String(), want)
	}
	if stderr.Len() != 0 {
		t.Errorf("red stderr = %q, want nothing", stderr.String())
	}
	if code != 1 {
		t.Errorf("red exit = %d, want 1", code)
	}
}

// BG02: a green phase beside a red one contributes no line. The filter is keyed on the
// phase, not on the run.
func TestGreenPhaseContributesNothingToARedRun(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	writePhaseStream(t, streams, "build", "build noise nobody asked for\n", "")
	writePhaseStream(t, streams, "vet", "vet finding\n", "")
	var stdout, stderr bytes.Buffer
	results := []phaseResult{
		{Name: "build", Argv: []string{"go", "build", "./..."}, Code: 0},
		{Name: "vet", Argv: []string{"go", "vet", "./..."}, Code: 1},
	}
	aggregateAndReport(results, false, streams, &stdout, &stderr)

	if strings.Contains(stdout.String(), "build noise nobody asked for") {
		t.Errorf("red stdout carried the green phase's stream: %q", stdout.String())
	}
	if got, want := tableRows(t, stdout.String()), []string{"vet finding"}; !slices.Equal(got, want) {
		t.Errorf("rows = %q, want %q", got, want)
	}
}

// BG06: a red phase that is not a Go test phase yields every non-empty line as a row. The
// blank line between the findings is padding and costs the reader no row.
func TestRedVetPhaseYieldsOneRowPerNonEmptyLine(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	writePhaseStream(t, streams, "vet", "first finding\n\nsecond finding\nthird finding\n", "")
	var stdout, stderr bytes.Buffer
	aggregateAndReport([]phaseResult{{Name: "vet", Argv: []string{"go", "vet", "./..."}, Code: 1}}, false, streams, &stdout, &stderr)

	want := []string{"first finding", "second finding", "third finding"}
	if got := tableRows(t, stdout.String()); !slices.Equal(got, want) {
		t.Errorf("rows = %q, want %q", got, want)
	}
	if !strings.HasPrefix(stdout.String(), "failures[3]{phase,line}:\n") {
		t.Errorf("red stdout header = %q, want a three-row failures header", stdout.String())
	}
}

// BG08: a red that said nothing still shows. Its exit code is otherwise the run's only
// evidence and stdout would carry none of it.
func TestRedPhaseWithAnEmptyStreamNamesItsExitCode(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	var stdout, stderr bytes.Buffer
	aggregateAndReport([]phaseResult{{Name: "shellcheck", Argv: []string{"shellcheck", "-S", "warning"}, Code: 1}}, false, streams, &stdout, &stderr)

	want := "failures[1]{phase,line}:\n  shellcheck,exit 1 with no output\ngate: red\n"
	if stdout.String() != want {
		t.Errorf("silent red stdout = %q, want %q", stdout.String(), want)
	}
}

// BG12: an interrupt is not a verdict. The stragglers say what the run was doing, and
// neither stream carries a gate line.
func TestCancelledRunPrintsStragglersAndNoVerdict(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	var stdout, stderr bytes.Buffer
	results := []phaseResult{{Name: "test", Argv: []string{"go", "test", "./..."}, Code: 130, Interrupted: true}}
	code := aggregateAndReport(results, true, streams, &stdout, &stderr)

	if code != 130 {
		t.Errorf("cancelled exit = %d, want 130", code)
	}
	if want := "gate: cancelled; still running: test\n"; stderr.String() != want {
		t.Errorf("cancelled stderr = %q, want %q", stderr.String(), want)
	}
	if stdout.Len() != 0 {
		t.Errorf("cancelled stdout = %q, want nothing", stdout.String())
	}
	if strings.Contains(stderr.String(), "gate: red") || strings.Contains(stderr.String(), "gate: green") {
		t.Errorf("cancelled run graded itself: %q", stderr.String())
	}
}

// BG15: a tool killed mid-line still contributes what it managed to say.
func TestRedStreamFlushesALastLineWithoutANewline(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	writePhaseStream(t, streams, "gofmt", "first line\ntruncated last line", "")
	var stdout, stderr bytes.Buffer
	aggregateAndReport([]phaseResult{{Name: "gofmt", Argv: []string{"gofmt", "-l", "."}, Code: 1}}, false, streams, &stdout, &stderr)

	want := []string{"first line", "truncated last line"}
	if got := tableRows(t, stdout.String()); !slices.Equal(got, want) {
		t.Errorf("rows = %q, want %q", got, want)
	}
}

// BG32: a red Go phase whose stream classified nothing — a race report — falls back to its
// own tail. An empty table would report the failure as no failure.
func TestRedGoPhaseWithNoClassifiedRowFallsBackToItsTail(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	var raw strings.Builder
	raw.WriteString("WARNING: DATA RACE\n")
	for i := 1; i <= 24; i++ {
		fmt.Fprintf(&raw, "race detail %d\n", i)
	}
	writePhaseStream(t, streams, "test", raw.String(), "")
	var stdout, stderr bytes.Buffer
	aggregateAndReport([]phaseResult{{Name: "test", Argv: []string{"go", "test", "-count=1", "./..."}, Code: 1}}, false, streams, &stdout, &stderr)

	rows := tableRows(t, stdout.String())
	if len(rows) != failureTailLines {
		t.Fatalf("race fallback yielded %d rows, want %d", len(rows), failureTailLines)
	}
	if rows[0] != "race detail 5" || rows[len(rows)-1] != "race detail 24" {
		t.Errorf("race fallback rows = %q, want the last twenty non-empty lines", rows)
	}
}

// A Go test phase whose stream does classify keeps the classifier's rows rather than its
// tail, so a failing test is one row and the rest of the run is not.
func TestRedGoPhaseUsesTheClassifiedRows(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	stream := "=== RUN   TestOne\n--- FAIL: TestOne (0.00s)\n    one_test.go:9: want 2 got 3\n=== RUN   TestTwo\n--- PASS: TestTwo (0.00s)\nFAIL\n"
	writePhaseStream(t, streams, "test", stream, "")
	var stdout, stderr bytes.Buffer
	aggregateAndReport([]phaseResult{{Name: "test", Argv: []string{"go", "test", "-count=1", "./..."}, Code: 1}}, false, streams, &stdout, &stderr)

	want := []string{"--- FAIL: TestOne (0.00s)", "one_test.go:9: want 2 got 3"}
	if got := tableRows(t, stdout.String()); !slices.Equal(got, want) {
		t.Errorf("rows = %q, want %q", got, want)
	}
}

// Two reds print in phase-table order, which is the order the operator reads the table in.
func TestTwoRedPhasesPrintInPhaseTableOrder(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	writePhaseStream(t, streams, "gofmt", "gofmt finding\n", "")
	writePhaseStream(t, streams, "vet", "vet finding\n", "")
	var stdout, stderr bytes.Buffer
	results := []phaseResult{
		{Name: "gofmt", Argv: []string{"gofmt", "-l", "."}, Code: 1},
		{Name: "vet", Argv: []string{"go", "vet", "./..."}, Code: 1},
	}
	aggregateAndReport(results, false, streams, &stdout, &stderr)

	want := "failures[2]{phase,line}:\n  gofmt,gofmt finding\n  vet,vet finding\ngate: red\n"
	if stdout.String() != want {
		t.Errorf("two-red stdout = %q, want %q", stdout.String(), want)
	}
}

// A tool that diagnoses on stderr is diagnosed all the same: one buffer holds both of a
// phase's streams, so the report never has to ask which one carried the failure.
func TestFailureArrivingOnStderrYieldsRows(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	writePhaseStream(t, streams, "prose", "", "prose refused a marker phrase\n")
	var stdout, stderr bytes.Buffer
	aggregateAndReport([]phaseResult{{Name: "prose", Argv: []string{"bench", "gate-prose"}, Code: 1}}, false, streams, &stdout, &stderr)

	if got, want := tableRows(t, stdout.String()), []string{"prose refused a marker phrase"}; !slices.Equal(got, want) {
		t.Errorf("rows = %q, want %q", got, want)
	}
}

// A phase a red need kept from launching is a consequence, not a failure. Naming it would
// read as a second thing to fix.
func TestPhaseSkippedByARedNeedContributesNoRow(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	writePhaseStream(t, streams, "build", "build finding\n", "")
	var stdout, stderr bytes.Buffer
	results := []phaseResult{
		{Name: "build", Argv: []string{"go", "build", "./..."}, Code: 1},
		{Name: "test", Argv: []string{"go", "test", "-count=1", "./..."}, SkippedBy: "build"},
	}
	aggregateAndReport(results, false, streams, &stdout, &stderr)

	if strings.Contains(stdout.String(), "test") {
		t.Errorf("the cascaded skip claimed a row: %q", stdout.String())
	}
	if got, want := tableRows(t, stdout.String()), []string{"build finding"}; !slices.Equal(got, want) {
		t.Errorf("rows = %q, want %q", got, want)
	}
}

// A red-reporting callback turns the run red on its own account, and its verdict still
// goes to stdout.
func TestARedReportCallbackTurnsTheRunRed(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	var stdout, stderr bytes.Buffer
	code := aggregateAndReport(nil, false, streams, &stdout, &stderr, func() bool { return true })

	if code != 1 {
		t.Errorf("callback red exit = %d, want 1", code)
	}
	if !strings.HasSuffix(stdout.String(), "gate: red\n") {
		t.Errorf("callback red stdout = %q, want it to end with the verdict", stdout.String())
	}
}

// A green run prints its verdict and no phase's stream. Ticket 04 adds the phase table
// above this line.
func TestGreenRunPrintsTheVerdictAlone(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	writePhaseStream(t, streams, "vet", "vet chatter\n", "")
	var stdout, stderr bytes.Buffer
	code := aggregateAndReport([]phaseResult{{Name: "vet", Argv: []string{"go", "vet", "./..."}, Code: 0}}, false, streams, &stdout, &stderr)

	if code != 0 {
		t.Errorf("green exit = %d, want 0", code)
	}
	if want := "gate: green\n"; stdout.String() != want {
		t.Errorf("green stdout = %q, want %q", stdout.String(), want)
	}
	if stderr.Len() != 0 {
		t.Errorf("green stderr = %q, want nothing", stderr.String())
	}
}
