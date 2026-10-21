package testreport

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
)

// halfTimeoutMS is the warning threshold, built here from the timeout by its own
// expression so that a changed production divisor turns these tests red.
func halfTimeoutMS() int64 { return bounds.GoTestTimeout.Milliseconds() / 2 }

// warnedPackageEvents is a stream in which pkg runs one test that ends with action, and
// the package then ends with action after elapsedMS milliseconds.
func warnedPackageEvents(pkg, action string, elapsedMS int64) []string {
	return []string{
		fmt.Sprintf(`{"Action":"run","Package":%q,"Test":"TestCanned"}`, pkg),
		fmt.Sprintf(`{"Action":"output","Package":%q,"Test":"TestCanned","Output":"    canned_test.go:9: boom\n"}`, pkg),
		testTerminal(action, pkg, "TestCanned", 0.01),
		testTerminal(action, pkg, "", float64(elapsedMS)/1000),
	}
}

// warningBlock is the expected `timeout_warnings` block with one threshold row for each
// package.
func warningBlock(packages ...string) string {
	block := fmt.Sprintf("timeout_warnings[%d]{package,elapsed_ms,timeout_ms}:\n", len(packages))
	for _, pkg := range packages {
		block += fmt.Sprintf("  %s,%d,%d\n", pkg, halfTimeoutMS(), bounds.GoTestTimeout.Milliseconds())
	}
	return block
}

// TestTimeoutWarningAtThreshold grades that a package at half of the timeout gets a row.
func TestTimeoutWarningAtThreshold(t *testing.T) {
	output, code := commandOverEvents(t, t.TempDir(), []string{"--package", "./..."}, warnedPackageEvents("canned", "pass", halfTimeoutMS()), 0)
	if want := warningBlock("canned"); !strings.Contains(output, want) || code != 0 {
		t.Fatalf("Command = (%d, %q), want (0, containing %q)", code, output, want)
	}
}

// TestTimeoutWarningBelowThreshold grades that a package one millisecond below the
// threshold prints no warning block.
func TestTimeoutWarningBelowThreshold(t *testing.T) {
	output, code := commandOverEvents(t, t.TempDir(), []string{"--package", "./..."}, warnedPackageEvents("canned", "pass", halfTimeoutMS()-1), 0)
	if strings.Contains(output, "timeout_warnings") || code != 0 {
		t.Fatalf("Command = (%d, %q), want (0, no timeout_warnings block)", code, output)
	}
}

// TestTimeoutWarningsOrderByPackage grades that warned packages, whose events arrive in
// reverse name order, print in package name order.
func TestTimeoutWarningsOrderByPackage(t *testing.T) {
	events := append(append(warnedPackageEvents("zeta", "pass", halfTimeoutMS()), warnedPackageEvents("mid", "pass", halfTimeoutMS())...), warnedPackageEvents("alpha", "pass", halfTimeoutMS())...)
	output, code := commandOverEvents(t, t.TempDir(), []string{"--package", "./..."}, events, 0)
	if want := warningBlock("alpha", "mid", "zeta"); !strings.Contains(output, want) || code != 0 {
		t.Fatalf("Command = (%d, %q), want (0, containing %q)", code, output, want)
	}
}

// TestTimeoutWarningsRenderFirst grades that the warning block leads a plain result, so a
// bounded response always shows it.
func TestTimeoutWarningsRenderFirst(t *testing.T) {
	output, code := commandOverEvents(t, t.TempDir(), []string{"--package", "./..."}, warnedPackageEvents("canned", "pass", halfTimeoutMS()), 0)
	if want := warningBlock("canned") + "packages[1]"; !strings.HasPrefix(output, want) || code != 0 {
		t.Fatalf("Command = (%d, %q), want (0, prefix %q)", code, output, want)
	}
}

// TestTimeoutWarningsFollowCheckRow grades that a named check prints the check row, then
// the warning block, then the packages table.
func TestTimeoutWarningsFollowCheckRow(t *testing.T) {
	output, code := commandOverEvents(t, t.TempDir(), []string{"--check", "line-routing"}, warnedPackageEvents("canned", "pass", halfTimeoutMS()), 0)
	want := checkHeader + "  line-routing,conformance,1,0\n" + warningBlock("canned") + "packages[1]"
	if !strings.HasPrefix(output, want) || code != 0 {
		t.Fatalf("Command = (%d, %q), want (0, prefix %q)", code, output, want)
	}
}

// TestTimeoutWarningKeepsPassingExit grades that a warning does not fail a passing run.
func TestTimeoutWarningKeepsPassingExit(t *testing.T) {
	output, code := commandOverEvents(t, t.TempDir(), []string{"--package", "./..."}, warnedPackageEvents("canned", "pass", halfTimeoutMS()), 0)
	if !strings.HasPrefix(output, warningBlock("canned")) || code != 0 {
		t.Fatalf("Command = (%d, %q), want (0, the warning block first)", code, output)
	}
}

// TestTimeoutWarningKeepsFailingExit grades that a warning does not hide a failure: the
// run exits 1 and still prints its failure row.
func TestTimeoutWarningKeepsFailingExit(t *testing.T) {
	output, code := commandOverEvents(t, t.TempDir(), []string{"--package", "./..."}, warnedPackageEvents("canned", "fail", halfTimeoutMS()), 1)
	failure := "failures[1]{package,test,line,lines}:\n  canned,TestCanned,\"canned_test.go:9: boom\",1\n"
	if !strings.HasPrefix(output, warningBlock("canned")) || !strings.Contains(output, failure) || code != 1 {
		t.Fatalf("Command = (%d, %q), want (1, the warning block first and %q)", code, output, failure)
	}
}
