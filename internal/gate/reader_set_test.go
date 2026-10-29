package gate

// The reader set. The phase-child tests drive the gate-phases process with a reader schedule
// over a kit-shaped fixture module. The run tests drive the reader run over the outcome
// fixture, whose gate exits 7 when the ignored file .gate-red exists.

import (
	"bytes"
	"context"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/durablefile"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// TestReaderScheduleRunsTheReaderSet proves that a red planted outside the reader
// set stays out of the run, the root-conformance entry runs under the root-conformance
// environment, the system suite runs, and the last line names the reader set.
func TestReaderScheduleRunsTheReaderSet(t *testing.T) {
	root, witness := readerSetFixture(t)
	code, out, diagnostic := runReaderSchedule(t, root, `{"packages":["./reader"],"race":false}`)
	if names := tablePhases(out); code != 0 || !slices.Equal(names, []string{canary.PhaseTest, SystemPhaseName}) || !strings.HasSuffix(out, "\ngate: green (reader set)\n") {
		t.Fatalf("reader set = %d %s %s; want the phases test and system and the reader verdict", code, out, diagnostic)
	}
	if got := string(outcomeRead(t, witness)); got != root {
		t.Fatalf("root-conformance witness = %q; want the entry run under the root %s", got, root)
	}
}

// TestReaderScheduleRunsTheRaceTable proves that the race flag adds the race phase, whose
// expected runs are every race-table sentinel.
func TestReaderScheduleRunsTheRaceTable(t *testing.T) {
	root, _ := readerSetFixture(t)
	code, out, diagnostic := runReaderSchedule(t, root, `{"packages":["./reader"],"race":true}`)
	if names := tablePhases(out); code != 0 || !slices.Equal(names, []string{canary.PhaseTest, canary.PhaseRace, SystemPhaseName}) {
		t.Fatalf("reader set with the race flag = %d %s %s; want the phases test, race, and system", code, out, diagnostic)
	}
}

// TestReaderScheduleRunsNoSourcePhase proves that the kit table of the fixture holds the
// gofmt, vet, and shellcheck phases, and the reader set holds none of them.
func TestReaderScheduleRunsNoSourcePhase(t *testing.T) {
	root, _ := readerSetFixture(t)
	source := []string{canary.PhaseGofmt, canary.PhaseVet, "shellcheck"}
	isSource := func(phase Phase) bool { return slices.Contains(source, phase.Name) }
	if whole := BenchkitPhases(root, root); slices.IndexFunc(whole, isSource) < 0 {
		t.Fatalf("kit table = %+v; want the source phases", whole)
	}
	t.Setenv(readerScheduleEnv, `{"packages":["./reader"],"race":true}`)
	phases, err := scheduledTable(root, root)
	if err != nil || slices.ContainsFunc(phases, isSource) {
		t.Fatalf("reader set = %+v, %v; want no source phase", phases, err)
	}
}

// TestReaderScheduleRefusesAMalformedSchedule proves that each malformed schedule refuses
// before any phase prints.
func TestReaderScheduleRefusesAMalformedSchedule(t *testing.T) {
	root, _ := readerSetFixture(t)
	for _, schedule := range []string{
		`./reader`,
		`{"packages":["./reader"]}`,
		`{"packages":["./reader"],"race":false,"extra":1}`,
		`{"packages":["./reader"],"race":false}{}`,
		`{"packages":["-run=x"],"race":false}`,
	} {
		if code, out, _ := runReaderSchedule(t, root, schedule); code == 0 || out != "" {
			t.Fatalf("malformed reader schedule %s = %d %s; want a refusal before any phase", schedule, code, out)
		}
	}
}

// TestReaderGreenWritesOnlyTheReaderRecord proves that after a whole green, a green reader
// run adds one reader record for its tree and leaves the verdict cache and the evidence store
// byte-identical.
func TestReaderGreenWritesOnlyTheReaderRecord(t *testing.T) {
	root := wholeGreenFixture(t)
	before := gateStoreBytes(t, root)
	result, tree := runReaderOn(t, root)
	if result.ActionExit != 0 || result.Schedule != ScheduleReader {
		t.Fatalf("reader run = %+v; want a green reader run", result)
	}
	if after := gateStoreBytes(t, root); len(after) == 0 || !maps.Equal(before, after) {
		t.Fatalf("verdict cache and evidence after a reader green = %q; want the bytes %q", after, before)
	}
	entries, err := os.ReadDir(filepath.Join(gittest.Output(t, root, "rev-parse", "--path-format=absolute", "--git-common-dir"), benchgit.ReaderRecordDir))
	if err != nil || len(entries) != 1 || entries[0].Name() != tree {
		t.Fatalf("reader records = %v, %v; want the one record of %s", entries, err, tree)
	}
	if recorded, err := ReaderRecorded(root, tree); !recorded || err != nil {
		t.Fatalf("ReaderRecorded after a reader green = %t, %v; want true", recorded, err)
	}
}

// TestReaderRecordIsNeverReusableGreen proves that a reader record as the verdict of the
// subject tree reads as a well-formed record that its class refuses to reuse.
func TestReaderRecordIsNeverReusableGreen(t *testing.T) {
	root := outcomeFixture(t)
	subject, err := buildSubject(root)
	if err != nil {
		t.Fatal(err)
	}
	writeInspectRecord(t, root, readerRecordBytes(t, subject.Tree))
	if got := Inspect(root); got.State != Ready || got.ReusableGreen || got.Reason == "" {
		t.Fatalf("Inspect of a reader record = %+v; want a ready record with a reuse refusal", got)
	}
}

// TestReaderRedWritesNoRecord proves that the planted red reds the reader run, which writes no
// reader record.
func TestReaderRedWritesNoRecord(t *testing.T) {
	root := wholeGreenFixture(t)
	outcomeWrite(t, root, ".gate-red", "\n", 0o644)
	result, tree := runReaderOn(t, root)
	if result.GateExit != 7 || result.Schedule != ScheduleReader {
		t.Fatalf("reader run = %+v; want the red reader run", result)
	}
	if recorded, err := ReaderRecorded(root, tree); recorded || err != nil {
		t.Fatalf("ReaderRecorded after a red reader run = %t, %v; want false", recorded, err)
	}
}

// TestReaderRecordedAnswersForItsTree proves that the query answers true for the tree of a
// valid record, false for another tree, and false for an invalid record of a tree.
func TestReaderRecordedAnswersForItsTree(t *testing.T) {
	root := gittest.Repo(t)
	gitdir, err := benchgit.CommonDir(root)
	if err != nil {
		t.Fatal(err)
	}
	recorded, other, invalid := strings.Repeat("a", 40), strings.Repeat("b", 40), strings.Repeat("c", 40)
	run := newReaderRun(subject{Tree: recorded}, readerSchedule{Packages: []string{"./reader"}}, other, other)
	if err := writeReaderGreen(durablefile.Native(), gitdir, run, strings.Repeat("d", 64), time.Now()); err != nil {
		t.Fatal(err)
	}
	run.plan.Tree = invalid
	if err := writeReaderGreen(durablefile.Native(), gitdir, run, "not a content address", time.Now()); err != nil {
		t.Fatal(err)
	}
	for tree, want := range map[string]bool{recorded: true, other: false, invalid: false} {
		if got, err := ReaderRecorded(root, tree); got != want || err != nil {
			t.Errorf("ReaderRecorded(%s) = %t, %v; want %t", tree, got, err, want)
		}
	}
}

// TestPhaseChildrenDropTheReaderSchedule proves that no phase child of a reader set runs a
// reader set of its own.
func TestPhaseChildrenDropTheReaderSchedule(t *testing.T) {
	t.Setenv(readerScheduleEnv, `{"packages":[],"race":false}`)
	env, err := gateEnv()
	if err != nil || slices.ContainsFunc(env, func(entry string) bool { return strings.HasPrefix(entry, readerScheduleEnv+"=") }) {
		t.Fatalf("phase child environment = %q, %v; want no reader schedule", env, err)
	}
}

// readerWitnessEnv names the file where the fixture's root-conformance entry writes the root
// that its environment names.
const readerWitnessEnv = "BENCH_TEST_READER_WITNESS"

// readerSetFixture is a kit-shaped fixture module: a reader package, a planted red outside every
// reader set, the root-conformance entry, the system suite, and each race-table sentinel. It
// answers the root and the witness path of the root-conformance entry.
func readerSetFixture(t *testing.T) (string, string) {
	t.Setenv(baselinePolicyEnv, "")
	witness := filepath.Join(t.TempDir(), "witness")
	t.Setenv(readerWitnessEnv, witness)
	root := gittest.Repo(t)
	operands, _ := SystemSuite(root)
	suite := strings.TrimPrefix(operands[slices.IndexFunc(operands, func(operand string) bool { return strings.HasPrefix(operand, "./") })], "./")
	files := map[string]string{
		"go.mod":                 "module fixture\n\ngo 1.21\n",
		"reader/reader_test.go":  goTest("reader", "TestReader", ""),
		"other/other_test.go":    goTest("other", "TestPlanted", `t.Fatal("planted red")`),
		suite + "/suite_test.go": goTest(path.Base(suite), "TestSuite", ""),
		conformancePackagePath + "/entry_test.go": "package conformance\n\nimport (\n\t\"os\"\n\t\"testing\"\n)\n\nfunc " + registry.RootConformanceTest +
			"(t *testing.T) {\n\tif err := os.WriteFile(os.Getenv(\"" + readerWitnessEnv + "\"), []byte(os.Getenv(\"" + registry.ConformanceRootEnv + "\")), 0o600); err != nil {\n\t\tt.Fatal(err)\n\t}\n}\n",
	}
	for _, test := range raceTests {
		dir := strings.TrimPrefix(test.PackagePath, "./")
		files[dir+"/"+test.Name+"_test.go"] = goTest(path.Base(dir), test.Name, "")
	}
	for name, body := range files {
		outcomeWrite(t, root, name, body, 0o644)
	}
	// The stub run binary exits 0 for every step.
	outcomeWrite(t, root, "bench", "#!/bin/sh\nexit 0\n", 0o755)
	return root, witness
}

// runReaderSchedule runs the phase command of root, as the kit, with the reader schedule. It
// answers the exit, stdout, and stderr.
func runReaderSchedule(t *testing.T, root, schedule string) (int, string, string) {
	t.Helper()
	t.Setenv(readerScheduleEnv, schedule)
	return runFixturePhases(context.Background(), t, root)
}

// wholeGreenFixture is the outcome fixture after a whole green, so its gate finds the verdict
// cache that it copies.
func wholeGreenFixture(t *testing.T) string {
	t.Helper()
	root := outcomeFixture(t)
	if result := Execute(context.Background(), root, &bytes.Buffer{}, &bytes.Buffer{}); result.ActionExit != 0 {
		t.Fatalf("whole green = %+v", result)
	}
	return root
}

// runReaderOn runs a reader set of the subject of root, with HEAD as the green marker and its
// tree as the destination. It answers the result and the graded tree.
func runReaderOn(t *testing.T, root string) (Result, string) {
	t.Helper()
	evaluation := newGateEvaluation(root)
	plan, err := evaluation.acceptPre()
	if err != nil {
		t.Fatal(err)
	}
	run := newReaderRun(plan, readerSchedule{Packages: []string{"./reader"}}, gittest.Output(t, root, "rev-parse", "HEAD^{tree}"), gittest.Output(t, root, "rev-parse", "HEAD"))
	ctx := context.Background()
	release := func() error { return nil }
	result := runReader(ctx, ctx, root, root, nil, run, evaluation, nil, stubSelection(t, ""), release, durablefile.Native(), &bytes.Buffer{}, &bytes.Buffer{})
	return result, plan.Tree
}

// stubSelection answers a run binary selection with kit as its source root, whose executable
// exits 0 for every step.
func stubSelection(t *testing.T, kit string) *runbinary.Selection {
	t.Helper()
	dir := t.TempDir()
	outcomeWrite(t, dir, "bench", "#!/bin/sh\nexit 0\n", 0o755)
	return &runbinary.Selection{Path: filepath.Join(dir, "bench"), SourceRoot: kit}
}

// readerRecordBytes answers the bytes of a valid reader record of tree, as the reader-record
// writer encodes them, without the final newline.
func readerRecordBytes(t *testing.T, tree string) []byte {
	t.Helper()
	gitdir := t.TempDir()
	run := newReaderRun(subject{Tree: tree}, readerSchedule{Packages: []string{"./reader"}}, tree, tree)
	if err := writeReaderGreen(durablefile.Native(), gitdir, run, strings.Repeat("a", 64), time.Now().Add(-time.Minute)); err != nil {
		t.Fatal(err)
	}
	return bytes.TrimSuffix(outcomeRead(t, filepath.Join(gitdir, benchgit.ReaderRecordDir, tree)), []byte("\n"))
}

// gateStoreBytes answers the bytes of the verdict cache and of each evidence entry of root,
// keyed by path.
func gateStoreBytes(t *testing.T, root string) map[string]string {
	t.Helper()
	common := gittest.Output(t, root, "rev-parse", "--path-format=absolute", "--git-common-dir")
	paths := []string{filepath.Join(gittest.Output(t, root, "rev-parse", "--absolute-git-dir"), benchgit.GateCacheFile)}
	evidence, err := filepath.Glob(filepath.Join(filepath.Dir(evidencePath(common, subject{})), "*"))
	if err != nil {
		t.Fatal(err)
	}
	store := map[string]string{}
	for _, file := range append(paths, evidence...) {
		store[file] = string(outcomeRead(t, file))
	}
	return store
}
