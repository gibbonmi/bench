package gate

// The wait of a chunk checkpoint on the background grade of its own key, and the steps that clear
// the background grades of its worktree before its own run. Each test has a private Bench home,
// which the lane passes and the checkpoint both read, so its queue is the one machine queue of the
// test.

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
)

// TestChunkCheckpointWaitsForItsKey is BG57. The checkpoint finds the live claim of a real grade
// on its key, waits until that grade ends, and then reuses its green with no oracle run of its own.
func TestChunkCheckpointWaitsForItsKey(t *testing.T) {
	w := holdGrade(t)
	w.awaitWait(t)
	w.release(t)
	w.end.await(t)
	if w.end.code != 0 || !strings.Contains(w.end.out, "reused") || outcomeRuns(t, w.root) != w.runs {
		t.Fatalf("checkpoint after the wait = %d %s after %d runs; want the reused green after the %d runs of the grade", w.end.code, w.end.out, outcomeRuns(t, w.root), w.runs)
	}
}

// TestWaitPromotesTheIOClass is BG58. When the wait starts, each thread of the grade and of each
// process that the grade started is in the best-effort class, the oracle group included.
func TestWaitPromotesTheIOClass(t *testing.T) {
	w := holdGrade(t)
	w.awaitWait(t)
	processtest.WaitFor(t, "each thread of the grade and of its oracle group to be best-effort", func() bool {
		threads, oracle := gradeThreads(t, w.grade.Owner)
		for _, tid := range threads {
			if class, err := ioClass(tid); err == nil && class != "best-effort" {
				return false
			}
		}
		return oracle
	})
}

// TestWaitKeepsTheNiceValue is BG59. During the wait, the grade keeps the nice value that it
// lowered itself to.
func TestWaitKeepsTheNiceValue(t *testing.T) {
	w := holdGrade(t)
	w.awaitWait(t)
	// The raw Linux syscall answers 20 minus the nice value.
	if priority, err := syscall.Getpriority(syscall.PRIO_PROCESS, w.grade.Owner); err != nil || 20-priority != backgroundNice {
		t.Fatalf("nice value of the waited grade = %d (%v); want %d", 20-priority, err, backgroundNice)
	}
}

// TestWaitCountsAgainstTheGateTimeout is BG60. A claim whose owner never ends holds the wait past
// the gate timeout, so the checkpoint ends with exit 124 and runs no oracle.
func TestWaitCountsAgainstTheGateTimeout(t *testing.T) {
	privateHome(t)
	f := backgroundFixture(t)
	claim, err := backgroundstate.ClaimFor(strippedEvidence(t, f).Grade, f.Root)
	if err != nil || backgroundstate.WriteClaim(durablefile.Native(), commonDir(t, f.Root), claim) != nil {
		t.Fatalf("claim = %+v (%v); want it written", claim, err)
	}
	runs := outcomeRuns(t, f.Root)
	previous := gateTimeout
	gateTimeout = time.Second
	t.Cleanup(func() { gateTimeout = previous })
	if code, out := runCheckpoint(t, f); code != 124 || outcomeRuns(t, f.Root) != runs {
		t.Fatalf("checkpoint on a claim that never ends = %d %s after %d runs; want exit 124 after %d runs", code, out, outcomeRuns(t, f.Root), runs)
	}
}

// TestWaitHoldsNoCheckoutLock is BG61. While the checkpoint waits, the checkout lock is free.
func TestWaitHoldsNoCheckoutLock(t *testing.T) {
	w := holdGrade(t)
	w.awaitWait(t)
	if held, err := ExecutionInProgress(w.root); err != nil || held {
		t.Fatalf("checkout lock during the wait held = %t (%v); want it free", held, err)
	}
}

// TestCheckpointRunsAfterABackgroundRed is BG62. The waited grade ends red, and the checkpoint then
// runs its own oracle and answers its own exit, not a refusal.
func TestCheckpointRunsAfterABackgroundRed(t *testing.T) {
	w := holdGrade(t, func(f *recordtest.Fixture) {
		f.Write(".gate-red", "")
		f.Git("add", "-f", ".gate-red")
		f.Commit("commit the red marker")
	})
	w.awaitWait(t)
	w.release(t)
	w.end.await(t)
	awaitExit(t, w.grade.Owner)
	if w.end.code != 7 || outcomeRuns(t, w.root) != w.runs+1 || backgroundRecord(t, w.f).Outcome != backgroundstate.Red {
		t.Fatalf("checkpoint after a background red = %d %s after %d runs, record %+v; want its own red after %d runs", w.end.code, w.end.out, outcomeRuns(t, w.root), backgroundRecord(t, w.f), w.runs+1)
	}
}

// TestQueuedSameKeyLeavesTheQueue is BG63. A grade of the checkpoint's key waits behind another
// worktree's grade, and the checkpoint's own run takes it out of the queue.
func TestQueuedSameKeyLeavesTheQueue(t *testing.T) {
	home := privateHome(t)
	first := standInWorktree(t, home)
	first.running(t)
	w := recordedWorktree(t, home)
	w.pass(t)
	requireQueue(t, home, entryOf(t, first), entryOf(t, w))
	requireOwnRun(t, w)
	requireQueue(t, home, entryOf(t, first))
}

// TestCheckpointStopsOtherKeyGrades is BG64. The running grade of an older tip of the worktree holds
// another key, and the checkpoint stops it and waits for its exit before its own run.
func TestCheckpointStopsOtherKeyGrades(t *testing.T) {
	w := recordedWorktree(t, privateHome(t))
	old := w.running(t)
	w.supersede(t)
	rerecord(t, w.f)
	requireOwnRun(t, w)
	requireAlive(t, old.PID, false)
	requireQueue(t, w.home, backgroundstate.Entry{})
}

// TestCheckpointDropsOtherKeyQueue is BG65. The queued grade of an older tip of the worktree holds
// another key, and the checkpoint drops it before its own run.
func TestCheckpointDropsOtherKeyQueue(t *testing.T) {
	home := privateHome(t)
	first := standInWorktree(t, home)
	first.running(t)
	w := recordedWorktree(t, home)
	w.pass(t)
	w.supersede(t)
	rerecord(t, w.f)
	requireOwnRun(t, w)
	requireQueue(t, home, entryOf(t, first))
}

// privateHome answers a Bench home for one test. The checkpoints of that test read it too.
func privateHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv(benchhome.Env, home)
	return home
}

// rerecord records the chunk of f again at its tip and commits the record, so the checkpoint of
// that tip passes its evidence pass after a commit that followed the record.
func rerecord(t *testing.T, f *recordtest.Fixture) {
	t.Helper()
	base := f.Record.Chunks[0].Base
	f.Record.Chunks = f.Record.Chunks[:0]
	f.RecordChunk(base)
	resave(t, f, "record the chunk at the tip")
}

// recordedWorktree is a keyed worktree whose chunk is recorded at its tip and whose worktree holds
// a gate cache, so its chunk checkpoint passes the evidence pass and runs its own oracle.
func recordedWorktree(t *testing.T, home string) queueWorktree {
	t.Helper()
	w := keyedWorktree(t, home)
	rerecord(t, w.f)
	seedGateCache(t, w.f)
	return w
}

// requireOwnRun requires that the chunk checkpoint of w runs its own oracle once and ends green.
func requireOwnRun(t *testing.T, w queueWorktree) {
	t.Helper()
	runs := outcomeRuns(t, w.root)
	if code, out := runCheckpoint(t, w.f); code != 0 || outcomeRuns(t, w.root) != runs+1 {
		t.Fatalf("checkpoint = %d %s after %d runs; want its own green after %d runs", code, out, outcomeRuns(t, w.root), runs+1)
	}
}

// waitedGrade is a real grade of a holding worktree whose oracle holds, the count of oracle runs
// once it holds, and the chunk checkpoint of the same tip that runs beside it.
type waitedGrade struct {
	queueWorktree
	grade backgroundstate.Slot
	runs  int
	end   *checkpointEnd
}

// checkpointEnd is the exit and the output of a chunk checkpoint that runs beside a test. done
// closes when the checkpoint ends.
type checkpointEnd struct {
	code int
	out  string
	done chan struct{}
}

// finished reports whether the checkpoint has ended.
func (e *checkpointEnd) finished() bool {
	select {
	case <-e.done:
		return true
	default:
		return false
	}
}

// await waits until the checkpoint has ended.
func (e *checkpointEnd) await(t *testing.T) {
	t.Helper()
	processtest.WaitFor(t, "the checkpoint to end", e.finished)
}

// holdGrade starts the real grade of a holding worktree whose tip each prepare step commits, waits
// until its oracle holds, and then starts the chunk checkpoint of that tip. The cleanup releases the
// oracle and waits for the checkpoint to end. The test observes the wait through the I/O class of
// the grade, which is a Linux priority.
func holdGrade(t *testing.T, prepare ...func(*recordtest.Fixture)) waitedGrade {
	t.Helper()
	if runtime.GOOS != "linux" {
		capability.Environment(t, "the test observes the wait through the I/O class, a Linux priority")
	}
	h := holdingWorktree(t, privateHome(t))
	for _, apply := range prepare {
		apply(h.f)
	}
	rerecord(t, h.f)
	h.pass(t)
	w := waitedGrade{queueWorktree: h, grade: h.awaitHolding(t), runs: outcomeRuns(t, h.root), end: &checkpointEnd{done: make(chan struct{})}}
	go func() {
		w.end.code, w.end.out = runCheckpoint(t, h.f)
		close(w.end.done)
	}()
	t.Cleanup(func() { h.release(t); w.end.await(t) })
	return w
}

// awaitWait waits until the checkpoint of w starts its wait: the main thread of the grade leaves
// the idle class. A checkpoint that ends first did not wait.
func (w waitedGrade) awaitWait(t *testing.T) {
	t.Helper()
	processtest.WaitFor(t, "the checkpoint to start its wait", func() bool {
		class, err := ioClass(w.grade.Owner)
		return w.end.finished() || (err == nil && class == "best-effort")
	})
	if w.end.finished() {
		t.Fatalf("checkpoint ended before the grade = %d %s; want a wait", w.end.code, w.end.out)
	}
}

// gradeThreads answers each thread of process pid and of each of its descendants, which the
// children list of each thread names, and whether a descendant is in a process group other than
// the group of pid. A process that exits during the walk names no thread.
func gradeThreads(t *testing.T, pid int) ([]int, bool) {
	t.Helper()
	group, err := syscall.Getpgid(pid)
	if err != nil {
		t.Fatal(err)
	}
	var threads []int
	other := false
	for pending := []int{pid}; len(pending) > 0; pending = pending[1:] {
		dir := filepath.Join("/proc", strconv.Itoa(pending[0]), "task")
		tasks, _ := os.ReadDir(dir)
		for _, task := range tasks {
			tid, err := strconv.Atoi(task.Name())
			if err != nil {
				continue
			}
			threads = append(threads, tid)
			children, _ := os.ReadFile(filepath.Join(dir, task.Name(), "children"))
			for _, field := range strings.Fields(string(children)) {
				child, err := strconv.Atoi(field)
				if err != nil {
					continue
				}
				if childGroup, err := syscall.Getpgid(child); err == nil && childGroup != group {
					other = true
				}
				pending = append(pending, child)
			}
		}
	}
	return threads, other
}
