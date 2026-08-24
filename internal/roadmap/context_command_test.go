// Tests for ContextCommand rendering: block list, schema modes, and body projection.
package roadmap

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/learnings"
)

func TestContextCommandEndsWithHelpBlock(t *testing.T) {
	newRepo(t)
	for _, args := range [][]string{{"--context"}, {"--context", "--full"}} {
		out, code := ContextCommand(args, func(string) GateCacheFact { return GateCacheFact{} })
		if code != 0 {
			t.Fatalf("args=%v exit: got %d, want 0; output=%q", args, code, out)
		}
		document, err := axitest.DecodeDocument(out)
		if err != nil {
			t.Fatalf("args=%v did not decode as one TOON document: %v", args, err)
		}
		if len(document.Blocks) == 0 || document.Blocks[len(document.Blocks)-1] != "help" {
			t.Fatalf("args=%v blocks = %q, want terminal help", args, document.Blocks)
		}
		rows, err := document.Rows("help")
		if err != nil {
			t.Fatalf("args=%v help block: %v", args, err)
		}
		wantRows := 0
		if len(args) == 1 {
			wantRows = 1
		}
		if len(rows) != wantRows {
			t.Fatalf("args=%v help rows = %d, want %d", args, len(rows), wantRows)
		}
	}
}

func TestContextCommandIndexDisclosesCompleteRoadmapQueries(t *testing.T) {
	newRepo(t)
	out, code := ContextCommand([]string{"--context"}, func(string) GateCacheFact { return GateCacheFact{} })
	if code != 0 {
		t.Fatalf("exit = %d, output=%q", code, out)
	}
	if !strings.Contains(out, "bench roadmap --context --row <ID,...>") || !strings.Contains(out, "bench roadmap --context --full") {
		t.Fatalf("index help = %q, want row selector and complete snapshot", out)
	}
}

func TestContextCommandIndexOmitsRoadmapBodiesWithTrueSizes(t *testing.T) {
	root := newRepo(t)
	const body = "complete roadmap evidence"
	writeBoard(t, root, Row{"**FT1 — first.**", body + "\n"})
	out, code := ContextCommand([]string{"--context"}, func(string) GateCacheFact { return GateCacheFact{} })
	if code != 0 {
		t.Fatalf("exit = %d, output=%q", code, out)
	}
	document, err := axitest.DecodeDocument(out)
	if err != nil {
		t.Fatal(err)
	}
	rows, err := document.Rows("roadmap_rows")
	if err != nil || len(rows) != 1 {
		t.Fatalf("roadmap rows = %#v, %v", rows, err)
	}
	row := rows[0].(map[string]any)
	if row["body"] != "" || row["body_bytes"] != float64(len(body)) {
		t.Fatalf("roadmap row = %#v", row)
	}
}

func TestContextCommandIndexOmitsCaptureBodiesWithTrueSizes(t *testing.T) {
	root := newRepo(t)
	files := map[string]string{
		IdeasFile:                     "- 2026-01-02  idea evidence\n",
		learnings.JournalPath:         "## 2026-01-03 — lesson  [open]\nlearning evidence\n",
		"capture/retros/completed.md": "retro evidence\n",
	}
	for path, body := range files {
		fullPath := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	out, code := ContextCommand([]string{"--context"}, func(string) GateCacheFact { return GateCacheFact{} })
	if code != 0 {
		t.Fatalf("exit = %d, output=%q", code, out)
	}
	document, err := axitest.DecodeDocument(out)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []struct {
		block, bodyField, bytesField, body string
	}{
		{"ideas", "text", "text_bytes", "idea evidence"},
		{"learnings", "body", "body_bytes", "learning evidence"},
		{"retros", "body", "body_bytes", "retro evidence\n"},
	} {
		rows, err := document.Rows(want.block)
		if err != nil || len(rows) != 1 {
			t.Fatalf("%s rows = %#v, %v", want.block, rows, err)
		}
		row := rows[0].(map[string]any)
		if row[want.bodyField] != "" || row[want.bytesField] != float64(len(want.body)) {
			t.Fatalf("%s row = %#v", want.block, row)
		}
	}
}

func TestContextCommandCarriesCaptureLineNumbersInEveryMode(t *testing.T) {
	root := newRepo(t)
	for path, body := range map[string]string{
		RoadmapFile:           "**FT1 — first.**\nroadmap body\n",
		IdeasFile:             "\n- 2026-01-02  idea evidence\n",
		learnings.JournalPath: "\n## 2026-01-03 — lesson  [open]\nlearning evidence\n",
	} {
		fullPath := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, args := range [][]string{{"--context"}, {"--context", "--row", "FT1"}, {"--context", "--full"}} {
		out, code := ContextCommand(args, func(string) GateCacheFact { return GateCacheFact{} })
		if code != 0 {
			t.Fatalf("args=%v exit = %d, output=%q", args, code, out)
		}
		document, err := axitest.DecodeDocument(out)
		if err != nil {
			t.Fatal(err)
		}
		for _, block := range []string{"ideas", "learnings"} {
			rows, err := document.Rows(block)
			if err != nil || len(rows) != 1 {
				t.Fatalf("args=%v %s rows = %#v, %v", args, block, rows, err)
			}
			if line := rows[0].(map[string]any)["line"]; line != float64(2) {
				t.Fatalf("args=%v %s line = %#v, want 2", args, block, line)
			}
		}
	}
}

func TestContextCommandFullCarriesCompleteBodiesAtSchemaFour(t *testing.T) {
	root := newRepo(t)
	writeBoard(t, root, Row{"**FT1 — first.**", "roadmap evidence\nNext: spec\n"})
	files := map[string]string{
		IdeasFile:                     "- 2026-01-02  idea evidence\nmalformed idea\n",
		learnings.JournalPath:         "## 2026-01-03 — lesson  [open]\nlearning evidence\n",
		"capture/retros/completed.md": "retro evidence\n",
	}
	for path, body := range files {
		fullPath := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	out, code := ContextCommand([]string{"--context", "--full"}, func(string) GateCacheFact { return GateCacheFact{} })
	if code != 0 {
		t.Fatalf("exit = %d, output=%q", code, out)
	}
	document, err := axitest.DecodeDocument(out)
	if err != nil {
		t.Fatal(err)
	}
	contextRows, err := document.Rows("context")
	if err != nil || len(contextRows) != 1 || contextRows[0].(map[string]any)["schema"] != float64(4) {
		t.Fatalf("context rows = %#v, %v", contextRows, err)
	}
	for _, want := range []struct{ block, field, body string }{
		{"roadmap_rows", "body", "roadmap evidence\nNext: spec"},
		{"ideas", "text", "idea evidence"},
		{"learnings", "body", "learning evidence"},
		{"retros", "body", "retro evidence\n"},
		{"parse_failures", "raw", "malformed idea"},
	} {
		rows, err := document.Rows(want.block)
		if err != nil || len(rows) != 1 || rows[0].(map[string]any)[want.field] != want.body {
			t.Fatalf("%s rows = %#v, %v", want.block, rows, err)
		}
	}
}

func TestContextCommandSchemaFourHasNoTruncatedColumn(t *testing.T) {
	root := newRepo(t)
	if err := os.WriteFile(roadmapPath(t, root), []byte("**FT1 — first.**\nbody\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"--context"}, {"--context", "--row", "FT1"}, {"--context", "--full"}} {
		out, code := ContextCommand(args, func(string) GateCacheFact { return GateCacheFact{} })
		if code != 0 {
			t.Fatalf("args=%v exit = %d, output=%q", args, code, out)
		}
		if strings.Contains(out, "truncated") {
			t.Fatalf("args=%v retained truncated column: %s", args, out)
		}
		document, err := axitest.DecodeDocument(out)
		if err != nil {
			t.Fatal(err)
		}
		rows, err := document.Rows("context")
		if err != nil || len(rows) != 1 || rows[0].(map[string]any)["schema"] != float64(4) {
			t.Fatalf("args=%v context rows = %#v, %v", args, rows, err)
		}
	}
}

func TestContextCommandIndexOmitsUnsupportedSchemaRaw(t *testing.T) {
	root := newRepo(t)
	const raw = "sectionless roadmap evidence"
	if err := os.WriteFile(roadmapPath(t, root), []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		args    []string
		wantRaw string
	}{
		{[]string{"--context"}, ""},
		{[]string{"--context", "--full"}, raw},
	} {
		out, code := ContextCommand(tc.args, func(string) GateCacheFact { return GateCacheFact{} })
		if code != 0 {
			t.Fatalf("args=%v exit = %d, output=%q", tc.args, code, out)
		}
		document, err := axitest.DecodeDocument(out)
		if err != nil {
			t.Fatal(err)
		}
		rows, err := document.Rows("parse_failures")
		if err != nil || len(rows) != 1 {
			t.Fatalf("args=%v parse failures = %#v, %v", tc.args, rows, err)
		}
		row := rows[0].(map[string]any)
		if row["raw"] != tc.wantRaw || row["raw_bytes"] != float64(len(raw)) {
			t.Fatalf("args=%v parse failure = %#v", tc.args, row)
		}
	}
}

func TestContextCommandEveryModeEnumeratesTheCompleteBlockList(t *testing.T) {
	root := newRepo(t)
	if err := os.WriteFile(roadmapPath(t, root), []byte("**FT1 — first.**\nbody\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	want := []string{
		"context", "sources", "roadmap_rows", "roadmap_sequence", "ideas", "learnings", "retros",
		"capture_occurrences", "occurrence_discrepancies", "structure", "specs", "spec_history",
		"git", "git_changes", "gate_cache", "parse_failures", "help",
	}
	for _, args := range [][]string{{"--context"}, {"--context", "--row", "FT1"}, {"--context", "--full"}} {
		out, code := ContextCommand(args, func(string) GateCacheFact { return GateCacheFact{} })
		if code != 0 {
			t.Fatalf("args=%v exit = %d, output=%q", args, code, out)
		}
		document, err := axitest.DecodeDocument(out)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(document.Blocks, want) {
			t.Fatalf("args=%v blocks = %q, want %q", args, document.Blocks, want)
		}
	}
}
