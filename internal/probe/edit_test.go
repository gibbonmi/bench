package probe

import (
	"path/filepath"
	"strings"
	"testing"
)

// PS1: a --with-file swap puts the exact multiline bytes of the file in place of the match.
func TestWithFileSwapsMultilineBytes(t *testing.T) {
	f := newFixture(t)
	old := "if n < 0 {\n\t\treturn 0\n\t}\n"
	replacement := "if n < 0 {\n\t\tn = -n\n\t}\n\t// widened\n"
	input := filepath.Join(scratchDir(t), "replacement")
	writeFixtureFile(t, input, replacement, 0o644)
	requireMutantBites(t, f, strings.Replace(clampSource, old, replacement, 1), "swap", "--swap", old, "--with-file", input)
}

// PS3, PS4, PS5: the insertion text goes immediately before the one anchor match, from
// --with or from --with-file, and an anchor that matches twice refuses as ambiguous.
func TestInsertBeforePutsTextBeforeTheAnchor(t *testing.T) {
	anchor := "return n\n"
	t.Run("with", func(t *testing.T) {
		f := newFixture(t)
		text := "n++\n\t"
		requireMutantBites(t, f, strings.Replace(clampSource, anchor, text+anchor, 1), "insert", "--insert-before", anchor, "--with", text)
	})
	t.Run("with-file", func(t *testing.T) {
		f := newFixture(t)
		text := "n++\n\tn *= 2\n\t"
		input := filepath.Join(scratchDir(t), "insertion")
		writeFixtureFile(t, input, text, 0o644)
		requireMutantBites(t, f, strings.Replace(clampSource, anchor, text+anchor, 1), "insert", "--insert-before", anchor, "--with-file", input)
	})
	t.Run("ambiguous anchor", func(t *testing.T) {
		f := refusalFixture(t)
		out, code := runProbe(t, "clamp.go", "--insert-before", "return", "--with", "n++\n\t", "--package", "./")
		requireRefused(t, f, out, code, "error: probe mutation ambiguous — the old string matches 2 times, want exactly 1\n", 1)
	})
}

// PS6, PS7, PS8: --occurrence selects one of several matches, an index past the matches
// refuses with the count, and a value that is not a positive decimal is a usage refusal.
func TestOccurrenceSelectsOneMatch(t *testing.T) {
	t.Run("second match", func(t *testing.T) {
		f := newFixture(t)
		want := strings.Replace(clampSource, "return n", "return - n", 1)
		requireMutantBites(t, f, want, "swap", "--swap", "return", "--with", "return -", "--occurrence", "2")
	})
	t.Run("out of range", func(t *testing.T) {
		f := refusalFixture(t)
		out, code := runProbe(t, "clamp.go", "--swap", "return", "--with", "x", "--occurrence", "3", "--package", "./")
		want := "error: probe mutation occurrence out of range — the old string matches 2 times, --occurrence is 3\n"
		requireRefused(t, f, out, code, want, 1)
	})
	for _, value := range []string{"0", "x", "-1"} {
		t.Run("usage "+value, func(t *testing.T) {
			f := refusalFixture(t)
			out, code := runProbe(t, "clamp.go", "--swap", "return", "--with", "x", "--occurrence", value, "--package", "./")
			requireRefused(t, f, out, code, grammar.Help+"\n", 2)
		})
	}
}

// PS19, PS20, PS21, PS22: each partner fault of the form matrix is a usage refusal, so no
// flag is dropped in silence.
func TestFormMatrixRefusesPartnerFaults(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
	}{
		{"with and with-file", []string{"--swap", "n < 0", "--with", "n > 0", "--with-file", "x"}},
		{"with-file without swap", []string{"--omit", "n < 0", "--with-file", "x"}},
		{"insert-before without text", []string{"--insert-before", "return n"}},
		{"occurrence with omit-file", []string{"--omit-file", "x", "--occurrence", "1"}},
		{"occurrence with delete", []string{"--delete", "--occurrence", "1"}},
		{"occurrence with combined", []string{"--omit", "n < 0", "--omit", "return n", "--occurrence", "1"}},
		{"delete with with", []string{"--delete", "--with", "x"}},
		{"delete with with-file", []string{"--delete", "--with-file", "x"}},
		{"swap without its with", []string{"--swap", "n < 0", "--with", "n > 0", "--swap", "return n"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := refusalFixture(t)
			out, code := runProbe(t, append(append([]string{"clamp.go"}, tc.args...), "--package", "./")...)
			requireRefused(t, f, out, code, grammar.Help+"\n", 2)
		})
	}
}

// PS24: an absent old string reads the substring miss under --occurrence and in a combined
// mutant, with the form's own mutation cell.
func TestZeroMatchesStayASubstringMiss(t *testing.T) {
	for _, tc := range []struct {
		name, mutation string
		args           []string
	}{
		{"occurrence", "swap", []string{"--swap", "absent", "--with", "x", "--occurrence", "2"}},
		{"combined", "swap+omit", []string{"--swap", "n < 0", "--with", "n > 0", "--omit", "absent"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := refusalFixture(t)
			argv := append(append([]string{"clamp.go"}, tc.args...), "--package", "./", "--run", "^TestClampNegative$")
			out, code := runProbe(t, argv...)
			want := probeRow(t, "invalid", "clamp.go", tc.mutation, "substring-miss", 0, "untouched") +
				selectionRow(t, "package", "./", "^TestClampNegative$", "", 0) + executionRow(t, 1, nil, 0)
			requireRefused(t, f, out, code, want, 1)
		})
	}
}
