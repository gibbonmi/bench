# Carry confirmed axis results at the checkpoint

Blocked by: 02-prove-test-only-tree-deltas.md
Writes: internal/reviewrecord/record.go, internal/reviewrecord/parse.go, internal/reviewrecord/check.go, internal/reviewrecord/coverage.go, internal/reviewrecord/completion.go, internal/reviewrecord/write.go, internal/reviewrecord/repair.go, internal/reviewrecord/fold.go, internal/reviewrecord/confirmation.go (new), internal/reviewrecord/confirmation_test.go (new), internal/reviewrecord/export_test.go, internal/reviewrecord/record_test.go, internal/reviewrecord/repair_test.go, internal/reviewrecord/source_test.go, internal/reviewrecord/fold_test.go, internal/reviewrecord/delegated_test.go, internal/reviewrecord/recordtest/fixture.go, internal/reviewrecord/recordtest/confirmation.go (new), internal/gate/checkpoint.go, internal/gate/review_checkpoint_confirmation_test.go (new), internal/gate/review_checkpoint_commits_test.go, internal/gate/review_checkpoint_test.go, internal/gate/complete_checkpoint_test.go, internal/gate/delegated_checkpoint_test.go
Covers: RC4, RC5, RC6, RC7, RC8, RC9, RC10, RC11, RC12, RC20, RC21, RC22, RC23, RC24, RC25, RC26, RC27, RC28, RC29, RC30, RC31, RC32, RC33, RC34, RC35, RC36, RC37, RC48, RC49, RC50, RC51, RC54

## What to build

The chunk entry gains a `confirmations` list with `omitempty`.
Each confirmation carries the shared evidence fields, the chunk base, the reviewed tip, and the repaired tip.
It also carries the closed finding IDs, a supersession link, and one probe.
Its role is `coordinator-confirmation`, and its source digest is the repaired digest.
The parser validates each confirmation and adds its ID to the one evidence-ID walk.
A record without confirmations keeps its current bytes.

The chunk checkpoint and the complete checkpoint read the latest confirmation of a chunk.
The confirmation is valid only when all of these conditions hold:

- Its base, repaired tip, and source digest equal the chunk pair.
- Its performer equals `completionReconciler`, and its state is completed.
- Its probe bites under the one probe-bite predicate.
- Its reviewed tip descends from the chunk base and is an ancestor of the chunk tip.
- `testgap.Prove` accepts the delta from the reviewed tip's tree digest to the chunk digest.

With a valid confirmation, the carry rule applies only when all three current axis results sit at the reviewed pair.
A carried result keeps the role, the performer-exclusion, and the `completed` state conditions, and it carries the derived reviewed digest.
Only the outcome and finding-ID condition relaxes.
The confirmation lists the union of the open finding IDs of the three carried results exactly.

A mixed state refuses through the review form before the strict per-axis refusals, with a reason that names the mixed pairs.
Each confirmation refusal names the chunk and the failed condition.
Fold reviews and the exported `CheckReviews` keep the strict rule.
Until ticket 04 adds the confirmation form, a confirmation-validity refusal routes through the review form.

One function in `internal/reviewrecord` finds the reviewed pair of a chunk from its three current axis results.
It refuses a mixed state, runs `testgap.Prove` from the reviewed tip's tree to the chunk digest, and returns the union of the open finding IDs.
The carry rule calls it, and ticket 04's writer and route call the same function.

Move the bite and restore conditions of the planned-probe check into one predicate that both checks call.
The planned-probe check keeps its mutation-name comparison and its two refusal texts.
Put the record-test confirmation builder in a new file, because the fixture file has 13 lines of headroom.
Each refusal row starts from the RC9 fixture, breaks exactly one condition, and asserts its reason.

This ticket owns the checkpoint behavior and builds each confirmation with the record-test builder.
The writer and the confirmation route come in the next ticket.

## Acceptance

- [ ] RC4–RC7: the parser accepts a valid confirmation and refuses a wrong role, a bad excerpt digest, and a duplicate ID.
- [ ] RC8: a record without confirmations renders with no `confirmations` key.
- [ ] RC9 and RC10: the chunk and complete checkpoints pass for a refrozen test-only repair. Standards is `fail` with F-1 at the reviewed pair, and the confirmation lists exactly F-1.
- [ ] RC11 and RC12: the checkpoint accepts an added `_test.go` file and refuses a delta that changes a non-test Go file.
- [ ] RC20–RC31 and RC48–RC51: each broken condition refuses with the reason that names it.
- [ ] RC32–RC34: the strict-evidence and comment-only cases keep their current outcomes.
- [ ] RC35 and RC36: the verification at the repaired digest and the distinct-session rule stay graded.
- [ ] RC37: the planned-probe check and the confirmation check call one probe-bite predicate.
- [ ] RC54: a mix of a fresh Coverage result and carried Standards and Spec results refuses through the review form.
- [ ] The carry rule gets the reviewed pair, the prover result, and the finding union from the one reviewed-pair function.
- [ ] The central-property mutation, which skips the prover in the carry branch, reds RC12, RC30, and RC31.
- [ ] A mutation that keeps the pass-and-no-findings condition on a carried result reds RC9 and RC10.
