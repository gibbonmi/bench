package toon

import "testing"

// Escape is the escaping rule five commands compose; every branch is pinned here at
// the pure-function seam. The leading/trailing-whitespace rows re-home the shell
// "AXI TOON leading/trailing-space escaping contract" that used to source
// bin/bench-query.sh — coverage the CLI contracts never exercise directly.
func TestEscape(t *testing.T) {
	cases := []struct{ in, want string }{
		{"plain", "plain"},
		{"", ""},
		{" padded ", `" padded "`},
		{" leading", `" leading"`},
		{"trailing ", `"trailing "`},
		{"\ttab-lead", "\"\ttab-lead\""},
		{"has,comma", `"has,comma"`},
		{`has"quote`, `"has""quote"`},
		{"line\nbreak", "\"line\nbreak\""},
		{`a, "b"`, `"a, ""b"""`},
		{"em—dash-and-hyphen", "em—dash-and-hyphen"},
	}
	for _, c := range cases {
		if got := Escape(c.in); got != c.want {
			t.Errorf("Escape(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestTable(t *testing.T) {
	// Empty input is the definitive empty table, not a blank line.
	if got := Table("learnings", []string{"date", "title"}, nil); got != "learnings[0]{date,title}:\n" {
		t.Errorf("empty Table = %q", got)
	}
	rows := [][]string{
		{"2026-01-01", "first"},
		{"2026-02-02", `a, "b"`},
	}
	want := "learnings[2]{date,title}:\n" +
		"  2026-01-01,first\n" +
		"  2026-02-02,\"a, \"\"b\"\"\"\n"
	if got := Table("learnings", []string{"date", "title"}, rows); got != want {
		t.Errorf("Table =\n%q\nwant\n%q", got, want)
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
