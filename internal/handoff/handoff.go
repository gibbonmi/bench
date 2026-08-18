package handoff

import (
	"fmt"
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
	Cmd:  "bench handoff",
	Help: "usage: bench handoff [--harness " + status.HarnessChoices() + "] [--next <command>]",
	Flags: []usage.Flag{
		{Name: "--harness", HasValue: true, NoEmptyValue: true},
		// An empty override names no command. Left to fall through it would read as a
		// clean board, which is a different — and false — statement about the tree.
		{Name: "--next", HasValue: true, NoEmptyValue: true},
	},
}

// Command implements `bench handoff`. It prints the pin block and rewrites
// capture/session-handoff.md from one derivation, preserving the reviewer-owned State section and
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
	harness := status.HarnessClaude
	if value, present := parsed.Flags["--harness"]; present {
		if !status.ValidHarness(value) {
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
		f.Action = override
	} else {
		applyRoute(&f, status.RouteFor(root, status.SignalsWith(root, status.Query{ExcludeDirtyPaths: []string{status.HandoffFile}}), harness))
	}
	if err := validate(f); err != nil {
		return err.Error() + "\n", 1
	}

	pin := render(f, state)
	if err := writeDocument(target, document(pin)); err != nil {
		return refusal{"cannot write the session handoff", err.Error()}.Error() + "\n", 1
	}
	return pin, 0
}

// writeDocument replaces the handoff atomically: a full write to a sibling temp file, then
// a rename over the target. A plain write truncates before it writes, so a disk that fills
// or a signal that lands between the two leaves the reviewer-owned State section destroyed
// — the one thing this command promises to pass through untouched. The rename either
// happens or does not, which is what makes "destroys never" true rather than likely.
func writeDocument(target, content string) error {
	// The capture directory is part of the target's identity, not a precondition the
	// caller can be assumed to have met: a fresh repo reaches this write before anything
	// has created it.
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		// Name the document, not just the directory: the refusal's whole job is to say
		// which artifact could not be written.
		return fmt.Errorf("%s: %w", target, err)
	}
	temp := target + ".tmp"
	if err := os.WriteFile(temp, []byte(content), 0o644); err != nil {
		return err
	}
	if err := os.Rename(temp, target); err != nil {
		os.Remove(temp)
		return err
	}
	return nil
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
