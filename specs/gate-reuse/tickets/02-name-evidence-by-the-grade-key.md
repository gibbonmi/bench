# Name each evidence record by the grade key

Blocked by: 01-compute-the-grade-identity.md
Writes: internal/gate/engine.go, internal/gate/run_durability.go, internal/gate/run_transaction.go, internal/gate/verdict.go, internal/gate/checkpoint.go, internal/gate/authorization/authorization.go, internal/gate/grade_key_test.go (new), internal/gate/complete_checkpoint_test.go, internal/gate/checkpoint_prepare_test.go, internal/gate/review_checkpoint_test.go, internal/gate/timeout_recovery_count_test.go, internal/gate/prospective_owner_test.go, specs/landing-test-efficiency/spec.md, internal/worktree/land_journey_test.go, internal/gate/authorization/authorization_test.go
Covers: GR12, GR13, GR10, GR11, GR14, GR15, GR16, GR17, GR18, GR19, GR20, GR21, GR65, GR53, GR57

## What to build

The evidence store under the common Git directory names each whole-gate green by the grade key: the tree framed with the grade identity of ticket 01.
`evidenceName` derives that name, and `evidencePath`, `retainGreen`, and `invalidateEvidence` keep their signatures.
A green that one checkout records answers a run of the same key from any checkout of the repository.

The store record names its key in a `grade` field. A reader refuses a record whose `grade` field differs from its name, as changed evidence.
`EvidenceInspection` carries the grade key in place of its `Oracle` field, so the landing authorization token binds the key.

The run purpose leaves the store key. `applyCheckpoint` keeps the purpose frames in the per-checkout oracle digest, and the gate cache keeps that digest.
A complete checkpoint still runs its completion evidence pass on every call, before it accepts a reused green.

A store record answers only inside the freshness window of its own recorded time. A red or timed-out run removes the record under its key, whichever checkout wrote it.
The delta base read and the reader record keep their current rules, so neither answers a grade lookup.
Update the comments of `TestCompleteCheckpointPurposeOmitsPrepare` and `TestReviewCheckpointReuse` that name the purpose as a key member.

`TestProspectiveEvidenceKeysToTheTreeAndBaselineRunnerIdentity` becomes `TestProspectiveEvidenceKeysToTheTreeAndRunBinary`. Two trees name two records, and a changed run binary member or running executable digest names another record.
`prospective_owner_test.go` is over its line limit, so its edit does not grow the file.
`checkpoint.go` holds 389 lines and `run_transaction.go` holds 370 lines, and each stays at or under its 400-line limit.

## Acceptance

- [ ] Two trees name two evidence records, and a changed run binary member or running executable digest names another record.
- [ ] A run of one tree from the primary checkout after a green in a linked worktree prints `gate: green (fresh verdict reused for this tree)` and runs no phase.
- [ ] A changed posture or baseline phase manifest runs the gate again.
- [ ] After an unpurposed green of the published tree, the complete checkpoint of that tree runs no phase, and a completion evidence defect still refuses it.
- [ ] A green older than its window, a green that a red run of the key removed, and a record with a wrong `grade` field each answer no reuse.
- [ ] A timed-out run removes the green, a delta checkpoint still reads its base without expiry, and neither a reader record nor a delta green answers.
