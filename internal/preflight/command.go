package preflight

import (
	"fmt"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// grammar is the declared argument shape usage.Parse enforces for this subcommand —
// two required positionals (mode, slug) and an optional explicit base. Both `review`
// and `build` are
// accepted modes; anything else is rejected by the mode-validity check below the same
// way an unknown word always is.
var grammar = usage.Grammar{
	Cmd:     "bench preflight",
	Help:    "usage: bench preflight review <slug> [--base <commit>]\n       bench preflight build <slug> [--base <commit>]\n",
	Flags:   []usage.Flag{{Name: "--base", HasValue: true, NoEmptyValue: true}},
	MinArgs: 2,
	MaxArgs: 2,
}

// Command implements `bench preflight review <slug>` and `bench preflight build
// <slug>`. It is the CLI-contract seam: grammar and usage errors ride usage.Parse
// (exit 2); a not-in-repo cwd or a bootstrap failure is one toon.Errorf line (exit
// 1); otherwise the five-check verdict renders as TOON and the exit code follows
// Verdict.Red (0 green, 1 red).
func Command(args []string) (string, int) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return line + "\n", code
	}
	mode, slug := parsed.Positionals[0], parsed.Positionals[1]
	base := parsed.Flags["--base"]
	if mode != "review" && mode != modeBuild {
		return toon.Usage(grammar.Cmd, mode) + "\n", 2
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}

	facts, bootErr := Gather(root, mode, slug, base)
	if bootErr != nil {
		if bootErr.Kind == "snapshot drift" {
			return snapshotDriftRefusal(args, bootErr.Hint), 1
		}
		return toon.Errorf(bootErr.Kind, bootErr.Hint) + "\n", 1
	}
	if err := unrepresentableChangedPath(facts.ChangedPaths); err != nil {
		return toon.RenderError(err) + "\n", 1
	}

	verdict := Decide(facts)
	rows := make([][]string, len(verdict.Checks))
	for i, c := range verdict.Checks {
		rows[i] = []string{c.Check, c.Verdict, c.Detail}
	}
	tbl, err := toon.Table("checks", []string{"check", "verdict", "detail"}, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}

	var b strings.Builder
	fmt.Fprintf(&b, "phase: %s\n", mode)
	fmt.Fprintf(&b, "spec: %s\n", facts.SpecPath)
	if facts.SourceBase != "" {
		source, err := toon.Table("source", []string{"base", "tip"}, [][]string{{facts.SourceBase, facts.SourceTip}})
		if err != nil {
			return toon.RenderError(err) + "\n", 1
		}
		b.WriteString(source)
	}
	b.WriteString(tbl)

	exit := 0
	if verdict.Red {
		exit = 1
	}
	return b.String(), exit
}

func snapshotDriftRefusal(args []string, hint string) string {
	invocation := make([]axi.InvocationArgument, 0, len(args)+1)
	invocation = append(invocation, axi.KnownArgument("preflight"))
	for _, arg := range args {
		invocation = append(invocation, axi.KnownArgument(arg))
	}
	help, err := axi.RenderHelp([]axi.Action{axi.ExecutableInvocation("retry after the repository stopped moving", invocation...)})
	if err != nil {
		return toon.RenderError(err) + "\n"
	}
	return toon.Errorf("snapshot drift", hint) + "\n" + help
}

// unrepresentableChangedPath refuses a changed path spec-TOON cannot render as a
// cell, before the verdict table is ever built. PF7's contract is unconditional: a
// path carrying a control byte exits 1 the same way whether it would land in a green
// row (never rendered) or a red one's detail cell — the refusal cannot depend on
// which row a later check happens to sort it into.
func unrepresentableChangedPath(paths []string) error {
	for _, p := range paths {
		if !toon.Representable(p) {
			return fmt.Errorf("changed path %q contains a byte spec-TOON cannot represent", p)
		}
	}
	return nil
}
