# Contract legacy aggregate and empty routes

Blocked by: migrate-guard-aggregates.md, migrate-outline-aggregates.md, migrate-roadmap-aggregates.md, migrate-worktree-aggregates.md, migrate-specbuild-aggregates-empty.md, migrate-publication-aggregates.md, migrate-status-dashboard-aggregates-empty.md, migrate-toon-query-empty-classes.md, declare-empty-dispositions.md
Ownership fence: `internal/conformance`, `internal/axi`, `internal/guards`, `internal/outline`, `internal/roadmap`, `internal/worktree`, `internal/specbuild`, `cmd/bench`, `internal/publication`, `internal/status`, `internal/dashboard`, `internal/toon`, `projects/benchkit.md`
Integration surfaces: shared aggregate/empty carriers→implemented prerequisite `axi-carriers-and-registry`; final contraction→contract-aggregate-empty-routes.md
Contracts: owner/member name, route identity, empty disposition, and legacy symbol census cross source census/production registry→`internal/conformance`; membership is every named aggregate and empty owner; order is production order; absence of observation or contraction refuses, asserted by CT1
Closure: CT1/all-owners, CT1/all-empty-members, CT1/routes, CT1/legacy-symbols, CT1/zero-consumer-exports

## What to build

Contract legacy aggregate and empty routes through the shared carriers without changing owner semantics or public bytes.

## Acceptance

- [ ] [CT1] (covers AE10) contract legacy aggregate and empty routes preserve all-owners, all-empty-members, routes, legacy-symbols, zero-consumer-exports.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CT1/all-owners | omit one named owner | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| CT1/all-empty-members | omit one empty member | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| CT1/routes | restore one bypass | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| CT1/legacy-symbols | restore one legacy symbol | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| CT1/zero-consumer-exports | leave one zero-consumer export | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |

