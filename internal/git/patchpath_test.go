package git

import "testing"

// TestPatchPathDecoder pins the shared patch path syntax: C-quoted fields decode to their raw
// bytes, a malformed escape keeps its spelling, side headers strip their display prefix and
// read /dev/null as absent, and a `diff --git` header splits unquoted fields only where both
// halves name one path.
func TestPatchPathDecoder(t *testing.T) {
	for _, test := range []struct{ field, want string }{
		{`plain name`, "plain name"},
		{`"caf\303\251 \"q\" back\\slash\ttab\nline\rret"`, "café \"q\" back\\slash\ttab\nline\rret"},
		{`"bad\q"`, `"bad\q"`},
		{`"short\30"`, `"short\30"`},
	} {
		if got := UnquotePath(test.field); got != test.want {
			t.Errorf("UnquotePath(%q) = %q, want %q", test.field, got, test.want)
		}
	}
	for _, test := range []struct{ field, want string }{
		{"a/dir/file", "dir/file"},
		{"b/space name.txt\t\n", "space name.txt"},
		// Git ends a side header whose path holds a space with a tab, so a trailing space
		// before that tab belongs to the path.
		{"a/trail \t\n", "trail "},
		{"b/dir/trail \t", "dir/trail "},
		{"\"b/q\\\"uote \"\t\n", "q\"uote "},
		{`"b/tab\there"`, "tab\there"},
		{"/dev/null", ""},
		{"c/no-prefix", "c/no-prefix"},
	} {
		if got := PatchSidePath(test.field); got != test.want {
			t.Errorf("PatchSidePath(%q) = %q, want %q", test.field, got, test.want)
		}
	}
	for _, test := range []struct {
		fields, base, tip string
		ok                bool
	}{
		{"a/x b/x", "x", "x", true},
		{"a/a b/c b/a b/c", "a b/c", "a b/c", true},
		{`"a/new\nline" "b/new\nline"`, "new\nline", "new\nline", true},
		{`"a/o\"ld" b/new name`, "o\"ld", "new name", true},
		{`a/old name "b/n\tew"`, "old name", "n\tew", true},
		{"a/old b/new", "", "", false},
		{"a/x b/xy", "", "", false},
		{`"a/unclosed b/x`, "", "", false},
	} {
		base, tip, ok := PatchGitPaths(test.fields)
		if base != test.base || tip != test.tip || ok != test.ok {
			t.Errorf("PatchGitPaths(%q) = (%q, %q, %v), want (%q, %q, %v)", test.fields, base, tip, ok, test.base, test.tip, test.ok)
		}
	}
}
