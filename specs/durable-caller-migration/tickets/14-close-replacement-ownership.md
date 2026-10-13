# Remove the unused replacement owner and close caller coverage

Blocked by: 02-preserve-gate-failures.md, 03-migrate-owner-records.md, 04-migrate-capture-documents.md, 05-migrate-handoff-documents.md, 06-preserve-ledger-effects.md, 07-recover-approval-stage.md, 08-migrate-publication-records.md, 09-migrate-broker-manifests.md, 10-migrate-dashboard-artifacts.md, 11-preserve-probe-recovery.md, 12-migrate-assessment-records.md, 13-migrate-repair-documents.md
Writes: internal/adopt/manifest.go, internal/adopt/lifecycle_preview_test.go, internal/adopt/lifecycle_consumption_test.go, internal/durablefile/durablefiletest/assert.go, internal/gate/durable_replacement_test.go, internal/gate/prospectiveartifact/replacement_test.go, internal/capturetx/replacement_test.go, internal/handoffdoc/replacement_test.go, internal/intent/replacement_test.go, internal/commitment/repository/replacement_test.go, internal/publication/replacement_test.go, internal/brokermanifest/replacement_test.go, internal/probe/replacement_test.go, internal/assessment/replacement_test.go, internal/assessment/record_test.go, internal/repairpilot/replacement_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: D30, D31, D33, D35, D36

## What to build

Delete only adopt.writeManifest, which has no production consumer.
Two adoption tests seed a manifest through it; seed those manifests from manifestBytes and manifestMode without a replacement algorithm.
Retain manifestBytes, manifestMode, stageBytes, and the adoption lifecycle transaction.
Reconcile every migrated caller against its recorded payload, mode, refusal, and bypass evidence.
Review the excluded immutable, directory, rollback, rotation, cancellation, and lease owners without migrating them.

Close D30 and D36 for every migrated caller family in its caller test.
Add one shared refusal assertion to durablefiletest and move each pasted refusal block to it.
Each caller family gets one payload or destination refusal witness and one pre-rename row.
Each pre-rename row reads back through RequireBytes or RequireAbsent.

The co-owned command registries satisfy the existing package binding; keep command membership, routes, help grammar, and inventories unchanged.
The source census also confirms that each exported caller entry binds the native replacer. No caller test can detect a plain write in its place.

Review chunk: D-D.
Complete acceptance and landing of durable-file-replacement supplies the shared leaf and its native qualification.
No ticket in this successor starts before that complete prerequisite lands.
Run this checkpoint with successor tickets unbuilt.
Preserve all existing assertions and successful behavior outside the approved delta.

## Acceptance

The source census finds no migrated temporary-write and rename algorithm outside durablefile.
The real adoption stage still uses its existing publication protocol.
All migration checkpoints retain their accepted payload, mode, pre-rename refusal, and sentinel evidence.

- [x] Each migrated caller retains its pre-write validation and sentinel category. (D30).
- [x] Migrated callers contain no independent temporary-write and rename algorithm. (D31).
- [x] Excluded protocols retain their current assertions and operation owners. (D33).
- [x] The unused adopt writeManifest helper no longer defines a replacement algorithm. (D35).
- [x] Differential fixtures preserve each migrated caller payload, mode, and pre-rename refusal result. (D36).

## Checkpoint verification

- `bench test --package ./internal/adopt`
- `bench coverage --check durable-caller-migration`

## Required omission evidence

Retain the unused writeManifest body or restore any migrated local algorithm.
The review-owned definition and consumer census must reject that source.
Each existing caller bypass probe supplies its independent runtime witness.
The adopt package command verifies its surviving protocol, not the review-owned global algorithm census.
Consume the prerequisite review-record preservation evidence beside each successor caller family.

Record the exact failing assertion, restored source, and passing result in the review pickup.
A shared-owner test alone cannot prove its caller or transaction outcome.
