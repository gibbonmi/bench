# Move the prose preparation step to its own file

Blocked by: none
Writes: internal/prose/prepare.go (new), internal/prose/parse.go, internal/prose/starts.go
Covers: none

## What to build

The `internal/prose` package prepares a document before it grades the sentences.
The `prepare` function is in `starts.go`, and its strip helpers are in `parse.go`.
This ticket moves `prepare`, `stripFrontmatter`, `stripComments`, `stripFences`, and
`fenceMarker` into a new `prepare.go` file. The move changes no behavior and no
exported name. After the move, `parse.go` is below its 400-line budget.

## Acceptance

- [ ] `prepare.go` holds `prepare` and the strip helpers, and no other file declares them.
- [ ] `parse.go` has 400 lines or fewer, and `bench structure` reports no violation for it.
- [ ] The package tests pass with no change to a test expectation.
