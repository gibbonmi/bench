package coverage

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/toon"
)

// stories declares three stories and, because most fixtures reference only story 1,
// carries the reasoned exception line the checker requires for the other two; the
// unreferenced-story rule has its own red cases below on bareStories.
const stories = "## User stories\n1. As a, I want b, so c.\n2. As d, I want e, so f.\n3. As g, I want h, so i.\n\nNot covered: story 2 — fixture\nNot covered: story 3 — fixture\n"
const bareStories = "## User stories\n1. As a, I want b, so c.\n2. As d, I want e, so f.\n3. As g, I want h, so i.\n4. As j, I want k, so l.\n5. As m, I want n, so o.\n"
const hdr = "| story | behavior | seam | red signal | why it catches the failure |\n|---|---|---|---|---|\n"
const hdr6 = "| row | story | behavior | seam | red signal | why it catches the failure |\n|---|---|---|---|---|---|\n"

// The reduced headers: the same map with the `red signal` column cut. hdrReducedID
// is the five-cell opted-in form, and it is deliberately the same width as hdr —
// the two differ only in whether they name `red signal`, which is what makes a
// schema chosen by cell count observably wrong.
const hdrReducedID = "| row | story | behavior | seam | why it catches the failure |\n|---|---|---|---|---|\n"
const hdrReduced = "| story | behavior | seam | why it catches the failure |\n|---|---|---|---|\n"

// bareStoriesExcused declares five stories and excuses the fifth, so a fixture
// covering exactly the four-story bound is clean rather than orphaning a story.
const bareStoriesExcused = bareStories + "\nNot covered: story 5 — out of scope\n"

// mapShape is one accepted header paired with the writer that places named cells
// under it. The four shapes let one violation table run against every schema, so a
// check reading a literal offset — which passes under the header it was written
// for and misreads under another — cannot hide behind a single-header fixture.
type mapShape struct {
	name  string
	hdr   string
	optIn bool
	row   func(id, story, behavior, seam string) string
}

var mapShapes = []mapShape{
	{"legacy", hdr, false, func(_, story, behavior, seam string) string {
		return "| " + story + " | " + behavior + " | " + seam + " | r | w |\n"
	}},
	{"legacy-opt-in", hdr6, true, func(id, story, behavior, seam string) string {
		return "| " + id + " | " + story + " | " + behavior + " | " + seam + " | r | w |\n"
	}},
	{"reduced-opt-in", hdrReducedID, true, func(id, story, behavior, seam string) string {
		return "| " + id + " | " + story + " | " + behavior + " | " + seam + " | w |\n"
	}},
	{"reduced", hdrReduced, false, func(_, story, behavior, seam string) string {
		return "| " + story + " | " + behavior + " | " + seam + " | w |\n"
	}},
}

// body renders a whole spec under this shape's header from named cells.
func (s mapShape) body(storyBlock string, rows [][4]string) string {
	b := "# b\n\n" + storyBlock + "\n### Acceptance coverage map\n" + s.hdr
	for _, r := range rows {
		b += s.row(r[0], r[1], r[2], r[3])
	}
	return b
}

func spec(body string) parsed { return parse([]byte(body)) }

func TestStateAndRows(t *testing.T) {
	p := spec("# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdr +
		"| 2–3 | does x \\| y | cli seam | cmd fails, loudly | catches z |\n" +
		"| edge of 1 | edge case | gate | already covered | catches w |\n")
	if State(p) != "mapped" {
		t.Fatalf("state = %q", State(p))
	}
	want := [][]string{{"2–3", `does x \| y`, "cli seam"}, {"edge of 1", "edge case", "gate"}}
	if got := Rows(p); !reflect.DeepEqual(got, want) {
		t.Errorf("Rows = %v, want %v", got, want)
	}

	if State(spec("# h\n<!-- coverage-map: historical -->\n### Acceptance coverage map\n|bad|\n")) != "historical" {
		t.Error("historical state not detected")
	}
	if State(spec("# n\nprose only\n")) != "no-map" {
		t.Error("no-map state not detected")
	}

	// The marker opts a spec out under a reduced header too. The row is deliberately
	// too narrow for its header, so a Check that still returns nil proves validation
	// was skipped rather than merely satisfied.
	reducedHistorical := "# h\n<!-- coverage-map: historical -->\n### Acceptance coverage map\n" + hdrReduced + "| 1 | b | s |\n"
	if State(spec(reducedHistorical)) != "historical" {
		t.Errorf("reduced historical state = %q, want historical", State(spec(reducedHistorical)))
	}
	if v := Check(spec(reducedHistorical)); v != nil {
		t.Errorf("reduced historical Check = %v, want nil", v)
	}
}

// Every validation phrasing is matched by substring downstream; pin each here.
func TestCheck(t *testing.T) {
	valid := spec("# v\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + "| 1, 2–3 | b | s | r | w |\n")
	if v := Check(valid); len(v) != 0 {
		t.Errorf("valid map violations = %v", v)
	}
	// Historical opts out of validation.
	if v := Check(spec("# h\n<!-- coverage-map: historical -->\n### Acceptance coverage map\n|bad|\n")); v != nil {
		t.Errorf("historical Check = %v, want nil", v)
	}
	cases := []struct{ body, want string }{
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n| a | b |\n|---|---|\n| 1 | x |\n", "coverage map missing the canonical header"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr, "coverage map has no data rows"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + "| 1 | b | s | r |\n", "coverage map row 1 has 4 cells (want 5)"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + "| 1 | b |  | r | w |\n", "coverage map row 1 has an empty 'seam' cell"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + "| 9 | b | s | r | w |\n", "references story 9, which the spec does not declare (has: 1, 2, 3)"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + "| x | b | s | r | w |\n", "has an unrecognized story reference 'x'"},
		// A 6-cell map is opt-in; its rows carry a leading row-ID cell that Check
		// validates in addition to the legacy fields.
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr6 + "| AB1 | 1 | b | s | r | w |\n| AB1 | 1 | b | s | r | w |\n",
			"coverage map row 2 has a duplicate row id 'AB1' (first used at row 1)"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr6 + "|  | 1 | b | s | r | w |\n",
			"coverage map row 1 has an empty 'row' cell"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr6 + "| ab-1 | 1 | b | s | r | w |\n",
			"coverage map row 1 has a malformed row id 'ab-1'"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr6 + "| AB1 | 1 | b | s | r | w |\n| not-an-id | 2 | b | s | r | w |\n",
			"coverage map row 2 has a malformed row id 'not-an-id'"},
		// The opt-in header's own cell resolution: each of these names a field at a
		// different offset than the legacy header puts it at, so a check reading a
		// literal index — or a schema whose field list has slipped against its
		// columns — reports the neighbouring column's name instead of these.
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr6 + "| AB1 | 1 | b | s | r |\n",
			"coverage map row 1 has 5 cells (want 6)"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr6 + "| AB1 | 1 | b |  | r | w |\n",
			"coverage map row 1 has an empty 'seam' cell"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr6 + "| AB1 | 1 | b | s |  | w |\n",
			"coverage map row 1 has an empty 'red signal' cell"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr6 + "| AB1 | 1 | prints x; removes y | s | r | w |\n",
			"coverage map row 1 behavior states more than one predicate (';' outside backticks); split the row"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr6 + "| AB1 | 9 | b | s | r | w |\n",
			"references story 9, which the spec does not declare (has: 1, 2, 3)"},
		// A row spanning more than bounds.CoverageRowStories stories is an outcome family.
		{"# b\n\n" + bareStories + "\n### Acceptance coverage map\n" + hdr + "| 1-5 | b | s | r | w |\n",
			"coverage map row 1 references 5 stories (max 4); an outcome family is not one red-capable row"},
		// A `;` outside backticks in the behavior cell is a second predicate; one inside
		// a code span is not.
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + "| 1 | prints x; removes y | s | r | w |\n",
			"coverage map row 1 behavior states more than one predicate (';' outside backticks); split the row"},
		// A declared story no row references needs a row or a reasoned exception line.
		{"# b\n\n" + bareStories + "\n### Acceptance coverage map\n" + hdr + "| 1-4 | b | s | r | w |\n",
			"coverage map leaves story 5 unreferenced; add a row or a `Not covered: story 5 — <reason>` line"},
		{"# b\n\n" + bareStories + "\nNot covered: story 5 — \n\n### Acceptance coverage map\n" + hdr + "| 1-4 | b | s | r | w |\n",
			"story 5 is marked Not covered without a reason"},
		// The reduced schemas answer with their own width and their own field names.
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReducedID, "coverage map has no data rows"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReduced, "coverage map has no data rows"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReducedID + "| AB1 | 1 | b | s |\n",
			"coverage map row 1 has 4 cells (want 5)"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReducedID + "| AB1 | 1 | b | s | r | w |\n",
			"coverage map row 1 has 6 cells (want 5)"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReduced + "| 1 | b | s | r | w |\n",
			"coverage map row 1 has 5 cells (want 4)"},
		// Index 3 is the only offset where the reduced and legacy five-cell field
		// lists disagree — `seam` there, `red signal` in the legacy list — so the
		// empty-cell case asserts the fourth cell. The last cell would prove nothing:
		// index 4 is `why it catches the failure` under both.
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReducedID + "| AB1 | 1 | b |  | w |\n",
			"coverage map row 1 has an empty 'seam' cell"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReduced + "| 1 | b |  | w |\n",
			"coverage map row 1 has an empty 'seam' cell"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReducedID + "| AB1 | 1 | prints x; removes y | s | w |\n",
			"coverage map row 1 behavior states more than one predicate (';' outside backticks); split the row"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReduced + "| 1 | prints x; removes y | s | w |\n",
			"coverage map row 1 behavior states more than one predicate (';' outside backticks); split the row"},
		// A header renaming only its last cell, at each accepted width. Every one is
		// refused rather than parsed under a guess: a schema resolved by cell count
		// whenever the names do not match would accept all three.
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n| story | behavior | seam | outcome |\n|---|---|---|---|\n| 1 | b | s | w |\n",
			"coverage map missing the canonical header"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n| row | story | behavior | seam | outcome |\n|---|---|---|---|---|\n| AB1 | 1 | b | s | w |\n",
			"coverage map missing the canonical header"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n| row | story | behavior | seam | red signal | outcome |\n|---|---|---|---|---|---|\n| AB1 | 1 | b | s | r | w |\n",
			"coverage map missing the canonical header"},
		// A spec with no map and no historical marker is a violation in itself.
		{"# n\nprose only\n", "coverage map missing and spec is not marked historical"},
	}
	// The controls for the three rules: exactly four stories, a `;` inside backticks,
	// and a reasoned exception all pass.
	for _, body := range []string{
		"# b\n\n" + bareStories + "\nNot covered: story 5 — out of scope\n\n### Acceptance coverage map\n" + hdr + "| 1-4 | runs `a; b` | s | r | w |\n",
	} {
		if v := Check(spec(body)); len(v) != 0 {
			t.Errorf("control fixture violations = %v", v)
		}
	}
	for _, c := range cases {
		v := Check(spec(c.body))
		if len(v) == 0 || !strings.Contains(strings.Join(v, "\n"), c.want) {
			t.Errorf("Check violations %v do not contain %q", v, c.want)
		}
	}
}

// TestRowsProjectsShortRowUnderBothSchemas drives the projection over a data row
// carrying fewer cells than its header declares. Check refuses such a row on width
// and moves on, but Rows walks every data row regardless, so the projection resolves
// fields that run past the row's last cell: those read as empty and the cells the row
// does carry still land in their named slots. Asserting the projected triple rather
// than absence-of-panic pins which cells those are.
func TestRowsProjectsShortRowUnderBothSchemas(t *testing.T) {
	legacy := spec("# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdr +
		"| 1 | does x |\n" +
		"| 2 | does y | cli seam |\n")
	if State(legacy) != "mapped" {
		t.Fatalf("legacy state = %q, want mapped", State(legacy))
	}
	// Row 1 carries a behavior but no seam; row 2 carries both.
	wantLegacy := [][]string{{"1", "does x", ""}, {"2", "does y", "cli seam"}}
	if got := Rows(legacy); !reflect.DeepEqual(got, wantLegacy) {
		t.Errorf("Rows = %v, want %v", got, wantLegacy)
	}

	optIn := spec("# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdr6 +
		"| AB1 | 1 | does x |\n" +
		"| CD2 | 2 | does y | gate |\n")
	if State(optIn) != "mapped" {
		t.Fatalf("opt-in state = %q, want mapped", State(optIn))
	}
	// The leading row-ID cell shifts every field one place right, so the same
	// shortfall lands on the same two names.
	wantOptIn := [][]string{{"1", "does x", ""}, {"2", "does y", "gate"}}
	if got := Rows(optIn); !reflect.DeepEqual(got, wantOptIn) {
		t.Errorf("Rows = %v, want %v", got, wantOptIn)
	}
}

// TestCommandRendersShortRow drives the same shortfall through Command, the surface a
// caller actually reads: a short row still renders as a full three-column TOON record
// with its absent cells empty, alongside the repair action its width violation earns.
func TestCommandRendersShortRow(t *testing.T) {
	t.Chdir(t.TempDir())
	mustWrite(t, "spec.md", "# t\n\n"+stories+"\n### Acceptance coverage map\n"+hdr+"| 1 | does x |\n")

	out, code := Command([]string{"spec.md"})
	want := "spec: spec.md\nstate: mapped\nrows[1]{story,behavior,seam}:\n  \"1\",does x,\"\"\nhelp[1]{cmd,why}:\n  bench coverage --check spec.md,retry after repairing coverage map\n"
	if code != 0 || out != want {
		t.Fatalf("Command = (%d, %q), want (0, %q)", code, out, want)
	}
}

// TestCheckOptIn covers the 6-cell canonical header on its own: a valid opt-in map
// has no violations, and its Rows/State behave exactly as a legacy map's.
func TestCheckOptIn(t *testing.T) {
	p := spec("# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdr6 +
		"| AB1 | 1 | does x | cli seam | cmd fails | catches z |\n" +
		"| CD2 | 2 | does y | gate | already covered | catches w |\n")
	if v := Check(p); len(v) != 0 {
		t.Errorf("valid opt-in map violations = %v, want none", v)
	}
	if State(p) != "mapped" {
		t.Fatalf("state = %q, want mapped", State(p))
	}
	want := [][]string{{"1", "does x", "cli seam"}, {"2", "does y", "gate"}}
	if got := Rows(p); !reflect.DeepEqual(got, want) {
		t.Errorf("Rows = %v, want %v", got, want)
	}
}

// TestParseSpecOptIn drives ParseSpec, the package's exported entry point, over a
// 6-cell map on disk: the opt-in verdict, the ordered row IDs, and no violations.
func TestParseSpecOptIn(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "spec.md")
	body := "# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdr6 +
		"| AB1 | 1 | does x | cli seam | cmd fails | catches z |\n" +
		"| CD2 | 2 | does y | gate | already covered | catches w |\n"
	mustWrite(t, path, body)

	optIn, ids, violations, err := ParseSpec(path)
	if err != nil {
		t.Fatalf("ParseSpec: %v", err)
	}
	if !optIn {
		t.Error("optIn = false, want true for a 6-cell map")
	}
	if want := []string{"AB1", "CD2"}; !reflect.DeepEqual(ids, want) {
		t.Errorf("ids = %v, want %v", ids, want)
	}
	if len(violations) != 0 {
		t.Errorf("violations = %v, want none", violations)
	}

	// A legacy 5-cell map reports not opted in, with nil IDs.
	legacyPath := filepath.Join(dir, "legacy.md")
	mustWrite(t, legacyPath, "# t\n\n"+stories+"\n### Acceptance coverage map\n"+hdr+"| 1 | b | s | r | w |\n")
	optIn, ids, _, err = ParseSpec(legacyPath)
	if err != nil {
		t.Fatalf("ParseSpec: %v", err)
	}
	if optIn {
		t.Error("optIn = true, want false for a 5-cell map")
	}
	if ids != nil {
		t.Errorf("ids = %v, want nil for a legacy map", ids)
	}
}

// mustWrite creates path (and any parent dirs) with content under the current CWD.
func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("MkdirAll(%q): %v", dir, err)
		}
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", path, err)
	}
}

// TestCommand drives Command through its public (args) -> (output, exit code) interface
// only, per the spec's testing decision — never parse/Check directly.
func TestCommand(t *testing.T) {
	mapped := func(row string) string {
		return "# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + row
	}

	t.Run("readable path argument resolves and round-trips", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("| 1 | b | s | r | w |\n")
		mustWrite(t, "spec.md", body)

		out, code := Command([]string{"spec.md"})
		want := "spec: spec.md\nstate: mapped\nrows[1]{story,behavior,seam}:\n  \"1\",b,s\nhelp[1]{cmd,why}:\n  bench coverage --check spec.md,check coverage for stories 1\n"
		if out != want || code != 0 {
			t.Errorf("Command = (%q, %d), want (%q, 0)", out, code, want)
		}
	})

	t.Run("appends one check action per unchecked row and deduplicates exact templates", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("| 1 | b | s | unchecked | catches one |\n| 2 | c | t | already covered | catches two |\n| 1 | b | s | unchecked | catches one |\n")
		mustWrite(t, "path with spaces/spec.md", body)
		out, code := Command([]string{"path with spaces/spec.md"})
		want := "spec: path with spaces/spec.md\nstate: mapped\nrows[3]{story,behavior,seam}:\n  \"1\",b,s\n  \"2\",c,t\n  \"1\",b,s\nhelp[2]{cmd,why}:\n  bench coverage --check 'path with spaces/spec.md',check coverage for stories 1\n  bench coverage --check 'path with spaces/spec.md',check coverage for stories 2\n"
		if code != 0 || out != want {
			t.Fatalf("Command = (%d, %q), want (%d, %q)", code, out, 0, want)
		}
	})

	t.Run("repair prose does not alter unchecked classification", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("| 1 | b | s | repair is not evidence | catches one |\n")
		mustWrite(t, "spec.md", body)
		out, code := Command([]string{"spec.md"})
		if code != 0 || !strings.Contains(out, "bench coverage --check spec.md,check coverage for stories 1") {
			t.Fatalf("Command = (%d, %q), want unchecked check action", code, out)
		}
	})

	t.Run("malformed map gets repair retry without changing extraction exit", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("| 9 | b | s | red | catches one |\n")
		mustWrite(t, "spec.md", body)
		out, code := Command([]string{"spec.md"})
		want := "spec: spec.md\nstate: mapped\nrows[1]{story,behavior,seam}:\n  \"9\",b,s\nhelp[1]{cmd,why}:\n  bench coverage --check spec.md,retry after repairing coverage map\n"
		if code != 0 || out != want {
			t.Fatalf("Command = (%d, %q), want (0, %q)", code, out, want)
		}
	})

	// An unrecognized header has no descriptor of its own, so its rows project through
	// the legacy field order: a caller still gets the story/behavior/seam it seeds tasks
	// from, next to the repair action the unknown header earns. Dropping the rows would
	// leave the map's content unreachable until someone fixes the header by hand.
	t.Run("unrecognized header projects rows through the legacy order", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := "# t\n\n" + stories + "\n### Acceptance coverage map\n" +
			"| story | behavior | seam | outcome |\n|---|---|---|---|\n| 1 | b | s | w |\n"
		mustWrite(t, "spec.md", body)

		out, code := Command([]string{"spec.md"})
		want := "spec: spec.md\nstate: mapped\nrows[1]{story,behavior,seam}:\n  \"1\",b,s\nhelp[1]{cmd,why}:\n  bench coverage --check spec.md,retry after repairing coverage map\n"
		if code != 0 || out != want {
			t.Fatalf("Command = (%d, %q), want (0, %q)", code, out, want)
		}
	})

	t.Run("canonical mapped zero-row table is terminal", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("")
		mustWrite(t, "spec.md", body)
		out, code := Command([]string{"spec.md"})
		want := "spec: spec.md\nstate: mapped\nrows[0]{story,behavior,seam}:\nhelp[0]{cmd,why}:\n"
		if out != want || code != 0 {
			t.Fatalf("Command = (%d, %q), want (%d, %q)", code, out, 0, want)
		}
	})

	// The behavior cell is author-controlled prose that reaches the TOON encoder, so a
	// control byte in it is refused whole: the AXI error contract replaces the response
	// rather than a table with one lossy cell in it.
	t.Run("control-bearing behavior cell refuses the whole response", func(t *testing.T) {
		t.Chdir(t.TempDir())
		mustWrite(t, "spec.md", mapped("| 1 | b\x1bx | s | r | w |\n"))

		out, code := Command([]string{"spec.md"})
		want := "error: unrepresentable TOON cell — toon: unsupported control character U+001B in string\n"
		if code != 1 || out != want {
			t.Fatalf("Command = (%d, %q), want (1, %q)", code, out, want)
		}
	})

	// A comma is the row delimiter and a tab is an escapable control, so a behavior
	// carrying either is the case where an unquoted cell would silently split or
	// corrupt the record. Both stay one three-field row, quoted by the encoder.
	for _, tc := range []struct{ name, behavior, row string }{
		{"comma", "a,b", "  \"1\",\"a,b\",s\n"},
		{"tab", "does\tx", "  \"1\",\"does\\tx\",s\n"},
	} {
		t.Run("delimiter-bearing behavior renders one quoted row: "+tc.name, func(t *testing.T) {
			t.Chdir(t.TempDir())
			mustWrite(t, "spec.md", mapped("| 1 | "+tc.behavior+" | s | r | w |\n"))

			out, code := Command([]string{"spec.md"})
			want := "spec: spec.md\nstate: mapped\nrows[1]{story,behavior,seam}:\n" + tc.row +
				"help[1]{cmd,why}:\n  bench coverage --check spec.md,check coverage for stories 1\n"
			if code != 0 || out != want {
				t.Fatalf("Command = (%d, %q), want (0, %q)", code, out, want)
			}
		})
	}

	t.Run("separator-free slug resolves specs/<slug>/spec.md", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("| 1 | b2 | s2 | r2 | w2 |\n")
		mustWrite(t, "specs/foo/spec.md", body)

		out, code := Command([]string{"foo"})
		want := "spec: specs/foo/spec.md\nstate: mapped\nrows[1]{story,behavior,seam}:\n  \"1\",b2,s2\nhelp[1]{cmd,why}:\n  bench coverage --check specs/foo/spec.md,check coverage for stories 1\n"
		if out != want || code != 0 {
			t.Errorf("Command = (%q, %d), want (%q, 0)", out, code, want)
		}
	})

	t.Run("slug already ending .md resolves folder spec", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("| 1 | b3 | s3 | r3 | w3 |\n")
		mustWrite(t, "specs/bar/spec.md", body)

		out, code := Command([]string{"bar.md"})
		want := "spec: specs/bar/spec.md\nstate: mapped\nrows[1]{story,behavior,seam}:\n  \"1\",b3,s3\nhelp[1]{cmd,why}:\n  bench coverage --check specs/bar/spec.md,check coverage for stories 1\n"
		if out != want || code != 0 {
			t.Errorf("Command = (%q, %d), want (%q, 0)", out, code, want)
		}
	})

	t.Run("slug matching no file names both forms tried, exit 1", func(t *testing.T) {
		t.Chdir(t.TempDir())

		out, code := Command([]string{"missing"})
		if code != 1 || !strings.Contains(out, "spec not found: missing, specs/missing/spec.md") {
			t.Errorf("Command = (%q, %d), want exit 1 naming both 'missing' and 'specs/missing/spec.md' in the not-found message", out, code)
		}
	})

	t.Run("separator-bearing argument gets no fallback", func(t *testing.T) {
		t.Chdir(t.TempDir())

		out, code := Command([]string{"sub/missing.md"})
		if code != 1 || !strings.Contains(out, "sub/missing.md") || strings.Contains(out, "specs/") {
			t.Errorf("Command = (%q, %d), want exit 1 naming only sub/missing.md with no specs/ form", out, code)
		}
	})

	t.Run("slug shadowed by a same-named CWD file resolves path-first", func(t *testing.T) {
		t.Chdir(t.TempDir())
		cwdBody := mapped("| 1 | cwd | cwd-seam | cwd-red | cwd-why |\n")
		specsBody := mapped("| 5 | specs | specs-seam | specs-red | specs-why |\n")
		mustWrite(t, "foo", cwdBody)
		mustWrite(t, "specs/foo/spec.md", specsBody)

		out, code := Command([]string{"foo"})
		want := "spec: foo\nstate: mapped\nrows[1]{story,behavior,seam}:\n  \"1\",cwd,cwd-seam\nhelp[1]{cmd,why}:\n  bench coverage --check foo,check coverage for stories 1\n"
		if out != want || code != 0 {
			t.Errorf("Command = (%q, %d), want (%q, 0) — the CWD file should shadow the specs/ fallback", out, code, want)
		}
	})

	t.Run("present-but-unreadable file surfaces the read error, not not-found", func(t *testing.T) {
		if os.Geteuid() == 0 {
			capability.Capability(t, capability.Privilege, "root reads 0000-mode files; permission case unobservable")
		}
		t.Chdir(t.TempDir())
		mustWrite(t, "sub/locked.md", mapped("| 1 | b | s | r | w |\n"))
		if err := os.Chmod("sub/locked.md", 0o000); err != nil {
			t.Fatalf("Chmod: %v", err)
		}

		out, code := Command([]string{"sub/locked.md"})
		if code != 1 || !strings.Contains(out, "spec not readable:") || strings.Contains(out, "spec not found") {
			t.Errorf("Command = (%q, %d), want exit 1 with a 'spec not readable:' error, never 'spec not found'", out, code)
		}
	})

	t.Run("slug matching a directory falls back to specs/<slug>/spec.md", func(t *testing.T) {
		t.Chdir(t.TempDir())
		body := mapped("| 1 | b | s | r | w |\n")
		mustWrite(t, "dist/keep", "x") // a directory named like the slug
		mustWrite(t, "specs/dist/spec.md", body)

		out, code := Command([]string{"dist"})
		want := "spec: specs/dist/spec.md\nstate: mapped\nrows[1]{story,behavior,seam}:\n  \"1\",b,s\nhelp[1]{cmd,why}:\n  bench coverage --check specs/dist/spec.md,check coverage for stories 1\n"
		if out != want || code != 0 {
			t.Errorf("Command = (%q, %d), want (%q, 0) — a directory is not a spec candidate", out, code, want)
		}
	})

	t.Run("--check resolves a slug and reports violations under the resolved label", func(t *testing.T) {
		t.Chdir(t.TempDir())
		mustWrite(t, "specs/chk/spec.md", mapped("| 1 | | s | r | w |\n")) // empty behavior cell

		out, code := Command([]string{"--check", "chk"})
		if code != 1 || !strings.Contains(out, "error: specs/chk/spec.md coverage map row 1 has an empty 'behavior' cell") {
			t.Errorf("Command = (%q, %d), want exit 1 with the violation under the resolved folder label", out, code)
		}
	})

	t.Run("flag-shaped argument stays a usage error, exit 2", func(t *testing.T) {
		t.Chdir(t.TempDir())

		out, code := Command([]string{"--bogus"})
		if want := toon.Usage("bench coverage", "--bogus") + "\n"; out != want || code != 2 {
			t.Errorf("Command = (%q, %d), want (%q, 2)", out, code, want)
		}
	})

	t.Run("no argument reports missing/required, not unknown argument, exit 2", func(t *testing.T) {
		t.Chdir(t.TempDir())

		out, code := Command(nil)
		if code != 2 {
			t.Errorf("Command = (%q, %d), want exit 2", out, code)
		}
		if !strings.Contains(out, "required") && !strings.Contains(out, "missing") {
			t.Errorf("Command = %q, want it to say the argument is missing/required", out)
		}
		if strings.Contains(out, "unknown argument") {
			t.Errorf("Command = %q, must not use the unknown-argument template for a missing argument", out)
		}
	})
}

// TestCommandProjectsOneRowShapeAcrossSchemas pins the whole primary response as a
// literal for every accepted header. The same named cells are written under each of
// the four, so an agent seeding tasks from `bench coverage` reads one row shape and
// never branches on the spec's schema — including the behavior cell, which every
// header carries and which names what to build. The expectation is spelled out here
// rather than recomputed from the implementation's own column list, which would
// follow a projection change silently instead of going red on it.
func TestCommandProjectsOneRowShapeAcrossSchemas(t *testing.T) {
	const want = "spec: spec.md\nstate: mapped\nrows[2]{story,behavior,seam}:\n" +
		"  \"1\",does x,cli seam\n" +
		"  \"2\",does y,gate\n" +
		"help[2]{cmd,why}:\n" +
		"  bench coverage --check spec.md,check coverage for stories 1\n" +
		"  bench coverage --check spec.md,check coverage for stories 2\n"
	rows := [][4]string{{"AB1", "1", "does x", "cli seam"}, {"CD2", "2", "does y", "gate"}}
	for _, shape := range mapShapes {
		t.Run(shape.name, func(t *testing.T) {
			t.Chdir(t.TempDir())
			mustWrite(t, "spec.md", shape.body(stories, rows))

			out, code := Command([]string{"spec.md"})
			if code != 0 || out != want {
				t.Fatalf("Command = (%d, %q), want (0, %q)", code, out, want)
			}
		})
	}
}

// TestCommandActionListIsUnchangedByTheSchema drives the other half of the action
// list — a map with violations — under every accepted header. Cutting a column
// changes what the rows block projects, not what the AXI block offers: a violating
// map still earns exactly one `retry after repairing coverage map` action, whichever
// schema it was written under. The clean per-row case is asserted literally in
// TestCommandProjectsOneRowShapeAcrossSchemas above.
func TestCommandActionListIsUnchangedByTheSchema(t *testing.T) {
	const want = "spec: spec.md\nstate: mapped\nrows[1]{story,behavior,seam}:\n" +
		"  \"9\",does x,cli seam\n" +
		"help[1]{cmd,why}:\n" +
		"  bench coverage --check spec.md,retry after repairing coverage map\n"
	rows := [][4]string{{"AB1", "9", "does x", "cli seam"}} // story 9 is undeclared
	for _, shape := range mapShapes {
		t.Run(shape.name, func(t *testing.T) {
			t.Chdir(t.TempDir())
			mustWrite(t, "spec.md", shape.body(stories, rows))

			out, code := Command([]string{"spec.md"})
			if code != 0 || out != want {
				t.Fatalf("Command = (%d, %q), want (0, %q)", code, out, want)
			}
		})
	}
}

func TestCommandPreservesCheckedInPreDisclosureResponses(t *testing.T) {
	mapped := func(row string) string {
		return "# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + row
	}
	for _, tc := range []struct {
		name, fixture, body, help string
	}{
		{"mapped actionable", "pre-disclosure-mapped.stdout", mapped("| 1 | b | s | unchecked | catches one |\n"), "help[1]{cmd,why}:\n  bench coverage --check spec.md,check coverage for stories 1\n"},
		{"repairable malformed", "pre-disclosure-malformed.stdout", mapped("| 9 | b | s | red | catches one |\n"), "help[1]{cmd,why}:\n  bench coverage --check spec.md,retry after repairing coverage map\n"},
		{"mapped zero-row", "pre-disclosure-zero-row.stdout", mapped(""), "help[0]{cmd,why}:\n"},
		{"historical terminal", "pre-disclosure-historical.stdout", "# historical\n<!-- coverage-map: historical -->\n### Acceptance coverage map\n|bad|\n", "help[0]{cmd,why}:\n"},
		{"no-map terminal", "pre-disclosure-no-map.stdout", "# no map\n", "help[0]{cmd,why}:\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			primary, err := os.ReadFile(filepath.Join("testdata", tc.fixture))
			if err != nil {
				t.Fatal(err)
			}
			t.Chdir(t.TempDir())
			mustWrite(t, "spec.md", tc.body)
			out, code := Command([]string{"spec.md"})
			if code != 0 || out != string(primary)+tc.help {
				t.Fatalf("Command = (%d, %q), want checked-in primary plus exactly one help block", code, out)
			}
		})
	}
}

func TestCommandControlBearingSpecPathPreservesPrimaryAndHonestFallback(t *testing.T) {
	mapped := "# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + "| 1 | b | s | unchecked | catches one |\n"
	for _, control := range []string{"\t", "\n", "\r", "\x1b"} {
		t.Run("control", func(t *testing.T) {
			t.Chdir(t.TempDir())
			path := "control" + control + "spec.md"
			mustWrite(t, path, mapped)
			out, code := Command([]string{path})
			primary := "spec: " + path + "\nstate: mapped\nrows[1]{story,behavior,seam}:\n  \"1\",b,s\n"
			if code != 0 || !strings.HasPrefix(out, primary) {
				t.Fatalf("Command = (%d, %q), want primary response and exit 0", code, out)
			}
			if control == "\x1b" {
				if out != primary+"help[0]{cmd,why}:\n" {
					t.Fatalf("Command fallback = %q, want primary plus empty help", out)
				}
				return
			}
			if !strings.Contains(out, "help[1]{cmd,why}:") {
				t.Fatalf("Command = %q, want one action", out)
			}
			argv, err := axitest.RecoverHelpCommandArgv(out)
			if err != nil {
				t.Fatal(err)
			}
			want := []string{"bench", "coverage", "--check", path}
			if !slices.Equal(argv, want) {
				t.Fatalf("shell argv = %q, want %q", argv, want)
			}
		})
	}
}

func TestCommandAngleBracketSpecPathPreservesPrimaryAndHonestFallback(t *testing.T) {
	mapped := "# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + "| 1 | b | s | unchecked | catches one |\n"
	for _, marker := range []string{"<", ">"} {
		t.Run(marker, func(t *testing.T) {
			t.Chdir(t.TempDir())
			path := "angle" + marker + "spec.md"
			mustWrite(t, path, mapped)
			out, code := Command([]string{path})
			primary := "spec: " + path + "\nstate: mapped\nrows[1]{story,behavior,seam}:\n  \"1\",b,s\n"
			want := primary + "help[0]{cmd,why}:\n"
			if code != 0 || out != want {
				t.Fatalf("Command = (%d, %q), want checked primary plus honest empty help", code, out)
			}
		})
	}
}

// TestSchemaSelectionByHeaderNames drives every accepted header through State and
// ParseSpec, the two verdicts a caller reads: each maps, and each reports the opt-in
// and ordered row IDs its own header declares.
func TestSchemaSelectionByHeaderNames(t *testing.T) {
	dir := t.TempDir()
	for _, s := range mapShapes {
		t.Run(s.name, func(t *testing.T) {
			body := s.body(stories, [][4]string{{"AB1", "1", "does x", "cli seam"}})
			if got := State(spec(body)); got != "mapped" {
				t.Fatalf("State = %q, want mapped", got)
			}
			path := filepath.Join(dir, s.name+".md")
			mustWrite(t, path, body)
			optIn, ids, violations, err := ParseSpec(path)
			if err != nil {
				t.Fatalf("ParseSpec: %v", err)
			}
			if optIn != s.optIn {
				t.Errorf("optIn = %v, want %v", optIn, s.optIn)
			}
			var wantIDs []string
			if s.optIn {
				wantIDs = []string{"AB1"}
			}
			if !reflect.DeepEqual(ids, wantIDs) {
				t.Errorf("ids = %v, want %v", ids, wantIDs)
			}
			if len(violations) != 0 {
				t.Errorf("violations = %v, want none", violations)
			}
		})
	}
}

// TestFiveCellHeadersSelectSchemaByName drives byte-identical rows through the two
// five-cell headers. They differ only in whether they name `red signal`, so a schema
// resolved by cell count would project both alike; resolved by name, the same cells
// land in different fields and the projections differ.
func TestFiveCellHeadersSelectSchemaByName(t *testing.T) {
	const row = "| 1 | b | s | r | w |\n"
	legacy := spec("# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + row)
	reduced := spec("# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReducedID + row)

	// Legacy reads story/behavior/seam off cells 0, 1, 2. The reduced header's leading
	// row-ID cell shifts every field one place right, so the same cells project as a
	// different triple.
	if want := [][]string{{"1", "b", "s"}}; !reflect.DeepEqual(Rows(legacy), want) {
		t.Errorf("legacy Rows = %v, want %v", Rows(legacy), want)
	}
	if want := [][]string{{"b", "s", "r"}}; !reflect.DeepEqual(Rows(reduced), want) {
		t.Errorf("reduced Rows = %v, want %v", Rows(reduced), want)
	}
	if reflect.DeepEqual(Rows(legacy), Rows(reduced)) {
		t.Errorf("both five-cell headers projected %v; the schema was chosen by cell count, not cell names", Rows(legacy))
	}
}

// TestViolationsAreIdenticalAcrossSchemas runs one violation table against every
// accepted header from one set of named cells. Cutting a column must not cut a
// check, so each case asserts the *same* violation strings under every schema —
// a check reading a literal offset would report a different message, or none, under
// one of the four. Behaviors carry an escaped pipe, since a behavior legitimately
// contains `|` and the parser's split sentinel has to survive the schema change.
func TestViolationsAreIdenticalAcrossSchemas(t *testing.T) {
	const behavior = `does x \| y`
	cases := []struct {
		name      string
		stories   string
		rows      [][4]string
		want      string // "" means the fixture must be clean under every schema
		optInOnly bool
	}{
		{name: "clean row with an escaped pipe", stories: stories,
			rows: [][4]string{{"AB1", "1", behavior, "cli seam"}}},
		{name: "story reference at the fan-out bound", stories: bareStoriesExcused,
			rows: [][4]string{{"AB1", "1-4", behavior, "cli seam"}}},
		{name: "unrecognized story reference", stories: stories,
			rows: [][4]string{{"AB1", "x", behavior, "cli seam"}},
			want: "coverage map row 1 has an unrecognized story reference 'x'"},
		{name: "undeclared story", stories: stories,
			rows: [][4]string{{"AB1", "9", behavior, "cli seam"}},
			want: "coverage map row 1 references story 9, which the spec does not declare (has: 1, 2, 3)"},
		{name: "story zero", stories: stories,
			rows: [][4]string{{"AB1", "0", behavior, "cli seam"}},
			want: "coverage map row 1 references story 0, which is not a valid story number"},
		{name: "range with end before start", stories: bareStoriesExcused,
			rows: [][4]string{{"AB1", "3-1", behavior, "cli seam"}},
			want: "coverage map row 1 has a story range with end before start '3-1'"},
		{name: "fan-out past the bound", stories: bareStories,
			rows: [][4]string{{"AB1", "1-5", behavior, "cli seam"}},
			want: "coverage map row 1 references 5 stories (max 4); an outcome family is not one red-capable row"},
		{name: "orphan story", stories: bareStories,
			rows: [][4]string{{"AB1", "1-4", behavior, "cli seam"}},
			want: "coverage map leaves story 5 unreferenced; add a row or a `Not covered: story 5 — <reason>` line"},
		{name: "duplicate row id", stories: stories, optInOnly: true,
			rows: [][4]string{{"AB1", "1", behavior, "cli seam"}, {"AB1", "2", behavior, "gate"}},
			want: "coverage map row 2 has a duplicate row id 'AB1' (first used at row 1)"},
		{name: "malformed row id", stories: stories, optInOnly: true,
			rows: [][4]string{{"ab-1", "1", behavior, "cli seam"}},
			want: "coverage map row 1 has a malformed row id 'ab-1'"},
		{name: "empty row id", stories: stories, optInOnly: true,
			rows: [][4]string{{"", "1", behavior, "cli seam"}},
			want: "coverage map row 1 has an empty 'row' cell"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var first []string
			firstName := ""
			for _, s := range mapShapes {
				if c.optInOnly && !s.optIn {
					continue
				}
				got := Check(spec(s.body(c.stories, c.rows)))
				if c.want == "" {
					if len(got) != 0 {
						t.Errorf("%s: violations = %v, want none", s.name, got)
					}
				} else if !strings.Contains(strings.Join(got, "\n"), c.want) {
					t.Errorf("%s: violations %v do not contain %q", s.name, got, c.want)
				}
				if firstName == "" {
					first, firstName = got, s.name
					continue
				}
				if !reflect.DeepEqual(got, first) {
					t.Errorf("%s violations = %v, but %s reported %v; the same map must be refused identically under every schema", s.name, got, firstName, first)
				}
			}
		})
	}
}
