# Migrate status and dashboard aggregates and empties

Blocked by: none
Ownership fence: `internal/status`, `internal/dashboard`
Integration surfaces: shared aggregate/empty carriers→implemented prerequisite `axi-carriers-and-registry`; final contraction→contract-aggregate-empty-routes.md
Contracts: status typed signal facts and empty-class enum cross `internal/status`→`internal/dashboard`; domain is clean/active/degraded; order is producer then section composition; absence distinguishes prose clean and refusal, asserted by SD1 and SD1E
Closure: SD1/signal-order, SD1/zero, SD1/unknown, SD1/composition, SD1/route, SD1/prose-clean, SD1/absent-refusal, SDE1/empty-class

## What to build

Migrate status and dashboard aggregates and empties through the shared carriers without changing owner semantics or public bytes.

## Acceptance

- [ ] [SD1] (covers AE7) migrate status and dashboard aggregates and empties preserve signal-order, zero, unknown, composition, route, prose-clean, absent-refusal.
- [ ] [SDE1] (covers AE8) migrate status and dashboard aggregates and empties preserve their exact empty/absent rendering class.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SD1/signal-order | locally recompute dashboard fact | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| SD1/zero | omit zero | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| SD1/unknown | coerce unknown | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| SD1/composition | reorder signals | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| SD1/route | bypass shared aggregate | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| SD1/prose-clean | normalize prose clean to table | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| SD1/absent-refusal | turn absent refusal into empty success | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| SDE1/empty-class | normalize the named empty class | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
