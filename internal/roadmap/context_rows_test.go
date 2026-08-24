// Tests for context row selectors, source listing, and parse-failure diagnostics.
package roadmap

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/learnings"
	"github.com/gibbonmi/bench/internal/roadmap/roadmaptest"
)

func TestContextCommandRowSelectorReturnsOnlyCompleteRows(t *testing.T) {
	root := newRepo(t)
	body := strings.Repeat("x", 4113)
	writeBoard(t, root, Row{"**FT1 — first.**", body + "\n"}, Row{"**FT2 — second.**", "other\n"})
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

// TestContextCommandSourcesListsRoadmapDirectory pins PR11 (stories 13 and 26): the
// sources block gains a roadmap/ row reporting the split directory's state and byte
// total. The context row still reads schema 4. A directory that is not there and one that
// holds nothing are authoritative answers, not degraded reads. Each is therefore pinned
// beside the parsed row, rather than left to the one healthy state.
func TestContextCommandSourcesListsRoadmapDirectory(t *testing.T) {
	const heading, body = "**FT1 — first.**", "row detail\n"
	for _, tc := range []struct {
		name, state string
		bytes       float64
		plant       func(*testing.T, string)
	}{
		{"parsed", "parsed", float64(len(heading + "\n" + body)), func(t *testing.T, root string) {
			writeBoard(t, root, Row{heading, body})
		}},
		{"absent", "absent", 0, func(t *testing.T, root string) {
			roadmaptest.WriteSplitBoard(t, root, heading+"\n", nil)
		}},
		{"empty", "empty", 0, func(t *testing.T, root string) {
			roadmaptest.WriteSplitBoard(t, root, heading+"\n", nil)
			if err := os.Mkdir(filepath.Join(root, RoadmapDir), 0o755); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newRepo(t)
			tc.plant(t, root)
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
			if found == nil || found["state"] != tc.state || found["bytes"] != tc.bytes {
				t.Fatalf("sources = %#v, want a roadmap/ row %s,%v", rows, tc.state, tc.bytes)
			}
			contextRows, err := document.Rows("context")
			if err != nil || len(contextRows) != 1 || contextRows[0].(map[string]any)["schema"] != float64(4) {
				t.Fatalf("context rows = %#v, %v", contextRows, err)
			}
		})
	}
}

// TestContextCommandDiagnosticRendersFailureAndFlipsRoadmapMalformed pins PR12 (story
// 14): a missing detail owner renders as a parse_failures row sourced at the offending
// row file. ROADMAP.md's own sources row flips to malformed. The context row's
// sequence_trusted goes false. This is a general predicate over any diagnostic, not only
// this fault class.
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

// TestContextCommandRowSelectorRendersRowFileBody pins PR13 (story 15): --row returns the
// requested row's body straight from its row file, with body_bytes.
func TestContextCommandRowSelectorRendersRowFileBody(t *testing.T) {
	root := newRepo(t)
	writeBoard(t, root, Row{"**FT7 (LOW) — x.**", "row detail\n"})
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

// TestContextCommandDirectoryRowFileRendersFailureAndUntrustsSequence pins PR27 (story
// 44, the render side): a row file the classifier reports wrong-type. That is a directory
// sitting where the row file belongs. It renders a parse_failures row naming it and drops
// sequence_trusted.
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

// TestContextCommandUnrecognizedFileColonInPathRendersFullSource pins the
// repair-typed-diagnostic ticket: an unrecognized-file basename that legally contains ":
// ". Nothing in the roadmap/ listing grammar forbids that basename. This basename must
// render its whole path as parse_failures.source.
//
// The old strings.Cut(d, ": ")-on-the-formatted-string approach cut at the basename's own
// ": " and reported a truncated, nonexistent path instead.
func TestContextCommandUnrecognizedFileColonInPathRendersFullSource(t *testing.T) {
	root := newRepo(t)
	writeBoard(t, root, Row{Heading: "**FT7 (LOW) — x.**", Body: ""})
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

func TestContextCommandRendersLostDatedLearningLineAsParseFailure(t *testing.T) {
	root := newRepo(t)
	journal := learnings.JournalSchemaHeading + "\n\n- 2026-08-21 — spec anchor drift\n"
	if err := os.WriteFile(filepath.Join(root, learnings.JournalPath), []byte(journal), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := ContextCommand([]string{"--context", "--full"}, func(string) GateCacheFact { return GateCacheFact{} })
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
		if row["source"] == learnings.JournalPath && row["reason"] == "dated learning entry is not a heading" && row["raw"] == "- 2026-08-21 — spec anchor drift" {
			found = true
		}
	}
	if !found {
		t.Fatalf("parse_failures = %#v, want the lost dated line sourced at %s", rows, learnings.JournalPath)
	}
}
