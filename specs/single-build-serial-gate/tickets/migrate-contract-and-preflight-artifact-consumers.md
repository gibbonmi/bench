# Migrate contract and preflight artifact consumers

Blocked by: share-artifacts-across-local-processes.md
Ownership fence: `internal/contract/`, `internal/preflight/`
Integration surfaces: artifact producer→share-artifacts-across-local-processes.md; ordinary contract/runtime/preflight consumer inventory→internal/contract and internal/preflight; migrated adapters→migrate-contract-preflight-release-proof-consumers.md; legacy helper retirement→enforce-the-executable-artifact-contract.md
Contracts: `ArtifactRecord` crosses `internal/artifactstore/`→adapters in `internal/contract/` and `internal/preflight/` as the unchanged host CLI class, membership is every registry consumer not observing construction, ordering is resolve before fixture materialization or command execution, absence/staleness refuses without local build, asserted by CP1-CP2 against the real producer
Closure: CP1/contract-runtime-reuse, CP1/contract-fixture-reuse, CP1/stale-refusal, CP2/preflight-current-subject-reuse, CP2/no-private-builder

## What to build

Migrate ordinary contract/runtime and preflight helpers that only need unchanged current-subject bytes. These local rows are interim migration landings whose universal enrollment is graded finally by SB4 in the contraction ticket. Preserve their fixture destinations and freshness assertions while removing their repeated current-root builds.

## Acceptance

- [ ] [CP1] (covers local) every ordinary contract/runtime consumer uses the selected current-subject artifact and retains stale/missing refusal.
- [ ] [CP2] (covers local) every ordinary preflight current-subject consumer materializes the selected artifact and invokes no private builder.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CP1/contract-runtime-reuse | restore a runtime helper's `scripts/go-build.sh` call | contract consumer enumeration/counting backend test | invoke the full registered ordinary contract family and expect zero direct builds |
| CP1/contract-fixture-reuse | point one cache/install fixture at a separately built output | fixture artifact-identity test | realize the destination and expect its executable digest to equal the selected record |
| CP1/stale-refusal | skip record/freshness verification in the thin contract adapter | stale-subject contract test | mutate a build input and expect refusal before the runtime assertion executes |
| CP2/preflight-current-subject-reuse | retain one per-test current-root build | preflight consumer enumeration/counting backend test | invoke all registered ordinary preflight helpers and expect one selected identity with no extra build |
| CP2/no-private-builder | create an unregistered helper around `scripts/go-build.sh` | structural executable-consumer audit | add/retain the adapter call and expect its consumer-without-registry diagnostic |
