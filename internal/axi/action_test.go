package axi

import (
	"strings"
	"testing"
)

func TestRenderHelpPreservesEveryKnownArgument(t *testing.T) {
	sha := strings.Repeat("a", 40)
	got, err := RenderHelp([]Action{InspectFull(sha)})
	if err != nil {
		t.Fatal(err)
	}
	want := "help[1]{cmd,why}:\n  bench diff --full --commit " + sha + ",inspect the complete patch\n"
	if got != want {
		t.Fatalf("RenderHelp = %q, want %q", got, want)
	}
}

func TestRenderHelpEmptyIsHonest(t *testing.T) {
	got, err := RenderHelp(nil)
	if err != nil {
		t.Fatal(err)
	}
	if want := "help[0]{cmd,why}:\n"; got != want {
		t.Fatalf("RenderHelp(nil) = %q, want %q", got, want)
	}
}

func TestRenderHelpRejectsNonExecutableOrLossyActions(t *testing.T) {
	sha := strings.Repeat("a", 40)
	for _, tc := range []struct {
		name   string
		action Action
	}{
		{name: "empty command", action: Action{}},
		{name: "known value placeholder", action: Action{kind: actionInspectFull, commit: "<sha>"}},
		{name: "guessed unknown value", action: Action{kind: actionInspectFull, commit: "unknown"}},
		{name: "dropped fixed full flag", action: Action{kind: actionInspectFull, commit: sha, invocation: []string{"diff", "--commit", sha}}},
		{name: "prose as command", action: Action{kind: actionRetryDiff, invocation: []string{"run bench diff --full"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := RenderHelp([]Action{tc.action}); err == nil {
				t.Fatal("RenderHelp accepted a non-executable or lossy action")
			}
		})
	}
}
