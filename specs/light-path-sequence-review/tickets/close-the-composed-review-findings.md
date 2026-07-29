# Close the composed review findings

Blocked by: none

## What to build

The five-ticket light-path sequence keeps one authoritative race-test registry
and carries durable regression coverage for the edge states found by semantic
review.

## Acceptance

- [x] Synthetic gate fixtures and the race canary generate their prerequisite
  declarations from the authoritative race-test registry.
- [x] `bench gate --help extra` refuses without running the oracle or replacing
  gate evidence.
- [x] A guard worker that continues cleanup after cancellation keeps `Scan`
  from returning until that worker finishes.
- [x] Long-form reversed branch-delete flags preserve the worktree branch
  carve-out at classifier and hook-envelope seams.
- [x] Wrong-type learnings input remains visible as degraded evidence through
  `bench roadmap`.
- [x] The capability fixture proves its nil override scrubs a hostile strict
  value already present in the fixture environment.
