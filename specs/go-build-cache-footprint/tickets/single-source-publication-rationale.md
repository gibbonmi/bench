# Single-source the publication rationale in the gate package

Blocked by: none
Ownership fence: `internal/gate/gate_go.go`, `internal/gate/phases.go`
Integration surfaces: none
Contracts: none crosses

## What to build

Repair for review finding `standards-01-stale-non-atomic-rationale`. Both
`GateGoArgv`'s comment and the `BenchkitPhases` comment claim `go build`
replaces `dist/bench` non-atomically — stale twice over: publication now goes
through the staged sealed operation, and the same fact is derived in two
places. Collapse to one source: state each comment's own surviving reason (the
build phase owns the only write to `dist/bench`; `go run` avoids depending on
the build phase) and defer publication semantics to the freshness owner rather
than restating them, so the transaction repair cannot strand a second stale
copy.

## Acceptance

- [ ] [SR1] No comment in `internal/gate` claims non-atomic replacement of `dist/bench`, and publication semantics are stated by at most one source in the package, which defers to the freshness owner rather than restating its mechanics.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SR1 | restore the stale non-atomic replacement sentence in either file | the review sweep over the gate package | run `rg -n 'non-atomically' internal/gate`, expect a match where the criterion requires none |
