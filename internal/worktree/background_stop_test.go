// The stop of background grades when work goes away. A verb that removes a worktree, or that
// lands it as a stack member, first stops the running grade of that worktree and drops its
// queued grade. Each fixture has its own Bench home, which holds the machine queue that the verb
// reads.
package worktree

import (
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/testrepo"
)

// queueGrades starts a stand-in grade and writes the machine queue under home: the stand-in runs
// the grade of the worktree running, and the worktree queued waits with a grade of another key.
func queueGrades(t *testing.T, home, running, queued string) standInGrade {
	t.Helper()
	grade := startStandInGrade(t)
	queueGradesOf(t, grade, home, running, queued)
	return grade
}

// queueGradesOf writes the machine queue of queueGrades for the stand-in grade.
func queueGradesOf(t *testing.T, grade standInGrade, home, running, queued string) {
	t.Helper()
	tip := strings.Repeat("c", 40)
	slot := grade.slot(t, otherKey, running, tip)
	entry := backgroundstate.Entry{Key: strings.Repeat("b", 64), Worktree: queued, Tip: tip}
	mustNoError(t, backgroundstate.WriteQueue(durablefile.Native(), home, backgroundstate.Queue{Running: &slot, Queued: []backgroundstate.Entry{entry}}))
}

// requireGradesStopped fails t unless the stand-in grade has exited, the machine queue under
// home holds no grade, and no queue step failed. A queued grade that the verb kept would start
// in a worktree with no wrapper, and that start fails.
func requireGradesStopped(t *testing.T, home string, grade standInGrade, r verbResult) {
	t.Helper()
	alive, err := subprocess.Alive(grade.pid)
	mustNoError(t, err)
	queue, err := backgroundstate.ReadQueue(home)
	mustNoError(t, err)
	if r.exit != 0 || alive || queue.Running != nil || len(queue.Queued) != 0 || strings.Contains(r.stderr, "background grade:") {
		t.Fatalf("verb beside the grades = (%d, %q, %q), grade alive %t, queue %+v; want exit 0 after the stop with an empty queue", r.exit, r.stdout, r.stderr, alive, queue)
	}
}

// TestReleaseStopsTheGrade is BG78. Before the release removes its worktree, it stops the
// running grade of that worktree and drops its queued grade.
func TestReleaseStopsTheGrade(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "background-stop-release")
	grade := queueGrades(t, f.home, f.creation.Path, f.creation.Path)
	r := runVerb(t, verbRelease, f.call("--request", "landed-background-stop-release", f.creation.Path))
	requireTest(t, strings.Contains(r.stdout, ",complete,removed\n"), "release = (%d, %q, %q), want a removed row", r.exit, r.stdout, r.stderr)
	requireGradesStopped(t, f.home, grade, r)
}

// TestReleaseKeepsTheGradeOfAPrefixWorktree is BG-C8-C1. A second worktree whose path starts
// with the path of the released worktree is not under it, so the release keeps its running grade
// and its queued grade.
func TestReleaseKeepsTheGradeOfAPrefixWorktree(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "background-keep-prefix")
	prefix := f.creation.Path + "2"
	grade := queueGrades(t, f.home, prefix, prefix)
	before, err := backgroundstate.ReadQueue(f.home)
	mustNoError(t, err)
	r := runVerb(t, verbRelease, f.call("--request", "landed-background-keep-prefix", f.creation.Path))
	requireTest(t, strings.Contains(r.stdout, ",complete,removed\n"), "release = (%d, %q, %q), want a removed row", r.exit, r.stdout, r.stderr)
	after, err := backgroundstate.ReadQueue(f.home)
	mustNoError(t, err)
	requireTest(t, grade.alive() && reflect.DeepEqual(after, before) && !strings.Contains(r.stderr, "background grade:"), "release beside the grade of %s = (%d, %q), grade alive %t, queue %+v; want the grade and the queue %+v kept", prefix, r.exit, r.stderr, grade.alive(), after, before)
}

// TestCleanStopsTheGrade is BG79. Before the clean removes a worktree that it selected, it stops
// the running grade of that worktree and drops its queued grade.
func TestCleanStopsTheGrade(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "background-stop-clean")
	grade := queueGrades(t, f.home, f.creation.Path, f.creation.Path)
	plan := runVerb(t, verbClean, f.call("--target", f.creation.Assignment.ID))
	r := runVerb(t, verbClean, f.call("--target", f.creation.Assignment.ID, "--apply", plan.mustFingerprint(t)))
	requireTest(t, strings.Count(r.stdout, ",removed,") == 1, "clean = (%d, %q, %q), want one removal", r.exit, r.stdout, r.stderr)
	requireGradesStopped(t, f.home, grade, r)
}

// TestReclaimStopsTheGrade is BG80. Before the reclaim removes a pool key, it stops the running
// grade of each worktree under that key and drops the queued grade of each.
func TestReclaimStopsTheGrade(t *testing.T) {
	t.Parallel()
	f := newReclaimPool(t)
	plantDeadChild(t, f.pool, "dead-key", "wt")
	worktree := filepath.Join(f.pool, "dead-key", "wt")
	grade := queueGrades(t, f.home, worktree, worktree)
	plan := runVerb(t, verbReclaim, f.call())
	r := runVerb(t, verbReclaim, f.call("--apply", plan.mustFingerprint(t)))
	_, verdicts := reclaimVerdicts(t, r.stdout)
	requireTest(t, verdicts["dead-key"][0] == poolVerdictRemoved, "reclaim = (%d, %q, %q), want the dead key removed", r.exit, r.stdout, r.stderr)
	requireGradesStopped(t, f.home, grade, r)
}

// TestStackLandingStopsMemberGrades is BG81. Before the stack grades, it stops the running grade
// of a member worktree and drops the queued grade of another member. The fixture gate goes red
// when the grade still runs.
func TestStackLandingStopsMemberGrades(t *testing.T) {
	t.Parallel()
	grade := startStandInGrade(t)
	step := func(*testrepo.GateFixture, string) string {
		return "if kill -0 " + strconv.Itoa(grade.pid) + " 2>/dev/null; then exit 9; fi\n"
	}
	f := stackLandingFixture(t, "stack-background-stop", 2, step)
	queueGradesOf(t, grade, f.home, f.members[1].creation.Path, f.members[0].creation.Path)
	_, args := f.stack(t)
	r := runVerb(t, verbLand, f.call(args...))
	requireGradesStopped(t, f.home, grade, r)
}
