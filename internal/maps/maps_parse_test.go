// Tests for map-index parsing, schema, candidate discovery, and tree validation.
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
	got, diagnostics := parseDecisionMap("split", []byte(splitIndex))
	if len(diagnostics) != 0 {
		t.Fatalf("parseDecisionMap diagnostics = %v", diagnostics)
	}
	if got.Title != "Split map" || got.Status != "shaping" {
		t.Fatalf("parseDecisionMap = %+v", got)
	}
	template := DecisionMapTemplate()
	for _, required := range []string{"# ", "Status: shaping", "## Destination", "## Notes", "## Decisions so far", "## Not yet specified", "## Spec-writer discretion", "## Out of scope", "## Sources"} {
		if !strings.Contains(template, required) {
			t.Errorf("index template missing schema token %q", required)
		}
	}
	if strings.Contains(template, "|") {
		t.Fatalf("template status is not paste-ready: %q", template)
	}
	if _, diagnostics := parseDecisionMap("template", []byte(template)); len(diagnostics) != 0 {
		t.Fatalf("template diagnostics = %v", diagnostics)
	}
	if _, diagnostics := parseDecisionMap("template", []byte(strings.Replace(template, "Status: shaping", "Status: ready", 1))); len(diagnostics) != 0 {
		t.Fatalf("ready status diagnostics = %v", diagnostics)
	}
	assertTemplateSourcesExample(t)
	assertTicketTemplateShape(t)
	assertTemplateValidatesClean(t)
}

// assertTemplateSourcesExample holds rows SAD1, SAD2, and SAD3. The Sources body
// teaches the record grammar, so the bullet count, the locator kind, and the
// two-physical-line shape each get an independent assertion.
func assertTemplateSourcesExample(t *testing.T) {
	t.Helper()
	template := DecisionMapTemplate()
	start := strings.Index(template, "## Sources")
	if start < 0 {
		t.Fatalf("template = %q, want a Sources heading", template)
	}
	lines := strings.Split(template[start:], "\n")
	bullets := 0
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "- ") {
			bullets++
		}
	}
	if bullets != 1 {
		t.Errorf("Sources bullets = %d, want 1", bullets)
	}
	if !strings.Contains(template[start:], "\n- URL: ") || strings.Contains(template[start:], "- Path: ") {
		t.Errorf("Sources example = %q, want a URL locator", template[start:])
	}
	supports := -1
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "Supports: ") {
			supports = i
		}
	}
	if supports < 0 || supports+1 >= len(lines) || !strings.HasPrefix(strings.TrimSpace(lines[supports+1]), "Drift: ") {
		t.Errorf("Sources record = %q, want Supports and Drift on two physical lines", template[start:])
	}
}

// assertTicketTemplateShape holds rows DS17, DS18, and DS20. The index skeleton
// renders no decision ticket, and the ticket skeleton carries the whole ticket.
func assertTicketTemplateShape(t *testing.T) {
	t.Helper()
	index := DecisionMapTemplate()
	if strings.Contains(index, "## #") {
		t.Errorf("index template = %q, want no inline decision ticket", index)
	}
	if !strings.Contains(index, "A map-owned asset stays in the map's assets folder, decisions/<topic>/assets/.") {
		t.Errorf("index template = %q, want the asset rule", index)
	}
	ticket := DecisionTicketTemplate()
	if !strings.HasPrefix(ticket, "# ") {
		t.Errorf("ticket template = %q, want a title line", ticket)
	}
	for _, required := range []string{"Blocked by: none", "Type: Research", "### Question", "### Answer", resolvedBlockedRule} {
		if !strings.Contains(ticket, required) {
			t.Errorf("ticket template missing %q", required)
		}
	}
}

// assertTemplateValidatesClean holds row SAD4 at the maps seam. The temporary root
// holds no repository files, so a Path locator would red here.
func assertTemplateValidatesClean(t *testing.T) {
	t.Helper()
	for _, status := range []string{"shaping", "ready"} {
		root := t.TempDir()
		index := strings.Replace(DecisionMapTemplate(), "Status: shaping", "Status: "+status, 1)
		writeSplitMap(t, root, "decisions/template.md", index, map[string]string{"1.md": DecisionTicketTemplate()})
		if diagnostics := ValidateDecisionMapTree(root); len(diagnostics) != 0 {
			t.Errorf("%s template diagnostics = %v", status, diagnostics)
		}
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
	if _, diagnostics := parseDecisionMap("template", []byte(template)); len(diagnostics) != 0 {
		t.Fatalf("parser drifted from schema syntax: %v", diagnostics)
	}
}

func TestParseDecisionMapRequiredShapeDiagnostics(t *testing.T) {
	cases := []struct {
		name, document, want string
	}{
		{"missing title", strings.TrimPrefix(splitIndex, "# Split map\n\n"), "missing title"},
		{"missing status", strings.Replace(splitIndex, "Status: shaping\n\n", "", 1), "missing Status"},
		{"unsupported status", strings.Replace(splitIndex, "Status: shaping", "Status: done", 1), `unsupported Status "done"`},
		{"missing destination", strings.Replace(splitIndex, "Ship one canonical parser.", "", 1), "missing Destination"},
		{"missing terminal section", strings.Replace(splitIndex, "## Sources\n", "", 1), "missing Sources section"},
		{"duplicate terminal section", splitIndex + "\n## Sources\n", "duplicate Sources section"},
		{"duplicate title", splitIndex + "\n# Another title\n", "duplicate title"},
		{"duplicate status", splitIndex + "\nStatus: ready\n", "duplicate Status"},
		{"duplicate destination", splitIndex + "\n## Destination\n\nAnother destination.\n", "duplicate Destination section"},
		{"duplicate index section", splitIndex + "\n## Notes\n", "duplicate Notes section"},
		{"unsupported handoff", splitIndex + "\n## Handoff\n", "unsupported Handoff section"},
		{"inline ticket", strings.Replace(splitIndex, "## Notes\n", "## #2: Old\n\nBlocked by: none\n\n", 1), "inline ticket #2: move it to split/tickets/2.md"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, diagnostics := parseDecisionMap("split", []byte(c.document))
			for _, diagnostic := range diagnostics {
				if diagnostic.Message == c.want {
					return
				}
			}
			t.Fatalf("diagnostics = %v, want %q", diagnostics, c.want)
		})
	}
}

func TestParseDecisionMapMarkdownEdges(t *testing.T) {
	const fenced = "# Map\r\n\r\nStatus: shaping\r\n\r\n## Destination\r\n\r\nParse maps.\r\n\r\n```md\r\n## Handoff\r\n## #1: Inline\r\n```\r\n\r\n## Notes\r\n\r\n## Decisions so far\r\n\r\n## Not yet specified\r\n\r\n## Spec-writer discretion\r\n\r\n## Out of scope\r\n\r\n## Sources"
	if _, diagnostics := parseDecisionMap("map", []byte(fenced)); len(diagnostics) != 0 {
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
		"decisions/active/assets/nested.md",
		"specs/compiled/decisions/compiled.md",
		"specs/compiled/decisions/.hidden.md",
		"specs/compiled/decisions/README.MD",
		"specs/compiled/decisions/uppercase.MD",
		"specs/compiled/decisions/compiled/assets/nested.md",
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
	active := DecisionMapTemplate()
	compiled := strings.Replace(active, "Status: shaping", "Status: ready", 1)
	writeSplitMap(t, root, "decisions/active.md", active, map[string]string{"1.md": DecisionTicketTemplate()})
	writeSplitMap(t, root, "specs/compiled/decisions/compiled.md", compiled, map[string]string{"1.md": DecisionTicketTemplate()})
	writeSplitMap(t, root, "specs/no-map/spec.md", "# No map\n", nil)

	if diagnostics := ValidateDecisionMapTree(root); len(diagnostics) != 0 {
		t.Fatalf("ValidateDecisionMapTree diagnostics = %v", diagnostics)
	}

	writeSplitMap(t, root, "specs/broken/decisions/broken.md", "# Broken\n", nil)
	writeSplitMap(t, root, "decisions/graph.md", active, map[string]string{
		"1.md": strings.Replace(DecisionTicketTemplate(), "Blocked by: none", "Blocked by: #1", 1),
	})
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

// TestDecisionMapDiagnosticsGolden pins the exact ordered diagnostic slice the
// decision-map grammar delivers. The lift of the shared field scan and graph
// walk must not reorder, drop, or reword one message.
func TestDecisionMapDiagnosticsGolden(t *testing.T) {
	const duplicated = `# Map

Status: shaping

## Destination

Ship it.

## Notes

## Decisions so far

## #1: First

Blocked by: none

## Not yet specified

## Spec-writer discretion

## Out of scope

## Sources

## Sources

# Second title

Status: ready
`
	const fenced = "# Map\n\nStatus: shaping\n\n## Destination\n\nShip it.\n\n```md\nStatus: ready\n## Handoff\n## Sources\n```\n\n## Notes\n\n## Decisions so far\n\n## Not yet specified\n\n## Spec-writer discretion\n\n## Out of scope\n\n## Sources\n"
	const bare = `# Map

## Notes
`
	for _, c := range []struct {
		name     string
		document string
		want     []string
	}{
		{"inline ticket and duplicate sections", duplicated, []string{
			"inline ticket #1: move it to map/tickets/1.md",
			"duplicate Sources section",
			"duplicate title",
			"duplicate Status",
		}},
		{"fenced lines grade nothing", fenced, nil},
		{"bare skeleton", bare, []string{
			"missing Status",
			"missing Destination",
			"missing Not yet specified section",
			"missing Spec-writer discretion section",
			"missing Out of scope section",
			"missing Sources section",
		}},
	} {
		t.Run(c.name, func(t *testing.T) {
			_, diagnostics := parseDecisionMap("map", []byte(c.document))
			if got := diagnosticMessages(diagnostics); !reflect.DeepEqual(got, c.want) {
				t.Fatalf("parseDecisionMap diagnostics =\n%#v\nwant\n%#v", got, c.want)
			}
		})
	}

	root := t.TempDir()
	gist := "## Decisions so far\n\n- [First](graph/tickets/1.md): Resolved.\n"
	ticket := func(blockedBy, answer string) string {
		return "# T\n\nBlocked by: " + blockedBy + "\nType: Grill\n\n### Question\n\nQ?\n\n### Answer\n\n" + answer + "\n"
	}
	writeSplitMap(t, root, "decisions/graph.md", strings.Replace(splitIndex, "## Decisions so far\n", gist, 1), map[string]string{
		"1.md": ticket("#2, #2", "Resolved."),
		"2.md": ticket("#2", "— (open)"),
		"3.md": ticket("#4", "— (open)"),
		"4.md": ticket("#3, #9", "— (open)"),
	})
	want := []string{
		"decisions/graph.md: resolved ticket #1: T depends on unresolved #2: T",
		"decisions/graph.md: ticket #1: T duplicate blocker #2",
		"decisions/graph.md: ticket #2: T self-edge #2 -> #2",
		"decisions/graph.md: ticket #4: T dangling blocker #9",
		"decisions/graph.md: cycle edge ticket #4: T -> ticket #3: T (#4 -> #3)",
	}
	if got := ValidateDecisionMapTree(root); !reflect.DeepEqual(got, want) {
		t.Fatalf("ValidateDecisionMapTree diagnostics =\n%#v\nwant\n%#v", got, want)
	}
}

func diagnosticMessages(diagnostics []Diagnostic) []string {
	var messages []string
	for _, diagnostic := range diagnostics {
		messages = append(messages, diagnostic.Message)
	}
	return messages
}

// liftedTable is a neutral field table: the lifted scan carries no decision-map vocabulary.
var liftedTable = []FieldSpec{
	{Name: "Owner", Syntax: "Owner: "},
	{Name: "Note", Syntax: "Note: ", Scoped: true},
}

func liftedScan() FieldScan {
	return FieldScan{
		Table: liftedTable,
		Scope: func(line string) (string, bool) {
			if strings.HasPrefix(line, "## ") {
				return strings.TrimPrefix(line, "## "), true
			}
			return "", false
		},
		Duplicate: func(spec FieldSpec, scope string) string {
			if spec.Scoped {
				return "section " + scope + ": duplicate " + spec.Name
			}
			return "duplicate " + spec.Name
		},
	}
}

func TestLiftedFieldScanReportsDuplicateField(t *testing.T) {
	const document = `Owner: first
Owner: second

## one

Note: a
Note: b

## two

Note: c
`
	lines, diagnostics := liftedScan().Scan([]byte(document))
	want := []string{"duplicate Owner", "section one: duplicate Note"}
	if !reflect.DeepEqual(diagnostics, want) {
		t.Fatalf("Scan diagnostics = %#v, want %#v", diagnostics, want)
	}
	if lines[0].Value != "first" || lines[1].Diagnostic != "duplicate Owner" {
		t.Fatalf("scanned owner lines = %#v", lines[:2])
	}
}

func TestLiftedFieldScanSkipsFencedLines(t *testing.T) {
	const document = "Owner: real\n\n```md\nOwner: quoted\nNote: quoted\n```\n\nOwner: second\n"
	lines, diagnostics := liftedScan().Scan([]byte(document))
	want := []string{"duplicate Owner"}
	if !reflect.DeepEqual(diagnostics, want) {
		t.Fatalf("Scan diagnostics = %#v, want %#v", diagnostics, want)
	}
	for _, line := range lines {
		if strings.Contains(line.Text, "quoted") && (!line.Fenced || line.Field != "") {
			t.Fatalf("fenced line parsed as a field: %#v", line)
		}
	}
}
