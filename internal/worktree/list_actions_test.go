package worktree

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/axitest"
)

type worktreeListArgvPair struct {
	Argv []string `json:"argv"`
	Old  struct {
		Stdout string `json:"stdout"`
		Stderr string `json:"stderr"`
		Exit   int    `json:"exit"`
	} `json:"old"`
	New struct {
		Stdout string `json:"stdout"`
		Stderr string `json:"stderr"`
		Exit   int    `json:"exit"`
	} `json:"new"`
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
	root := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", root).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	t.Chdir(root)
	for _, pair := range pairs {
		out, code := ListCommand(pair.Argv)
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
	root := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", root).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	t.Chdir(root)
	for _, arg := range []string{"--help", "-h", "help"} {
		out, code := ListCommand([]string{arg})
		if code != 0 || out != "usage: bench worktree list\n" {
			t.Errorf("ListCommand(%q) = (%d, %q)", arg, code, out)
		}
	}
	for _, args := range [][]string{{"--unknown"}, {"extra"}, {"--"}} {
		out, code := ListCommand(args)
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
	root := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", root).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	t.Chdir(root)
	out, code := ListCommand(nil)
	if code != 0 || out != string(primary)+"help[0]{cmd,why}:\n" {
		t.Fatalf("ListCommand = (%d, %q), want checked-in primary plus exactly one help block", code, out)
	}
}

func TestActionsForRowsEnumeratesActiveAndOrphanRows(t *testing.T) {
	rows := [][]any{
		{"a", "active", "active"},
		{"done", "complete", "complete"},
		{"foreign", "one", "foreign", "foreign", "missing"},
		{"b", "active", "active"},
		{"foreign", "present", "foreign", "foreign", "present"},
		{"foreign", "two", "foreign", "foreign", "missing"},
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
	t.Setenv("BENCH_HOME", filepath.Join(t.TempDir(), "bench-home"))
	a := mustCreate(t, root, "request-a", "alpha")
	b := mustCreate(t, root, "request-b", "beta")
	present := filepath.Join(t.TempDir(), "present foreign")
	missing := filepath.Join(t.TempDir(), "missing foreign * path")
	gitRun(t, root, "worktree", "add", "-q", "--detach", present, "HEAD")
	gitRun(t, root, "worktree", "add", "-q", "--detach", missing, "HEAD")
	if err := os.RemoveAll(missing); err != nil {
		t.Fatal(err)
	}
	chdir(t, root)
	out, code := ListCommand(nil)
	first, second := a, b
	if first.Assignment.ID > second.Assignment.ID {
		first, second = second, first
	}
	primary := strings.NewReplacer(
		"{{ID1}}", first.Assignment.ID,
		"{{LABEL1}}", first.Assignment.Label,
		"{{ID2}}", second.Assignment.ID,
		"{{LABEL2}}", second.Assignment.Label,
		"{{PRESENT}}", present,
		"{{MISSING}}", missing,
	).Replace(string(primaryTemplate))
	help := fmt.Sprintf("help[5]{cmd,why}:\n  bench worktree path %s,inspect active worktree\n  bench worktree exec %s -- <command>,run a command in the active worktree\n  bench worktree path %s,inspect active worktree\n  bench worktree exec %s -- <command>,run a command in the active worktree\n  bench worktree clean '%s',clean the orphaned worktree\n", first.Assignment.ID, first.Assignment.ID, second.Assignment.ID, second.Assignment.ID, missing)
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
			chdir(t, root)
			out, code := ListCommand(nil)
			if code != 0 || !strings.HasPrefix(out, "worktrees[1]{id,label,state,source,tree,lease,landed,ignored}:\n") {
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
			chdir(t, root)
			out, code := ListCommand(nil)
			primary := "worktrees[1]{id,label,state,source,tree,lease,landed,ignored}:\n  foreign," + missing + ",foreign,foreign,missing,none,unknown,unknown\n"
			want := primary + "help[0]{cmd,why}:\n"
			if code != 0 || out != want {
				t.Fatalf("ListCommand = (%d, %q), want checked primary plus honest empty help", code, out)
			}
		})
	}
}
