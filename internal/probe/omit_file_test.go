package probe

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
)

func TestOmitFileUsesExactBytesAndReportsOmit(t *testing.T) {
	f := newFixture(t)
	input := filepath.Join(t.TempDir(), "substring")
	needle := "if n < 0 {\n\t\treturn 0\n\t}"
	writeFixtureFile(t, input, needle, 0o644)
	requireMutantBites(t, f, strings.Replace(clampSource, needle, "", 1), "omit", "--omit-file", input)
}

// mutationFileForms is the argv of each flag that reads its text through the one mutation
// file reader, so each refusal row runs over both flags.
var mutationFileForms = []struct {
	flag string
	args func(path string) []string
}{
	{"--omit-file", func(path string) []string { return []string{"clamp.go", "--omit-file", path, "--package", "./"} }},
	{"--with-file", func(path string) []string {
		return []string{"clamp.go", "--swap", "n < 0", "--with-file", path, "--package", "./"}
	}},
}

func TestOmitFileRefusesInvalidInputsBeforeRunning(t *testing.T) {
	for _, tc := range []struct {
		name   string
		make   func(string) error
		reason string
	}{
		{"absent", func(string) error { return nil }, "absent"},
		{"empty", func(path string) error { return os.WriteFile(path, nil, 0o644) }, "empty"},
		{"symlink", func(path string) error { return os.Symlink("missing", path) }, "a symlink"},
		{"directory", func(path string) error { return os.Mkdir(path, 0o755) }, "a directory"},
		{"special", func(path string) error { return syscall.Mkfifo(path, 0o600) }, "a special file"},
		{"oversized", func(path string) error {
			return os.WriteFile(path, []byte(strings.Repeat("x", int(bounds.ControlRecordLimit)+1)), 0o644)
		}, "oversized"},
	} {
		for _, form := range mutationFileForms {
			t.Run(form.flag+" "+tc.name, func(t *testing.T) {
				f := refusalFixture(t)
				path := filepath.Join(t.TempDir(), "mutation")
				if tc.name != "absent" {
					if err := tc.make(path); err != nil {
						t.Fatal(err)
					}
				}
				out, code := runProbe(t, form.args(path)...)
				requireRefused(t, f, out, code, "error: probe mutation file unavailable — "+path+" is "+tc.reason+"\n", 1)
			})
		}
	}
}

func TestOmitFileRefusesAnUnreadableInputBeforeRunning(t *testing.T) {
	for _, form := range mutationFileForms {
		t.Run(form.flag, func(t *testing.T) {
			f := refusalFixture(t)
			path := filepath.Join(t.TempDir(), "mutation")
			writeFixtureFile(t, path, "return 0", 0o000)
			t.Cleanup(func() { _ = os.Chmod(path, 0o600) })
			if _, err := os.ReadFile(path); err == nil {
				capability.Capability(t, capability.Privilege, "the test process reads mode 0000 files, so an unreadable mutation file is unobservable")
			}
			out, code := runProbe(t, form.args(path)...)
			requireRefused(t, f, out, code, "error: probe mutation file unavailable — "+path+" is unreadable\n", 1)
		})
	}
}

func TestOmitFileRefusesZeroAndMultipleMatchesSeparately(t *testing.T) {
	t.Run("zero matches", func(t *testing.T) {
		f := refusalFixture(t)
		path := filepath.Join(t.TempDir(), "mutation")
		writeFixtureFile(t, path, "absent", 0o644)
		out, code := runProbe(t, "clamp.go", "--omit-file", path, "--package", "./", "--run", "^TestClampNegative$")
		want := probeRow(t, "invalid", "clamp.go", "omit", "substring-miss", 0, nil, "untouched") +
			selectionRow(t, "package", "./", "^TestClampNegative$", "", 0) + executionRow(t, 1, nil, 0)
		requireRefused(t, f, out, code, want, 1)
	})
	t.Run("multiple matches", func(t *testing.T) {
		f := refusalFixture(t)
		path := filepath.Join(t.TempDir(), "mutation")
		writeFixtureFile(t, path, "return", 0o644)
		out, code := runProbe(t, "clamp.go", "--omit-file", path, "--package", "./", "--run", "^TestClampNegative$")
		requireRefused(t, f, out, code, "error: probe mutation ambiguous — the old string matches 2 times, want exactly 1\n", 1)
	})
}
