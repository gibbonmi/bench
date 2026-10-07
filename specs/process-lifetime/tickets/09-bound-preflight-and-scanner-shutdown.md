# Bound preflight and scanner teardown without changing findings

Blocked by: 01-own-lifetime-through-bounds.md, 02-protect-nested-bounds-users.md, 03-migrate-gate-process-policy.md, 04-retain-selected-build-binaries.md, 05-protect-private-test-directories.md, 06-refuse-clean-of-protected-cache.md, 07-retain-prospective-gate-artifacts.md, 08-preserve-release-step-evidence.md
Writes: internal/conformance/injected_ports_registry_test.go, internal/releasepreflight/command.go, internal/releasepreflight/lifetime_test.go (new), internal/releasepreflight/vulnerability.go
Covers: PL62, PL63

## What to build

A TERM-resistant preflight or scanner child receives escalation and returns retained failure by the final deadline.
Migrate both existing controlled-group callers through Lifetime.
Use TERM, the accepted missing two-second grace, and the shared final cleanup window.
Preserve normal phase results, scanner validation, exit-3 findings, and input refusal order.

Consume the release-evidence Descriptor and publication disposition from ticket 8.
External preflight refuses promotion while any applicable user remains protected or uncertain.
Keep process, decode, cleanup, and evidence results separate.
Do not reconstruct successful findings from truncated output or child zero.

Before implementation starts, confirm complete accepted and landed durable-file-replacement delivery.
Its leaf, review-record caller, and native qualification must be complete.
Durable-caller-migration is not a prerequisite.
Use one retained integration source and a fresh author for this ticket.

Keep each caller's selected normal-exit policy.
A raw child status cannot certify required cleanup or publication.
Retain uncertain resources and every unresolved obligation.
Never signal from recovered records or add a blanket normal-exit kill.

This ticket belongs to PL-C4.
Finish the preceding chunk review before this chunk starts.
Its blockers include every prior overlapping Writes owner.
No successor caller must exist for this ticket's owned predicates.

## Acceptance

- [ ] PL62: External preflight refuses promotion after incomplete cleanup.
- [ ] PL63: Scanner cancellation escalates a TERM-resistant child.

## Verification

Run the named checks through their production entries.
For each owned row, record a behavioral omission or swap that fails, restoration, and the passing result.
A compilation failure is not a behavioral red.
Preserve existing assertions and fixture mutation purpose.
Co-owned closure holders receive only necessary reference or fixture updates.

- `bench test --package ./internal/releasepreflight`

Mutation witness: Restore an unbounded scanner wait or accept promotion from child zero with unresolved cleanup. The real resistant scanner or protected promotion assertion must fail.
