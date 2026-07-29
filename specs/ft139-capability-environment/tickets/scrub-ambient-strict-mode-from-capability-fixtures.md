# Scrub ambient strict mode from capability fixtures

Blocked by: none

## What to build

The runtime capability-skip contract fixture gives each subtest exactly the
strict-mode environment it declares, independent of the parent test process.

## Acceptance

- [x] The informational capability-skip subtest stays green when the parent
  process exports `BENCH_REQUIRE_CAPABILITIES=1`.
- [x] Subtests that explicitly enable strict mode retain their existing red or
  green verdicts for capability and environment skips.
- [x] The fixture uses the existing environment override seam rather than
  introducing a second command runner or environment merger.
