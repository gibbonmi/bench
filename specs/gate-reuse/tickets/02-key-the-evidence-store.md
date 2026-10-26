# Key the prospective evidence store by the authorization key

Blocked by: 01-leave-the-checkout-out-of-the-identity.md
Writes: internal/gate/subject.go, internal/gate/engine.go, internal/gate/run_durability.go, internal/gate/run_transaction.go, internal/gate/checkpoint.go, internal/gate/verdict.go, internal/gate/verdict_registry.go, internal/gate/verdict_registry_guard_test.go, internal/gate/prospective_owner_test.go, internal/gate/authorization_key_test.go (new)
Covers: GR10, GR11, GR12, GR13, GR14, GR15, GR16, GR17, GR18, GR19, GR20, GR21

## What to build

The subject carries the authorization key beside the oracle digest, in a new field of the `subject` type. The key is the digest of the tree and the oracle identity of ticket 01.
`buildSubjectForTree` sets the unpurposed key in the same pass that sets `Oracle` and `Identity`, so every subject path has a key: `buildSubject`, `ValidateProjectGreen`, and the prospective owner. Count the headroom of `internal/gate/subject.go` again after ticket 01, against its grant of 513 lines.

`applyCheckpoint` frames a run's purpose over the key as it frames it over the oracle digest today. The purpose is the checkpoint value, the completion folder, and the source tip.
One function frames the purpose for both digests, so the two framings cannot drift. A chunk checkpoint keeps the unpurposed key, as it keeps the unpurposed oracle digest.

`evidenceName` names each store record by the key. The store record carries the key in an `authorization` field and no `oracle` field, in its own record class of the verdict registry. The registry guard test lists that class.
A reader refuses a record whose `authorization` value differs from its own key, as changed evidence.

`retainGreen`, `invalidateEvidence`, `inspectEvidenceWindowed`, and `reusableEvidence` reach the store only through `evidencePath`. Keep the signature `evidencePath(gitdir, plan)`, because five gate test files call it for the store directory alone.

The gate cache record (`bench-last-gate`) keeps its `oracle` field and its class. This ticket writes no gate cache on a reuse, and ticket 03 adds that write.

The capability posture now binds every reuse, because the key frames the oracle identity. This reverses FT449's posture rule, by reviewer approval on 2026-10-26.

The comment of `TestProspectiveEvidenceKeysToTheTreeAndBaselineRunnerIdentity` says that the subject binds the identity-root path. Correct that comment with no new line, because the file holds 480 lines over its 400-line limit.

The new tests go in `authorization_key_test.go`. They use `outcomeFixture`, whose gate counts runs in the common Git directory, and a linked worktree that `runGitIn` creates.
GR16 goes in the same file, because `complete_checkpoint_test.go` holds 357 lines. It calls `closureSource` and the published-tree derivation of `publishedEvidence`, which the package's test files share.

`run_transaction.go` holds 370 of 400 lines, and `checkpoint.go` holds 389. Put a new helper in `engine.go` or in a new file of `internal/gate`.

## Acceptance

- [ ] After a green prospective run of one tree in a linked worktree, a prospective run of that tree from the primary checkout keeps the run count at 1 (GR10).
- [ ] That primary run prints exactly `gate: green (fresh verdict reused for this tree)` (GR11).
- [ ] After a green under the relaxed posture, a run with `BENCH_REQUIRE_CAPABILITIES=1` runs the gate, and the run count reads 2 (GR14).
- [ ] After a linked-worktree green, a primary run whose baseline phase manifest differs runs the gate (GR15).
- [ ] After an unpurposed green prospective run of the published tree, the complete checkpoint of that tree runs the gate, and the run count reads 2 (GR16).
- [ ] A linked-worktree green whose recorded time is older than the freshness window does not answer the primary run (GR17).
- [ ] A red primary run removes a linked-worktree green of the same key, and the next worktree run runs the gate (GR18).
- [ ] An evidence record whose `authorization` field differs from its file name does not answer, and the run runs the gate (GR20).
- [ ] `TestProspectiveEvidenceKeysToTheTreeAndBaselineRunnerIdentity` (GR12, GR13), `TestGateRunTimeoutInvalidatesOldEvidence` (GR19), and `TestDeltaCheckpointReadsTheBaseWithoutExpiry` (GR21) stay green.
