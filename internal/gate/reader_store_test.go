package gate

// The reader-table store, driven through whole runs of a fixture kit. The fixture gate
// routes through the phase table and copies canned logs into the test log directory of its
// run, so a run derives a table with no Go test run.

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
	"github.com/gibbonmi/bench/internal/gate/greenmarker"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/testrepo"
)

const (
	// readerFixtureLogs holds the canned logs that the reader gate copies into the test log
	// directory of each run.
	readerFixtureLogs outcomeWitness = ".gate-fixture-logs"
	// readerLogDir records the test log directory of the last run.
	readerLogDir outcomeWitness = ".gate-testlog"
	// readerLogRoot stands for the graded root in a canned log. The reader gate writes its own
	// working directory in its place, so a log fits a prospective checkout as well as the root.
	readerLogRoot = "@ROOT@"
)

// readerFixture is a fixture kit whose gate routes through the phase table and copies the
// canned logs, with its own root for readerLogRoot, into the test log directory of its run.
// The kit declares the root-conformance entry and tracks the link .claude/commands to
// ../.agents/commands.
func readerFixture(t *testing.T) string {
	t.Helper()
	t.Setenv(baselinePolicyEnv, "")
	root := outcomeFixture(t, func(f *testrepo.GateFixture, body string) string {
		return body + phaseTableHandoff + `if [ -n "${` + testLogEnv + `:-}" ]; then
  printf '%s\n' "$` + testLogEnv + `" > "$common/` + string(readerLogDir) + `"
  for log in "$common/` + string(readerFixtureLogs) + `"/*; do
    if [ -e "$log" ]; then ` + f.Command("sed") + ` "s|` + readerLogRoot + `|$(pwd)|g" "$log" > "$` + testLogEnv + `/${log##*/}"; fi
  done
fi
`
	})
	outcomeWrite(t, root, "go.mod", "module fixture\n\ngo 1.21\n", 0o644)
	outcomeWrite(t, root, conformancePackagePath+"/entry_test.go", goTest("conformance", registry.RootConformanceTest, ""), 0o644)
	if err := os.MkdirAll(filepath.Join(root, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("../.agents/commands", filepath.Join(root, ".claude", "commands")); err != nil {
		t.Fatal(err)
	}
	outcomeCommit(t, root, "kit")
	return root
}

// writeWitnessLog cans the root-conformance log of root, which reads .agents/ only through
// the tracked link.
func writeWitnessLog(t *testing.T, root string) {
	t.Helper()
	dir := readerFixtureLogs.path(t, root)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	file, err := testLogFile(dir, root, filepath.Join(root, conformancePackagePath))
	if err != nil {
		t.Fatal(err)
	}
	log := witnessLog(readerLogRoot, ".agents/") + "open " + readerLogRoot + "/.claude/commands/x.md\n"
	if err := os.WriteFile(file, []byte(log), 0o644); err != nil {
		t.Fatal(err)
	}
}

// runReaderGate runs the whole gate of root once under the fixture selection.
func runReaderGate(t *testing.T, root string) Result {
	t.Helper()
	return runReaderGateThrough(t, root, durablefile.Native())
}

// runReaderGateThrough runs the whole gate of root once under the fixture selection, and
// persists each record through replacer.
func runReaderGateThrough(t *testing.T, root string, replacer durablefile.Replacer) Result {
	t.Helper()
	var stdout, stderr bytes.Buffer
	owner := func(context.Context, string) (*runbinary.Selection, error) { return fixtureSelection(root), nil }
	result := executeSubjectReplacing(context.Background(), root, root, &stdout, &stderr, nil, forceRun, newGateEvaluation(root), owner, "", replacer)
	t.Logf("gate exit %d\nstdout:\n%s\nstderr:\n%s", result.ActionExit, stdout.String(), stderr.String())
	return result
}

// currentSubject answers the subject of root.
func currentSubject(t *testing.T, root string) subject {
	t.Helper()
	plan, err := buildSubject(root)
	if err != nil {
		t.Fatal(err)
	}
	return plan
}

// storedTables answers the trees of the tables in the store of root.
func storedTables(t *testing.T, root string) []string {
	t.Helper()
	common, err := benchgit.CommonDir(root)
	if err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(filepath.Join(common, benchgit.ReaderTableDir))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	var trees []string
	for _, entry := range entries {
		trees = append(trees, strings.TrimSuffix(entry.Name(), readerTableSuffix))
	}
	return trees
}

// TestWholeGreenPersistsTheReaderTable requires that a green whole run stores one table keyed
// by its tree, with the identity digest of the run and the readers that its logs show.
func TestWholeGreenPersistsTheReaderTable(t *testing.T) {
	root := readerFixture(t)
	writeWitnessLog(t, root)
	if result := runReaderGate(t, root); result.ActionExit != 0 {
		t.Fatalf("whole run = %d; want green", result.ActionExit)
	}
	plan := currentSubject(t, root)
	record, err := loadReaderTable(root, plan.Tree)
	if err != nil || record.Identity != plan.Identity || !slices.Contains(record.Readers[".agents/"], conformancePackagePath) {
		t.Fatalf("table of %s = %+v, %v; want identity %s and %s as a reader of .agents/", plan.Tree, record, err, plan.Identity, conformancePackagePath)
	}
	if tables := storedTables(t, root); !slices.Equal(tables, []string{plan.Tree}) {
		t.Fatalf("stored tables = %q; want the one table of %s", tables, plan.Tree)
	}
}

// TestWholeGreenLandingPersistsTheReaderTable is NG84. A landing grades a prospective
// checkout that the run releases before the verdict persists, so a table that derives after
// the release reads a removed root and stores nothing.
func TestWholeGreenLandingPersistsTheReaderTable(t *testing.T) {
	root := readerFixture(t)
	writeWitnessLog(t, root)
	tempRoot := t.TempDir()
	t.Setenv("TMPDIR", tempRoot)
	outcomeWrite(t, root, prospectiveGatePath, "#!/bin/sh\nexec \"$1/.bench/gate.sh\"\n", 0o755)
	outcomeWrite(t, root, "scripts/go-build.sh", "#!/bin/sh\n", 0o755)
	outcomeCommit(t, root, "prospective gate")
	tree := gittest.Output(t, root, "rev-parse", "HEAD^{tree}")
	// The prospective gate selects its executable from the checkout, so the kit is the checkout.
	owner := func(_ context.Context, source string) (*runbinary.Selection, error) {
		return fixtureSelection(source), nil
	}
	var stdout, stderr bytes.Buffer
	if result := executeTreeWithOwner(context.Background(), root, tree, &stdout, &stderr, nil, forceRun, owner); result.ActionExit != 0 {
		t.Fatalf("landing run = %+v; want green\nstdout:\n%s\nstderr:\n%s", result, stdout.String(), stderr.String())
	}
	requireNoProspectiveBundles(t, tempRoot)
	if tables := storedTables(t, root); !slices.Equal(tables, []string{tree}) {
		t.Fatalf("stored tables after a released landing run = %q; want the one table of %s", tables, tree)
	}
}

// TestOpenWholeGreenPersistsNoReaderTable requires that a green whole run of an open oracle
// subject, here for want of the input manifest, stores no table.
func TestOpenWholeGreenPersistsNoReaderTable(t *testing.T) {
	root := readerFixture(t)
	writeWitnessLog(t, root)
	gittest.MustRun(t, root, "rm", "-q", ".bench/gate-inputs.json")
	outcomeCommit(t, root, "open the subject")
	if plan := currentSubject(t, root); plan.Closed {
		t.Fatalf("subject without its manifest = %+v; want open", plan)
	}
	if result := runReaderGate(t, root); result.ActionExit != 0 {
		t.Fatalf("open whole run = %d; want green", result.ActionExit)
	}
	if tables := storedTables(t, root); len(tables) != 0 {
		t.Fatalf("stored tables after an open whole green = %q; want none", tables)
	}
}

// TestRedWholeRunPersistsNoReaderTable requires that a red whole run stores no table.
func TestRedWholeRunPersistsNoReaderTable(t *testing.T) {
	root := readerFixture(t)
	writeWitnessLog(t, root)
	outcomeWrite(t, root, ".gate-red", "\n", 0o644)
	if result := runReaderGate(t, root); result.ActionExit != 7 {
		t.Fatalf("whole run = %d; want the planted red", result.ActionExit)
	}
	if tables := storedTables(t, root); len(tables) != 0 {
		t.Fatalf("stored tables after a red = %q; want none", tables)
	}
}

// TestFailedVerdictPersistenceLeavesNoReaderTable requires that a green whole run whose final
// verdict fails to persist stores no table, because the table persists only after the
// verdict.
func TestFailedVerdictPersistenceLeavesNoReaderTable(t *testing.T) {
	root := readerFixture(t)
	writeWitnessLog(t, root)
	faults := newFaultPlan(map[durablefiletest.Fault]map[string]error{
		{Destination: benchgit.GateCacheFile, Ordinal: 2}: {"rename": errors.New("injected verdict rename failure")},
	})
	if result := runReaderGateThrough(t, root, faults.Replacer()); result.ActionExit != 1 {
		t.Fatalf("whole run = %d; want the failed verdict persistence", result.ActionExit)
	}
	if tables := storedTables(t, root); len(tables) != 0 {
		t.Fatalf("stored tables after a failed verdict persistence = %q; want none", tables)
	}
}

// TestReaderTableListsUnderTheClosedEnvironment requires that the package listing of a whole
// run reads the closed environment of the run, so an ambient GOFLAGS that adds a test package
// with no log refuses no table.
func TestReaderTableListsUnderTheClosedEnvironment(t *testing.T) {
	root := readerFixture(t)
	writeWitnessLog(t, root)
	outcomeWrite(t, root, "tagged/tagged.go", "package tagged\n", 0o644)
	outcomeWrite(t, root, "tagged/tagged_test.go", "//go:build ambient\n\npackage tagged\n", 0o644)
	outcomeCommit(t, root, "tagged")
	t.Setenv("GOFLAGS", "-tags=ambient")
	if result := runReaderGate(t, root); result.ActionExit != 0 {
		t.Fatalf("whole run = %d; want green", result.ActionExit)
	}
	if tables := storedTables(t, root); !slices.Equal(tables, []string{currentSubject(t, root).Tree}) {
		t.Fatalf("stored tables under an ambient GOFLAGS = %q; want the table of the run's tree", tables)
	}
}

// TestFailedDerivationKeepsTheGreen requires that a whole run whose derivation fails, here
// for want of any log, stays green and stores no table.
func TestFailedDerivationKeepsTheGreen(t *testing.T) {
	root := readerFixture(t)
	if result := runReaderGate(t, root); result.ActionExit != 0 || !Inspect(root).ReusableGreen {
		t.Fatalf("whole run = %d with the inspection %+v; want a reusable green", result.ActionExit, Inspect(root))
	}
	if tables := storedTables(t, root); len(tables) != 0 {
		t.Fatalf("stored tables after a failed derivation = %q; want none", tables)
	}
}

// TestReaderTableStoreBoundsItsTables requires that the store keeps at most readerTableLimit
// tables, removes the oldest first, and keeps the table of the green marker's tree.
func TestReaderTableStoreBoundsItsTables(t *testing.T) {
	root := gittest.RepoOnBranch(t, "main")
	outcomeWrite(t, root, "tracked.txt", "tracked\n", 0o644)
	outcomeCommit(t, root, "fixture")
	if err := greenmarker.Advance(root, "main", gittest.Output(t, root, "rev-parse", "HEAD"), ""); err != nil {
		t.Fatal(err)
	}
	marker := gittest.Output(t, root, "rev-parse", "HEAD^{tree}")
	store := func(tree string) string {
		t.Helper()
		path, err := storeReaderTable(durablefile.Native(), root, readerTableRecord{Schema: readerTableSchema, Tree: tree, Identity: strings.Repeat("a", 64), Readers: readerTable{}})
		if err != nil {
			t.Fatal(err)
		}
		return path
	}
	trees := []string{marker}
	for i := 1; i < readerTableLimit; i++ {
		trees = append(trees, fmt.Sprintf("%040x", i))
	}
	oldest := time.Now().Add(-time.Hour)
	for i, tree := range trees {
		at := oldest.Add(time.Duration(i) * time.Second)
		if err := os.Chtimes(store(tree), at, at); err != nil {
			t.Fatal(err)
		}
	}
	newest := fmt.Sprintf("%040x", readerTableLimit)
	store(newest)
	tables := storedTables(t, root)
	if len(tables) != readerTableLimit || !slices.Contains(tables, marker) || !slices.Contains(tables, newest) || slices.Contains(tables, trees[1]) {
		t.Fatalf("stored tables = %d %q; want %d, with the marker's %s and the newest %s, without the oldest other %s", len(tables), tables, readerTableLimit, marker, newest, trees[1])
	}
}

// TestWholeRunRemovesItsTestLogs requires that a whole run leaves no test log directory.
func TestWholeRunRemovesItsTestLogs(t *testing.T) {
	root := readerFixture(t)
	writeWitnessLog(t, root)
	if result := runReaderGate(t, root); result.ActionExit != 0 {
		t.Fatalf("whole run = %d; want green", result.ActionExit)
	}
	dir := strings.TrimSpace(string(outcomeRead(t, readerLogDir.path(t, root))))
	if _, err := os.Stat(dir); dir == "" || !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("test log directory %q after the run: %v; want it removed", dir, err)
	}
}
