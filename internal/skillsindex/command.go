package skillsindex

import (
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// grammar is the declared argument shape usage.Parse enforces for `bench skills-index`:
// `--check` (the default) prints the diagnostics Check would emit, `--write` regenerates
// the block in place. Declaring both as plain flags rather than encoding a mode string
// keeps a mistyped third spelling a usage error rather than a silently ignored no-op.
var grammar = usage.Grammar{
	Cmd:   "bench skills-index",
	Help:  "usage: bench skills-index [--check|--write]",
	Flags: []usage.Flag{{Name: "--check"}, {Name: "--write"}},
}

// Command implements `bench skills-index [--check|--write]`, the operator's one
// regenerator and drift check now that the retired shell regenerator is gone. `--check`
// (the default) prints Check's diagnostics one per line and exits 1 if any, 0 if
// clean. `--write` regenerates the block via Write; a blocking refusal (the
// allowlist unparseable, or the reference file missing or unmarked) prints Write's
// own error and exits 1 without the file being touched — Write refuses before any
// bytes are written, so the caller and the file agree on nothing having changed.
func Command(args []string) (string, int) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return line + "\n", code
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	if _, write := parsed.Flags["--write"]; write {
		if werr := Write(root); werr != nil {
			return werr.Error() + "\n", 1
		}
		return "", 0
	}
	diags := Check(root)
	if len(diags) == 0 {
		return "", 0
	}
	return strings.Join(diags, "\n") + "\n", 1
}
