package gate

// The machine queue of the background grade. The worktrees of a test share one private Bench
// home, so its queue is the one machine queue of the test. A stand-in worktree holds no open
// build, so its tip has no grade key, and a keyed worktree is a checkpoint fixture.

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
	"github.com/gibbonmi/bench/internal/gate/prospectiveartifact"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
)

// queueWorktree is one worktree: its root, its queue home, the temporary directory of a real
// grade, the stand-ins that its wrapper records, and the checkpoint fixture of a keyed worktree.
type queueWorktree struct {
	root, home, tmp string
	launches        *processtest.Launches
	f               *recordtest.Fixture
}

// standInWorktree answers a worktree with no open build under the queue of home.
func standInWorktree(t *testing.T, home string) queueWorktree {
	t.Helper()
	f := newLaunchFixture(t, true)
	return queueWorktree{root: f.root, home: home, launches: f.launches}
}

// keyedWorktree answers a checkpoint fixture under the queue of home, with its committed wrapper.
func keyedWorktree(t *testing.T, home string) queueWorktree {
	t.Helper()
	f := checkpointFixture(t)
	w := queueWorktree{root: f.Root, home: home, launches: processtest.NewLaunches(t), f: f}
	w.wrap(t, w.launches.StandIn())
	f.Commit("wrapper")
	commitmenttest.Register(t, f.Root, "keyed")
	return w
}

// linkedWorktree answers a second worktree of the repository of f at its tip, under the queue of
// home. Its wrapper is the committed wrapper of f.
func linkedWorktree(t *testing.T, home string, f *recordtest.Fixture) queueWorktree {
	t.Helper()
	root := filepath.Join(t.TempDir(), "linked")
	f.Git("worktree", "add", "-q", "--detach", root, "HEAD")
	commitmenttest.Register(t, root, "linked")
	return queueWorktree{root: root, home: home, launches: processtest.NewLaunches(t)}
}

// wrap writes body as the executable wrapper of the keyed worktree w.
func (w queueWorktree) wrap(t *testing.T, body string) {
	t.Helper()
	w.f.Write("bin/bench.sh", body)
	if err := os.Chmod(backgroundWrapper(w.root), 0o755); err != nil {
		t.Fatal(err)
	}
}

// pass runs the lane pass of the HEAD tip of w in this process, as the detached pass runs it. It
// answers w.
func (w queueWorktree) pass(t *testing.T) queueWorktree {
	t.Helper()
	w.passTip(t, gittest.Output(t, w.root, "rev-parse", "HEAD"))
	return w
}

// passTip runs the lane pass of tip in w, and the pass fails nothing.
func (w queueWorktree) passTip(t *testing.T, tip string) {
	t.Helper()
	if err := passLane(w.home, w.root, tip); err != nil {
		t.Fatalf("lane pass of %s in %s: %v", tip, w.root, err)
	}
}

// running makes the grade of w the running grade of its queue, and answers its stand-in.
func (w queueWorktree) running(t *testing.T) processtest.Launch {
	t.Helper()
	w.pass(t)
	return w.launches.Await(t)
}

// supersede commits a change to the tracked tree of w, so its new tip has a new grade key.
func (w queueWorktree) supersede(t *testing.T) {
	t.Helper()
	w.f.Write("superseded.txt", w.f.Tip()+"\n")
	w.f.Commit("supersede the tip")
}

// entryOf answers the queue entry of the HEAD tip of w, with the key of its stripped tree if any.
func entryOf(t *testing.T, w queueWorktree) backgroundstate.Entry {
	t.Helper()
	entry := backgroundstate.Entry{Worktree: w.root, Tip: gittest.Output(t, w.root, "rev-parse", "HEAD")}
	if w.f != nil {
		entry.Key = strippedEvidence(t, w.f).Grade
	}
	return entry
}

// requireQueue fails t unless the queue under home runs the grade of running, or none when
// running is empty, and queues the grades of queued in order.
func requireQueue(t *testing.T, home string, running backgroundstate.Entry, queued ...backgroundstate.Entry) {
	t.Helper()
	queue, err := backgroundstate.ReadQueue(home)
	var got backgroundstate.Entry
	if queue.Running != nil {
		got = queue.Running.Entry
	}
	if err != nil || got != running || !slices.Equal(queue.Queued, queued) {
		t.Fatalf("queue = running %+v, queued %+v (%v); want running %+v, queued %+v", got, queue.Queued, err, running, queued)
	}
}

// requireNoGrade fails t unless the queue under home holds no grade and w started none.
func requireNoGrade(t *testing.T, w queueWorktree) {
	t.Helper()
	requireQueue(t, w.home, backgroundstate.Entry{})
	w.launches.End(t)
	if got := w.launches.Launched(t); len(got) != 0 {
		t.Fatalf("launches = %+v; want none", got)
	}
}

// requireAlive fails t unless process pid runs exactly when want holds.
func requireAlive(t *testing.T, pid int, want bool) {
	t.Helper()
	if alive, err := subprocess.Alive(pid); err != nil || alive != want {
		t.Fatalf("grade %d alive = %t (%v); want %t", pid, alive, err, want)
	}
}

// awaitExit waits until process pid has exited.
func awaitExit(t *testing.T, pid int) {
	t.Helper()
	processtest.WaitFor(t, "grade "+strconv.Itoa(pid)+" to exit", func() bool { alive, err := subprocess.Alive(pid); return err == nil && !alive })
}

// TestQueueServesLanePassOrder is BG46. Three lane passes queue behind a running grade, and each
// later free slot starts the oldest of them.
func TestQueueServesLanePassOrder(t *testing.T) {
	home := t.TempDir()
	first := standInWorktree(t, home)
	first.running(t)
	a, b, c := standInWorktree(t, home), standInWorktree(t, home), standInWorktree(t, home)
	a.pass(t)
	b.pass(t)
	c.pass(t)
	requireQueue(t, home, entryOf(t, first), entryOf(t, a), entryOf(t, b), entryOf(t, c))
	first.launches.End(t)
	d := standInWorktree(t, home).pass(t)
	requireQueue(t, home, entryOf(t, a), entryOf(t, b), entryOf(t, c), entryOf(t, d))
	a.launches.Await(t)
	a.launches.End(t)
	e := standInWorktree(t, home).pass(t)
	requireQueue(t, home, entryOf(t, b), entryOf(t, c), entryOf(t, d), entryOf(t, e))
	b.launches.Await(t)
}

// TestFreshGreenKeyStartsNothing is BG47. A background green holds the key of the tip, so its
// lane pass queues nothing and starts no grade.
func TestFreshGreenKeyStartsNothing(t *testing.T) {
	w := keyedWorktree(t, t.TempDir())
	seedGateCache(t, w.f)
	requireBackgroundGreen(t, w.f, nil)
	w.pass(t)
	requireNoGrade(t, w)
}

// TestRunningKeyStartsNothing is BG48. A live claim holds the key of the tip, so its lane pass
// queues nothing and starts no grade.
func TestRunningKeyStartsNothing(t *testing.T) {
	w := keyedWorktree(t, t.TempDir())
	claim, err := backgroundstate.ClaimFor(entryOf(t, w).Key, w.root)
	if err != nil || backgroundstate.WriteClaim(durablefile.Native(), commonDir(t, w.root), claim) != nil {
		t.Fatalf("claim = %+v (%v); want it written", claim, err)
	}
	w.pass(t)
	requireNoGrade(t, w)
}

// TestQueuedKeyQueuesNothing is BG49. A record-only commit keeps the grade key of the queued tip,
// so its lane pass queues no second grade of that key.
func TestQueuedKeyQueuesNothing(t *testing.T) {
	home := t.TempDir()
	first := standInWorktree(t, home)
	first.running(t)
	w := keyedWorktree(t, home)
	w.pass(t)
	queued := entryOf(t, w)
	w.f.Record.Chunks[0].ID = "2"
	resave(t, w.f, "a record-only commit")
	w.pass(t)
	requireQueue(t, home, entryOf(t, first), queued)
}

// TestNewerTipReplacesTheQueuedGrade is BG50. A newer tip with a new key replaces the queued grade
// of its worktree. A late lane pass of the older tip is no longer the head of the worktree, so it
// drops out and does not replace the grade of the newer tip.
func TestNewerTipReplacesTheQueuedGrade(t *testing.T) {
	home := t.TempDir()
	first := standInWorktree(t, home)
	first.running(t)
	w := keyedWorktree(t, home)
	older := w.pass(t).f.Tip()
	w.supersede(t)
	w.pass(t)
	w.passTip(t, older)
	requireQueue(t, home, entryOf(t, first), entryOf(t, w))
}

// TestLanePassKeyFramesTheTargetExecutable is BG115. The launching verb runs an executable other
// than the one that the wrapper of the worktree resolves. The grade that its detached lane pass
// queues holds the key of the wrapper's executable, which is the key that the grade runs.
func TestLanePassKeyFramesTheTargetExecutable(t *testing.T) {
	w := keyedWorktree(t, t.TempDir())
	w.wrap(t, passWrapper(t, w.root, w.launches))
	w.f.Commit("pass through the wrapper")
	want := entryOf(t, w)
	setRunningExecutable(t, "the launching executable")
	if launcher := entryOf(t, w); launcher.Key == want.Key {
		t.Fatalf("key under the launching executable = %s; want a key other than the key of the wrapper's executable", launcher.Key)
	}
	var stderr bytes.Buffer
	if StartBackground(w.home, w.root, w.root, want.Tip, &stderr); stderr.Len() != 0 {
		t.Fatalf("launch printed %q", stderr.String())
	}
	w.launches.Await(t)
	awaitQueue(t, w.home, running)
	requireQueue(t, w.home, want)
}

// TestNewerTipStopsTheRunningGrade is BG51. A newer tip with a new key stops the running grade of
// its worktree, and the grade of the newer tip takes the slot.
func TestNewerTipStopsTheRunningGrade(t *testing.T) {
	w := keyedWorktree(t, t.TempDir())
	old := w.running(t)
	w.supersede(t)
	w.pass(t)
	requireQueue(t, w.home, entryOf(t, w))
	requireAlive(t, old.PID, false)
}

// TestWantedKeyKeepsRunning is BG52. A second worktree at the running tip has the running key as
// its latest key, so a newer tip of the first worktree queues and does not stop that grade.
func TestWantedKeyKeepsRunning(t *testing.T) {
	w := keyedWorktree(t, t.TempDir())
	wanted, grade := entryOf(t, w), w.running(t)
	linkedWorktree(t, w.home, w.f).pass(t)
	w.supersede(t)
	w.pass(t)
	requireQueue(t, w.home, wanted, entryOf(t, w))
	requireAlive(t, grade.PID, true)
}

// TestWantedKeyStaysQueued is BG53. A second worktree at the queued tip has the queued key as its
// latest key, so a newer tip of the first worktree queues beside that grade.
func TestWantedKeyStaysQueued(t *testing.T) {
	home := t.TempDir()
	first := standInWorktree(t, home)
	first.running(t)
	w := keyedWorktree(t, home)
	w.pass(t)
	wanted := entryOf(t, w)
	linkedWorktree(t, home, w.f).pass(t)
	w.supersede(t)
	w.pass(t)
	requireQueue(t, home, entryOf(t, first), wanted, entryOf(t, w))
}

// TestEndedGradeStartsTheOldestQueued is BG54. Two grades queue behind a real grade. When it ends,
// its own process starts the oldest of them and then exits, and that grade still runs.
func TestEndedGradeStartsTheOldestQueued(t *testing.T) {
	home := t.TempDir()
	h := holdingWorktree(t, home)
	h.pass(t)
	grade := h.awaitHolding(t)
	a, b := standInWorktree(t, home), standInWorktree(t, home)
	a.pass(t)
	b.pass(t)
	h.release(t)
	awaitExit(t, grade.Owner)
	requireQueue(t, home, entryOf(t, a), entryOf(t, b))
	requireAlive(t, a.launches.Await(t).PID, true)
}

// TestLanePassResumesADeadChain is BG55. The running grade ends with no chain, and the next lane
// pass starts the oldest queued grade, not its own.
func TestLanePassResumesADeadChain(t *testing.T) {
	home := t.TempDir()
	first := standInWorktree(t, home)
	first.running(t)
	a := standInWorktree(t, home)
	a.pass(t)
	first.launches.End(t)
	b := standInWorktree(t, home).pass(t)
	requireQueue(t, home, entryOf(t, a), entryOf(t, b))
	a.launches.Await(t)
}

// TestDeadOwnerRecordsAborted is BG56. The running grade dies with no record, and the next lane
// pass records aborted in the background record of its worktree.
func TestDeadOwnerRecordsAborted(t *testing.T) {
	home := t.TempDir()
	dead := standInWorktree(t, home)
	dead.running(t)
	dead.launches.End(t)
	standInWorktree(t, home).pass(t)
	gitdir := adminDir(t, dead.root)
	want := backgroundstate.Record{Outcome: backgroundstate.Aborted, Tip: entryOf(t, dead).Tip, Log: backgroundstate.LogPath(gitdir)}
	if record, err := backgroundstate.ReadRecord(gitdir); err != nil || record != want {
		t.Fatalf("record of the dead grade = %+v (%v); want %+v", record, err, want)
	}
}

// TestStoppedGradeRecordsStopped is BG83. A newer tip stops a real grade during its oracle. The
// grade writes no green, and its background record names stopped.
func TestStoppedGradeRecordsStopped(t *testing.T) {
	s := stopHoldingGrade(t)
	want := backgroundstate.Record{Outcome: backgroundstate.Stopped, Key: s.evidence.Grade, Tip: s.tip, Log: backgroundstate.LogPath(adminDir(t, s.root))}
	if record, err := backgroundstate.ReadRecord(adminDir(t, s.root)); err != nil || record != want {
		t.Fatalf("record of the stopped grade = %+v (%v); want %+v", record, err, want)
	}
	if evidence := InspectTree(s.root, s.evidence.Tree); evidence.Grade != s.evidence.Grade || evidence.ReusableGreen {
		t.Fatalf("evidence of the stopped subject = %+v; want no green under %s", evidence, s.evidence.Grade)
	}
}

// TestStoppedGradeReleasesItsResources is BG108. After a stopped grade exits, its private
// prospective bundle is gone, so no protection registration of that grade is left to resolve.
func TestStoppedGradeReleasesItsResources(t *testing.T) {
	s := stopHoldingGrade(t)
	requireAlive(t, s.grade.Owner, false)
	if retained, err := filepath.Glob(filepath.Join(s.tmp, prospectiveartifact.BundlePrefix+"*")); err != nil || len(retained) != 0 {
		t.Fatalf("bundles of the stopped grade = %v (%v); want none", retained, err)
	}
}

// stoppedGrade is a real grade that a newer tip stopped during its oracle: its slot, the tip that
// it graded, and the evidence of the subject of that tip before the grade.
type stoppedGrade struct {
	queueWorktree
	grade    backgroundstate.Slot
	tip      string
	evidence EvidenceInspection
}

// stopHoldingGrade starts the real grade of a holding worktree. While it holds, the wrapper becomes
// a stand-in, and the lane pass of a newer tip with a new key stops the grade and starts that one.
func stopHoldingGrade(t *testing.T) stoppedGrade {
	t.Helper()
	h := holdingWorktree(t, t.TempDir())
	s := stoppedGrade{queueWorktree: h, tip: h.f.Tip(), evidence: strippedEvidence(t, h.f)}
	h.pass(t)
	s.grade = h.awaitHolding(t)
	h.wrap(t, h.launches.StandIn())
	h.f.Commit("supersede the held tip")
	h.pass(t)
	requireQueue(t, h.home, entryOf(t, h))
	h.launches.Await(t)
	return s
}
