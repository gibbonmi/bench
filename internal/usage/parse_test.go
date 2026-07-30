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

// variadicGrammar is the unbounded-positional fixture: the shape `bench idea` and
// `bench commit` take, where free text or a path list can legitimately contain the
// word "help".
func variadicGrammar() Grammar {
	return Grammar{
		Cmd:     "bench park",
		Help:    "usage: bench park \"<text>\"\n",
		MaxArgs: -1,
	}
}

// TestParseBareHelpOnlyWhenSole pins the boundary between the help request and
// ordinary input: bare "help" is a help request only as the sole argument, because a
// variadic grammar's free text or path list may legitimately contain the word, and
// recognizing it anywhere silently discards the rest of the invocation.
func TestParseBareHelpOnlyWhenSole(t *testing.T) {
	v := variadicGrammar()
	res, line, code := Parse(v, []string{"help"})
	if line != v.Help || code != 0 {
		t.Errorf("Parse(help) = (%+v, %q, %d), want help text and exit 0", res, line, code)
	}

	res, line, code = Parse(v, []string{"help", "me", "remember", "the", "parser"})
	if line != "" || code != 0 {
		t.Fatalf("Parse(help me remember the parser) errored: line=%q code=%d", line, code)
	}
	want := []string{"help", "me", "remember", "the", "parser"}
	if len(res.Positionals) != len(want) {
		t.Fatalf("positionals = %v, want %v", res.Positionals, want)
	}
	for i := range want {
		if res.Positionals[i] != want[i] {
			t.Fatalf("positionals = %v, want %v", res.Positionals, want)
		}
	}

	g := testGrammar()
	res, line, code = Parse(g, []string{"target", "help"})
	if line != "" || code != 0 {
		t.Fatalf("Parse(target help) errored: line=%q code=%d", line, code)
	}
	if len(res.Positionals) != 2 || res.Positionals[1] != "help" {
		t.Errorf("positionals = %v, want [target help]", res.Positionals)
	}

	// The flag-spelled requests are unambiguous, so they keep their anywhere-before-`--`
	// recognition even alongside other arguments.
	for _, spelling := range []string{"--help", "-h"} {
		_, line, code := Parse(v, []string{"some", "text", spelling})
		if line != v.Help || code != 0 {
			t.Errorf("Parse(some text %s) = (%q, %d), want help text and exit 0", spelling, line, code)
		}
	}
}

// TestParseEmptyPositionalIsUsageError pins the shape an unset shell variable produces
// (`bench commit -m "$MSG" "$FILE"` with FILE unset): an empty positional names no
// path, and a path-taking subcommand that resolves it against the filesystem widens
// silently to the cwd. Rejecting it here gives every grammar the guard at once.
func TestParseEmptyPositionalIsUsageError(t *testing.T) {
	for _, tc := range []struct {
		name string
		g    Grammar
		args []string
	}{
		{"sole empty positional", testGrammar(), []string{""}},
		{"empty among others", variadicGrammar(), []string{"a", "", "b"}},
		{"empty after the flag terminator", variadicGrammar(), []string{"--", ""}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res, line, code := Parse(tc.g, tc.args)
			if want := toon.Usage(tc.g.Cmd, `""`); line != want || code != 2 {
				t.Errorf("Parse(%q) = (%+v, %q, %d), want (%q, 2)", tc.args, res, line, code, want)
			}
		})
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
	if res.PositionalsBeforeTerminator != 0 {
		t.Errorf("Parse(-- -x) positionals before terminator = %d, want 0", res.PositionalsBeforeTerminator)
	}

	res, line, code = Parse(g, []string{"target", "--", "child"})
	if line != "" || code != 0 {
		t.Fatalf("Parse(target -- child) errored: line=%q code=%d", line, code)
	}
	if res.PositionalsBeforeTerminator != 1 {
		t.Errorf("Parse(target -- child) positionals before terminator = %d, want 1", res.PositionalsBeforeTerminator)
	}
}

func TestParseReservedPositionalsPrecedeHelpFlagsAndTerminator(t *testing.T) {
	g := Grammar{
		Cmd:                                 "bench worktree exec",
		Help:                                "usage: bench worktree exec <target> -- <command> [args...]\n",
		MinArgs:                             2,
		MaxArgs:                             -1,
		ReservedPositionalsBeforeTerminator: 1,
	}
	for _, target := range []string{"--help", "-label", "--"} {
		t.Run(target, func(t *testing.T) {
			res, line, code := Parse(g, []string{target, "--", "child"})
			if line != "" || code != 0 {
				t.Fatalf("Parse(%q -- child) = (%+v, %q, %d), want success", target, res, line, code)
			}
			if got := res.Positionals; len(got) != 2 || got[0] != target || got[1] != "child" {
				t.Errorf("positionals = %q, want [%q child]", got, target)
			}
			if !res.EndedFlags || res.PositionalsBeforeTerminator != 1 {
				t.Errorf("terminator result = %+v, want one positional before terminator", res)
			}
		})
	}

	_, line, code := Parse(g, []string{"--help"})
	if line != g.Help || code != 0 {
		t.Errorf("Parse(--help) = (%q, %d), want help text and exit 0", line, code)
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

// TestParseRepeatedFlagIsUsageError pins the duplicate-flag rule and every
// interaction it has with the rules around it. Last-one-wins is the cheapest wrong
// implementation and it exits 0, so each case asserts the rendered line and code
// rather than only the outcome.
func TestParseRepeatedFlagIsUsageError(t *testing.T) {
	g := testGrammar()
	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		// One rule covers both flag kinds: the value flag reports the flag itself
		// rather than consuming the repeat's value, and the boolean flag renders
		// through the same toon.Usage call.
		{"repeated value flag", []string{"-m", "a", "-m", "b", "t"}, toon.Usage(g.Cmd, "-m")},
		{"repeated boolean flag", []string{"--all", "--all", "t"}, toon.Usage(g.Cmd, "--all")},
		// A repeated value flag missing its second value is still the duplicate,
		// which is the earlier fault.
		{"repeated value flag with no second value", []string{"-m", "a", "-m"}, toon.Usage(g.Cmd, "-m")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res, line, code := Parse(g, tc.args)
			if line != tc.want || code != 2 {
				t.Errorf("Parse(%v) = (%+v, %q, %d), want (%q, 2)", tc.args, res, line, code, tc.want)
			}
		})
	}
}

// TestParseRepeatedFlagAfterDoubleDashIsPositional pins the interaction with the
// "--" rule: flag parsing has ended, so a repeat of a declared flag is an ordinary
// positional and the duplicate rule must not reach it. A duplicate check applied to
// the raw argv rather than to the flags actually parsed would fail here.
func TestParseRepeatedFlagAfterDoubleDashIsPositional(t *testing.T) {
	g := testGrammar()
	res, line, code := Parse(g, []string{"--all", "--", "--all"})
	if line != "" || code != 0 {
		t.Fatalf("Parse(--all -- --all) errored: line=%q code=%d", line, code)
	}
	if len(res.Positionals) != 1 || res.Positionals[0] != "--all" {
		t.Errorf("positionals = %v, want [--all]", res.Positionals)
	}
	if _, ok := res.Flags["--all"]; !ok {
		t.Errorf("flags = %v, want the pre-`--` --all recorded", res.Flags)
	}
}

// TestParseRepeatedHelpIsHelp pins the interaction with the help rule: help is
// checked before any flag bookkeeping and is always success, so a repeated help
// spelling — or help repeated alongside a repeated flag — still exits 0 with the
// declared help text rather than becoming a duplicate-flag error.
func TestParseRepeatedHelpIsHelp(t *testing.T) {
	g := testGrammar()
	for _, args := range [][]string{
		{"--help", "--help"},
		{"-h", "--help", "help"},
		{"--all", "--help", "--all"},
	} {
		res, line, code := Parse(g, args)
		if line != g.Help || code != 0 {
			t.Errorf("Parse(%v) = (%+v, %q, %d), want help text and exit 0", args, res, line, code)
		}
	}
}
