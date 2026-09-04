package handoff

import (
	"github.com/gibbonmi/bench/internal/toon"
)

// ShapeHeading is the heading above the document's Shape body. The conformance check
// reads the Shape text back out of the artifact, so it must locate the heading by the
// same string the writer emits rather than restating it.
const ShapeHeading = "## Shape"

// refusal is a structured error that carries the AXI kind/hint pair. It lets the command
// surface one rendering of a refusal rather than re-deriving the line at each exit.
type refusal struct{ kind, hint string }

func (r refusal) Error() string { return toon.Errorf(r.kind, r.hint) }
