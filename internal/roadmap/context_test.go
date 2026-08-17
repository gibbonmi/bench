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
	writeBoard(t, root, [2]string{"**FT1 — first.**", body + "\n"})
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
	writeBoard(t, root, [2]string{"**FT1 — first.**", "roadmap evidence\n"})
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
		{"roadmap_rows", "body", "roadmap evidence"},
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

func TestContextCommandRowSelectorReturnsOnlyCompleteRows(t *testing.T) {
	root := newRepo(t)
	body := strings.Repeat("x", 4113)
	writeBoard(t, root, [2]string{"**FT1 — first.**", body + "\n"}, [2]string{"**FT2 — second.**", "other\n"})
	if err := os.MkdirAll(filepath.Join(root, "capture"), 0o755); err != nil {
		t.Fatal(err)
	}
	longCapture := strings.Repeat("capture-body", 500)
	if err := os.WriteFile(filepath.Join(root, "capture", "IDEAS.md"), []byte("- 2026-01-01  "+longCapture+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := ContextCommand([]string{"--context", "--row", "FT2,FT1,FT1"}, func(string) GateCacheFact { return GateCacheFact{} })
	if code != 0 {
		t.Fatalf("selector exit = %d, output=%q", code, out)
	}
	document, err := axitest.DecodeDocument(out)
	if err != nil {
		t.Fatal(err)
	}
	rows, err := document.Rows("roadmap_rows")
	if err != nil || len(rows) != 2 {
		t.Fatalf("selected rows = %#v, %v", rows, err)
	}
	for i, want := range []struct{ id, body string }{{"FT2", "other"}, {"FT1", body}} {
		row, ok := rows[i].(map[string]any)
		if !ok || row["id"] != want.id || row["body"] != want.body || row["body_bytes"] != float64(len(want.body)) {
			t.Fatalf("selected row = %#v", rows[i])
		}
	}
	if strings.Contains(out, longCapture) {
		t.Fatal("complete capture body leaked into selector output")
	}
	if document.Blocks[len(document.Blocks)-1] != "help" {
		t.Fatalf("blocks = %q", document.Blocks)
	}
	if _, err := document.HelpActions(); err != nil {
		t.Fatalf("help = %v", err)
	}
}

// TestContextCommandSourcesListsRoadmapDirectory covers PR11 (stories 13, 26): the
// sources block gains a roadmap/ row reporting the split directory's state and byte
// total, and the context row still reads schema 4.
func TestContextCommandSourcesListsRoadmapDirectory(t *testing.T) {
	root := newRepo(t)
	const body = "row detail\n"
	writeBoard(t, root, [2]string{"**FT1 — first.**", body})
	out, code := ContextCommand([]string{"--context"}, func(string) GateCacheFact { return GateCacheFact{} })
	if code != 0 {
		t.Fatalf("exit = %d, output=%q", code, out)
	}
	document, err := axitest.DecodeDocument(out)
	if err != nil {
		t.Fatal(err)
	}
	rows, err := document.Rows("sources")
	if err != nil {
		t.Fatal(err)
	}
	var found map[string]any
	for _, r := range rows {
		row := r.(map[string]any)
		if row["source"] == "roadmap/" {
			found = row
		}
	}
	wantBytes := float64(len("**FT1 — first.**\n" + body))
	if found == nil || found["state"] != "parsed" || found["bytes"] != wantBytes {
		t.Fatalf("sources = %#v, want a roadmap/ row parsed,%v", rows, wantBytes)
	}
	contextRows, err := document.Rows("context")
	if err != nil || len(contextRows) != 1 || contextRows[0].(map[string]any)["schema"] != float64(4) {
		t.Fatalf("context rows = %#v, %v", contextRows, err)
	}
}

// TestContextCommandDiagnosticRendersFailureAndFlipsRoadmapMalformed covers PR12
// (story 14): a missing detail owner renders as a parse_failures row sourced at the
// offending row file, ROADMAP.md's own sources row flips to malformed, and the
// context row's sequence_trusted goes false — a general predicate over any
// diagnostic, not only this fault class.
func TestContextCommandDiagnosticRendersFailureAndFlipsRoadmapMalformed(t *testing.T) {
	root := newRepo(t)
	if err := os.WriteFile(roadmapPath(t, root), []byte("**FT7 (LOW) — x.**\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := ContextCommand([]string{"--context"}, func(string) GateCacheFact { return GateCacheFact{} })
	if code != 0 {
		t.Fatalf("exit = %d, output=%q", code, out)
	}
	document, err := axitest.DecodeDocument(out)
	if err != nil {
		t.Fatal(err)
	}
	rows, err := document.Rows("parse_failures")
	if err != nil {
		t.Fatal(err)
	}
	var found bool
	for _, r := range rows {
		row := r.(map[string]any)
		if row["source"] == "roadmap/FT7.md" && row["reason"] == "missing detail owner for ROADMAP.md row FT7" {
			found = true
		}
	}
	if !found {
		t.Fatalf("parse_failures = %#v, want roadmap/FT7.md missing detail owner", rows)
	}
	sourceRows, err := document.Rows("sources")
	if err != nil {
		t.Fatal(err)
	}
	var roadmapState any
	for _, r := range sourceRows {
		row := r.(map[string]any)
		if row["source"] == RoadmapFile {
			roadmapState = row["state"]
		}
	}
	if roadmapState != "malformed" {
		t.Fatalf("ROADMAP.md source state = %v, want malformed", roadmapState)
	}
	contextRows, err := document.Rows("context")
	if err != nil || len(contextRows) != 1 || contextRows[0].(map[string]any)["sequence_trusted"] != false {
		t.Fatalf("context rows = %#v, %v", contextRows, err)
	}
}

// TestContextCommandRowSelectorRendersRowFileBody covers PR13 (story 15): --row
// returns the requested row's body straight from its row file, with body_bytes.
func TestContextCommandRowSelectorRendersRowFileBody(t *testing.T) {
	root := newRepo(t)
	writeBoard(t, root, [2]string{"**FT7 (LOW) — x.**", "row detail\n"})
	out, code := ContextCommand([]string{"--context", "--row", "FT7"}, func(string) GateCacheFact { return GateCacheFact{} })
	if code != 0 {
		t.Fatalf("exit = %d, output=%q", code, out)
	}
	document, err := axitest.DecodeDocument(out)
	if err != nil {
		t.Fatal(err)
	}
	rows, err := document.Rows("roadmap_rows")
	if err != nil || len(rows) != 1 {
		t.Fatalf("roadmap_rows = %#v, %v", rows, err)
	}
	row := rows[0].(map[string]any)
	if row["body"] != "row detail" || row["body_bytes"] != float64(len("row detail")) {
		t.Fatalf("row = %#v", row)
	}
}

// TestContextCommandDirectoryRowFileRendersFailureAndUntrustsSequence covers PR27
// (story 44, the render-side half shared with ticket 1's loader): a row file the
// classifier reports wrong-type — a directory sitting where the row file belongs —
// renders a parse_failures row naming it and drops sequence_trusted.
func TestContextCommandDirectoryRowFileRendersFailureAndUntrustsSequence(t *testing.T) {
	root := newRepo(t)
	if err := os.WriteFile(roadmapPath(t, root), []byte("**FT7 (LOW) — x.**\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, RoadmapDir, "FT7.md"), 0o755); err != nil {
		t.Fatal(err)
	}
	out, code := ContextCommand([]string{"--context"}, func(string) GateCacheFact { return GateCacheFact{} })
	if code != 0 {
		t.Fatalf("exit = %d, output=%q", code, out)
	}
	document, err := axitest.DecodeDocument(out)
	if err != nil {
		t.Fatal(err)
	}
	rows, err := document.Rows("parse_failures")
	if err != nil {
		t.Fatal(err)
	}
	var found bool
	for _, r := range rows {
		row := r.(map[string]any)
		if row["source"] == "roadmap/FT7.md" {
			found = true
		}
	}
	if !found {
		t.Fatalf("parse_failures = %#v, want roadmap/FT7.md named", rows)
	}
	contextRows, err := document.Rows("context")
	if err != nil || len(contextRows) != 1 || contextRows[0].(map[string]any)["sequence_trusted"] != false {
		t.Fatalf("context rows = %#v, %v", contextRows, err)
	}
}

// TestContextCommandUnrecognizedFileColonInPathRendersFullSource covers the
// repair-typed-diagnostic ticket: an unrecognized-file basename that legally
// contains ": " (nothing in the roadmap/ listing grammar forbids it) must render its
// whole path as parse_failures.source. The old strings.Cut(d, ": ")-on-the-formatted-
// string approach cut at the basename's own ": " and reported a truncated,
// nonexistent path instead.
func TestContextCommandUnrecognizedFileColonInPathRendersFullSource(t *testing.T) {
	root := newRepo(t)
	writeBoard(t, root, [2]string{"**FT7 (LOW) — x.**", ""})
	if err := os.WriteFile(filepath.Join(root, RoadmapDir, "x: y.md"), []byte("scratch\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := ContextCommand([]string{"--context"}, func(string) GateCacheFact { return GateCacheFact{} })
	if code != 0 {
		t.Fatalf("exit = %d, output=%q", code, out)
	}
	document, err := axitest.DecodeDocument(out)
	if err != nil {
		t.Fatal(err)
	}
	rows, err := document.Rows("parse_failures")
	if err != nil {
		t.Fatal(err)
	}
	const wantSource = "roadmap/x: y.md"
	const wantReason = "unrecognized file under roadmap/; expected <row ID>.md"
	var found, truncated bool
	for _, r := range rows {
		row := r.(map[string]any)
		if row["source"] == wantSource && row["reason"] == wantReason {
			found = true
		}
		if row["source"] == "roadmap/x" {
			truncated = true
		}
	}
	if truncated {
		t.Fatalf("parse_failures = %#v, source truncated at the basename's own \": \"", rows)
	}
	if !found {
		t.Fatalf("parse_failures = %#v, want source %q reason %q", rows, wantSource, wantReason)
	}
}

func TestContextCommandRowSelectorRefusesMalformedAndMissing(t *testing.T) {
	root := newRepo(t)
	if err := os.WriteFile(roadmapPath(t, root), []byte("**FT1 — first.**\nbody\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		args []string
		code int
		want []string
	}{
		{[]string{"--context", "--row", ""}, 2, []string{"usage:"}},
		{[]string{"--context", "--row", "FT1,"}, 2, []string{"usage:"}},
		{[]string{"--context", "--row", "not-an-id"}, 2, []string{"usage:"}},
		{[]string{"--context", "--row", "FT1,NOPE2,NOPE3"}, 1, []string{"NOPE2", "NOPE3"}},
		{[]string{"--context", "--row", "FT1", "--full"}, 2, []string{"usage:"}},
	} {
		out, code := ContextCommand(tc.args, func(string) GateCacheFact { return GateCacheFact{} })
		missing := false
		for _, want := range tc.want {
			if !strings.Contains(out, want) {
				missing = true
			}
		}
		if code != tc.code || missing {
			t.Fatalf("args=%v = %q/%d, want %q/%d", tc.args, out, code, tc.want, tc.code)
		}
		if tc.code == 1 && (strings.Contains(out, "roadmap_rows[") || strings.Contains(out, "body_bytes")) {
			t.Fatalf("partial selector output = %q", out)
		}
	}
}

func TestBuildContextCarriesRetrosAndDegradedEvidence(t *testing.T) {
	root := newRepo(t)
	dir := filepath.Join(root, "capture", "retros")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "b.md"), []byte("second"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.md"), []byte("first"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("missing.md", filepath.Join(dir, "bad.md")); err != nil {
		t.Fatal(err)
	}
	s, err := BuildContext(root, false, GateCacheFact{})
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Retros) != 3 || s.Retros[0].Path != "capture/retros/a.md" || s.Retros[1].Path != "capture/retros/b.md" || s.Retros[2].State != "unreadable" {
		t.Fatalf("retros = %#v", s.Retros)
	}
	if len(s.Sources) < 2 || s.Sources[len(s.Sources)-2].Source != "capture/retros/" {
		t.Fatalf("sources = %#v", s.Sources)
	}
	if got, err := renderContext(s); err != nil || !strings.Contains(got, "retros[3]{path,state,body,body_bytes}:") || !strings.Contains(got, "capture/retros/bad.md,unreadable") || !strings.Contains(got, "parse_failures[1]{") {
		t.Fatalf("context = %q, %v", got, err)
	}
}

func TestParseDocumentOccurrenceLedgers(t *testing.T) {
	index, files := board(
		[2]string{"**FT1 — one.**", "Body\nOccurrences: alpha-1, beta-2\n"},
		[2]string{"**FT2 — two.**", "Body\nOccurrences: bad, bad\n"},
	)
	doc, failures, _ := ParseDocument(splitTree(index, files), nil, false)
	if got := doc.Rows[0]; got.OccurrenceKeys != "alpha-1, beta-2" || got.OccurrenceCount != 2 {
		t.Fatalf("valid ledger = %#v", got)
	}
	if len(failures) != 1 || failures[0].Reason != "malformed-ledger" || failures[0].Source != "roadmap/FT2.md" {
		t.Fatalf("failures = %#v", failures)
	}
}

func TestOccurrenceIncidentGrammar(t *testing.T) {
	for _, key := range []string{"a", strings.Repeat("a", 64), "a-1"} {
		if !ValidOccurrenceIncident(key) {
			t.Fatalf("valid incident rejected: %q", key)
		}
	}
	for _, key := range []string{"", strings.Repeat("a", 65), "A", "é", "a_b", "-a", "a-", "a\tb", "a\rb", "a\nb"} {
		if ValidOccurrenceIncident(key) {
			t.Fatalf("invalid incident accepted: %q", key)
		}
	}
}

func TestOccurrenceLedgerMalformedAndLineEndings(t *testing.T) {
	const heading = "**FT1 — one.**"
	valid, failures, _ := ParseDocument(splitTree(heading+"\n", map[string]string{"FT1.md": heading + "\r\nOccurrences: alpha-1, beta-2"}), nil, false)
	if len(failures) != 0 || valid.Rows[0].OccurrenceCount != 2 {
		t.Fatalf("CRLF newline-less ledger = %#v, %#v", valid, failures)
	}
	for _, ledger := range []string{"Occurrences:", "Occurrences: alpha_1", "Occurrences: beta, alpha", "Occurrences: alpha, alpha", "Occurrences: alpha\nOccurrences: beta"} {
		doc, got, _ := ParseDocument(splitTree(heading+"\n", map[string]string{"FT1.md": heading + "\n" + ledger + "\n"}), nil, false)
		if len(got) != 1 || got[0].Reason != "malformed-ledger" || len(doc.OccurrenceDiscrepancies) != 1 {
			t.Fatalf("ledger %q accepted: %#v, %#v", ledger, doc, got)
		}
	}
}

func TestBuildContextProjectsPendingCaptureOccurrences(t *testing.T) {
	root := newRepo(t)
	for path, body := range map[string]string{
		RoadmapFile:             "**FT1 — one.**\n\n**FT2 — two.**\n",
		IdeasFile:               "- 2026-07-10  later [occurrence:FT2/idea-2]\n- 2026-07-10  first [occurrence:FT1/idea-1]\n",
		learnings.JournalPath:   "## 2026-07-10 — lesson  [open]\nbody [occurrence:FT1/learning-1]\n",
		"capture/retros/one.md": "## Agent-experience improvements\n\nParagraph [occurrence:FT2/retro-1]\n",
	} {
		if err := os.MkdirAll(filepath.Dir(filepath.Join(root, path)), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, path), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	s, err := BuildContext(root, false, GateCacheFact{})
	if err != nil {
		t.Fatal(err)
	}
	out, err := renderContext(s)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"FT1,idea-1,capture/IDEAS.md,line 2,pending",
		"FT1,learning-1,capture/learnings.md,line 1,pending",
		"FT2,idea-2,capture/IDEAS.md,line 1,pending",
		"FT2,retro-1,capture/retros/one.md,line 3,pending",
	}
	last := -1
	for _, row := range want {
		at := strings.Index(out, row)
		if at < 0 || at <= last {
			t.Fatalf("capture rows = %s", out)
		}
		last = at
	}
}

func TestBuildContextClassifiesOccurrenceDiscrepancies(t *testing.T) {
	root := newRepo(t)
	writeBoard(t, root, [2]string{"**FT1 — one.**", "Occurrences: recorded\n"}, [2]string{"**FT2 — two.**", "Occurrences: bad, bad\n"})
	ideas := strings.Join([]string{
		"- 2026-07-10  duplicate [occurrence:FT1/recorded]",
		"- 2026-07-10  malformed [occurrence:FT1/unterminated",
		"- 2026-07-10  unknown [occurrence:FT9/new]",
		"- 2026-07-10  multiple [occurrence:FT1/one] [occurrence:FT1/two]",
	}, "\n")
	if err := os.WriteFile(filepath.Join(root, IdeasFile), []byte(ideas), 0o644); err != nil {
		t.Fatal(err)
	}

	s, err := BuildContext(root, false, GateCacheFact{})
	if err != nil {
		t.Fatal(err)
	}
	if s.SequenceTrusted {
		t.Fatal("structural discrepancies left sequence trusted")
	}
	out, err := renderContext(s)
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range []string{
		"FT1,recorded,capture/IDEAS.md,line 1,already-recorded",
		"capture/IDEAS.md,line 1,already-recorded,FT1,recorded,false",
		"capture/IDEAS.md,line 2,malformed-token,FT1,unterminated,true",
		"capture/IDEAS.md,line 3,unknown-owner,FT9,new,true",
		"capture/IDEAS.md,line 4,multiple-tokens,FT1,two,true",
		"roadmap/FT2.md,FT2,malformed-ledger,FT2,\"\",true",
	} {
		if !strings.Contains(out, row) {
			t.Fatalf("missing %q in %s", row, out)
		}
	}
}

func TestBuildContextTreatsSeparatedValidTokensAsMultiple(t *testing.T) {
	root := newRepo(t)
	if err := os.WriteFile(filepath.Join(root, RoadmapFile), []byte("**FT1 — one.**\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, IdeasFile), []byte("- 2026-07-10  first [occurrence:FT1/one] with prose between [occurrence:FT1/two]\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	s, err := BuildContext(root, false, GateCacheFact{})
	if err != nil {
		t.Fatal(err)
	}
	if s.SequenceTrusted || len(s.CaptureOccurrences) != 0 {
		t.Fatalf("multiple tokens = trusted:%v occurrences:%#v", s.SequenceTrusted, s.CaptureOccurrences)
	}
	if got, want := s.Roadmap.OccurrenceDiscrepancies, []OccurrenceDiscrepancy{{Source: IdeasFile, CaptureUnit: "line 1", Kind: "multiple-tokens", Owner: "FT1", Incident: "two", Structural: true}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("discrepancies = %#v, want %#v", got, want)
	}
}

func TestBuildContextKeepsLearningAndRetroOccurrencesInTheirUnits(t *testing.T) {
	root := newRepo(t)
	for path, body := range map[string]string{
		RoadmapFile:           "**FT1 — one.**\n\n**FT2 — two.**\n",
		learnings.JournalPath: "## 2026-07-10 — lesson  [open]\nProse [occurrence:FT2/ignored_] does not end the entry.\nFinal [occurrence:FT1/learning-final]\n",
		"capture/retros/one.md": strings.Join([]string{
			"Outside [occurrence:FT2/ignored]",
			"",
			"## Agent-experience improvements",
			"",
			"Paragraph [occurrence:FT1/retro-paragraph]",
			"",
			"- First item [occurrence:FT1/retro-list]",
			"- Second item [occurrence:FT2/retro-other-owner]",
			"",
			"## Other",
			"Ignored [occurrence:FT2/ignored-too]",
		}, "\n"),
	} {
		if err := os.MkdirAll(filepath.Dir(filepath.Join(root, path)), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, path), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	s, err := BuildContext(root, false, GateCacheFact{})
	if err != nil {
		t.Fatal(err)
	}
	out, err := renderContext(s)
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range []string{
		"FT1,learning-final,capture/learnings.md,line 1,pending",
		"FT1,retro-list,capture/retros/one.md,line 7,pending",
		"FT1,retro-paragraph,capture/retros/one.md,line 5,pending",
		"FT2,retro-other-owner,capture/retros/one.md,line 8,pending",
	} {
		if !strings.Contains(out, row) {
			t.Fatalf("missing %q in %s", row, out)
		}
	}
	for _, incident := range []string{"ignored", "ignored-too"} {
		if strings.Contains(out, "FT2,"+incident+",") {
			t.Fatalf("non-final token %q became an occurrence: %s", incident, out)
		}
	}
}

func TestBuildContextNormalizesPendingOccurrencePairs(t *testing.T) {
	root := newRepo(t)
	if err := os.WriteFile(filepath.Join(root, RoadmapFile), []byte("**FT1 — one.**\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, IdeasFile), []byte("- 2026-07-10  first [occurrence:FT1/shared]\n- 2026-07-10  second [occurrence:FT1/shared]\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	s, err := BuildContext(root, false, GateCacheFact{})
	if err != nil {
		t.Fatal(err)
	}
	if len(s.CaptureOccurrences) != 2 {
		t.Fatalf("capture occurrences = %#v", s.CaptureOccurrences)
	}
	if got, want := s.PendingOccurrences, []OccurrencePair{{Owner: "FT1", Incident: "shared"}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("pending pairs = %#v, want %#v", got, want)
	}
}

func TestBuildContextKeepsSameIncidentForSeparateOwners(t *testing.T) {
	root := newRepo(t)
	if err := os.WriteFile(filepath.Join(root, RoadmapFile), []byte("**FT1 — one.**\n\n**FT2 — two.**\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, IdeasFile), []byte("- 2026-07-10  first [occurrence:FT2/shared]\n- 2026-07-10  second [occurrence:FT1/shared]\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	s, err := BuildContext(root, false, GateCacheFact{})
	if err != nil {
		t.Fatal(err)
	}
	want := []OccurrencePair{{Owner: "FT1", Incident: "shared"}, {Owner: "FT2", Incident: "shared"}}
	if !reflect.DeepEqual(s.PendingOccurrences, want) {
		t.Fatalf("pending pairs = %#v, want %#v", s.PendingOccurrences, want)
	}
}

func TestBuildContextProjectsEveryRecordedSourceWithoutPendingPair(t *testing.T) {
	root := newRepo(t)
	writeBoard(t, root, [2]string{"**FT1 — one.**", "Occurrences: recorded\n"})
	if err := os.WriteFile(filepath.Join(root, IdeasFile), []byte("- 2026-07-10  first [occurrence:FT1/recorded]\n- 2026-07-10  second [occurrence:FT1/recorded]\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	s, err := BuildContext(root, false, GateCacheFact{})
	if err != nil {
		t.Fatal(err)
	}
	if !s.SequenceTrusted {
		t.Fatal("only advisory discrepancies made the sequence untrusted")
	}
	if len(s.PendingOccurrences) != 0 {
		t.Fatalf("recorded pair remained pending: %#v", s.PendingOccurrences)
	}
	if got, want := s.CaptureOccurrences, []CaptureOccurrence{
		{Owner: "FT1", Incident: "recorded", Source: IdeasFile, CaptureUnit: "line 1", State: "already-recorded"},
		{Owner: "FT1", Incident: "recorded", Source: IdeasFile, CaptureUnit: "line 2", State: "already-recorded"},
	}; !reflect.DeepEqual(got, want) {
		t.Fatalf("capture occurrences = %#v, want %#v", got, want)
	}
	if got, want := s.Roadmap.OccurrenceDiscrepancies, []OccurrenceDiscrepancy{
		{Source: IdeasFile, CaptureUnit: "line 1", Kind: "already-recorded", Owner: "FT1", Incident: "recorded"},
		{Source: IdeasFile, CaptureUnit: "line 2", Kind: "already-recorded", Owner: "FT1", Incident: "recorded"},
	}; !reflect.DeepEqual(got, want) {
		t.Fatalf("discrepancies = %#v, want %#v", got, want)
	}
	out, err := renderContext(s)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "4,false,true") {
		t.Fatalf("trusted context header missing: %s", out)
	}
	if !strings.Contains(out, "capture/IDEAS.md,line 1,already-recorded,FT1,recorded,false") {
		t.Fatalf("advisory discrepancy missing from context: %s", out)
	}
}

func TestBuildContextRequiresUsableCaptureSourcesForOccurrenceTrust(t *testing.T) {
	for _, tc := range []struct {
		name    string
		prepare func(t *testing.T, root string)
	}{
		{
			name: "ideas wrong type",
			prepare: func(t *testing.T, root string) {
				t.Helper()
				if err := os.Mkdir(filepath.Join(root, IdeasFile), 0o755); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "learnings malformed",
			prepare: func(t *testing.T, root string) {
				t.Helper()
				path := filepath.Join(root, learnings.JournalPath)
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte{0xff}, 0o644); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "retros wrong type",
			prepare: func(t *testing.T, root string) {
				t.Helper()
				path := filepath.Join(root, "capture", "retros")
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte("not a directory"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newRepo(t)
			if err := os.WriteFile(filepath.Join(root, RoadmapFile), []byte("**FT1 — one.**\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			tc.prepare(t, root)

			s, err := BuildContext(root, false, GateCacheFact{})
			if err != nil {
				t.Fatal(err)
			}
			if s.SequenceTrusted {
				t.Fatal("degraded capture source left recurrence sequence trusted")
			}
			if _, err := renderContext(s); err != nil {
				t.Fatalf("complete context did not render: %v", err)
			}
		})
	}
}
