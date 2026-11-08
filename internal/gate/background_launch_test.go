package gate

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/durablefile"
	benchfreshness "github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
)

// launchFixture is a kit repository whose own checkout is the worktree that a launch grades,
// under a private Bench home. Its committed wrapper is a stand-in grade of launches.
type launchFixture struct {
	root, home string
	launches   *processtest.Launches
}

// newLaunchFixture builds the fixture. An assigned fixture registers an active assignment for its
// checkout.
func newLaunchFixture(t *testing.T, assigned bool) launchFixture {
	t.Helper()
	f := launchFixture{root: gittest.RepoOnBranch(t, "main"), home: t.TempDir(), launches: processtest.NewLaunches(t)}
	if err := os.MkdirAll(filepath.Join(f.root, "bin"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(backgroundWrapper(f.root), []byte(f.launches.StandIn()), 0o755); err != nil {
		t.Fatal(err)
	}
	gittest.Output(t, f.root, "add", "bin/bench.sh")
	gittest.Output(t, f.root, "commit", "-qm", "wrapper")
	if assigned {
		commitmenttest.Register(t, f.root, "launch")
	}
	return f
}

// launch runs one launch of tip in the fixture and requires that it prints nothing.
func (f launchFixture) launch(t *testing.T, tip string) {
	t.Helper()
	var stderr bytes.Buffer
	StartBackground(f.home, f.root, f.root, tip, &stderr)
	if stderr.Len() != 0 {
		t.Fatalf("launch of %s printed %q", tip, stderr.String())
	}
}

// passing makes the wrapper of the fixture a pass wrapper, and answers the HEAD tip.
func (f launchFixture) passing(t *testing.T) string {
	t.Helper()
	if err := os.WriteFile(backgroundWrapper(f.root), []byte(passWrapper(t, f.root, f.launches)), 0o755); err != nil {
		t.Fatal(err)
	}
	return gittest.Output(t, f.root, "rev-parse", "HEAD")
}

// passWrapper answers a wrapper of the worktree root that runs the real entry of this test binary
// for a lane pass, and the stand-in grade of grades for each other run. So a test sees the grade
// that the detached lane pass starts.
func passWrapper(t *testing.T, root string, grades *processtest.Launches) string {
	t.Helper()
	return processtest.Branch(backgroundPassEnv, processtest.Script(t, backgroundGradeHelper, kitRunRootEnv+"="+root, backgroundResolvedEnv+"=1"), grades.StandIn())
}

// awaitQueue waits until the queue under home holds what done reports, so a detached lane pass has
// written it.
func awaitQueue(t *testing.T, home string, done func(backgroundstate.Queue) bool) {
	t.Helper()
	processtest.WaitFor(t, "the lane pass to write the queue", func() bool { queue, err := backgroundstate.ReadQueue(home); return err == nil && done(queue) })
}

// running reports whether the queue runs a grade.
func running(queue backgroundstate.Queue) bool { return queue.Running != nil }

// sessionOf answers the session id of process pid from its stat record. A host with no stat
// record skips the test.
func sessionOf(t *testing.T, pid int) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "stat"))
	if err != nil {
		capability.Capability(t, capability.PID, "no process stat record: "+err.Error())
	}
	// The fields after the parenthesized command are the state, the parent, the group, and
	// then the session.
	fields := strings.Fields(string(data[bytes.LastIndexByte(data, ')')+1:]))
	if len(fields) < 4 {
		t.Fatalf("stat record of %d = %q, want a session field", pid, data)
	}
	return fields[3]
}

// TestLaunchedGradeLeavesTheVerbSession is BG43. The grade that the detached lane pass starts runs
// in a session other than the session of the launching verb, so the end of the verb's session does
// not end the grade.
func TestLaunchedGradeLeavesTheVerbSession(t *testing.T) {
	f := newLaunchFixture(t, true)
	verb := sessionOf(t, os.Getpid())
	f.launch(t, f.passing(t))
	if grade := sessionOf(t, f.launches.Await(t).PID); grade == verb {
		t.Fatalf("grade session = %s, want a session other than the verb session %s", grade, verb)
	}
}

// TestSecondGradeStartsNoSecondOracle is BG45. While one grade runs, the lane pass of a second tip
// queues its grade and starts no second grade, so the machine runs one oracle.
func TestSecondGradeStartsNoSecondOracle(t *testing.T) {
	f := newLaunchFixture(t, true)
	first := f.passing(t)
	f.launch(t, first)
	f.launches.Await(t)
	awaitQueue(t, f.home, running)
	gittest.Output(t, f.root, "commit", "-q", "--allow-empty", "-m", "second")
	f.launch(t, gittest.Output(t, f.root, "rev-parse", "HEAD"))
	awaitQueue(t, f.home, func(queue backgroundstate.Queue) bool { return len(queue.Queued) == 1 })
	f.launches.End(t)
	want := []string{backgroundVerb, f.root, first}
	if got := f.launches.Launched(t); len(got) != 1 || !slices.Equal(got[0].Args, want) {
		t.Fatalf("launches = %+v, want the one grade %q", got, want)
	}
}

// TestStopAllRecordsStopped is BG83 through the stop of every grade. The stop ends a real grade
// during its oracle, and that grade records stopped itself.
func TestStopAllRecordsStopped(t *testing.T) {
	h := holdingWorktree(t, t.TempDir())
	h.pass(t)
	h.awaitHolding(t)
	if err := backgroundstate.StopAll(durablefile.Native(), h.home); err != nil {
		t.Fatal(err)
	}
	requireStoppedRecord(t, h)
}

// TestSlowTeardownRecordsStopped is BG83 under a slow teardown. The stop window ends before the
// real grade exits, so the stopper returns first, and the grade records stopped when it ends.
func TestSlowTeardownRecordsStopped(t *testing.T) {
	h := holdingWorktree(t, t.TempDir())
	h.pass(t)
	grade := h.awaitHolding(t)
	window, end := context.WithCancel(context.Background())
	end()
	if err := backgroundstate.Stop(window, durablefile.Native(), grade.Claim()); err != nil && !errors.Is(err, context.Canceled) {
		t.Fatalf("stop after the end of its window = %v; want no error or the end of the window", err)
	}
	awaitExit(t, grade.Owner)
	requireStoppedRecord(t, h)
}

// TestReapedStoppedGradeRecordsStopped is BG83 through the reap path. A stop marks the running
// grade and ends it before it writes any record, and the next lane pass reaps it as stopped.
func TestReapedStoppedGradeRecordsStopped(t *testing.T) {
	home := t.TempDir()
	dead := keyedWorktree(t, home)
	dead.running(t)
	queue, err := backgroundstate.ReadQueue(home)
	if err != nil || queue.Running == nil {
		t.Fatalf("running grade = %+v (%v); want the grade of %s", queue.Running, err, dead.root)
	}
	if err := backgroundstate.StopSlot(durablefile.Native(), *queue.Running); err != nil {
		t.Fatal(err)
	}
	standInWorktree(t, home).pass(t)
	gitdir, entry := adminDir(t, dead.root), entryOf(t, dead)
	want := backgroundstate.Record{Outcome: backgroundstate.Stopped, Key: entry.Key, Tip: entry.Tip, Log: backgroundstate.LogPath(gitdir)}
	if record, err := backgroundstate.ReadRecord(gitdir); err != nil || record != want {
		t.Fatalf("record of the reaped grade = %+v (%v); want %+v", record, err, want)
	}
}

// requireStoppedRecord fails t unless the background record of h names stopped for the key and
// the tip of its HEAD.
func requireStoppedRecord(t *testing.T, h queueWorktree) {
	t.Helper()
	gitdir := adminDir(t, h.root)
	want := backgroundstate.Record{Outcome: backgroundstate.Stopped, Key: strippedEvidence(t, h.f).Grade, Tip: h.f.Tip(), Log: backgroundstate.LogPath(gitdir)}
	if record, err := backgroundstate.ReadRecord(gitdir); err != nil || record != want {
		t.Fatalf("record of the stopped grade = %+v (%v); want %+v", record, err, want)
	}
}

// write writes body at name in the checkout of the fixture and stages it.
func (f launchFixture) write(t *testing.T, name, body string) {
	t.Helper()
	path := filepath.Join(f.root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	gittest.Output(t, f.root, "add", "--", name)
}

// tip commits the staged change of the fixture and answers the new tip.
func (f launchFixture) tip(t *testing.T) string {
	t.Helper()
	gittest.Output(t, f.root, "commit", "-qm", "change")
	return gittest.Output(t, f.root, "rev-parse", "HEAD")
}

// requirePasses waits until a lane pass of each tip of want has started, ends the stand-ins,
// and fails t unless the started passes are exactly those of want, in order.
func (f launchFixture) requirePasses(t *testing.T, want ...string) {
	t.Helper()
	tips := func() []string {
		var tips []string
		for _, launch := range f.launches.Launched(t) {
			tips = append(tips, launch.Args[len(launch.Args)-1])
		}
		return tips
	}
	processtest.WaitFor(t, "the wanted lane passes", func() bool {
		got := tips()
		return !slices.ContainsFunc(want, func(tip string) bool { return !slices.Contains(got, tip) })
	})
	f.launches.End(t)
	if got := tips(); !slices.Equal(got, want) {
		t.Fatalf("lane passes of %q, want %q", got, want)
	}
}

// requireNoPass fails t if the launch before it started a lane pass. The later launch of a tip
// that names no commit starts one pass, and the test waits for it, so a pass of the earlier
// launch, which started first, is in the record too.
func (f launchFixture) requireNoPass(t *testing.T) {
	t.Helper()
	f.launch(t, "tip")
	f.requirePasses(t, "tip")
}

// TestMarkdownDeletionStartsNoGrade is SP24. A tip that only deletes a Markdown file starts no
// lane pass, because the deleted path is Markdown too.
func TestMarkdownDeletionStartsNoGrade(t *testing.T) {
	f := newLaunchFixture(t, true)
	f.write(t, "notes.md", "# Notes\n")
	f.tip(t)
	gittest.Output(t, f.root, "rm", "-q", "notes.md")
	f.launch(t, f.tip(t))
	f.requireNoPass(t)
}

// TestRenameToMarkdownStartsAGrade is SP25. A tip that renames notes.txt to notes.md starts
// one lane pass, because the change list names the old text path as well as the new one.
func TestRenameToMarkdownStartsAGrade(t *testing.T) {
	f := newLaunchFixture(t, true)
	f.write(t, "notes.txt", "notes\n")
	f.tip(t)
	gittest.Output(t, f.root, "mv", "notes.txt", "notes.md")
	tip := f.tip(t)
	f.launch(t, tip)
	f.requirePasses(t, tip)
}

// TestQuotedMarkdownPathStartsNoGrade is SP26. A tip whose only change is a Markdown file whose
// name holds a double quote starts no lane pass, because the change list keeps each path byte.
func TestQuotedMarkdownPathStartsNoGrade(t *testing.T) {
	f := newLaunchFixture(t, true)
	f.write(t, `a "b".md`, "# Notes\n")
	f.launch(t, f.tip(t))
	f.requireNoPass(t)
}

// TestMarkdownBuildInputStartsAGrade is SP27. A tip whose only change is a Markdown file that
// the auxiliary build-input manifest names starts one lane pass, because it changes the build.
func TestMarkdownBuildInputStartsAGrade(t *testing.T) {
	f := newLaunchFixture(t, true)
	f.write(t, "go.mod", "module fixture\n\ngo 1.22\n")
	f.write(t, "cmd/bench/main.go", "package main\n\nfunc main() {}\n")
	f.write(t, benchfreshness.BuildInputsManifest, benchfreshness.BuildInputLine("notes", "notes.md"))
	f.write(t, "notes.md", "# Notes\n")
	f.tip(t)
	f.write(t, "notes.md", "# Notes\n\nA build input.\n")
	tip := f.tip(t)
	f.launch(t, tip)
	f.requirePasses(t, tip)
}

// TestRootCommitStartsAGrade is SP28. A root commit that holds only Markdown has no first
// parent to compare, so it starts one lane pass.
func TestRootCommitStartsAGrade(t *testing.T) {
	f := newLaunchFixture(t, true)
	gittest.Output(t, f.root, "checkout", "-q", "--orphan", "root")
	gittest.Output(t, f.root, "rm", "-rq", "--cached", ".")
	f.write(t, "notes.md", "# Notes\n")
	tip := f.tip(t)
	f.launch(t, tip)
	f.requirePasses(t, tip)
}

// TestUnreadableTipStartsAGrade is SP29. A tip that names no commit has no change list that
// reads, so it starts one lane pass.
func TestUnreadableTipStartsAGrade(t *testing.T) {
	f := newLaunchFixture(t, true)
	f.launch(t, "tip")
	f.requirePasses(t, "tip")
}

// TestUnassignedWorktreeStartsNoGrade holds the assignment predicate of the launch. A kit
// checkout that no active assignment owns starts no grade and takes no slot.
func TestUnassignedWorktreeStartsNoGrade(t *testing.T) {
	f := newLaunchFixture(t, false)
	f.launch(t, "tip")
	f.launches.End(t)
	queue, err := backgroundstate.ReadQueue(f.home)
	if err != nil || queue.Running != nil || len(f.launches.Launched(t)) != 0 {
		t.Fatalf("running grade = %+v (%v), launches = %+v; want no grade", queue.Running, err, f.launches.Launched(t))
	}
}
