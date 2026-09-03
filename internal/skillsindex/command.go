package skillsindex

import (
	"context"
	"errors"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// grammar is the declared argument shape usage.Parse enforces for `bench skills-index`.
// `--check`, the default, prints the diagnostics Check would emit; `--write` regenerates
// the block in place. Declaring both as plain flags, rather than encoding a mode string,
// keeps a mistyped third spelling a usage error rather than a silently ignored no-op.
var grammar = usage.Grammar{
	Cmd:   "bench skills-index",
	Help:  "usage: bench skills-index [--check|--write]",
	Flags: []usage.Flag{{Name: "--check"}, {Name: "--write"}},
}

// Command implements `bench skills-index [--check|--write]`, the operator's regenerator
// and drift check. `--check`, the default, prints Check's diagnostics one per line and
// exits 1 if any, 0 if clean; `--write` regenerates the block via Write. A blocking
// refusal, the allowlist unparseable or the reference file missing or unmarked, prints
// Write's own error and exits 1 without the file being touched. Write refuses before
// any bytes are written, so the caller and the file agree on nothing having changed.
// The two modes name opposite intents, so asking for both is a usage error rather than
// a silent precedence rule. This is refused before repository discovery, so the verdict
// on the arguments does not depend on where the caller stands.
func Command(args []string) (string, int) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return line + "\n", code
	}
	_, write := parsed.Flags["--write"]
	if _, check := parsed.Flags["--check"]; check && write {
		return grammar.Help + " (--check and --write are mutually exclusive)\n", 2
	}
	root, err := git.Root()
	if err != nil {
		return discoveryRefusal(err) + "\n", 1
	}
	if write {
		// The verb owns the termination signals for exactly as long as it is replacing bytes.
		// An operator's Ctrl-C during a write is an instruction to abandon it, and the default
		// handler would take the process down mid-replacement instead. Session teardown sends
		// SIGTERM or SIGHUP rather than SIGINT, so the trapped set is subprocess.CancelSignals.
		ctx, stop := signal.NotifyContext(context.Background(), subprocess.CancelSignals...)
		defer stop()
		if werr := Write(ctx, root); werr != nil {
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

// discoveryRefusal names the recovery action the caller actually has. A discovery error
// that never launched Git is an environment defect, not a location one. Telling an
// operator with no `git` on PATH to stand somewhere else sends them after the wrong fix.
// Only an executed probe reports the position it measured.
func discoveryRefusal(err error) string {
	var launch *exec.Error
	if errors.As(err, &launch) {
		return toon.Errorf("required tool is missing or not executable: "+launch.Name, "install git and re-run")
	}
	return toon.NotInRepo()
}
