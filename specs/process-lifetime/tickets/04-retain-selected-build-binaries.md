# Retain selected binaries for unresolved users

Blocked by: 01-own-lifetime-through-bounds.md, 02-protect-nested-bounds-users.md, 03-migrate-gate-process-policy.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, cmd/bench/process_resources.go (new), internal/conformance/axi_query_registry_test.go, internal/conformance/injected_ports_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/runbinary/lifetime_resources_test.go (new), internal/runbinary/lifetime_test.go (new), internal/runbinary/pdeathsig_linux_test.go, internal/runbinary/protection.go (new), internal/runbinary/runbinary.go, internal/runbinary/runbinary_test.go, internal/runbinary/sysprocattr_darwin.go, internal/runbinary/sysprocattr_linux.go, internal/runbinary/sysprocattr_other.go
Covers: PL34, PL35, PL36

## What to build

A real selected build drains its builder group and retains the executable while any registered user remains unresolved.
Migrate the builder through Lifetime with TERM and its existing two-second grace.
Preserve executable validation and the selected command inputs.

Selection creates its private resource generation before builder or executable use.
Expose its Descriptor through the existing selection composition path.
Selection.Close consumes SealAndObserve and returns unresolved release errors.
Build-failure cleanup consumes the same result before removing the selection directory.
The CLI table adds the selection owner's identity validator.

Consume the protection API and schema from ticket 2 without copying registration or recovery rules.
Preserve Linux Pdeathsig and Darwin Setpgid through the shared command adapter.
Do not claim Linux parent-death behavior on Darwin.
The final chunk performs native reconciliation; this ticket retains and exercises existing platform assertions.

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

- [ ] PL34: Builder normal completion drains its remaining group.
- [ ] PL35: Builder cancellation keeps TERM and its existing grace.
- [ ] PL36: Selection.Close retains a protected private executable.

## Verification

Run the named checks through their production entries.
For each owned row, record a behavioral omission or swap that fails, restoration, and the passing result.
A compilation failure is not a behavioral red.
Preserve existing assertions and fixture mutation purpose.
Co-owned closure holders receive only necessary reference or fixture updates.

- `bench test --package ./internal/runbinary`

Mutation witness: Remove normal builder drain or delete the selection before all registered groups are absent. The live descendant or executable-retention assertion must fail.
