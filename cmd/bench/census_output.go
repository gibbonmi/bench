package main

import (
	"sync"
	"time"

	"github.com/gibbonmi/bench/internal/census"
	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/worktree"
)

// runBounded runs one bounded call. The command writes both streams into one owner, which
// prints the response after the command returns, and the command's own exit code stays the
// verb's exit. The spill scope and the census output record share one root lookup. It runs
// at most once and only when first needed, so a spill that opens after the verb removed its
// own tree finds no root and takes the capped store outside any repository.
func (c Command) runBounded(definition commandDefinition, args []string) int {
	leaf, _ := leafNamed(definition.Leaves, args)
	root := sync.OnceValue(boundaryRoot)
	owner := responsebound.New(worktree.Home(), c.Stdout, c.Stderr, root, leaf.Retires)
	var resolved string
	c.Stdout, c.Stderr, c.resolved = owner.Stdout(), owner.Stderr(), &resolved
	exit := definition.run(c, args)
	owner.Finish()
	_ = recordOutput(owner.Size(), outputHead(definition, leaf), resolved, root())
	return exit
}

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
// the working tree at root. A process with neither writes no record, and neither does an
// assignment that the ledger does not hold as active after the verb returns. The dispatcher
// discards the error: the census is evidence beside the verb, so a failed write changes
// neither the response nor the exit code.
func recordOutput(size responsebound.Size, head, resolved, root string) error {
	if root == "" {
		return nil
	}
	home := worktree.Home()
	assignment := resolved
	if assignment == "" {
		scope, ok := responsebound.AssignmentScope(home, root)
		if !ok {
			return nil
		}
		assignment = scope
	}
	// A verb can retire the assignment it ran in or reported, and the retirement drops the
	// census files. So the record goes only to an assignment the ledger still holds active.
	if !worktree.AssignmentActive(root, assignment) {
		return nil
	}
	output := census.Output{Head: head, Lines: size.Lines, Bytes: size.Bytes, Spilled: size.Spilled}
	return census.RecordOutput(home, root, assignment, output, time.Now())
}
