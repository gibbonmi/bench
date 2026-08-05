# Route optimistic reuse through the generation owner

Blocked by: none
Ownership fence: `internal/gate/`
Contracts: none crosses — the reuse subject moves from the root-based capture to the evaluation-owned generation entirely inside `internal/gate/`, asserted by OR1-OR2 through the public ExecuteReusingFreshGreen journey

## What to build

Close review finding S1-optimistic-reuse-bypasses-generation-owner: `ExecuteReusingFreshGreen` builds its reuse subject through the root-based `buildSubject` capture instead of the evaluation-owned generation lifecycle, so the optimistic exact-green answer derives identity outside the one snapshot contract. Route the pre-lock reuse decision through a working-tree evaluation's accepted pre generation; a miss still falls through to the real execution unchanged.

## Acceptance

- [ ] [OR1] A reuse hit derives its subject from an evaluation-owned pre generation: it announces the reused verdict, performs exactly one parsed listing, at most one working-tree materialization, and creates no post generation.
- [ ] [OR2] A subject with no reusable verdict still falls through to a real execution that runs the gate and records its verdict.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| OR1 | restore the root-based buildSubject capture on the reuse path | the evaluation-owned reuse journey test | apply the mutation, run `go test ./internal/gate -run ExecuteReusingFreshGreen`, expect the listing-count failure |
| OR2 | return the optimistic result unconditionally instead of falling through | the reuse fall-through journey test | apply the mutation, run `go test ./internal/gate -run ExecuteReusingFreshGreen`, expect the missing-real-run failure |
