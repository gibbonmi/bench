package roadmap

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
)

// TestTreeParseAcceptsHeadingOnlyRowFile covers PR1: an index of one physical heading
// line per row plus a row file repeating that heading with nothing after it is a clean
// board — the row carries the index title and an empty body, and nothing is reported.
func TestTreeParseAcceptsHeadingOnlyRowFile(t *testing.T) {
	const heading = "**FT7 (LOW) — x.**"
	tree := splitTree("# Roadmap\n\n"+heading+"\n", map[string]string{"FT7.md": heading + "\n"})
	doc, failures, diagnostics := ParseDocument(tree, nil, true)
	if len(failures) != 0 || len(diagnostics) != 0 {
		t.Fatalf("clean split board reported failures=%#v diagnostics=%#v", failures, diagnostics)
	}
	if got, want := doc.Rows, []RoadmapRow{{ID: "FT7", Title: "(LOW) — x."}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("rows = %#v, want %#v", got, want)
	}
}

// TestTreeParseProjectsRowFileLedger covers PR2: the occurrence ledger is read from the
// row-file body under today's rules, and a descending ledger yields the same
// malformed-ledger discrepancy, now sourced at the row file that carries it.
func TestTreeParseProjectsRowFileLedger(t *testing.T) {
	const heading = "**FT1 — one.**"
	body := "Occurrence: 2026-08-14 FT189 refused a malformed admin file.\nOccurrences: baseline-01, baseline-02\n"
	doc, failures, diagnostics := ParseDocument(splitTree(heading+"\n", map[string]string{"FT1.md": heading + "\n" + body}), nil, true)
	if len(failures) != 0 || len(diagnostics) != 0 {
		t.Fatalf("row-file ledger reported failures=%#v diagnostics=%#v", failures, diagnostics)
	}
	if got := doc.Rows[0]; got.OccurrenceCount != 2 || got.OccurrenceKeys != "baseline-01, baseline-02" {
		t.Fatalf("ledger row = %#v, want count 2 and both keys", got)
	}

	descending := "Occurrences: baseline-02, baseline-01\n"
	doc, failures, _ = ParseDocument(splitTree(heading+"\n", map[string]string{"FT1.md": heading + "\n" + descending}), nil, true)
	if len(failures) != 1 || failures[0].Reason != "malformed-ledger" || failures[0].Source != "roadmap/FT1.md" {
		t.Fatalf("descending ledger failures = %#v, want one malformed-ledger sourced roadmap/FT1.md", failures)
	}
	want := []OccurrenceDiscrepancy{{Source: "roadmap/FT1.md", CaptureUnit: "FT1", Kind: "malformed-ledger", Owner: "FT1", Structural: true}}
	if !reflect.DeepEqual(doc.OccurrenceDiscrepancies, want) {
		t.Fatalf("discrepancies = %#v, want %#v", doc.OccurrenceDiscrepancies, want)
	}
}

// TestTreeParseDerivesSpecAndTriggerFromRowFile covers PR3: the spec slug, its status,
// and the external-trigger words come from the row-file body, which is where the prose
// that names them now lives.
func TestTreeParseDerivesSpecAndTriggerFromRowFile(t *testing.T) {
	const heading = "**FT1 — one.**"
	body := "Blocked until the deploy is scheduled; the spec is `specs/foo/spec.md`.\n"
	doc, _, _ := ParseDocument(splitTree(heading+"\n", map[string]string{"FT1.md": heading + "\n" + body}), map[string]string{"foo": "staged"}, true)
	if got := doc.Rows[0]; got.Spec != "foo" || got.SpecStatus != "staged" || !got.ExternalTrigger {
		t.Fatalf("row = %#v, want spec foo/staged and external_trigger true", got)
	}
}

// TestTreeParseReportsMissingDetailOwner covers PR4 and the index-side half of PR10: an
// index row with no detail owner is named, and the row survives so rows_total still
// counts the work the board is tracking.
func TestTreeParseReportsMissingDetailOwner(t *testing.T) {
	doc, _, diagnostics := ParseDocument(indexTree("**FT7 (LOW) — x.**\n"), nil, true)
	want := []string{"roadmap/FT7.md: missing detail owner for ROADMAP.md row FT7"}
	if !reflect.DeepEqual(diagnostics, want) {
		t.Fatalf("diagnostics = %#v, want %#v", diagnostics, want)
	}
	if len(doc.Rows) != 1 || doc.Rows[0].ID != "FT7" || doc.Rows[0].Body != "" {
		t.Fatalf("rows = %#v, want the FT7 row kept with an empty body", doc.Rows)
	}
}

// TestTreeParseReportsOrphanRowFile covers PR5: a detail file with no index row is an
// orphan, whether the index merely omits it or is absent from the tree entirely.
func TestTreeParseReportsOrphanRowFile(t *testing.T) {
	const orphan = "**FT8 (LOW) — y.**\n"
	want := "roadmap/FT8.md: orphan detail file with no ROADMAP.md row FT8"
	for _, tc := range []struct{ name, index string }{
		{"index present", "# Roadmap\n"},
		{"index absent", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			doc, _, diagnostics := ParseDocument(splitTree(tc.index, map[string]string{"FT8.md": orphan}), nil, true)
			if !reflect.DeepEqual(diagnostics, []string{want}) {
				t.Fatalf("diagnostics = %#v, want %#v", diagnostics, []string{want})
			}
			if len(doc.Rows) != 0 {
				t.Fatalf("orphan produced rows = %#v", doc.Rows)
			}
		})
	}
}

// TestTreeParseReportsInlineBody covers PR6: a non-blank line under an index line is a
// body that belongs in the row file, and the index text is not the row's body.
func TestTreeParseReportsInlineBody(t *testing.T) {
	const heading = "**FT7 (LOW) — x.**"
	tree := splitTree(heading+"\nBody text.\n\n**FT8 (LOW) — y.**\n", map[string]string{
		"FT7.md": heading + "\nOwned body.\n",
		"FT8.md": "**FT8 (LOW) — y.**\n",
	})
	doc, _, diagnostics := ParseDocument(tree, nil, true)
	want := []string{"ROADMAP.md: row FT7 carries an inline body; move it to roadmap/FT7.md"}
	if !reflect.DeepEqual(diagnostics, want) {
		t.Fatalf("diagnostics = %#v, want %#v", diagnostics, want)
	}
	if len(doc.Rows) != 2 || doc.Rows[0].Body != "Owned body." {
		t.Fatalf("rows = %#v, want both rows kept and FT7's body from its file", doc.Rows)
	}
}

// TestTreeParseReportsHeadingMismatch covers PR7: a row file whose first line has
// drifted from its index line is named with the file and the row, so the two copies
// cannot diverge unnoticed.
func TestTreeParseReportsHeadingMismatch(t *testing.T) {
	tree := splitTree("**FT7 (LOW) — x.**\n", map[string]string{"FT7.md": "**FT7 (LOW) — y.**\nbody\n"})
	doc, _, diagnostics := ParseDocument(tree, nil, true)
	want := []string{"roadmap/FT7.md: heading does not match ROADMAP.md row FT7"}
	if !reflect.DeepEqual(diagnostics, want) {
		t.Fatalf("diagnostics = %#v, want %#v", diagnostics, want)
	}
	if len(doc.Rows) != 1 || doc.Rows[0].Title != "(LOW) — x." || doc.Rows[0].Body != "body" {
		t.Fatalf("rows = %#v, want the index title over the file body", doc.Rows)
	}
}

// TestTreeParseReportsUnrecognizedFiles covers PR8: a basename that is not `<row ID>.md`
// is loud, and the row-ID grammar is today's, so a non-FT ID is simply a row.
func TestTreeParseReportsUnrecognizedFiles(t *testing.T) {
	tree := splitTree("**AB1 (LOW) — z.**\n", map[string]string{
		"AB1.md":   "**AB1 (LOW) — z.**\n",
		"FT7.txt":  "**FT7 (LOW) — x.**\n",
		"notes.md": "scratch\n",
	})
	doc, _, diagnostics := ParseDocument(tree, nil, true)
	want := []string{
		"roadmap/FT7.txt: unrecognized file under roadmap/; expected <row ID>.md",
		"roadmap/notes.md: unrecognized file under roadmap/; expected <row ID>.md",
	}
	if !reflect.DeepEqual(diagnostics, want) {
		t.Fatalf("diagnostics = %#v, want %#v", diagnostics, want)
	}
	if len(doc.Rows) != 1 || doc.Rows[0].ID != "AB1" {
		t.Fatalf("rows = %#v, want the AB1 row", doc.Rows)
	}
}

// TestTreeParseReportsDuplicateIndexID covers PR9: an ID on two index lines names the
// duplicate and its position, and only the first position keeps a row.
func TestTreeParseReportsDuplicateIndexID(t *testing.T) {
	const heading = "**FT7 (LOW) — x.**"
	tree := splitTree(heading+"\n\n"+heading+"\n", map[string]string{"FT7.md": heading + "\n"})
	doc, _, diagnostics := ParseDocument(tree, nil, true)
	want := []string{"ROADMAP.md: duplicate row FT7 at line 3"}
	if !reflect.DeepEqual(diagnostics, want) {
		t.Fatalf("diagnostics = %#v, want %#v", diagnostics, want)
	}
	if len(doc.Rows) != 1 {
		t.Fatalf("rows = %#v, want one", doc.Rows)
	}
}

// TestTreeParseDropsWrappedHeadingAndKeepsMissingOwnerRow covers PR10: a heading whose
// close is on the next physical line yields no row at all, while the index-side
// missing-owner fault beside it keeps its row and its place in rows_total.
func TestTreeParseDropsWrappedHeadingAndKeepsMissingOwnerRow(t *testing.T) {
	index := "**FT7 (LOW) — a\nlong title.**\n\n**FT8 (LOW) — y.**\n"
	doc, _, diagnostics := ParseDocument(indexTree(index), nil, true)
	want := []string{
		"ROADMAP.md: wrapped heading at line 1; a row heading is one physical line",
		"roadmap/FT8.md: missing detail owner for ROADMAP.md row FT8",
	}
	if !reflect.DeepEqual(diagnostics, want) {
		t.Fatalf("diagnostics = %#v, want %#v", diagnostics, want)
	}
	if len(doc.Rows) != 1 || doc.Rows[0].ID != "FT8" || doc.Rows[0].Body != "" {
		t.Fatalf("rows = %#v, want only the FT8 row with an empty body", doc.Rows)
	}
}

// TestTreeParseReportsUnreadRowFile covers story 44: a row file the classifier could not
// read is named, so its unread ledger is never silently counted as zero.
func TestTreeParseReportsUnreadRowFile(t *testing.T) {
	tree := indexTree("**FT7 (LOW) — x.**\n")
	tree.DirState = bounds.StateParsed
	tree.Files = []RowFile{{Name: "FT7.md", State: bounds.StateWrongType, Reason: "not a regular file: d---------"}}
	doc, _, diagnostics := ParseDocument(tree, nil, true)
	want := []string{"roadmap/FT7.md: wrong-type detail file: not a regular file: d---------"}
	if !reflect.DeepEqual(diagnostics, want) {
		t.Fatalf("diagnostics = %#v, want %#v", diagnostics, want)
	}
	if len(doc.Rows) != 1 || doc.Rows[0].Body != "" {
		t.Fatalf("rows = %#v, want the row kept with an unread body", doc.Rows)
	}
}

// TestTreeParseAbsentBoardIsQuiet covers PR16 and PR19 at the parse: a repository with
// neither the index nor the directory is a repository without a board, not a broken one.
func TestTreeParseAbsentBoardIsQuiet(t *testing.T) {
	doc, failures, diagnostics := ParseDocument(Tree{Index: bounds.Classified{State: bounds.StateAbsent}}, nil, true)
	if len(doc.Rows) != 0 || len(failures) != 0 || len(diagnostics) != 0 {
		t.Fatalf("absent board = rows %#v failures %#v diagnostics %#v", doc.Rows, failures, diagnostics)
	}
}

// TestIdeaOwnerValidatesThroughTheSplitTree covers PR16: `bench idea --owner` refuses an
// owner whose detail file is missing with the untrusted-tree line, and appends once the
// tree is whole.
func TestIdeaOwnerValidatesThroughTheSplitTree(t *testing.T) {
	const heading = "**FT7 (LOW) — x.**"
	writeBoard := func(t *testing.T, root string, rowFile bool) {
		t.Helper()
		if err := os.WriteFile(roadmapPath(t, root), []byte(heading+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if !rowFile {
			return
		}
		if err := os.MkdirAll(filepath.Join(root, RoadmapDir), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, RoadmapDir, "FT7.md"), []byte(heading+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	root := newRepo(t)
	writeBoard(t, root, false)
	out, code := IdeaCommand([]string{"--owner", "FT7", "--incident", "signal", "text"})
	if code != 1 || !strings.Contains(out, "ROADMAP.md is structurally untrusted") {
		t.Fatalf("missing detail owner = %q/%d, want exit 1 naming the untrusted tree", out, code)
	}
	if _, err := os.Stat(ideasPath(t, root)); err == nil {
		t.Fatal("refused owner still appended to the inbox")
	}

	root = newRepo(t)
	writeBoard(t, root, true)
	if out, code := IdeaCommand([]string{"--owner", "FT7", "--incident", "signal", "text"}); code != 0 {
		t.Fatalf("whole tree = %q/%d, want exit 0", out, code)
	}
	got, err := os.ReadFile(ideasPath(t, root))
	if err != nil {
		t.Fatal(err)
	}
	if !regexp.MustCompile(`(?m)^- [0-9-]{10}  text \[occurrence:FT7/signal\]$`).Match(got) {
		t.Fatalf("appended entry = %q", got)
	}
}
