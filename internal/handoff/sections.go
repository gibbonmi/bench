package handoff

import (
	"github.com/gibbonmi/bench/internal/toon"
)

// refusal is a structured error that carries the AXI kind/hint pair. It lets the command
// surface one rendering of a refusal rather than re-deriving the line at each exit.
type refusal struct{ kind, hint string }

func (r refusal) Error() string { return toon.Errorf(r.kind, r.hint) }
