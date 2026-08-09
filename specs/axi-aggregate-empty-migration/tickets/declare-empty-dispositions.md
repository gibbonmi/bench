# Declare every production empty disposition

Blocked by: migrate-specbuild-aggregates-empty.md, migrate-status-dashboard-aggregates-empty.md, migrate-toon-query-empty-classes.md
Ownership fence: `cmd/bench`, `internal/conformance`, `internal/axi`, `projects/benchkit.md`
Integration surfaces: shared aggregate/empty carriers→implemented prerequisite `axi-carriers-and-registry`; final contraction→contract-aggregate-empty-routes.md
Contracts: member name and empty-class enum cross production registry→`internal/conformance`; membership is every root/nested member; order is registry order; absence is explicit no-semantic-empty, asserted by ED1
Closure: ED1/every-member, ED1/exact-class, ED1/explicit-no-empty, ED1/route

## What to build

Declare every production empty disposition through the shared carriers without changing owner semantics or public bytes.

## Acceptance

- [ ] [ED1] (covers AE9) declare every production empty disposition preserve every-member, exact-class, explicit-no-empty, route.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| ED1/every-member | omit one member | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| ED1/exact-class | default one class | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| ED1/explicit-no-empty | replace explicit no-empty | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| ED1/route | bypass one empty-capable producer | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |

