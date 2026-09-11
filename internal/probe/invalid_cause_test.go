package probe

import (
	"os"
	"strings"
	"testing"
)

func TestZeroMatchMutationsAreSubstringMissesWithoutRunning(t *testing.T) {
	f := newFixture(t)
	before, err := os.ReadFile(f.subject)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		args []string
	}{
		{"swap", []string{"clamp.go", "--swap", "absent text", "--with", "replacement", "--package", "./", "--run", "^TestClampNegative$"}},
		{"omit", []string{"clamp.go", "--omit", "absent text", "--package", "./", "--run", "^TestClampNegative$"}},
		{"unwrap", []string{"clamp.go", "--unwrap", "absentCall(n)", "--package", "./", "--run", "^TestClampNegative$"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, code := runProbe(t, tc.args...)
			if code != 1 || !strings.Contains(out, ",substring-miss,0,untouched") {
				t.Fatalf("output = (%q, %d)", out, code)
			}
			if _, err := os.Stat(f.marker); !os.IsNotExist(err) {
				t.Fatalf("test marker exists after zero-match mutation")
			}
			after, err := os.ReadFile(f.subject)
			if err != nil {
				t.Fatal(err)
			}
			if string(after) != string(before) {
				t.Fatal("zero-match mutation changed the subject")
			}
		})
	}
}
