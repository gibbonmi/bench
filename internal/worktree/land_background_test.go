// The landing wait on a background grade. A landing reuses the background green of its key, it
// waits before its destination lease on a live claim of its key, and before its own run it clears
// the background grades of its worktree. Each fixture has its own Bench home, which holds the
// machine queue that the landing reads.
package worktree

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/landing/published"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
	"github.com/gibbonmi/bench/internal/testrepo"
)

// specLessUnbuiltFixture is the spec-less landing fixture whose source tip drops the commit of
// its review record, so the tip holds no open build and its background grade grades the landing.
func specLessUnbuiltFixture(t *testing.T, request string) landingFixture {
	t.Helper()
	f := specLessLandingFixture(t, request)
	gitRun(t, f.creation.Path, "reset", "-q", "--hard", "HEAD~1")
	f.tip = gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	return f
}

// gradeInTheBackground runs one real background grade of the fixture tip as the package's Bench
// executable, through a wrapper that the worktree holds only for the grade, and requires its green.
func gradeInTheBackground(t *testing.T, binary string, f landingFixture) {
	t.Helper()
	wrapper := filepath.Join(f.creation.Path, "bin", "bench.sh")
	mustMkdirAll(t, filepath.Dir(wrapper), 0o755)
	mustWrite(t, wrapper, []byte("#!/bin/sh\nexec "+sanitize.ShellQuote(binary)+" \"$@\"\n"), 0o755)
	var out bytes.Buffer
	cmd := descendant(t, binary, "gate-background", f.creation.Path, f.tip)
	cmd.Dir, cmd.Stdout, cmd.Stderr, cmd.Env = f.creation.Path, &out, &out, journeyChildEnv(t, f.home)
	code := exitCode(cmd.Run())
	mustNoError(t, os.RemoveAll(filepath.Dir(wrapper)))
	record, err := backgroundstate.ReadRecord(gitOutput(t, f.creation.Path, "rev-parse", "--absolute-git-dir"))
	if code != 0 || err != nil || record.Outcome != backgroundstate.Green {
		t.Fatalf("background grade = (%d, %q), record %+v (%v); want a green", code, out.String(), record, err)
	}
	requireOneRun(t, f.tally, "after the background grade")
}

// landAsTheBinary runs the landing of args as the package's Bench executable, so the landing and
// the background grade frame one running executable, and requires a publication.
func landAsTheBinary(t *testing.T, binary string, f landingFixture, args []string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	cmd := descendant(t, binary, append([]string{"worktree", "land"}, args...)...)
	cmd.Dir, cmd.Stdout, cmd.Stderr, cmd.Env = f.root, &stdout, &stderr, journeyChildEnv(t, f.home)
	if code := exitCode(cmd.Run()); code != 0 || gitOutput(t, f.root, "rev-parse", "main") == f.base {
		t.Fatalf("landing = (%d, %q, %q), want a publication", code, stdout.String(), stderr.String())
	}
}

// TestSpecLessLandingReusesTheBackgroundGreen is BG2. The landing of a spec-less tip after the
// background green of that tip runs no oracle.
func TestSpecLessLandingReusesTheBackgroundGreen(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	request := "land-background-spec-less"
	f := specLessUnbuiltFixture(t, request)
	gradeInTheBackground(t, binary, f)
	landAsTheBinary(t, binary, f, specLessLandArgs(request, f.base, f.tip, f.creation.Path))
	requireOneRun(t, f.tally, "after the landing of the graded tip")
}

// TestTicketsOnlyLandingReusesTheBackgroundGreen is BG3. The landing of a tip that adds one
// tickets-only folder after the background green of that tip runs no oracle.
func TestTicketsOnlyLandingReusesTheBackgroundGreen(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	request := "land-background-tickets-only"
	f := lightLandingFixture(t, request)
	gradeInTheBackground(t, binary, f)
	landAsTheBinary(t, binary, f, ticketsOnlyLandArgs(request, f.base, f.tip, lightLandingSlug, f.creation.Path))
	requireOneRun(t, f.tally, "after the landing of the graded tip")
}

// standInGrade is a process that stands in for the owner of a background grade: it leads its own
// session in the idle I/O class, as a grade does, until it is stopped. exited closes once the test
// has reaped it, so a stop that waits for its exit returns.
type standInGrade struct {
	pid    int
	exited chan struct{}
}

// startStandInGrade starts a stand-in grade. The test observes the I/O class through ionice, a
// Linux tool, so another platform skips.
func startStandInGrade(t *testing.T) standInGrade {
	t.Helper()
	ionice, err := exec.LookPath("ionice")
	if runtime.GOOS != "linux" || err != nil {
		capability.Environment(t, "the test reads the I/O class of a grade through ionice, a Linux tool")
	}
	cmd := descendant(t, ionice, "-c", "3", "sleep", "300")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	mustNoError(t, cmd.Start())
	g := standInGrade{pid: cmd.Process.Pid, exited: make(chan struct{})}
	go func() {
		_ = cmd.Wait()
		close(g.exited)
	}()
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		<-g.exited
	})
	processtest.WaitFor(t, "the stand-in grade to sleep in the idle class", func() bool { return g.ioClass(t) == "idle" })
	return g
}

// ioClass answers the I/O class that ionice reads for the stand-in grade.
func (g standInGrade) ioClass(t *testing.T) string {
	t.Helper()
	out, _ := descendant(t, "ionice", "-p", strconv.Itoa(g.pid)).Output()
	class, _, _ := strings.Cut(strings.TrimSpace(string(out)), ":")
	return class
}

// alive reports whether the stand-in grade still runs.
func (g standInGrade) alive() bool {
	select {
	case <-g.exited:
		return false
	default:
		return true
	}
}

// slot answers the slot of the stand-in grade for key in worktree at tip.
func (g standInGrade) slot(t *testing.T, key, worktree, tip string) backgroundstate.Slot {
	t.Helper()
	slot, err := backgroundstate.SlotFor(backgroundstate.Entry{Key: key, Worktree: worktree, Tip: tip}, g.pid)
	mustNoError(t, err)
	return slot
}

// claim writes the claim of the stand-in grade on key for the fixture worktree.
func (g standInGrade) claim(t *testing.T, f landingFixture, key string) {
	t.Helper()
	common := gitOutput(t, f.root, "rev-parse", "--path-format=absolute", "--git-common-dir")
	mustNoError(t, backgroundstate.WriteClaim(durablefile.Native(), common, g.slot(t, key, f.creation.Path, f.tip).Claim()))
}

// otherKey is a grade key that no landing of a fixture holds.
var otherKey = strings.Repeat("a", 64)

// lightLandingTree answers the tree that the tickets-only landing of the fixture publishes onto
// main, and the grade key of that tree.
func lightLandingTree(t *testing.T, f landingFixture) (tree, key string) {
	t.Helper()
	tree, err := published.Subject(f.root, "main", f.tip, "", lightPathTicketsFolder, nil)
	mustNoError(t, err)
	return tree, gate.InspectTree(f.root, tree).Grade
}

// lightPathTicketsFolder is the tickets-only folder that the light landing fixture delivers.
var lightPathTicketsFolder = landing.ClosedFolderPath(lightLandingSlug)

// landingEnd is the result of a landing that runs beside a test. done closes when it ends.
type landingEnd struct {
	r    verbResult
	done chan struct{}
}

// backgroundLightLanding is lightLandingOn of f after main commits an executable wrapper, so the
// source of request is a background worktree that holds its wrapper. The landing base is that
// commit.
func backgroundLightLanding(t *testing.T, f landingFixture, request string) landingFixture {
	t.Helper()
	wrapper := filepath.Join(f.root, "bin", "bench.sh")
	mustMkdirAll(t, filepath.Dir(wrapper), 0o755)
	mustWrite(t, wrapper, []byte("#!/bin/sh\nexit 0\n"), 0o755)
	gitRun(t, f.root, "add", "bin/bench.sh")
	gitRun(t, f.root, withIdentity("commit", "-qm", "wrapper")...)
	f.base = gitOutput(t, f.root, "rev-parse", "HEAD")
	return lightLandingOn(t, f, request)
}

// backgroundLand answers the tickets-only landing call of a backgroundLightLanding fixture, which
// reads the fixture repository as its kit.
func backgroundLand(t *testing.T, f landingFixture, request string) verbCall {
	t.Helper()
	call := f.call(ticketsOnlyLandArgs(request, f.base, f.tip, lightLandingSlug, f.creation.Path)...)
	call.kit = f.root
	return call
}

// landBeside starts the tickets-only landing of the fixture beside the test.
func landBeside(t *testing.T, f landingFixture, request string) *landingEnd {
	t.Helper()
	end := &landingEnd{done: make(chan struct{})}
	call := backgroundLand(t, f, request)
	go func() {
		end.r = runVerb(t, verbLand, call)
		close(end.done)
	}()
	t.Cleanup(func() { <-end.done })
	return end
}

// finished reports whether the landing has ended.
func (e *landingEnd) finished() bool {
	select {
	case <-e.done:
		return true
	default:
		return false
	}
}

// waitedLanding is a light landing fixture whose landing runs beside the test and waits on the
// claim of a stand-in grade on its own key.
type waitedLanding struct {
	landingFixture
	grade standInGrade
	tree  string
	end   *landingEnd
}

// waitOnStandIn starts the landing of a light fixture while a stand-in grade claims its key, and
// returns once the landing starts its wait: the I/O class of the grade leaves the idle class. A
// landing that ends first did not wait.
func waitOnStandIn(t *testing.T, request string) waitedLanding {
	t.Helper()
	grade := startStandInGrade(t)
	f := backgroundLightLanding(t, specLessLandingFixture(t, request+"-bound"), request)
	tree, key := lightLandingTree(t, f)
	grade.claim(t, f, key)
	w := waitedLanding{landingFixture: f, grade: grade, tree: tree, end: landBeside(t, f, request)}
	// A cleanup runs before each earlier one, so a failed test ends the grade before it waits for the landing.
	t.Cleanup(func() { _ = syscall.Kill(grade.pid, syscall.SIGKILL) })
	processtest.WaitFor(t, "the landing to start its wait", func() bool { return w.end.finished() || grade.ioClass(t) == "best-effort" })
	if w.end.finished() {
		t.Fatalf("landing ended before the grade = (%d, %q, %q), want a wait", w.end.r.exit, w.end.r.stdout, w.end.r.stderr)
	}
	return w
}

// endGrade stops the stand-in grade, so its claim is no longer live, and requires that the landing
// then publishes.
func (w waitedLanding) endGrade(t *testing.T) {
	t.Helper()
	mustNoError(t, syscall.Kill(w.grade.pid, syscall.SIGKILL))
	processtest.WaitFor(t, "the landing to end", w.end.finished)
	if w.end.r.exit != 0 || gitOutput(t, w.root, "rev-parse", "main") == w.base {
		t.Fatalf("landing after the wait = (%d, %q, %q), want a publication", w.end.r.exit, w.end.r.stdout, w.end.r.stderr)
	}
}

// TestLandingWaitsBeforeTheLease is BG75. While the landing waits on the live claim of its key,
// the destination lease is free, so another landing can take it.
func TestLandingWaitsBeforeTheLease(t *testing.T) {
	t.Parallel()
	w := waitOnStandIn(t, "land-background-wait")
	lease, err := commitrepo.Store{Root: w.root}.AcquireLease(leasedMain, otherLandingID)
	if err != nil {
		t.Fatalf("lease during the wait: %v; want it free", err)
	}
	lease.Release()
	if w.end.finished() {
		t.Fatalf("landing ended during the wait = (%d, %q, %q), want it to wait", w.end.r.exit, w.end.r.stdout, w.end.r.stderr)
	}
	w.endGrade(t)
	requireOneRun(t, w.tally, "after the landing with no green")
}

// TestLandingReusesTheGreenAfterTheWait is BG76. The waited grade writes its green, and after the
// wait the landing takes its lease and reuses that green with no oracle run of its own.
func TestLandingReusesTheGreenAfterTheWait(t *testing.T) {
	t.Parallel()
	w := waitOnStandIn(t, "land-background-reuse")
	var stdout, stderr bytes.Buffer
	if result := gate.ExecuteTree(context.Background(), w.root, w.tree, &stdout, &stderr); result.ActionExit != 0 {
		t.Fatalf("green of the waited key = %+v, stderr=%q", result, stderr.String())
	}
	requireOneRun(t, w.tally, "after the green of the waited grade")
	w.endGrade(t)
	requireOneRun(t, w.tally, "after the landing that waited")
}

// TestLandingOfAnotherKeyDoesNotWait is BG110. A live grade of the worktree claims another key, and
// the landing does not wait on it: it publishes while that grade still runs in the idle class.
func TestLandingOfAnotherKeyDoesNotWait(t *testing.T) {
	t.Parallel()
	grade := startStandInGrade(t)
	request := "land-background-other-key"
	f := backgroundLightLanding(t, specLessLandingFixture(t, request+"-bound"), request)
	grade.claim(t, f, otherKey)
	r := runVerb(t, verbLand, backgroundLand(t, f, request))
	if r.exit != 0 || !grade.alive() || grade.ioClass(t) != "idle" {
		t.Fatalf("landing beside a grade of another key = (%d, %q, %q), grade alive %t in class %q; want a publication beside the idle grade", r.exit, r.stdout, r.stderr, grade.alive(), grade.ioClass(t))
	}
}

// TestLandingOfAnotherKeyStopsTheGrade is BG77. The running grade of the worktree holds another
// key, and a queued grade of the worktree holds a third key. Before its own oracle runs, the
// landing stops the running grade and waits for its exit, and it drops the queued grade. The
// fixture gate goes red when the grade still runs.
func TestLandingOfAnotherKeyStopsTheGrade(t *testing.T) {
	t.Parallel()
	grade := startStandInGrade(t)
	request := "land-background-stop"
	step := func(*testrepo.GateFixture, string) string {
		return "if kill -0 " + strconv.Itoa(grade.pid) + " 2>/dev/null; then exit 9; fi\n"
	}
	f := backgroundLightLanding(t, landingFixtureWithGateStep(t, request+"-bound", "", "", filepath.Join(t.TempDir(), "bench-home"), false, step), request)
	running := grade.slot(t, otherKey, f.creation.Path, f.tip)
	queued := backgroundstate.Entry{Key: strings.Repeat("b", 64), Worktree: f.creation.Path, Tip: f.tip}
	mustNoError(t, backgroundstate.WriteQueue(durablefile.Native(), f.home, backgroundstate.Queue{Running: &running, Queued: []backgroundstate.Entry{queued}}))
	grade.claim(t, f, otherKey)
	r := runVerb(t, verbLand, backgroundLand(t, f, request))
	queue, err := backgroundstate.ReadQueue(f.home)
	mustNoError(t, err)
	if r.exit != 0 || grade.alive() || queue.Running != nil || len(queue.Queued) != 0 {
		t.Fatalf("landing beside grades of other keys = (%d, %q, %q), grade alive %t, queue %+v; want a publication after the stop with an empty queue", r.exit, r.stdout, r.stderr, grade.alive(), queue)
	}
	requireOneRun(t, f.tally, "after the landing that stopped the grade")
}

// TestLandingDropsItsQueuedKey is BG111. Another worktree runs a grade, and a grade of the
// landing's own key waits in the queue behind it. The landing's own run takes that grade out of
// the queue and leaves the running grade of the other worktree.
func TestLandingDropsItsQueuedKey(t *testing.T) {
	t.Parallel()
	grade := startStandInGrade(t)
	request := "land-background-queued-key"
	f := backgroundLightLanding(t, specLessLandingFixture(t, request+"-bound"), request)
	_, key := lightLandingTree(t, f)
	other := filepath.Join(t.TempDir(), "other")
	running := grade.slot(t, otherKey, other, f.tip)
	queued := backgroundstate.Entry{Key: key, Worktree: other, Tip: f.tip}
	mustNoError(t, backgroundstate.WriteQueue(durablefile.Native(), f.home, backgroundstate.Queue{Running: &running, Queued: []backgroundstate.Entry{queued}}))
	r := runVerb(t, verbLand, backgroundLand(t, f, request))
	queue, err := backgroundstate.ReadQueue(f.home)
	mustNoError(t, err)
	if r.exit != 0 || !grade.alive() || queue.Running == nil || *queue.Running != running || len(queue.Queued) != 0 {
		t.Fatalf("landing beside a queued grade of its key = (%d, %q, %q), queue %+v; want a publication that drops the queued grade and keeps %+v", r.exit, r.stdout, r.stderr, queue, running)
	}
	requireOneRun(t, f.tally, "after the landing's own run")
}
