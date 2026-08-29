package worktree

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/usage"
)

// ShowCommand prints one blob from a revision of an active Bench-owned worktree. Git's
// stdout, stderr, and exit code pass through unchanged, so a NUL byte survives and a
// missing object reads as Git's own error rather than a Bench sentence. The operand is
// one `<rev>:<path>` token: it must hold a `:` and must not start with `-`, so no bare
// path reads as a revision and no Git option reaches the child.
func ShowCommand(root, _ string, args []string, stdout, stderr io.Writer) int {
	parsed, line, code := usage.Parse(worktreeShowGrammar, args)
	if line != "" {
		fmt.Fprintln(stderr, line)
		return code
	}
	if !revisionOperand(parsed.Positionals[1]) {
		fmt.Fprintln(stderr, worktreeShowGrammar.Help)
		return 2
	}
	path, err := resolveWorktree(root, parsed.Positionals[0])
	if err != nil {
		return printTargetRefusal(stderr, "bench worktree show", err)
	}
	return runShowChild(path, parsed.Positionals[1], stdout, stderr)
}

// worktreeShowGrammar reserves both positional slots, so a `--output=...` operand reads
// as the literal value Git would receive rather than as a flag. The sole `--help`
// spelling still answers with the grammar line.
var worktreeShowGrammar = usage.Grammar{
	Cmd:                                 "bench worktree show",
	Help:                                "usage: " + usage.WorktreeShow,
	MinArgs:                             2,
	MaxArgs:                             2,
	ReservedPositionalsBeforeTerminator: 2,
}

// revisionOperand reports whether the operand is the one shape Git may receive. The
// check runs before the target resolves, so a bare path, a dash operand, or a control
// byte never starts a child.
func revisionOperand(operand string) bool {
	return strings.Contains(operand, ":") && !strings.HasPrefix(operand, "-") && lineSafe(operand)
}

// runShowChild reads the blob through Git itself, and writes the child's own streams.
// The caller's writers take the bytes directly, because a line-oriented copy would alter
// a binary blob. The exit code follows the shared child rule.
func runShowChild(dir, operand string, stdout, stderr io.Writer) int {
	ctx, stop := subprocess.NotifyCancel(context.Background())
	defer stop()
	cmd := exec.CommandContext(ctx, "git", "cat-file", "blob", operand)
	cmd.Dir, cmd.Stdout, cmd.Stderr = dir, stdout, stderr
	return childExitCode(cmd, cmd.Run())
}
