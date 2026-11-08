package commit

// The background grade that a commit starts. A commit only detaches the lane pass of its tip,
// which runs as the wrapper of the worktree. Each fixture worktree holds a wrapper that runs the
// stand-in of this test binary, so each lane pass that a commit starts is observable, and each one
// ends at the cleanup of its test. The gate package tests what a pass then queues.

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/poolkey"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
)

// standInGrade is the helper that the fixture wrapper runs.
const standInGrade = "stand-in-grade"

func TestMain(m *testing.M) {
	processtest.Main(m, map[string]func() int{standInGrade: runStandInGrade})
}

// runStandInGrade is the grade that the fixture wrapper runs as `bench.sh gate-background
// <worktree> <tip>`. It holds a claim on the tip as its grade key, records its launch, and
// releases the claim when the test releases it.
func runStandInGrade() int {
	args := os.Args[1:]
	if len(args) != 3 {
		return 2
	}
	common, err := git.CommonDir(args[1])
	if err != nil {
		return 1
	}
	claim, err := backgroundstate.ClaimFor(args[2], args[1])
	if err != nil || backgroundstate.WriteClaim(durablefile.Native(), common, claim) != nil {
		return 1
	}
	defer backgroundstate.ReleaseClaim(common, claim)
	if processtest.RecordLaunch() != nil {
		return 1
	}
	processtest.AwaitRelease()
	return 0
}

// launchFixture is the worktree of an assignment, at a path with a space in the pool of a
// private Bench home, whose wrapper runs the stand-in grade and whose lane runs check. The
// commit reads a kit apart from the fixture until asKit makes the primary checkout its kit.
type launchFixture struct {
	root, primary, common, home string
	launches                    *processtest.Launches
}

// newLaunchFixture builds the fixture.
func newLaunchFixture(t *testing.T, check string) launchFixture {
	t.Helper()
	return newLaunchFixtureWith(t, check, func(*testing.T, string) {})
}

// newLaunchFixtureWith builds the fixture with the files that tree writes in the worktree
// before the wrapper, so the wrapper of the fixture replaces any wrapper that tree writes.
func newLaunchFixtureWith(t *testing.T, check string, tree func(t *testing.T, root string)) launchFixture {
	t.Helper()
	notTheKitRoot(t)
	f := launchFixture{primary: t.TempDir(), home: t.TempDir()}
	t.Setenv("BENCH_HOME", f.home)
	initializeLandingRepo(t, f.primary, 0)
	linked := filepath.Join(poolkey.Pools(f.home), "linked work tree")
	gittest.MustRun(t, f.primary, "worktree", "add", "-q", "-b", "topic", linked)
	commitmenttest.Register(t, linked, landingLabel)
	commitmenttest.Admit(t, linked, landingLabel, "specs/commit-fixture/spec.md")
	// The record of the stand-ins comes last, so its cleanup ends them before any fixture
	// directory goes away.
	f.launches = processtest.NewLaunches(t)
	prepareLandingCheckout(t, linked, func(t *testing.T, root string) {
		tree(t, root)
		mustMkdirAll(t, filepath.Join(root, "bin"))
		mustWrite(t, filepath.Join(root, "bin", "bench.sh"), processtest.Script(t, standInGrade, f.launches.Env()...), 0o755)
		mustWrite(t, filepath.Join(root, ".bench", "check.sh"), check, 0o755)
		mustWrite(t, filepath.Join(root, ".bench", "phases.json"), `{"lane":[{"name":"check","argv":["sh",".bench/check.sh"]}]}`, 0o644)
	})
	f.root = gittest.Output(t, linked, "rev-parse", "--show-toplevel")
	var err error
	if f.common, err = git.CommonDir(f.root); err != nil {
		t.Fatal(err)
	}
	return f
}

// asKit makes the primary checkout the kit that the commit reads, as the wrapper of a kit
// worktree does.
func (f launchFixture) asKit(t *testing.T) {
	t.Helper()
	t.Setenv("BENCH_KIT", f.primary)
}

// kitLaunch is the launch fixture over a stale kit tree, marked as the kit, whose commits
// run build as the worktree build.
func kitLaunch(t *testing.T, check string, build *fakeBuild) launchFixture {
	t.Helper()
	f := newLaunchFixtureWith(t, check, staleKit)
	f.asKit(t)
	swapWorktreeBuild(t, build)
	return f
}

// swapWorktreeBuild makes build the worktree build of each commit until t ends.
func swapWorktreeBuild(t *testing.T, build *fakeBuild) {
	t.Helper()
	old := worktreeBuild
	t.Cleanup(func() { worktreeBuild = old })
	worktreeBuild = build.build
}

// seedRecord writes a red background record of the worktree, so a commit prints the latest
// background outcome before it starts its grade. It returns that record.
func (f launchFixture) seedRecord(t *testing.T) backgroundstate.Record {
	t.Helper()
	gitdir, err := git.AdminDir(f.root)
	if err != nil {
		t.Fatal(err)
	}
	record := backgroundstate.Record{Outcome: backgroundstate.Red, Key: "key", Tip: head(t, f.root), Log: backgroundstate.LogPath(gitdir)}
	if err := backgroundstate.WriteRecord(durablefile.Native(), gitdir, record); err != nil {
		t.Fatal(err)
	}
	return record
}

// requireNoGrade fails t unless the commit took no slot and no stand-in started.
func (f launchFixture) requireNoGrade(t *testing.T) {
	t.Helper()
	f.launches.End(t)
	queue, err := backgroundstate.ReadQueue(f.home)
	if err != nil || queue.Running != nil || len(f.launches.Launched(t)) != 0 {
		t.Fatalf("running grade = %+v (%v), launches = %+v; want no grade", queue.Running, err, f.launches.Launched(t))
	}
}

// TestCommitStartsABackgroundGrade is BG5. A commit that publishes in a kit worktree starts one
// lane pass, and the test waits for that pass. The pass names the worktree and the published tip.
func TestCommitStartsABackgroundGrade(t *testing.T) {
	f := newLaunchFixture(t, passingCheck)
	f.asKit(t)
	if code, stdout, stderr := runCommand(t, f.root, "-m", "m", "tracked.txt"); code != 0 {
		t.Fatalf("commit = (%d, %q, %q), want exit 0", code, stdout, stderr)
	}
	want := []string{f.root, head(t, f.root)}
	f.launches.Await(t)
	f.launches.End(t)
	if got := f.launches.Launched(t); len(got) != 1 || len(got[0].Args) != 3 || !slices.Equal(got[0].Args[1:], want) {
		t.Fatalf("launches = %+v, want one lane pass of %q", got, want)
	}
}

// markdownBody is the body of the Markdown file that a commit test writes. It is clean prose,
// so the lane passes it.
const markdownBody = "# Notes\n\nThis file holds notes.\n"

// requireChanges fails t unless the tip of root changes exactly paths against its first parent.
func requireChanges(t *testing.T, root string, paths ...string) {
	t.Helper()
	if got := strings.Fields(gittest.Output(t, root, "diff", "--name-only", "HEAD^1", "HEAD")); !slices.Equal(got, paths) {
		t.Fatalf("tip changes %q, want %q", got, paths)
	}
}

// TestMarkdownOnlyCommitStartsNoGrade is SP21. A kit-worktree commit whose only change is one
// Markdown file starts no lane pass.
func TestMarkdownOnlyCommitStartsNoGrade(t *testing.T) {
	f := newLaunchFixture(t, passingCheck)
	f.asKit(t)
	mustWrite(t, filepath.Join(f.root, "notes.md"), markdownBody, 0o644)
	if code, stdout, stderr := runCommand(t, f.root, "-m", "m", "notes.md"); code != 0 {
		t.Fatalf("commit = (%d, %q, %q), want exit 0", code, stdout, stderr)
	}
	requireChanges(t, f.root, "notes.md")
	// The next commit, of tracked.txt, starts one pass, so a pass of the Markdown commit, which
	// started first, is in the record when the test reads it.
	if code, stdout, stderr := runCommand(t, f.root, "-m", "next", "tracked.txt"); code != 0 {
		t.Fatalf("next commit = (%d, %q, %q), want exit 0", code, stdout, stderr)
	}
	f.requireOnlyHeadPass(t)
}

// requireOnlyHeadPass waits for the lane pass of the HEAD tip of the fixture, ends the
// stand-ins, and fails t unless that pass is the only one that started.
func (f launchFixture) requireOnlyHeadPass(t *testing.T) {
	t.Helper()
	want := []string{f.root, head(t, f.root)}
	passOf := func(launch processtest.Launch) bool {
		return len(launch.Args) == 3 && slices.Equal(launch.Args[1:], want)
	}
	processtest.WaitFor(t, "the lane pass of the HEAD tip", func() bool { return slices.ContainsFunc(f.launches.Launched(t), passOf) })
	f.launches.End(t)
	if got := f.launches.Launched(t); len(got) != 1 || !passOf(got[0]) {
		t.Fatalf("launches = %+v, want only the lane pass of %q", got, want)
	}
}

// TestMixedCommitStartsAGrade is SP22. A kit-worktree commit of one Markdown file and one text
// file starts one lane pass of its tip.
func TestMixedCommitStartsAGrade(t *testing.T) {
	f := newLaunchFixture(t, passingCheck)
	f.asKit(t)
	mustWrite(t, filepath.Join(f.root, "notes.md"), markdownBody, 0o644)
	if code, stdout, stderr := runCommand(t, f.root, "-m", "m", "notes.md", "tracked.txt"); code != 0 {
		t.Fatalf("commit = (%d, %q, %q), want exit 0", code, stdout, stderr)
	}
	requireChanges(t, f.root, "notes.md", "tracked.txt")
	f.requireOnlyHeadPass(t)
}

// TestLinkedRepoCommitStartsNoGrade is BG7. The kit of a linked repository is not the worktree,
// so a commit that publishes there starts no grade.
func TestLinkedRepoCommitStartsNoGrade(t *testing.T) {
	f := newLaunchFixture(t, passingCheck)
	if code, stdout, stderr := runCommand(t, f.root, "-m", "m", "tracked.txt"); code != 0 {
		t.Fatalf("commit = (%d, %q, %q), want exit 0", code, stdout, stderr)
	}
	f.requireNoGrade(t)
}

// TestDryRunStartsNoGrade is BG8. A dry run publishes nothing, so it starts no grade.
func TestDryRunStartsNoGrade(t *testing.T) {
	f := newLaunchFixture(t, passingCheck)
	f.asKit(t)
	if code, stdout, stderr := runCommand(t, f.root, "--dry-run", "-m", "m", "tracked.txt"); code != 0 {
		t.Fatalf("dry run = (%d, %q, %q), want exit 0", code, stdout, stderr)
	}
	f.requireNoGrade(t)
}

// TestRefusedCommitStartsNoGrade is BG9. A commit that the lane refuses publishes nothing, so it
// starts no grade.
func TestRefusedCommitStartsNoGrade(t *testing.T) {
	f := newLaunchFixture(t, failingCheck)
	f.asKit(t)
	if code, stdout, stderr := runCommand(t, f.root, "-m", "m", "tracked.txt"); code != 1 {
		t.Fatalf("commit = (%d, %q, %q), want the lane refusal at exit 1", code, stdout, stderr)
	}
	f.requireNoGrade(t)
}

// TestCommitPrintsTheLastBackgroundOutcome is BG85. The next commit in a worktree prints the
// outcome of the latest background grade of that worktree and the log path of that grade, and
// then it starts the next grade. A commit with no record prints no background outcome.
func TestCommitPrintsTheLastBackgroundOutcome(t *testing.T) {
	f := newLaunchFixture(t, passingCheck)
	f.asKit(t)
	record := f.seedRecord(t)
	code, stdout, stderr := runCommand(t, f.root, "-m", "m", "tracked.txt")
	if code != 0 {
		t.Fatalf("commit = (%d, %q, %q), want exit 0", code, stdout, stderr)
	}
	f.launches.Await(t)
	f.launches.End(t)
	if want := backgroundOutcomeLine(record); !strings.Contains(stderr, want) || !strings.Contains(want, string(record.Outcome)) || !strings.Contains(want, record.Log) {
		t.Fatalf("commit stderr = %q, want the line %q that names the red outcome and %s", stderr, want, record.Log)
	}
}

// TestCommitPrintsNoOutcomeWithoutARecord keeps a commit in a worktree with no background grade
// free of a background outcome line.
func TestCommitPrintsNoOutcomeWithoutARecord(t *testing.T) {
	f := newLaunchFixture(t, passingCheck)
	if code, stdout, stderr := runCommand(t, f.root, "-m", "m", "tracked.txt"); code != 0 || strings.Contains(stderr, backgroundOutcomePrefix) {
		t.Fatalf("commit = (%d, %q, %q), want exit 0 and no background outcome", code, stdout, stderr)
	}
	f.requireNoGrade(t)
}

// TestCommitReturnsBeforeItsGrade is BG11. The commit exits 0 while its lane pass still runs:
// after the commit returns, the stand-in pass holds a live claim. A commit that waited for its pass
// would return only after the pass released that claim.
func TestCommitReturnsBeforeItsGrade(t *testing.T) {
	f := newLaunchFixture(t, passingCheck)
	f.asKit(t)
	if code, stdout, stderr := runCommand(t, f.root, "-m", "m", "tracked.txt"); code != 0 {
		t.Fatalf("commit = (%d, %q, %q), want exit 0", code, stdout, stderr)
	}
	grade := f.launches.Await(t)
	claim, err := backgroundstate.ReadClaim(f.common, head(t, f.root))
	if err != nil {
		t.Fatal(err)
	}
	if live, err := claim.Live(); err != nil || !live || claim.Owner != grade.PID {
		t.Fatalf("claim after the commit = %+v (live %t, %v), want the live claim of grade %d", claim, live, err, grade.PID)
	}
}
