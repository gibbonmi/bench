package main

import (
	"cmp"
	"fmt"
	"io"

	"github.com/gibbonmi/bench/internal/commit"
	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/preflight/evidencecmd"
	"github.com/gibbonmi/bench/internal/worktree"
)

// chainSteps are the steps of `bench commit --preflight-build`, each with its owner's own
// signature, so a test replaces one step without touching the others.
type chainSteps struct {
	Commit    func(args []string, stdout, stderr io.Writer) (commit.Outcome, int)
	Build     func(root, home string, args []string, stdout, stderr io.Writer) int
	Preflight func(args []string) (string, int)
}

// commitChain binds each step to its owner.
var commitChain = chainSteps{Commit: commit.Run, Build: worktree.BuildCommand, Preflight: preflight.CommandWithVersion(version)}

// chainHelp is the help line for the chain that commitChainCommand runs.
const chainHelp = commit.PreflightBuildFlag + " <slug>: after a commit that publishes, run bench worktree build on this worktree, then bench preflight build <slug> at the published commit; end with commit-chain{commit,build,preflight} and exit with the first non-zero step exit"

// commitChainCommand runs the commit and, when --preflight-build names a slug, the worktree
// build and the build preflight at the published commit. Each step prints its own
// response, then one commit-chain line states every step. The exit is the first non-zero
// step exit. A help answer gets the chain's help line after the commit's help.
func commitChainCommand(c Command, args []string) int {
	steps := commitChain
	outcome, exit := steps.Commit(args, c.Stdout, c.Stderr)
	if outcome.Help {
		fmt.Fprintln(c.Stdout, chainHelp)
	}
	if outcome.PreflightBuild == "" {
		return exit
	}
	build, preflight := "skipped", "skipped"
	// Exit 3 published a commit that the checkout does not match, so a build would grade
	// an unreconciled tree. It stops the chain the way a refusal does.
	if exit == 0 {
		exit = steps.Build(outcome.Root, worktree.Home(), []string{outcome.Root}, c.Stdout, c.Stderr)
		build = stepState[exit == 0]
		if exit == 0 {
			out, code := steps.Preflight([]string{evidencecmd.ModeBuild, outcome.PreflightBuild, evidencecmd.FlagTip, outcome.Published})
			fmt.Fprint(c.Stdout, out)
			exit, preflight = code, stepState[code == 0]
		}
	}
	fmt.Fprintf(c.Stdout, "commit-chain{commit=%s,build=%s,preflight=%s}\n", cmp.Or(outcome.Published, "none"), build, preflight)
	return exit
}

// stepState names a chained step that ran by whether it exited 0.
var stepState = map[bool]string{true: "green", false: "red"}
