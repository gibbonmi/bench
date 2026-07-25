package handoff

import (
	"os"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/status"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// grammar is the declared argument shape usage.Parse enforces for this subcommand —
// arity, flag recognition, `--`, and help all come from there rather than a local switch.
// The `--harness` *value* is this package's to check: the shared parser validates flag
// spelling, repetition, and arity, and accepts anything as a declared flag's value.
var grammar = usage.Grammar{
	Cmd:   "bench handoff",
	Help:  "usage: bench handoff [--harness claude|codex] [--next <command>]",
	Flags: []usage.Flag{{Name: "--harness", HasValue: true}, {Name: "--next", HasValue: true}},
}

// Command implements `bench handoff`. It prints the pin block and rewrites
// session-handoff.md from one derivation, preserving the reviewer-owned State section and
// regenerating everything else. It creates freely and destroys never: a missing file is
// scaffolded, and a file whose State section cannot be located unambiguously is left
// exactly as it was behind a non-zero exit.
//
// Exits: 2 for a usage error, 1 outside a repository or when the document cannot be parsed,
// rendered, or written, 0 otherwise. Nothing is written on any non-zero path.
func Command(args []string) (string, int) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return line + "\n", code
	}
	harness := harnessClaude
	if value, present := parsed.Flags["--harness"]; present {
		if _, known := harnessPrefix[value]; !known {
			return toon.Usage(grammar.Cmd, "--harness "+value) + "\n", 2
		}
		harness = value
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}

	target := filepath.Join(root, status.HandoffFile)
	state, err := preservedState(target)
	if err != nil {
		return err.Error() + "\n", 1
	}

	f := collect(root)
	if override, present := parsed.Flags["--next"]; present {
		// An empty override names no command. Left to fall through it would read as a
		// clean board, which is a different — and false — statement about the tree.
		if override == "" {
			return toon.Usage(grammar.Cmd, `--next ""`) + "\n", 2
		}
		f.Action = override
	} else {
		action, signal, noneInvocable := nextAction(root)
		f.Action, f.Signal, f.NoInvocable = translate(action, harness), signal, noneInvocable
	}
	if err := validate(f); err != nil {
		return err.Error() + "\n", 1
	}

	pin := render(f, state)
	if err := os.WriteFile(target, []byte(document(pin)), 0o644); err != nil {
		return refusal{"cannot write the session handoff", err.Error()}.Error() + "\n", 1
	}
	return pin, 0
}

// preservedState reads the State section a run must carry forward. An absent file is the
// scaffold case — the first run in a repo owes it a skeleton, not an error — while a file
// that exists but cannot be read or split is left untouched.
func preservedState(target string) (string, error) {
	content, err := os.ReadFile(target)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", refusal{"cannot read the session handoff", err.Error()}
	}
	return splitState(content)
}
