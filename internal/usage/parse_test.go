package usage

import (
	"testing"

	"github.com/gibbonmi/bench/internal/toon"
)

// testGrammar is the shared fixture for the parse tests: a small subcommand
// taking a value flag, a boolean flag, and one required plus one optional
// positional — enough shape to exercise every grammar rule below.
func testGrammar() Grammar {
	return Grammar{
		Cmd:  "bench frobnicate",
		Help: "usage: bench frobnicate <target> [extra]\n",
		Flags: []Flag{
			{Name: "-m", HasValue: true},
			{Name: "--all"},
		},
		MinArgs: 1,
		MaxArgs: 2,
	}
}

// TestParseHelpSpellings pins that each of the three help spellings
// short-circuits to the declared help text at exit 0, ahead of an otherwise
// satisfied grammar — a partial implementation that only recognizes --help
// must fail on the other two.
func TestParseHelpSpellings(t *testing.T) {
	g := testGrammar()
	for _, spelling := range []string{"help", "--help", "-h"} {
		res, line, code := Parse(g, []string{spelling})
		if line != g.Help || code != 0 {
			t.Errorf("Parse(%q) = (%+v, %q, %d), want help text and exit 0", spelling, res, line, code)
		}
	}
}

// TestParseDoubleDashEndsFlags pins the leading-dash-positional rule the
// grammar exists to make expressible: once "--" is seen, a later dash-prefixed
// token is an ordinary positional, never an unknown-flag error.
func TestParseDoubleDashEndsFlags(t *testing.T) {
	g := testGrammar()
	res, line, code := Parse(g, []string{"--", "-x"})
	if line != "" || code != 0 {
		t.Fatalf("Parse(-- -x) errored: line=%q code=%d", line, code)
	}
	if len(res.Positionals) != 1 || res.Positionals[0] != "-x" {
		t.Errorf("Parse(-- -x) positionals = %v, want [-x]", res.Positionals)
	}
}

// TestParseBareDashAndSecondDoubleDash pins the two tokens a naive
// dash-prefix check mis-parses as flags: a lone "-" is a positional (a legal
// filename) both in ordinary flag position and after flag parsing has ended,
// and a second "--" — once flag parsing has already ended — is an ordinary
// positional rather than being treated specially again.
func TestParseBareDashAndSecondDoubleDash(t *testing.T) {
	g := testGrammar()
	// In flag position (no "--" seen yet), a naive `strings.HasPrefix(a, "-")`
	// check would misfire on this exact token.
	res, line, code := Parse(g, []string{"-", "target"})
	if line != "" || code != 0 {
		t.Fatalf("Parse(- target) errored: line=%q code=%d", line, code)
	}
	if len(res.Positionals) != 2 || res.Positionals[0] != "-" || res.Positionals[1] != "target" {
		t.Errorf("Parse(- target) positionals = %v, want [- target]", res.Positionals)
	}

	// The first "--" ends flag parsing; "-" and the second "--" that follow
	// are then ordinary positionals rather than a separator or a flag.
	res, line, code = Parse(g, []string{"--", "-", "--"})
	if line != "" || code != 0 {
		t.Fatalf("Parse(-- - --) errored: line=%q code=%d", line, code)
	}
	want := []string{"-", "--"}
	if len(res.Positionals) != 2 || res.Positionals[0] != want[0] || res.Positionals[1] != want[1] {
		t.Errorf("Parse(-- - --) positionals = %v, want %v", res.Positionals, want)
	}
}

// TestParseMissingFlagValueDistinctFromUnknownFlag pins that a value flag with
// nothing following renders MissingArg, not Usage — collapsing the two labels
// is the cheapest wrong implementation and it still exits 2, so the rendered
// line is asserted alongside the code.
func TestParseMissingFlagValueDistinctFromUnknownFlag(t *testing.T) {
	g := testGrammar()
	_, line, code := Parse(g, []string{"target", "-m"})
	if want := toon.MissingArg(g.Cmd, "-m"); line != want || code != 2 {
		t.Errorf("Parse(target -m) = (%q, %d), want (%q, 2)", line, code, want)
	}
	_, line, code = Parse(g, []string{"target", "-z"})
	if want := toon.Usage(g.Cmd, "-z"); line != want || code != 2 {
		t.Errorf("Parse(target -z) = (%q, %d), want (%q, 2)", line, code, want)
	}
}

// TestParseArityUnmetIsMissingArg pins fewer-than-required positionals as
// MissingArg, distinct from the excess-positional Usage case below.
func TestParseArityUnmetIsMissingArg(t *testing.T) {
	g := testGrammar()
	_, line, code := Parse(g, nil)
	if want := toon.MissingArg(g.Cmd, "argument"); line != want || code != 2 {
		t.Errorf("Parse(nil) = (%q, %d), want (%q, 2)", line, code, want)
	}
}

// TestParseExcessPositionalNamesFirstExcess pins the trailing-garbage rule:
// the error names the first argument beyond the declared arity, not a generic
// message, so a mistyped invocation reports which token was unexpected.
func TestParseExcessPositionalNamesFirstExcess(t *testing.T) {
	g := testGrammar()
	_, line, code := Parse(g, []string{"target", "extra", "surplus"})
	if want := toon.Usage(g.Cmd, "surplus"); line != want || code != 2 {
		t.Errorf("Parse(target extra surplus) = (%q, %d), want (%q, 2)", line, code, want)
	}
}

// TestParseSuccess pins the ordinary satisfied-grammar path: flags recorded,
// positionals in argv order, empty line and exit 0.
func TestParseSuccess(t *testing.T) {
	g := testGrammar()
	res, line, code := Parse(g, []string{"-m", "msg", "--all", "target", "extra"})
	if line != "" || code != 0 {
		t.Fatalf("Parse success case errored: line=%q code=%d", line, code)
	}
	if res.Flags["-m"] != "msg" {
		t.Errorf("Flags[-m] = %q, want msg", res.Flags["-m"])
	}
	if _, ok := res.Flags["--all"]; !ok {
		t.Errorf("Flags[--all] missing")
	}
	want := []string{"target", "extra"}
	if len(res.Positionals) != 2 || res.Positionals[0] != want[0] || res.Positionals[1] != want[1] {
		t.Errorf("Positionals = %v, want %v", res.Positionals, want)
	}
}
