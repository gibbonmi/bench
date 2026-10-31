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

// eligibleRead is the one eligibility rule for a raw read, and it decides both the advice and
// the owned census record. A read is eligible when census.ObserveRead observed its shape and
// its file resolves inside the repository at root. The path owner is the one that bench read
// refuses an outside operand with. A read outside the repository, such as a pool path read
// from another checkout, gets no advice and no owned record.
func eligibleRead(root string, read census.Read, observed bool) bool {
	if !observed {
		return false
	}
	_, _, ok := repositoryPath(root, read.Path, "")
	return ok
}

// readOwner answers the verified assignment that owns the tree at root, for the census
// record of an eligible read. A read that eligibleRead refuses, and a read in the primary
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
