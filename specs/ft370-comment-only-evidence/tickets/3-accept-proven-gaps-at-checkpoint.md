# Accept a proven comment-only gap at the checkpoint

Blocked by: 2-prove-comment-only-gaps.md
Writes: internal/reviewrecord/recordcmd/chunk_test.go, internal/reviewrecord/recordcmd/verification_test.go, internal/reviewrecord/coverage.go, internal/gate/review_checkpoint_commits_test.go, internal/reviewrecord/recordcmd/completion_test.go, internal/landing/completion_evidence_test.go
Covers: CG1, CG2, CG3, CG4, CG5, CG6, CG7, CG8, CG33, CG34, CG35, CG36, CG44, CG45

## What to build

Change `reviewrecord.checkSource` at one place: the comparison of the last
matched chunk's source digest with the graded source digest. When the two
digests differ, call `commentgap.Prove` from ticket 2. The chunk's source
digest is the `reviewed` operand, and the graded source digest is the `later`
operand. A nil answer accepts the gap, and the check continues as today.

A refusal keeps the current error text, its suffix, and its range rule:
`chunk <id>: stale reviewed source: no chunk review covers <range>`. The
checkpoint does not print the rule sentinel. The per-chunk recompute of the
digest from the chunk tip commit stays before the proof. So the proof reads
only a reviewed tree that the chunk tip authenticates.

Every other row of the spec's `The staleness sites` table stays strict. Change
the chain-gap refusal clause `only record commits follow a chunk tip` to
`only record commits and comment-only corrections follow a chunk tip`. The rest
of that message stays.

The one change serves `--chunk`, `--complete`, the landing broker through
`gate.WithCompletion`, and `bench record completion` through
`RecordCompletion`. Ticket 4 states the chain-gap clause in the implementation
phase, so the two texts agree.

Keep each FT358 ticket 7 contract from the spec's `Overlap with FT358` note:

- Add no file to `internal/reviewrecord`.
- Add no use of `fenceMarker`, and no string literal with a run of three backticks to `internal/reviewrecord`.
- Do not edit `parse.go`, `write.go`, `record_test.go`, or `recordtest/fixture.go`.

Build the tests that the spec's `Testing decisions` state:

- `TestReviewCheckpointCommentOnlyGap` holds CG1, CG2, CG4, and CG5. The fixture gets its Go file through the `prepare` hook of `attachedCheckpointFixture`.
- `TestReviewCheckpointKeepsStrictEvidence` holds CG6, CG7, CG8, CG33, CG34, and CG36.
- `TestReviewCheckpointChainGapNamesTheExpectedBase` gains the new clause in its second assertion (CG35).
- `TestLandingCompletionEvidence` gains one case (CG3). `attachedCompletionFixture` gains a variadic prepare hook. `newCompletionFixture` forwards the hook, and the delegated caller stays unchanged.
- `TestRecordCompletionAcceptsCommentOnlyGap` holds CG44 and CG45.
- Move `recorded` from `verification_test.go` beside the chunk fixture helpers in `chunk_test.go`. Add a variadic prepare hook, which `readyCompletion` forwards.
  This seeds Go before the chunk without copying setup or exceeding the verification file budget.

The three checkpoint tests go in
`internal/gate/review_checkpoint_commits_test.go`. Each gap test edits the Go
file with `Write` and `Commit`. CG34 needs two chunks, and
`attachedCheckpointFixture` records one. So build CG34 on `outcomeFixture` in
`run_outcomes_test.go`, as `TestReviewCheckpointChainGapNamesTheExpectedBase`
does.

## Acceptance

- [ ] A `--chunk 1` checkpoint exits 0 after a commit that changes only a comment in a `_test.go` file.
- [ ] An accepted gap raises the oracle run count by one.
- [ ] A `--complete` checkpoint exits 0 when the completion evidence names the corrected source.
- [ ] A chunk whose Standards result carried a finding passes after the correction and one superseding Standards pass at the reviewed frozen pair.
- [ ] A Go statement change after the reviewed tip refuses with `chunk 1: stale reviewed source: no chunk review covers`, and the oracle does not run.
- [ ] A committed comment fix beside an uncommitted Go statement change refuses with `chunk 1: stale reviewed source`.
- [ ] A chunk entry whose digest names the corrected tree, with its tip at the reviewed commit, refuses with `chunk 1: stale source digest`.
- [ ] Completion evidence recorded before a comment-only correction refuses with `completion is incomplete or stale`.
- [ ] A second chunk whose base follows a comment-only correction of the first chunk refuses with `expected base`.
- [ ] A chunk re-recorded at the corrected tip, with its three axis results at the reviewed pair, refuses with `chunk 1: stale Standards`.
- [ ] The chain-gap refusal holds `and only record commits and comment-only corrections follow a chunk tip`.
- [ ] A reviewed landing publishes `Status: implemented` over a comment-only correction after the last chunk.
- [ ] `bench record completion` writes a completion at the corrected digest, and it refuses a statement gap with exit 1 and unchanged record bytes.
- [ ] An acceptance that passes every differing digest without a call to `commentgap.Prove` makes the CG6 row fail.
- [ ] Each file in `Writes:` stays under 400 lines.
