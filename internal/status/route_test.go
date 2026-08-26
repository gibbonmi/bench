package status

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/harnesses"
)

func testSignal(severity int, name, detail, text string) Signal {
	action := parseAction(text)
	if !action.invocable() && text != "" {
		action = advisoryAction(text)
	}
	return newSignal(severity, name, detail, action)
}

func TestIsInvocable(t *testing.T) {
	cases := []struct {
		action string
		want   bool
	}{
		{"bench gate --fresh for a whole-tree verdict", false},
		{"/bench-final-check / push", false},
		{"benchmark the loop", false},
		{"git push", true},
		{"bench spec retire <slug>", true},
		{"/bench-drain", true},
		{"/bench-review-implementation", true},
		{"/bench-implement-spec specs/my [draft]/spec.md", true},
		{"/bench-write-spec decisions/my * map.md", true},
	}
	for _, tc := range cases {
		if got := IsInvocable(tc.action); got != tc.want {
			t.Errorf("IsInvocable(%q) = %v, want %v", tc.action, got, tc.want)
		}
	}
}

func TestActionDefinitionsRenderAndParseTheSameCommand(t *testing.T) {
	for id := actionID(1); id < actionCount; id++ {
		definition := actionDefinitions[id]
		if definition.argument == anyPhaseCommand {
			continue
		}
		action := commandAction(id)
		switch definition.argument {
		case oneWordArgument:
			action = commandActionWithArgument(id, "<value>")
		case optionalSpecPath:
			action = commandActionWithArgument(id, "specs/my draft/spec.md")
		case optionalDecisionPath:
			action = commandActionWithArgument(id, "decisions/my map.md")
		}
		rendered := action.render()
		parsed := parseAction(rendered)
		if parsed.id != id || parsed.render() != rendered || !parsed.invocable() {
			t.Errorf("definition %d roundtrip = %#v from %q", id, parsed, rendered)
		}
	}
}

func TestRouteDoesNotReparseAnUnconstructedAction(t *testing.T) {
	signals := []Signal{{Name: "gate", Detail: "locked-pending", Action: "/bench-drain"}}

	got := Route(signals, HarnessClaude)
	if !got.NoCommand || got.Lead != signals[0] {
		t.Fatalf("Route(raw action) = %#v, want the board signal retained as non-invocable", got)
	}
}

func TestProducedDrainRoutesFromItsConstruction(t *testing.T) {
	root := initRepo(t)
	if err := os.WriteFile(filepath.Join(root, "capture/IDEAS.md"), []byte("- 2026-08-18  pending\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")

	got := RouteFor(root, SignalsWith(root, Query{}), HarnessClaude)
	if got.NoCommand || got.Lead.Name != "drain" || got.Lead.Action != "/bench-drain" {
		t.Fatalf("produced drain route = %#v, want /bench-drain as the lead command", got)
	}
}

// HC13. The choice list names every record row, with claude first and the rest sorted.
func TestHarnessChoicesNamesTheRecord(t *testing.T) {
	if got, want := HarnessChoices(), "claude|codex|none|opencode"; got != want {
		t.Errorf("HarnessChoices() = %q, want %q", got, want)
	}
	for _, name := range []string{"claude", "codex", "opencode", "none"} {
		if !ValidHarness(name) {
			t.Errorf("ValidHarness(%q) = false, want true", name)
		}
	}
	if ValidHarness("cursor") {
		t.Error(`ValidHarness("cursor") = true, want false`)
	}
}

// The table holds one entry per record row, and a formless row holds an empty form.
func TestHarnessPrefixMirrorsTheRecord(t *testing.T) {
	if len(harnessPrefix) != len(harnesses.Rows) {
		t.Fatalf("harnessPrefix has %d entries, want %d", len(harnessPrefix), len(harnesses.Rows))
	}
	for _, row := range harnesses.Rows {
		form, ok := harnessPrefix[row.Harness]
		if !ok || form != row.PhaseForm {
			t.Errorf("harnessPrefix[%q] = (%q, %v), want %q", row.Harness, form, ok, row.PhaseForm)
		}
	}
	if HarnessClaude != "claude" {
		t.Errorf("HarnessClaude = %q, want %q", HarnessClaude, "claude")
	}
}

func TestRouteTranslatesHarnessPrefix(t *testing.T) {
	signals := []Signal{testSignal(4, "drain", "1 idea", "/bench-drain")}
	for _, tc := range []struct {
		harness, want string
	}{
		{HarnessClaude, "/bench-drain"},
		{"codex", "$bench-drain"},
	} {
		if got := Route(signals, tc.harness).Lead.Action; got != tc.want {
			t.Errorf("Route(%q).Lead.Action = %q, want %q", tc.harness, got, tc.want)
		}
	}
}

func TestFirstInvocable(t *testing.T) {
	prose := []Signal{
		{Severity: 0, Name: "gate", Detail: "red", Action: "fix before commit"},
		{Severity: 1, Name: "git", Detail: "1 dirty path", Action: "commit on green"},
		{Severity: 5, Name: "structure", Detail: "2 issue(s)", Action: "split (craft-seams)"},
	}
	cases := []struct {
		what        string
		signals     []Signal
		action      string
		wantName    string
		wantPresent bool
	}{
		{what: "no board at all", signals: nil},
		{what: "empty board", signals: []Signal{}},
		{what: "every action is prose", signals: prose},
		{
			what: "prose ahead of a command",
			signals: append(append([]Signal{}, prose...),
				testSignal(4, "drain", "1 idea(s)", "/bench-drain")),
			action: "/bench-drain", wantName: "drain", wantPresent: true,
		},
		{
			what: "a command leads",
			signals: []Signal{
				testSignal(3, "guards", "pre-push hook missing", "bench link"),
				testSignal(4, "drain", "1 idea(s)", "/bench-drain"),
			},
			action: "bench link", wantName: "guards", wantPresent: true,
		},
		{
			what: "a bench subcommand carrying an argument placeholder",
			signals: []Signal{
				{Severity: 5, Name: "structure", Detail: "2 issue(s)", Action: "split (craft-seams)"},
				testSignal(8, "specs", "1 merged spec(s)", "bench spec retire <slug>"),
			},
			action: "bench spec retire <slug>", wantName: "specs", wantPresent: true,
		},
		{
			what: "a compound action naming several steps is not one invocation",
			signals: []Signal{
				{Severity: 1, Name: "git", Detail: "7 unpushed commits", Action: "/bench-final-check / push"},
				testSignal(4, "drain", "1 idea(s)", "/bench-drain"),
			},
			action: "/bench-drain", wantName: "drain", wantPresent: true,
		},
		{
			what:    "a word beginning with bench is not a command",
			signals: []Signal{{Severity: 11, Name: "handoff", Detail: "stale", Action: "benchmark the loop"}},
		},
	}
	for _, tc := range cases {
		_, signal, ok := firstInvocable(tc.signals)
		if ok != tc.wantPresent || signal.Action != tc.action || signal.Name != tc.wantName {
			t.Fatalf("%s: firstInvocable = (%q, %q, %v), want (%q, %q, %v)",
				tc.what, signal.Action, signal.Name, ok, tc.action, tc.wantName, tc.wantPresent)
		}
	}
}

func TestRouteForUsesCleanFallbackOnlyForAnEmptyBoard(t *testing.T) {
	root := t.TempDir()
	drain := RouteFor(root, nil, HarnessClaude)
	if drain.Lead != testSignal(0, "clean", "nothing pending", "/bench-drain") || drain.NoCommand {
		t.Fatalf("drain fallback = %#v", drain)
	}
	if got, code := renderRoute(drain); code != 0 || got != "next[1]{state,why,command}:\n  clean,nothing pending,/bench-drain\nalso: none\n" {
		t.Fatalf("drain fallback output = (%q, %d)", got, code)
	}
	if err := os.WriteFile(filepath.Join(root, "ROADMAP.md"), []byte("# Roadmap\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	roadmap := RouteFor(root, nil, HarnessClaude)
	if roadmap.Lead != testSignal(0, "clean", "nothing pending", "bench roadmap") || roadmap.NoCommand {
		t.Fatalf("roadmap fallback = %#v", roadmap)
	}
	if got, code := renderRoute(roadmap); code != 0 || got != "next[1]{state,why,command}:\n  clean,nothing pending,bench roadmap\nalso: none\n" {
		t.Fatalf("roadmap fallback output = (%q, %d)", got, code)
	}

	locked := RouteFor(root, []Signal{{Name: "gate", Detail: "locked-pending"}}, HarnessClaude)
	if locked.Lead != (Signal{Name: "gate", Detail: "locked-pending"}) || !locked.NoCommand || len(locked.RunnersUp) != 0 {
		t.Fatalf("locked board = %#v, want no-command lead without fallback", locked)
	}
	if got, code := renderRoute(locked); code != 0 || got != "next[1]{state,why,command}:\n  gate,locked-pending,\"\"\nalso: none\n" {
		t.Fatalf("locked route = (%q, %d)", got, code)
	}
}

func TestRouteForTranslatesLeadAndRunnersUpForCodex(t *testing.T) {
	route := RouteFor(t.TempDir(), []Signal{
		testSignal(0, "drain", "1 idea", "/bench-drain"),
		testSignal(0, "specs", "1 staged spec", "/bench-implement-spec"),
	}, "codex")
	if got, code := renderRoute(route); code != 0 || got != "next[1]{state,why,command}:\n  drain,1 idea,$bench-drain\nalso: specs (1 staged spec) → $bench-implement-spec\n" {
		t.Fatalf("codex route = (%q, %d)", got, code)
	}
}
