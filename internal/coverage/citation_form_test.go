// Tests for the citation-form rule: a Go test name in a seam cell resolves only inside
// the citation form, so a name written anywhere else in the cell is a violation.
package coverage

import (
	"os"
	"strings"
	"testing"
)

// strayWant renders the one diagnostic the citation-form rule gives for name in row 1.
func strayWant(name string) string {
	return "coverage map row 1 names '" + name + "' outside the citation form; cite it as `<path>_test.go` (`" + name + "`)"
}

// TestSeamCellTestNameOutsideTheCitationForm grades the citation-form rule at CheckFiles:
// each shape that names a test outside the form reds, and the form itself keeps its
// resolution.
func TestSeamCellTestNameOutsideTheCitationForm(t *testing.T) {
	for _, c := range []struct{ name, seam string }{
		{"a bare name is a violation", "covered by TestMissing"},
		{"a backticked name with no path is a violation", "`TestMissing`"},
		{"a repeated name reports once", "TestMissing, then TestMissing again"},
		{"an unbackticked list name is a violation", "`internal/x/foo_test.go` (TestMissing)"},
		{"a marker inside a word does not exempt", "unplanned TestMissing"},
	} {
		t.Run(c.name, func(t *testing.T) {
			root, specPath := citedSpec(t, "s", c.seam, "")
			mentionSpecFile(t, root, "func TestPresent(t *testing.T) {}\n")

			want := strayWant("TestMissing")
			if v := checkFilesOf(t, specPath); len(v) != 1 || v[0] != want {
				t.Fatalf("CheckFiles = %#v, want exactly %q", v, want)
			}
		})
	}

	t.Run("a name beside an unpaired path is a violation", func(t *testing.T) {
		root, specPath := citedSpec(t, "s", "`internal/x/foo_test.go` TestPresent", "")
		mentionSpecFile(t, root, "func TestPresent(t *testing.T) {}\n")

		v := checkFilesOf(t, specPath)
		if len(v) != 2 || !hasViolation(v, strayWant("TestPresent")) || !hasViolation(v, "mentions 'internal/x/foo_test.go'") {
			t.Fatalf("CheckFiles = %#v, want the mention and the stray name TestPresent", v)
		}
	})

	t.Run("the citation form keeps its resolution", func(t *testing.T) {
		root, specPath := citedSpec(t, "s", "`internal/x/foo_test.go` (`TestPresent`, `TestPresent/a case`)", "")
		mentionSpecFile(t, root, "func TestPresent(t *testing.T) {\n\tt.Run(\"a case\", func(t *testing.T) {})\n}\n")

		if v := checkFilesOf(t, specPath); len(v) != 0 {
			t.Fatalf("CheckFiles = %#v, want no violation for the citation form", v)
		}
	})

	t.Run("prose that is no test name is not graded", func(t *testing.T) {
		_, specPath := citedSpec(t, "s", "Tests and Testing stay prose; retestX and TestlowerCase too", "")

		if v := checkFilesOf(t, specPath); len(v) != 0 {
			t.Fatalf("CheckFiles = %#v, want no violation for prose", v)
		}
	})

	t.Run("a why-cell name is not graded", func(t *testing.T) {
		_, specPath := citedCellSpec(t, "s", "review-owned: the Standards axis reads the type", "TestMissing names the shape", "")

		if v := checkFilesOf(t, specPath); len(v) != 0 {
			t.Fatalf("CheckFiles = %#v, want no violation for a why-cell name", v)
		}
	})

	// The planned cell is the control for the bare-name case above: the same name, beside
	// the marker of a test the build has yet to write, keeps the behavior it had before the
	// rule. Both cells stay rows the uncited report names.
	for _, c := range []struct {
		name, seam string
		violations int
	}{
		{"a stray name reports its row as uncited", "covered by TestMissing", 1},
		{"a planned cell stays uncited and passes", "planned TestMissing in internal/x", 0},
	} {
		t.Run(c.name, func(t *testing.T) {
			_, specPath := citedSpec(t, "s", c.seam, "")
			content, err := os.ReadFile(specPath)
			if err != nil {
				t.Fatalf("ReadFile(%q): %v", specPath, err)
			}

			if got := uncitedRows(parse(content)); len(got) != 1 || got[0] != "1" {
				t.Fatalf("uncitedRows = %#v, want row 1", got)
			}
			if v := checkFilesOf(t, specPath); len(v) != c.violations {
				t.Fatalf("CheckFiles = %#v, want %d violation(s)", v, c.violations)
			}
		})
	}
}

// TestCommandCheckExitsOneOnAStrayTestName pins the citation-form rule at the surface a
// caller reads: `bench coverage --check` exits 1 and prints the diagnostic.
func TestCommandCheckExitsOneOnAStrayTestName(t *testing.T) {
	root, _ := citedSpec(t, "s", "covered by TestMissing", "")
	t.Chdir(root)

	out, code := Command([]string{"--check", "specs/s/spec.md"})
	if code != 1 || !strings.Contains(out, strayWant("TestMissing")) {
		t.Fatalf("Command = (%d, %q), want exit 1 with %q", code, out, strayWant("TestMissing"))
	}
}
