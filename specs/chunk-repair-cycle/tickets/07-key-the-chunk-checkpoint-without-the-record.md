# Key the chunk checkpoint on the committed tree without the record

Blocked by: 06-refuse-a-dirty-chunk-checkout.md
Writes: internal/gate/checkpoint.go, internal/gate/gate.go, internal/gate/engine.go, internal/gate/evaluation.go, internal/gate/run_transaction.go, internal/gate/gate_prose.go, internal/gate/subject.go, internal/gate/complete_checkpoint.go, internal/gate/lane.go, internal/gate/checkpoint_prepare_test.go (new), internal/gate/review_checkpoint_test.go, internal/gate/refusal_route_test.go, internal/gate/complete_checkpoint_test.go, internal/gate/commitment_completion_test.go, internal/gate/cache_env_test.go, internal/gate/run_outcomes_test.go, internal/gate/prospective_owner_test.go, internal/reviewrecord/recordcmd/stage_test.go, .bench/BENCH-reference.md, docs/adr/0002-accepted-enforcement-postures.md, internal/refusalroute/faces_gate.go, internal/refusalroute/route_test.go, internal/worktree/delegated_integration_test.go
Covers: RC39, RC40, RC41, RC46, RC48

## What to build

Review chunk: RC-C4.

The chunk checkpoint grades the HEAD commit, and its gate key excludes the review record. A chunk checkpoint after a record-only commit then reuses a whole green of the same subject and oracle. It still grades the evidence and the record prose on every call.

The subject is the HEAD tree without `reviews/<slug>.md`. `reviewrecord.SourceDigest` (`internal/reviewrecord/plan.go:201-210`) derives it, so the record exclusion has one source. Run that tree in a private checkout through the prospective artifact owner (`executeTreeWithOwner`, `internal/gate/engine.go:53-66,159-172`), with the worktree as the storage and identity root. The chunk path today runs `executeAfterAcquire` (`internal/gate/gate.go:209,335`). The oracle is the base oracle of that subject: the purpose JSON and the HEAD tip leave the chunk oracle (`internal/gate/checkpoint.go:277-285`). The complete form keeps its published-tree oracle and its completion frames (`internal/gate/complete_checkpoint.go:42`). Prospective landing evidence keeps its key (RC46, `internal/gate/prospective_owner_test.go`, at 480 lines over its limit, so keep that file line-neutral).

The tip-moved guard and the evidence pass still run on every chunk-checkpoint call (`internal/gate/evaluation.go:78-127`, `internal/gate/run_transaction.go:58-80,112-114`). The evidence pass reads the record from the HEAD tree, not from the stripped subject (`reviewrecord.CheckTrees` already takes separate trees, `checkpoint.go:273`).

Grade the prose of `reviews/<slug>.md` on every call. Read the record blob and the exclusion-policy blob from the HEAD tree and grade them through `prose.NewGraderFromIndex` and `GradeBytes`. The staged form of `bench gate-prose` is the precedent (`internal/gate/gate_prose.go:132-176`). Extract the shared bounded read and classify step there, so the two callers share one source. Add no second splitter and no exclusion row. The grade must run before the reuse return at `run_transaction.go:75-80`, or a reuse skips it. `run_transaction.go` is at 381 of 400 lines, and tickets 08, 10, and 11 also add to it. Add only the call there, and put the grade step in `gate_prose.go`.

`RunCheckpointArmed` (`internal/gate/review_checkpoint_test.go:52-57`, used at `internal/gate/refusal_route_test.go:251`) must follow the new chunk executor, so the tip-moved face test still exercises the real path. The fakes of the `gateEvaluation` interface (`internal/gate/cache_env_test.go:151-153`, `internal/gate/run_outcomes_test.go:26-34`) and `internal/gate/commitment_completion_test.go:145-147` change only if the interface changes. `internal/gate/lane.go:109` changes only if the owner signature changes.

Rewrite the first case of `TestReviewCheckpointReuse` (`review_checkpoint_test.go:112-125`): under the new key, a whole green over the same stripped subject and oracle can serve, because the evidence pass still runs. Keep its record-mutation case as RC39.

A recorded verification of a gate command stays a staged native excerpt with its digest and its source digest (`internal/reviewrecord/recordcmd/stage.go:56-82`). Reuse writes no verification entry (RC48).

Update the sentences that this key makes stale: `.bench/BENCH-reference.md:559` ("Checkpoint purpose and record bytes participate in the existing gate verdict identity.") and ADR 0002 posture 5 ("The subject binds the working tree"). Ticket 08 adds the prepared run to the same documents. An ADR holds no file path and no code.

## Acceptance

- [ ] RC39: a reused green with a stale verification entry refuses through the evidence pass.
- [ ] RC40: a record with a sentence over the prose bound reds the chunk checkpoint on a reused green.
- [ ] RC41: the same record reds the chunk checkpoint on a fresh run.
- [ ] RC46: prospective landing evidence still keys to the tree and the baseline runner identity (`TestProspectiveEvidenceKeysToTheTreeAndBaselineRunnerIdentity`).
- [ ] RC48: a plan requirement whose command runs the gate accepts only a staged native excerpt with a matching digest.
- [ ] A chunk checkpoint in a second process after a record-only commit reuses the earlier chunk green and runs no phase.
- [ ] The probe "frame the HEAD tip into the chunk-checkpoint oracle" bites, and the probe "skip the record prose grade after a reuse" bites, each in `bench test --package ./internal/gate`.
