# Run binary provenance in a `bench test` result (FT290)

Status: shaping

## Destination

A `bench test` result shows which run binary and which source graded the run.
A caller can then see when a named check graded test code that is not the
source in the calling checkout.

## Notes

Domain: the run binary selection and the `bench test` result projection.
Charge `bench-craft-research` for ticket #1.

This map is the second half of the FT290 split. The first half is
`decisions/ft290-test-projection.md`, and it owns the `check` row that a
provenance fact can join.

A map-owned asset stays in the map's assets folder,
decisions/run-binary-provenance/assets/.

## Decisions so far

## Not yet specified

## Spec-writer discretion

## Out of scope

- The identity in the unknown-check refusal; `decisions/ft290-test-projection.md` ticket #2 owns it.

## Sources

- Path: `roadmap/FT290.md`
  Supports: the 2026-09-17 occurrence that ticket #1 must reproduce.
  Drift: a new occurrence or a body edit on the row.
- Path: `internal/runbinary/runbinary.go`
  Supports: the entry state of ticket #1. `Own` builds from the source root. `Inherit` refuses a seal that does not agree with the source digest.
  Drift: a change to `ReuseOrOwn`, `Own`, `Inherit`, or `canonicalVerify`.
