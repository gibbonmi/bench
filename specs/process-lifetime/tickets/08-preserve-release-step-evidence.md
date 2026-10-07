# Retain release evidence after incomplete phase shutdown

Blocked by: 01-own-lifetime-through-bounds.md, 02-protect-nested-bounds-users.md, 03-migrate-gate-process-policy.md, 04-retain-selected-build-binaries.md, 05-protect-private-test-directories.md, 06-refuse-clean-of-protected-cache.md, 07-retain-prospective-gate-artifacts.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, cmd/bench/process_resources.go (new), internal/conformance/axi_query_registry_test.go, internal/conformance/injected_ports_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/preprelease/lifetime_test.go (new), internal/preprelease/preprelease.go, internal/releaseevidence/evidence_promotion.go, internal/releaseevidence/protection.go (new), internal/releaseevidence/protection_test.go (new), internal/releaseevidence/release_evidence.go
Covers: PL61, PL119

## What to build

A real prep-release step retains its evidence after resistant descendant teardown is incomplete.
Migrate the step through Lifetime with TERM and the accepted missing two-second grace.
Preserve step order, attribution, normal policy, and existing publication authority.

Create the release-evidence store in the repository common directory before any phase launches.
Bind its generation to that repository and the canonical dist/preflight target.
The store stays outside exchanged evidence and temporary stage directories.
Expose its Descriptor and install the owner identity adapter in the existing CLI table.

FinalizeEvidence, promotion, and abandoned-stage cleanup consume SealAndObserve before replacement or deletion.
Unresolved protection cannot authorize publication.
Keep the existing atomic directory exchange; this ticket adds no replacement algorithm.
Later preflight callers consume this same owner.

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

- [ ] PL61: Prep-release retains evidence after incomplete phase shutdown.
- [ ] PL119: Abandoned release-stage cleanup retains protected evidence.

## Verification

Run the named checks through their production entries.
For each owned row, record a behavioral omission or swap that fails, restoration, and the passing result.
A compilation failure is not a behavioral red.
Preserve existing assertions and fixture mutation purpose.
Co-owned closure holders receive only necessary reference or fixture updates.

- `bench test --package ./internal/preprelease`
- `bench test --package ./internal/releaseevidence`

Mutation witness: Keep the store inside a swapped stage or delete an abandoned protected stage from its owner marker alone. The resistant-child evidence sentinel must survive and the erroneous success assertion must fail.
