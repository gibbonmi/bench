package gate

// The landing plan. Each test builds a fixture kit with a green marker commit, a stored reader
// table for the marker's tree, and an optional chain of reader-record commits. The test then
// commits one change after the destination and asks the plan for its verdict. Each fixture
// change reaches exactly one case of the plan.

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/durablefile"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// planReaders is the stored table of every landing fixture. The internal/worktree reader of
// roadmap/ is a race-table package, and CHANGELOG.md has no reader beyond the witness, so
// only the child-reader declaration adds internal/releaseevidence.
var planReaders = readerTable{
	"docs/":        {"internal/conformance", "internal/docsreader"},
	"roadmap/":     {"internal/conformance", "internal/worktree"},
	"CHANGELOG.md": {"internal/conformance"},
}

// landingFixture is a fixture kit whose green marker commit has a stored reader table with
// the identity of the kit's subject. destination starts at the marker.
type landingFixture struct {
	t                         *testing.T
	root, marker, destination string
	plan                      subject
}

// newLandingFixture builds the kit, applies base, and commits it as the green marker.
func newLandingFixture(t *testing.T, base func(*landingFixture)) *landingFixture {
	t.Helper()
	f := &landingFixture{t: t, root: readerFixture(t)}
	f.plan = currentSubject(t, f.root)
	if base != nil {
		base(f)
	}
	f.marker = f.commit("marker")
	f.destination = f.marker
	f.storeTable(f.plan.Identity)
	return f
}

// storeTable stores planReaders with identity as the table of the marker's tree.
func (f *landingFixture) storeTable(identity string) {
	f.t.Helper()
	if _, err := storeReaderTable(durablefile.Native(), f.root, readerTableRecord{Schema: readerTableSchema, Tree: gittest.Output(f.t, f.root, "rev-parse", f.marker+"^{tree}"), Identity: identity, Readers: planReaders}); err != nil {
		f.t.Fatal(err)
	}
}

// put writes body at rel and stages it.
func (f *landingFixture) put(rel, body string) {
	f.t.Helper()
	outcomeWrite(f.t, f.root, rel, body, 0o644)
	gittest.MustRun(f.t, f.root, "add", "--", rel)
}

// link replaces whatever rel holds with a symbolic link to target, and stages it.
func (f *landingFixture) link(rel, target string) {
	f.t.Helper()
	full := filepath.Join(f.root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		f.t.Fatal(err)
	}
	if err := os.RemoveAll(full); err != nil {
		f.t.Fatal(err)
	}
	if err := os.Symlink(target, full); err != nil {
		f.t.Fatal(err)
	}
	gittest.MustRun(f.t, f.root, "add", "--", rel)
}

// gitlink stages a gitlink at rel that names the HEAD commit.
func (f *landingFixture) gitlink(rel string) {
	f.t.Helper()
	gittest.MustRun(f.t, f.root, "update-index", "--add", "--cacheinfo", benchgit.GitlinkMode+","+gittest.Output(f.t, f.root, "rev-parse", "HEAD")+","+rel)
}

// commit commits the index and answers the new commit.
func (f *landingFixture) commit(message string) string {
	f.t.Helper()
	gittest.Commit(f.t, f.root, message, "--allow-empty")
	return gittest.Output(f.t, f.root, "rev-parse", "HEAD")
}

// chain adds n first-parent commits after the destination and moves the destination to the
// last. Each commit has a reader record of its tree, except the commit at index gap.
func (f *landingFixture) chain(n, gap int) {
	f.t.Helper()
	gitdir, err := benchgit.CommonDir(f.root)
	if err != nil {
		f.t.Fatal(err)
	}
	for i := range n {
		f.put(fmt.Sprintf("roadmap/FT%d.md", i), "chain\n")
		f.destination = f.commit(fmt.Sprintf("reader landing %d", i))
		if i == gap {
			continue
		}
		tree := gittest.Output(f.t, f.root, "rev-parse", f.destination+"^{tree}")
		run := newReaderRun(subject{Tree: tree}, readerSchedule{Packages: []string{"./internal/conformance"}}, tree, f.marker)
		if err := writeReaderGreen(durablefile.Native(), gitdir, run, strings.Repeat("d", 64), time.Now()); err != nil {
			f.t.Fatal(err)
		}
	}
}

// decide commits change after the destination and answers the plan's verdict over the
// committed tree, under ctx, with the inventory as the shellcheck inventory and a selection
// whose kit is kit. It answers the reader run, whether the run is the reader set, and stdout.
func (f *landingFixture) decide(ctx context.Context, change func(*landingFixture), inventory []string, kit string) (readerRun, bool, string) {
	f.t.Helper()
	change(f)
	plan := f.plan
	plan.Tree = gittest.Output(f.t, f.root, "rev-parse", f.commit("change")+"^{tree}")
	var stdout bytes.Buffer
	planner := readerPlanner{shellcheck: func(string) []string { return inventory }}
	run, eligible := planner.plan(ctx, &stdout, f.root, f.root, plan, &runbinary.Selection{Path: filepath.Join(f.root, "bench"), SourceRoot: kit})
	return run, eligible, stdout.String()
}

// landing answers ctx with the fixture's landing value.
func (f *landingFixture) landing() context.Context {
	return WithLanding(context.Background(), f.destination, f.marker)
}

// refuses requires that change, after the fixture's chain, sends the landing to the whole
// gate and names refusal.
func (f *landingFixture) refuses(change func(*landingFixture), inventory []string, refusal string) {
	f.t.Helper()
	if _, eligible, out := f.decide(f.landing(), change, inventory, f.root); eligible || out != fmt.Sprintf(readerRefusedFormat, refusal) {
		f.t.Fatalf("plan = %t %q; want the whole gate with the refusal %q", eligible, out, refusal)
	}
}

// selects requires that change, after the fixture's chain, runs the reader set, and answers
// its schedule.
func (f *landingFixture) selects(change func(*landingFixture)) readerSchedule {
	f.t.Helper()
	run, eligible, out := f.decide(f.landing(), change, nil, f.root)
	if !eligible || out != "" {
		f.t.Fatalf("plan = %t %q; want the reader set", eligible, out)
	}
	schedule, err := parseReaderSchedule(run.schedule)
	if err != nil {
		f.t.Fatal(err)
	}
	return schedule
}

// docsChange changes one covered path.
func docsChange(f *landingFixture) { f.put("docs/a file.md", "change\n") }

// TestLandingRunsTheReaderSet proves that a landing over a covered change, with a valid
// table, an equal identity, and a clean chain, runs the reader set and records its tree.
func TestLandingRunsTheReaderSet(t *testing.T) {
	f := newLandingFixture(t, nil)
	// The fixture gate copies the verdict cache, so a whole green comes first.
	if result := Execute(context.Background(), f.root, &bytes.Buffer{}, &bytes.Buffer{}); result.ActionExit != 0 {
		t.Fatalf("whole green = %+v", result)
	}
	f.chain(1, -1)
	docsChange(f)
	f.commit("change")
	stub := t.TempDir()
	outcomeWrite(t, stub, "bench", "#!/bin/sh\nexit 0\n", 0o755)
	owner := func(context.Context, string) (*runbinary.Selection, error) {
		return &runbinary.Selection{Path: filepath.Join(stub, "bench"), SourceRoot: f.root}, nil
	}
	var stdout, stderr bytes.Buffer
	result := executeSubjectReplacing(f.landing(), f.root, f.root, &stdout, &stderr, nil, forceRun, newGateEvaluation(f.root), owner, "", durablefile.Native())
	tree := gittest.Output(t, f.root, "rev-parse", "HEAD^{tree}")
	if recorded, err := ReaderRecorded(f.root, tree); result.ActionExit != 0 || result.Schedule != ScheduleReader || !recorded || err != nil {
		t.Fatalf("landing = %+v, recorded %t %v\nstdout:\n%s\nstderr:\n%s\nwant a green reader set that records %s", result, recorded, err, stdout.String(), stderr.String(), tree)
	}
}

// TestReaderPlanSelectsTheChangedFamilyReaders proves that the reader packages are the
// table's readers of each changed family, the declared child readers, and the
// root-conformance package, each once, with no race flag.
func TestReaderPlanSelectsTheChangedFamilyReaders(t *testing.T) {
	schedule := newLandingFixture(t, nil).selects(func(f *landingFixture) {
		docsChange(f)
		f.put("docs/b.md", "change\n")
		f.put("CHANGELOG.md", "change\n")
	})
	want := []string{"./internal/conformance", "./internal/docsreader", "./internal/releaseevidence"}
	if !slices.Equal(schedule.Packages, want) || schedule.Race {
		t.Fatalf("reader schedule = %+v; want the packages %q without the race flag", schedule, want)
	}
}

// TestReaderPlanFlagsTheRaceTable proves that a race-table reader of a changed family sets
// the race flag.
func TestReaderPlanFlagsTheRaceTable(t *testing.T) {
	schedule := newLandingFixture(t, nil).selects(func(f *landingFixture) { f.put("roadmap/FT9.md", "change\n") })
	if !schedule.Race || !slices.Contains(schedule.Packages, "./internal/worktree") {
		t.Fatalf("reader schedule = %+v; want the race flag for ./internal/worktree", schedule)
	}
}

// TestReaderPlanRunsWholeForAnUncoveredPath proves that a change to AGENTS.md runs the whole
// gate.
func TestReaderPlanRunsWholeForAnUncoveredPath(t *testing.T) {
	newLandingFixture(t, nil).refuses(func(f *landingFixture) { f.put("AGENTS.md", "change\n") }, nil, refusedUncovered)
}

// TestReaderPlanRunsWholeForARenameOutOfAFamily proves that a rename of docs/x.md to x.md
// runs the whole gate.
func TestReaderPlanRunsWholeForARenameOutOfAFamily(t *testing.T) {
	f := newLandingFixture(t, func(f *landingFixture) { f.put("docs/x.md", "moved\n") })
	f.refuses(func(f *landingFixture) { gittest.MustRun(t, f.root, "mv", "docs/x.md", "x.md") }, nil, refusedUncovered)
}

// embedFixture is a landing fixture whose base holds the Go file name with body and the
// covered file target.
func embedFixture(t *testing.T, name, body, target string) *landingFixture {
	return newLandingFixture(t, func(f *landingFixture) {
		f.put(name, "package fixture\n\nimport _ \"embed\"\n\n"+body+"\nvar embedded string\n")
		f.put(target, "embedded\n")
	})
}

// TestReaderPlanRunsWholeForAnEmbedTarget proves that a change to an exact embed target runs
// the whole gate.
func TestReaderPlanRunsWholeForAnEmbedTarget(t *testing.T) {
	embedFixture(t, "embed.go", "//go:embed docs/x.txt", "docs/x.txt").refuses(func(f *landingFixture) { f.put("docs/x.txt", "change\n") }, nil, refusedBuildInput)
}

// TestReaderPlanRunsWholeBelowAnEmbeddedDirectory proves that deleting a file below an
// embedded directory runs the whole gate.
func TestReaderPlanRunsWholeBelowAnEmbeddedDirectory(t *testing.T) {
	embedFixture(t, "embed.go", "//go:embed docs/assets", "docs/assets/a.txt").refuses(func(f *landingFixture) { gittest.MustRun(t, f.root, "rm", "-q", "docs/assets/a.txt") }, nil, refusedBuildInput)
}

// TestReaderPlanRunsWholeForAnEmbedGlob proves that a change to a glob match runs the whole
// gate.
func TestReaderPlanRunsWholeForAnEmbedGlob(t *testing.T) {
	embedFixture(t, "embed.go", "//go:embed docs/*.txt", "docs/a.txt").refuses(func(f *landingFixture) { f.put("docs/a.txt", "change\n") }, nil, refusedBuildInput)
}

// TestReaderPlanRunsWholeForATestFileEmbed proves that a change to the target of a test
// file's embed directive runs the whole gate.
func TestReaderPlanRunsWholeForATestFileEmbed(t *testing.T) {
	embedFixture(t, "embed_test.go", "//go:embed docs/t.txt", "docs/t.txt").refuses(func(f *landingFixture) { f.put("docs/t.txt", "change\n") }, nil, refusedBuildInput)
}

// TestReaderPlanRunsWholeForAShellcheckPath proves that a change to a path of an injected
// shellcheck inventory runs the whole gate.
func TestReaderPlanRunsWholeForAShellcheckPath(t *testing.T) {
	newLandingFixture(t, nil).refuses(func(f *landingFixture) { f.put("docs/run.sh", "#!/bin/sh\n") }, []string{"docs/run.sh"}, refusedBuildInput)
}

// TestReaderPlanRunsWholeForAnAddedSymlink proves that adding a link in docs/ runs the whole
// gate.
func TestReaderPlanRunsWholeForAnAddedSymlink(t *testing.T) {
	newLandingFixture(t, nil).refuses(func(f *landingFixture) { f.link("docs/link", "a.md") }, nil, refusedLink)
}

// TestReaderPlanRunsWholeForARemovedSymlink proves that replacing a link in docs/ with a
// regular file runs the whole gate.
func TestReaderPlanRunsWholeForARemovedSymlink(t *testing.T) {
	f := newLandingFixture(t, func(f *landingFixture) { f.link("docs/link", "a.md") })
	f.refuses(func(f *landingFixture) {
		if err := os.Remove(filepath.Join(f.root, "docs", "link")); err != nil {
			t.Fatal(err)
		}
		f.put("docs/link", "regular\n")
	}, nil, refusedLink)
}

// TestReaderPlanRunsWholeForAnAddedGitlink proves that adding a gitlink in docs/ runs the
// whole gate, even when the repository's configuration hides submodule changes from a diff.
func TestReaderPlanRunsWholeForAnAddedGitlink(t *testing.T) {
	f := newLandingFixture(t, nil)
	gittest.MustRun(t, f.root, "config", "diff.ignoreSubmodules", "all")
	f.refuses(func(f *landingFixture) { f.gitlink("docs/module") }, nil, refusedLink)
}

// TestReaderPlanRunsWholeForARemovedGitlink proves that removing a gitlink from docs/ runs the
// whole gate.
func TestReaderPlanRunsWholeForARemovedGitlink(t *testing.T) {
	f := newLandingFixture(t, func(f *landingFixture) { f.gitlink("docs/module") })
	f.refuses(func(f *landingFixture) { gittest.MustRun(t, f.root, "rm", "-q", "--cached", "docs/module") }, nil, refusedLink)
}

// TestReaderPlanRunsWholeWithoutATable proves that a store without the table of the marker's
// tree runs the whole gate.
func TestReaderPlanRunsWholeWithoutATable(t *testing.T) {
	f := newLandingFixture(t, nil)
	gitdir, err := benchgit.CommonDir(f.root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(filepath.Join(gitdir, benchgit.ReaderTableDir)); err != nil {
		t.Fatal(err)
	}
	f.refuses(docsChange, nil, refusedNoTable)
}

// TestReaderPlanRunsWholeForANewGradeIdentity proves that a table of another grade identity
// runs the whole gate under a case that names the grade identity.
func TestReaderPlanRunsWholeForANewGradeIdentity(t *testing.T) {
	f := newLandingFixture(t, nil)
	f.storeTable(strings.Repeat("a", 64))
	f.refuses(docsChange, nil, "new grade identity")
}

// TestReaderPlanRunsWholeForAnOpenSubject proves that an open oracle subject runs the whole
// gate under the identity case, even with a table of its own identity.
func TestReaderPlanRunsWholeForAnOpenSubject(t *testing.T) {
	f := newLandingFixture(t, nil)
	f.plan.Closed = false
	f.refuses(docsChange, nil, refusedIdentity)
}

// TestReaderPlanRunsWholeAfterFiveReaderLandings proves that five reader-record commits after
// the marker run the whole gate.
func TestReaderPlanRunsWholeAfterFiveReaderLandings(t *testing.T) {
	f := newLandingFixture(t, nil)
	f.chain(5, -1)
	f.refuses(docsChange, nil, refusedChainBound)
}

// TestReaderPlanRunsTheReaderSetAfterFourReaderLandings proves that four reader-record
// commits after the marker run the reader set.
func TestReaderPlanRunsTheReaderSetAfterFourReaderLandings(t *testing.T) {
	f := newLandingFixture(t, nil)
	f.chain(4, -1)
	f.selects(docsChange)
}

// TestReaderPlanRunsWholeForAChainGap proves that a commit after the marker with no reader
// record runs the whole gate, at the first commit and at the last before the destination.
func TestReaderPlanRunsWholeForAChainGap(t *testing.T) {
	for _, gap := range []int{0, 1} {
		f := newLandingFixture(t, nil)
		f.chain(2, gap)
		f.refuses(docsChange, nil, refusedChainGap)
	}
}

// TestReaderPlanRunsWholeWithoutAMarker proves that a landing value with no green marker, or
// with a marker that is not an ancestor of the destination, runs the whole gate.
func TestReaderPlanRunsWholeWithoutAMarker(t *testing.T) {
	f := newLandingFixture(t, nil)
	if _, eligible, out := f.decide(WithLanding(context.Background(), f.destination, ""), docsChange, nil, f.root); eligible || out != fmt.Sprintf(readerRefusedFormat, refusedNoMarker) {
		t.Fatalf("plan without a marker = %t %q; want the whole gate with the refusal %q", eligible, out, refusedNoMarker)
	}
	f.destination = gittest.Output(t, f.root, "commit-tree", "-m", "unrelated destination", f.marker+"^{tree}")
	f.refuses(docsChange, nil, refusedNoMarker)
}

// TestReaderPlanNeedsALandingContext proves that the covered change without a landing value
// runs the whole gate and prints no refusal, because only a landing reads the plan.
func TestReaderPlanNeedsALandingContext(t *testing.T) {
	f := newLandingFixture(t, nil)
	if _, eligible, out := f.decide(context.Background(), docsChange, nil, f.root); eligible || out != "" {
		t.Fatalf("plan without a landing value = %t %q; want the whole gate with no line", eligible, out)
	}
}

// TestReaderPlanRunsWholeForALinkedRoot proves that a covered landing of a root that is not
// the kit of its selection runs the whole gate and prints no line, so a linked repository
// keeps its landing output.
func TestReaderPlanRunsWholeForALinkedRoot(t *testing.T) {
	f := newLandingFixture(t, nil)
	if _, eligible, out := f.decide(f.landing(), docsChange, nil, t.TempDir()); eligible || out != "" {
		t.Fatalf("plan of a linked root = %t %q; want the whole gate with no line", eligible, out)
	}
}

// TestReaderPlanNamesTheFirstRefusingCase proves that a landing that leaves the covered
// families, adds a link, and meets a table of another identity names the uncovered path, the
// first of the three cases, in the exact line.
func TestReaderPlanNamesTheFirstRefusingCase(t *testing.T) {
	f := newLandingFixture(t, nil)
	f.storeTable(strings.Repeat("a", 64))
	_, eligible, out := f.decide(f.landing(), func(f *landingFixture) {
		f.put("AGENTS.md", "change\n")
		f.link("docs/link", "a.md")
	}, nil, f.root)
	if want := "gate: whole gate (reader set refused: uncovered path)\n"; eligible || out != want {
		t.Fatalf("plan with three refusing cases = %t %q; want %q", eligible, out, want)
	}
}
