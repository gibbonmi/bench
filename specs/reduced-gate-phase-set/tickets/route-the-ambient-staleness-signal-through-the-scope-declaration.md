# Route the ambient staleness signal through the scope declaration

Blocked by: Declare the reduced gate scope

Ownership fence: `internal/status`
Assumptions: the declaration exports a confinement predicate `internal/status` can call, and the board's capture-only softening rule is otherwise unchanged

## What to build

`bench status` softens a stale gate to `capture-only drift` through the shared
declaration instead of its own `captureOnlyStalePaths` literal. The board's advice
and the oracle's behavior then cannot name different files, which is the whole point
of consolidating: today a path could be on the board's fast list and off the gate's,
or the reverse, with nothing to notice.

Delete the private map rather than leaving it beside the declaration. A literal that
survives as a second derivation of the same fact is the defect this ticket exists to
remove, and the acceptance criterion is written to fail if it stays.

## Acceptance

- [ ] [R01] Softening behavior follows the declaration: a path the declaration no longer carries stops softening, with no private literal left in `internal/status` still answering the question.
