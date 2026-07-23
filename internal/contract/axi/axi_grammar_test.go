package axi

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestAXIGrammarContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "AXI trailing-garbage contract", testAXITrailingGarbage)
	contract.RunParallel(t, "AXI satisfied-grammar contract", testAXISatisfiedGrammar)
	contract.RunParallel(t, "AXI help-is-success contract", testAXIHelpIsSuccess)
	contract.RunParallel(t, "AXI routed flat-subcommand contract", testAXIRoutedFlatSubcommands)
	contract.RunParallel(t, "AXI repeated-flag contract", testAXIRepeatedFlag)
}

// routedFlatCases cover the remaining flat subcommands the routing registry records as
// routed. They are separate from the four the trailing-garbage story enumerates: these
// already rejected an unknown token, so what routing buys them is a usage line naming the
// offending token rather than the flag before it, plus help at exit 0 everywhere.
var routedFlatCases = []struct {
	name    string
	garbage []string // invocation whose last token is excess; empty skips the check
	usage   string
	help    []string
	helpFor string
}{
	{name: "learnings", garbage: []string{"learnings", "x"}, usage: "usage: bench learnings (unknown argument: x)", help: []string{"learnings", "-h"}, helpFor: "usage: bench learnings"},
	{name: "status", garbage: []string{"status", "--all", "x"}, usage: "usage: bench status (unknown argument: x)", help: []string{"status", "-h"}, helpFor: "usage: bench status [--all]"},
	{name: "models", garbage: []string{"models", "x"}, usage: "usage: bench models (unknown argument: x)", help: []string{"models", "-h"}, helpFor: "usage: bench models"},
	{name: "outline", garbage: []string{"outline", "--full", "a", "x"}, usage: "usage: bench outline (unknown argument: x)", help: []string{"outline", "-h"}, helpFor: "usage: bench outline"},
	{name: "idea", help: []string{"idea", "-h"}, helpFor: `usage: bench idea "<text>"`},
	{name: "diff", garbage: []string{"diff", "--full", "x"}, usage: "usage: bench diff (unknown argument: x)", help: []string{"diff", "-h"}, helpFor: "usage: bench diff [--full] [--commit <sha>]"},
	{name: "roadmap", garbage: []string{"roadmap", "--context", "x"}, usage: "usage: bench roadmap --context (unknown argument: x)", help: []string{"roadmap", "-h"}, helpFor: "usage: bench roadmap --context [--full]"},
}

// repeatedFlagCases are the invocations whose contract is that a declared flag given
// twice is a usage error. diff and roadmap could not route until the grammar owned that
// rule; status had the posture before it routed and regains it here.
var repeatedFlagCases = []struct {
	args  []string
	usage string
}{
	{[]string{"diff", "--commit", "HEAD", "--commit", "HEAD"}, "usage: bench diff (unknown argument: --commit)"},
	{[]string{"roadmap", "--context", "--context"}, "usage: bench roadmap --context (unknown argument: --context)"},
	{[]string{"status", "--all", "--all"}, "usage: bench status (unknown argument: --all)"},
}

func testAXIRepeatedFlag(t *testing.T) {
	for _, tc := range repeatedFlagCases {
		t.Run(strings.Join(tc.args, " "), func(t *testing.T) {
			t.Parallel()
			f := grammarFixture(t)

			probe := f.Bench(tc.args...)

			probe.RequireExit(2)
			requireOutputLine(t, probe, tc.usage)
		})
	}
}

func testAXIRoutedFlatSubcommands(t *testing.T) {
	for _, tc := range routedFlatCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := grammarFixture(t)

			if len(tc.garbage) > 0 {
				garbage := f.Bench(tc.garbage...)
				garbage.RequireExit(2)
				requireOutputLine(t, garbage, tc.usage)
			}

			help := f.Bench(tc.help...)
			help.RequireExit(0)
			if !strings.Contains(help.Stdout, tc.helpFor) {
				t.Fatalf("bench %s -h printed no help\nstdout:\n%s\nstderr:\n%s", tc.name, help.Stdout, help.Stderr)
			}
		})
	}
}

// grammarFixture is a repo with one commit and a git identity, so `--since HEAD`
// resolves a real range and the trailing-garbage probes fail on the grammar rather
// than on a missing ref.
func grammarFixture(t *testing.T) contract.Fixture {
	t.Helper()
	f := contract.NewFixture(t)
	f.Git("config", "user.email", "bench@local")
	f.Git("config", "user.name", "bench")
	f.WriteFile("README.md", "prose\n")
	f.CommitAll("seed")
	return f
}

// trailingGarbageCases enumerate the complete set of subcommands that accepted a
// satisfied grammar and then ignored whatever followed it. Naming all four is the
// quantifier: a fix applied to one of them still leaves this table red.
var trailingGarbageCases = []struct {
	name  string
	args  []string
	usage string
	// satisfied is the marker the same invocation without the trailing token must
	// still print, so tightening arity by rejecting the flag outright fails too.
	satisfied string
}{
	{"maps --count", []string{"maps", "--count"}, "usage: bench maps (unknown argument: x)", "0"},
	{"guards --brief", []string{"guards", "--brief"}, "usage: bench guards (unknown argument: x)", "full manifests: bench guards"},
	{"dashboard --stdout", []string{"dashboard", "--stdout"}, "usage: bench dashboard (unknown argument: x)", "<html lang=\"en\">"},
	{"structure --since", []string{"structure", "--since", "HEAD"}, "usage: bench structure (unknown argument: x)", "no tracked source files to check"},
}

func testAXITrailingGarbage(t *testing.T) {
	for _, tc := range trailingGarbageCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := grammarFixture(t)

			probe := f.Bench(append(append([]string(nil), tc.args...), "x")...)

			probe.RequireExit(2)
			requireOutputLine(t, probe, tc.usage)
		})
	}
}

func testAXISatisfiedGrammar(t *testing.T) {
	for _, tc := range trailingGarbageCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := grammarFixture(t)

			probe := f.Bench(tc.args...)

			probe.RequireExit(0)
			probe.RequireContains(probe.Stdout, tc.satisfied)
		})
	}
}

// helpIsSuccessCases are the subcommands that reported `-h` as misuse: commit and
// dashboard exited 2, and `commands` had no help form at all. coverage already
// complied and rides along as the boundary the rule must not break.
var helpIsSuccessCases = []struct {
	subcommand string
	help       string
}{
	{"commit", "usage: bench commit"},
	{"dashboard", "usage: bench dashboard"},
	{"commands", "usage: bench commands"},
	{"coverage", "usage: bench coverage"},
}

func testAXIHelpIsSuccess(t *testing.T) {
	for _, tc := range helpIsSuccessCases {
		t.Run(tc.subcommand, func(t *testing.T) {
			t.Parallel()
			f := grammarFixture(t)

			probe := f.Bench(tc.subcommand, "-h")

			probe.RequireExit(0)
			if !strings.Contains(probe.Stdout, tc.help) {
				t.Fatalf("bench %s -h printed no help\nstdout:\n%s\nstderr:\n%s", tc.subcommand, probe.Stdout, probe.Stderr)
			}
		})
	}
}
