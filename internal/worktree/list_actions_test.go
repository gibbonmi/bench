package worktree

import (
	"bytes"
	cryptorand "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	root := journeyRepoOnBranch(t, "main")
	journeyFIFOWorktreeAdmin(t, root, "typed")
	out, code := ListCommand(root, nil)
	if code == 0 || !strings.Contains(out, "worktrees/typed/gitdir") || !strings.Contains(out, "fifo") || !strings.Contains(out, "inspect and remove it") {
		t.Fatalf("typed list output code=%d out=%q", code, out)
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
			out, code := ListCommand(root, nil)
			if code != 1 || !strings.Contains(out, tc.detail) || !strings.Contains(out, tc.action) {
				t.Fatalf("%s list output code=%d out=%q", tc.mode, code, out)
			}
		})
	}
}

func TestListCommandRendersBoundExpiryAsTypedFailure(t *testing.T) {
	restore := git.SetWorktreeListTimeoutForTest(100 * time.Millisecond)
	t.Cleanup(restore)
	root := journeyRepoOnBranch(t, "main")
	journeyStubGit(t, root, "block-worktree", filepath.Join(t.TempDir(), "argv"))
	out, code := ListCommand(root, nil)
	if code != 1 || !strings.Contains(out, "worktree list") || !strings.Contains(out, "investigate the git failure") || strings.Contains(out, "inspect and remove it") || strings.Contains(out, "retry") {
		t.Fatalf("bound list output code=%d out=%q", code, out)
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
		out, code := ListCommand(root, pair.Argv)
		if out != pair.New.Stdout || pair.New.Stderr != "" || code != pair.New.Exit {
			t.Fatalf("ListCommand(%q) = stdout=%q stderr=%q exit=%d, want checked-in new response", pair.Argv, out, "", code)
		}
		if pair.Old != pair.New {
			if len(pair.Argv) != 1 || (pair.Argv[0] != "--help" && pair.Argv[0] != "-h" && pair.Argv[0] != "help") {
				t.Fatalf("paired fixture admits an unapproved argv delta: %#v", pair)
			}
		}
	}
}

func TestListCommandHelpAndArgumentMatrix(t *testing.T) {
	root := journeyRepo(t)
	for _, arg := range []string{"--help", "-h", "help"} {
		out, code := ListCommand(root, []string{arg})
		if code != 0 || out != "usage: bench worktree list\n" {
			t.Errorf("ListCommand(%q) = (%d, %q)", arg, code, out)
		}
	}
	for _, args := range [][]string{{"--unknown"}, {"extra"}, {"--"}} {
		out, code := ListCommand(root, args)
		want := "usage: bench worktree list (unknown argument: " + args[0] + ")\n"
		if code != 2 || out != want {
			t.Errorf("ListCommand(%q) = (%d, %q), want usage exit 2", args, code, out)
		}
	}
}

func TestListCommandPreservesCheckedInEmptyPrimaryResponse(t *testing.T) {
	primary, err := os.ReadFile("testdata/pre-disclosure-empty.stdout")
	if err != nil {
		t.Fatal(err)
	}
	root := journeyRepo(t)
	out, code := ListCommand(root, nil)
	if code != 0 || out != string(primary)+"help[0]{cmd,why}:\n" {
		t.Fatalf("ListCommand = (%d, %q), want checked-in primary plus exactly one help block", code, out)
	}
}

func TestListCommandCheckedInPresentForeignTerminalPair(t *testing.T) {
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
	out, code := ListCommand(root, nil)
	pair.Old.Stdout = strings.ReplaceAll(pair.Old.Stdout, "{{PRESENT}}", present)
	pair.New.Stdout = strings.ReplaceAll(pair.New.Stdout, "{{PRESENT}}", present)
	if pair.Old.Stderr != "" || pair.New.Stderr != "" || pair.Old.Exit != 0 || pair.New.Exit != 0 || pair.New.Stdout != pair.Old.Stdout+"help[0]{cmd,why}:\n" {
		t.Fatalf("terminal pair admits a response change beyond exactly one empty help block: %#v", pair)
	}
	if code != pair.New.Exit || out != pair.New.Stdout {
		t.Fatalf("ListCommand = stdout=%q stderr=%q exit=%d, want checked-in terminal primary plus exactly one empty help block", out, "", code)
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
	creation := mustCreate(t, root, "landed-complete-assignment", "complete assignment")
	boundary := cleanupTransactionBoundary
	t.Cleanup(func() { cleanupTransactionBoundary = boundary })
	cleanupTransactionBoundary = func(step LifecycleStep) error {
		if step == StepTerminalReceipt {
			return errors.New("stop before the completed record is compacted")
		}
		return nil
	}
	if code := ReleaseCommand(root, []string{"--request", "landed-complete-assignment", creation.Path}, io.Discard, io.Discard); code == 0 {
		t.Fatalf("interrupted release exit = %d, want non-zero", code)
	}
	cleanupTransactionBoundary = boundary
	assignments, err := intent.Assignments(root)
	if err != nil || len(assignments) != 1 || assignments[0].State != intent.StateComplete {
		t.Fatalf("Assignments = %#v, %v, want one completed assignment", assignments, err)
	}
	out, code := ListCommand(root, nil)
	materialize := strings.NewReplacer("{{LABEL}}", assignments[0].Label)
	pair.Old.Stdout = materialize.Replace(pair.Old.Stdout)
	pair.New.Stdout = materialize.Replace(pair.New.Stdout)
	if pair.Old.Stderr != "" || pair.New.Stderr != "" || pair.Old.Exit != 0 || pair.New.Exit != 0 || pair.New.Stdout != pair.Old.Stdout+"help[0]{cmd,why}:\n" {
		t.Fatalf("terminal pair admits a response change beyond exactly one empty help block: %#v", pair)
	}
	if code != pair.New.Exit || out != pair.New.Stdout {
		t.Fatalf("ListCommand = stdout=%q stderr=%q exit=%d, want checked-in completed-assignment primary plus exactly one empty help block", out, "", code)
	}
}

func TestActionsForRowsEnumeratesActiveAndOrphanRows(t *testing.T) {
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
	want := "help[6]{cmd,why}:\n  bench worktree path a,inspect active worktree\n  bench worktree exec a -- <command>,run a command in the active worktree\n  bench worktree clean '/tmp/orphan one',clean the orphaned worktree\n  bench worktree path b,inspect active worktree\n  bench worktree exec b -- <command>,run a command in the active worktree\n  bench worktree clean /tmp/orphan-two,clean the orphaned worktree\n"
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
	mustCreate(t, root, "request-a", "alpha")
	mustCreate(t, root, "request-b", "beta")
	present := filepath.Join(t.TempDir(), "present foreign")
	missing := filepath.Join(t.TempDir(), "missing foreign * path")
	gitRun(t, root, "worktree", "add", "-q", "--detach", present, "HEAD")
	gitRun(t, root, "worktree", "add", "-q", "--detach", missing, "HEAD")
	if err := os.RemoveAll(missing); err != nil {
		t.Fatal(err)
	}
	out, code := ListCommand(root, nil)
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
	help := fmt.Sprintf("help[6]{cmd,why}:\n  bench worktree path %s,inspect active worktree\n  bench worktree exec %s -- <command>,run a command in the active worktree\n  bench worktree path %s,inspect active worktree\n  bench worktree exec %s -- <command>,run a command in the active worktree\n  bench worktree clean '%s',clean the orphaned worktree\n  bench worktree clean --landed,clean landed assignments\n", assignments[0].ID, assignments[0].ID, assignments[1].ID, assignments[1].ID, missing)
	if code != 0 || out != primary+help {
		t.Fatalf("ListCommand = (%d, %q), want materialized checked-in primary plus exactly one help block", code, out)
	}
}

func TestListCommandControlBearingOrphanPathPreservesPrimaryAndAction(t *testing.T) {
	for _, control := range []string{"\t", "\n", "\r"} {
		t.Run("control", func(t *testing.T) {
			root := newWorktreeRepo(t)
			missing := filepath.Join(t.TempDir(), "orphan"+control+"path")
			gitRun(t, root, "worktree", "add", "-q", "--detach", missing, "HEAD")
			if err := os.RemoveAll(missing); err != nil {
				t.Fatal(err)
			}
			out, code := ListCommand(root, nil)
			if code != 0 || !strings.HasPrefix(out, "worktrees[1]{id,label,request,state,source,tree,lease,landed,ignored}:\n") {
				t.Fatalf("ListCommand = (%d, %q), want primary worktree response and exit 0", code, out)
			}
			if !strings.Contains(out, "help[1]{cmd,why}:\n") {
				t.Fatalf("ListCommand = %q, want one orphan-clean action", out)
			}
			argv, err := axitest.RecoverHelpCommandArgv(out)
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
	for _, marker := range []string{"<", ">"} {
		t.Run(marker, func(t *testing.T) {
			root := newWorktreeRepo(t)
			missing := filepath.Join(t.TempDir(), "orphan"+marker+"path")
			gitRun(t, root, "worktree", "add", "-q", "--detach", missing, "HEAD")
			if err := os.RemoveAll(missing); err != nil {
				t.Fatal(err)
			}
			out, code := ListCommand(root, nil)
			primary := "worktrees[1]{id,label,request,state,source,tree,lease,landed,ignored}:\n  foreign," + missing + ",\"\",foreign,foreign,missing,none,unknown,unknown\n"
			want := primary + "help[0]{cmd,why}:\n"
			if code != 0 || out != want {
				t.Fatalf("ListCommand = (%d, %q), want checked primary plus honest empty help", code, out)
			}
		})
	}
}
