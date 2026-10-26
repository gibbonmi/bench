# Admit a stack publication in one commitment decision

Blocked by: none
Writes: internal/commitment/repository/publication.go, internal/commitment/repository/stack_publication_test.go (new), internal/commitment/repository/export_test.go
Covers: LS18, LS19, LS20, LS21, LS22, LS23

## What to build

Add two exported operations to the commitment store in `internal/commitment/repository/publication.go`, beside `PublishAdmitted`. The directory already holds 16 source files, so they take no new production file.

`Store.PublishStackAdmitted(members []Publication, trees []string, publish func() error) error` decides a stack under one intent transaction, as `PublishAdmitted` decides one landing:

- It decides the build window for each published member with `store.window`.
- It decides the admission of each published member with `store.admitPublication` over that member's own tree, which is `trees[i]`.
- When every decision admits, it runs `publish` once.
- After `publish` succeeds, it consumes the override of each published member and drops the one lease that the first member's `Publication` holds.
- A refused member stops the decision before `publish`, and the refusal keeps its own face.

`Store.HoldsDeliveryBinding(assignment string) (bool, error)` answers whether the ledger holds a delivery binding for that assignment. A binding of any deliverable kind counts, a started tickets-only deliverable included. It reads the bindings of the ledger and not the open builds, because `openBuilds` counts only a staged-spec deliverable (`internal/commitment/repository/window.go:246`).

Ticket 07 calls `PublishStackAdmitted` from the stack admission of the verb. Ticket 10 calls `HoldsDeliveryBinding` from the stack preflight.

## Acceptance

- [ ] The stack publication refuses, and runs no publish, when the build window refuses one member (LS18).
- [ ] A stack publication consumes the override of each published member (LS19).
- [ ] A stack publication keeps the override of a member that it does not publish (LS20).
- [ ] The stack publication drops the one landing lease in its own decision (LS21).
- [ ] The binding query answers true for an assignment that holds the delivery binding of a started tickets-only deliverable (LS22).
- [ ] The binding query answers false for an unbound light-path assignment during an open build (LS23).

The tickets-only fixture follows the `tickets-only` case of `TestOpenBuildsDeriveFromActiveOwners` (`internal/commitment/repository/window_test.go:96`). The override fixture follows `TestWindowNeedsAGapOfEveryOtherBuild`, which writes an override into the ledger state.

## Checkpoint verification

- `bench test --package ./internal/commitment/repository`
- `bench test --check root-conformance`

## Mutation proof

Make the stack decision check the window of the first member only. The LS18 test must fail. Spell the probe against the code that the author writes, and require `bench probe` to report `bit` with a byte-exact restore.
