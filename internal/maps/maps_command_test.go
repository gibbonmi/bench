// Tests for the maps command surface and its active-rows projection.
package maps

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/gittest"
)

func TestCommandAppendsOnlyMapActionsToTheCapturedPrimaryResponse(t *testing.T) {
	root := gittest.Repo(t)
	writeSplitMap(t, root, "decisions/alpha.md", splitIndex, map[string]string{
		"1.md": strings.Replace(splitTicket, "# Which parser owns the map?", "# First", 1),
		"2.md": strings.NewReplacer("# Which parser owns the map?", "# Second", "Type: Research", "Type: Task").Replace(splitTicket),
	})
	if err := os.WriteFile(filepath.Join(root, DecisionsDir, "broken.md"), []byte("# Broken\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	primary, err := os.ReadFile("testdata/pre-disclosure-frontier-invalid.stdout")
	if err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)

	out, code := Command(nil)
	const help = "help[3]{cmd,why}:\n  /bench-shape-idea,\"shape alpha: First\"\n  /bench-shape-idea,\"shape alpha: Second\"\n  bench maps --template,repair decisions/broken.md\n"
	if code != 1 || out != string(primary)+help {
		t.Fatalf("Command(%v) = (exit %d, %q), want captured primary plus map actions", []string(nil), code, out)
	}
}

func TestCommandAppendsHonestEmptyHelpForEmptyAndCompleteMaps(t *testing.T) {
	for _, tc := range []struct {
		name, fixture string
		write         func(t *testing.T, root string)
	}{
		{name: "empty", fixture: "pre-disclosure-terminal.stdout", write: func(t *testing.T, root string) {}},
		{
			name: "complete (aliases terminal empty)", fixture: "pre-disclosure-terminal.stdout",
			write: func(t *testing.T, root string) {
				t.Helper()
				index := strings.Replace(DecisionMapTemplate(), "Status: shaping", "Status: ready", 1)
				writeSplitMap(t, root, "decisions/complete.md", index, map[string]string{"1.md": DecisionTicketTemplate()})
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			primary, err := os.ReadFile(filepath.Join("testdata", tc.fixture))
			if err != nil {
				t.Fatal(err)
			}
			root := gittest.Repo(t)
			tc.write(t, root)
			t.Chdir(root)

			out, code := Command(nil)
			if code != 0 || out != string(primary)+"help[0]{cmd,why}:\n" {
				t.Fatalf("Command(%v) = (exit %d, %q), want terminal empty help", []string(nil), code, out)
			}
		})
	}
}

func TestActionsForRowsCarriesTheInvalidDiagnosticPath(t *testing.T) {
	row := []any{"broken", "invalid", "map", "invalid", "decisions/broken: map.md: missing Status"}
	paths := map[string]string{invalidRowKey(row): "decisions/broken: map.md"}
	help, err := axi.RenderHelp(actionsForRows([][]any{row}, paths))
	if err != nil {
		t.Fatal(err)
	}
	const want = "help[1]{cmd,why}:\n  bench maps --template,\"repair decisions/broken: map.md\"\n"
	if help != want {
		t.Fatalf("actionsForRows invalid help = %q, want %q", help, want)
	}
}

func TestCommandDisclosesTheFullPathForABoundsInvalidMap(t *testing.T) {
	root := gittest.Repo(t)
	if err := os.MkdirAll(filepath.Join(root, DecisionsDir), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, DecisionsDir, "hollow.md"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)

	out, code := Command(nil)
	const help = "help[1]{cmd,why}:\n  bench maps --template,repair decisions/hollow.md\n"
	if code != 1 || !strings.Contains(out, "hollow,invalid,map,invalid,\"empty: \"") || !strings.HasSuffix(out, help) {
		t.Fatalf("Command(%v) = (exit %d, %q), want exit 1 with the bounds diagnostic and %q", []string(nil), code, out, help)
	}
}

func TestActiveRowsProjectUnresolvedTicketsAndFog(t *testing.T) {
	root := t.TempDir()
	ticket := func(title, typ, blockedBy, answer string) string {
		return "# " + title + "\n\nBlocked by: " + blockedBy + "\nType: " + typ + "\n\n### Question\n\nWhat?\n\n### Answer\n\n" + answer + "\n"
	}
	writeSplitMap(t, root, "decisions/model.md", splitIndex, map[string]string{
		"1.md": ticket("First", "Research", "none", "— (open)"),
		"2.md": ticket("Second", "Task", "#1", "— (open)"),
		"3.md": ticket("Third", "Prototype", "#1", "— (deferred)"),
	})
	rows, count, state := ActiveRows(root)
	if state != bounds.StateParsed {
		t.Fatalf("ActiveRows state = %s, want parsed", state)
	}
	want := [][]any{
		{"model", "First", "Research", "frontier", ""},
		{"model", "Second", "Task", "blocked", "First"},
		{"model", "Third", "Prototype", "deferred", "First"},
	}
	if !reflect.DeepEqual(rows, want) {
		t.Fatalf("ActiveRows rows = %#v, want %#v", rows, want)
	}
	if count != 1 {
		t.Fatalf("ActiveRows count = %d, want 1", count)
	}
}

func TestActiveRowsProjectFogOnlyShapingMap(t *testing.T) {
	root := t.TempDir()
	gist := "## Decisions so far\n\n- [Settled](model/tickets/1.md): Resolved.\n"
	writeSplitMap(t, root, "decisions/model.md", strings.Replace(splitIndex, "## Decisions so far\n", gist, 1), map[string]string{
		"1.md": strings.NewReplacer("# Which parser owns the map?", "# Settled", "— (open)", "Resolved.").Replace(splitTicket),
	})
	rows, count, state := ActiveRows(root)
	if state != bounds.StateParsed {
		t.Fatalf("ActiveRows state = %s, want parsed", state)
	}
	want := [][]any{{"model", "Not yet specified", "fog", "shaping", ""}}
	if !reflect.DeepEqual(rows, want) {
		t.Fatalf("ActiveRows rows = %#v, want %#v", rows, want)
	}
	if count != 1 {
		t.Fatalf("ActiveRows count = %d, want 1", count)
	}
}

func TestCommandRejectsExclusiveFlagsTogether(t *testing.T) {
	for _, testCase := range []struct {
		argv []string
		want string
	}{
		{[]string{"--count", "--template"}, "--count and --template are mutually exclusive"},
		{[]string{"--template", "--ticket-template"}, "--template and --ticket-template are mutually exclusive"},
		{[]string{"--count", "--ticket-template"}, "--count and --ticket-template are mutually exclusive"},
	} {
		out, code := Command(testCase.argv)
		if code != 2 || !strings.Contains(out, testCase.want) {
			t.Fatalf("Command(%v) = (%q, %d), want usage exit 2 with %q", testCase.argv, out, code, testCase.want)
		}
	}
}

func TestCommandPrintsTheTicketTemplate(t *testing.T) {
	out, code := Command([]string{"--ticket-template"})
	if code != 0 || out != DecisionTicketTemplate() {
		t.Fatalf("Command(--ticket-template) = (%q, %d), want the ticket skeleton", out, code)
	}
}

func TestInvalidMapActionWithUnsupportedWhyUsesHonestEmptyHelp(t *testing.T) {
	actions := actionsForRows(
		[][]any{{"broken", "invalid", "map", "invalid", "diagnostic"}},
		map[string]string{invalidRowKey([]any{"broken", "invalid", "map", "invalid", "diagnostic"}): "decisions/broken\x1b.md"},
	)
	help, err := axi.RenderHelp(actions)
	if err != nil {
		t.Fatal(err)
	}
	if help != "help[0]{cmd,why}:\n" {
		t.Fatalf("RenderHelp = %q, want honest empty help", help)
	}
}

func TestActiveRowsCountsSilentShapingMap(t *testing.T) {
	root := t.TempDir()
	writeSplitMap(t, root, "decisions/silent.md", DecisionMapTemplate(), map[string]string{"1.md": DecisionTicketTemplate()})
	rows, count, state := ActiveRows(root)
	want := [][]any{{"silent", "Not yet specified", "fog", "shaping", ""}}
	if state != bounds.StateParsed || !reflect.DeepEqual(rows, want) || count != 1 {
		t.Fatalf("ActiveRows = (%#v, %d, %s), want (%#v, 1, parsed)", rows, count, state, want)
	}
}

func TestActiveCountsSeparatesReadyMaps(t *testing.T) {
	root := t.TempDir()
	shaping := DecisionMapTemplate()
	ready := strings.Replace(shaping, "Status: shaping", "Status: ready", 1)
	for name, index := range map[string]string{"shaping.md": shaping, "ready.md": ready} {
		writeSplitMap(t, root, "decisions/"+name, index, map[string]string{"1.md": DecisionTicketTemplate()})
	}

	unresolved, readyCount, state := ActiveCounts(root)
	if unresolved != 1 || readyCount != 1 || state != bounds.StateParsed {
		t.Fatalf("ActiveCounts = (%d, %d, %s), want (1, 1, parsed)", unresolved, readyCount, state)
	}
}
