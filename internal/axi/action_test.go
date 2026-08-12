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

func TestRenderHelpRendersReusableInvocationAndHarnessPhase(t *testing.T) {
	got, err := RenderHelp([]Action{
		ExecutableInvocation("repair decisions/model.md", KnownArgument("maps"), KnownArgument("--template")),
		ExecutableInvocation("inspect active worktree", KnownArgument("worktree"), KnownArgument("exec"), KnownArgument("alpha"), KnownArgument("--"), FutureInput("command")),
		HarnessPhase("/bench-shape-idea", "shape model"),
		HarnessPhase("/bench-what-next", "drain the capture inbox"),
	})
	if err != nil {
		t.Fatal(err)
	}
	const want = "help[4]{cmd,why}:\n  bench maps --template,repair decisions/model.md\n  bench worktree exec alpha -- <command>,inspect active worktree\n  /bench-shape-idea,shape model\n  /bench-what-next,drain the capture inbox\n"
	if got != want {
		t.Fatalf("RenderHelp = %q, want %q", got, want)
	}
}

func TestRenderHelpRendersKnownArgumentNamedUnknown(t *testing.T) {
	got, err := RenderHelp([]Action{ExecutableInvocation("inspect an unknown-named map", KnownArgument("maps"), KnownArgument("unknown"))})
	if err != nil {
		t.Fatal(err)
	}
	const want = "help[1]{cmd,why}:\n  bench maps unknown,inspect an unknown-named map\n"
	if got != want {
		t.Fatalf("RenderHelp = %q, want %q", got, want)
	}
}

func TestRenderHelpShellQuotesKnownArgumentsWithoutChangingFutureInputs(t *testing.T) {
	got, err := RenderHelp([]Action{ExecutableInvocation(
		"clean the orphaned worktree",
		KnownArgument("worktree"),
		KnownArgument("clean"),
		KnownArgument("orphan pool/a*.wt"),
		KnownArgument("owner's copy"),
		FutureInput("confirmation"),
	)})
	if err != nil {
		t.Fatal(err)
	}
	const want = "help[1]{cmd,why}:\n  \"bench worktree clean 'orphan pool/a*.wt' 'owner'\\\"'\\\"'s copy' <confirmation>\",clean the orphaned worktree\n"
	if got != want {
		t.Fatalf("RenderHelp = %q, want %q", got, want)
	}
}

func TestRenderHelpRejectsUnapprovedHarnessPhase(t *testing.T) {
	if _, err := RenderHelp([]Action{HarnessPhase("/foo", "do something")}); err == nil {
		t.Fatal("RenderHelp accepted an unapproved harness phase")
	}
}

func TestRenderHelpRejectsLossyReusableActions(t *testing.T) {
	for _, tc := range []struct {
		name   string
		action Action
	}{
		{
			name: "dropped known argument",
			action: Action{
				kind:       actionInvocation,
				arguments:  []InvocationArgument{KnownArgument("maps"), KnownArgument("--template")},
				invocation: []string{"maps"},
				why:        "repair decisions/model.md",
			},
		},
		{name: "empty known value", action: ExecutableInvocation("repair", KnownArgument("maps"), KnownArgument(""))},
		{name: "undeclared placeholder", action: ExecutableInvocation("repair", KnownArgument("maps"), KnownArgument("<path>"))},
		{name: "control byte", action: ExecutableInvocation("repair", KnownArgument("worktree"), KnownArgument("clean"), KnownArgument("orphan\x1bpool"))},
		{name: "prose as command", action: ExecutableInvocation("repair", KnownArgument("run bench maps --template"))},
		{name: "unsafe executable name", action: ExecutableInvocation("repair", KnownArgument("maps!"))},
		{name: "future command", action: ExecutableInvocation("repair", FutureInput("command"))},
		{name: "shell command phase", action: HarnessPhase("bench /bench-shape-idea", "shape model")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := RenderHelp([]Action{tc.action}); err == nil {
				t.Fatal("RenderHelp accepted a lossy reusable action")
			}
		})
	}
}

func TestRenderHelpDeduplicatesExactTemplatesWithoutReordering(t *testing.T) {
	got, err := RenderHelp([]Action{
		HarnessPhase("/bench-shape-idea", "shape alpha"),
		HarnessPhase("/bench-shape-idea", "shape alpha"),
		ExecutableInvocation("repair decisions/broken.md", KnownArgument("maps"), KnownArgument("--template")),
	})
	if err != nil {
		t.Fatal(err)
	}
	const want = "help[2]{cmd,why}:\n  /bench-shape-idea,shape alpha\n  bench maps --template,repair decisions/broken.md\n"
	if got != want {
		t.Fatalf("RenderHelp = %q, want %q", got, want)
	}
}
