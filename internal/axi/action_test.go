package axi

import (
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
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
		HarnessPhase("/bench-write-spec", "spec the ready map"),
		HarnessPhaseOn("/bench-write-spec", "decisions/my map.md", "spec my map: Model"),
	})
	if err != nil {
		t.Fatal(err)
	}
	const want = "help[6]{cmd,why}:\n  bench 'maps' '--template',repair decisions/model.md\n  bench 'worktree' 'exec' 'alpha' '--' <command>,inspect active worktree\n  /bench-shape-idea,shape model\n  /bench-what-next,drain the capture inbox\n  /bench-write-spec,spec the ready map\n  /bench-write-spec decisions/my map.md,\"spec my map: Model\"\n"
	if got != want {
		t.Fatalf("RenderHelp = %q, want %q", got, want)
	}
}

func TestRenderHelpRendersKnownArgumentNamedUnknown(t *testing.T) {
	got, err := RenderHelp([]Action{ExecutableInvocation("inspect an unknown-named map", KnownArgument("maps"), KnownArgument("unknown"))})
	if err != nil {
		t.Fatal(err)
	}
	const want = "help[1]{cmd,why}:\n  bench 'maps' 'unknown',inspect an unknown-named map\n"
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
	const want = "help[1]{cmd,why}:\n  \"bench 'worktree' 'clean' 'orphan pool/a*.wt' 'owner'\\\\''s copy' <confirmation>\",clean the orphaned worktree\n"
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
		{name: "non-canonical phase with an argument", action: HarnessPhaseOn("/bench-compile-map", "decisions/model.md", "spec model")},
		{name: "placeholder phase argument", action: HarnessPhaseOn("/bench-write-spec", "decisions/<map>.md", "spec model")},
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
	const want = "help[2]{cmd,why}:\n  /bench-shape-idea,shape alpha\n  bench 'maps' '--template',repair decisions/broken.md\n"
	if got != want {
		t.Fatalf("RenderHelp = %q, want %q", got, want)
	}
}

// TestRenderHelpAlwaysQuoted is S05: every accepted known argument renders as one
// always-quoted token, a safe token included, so the former bare-token spelling reds the
// exact bytes. The expected cell is authored by hand, and a POSIX shell then reads the
// decoded cell back to the declared values. A sentinel file reds a value that the shell runs.
func TestRenderHelpAlwaysQuoted(t *testing.T) {
	sentinel := axitest.NewShellSentinel(t)
	hostile := "it's $(" + sentinel.Touch() + ")"
	action := ExecutableInvocation("clean the orphaned worktree", KnownArgument("worktree"), KnownArgument("clean"), KnownArgument("pool/a.wt"), KnownArgument(hostile))
	got, err := RenderHelp([]Action{action})
	if err != nil {
		t.Fatal(err)
	}
	command := `bench 'worktree' 'clean' 'pool/a.wt' 'it'\''s $(` + sentinel.Touch() + `)'`
	want := "help[1]{cmd,why}:\n  \"" + strings.ReplaceAll(command, `\`, `\\`) + "\",clean the orphaned worktree\n"
	if got != want {
		t.Fatalf("RenderHelp = %q, want %q", got, want)
	}
	if line, err := action.Command(); err != nil || line != command {
		t.Fatalf("Command = (%q, %v), want %q", line, err, command)
	}
	argv, err := axitest.RecoverHelpCommandArgv(got)
	if wantArgv := []string{"bench", "worktree", "clean", "pool/a.wt", hostile}; err != nil || !slices.Equal(argv, wantArgv) {
		t.Fatalf("shell argv = (%q, %v), want %q", argv, err, wantArgv)
	}
	sentinel.RequireUnexecuted(t)
}

// TestRenderHelpFutureInputs is S06: a future-input slot keeps its exact placeholder bytes
// beside always-quoted known arguments. The placeholder is never executed; only the known
// prefix before the first slot is read back through a POSIX shell.
func TestRenderHelpFutureInputs(t *testing.T) {
	got, err := RenderHelp([]Action{ExecutableInvocation("run a command in an active worktree by its id", KnownArgument("worktree"), KnownArgument("exec"), FutureInput("target"), KnownArgument("--"), FutureInput("command"))})
	if err != nil {
		t.Fatal(err)
	}
	const want = "help[1]{cmd,why}:\n  bench 'worktree' 'exec' <target> '--' <command>,run a command in an active worktree by its id\n"
	if got != want {
		t.Fatalf("RenderHelp = %q, want %q", got, want)
	}
	argv, err := axitest.HelpAction{Cmd: "bench 'worktree' 'exec'"}.Argv()
	if wantArgv := []string{"bench", "worktree", "exec"}; err != nil || !slices.Equal(argv, wantArgv) {
		t.Fatalf("known prefix argv = (%q, %v), want %q", argv, err, wantArgv)
	}
}

// TestRenderHelpRejectsQuotedCommandSpoof is S10: a first known value of bench stays refused
// although its rendered word now carries quotes, so the check reads the declared value. A
// first value that only starts with bench is a different word and still renders.
func TestRenderHelpRejectsQuotedCommandSpoof(t *testing.T) {
	spoof := ExecutableInvocation("spoof", KnownArgument("bench"), KnownArgument("status"))
	if got, err := RenderHelp([]Action{spoof}); err == nil {
		t.Fatalf("RenderHelp accepted a bench first word: %q", got)
	}
	if line, err := spoof.Command(); err == nil {
		t.Fatalf("Command accepted a bench first word: %q", line)
	}
	got, err := RenderHelp([]Action{ExecutableInvocation("near", KnownArgument("benchmark"))})
	if want := "help[1]{cmd,why}:\n  bench 'benchmark',near\n"; err != nil || got != want {
		t.Fatalf("RenderHelp = (%q, %v), want %q", got, err, want)
	}
}
