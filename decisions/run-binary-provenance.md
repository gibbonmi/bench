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
`specs/ft290-test-projection/decisions/ft290-test-projection.md`, and it owns the `check` row that a
provenance fact can join.

A map-owned asset stays in the map's assets folder,
decisions/run-binary-provenance/assets/.

## Decisions so far

## Not yet specified

## Spec-writer discretion

## Out of scope

- The identity in the unknown-check refusal; ticket #2 of `specs/ft290-test-projection/decisions/ft290-test-projection.md` owns it.

## Sources

- Path: `specs/ft290-test-projection/spec.md`
  Supports: two lines name this map as the owner of the run binary provenance. They are the line "Not covered: story 36" and the Out of scope line "The run binary provenance in a result".
  Drift: a change to either line.
- Path: `internal/runbinary/runbinary.go`
  Supports: the entry state of ticket #1. `Own` builds from the source root. `Inherit` refuses a seal that does not agree with the source digest.
  Drift: a change to `ReuseOrOwn`, `Own`, `Inherit`, or `canonicalVerify`.
