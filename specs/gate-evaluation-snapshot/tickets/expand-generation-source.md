# Expand immutable generation sources

Blocked by: none
Ownership fence: `internal/gate/`
Contracts: the immutable generation produced in `internal/gate/tree_snapshot.go` crosses source adapters→identity consumers inside `internal/gate/`, asserted by EG1-EG4 against real Git objects

## What to build

Add the working-tree and prospective-tree source adapters plus the immutable parsed generation and object-keyed blob cache beside the current capture paths, without changing gate selection or evidence behavior yet.

## Acceptance

- [ ] [EG1] Working and prospective sources produce the same ordered entry, path-lookup, and blob-result contract from real Git tree objects.
- [ ] [EG2] Repeated requests for one blob object, including requests through different paths, perform one read and return immutable views.
- [ ] [EG3] A blob-read failure is memoized by object identity and returned identically without retrying inside the generation.
- [ ] [EG4] Two captures of the same unchanged tree produce distinct generation instances rather than sharing mutable cache state.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| EG1 | route the prospective source through working-tree materialization | the source-adapter contract test | build one unpublished tree, capture through both adapters, run the focused `internal/gate` test, expect the prospective source-count mismatch |
| EG2 | key the blob cache by path instead of object identity | the shared-object memoization test | request two paths with one object, run the focused test, expect two reads |
| EG3 | retry a cached blob error | the memoized-error test | make the object reader fail, request the object twice, run the focused test, expect the read-count mismatch |
| EG4 | reuse one generation pointer across two captures | the generation-lifecycle test | capture the same tree twice, run the focused test, expect the distinct-instance assertion |
