package toon

import "testing"

// The adapter's spec-TOON quoting and escaping is pinned here at the pure-function seam:
// one row per special-value trigger class the library quotes, plus the carried-over
// cases (leading/trailing space, comma, inner quote — now backslash — and newline —
// now \n). A one-field table isolates the cell bytes; the header wrapping is exercised
// by TestTable below.
func TestTableCellEscaping(t *testing.T) {
	cases := []struct{ in, wantCell string }{
		// bare — no trigger
		{"plain", "plain"},
		{"em—dash-and-hyphen", "em—dash-and-hyphen"},
		// spec quoting-trigger classes
		{"", `""`},
		{"true", `"true"`},
		{"false", `"false"`},
		{"null", `"null"`},
		{"3", `"3"`},
		{"-1", `"-1"`},
		{"1.5", `"1.5"`},
		{"a:b", `"a:b"`},
		{`a\b`, `"a\\b"`},
		{"[x]", `"[x]"`},
		{"{y}", `"{y}"`},
		{"-lead", `"-lead"`},
		// carried-over cases
		{" leading", `" leading"`},
		{"trailing ", `"trailing "`},
		{" padded ", `" padded "`},
		{"\ttab-lead", `"\ttab-lead"`}, // control char: raw tab → \t escape
		{"has,comma", `"has,comma"`},
		{`has"quote`, `"has\"quote"`},
		{"line\nbreak", `"line\nbreak"`},
		{`a, "b"`, `"a, \"b\""`},
	}
	for _, c := range cases {
		want := "t[1]{v}:\n  " + c.wantCell + "\n"
		got, err := Table("t", []string{"v"}, [][]string{{c.in}})
		if err != nil {
			t.Errorf("Table cell %q errored: %v", c.in, err)
		}
		if got != want {
			t.Errorf("Table cell %q = %q, want %q", c.in, got, want)
		}
	}
}

// A control byte spec-TOON cannot represent (anything below 0x20 that is not the
// escapable tab/newline/return) has no valid cell form, so the library refuses and the
// adapter returns an error instead of panicking or forging a lossy block. A form-feed
// is the concrete carrier — it can ride in a git path through `bench diff`.
func TestTableUnrepresentableCellErrors(t *testing.T) {
	for _, in := range []string{"a\x0cb", "esc\x1b", "\x00nul"} {
		got, err := Table("t", []string{"v"}, [][]string{{in}})
		if err == nil {
			t.Errorf("Table(%q) = %q, want an error for an unrepresentable cell", in, got)
		}
		if got != "" {
			t.Errorf("Table(%q) returned block %q alongside the error; want empty", in, got)
		}
	}
}

// Representable is the one predicate row-filtering callers consult instead of failing
// a whole table on one cell, so it must agree with the encoder's actual refusal
// behavior byte-for-byte. This pin sweeps every C0 control, DEL, and a spread of
// ordinary cells, asserting predicate-true == Table-renders: if the library's refusal
// rule ever moves (say, it starts refusing DEL), this goes red instead of the two
// rules drifting apart silently.
func TestRepresentableMatchesEncoder(t *testing.T) {
	probes := []string{"plain", "em—dash", "🚀 emoji", `quo"te`, "has,comma", "del\x7fbyte"}
	for b := byte(0); b < 0x20; b++ {
		probes = append(probes, "ctl"+string([]byte{b})+"byte")
	}
	for _, in := range probes {
		_, err := Table("t", []string{"v"}, [][]string{{in}})
		if got, want := Representable(in), err == nil; got != want {
			t.Errorf("Representable(%q) = %v but Table error = %v — predicate and encoder disagree", in, got, err)
		}
	}
}

func TestTable(t *testing.T) {
	// Empty input is the definitive empty table with its schema, not a blank line.
	if got, _ := Table("learnings", []string{"date", "title"}, nil); got != "learnings[0]{date,title}:\n" {
		t.Errorf("empty Table = %q", got)
	}
	rows := [][]string{
		{"2026-01-01", "first"},
		{"2026-02-02", `a, "b"`},
	}
	want := "learnings[2]{date,title}:\n" +
		"  2026-01-01,first\n" +
		"  2026-02-02,\"a, \\\"b\\\"\"\n"
	if got, _ := Table("learnings", []string{"date", "title"}, rows); got != want {
		t.Errorf("Table =\n%q\nwant\n%q", got, want)
	}
}

// TableTyped keeps a genuine int cell bare while quoting a numeric-looking string in
// the same column — the mixed-column contract the maps `ticket` field relies on.
func TestTableTyped(t *testing.T) {
	// Empty typed table keeps its schema too.
	if got, _ := TableTyped("maps", []string{"map", "ticket"}, nil); got != "maps[0]{map,ticket}:\n" {
		t.Errorf("empty TableTyped = %q", got)
	}
	rows := [][]any{
		{"grill", 6},         // real ticket: int stays bare
		{"grill", "handoff"}, // close-readiness row: string, not a trigger, stays bare
	}
	want := "maps[2]{map,ticket}:\n" +
		"  grill,6\n" +
		"  grill,handoff\n"
	if got, _ := TableTyped("maps", []string{"map", "ticket"}, rows); got != want {
		t.Errorf("TableTyped =\n%q\nwant\n%q", got, want)
	}
	// A numeric-looking STRING in the same slot quotes — only a real int survives bare.
	strWant := "maps[1]{map,ticket}:\n  grill,\"6\"\n"
	if got, _ := TableTyped("maps", []string{"map", "ticket"}, [][]any{{"grill", "6"}}); got != strWant {
		t.Errorf("TableTyped numeric string = %q, want %q", got, strWant)
	}
}

func TestErrorfUsage(t *testing.T) {
	if got := Errorf("not in a git repository", "run inside a Bench-linked repo"); got != "error: not in a git repository — run inside a Bench-linked repo" {
		t.Errorf("Errorf = %q", got)
	}
	if got := Usage("bench learnings", "bogusarg"); got != "usage: bench learnings (unknown argument: bogusarg)" {
		t.Errorf("Usage = %q", got)
	}
}
