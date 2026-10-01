package main

import (
	"fmt"
	"strings"

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
	{Name: "exec", Grammar: usage.WorktreeExec, Root: rootRequired, Bound: boundResponse, Scope: scopeRepository, Run: func(c Command, root string, args []string) int {
		exit, assignment := worktree.ExecCommandResolving(root, worktree.Home(), args, c.Stdin, c.Stdout, c.Stderr)
		c.reportAssignment(assignment)
		return exit
	}},
	{Name: "shell", Root: rootNone, Bound: boundExempt(boundReasonTerminal), Scope: scopeRepository, Run: func(c Command, _ string, args []string) int {
		return worktree.Subshell(worktree.Home(), args, c.Stdin, c.Stdout, c.Stderr)
	}},
	{Name: "path", Root: rootRequired, Bound: boundResponse, Scope: scopeRepository, Run: func(c Command, root string, args []string) int {
		return worktree.PathCommand(root, worktree.Home(), args, c.Stdout, c.Stderr)
	}},
	{Name: "show", Grammar: usage.WorktreeShow, Root: rootRequired, Bound: boundResponse, Scope: scopeRepository, Run: func(c Command, root string, args []string) int {
		return worktree.ShowCommand(root, worktree.Home(), args, c.Stdout, c.Stderr)
	}},
	{Name: "build", Grammar: usage.WorktreeBuild, Root: rootRequired, Bound: boundResponse, Scope: scopeRepository, Run: func(c Command, root string, args []string) int {
		return worktree.BuildCommand(root, worktree.Home(), args, c.Stdout, c.Stderr)
	}},
	{Name: "list", Grammar: usage.WorktreeList, Root: rootBoundary, Bound: boundResponse, Scope: scopeRepository, Run: func(c Command, root string, args []string) int {
		out, code := worktree.ListCommand(root, worktree.Home(), args)
		fmt.Fprint(c.Stdout, out)
		return code
	}},
	{Name: "create", Grammar: usage.WorktreeCreate, Root: rootRequired, Bound: boundResponse, Scope: scopeRepository, Run: func(c Command, root string, args []string) int {
		return worktree.CreateCommand(root, worktree.Home(), args, c.Stdout, c.Stderr)
	}},
	{Name: "release", Grammar: usage.WorktreeRelease, Root: rootRequired, Bound: boundResponse, Retires: true, Scope: scopeRepository, Run: func(c Command, root string, args []string) int {
		return worktree.ReleaseCommand(root, worktree.Home(), args, c.Stdout, c.Stderr)
	}},
	{Name: "clean", Grammar: usage.WorktreeClean, Root: rootBoundary, Bound: boundResponse, Retires: true, Scope: scopeRepository, Run: func(c Command, root string, args []string) int {
		return worktree.CleanCommand(root, worktree.Home(), args, c.Stdout, c.Stderr)
	}},
	{Name: "reclaim", Root: rootBoundary, Bound: boundResponse, Retires: true, Scope: scopeRepository, Run: func(c Command, root string, args []string) int {
		return worktree.ReclaimCommand(root, worktree.Home(), args, c.Stdout, c.Stderr)
	}},
	{Name: "reauthorize", Grammar: usage.WorktreeReauthorize, Root: rootRequired, Bound: boundResponse, Scope: scopeRepository, Run: func(c Command, root string, args []string) int {
		return worktree.ReauthorizeCommand(root, worktree.Home(), args, c.Stdout, c.Stderr)
	}},
	{Name: "merge", Grammar: usage.WorktreeMerge, Root: rootRequired, Bound: boundResponse, Scope: scopeRepository, Run: func(c Command, root string, args []string) int {
		return worktree.MergeCommand(root, worktree.Home(), args, c.Stdout, c.Stderr)
	}},
	{Name: "reset", Grammar: usage.WorktreeReset, Root: rootRequired, Bound: boundResponse, Scope: scopeRepository, Run: func(c Command, root string, args []string) int {
		return worktree.ResetCommand(root, worktree.Home(), args, c.Stdout, c.Stderr)
	}},
	{Name: "land", Grammar: usage.WorktreeLand, Root: rootRequired, Bound: boundResponse, Retires: true, Scope: scopeRepository, Run: func(c Command, root string, args []string) int {
		return worktree.LandCommand(root, worktree.Home(), args, c.Stdout, c.Stderr)
	}},
}

// worktreeSuffix derives a leaf's help suffix from its usage grammar, so the inventory
// row and the grammar the verb refuses with have one source.
func worktreeSuffix(grammar string) string { return strings.TrimPrefix(grammar, "bench worktree") }

// poolCommand supplies the resolved Bench home to the argv-only handler.
func poolCommand(args []string) (string, int) {
	return worktree.PoolCommand(worktree.Home(), args)
}

func resumeCleanCommand(c Command, args []string) int {
	return worktree.ResumeCleanCommand(boundaryRoot(), worktree.Home(), args, c.Stdout, c.Stderr)
}
