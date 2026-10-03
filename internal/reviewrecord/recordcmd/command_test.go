package recordcmd_test

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/reviewrecord/recordcmd"
)

func TestRecordGrammarRefusals(t *testing.T) {
	f := linked(t, 1)
	for _, args := range [][]string{nil, {"nosuch"}, {"chunk"}, {"chunk", slug(t), "--chunk", "1", "--base", "HEAD"}, {"completion", slug(t)}} {
		out, code := recordcmd.Command(f.Root, args)
		if code != 2 || !strings.HasPrefix(out, "usage: bench record") {
			t.Errorf("%q = exit %d, output %q; want exit 2 with a usage: bench record line", args, code, out)
		}
	}
}

// help is the `bench record --help` answer, which the other spellings must repeat.
func help(t *testing.T) string {
	t.Helper()
	out, code := recordcmd.Command("", []string{"--help"})
	if code != 0 || !strings.HasPrefix(out, "usage: bench record") {
		t.Fatalf("--help = exit %d, output %q; want exit 0 with usage: bench record lines", code, out)
	}
	return out
}

func TestRecordHelpPrintsEachForm(t *testing.T) {
	rows := recordcmd.HelpRows()
	lines := strings.Split(strings.TrimSuffix(help(t), "\n"), "\n")
	if len(rows) == 0 || len(lines) != len(rows) {
		t.Fatalf("help lines %q, want one for each of %d forms", lines, len(rows))
	}
	for i, row := range rows {
		if want := "usage: bench record" + row.Suffix; lines[i] != want {
			t.Errorf("help line %d = %q, want %q", i, lines[i], want)
		}
		// Each form's own grammar answers its help with the same line.
		form := strings.Fields(row.Suffix)[0]
		if out, code := recordcmd.Command("", []string{form, "--help"}); code != 0 || out != lines[i]+"\n" {
			t.Errorf("%s --help = exit %d, output %q; want %q", form, code, out, lines[i])
		}
	}
}

func TestRecordHelpSpellings(t *testing.T) {
	want := help(t)
	for _, spelling := range []string{"-h", "help"} {
		if out, code := recordcmd.Command("", []string{spelling}); code != 0 || out != want {
			t.Errorf("%s = exit %d, output %q; want exit 0 and %q", spelling, code, out, want)
		}
	}
}
