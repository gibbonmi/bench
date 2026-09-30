package treetarget

import (
	"errors"
	"fmt"
	"io"

	"github.com/gibbonmi/bench/internal/canonicalpath"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/worktree"
)

// Call is one tree-scoped call whose first argument is the tree-target flag.
type Call struct {
	// Name is the verb, and Flag is the tree-target flag that the dispatcher consumed.
	Name, Flag string
	// Args holds the flag value and then the verb's own arguments.
	Args []string
	// Root is the repository root of the calling directory, or empty outside a repository.
	Root string
	// Home is the Bench home that the child runs under.
	Home string
	// Wrapper is the invoking wrapper that the parent's environment names, or empty for a
	// direct run of the executable.
	Wrapper string
	// Running is the absolute path of the running executable.
	Running        string
	Stdin          io.Reader
	Stdout, Stderr io.Writer
}

// Run resolves the tree target of call and runs the verb there as one child. It answers the
// child's exit, or the exit of the refusal that starts no child. The parent prints no row
// and applies no response bound: the child names its own tree and bounds its own response.
func Run(call Call) int {
	command := "bench " + call.Name + " " + call.Flag
	value, rest, usage := targetValue(command, call.Args)
	if usage != "" {
		fmt.Fprintln(call.Stdout, usage)
		return 2
	}
	if call.Root == "" {
		fmt.Fprintln(call.Stderr, toon.NotInRepo())
		return 1
	}
	target, err := targetRoot(call.Root, value)
	if errors.Is(err, worktree.ErrTreeTargetPath) {
		fmt.Fprintln(call.Stdout, valueUsage(command, value))
		return 2
	}
	if err != nil {
		return worktree.PrintTreeTargetRefusal(call.Stderr, command, err)
	}
	// A root that the ledger recorded through a link runs in its physical spelling, so the
	// child reads one tree whatever spelling reached it.
	dir, err := canonicalpath.Resolve(target)
	if err != nil {
		return worktree.PrintTreeTargetRefusal(call.Stderr, command, err)
	}
	argv := append([]string{childExecutable(call), call.Name}, rest...)
	return worktree.RunTreeChild(command, argv, dir, call.Home, call.Stdin, call.Stdout, call.Stderr)
}

// targetRoot answers the root of the tree that value names. The keyword names the primary
// checkout, the first registered worktree, and every other value names a worktree label.
func targetRoot(root, value string) (string, error) {
	if value != primaryTarget {
		return worktree.TreeTarget(root, value)
	}
	worktrees, err := git.Worktrees(root)
	if err != nil {
		return "", err
	}
	if len(worktrees) == 0 {
		return "", errors.New("the repository registers no primary checkout")
	}
	return worktrees[0].Path, nil
}

// childExecutable answers the executable that the child runs. The invoking wrapper resolves
// the kit of the tree that it runs in, so it wins. A direct run of the executable has no
// wrapper, and the child runs the same executable.
func childExecutable(call Call) string {
	if call.Wrapper != "" {
		return call.Wrapper
	}
	return call.Running
}
