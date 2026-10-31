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
	{Name: "read", Grammar: usage.WorktreeRead, Root: rootBoundary, Bound: boundResponse, Scope: scopeTree, Run: func(c Command, root string, args []string) int {
		out, code := worktreeReadCommand(root, args)
		fmt.Fprint(c.Stdout, out)
		return code
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

// worktreeHelpRows declares the root help rows of the `bench worktree` family beside its
// leaves, so the registry row carries no per-leaf help text.
func worktreeHelpRows() []helpRow {
	return []helpRow{
		{Order: 31, Suffix: " shell [--refresh] [objective]", Gap: 1, Description: "create an owned worktree subshell and release it on exit"},
		{Order: 32, Suffix: " list", Description: "list assignments and registered worktrees as TOON"},
		{Order: 32, Suffix: worktreeSuffix(usage.WorktreeListRelated), Description: "the assignments that commitment scopes relate to one spec as bounded TOON pages"},
		{Order: 33, Suffix: worktreeSuffix(usage.WorktreePath), Description: "print one active owned worktree's absolute path for the file tools"},
		{Order: 34, Suffix: worktreeSuffix(usage.WorktreeExec), Description: "run a child directly in an active owned worktree"},
		{Order: 34, Suffix: worktreeSuffix(usage.WorktreeShow), Description: "print one blob from a revision of an active owned worktree"},
		{Order: 34, Suffix: worktreeSuffix(usage.WorktreeRead), Description: "exact file, line-range, or H2 section bytes of one owned worktree as bounded TOON pages"},
		{Order: 34, Suffix: worktreeSuffix(usage.WorktreeBuild), Description: "build an active owned worktree's tree into its own dist/bench"},
		{Order: 34, Suffix: worktreeSuffix(usage.WorktreeExecGate), Description: "run one active owned worktree's gate"},
		{Order: 35, Suffix: " reauthorize --assignment <id> --request <token> --base <commit> --source-tip <commit> <path>", Description: "replace one lost request token after identity proof"},
		{Order: 36, Suffix: " merge --from <commit|target> <target>", Description: "merge a default-branch commit or a sibling's tip into an owned worktree"},
		{Order: 36, Suffix: worktreeSuffix(usage.WorktreeReset), Description: "plan or apply a recoverable reset or restore"},
		{Order: 36, Suffix: worktreeSuffix(usage.WorktreeLand), Description: "compose, gate, and publish one owned worktree"},
		{Order: 36, Suffix: worktreeSuffix(usage.WorktreeLandResume), Description: "resume incomplete post-publication landing work"},
		{Order: 36, Suffix: worktreeSuffix(usage.WorktreeLandStack), Description: "compose, gate once, and publish two or more owned worktrees in order"},
		{Order: 36, Suffix: worktreeSuffix(usage.WorktreeLandStackResume), Description: "resume incomplete post-publication stack work"},
		{Order: 37, Suffix: " --help", Description: "show exact list, path, exec, show, read, build, create, release, clean, reclaim, reauthorize, merge, reset, and land grammar"},
	}
}

// poolCommand supplies the resolved Bench home to the argv-only handler.
func poolCommand(args []string) (string, int) {
	return worktree.PoolCommand(worktree.Home(), args)
}

func resumeCleanCommand(c Command, args []string) int {
	return worktree.ResumeCleanCommand(boundaryRoot(), worktree.Home(), args, c.Stdout, c.Stderr)
}
