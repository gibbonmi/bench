# Align the projection ticket fence

Blocked by: 06-pin-the-projection-flag-exclusion.md
Writes: specs/bench-shell-follow-on-guard/tickets/06-pin-the-projection-flag-exclusion.md

## What to build

Make ticket 06 name the authorized top-level command seam that carries its proof.

## Acceptance

- [ ] Ticket 06 names `cmd/bench`, not `internal/gate`, as its write fence.
- [ ] The ticket still describes `bench gate --brief` as an exit-2 usage error.
