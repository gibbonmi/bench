package main

import (
	"fmt"

	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/treetarget"
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

// scope answers the scope of one call. A call that names a family leaf takes the scope of
// that leaf, and every other call takes the scope of the definition.
func (definition commandDefinition) scope(args []string) treeScope {
	if leaf, ok := leafNamed(definition.Leaves, args); ok {
		return leaf.Scope
	}
	return definition.Scope
}

// treeRow answers the identity row block of one call, or the empty string for a call that
// prints no row. Only a tree-scoped call that is not a help form and that runs inside a
// repository prints the row. The dispatcher computes it before the verb runs, so the row
// names the tree that the verb read. A label that TOON cannot carry prints no row.
func (definition commandDefinition) treeRow(args []string) string {
	if definition.scope(args) != scopeTree || definition.helpForm(args) {
		return ""
	}
	root := boundaryRoot()
	if root == "" {
		return ""
	}
	row, err := treetarget.Row(treetarget.Identify(root))
	if err != nil {
		return ""
	}
	return row
}

// shownRow answers the row that a call with exit prints. A grammar refusal at exit 2
// prints only its usage line, so it prints no row.
func shownRow(row string, exit int) string {
	if exit == 2 {
		return ""
	}
	return row
}

// finishExempt prints the row of an exempt call on stderr after the verb returns, so an
// artifact on stdout stays byte-clean. It answers the verb's exit.
func (c Command) finishExempt(row string, exit int) int {
	fmt.Fprint(c.Stderr, shownRow(row, exit))
	return exit
}
