package gocache

import (
	"os"

	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// grammar is the declared argument shape usage.Parse enforces for this subcommand.
// Arity, flag recognition, `--`, and help all come from there, instead of a local switch.
var grammar = usage.Grammar{
	Cmd:  "bench cache",
	Help: "usage: bench cache",
}

// tableName is the one block `bench cache` prints, and tableFields is its schema.
const tableName = "go_build_cache"

var tableFields = []string{"dir", "bytes", "files", "last_trim", "bound", "over_bound"}

// Command implements `bench cache`: a read-only report of the Bench build cache
// footprint. It derives the directory from the process environment, so it reads no
// repository and resolves no git root. An operator therefore runs it anywhere on the
// machine, a directory outside a repository included. An absent or empty directory is a
// zero row at exit 0, because a first run must pass.
func Command(args []string) (string, int) {
	if _, line, code := usage.Parse(grammar, args); line != "" {
		return line + "\n", code
	}
	return report(os.Environ())
}

// report is the command body over an explicit environment slice, which is the seam the
// command tests drive.
func report(env []string) (string, int) {
	dir, err := Dir(env)
	if err != nil {
		return toon.Errorf("cache directory not derived", err.Error()) + "\n", 1
	}
	// A control byte in the path is refused before the walk, so no table renders it and
	// no cell reaches the encoder that the encoder would refuse for a second reason.
	if !toon.Representable(dir) {
		return toon.Errorf("unrepresentable cache directory", "the Bench build cache path holds a control byte; clear it from HOME") + "\n", 1
	}
	footprint := Measure(dir)
	row := []any{footprint.Dir, footprint.Bytes, footprint.Files, footprint.LastTrim, Bound, footprint.OverBound()}
	block, err := toon.TableTyped(tableName, tableFields, [][]any{row})
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return block, 0
}
