package gate

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
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
