# Migrate spec-build aggregates and empty state

Blocked by: none
Ownership fence: `internal/specbuild`, `cmd/bench/specbuild.go`, `cmd/bench/specbuild_test.go`
Integration surfaces: shared aggregate/empty carriers→implemented prerequisite `axi-carriers-and-registry`; final contraction→contract-aggregate-empty-routes.md
Contracts: spec-build integer disposition facts and empty-class enum cross `internal/specbuild`→`cmd/bench/specbuild.go`; domain is abandon/reclaim/no-run; order is current record order; absence is one-row state=empty plus full zero detail, asserted by SA1 and SA1E
Closure: SA1/ordered-counts, SA1/zero-classes, SA1/route, SA1/one-row-empty, SA1/full-detail-empty, SBE1/empty-class

## What to build

Migrate spec-build aggregates and empty state through the shared carriers without changing owner semantics or public bytes.

## Acceptance

- [ ] [SA1] (covers AE5) migrate spec-build aggregates and empty state preserve ordered counts, zero classes, route, one-row-empty, full-detail-empty.
- [ ] [SBE1] (covers AE8) migrate spec-build aggregates and empty state preserve their exact empty/absent rendering class.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SA1/ordered-counts | reorder one reclamation count | the independent spec-build renderer test | invoke reclaim and require current order |
| SA1/zero-classes | omit one zero class | the independent reclaim receipt test | invoke the empty-class fixture and require every disposition |
| SA1/route | reorder fields | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| SA1/one-row-empty | bypass shared aggregate | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| SA1/full-detail-empty | normalize one-row empty to zero rows | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| SBE1/empty-class | omit full detail empties | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| undefined | normalize the named empty class | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
