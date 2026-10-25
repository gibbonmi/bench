# Grade each verification entry at its own commit

Blocked by: 01-exempt-the-assignment-block.md
Writes: internal/reviewrecord/record.go, internal/reviewrecord/parse.go, internal/reviewrecord/check.go, internal/reviewrecord/coverage.go, internal/reviewrecord/confirmation.go, internal/reviewrecord/record_test.go, internal/reviewrecord/source_test.go, internal/reviewrecord/fold_test.go, internal/reviewrecord/recordtest/fixture.go, internal/reviewrecord/recordtest/stage.go, internal/reviewrecord/recordtest/confirmation.go, internal/gate/checkpoint.go, internal/gate/review_checkpoint_test.go, internal/gate/review_checkpoint_commits_test.go, internal/gate/review_checkpoint_confirmation_test.go, internal/gate/delegated_checkpoint_test.go, internal/landing/completion_evidence_test.go, internal/landing/delegated_completion_test.go, internal/worktree/delegated_integration_test.go, internal/preflight/ownership/evidence.go, internal/preflight/ownership/premises_test.go, internal/preflight/ownership/premises_fixture_test.go, internal/reviewrecord/fold.go
Covers: RC8, RC13, RC14, RC15, RC17, RC18, RC19, RC20, RC21, RC22, RC23

## What to build

Review chunk: RC-C2. This chunk is a declared build (`Persisted fields: review-record:chunks.verification.source_commit`), so the orchestrator installs the candidate before the first dispatch.

The checkpoint and the reused-test premise accept a chunk verification entry at the commit that it graded, so a later ticket's commit does not make an earlier entry stale.

Add the field `source_commit` to each chunk verification entry (`internal/reviewrecord/record.go:32-38`). The parser (`validateVerification`, `internal/reviewrecord/parse.go:126`) refuses a value that is not a full object id. An absent value and an empty value are the same to the decoder, so both take the legacy rule. Completion entries keep their tree binding and carry no commit.

Add one exported `reviewrecord` predicate that decides whether an entry is current for a chunk. With a commit C, the entry is current when C strictly descends from the chunk base, C is an ancestor of the chunk tip, and the entry digest equals `SourceDigest` of C's tree. The ancestry query is reflexive, so refuse C equal to the chunk base by a separate test, as `FindReviewedPair` does (`internal/reviewrecord/confirmation.go:95-98`). With no commit, the entry is current when its digest equals the chunk digest. The command, performer, role, outcome, exit, and probe-bite rules stay as they are (`internal/reviewrecord/check.go:54-68`).

Adopt the predicate at two sites: the chunk path of `checkVerification` (`internal/reviewrecord/coverage.go:76-80`, `internal/reviewrecord/check.go:54`) and the reused-test premise (`internal/preflight/ownership/evidence.go:63,208`). The completion path at `check.go:96` keeps its digest rule and its refusal text, because `internal/reviewrecord/recordcmd/completion_test.go:68` pins "is incomplete, failed, or stale". Ticket 03 adopts the predicate at the writer and the stage check. Keep `Check`, `CheckTrees`, and `CheckSource` signatures as they are. `internal/preflight/ownership/evidence.go:23` already names a premise field `source_commit`; keep the two meanings apart in names and messages.

Bind the shared fixture's entries to their freeze commit (`Fixture.Verification` and `Fixture.freeze`, `internal/reviewrecord/recordtest/fixture.go:307-351`), so the fixture models the new record shape. An entry at the freeze tip is current under both rules, so a fixture test that passes today still passes. A test that stales a digest by hand must still refuse: `internal/gate/review_checkpoint_test.go:88`, `internal/gate/delegated_checkpoint_test.go:55`, `internal/landing/completion_evidence_test.go:181`, `internal/landing/delegated_completion_test.go:36`, and the relabel case `internal/worktree/delegated_integration_test.go:203-210` (RC15). `fixture.go` is at 397 of 400 lines and the commit lane grades growth, so move bytes to `recordtest/stage.go` or `recordtest/confirmation.go` when the change adds lines. The `Refreeze` change belongs to ticket 03.

The legacy rule keeps the two landed records valid (spec Implementation decisions, RC-C2), so no file under `reviews/` changes.

The planned checkpoint rows go in `internal/gate/review_checkpoint_test.go` (313 lines). Eight rows can pass the 400-line limit. Use one table where the rows share a fixture, or move an existing test to `internal/gate/review_checkpoint_commits_test.go` (195 lines). `internal/preflight/ownership/premises_test.go` is at 375 lines; move a helper or a test to `premises_fixture_test.go` when RC20 and RC21 need the room. The RC20 fixture needs an earlier-commit variant of `seedEvidenceCommits` (`premises_fixture_test.go:157-208`).

When the fixture change flips the "verification at the reviewed digest" case of `TestReviewCheckpointConfirmationRefusals` (`internal/gate/review_checkpoint_confirmation_test.go:255`), invert it to the RC17 behavior in this ticket. The checkpoint still runs the whole gate when every entry is current (RC22); this ticket does not change `internal/gate/checkpoint.go` unless RC22 shows otherwise.

## Acceptance

- [ ] RC8: the checkpoint accepts an earlier ticket's entry at that ticket's commit after a later ticket moves the tip.
- [ ] RC13: the checkpoint refuses an entry whose commit is not an ancestor of the chunk tip.
- [ ] RC14: the checkpoint refuses an entry whose commit is the chunk base.
- [ ] RC15: the checkpoint refuses an entry whose digest differs from the digest of its commit's tree.
- [ ] RC17: the checkpoint accepts another ticket's probe entry at its own commit after a repair moved the tip.
- [ ] RC18: the checkpoint accepts an entry with no `source_commit` at the chunk digest.
- [ ] RC19: the checkpoint refuses an entry with no `source_commit` at a digest other than the chunk digest.
- [ ] RC20: the reused-test premise accepts an entry at an earlier commit on the chunk line.
- [ ] RC21: the reused-test premise refuses an entry whose commit is not on the chunk line.
- [ ] RC22: the checkpoint goes red on a red planted at the chunk tip when every entry is current.
- [ ] RC23: the record parser refuses a `source_commit` that is not a full object id.
- [ ] The probe "drop the ancestry test of the entry commit" bites in `bench test --package ./internal/gate`.
