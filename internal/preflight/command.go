package preflight

import (
	"fmt"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// grammar is the declared argument shape usage.Parse enforces for this
// subcommand. Two required positionals (mode, slug), an optional explicit
// base, and the optional frozen source tip the review phase pins. Both
// `review` and `build` are accepted modes. The mode-validity check below
// rejects anything else the same way it rejects any unknown word.
var grammar = usage.Grammar{
	Cmd:  "bench preflight",
	Help: "usage: bench preflight review <slug> [--base <commit>] [--source-tip <commit>]\n       bench preflight build <slug> [--base <commit>] [--source-tip <commit>]\n",
	Flags: []usage.Flag{
		{Name: "--base", HasValue: true, NoEmptyValue: true},
		{Name: "--source-tip", HasValue: true, NoEmptyValue: true},
	},
	MinArgs: 2,
	MaxArgs: 2,
}

// Command implements `bench preflight review <slug>` and `bench preflight build
// <slug>`. It is the CLI-contract seam. Grammar and usage errors ride
// usage.Parse (exit 2). A not-in-repo cwd or a bootstrap failure is one
// toon.Errorf line (exit 1). Otherwise the verdict renders as TOON and the
// exit code follows Verdict.Red (0 green, 1 red).
func Command(args []string) (string, int) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return line + "\n", code
	}
	mode, slug := parsed.Positionals[0], parsed.Positionals[1]
	base := parsed.Flags["--base"]
	sourceTip := parsed.Flags["--source-tip"]
	if mode != "review" && mode != modeBuild {
		return toon.Usage(grammar.Cmd, mode) + "\n", 2
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}

	if err := unrepresentableCell("--source-tip", sourceTip); err != nil {
		return toon.RenderError(err) + "\n", 1
	}

	facts, bootErr := GatherPinned(root, mode, slug, base, sourceTip)
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
		rows[i] = []string{c.Check, c.Verdict, c.Detail, c.Next}
	}
	tbl, err := toon.Table("checks", []string{"check", "verdict", "detail", "next"}, rows)
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
	help, err := axi.RenderHelp([]axi.Action{axi.RetryInvocation(invocation...)})
	if err != nil {
		return toon.RenderError(err) + "\n"
	}
	return toon.Errorf("snapshot drift", hint) + "\n" + help
}

// unrepresentableChangedPath refuses a changed path spec-TOON cannot render as a
// cell, before the verdict table is ever built. PF7's contract is unconditional.
// A path carrying a control byte exits 1 the same way in every case: a
// green row (never rendered) or a red row's detail cell. The refusal
// never depends on which row a later check sorts it into.
func unrepresentableChangedPath(paths []string) error {
	for _, p := range paths {
		if err := unrepresentableCell("changed path", p); err != nil {
			return err
		}
	}
	return nil
}

// unrepresentableCell is that refusal for one value. The changed-path
// sweep and --source-tip both share it: a pin reaches a detail cell and
// the snapshot-drift retry action. So a control byte in it gets refused,
// not rendered. The %q quoting keeps the offending byte out of the
// message it explains.
func unrepresentableCell(what, value string) error {
	if toon.Representable(value) {
		return nil
	}
	return fmt.Errorf("%s %q contains a byte spec-TOON cannot represent", what, value)
}
