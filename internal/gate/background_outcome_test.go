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

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
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
	set, err := json.Marshal([]string{"./."})
	if err != nil {
		t.Fatal(err)
	}
	if schedules := gateSchedules(t, f); !slices.Equal(schedules, []string{wholeSchedule, string(set)}) {
		t.Fatalf("gate schedules = %q; want the prepared whole run and then the delta set %s", schedules, set)
	}
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
		t.Skip("the I/O class is a Linux priority")
	}
	if class := backgroundPriority(t)[1]; class != "idle" {
		t.Fatalf("I/O class of the background oracle = %q; want idle", class)
	}
}

// TestBackgroundRedWritesNoHold is BG36. A red whole run under a landing value writes the
// reader hold of its marker, and the same red as a background run writes none.
func TestBackgroundRedWritesNoHold(t *testing.T) {
	root := outcomeFixture(t)
	marker := outcomeGit(t, root, "rev-parse", "HEAD")
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
	f := markedBackgroundFixture(t, ".gate-sleep")
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

// backgroundPriority runs one background grade of a fixture in a child process, and answers
// the nice value and the I/O class that its oracle wrote.
func backgroundPriority(t *testing.T) []string {
	t.Helper()
	f := recordtest.Attach(t, outcomeFixture(t, func(g *testrepo.GateFixture, body string) string {
		class := ""
		if runtime.GOOS == "linux" {
			class = "$(" + g.Command("ionice") + " -p $$)"
		}
		return body + `printf '%s\n' "$(` + g.Command("ps") + ` -o nice= -p $$)" "` + class + `" > "$common/` + string(priorityWitness) + "\"\n"
	}), 1)
	f.AddChunk()
	f.Save()
	f.Commit("retain review evidence")
	seedGateCache(t, f)
	child := processtest.Helper(t, backgroundGradeHelper, kitRunRootEnv+"="+f.Root, backgroundResolvedEnv+"=1")
	if out, err := child.CombinedOutput(); err != nil {
		t.Fatalf("background grade child = %v\n%s", err, out)
	}
	lines := strings.Split(string(outcomeRead(t, priorityWitness.path(t, f.Root))), "\n")
	return []string{strings.TrimSpace(lines[0]), strings.TrimSpace(lines[1])}
}

// markedBackgroundFixture is a background fixture over the failure gate whose tip commits the
// gate marker.
func markedBackgroundFixture(t *testing.T, marker string) *recordtest.Fixture {
	t.Helper()
	f := recordtest.Attach(t, failureOutcomeFixture(t), 1)
	f.AddChunk()
	f.Save()
	f.Commit("retain review evidence")
	seedGateCache(t, f)
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
