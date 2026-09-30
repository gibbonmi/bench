// Package treetarget names the tree that one tree-scoped Bench call reads. It answers the
// identity row that each tree-scoped response prints as its first block, so a reader
// knows which checkout the response graded without a path in the output. It also runs a
// call that names its tree target as one child in that tree.
package treetarget

import (
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
)

// primaryTarget names the primary checkout, both as the target cell of the identity row
// and as the tree-target keyword.
const primaryTarget = "primary"

// Identity holds the three cells of one identity row.
type Identity struct {
	// Target is `primary` for the primary checkout, the label of the active assignment
	// that owns the tree, or `unassigned` for any other tree.
	Target string
	// Head is the full HEAD commit, or `none` when HEAD does not resolve.
	Head string
	// Dirty is the bool that the porcelain status query answers, or the string `unknown`
	// when that query fails. The bool prints bare, where a string `true` prints quoted.
	Dirty any
}

// Identify answers the identity of the tree at root. Each cell has a fallback value, so a
// failed query names its failure in the row and never fails the call.
func Identify(root string) Identity {
	identity := Identity{Target: "unassigned", Head: "none", Dirty: "unknown"}
	if primary, err := git.IsPrimaryCheckout(root); err == nil && primary {
		identity.Target = primaryTarget
	} else if assignment, ok := intent.AssignmentForWorktree(root); ok {
		identity.Target = assignment.Label
	}
	if head, err := git.ResolveCommit(root, "HEAD"); err == nil {
		identity.Head = head
	}
	if dirty, err := git.WorktreeDirty(root); err == nil {
		identity.Dirty = dirty
	}
	return identity
}

// Row renders the identity as the block `tree[1]{target,head,dirty}:` with its one row.
// The ledger accepts a label with a control byte, and TOON refuses some of those bytes. So
// the label renders through sanitize.Strip, the duty for a table cell, and the row always
// prints. The other cells come from Git or from fixed words, so the encoder refuses none.
func Row(identity Identity) (string, error) {
	cells := [][]any{{sanitize.Strip(identity.Target), identity.Head, identity.Dirty}}
	return toon.TableTyped("tree", []string{"target", "head", "dirty"}, cells)
}
