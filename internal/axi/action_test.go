package axi

import (
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axitest"
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

func TestKnownArgumentControlValuesRoundTripThroughTOONAndPOSIXShell(t *testing.T) {
	for _, source := range []string{"tab\tvalue", "line\nvalue", "return\rvalue"} {
		t.Run(source, func(t *testing.T) {
			rendered, err := RenderHelp([]Action{ExecutableInvocation("clean", KnownArgument("worktree"), KnownArgument("clean"), KnownArgument(source))})
			if err != nil {
				t.Fatal(err)
			}
			argv, err := axitest.RecoverHelpCommandArgv(rendered)
			if err != nil {
				t.Fatal(err)
			}
			want := []string{"bench", "worktree", "clean", source}
			if !slices.Equal(argv, want) {
				t.Fatalf("shell argv = %q, want %q", argv, want)
			}
		})
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

func TestRenderHelpUsesHonestEmptyForUnsupportedDisclosureValue(t *testing.T) {
	for _, action := range []Action{
		ExecutableInvocation("repair", KnownArgument("maps"), KnownArgument("")),
		ExecutableInvocation("repair", KnownArgument("maps"), KnownArgument("decisions/<broken>.md")),
		ExecutableInvocation("repair", KnownArgument("worktree"), KnownArgument("clean"), KnownArgument("orphan\x1bpool")),
		ExecutableInvocation("repair malformed\x1bmap", KnownArgument("maps"), KnownArgument("--template")),
	} {
		got, err := RenderHelp([]Action{action})
		if err != nil {
			t.Fatal(err)
		}
		if got != "help[0]{cmd,why}:\n" {
			t.Fatalf("RenderHelp = %q, want honest empty help", got)
		}
	}
}

func TestKnownArgumentRefusesUnsupportedControls(t *testing.T) {
	controls := make([]rune, 0, 61)
	for r := rune(0); r < 0x20; r++ {
		if r != '\t' && r != '\n' && r != '\r' {
			controls = append(controls, r)
		}
	}
	controls = append(controls, 0x7f)
	for r := rune(0x80); r <= 0x9f; r++ {
		controls = append(controls, r)
	}
	for _, r := range controls {
		got, err := RenderHelp([]Action{ExecutableInvocation("repair", KnownArgument("maps"), KnownArgument(string(r)))})
		if err != nil {
			t.Fatalf("RenderHelp(%U): %v", r, err)
		}
		if got != "help[0]{cmd,why}:\n" {
			t.Fatalf("RenderHelp(%U) = %q, want honest empty help", r, got)
		}
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
