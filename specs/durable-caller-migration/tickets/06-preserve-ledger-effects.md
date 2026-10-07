# Preserve matching effects after ledger publication

Blocked by: none
Writes: internal/intent/intent.go, internal/intent/transaction.go, internal/intent/assignment.go, internal/intent/replacement_test.go (new), tests/canary/package-core-guard/bounds-intent-window-fixed
Covers: D22, D44, D45, D52

## What to build

Add the accepted TransactWithReplacer form and migrate writePath through D-A.
Keep one address, lock, decision, replacement, compensation, and release implementation.
Compensate exactly once only when ledger rename did not complete.
A published ledger failure retains the decision effects and returns classified failure.
Exercise the real reauthorization request transition through the same invocation seam.

Review chunk: D-B5.
Complete acceptance and landing of durable-file-replacement supplies the shared leaf and its native qualification.
No ticket in this successor starts before that complete prerequisite lands.
Run this checkpoint with successor tickets unbuilt.
Preserve all existing assertions and successful behavior outside the approved delta.

## Acceptance

A reversible transition changes an external witness and the request digest.
Fail ledger directory sync after completed rename.
The new ledger and matching witness remain, compensation does not run, the lock releases, and the call returns uncertainty.
Fail before rename separately and keep the existing old-byte and compensation-count assertions.

- [ ] The intent ledger caller exposes an injected directory-sync failure. (D22).
- [ ] A post-rename ledger failure skips compensation. (D44).
- [ ] A pre-rename ledger failure compensates exactly once. (D45).
- [ ] A post-rename reauthorization failure retains the next request and its matching external effect. (D52).

## Checkpoint verification

- `bench test --package ./internal/intent`

- `bench test --check bounds-policy`

## Required omission evidence

Compensate every write error or omit the reauthorization effect check.
The new-ledger and external-witness test must detect their mismatch.

Record the exact failing assertion, restored source, and passing result in the review pickup.
A shared-owner test alone cannot prove its caller or transaction outcome.

## Existing fixture closure

Keep bounds.VerdictWindow(bounds.IntentLockTimeout) in its existing intent-lock owner.
The exact canary directory joins this ticket because its BASE pins the edited source.
Preserve CHECK, EXPECT, and MUTATE.json policy and their existing red purpose.
Do not rewrite the fixture to permit a wrong bounds class or copied Git administration flag.
