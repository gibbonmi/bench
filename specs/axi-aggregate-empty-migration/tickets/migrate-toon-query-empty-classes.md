# Migrate TOON and query empty classes

Blocked by: none
Ownership fence: `internal/toon`, `cmd/bench`
Integration surfaces: shared aggregate/empty carriers→implemented prerequisite `axi-carriers-and-registry`; final contraction→contract-aggregate-empty-routes.md
Contracts: empty-class enum and exact raw output cross `internal/toon` and cmd query producers→existing renderers; domain is zero-row success/prose clean/absent refusal; order is current rendering; absence never defaults across classes, asserted by TE1
Closure: TE1/zero-row-toon, TE1/prose-clean, TE1/absent-refusal, TE1/route

## What to build

Migrate TOON and query empty classes through the shared carriers without changing owner semantics or public bytes.

## Acceptance

- [ ] [TE1] (covers AE8) migrate toon and query empty classes preserve zero-row-toon, prose-clean, absent-refusal, route.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| TE1/zero-row-toon | normalize zero-row TOON to prose | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| TE1/prose-clean | normalize prose clean to zero rows | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| TE1/absent-refusal | convert absent refusal to success | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| TE1/route | bypass typed empty route | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |

