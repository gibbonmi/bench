# Preserve probe recovery after failed mutation

Blocked by: 07-recover-approval-stage.md
Writes: internal/probe/subject.go, internal/probe/probe.go, internal/probe/replacement_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: D27, D34, D41, D42, D43

## What to build

Use D-A for subject mutation and restore through the real probe invocation.
Classify failed mutation before preservation release.
After mutation rename completes, attempt restore before a focused mutation run.
Release preservation only after successful synchronization and verified original-byte read-back.
Failed restore retains the copy, exits 2, and reports both causes with the preserved path.

Review chunk: D-C4.
Complete acceptance and landing of durable-file-replacement supplies the shared leaf and its native qualification.
No ticket in this successor starts before that complete prerequisite lands.
Run this checkpoint with successor tickets unbuilt.
Preserve all existing assertions and successful behavior outside the approved delta.

## Acceptance

Complete the baseline, publish mutated bytes, then fail directory sync.
No focused mutation run starts.
A successful restore releases preservation and returns the mutation error at exit 1.
An independent restore fault retains exact preserved bytes at exit 2.
The pre-rename fault keeps original bytes and the previous cleanup result.

- [ ] The probe subject caller exposes an injected directory-sync failure. (D27).
- [ ] A failed probe restore retains its preserved copy. (D34).
- [ ] A post-rename mutation failure with failed restore retains the preserved copy at exit 2. (D41).
- [ ] A post-rename mutation failure releases preservation only after verified restore. (D42).
- [ ] A pre-rename mutation failure retains the original subject and existing cleanup result. (D43).

## Checkpoint verification

- `bench test --package ./internal/probe`
- `bench test --package ./cmd/bench --run 'CommandRegistry'`
- `bench test --check axi-query-registry`
- `bench test --check subcommand-routing`

## Required omission evidence

Release preservation on every mutation error or treat failed restore sync as success.
The preservation-path and focused-start witnesses must detect each mutation.

Record the exact failing assertion, restored source, and passing result in the review pickup.
A shared-owner test alone cannot prove its caller or transaction outcome.

The co-owned command registries satisfy the existing package binding.
Keep command membership, routes, help grammar, and declared inventories unchanged.
Modify a registry only when the approved caller change requires its exact disposition.
