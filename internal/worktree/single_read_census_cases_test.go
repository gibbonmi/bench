package worktree

// This file holds the synthetic cases of the single-read census. Each case plants an
// effects file, a gate file, and one reader file, and compares the reports with the
// census formatter's output for literal inputs.

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// syntheticEffectsFile mirrors effects.go: currentTime reads the clock, newAmbient the
// kit and the clock, and Home the Bench home. Its os import serves an added read.
const syntheticEffectsFile = `package worktree

import (
	"os"
	"time"
	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/gate"
)

func currentTime() time.Time { return time.Now() }
func newAmbient(home string) ambient { return ambient{kit: gate.KitValue(), now: currentTime()} }
func Home() string { return benchhome.Dir() }
`

// syntheticGateFile is the gate directory of a synthetic package. KitDir calls
// KitValue, KitRoot reaches it through kitRoot, and LaneForCommitAtKit reaches no read.
const syntheticGateFile = `package gate

import "os"

func KitValue() string { return os.Getenv("BENCH_KIT") }
func KitDir() string { return KitValue() }
func KitRoot() string { return kitRoot() }
func kitRoot() string { return KitValue() }
func LaneForCommitAtKit(root, kit string) string { return kit }
`

// wantReadCensus runs the census over effects, syntheticGateFile, and a reader.go whose
// source begins on line 3, after the package clause. The reports must equal want.
func wantReadCensus(t *testing.T, effects, source string, want ...string) {
	t.Helper()
	dir := plantTestFiles(t, map[string]string{effectsFile: effects, "reader.go": "package worktree\n\n" + source})
	reports, err := singleReadCensus(dir, plantTestFiles(t, map[string]string{"gate.go": syntheticGateFile}))
	if err != nil {
		t.Fatalf("single-read census: %v", err)
	}
	if !slices.Equal(reports, want) {
		t.Fatalf("single-read census = %q, want exactly %q", reports, want)
	}
}

func TestSingleReadCensusAcceptsOneReadInAnEntry(t *testing.T) {
	t.Parallel()
	wantReadCensus(t, syntheticEffectsFile, "func Run(paths []string) {\n\tnow := currentTime()\n\tfor range paths {\n\t\t_ = now\n\t}\n\t_ = func() { _ = now }\n}\n")
}

func TestSingleReadCensusRefusesAReadInAnUnexportedFunction(t *testing.T) {
	t.Parallel()
	wantReadCensus(t, syntheticEffectsFile, "func stamp() { _ = currentTime() }\n",
		singleReadReport("reader.go", 3, "stamp", readBelowEntry, "currentTime", ""))
}

// TestSingleReadCensusRefusesAReadInAFunctionLiteral proves a literal is no shelter.
func TestSingleReadCensusRefusesAReadInAFunctionLiteral(t *testing.T) {
	t.Parallel()
	wantReadCensus(t, syntheticEffectsFile, "func Run() {\n\tstamp := func() { _ = currentTime() }\n\tstamp()\n}\n",
		singleReadReport("reader.go", 4, "Run", readBelowEntry, "currentTime", ""))
}

// TestSingleReadCensusRefusesAReadInALoopBody proves a range body is no shelter, nor
// the condition, the post statement, or the body of a three-clause loop. The init
// statement runs once, so a read there is the entry's own.
func TestSingleReadCensusRefusesAReadInALoopBody(t *testing.T) {
	t.Parallel()
	wantReadCensus(t, syntheticEffectsFile, "func Run(paths []string) {\n\tfor range paths {\n\t\t_ = Home()\n\t}\n}\n",
		singleReadReport("reader.go", 5, "Run", readBelowEntry, "Home", ""))
	for _, loop := range []struct {
		place, clause string
		want          []string
	}{
		{"init", "for i := len(Home()); i < 1; i++ {\n\t}", nil},
		{"condition", "for i := 0; i < len(Home()); i++ {\n\t}", []string{singleReadReport("reader.go", 4, "Run", readBelowEntry, "Home", "")}},
		{"post", "for i := 0; i < 1; i += len(Home()) {\n\t}", []string{singleReadReport("reader.go", 4, "Run", readBelowEntry, "Home", "")}},
		{"body", "for i := 0; i < 1; i++ {\n\t\t_ = Home()\n\t}", []string{singleReadReport("reader.go", 5, "Run", readBelowEntry, "Home", "")}},
	} {
		t.Run(loop.place, func(t *testing.T) {
			t.Parallel()
			wantReadCensus(t, syntheticEffectsFile, "func Run() {\n\t"+loop.clause+"\n}\n", loop.want...)
		})
	}
}

func TestSingleReadCensusRefusesASecondRead(t *testing.T) {
	t.Parallel()
	wantReadCensus(t, syntheticEffectsFile, "func Run() {\n\t_ = currentTime()\n\t_ = currentTime()\n}\n",
		singleReadReport("reader.go", 5, "Run", readSecondTime, "", "time.Now"))
}

// TestSingleReadCensusRefusesAHelperCallToAReadingEntry proves a helper call to a
// reading exported function or exported method is refused.
func TestSingleReadCensusRefusesAHelperCallToAReadingEntry(t *testing.T) {
	t.Parallel()
	wantReadCensus(t, syntheticEffectsFile, "func Run() { _ = currentTime() }\n\nfunc helper() {\n\tRun()\n}\n",
		singleReadReport("reader.go", 6, "helper", readThroughEntry, "Run", "time.Now"))
	wantReadCensus(t, syntheticEffectsFile, "type runner struct{}\n\nfunc (runner) Run() { _ = currentTime() }\n\nfunc helper(r runner) {\n\tr.Run()\n}\n",
		singleReadReport("reader.go", 8, "helper", readThroughEntry, "Run", "time.Now"))
}

// TestSingleReadCensusRefusesAFunctionValue proves a reference that is not a call reads.
func TestSingleReadCensusRefusesAFunctionValue(t *testing.T) {
	t.Parallel()
	wantReadCensus(t, syntheticEffectsFile, "func clocks() { _ = []any{currentTime} }\n",
		singleReadReport("reader.go", 3, "clocks", readBelowEntry, "currentTime", ""))
}

func TestSingleReadCensusRefusesAPackageLevelRead(t *testing.T) {
	t.Parallel()
	wantReadCensus(t, syntheticEffectsFile, "var defaultHome = Home()\n",
		singleReadReport("reader.go", 3, "defaultHome", readBelowEntry, "Home", ""))
}

func TestSingleReadCensusDerivesTheReadSetFromEffects(t *testing.T) {
	t.Parallel()
	effects := syntheticEffectsFile + "func shellPath() string { return os.Getenv(\"SHELL\") }\n"
	wantReadCensus(t, effects, "func launch() { _ = shellPath() }\n",
		singleReadReport("reader.go", 3, "launch", readBelowEntry, "shellPath", ""))
}

// TestSingleReadCensusRefusesAQualifiedRead proves a qualified read is refused under
// any import alias, in the reader or in the effects file.
func TestSingleReadCensusRefusesAQualifiedRead(t *testing.T) {
	t.Parallel()
	wantReadCensus(t, syntheticEffectsFile, "import \"github.com/gibbonmi/bench/internal/benchhome\"\n\nfunc poolHome() string { return benchhome.Dir() }\n",
		singleReadReport("reader.go", 5, "poolHome", readBelowEntry, "benchhome.Dir", ""))
	wantReadCensus(t, syntheticEffectsFile, "import bh \"github.com/gibbonmi/bench/internal/benchhome\"\n\nfunc poolHome() string { return bh.Dir() }\n",
		singleReadReport("reader.go", 5, "poolHome", readBelowEntry, "bh.Dir", ""))
	aliased := strings.NewReplacer("\t\"github.com/gibbonmi/bench/internal/gate\"", "\tkit \"github.com/gibbonmi/bench/internal/gate\"", "gate.KitValue()", "kit.KitValue()").Replace(syntheticEffectsFile)
	wantReadCensus(t, aliased, "import \"github.com/gibbonmi/bench/internal/gate\"\n\nfunc kitDir() string { return gate.KitDir() }\n",
		singleReadReport("reader.go", 5, "kitDir", readBelowEntry, "gate.KitDir", ""))
}

// TestSingleReadCensusRefusesAKitWrapperCall proves an exported gate function that
// calls the kit read is a read, and a gate form that takes the kit value is not.
func TestSingleReadCensusRefusesAKitWrapperCall(t *testing.T) {
	t.Parallel()
	wantReadCensus(t, syntheticEffectsFile, "import \"github.com/gibbonmi/bench/internal/gate\"\n\nfunc kitDir(root, kit string) string {\n\t_ = gate.LaneForCommitAtKit(root, kit)\n\treturn gate.KitDir()\n}\n",
		singleReadReport("reader.go", 7, "kitDir", readBelowEntry, "gate.KitDir", ""))
}

func TestSingleReadCensusRefusesAnIndirectKitRead(t *testing.T) {
	t.Parallel()
	wantReadCensus(t, syntheticEffectsFile, "import \"github.com/gibbonmi/bench/internal/gate\"\n\nfunc kitRoot() string { return gate.KitRoot() }\n",
		singleReadReport("reader.go", 5, "kitRoot", readBelowEntry, "gate.KitRoot", ""))
}

// TestSingleReadCensusCountsEachKindOfAConstructor proves newAmbient reads time.Now too.
func TestSingleReadCensusCountsEachKindOfAConstructor(t *testing.T) {
	t.Parallel()
	wantReadCensus(t, syntheticEffectsFile, "func Run(home string) {\n\ta := newAmbient(home)\n\t_, _ = a, currentTime()\n}\n",
		singleReadReport("reader.go", 5, "Run", readSecondTime, "", "time.Now"))
}

// TestSingleReadCensusRefusesAnEmptyReadSet proves a package without effectsFile is an
// error, so the census never passes with nothing to grade.
func TestSingleReadCensusRefusesAnEmptyReadSet(t *testing.T) {
	t.Parallel()
	dir := plantTestFiles(t, map[string]string{"reader.go": "package worktree\n\nfunc Run() {}\n"})
	reports, err := singleReadCensus(dir, plantTestFiles(t, map[string]string{"gate.go": syntheticGateFile}))
	if err == nil {
		t.Fatalf("single-read census without %s = %q, want an empty read set error", effectsFile, reports)
	}
}

// TestSingleReadCensusOnTheLiveTree proves each read in the package sits at an entry.
func TestSingleReadCensusOnTheLiveTree(t *testing.T) {
	t.Parallel()
	reports, err := singleReadCensus(".", filepath.Join("..", "gate"))
	if err != nil {
		t.Fatalf("single-read census the package: %v", err)
	}
	if len(reports) != 0 {
		t.Fatalf("the single-read census reports %d reads:\n%s", len(reports), strings.Join(reports, "\n"))
	}
}
