package main

import (
	"fmt"

	"github.com/gibbonmi/bench/internal/usage"
	"github.com/gibbonmi/bench/internal/worktree"
)

// worktreeLeaves declares the `bench worktree` family once. A slice that adds or moves a
// leaf edits this row and nothing else: the row is the only place that pairs the leaf's
// name with its grammar, its root need, and its handler.
//
// The `path` and `reclaim` rows name no grammar because neither answers `--help` with
// its own grammar today: `path` reads `--help` as a target operand, and `reclaim`
// refuses it as an unknown argument. `shell` has no grammar constant at all.
var worktreeLeaves = []commandLeaf{
	{Name: "exec", Grammar: usage.WorktreeExec, Root: rootRequired, Run: func(c Command, root string, args []string) int {
		return worktree.ExecCommand(root, worktree.Home(), args, c.Stdin, c.Stdout, c.Stderr)
	}},
	{Name: "shell", Root: rootNone, Run: func(c Command, _ string, args []string) int {
		return worktree.Subshell(worktree.Home(), args, c.Stdin, c.Stdout, c.Stderr)
	}},
	{Name: "path", Root: rootRequired, Run: func(c Command, root string, args []string) int {
		return worktree.PathCommand(root, worktree.Home(), args, c.Stdout, c.Stderr)
	}},
	{Name: "show", Grammar: usage.WorktreeShow, Root: rootRequired, Run: func(c Command, root string, args []string) int {
		return worktree.ShowCommand(root, worktree.Home(), args, c.Stdout, c.Stderr)
	}},
	{Name: "build", Grammar: usage.WorktreeBuild, Root: rootRequired, Run: func(c Command, root string, args []string) int {
		return worktree.BuildCommand(root, worktree.Home(), args, c.Stdout, c.Stderr)
	}},
	{Name: "list", Grammar: usage.WorktreeList, Root: rootBoundary, Run: func(c Command, root string, args []string) int {
		out, code := worktree.ListCommand(root, worktree.Home(), args)
		fmt.Fprint(c.Stdout, out)
		return code
	}},
	{Name: "create", Grammar: usage.WorktreeCreate, Root: rootRequired, Run: func(c Command, root string, args []string) int {
		return worktree.CreateCommand(root, worktree.Home(), args, c.Stdout, c.Stderr)
	}},
	{Name: "release", Grammar: usage.WorktreeRelease, Root: rootRequired, Run: func(c Command, root string, args []string) int {
		return worktree.ReleaseCommand(root, worktree.Home(), args, c.Stdout, c.Stderr)
	}},
	{Name: "clean", Grammar: usage.WorktreeClean, Root: rootBoundary, Run: func(c Command, root string, args []string) int {
		return worktree.CleanCommand(root, worktree.Home(), args, c.Stdout, c.Stderr)
	}},
	{Name: "reclaim", Root: rootBoundary, Run: func(c Command, root string, args []string) int {
		return worktree.ReclaimCommand(root, worktree.Home(), args, c.Stdout, c.Stderr)
	}},
	{Name: "reauthorize", Grammar: usage.WorktreeReauthorize, Root: rootRequired, Run: func(c Command, root string, args []string) int {
		return worktree.ReauthorizeCommand(root, worktree.Home(), args, c.Stdout, c.Stderr)
	}},
	{Name: "merge", Grammar: usage.WorktreeMerge, Root: rootRequired, Run: func(c Command, root string, args []string) int {
		return worktree.MergeCommand(root, worktree.Home(), args, c.Stdout, c.Stderr)
	}},
	{Name: "land", Grammar: usage.WorktreeLand, Root: rootRequired, Run: func(c Command, root string, args []string) int {
		return worktree.LandCommand(root, worktree.Home(), c.Executable, args, c.Stdout, c.Stderr)
	}},
}

// worktreeCommand routes the worktree family through the shared leaf dispatcher. The
// family answers its own bare, help, and unknown-leaf forms before any leaf runs, so a
// missing or unknown leaf never acquires or creates a worktree.
func worktreeCommand(c Command, args []string) int {
	return dispatchLeafFamily(c, "bench worktree", usage.WorktreeUsage(), worktreeLeaves, args)
}
