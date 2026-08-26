package gate

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
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

// tableCells answers every row of a rendered failures block as a phase cell and a line
// cell, so a test asserts the rows it cares about without restating the TOON escaping
// rules. It is the one parser the row helpers below read through.
func tableCells(t *testing.T, stdout string) [][2]string {
	t.Helper()
	var rows [][2]string
	for _, line := range strings.Split(stdout, "\n") {
		if !strings.HasPrefix(line, "  ") {
			continue
		}
		cells := strings.SplitN(strings.TrimPrefix(line, "  "), ",", 2)
		if len(cells) != 2 {
			t.Fatalf("row %q is not a two-cell row", line)
		}
		rows = append(rows, [2]string{unquoteCell(cells[0]), unquoteCell(cells[1])})
	}
	return rows
}

// tableRows answers the `line` cell of every row.
func tableRows(t *testing.T, stdout string) []string {
	t.Helper()
	var lines []string
	for _, row := range tableCells(t, stdout) {
		lines = append(lines, row[1])
	}
	return lines
}

// rowsForPhase answers the `line` cell of every row filed under one phase, which is how a
// test asks what a phase contributed without depending on what the others did.
func rowsForPhase(t *testing.T, stdout, phase string) []string {
	t.Helper()
	var lines []string
	for _, row := range tableCells(t, stdout) {
		if row[0] == phase {
			lines = append(lines, row[1])
		}
	}
	return lines
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
	if len(rows) != failureRowCap {
		t.Fatalf("race fallback yielded %d rows, want %d", len(rows), failureRowCap)
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

// BG09: a phase a red need kept from launching is a consequence, not a failure. Naming it
// would read as a second thing to fix.
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
	code := aggregateAndReport(nil, false, streams, &stdout, &stderr, func(io.Writer) ([]string, bool) { return nil, true })

	if code != 1 {
		t.Errorf("callback red exit = %d, want 1", code)
	}
	if !strings.HasSuffix(stdout.String(), "gate: red\n") {
		t.Errorf("callback red stdout = %q, want it to end with the verdict", stdout.String())
	}
}

// BG07: a phase past the cap prints exactly the cap's worth of rows and one row that
// counts what it dropped. Without the cap a single runaway tool owns the whole answer, so
// the count of rows, not only their content, is the property under test.
func TestRedPhasePastTheCapPrintsTwentyRowsAndACountOfTheRest(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	var raw strings.Builder
	for i := 1; i <= 50; i++ {
		fmt.Fprintf(&raw, "finding %d\n", i)
	}
	writePhaseStream(t, streams, "vet", raw.String(), "")
	var stdout, stderr bytes.Buffer
	aggregateAndReport([]phaseResult{{Name: "vet", Argv: []string{"go", "vet", "./..."}, Code: 1}}, false, streams, &stdout, &stderr)

	rows := rowsForPhase(t, stdout.String(), "vet")
	if len(rows) != failureRowCap+1 {
		t.Fatalf("fifty findings yielded %d rows, want %d", len(rows), failureRowCap+1)
	}
	if rows[0] != "finding 1" || rows[failureRowCap-1] != "finding 20" {
		t.Errorf("capped rows = %q, want the first twenty findings", rows[:failureRowCap])
	}
	// This run retained no stream file, so the row counts the drop and says where the
	// rest is not. Ticket 05 opens the file and the same row names it.
	if want := "+30 more lines (stream unavailable)"; rows[failureRowCap] != want {
		t.Errorf("more-row = %q, want %q", rows[failureRowCap], want)
	}
}

// The same row, on a run that retained its stream: the cap names the file where the rest
// of the phase is readable. This is the whole point of the bound, so the path in the row
// must be the file the buffer actually wrote to.
func TestRedPhasePastTheCapNamesTheRetainedStreamFile(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	file := openTestStreamFile(t)
	streams.retain(file)
	var raw strings.Builder
	for i := 1; i <= 50; i++ {
		fmt.Fprintf(&raw, "finding %d\n", i)
	}
	writePhaseStream(t, streams, "vet", raw.String(), "")
	var stdout, stderr bytes.Buffer
	aggregateAndReport([]phaseResult{{Name: "vet", Argv: []string{"go", "vet", "./..."}, Code: 1}}, false, streams, &stdout, &stderr)

	rows := rowsForPhase(t, stdout.String(), "vet")
	if want := "+30 more lines: " + file.Name(); rows[len(rows)-1] != want {
		t.Errorf("more-row = %q, want %q", rows[len(rows)-1], want)
	}
}

// A phase at exactly the cap dropped nothing, so it gets no more-row. An off-by-one that
// counts the cap itself as an omission reports a loss that did not happen.
func TestRedPhaseAtExactlyTheCapPrintsNoMoreRow(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	var raw strings.Builder
	for i := 1; i <= failureRowCap; i++ {
		fmt.Fprintf(&raw, "finding %d\n", i)
	}
	writePhaseStream(t, streams, "vet", raw.String(), "")
	var stdout, stderr bytes.Buffer
	aggregateAndReport([]phaseResult{{Name: "vet", Argv: []string{"go", "vet", "./..."}, Code: 1}}, false, streams, &stdout, &stderr)

	rows := rowsForPhase(t, stdout.String(), "vet")
	if len(rows) != failureRowCap {
		t.Fatalf("twenty findings yielded %d rows, want %d", len(rows), failureRowCap)
	}
	if strings.Contains(stdout.String(), "more lines") {
		t.Errorf("a phase at the cap claimed dropped lines: %q", stdout.String())
	}
}

// BG07, second half: when the run does retain a stream file, the more-row names it, so a
// reader is sent to the complete output rather than told it is gone. Ticket 05 wires the
// path; this pins the row the wiring must produce.
func TestMoreRowNamesTheStreamFileWhenTheRunRetainedOne(t *testing.T) {
	lines := make([]string, failureRowCap+3)
	for i := range lines {
		lines[i] = fmt.Sprintf("finding %d", i+1)
	}
	bounded := boundRows(lines, "/tmp/.logs/gate-42.out")

	if len(bounded) != failureRowCap+1 {
		t.Fatalf("bounded %d lines to %d rows, want %d", len(lines), len(bounded), failureRowCap+1)
	}
	if want := "+3 more lines: /tmp/.logs/gate-42.out"; bounded[failureRowCap] != want {
		t.Errorf("more-row = %q, want %q", bounded[failureRowCap], want)
	}
}

// BG10: an environment skip is a check the oracle asked for and did not get, so it enters
// the table as its own row under `capability`. The reporter is the real one and the log is
// a real fixture, so a drift between what the reporter diagnoses and what the table files
// turns this red.
func TestEnvironmentSkipEntersARowUnderTheCapabilityPhase(t *testing.T) {
	t.Setenv(requireCapabilitiesEnv, "")
	path := writeSkipLog(t, capability.Skip{Kind: capability.KindEnvironment, Name: "TestRootConformance", Reason: "BENCH_CONFORMANCE_ROOT not set"})
	streams := newPhaseStreams(io.Discard)
	var stdout, stderr bytes.Buffer
	code := aggregateAndReport(nil, false, streams, &stdout, &stderr, func(w io.Writer) ([]string, bool) {
		return reportCapabilitySkips(path, w)
	})

	if code != 1 {
		t.Errorf("environment skip exit = %d, want 1", code)
	}
	want := []string{"TestRootConformance: BENCH_CONFORMANCE_ROOT not set"}
	if got := rowsForPhase(t, stdout.String(), capabilityPhase); !slices.Equal(got, want) {
		t.Errorf("capability rows = %q, want %q", got, want)
	}
}

// BG33: the two strict-mode diagnoses are rows as well. A diagnosis left on stderr leaves
// the reader a red table that names none of what went wrong.
func TestStrictDiagnosesEachEnterARowUnderTheCapabilityPhase(t *testing.T) {
	report := func(t *testing.T, path string) string {
		t.Helper()
		streams := newPhaseStreams(io.Discard)
		var stdout, stderr bytes.Buffer
		if code := aggregateAndReport(nil, false, streams, &stdout, &stderr, func(w io.Writer) ([]string, bool) {
			return reportCapabilitySkips(path, w)
		}); code != 1 {
			t.Fatalf("strict exit = %d, want 1", code)
		}
		return stdout.String()
	}

	t.Run("capability skip", func(t *testing.T) {
		t.Setenv(requireCapabilitiesEnv, "1")
		path := writeSkipLog(t, capability.Skip{Kind: capability.KindCapability, Class: capability.Fifo, Name: "TestFifoRefusal", Reason: "requires a host fifo"})
		want := []string{"capability skips are fatal under " + requireCapabilitiesEnv + "=1: fifo"}
		if got := rowsForPhase(t, report(t, path), capabilityPhase); !slices.Equal(got, want) {
			t.Errorf("capability rows = %q, want %q", got, want)
		}
	})

	t.Run("unreadable log", func(t *testing.T) {
		t.Setenv(requireCapabilitiesEnv, "1")
		// A directory is a log the reader cannot read and did not fail to find, which is
		// the case an absent log must not be confused with.
		path := t.TempDir()
		rows := rowsForPhase(t, report(t, path), capabilityPhase)
		if len(rows) != 1 {
			t.Fatalf("capability rows = %q, want exactly one", rows)
		}
		for _, want := range []string{path, "is unreadable, so the counts above prove nothing"} {
			if !strings.Contains(rows[0], want) {
				t.Errorf("unreadable-log row %q does not name %q", rows[0], want)
			}
		}
	})
}

// A capability skip outside strict mode is a fact about the host, not something the reader
// is asked to fix, so it enters no row even when the run is red for another reason.
func TestCapabilitySkipOutsideStrictModeEntersNoRow(t *testing.T) {
	t.Setenv(requireCapabilitiesEnv, "")
	path := writeSkipLog(t, capability.Skip{Kind: capability.KindCapability, Class: capability.Fifo, Name: "TestFifoRefusal", Reason: "requires a host fifo"})
	streams := newPhaseStreams(io.Discard)
	writePhaseStream(t, streams, "vet", "vet finding\n", "")
	var stdout, stderr bytes.Buffer
	aggregateAndReport([]phaseResult{{Name: "vet", Argv: []string{"go", "vet", "./..."}, Code: 1}}, false, streams, &stdout, &stderr, func(w io.Writer) ([]string, bool) {
		return reportCapabilitySkips(path, w)
	})

	if got := rowsForPhase(t, stdout.String(), capabilityPhase); len(got) != 0 {
		t.Errorf("capability rows = %q, want none", got)
	}
	if got, want := rowsForPhase(t, stdout.String(), "vet"), []string{"vet finding"}; !slices.Equal(got, want) {
		t.Errorf("rows = %q, want %q", got, want)
	}
}

// BG14: a cell the encoder refuses costs the whole table, so a control byte is removed
// before the line becomes a cell. The ESC is the byte that matters: it drives the terminal
// that prints the row, and a tool that colors its output emits one on every finding.
func TestFailureLineWithAnEscapeByteRendersWithoutIt(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	writePhaseStream(t, streams, "vet", "vet found "+string(rune(0x1b))+"[31ma shadowed err"+string(rune(0x1b))+"[0m\n", "")
	var stdout, stderr bytes.Buffer
	aggregateAndReport([]phaseResult{{Name: "vet", Argv: []string{"go", "vet", "./..."}, Code: 1}}, false, streams, &stdout, &stderr)

	want := []string{"vet found [31ma shadowed err[0m"}
	if got := tableRows(t, stdout.String()); !slices.Equal(got, want) {
		t.Errorf("rows = %q, want %q", got, want)
	}
	if strings.ContainsRune(stdout.String(), rune(0x1b)) {
		t.Errorf("the table carried a raw escape byte: %q", stdout.String())
	}
}

// A backslash reaches the reader once. The filter strips rather than escapes for exactly
// this reason: the encoder escapes the cell itself, and a filter that escaped first would
// show a Windows path with four backslashes.
func TestFailureLineWithOneBackslashRendersWithOneBackslash(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	writePhaseStream(t, streams, "gofmt", `needs a rewrite: src\main.go`+"\n", "")
	var stdout, stderr bytes.Buffer
	aggregateAndReport([]phaseResult{{Name: "gofmt", Argv: []string{"gofmt", "-l", "."}, Code: 1}}, false, streams, &stdout, &stderr)

	want := []string{`needs a rewrite: src\main.go`}
	if got := tableRows(t, stdout.String()); !slices.Equal(got, want) {
		t.Errorf("decoded rows = %q, want %q", got, want)
	}
}

// A green run prints its phase table and its verdict, and no phase's stream.
func TestGreenRunPrintsTheTableAndTheVerdict(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	writePhaseStream(t, streams, "vet", "vet chatter\n", "")
	var stdout, stderr bytes.Buffer
	code := aggregateAndReport([]phaseResult{{Name: "vet", Argv: []string{"go", "vet", "./..."}, Code: 0, ElapsedMS: 12}}, false, streams, &stdout, &stderr)

	if code != 0 {
		t.Errorf("green exit = %d, want 0", code)
	}
	if want := "phases[1]{phase,verdict,elapsed_ms}:\n  vet,green,12\ngate: green\n"; stdout.String() != want {
		t.Errorf("green stdout = %q, want %q", stdout.String(), want)
	}
	if stderr.Len() != 0 {
		t.Errorf("green stderr = %q, want nothing", stderr.String())
	}
}

// greenPhaseResults builds n settled-green phases with a distinct elapsed time each, so a
// test that pins the table also pins that each row carries its own phase's number.
func greenPhaseResults(n int) []phaseResult {
	results := make([]phaseResult, 0, n)
	for i := 1; i <= n; i++ {
		results = append(results, phaseResult{Name: fmt.Sprintf("phase%d", i), Argv: []string{"go", "build"}, ElapsedMS: int64(i)})
	}
	return results
}

// stdoutLines answers the printed lines of a report, which is the unit the bound counts.
func stdoutLines(stdout string) []string {
	return strings.Split(strings.TrimSuffix(stdout, "\n"), "\n")
}

// BG16, BG30: a green six-phase run is the table, one capability-skips line, and the
// verdict — nine lines, with no phase's own stream among them. The line count is asserted
// beside the bytes because the bound, not only the wording, is the contract: a build that
// added the table and kept relaying the phase streams would satisfy BG16 alone.
func TestGreenSixPhaseRunPrintsNineLines(t *testing.T) {
	t.Setenv(requireCapabilitiesEnv, "")
	// An absent log is a run whose phases skipped nothing, which is the zero-class case.
	path := filepath.Join(t.TempDir(), "absent.log")
	streams := newPhaseStreams(io.Discard)
	results := greenPhaseResults(6)
	for _, result := range results {
		writePhaseStream(t, streams, result.Name, result.Name+" chatter\n", "")
	}
	var stdout, stderr bytes.Buffer
	code := aggregateAndReport(results, false, streams, &stdout, &stderr, func(w io.Writer) ([]string, bool) {
		return reportCapabilitySkips(path, w)
	})

	if code != 0 {
		t.Errorf("green exit = %d, want 0", code)
	}
	want := "phases[6]{phase,verdict,elapsed_ms}:\n" +
		"  phase1,green,1\n" +
		"  phase2,green,2\n" +
		"  phase3,green,3\n" +
		"  phase4,green,4\n" +
		"  phase5,green,5\n" +
		"  phase6,green,6\n" +
		"capability-skips: 0 (capability=0 environment=0)\n" +
		"gate: green\n"
	if stdout.String() != want {
		t.Errorf("green stdout =\n%q\nwant\n%q", stdout.String(), want)
	}
	lines := stdoutLines(stdout.String())
	if len(lines) != 9 {
		t.Errorf("green stdout = %d lines, want 9: %q", len(lines), lines)
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "[") {
			t.Errorf("green stdout carried a phase-prefixed line: %q", line)
		}
		if strings.HasPrefix(line, "phase ") {
			t.Errorf("green stdout carried a per-phase summary line: %q", line)
		}
	}
	if strings.Contains(stdout.String(), "chatter") {
		t.Errorf("green stdout carried a phase's stream: %q", stdout.String())
	}
}

// The table holds at exactly the cap. Eight lines of table, the skip-count line, and the
// verdict are the ten the bound allows, so this is the widest run that still names its
// phases. It pins the threshold from below; the eight-phase test pins it from above.
func TestGreenSevenPhaseRunStillPrintsTheTable(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	var stdout, stderr bytes.Buffer
	aggregateAndReport(greenPhaseResults(7), false, streams, &stdout, &stderr)

	if !strings.HasPrefix(stdout.String(), "phases[7]{phase,verdict,elapsed_ms}:\n") {
		t.Errorf("seven-phase stdout = %q, want a seven-row phases table", stdout.String())
	}
	if got := len(stdoutLines(stdout.String())); got != 9 {
		t.Errorf("seven-phase stdout = %d lines, want 9", got)
	}
}

// BG17: above seven phases the table collapses to one counted row, so a wide gate cannot
// spend the bound on rows. Every phase graded green here, so the count is 8/8.
func TestGreenEightPhaseRunCollapsesToOneRow(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	var stdout, stderr bytes.Buffer
	code := aggregateAndReport(greenPhaseResults(8), false, streams, &stdout, &stderr)

	if code != 0 {
		t.Errorf("collapsed exit = %d, want 0", code)
	}
	if want := "phases: 8/8 green\ngate: green\n"; stdout.String() != want {
		t.Errorf("eight-phase stdout = %q, want %q", stdout.String(), want)
	}
}

// The collapsed row counts what graded, not what the run contained. A skip leaves the run
// green and still costs the numerator, because the table this row stands in for would
// have shown that phase as `skipped`.
func TestCollapsedRowCountsOnlyTheGradedPhases(t *testing.T) {
	results := greenPhaseResults(8)
	results[2].Skipped = true
	results[3].SkippedBy = "phase3"
	streams := newPhaseStreams(io.Discard)
	var stdout, stderr bytes.Buffer
	code := aggregateAndReport(results, false, streams, &stdout, &stderr)

	if code != 0 {
		t.Errorf("collapsed exit = %d, want 0", code)
	}
	if want := "phases: 6/8 green\ngate: green\n"; stdout.String() != want {
		t.Errorf("collapsed stdout = %q, want %q", stdout.String(), want)
	}
}

// BG18: a phase its own absent tool skipped, and a phase a skip kept from launching, both
// read `skipped`. Neither is a red, so the run stays green. A phase that never launched
// carries no elapsed time, and the zero is the honest cell for it.
func TestSkippedPhaseReadsSkippedAndKeepsTheRunGreen(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	results := []phaseResult{
		{Name: "build", Argv: []string{"go", "build", "./..."}, ElapsedMS: 40},
		{Name: "shellcheck", Argv: []string{"shellcheck", "-S", "warning"}, Skipped: true},
		{Name: "shellfmt", Argv: []string{"shfmt", "-l", "."}, SkippedBy: "shellcheck"},
	}
	var stdout, stderr bytes.Buffer
	code := aggregateAndReport(results, false, streams, &stdout, &stderr)

	if code != 0 {
		t.Errorf("skipped-phase exit = %d, want 0", code)
	}
	want := "phases[3]{phase,verdict,elapsed_ms}:\n" +
		"  build,green,40\n" +
		"  shellcheck,skipped,0\n" +
		"  shellfmt,skipped,0\n" +
		"gate: green\n"
	if stdout.String() != want {
		t.Errorf("skipped-phase stdout = %q, want %q", stdout.String(), want)
	}
}

// BG19: every class folds into the one capability-skips line, in the package's declared
// order. Three skips in two classes used to cost three lines; they cost one now, and the
// line still names both classes so the reader learns what the host could not run.
func TestCapabilitySkipsInTwoClassesPrintOneLine(t *testing.T) {
	t.Setenv(requireCapabilitiesEnv, "")
	path := writeSkipLog(t,
		capability.Skip{Kind: capability.KindCapability, Class: capability.Privilege, Name: "TestDropPrivilege", Reason: "requires a second uid"},
		capability.Skip{Kind: capability.KindCapability, Class: capability.Fifo, Name: "TestFifoRefusal", Reason: "requires a host fifo"},
		capability.Skip{Kind: capability.KindCapability, Class: capability.Fifo, Name: "TestFifoTimeout", Reason: "requires a host fifo"},
	)
	streams := newPhaseStreams(io.Discard)
	var stdout, stderr bytes.Buffer
	code := aggregateAndReport(greenPhaseResults(1), false, streams, &stdout, &stderr, func(w io.Writer) ([]string, bool) {
		return reportCapabilitySkips(path, w)
	})

	if code != 0 {
		t.Errorf("capability-skip exit = %d, want 0", code)
	}
	want := "phases[1]{phase,verdict,elapsed_ms}:\n" +
		"  phase1,green,1\n" +
		"capability-skips: 3 (capability=3 environment=0; fifo=2 privilege=1)\n" +
		"gate: green\n"
	if stdout.String() != want {
		t.Errorf("capability-skip stdout = %q, want %q", stdout.String(), want)
	}
}
