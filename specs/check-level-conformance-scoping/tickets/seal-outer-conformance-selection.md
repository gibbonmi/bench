# Seal outer conformance selection

Blocked by: none
Ownership fence: `internal/gate/gate.go`, `internal/gate/phases.go`, `internal/gate/runner.go`, `internal/conformance/gate_entry_test.go`, `internal/contract/runtime`
Assumptions: inner canary fixtures still need one authenticated check selector; the outer gate owns selection and accepts no ambient narrowing; claims re-derived from the tree at pickup

## What to build

Outer conformance ignores ambient check selection while authenticated inner-canary
execution retains its one-check control, with unknown or misplaced selectors failing
closed.

## Acceptance

- [ ] [OS1] A valid ambient conformance-check name does not narrow an outer gate's dev check set.
- [ ] [OS2] An inner canary fixture can still execute exactly its registered owning check.
- [ ] [OS3] An unknown, tier-invalid, or outer-supplied selector cannot produce an empty or partial green.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| OS1 | restore unconditional propagation of the ambient selector into the outer phase | runtime gate contract | seed a valid selector in the fixture environment, drive the public gate entry, expect the missing-check-set failure |
| OS2 | scrub the selector from authenticated inner execution too | canary ownership contract | drive one owned inner fixture, expect the timing record to contain more or fewer than its one owner |
| OS3 | turn an invalid selector into an empty selection | conformance driver contract | inject unknown and ship-only names at dev tier, expect the distinct fail-closed diagnostics |
