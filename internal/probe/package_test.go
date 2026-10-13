package probe

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/probe/probetest"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// bothPackages is the two expressions of twoPackages, joined in argv order.
const bothPackages = "./clamp ./other"

// twoPackages is the swap of clamp.go under two package expressions in argv order.
var twoPackages = []string{"clamp.go", "--swap", "n < 0", "--with", "n > 0", "--package", "./clamp", "--package", "./other", "--run", "^TestClampNegative$"}

// A repeated --package reaches each Go start as separate arguments in argv order, and the
// target cell names every expression in that order. The stub writes one argument per line
// to a file for each start, so a joined value reads as one line.
func TestProbeRepeatsThePackage(t *testing.T) {
	f := newFixture(t)
	argvs := scratchDir(t)
	record := "printf '%s\\n' \"$@\" > " + sanitize.ShellQuote(argvs) + "/$n\n"
	installStubStarts(t, f, probetest.Start{Hook: record, Events: cannedPass}, probetest.Start{Hook: record, Events: cannedFailure("caught"), Exit: 1})
	out, code := runProbe(t, twoPackages...)
	want := probeRow(t, "bit", "clamp.go", "swap", "failed", 1, "TestClampNegative", "yes") +
		selectionRow(t, "package", bothPackages, "^TestClampNegative$", "passed", 1)
	requireRow(t, out, code, 0, want)
	starts, err := os.ReadDir(argvs)
	if err != nil {
		t.Fatal(err)
	}
	if len(starts) != 2 {
		t.Fatalf("Go test starts = %d, want the baseline and the mutated run", len(starts))
	}
	packages := strings.Fields(bothPackages)
	for _, start := range starts {
		body, err := os.ReadFile(filepath.Join(argvs, start.Name()))
		if err != nil {
			t.Fatal(err)
		}
		args := strings.Split(strings.TrimSuffix(string(body), "\n"), "\n")
		if len(args) < len(packages) || !slices.Equal(args[len(args)-len(packages):], packages) {
			t.Fatalf("go argv = %q, want %q as its last arguments", args, packages)
		}
	}
	requireSubjectBytes(t, f, clampSource)
}

// Every package expression shares one Go start for each run, and each run owns one build.
func TestProbePackagesShareOneBuildPerRun(t *testing.T) {
	f := newFixture(t)
	records := buildEachRun(t, f, "clamp.go")
	installStubGo(t, f, cannedFailure("caught"), 1, "")
	out, code := runProbe(t, twoPackages...)
	requireRow(t, out, code, 0, probeRow(t, "bit", "clamp.go", "swap", "failed", 1, "TestClampNegative", "yes"))
	if got := startCount(t, f); got != 2 {
		t.Fatalf("Go test starts = %d, want 2", got)
	}
	if got := len(buildRecords(t, records)); got != 2 {
		t.Fatalf("builds = %d, want 2", got)
	}
	requireSubjectBytes(t, f, clampSource)
}
