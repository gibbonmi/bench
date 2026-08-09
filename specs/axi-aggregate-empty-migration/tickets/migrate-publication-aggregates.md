# Migrate publication aggregates

Blocked by: none
Ownership fence: `internal/publication`
Integration surfaces: shared aggregate/empty carriers→implemented prerequisite `axi-carriers-and-registry`; final contraction→contract-aggregate-empty-routes.md
Contracts: publication record typed facts cross `internal/publication` state machine→renderer; domain is durable release record; order is record order; absence is explicit no-record with zero/unknown preserved, asserted by PA1
Closure: PA1/order, PA1/typed-facts, PA1/durable-next, PA1/zero, PA1/unknown, PA1/route

## What to build

Migrate publication aggregates through the shared carriers without changing owner semantics or public bytes.

## Acceptance

- [ ] [PA1] (covers AE6) migrate publication aggregates preserve order, typed-facts, durable-next, zero, unknown, route.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PA1/order | derive next state outside record | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| PA1/typed-facts | stringify numeric fact | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| PA1/durable-next | omit zero | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| PA1/zero | coerce unknown | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| PA1/unknown | reorder fields | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| PA1/route | bypass shared aggregate | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |

