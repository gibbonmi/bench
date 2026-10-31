package main

import (
	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/census"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/responsebound"
)

// rawReadAdvice renders the advisory line for one eligible raw read. It names the native
// file read, or the line-range read when the raw read prints one range, from the read
// grammar and the line-range owner. The line advises only: the raw read still runs.
func rawReadAdvice(read census.Read) string {
	form := readGrammar.Cmd + " " + axi.ShellQuote(read.Path)
	if read.Ranged() {
		form += " " + responsebound.LinesFlag + " " + responsebound.LineRange(read.First, read.Last)
	}
	return "ADVICE: " + form + " reads the same text through the native reader. This command still runs."
}

// readOwner answers the verified assignment that owns the tree at root, for the census
// record of an eligible read. A read that is not eligible, and a read in the primary
// checkout or in an unowned tree, get no assignment.
func readOwner(root string, eligible bool) string {
	if !eligible {
		return ""
	}
	owner, ok := intent.AssignmentForWorktree(root)
	if !ok {
		return ""
	}
	return owner.ID
}
