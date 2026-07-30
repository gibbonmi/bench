# Close semantic review findings

Blocked by: Attribute Go test aborts, Attribute process aborts

## What to build

The abort classifier uses one fixture-owner classification for both test and
process diagnostics, verifies runtime-fatal grammar against authentic Go test
output, and covers process-abort attribution for contract, conformance-scope,
and unscoped fixtures.

## Acceptance

- [x] Test-abort and process-abort diagnostics derive package, conformance
  scope, and inner-gate ownership through one implementation source.
- [x] A helper subprocess produces authentic Go runtime-fatal output that
  `Sweep` classifies as an attributed inner test abort.
- [x] Process aborts from conformance-scoped and unscoped fixtures report their
  exact owners without falling back to contract-package wording.
