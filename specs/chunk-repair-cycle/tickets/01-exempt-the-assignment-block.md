# Exempt the assignment block from the test-only proof

Blocked by: none
Writes: internal/reviewrecord/confirmation.go, internal/reviewrecord/coverage.go, internal/reviewrecord/assignment.go, internal/reviewrecord/write.go, internal/reviewrecord/parse.go, internal/reviewrecord/plan.go, internal/reviewrecord/recordcmd/assignment.go, internal/reviewrecord/replacement_test.go, internal/reviewrecord/recordtest/confirmation.go, internal/reviewrecord/recordtest/fixture.go, internal/testgap/testgap.go, internal/testgap/testgap_test.go, internal/reviewrecord/recordcmd/confirmation_test.go, internal/gate/review_checkpoint_confirmation_test.go, docs/known-issues.md, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, .bench/BENCH-reference.md, projects/benchkit.md, .agents/skills/bench-craft-line/references/bounded-repair-policy.md, .agents/commands/bench-review-implementation.md, .bench/commitment.json
Covers: RC1, RC2, RC3, RC4, RC5, RC6, RC7

## What to build

Review chunk: RC-C1.

A test-only repair of a delegated build closes through `bench record confirmation` when the only other change in its delta is an assignment-only plan commit. The chunk checkpoint accepts the same delta, because both reach one rule through `reviewrecord.FindReviewedPair`.

Put the rule at the one call site of `testgap.Prove` (`internal/reviewrecord/confirmation.go:121`). The rule applies when the spec file is the only non-test path in the delta. It decodes the plan fence of the reviewed spec and of the repaired spec. It clears `execution.assignments` in both and re-encodes both through the payload, locate, and splice path that `RecordAssignment` uses (`internal/reviewrecord/assignment.go:65-77`, `internal/reviewrecord/write.go:45,57`, `internal/reviewrecord/parse.go:292`). Then it requires the two spec documents to be byte-equal. An absent assignment block and an empty assignment list clear to the same fence. Each other changed path goes through the test-gap proof unchanged. A ticket file stays a non-test path.

Add one exported `testgap` function that returns the proven test-file list, and make `Prove` compose it. RC-C6 reads that list, so keep its signature small and stable. Keep `Prove(root, reviewed, later)` as it is, because `TestProveTestGap` (`internal/testgap/testgap_test.go:79`) and the call site read it.

The fixture in `internal/reviewrecord/recordtest/fixture.go` is at 397 of 400 lines, and the commit lane grades file growth. Put a new helper for an assignment-only plan commit in `internal/reviewrecord/recordtest/confirmation.go`, or move bytes out of `fixture.go` in this ticket.

Remove the FT438 row from `docs/known-issues.md`, because this ticket makes its retire condition hold. Add one Forbid row on the needle `FT438` for `docs/known-issues.md` in `internal/anchors/registry_ticket_passes.go`, beside the FT405 row at line 15. `TestTicketSlicingAnchors` (`internal/anchors/registry_ticket_passes_test.go`) lists the needles a second time. Add the mirrored row there, and record the demonstrated red that the AGENTS.md exception for an independent expectation requires.

`.bench/BENCH-reference.md:549`, `projects/benchkit.md:44-45`, `.agents/commands/bench-review-implementation.md:47-49`, and `.agents/skills/bench-craft-line/references/bounded-repair-policy.md:56` state that a coordinator confirmation closes a test-only repair. Each sentence stays true. Edit one only if it states that any spec change takes a confirming round. The RC6 refusal text "the delta to the chunk tip is not test-only" stays, because `internal/reviewrecord/recordcmd/confirmation_test.go:105` and `internal/gate/review_checkpoint_confirmation_test.go:171-173` pin it.

`.bench/commitment.json` is in this fence because the spec-stage close commits the deliverable binding inside the range that the plan-only preflight grades. The author does not edit it.

## Acceptance

- [ ] RC1: a repair delta of one test file and one assignment-only plan commit records a confirmation through `bench record confirmation`.
- [ ] RC2: a repair delta that changes a plan-fence field outside the assignment block refuses the confirmation.
- [ ] RC3: a repair delta that widens a ticket `Writes:` line refuses the confirmation.
- [ ] RC4: a repair delta that changes a seam cell outside the plan fence refuses the confirmation.
- [ ] RC5: the chunk checkpoint accepts the assignment-only confirmation of RC1.
- [ ] RC6: a repair delta with a non-test Go file still refuses the confirmation (`TestRecordConfirmationRefusals`).
- [ ] RC7: `docs/known-issues.md` holds no row that names FT438, and a restored row reds the `docs-currency-workflow` anchors.
- [ ] The probe "clear the whole plan fence instead of the assignment block" bites in `bench test --package ./internal/reviewrecord/recordcmd`.
