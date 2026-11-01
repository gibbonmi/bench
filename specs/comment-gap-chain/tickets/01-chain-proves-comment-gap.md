# Prove a comment gap in the review chain

Blocked by: none
Writes: internal/reviewrecord/coverage.go, internal/reviewrecord/fold.go, internal/reviewrecord/fold_test.go, internal/reviewrecord/recordtest/fixture.go, internal/reviewrecord/recordcmd/completion_test.go, internal/gate/review_checkpoint_commits_test.go, internal/gate/delta_checkpoint_test.go
Covers: CG01, CG02, CG03, CG04, CG05, CG06, CG07, CG08, CG09, CG10, CG11, CG12, CG13, CG14, CG15, CG16, CG17

## What to build

A later chunk base after a comment-only Go correction passes the review chain. The orchestrator then commits the correction with no refreeze, and the frozen chunk keeps its reviews and its verifications.

Add one unexported predicate in `internal/reviewrecord`. It answers whether a later digest follows a covered digest: the digests are equal, or `commentgap.Prove` accepts the gap from the covered digest to the later digest. The `follows` helper already holds this rule (`internal/reviewrecord/fold.go:136-141`), so the new predicate takes its digest form, and `follows` reads it. Put the predicate in `fold.go` or `coverage.go`, because `bench structure` grades `internal/reviewrecord` at 27 source files against a maximum of 12.

Four sites read the predicate:

- the equality branch of the chain in `checkSource` (`coverage.go:88`);
- `foldCoversBase` (`fold.go:115`);
- the covered-digest check of `checkSource` (`coverage.go:150`);
- `follows` (`fold.go:138-141`).

After the change, `commentgap.Prove` has one call site in `internal/reviewrecord`. Keep both ancestry conjuncts: `coverage.go:89` for the chunk and `fold.go:118` for the fold commit. Keep the order of the chain: the previous chunk digest, then the covering fold, then the `FormFold` refusal.

Change the rule clause of the refusal at `coverage.go:99` to `only record commits, recorded folds, and proved comment-only Go corrections follow a chunk tip until its checkpoint`. Keep the rest of the sentence. `TestReviewCheckpointChainGapNamesTheExpectedBase` pins the new clause (`internal/gate/review_checkpoint_commits_test.go:99`). `CHANGELOG.md:238` keeps the old clause as history, and ticket 02 owns `CHANGELOG.md`.

Add `TestChainAcceptsAProvedCommentGap` to `internal/reviewrecord/fold_test.go`. Mirror `foldedAfterChunk1` and `foldRepair` (`fold_test.go:16-33`), and call `reviewrecord.CheckSource` for chunk 2. Its cases are CG04 to CG09. For CG09, make the chunk 2 tip a two-parent commit of the off-line base and the chunk 1 tip that keeps the base tree. That tip passes the pair guard at `coverage.go:45-53`, so only the chain ancestry conjunct refuses.

Add `TestReviewCheckpointChainAcceptsACommentGap` to `internal/gate/review_checkpoint_commits_test.go` for CG01 to CG03. Mirror `TestReviewCheckpointCommentOnlyGap` with `seedCommentGap` and `correctCommentGap` (`review_checkpoint_commits_test.go:104-178`). For CG02, compare each field of the retained chunk 1 entry with the frozen entry. Change the gap of the `chain base` case of `TestReviewCheckpointKeepsStrictEvidence` from a comment edit to a Go statement edit, so that the case keeps its `expected base` refusal.

Write each expected refusal and each expected exit code as a literal. No expectation calls the predicate under test.

`recordtest/fixture.go` has 5 lines of headroom under the 400-line maximum, and `internal/gate/delta_checkpoint_test.go` has 2. Prefer a local helper in the test file. `recordcmd/completion_test.go` (CG13) and `delta_checkpoint_test.go` (CG17) hold rows that must stay green with no edit.

## Acceptance

- [ ] `bench gate --checkpoint <spec> --chunk 2` exits 0 when the chunk 2 base follows the chunk 1 tip by a Go comment edit.
- [ ] After that checkpoint, the chunk 1 entry keeps its frozen Base, Tip, SourceDigest, Reviews, and Verification.
- [ ] `bench gate --checkpoint <spec> --complete` exits 0 over the same two chunks.
- [ ] `reviewrecord.CheckSource` accepts a comment gap after the chunk 1 tip and after a recorded fold.
- [ ] `reviewrecord.CheckSource` refuses with `stale review chain gap` for a statement gap, a directive gap, a Markdown gap, and an off-line base.
- [ ] The chain refusal holds the new rule clause, and `commentgap.Prove` has one call site in `internal/reviewrecord`.
- [ ] The covered-digest check, the failed ancestry query, the fold refusals, the review-pair refusal, and the delta-checkpoint whole run keep their current behavior.
- [ ] Each mutation of the completion-plan probes for CG-C1 turns its named test red.
