# Delete the three partial verdict classes

Blocked by: none
Writes: internal/gate/verdict_registry.go, internal/gate/verdict.go, internal/gate/composed_green.go, internal/gate/verdict_registry_guard_test.go, internal/gate/verdict_reason_test.go, internal/status/status.go, internal/status/status_gatecache_test.go, internal/status/status_fixtures_test.go, internal/status/status_signals_test.go, internal/status/status_producible_test.go, internal/preprelease/preprelease.go
Covers: none

## What to build

The quality survey of 2026-09-29, card 03, found that the gate verdict store keeps three partial record classes that nothing writes. These classes are the partial verdict, the check-partial verdict, and the combined-partial verdict. The component partition that wrote them is retired. Only readers remain, and the status tests make partial records by hand to keep those readers alive.

Delete the three classes, their validators, and the skip-evidence validators that only they use. Delete the partition types on the gate inspection and the partial fields of the record. Delete the partial-green row in the status board and the partial branch of the prep-release refusal.

After the change, the store reads an old partial record as an invalid cache record, because field matching is exact. Such a record still fails closed: it is never a reusable green, and the status board and prep-release send the reader to `bench gate`. The profile sentence about legacy input classes that fail closed stays true. The composed-green check asks the class for its own reuse posture instead of the partial fields.

The reviewer approved this light path on 2026-09-30, although the change alters what the status board and prep-release show for an old partial record.

## Acceptance

- [ ] The verdict store has three record classes: full verdict, pending, and lane record.
- [ ] No Go source names `Partition`, `CheckPartition`, `ComponentSkip`, `skipEvidence`, or a partial verdict class.
- [ ] A legacy partial, check-partial, or combined-partial record reads as an invalid cache record in the gate inspection.
- [ ] The status board reads a legacy partial record as an invalid cache, and it has no partial-green row.
- [ ] The prep-release refusal has no partial branch.
- [ ] The gate, status, and preprelease package tests and root conformance pass.
