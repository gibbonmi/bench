# Migrate gate fixture artifact consumers

Blocked by: share-artifacts-across-local-processes.md, prepare-gate-artifacts-before-scheduling.md
Ownership fence: `internal/gate/`
Integration surfaces: selected artifact producer→prepare-gate-artifacts-before-scheduling.md; kit-shaped/current-subject consumer inventory→internal/gate; migrated fixture seams→migrate-gate-real-build-proof-consumers.md; legacy helper retirement→enforce-the-executable-artifact-contract.md
Contracts: `ArtifactRecord` crosses gate preparation→`internal/gate` fixture/current-subject adapters, membership is every non-compiler-observing gate consumer in the executable registry, consumers preserve registry order and request before materialization, absence refuses rather than building locally, asserted by GF1-GF2 against the real store artifact
Closure: GF1/kit-shaped-reuse, GF1/current-subject-reuse, GF1/root-local-seals, GF2/missing-record-refusal, GF2/no-process-template

## What to build

Move every ordinary `internal/gate` fixture and current-subject helper onto the selected artifact while retaining private roots, seals, Git state, and mutable-publication behavior. These local rows are interim migration landings whose universal enrollment is graded finally by SB4 in the contraction ticket. Leave compiler-observing cases for the registered-proof migration.

## Acceptance

- [ ] [GF1] (covers local) deterministic kit-shaped and current-subject consumers reuse the real store artifact while each root retains its own seal, Git directory, and evidence.
- [ ] [GF2] (covers local) a gate fixture missing its selected artifact refuses and no process-local `sync.Once` or direct canonical builder remains as fallback.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| GF1/kit-shaped-reuse | restore the process-scoped template builder | kit-shaped backend-count test | construct fixtures across re-exec processes and expect one artifact identity/build |
| GF1/current-subject-reuse | restore `currentBenchBinary` construction | current-subject helper enumeration test | invoke every registered current-subject helper and expect zero direct backend calls |
| GF1/root-local-seals | publish the store path directly into a fixture | existing freshness/fixture isolation test | construct two roots, mutate one local publication, and expect the other/store digest to remain fixed |
| GF2/missing-record-refusal | fall back to `buildFixtureBinaryTo` on absence | missing-artifact fixture test | remove the selected record, construct the fixture, and expect refusal plus no builder marker |
| GF2/no-process-template | retain `kitShapedTemplateState.once` as an alternate route | structural gate consumer audit | enumerate template/build helper symbols and expect no unregistered ordinary constructor |
