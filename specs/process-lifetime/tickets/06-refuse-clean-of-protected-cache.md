# Refuse cache clean after unresolved CLI exit

Blocked by: 01-own-lifetime-through-bounds.md, 02-protect-nested-bounds-users.md, 03-migrate-gate-process-policy.md, 04-retain-selected-build-binaries.md, 05-protect-private-test-directories.md
Writes: internal/conformance/injected_ports_registry_test.go, internal/gate/lane.go, internal/gate/lifetime_resources_test.go (new), internal/gate/phases.go, internal/gate/run_transaction.go, internal/gate/runner.go, internal/gocache/clean.go, internal/gocache/lock.go, internal/gocache/protection.go (new), internal/gocache/protection_test.go (new), internal/systemtest/process_lifetime_cache_test.go (new)
Covers: PL55, PL56, PL57, PL76, PL77, PL78

## What to build

A second real cache-clean CLI refuses after the initiating CLI exits and drops its process-owned shared lock.
Complete cache-holder integration for gate and lane callers through the existing protection descriptor.
Hold shared exclusion before registration and clear ordinary complete protection before releasing the lock.
Unresolved records survive lock release.

The exclusive cleaner seals and observes the complete durable user set before go clean starts.
Keep absent-cache zero behavior without creating state.
Keep an empty cache cleanable when complete protection permits release.
Unreadable state remains uncertain.

Consume the typed gate outcome from ticket 3 and close results from existing resource owners.
A failed gate terminal protection publication permits no retained green.
Retain pending evidence and both process and persistence causes.
Use the shared protection predicate rather than a second cache-specific completion rule.

Keep over-budget source files from gaining lines at this checkpoint.
Move only cohesive outcome or phase composition into the co-owned runner or phases source when needed.
Do not defer headroom debt to a successor.

Before implementation starts, confirm complete accepted and landed durable-file-replacement delivery.
Its leaf, review-record caller, and native qualification must be complete.
Durable-caller-migration is not a prerequisite.
Use one retained integration source and a fresh author for this ticket.

Keep each caller's selected normal-exit policy.
A raw child status cannot certify required cleanup or publication.
Retain uncertain resources and every unresolved obligation.
Never signal from recovered records or add a blanket normal-exit kill.

This ticket belongs to PL-C3.
Finish the preceding chunk review before this chunk starts.
Its blockers include every prior overlapping Writes owner.
No successor caller must exist for this ticket's owned predicates.

BENCH_KIT is supplied by bench test --check system through the existing sealed system owner.
Do not introduce another subprocess fixture runner or private binary publisher.

## Acceptance

- [ ] PL55: A second-process cache clean refuses unresolved protection.
- [ ] PL56: The exclusive cleaner retains unreadable protection.
- [ ] PL57: Completed cache protection clears before shared-lock release.
- [ ] PL76: An absent cache clean creates no directory.
- [ ] PL77: An empty cache with complete protection remains cleanable.
- [ ] PL78: Gate terminal persistence failure permits no retained green.

## Verification

Run the named checks through their production entries.
For each owned row, record a behavioral omission or swap that fails, restoration, and the passing result.
A compilation failure is not a behavioral red.
Preserve existing assertions and fixture mutation purpose.
Co-owned closure holders receive only necessary reference or fixture updates.

- `bench test --package ./internal/gocache`
- `bench test --package ./internal/gate`
- `bench test --check system`

Mutation witness: Use only the dropped owner lock, treat unreadable records as empty, or publish green before terminal persistence. The second CLI cleaner and real gate record witnesses must fail.
