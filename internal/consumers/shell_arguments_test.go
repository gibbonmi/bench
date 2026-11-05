package consumers

import (
	"slices"
	"testing"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/axi/axitest"
)

// requireHelpArguments decodes the terminal help block of out and requires each row's command
// to equal its hand-authored spelling and a POSIX shell to read it back to its argv.
func requireHelpArguments(t *testing.T, out string, want []string, argv [][]string) {
	t.Helper()
	document, err := axitest.DecodeDocument(out)
	if err != nil {
		t.Fatal(err)
	}
	actions, err := document.HelpActions()
	if err != nil {
		t.Fatal(err)
	}
	var commands []string
	for _, action := range actions {
		commands = append(commands, action.Cmd)
	}
	if !slices.Equal(commands, want) {
		t.Fatalf("help commands = %q, want %q", commands, want)
	}
	for i, action := range actions {
		got, err := action.Argv()
		if err != nil || !slices.Equal(got, argv[i]) {
			t.Fatalf("argv of %q = (%q, %v), want %q", action.Cmd, got, err, argv[i])
		}
	}
}

// TestConsumersQueryHelpArguments is S21: the over-cap replay and each candidate re-query
// quote every known argument, the filter included, and recover the replayed argv.
func TestConsumersQueryHelpArguments(t *testing.T) {
	t.Run("over-cap replay", func(t *testing.T) {
		stubLoad(t, splitFixture(101, 100))
		out, code := run(t, "target.Symbol", "--production")
		if code != 0 {
			t.Fatalf("exit = %d, want 0; out=%q", code, out)
		}
		requireHelpArguments(t, out,
			[]string{"bench 'consumers' 'target.Symbol' '--production' '--full'"},
			[][]string{{"bench", "consumers", "target.Symbol", "--production", "--full"}})
	})
	t.Run("candidates", func(t *testing.T) {
		stubLoad(t, ambiguousFixture)
		out, code := run(t, "Symbol")
		if code != 0 {
			t.Fatalf("exit = %d, want 0; out=%q", code, out)
		}
		requireHelpArguments(t, out,
			[]string{"bench 'consumers' 'alpha.Symbol'", "bench 'consumers' 'beta.Symbol'"},
			[][]string{{"bench", "consumers", "alpha.Symbol"}, {"bench", "consumers", "beta.Symbol"}})
	})
}

// TestConsumersBlastHelpArguments is S22: the per-symbol walk from a real blast, and the
// over-cap replay of hostile original arguments, quote every known argument. A sentinel file
// reds a replayed value that the shell runs.
func TestConsumersBlastHelpArguments(t *testing.T) {
	t.Run("per-symbol walk", func(t *testing.T) {
		base := blastRepo(t)
		stubLoad(t, tipFixture)
		out, code := run(t, "--changed", "--base", base, "--source-tip", head(t), "--full")
		if code != 0 {
			t.Fatalf("exit = %d, want 0; out=%q", code, out)
		}
		requireHelpArguments(t, out,
			[]string{"bench 'consumers' 'target.Changed' '--full'"},
			[][]string{{"bench", "consumers", "target.Changed", "--full"}})
	})
	t.Run("over-cap replay", func(t *testing.T) {
		sentinel := axitest.NewShellSentinel(t)
		base := "it's $(" + sentinel.Touch() + ")"
		help, err := axi.RenderHelp(blastActions(nil, true, []string{"--changed", "--base", base}))
		if err != nil {
			t.Fatal(err)
		}
		requireHelpArguments(t, help,
			[]string{`bench 'consumers' '--changed' '--base' 'it'\''s $(` + sentinel.Touch() + `)' '--full'`},
			[][]string{{"bench", "consumers", "--changed", "--base", base, "--full"}})
		sentinel.RequireUnexecuted(t)
	})
}
