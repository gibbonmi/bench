# Record each ticket's entries at its own commit

Blocked by: 02-grade-entries-at-their-own-commit.md
Writes: internal/reviewrecord/change.go, internal/reviewrecord/ticket.go, internal/reviewrecord/write.go, internal/reviewrecord/repair_test.go, internal/reviewrecord/fold_test.go, internal/reviewrecord/recordcmd/ticket.go, internal/reviewrecord/recordcmd/ticket_test.go, internal/reviewrecord/recordcmd/chunk_test.go, internal/reviewrecord/recordcmd/verification.go, internal/reviewrecord/recordcmd/verification_test.go, internal/reviewrecord/recordcmd/verification_import_test.go, internal/reviewrecord/recordtest/confirmation.go, internal/reviewrecord/recordtest/fixture.go, internal/gate/review_checkpoint_confirmation_test.go, cmd/bench/record_operations_test.go, .agents/commands/bench-implement-spec.md, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, .bench/BENCH-reference.md, .agents/skills/bench-craft-line/references/bounded-repair-policy.md
Covers: RC9, RC10, RC11, RC12, RC16, RC24

## What to build

Review chunk: RC-C2.

`bench record ticket` writes each ticket's verification and probe entries at the commit of the staged run, so a tip move makes no ticket stage again. A repair session runs each owed check once at its repair commit, and its record call refreezes the chunk at that commit.

The writer (`changeVerification`, `internal/reviewrecord/change.go:106-160`) and the stage check (`ownedRuns.calls`, `internal/reviewrecord/recordcmd/ticket.go:123-137`) adopt the predicate that ticket 02 exports. Each written entry carries `source_commit` equal to the staged run's HEAD (`internal/responsebound/retain.go:94-97`). A record-only commit never moves the chunk tip, so a staged run at a record-only descendant of the tip binds the tip's commit and leaves the tip unchanged. The stage check refuses a staged run whose commit is not on the chunk line. The dirty-path refusal stays (`recordcmd/ticket.go:133-137`).

`RecordTicket` (`internal/reviewrecord/ticket.go:86-181`) drops the rule that only the last ticket freezes. Any ticket's record call at a commit that descends from the chunk tip refreezes the chunk at that commit, and the first record call of a chunk creates the chunk entry. `lastFreezes`, the `create` and `bind` closures, the `FormChunk` step, and their doc comments (`ticket.go:27-28,79-85,98-143,177-181`) change with it. `refusePlanned` runs only for the last ticket today (`recordcmd/ticket.go:52-56`); keep the planned-cell refusal on every call that freezes. `changeChunk` keeps earlier results on a refreeze (`change.go:73-77`), and `bench record chunk` also calls it, so `recordcmd/chunk_test.go` must stay green.

`bench record verification` reaches the same writer through `RecordVerification` (`internal/reviewrecord/write.go:168-172`, `recordcmd/verification.go:78`). Bind its chunk entries to the commit of their source by the same rule, so that one writer rule serves both forms. `VerificationCall` (`write.go:159-164`) has consumers in `fold_test.go:161` and `repair_test.go:29`.

The shared fixture's `Refreeze` (`internal/reviewrecord/recordtest/confirmation.go:17-26`) stops re-retaining entries, and its doc comment changes. Its four callers in `internal/gate/review_checkpoint_confirmation_test.go` (lines 52, 92, 122, 176) then rely on the commit binding of ticket 02. If this change flips the "verification at the reviewed digest" case (line 255), invert it here to the RC17 behavior.

These refusal texts change, and each pin moves in this ticket:

- "the staged run is at %s, not the chunk source %s" (`recordcmd/ticket.go:130`), pinned at `recordcmd/ticket_test.go:146,155`. Rewrite both cases to the off-line refusal.
- "chunk N has no entry: the result needs a current chunk entry" for an earlier ticket with no chunk entry, pinned at `recordcmd/ticket_test.go:169-175` and through `ErrNoChunkEntry` at `internal/reviewrecord/repair_test.go:39,48`. The first record call now creates the entry, so rewrite or remove each case.
- Check the count pin `len(chunk.Verification) != 2` at `recordcmd/ticket_test.go:302-303`.

`internal/reviewrecord/recordcmd/ticket_test.go` is at 400 of 400 lines, and the commit lane grades growth. The planned rows RC9 to RC12 and RC16 go in that file, so move existing tests to `recordcmd/chunk_test.go` (323 lines) or another existing sibling first. Do not add a file to the crowded `recordcmd/` directory. `recordtest/fixture.go` is at 397 lines.

Remove the restage sentences from the Land section of `.agents/commands/bench-implement-spec.md`: at line 66, "Each ticket author then stages its verification and probe record at that chunk source" and "In a multi-ticket chunk, the last author commits no record before each earlier author stages its record"; at line 68, "A repair session also retains its ticket verification entries at the final chunk source before the checkpoint." State that each author stages at its own commit and that a repair's record call refreezes the chunk. The file is at 95 of 95 lines, and each paragraph is one line, so keep the paragraph count. Keep these sentences verbatim: line 21 (the charge), line 52 (the post-ticket preflight), and the line 68 sentence "Before a successor chunk starts, run `bench gate --checkpoint specs/<slug>/spec.md --chunk <id>`." New wording must avoid the Forbid needles of `internal/anchors/registry_chunk_chain.go:18-24`: "the last chunk refreezes", "review freezes", "When the orchestrator freezes a chunk after its last ticket", and "--chunk <id>` before the successor". "The last ticket's author freezes" stays true (`.bench/BENCH-reference.md:339`), because the last ticket's call still freezes before the review.

## Acceptance

- [ ] RC9: an earlier ticket's `bench record ticket` at its own commit writes its entries after the tip moved past it.
- [ ] RC10: a repair session's record call writes entries for its own ticket only and leaves the other entries unchanged.
- [ ] RC11: a non-last ticket's record call at a commit that descends from the chunk tip refreezes the chunk at that commit.
- [ ] RC12: a staged run at a record-only descendant of the chunk tip binds the tip commit and leaves the tip unchanged.
- [ ] RC16: a written entry carries `source_commit` equal to the staged run's HEAD.
- [ ] RC24: a staged run with a dirty path other than the record still refuses (`TestTicketRefusals`).
- [ ] `.agents/commands/bench-implement-spec.md` holds no restage sentence and stays at 95 lines or fewer.
- [ ] The probe "bind an earlier ticket to the chunk tip" bites in `bench test --package ./internal/reviewrecord/recordcmd`.
