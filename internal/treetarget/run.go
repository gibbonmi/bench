package treetarget

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gibbonmi/bench/internal/canonicalpath"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/worktree"
)

// Call is one tree-scoped call whose first argument is the tree-target flag.
type Call struct {
	// Name is the verb, and Flag is the tree-target flag that the dispatcher consumed.
	Name, Flag string
	// Leaf is the family leaf that the call names before the flag, or empty for a root
	// verb. The child argv holds the verb and the leaf as two tokens.
	Leaf string
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
	words := []string{call.Name}
	if call.Leaf != "" {
		words = append(words, call.Leaf)
	}
	command := "bench " + strings.Join(words, " ") + " " + call.Flag
	value, rest, usage := targetValue(command, call.Args)
	if usage != "" {
		fmt.Fprintln(call.Stdout, usage)
		return 2
	}
	if call.Root == "" {
		fmt.Fprintln(call.Stderr, toon.NotInRepo())
		return 1
	}
	dir, err := Resolve(call.Root, value)
	if errors.Is(err, worktree.ErrTreeTargetPath) {
		fmt.Fprintln(call.Stdout, valueUsage(command, value))
		return 2
	}
	if err != nil {
		return worktree.PrintTreeTargetRefusal(call.Stderr, command, err)
	}
	executable, refusal := childExecutable(call, value, dir)
	if refusal != "" {
		return worktree.PrintTreeBuildRefusal(call.Stderr, command, value, refusal)
	}
	argv := append(append([]string{executable}, words...), rest...)
	return worktree.RunTreeChild(command, argv, dir, call.Home, call.Stdin, call.Stdout, call.Stderr)
}

// Resolve answers the root of the tree that the tree-target value names, from root, the
// repository root of any checkout. The keyword names the primary checkout, the first
// registered worktree, and every other value names a worktree label. A root that the ledger
// recorded through a link answers its physical spelling, so a child reads one tree whatever
// spelling reached it, and a caller compares two roots by their resolved spellings.
func Resolve(root, value string) (string, error) {
	var target string
	var err error
	if value != primaryTarget {
		target, err = worktree.TreeTarget(root, value)
	} else {
		target, err = git.PrimaryCheckout(root)
	}
	if err != nil {
		return "", err
	}
	return canonicalpath.Resolve(target)
}
