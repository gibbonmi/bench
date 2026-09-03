// Tests for decision-map graph edges, readiness, and source diagnostics.
package maps

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
)

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
		"second path":              {"- Path: evidence.md\n  Path: evidence.md\n  Supports: support\n  Drift: drift\n", "unexpected field Path"},
		"second url":               {"- URL: https://example.invalid/one\n  URL: https://example.invalid/two\n  Supports: support\n  Drift: drift\n", "unexpected field URL"},
		"mixed locator":            {"- Path: evidence.md\n  URL: https://example.invalid/two\n  Supports: support\n  Drift: drift\n", "unexpected field URL"},
		"unknown field":            {"- Path: evidence.md\n  Owner: team\n  Supports: support\n  Drift: drift\n", "unexpected field Owner"},
		"duplicate support":        {"- Path: evidence.md\n  Supports: support\n  Supports: repeated\n  Drift: drift\n", "duplicate field Supports"},
		"reordered fields":         {"- Path: evidence.md\n  Drift: drift\n  Supports: support\n", "field Drift is out of order; expected Supports"},
		"empty field":              {"- Path: evidence.md\n  Supports: \n  Drift: drift\n", "field Supports must be non-empty"},
		"wrapped field":            {"- Path: evidence.md\n  Supports: support\n  that continues here.\n  Drift: drift\n", "Sources evidence.md line \"that continues here.\" has no field name; write each Sources record field on one physical line"},
		"wrapped field with colon": {"- Path: evidence.md\n  Supports: support\n  and it continues: here.\n  Drift: drift\n", "Sources evidence.md line \"and it continues: here.\" has no field name; write each Sources record field on one physical line"},
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

// TestLiftedGraphWalkReportsEveryEdgeFault drives the lifted walk with a
// neutral vocabulary, so no decision-map word reaches the shared symbol.
func TestLiftedGraphWalkReportsEveryEdgeFault(t *testing.T) {
	walk := GraphWalk{
		Names: []string{"a", "b", "c"},
		Edges: [][]string{{"b", "b", "z"}, {"b", "c"}, {"b"}},
		Fault: func(fault GraphFault, node int, target string) string {
			return string(fault) + " " + []string{"a", "b", "c"}[node] + "->" + target
		},
	}
	want := []string{
		"duplicate a->b",
		"dangling a->z",
		"self b->b",
		"cycle c->b",
	}
	if got := walk.Diagnostics(); !reflect.DeepEqual(got, want) {
		t.Fatalf("GraphWalk diagnostics = %#v, want %#v", got, want)
	}
	if got := FieldList("none", "none", "#"); got != nil {
		t.Fatalf("FieldList(none) = %#v", got)
	}
	if got := FieldList("#1, #2", "none", "#"); !reflect.DeepEqual(got, []string{"1", "2"}) {
		t.Fatalf("FieldList = %#v", got)
	}
}
