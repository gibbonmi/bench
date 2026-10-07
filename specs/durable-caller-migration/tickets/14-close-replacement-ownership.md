# Remove the unused replacement owner and close caller coverage

Blocked by: 02-preserve-gate-failures.md, 03-migrate-owner-records.md, 04-migrate-capture-documents.md, 05-migrate-handoff-documents.md, 06-preserve-ledger-effects.md, 07-recover-approval-stage.md, 08-migrate-publication-records.md, 09-migrate-broker-manifests.md, 10-migrate-dashboard-artifacts.md, 11-preserve-probe-recovery.md, 12-migrate-assessment-records.md, 13-migrate-repair-documents.md
Writes: internal/adopt/manifest.go
Covers: D30, D31, D33, D35, D36

## What to build

Delete only unused adopt.writeManifest after confirming its absent static consumers.
Retain manifestBytes, manifestMode, stageBytes, and the adoption lifecycle transaction.
Reconcile every migrated caller against its recorded payload, mode, refusal, and bypass evidence.
Review the excluded immutable, directory, rollback, rotation, cancellation, and lease owners without migrating them.

Review chunk: D-D.
Complete acceptance and landing of durable-file-replacement supplies the shared leaf and its native qualification.
No ticket in this successor starts before that complete prerequisite lands.
Run this checkpoint with successor tickets unbuilt.
Preserve all existing assertions and successful behavior outside the approved delta.

## Acceptance

The source census finds no migrated temporary-write and rename algorithm outside durablefile.
The real adoption stage still uses its existing publication protocol.
All migration checkpoints retain their accepted payload, mode, pre-rename refusal, and sentinel evidence.

- [ ] Each migrated caller retains its pre-write validation and sentinel category. (D30).
- [ ] Migrated callers contain no independent temporary-write and rename algorithm. (D31).
- [ ] Excluded protocols retain their current assertions and operation owners. (D33).
- [ ] The unused adopt writeManifest helper no longer defines a replacement algorithm. (D35).
- [ ] Differential fixtures preserve each migrated caller payload, mode, and pre-rename refusal result. (D36).

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
