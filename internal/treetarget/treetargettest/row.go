// Package treetargettest holds the test support that reads the identity row of a response.
// Only tests import this package.
package treetargettest

import (
	"strings"

	"github.com/gibbonmi/bench/internal/treetarget"
)

// WithoutRow answers out without its leading identity row block, the header line and its
// one row. A response that does not start with the block comes back unchanged. The header
// comes from the renderer, so a fixture that compares the rest of a tree-scoped response
// follows a changed header with no edit.
func WithoutRow(out string) string {
	block, err := treetarget.Row(treetarget.Identity{Dirty: false})
	if err != nil {
		panic(err)
	}
	header, _, _ := strings.Cut(block, "\n")
	rest, ok := strings.CutPrefix(out, header+"\n")
	if !ok {
		return out
	}
	_, after, _ := strings.Cut(rest, "\n")
	return after
}
