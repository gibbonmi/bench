# Fast-forward an empty run's recorded base on checkpoint and start

Blocked by: restart-reads-live-marker.md, classify-decayed-checkouts-delete-exemption.md
Ownership fence: `internal/specbuild`
Assumptions: the recompose refusal returns a populated subject from `preconditions` (`internal/specbuild/precondition.go:103-108`) and must keep doing so for promote; `recognizedAdvance` (precondition.go:111) and `updateRef` (state.go:290) are the reused guards; the CAS-then-save sequence to mirror is `finishRecomposition` (`internal/specbuild/recompose.go:34-48`); "empty" is `run.CandidateTip == run.Base` plus no assignment carrying checkpoint evidence (`Checkpoint == "" && CheckpointRef == ""`); claims re-derived from the tree at pickup

## What to build

For a non-terminal run where no assignment carries checkpoint evidence and the
candidate tip still equals the recorded base, when the working tip is a
recognized descendant of the recorded base, **checkpoint and start** — a
closed list — fast-forward the recorded base and candidate (durable candidate
ref moves by compare-and-swap from its old tip, then record save) and proceed;
start then reports status instead of erroring. The fast-forward lives behind
the closed operation list, never in the shared precondition return path:
promote keeps routing `errRecompose` into recomposition with its Bootstrap,
review refuses without mutation, assign keeps today's refusal, terminal runs
keep restarting. No gate runs on the fast-forward. A run with any checkpoint
evidence keeps today's refusal routing to promote; a non-ancestor tip stays a
subject-mismatch refusal everywhere. Re-scope exactly the four tests the spec
authorizes (`TestNonAbandonMutationsStillRecomposeOnMovedTip` start and
checkpoint cases; `TestLifecycleMutatorsRefuseSharedPreconditionDriftWithoutMutation`
working-advance checkpoint case; `TestStartResumeAndConflictsDoNotDuplicateRun`
moved-tip subtest; `TestRecompositionErrorIsStable` stays green unchanged as
the assign control) — each keeps its scenario and gets the new expectation;
no other existing test may be edited. RM5 in
`restart-reads-live-marker.md`'s terminal-restart control must stay green.

## Acceptance

- [ ] [FF1] On an empty non-terminal run whose tip advanced, Checkpoint fast-forwards base and candidate and proceeds (today: recompose refusal).
- [ ] [FF2] Start on the same state reports status instead of erroring, with base and candidate advanced.
- [ ] [FF3] After a fresh Service reload, record base, candidate tip, and the durable ref agree with the fast-forwarded values.
- [ ] [FF4] A fast-forward against a candidate ref that moved externally refuses without mutating the record (CAS, not overwrite).
- [ ] [FF5] A run with real checkpoint evidence and a moved tip still refuses checkpoint with the promote route (new checkpointed fixture).
- [ ] [FF6] Promote on an empty moved-tip run still routes through recomposition and its Bootstrap (existing tests unchanged).
- [ ] [FF7] Review on an empty moved-tip run still refuses without mutation (new fixture, state snapshot identical).
- [ ] [FF8] Assign on an empty moved-tip run still refuses with the promote route (`TestRecompositionErrorIsStable` unchanged).
- [ ] [FF9] A non-ancestor tip still refuses as subject mismatch, empty run or not.
- [ ] [FF10] Zero Bootstrap calls across the fast-forward (counting gate).
- [ ] [FF11] A second mutation after the fast-forward is a no-op advance: one ref move, stable state across the second call.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| FF1 | drop the fast-forward, keep the refusal | the empty-run checkpoint test | apply, run `go test ./internal/specbuild -run Checkpoint`, expect the recompose-refusal failure |
| FF4 | update the candidate ref unconditionally (no old value) | the moved-ref CAS test | apply, run it, expect the missing-refusal failure |
| FF5 | fast-forward every descendant tip regardless of checkpoint evidence | the checkpointed-fixture test | apply, run it, expect the missing-refusal failure |
| FF6 | place the fast-forward in the shared precondition return path | the promote Bootstrap control and FF7's review fixture | apply, run `go test ./internal/specbuild`, expect the promote/review controls to fail |
| FF9 | key the fast-forward on "tip differs" instead of recognized descendant | the non-ancestor test | apply, run it, expect the adopted-rewrite failure |
| FF10 | call Bootstrap inside the fast-forward | the counting-gate test | apply, run it, expect the nonzero-call failure |
| FF11 | re-fire the fast-forward on equal tips | the idempotency test | apply, run it, expect the second-ref-move failure |
