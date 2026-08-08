# Rendezvous artifact consumers across processes

Blocked by: share-artifacts-across-local-processes.md
Ownership fence: `internal/artifactstore/`
Integration surfaces: durable store API→share-artifacts-across-local-processes.md; target-aware resolver→resolve-target-aware-artifact-identities.md; focused and multi-package process harnesses→owned paths
Contracts: one repository-common store root crosses focused re-exec and sibling-package consumers→`internal/artifactstore/`, membership is every consumer process for one common Git directory, ordering is resolve through the durable store before any backend call, absence of shared rendezvous is a duplicate-build red, asserted by AS4 against raw multiprocess package execution
Closure: AS4/focused-reexec, AS4/multipackage-rendezvous

## What to build

Prove the expanded durable store is the rendezvous used by focused re-exec children and independent package processes, without a process-local or package-private fallback.

## Acceptance

- [ ] [AS4] (covers SB6) focused re-exec children and sibling consumer packages from raw `go test ./...` rendezvous through the common repository store.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AS4/focused-reexec | create a process-local cache in the re-exec child | focused re-exec counter test | resolve in parent and child, then expect one shared identity and one backend call |
| AS4/multipackage-rendezvous | derive a package-private store root | multi-package subprocess fixture | run the two consumer packages together and expect the duplicate-build failure |
