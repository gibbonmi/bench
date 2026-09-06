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

// readyIndex renders one ready-map index. The graph and readiness rules read the
// ticket files, so each caller writes the index beside the ticket it grades.
func readyIndex(status, fog, sources string) string {
	return "# Ready\n\nStatus: " + status +
		"\n\n## Destination\n\nSettle it.\n\n## Notes\n\n## Decisions so far\n\n- [Answer](ready/tickets/1.md): Resolved.\n" +
		"\n## Not yet specified\n" + fog +
		"\n## Spec-writer discretion\n\n- A bounded choice.\n\n## Out of scope\n\n- Not included.\n\n## Sources\n" + sources
}

func TestMapReadinessAndStructuredSources(t *testing.T) {
	sources := strings.ReplaceAll("\n- Path: TICKevidence.mdTICK\n  Supports: The settled answer.\n  Drift: Update when evidence changes.\n"+
		"- URL: TICKhttps://example.invalid/sourceTICK\n  Supports: External context.\n  Drift: Update when it changes.\n", "TICK", "`")
	const resolvedTicket = "# Answer\n\nBlocked by: none\nType: Grill\n\n### Question\n\nAnswer?\n\n### Answer\n\nRESULT\n"
	validate := func(t *testing.T, index, answer string) []string {
		t.Helper()
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "evidence.md"), []byte("evidence"), 0o644); err != nil {
			t.Fatal(err)
		}
		writeSplitMap(t, root, "decisions/ready.md", index, map[string]string{"1.md": strings.Replace(resolvedTicket, "RESULT", answer, 1)})
		return ValidateDecisionMapTree(root)
	}
	if diagnostics := validate(t, readyIndex("shaping", "\n- Honest fog.\n", sources), "Resolved."); len(diagnostics) != 0 {
		t.Fatalf("shaping diagnostics = %v", diagnostics)
	}
	if diagnostics := validate(t, readyIndex("ready", "\n- Honest fog.\n", sources), "Resolved."); !hasMessage(diagnostics, "decisions/ready.md: ready map has non-empty Not yet specified") {
		t.Fatalf("ready diagnostics = %v", diagnostics)
	}
	for _, marker := range []string{"— (open)", "— (deferred)", "GRILL DEFERRED"} {
		// An unresolved ticket also drops its gist, so the index under proof carries none.
		index := strings.Replace(readyIndex("ready", "", sources), "\n- [Answer](ready/tickets/1.md): Resolved.\n", "", 1)
		diagnostics := validate(t, index, marker)
		if !hasTreeDiagnostic(diagnostics, "ready map has unresolved ticket") {
			t.Errorf("ready marker %q diagnostics = %v", marker, diagnostics)
		}
	}
	if diagnostics := validate(t, readyIndex("ready", "", ""), "Resolved."); len(diagnostics) != 0 {
		t.Fatalf("empty sources diagnostics = %v", diagnostics)
	}
	if diagnostics := validate(t, readyIndex("ready", "", "\nn/a\n"), "Resolved."); !hasTreeDiagnostic(diagnostics, "Sources entry") {
		t.Fatalf("sentinel diagnostics = %v", diagnostics)
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
	index := func(status, gist, fog, discretion string) string {
		return "# Map\n\nStatus: " + status + "\n\n## Destination\n\nSettle it.\n\n## Notes\n\n## Decisions so far\n" + gist +
			"\n## Not yet specified\n" + fog + "\n## Spec-writer discretion\n" + discretion + "\n## Out of scope\n\n- Excluded.\n\n## Sources\n"
	}
	ticket := func(answer string) string {
		return "# Decision\n\nBlocked by: none\nType: Grill\n\n### Question\n\nQuestion?\n\n### Answer\n" + answer
	}
	validate := func(t *testing.T, index, ticket string) []string {
		t.Helper()
		root := t.TempDir()
		writeSplitMap(t, root, "decisions/map.md", index, map[string]string{"1.md": ticket})
		return ValidateDecisionMapTree(root)
	}
	const gist = "\n- [Decision](map/tickets/1.md): Resolved.\n"
	wrapped := index("shaping", gist, "\n- A fog item\n  that continues on the next line.\n\n", "\n- A bounded choice.\n\n")
	if diagnostics := validate(t, wrapped, ticket("\nResolved.\n")); len(diagnostics) != 0 {
		t.Fatalf("wrapped list diagnostics = %v", diagnostics)
	}
	prose := strings.Replace(wrapped, "- A bounded choice.", "A bounded choice.", 1)
	if diagnostics := validate(t, prose, ticket("\nResolved.\n")); !hasTreeDiagnostic(diagnostics, "Spec-writer discretion must be a Markdown bullet list") {
		t.Fatalf("prose diagnostics = %v", diagnostics)
	}
	empty := index("shaping", "", "", "")
	if diagnostics := validate(t, empty, ticket("\n")); len(diagnostics) != 0 {
		t.Fatalf("shaping empty-answer diagnostics = %v", diagnostics)
	}
	ready := strings.Replace(empty, "Status: shaping", "Status: ready", 1)
	if diagnostics := validate(t, ready, ticket("\n")); !hasTreeDiagnostic(diagnostics, "ready map has unresolved ticket") {
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
