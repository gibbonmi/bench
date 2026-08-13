# Remove production canary dispatch

Blocked by: report-truthful-canary-inventory.md, stop-seeding-linked-proof.md
Writes: `internal/canary`, `internal/preprelease`

## What to build

Collapse the production canary package onto its one selection result and remove `Dispatch`, `DispatchResult`, the dispatched-owner field, and every function-typed production parameter. Move the ship consumer onto inventory validation in the same landing: separating it strands a compile red because `internal/preprelease` calls the removed `SweepShip` API.

## Acceptance

- [ ] (covers CI2) Production `internal/canary` has no function-typed parameter, dispatch result, dispatched field, or second inventory derivation, and every production caller compiles against the one inventory decision.
