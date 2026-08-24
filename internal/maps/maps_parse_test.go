// Tests for decision-map parsing, schema, candidate discovery, and tree validation.
package maps

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
)

func TestParseDecisionMapSchemaAndTemplate(t *testing.T) {
	const document = `# Schema map

Status: shaping

## Destination

Ship one canonical parser.

## #1: Which parser owns the map?

Blocked by: none
Type: Research

### Question

Which parser owns the map?

### Answer

The maps package.

## Not yet specified

- None.

## Spec-writer discretion

- None.

## Out of scope

- None.

## Sources
`
	got, diagnostics := ParseDecisionMap([]byte(document))
	if len(diagnostics) != 0 {
		t.Fatalf("ParseDecisionMap diagnostics = %v", diagnostics)
	}
	if got.Title != "Schema map" || got.Status != "shaping" {
		t.Fatalf("ParseDecisionMap = %+v", got)
	}
	for _, required := range []string{"# ", "Status: shaping", "## Destination", "## #", "Type: Research", "### Question", "### Answer", "## Not yet specified", "## Spec-writer discretion", "## Out of scope", "## Sources"} {
		if !strings.Contains(DecisionMapTemplate(), required) {
			t.Errorf("template missing schema token %q", required)
		}
	}
	if strings.Contains(DecisionMapTemplate(), "|") {
		t.Fatalf("template status is not paste-ready: %q", DecisionMapTemplate())
	}
	if _, diagnostics := ParseDecisionMap([]byte(DecisionMapTemplate())); len(diagnostics) != 0 {
		t.Fatalf("template diagnostics = %v", diagnostics)
	}
	if _, diagnostics := ParseDecisionMap([]byte(strings.Replace(DecisionMapTemplate(), "Status: shaping", "Status: ready", 1))); len(diagnostics) != 0 {
		t.Fatalf("ready status diagnostics = %v", diagnostics)
	}
}

func TestDecisionMapSchemaSyntaxDrivesParserAndTemplate(t *testing.T) {
	status := -1
	for i, field := range canonicalDecisionMapSchema.fields {
		if field.name == "Status" {
			status = i
			break
		}
	}
	if status < 0 {
		t.Fatal("schema has no Status field")
	}
	original := canonicalDecisionMapSchema.fields[status].syntax
	canonicalDecisionMapSchema.fields[status].syntax = "Phase: "
	t.Cleanup(func() { canonicalDecisionMapSchema.fields[status].syntax = original })
	template := DecisionMapTemplate()
	if !strings.Contains(template, "Phase: shaping") {
		t.Fatalf("template = %q, want schema status syntax", template)
	}
	if _, diagnostics := ParseDecisionMap([]byte(template)); len(diagnostics) != 0 {
		t.Fatalf("parser drifted from schema syntax: %v", diagnostics)
	}
}

func TestParseDecisionMapRequiredShapeDiagnostics(t *testing.T) {
	const valid = `# Schema map

Status: shaping

## Destination

Ship one canonical parser.

## #1: Which parser owns the map?

Blocked by: none
Type: Research

### Question

Which parser owns the map?

### Answer

The maps package.

## Not yet specified

- None.

## Spec-writer discretion

- None.

## Out of scope

- None.

## Sources
`
	cases := []struct {
		name, document, want string
	}{
		{"missing title", strings.TrimPrefix(valid, "# Schema map\n\n"), "missing title"},
		{"missing status", strings.Replace(valid, "Status: shaping\n\n", "", 1), "missing Status"},
		{"missing destination", strings.Replace(valid, "Ship one canonical parser.", "", 1), "missing Destination"},
		{"missing ticket", valid[:strings.Index(valid, "## #1:")] + "## Not yet specified\n\n## Spec-writer discretion\n\n## Out of scope\n\n## Sources\n", "missing decision ticket"},
		{"missing terminal section", strings.Replace(valid, "## Sources\n", "", 1), "missing Sources section"},
		{"duplicate terminal section", valid + "\n## Sources\n", "duplicate Sources section"},
		{"duplicate title", valid + "\n# Another title\n", "duplicate title"},
		{"duplicate status", valid + "\nStatus: ready\n", "duplicate Status"},
		{"duplicate destination", valid + "\n## Destination\n\nAnother destination.\n", "duplicate Destination section"},
		{"duplicate ticket field", strings.Replace(valid, "Type: Research", "Type: Research\nType: Grill", 1), "ticket #1: duplicate Type"},
		{"malformed blockers", strings.Replace(valid, "Blocked by: none", "Blocked by: #0", 1), "ticket #1: malformed Blocked by"},
		{"unsupported handoff", valid + "\n## Handoff\n", "unsupported Handoff section"},
		{"unsupported type", strings.Replace(valid, "Type: Research", "Type: Guess", 1), "ticket #1: unsupported Type \"Guess\""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, diagnostics := ParseDecisionMap([]byte(c.document))
			for _, diagnostic := range diagnostics {
				if diagnostic.Message == c.want {
					return
				}
			}
			t.Fatalf("diagnostics = %v, want %q", diagnostics, c.want)
		})
	}
}

func TestParseDecisionMapTicketTypesAndMarkdownEdges(t *testing.T) {
	for _, typ := range []string{"Research", "Prototype", "Grill", "Task"} {
		t.Run(typ, func(t *testing.T) {
			document := strings.Replace(`# Map

Status: shaping

## Destination

Parse maps.

## #1: A decision?

Blocked by: none
Type: TYPE

### Question

A decision?

### Answer

Decided.

## Not yet specified

## Spec-writer discretion

## Out of scope

## Sources
`, "TYPE", typ, 1)
			if _, diagnostics := ParseDecisionMap([]byte(document)); len(diagnostics) != 0 {
				t.Fatalf("diagnostics = %v", diagnostics)
			}
		})
	}
	const fenced = "# Map\r\n\r\nStatus: shaping\r\n\r\n## Destination\r\n\r\nParse maps.\r\n\r\n```md\r\n## Handoff\r\nType: Guess\r\n```\r\n\r\n## #1: A decision?\r\n\r\nBlocked by: none\r\nType: Grill\r\n\r\n### Question\r\n\r\nA decision?\r\n\r\n### Answer\r\n\r\nDecided.\r\n\r\n## Not yet specified\r\n\r\n## Spec-writer discretion\r\n\r\n## Out of scope\r\n\r\n## Sources"
	if _, diagnostics := ParseDecisionMap([]byte(fenced)); len(diagnostics) != 0 {
		t.Fatalf("CRLF/fenced/no-final-newline diagnostics = %v", diagnostics)
	}
}

func TestDiscoverDecisionMapCandidatesDirectChildren(t *testing.T) {
	root := t.TempDir()
	for _, path := range []string{
		"decisions/active.md",
		"decisions/.hidden.md",
		"decisions/README.md",
		"decisions/readme.MD",
		"decisions/uppercase.MD",
		"decisions/assets/nested.md",
		"specs/compiled/decisions/compiled.md",
		"specs/compiled/decisions/.hidden.md",
		"specs/compiled/decisions/README.MD",
		"specs/compiled/decisions/uppercase.MD",
		"specs/compiled/decisions/assets/nested.md",
		"specs/no-map/notes.md",
	} {
		full := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte("# map\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	got, err := DiscoverDecisionMapCandidates(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []DecisionMapCandidate{
		{Path: "decisions/active.md"},
		{Path: "specs/compiled/decisions/compiled.md", Compiled: true},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DiscoverDecisionMapCandidates = %#v, want %#v", got, want)
	}
}

func TestValidateDecisionMapTreeValidatesActiveAndCompiledCandidates(t *testing.T) {
	root := t.TempDir()
	writeMap := func(path, document string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(document), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	active := strings.Replace(DecisionMapTemplate(), "<answer>", "Resolved.", 1)
	compiled := strings.Replace(active, "Status: shaping", "Status: ready", 1)
	writeMap(filepath.Join(root, "decisions", "active.md"), active)
	writeMap(filepath.Join(root, "specs", "compiled", "decisions", "compiled.md"), compiled)
	writeMap(filepath.Join(root, "specs", "no-map", "spec.md"), "# No map\n")

	if diagnostics := ValidateDecisionMapTree(root); len(diagnostics) != 0 {
		t.Fatalf("ValidateDecisionMapTree diagnostics = %v", diagnostics)
	}

	writeMap(filepath.Join(root, "specs", "broken", "decisions", "broken.md"), "# Broken\n")
	writeMap(filepath.Join(root, "decisions", "graph.md"), strings.Replace(active, "Blocked by: none", "Blocked by: #1", 1))
	diagnostics := ValidateDecisionMapTree(root)
	for _, want := range []string{
		"specs/broken/decisions/broken.md: missing Status",
		"decisions/graph.md: ticket #1: <decision question> self-edge #1 -> #1",
	} {
		if !hasTreeDiagnostic(diagnostics, want) {
			t.Fatalf("ValidateDecisionMapTree diagnostics = %v, want %q", diagnostics, want)
		}
	}
}

func TestActiveRowsExcludeInvalidCompiledCandidates(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "specs", "compiled", "decisions", "broken.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("# Broken\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	rows, count, state := ActiveRows(root)
	if state != bounds.StateParsed || len(rows) != 0 || count != 0 {
		t.Fatalf("ActiveRows = (%v, %d, %s), want no active rows or count for an invalid compiled candidate", rows, count, state)
	}
}

func hasTreeDiagnostic(diagnostics []string, want string) bool {
	for _, diagnostic := range diagnostics {
		if strings.Contains(diagnostic, want) {
			return true
		}
	}
	return false
}

func TestActiveScanSharesDirectCandidateDiscovery(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"active.md", ".hidden.md", "README.md", "uppercase.MD"} {
		path := filepath.Join(root, DecisionsDir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("# invalid\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	candidates, err := DiscoverDecisionMapCandidates(root)
	if err != nil {
		t.Fatal(err)
	}
	rows, count, state := ActiveRows(root)
	if state != bounds.StateParsed || count != 1 || len(rows) != 1 || candidates[0].Path != "decisions/active.md" || rows[0][0] != "active" {
		t.Fatalf("discovery and active scan diverged: candidates=%v rows=%v count=%d state=%s", candidates, rows, count, state)
	}
}

func TestTicketAnswerStateIsSharedByReadinessAndProjection(t *testing.T) {
	for _, tc := range []struct{ answer, want string }{
		{"— (open)", "frontier"}, {"— (deferred)", "deferred"}, {"GRILL DEFERRED — wait", "deferred"}, {"", "frontier"}, {"Resolved.", "resolved"},
	} {
		ticket := DecisionTicket{Answer: tc.answer}
		if got := ticketAnswerState(ticket); got != tc.want || resolved(ticket) != (tc.want == "resolved") {
			t.Errorf("answer %q = (%q, resolved=%v), want %q", tc.answer, got, resolved(ticket), tc.want)
		}
	}
}
