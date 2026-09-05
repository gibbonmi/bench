package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/toon"
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
	want := []string{"exec", "path", "show", "build", "create", "release", "reauthorize", "merge", "land"}
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
