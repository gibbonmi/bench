# Converge symlinked hook lifecycle

Blocked by: expose-hook-health-record.md
Ownership fence: `internal/adopt/link_transaction.go`, `internal/adopt/link_stage.go`, `internal/contract/surface/link_test.go`, `internal/contract/surface/link_lifecycle_test.go`
Integration surfaces: hook-health classifier and prospective bytes→expose-hook-health-record.md; manifest match and symlink-parent ordering→`internal/adopt/link_transaction.go` + L1/L2/L3/L5; hook staging→`internal/adopt/link_stage.go` + L4; lifecycle contracts→`internal/contract/surface/link_test.go` and `internal/contract/surface/link_lifecycle_test.go` + L1/L2/L3/L4/L5
Contracts: manifest ownership and prospective hook bytes cross `internal/adopt/link_transaction.go`→`internal/adopt/link_stage.go`, asserted by L4 against the real hook-health producer
Closure: L1/converged-symlink-skip, L2/drifted-conflict, L3/new-conflict, L4/upgrade-hook-refresh, L5/link-idempotency

## What to build

Let link and upgrade traverse a symlinked managed directory only after each prospective entry is proven manifest-clean, then refresh the hook and retain idempotency. Keep link and upgrade in this ticket: a link-only cut strands the lifecycle contract’s converged `bench upgrade` fixture red because upgrade reaches the same parent-ordering refusal before it can refresh the hook.

## Acceptance

- [ ] [L1] `bench link` completes when a symlinked managed directory contains only entries already matching the manifest.
- [ ] [L2] A drifted entry under a symlink parent still aborts.
- [ ] [L3] A new entry under a symlink parent still aborts.
- [ ] [L4] `bench upgrade` completes on the converged symlink fixture and refreshes the hook.
- [ ] [L5] Re-running link over the converged symlink fixture is idempotent in tree and manifest.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| L1/converged-symlink-skip | check the symlink parent before manifest equality | lifecycle surface contract | link the converged symlink fixture, expect success |
| L2/drifted-conflict | skip every existing destination | lifecycle surface contract | drift a manifest-owned entry under the symlink, run link, expect conflict |
| L3/new-conflict | skip every path beneath a known symlink | lifecycle surface contract | add a new planned entry beneath the symlink, run link, expect conflict |
| L4/upgrade-hook-refresh | fix link but leave upgrade on the old ordering | lifecycle surface contract | upgrade the converged fixture, expect success and current hook bytes |
| L5/link-idempotency | skip manifest bookkeeping with the write | lifecycle surface contract | run link twice, compare tree and manifest, expect equality |
