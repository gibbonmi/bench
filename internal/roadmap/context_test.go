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
		if len(rows) != 0 || !strings.HasSuffix(out, "help[0]{cmd,why}:\n") {
			t.Fatalf("args=%v help rows/output = %d/%q", args, len(rows), out)
		}
	}
}

func TestContextCommandIndexOmitsRoadmapBodiesWithTrueSizes(t *testing.T) {
	root := newRepo(t)
	const body = "complete roadmap evidence"
	if err := os.WriteFile(roadmapPath(t, root), []byte("**FT1 — first.**\n"+body+"\n"), 0o644); err != nil {
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
	files := map[string]string{
		RoadmapFile:                   "**FT1 — first.**\nroadmap evidence\n",
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
	if err := os.WriteFile(roadmapPath(t, root), []byte("**FT1 — first.**\n"+body+"\n\n**FT2 — second.**\nother\n"), 0o644); err != nil {
		t.Fatal(err)
	}
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
	content := []byte("**FT1 — one.** Body\nOccurrences: alpha-1, beta-2\n\n**FT2 — two.** Body\nOccurrences: bad, bad\n")
	doc, failures := ParseDocument(content, nil, false)
	if got := doc.Rows[0]; got.OccurrenceKeys != "alpha-1, beta-2" || got.OccurrenceCount != 2 {
		t.Fatalf("valid ledger = %#v", got)
	}
	if len(failures) != 1 || failures[0].Reason != "malformed-ledger" {
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
	valid, failures := ParseDocument([]byte("**FT1 — one.**\r\nOccurrences: alpha-1, beta-2"), nil, false)
	if len(failures) != 0 || valid.Rows[0].OccurrenceCount != 2 {
		t.Fatalf("CRLF newline-less ledger = %#v, %#v", valid, failures)
	}
	for _, ledger := range []string{"Occurrences:", "Occurrences: alpha_1", "Occurrences: beta, alpha", "Occurrences: alpha, alpha", "Occurrences: alpha\nOccurrences: beta"} {
		doc, got := ParseDocument([]byte("**FT1 — one.**\n"+ledger+"\n"), nil, false)
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
	if err := os.WriteFile(filepath.Join(root, RoadmapFile), []byte("**FT1 — one.**\nOccurrences: recorded\n\n**FT2 — two.**\nOccurrences: bad, bad\n"), 0o644); err != nil {
		t.Fatal(err)
	}
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
		"ROADMAP.md,FT2,malformed-ledger,FT2,\"\",true",
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
	if err := os.WriteFile(filepath.Join(root, RoadmapFile), []byte("**FT1 — one.**\nOccurrences: recorded\n"), 0o644); err != nil {
		t.Fatal(err)
	}
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
