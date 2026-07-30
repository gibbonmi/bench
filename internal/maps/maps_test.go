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

// parseFile is the shared engine; these pin the marker/fence/handoff edges at the
// pure seam the two-derivations bug class breeds.
func TestParseFileMarkers(t *testing.T) {
	r := parseFile([]byte("## #1: a?\nType: Grill\n### Answer\n— (open)\n"))
	if len(r.tickets) != 1 || r.tickets[0] != (ticket{"1", "Grill", "open"}) {
		t.Fatalf("open ticket = %+v", r.tickets)
	}
	if !r.preHandoffMarker || !r.notCloseReady() {
		t.Errorf("open ticket file should be not-close-ready")
	}

	// A mid-line GRILL DEFERRED mention and a fenced placeholder are not markers.
	r = parseFile([]byte("## #1: a?\nType: Grill\n### Answer\nDecided: mid-line GRILL DEFERRED mention.\n\n```\n— (open)\n```\n"))
	if len(r.tickets) != 0 {
		t.Errorf("over-match: unexpected tickets %+v", r.tickets)
	}

	// No Type line → unknown.
	r = parseFile([]byte("## #1: a?\n### Answer\n— (open)\n"))
	if r.tickets[0].typ != "unknown" {
		t.Errorf("typeless ticket type = %q, want unknown", r.tickets[0].typ)
	}
}

func TestParseFileHandoff(t *testing.T) {
	cases := []struct {
		name         string
		in           string
		wantRows     [][]any
		wantNotReady bool
	}{
		{"missing handoff", "## #1: q?\nType: Grill\n### Answer\nDecided: yes.\n",
			[][]any{{"m", "handoff", "handoff", "missing"}}, true},
		{"filled handoff silent", "## #1: q?\nType: Grill\n### Answer\nDecided: yes.\n\n## Handoff\n1. a. n/a\n2. b. n/a\n",
			nil, false},
		{"placeholder in handoff", "## #1: q?\nType: Grill\n### Answer\nDecided: yes.\n\n## Handoff\n1. a.\n— (open)\n",
			[][]any{{"m", "handoff", "handoff", "open"}}, true},
		{"fenced handoff is missing", "## #1: q?\nType: Grill\n### Answer\nDecided: yes.\n\n```\n## Handoff\n1. a.\n```\n",
			[][]any{{"m", "handoff", "handoff", "missing"}}, true},
		{"open ticket no handoff row", "## #1: q?\nType: Grill\n### Answer\n— (open)\n",
			[][]any{{"m", 1, "Grill", "open"}}, true},
		{"non-map never nagged", "# Index\nprose, not a map.\n", nil, false},
	}
	for _, c := range cases {
		r := parseFile([]byte(c.in))
		if got := fileRows("m", r); !reflect.DeepEqual(got, c.wantRows) {
			t.Errorf("%s: fileRows = %v, want %v", c.name, got, c.wantRows)
		}
		if r.notCloseReady() != c.wantNotReady {
			t.Errorf("%s: notCloseReady = %v, want %v", c.name, r.notCloseReady(), c.wantNotReady)
		}
	}
}

// UnresolvedCount is DISTINCT not-close-ready files — re-homes the shell
// "AXI maps_unresolved_count distinct-file contract" and the close-readiness count
// tail, both of which used to source bin/bench-query.sh.
func TestUnresolvedCount(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "decisions")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(name, body string) {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// Two unresolved tickets in one file → 3 ticket rows but distinct-file count of the
	// two files below is 2.
	write("multi.md", "## #1: a?\nType: Grill\n### Answer\n— (open)\n\n## #2: b?\nType: Grill\n### Answer\n— (deferred)\n")
	write("solo.md", "## #1: c?\nType: Grill\n### Answer\n— (open)\n")
	// A dotfile the shell glob would never expand must stay invisible to rows and count.
	write(".hidden.md", "## #1: h?\nType: Grill\n### Answer\n— (open)\n")
	if got := len(Rows(root)); got != 3 {
		t.Errorf("Rows count = %d, want 3", got)
	}
	if got, state := UnresolvedCount(root); got != 2 || state != bounds.StateParsed {
		t.Errorf("UnresolvedCount = (%d, %s), want (2, %s)", got, state, bounds.StateParsed)
	}
}

// The close-readiness aggregate re-homes the shell contract's count tail: of the six
// files, hm/hx/hp/ho are not-close-ready, hf is ready, and README documents the
// directory rather than claiming to be a map, so the tally never sees it → 4.
func TestUnresolvedCountCloseReadiness(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "decisions")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"hm.md":     "# HM\n## #1: q?\nType: Grill\n### Answer\nDecided: yes.\n",
		"hf.md":     "# HF\n## #1: q?\nType: Grill\n### Answer\nDecided: yes.\n\n## Handoff\n1. a. n/a\n",
		"ho.md":     "# HO\n## #1: q?\nType: Grill\n### Answer\n— (open)\n",
		"hx.md":     "# HX\n## #1: q?\nType: Grill\n### Answer\nDecided: yes.\n\n```\n## Handoff\n1. a.\n```\n",
		"hp.md":     "# HP\n## #1: q?\nType: Grill\n### Answer\nDecided: yes.\n\n## Handoff\n1. a.\n— (open)\n",
		"README.md": "# Index\nnot a map.\n",
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if got, _ := UnresolvedCount(root); got != 4 {
		t.Errorf("close-readiness UnresolvedCount = %d, want 4", got)
	}
	// A file-scope marker without any ticket heading is a recognized shape (a marker),
	// not unsupported-schema: no listed row, but counted via preHandoffMarker.
	if err := os.WriteFile(filepath.Join(dir, "scope.md"), []byte("### Answer\n— (deferred)\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := len(fileRows("scope", parseFile([]byte("### Answer\n— (deferred)\n")))); got != 0 {
		t.Errorf("file-scope marker emitted %d rows, want 0", got)
	}
	if got, _ := UnresolvedCount(root); got != 5 {
		t.Errorf("UnresolvedCount with file-scope marker = %d, want 5", got)
	}
}
