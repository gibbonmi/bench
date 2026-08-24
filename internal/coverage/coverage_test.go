// Tests for the parse and Check verdicts, plus the shared spec fixtures the suite maps rows from.
package coverage

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// stories declares three stories. Because most fixtures reference only story 1, it
// carries the reasoned exception line the checker requires for the other two. The
// unreferenced-story rule has its own red cases below on bareStories.
const stories = "## User stories\n1. As a, I want b, so c.\n2. As d, I want e, so f.\n3. As g, I want h, so i.\n\nNot covered: story 2 — fixture\nNot covered: story 3 — fixture\n"
const bareStories = "## User stories\n1. As a, I want b, so c.\n2. As d, I want e, so f.\n3. As g, I want h, so i.\n4. As j, I want k, so l.\n5. As m, I want n, so o.\n"

// hdrReducedID and hdrReduced are the only two accepted headers. hdrReducedID opts
// into per-row IDs and has five cells; hdrReduced does not, and has four cells. The
// two widths no longer collide, so no accepted header shares a width with another.
const hdrReducedID = "| row | story | behavior | seam | why it catches the failure |\n|---|---|---|---|---|\n"
const hdrReduced = "| story | behavior | seam | why it catches the failure |\n|---|---|---|---|\n"

// hdrRedSignal and hdrRedSignalID are the two retired headers: they carry the
// `red signal` column no schema names anymore. The tests use both only to prove
// Check rejects them outright, naming them as missing the canonical header rather
// than parsing them.
const hdrRedSignal = "| story | behavior | seam | red signal | why it catches the failure |\n|---|---|---|---|---|\n"
const hdrRedSignalID = "| row | story | behavior | seam | red signal | why it catches the failure |\n|---|---|---|---|---|---|\n"

// bareStoriesExcused declares five stories and excuses the fifth, so a fixture
// covering exactly the four-story bound is clean rather than orphaning a story.
const bareStoriesExcused = bareStories + "\nNot covered: story 5 — out of scope\n"

// mapShape is one accepted header paired with the writer that places named cells
// under it. The two shapes let one violation table run against both schemas. A check
// reading a literal offset passes under the header it was written for and misreads
// under the other, so it cannot hide behind a single-header fixture.
type mapShape struct {
	name  string
	hdr   string
	optIn bool
	row   func(id, story, behavior, seam string) string
}

var mapShapes = []mapShape{
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
	p := spec("# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReduced +
		"| 2–3 | does x \\| y | cli seam | catches z |\n" +
		"| edge of 1 | edge case | gate | catches w |\n")
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

// Downstream consumers match every validation phrasing by substring; this test pins
// each here.
func TestCheck(t *testing.T) {
	valid := spec("# v\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReduced + "| 1, 2–3 | b | s | w |\n")
	if v := Check(valid); len(v) != 0 {
		t.Errorf("valid map violations = %v", v)
	}
	// Historical opts out of validation.
	if v := Check(spec("# h\n<!-- coverage-map: historical -->\n### Acceptance coverage map\n|bad|\n")); v != nil {
		t.Errorf("historical Check = %v, want nil", v)
	}
	cases := []struct{ body, want string }{
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n| a | b |\n|---|---|\n| 1 | x |\n", "coverage map missing the canonical header"},
		// Check rejects both retired red-signal headers outright, naming them as missing
		// the canonical header rather than parsed. One row covers each accepted width they
		// used to share: five cells and six cells.
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdrRedSignal + "| 1 | b | s | r | w |\n",
			"coverage map missing the canonical header"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdrRedSignalID + "| AB1 | 1 | b | s | r | w |\n",
			"coverage map missing the canonical header"},
		// A declared story no row references needs a row or a reasoned exception line;
		// an explicit but empty reason is its own distinct violation.
		{"# b\n\n" + bareStories + "\nNot covered: story 5 — \n\n### Acceptance coverage map\n" + hdrReduced + "| 1-4 | b | s | w |\n",
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
		// The empty-cell case asserts the fourth cell, `seam`, rather than the last cell.
		// The last cell — `why it catches the failure` — is identical in both schemas, so
		// it could not catch an offset slip.
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReducedID + "| AB1 | 1 | b |  | w |\n",
			"coverage map row 1 has an empty 'seam' cell"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReduced + "| 1 | b |  | w |\n",
			"coverage map row 1 has an empty 'seam' cell"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReducedID + "| AB1 | 1 | prints x; removes y | s | w |\n",
			"coverage map row 1 behavior states more than one predicate (';' outside backticks); split the row"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReduced + "| 1 | prints x; removes y | s | w |\n",
			"coverage map row 1 behavior states more than one predicate (';' outside backticks); split the row"},
		// Each case renames only the header's last cell, at every accepted width. Check
		// refuses every one rather than parsing it under a guess. A schema resolved by
		// cell count, when names do not match, would accept all three.
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
		"# b\n\n" + bareStories + "\nNot covered: story 5 — out of scope\n\n### Acceptance coverage map\n" + hdrReduced + "| 1-4 | runs `a; b` | s | w |\n",
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
// and moves on. Rows walks every data row regardless, so the projection resolves
// fields that run past the row's last cell. Those fields read as empty, and the
// cells the row does carry still land in their named slots. This test asserts the
// projected triple, rather than checking only for the absence of a panic, to pin
// down which cells those are.
func TestRowsProjectsShortRowUnderBothSchemas(t *testing.T) {
	nonOptIn := spec("# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReduced +
		"| 1 | does x |\n" +
		"| 2 | does y | cli seam |\n")
	if State(nonOptIn) != "mapped" {
		t.Fatalf("non-opt-in state = %q, want mapped", State(nonOptIn))
	}
	// Row 1 carries a behavior but no seam; row 2 carries both.
	wantNonOptIn := [][]string{{"1", "does x", ""}, {"2", "does y", "cli seam"}}
	if got := Rows(nonOptIn); !reflect.DeepEqual(got, wantNonOptIn) {
		t.Errorf("Rows = %v, want %v", got, wantNonOptIn)
	}

	optIn := spec("# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReducedID +
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

// TestCheckOptIn covers the 5-cell opt-in canonical header on its own. A valid
// opt-in map has no violations, and its Rows and State behave exactly as a
// non-opt-in map's.
func TestCheckOptIn(t *testing.T) {
	p := spec("# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReducedID +
		"| AB1 | 1 | does x | cli seam | catches z |\n" +
		"| CD2 | 2 | does y | gate | catches w |\n")
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
// 5-cell opt-in map on disk. It checks the opt-in verdict, the ordered row IDs, and
// the absence of violations.
func TestParseSpecOptIn(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "spec.md")
	body := "# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReducedID +
		"| AB1 | 1 | does x | cli seam | catches z |\n" +
		"| CD2 | 2 | does y | gate | catches w |\n"
	mustWrite(t, path, body)

	optIn, ids, violations, err := ParseSpec(path)
	if err != nil {
		t.Fatalf("ParseSpec: %v", err)
	}
	if !optIn {
		t.Error("optIn = false, want true for a 5-cell opt-in map")
	}
	if want := []string{"AB1", "CD2"}; !reflect.DeepEqual(ids, want) {
		t.Errorf("ids = %v, want %v", ids, want)
	}
	if len(violations) != 0 {
		t.Errorf("violations = %v, want none", violations)
	}

	// A non-opt-in reduced map reports not opted in, with nil IDs.
	nonOptInPath := filepath.Join(dir, "non-opt-in.md")
	mustWrite(t, nonOptInPath, "# t\n\n"+stories+"\n### Acceptance coverage map\n"+hdrReduced+"| 1 | b | s | w |\n")
	optIn, ids, _, err = ParseSpec(nonOptInPath)
	if err != nil {
		t.Fatalf("ParseSpec: %v", err)
	}
	if optIn {
		t.Error("optIn = true, want false for a 4-cell non-opt-in map")
	}
	if ids != nil {
		t.Errorf("ids = %v, want nil for a non-opt-in map", ids)
	}
}
