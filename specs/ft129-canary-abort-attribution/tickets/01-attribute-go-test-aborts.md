# Attribute Go test aborts

Blocked by: none

## What to build

Canary sweeps classify authentic Go panic and runtime-fatal output before bite
classification, attribute the deepest abort-associated test or subtest when the
runner grammar provides one, and emit a bounded control-safe diagnostic when it
does not.

## Acceptance

- [x] A root-test panic is reported as an inner test abort naming the failing
  test.
- [x] A subtest panic reports the deepest failing `Test.../...` name associated
  with the abort.
- [x] Optional gate-phase prefixes are accepted, but incidental text containing
  `panic:` is not treated as an abort.
- [x] Truncated panic and runtime-fatal output remain attributed aborts without
  copying raw or unbounded child output.
- [x] Competing failure headers do not displace the header associated with the
  abort.
- [x] Panic classification wins when partial output already contains the
  fixture's `EXPECT`.
