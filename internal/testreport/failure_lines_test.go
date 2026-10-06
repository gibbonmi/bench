package testreport

import (
	"strings"
	"testing"
)

// threeLineFailureEvents is a failing stream of package `canned` in which one test emits
// three diagnostic lines, two in one output event and one in a second.
var threeLineFailureEvents = []string{
	`{"Action":"run","Package":"canned","Test":"TestCanned"}`,
	`{"Action":"output","Package":"canned","Test":"TestCanned","Output":"    alpha one\n    beta two\n"}`,
	`{"Action":"output","Package":"canned","Test":"TestCanned","Output":"    gamma three\n"}`,
	`{"Action":"fail","Package":"canned","Test":"TestCanned","Elapsed":0.01}`,
	`{"Action":"fail","Package":"canned","Elapsed":0.5}`,
}

const failedPackageRow = "packages[1]{package,status,elapsed_ms,tests_run}:\n  canned,fail,500,1\n"

// TestFailuresRowCountsLines grades that the default row shows the first line and counts
// every diagnostic line of the test.
func TestFailuresRowCountsLines(t *testing.T) {
	output, code := commandOverEvents(t, t.TempDir(), []string{"--package", "./..."}, threeLineFailureEvents, 1)
	want := failedPackageRow + "failures[1]{package,test,line,lines}:\n  canned,TestCanned,alpha one,3\nskips[0]{package,test,reason}:\n"
	if output != want || code != 1 {
		t.Fatalf("Command = (%d, %q), want (1, %q)", code, output, want)
	}
}

// TestFailuresRowWithNoDiagnosticCountsZero grades that the placeholder cell is not a
// counted line.
func TestFailuresRowWithNoDiagnosticCountsZero(t *testing.T) {
	events := []string{threeLineFailureEvents[0], threeLineFailureEvents[3], threeLineFailureEvents[4]}
	for _, args := range [][]string{{"--package", "./..."}, {"--package", "./...", "--full"}} {
		output, code := commandOverEvents(t, t.TempDir(), args, events, 1)
		want := failedPackageRow + "failures[1]{package,test,line,lines}:\n  canned,TestCanned,no diagnostic emitted,0\nskips[0]{package,test,reason}:\n"
		if output != want || code != 1 {
			t.Fatalf("Command %v = (%d, %q), want (1, %q)", args, code, output, want)
		}
	}
}

// TestFullFailuresPrintOneRowPerLine grades that `--full` prints each diagnostic line as
// its own row in emitted order, and that each row carries the count of the test.
func TestFullFailuresPrintOneRowPerLine(t *testing.T) {
	output, code := commandOverEvents(t, t.TempDir(), []string{"--package", "./...", "--full"}, threeLineFailureEvents, 1)
	want := failedPackageRow + "failures[3]{package,test,line,lines}:\n  canned,TestCanned,alpha one,3\n  canned,TestCanned,beta two,3\n  canned,TestCanned,gamma three,3\nskips[0]{package,test,reason}:\n"
	if output != want || code != 1 {
		t.Fatalf("Command = (%d, %q), want (1, %q)", code, output, want)
	}
}

// TestFullFailuresHoldNoJoinedCell grades that no `--full` failures cell joins two lines.
func TestFullFailuresHoldNoJoinedCell(t *testing.T) {
	output, _ := commandOverEvents(t, t.TempDir(), []string{"--package", "./...", "--full"}, threeLineFailureEvents, 1)
	if strings.Contains(output, `\n`) {
		t.Fatalf("a --full failures cell holds an escaped newline:\n%s", output)
	}
}

// TestFullPackageFailurePrintsOneRowPerLogLine grades that a package failure splits its
// package log the same way, with an empty test cell.
func TestFullPackageFailurePrintsOneRowPerLogLine(t *testing.T) {
	output, code := commandOverEvents(t, t.TempDir(), []string{"--package", "./...", "--full"}, []string{
		`{"Action":"build-output","ImportPath":"canned","Output":"bad one\nbad two\n"}`,
		`{"Action":"build-fail","ImportPath":"canned"}`,
	}, 1)
	want := "failures[2]{package,test,line,lines}:\n  canned,\"\",bad one,2\n  canned,\"\",bad two,2\n"
	if !strings.Contains(output, want) || code != 1 {
		t.Fatalf("Command = (%d, %q), want (1, containing %q)", code, output, want)
	}
}

// TestFullFailedTestsCountsTests grades that the typed count is distinct failed tests,
// not printed rows.
func TestFullFailedTestsCountsTests(t *testing.T) {
	installCannedGo(t, cannedSet{name: t.Name(), events: threeLineFailureEvents, exit: 1})
	installCannedSelection(t)
	outcome := executeSelection(t, t.TempDir(), []string{"--package", "./...", "--full"})
	if want := (Outcome{Kind: OutcomeFailed, FailedTests: 1, Ran: 1}); outcome != want {
		t.Fatalf("outcome = %+v, want %+v", outcome, want)
	}
}
