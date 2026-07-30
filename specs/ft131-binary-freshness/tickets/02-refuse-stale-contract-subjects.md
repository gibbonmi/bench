# Refuse stale contract subjects

Blocked by: Publish and verify a content-sealed binary

## What to build

The contract launch adapter verifies the repository's selected `dist/bench`
through the current Go sources before AXI wrapper selection or runtime direct
execution, so assertions never credit an unverified subject.

## Acceptance

- [x] Every current AXI suite using `runBenchInDir` rejects a changed source
  before an old executable can emit otherwise-acceptable output.
- [x] Paired runtime fixtures prove stale subjects cannot create either a false
  green or a false red; both report freshness instead of assertion output.
- [x] Missing, malformed, legacy, unreadable, partial, or mismatched artifacts
  produce the same actionable rebuild diagnostic before subject execution.
- [x] Coarse or misleading mtimes do not change the contract verdict.
- [x] Root, nested-cwd, and repeated unchanged launches agree without tracked
  writes or self-invalidation.
