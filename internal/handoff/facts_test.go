package handoff

import (
	"testing"

	"github.com/gibbonmi/bench/internal/status"
)

// TestFirstInvocable enumerates the selection's edges, including the all-prose board a live
// fixture cannot be relied on to produce. The prose actions are the board's own, so a rule
// that admitted any of them would show up here rather than in the rendered field.
func TestFirstInvocable(t *testing.T) {
	prose := []status.Signal{
		{Severity: 0, Name: "gate", Detail: "red", Action: "fix before commit"},
		{Severity: 1, Name: "git", Detail: "1 dirty path", Action: "commit on green"},
		{Severity: 5, Name: "structure", Detail: "2 issue(s)", Action: "split (craft-seams)"},
	}
	cases := []struct {
		what        string
		signals     []status.Signal
		action      string
		wantName    string
		wantPresent bool
	}{
		{what: "no board at all", signals: nil},
		{what: "empty board", signals: []status.Signal{}},
		{what: "every action is prose", signals: prose},
		{
			what: "prose ahead of a command",
			signals: append(append([]status.Signal{}, prose...),
				status.Signal{Severity: 4, Name: "drain", Detail: "1 idea(s)", Action: "/bench-what-next"}),
			action: "/bench-what-next", wantName: "drain", wantPresent: true,
		},
		{
			what: "a command leads",
			signals: []status.Signal{
				{Severity: 3, Name: "guards", Detail: "pre-push hook missing", Action: "bench link"},
				{Severity: 4, Name: "drain", Detail: "1 idea(s)", Action: "/bench-what-next"},
			},
			action: "bench link", wantName: "guards", wantPresent: true,
		},
		{
			what: "a `bench` subcommand carrying an argument placeholder",
			signals: []status.Signal{
				{Severity: 5, Name: "structure", Detail: "2 issue(s)", Action: "split (craft-seams)"},
				{Severity: 8, Name: "specs", Detail: "1 merged spec(s)", Action: "bench spec retire <slug>"},
			},
			action: "bench spec retire <slug>", wantName: "specs", wantPresent: true,
		},
		{
			// The git row takes this exact shape once the tree is clean and commits are
			// unpushed. It opens a phase invocation, so a bare prefix test admits it and
			// the field renders two commands where it promised one.
			what: "a compound action naming several steps is not one invocation",
			signals: []status.Signal{
				{Severity: 1, Name: "git", Detail: "7 unpushed commits", Action: "/bench-final-check / push"},
				{Severity: 4, Name: "drain", Detail: "1 idea(s)", Action: "/bench-what-next"},
			},
			action: "/bench-what-next", wantName: "drain", wantPresent: true,
		},
		{
			what: "a word beginning with bench is not a command",
			signals: []status.Signal{
				{Severity: 11, Name: "handoff", Detail: "stale", Action: "benchmark the loop"},
			},
		},
	}
	for _, tc := range cases {
		action, name, ok := firstInvocable(tc.signals)
		if ok != tc.wantPresent || action != tc.action || name != tc.wantName {
			t.Fatalf("%s: firstInvocable = (%q, %q, %v), want (%q, %q, %v)",
				tc.what, action, name, ok, tc.action, tc.wantName, tc.wantPresent)
		}
	}
}
