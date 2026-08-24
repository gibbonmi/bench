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
	if err := os.MkdirAll(filepath.Join(root, DecisionsDir), 0o755); err != nil {
		t.Fatal(err)
	}
	const frontier = `# Alpha

Status: shaping

## Destination

Settle it.

## #1: First

Blocked by: none
Type: Research

### Question

What first?

### Answer

— (open)

## #2: Second

Blocked by: none
Type: Task

### Question

What second?

### Answer

— (open)

## Not yet specified

## Spec-writer discretion

## Out of scope

## Sources
`
	if err := os.WriteFile(filepath.Join(root, DecisionsDir, "alpha.md"), []byte(frontier), 0o644); err != nil {
		t.Fatal(err)
	}
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
				if err := os.MkdirAll(filepath.Join(root, DecisionsDir), 0o755); err != nil {
					t.Fatal(err)
				}
				document := strings.NewReplacer("Status: shaping", "Status: ready", "<answer>", "Resolved.").Replace(DecisionMapTemplate())
				if err := os.WriteFile(filepath.Join(root, DecisionsDir, "complete.md"), []byte(document), 0o644); err != nil {
					t.Fatal(err)
				}
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
	if err := os.MkdirAll(filepath.Join(root, DecisionsDir), 0o755); err != nil {
		t.Fatal(err)
	}
	document := `# Model

Status: shaping

## Destination

Settle it.

## #1: First

Blocked by: none
Type: Research

### Question

What?

### Answer

— (open)

## #2: Second

Blocked by: #1
Type: Task

### Question

What next?

### Answer

— (open)

## #3: Third

Blocked by: #1
Type: Prototype

### Question

What later?

### Answer

— (deferred)

## Not yet specified

- Honest fog.

## Spec-writer discretion

## Out of scope

## Sources
`
	if err := os.WriteFile(filepath.Join(root, DecisionsDir, "model.md"), []byte(document), 0o644); err != nil {
		t.Fatal(err)
	}
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
	if err := os.MkdirAll(filepath.Join(root, DecisionsDir), 0o755); err != nil {
		t.Fatal(err)
	}
	document := `# Model

Status: shaping

## Destination

Settle it.

## #1: Settled

Blocked by: none
Type: Grill

### Question

What?

### Answer

Resolved.

## Not yet specified

- Honest fog.

## Spec-writer discretion

## Out of scope

## Sources
`
	if err := os.WriteFile(filepath.Join(root, DecisionsDir, "model.md"), []byte(document), 0o644); err != nil {
		t.Fatal(err)
	}
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

func TestCommandRejectsCountAndTemplateTogether(t *testing.T) {
	out, code := Command([]string{"--count", "--template"})
	if code != 2 || !strings.Contains(out, "--count and --template are mutually exclusive") {
		t.Fatalf("Command(--count --template) = (%q, %d), want usage exit 2", out, code)
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
	if err := os.MkdirAll(filepath.Join(root, DecisionsDir), 0o755); err != nil {
		t.Fatal(err)
	}
	document := strings.Replace(DecisionMapTemplate(), "<answer>", "Resolved.", 1)
	if err := os.WriteFile(filepath.Join(root, DecisionsDir, "silent.md"), []byte(document), 0o644); err != nil {
		t.Fatal(err)
	}
	rows, count, state := ActiveRows(root)
	want := [][]any{{"silent", "Not yet specified", "fog", "shaping", ""}}
	if state != bounds.StateParsed || !reflect.DeepEqual(rows, want) || count != 1 {
		t.Fatalf("ActiveRows = (%#v, %d, %s), want (%#v, 1, parsed)", rows, count, state, want)
	}
}

func TestActiveCountsSeparatesReadyMaps(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, DecisionsDir), 0o755); err != nil {
		t.Fatal(err)
	}
	shaping := strings.Replace(DecisionMapTemplate(), "<answer>", "Resolved.", 1)
	ready := strings.Replace(shaping, "Status: shaping", "Status: ready", 1)
	for name, document := range map[string]string{"shaping.md": shaping, "ready.md": ready} {
		if err := os.WriteFile(filepath.Join(root, DecisionsDir, name), []byte(document), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	unresolved, readyCount, state := ActiveCounts(root)
	if unresolved != 1 || readyCount != 1 || state != bounds.StateParsed {
		t.Fatalf("ActiveCounts = (%d, %d, %s), want (1, 1, parsed)", unresolved, readyCount, state)
	}
}
