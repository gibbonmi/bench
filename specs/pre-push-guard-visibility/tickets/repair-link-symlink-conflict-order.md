# Repair link symlink conflict order

Blocked by: none
Ownership fence: `internal/adopt/link_transaction.go`, `internal/contract/surface/link_lifecycle_test.go`
Integration surfaces: lifecycle gate ordering→`internal/adopt/link_transaction.go` + LSO1; clean-entry skip→`internal/adopt/link_transaction.go` + LSO2; link lifecycle contract→`internal/contract/surface/link_lifecycle_test.go` + LSO1
Contracts: none crosses
Closure: LSO1/symlink-order, LSO2/clean-skip

## What to build

Restore the spec's lifecycle-gate semantics in the link transaction. A *drifted*
entry (exists, manifest-owned, content differs) under a symlink parent directory
must hard-abort the whole transaction with exit 1 and the
`conflict: … has a symlink parent directory` message — never join the soft
`conflicts` list, never let the transaction proceed to promote other entries. A
*clean* entry (exists and matches the manifest) must be skipped entirely: no
staging, no backup rename, no inode replacement at the destination — the spec's
"skipped, no write, no conflict". Keep the new-entry symlink abort green. Update
`link_lifecycle_test.go` so the drifted-under-symlink case asserts exit 1 (the
current test asserts the defective exit 3).

## Acceptance

- [ ] [LSO1] (covers local) `bench link` with a drifted manifest-owned entry under a symlink parent exits 1 with the symlink-parent conflict message and promotes nothing.
- [ ] [LSO2] (covers local) `bench link` leaves a clean entry's inode untouched (no backup, no rename-over), including when its parent is a symlink.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| LSO1/symlink-order | evaluate the existence/clean branch before the symlink-parent abort so the drifted entry lands in the soft conflicts list | link lifecycle contract | seed a drifted managed entry under a symlinked parent, run `bench link`, expect exit-1 assertion red on exit 3 |
| LSO2/clean-skip | let a clean entry fall through to `accepted` and be promoted | link lifecycle contract | record the destination inode, run `bench link` with an all-clean manifest, expect the same-inode assertion red |
