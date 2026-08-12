package maps

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/bounds"
)

func TestCommandAppendsOnlyMapActionsToTheCapturedPrimaryResponse(t *testing.T) {
	root := mapsRepo(t)
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
	t.Chdir(root)

	out, code := Command(nil)
	const primary = "maps[3]{map,title,type,state,blockers}:\n  alpha,First,Research,frontier,\"\"\n  alpha,Second,Task,frontier,\"\"\n  broken,invalid,map,invalid,\"decisions/broken.md: missing Status\"\n"
	const help = "help[3]{cmd,why}:\n  /bench-shape-idea,\"shape alpha: First\"\n  /bench-shape-idea,\"shape alpha: Second\"\n  bench maps --template,repair decisions/broken.md\n"
	if code != 1 || out != primary+help {
		t.Fatalf("Command(%v) = (exit %d, %q), want captured primary plus map actions", []string(nil), code, out)
	}
}

func TestCommandAppendsHonestEmptyHelpForEmptyAndCompleteMaps(t *testing.T) {
	for _, tc := range []struct {
		name  string
		write func(t *testing.T, root string)
	}{
		{name: "empty", write: func(t *testing.T, root string) {}},
		{
			name: "complete",
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
			root := mapsRepo(t)
			tc.write(t, root)
			t.Chdir(root)

			out, code := Command(nil)
			const want = "maps[0]{map,title,type,state,blockers}:\nhelp[0]{cmd,why}:\n"
			if code != 0 || out != want {
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

func mapsRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", root).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	return root
}

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

func TestMapGraphRejectsInvalidEdges(t *testing.T) {
	const document = `# Graph

Status: shaping

## Destination

Settle the graph.

## #1: First

Blocked by: #2, #2
Type: Grill

### Question

First?

### Answer

Resolved.

## #1: Duplicate

Blocked by: #9
Type: Grill

### Question

Duplicate?

### Answer

Resolved.

## #2: Second

Blocked by: #2
Type: Grill

### Question

Second?

### Answer

— (open)

## #3: Third

Blocked by: #4
Type: Grill

### Question

Third?

### Answer

— (open)

## #4: Fourth

Blocked by: #3
Type: Grill

### Question

Fourth?

### Answer

— (open)

## Not yet specified

## Spec-writer discretion

## Out of scope

## Sources
`
	_, diagnostics := ValidateDecisionMap(t.TempDir(), "decisions/graph.md", false, []byte(document))
	for _, want := range []string{"duplicate ID #1", "duplicate blocker #2", "dangling blocker #9", "self-edge #2 -> #2", "cycle edge ticket #4: Fourth -> ticket #3: Third", "resolved ticket #1: First depends on unresolved #2: Second"} {
		if !hasDiagnostic(diagnostics, want) {
			t.Errorf("diagnostics = %v, want %q", diagnostics, want)
		}
	}
}

func TestMapReadinessAndStructuredSources(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "evidence.md"), []byte("evidence"), 0o644); err != nil {
		t.Fatal(err)
	}
	body := `# Ready

Status: STATUS

## Destination

Settle it.

## #1: Answer

Blocked by: none
Type: Grill

### Question

Answer?

### Answer

Resolved.

## Not yet specified

FOG

## Spec-writer discretion

- A bounded choice.

## Out of scope

- Not included.

## Sources

- Path: TICKevidence.mdTICK
  Supports: The settled answer.
  Drift: Update when evidence changes.
- URL: TICKhttps://example.invalid/sourceTICK
  Supports: External context.
  Drift: Update when it changes.
`
	body = strings.ReplaceAll(body, "TICK", "`")
	_, shaping := ValidateDecisionMap(root, "decisions/ready.md", false, []byte(strings.ReplaceAll(strings.ReplaceAll(body, "STATUS", "shaping"), "FOG", "- Honest fog.")))
	if len(shaping) != 0 {
		t.Fatalf("shaping diagnostics = %v", shaping)
	}
	_, ready := ValidateDecisionMap(root, "decisions/ready.md", false, []byte(strings.ReplaceAll(strings.ReplaceAll(body, "STATUS", "ready"), "FOG", "- Honest fog.")))
	if !hasDiagnostic(ready, "ready map has non-empty Not yet specified") {
		t.Fatalf("ready diagnostics = %v", ready)
	}
	for _, marker := range []string{"— (open)", "— (deferred)", "GRILL DEFERRED"} {
		unresolved := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(body, "STATUS", "ready"), "FOG", ""), "Resolved.", marker)
		_, diagnostics := ValidateDecisionMap(root, "decisions/ready.md", false, []byte(unresolved))
		if !hasDiagnostic(diagnostics, "ready map has unresolved ticket") {
			t.Errorf("ready marker %q diagnostics = %v", marker, diagnostics)
		}
	}
	emptySourceBody := body[:strings.Index(body, "## Sources")+len("## Sources\n")]
	emptySourceBody = strings.ReplaceAll(strings.ReplaceAll(emptySourceBody, "STATUS", "ready"), "FOG", "")
	_, emptySources := ValidateDecisionMap(root, "decisions/ready.md", false, []byte(emptySourceBody))
	if len(emptySources) != 0 {
		t.Fatalf("empty sources diagnostics = %v", emptySources)
	}
	_, sentinel := ValidateDecisionMap(root, "decisions/ready.md", false, []byte(emptySourceBody+"n/a\n"))
	if !hasDiagnostic(sentinel, "Sources entry") {
		t.Fatalf("sentinel diagnostics = %v", sentinel)
	}
}

func TestMapSourcesAndTerminalListsRejectHostileShapes(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "valid file.md"), []byte("ok"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, "directory"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("valid file.md", filepath.Join(root, "inside-link")); err != nil {
		t.Fatal(err)
	}
	rootLink := filepath.Join(t.TempDir(), "repository")
	if err := os.Symlink(root, rootLink); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("/dev/null", filepath.Join(root, "outside-link")); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(filepath.Join(root, "fifo"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, source := range []string{"valid file.md", "inside-link"} {
		if got := validateSourcePath(root, source); got != "" {
			t.Errorf("validateSourcePath(%q) = %q", source, got)
		}
	}
	if got := validateSourcePath(rootLink, "inside-link"); got != "" {
		t.Errorf("validateSourcePath through root symlink = %q", got)
	}
	for _, source := range []string{"", "/tmp/no", "missing", "directory", "outside-link", "fifo", "../escape"} {
		if got := validateSourcePath(root, source); got == "" {
			t.Errorf("validateSourcePath(%q) accepted hostile target", source)
		}
	}
	for _, raw := range []string{"https://example.invalid/x", "http://example.invalid"} {
		if !validSourceURL(raw) {
			t.Errorf("validSourceURL(%q) = false", raw)
		}
	}
	for _, raw := range []string{"", "example.invalid", "ftp://example.invalid", "https:///no-host", "://broken"} {
		if validSourceURL(raw) {
			t.Errorf("validSourceURL(%q) = true", raw)
		}
	}
	if markdownBullets("prose") || !markdownBullets("- one\n- two") {
		t.Fatal("Markdown list classification drifted")
	}
}

func TestMapTerminalContinuationAndEmptyAnswer(t *testing.T) {
	root := t.TempDir()
	document := func(status, answer, fog, discretion, sources string) string {
		return "# Map\n\nStatus: " + status + "\n\n## Destination\n\nSettle it.\n\n## #1: Decision\n\nBlocked by: none\nType: Grill\n\n### Question\n\nQuestion?\n\n### Answer\n" + answer + "\n## Not yet specified\n" + fog + "\n## Spec-writer discretion\n" + discretion + "\n## Out of scope\n\n- Excluded.\n\n## Sources\n" + sources
	}
	wrapped := document("shaping", "\nResolved.\n\n", "\n- A fog item\n  that continues on the next line.\n\n", "\n- A bounded choice.\n\n", "")
	if _, diagnostics := ValidateDecisionMap(root, "decisions/map.md", false, []byte(wrapped)); len(diagnostics) != 0 {
		t.Fatalf("wrapped list diagnostics = %v", diagnostics)
	}
	prose := strings.Replace(wrapped, "- A bounded choice.", "A bounded choice.", 1)
	if _, diagnostics := ValidateDecisionMap(root, "decisions/map.md", false, []byte(prose)); !hasDiagnostic(diagnostics, "Spec-writer discretion must be a Markdown bullet list") {
		t.Fatalf("prose diagnostics = %v", diagnostics)
	}
	empty := document("shaping", "\n", "", "", "")
	if _, diagnostics := ValidateDecisionMap(root, "decisions/map.md", false, []byte(empty)); len(diagnostics) != 0 {
		t.Fatalf("shaping empty-answer diagnostics = %v", diagnostics)
	}
	ready := strings.Replace(empty, "Status: shaping", "Status: ready", 1)
	if _, diagnostics := ValidateDecisionMap(root, "decisions/map.md", false, []byte(ready)); !hasDiagnostic(diagnostics, "ready map has unresolved ticket") {
		t.Fatalf("ready empty-answer diagnostics = %v", diagnostics)
	}
}

func TestMapSourcesCollectIndependentRecordFailures(t *testing.T) {
	root := t.TempDir()
	body := "- Path: `missing-one`\n  Supports: first support.\n- URL: `https://example.invalid/two`\n  Drift: second drift.\n"
	diagnostics := sourceDiagnostics(root, body)
	for _, want := range []string{"missing-one", "https://example.invalid/two"} {
		if !hasDiagnostic(diagnostics, want) {
			t.Errorf("diagnostics = %v, want source identity %q", diagnostics, want)
		}
	}
	if !validSourceURL(sourceLocator("`https://example.invalid`")) || validSourceURL(sourceLocator("`https://example.invalid")) {
		t.Fatal("source locator backtick handling drifted")
	}
}

func TestMapSourcesRequireExactRecordShape(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "evidence.md"), []byte("evidence"), 0o644); err != nil {
		t.Fatal(err)
	}
	for name, test := range map[string]struct{ body, want string }{
		"second path":       {"- Path: evidence.md\n  Path: evidence.md\n  Supports: support\n  Drift: drift\n", "unexpected field Path"},
		"second url":        {"- URL: https://example.invalid/one\n  URL: https://example.invalid/two\n  Supports: support\n  Drift: drift\n", "unexpected field URL"},
		"mixed locator":     {"- Path: evidence.md\n  URL: https://example.invalid/two\n  Supports: support\n  Drift: drift\n", "unexpected field URL"},
		"unknown field":     {"- Path: evidence.md\n  Owner: team\n  Supports: support\n  Drift: drift\n", "unexpected field Owner"},
		"duplicate support": {"- Path: evidence.md\n  Supports: support\n  Supports: repeated\n  Drift: drift\n", "duplicate field Supports"},
		"reordered fields":  {"- Path: evidence.md\n  Drift: drift\n  Supports: support\n", "field Drift is out of order; expected Supports"},
		"empty field":       {"- Path: evidence.md\n  Supports: \n  Drift: drift\n", "field Supports must be non-empty"},
	} {
		t.Run(name, func(t *testing.T) {
			if diagnostics := sourceDiagnostics(root, test.body); !hasDiagnostic(diagnostics, test.want) {
				t.Fatalf("diagnostics = %v, want %q", diagnostics, test.want)
			}
		})
	}
}

func hasDiagnostic(diagnostics []Diagnostic, want string) bool {
	for _, diagnostic := range diagnostics {
		if strings.Contains(diagnostic.Message, want) {
			return true
		}
	}
	return false
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
