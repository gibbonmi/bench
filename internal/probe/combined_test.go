package probe

import (
	"path/filepath"
	"strings"
	"testing"
)

// PS13, PS17: two --swap pairs and one --omit make one mutant that holds all three edits,
// and the cell lists the swaps first even when the argv leads with the omission.
func TestCombinedEditsMakeOneMutant(t *testing.T) {
	f := newFixture(t)
	comment := "// Clamp keeps n at or above zero.\n"
	want := strings.NewReplacer(comment, "", "n < 0", "n > 0", "return n", "return -n").Replace(clampSource)
	requireMutantBites(t, f, want, "swap+swap+omit",
		"--omit", comment, "--swap", "n < 0", "--with", "n > 0", "--swap", "return n", "--with", "return -n")
}

// PS14, PS15, PS23: each combined old string matches exactly once in the start bytes, not
// in the result of an earlier edit, and every pair whose --with equals its --swap refuses.
func TestCombinedEditsMatchOnce(t *testing.T) {
	t.Run("start bytes", func(t *testing.T) {
		f := newFixture(t)
		// The first edit adds a second "return n", so only a match over the start bytes
		// finds the second old string exactly once.
		want := strings.NewReplacer("return 0", "return n", "return n", "return 1").Replace(clampSource)
		requireMutantBites(t, f, want, "swap+swap", "--swap", "return 0", "--with", "return n", "--swap", "return n", "--with", "return 1")
	})
	t.Run("ambiguous", func(t *testing.T) {
		f := refusalFixture(t)
		out, code := runProbe(t, "clamp.go", "--swap", "n < 0", "--with", "n > 0", "--omit", "return", "--package", "./")
		requireRefused(t, f, out, code, "error: probe mutation ambiguous — the old string matches 2 times, want exactly 1\n", 1)
	})
	t.Run("empty pair", func(t *testing.T) {
		f := refusalFixture(t)
		out, code := runProbe(t, "clamp.go", "--swap", "n < 0", "--with", "n > 0", "--swap", "return n", "--with", "return n", "--package", "./")
		requireRefused(t, f, out, code, "error: probe mutation empty — --with equals --swap\n", 1)
	})
}

// PS16: two edits whose byte ranges overlap refuse before any write.
func TestCombinedEditsRefuseOverlap(t *testing.T) {
	f := refusalFixture(t)
	out, code := runProbe(t, "clamp.go", "--swap", "n < 0", "--with", "n > 0", "--omit", "if n < 0 {", "--package", "./")
	requireRefused(t, f, out, code, "error: probe mutations overlap — two edits touch the same bytes\n", 1)
}

// PS18: a combined mutant takes no single-edit form, so each one is a usage refusal.
func TestCombinedEditsRefuseSingleForms(t *testing.T) {
	combined := []string{"clamp.go", "--swap", "n < 0", "--with", "n > 0", "--omit", "return 0"}
	for _, single := range [][]string{
		{"--with-file", "x"},
		{"--insert-before", "return n"},
		{"--unwrap", "ensureNonNegative(n)"},
		{"--omit-file", "x"},
		{"--delete"},
	} {
		t.Run(single[0], func(t *testing.T) {
			f := refusalFixture(t)
			argv := append(append(append([]string(nil), combined...), single...), "--package", "./")
			out, code := runProbe(t, argv...)
			requireRefused(t, f, out, code, grammar.Help+"\n", 2)
		})
	}
	// The one swap takes its text from the file, so no other partner rule refuses this argv.
	t.Run("--with-file without --with", func(t *testing.T) {
		f := refusalFixture(t)
		input := filepath.Join(scratchDir(t), "replacement")
		writeFixtureFile(t, input, "n > 0", 0o644)
		out, code := runProbe(t, "clamp.go", "--swap", "n < 0", "--with-file", input, "--omit", "return 0", "--package", "./")
		requireRefused(t, f, out, code, grammar.Help+"\n", 2)
	})
}
