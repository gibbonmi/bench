package sanitize_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// requireOneArgument requires that ShellQuote renders value as want, the quoted bytes
// authored from value, and that a POSIX shell reads the rendered value between two fixed
// words back as exactly one argument with the bytes of value. The fixed words make a
// dropped argument change the count.
func requireOneArgument(t *testing.T, value, want string) {
	t.Helper()
	got := sanitize.ShellQuote(value)
	if got != want {
		t.Errorf("ShellQuote(%q) = %q, want %q", value, got, want)
	}
	argv, err := axitest.HelpAction{Cmd: "before " + got + " after"}.Argv()
	if wantArgv := []string{"before", value, "after"}; err != nil || !slices.Equal(argv, wantArgv) {
		t.Errorf("argv of ShellQuote(%q) = (%q, %v), want %q", value, argv, err, wantArgv)
	}
}

// TestShellQuoteRoundTrip: an empty value, embedded quotes, and shell-active bytes each
// survive a POSIX shell as one literal argument. A sentinel file reds a value that the shell
// executes. The test runs in a temporary directory, so a redirect that an unquoted value
// performs cannot write into the source tree.
func TestShellQuoteRoundTrip(t *testing.T) {
	t.Chdir(t.TempDir())
	sentinel := filepath.Join(t.TempDir(), "executed")
	for _, row := range []struct {
		name, value, want string
	}{
		{name: "empty", value: "", want: "''"},
		{name: "an embedded quote", value: "it's", want: `'it'\''s'`},
		{name: "adjacent quotes", value: "''", want: `''\'''\'''`},
		{name: "a leading quote", value: "'x", want: `''\''x'`},
		{name: "spaces", value: "my  alpha ", want: "'my  alpha '"},
		{name: "glob characters", value: "*.go ?[ab] ~", want: "'*.go ?[ab] ~'"},
		{name: "dollar signs", value: "$HOME ${PATH} $1", want: "'$HOME ${PATH} $1'"},
		{name: "shell operators", value: "a;b|c&d>e<f #g", want: "'a;b|c&d>e<f #g'"},
		{name: "a backslash", value: `a\b\n`, want: `'a\b\n'`},
		{name: "Unicode", value: "héllo wörld", want: "'héllo wörld'"},
		{name: "a command substitution", value: "$(touch " + sentinel + ")", want: "'$(touch " + sentinel + ")'"},
		{name: "a backquote substitution", value: "`touch " + sentinel + "`", want: "'`touch " + sentinel + "`'"},
		{name: "a command list", value: "x; touch " + sentinel, want: "'x; touch " + sentinel + "'"},
	} {
		t.Run(row.name, func(t *testing.T) { requireOneArgument(t, row.value, row.want) })
	}
	if _, err := os.Stat(sentinel); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("a quoted value ran as a command: %v", err)
	}
}

// TestShellQuoteAlwaysQuoted: a value that needs no quoting still renders inside single
// quotes, so ShellQuote has no safe-token shortcut.
func TestShellQuoteAlwaysQuoted(t *testing.T) {
	t.Parallel()
	for _, row := range []struct{ value, want string }{
		{value: "alpha", want: "'alpha'"},
		{value: "a-b_c.d/e", want: "'a-b_c.d/e'"},
		{value: "x=1,y:2@%+", want: "'x=1,y:2@%+'"},
	} {
		t.Run(row.value, func(t *testing.T) { requireOneArgument(t, row.value, row.want) })
	}
}
