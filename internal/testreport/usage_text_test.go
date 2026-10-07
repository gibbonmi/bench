package testreport

import (
	"strings"
	"testing"
)

// systemRunForm is authored apart from the grammar, so a grammar that drops the
// filtered system form reds this test.
const systemRunForm = "bench test [--full] --check system --run <go-regex>"

// TestUsageTextNamesSystemRunForm grades that a usage refusal and the help text
// both name the filtered system form.
func TestUsageTextNamesSystemRunForm(t *testing.T) {
	for _, args := range [][]string{
		{"--check", "line-routing", "--run", "^TestX$"},
		{"--help"},
	} {
		output, _ := Command(t.TempDir(), args)
		if !strings.Contains(output, systemRunForm) {
			t.Errorf("Command(%q) omits %q:\n%s", args, systemRunForm, output)
		}
	}
}
