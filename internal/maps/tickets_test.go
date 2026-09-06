// Tests for the split decision-map shape: an index file beside its ticket files.
package maps

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const splitIndex = `# Split map

Status: shaping

## Destination

Ship one canonical parser.

## Notes

## Decisions so far

## Not yet specified

- None.

## Spec-writer discretion

- None.

## Out of scope

- None.

## Sources
`

const splitTicket = `# Which parser owns the map?

Blocked by: none
Type: Research

### Question

Which parser owns the map?

### Answer

— (open)
`

func hasMessage(diagnostics []string, want string) bool {
	for _, diagnostic := range diagnostics {
		if diagnostic == want {
			return true
		}
	}
	return false
}

// writeSplitMap writes one index and its ticket files, and returns the index's
// repository-relative path.
func writeSplitMap(t *testing.T, root, rel, index string, tickets map[string]string) string {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(index), 0o644); err != nil {
		t.Fatal(err)
	}
	folder := filepath.Join(strings.TrimSuffix(path, ".md"), "tickets")
	if tickets != nil {
		if err := os.MkdirAll(folder, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	for name, body := range tickets {
		if err := os.WriteFile(filepath.Join(folder, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return rel
}

func TestValidateDecisionMapTreeAcceptsTheSplitShape(t *testing.T) {
	root := t.TempDir()
	writeSplitMap(t, root, "decisions/split.md", splitIndex, map[string]string{"1.md": splitTicket})
	if diagnostics := ValidateDecisionMapTree(root); len(diagnostics) != 0 {
		t.Fatalf("split map diagnostics = %v", diagnostics)
	}
}

func TestValidateDecisionMapTreeReadsTheBasenameAsTheTicketID(t *testing.T) {
	root := t.TempDir()
	writeSplitMap(t, root, "decisions/split.md", splitIndex, map[string]string{
		"7.md": splitTicket,
		"8.md": strings.Replace(splitTicket, "Blocked by: none", "Blocked by: #7", 1),
	})
	if diagnostics := ValidateDecisionMapTree(root); len(diagnostics) != 0 {
		t.Fatalf("basename-keyed blocker diagnostics = %v", diagnostics)
	}
}

func TestValidateDecisionMapTreeGradesTicketFileFields(t *testing.T) {
	for _, testCase := range []struct{ name, ticket, want string }{
		{"missing Type", strings.Replace(splitTicket, "Type: Research\n", "", 1), "ticket #1: missing Type"},
		{"unsupported Type", strings.Replace(splitTicket, "Type: Research", "Type: Guess", 1), `ticket #1: unsupported Type "Guess"`},
		{"missing Blocked by", strings.Replace(splitTicket, "Blocked by: none\n", "", 1), "ticket #1: missing Blocked by"},
		{"missing Question", strings.Replace(splitTicket, "### Question\n", "", 1), "ticket #1: missing Question"},
		{"missing Answer", strings.Replace(splitTicket, "### Answer\n", "", 1), "ticket #1: missing Answer"},
		{"duplicate Type", strings.Replace(splitTicket, "Type: Research", "Type: Research\nType: Task", 1), "ticket #1: duplicate Type"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			root := t.TempDir()
			writeSplitMap(t, root, "decisions/split.md", splitIndex, map[string]string{"1.md": testCase.ticket})
			if diagnostics := ValidateDecisionMapTree(root); !hasMessage(diagnostics, "decisions/split.md: "+testCase.want) {
				t.Fatalf("ticket-field diagnostics = %v, want %q", diagnostics, testCase.want)
			}
		})
	}
}

func TestValidateDecisionMapTreeRequiresTheIndexSectionsOfASplitMap(t *testing.T) {
	for _, testCase := range []struct{ old, want string }{
		{"## Notes\n", "missing Notes section"},
		{"## Decisions so far\n", "missing Decisions so far section"},
	} {
		t.Run(testCase.want, func(t *testing.T) {
			root := t.TempDir()
			index := strings.Replace(splitIndex, testCase.old, "", 1)
			writeSplitMap(t, root, "decisions/split.md", index, map[string]string{"1.md": splitTicket})
			if diagnostics := ValidateDecisionMapTree(root); !hasMessage(diagnostics, "decisions/split.md: "+testCase.want) {
				t.Fatalf("index-section diagnostics = %v, want %q", diagnostics, testCase.want)
			}
		})
	}
}

func TestValidateDecisionMapTreeRefusesAnInlineDecisionTicket(t *testing.T) {
	root := t.TempDir()
	inline := strings.Replace(splitIndex, "## Not yet specified\n", "## #2: Old\n\nBlocked by: none\nType: Research\n\n### Question\n\nOld?\n\n### Answer\n\nOld.\n\n## Not yet specified\n", 1)
	writeSplitMap(t, root, "decisions/alpha.md", inline, map[string]string{"1.md": splitTicket})
	want := "decisions/alpha.md: inline ticket #2: move it to alpha/tickets/2.md"
	if diagnostics := ValidateDecisionMapTree(root); !hasMessage(diagnostics, want) {
		t.Fatalf("inline-heading diagnostics = %v, want %q", diagnostics, want)
	}
}

func TestValidateDecisionMapTreeGradesTicketTypeAndBlockerValues(t *testing.T) {
	for _, typ := range []string{"Research", "Prototype", "Grill", "Task"} {
		root := t.TempDir()
		writeSplitMap(t, root, "decisions/split.md", splitIndex, map[string]string{
			"1.md": strings.Replace(splitTicket, "Type: Research", "Type: "+typ, 1),
		})
		if diagnostics := ValidateDecisionMapTree(root); len(diagnostics) != 0 {
			t.Fatalf("type %q diagnostics = %v", typ, diagnostics)
		}
	}
	root := t.TempDir()
	writeSplitMap(t, root, "decisions/split.md", splitIndex, map[string]string{
		"1.md": strings.Replace(splitTicket, "Blocked by: none", "Blocked by: #0", 1),
	})
	if diagnostics := ValidateDecisionMapTree(root); !hasMessage(diagnostics, "decisions/split.md: ticket #1: malformed Blocked by") {
		t.Fatalf("malformed blocker diagnostics = %v", diagnostics)
	}
}

func TestValidateDecisionMapTreeGradesGistDrift(t *testing.T) {
	resolved := strings.Replace(splitTicket, "— (open)", "The maps package.", 1)
	gist := "## Decisions so far\n\n- [Which parser](split/tickets/1.md): The maps package.\n"
	for _, testCase := range []struct{ name, index, ticket, want string }{
		{"resolved without a gist", splitIndex, resolved, "ticket #1: resolved without a gist in Decisions so far"},
		{"unresolved gist", strings.Replace(splitIndex, "## Decisions so far\n", gist, 1), splitTicket, "Decisions so far links unresolved ticket #1"},
		{"missing ticket file", strings.Replace(splitIndex, "## Decisions so far\n", strings.Replace(gist, "/1.md", "/9.md", 1), 1), resolved, "Decisions so far links missing ticket #9"},
		{"malformed gist", strings.Replace(splitIndex, "## Decisions so far\n", "## Decisions so far\n\n- The maps package owns it.\n", 1), resolved, "Decisions so far line has no ticket link"},
		{"clean gist", strings.Replace(splitIndex, "## Decisions so far\n", gist, 1), resolved, ""},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			root := t.TempDir()
			writeSplitMap(t, root, "decisions/split.md", testCase.index, map[string]string{"1.md": testCase.ticket})
			diagnostics := ValidateDecisionMapTree(root)
			if testCase.want == "" {
				if len(diagnostics) != 0 {
					t.Fatalf("clean gist diagnostics = %v", diagnostics)
				}
				return
			}
			if !hasMessage(diagnostics, "decisions/split.md: "+testCase.want) {
				t.Fatalf("gist-drift diagnostics = %v, want %q", diagnostics, testCase.want)
			}
		})
	}
}

func TestValidateDecisionMapTreeGradesTicketFolderShape(t *testing.T) {
	t.Run("empty folder", func(t *testing.T) {
		root := t.TempDir()
		writeSplitMap(t, root, "decisions/split.md", splitIndex, map[string]string{})
		if diagnostics := ValidateDecisionMapTree(root); !hasMessage(diagnostics, "decisions/split.md: missing decision ticket") {
			t.Fatalf("empty tickets folder diagnostics = %v", diagnostics)
		}
	})
	t.Run("basename", func(t *testing.T) {
		for _, name := range []string{"07.md", "0.md", "a.md"} {
			root := t.TempDir()
			writeSplitMap(t, root, "decisions/split.md", splitIndex, map[string]string{name: splitTicket, "1.md": splitTicket})
			want := "decisions/split.md: tickets/" + name + ": ticket file name must be a decision ticket number"
			if diagnostics := ValidateDecisionMapTree(root); !hasMessage(diagnostics, want) {
				t.Fatalf("basename diagnostics = %v, want %q", diagnostics, want)
			}
		}
	})
	t.Run("orphan folder", func(t *testing.T) {
		root := t.TempDir()
		folder := filepath.Join(root, "decisions", "beta", "tickets")
		if err := os.MkdirAll(folder, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(folder, "1.md"), []byte(splitTicket), 0o644); err != nil {
			t.Fatal(err)
		}
		if diagnostics := ValidateDecisionMapTree(root); !hasMessage(diagnostics, "decisions/beta/tickets: no map index at beta.md") {
			t.Fatalf("orphan tickets folder diagnostics = %v", diagnostics)
		}
	})
}

func TestValidateDecisionMapTreeScansCompiledAndSpacedSplitMaps(t *testing.T) {
	root := t.TempDir()
	ready := strings.Replace(strings.Replace(splitIndex, "Status: shaping", "Status: ready", 1), "## Not yet specified\n\n- None.\n", "## Not yet specified\n", 1)
	readyTicket := strings.Replace(splitTicket, "— (open)", "The maps package.", 1)
	gist := "## Decisions so far\n\n- [Which parser](compiled/tickets/1.md): The maps package.\n"
	writeSplitMap(t, root, "specs/x/decisions/compiled.md", strings.Replace(ready, "## Decisions so far\n", gist, 1), map[string]string{"1.md": readyTicket})
	writeSplitMap(t, root, "decisions/my map.md", splitIndex, map[string]string{"1.md": splitTicket})
	if diagnostics := ValidateDecisionMapTree(root); len(diagnostics) != 0 {
		t.Fatalf("compiled and spaced split map diagnostics = %v", diagnostics)
	}

	shaping := strings.Replace(strings.Replace(ready, "## Decisions so far\n", gist, 1), "Status: ready", "Status: shaping", 1)
	writeSplitMap(t, root, "specs/x/decisions/compiled.md", shaping, map[string]string{"1.md": readyTicket})
	if diagnostics := ValidateDecisionMapTree(root); !hasMessage(diagnostics, "specs/x/decisions/compiled.md: compiled map must be ready") {
		t.Fatalf("compiled shaping diagnostics = %v", diagnostics)
	}
}
