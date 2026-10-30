package main

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/usage"
)

// TestHelpSpecRowsNameRetireHistoryAndWiden pins FA9: the retired subcommand keeps no
// help row, while the live retire, history, and widen rows stay advertised.
func TestHelpSpecRowsNameRetireHistoryAndWiden(t *testing.T) {
	help := renderCommandHelp()
	if strings.Contains(help, "bench spec implemented") {
		t.Errorf("bench help still advertises a retired subcommand:\n%s", help)
	}
	for _, want := range []string{"bench spec [--in <label|primary>] retire <slug>", "bench spec [--in <label|primary>] history <slug>", "bench spec [--in <label|primary>] widen <slug>"} {
		if !strings.Contains(help, want) {
			t.Errorf("bench help is missing %q:\n%s", want, help)
		}
	}
}

func TestTestHelpNamesOnlyRunnableFocusedForms(t *testing.T) {
	help := renderCommandHelp()
	if !strings.Contains(help, "bench test [--in <label|primary>] [--full] [--package <expr> [--package <expr>]... | <legacy-package> | --changed] [--base <commit> [--source-tip <commit>]] [--run <go-regex>] | bench test [--full] --check <name> | bench test [--full] --check system --run <go-regex> | bench test --check <name> --fixtures | bench test --checks") {
		t.Errorf("bench help is missing the focused test grammar:\n%s", help)
	}
}

func TestFocusedHelpNamesCheckAsEvidenceWithoutVerdict(t *testing.T) {
	help := renderCommandHelp()
	if !strings.Contains(help, "bench test [--full] --check <name>") {
		t.Fatalf("bench help is missing the named-check form:\n%s", help)
	}
	if strings.Contains(strings.ToLower(help), "test --check <name>  gate") {
		t.Fatalf("bench help presents a focused check as a verdict:\n%s", help)
	}
}

func TestWorktreeCleanHelpUsesItsGrammar(t *testing.T) {
	result := runAXICommandAt(t, t.TempDir(), []string{"worktree", "clean", "--help"})
	want := "usage: " + usage.WorktreeClean + "\n"
	if result.stdout != want || result.stderr != "" || result.code != 0 {
		t.Fatalf("worktree clean --help = stdout=%q stderr=%q exit=%d, want stdout=%q stderr=\"\" exit=0", result.stdout, result.stderr, result.code, want)
	}
}

// TestWorktreeHelpNamesTheExecGateForm pins WF14. The worktree usage trailer names the
// gate as a worktree-native exec form, and it names no raw wrapper path. An appended
// line beside the old trailer would keep the raw path alive, so the assertion demands
// the trailer be the last line and demands the whole text hold no bin/bench.sh.
func TestWorktreeHelpNamesTheExecGateForm(t *testing.T) {
	out, code := runKeptRoute([]string{"worktree", "--help"})
	if code != 0 {
		t.Fatalf("worktree --help exit = %d, want 0; output=%q", code, out)
	}
	const trailer = "bench worktree exec <target> -- bench gate"
	lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
	if last := strings.TrimSpace(lines[len(lines)-1]); last != trailer {
		t.Errorf("worktree help last line = %q, want %q", last, trailer)
	}
	if strings.Contains(out, "bin/bench.sh") {
		t.Errorf("worktree help = %q, want no bin/bench.sh", out)
	}
}
