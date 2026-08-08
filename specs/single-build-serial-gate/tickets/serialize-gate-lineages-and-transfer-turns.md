# Serialize gate lineages and transfer turns

Blocked by: share-artifacts-across-local-processes.md, enforce-the-executable-artifact-contract.md, admit-one-gate-per-common-repository.md
Ownership fence: `internal/artifactstore/`, `internal/gate/`
Integration surfaces: repository-common admission→admit-one-gate-per-common-repository.md; artifact coordinator/registry→enforce-the-executable-artifact-contract.md; transferable turn descriptor→reclaim-interrupted-resource-turns.md; transferable turn descriptor→pin-go-children-and-settle-core-packages.md; transferable turn descriptor→serialize-canary-stages.md; transferable turn descriptor→close-all-resource-concurrency-routes.md; turn settlement/watchdog junction→close-all-resource-concurrency-routes.md
Contracts: `TurnDescriptor` (authority schema, lineage ID, generation, authenticated descriptor, owner state) crosses gate parent→synchronous child in `internal/gate/`, membership is one lineage, ordering is parent-transfer/quiesce then child-work/reap then return, absence, forgery, duplication, closed state, malformed state, or foreign lineage refuses, asserted by TL1-TL2 against the real scheduler and process harness
Closure: TL1/one-sibling-lineage, TL1/split-table-order, TL1/quiescent-parent, TL2/exact-transfer, TL2/authentication-inputs, TL2/missing-transfer-refusal, TL2/forged-transfer-refusal, TL2/duplicated-transfer-refusal, TL2/closed-transfer-refusal, TL2/malformed-schema-refusal, TL2/malformed-field-set-refusal, TL2/foreign-lineage-refusal

## What to build

Replace concurrent outer/split scheduling with one sibling phase lineage and add exact parent-to-child turn transfer. Keep scheduling and transfer together because a capacity-one scheduler without transfer leaves the nested canary parent holding the turn while its child waits, making the existing nested gate timeout red.

## Acceptance

- [ ] [TL1] (covers ZC1) primary and stripped tables launch one sibling lineage in stable dependency order while a live waiting parent is quiescent and one descendant may own the turn.
- [ ] [TL2] (covers ZC4) synchronous parent/child transfer has one owner and refuses missing, forged, duplicated, closed, malformed, or out-of-lineage state without local fallback.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| TL1/one-sibling-lineage | restore launch of every ready sibling | blocking scheduler active-recorder test | release two independent phases and expect maximum sibling lineages of one |
| TL1/split-table-order | run primary and stripped schedules in separate goroutines | split-schedule marker test | block both first phases and expect no simultaneous markers plus stable merged order |
| TL1/quiescent-parent | let the outer canary parent perform registered work while its inner child owns the turn | parent/child active-recorder test | transfer to the child, trigger parent work, and expect the double-actor failure |
| TL2/exact-transfer | make the child acquire a fresh turn instead of inheriting | re-exec transfer test | spawn the nested child and expect one lineage/generation with no second acquisition |
| TL2/authentication-inputs | remove the generation from the authenticated descriptor preimage | cross-generation replay test | mint a descriptor, advance its generation while retaining the old authentication, and expect refusal before the resource marker |
| TL2/missing-transfer-refusal | accept a missing descriptor | hostile descriptor table | invoke without a descriptor and expect refusal before its resource marker |
| TL2/forged-transfer-refusal | accept a descriptor with invalid authentication | hostile descriptor table | invoke the forged descriptor and expect refusal before its resource marker |
| TL2/duplicated-transfer-refusal | accept a descriptor already consumed once | hostile descriptor table | replay the descriptor and expect refusal before its resource marker |
| TL2/closed-transfer-refusal | accept a descriptor whose turn is closed | hostile descriptor table | invoke the closed descriptor and expect refusal before its resource marker |
| TL2/malformed-schema-refusal | accept a descriptor with an unknown schema | hostile descriptor table | invoke the unknown schema and expect refusal before its resource marker |
| TL2/malformed-field-set-refusal | accept a descriptor with a missing or extra field | hostile descriptor table | invoke the malformed field set and expect refusal before its resource marker |
| TL2/foreign-lineage-refusal | accept a valid descriptor from another lineage | hostile descriptor table | invoke the foreign descriptor and expect refusal before its resource marker |
