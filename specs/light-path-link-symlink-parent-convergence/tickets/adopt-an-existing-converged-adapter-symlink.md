# Adopt an existing converged adapter symlink

Blocked by: none
Ownership fence: `internal/adopt/link_transaction.go`, `internal/adopt/link_transaction_test.go`, `CHANGELOG.md`, `ROADMAP.md`, `specs/light-path-link-symlink-parent-convergence/tickets/adopt-an-existing-converged-adapter-symlink.md`
Integration surfaces: staged adapter classification and transaction refusal→`internal/adopt/link_transaction.go`; first-adoption regression→`internal/adopt/link_transaction_test.go`; public `bench link` repro→unchanged command surface + LC1
Contracts: none crosses
Closure: LC1/unowned-converged-parent-symlink, LC2/unowned-divergent-parent-symlink, LC3/foreign-identical-mirror

## What to build

`bench link` adopts an existing adapter destination reached through a directory symlink when the staged adapter and destination resolve to the same regular file, even before a manifest row exists. A genuinely divergent or byte-identical foreign destination beneath that symlink still triggers the transaction-wide hard refusal before any planned write lands.

## Acceptance

- [ ] [LC1] With no prior manifest ownership, a staged adapter whose destination resolves through a parent-directory symlink to the canonical target's same regular file is classified as converged and the transaction succeeds.
- [ ] [LC2] With no prior manifest ownership, a staged adapter whose destination resolves through a parent-directory symlink to different regular-file content is refused with the symlink-parent conflict and no accepted write is promoted.
- [ ] [LC3] A byte-identical foreign mirror beneath a parent-directory symlink is refused with the symlink-parent conflict and no accepted write is promoted.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| LC1/unowned-converged-parent-symlink | require a nonempty old-manifest hash before checking convergence | the first-adoption transaction test | construct a root with the adapter-directory symlink and no manifest, run the focused `internal/adopt` test, expect exit 0 to become the symlink-parent refusal |
| LC2/unowned-divergent-parent-symlink | skip the symlink-parent refusal after convergence fails | the divergent negative-control leg | point the adapter directory at different regular-file bytes, include one otherwise accepted plan entry, run the focused test, expect exit 1 and assert the accepted write remains absent |
| LC3/foreign-identical-mirror | admit byte equality without same-file identity | the foreign-identical negative-control leg | point the adapter directory at a distinct file with identical bytes, include one otherwise accepted plan entry, run the focused test, expect exit 1 and assert the accepted write remains absent |
