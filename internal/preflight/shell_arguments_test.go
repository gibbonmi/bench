package preflight

import (
	"slices"
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
	document, err := axitest.DecodeDocument(help)
	if err != nil {
		t.Fatal(err)
	}
	actions, err := document.HelpActions()
	want := []axitest.HelpAction{{Cmd: `bench 'preflight' 'build' 'it'\''s; ` + sentinel.Touch() + `' '--base' 'HEAD~1'`, Why: axi.RetryAfterMovement}}
	if err != nil || !slices.Equal(actions, want) {
		t.Fatalf("help actions = (%q, %v), want %q", actions, err, want)
	}
	argv, err := actions[0].Argv()
	if wantArgv := []string{"bench", "preflight", "build", slug, "--base", "HEAD~1"}; err != nil || !slices.Equal(argv, wantArgv) {
		t.Fatalf("argv of %q = (%q, %v), want %q", actions[0].Cmd, argv, err, wantArgv)
	}
	sentinel.RequireUnexecuted(t)
}
