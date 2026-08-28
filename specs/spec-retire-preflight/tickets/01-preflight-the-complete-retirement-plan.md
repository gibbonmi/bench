# Preflight the complete retirement plan

Blocked by: none
Writes: internal/spec/spec.go, internal/spec/spec_test.go

## What to build

Before `bench spec retire` removes its first target, it resolves the complete
retirement plan and verifies that each planned removal can proceed. A refusal
names every blocked repository-relative path and leaves the pickup, tickets,
spec file, and spec folder unchanged.

Tests attach at the existing `internal/spec.Command` seam. They exercise the
CLI result and the checkout state, not a deletion helper.

## Acceptance

- [ ] A blocked target makes retirement exit 1 before any planned target is removed.
- [ ] The refusal names each blocked repository-relative path in deterministic order.
- [ ] A writable complete plan still removes the pickup and complete spec folder.
