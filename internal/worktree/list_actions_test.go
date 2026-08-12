package worktree

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/axi"
)

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

func TestActionsForRowsEnumeratesActiveAndOrphanRows(t *testing.T) {
	rows := [][]any{
		{"a", "active", "active"},
		{"done", "complete", "complete"},
		{"foreign", "one", "foreign", "foreign", "missing"},
		{"b", "active", "active"},
		{"foreign", "present", "foreign", "foreign", "present"},
		{"foreign", "two", "foreign", "foreign", "missing"},
	}
	help, err := axi.RenderHelp(actionsForRows(rows, []string{"/tmp/orphan one", "/tmp/orphan-two"}))
	if err != nil {
		t.Fatal(err)
	}
	want := "help[6]{cmd,why}:\n  bench worktree path a,inspect active worktree\n  bench worktree exec a -- <command>,run a command in the active worktree\n  bench worktree clean '/tmp/orphan one',clean the orphaned worktree\n  bench worktree path b,inspect active worktree\n  bench worktree exec b -- <command>,run a command in the active worktree\n  bench worktree clean /tmp/orphan-two,clean the orphaned worktree\n"
	if help != want {
		t.Fatalf("help = %q, want %q", help, want)
	}
}

func TestListCommandPublicRowsAndDisclosure(t *testing.T) {
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
	want := fmt.Sprintf("worktrees[4]{id,label,state,source,tree,lease,landed,ignored}:\n  %s,%s,active,assignment,present,none,true,0\n  %s,%s,active,assignment,present,none,true,0\n  foreign,%s,foreign,foreign,present,none,unknown,0\n  foreign,%s,foreign,foreign,missing,none,unknown,unknown\nhelp[5]{cmd,why}:\n  bench worktree path %s,inspect active worktree\n  bench worktree exec %s -- <command>,run a command in the active worktree\n  bench worktree path %s,inspect active worktree\n  bench worktree exec %s -- <command>,run a command in the active worktree\n  bench worktree clean '%s',clean the orphaned worktree\n", first.Assignment.ID, first.Assignment.Label, second.Assignment.ID, second.Assignment.Label, present, missing, first.Assignment.ID, first.Assignment.ID, second.Assignment.ID, second.Assignment.ID, missing)
	if code != 0 || out != want {
		t.Fatalf("ListCommand = (%d, %q), want %q", code, out, want)
	}
}
