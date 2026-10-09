// Package treetarget names the tree that one tree-scoped Bench call reads. It answers the
// identity row that each tree-scoped response prints as its first block, so a reader
// knows which checkout the response graded without a path in the output. It also runs a
// call that names its tree target as one child in that tree, and a kit worktree target
// runs its own current build.
package treetarget

import (
	"bytes"
	"errors"
	"fmt"

	toonlib "github.com/toon-format/toon-go"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
)

// primaryTarget names the primary checkout, both as the target cell of the identity row
// and as the tree-target keyword.
const primaryTarget = "primary"

// rowName and rowFields are the block name and the field set of the identity row.
const rowName = "tree"

// unknownDirty is the dirty cell of a tree whose status query failed.
const unknownDirty = "unknown"

var rowFields = []string{"target", "head", "dirty"}

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
	identity := Identity{Target: "unassigned", Head: "none", Dirty: unknownDirty}
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
	return toon.TableTyped(rowName, rowFields, cells)
}

// ReadLead reads the identity row that a tree-scoped response may print first. Data that
// does not open with the row answers its zero identity, data itself, and false. Data that
// opens with the row answers the identity, the bytes after the row, and true. The body is a
// view of data, so the caller keeps the exact selected bytes. The row must be complete and
// canonical, and a second row must not follow it. The row names an observed scope only; it
// grants no source or execution authority.
func ReadLead(data []byte) (Identity, []byte, bool, error) {
	if !bytes.HasPrefix(data, []byte(rowName+"[")) {
		return Identity{}, data, false, nil
	}
	end := 0
	for range 2 {
		next := bytes.IndexByte(data[end:], '\n')
		if next < 0 {
			return Identity{}, nil, false, errors.New("the identity row is incomplete")
		}
		end += next + 1
	}
	lead, body := data[:end], data[end:]
	identity, err := decodeRow(lead)
	if err != nil {
		return Identity{}, nil, false, fmt.Errorf("the identity row is malformed: %w", err)
	}
	if canonical, err := Row(identity); err != nil || canonical != string(lead) {
		return Identity{}, nil, false, errors.New("the identity row is not canonical")
	}
	if bytes.HasPrefix(body, []byte(rowName+"[")) {
		return Identity{}, nil, false, errors.New("a second identity row follows the first")
	}
	return identity, body, true, nil
}

// decodeRow decodes one identity block through the pinned strict TOON decoder and checks
// each cell's type: a string target and head, and a bool or `unknown` dirty value.
func decodeRow(lead []byte) (Identity, error) {
	decoded, err := toonlib.Decode(lead)
	if err != nil {
		return Identity{}, err
	}
	document, ok := decoded.(map[string]any)
	if !ok || len(document) != 1 {
		return Identity{}, fmt.Errorf("decodes as %T, want one %s block", decoded, rowName)
	}
	rows, ok := document[rowName].([]any)
	if !ok || len(rows) != 1 {
		return Identity{}, fmt.Errorf("want one %s row", rowName)
	}
	cells, ok := rows[0].(map[string]any)
	if !ok || len(cells) != len(rowFields) {
		return Identity{}, fmt.Errorf("want the fields %v", rowFields)
	}
	target, okTarget := cells[rowFields[0]].(string)
	head, okHead := cells[rowFields[1]].(string)
	dirty := cells[rowFields[2]]
	if _, isBool := dirty.(bool); !okTarget || !okHead || !isBool && dirty != unknownDirty {
		return Identity{}, fmt.Errorf("cells %v have the wrong types", cells)
	}
	return Identity{Target: target, Head: head, Dirty: dirty}, nil
}
