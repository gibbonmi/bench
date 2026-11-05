package coverage

import (
	"slices"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
)

// TestCoverageHelpArguments is S20: the real coverage producer quotes the verb, the flag, and
// the spec path of each check and retry action. Each command is authored by hand, a POSIX
// shell reads it back to the typed path, and a sentinel file reds a path that the shell runs.
func TestCoverageHelpArguments(t *testing.T) {
	sentinel := axitest.NewShellSentinel(t)
	mapped := func(row string) string {
		return "# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReduced + row
	}
	path := "it's; " + sentinel.Touch() + "/spec.md"
	quoted := `'it'\''s; ` + sentinel.Touch() + `/spec.md'`
	for _, row := range []struct {
		name, body, why string
	}{
		{name: "check", body: mapped("| 1 | b | s | catches one |\n"), why: "check coverage for stories 1"},
		{name: "retry", body: mapped("| 9 | b | s | catches one |\n"), why: "retry after repairing coverage map"},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Chdir(t.TempDir())
			mustWrite(t, path, row.body)
			out, code := Command([]string{path})
			if code != 0 {
				t.Fatalf("Command = (%d, %q), want exit 0", code, out)
			}
			document, err := axitest.DecodeDocument(out)
			if err != nil {
				t.Fatal(err)
			}
			actions, err := document.HelpActions()
			want := []axitest.HelpAction{{Cmd: "bench 'coverage' '--check' " + quoted, Why: row.why}}
			if err != nil || !slices.Equal(actions, want) {
				t.Fatalf("help actions = (%q, %v), want %q", actions, err, want)
			}
			argv, err := actions[0].Argv()
			if wantArgv := []string{"bench", "coverage", "--check", path}; err != nil || !slices.Equal(argv, wantArgv) {
				t.Fatalf("argv of %q = (%q, %v), want %q", actions[0].Cmd, argv, err, wantArgv)
			}
		})
	}
	sentinel.RequireUnexecuted(t)
}
