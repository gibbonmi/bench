package probe

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/probe/probetest"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// bothPackages is the two expressions of twoPackages, joined in argv order.
const bothPackages = "./clamp ./other"

// twoPackages is the swap of clamp.go under two package expressions in argv order.
var twoPackages = []string{"clamp.go", "--swap", "n < 0", "--with", "n > 0", "--package", "./clamp", "--package", "./other", "--run", "^TestClampNegative$"}

// A repeated --package reaches each Go start in argv order, and the target cell names
// every expression in that order.
func TestProbeRepeatsThePackage(t *testing.T) {
	f := newFixture(t)
	argv := filepath.Join(scratchDir(t), "argv")
	record := "printf '%s\\n' \"$*\" >> " + sanitize.ShellQuote(argv) + "\n"
	installStubStarts(t, f, probetest.Start{Hook: record, Events: cannedPass}, probetest.Start{Hook: record, Events: cannedFailure("caught"), Exit: 1})
	out, code := runProbe(t, twoPackages...)
	want := probeRow(t, "bit", "clamp.go", "swap", "failed", 1, "TestClampNegative", "yes") +
		selectionRow(t, "package", bothPackages, "^TestClampNegative$", "passed", 1)
	requireRow(t, out, code, 0, want)
	body, err := os.ReadFile(argv)
	if err != nil {
		t.Fatal(err)
	}
	starts := strings.Split(strings.TrimSuffix(string(body), "\n"), "\n")
	if len(starts) != 2 {
		t.Fatalf("Go test starts = %q, want the baseline and the mutated run", starts)
	}
	for _, start := range starts {
		if !strings.HasSuffix(start, " "+bothPackages) {
			t.Fatalf("go argv = %q, want ./clamp and then ./other at its end", start)
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
