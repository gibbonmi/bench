package testreport

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
)

// systemRunArgs selects the system suite filtered to one test.
var systemRunArgs = []string{"--check", gate.SystemPhaseName, "--run", "^TestX$"}

// TestSystemRunPatternReachesGoArgv grades that the Go child of a filtered system check
// receives the pattern after the system suite operands.
func TestSystemRunPatternReachesGoArgv(t *testing.T) {
	root := canonicalTestDir(t)
	marker := filepath.Join(t.TempDir(), "environment")
	goDir := t.TempDir()
	writeCheckGo(t, filepath.Join(goDir, "go"), marker)
	t.Setenv("PATH", goDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	installCannedSelection(t)
	t.Setenv("BENCH_KIT", root)

	output, code := Command(root, systemRunArgs)
	if code != 0 {
		t.Fatalf("filtered system check = %d, want 0\n%s", code, output)
	}
	operands, _ := gate.SystemSuite(root)
	want := " " + strings.Join(operands, " ") + " -run ^TestX$\n"
	if environment := readTestReportFile(t, marker); !strings.Contains(environment, want) {
		t.Fatalf("system-check Go argv missing %q:\n%s", strings.TrimSpace(want), environment)
	}
}

// TestSystemRequestRunFact grades that a filtered system request answers its own pattern.
func TestSystemRequestRunFact(t *testing.T) {
	request, line, code := Prepare(t.TempDir(), systemRunArgs)
	if line != "" {
		t.Fatalf("Prepare refused %v: (%d, %q)", systemRunArgs, code, line)
	}
	if got := request.Run(); got != "^TestX$" {
		t.Fatalf("Run() = %q, want %q", got, "^TestX$")
	}
}

// TestSystemRunPatternNoMatchRefusalWins grades that a pattern with no match keeps the
// run-pattern refusal, and the zero-rule title does not print.
func TestSystemRunPatternNoMatchRefusalWins(t *testing.T) {
	root := canonicalTestDir(t)
	t.Setenv("BENCH_KIT", root)
	output, code := commandOverEvents(t, root, []string{"--check", gate.SystemPhaseName, "--run", "^TestNone$"}, []string{
		`{"Action":"pass","Package":"canned","Elapsed":0.01}`,
	}, 0)
	if code != 1 || !strings.Contains(output, "go test reported no test runs") || strings.Contains(output, "named check ran nothing") {
		t.Fatalf("Command = (%d, %q), want exit 1 with the run-pattern refusal and no zero-rule title", code, output)
	}
}

// TestProseRefusesRunPattern grades that the prose check refuses a run pattern as usage.
// The prose grade returns before the Go path, so only the grammar can refuse it.
func TestProseRefusesRunPattern(t *testing.T) {
	output, code := Command(t.TempDir(), []string{"--check", proseCheckName, "--run", "^TestX$"})
	if code != 2 || !strings.HasPrefix(output, "usage: bench test") {
		t.Fatalf("Command = (%d, %q), want usage exit 2", code, output)
	}
}
