package roadmap

import (
	"os"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/retros"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

var retroGrammar = usage.Grammar{
	Cmd:     "bench retro",
	Help:    "usage: bench retro <slug> --body <markdown>",
	Flags:   []usage.Flag{{Name: "--body", HasValue: true, NoEmptyValue: true}},
	MaxArgs: 1,
}

// RetroCommand validates and replaces one primary-local implementation retrospective.
func RetroCommand(args []string) (string, int) {
	parsed, line, code := usage.Parse(retroGrammar, args)
	if line != "" {
		return line + "\n", code
	}
	if len(parsed.Positionals) != 1 || !retros.ValidSlug(parsed.Positionals[0]) {
		return retroGrammar.Help + "\n", 2
	}
	body, present := parsed.Flags["--body"]
	if !present {
		return toon.MissingArg(retroGrammar.Cmd, "--body") + "\n", 2
	}
	if err := retros.Parse([]byte(body)); err != nil {
		return toon.Errorf("invalid retrospective", err.Error()) + "\n", 1
	}
	path := retros.Path(parsed.Positionals[0])
	root, refusal, code := inboxRoot(path)
	if refusal != "" {
		return refusal, code
	}
	file := filepath.Join(root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		return cannotWrite(path, err), 1
	}
	if err := os.WriteFile(file, []byte(body), 0o644); err != nil {
		return cannotWrite(path, err), 1
	}
	return "captured: " + parsed.Positionals[0] + "\n", 0
}
