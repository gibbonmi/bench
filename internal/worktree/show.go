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
	if len(args) != 2 || !revisionOperand(args[1]) {
		fmt.Fprintln(stderr, "usage: "+usage.WorktreeShow)
		return 2
	}
	path, err := resolveWorktree(root, args[0])
	if err != nil {
		return printTargetRefusal(stderr, "bench worktree show", err)
	}
	return runShowChild(path, args[1], stdout, stderr)
}

// revisionOperand reports whether the operand is the one shape Git may receive. The
// check runs before the target resolves, so a bare path or a dash operand never starts
// a child.
func revisionOperand(operand string) bool {
	return strings.Contains(operand, ":") && !strings.HasPrefix(operand, "-")
}

// runShowChild reads the blob through Git itself, and writes the child's own streams.
// The caller's writers take the bytes directly, because a line-oriented copy would alter
// a binary blob. The exit code follows the shared child rule, so an interrupt and a Git
// refusal report the same way exec does.
func runShowChild(dir, operand string, stdout, stderr io.Writer) int {
	ctx, stop := subprocess.NotifyCancel(context.Background())
	defer stop()
	cmd := exec.CommandContext(ctx, "git", "cat-file", "blob", operand)
	cmd.Dir, cmd.Stdout, cmd.Stderr = dir, stdout, stderr
	return childExitCode(cmd, cmd.Run())
}
