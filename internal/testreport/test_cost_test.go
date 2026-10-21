package testreport

import (
	"strconv"
	"strings"
	"testing"
)

// testTerminal is the terminal event of test in package pkg with the elapsed seconds.
func testTerminal(action, pkg, test string, seconds float64) string {
	return `{"Action":"` + action + `","Package":"` + pkg + `","Test":"` + test + `","Elapsed":` + strconv.FormatFloat(seconds, 'f', -1, 64) + `}`
}

// renderEvents decodes the go test -json events and renders the report with --full.
func renderEvents(t *testing.T, events ...string) string {
	t.Helper()
	report, err := decode(strings.NewReader(strings.Join(events, "\n")))
	if err != nil {
		t.Fatal(err)
	}
	output, err := report.render(true)
	if err != nil {
		t.Fatal(err)
	}
	return output
}

// wantSlowest grades the rendered tail from the slowest header on. Each row is the
// comma-joined package, test, and elapsed_ms cells.
func wantSlowest(t *testing.T, events []string, rows ...string) {
	t.Helper()
	want := "slowest[" + strconv.Itoa(len(rows)) + "]{package,test,elapsed_ms}:\n"
	for _, row := range rows {
		want += "  " + row + "\n"
	}
	output := renderEvents(t, events...)
	at := strings.Index(output, "slowest[")
	if at < 0 || output[at:] != want {
		t.Fatalf("report = %q, want the tail %q", output, want)
	}
}

func TestSlowestCarriesTestElapsed(t *testing.T) {
	wantSlowest(t, []string{
		`{"Action":"run","Package":"canned","Test":"TestA"}`,
		testTerminal("pass", "canned", "TestA", 1.5),
		`{"Action":"pass","Package":"canned","Elapsed":2}`,
	}, "canned,TestA,1500")
}

func TestSlowestRoundsToNearestMillisecond(t *testing.T) {
	// A fraction above one half rounds up and one below rounds down, so truncation and
	// ceiling both go red.
	wantSlowest(t, []string{
		testTerminal("pass", "canned", "TestA", 0.0026),
		testTerminal("pass", "canned", "TestB", 0.0024),
	}, "canned,TestA,3", "canned,TestB,2")
}

// TestSlowestLastTerminalEventWins sends the smaller elapsed value last and then the
// larger value last, so a first-wins, a max-wins, or a min-wins reading goes red.
func TestSlowestLastTerminalEventWins(t *testing.T) {
	wantSlowest(t, []string{
		testTerminal("pass", "canned", "TestA", 0.5),
		testTerminal("fail", "canned", "TestA", 0.3),
	}, "canned,TestA,300")
	wantSlowest(t, []string{
		testTerminal("fail", "canned", "TestA", 0.2),
		testTerminal("pass", "canned", "TestA", 0.6),
	}, "canned,TestA,600")
}

func TestSlowestKeepsThreePerPackage(t *testing.T) {
	wantSlowest(t, []string{
		testTerminal("pass", "canned", "TestD", 0.1),
		testTerminal("pass", "canned", "TestA", 0.4),
		testTerminal("pass", "canned", "TestB", 0.3),
		testTerminal("pass", "canned", "TestC", 0.2),
	}, "canned,TestA,400", "canned,TestB,300", "canned,TestC,200")
}

func TestSlowestOrdersByElapsed(t *testing.T) {
	wantSlowest(t, []string{
		testTerminal("pass", "canned", "TestA", 0.1),
		testTerminal("pass", "canned", "TestB", 0.2),
		testTerminal("pass", "canned", "TestC", 0.3),
	}, "canned,TestC,300", "canned,TestB,200", "canned,TestA,100")
}

// TestSlowestBreaksTiesByName renders the stream many times, because the tests map
// iterates in a different order each time and a missing tie-break holds by chance.
func TestSlowestBreaksTiesByName(t *testing.T) {
	for range 32 {
		wantSlowest(t, []string{
			testTerminal("pass", "canned", "TestB", 0.2),
			testTerminal("pass", "canned", "TestA", 0.2),
		}, "canned,TestA,200", "canned,TestB,200")
	}
}

func TestSlowestOrdersByPackage(t *testing.T) {
	events := make([]string, 0, 8)
	for i, name := range []string{"TestA", "TestB", "TestC", "TestD"} {
		events = append(events,
			testTerminal("pass", "zeta", name, 0.5+float64(i)/10),
			testTerminal("pass", "alpha", name, 0.1+float64(i)/10))
	}
	wantSlowest(t, events,
		"alpha,TestD,400", "alpha,TestC,300", "alpha,TestB,200",
		"zeta,TestD,800", "zeta,TestC,700", "zeta,TestB,600")
}

func TestSlowestExcludesSubtests(t *testing.T) {
	wantSlowest(t, []string{
		testTerminal("pass", "canned", "TestA/sub", 0.9),
		testTerminal("pass", "canned", "TestA", 0.2),
	}, "canned,TestA,200")
}

func TestSlowestExcludesSkippedTests(t *testing.T) {
	wantSlowest(t, []string{
		testTerminal("skip", "canned", "TestA", 0.3),
		testTerminal("pass", "canned", "TestB", 0.1),
	}, "canned,TestB,100")
}

func TestSlowestIncludesFailedTests(t *testing.T) {
	wantSlowest(t, []string{testTerminal("fail", "canned", "TestA", 0.3)}, "canned,TestA,300")
}

func TestSlowestExcludesUnfinishedTests(t *testing.T) {
	wantSlowest(t, []string{
		`{"Action":"run","Package":"canned","Test":"TestA"}`,
		`{"Action":"output","Package":"canned","Test":"TestA","Output":"    a_test.go:9: started\n"}`,
		testTerminal("pass", "canned", "TestB", 0.1),
	}, "canned,TestB,100")
}

func TestSlowestMissingElapsedIsZero(t *testing.T) {
	wantSlowest(t, []string{`{"Action":"pass","Package":"canned","Test":"TestA"}`}, "canned,TestA,0")
}

func TestSlowestExcludesPackageEvent(t *testing.T) {
	wantSlowest(t, []string{`{"Action":"pass","Package":"canned","Elapsed":0.5}`})
}

// fullRun is the argument list of a `--full` run over every package.
var fullRun = []string{"--package", "./...", "--full"}

// TestSlowestEmptyOnFullRun grades that a `--full` run with no test events prints the
// empty slowest table.
func TestSlowestEmptyOnFullRun(t *testing.T) {
	output, code := commandOverEvents(t, t.TempDir(), fullRun, cannedSetNamed(t, "no-run").events, 0)
	want := "packages[1]{package,status,elapsed_ms,tests_run}:\n  canned,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nslowest[0]{package,test,elapsed_ms}:\n"
	if output != want || code != 0 {
		t.Fatalf("Command = (%d, %q), want (0, %q)", code, output, want)
	}
}

// TestSlowestFollowsSkipsOnFullRun grades that on a `--full` run the slowest block comes
// right after a skips block that holds a row.
func TestSlowestFollowsSkipsOnFullRun(t *testing.T) {
	output, code := commandOverEvents(t, t.TempDir(), fullRun, []string{
		`{"Action":"run","Package":"canned","Test":"TestA"}`,
		`{"Action":"output","Package":"canned","Test":"TestA","Output":"    not here\n"}`,
		testTerminal("skip", "canned", "TestA", 0),
		`{"Action":"run","Package":"canned","Test":"TestB"}`,
		testTerminal("pass", "canned", "TestB", 0.02),
		`{"Action":"pass","Package":"canned","Elapsed":0.1}`,
	}, 0)
	want := "packages[1]{package,status,elapsed_ms,tests_run}:\n  canned,pass,100,2\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  canned,TestA,not here\nslowest[1]{package,test,elapsed_ms}:\n  canned,TestB,20\n"
	if output != want || code != 0 {
		t.Fatalf("Command = (%d, %q), want (0, %q)", code, output, want)
	}
}

// TestSlowestPrecedesRootConformance grades that the conditional root_conformance block
// stays last, right after the slowest block.
func TestSlowestPrecedesRootConformance(t *testing.T) {
	const pkg = "github.com/gibbonmi/bench/internal/conformance"
	output := renderEvents(t,
		testTerminal("skip", pkg, "TestRootConformance", 0),
		`{"Action":"pass","Package":"`+pkg+`","Elapsed":0.1}`)
	if want := "slowest[0]{package,test,elapsed_ms}:\nroot_conformance[1]{package,status,route}:\n"; !strings.Contains(output, want) {
		t.Fatalf("report = %q, want %q", output, want)
	}
}

// TestPackageElapsedIgnoresTestElapsed grades that each cell reads its own terminal
// event: the package cell reads the package event and the slowest cell the test event.
func TestPackageElapsedIgnoresTestElapsed(t *testing.T) {
	output := renderEvents(t,
		`{"Action":"run","Package":"canned","Test":"TestA"}`,
		testTerminal("pass", "canned", "TestA", 2.0),
		`{"Action":"pass","Package":"canned","Elapsed":0.5}`)
	want := "packages[1]{package,status,elapsed_ms,tests_run}:\n  canned,pass,500,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nslowest[1]{package,test,elapsed_ms}:\n  canned,TestA,2000\n"
	if output != want {
		t.Fatalf("report = %q, want %q", output, want)
	}
}
