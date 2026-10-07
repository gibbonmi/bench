# Preserve gate signals through the lifetime owner

Blocked by: 01-own-lifetime-through-bounds.md, 02-protect-nested-bounds-users.md
Writes: internal/conformance/injected_ports_registry_test.go, internal/gate/engine.go, internal/gate/lifetime_resources_test.go (new), internal/gate/lifetime_test.go (new), internal/gate/phases.go, internal/gate/run_transaction.go, internal/gate/runner.go
Covers: PL29, PL30, PL31

## What to build

A real gate child keeps INT/130 interruption, TERM/124 timeout, and normal group drain.
Replace its private start, wait, escalation, and drain loops with the established Lifetime.
Use the current two-second grace and accepted final window.

Add an outcome-returning gate execution form and pass its full result to required cleanup and certification callers.
Keep compatibility entry points only where existing callers need scalar presentation.
No required owner may consume that scalar instead of Outcome.
The internal form retains child status, cancellation attribution, stream facts, and unresolved obligations.
Later gate resource tickets consume this form without creating another process owner.

Apply the cache descriptor from ticket 2 and the current owner-composed inherited vector.
Leave later resource kinds unpublished until their owner migration exists.
Retain existing phase order and evidence authority.
Normal gate completion still drains its own group before success.

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

## Acceptance

- [ ] PL29: Gate interruption retains INT and exit 130.
- [ ] PL30: Gate timeout retains TERM and exit 124.
- [ ] PL31: Gate normal completion drains the remaining group.

## Verification

Run the named checks through their production entries.
For each owned row, record a behavioral omission or swap that fails, restoration, and the passing result.
A compilation failure is not a behavioral red.
Preserve existing assertions and fixture mutation purpose.
Co-owned closure holders receive only necessary reference or fixture updates.

- `bench test --package ./internal/gate`

Mutation witness: Replace INT with TERM, map timeout to interruption, or omit normal group drain. The actual gate signal and descendant-survival fixtures must fail.
