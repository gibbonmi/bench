package gocache

import (
	"os"

	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// grammar is the declared argument shape usage.Parse enforces for this subcommand.
// Arity, flag recognition, `--`, and help all come from there, instead of a local switch.
var grammar = usage.Grammar{
	Cmd:     "bench cache [clean]",
	Help:    "usage: bench cache\n       bench cache clean\n",
	MaxArgs: 1,
}

// cleanChild is the one word the verb accepts as a subcommand. Anything else is an unknown
// argument, answered by the same usage line an unknown flag gets.
const cleanChild = "clean"

// tableName is the one block `bench cache` prints, and tableFields is its schema.
const tableName = "go_build_cache"

var tableFields = []string{"dir", "bytes", "files", "last_trim", "bound", "over_bound"}

// Command implements `bench cache` and its one mutating child, `bench cache clean`. Both
// derive the directory from the process environment, so they read no repository and
// resolve no git root; an operator runs them anywhere on the machine, a directory outside
// a repository included. The bare verb is a read-only report, and an absent or empty
// directory is a zero row at exit 0, because a first run must pass.
func Command(args []string) (string, int) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return line + "\n", code
	}
	if len(parsed.Positionals) == 1 {
		if parsed.Positionals[0] != cleanChild {
			return toon.Usage(grammar.Cmd, parsed.Positionals[0]) + "\n", 2
		}
		return clean(os.Environ())
	}
	return report(os.Environ())
}

// derivedDir answers the build cache directory for env, or the operator's refusal line
// when the environment gives no directory or the path holds a control byte. A control
// byte is refused before any walk, so no cell reaches the encoder that the encoder would
// refuse for a second reason.
func derivedDir(env []string) (string, string) {
	dir, err := Dir(env)
	if err != nil {
		return "", notDerived(err)
	}
	if !toon.Representable(dir) {
		return "", toon.Errorf("unrepresentable cache directory", "the Bench build cache path holds a control byte; clear it from HOME") + "\n"
	}
	return dir, ""
}

// notDerived is the one refusal a failed derivation produces, for the report, the clean,
// and the clean's child environment alike.
func notDerived(err error) string {
	return toon.Errorf("cache directory not derived", err.Error()) + "\n"
}

// report is the command body over an explicit environment slice, which is the seam the
// command tests drive.
func report(env []string) (string, int) {
	dir, refusal := derivedDir(env)
	if refusal != "" {
		return refusal, 1
	}
	footprint := Measure(dir)
	row := []any{footprint.Dir, footprint.Bytes, footprint.Files, footprint.LastTrim, Bound, footprint.OverBound()}
	block, err := toon.TableTyped(tableName, tableFields, [][]any{row})
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return block, 0
}
