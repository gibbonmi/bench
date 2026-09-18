# Derive TagOf from the row-ID grammar

Blocked by: none
Writes: internal/tickets/tickets.go, internal/tickets/enumerate.go
Covers: none

## What to build

`tickets.RowIDPattern` is the one row-ID grammar, and its group captures the tag.
`tickets.TagOf` scans the same uppercase tag with a hand-written byte loop. The two
sources of the tag grammar can drift. This ticket moves the tag grammar into one
constant. `RowIDPattern` composes that constant, and `TagOf` matches the same
constant as a prefix. The change keeps every `TagOf` result, and it keeps the
`RowIDPattern` value.

## Acceptance

- [ ] One constant holds the tag grammar, and both `RowIDPattern` and `TagOf` read it.
- [ ] `TagOf` has no hand-written character scan.
- [ ] The `internal/tickets`, `internal/preflight`, and `internal/coverage` tests pass with no change to a test expectation.
