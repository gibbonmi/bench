package main

import (
	"bytes"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// TestWorktreeLeafHelpAnswersItsOwnGrammar grades SR1. Every row that names a help
// grammar must answer that grammar, so a row paired with another leaf's handler shows
// up as a foreign grammar rather than as a silent reroute. The two sinks are joined
// because a leaf answers its grammar on stdout or on stderr and the row does not grade
// which.
func TestWorktreeLeafHelpAnswersItsOwnGrammar(t *testing.T) {
	for _, leaf := range worktreeLeaves {
		if leaf.Grammar == "" {
			continue
		}
		t.Run(leaf.Name, func(t *testing.T) {
			root := newAXIEnvelopeRepo(t)
			t.Setenv("BENCH_HOME", t.TempDir())
			result := runAXICommandAt(t, root, []string{"worktree", leaf.Name, "--help"})
			out := result.stdout + result.stderr
			want := "usage: " + leaf.Grammar
			if result.code != 0 || !strings.Contains(out, want) {
				t.Fatalf("worktree %s --help = output=%q exit=%d, want exit 0 and %q", leaf.Name, out, result.code, want)
			}
		})
	}
}

// TestWorktreeRootRequiredLeavesRefuseOutsideARepository grades SR2. A leaf whose root
// need was lost would run its handler against an empty root and answer that handler's
// own refusal, so the row demands the one not-in-repo line on stderr at exit 1.
func TestWorktreeRootRequiredLeavesRefuseOutsideARepository(t *testing.T) {
	names := make([]string, 0, len(worktreeLeaves))
	for _, leaf := range worktreeLeaves {
		if leaf.Root == rootRequired {
			names = append(names, leaf.Name)
		}
	}
	want := []string{"exec", "path", "show", "build", "create", "release", "reauthorize", "merge", "reset", "land"}
	if !equalStringSets(names, want) {
		t.Fatalf("root-required leaves = %q, want %q", names, want)
	}
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			t.Setenv("BENCH_HOME", t.TempDir())
			result := runAXICommandAt(t, t.TempDir(), []string{"worktree", name})
			wantErr := toon.NotInRepo() + "\n"
			if result.stdout != "" || result.stderr != wantErr || result.code != 1 {
				t.Fatalf("worktree %s outside a repository = stdout=%q stderr=%q exit=%d, want stderr=%q and exit 1", name, result.stdout, result.stderr, result.code, wantErr)
			}
		})
	}
}

// ER29: the worktree reader is a tree-scoped query leaf. A target that names no active
// assignment refuses through the tree resolver before any child starts. In the caller's
// tree the identity row leads, the framed page fits the response bound, and the
// continuation names the tree and the root-relative path. No call here follows the
// continuation, because its tree target starts a child.
func TestWorktreeReadIsATreeScopedQuery(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	writeAXIFixture(t, filepath.Join(root, "page.txt"), strings.Repeat("page text\n", 600))
	writeAXIFixture(t, filepath.Join(root, "sub", "other.txt"), "")
	t.Chdir(filepath.Join(root, "sub"))
	var stdout, stderr bytes.Buffer
	if code := (Command{Stdout: &stdout, Stderr: &stderr}).Run([]string{"worktree", "read", "--in", "beta", "page.txt"}); code != 1 || stdout.String() != "" || stderr.String() != "bench worktree read --in: target is unassigned\nnext=bench worktree list\n" {
		t.Fatalf("worktree read --in beta = %q %q exit %d, want the unassigned target refusal", stdout.String(), stderr.String(), code)
	}
	stdout.Reset()
	code := (Command{Stdout: &stdout}).Run([]string{"worktree", "read", "../page.txt"})
	if out := stdout.String(); code != 0 || !strings.HasPrefix(out, "tree[1]{target,head,dirty}:\n  primary,") || !responsebound.Fits(out) {
		t.Fatalf("worktree read = %q exit %d, want the identity row first and the framed page within the bound", out, code)
	}
	next, err := axitest.RecoverHelpCommandArgv(withoutTreeRow(stdout.String()))
	if want := []string{"bench", "worktree", "read", "--in", "primary", "page.txt", "--offset"}; err != nil || len(next) < len(want) || !slices.Equal(next[:len(want)], want) {
		t.Fatalf("worktree read continuation = %q (%v), want it to begin %q", next, err, want)
	}
	if !slices.Contains(axiRegistryMemberNames(t), "worktree read") {
		t.Fatal("worktree read is no AXI query member")
	}
}

// TestWorktreeHelpWithAnExtraArgumentRefuses grades SR6. The help form matches on
// exactly one argument, so a table that matched `help` by name alone would answer the
// family usage at exit 0 for a two-word call.
func TestWorktreeHelpWithAnExtraArgumentRefuses(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Command{Stdout: &stdout, Stderr: &stderr, Executable: "bench"}.Run([]string{"worktree", "help", "extra"})
	want := toon.Usage("bench worktree", "help") + "\n"
	if stdout.String() != "" || stderr.String() != want || code != 2 {
		t.Fatalf("worktree help extra = stdout=%q stderr=%q exit=%d, want stderr=%q and exit 2", stdout.String(), stderr.String(), code, want)
	}
}

func equalStringSets(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	seen := make(map[string]int, len(got))
	for _, name := range got {
		seen[name]++
	}
	for _, name := range want {
		seen[name]--
	}
	for _, count := range seen {
		if count != 0 {
			return false
		}
	}
	return true
}

func TestResetHelpIsAvailableAtEveryLevel(t *testing.T) {
	for _, argv := range [][]string{{"worktree", "reset", "--help"}, {"worktree", "--help"}, {"help"}} {
		result := runAXICommandAt(t, newAXIEnvelopeRepo(t), argv)
		if result.code != 0 || !strings.Contains(result.stdout+result.stderr, usage.WorktreeReset) {
			t.Fatalf("%v = %d %s %s, want %s", argv, result.code, result.stdout, result.stderr, usage.WorktreeReset)
		}
	}
}

// keptWorktreeGrammars are the pool operations the worktree family help must keep naming.
// The surviving family route says nothing about the operations under it. Each operation is
// reachable only through that dispatcher, so the grammar line is where its survival shows.
var keptWorktreeGrammars = []string{
	// WF42: the full create grammar, so the family help pins `--from` rather than the
	// bare verb.
	usage.WorktreeCreate,
	"bench worktree path",
	"bench worktree exec",
	usage.WorktreeShow,
	usage.WorktreeBuild,
	"bench worktree release",
	"bench worktree clean",
	"bench worktree reauthorize",
	"bench worktree merge",
	usage.WorktreeReset,
	"bench worktree land",
}
