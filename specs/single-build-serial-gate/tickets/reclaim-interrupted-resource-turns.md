# Reclaim interrupted resource turns

Blocked by: serialize-gate-lineages-and-transfer-turns.md
Ownership fence: `internal/gate/`
Integration surfaces: transferable turn lifecycle→serialize-gate-lineages-and-transfer-turns.md; interruption settlement→pin-go-children-and-settle-core-packages.md; interruption settlement→serialize-canary-stages.md; exact-once reclamation and verdict junction→close-all-resource-concurrency-routes.md
Contracts: interruption events cross current owner, delegated child, or queued waiter→turn coordinator, membership is every live turn state, ordering is cancel descendants-reap child-return or remove waiter-publish terminal settlement, absence cannot release a sibling or publish ready, asserted by TL3 against real cancellation and accounting paths
Closure: TL3/owner-interrupt, TL3/child-interrupt, TL3/waiter-interrupt, TL3/exact-once-return, TL3/no-ready-after-owner-interrupt, TL3/no-ready-after-child-interrupt

## What to build

Add exact-once interruption teardown to the established turn lifecycle. Owners, delegated children, and waiters must leave no runnable descendant, leaked queue entry, duplicate return, or false ready verdict.

## Acceptance

- [ ] [TL3] (covers ZC6) owner, delegated-child, and waiter interruption tears down descendants, reclaims or returns the turn exactly once, permits the next acquisition, and cannot publish a false ready verdict.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| TL3/owner-interrupt | omit owner teardown/reclamation | owner cancellation test | interrupt the owner, bound the wait, and expect descendants gone and next acquisition live |
| TL3/child-interrupt | return the turn before the delegated child is reaped | child cancellation test | interrupt the child and expect no sibling start until reap plus one return |
| TL3/waiter-interrupt | leave a cancelled waiter queued | waiter cancellation test | cancel the waiter and expect it absent before the next transfer |
| TL3/exact-once-return | return on both signal and wait paths | turn accounting test | force the race and expect the token and owner counts to remain exactly one |
| TL3/no-ready-after-owner-interrupt | publish ready while the owner is interrupted | interrupted-lineage verdict test | snapshot the gate record, interrupt the owner, and expect no new ready verdict before the next run acquires the turn |
| TL3/no-ready-after-child-interrupt | publish ready while the delegated child is interrupted | interrupted-lineage verdict test | snapshot the gate record, interrupt the child, and expect no new ready verdict before the next run acquires the turn |
