package testreport

import (
	"strings"
	"testing"
)

// commandOverEvents runs `Command` at root with args over one canned stream.
func commandOverEvents(t *testing.T, root string, args, events []string, exit int) (string, int) {
	t.Helper()
	installCannedGo(t, cannedSet{name: t.Name(), events: events, exit: exit})
	installCannedSelection(t)
	return Command(root, args)
}

// twoTestsRunEvents is a passing stream of package `canned` in which two distinct tests
// emit a run event and one of them emits it twice.
var twoTestsRunEvents = []string{
	`{"Action":"run","Package":"canned","Test":"TestCanned"}`,
	`{"Action":"run","Package":"canned","Test":"TestCanned/sub"}`,
	`{"Action":"run","Package":"canned","Test":"TestCanned"}`,
	`{"Action":"pass","Package":"canned","Test":"TestCanned/sub","Elapsed":0.01}`,
	`{"Action":"pass","Package":"canned","Test":"TestCanned","Elapsed":0.01}`,
	`{"Action":"pass","Package":"canned","Elapsed":0.25}`,
}

// TestPackagesRowCountsRunEvents grades the `tests_run` cell: the count of distinct tests
// and subtests of the package that emitted a run event. A repeated run event counts once.
func TestPackagesRowCountsRunEvents(t *testing.T) {
	output, code := commandOverEvents(t, t.TempDir(), []string{"--package", "./..."}, twoTestsRunEvents, 0)
	want := "packages[1]{package,status,elapsed_ms,tests_run}:\n  canned,pass,250,2\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
	if output != want || code != 0 {
		t.Fatalf("Command = (%d, %q), want (0, %q)", code, output, want)
	}
}

// TestPackagesRowNoTestsCountsZero grades that each row counts the run events of its own
// package: a `no-tests` package prints 0 beside a package that ran one test.
func TestPackagesRowNoTestsCountsZero(t *testing.T) {
	output, code := commandOverEvents(t, t.TempDir(), []string{"--package", "./..."}, []string{
		`{"Action":"output","Package":"canned","Output":"?   \tcanned\t[no test files]\n"}`,
		`{"Action":"skip","Package":"canned","Elapsed":0}`,
		`{"Action":"run","Package":"other","Test":"TestOther"}`,
		`{"Action":"pass","Package":"other","Test":"TestOther","Elapsed":0.01}`,
		`{"Action":"pass","Package":"other","Elapsed":0.1}`,
	}, 0)
	want := "packages[2]{package,status,elapsed_ms,tests_run}:\n  canned,no-tests,0,0\n  other,pass,100,1\n"
	if !strings.HasPrefix(output, want) || code != 0 {
		t.Fatalf("Command = (%d, %q), want (0, prefix %q)", code, output, want)
	}
}

// TestPackagesRowIgnoresUnnamedRunEvent grades that a run event with no test name counts
// no test that ran: the package that emits only that event prints 0.
func TestPackagesRowIgnoresUnnamedRunEvent(t *testing.T) {
	output, code := commandOverEvents(t, t.TempDir(), []string{"--package", "./..."}, []string{
		`{"Action":"run","Package":"canned"}`,
		`{"Action":"pass","Package":"canned","Elapsed":0.25}`,
	}, 0)
	want := "packages[1]{package,status,elapsed_ms,tests_run}:\n  canned,pass,250,0\n"
	if !strings.HasPrefix(output, want) {
		t.Fatalf("Command = (%d, %q), want prefix %q", code, output, want)
	}
}

// TestPackageFormHasNoSelectedBy grades that a `--package` result prints the packages
// header with no cause cell.
func TestPackageFormHasNoSelectedBy(t *testing.T) {
	output, code := commandOverEvents(t, t.TempDir(), []string{"--package", "./..."}, cannedSetNamed(t, "passing").events, 0)
	header := "packages[1]{package,status,elapsed_ms,tests_run}:\n"
	if !strings.HasPrefix(output, header) || strings.Contains(output, "selected_by") || code != 0 {
		t.Fatalf("Command = (%d, %q), want (0, header %q and no selected_by)", code, output, header)
	}
}
