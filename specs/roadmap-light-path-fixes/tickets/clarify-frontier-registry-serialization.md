# Clarify frontier registry serialization

Blocked by: none
Writes: specs/roadmap-light-path-fixes/spec.md
Covers: none

## What to build

Make the ticket concurrency rule match the required registry and fixture
closure metadata. Blockers continue to order behavioral dependencies. A pair
of otherwise-frontier tickets with overlapping `Writes:` entries is instead an
explicit coordinator serialization conflict, even when no blocker edge exists.

## Acceptance

- [ ] The spec no longer claims blockers are the only shared-file ordering mechanism.
- [ ] Registry and fixture closure paths remain declared by every bound ticket.
- [ ] The rule matches the repository's ticket discipline without inventing dependency edges.
