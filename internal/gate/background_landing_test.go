package gate

// The background landing grade. Each fixture is a kit repository whose primary checkout holds
// main and a linked worktree whose tip holds no open build of a staged spec, so a background
// grade of that tip grades the tree that its landing would publish onto main.

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/durablefile"
	benchfreshness "github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
	"github.com/gibbonmi/bench/internal/gate/greenmarker"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/landing/published"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
)

// landingWorktree is a linked worktree of the repository whose primary checkout holds main.
type landingWorktree struct {
	t             *testing.T
	primary, root string
}

// newLandingWorktree seeds the gate cache of primary, which the fixture gate copies on each
// run, and links a worktree of primary.
func newLandingWorktree(t *testing.T, primary string) landingWorktree {
	t.Helper()
	if result := Execute(context.Background(), primary, io.Discard, io.Discard); result.ActionExit != 0 {
		t.Fatalf("seed gate = %+v; want a green", result)
	}
	return linkWorktree(t, primary)
}

// linkWorktree adds a linked worktree of primary on its own branch at main.
func linkWorktree(t *testing.T, primary string) landingWorktree {
	t.Helper()
	root := filepath.Join(t.TempDir(), "linked")
	gittest.MustRun(t, primary, "worktree", "add", "-q", "-b", "work", root, "main")
	return landingWorktree{t: t, primary: primary, root: root}
}

// commit writes each path of files with its body in the worktree, commits only those paths,
// and answers the tip.
func (w landingWorktree) commit(files map[string]string) string {
	w.t.Helper()
	for path, body := range files {
		outcomeWrite(w.t, w.root, path, body, 0o644)
		gittest.MustRun(w.t, w.root, "add", "--", path)
	}
	gittest.Commit(w.t, w.root, "work")
	return gittest.Output(w.t, w.root, "rev-parse", "HEAD")
}

// subject answers the tree that the landing of tip with deliverable publishes onto main.
func (w landingWorktree) subject(tip, deliverable string) string {
	w.t.Helper()
	tree, err := published.Subject(w.primary, "main", tip, "", deliverable, nil)
	if err != nil {
		w.t.Fatal(err)
	}
	return tree
}

// grade runs one background grade of tip in the worktree, and answers its exit, its output,
// and its record.
func (w landingWorktree) grade(tip string) (int, string, backgroundstate.Record) {
	w.t.Helper()
	var out bytes.Buffer
	code := gradeInBackground(context.Background(), w.root, tip, nil, &out, &out)
	record, err := backgroundstate.ReadRecord(adminDir(w.t, w.root))
	if err != nil {
		w.t.Fatal(err)
	}
	log, _ := os.ReadFile(backgroundstate.LogPath(adminDir(w.t, w.root)))
	return code, out.String() + string(log), record
}

// requireGreen requires that one background grade of tip exits 0 after one more oracle run
// and records green, and answers its record.
func (w landingWorktree) requireGreen(tip string) backgroundstate.Record {
	w.t.Helper()
	runs := outcomeRuns(w.t, w.primary)
	code, out, record := w.grade(tip)
	if code != 0 || record.Outcome != backgroundstate.Green || outcomeRuns(w.t, w.primary) != runs+1 {
		w.t.Fatalf("landing grade = %d %s, record %+v after %d runs; want a green after %d runs", code, out, record, outcomeRuns(w.t, w.primary), runs+1)
	}
	return record
}

// declareBroker makes the primary checkout declare Bench build inputs and binds a broker that
// is not the running executable. It answers the bound broker.
func declareBroker(t *testing.T, primary string) string {
	t.Helper()
	outcomeWrite(t, primary, benchfreshness.BuildInputsManifest, "cmd\n", 0o644)
	return bindBroker(t, primary)
}

// TestLandingGradeKeysThePublishedComposition is BG67. A tip that adds one tickets-only folder
// grades the tree that its landing publishes: the composition with the folder closed.
func TestLandingGradeKeysThePublishedComposition(t *testing.T) {
	w := newLandingWorktree(t, outcomeFixture(t))
	tip := w.commit(map[string]string{commitmenttest.TicketsFolder + "/tickets/one.md": "Light path ticket.\n", "work.txt": "work\n"})
	record := w.requireGreen(tip)
	closed, open := InspectTree(w.primary, w.subject(tip, commitmenttest.TicketsFolder)), InspectTree(w.primary, w.subject(tip, ""))
	if !closed.ReusableGreen || record.Key != closed.Grade || record.Key == open.Grade {
		t.Fatalf("record = %+v, published = %+v, unclosed = %+v; want the green of the published tree", record, closed, open)
	}
}

// TestSpecLessLandingGradeHasNoTransform is BG68. main moves after the worktree branches, and
// the grade of a spec-less tip keys the tip composed on main, with no transform.
func TestSpecLessLandingGradeHasNoTransform(t *testing.T) {
	w := newLandingWorktree(t, outcomeFixture(t))
	tip := w.commit(map[string]string{"work.txt": "work\n"})
	outcomeWrite(t, w.primary, "main.txt", "main\n", 0o644)
	outcomeCommit(t, w.primary, "move main")
	record := w.requireGreen(tip)
	composed, raw := InspectTree(w.primary, w.subject(tip, "")), InspectTree(w.primary, gittest.Output(t, w.root, "rev-parse", tip+"^{tree}"))
	if !composed.ReusableGreen || record.Key != composed.Grade || record.Key == raw.Grade {
		t.Fatalf("record = %+v, composed = %+v, raw tip = %+v; want the green of the composition", record, composed, raw)
	}
}

// TestLandingGradeRunsAsTheBroker is BG69. The entry of a landing tip execs the installed broker
// with its own arguments, and a grade that runs as that broker records its green under the key
// that the landing reads.
func TestLandingGradeRunsAsTheBroker(t *testing.T) {
	w := newLandingWorktree(t, outcomeFixture(t))
	tip := w.commit(map[string]string{"work.txt": "work\n"})
	broker := declareBroker(t, w.primary)
	calls := StubBrokerExec(t)
	t.Setenv(backgroundResolvedEnv, "1")
	runs := outcomeRuns(t, w.primary)
	var out bytes.Buffer
	code := BackgroundCommand([]string{w.root, tip}, &out, &out)
	want := []string{backgroundVerb, w.root, tip}
	if code != 0 || len(*calls) != 1 || (*calls)[0].Broker != broker || !slices.Equal((*calls)[0].Args, want) || outcomeRuns(t, w.primary) != runs {
		t.Fatalf("entry = %d %s with broker execs %+v after %d runs; want one exec of %s %q and no oracle", code, out.String(), *calls, outcomeRuns(t, w.primary), broker, want)
	}
	commitmenttest.BindRunning(t, w.primary)
	if record := w.requireGreen(tip); record.Key != InspectTree(w.primary, w.subject(tip, "")).Grade {
		t.Fatalf("record under the broker = %+v; want the key of the landing", record)
	}
}

// TestUnauthenticatedBrokerAbortsTheGrade is BG70. A broker whose bytes changed after its
// manifest bound them makes the grade record aborted, and no oracle runs.
func TestUnauthenticatedBrokerAbortsTheGrade(t *testing.T) {
	w := newLandingWorktree(t, outcomeFixture(t))
	tip := w.commit(map[string]string{"work.txt": "work\n"})
	broker := declareBroker(t, w.primary)
	if err := os.WriteFile(broker, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	runs := outcomeRuns(t, w.primary)
	code, out, record := w.grade(tip)
	want := backgroundstate.Record{Outcome: backgroundstate.Aborted, Tip: tip, Log: backgroundstate.LogPath(adminDir(t, w.root))}
	if record != want || outcomeRuns(t, w.primary) != runs {
		t.Fatalf("grade = %d %s, record = %+v after %d runs; want %+v and no oracle", code, out, record, outcomeRuns(t, w.primary), want)
	}
}

// TestLandingGradeTakesTheDestinationSchedule is BG71. The tip adds a phase manifest, and the
// grade key frames the schedule of main, not the schedule of the tip.
func TestLandingGradeTakesTheDestinationSchedule(t *testing.T) {
	w := newLandingWorktree(t, outcomeFixture(t))
	tip := w.commit(map[string]string{canary.PhaseManifestPath: `{"phases":[{"name":"tip","argv":["true"]}]}` + "\n"})
	record := w.requireGreen(tip)
	tree := w.subject(tip, "")
	destination, candidate := InspectTree(w.primary, tree), InspectTree(w.root, tree)
	if record.Key != destination.Grade || record.Key == candidate.Grade {
		t.Fatalf("record = %+v, destination schedule = %+v, tip schedule = %+v; want the key of the destination schedule", record, destination, candidate)
	}
}

// readerLanding is a kit whose prospective run builds its own run binary, whose main is its
// green marker with a stored reader table, and a linked worktree with a stand-in wrapper.
type readerLanding struct {
	landingWorktree
	home     string
	launches *processtest.Launches
}

func newReaderLanding(t *testing.T) readerLanding {
	t.Helper()
	primary := readerFixture(t)
	outcomeWrite(t, primary, prospectiveGatePath, "#!/bin/sh\nexit 0\n", 0o755)
	outcomeWrite(t, primary, "scripts/go-build.sh", "#!/bin/sh\n", 0o644)
	outcomeCommit(t, primary, "prospective kit")
	marker := gittest.Output(t, primary, "rev-parse", "HEAD")
	if err := greenmarker.Advance(primary, "main", marker, ""); err != nil {
		t.Fatal(err)
	}
	tree := gittest.Output(t, primary, "rev-parse", marker+"^{tree}")
	var identity string
	if err := withLandingPlan(primary, tree, func(_ string, plan subject) { identity = plan.Identity }); err != nil {
		t.Fatal(err)
	}
	if _, err := storeReaderTable(durablefile.Native(), primary, readerTableRecord{Schema: readerTableSchema, Tree: tree, Identity: identity, Readers: planReaders}); err != nil {
		t.Fatal(err)
	}
	// The pass runs no oracle, so the fixture seeds no gate cache.
	f := readerLanding{landingWorktree: linkWorktree(t, primary), home: t.TempDir(), launches: processtest.NewLaunches(t)}
	outcomeWrite(t, f.root, "bin/bench.sh", f.launches.StandIn(), 0o755)
	return f
}

// pass runs the lane pass of tip in the worktree.
func (f readerLanding) pass(tip string) {
	f.t.Helper()
	if err := passLane(f.home, f.root, tip); err != nil {
		f.t.Fatalf("lane pass of %s: %v", tip, err)
	}
}

// TestReaderSetLandingStartsNoGrade is BG72. A covered change whose landing plan selects the
// reader set queues no grade and starts none.
func TestReaderSetLandingStartsNoGrade(t *testing.T) {
	f := newReaderLanding(t)
	f.pass(f.commit(map[string]string{"docs/a.md": "change\n"}))
	requireQueue(t, f.home, backgroundstate.Entry{})
	f.launches.End(t)
	if got := f.launches.Launched(t); len(got) != 0 {
		t.Fatalf("launches = %+v; want none", got)
	}
}

// TestWholeGateLandingStartsAGrade is BG73. An uncovered change sends the landing plan to the
// whole gate, so the pass queues one landing grade under the key of the landing.
func TestWholeGateLandingStartsAGrade(t *testing.T) {
	f := newReaderLanding(t)
	tip := f.commit(map[string]string{"uncovered.txt": "change\n"})
	f.pass(tip)
	f.launches.Await(t)
	requireQueue(t, f.home, backgroundstate.Entry{Key: InspectTree(f.primary, f.subject(tip, "")).Grade, Worktree: f.root, Tip: tip})
	f.pass(tip)
	f.launches.End(t)
	if got := f.launches.Launched(t); len(got) != 1 || !slices.Equal(got[0].Args, []string{backgroundVerb, f.root, tip}) {
		t.Fatalf("launches = %+v; want the one landing grade of %s", got, tip)
	}
}

// TestConflictedLandingStartsNoGrade is a lane pass of a tip whose landing composition conflicts
// on main. The landing refuses before its gate, so the pass queues no grade and starts none.
func TestConflictedLandingStartsNoGrade(t *testing.T) {
	w := linkWorktree(t, outcomeFixture(t))
	launches, home := processtest.NewLaunches(t), t.TempDir()
	outcomeWrite(t, w.root, "bin/bench.sh", launches.StandIn(), 0o755)
	tip := w.commit(map[string]string{"work.txt": "work\n"})
	outcomeWrite(t, w.primary, "work.txt", "main\n", 0o644)
	outcomeCommit(t, w.primary, "conflict on main")
	if err := passLane(home, w.root, tip); err != nil {
		t.Fatalf("lane pass of the conflicted tip: %v", err)
	}
	requireQueue(t, home, backgroundstate.Entry{})
	launches.End(t)
	if got := launches.Launched(t); len(got) != 0 {
		t.Fatalf("launches = %+v; want none", got)
	}
}

// TestNonKitLandingDerivesNoSubject is story 6 for the landing wait. The worktree holds its
// wrapper and an active assignment, and a live claim holds the key of its landing, so a landing
// that derives its subject waits on that claim. The landing of a kit worktree waits, and the
// landing of a worktree that is not its kit derives no subject: it carries no wait and no clear.
func TestNonKitLandingDerivesNoSubject(t *testing.T) {
	w := linkWorktree(t, outcomeFixture(t))
	tip := w.commit(map[string]string{"work.txt": "work\n"})
	outcomeWrite(t, w.root, "bin/bench.sh", "#!/bin/sh\nexit 0\n", 0o755)
	commitmenttest.Register(t, w.root, "non-kit-landing")
	claim, err := backgroundstate.ClaimFor(InspectTree(w.primary, w.subject(tip, "")).Grade, w.root)
	if err != nil {
		t.Fatal(err)
	}
	if err := backgroundstate.WriteClaim(durablefile.Native(), gittest.Output(t, w.root, "rev-parse", "--path-format=absolute", "--git-common-dir"), claim); err != nil {
		t.Fatal(err)
	}
	// The wait ends at once, so a landing that derives its subject answers the budget of a wait.
	ended, cancel := context.WithCancel(context.Background())
	cancel()
	for kit, derives := range map[string]bool{w.root: true, gittest.RepoOnBranch(t, "main"): false} {
		ctx, release := AwaitLanding(ended, t.TempDir(), kit, w.root, tip, io.Discard)
		wait, carried := landingWaitOf(ctx)
		release()
		if waited := wait.budget.waited != nil; carried != derives || waited != derives {
			t.Fatalf("landing with kit %s carries a wait %t that waited %t; want %t for both", kit, carried, waited, derives)
		}
	}
}

// TestWrongDeliverableOnlyMisses is BG74. A tip that adds two tickets-only folders derives no
// deliverable, so its grade writes one green, only under the spec-less key that it ran.
func TestWrongDeliverableOnlyMisses(t *testing.T) {
	w := newLandingWorktree(t, outcomeFixture(t))
	tip := w.commit(map[string]string{"specs/a/tickets/one.md": "a\n", "specs/b/tickets/one.md": "b\n"})
	before := storeEntries(t, w.primary)
	record := w.requireGreen(tip)
	added := slices.DeleteFunc(storeEntries(t, w.primary), func(name string) bool { return slices.Contains(before, name) })
	ran, closed := InspectTree(w.primary, w.subject(tip, "")), InspectTree(w.primary, w.subject(tip, "specs/a"))
	if !slices.Equal(added, []string{record.Key}) || record.Key != ran.Grade || closed.ReusableGreen {
		t.Fatalf("new store entries = %v, record = %+v, closed = %+v; want the one key %s that the grade ran", added, record, closed, ran.Grade)
	}
}
