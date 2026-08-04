# Single-source the parser's empty-argument operation list

Blocked by: none
Ownership fence: `internal/spec`
Contracts: the spec-build operation set crosses `internal/spec/build.go`'s `buildOperations` table→the empty-argument diagnostic in the same file, and is asserted by PO1 by adding an operation to the table and observing the diagnostic follow it; ordering is the diagnostic's own declared order and is asserted by PO3 so the rendered list stays stable rather than map-random
Assumptions: `buildOperations` remains the one declaration of which operations exist, as the existing bidirectional check in `build_test.go` already asserts against the test's invocation table; the diagnostic's wording and its `|` separator are the surface being preserved, not redesigned; claims re-derived from the tree at pickup

## What to build

`ParseBuild`'s empty-argument branch hand-writes the operation set as a literal:
`"start|assign|checkpoint|integrate|review|status|promote|abandon|reclaim"`. That is a
third derivation of a fact `buildOperations` already declares. The previous repair round
restored the single source between the grammar table and the test's invocation table with a
bidirectional check; this literal was deferred with evidence and is what remains.

The failure it produces is silent and one-directional: an operation added to
`buildOperations` gets a working grammar, a working parse, and a working command, while
`bench spec build` with no arguments keeps advertising the old set. Nothing goes red,
because the existing check compares the grammar table against the test's invocations and
never reads this string. The operator is told the operation does not exist by the one
surface whose entire job is to say which ones do.

Derive the diagnostic's list from `buildOperations`. Ordering has to be decided rather than
inherited: Go map iteration is randomized, so a naive derivation would make the message
differ run to run and turn any assertion on it flaky. Give the operations a declared order —
the lifecycle order the list already reads in, not alphabetical, since that order is what
makes the message legible — and derive both the table and the message from it.

## Acceptance

- [ ] [PO1] adding an operation to `buildOperations` makes it appear in the no-argument diagnostic without editing any second list.
- [ ] [PO2] removing an operation from `buildOperations` removes it from that diagnostic the same way.
- [ ] [PO3] the diagnostic renders its operations in a stable declared order across repeated runs, not in map order.
- [ ] [PO4] the no-argument diagnostic keeps its existing wording, separator, and exit code, so the change is single-sourcing and not a surface redesign.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PO1 | restore the hand-written literal in the empty-argument branch | the derived-operation-list test | inline the literal, run `go test ./internal/spec -count=1 -timeout 120s`, expect the added-operation assertion to fail |
| PO2 | derive the list from a fixed slice that outlives a removed table row | the derived-operation-list test | pin the slice, run `go test ./internal/spec -count=1 -timeout 120s`, expect the removed-operation assertion to fail |
| PO3 | derive the list by ranging the map directly | the stable-order test | range `buildOperations` without sorting, run `go test ./internal/spec -count=20 -timeout 120s`, expect the repeated-render assertion to fail |
| PO4 | change the separator from `\|` to `, ` | the diagnostic-surface test | swap the separator, run `go test ./internal/spec -count=1 -timeout 120s`, expect the wording assertion to fail |
