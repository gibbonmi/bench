# Record the cache contract

Blocked by: none
Writes: docs/adr/0018-bench-owns-one-bounded-go-build-cache.md (new)

## What to build

Record the accepted Bench Go build-cache contract for a teammate who has no
build history. The record states ownership, lock behavior, eviction posture,
and the capacity bound. It records the current decision without file paths or
implementation examples.

## Acceptance

- [ ] The ADR states that Bench owns one shared Go build cache.
- [ ] The ADR states that active runs share the cache and cleaning takes exclusive ownership.
- [ ] The ADR states that no run evicts the cache and that the footprint is bounded at 10 GiB.
