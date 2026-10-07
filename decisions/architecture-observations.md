# Conformance observations and test purpose (C03; FT365)

Status: ready

## Destination

Use the approved landing-test-efficiency spec for shared source observations.
Define the separate residual FT365 scope without duplicating that implementation.

## Notes

The reviewer delegated architecture recommendations on 2026-10-06.
An observation is file data or a parse result for one invocation.
Avoid: policy verdict or reusable green evidence.
Use craft-spec, craft-seams, and craft-gate for the successor spec.

## Decisions so far

- [Which owner and lifetime apply?](architecture-observations/tickets/1.md): use the approved invocation-owned snapshot.
- [What remains after the approved migration?](architecture-observations/tickets/2.md): migrate residual visitors and audit expectations separately.
- [What evidence permits a cheaper test?](architecture-observations/tickets/3.md): retain dispatch and omission detection.

## Not yet specified

## Spec-writer discretion

- Visitor migration order after the approved observation owner exists.
- Private data layout that preserves subject identity and invocation lifetime.

## Out of scope

- A process-global or cross-run source cache.
- One common exclusion policy for every visitor.
- Deletion of an independent expectation without its required red evidence.
- Gate phase overlap under the closed ADR 0024 decision.

## Sources

- Path: `specs/landing-test-efficiency/spec.md`
  Supports: the approved dispatch and observation contracts.
  Drift: a change to LTE-C3 through LTE-C6.
- Path: `roadmap/FT365.md`
  Supports: the residual scanner and expectation scope.
  Drift: a change to the remaining caller set.
- Path: `docs/adr/0006-independent-test-expectations-preserve-omission-oracles.md`
  Supports: the condition for independent expectations.
  Drift: a new approved decision.
- Path: `internal/conformance/canonical_path_owner_test.go`
  Supports: a residual source visitor.
  Drift: a change to its source acquisition.
- Path: `internal/conformance/cancel_signal_registrations_test.go`
  Supports: another residual source visitor.
  Drift: a change to its source acquisition.

