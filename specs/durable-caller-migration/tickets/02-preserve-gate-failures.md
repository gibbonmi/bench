# Preserve durability failures through gate execution

Blocked by: none
Writes: internal/gate/verdict.go, internal/gate/engine.go, internal/gate/run_transaction.go, internal/gate/durable_replacement_test.go (new), internal/gate/run_durability.go (new)
Covers: D16, D17, D37, D38, D39, D40

## What to build

Migrate verdict, evidence, and lane replacement through accepted D-A.
Pass a per-invocation replacer through complete gate execution rather than fabricated results.
Preserve pending, timeout, evidence, and final categories while retaining the typed cause.
Interrupted-green demotion returns persistence failure to its diagnostic caller.
Keep pending recovery and report both causes when recovery also fails.

Review chunk: D-B1.
Complete acceptance and landing of durable-file-replacement supplies the shared leaf and its native qualification.
No ticket in this successor starts before that complete prerequisite lands.
Run this checkpoint with successor tickets unbuilt.
Preserve all existing assertions and successful behavior outside the approved delta.

## Acceptance

Fail the final directory sync after a real green or red oracle result, then independently fail pending recovery.
ActionExit is 1, the returned Inspection is not reusable green, and diagnostics retain both causes without claiming old pending bytes.
A failed initial pending write starts no oracle; a failed evidence write retains its classified cause.

- [ ] The gate verdict caller exposes an injected directory-sync failure. (D16).
- [ ] The gate evidence and lane caller exposes an injected directory-sync failure. (D17).
- [ ] Pending persistence failure prints its cause and prevents oracle execution. (D37).
- [ ] Final persistence failure retains its diagnostic cause at ActionExit 1. (D38).
- [ ] Green-evidence persistence failure prints the classified cause at ActionExit 1. (D39).
- [ ] Failed pending recovery reports both persistence errors. (D40).

## Checkpoint verification

- `bench test --package ./internal/gate`

## Required omission evidence

Discard the leaf cause, allow oracle start after pending failure, hide evidence failure, or discard the recovery error.
Each mutation must fail its real execution witness.

Record the exact failing assertion, restored source, and passing result in the review pickup.
A shared-owner test alone cannot prove its caller or transaction outcome.

## Source headroom

Move the cohesive persistence-failure and pending-recovery helpers into run_durability.go in this same checkpoint.
Keep run_transaction.go below its current growth boundary and preserve its execution ordering.
The extraction and new behavior land together; a later ticket cannot pay this headroom debt.
