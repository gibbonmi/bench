package worktree

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/gibbonmi/bench/internal/git"
)

// Subshell implements `bench worktree`: it acquires a warm, isolated pool worktree,
// drops the user into an interactive shell there inheriting the caller's stdio, and
// releases the worktree on any exit status. The reset is "soft" so a stale resetRef
// (origin/<branch> that no longer resolves) still yields a usable worktree rather than
// failing the drop-in.
func Subshell(stdin io.Reader, stdout, stderr io.Writer) int {
	root, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, "not in a git repo")
		return 1
	}
	wt, err := Acquire(root, "origin/"+git.DefaultBranch(root), "soft")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stderr, "🪵 worktree: %s  (exit to release)\n", wt)
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "bash"
	}
	cmd := exec.Command(shell)
	cmd.Dir = wt
	cmd.Stdin, cmd.Stdout, cmd.Stderr = stdin, stdout, stderr
	_ = cmd.Run() // release regardless of the subshell's exit status
	Release(wt)
	fmt.Fprintln(stderr, "🪵 released")
	return 0
}
