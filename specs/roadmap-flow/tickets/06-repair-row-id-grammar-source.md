# Repair: read the row-ID grammar from package roadmap

Blocked by: 01-report-the-board-flow.md
Writes: internal/roadmapflow/flow.go

## What to build

The flow derivation decides whether a history path names a detail owner
through `roadmap.ValidOccurrenceOwner`, not through a second regular
expression of the `FT<n>` grammar. Review finding (Standards F1, judgment
call): one source per fact.

## Acceptance

- [ ] `internal/roadmapflow` contains no literal of the `FT[1-9][0-9]*` grammar.
- [ ] `go test ./internal/roadmapflow` stays green; RF1 and RF5 still red a path outside the grammar.
