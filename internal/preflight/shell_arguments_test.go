package preflight

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/axi/axitest"
)

// TestPreflightHelpArguments is S26: the snapshot-drift refusal retries the exact invocation
// with every original argument quoted, in its original order. The command is authored by
// hand, a POSIX shell reads it back, and a sentinel file reds an argument that the shell runs.
func TestPreflightHelpArguments(t *testing.T) {
	sentinel := axitest.NewShellSentinel(t)
	slug := "it's; " + sentinel.Touch()
	out := snapshotDriftRefusal([]string{"build", slug, "--base", "HEAD~1"}, "the dirty content changed while reading")
	_, help, ok := strings.Cut(out, "\n")
	if !ok {
		t.Fatalf("refusal = %q, want an error line and a help block", out)
	}
	actions := axitest.RequireHelpCommands(t, help,
		[]string{`bench 'preflight' 'build' 'it'\''s; ` + sentinel.Touch() + `' '--base' 'HEAD~1'`},
		[][]string{{"bench", "preflight", "build", slug, "--base", "HEAD~1"}})
	if actions[0].Why != axi.RetryAfterMovement {
		t.Fatalf("help why = %q, want %q", actions[0].Why, axi.RetryAfterMovement)
	}
	sentinel.RequireUnexecuted(t)
}
