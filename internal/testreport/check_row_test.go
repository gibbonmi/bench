package testreport

import (
	"strings"
	"testing"
)

const checkHeader = "check[1]{name,kind,tests_run,subjects}:\n"

// TestNamedCheckPrintsCheckRowFirst grades that the check row is the first block of the
// result and that the packages table follows it.
func TestNamedCheckPrintsCheckRowFirst(t *testing.T) {
	output, code := commandOverEvents(t, t.TempDir(), []string{"--check", "line-routing"}, cannedSetNamed(t, "passing").events, 0)
	want := checkHeader + "  line-routing,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n"
	if !strings.HasPrefix(output, want) || code != 0 {
		t.Fatalf("Command = (%d, %q), want (0, prefix %q)", code, output, want)
	}
}

// TestNamedCheckFailurePrintsCheckRowFirst grades that a failed check reached a verdict,
// so the check row leads its result too.
func TestNamedCheckFailurePrintsCheckRowFirst(t *testing.T) {
	set := cannedSetNamed(t, "failing")
	output, code := commandOverEvents(t, t.TempDir(), []string{"--check", "line-routing"}, set.events, set.exit)
	want := checkHeader + "  line-routing,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n"
	if !strings.HasPrefix(output, want) || code != 1 {
		t.Fatalf("Command = (%d, %q), want (1, prefix %q)", code, output, want)
	}
}

// TestNamedCheckRowCountsEachRunTest grades that the check row's tests_run cell is the
// count of distinct tests that emitted a run event, not a ran-or-not flag.
func TestNamedCheckRowCountsEachRunTest(t *testing.T) {
	output, code := commandOverEvents(t, t.TempDir(), []string{"--check", "line-routing"}, twoTestsRunEvents, 0)
	want := checkHeader + "  line-routing,conformance,2,0\n"
	if !strings.HasPrefix(output, want) || code != 0 {
		t.Fatalf("Command = (%d, %q), want (0, prefix %q)", code, output, want)
	}
}

// TestNamedCheckInterruptKeepsItsBytes grades that an interrupted check reached no
// verdict, so it prints the interrupt error and no check row.
func TestNamedCheckInterruptKeepsItsBytes(t *testing.T) {
	installSignallingGo(t)
	installCannedSelection(t)
	output, code := Command(t.TempDir(), []string{"--check", "line-routing"})
	want := "error: go test interrupted — child process group cancelled\n"
	if output != want || code != 1 {
		t.Fatalf("Command = (%d, %q), want (1, %q)", code, output, want)
	}
}

// TestNamedCheckFailureWithNoRunEventIsNotRanNothing grades that a test failure is a
// verdict even when no run event came first: the failures table prints, and the
// zero-rule title does not.
func TestNamedCheckFailureWithNoRunEventIsNotRanNothing(t *testing.T) {
	output, code := commandOverEvents(t, t.TempDir(), []string{"--check", "line-routing"}, []string{
		`{"Action":"output","Package":"canned","Test":"TestCanned","Output":"    canned_test.go:9: boom\n"}`,
		`{"Action":"fail","Package":"canned","Test":"TestCanned","Elapsed":0.01}`,
		`{"Action":"fail","Package":"canned","Elapsed":0.5}`,
	}, 1)
	failures := "failures[1]{package,test,line,lines}:\n  canned,TestCanned,\"canned_test.go:9: boom\",1\n"
	if !strings.Contains(output, failures) || strings.Contains(output, "named check ran nothing") || code != 1 {
		t.Fatalf("Command = (%d, %q), want (1, %q and no zero-rule title)", code, output, failures)
	}
}

// TestSystemCheckRowKind grades the kind cell of the system check.
func TestSystemCheckRowKind(t *testing.T) {
	root := canonicalTestDir(t)
	t.Setenv("BENCH_KIT", root)
	output, code := commandOverEvents(t, root, []string{"--check", "system"}, cannedSetNamed(t, "passing").events, 0)
	want := checkHeader + "  system,system,1,0\n"
	if !strings.HasPrefix(output, want) || code != 0 {
		t.Fatalf("Command = (%d, %q), want (0, prefix %q)", code, output, want)
	}
}

// TestNamedCheckRanNothingExitsOne grades the zero rule: a package pass with no run
// event is not evidence that the check ran.
func TestNamedCheckRanNothingExitsOne(t *testing.T) {
	output, code := commandOverEvents(t, t.TempDir(), []string{"--check", "line-routing"}, []string{`{"Action":"pass","Package":"canned","Elapsed":0.01}`}, 0)
	want := checkHeader + "  line-routing,conformance,0,0\nerror: named check ran nothing"
	if !strings.HasPrefix(output, want) || code != 1 {
		t.Fatalf("Command = (%d, %q), want (1, prefix %q)", code, output, want)
	}
}

// TestNamedCheckBuildFailureWinsOverZeroRule grades that a compile failure with no run
// event prints its diagnostic and not the zero-rule title.
func TestNamedCheckBuildFailureWinsOverZeroRule(t *testing.T) {
	set := cannedSetNamed(t, "build-fail")
	output, code := commandOverEvents(t, t.TempDir(), []string{"--check", "line-routing"}, set.events, set.exit)
	want := checkHeader + "  line-routing,conformance,0,0\n"
	if !strings.HasPrefix(output, want) || code != 1 {
		t.Fatalf("Command = (%d, %q), want (1, prefix %q)", code, output, want)
	}
	if !strings.Contains(output, "syntax error: unexpected }") || strings.Contains(output, "named check ran nothing") {
		t.Fatalf("Command output = %q, want the compile diagnostic and no zero-rule title", output)
	}
}

// TestPackageRunWithNoTestKeepsExitZero grades that the zero rule is a named-check rule
// only.
func TestPackageRunWithNoTestKeepsExitZero(t *testing.T) {
	if output, code := commandOverEvents(t, focusedTestModule(t), []string{"--package", "chosen"}, []string{`{"Action":"pass","Package":"focusedfixture/chosen","Elapsed":0.01}`}, 0); code != 0 {
		t.Fatalf("Command = (%d, %q), want exit 0", code, output)
	}
}

// TestChangedRunWithNoTestKeepsExitZero grades that a `--changed` run with no test keeps
// exit 0, because the merge reads that exit as its retry signal.
func TestChangedRunWithNoTestKeepsExitZero(t *testing.T) {
	if output, code := changedCommandOverOnePackage(t); code != 0 {
		t.Fatalf("Command = (%d, %q), want exit 0", code, output)
	}
}
