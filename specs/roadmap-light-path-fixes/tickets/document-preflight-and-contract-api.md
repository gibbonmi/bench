# Document preflight and marker-wait contracts

Blocked by: none
Writes: internal/preflight/decision.go, internal/preflight/gather.go, internal/contract/marker_wait.go
Covers: LF17

## What to build

Document Decide, Gather, GatherPinned, AuthorizeReviewedSource, and
WaitForTwoLegMarkers. State accepted inputs, snapshot guarantees, terminal
results, and interruption behavior.

## Acceptance

- [ ] All five named functions state their caller-facing contracts.
- [ ] Preflight comments distinguish gathering from pure decision.
- [ ] Marker waiting documents deadlines, exit signaling, and terminal misses.
