// Tests that pin one verdict across every accepted coverage-map header schema.
package coverage

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// TestCommandProjectsOneRowShapeAcrossSchemas pins the whole primary response as a
// literal for every accepted header. The test writes the same named cells under each
// of the two. This lets an agent seeding tasks from `bench coverage` read one row
// shape and never branch on the spec's schema. This includes the behavior cell,
// which every header carries and which names what to build. The test spells out the
// expectation here rather than recomputing it from the implementation's own column
// list. That list would follow a projection change silently instead of going red on
// it.
func TestCommandProjectsOneRowShapeAcrossSchemas(t *testing.T) {
	const want = "spec: spec.md\nstate: mapped\nrows[2]{story,behavior,seam}:\n" +
		"  \"1\",does x,cli seam\n" +
		"  \"2\",does y,gate\n" +
		"help[2]{cmd,why}:\n" +
		"  bench coverage --check spec.md,check coverage for stories 1\n" +
		"  bench coverage --check spec.md,check coverage for stories 2\n"
	rows := [][4]string{{"AB1", "1", "does x", "cli seam"}, {"AB2", "2", "does y", "gate"}}
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
// changes what the rows block projects, not what the AXI block offers. A violating
// map still earns exactly one `retry after repairing coverage map` action, whichever
// schema it was written under. TestCommandProjectsOneRowShapeAcrossSchemas above
// asserts the clean per-row case literally.
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

// TestSchemaSelectionByHeaderNames drives every accepted header through State and
// ParseSpec, the two verdicts a caller reads. Each maps, and each reports the opt-in
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

// TestFiveCellHeadersSelectSchemaByName is retired. It drove byte-identical rows
// through the two legacy five-cell headers: the no-row-ID form and its row-ID
// variant, which differed only in whether they named `red signal`. It proved the
// parser resolves a schema by cell *names*, never by cell count. The legacy pair's
// removal left widths 5, hdrReducedID, and 4, hdrReduced, as the only accepted
// forms, so no two accepted headers share a width anymore. The confusion class
// this test guarded is now structurally impossible rather than merely untested.
// This suite does not invent a synthetic same-width header to keep it alive.

// TestViolationsAreIdenticalAcrossSchemas runs one violation table against every
// accepted header from one set of named cells. A cut column must not cut a check,
// so each case asserts the *same* violation strings under every schema. A check
// reading a literal offset would report a different message, or none, under one of
// the two. Behaviors carry an escaped pipe, since a behavior legitimately contains
// `|` and the parser's split sentinel has to survive the schema change.
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
		// The check must still reject a malformed row ID at a later row, not just row 0.
		// The well-formed first row proves the check runs per row rather than only at the
		// start of the loop.
		{name: "malformed row id at a later row", stories: stories, optInOnly: true,
			rows: [][4]string{{"AB1", "1", behavior, "cli seam"}, {"not-an-id", "2", behavior, "gate"}},
			want: "coverage map row 2 has a malformed row id 'not-an-id'"},
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

// TestMixedTagRowIDs pins the mixed-tag refusal. A coverage map declares one row-ID
// tag, and the preflight membership check scopes to the tag of the first row, so a
// foreign-tag row escapes that check. CE12 asserts the two-tag map is a violation and
// CE13 asserts the message names every tag it found, in a stable order. The one-tag
// arm is the anti-regression: the ordinary map must stay silent.
func TestMixedTagRowIDs(t *testing.T) {
	const want = "coverage map row ids carry more than one tag (AB, CD); a map declares one tag"
	for _, s := range mapShapes {
		if !s.optIn {
			continue
		}
		t.Run(s.name+"/two tags", func(t *testing.T) {
			rows := [][4]string{{"AB1", "1", "does x", "cli seam"}, {"CD2", "2", "does y", "gate"}}
			got := Check(spec(s.body(stories, rows)))
			if !reflect.DeepEqual(got, []string{want}) {
				t.Errorf("violations = %v, want exactly [%q]", got, want)
			}
		})
		t.Run(s.name+"/one tag", func(t *testing.T) {
			rows := [][4]string{{"AB1", "1", "does x", "cli seam"}, {"AB2", "2", "does y", "gate"}}
			if got := Check(spec(s.body(stories, rows))); len(got) != 0 {
				t.Errorf("violations = %v, want none", got)
			}
		})
	}
}
