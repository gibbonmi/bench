package main

import (
	"cmp"
	"errors"
	"fmt"
	"io"

	"github.com/gibbonmi/bench/internal/commit"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/worktree"
)

// chainSteps are the steps of `bench commit --preflight-build`, each with its owner's own
// signature, so a test replaces one step without touching the others. BuildInputs decides
// whether the root has a worktree build at all.
type chainSteps struct {
	Commit      func(args []string, stdout, stderr io.Writer, check func(root, candidate string) error) (commit.Outcome, int)
	BuildInputs func(root string) bool
	Build       func(root, home string, args []string, stdout, stderr io.Writer) int
	Preflight   func(root, slug, base, candidate string) (string, int)
}

// commitChain binds each step to its owner.
var commitChain = chainSteps{Commit: commit.RunChecked, BuildInputs: freshness.DeclaresBuildInputs, Build: worktree.BuildCommand, Preflight: preflight.BuildVerdictAt}

// chainHelp is the help line for the chain that commitChainCommand runs.
var chainHelp = commit.PreflightBuildGroup + ": before publication, run bench worktree build on this worktree when its root declares Bench build inputs, then bench preflight build <slug> at the candidate commit, with " + commit.BaseFlag + " when the commit names one; a red step publishes nothing; end with commit-chain{commit,build,preflight}"

// commitChainCommand runs the commit and, when --preflight-build names a slug, hands it a
// check that runs the worktree build and then the build preflight at the candidate commit
// before publication. A root that declares no Bench build inputs has no worktree build, so
// the check skips the build there and runs the preflight alone. Each step prints its own
// response, then one commit-chain line states every step. A red step refuses the commit, so
// the commit's exit is the chain's exit. A help answer gets the chain's help line after the
// commit's help.
func commitChainCommand(c Command, args []string) int {
	steps := commitChain
	slug, base := commit.ChainValues(args)
	if slug == "" {
		outcome, exit := steps.Commit(args, c.Stdout, c.Stderr, nil)
		if outcome.Help {
			fmt.Fprintln(c.Stdout, chainHelp)
		}
		return exit
	}
	build, preflight := "skipped", "skipped"
	check := func(root, candidate string) error {
		if steps.BuildInputs(root) {
			code := steps.Build(root, worktree.Home(), []string{root}, c.Stdout, c.Stderr)
			if build = stepState[code == 0]; code != 0 {
				return errors.New("the worktree build is red")
			}
		}
		out, code := steps.Preflight(root, slug, base, candidate)
		fmt.Fprint(c.Stdout, out)
		if preflight = stepState[code == 0]; code != 0 {
			return fmt.Errorf("the build preflight of %s is red at the candidate commit", slug)
		}
		return nil
	}
	outcome, exit := steps.Commit(args, c.Stdout, c.Stderr, check)
	fmt.Fprintf(c.Stdout, "commit-chain{commit=%s,build=%s,preflight=%s}\n", cmp.Or(outcome.Published, "none"), build, preflight)
	return exit
}

// stepState names a chained step that ran by whether it exited 0.
var stepState = map[bool]string{true: "green", false: "red"}
