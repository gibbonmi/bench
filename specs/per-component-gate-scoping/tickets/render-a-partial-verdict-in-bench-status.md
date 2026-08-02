# Render a partial verdict in bench status

Blocked by: Carry the partition in the verdict and refuse its reuse
Ownership fence: `internal/status/status.go`, `internal/status/status_test.go`
Assumptions: `status.go` today renders a reduced green as its own row
(`reduced green (capture-only scope)`), derives `Stale` from
`!ReusableGreen`, and reads narrowness off `gate.Inspection`. Re-derive from the
tree at pickup.

## What to build

A partial verdict is narrow, not drifted — the same distinction the reduced row
already draws, at component granularity. The board renders it as its own class,
naming what was skipped, so a reader cannot mistake a partition for a whole-tree
green and cannot mistake it for tree drift. Its action stays the operator's one
lever: `bench gate --fresh`.

A partial verdict whose tree has since moved is still drift, and falls through to
the drift row exactly as a reduced one does — narrowness and staleness stay
independent.

## Acceptance

- [ ] [PC17a] a partial green over its own tree renders as a partial row naming the skipped components, not as stale and not as full green.
- [ ] [PC17b] a partial green whose tree has moved renders the drift row, not the partial row.
- [ ] [PS25] the partial row's action is `bench gate --fresh`.
- [ ] [PS26] the reduced row still renders unchanged for a reduced verdict.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PC17a | fall through to the stale row when the record is partial | `TestStatusRendersAPartialVerdict` | seed a partial green for the work tree, read the board, assert the partial row and its skipped-component names |
| PC17b | return the partial row before comparing cached and work trees | `TestPartialVerdictOnAMovedTreeIsDrift` | seed a partial green, edit a tracked file, read the board, assert the drift row |
| PS25 | name `bench gate` as the action | `TestPartialRowActionIsFresh` | seed a partial green, read the board, compare the action string |
| PS26 | replace the reduced row with the partial row | existing reduced-row status test | seed a reduced green, read the board, assert the reduced row unchanged |
