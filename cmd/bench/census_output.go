package main

import (
	"time"

	"github.com/gibbonmi/bench/internal/census"
	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/worktree"
)

// outputHead names the verb of one call in its census output record: `bench`, the command
// name, and the leaf word when the call names a leaf of the command's family.
func outputHead(definition commandDefinition, leaf commandLeaf) string {
	if leaf.Name == "" {
		return "bench " + definition.Name
	}
	return "bench " + definition.Name + " " + leaf.Name
}

// reportAssignment hands the dispatcher the assignment that a verb resolved for its own
// target. A call outside the bounded dispatch has no receiver, and the report is dropped.
func (c Command) reportAssignment(id string) {
	if c.resolved != nil {
		*c.resolved = id
	}
}

// recordOutput appends the census output record of one finished bounded response. The
// record takes the assignment the verb reported for its target, or else the assignment of
// the working tree. A process with neither writes no record. The dispatcher discards the
// error: the census is evidence beside the verb, so a failed write changes neither the
// response nor the exit code.
func recordOutput(size responsebound.Size, head, resolved string, retiring bool) error {
	root := boundaryRoot()
	if root == "" {
		return nil
	}
	home := worktree.Home()
	assignment := resolved
	if assignment == "" {
		scope, ok := responsebound.AssignmentScope(home, root, retiring)
		if !ok {
			return nil
		}
		assignment = scope
	}
	output := census.Output{Head: head, Lines: size.Lines, Bytes: size.Bytes, Spilled: size.Spilled}
	return census.RecordOutput(home, root, assignment, output, time.Now())
}
