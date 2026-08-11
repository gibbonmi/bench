package preflight

import (
	"fmt"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// grammar is the declared argument shape usage.Parse enforces for this subcommand —
// two required positionals (mode, slug), no flags yet. Both `review` and `build` are
// accepted modes; anything else is rejected by the mode-validity check below the same
// way an unknown word always is.
var grammar = usage.Grammar{
	Cmd:     "bench preflight",
	Help:    "usage: bench preflight review <slug>\n       bench preflight build <slug>\n",
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
	if mode != "review" && mode != modeBuild {
		return toon.Usage(grammar.Cmd, mode) + "\n", 2
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}

	facts, bootErr := Gather(root, mode, slug)
	if bootErr != nil {
		return toon.Errorf(bootErr.Kind, bootErr.Hint) + "\n", 1
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
	b.WriteString(tbl)

	exit := 0
	if verdict.Red {
		exit = 1
	}
	return b.String(), exit
}
