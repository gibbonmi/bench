# Single-source the canary check marker

Blocked by: none
Writes: internal/canary/inventory.go, internal/canary/shape_test.go
Covers: none

## What to build

The canary inventory reader uses the declared check-marker filename instead of
a second literal spelling. The production shape check refuses a marker reader
that embeds a filename literal.

## Acceptance

- [ ] The canary check-marker reader gets its filename from `checkFileName`.
- [ ] The canary shape test fails when a production `readMarker` call embeds a
      filename literal.
