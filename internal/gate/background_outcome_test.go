package gate

// The outputs of one background run: the delta record, the reader table, no reader hold, the
// priority of its oracle, and the aborted outcome of each end without a verdict.

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
	"github.com/gibbonmi/bench/internal/testrepo"
)

// backgroundGradeHelper names the helper that runs one background grade of the HEAD tip of
// the root that kitRunRootEnv names, through the resolved branch of the entry. The grade lowers
// the priority of its own process, so it runs in a child and not in this test process.
const backgroundGradeHelper = "background-grade"

// priorityWitness is the file where the fixture gate writes its own nice value and, on Linux,
// its I/O class.
const priorityWitness outcomeWitness = ".gate-priority"

// backgroundGradeChild runs the background grade entry as the wrapper of the worktree resolved
// it.
func backgroundGradeChild() int {
	return BackgroundCommand([]string{os.Getenv(kitRunRootEnv), "HEAD"}, os.Stdout, os.Stderr)
}

// TestBackgroundGradeTakesTheDelta is BG25. After a whole green of the reviewed tip, a
// test-only change makes the background grade run the delta set and write the delta record of
// the worktree.
func TestBackgroundGradeTakesTheDelta(t *testing.T) {
	f := deltaFixture(t, graded(t, modifyTest))
	if code, out := gradeBackground(t, f, nil); code != 0 {
		t.Fatalf("background grade = %d %s; want green", code, out)
	}
	requireDeltaSet(t, f, "./.")
	if _, err := os.Stat(deltaRecordPath(t, f)); err != nil {
		t.Fatalf("delta record after a background delta green: %v; want the record", err)
	}
}

// TestBackgroundOracleRunsAtNiceNineteen is BG30.
func TestBackgroundOracleRunsAtNiceNineteen(t *testing.T) {
	if nice := backgroundPriority(t)[0]; nice != "19" {
		t.Fatalf("nice value of the background oracle = %q; want 19", nice)
	}
}

// TestBackgroundOracleRunsInTheIdleClass is BG31.
func TestBackgroundOracleRunsInTheIdleClass(t *testing.T) {
	if runtime.GOOS != "linux" {
		capability.Environment(t, "the I/O class is a Linux priority")
	}
	if class := backgroundPriority(t)[1]; class != "idle" {
		t.Fatalf("I/O class of the background oracle = %q; want idle", class)
	}
}

// TestBackgroundRedWritesNoHold is BG36. A red whole run under a landing value writes the
// reader hold of its marker, and the same red as a background run writes none.
func TestBackgroundRedWritesNoHold(t *testing.T) {
	root := outcomeFixture(t)
	marker := gittest.Output(t, root, "rev-parse", "HEAD")
	outcomeWrite(t, root, ".gate-red", "\n", 0o644)
	if result := runBackgroundSubject(t, WithLanding(context.Background(), marker, marker), root, nil); result.ActionExit != 7 {
		t.Fatalf("background landing run = %+v; want the planted red", result)
	}
	if held, err := readerHeld(root, marker); err != nil || held {
		t.Fatalf("reader hold of %s after a background red = %t %v; want no hold", marker, held, err)
	}
}

// TestBackgroundGreenWritesTheReaderTable is BG37.
func TestBackgroundGreenWritesTheReaderTable(t *testing.T) {
	root := readerFixture(t)
	writeWitnessLog(t, root)
	owner := func(context.Context, string) (*runbinary.Selection, error) { return fixtureSelection(root), nil }
	if result := runBackgroundSubject(t, context.Background(), root, owner); result.ActionExit != 0 {
		t.Fatalf("background whole run = %+v; want green", result)
	}
	if tree, tables := currentSubject(t, root).Tree, storedTables(t, root); !slices.Equal(tables, []string{tree}) {
		t.Fatalf("stored tables after a background whole green = %q; want the one table of %s", tables, tree)
	}
}

// TestTimedOutGradeRecordsAborted is BG38.
func TestTimedOutGradeRecordsAborted(t *testing.T) {
	f := markedBackgroundFixture(t, failureOutcomeFixture(t), ".gate-sleep")
	previous := gateTimeout
	gateTimeout = time.Second
	t.Cleanup(func() { gateTimeout = previous })
	code, out := gradeBackground(t, f, nil)
	if got := backgroundRecord(t, f); code != 124 || got.Outcome != backgroundstate.Aborted || got.Key != strippedEvidence(t, f).Grade {
		t.Fatalf("timed-out grade = %d %s, record = %+v; want exit 124 and an aborted record of its key", code, out, got)
	}
}

// TestInterruptedGradeRecordsAborted is BG39. The grade is cancelled once it holds its claim.
func TestInterruptedGradeRecordsAborted(t *testing.T) {
	f := backgroundFixture(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	arm := func(ctx context.Context) (context.Context, func()) { cancel(); return ctx, func() {} }
	var out bytes.Buffer
	code := gradeInBackground(ctx, f.Root, f.Tip(), arm, &out, &out)
	if got := backgroundRecord(t, f); got.Outcome != backgroundstate.Aborted || got.Key != strippedEvidence(t, f).Grade {
		t.Fatalf("interrupted grade = %d %s, record = %+v; want an aborted record of its key", code, out.String(), got)
	}
}

// TestRefusedGradeRecordsAborted is BG40. A directory at the claim path of the grade key makes
// the claim unwritable, so the grade cannot hold its subject, records aborted, and runs no
// oracle.
func TestRefusedGradeRecordsAborted(t *testing.T) {
	f := backgroundFixture(t)
	key := strippedEvidence(t, f).Grade
	if err := os.MkdirAll(filepath.Join(backgroundstate.ClaimPath(commonDir(t, f.Root), key), "held"), 0o755); err != nil {
		t.Fatal(err)
	}
	runs := outcomeRuns(t, f.Root)
	code, out := gradeBackground(t, f, nil)
	want := backgroundstate.Record{Outcome: backgroundstate.Aborted, Key: key, Tip: f.Tip(), Log: backgroundstate.LogPath(adminDir(t, f.Root))}
	if got := backgroundRecord(t, f); got != want || outcomeRuns(t, f.Root) != runs {
		t.Fatalf("refused grade = %d %s, record = %+v after %d runs; want %+v after %d runs", code, out, got, outcomeRuns(t, f.Root), want, runs)
	}
}

// TestKeylessGradeRecordsAborted is BG-C5-P1. A grade that ends before it knows its grade key
// records aborted with an empty key, its tip, and its log: here a tip with no open build, and a
// tip that does not resolve.
func TestKeylessGradeRecordsAborted(t *testing.T) {
	root := outcomeFixture(t)
	head := gittest.Output(t, root, "rev-parse", "HEAD")
	for _, tip := range []string{head, "no-such-tip"} {
		var out bytes.Buffer
		code := gradeInBackground(context.Background(), root, tip, nil, &out, &out)
		record, err := backgroundstate.ReadRecord(adminDir(t, root))
		want := backgroundstate.Record{Outcome: backgroundstate.Aborted, Tip: tip, Log: backgroundstate.LogPath(adminDir(t, root))}
		if err != nil || record != want {
			t.Fatalf("keyless grade of %s = %d %s, record = %+v %v; want %+v", tip, code, out.String(), record, err, want)
		}
	}
}

// TestMovedRedGradeRecordsAborted is BG-C5-P2. The oracle moves its subject and then exits red,
// so the settle step refuses the run, and the grade records aborted under its key.
func TestMovedRedGradeRecordsAborted(t *testing.T) {
	f := markedBackgroundFixture(t, outcomeFixture(t, func(_ *testrepo.GateFixture, body string) string {
		return body + "if [ -e .gate-moved-red ]; then printf 'moved\\n' >> tracked.txt; exit 7; fi\n"
	}), ".gate-moved-red")
	code, out := gradeBackground(t, f, nil)
	if got := backgroundRecord(t, f); got.Outcome != backgroundstate.Aborted || got.Key != strippedEvidence(t, f).Grade {
		t.Fatalf("moved red grade = %d %s, record = %+v; want an aborted record of its key", code, out, got)
	}
}

// TestLaterGradeKeepsTheRecordedLog is BG-C8-P2 under story 41. A red grade records its log.
// While a later grade of the worktree runs, the log that the red record names still holds the
// output of the red grade. The fixture oracle prints its run count, so each grade prints its own
// output.
func TestLaterGradeKeepsTheRecordedLog(t *testing.T) {
	f := markedBackgroundFixture(t, outcomeFixture(t, func(_ *testrepo.GateFixture, body string) string {
		return body + "printf 'oracle run %s\\n' \"$((count + 1))\"\n"
	}), ".gate-red")
	gradeBackground(t, f, nil)
	record := backgroundRecord(t, f)
	red := outcomeRead(t, record.Log)
	var during []byte
	gradeBackground(t, f, func() { during, _ = os.ReadFile(record.Log) })
	if record.Outcome != backgroundstate.Red || len(red) == 0 || !bytes.Equal(during, red) {
		t.Fatalf("log of the %s record during a later grade = %q; want the output of that grade %q", record.Outcome, during, red)
	}
}

// backgroundPriority runs one background grade of a fixture in a child process, and answers
// the nice value and the I/O class that its oracle wrote.
func backgroundPriority(t *testing.T) []string {
	t.Helper()
	f := chunkBackgroundFixture(t, outcomeFixture(t, func(g *testrepo.GateFixture, body string) string {
		class := ""
		if runtime.GOOS == "linux" {
			class = "$(" + g.Command("ionice") + " -p $$)"
		}
		return body + `printf '%s\n' "$(` + g.Command("ps") + ` -o nice= -p $$)" "` + class + `" > "$common/` + string(priorityWitness) + "\"\n"
	}))
	child := processtest.Helper(t, backgroundGradeHelper, kitRunRootEnv+"="+f.Root, backgroundResolvedEnv+"=1")
	if out, err := child.CombinedOutput(); err != nil {
		t.Fatalf("background grade child = %v\n%s", err, out)
	}
	lines := strings.Split(string(outcomeRead(t, priorityWitness.path(t, f.Root))), "\n")
	return []string{strings.TrimSpace(lines[0]), strings.TrimSpace(lines[1])}
}

// chunkBackgroundFixture attaches one chunk entry of a review record to the outcome fixture root,
// and seeds its gate cache, so a background grade of its tip grades a chunk subject.
func chunkBackgroundFixture(t *testing.T, root string) *recordtest.Fixture {
	t.Helper()
	f := recordtest.Attach(t, root, 1)
	f.AddChunk()
	f.Save()
	f.Commit("retain review evidence")
	seedGateCache(t, f)
	return f
}

// requireDeltaSet requires that the gate of f ran one whole run and then one run of the delta set
// of packages.
func requireDeltaSet(t *testing.T, f *recordtest.Fixture, packages ...string) {
	t.Helper()
	set, err := json.Marshal(packages)
	if err != nil {
		t.Fatal(err)
	}
	if schedules := gateSchedules(t, f); !slices.Equal(schedules, []string{wholeSchedule, string(set)}) {
		t.Fatalf("gate schedules = %q; want one whole run and then the delta set %s", schedules, set)
	}
}

// markedBackgroundFixture is a background fixture over the outcome fixture root whose tip commits
// the gate marker.
func markedBackgroundFixture(t *testing.T, root, marker string) *recordtest.Fixture {
	t.Helper()
	f := chunkBackgroundFixture(t, root)
	f.Write(marker, "")
	f.Git("add", "-f", marker)
	f.Commit("commit the gate marker")
	return f
}

// runBackgroundSubject runs the worktree subject of root once as a background run under ctx,
// with owner as the run-binary owner, and answers its result. The fixture gate copies the gate
// cache, and a background run writes none, so a placeholder cache stands in for the cache of a
// first foreground gate. A foreground green would store the reader table that BG37 reads.
func runBackgroundSubject(t *testing.T, ctx context.Context, root string, owner runBinaryOwner) Result {
	t.Helper()
	outcomeWrite(t, filepath.Dir(outcomeCachePath(t, root)), filepath.Base(outcomeCachePath(t, root)), "{}\n", 0o600)
	var out bytes.Buffer
	ctx = WithCheckpoint(ctx, Checkpoint{Background: &backgroundGrade{}})
	result := executeSubjectReplacing(ctx, root, root, &out, &out, nil, reuseFreshGreen, newGateEvaluation(root), owner, "", durablefile.Native())
	t.Logf("background run exit %d: %q", result.ActionExit, out.String())
	return result
}
