# Run binary provenance in a `bench test` result (FT290)

Status: shaping

## Destination

A `bench test` result shows which run binary and which source graded the run.
A caller can then see when a named check graded test code that is not the
source in the calling checkout.

## Notes

Domain: the run binary selection and the `bench test` result projection.
Charge `bench-craft-research` for ticket #1.

This map is the second half of the FT290 split. The first half is delivered.
It added the `check` row of a named-check result, and a provenance fact can
join that row.

A map-owned asset stays in the map's assets folder,
decisions/run-binary-provenance/assets/.

## Decisions so far

## Not yet specified

## Spec-writer discretion

## Out of scope

- The identity in the unknown-check refusal; FT290 delivered it.

## Sources

- Path: `internal/runbinary/runbinary.go`
  Supports: the entry state of ticket #1. `Own` builds from the source root. `Inherit` refuses a seal that does not agree with the source digest.
  Drift: a change to `ReuseOrOwn`, `Own`, `Inherit`, or `canonicalVerify`.
