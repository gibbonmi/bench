package main

import (
	"fmt"

	"github.com/gibbonmi/bench/internal/toon"
)

// treeScope states which tree a public command leaf serves. Each public definition that is
// not a family declares one, and so does each leaf row of a family. A family and a plumbing
// definition declare none. The zero value means undeclared: such a definition prints no
// tree row and takes no tree target, so a planted test registry keeps its output.
type treeScope uint8

const (
	// scopeTree marks a leaf that reads, grades, or writes tracked content of one checkout.
	scopeTree treeScope = iota + 1
	// scopeRepository marks a leaf that reads or writes only Git refs, the worktree
	// ledger, git-ignored capture files, or state under the Bench home.
	scopeRepository
)

// treeTargetFlag is the one spelling of the tree-target flag.
const treeTargetFlag = "--in"

// refusesTreeTarget answers the repository refusal. A repository-scoped definition that
// gets the tree-target flag as its first argument prints the usage line and does not run,
// so an ignored target never looks accepted. A family declares no scope, so its first
// argument still routes as a leaf name.
func (definition commandDefinition) refusesTreeTarget(c Command, args []string) bool {
	if definition.Scope != scopeRepository || len(args) == 0 || args[0] != treeTargetFlag {
		return false
	}
	fmt.Fprintln(c.Stdout, toon.Usage("bench "+definition.Name, treeTargetFlag))
	return true
}
