# Generic Git fixtures (FT360; C04)

Status: ready

## Destination

Refresh the staged shared-test-fixtures spec around the current helper contract.
Keep generic fixture policy separate from specialized process probes.

## Notes

The reviewer delegated architecture recommendations on 2026-10-06.
The staged spec still requires implementation approval.
A generic fixture prepares repository state without grading its process transport.
Avoid: production Git adapter.
Use craft-spec and craft-seams for the amendment.

## Decisions so far

- [How should the existing Run name be handled?](architecture-fixtures/tickets/1.md): preserve error-returning Run and amend the fatal helper name.
- [Which fixture knowledge moves?](architecture-fixtures/tickets/2.md): move generic execution, identity, and snapshots only.
- [Which tests and ordering constrain the amendment?](architecture-fixtures/tickets/3.md): preserve fixture semantics and coordinate LTE edits.

## Not yet specified

## Spec-writer discretion

- The private execution helper shape after the public helper contracts remain explicit.
- Migration batches that preserve operation-specific exceptions.

## Out of scope

- Replacement of real Git journeys with mocked repository state.
- Implicit staging in the index-only commit helper.
- Removal of specialized process supervision.
- The minimal adoption payload work already owned by landing-test-efficiency.

## Sources

- Path: `roadmap/FT360.md`
  Supports: the staged owner and acknowledged Run collision.
  Drift: a change to the row or helper contract.
- Path: `specs/shared-test-fixtures/spec.md`
  Supports: the existing migration and capability requirements.
  Drift: an approved spec amendment.
- Path: `internal/gittest/gittest.go`
  Supports: the current Output and error-returning Run forms.
  Drift: a change to those exported forms.
- Path: `internal/testrepo/working_tree.go`
  Supports: the second snapshot and identity implementation.
  Drift: a change to snapshot behavior.

