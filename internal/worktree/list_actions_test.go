package worktree

import (
	"bytes"
	cryptorand "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

func TestListCommandRendersTypedAdminRefusal(t *testing.T) {
	t.Parallel()
	root := journeyRepoOnBranch(t, "main")
	journeyFIFOWorktreeAdmin(t, root, "typed")
	r := runVerb(t, verbList, repoHome{root, Home()}.call())
	if r.exit == 0 || !strings.Contains(r.stdout, "worktrees/typed/gitdir") || !strings.Contains(r.stdout, "fifo") || !strings.Contains(r.stdout, "inspect and remove it") {
		t.Fatalf("typed list output code=%d out=%q", r.exit, r.stdout)
	}
}

func TestListCommandKeepsTypedAndPorcelainFailureActionsDistinct(t *testing.T) {
	for _, tc := range []struct {
		mode, detail, action string
	}{
		{"bad-rev-parse", "missing-common", "investigate the git failure"},
		{"fail-worktree", "cannot read registered worktrees", "run git worktree list and retry"},
	} {
		t.Run(tc.mode, func(t *testing.T) {
			root := journeyRepoOnBranch(t, "main")
			journeyStubGit(t, root, tc.mode, filepath.Join(t.TempDir(), "argv"))
			r := runVerb(t, verbList, repoHome{root, Home()}.call())
			if r.exit != 1 || !strings.Contains(r.stdout, tc.detail) || !strings.Contains(r.stdout, tc.action) {
				t.Fatalf("%s list output code=%d out=%q", tc.mode, r.exit, r.stdout)
			}
		})
	}
}

func TestListCommandRendersBoundExpiryAsTypedFailure(t *testing.T) {
	restore := git.SetWorktreeListTimeoutForTest(100 * time.Millisecond)
	t.Cleanup(restore)
	root := journeyRepoOnBranch(t, "main")
	journeyStubGit(t, root, "block-worktree", filepath.Join(t.TempDir(), "argv"))
	r := runVerb(t, verbList, repoHome{root, Home()}.call())
	if r.exit != 1 || !strings.Contains(r.stdout, "worktree list") || !strings.Contains(r.stdout, "investigate the git failure") || strings.Contains(r.stdout, "inspect and remove it") || strings.Contains(r.stdout, "retry") {
		t.Fatalf("bound list output code=%d out=%q", r.exit, r.stdout)
	}
}

type worktreeListResponse struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Exit   int    `json:"exit"`
}

type worktreeListArgvPair struct {
	Argv []string             `json:"argv"`
	Old  worktreeListResponse `json:"old"`
	New  worktreeListResponse `json:"new"`
}

type worktreeListTerminalPair struct {
	Old worktreeListResponse `json:"old"`
	New worktreeListResponse `json:"new"`
}

func setRandomReader(t *testing.T, data []byte) {
	t.Helper()
	randomReader := cryptorand.Reader
	cryptorand.Reader = bytes.NewReader(data)
	t.Cleanup(func() { cryptorand.Reader = randomReader })
}

func TestListCommandCheckedInOldNewArgvCompatibility(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile("testdata/pre-disclosure-argv-pairs.json")
	if err != nil {
		t.Fatal(err)
	}
	var pairs []worktreeListArgvPair
	if err := json.Unmarshal(data, &pairs); err != nil {
		t.Fatal(err)
	}
	root := journeyRepo(t)
	for _, pair := range pairs {
		r := runVerb(t, verbList, repoHome{root, Home()}.call(pair.Argv...))
		if r.stdout != pair.New.Stdout || pair.New.Stderr != "" || r.exit != pair.New.Exit {
			t.Fatalf("list %q = stdout=%q stderr=%q exit=%d, want checked-in new response", pair.Argv, r.stdout, "", r.exit)
		}
		if pair.Old != pair.New {
			if len(pair.Argv) != 1 || (pair.Argv[0] != "--help" && pair.Argv[0] != "-h" && pair.Argv[0] != "help") {
				t.Fatalf("paired fixture admits an unapproved argv delta: %#v", pair)
			}
		}
	}
}

func TestListCommandHelpAndArgumentMatrix(t *testing.T) {
	t.Parallel()
	root := journeyRepo(t)
	for _, arg := range []string{"--help", "-h", "help"} {
		r := runVerb(t, verbList, repoHome{root, Home()}.call(arg))
		if r.exit != 0 || r.stdout != "usage: bench worktree list\n" {
			t.Errorf("list %q = (%d, %q)", arg, r.exit, r.stdout)
		}
	}
	for _, args := range [][]string{{"--unknown"}, {"extra"}, {"--"}} {
		r := runVerb(t, verbList, repoHome{root, Home()}.call(args...))
		want := "usage: bench worktree list (unknown argument: " + args[0] + ")\n"
		if r.exit != 2 || r.stdout != want {
			t.Errorf("list %q = (%d, %q), want usage exit 2", args, r.exit, r.stdout)
		}
	}
}

func TestListCommandPreservesCheckedInEmptyPrimaryResponse(t *testing.T) {
	t.Parallel()
	primary, err := os.ReadFile("testdata/pre-disclosure-empty.stdout")
	if err != nil {
		t.Fatal(err)
	}
	root := journeyRepo(t)
	r := runVerb(t, verbList, repoHome{root, Home()}.call())
	if r.exit != 0 || r.stdout != string(primary)+"help[0]{cmd,why}:\n" {
		t.Fatalf("ListCommand = (%d, %q), want checked-in primary plus exactly one help block", r.exit, r.stdout)
	}
}

func TestListCommandCheckedInPresentForeignTerminalPair(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile("testdata/pre-disclosure-present-foreign-pair.json")
	if err != nil {
		t.Fatal(err)
	}
	var pair worktreeListTerminalPair
	if err := json.Unmarshal(data, &pair); err != nil {
		t.Fatal(err)
	}
	root := newWorktreeRepo(t)
	present := filepath.Join(t.TempDir(), "present foreign")
	gitRun(t, root, "worktree", "add", "-q", "--detach", present, "HEAD")
	r := runVerb(t, verbList, repoHome{root, Home()}.call())
	pair.Old.Stdout = strings.ReplaceAll(pair.Old.Stdout, "{{PRESENT}}", present)
	pair.New.Stdout = strings.ReplaceAll(pair.New.Stdout, "{{PRESENT}}", present)
	if pair.Old.Stderr != "" || pair.New.Stderr != "" || pair.Old.Exit != 0 || pair.New.Exit != 0 || pair.New.Stdout != pair.Old.Stdout+"help[0]{cmd,why}:\n" {
		t.Fatalf("terminal pair admits a response change beyond exactly one empty help block: %#v", pair)
	}
	if r.exit != pair.New.Exit || r.stdout != pair.New.Stdout {
		t.Fatalf("ListCommand = stdout=%q stderr=%q exit=%d, want checked-in terminal primary plus exactly one empty help block", r.stdout, "", r.exit)
	}
}

// TestListCommandCheckedInCompletedAssignmentTerminalPair pins the owned completed
// assignment row, the state the present-foreign pair above never reaches. The release
// transaction compacts a completed record on its way out. A listing can then observe
// only the state a release interrupted at its terminal-receipt boundary leaves. This
// fixture reaches that state through ReleaseCommand rather than a hand-written ledger
// entry. A completed assignment is non-actionable, so its disclosure is exactly one
// empty help block.
func TestListCommandCheckedInCompletedAssignmentTerminalPair(t *testing.T) {
	data, err := os.ReadFile("testdata/pre-disclosure-complete-assignment-pair.json")
	if err != nil {
		t.Fatal(err)
	}
	var pair worktreeListTerminalPair
	if err := json.Unmarshal(data, &pair); err != nil {
		t.Fatal(err)
	}
	root := newWorktreeRepo(t)
	bindEnv(t, "BENCH_HOME", filepath.Join(t.TempDir(), "bench-home"))
	setRandomReader(t, []byte(strings.Repeat("\x10", 16)+strings.Repeat("\x01", 16)))
	creation := mustCreate(t, root, Home(), "landed-complete-assignment", "complete assignment")
	j := defaultJoins()
	j.cleanupBoundary = func(step LifecycleStep) error {
		if step == StepTerminalReceipt {
			return errors.New("stop before the completed record is compacted")
		}
		return nil
	}
	if released := runVerb(t, verbRelease, repoHome{root, Home()}.callWith(j, "--request", "landed-complete-assignment", creation.Path)); released.exit == 0 {
		t.Fatalf("interrupted release exit = %d, want non-zero", released.exit)
	}
	assignments, err := intent.Assignments(root)
	if err != nil || len(assignments) != 1 || assignments[0].State != intent.StateComplete {
		t.Fatalf("Assignments = %#v, %v, want one completed assignment", assignments, err)
	}
	r := runVerb(t, verbList, repoHome{root, Home()}.call())
	materialize := strings.NewReplacer("{{LABEL}}", assignments[0].Label)
	pair.Old.Stdout = materialize.Replace(pair.Old.Stdout)
	pair.New.Stdout = materialize.Replace(pair.New.Stdout)
	if pair.Old.Stderr != "" || pair.New.Stderr != "" || pair.Old.Exit != 0 || pair.New.Exit != 0 || pair.New.Stdout != pair.Old.Stdout+"help[0]{cmd,why}:\n" {
		t.Fatalf("terminal pair admits a response change beyond exactly one empty help block: %#v", pair)
	}
	if r.exit != pair.New.Exit || r.stdout != pair.New.Stdout {
		t.Fatalf("ListCommand = stdout=%q stderr=%q exit=%d, want checked-in completed-assignment primary plus exactly one empty help block", r.stdout, "", r.exit)
	}
}

func TestActionsForRowsEnumeratesActiveAndOrphanRows(t *testing.T) {
	t.Parallel()
	rows := [][]any{
		{"a", "active", "", "active"},
		{"done", "complete", "", "complete"},
		{"foreign", "one", "", "foreign", "foreign", "missing"},
		{"b", "active", "", "active"},
		{"foreign", "present", "", "foreign", "foreign", "present"},
		{"foreign", "two", "", "foreign", "foreign", "missing"},
	}
	owned := make([]listRow, len(rows))
	for i, row := range rows {
		owned[i] = listRow{values: row}
	}
	owned[2].orphanPath, owned[5].orphanPath = "/tmp/orphan one", "/tmp/orphan-two"
	help, err := axi.RenderHelp(actionsForRows(owned))
	if err != nil {
		t.Fatal(err)
	}
	want := "help[4]{cmd,why}:\n" + activeHelpRows + "  bench worktree clean '/tmp/orphan one',clean the orphaned worktree\n  bench worktree clean /tmp/orphan-two,clean the orphaned worktree\n"
	if help != want {
		t.Fatalf("help = %q, want %q", help, want)
	}
}

func TestListCommandPublicRowsAndDisclosure(t *testing.T) {
	primaryTemplate, err := os.ReadFile("testdata/pre-disclosure-active-orphan.stdout")
	if err != nil {
		t.Fatal(err)
	}
	root := newWorktreeRepo(t)
	bindEnv(t, "BENCH_HOME", filepath.Join(t.TempDir(), "bench-home"))
	setRandomReader(t, []byte(strings.Repeat("\x10", 16)+strings.Repeat("\x01", 16)+strings.Repeat("\x20", 16)+strings.Repeat("\x02", 16)))
	mustCreate(t, root, Home(), "request-a", "alpha")
	mustCreate(t, root, Home(), "request-b", "beta")
	present := filepath.Join(t.TempDir(), "present foreign")
	missing := filepath.Join(t.TempDir(), "missing foreign * path")
	gitRun(t, root, "worktree", "add", "-q", "--detach", present, "HEAD")
	gitRun(t, root, "worktree", "add", "-q", "--detach", missing, "HEAD")
	if err := os.RemoveAll(missing); err != nil {
		t.Fatal(err)
	}
	r := runVerb(t, verbList, repoHome{root, Home()}.call())
	assignments, err := intent.Assignments(root)
	if err != nil || len(assignments) != 2 {
		t.Fatalf("Assignments = %#v, %v, want two producer-ordered rows", assignments, err)
	}
	primary := strings.NewReplacer(
		"{{LABEL1}}", assignments[0].Label,
		"{{LABEL2}}", assignments[1].Label,
		"{{PRESENT}}", present,
		"{{MISSING}}", missing,
	).Replace(string(primaryTemplate))
	help := "help[4]{cmd,why}:\n" + activeHelpRows + fmt.Sprintf("  bench worktree clean '%s',clean the orphaned worktree\n  bench worktree clean --landed,clean landed assignments\n", missing)
	if r.exit != 0 || r.stdout != primary+help {
		t.Fatalf("ListCommand = (%d, %q), want materialized checked-in primary plus exactly one help block", r.exit, r.stdout)
	}
}

func TestListCommandControlBearingOrphanPathPreservesPrimaryAndAction(t *testing.T) {
	t.Parallel()
	for _, control := range []string{"\t", "\n", "\r"} {
		t.Run("control", func(t *testing.T) {
			root := newWorktreeRepo(t)
			missing := filepath.Join(t.TempDir(), "orphan"+control+"path")
			gitRun(t, root, "worktree", "add", "-q", "--detach", missing, "HEAD")
			if err := os.RemoveAll(missing); err != nil {
				t.Fatal(err)
			}
			r := runVerb(t, verbList, repoHome{root, Home()}.call())
			if r.exit != 0 || !strings.HasPrefix(r.stdout, listTable+"[1]{id,label,request,state,source,tree,lease,landed,ignored}:\n") {
				t.Fatalf("ListCommand = (%d, %q), want primary worktree response and exit 0", r.exit, r.stdout)
			}
			if !strings.Contains(r.stdout, "help[1]{cmd,why}:\n") {
				t.Fatalf("ListCommand = %q, want one orphan-clean action", r.stdout)
			}
			argv, err := axitest.RecoverHelpCommandArgv(r.stdout)
			if err != nil {
				t.Fatal(err)
			}
			want := []string{"bench", "worktree", "clean", missing}
			if !slices.Equal(argv, want) {
				t.Fatalf("shell argv = %q, want %q", argv, want)
			}
		})
	}
}

func TestListCommandAngleBracketOrphanPathPreservesPrimaryAndHonestFallback(t *testing.T) {
	t.Parallel()
	for _, marker := range []string{"<", ">"} {
		t.Run(marker, func(t *testing.T) {
			root := newWorktreeRepo(t)
			missing := filepath.Join(t.TempDir(), "orphan"+marker+"path")
			gitRun(t, root, "worktree", "add", "-q", "--detach", missing, "HEAD")
			if err := os.RemoveAll(missing); err != nil {
				t.Fatal(err)
			}
			r := runVerb(t, verbList, repoHome{root, Home()}.call())
			primary := listTable + "[1]{id,label,request,state,source,tree,lease,landed,ignored}:\n  foreign," + missing + ",\"\",foreign,foreign,missing,none,unknown,unknown\n"
			want := primary + "help[0]{cmd,why}:\n"
			if r.exit != 0 || r.stdout != want {
				t.Fatalf("ListCommand = (%d, %q), want checked primary plus honest empty help", r.exit, r.stdout)
			}
		})
	}
}

// The help rows that every active row with a present tree shares. The list tests read
// them here, so one expectation pins the text that actionsForRows renders.
const (
	activePathHelpRow = "  bench worktree path <target>,inspect an active worktree by its id\n"
	activeExecHelpRow = "  bench worktree exec <target> -- <command>,run a command in an active worktree by its id\n"
	activeHelpRows    = activePathHelpRow + activeExecHelpRow
)

// activeListRow builds one owned assignment row in the field order the list response
// declares, so an action test reads the same cells the command produces.
func activeListRow(id, request, tree string, landed any, path string) listRow {
	return listRow{
		values:         []any{id, "label", request, string(intent.StateActive), "assignment", tree, "none", landed, "unknown"},
		assignmentPath: path,
	}
}

// TestActionsForRowsReadsTheTreeCell is F10, F11, and F12: an advertised action must be
// one the operator can run, so a row whose tree is gone offers its recovery verb alone,
// and a present row offers the two target-slot actions.
func TestActionsForRowsReadsTheTreeCell(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		row  listRow
		want string
	}{
		{
			name: "missing tree, not landed",
			row:  activeListRow("gone", "req-gone", "missing", false, "/tmp/gone one"),
			want: "help[1]{cmd,why}:\n  bench worktree release --request req-gone '/tmp/gone one',release the assignment whose worktree tree is missing\n",
		},
		{
			name: "missing tree, landed",
			row:  activeListRow("done", "req-done", "missing", true, "/tmp/done"),
			want: "help[1]{cmd,why}:\n  bench worktree clean --landed,clean landed assignments\n",
		},
		{
			name: "present tree",
			row:  activeListRow("here", "req-here", "present", false, "/tmp/here"),
			want: "help[2]{cmd,why}:\n" + activeHelpRows,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			help, err := axi.RenderHelp(actionsForRows([]listRow{tc.row}))
			if err != nil {
				t.Fatal(err)
			}
			if help != tc.want {
				t.Fatalf("help = %q, want %q", help, tc.want)
			}
		})
	}
}

// TestListCommandNamesOneCleanLandedRowForAMissingTree is F13: the landed recovery verb
// is one route, so the response advertises it once however many rows reach it.
func TestListCommandNamesOneCleanLandedRowForAMissingTree(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "list-missing-landed")
	landAssignment(t, f.root, f.creation, "landed.txt")
	if err := os.RemoveAll(f.creation.Path); err != nil {
		t.Fatal(err)
	}
	r := runVerb(t, verbList, repoHome{f.root, Home()}.call())
	if r.exit != 0 {
		t.Fatalf("ListCommand = (%d, %q), want exit 0", r.exit, r.stdout)
	}
	if got := strings.Count(r.stdout, "bench worktree clean --landed"); got != 1 {
		t.Fatalf("ListCommand printed %d clean --landed rows, want 1: %q", got, r.stdout)
	}
	if strings.Contains(r.stdout, "bench worktree path") || strings.Contains(r.stdout, "bench worktree exec") {
		t.Fatalf("ListCommand advertised an action on a missing tree: %q", r.stdout)
	}
}
