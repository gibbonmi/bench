# Preserve review documents through the shared reader

Blocked by: 2-jsonfile-gate.md
Writes: internal/reviewrecord/parse.go, internal/reviewrecord/json_compatibility_test.go (new), internal/reviewrecord/testdata/strict-json-baseline.json (new), internal/reviewrecord/recordtest/fixture.go
Covers: SJ32, SJ33, SJ34, SJ35

## What to build

Replace the review decoder grammar with jsonfile.DecodeDocument. Remove uniqueJSON when this ticket migrates its last caller.
Keep bounds.ClassifyBytes before the shared decode. Preserve review record and completion-plan domain checks and caller error context.

Capture the accepted spec's review baseline family before the production edit. Include null and competing invalid inputs.
Compare Parse and ReadPlan with that permanent baseline. Preserve their observable diagnostic precedence.

Render a valid completion-plan fence with the existing fixture owner. Read its payload through ReadPlan without a final newline.
Then change a required domain field and confirm the same caller refusal as the baseline.

This ticket consumes the shared grammar and errors accepted at SJ-C2. It delivers both review consumers before its checkpoint.
Its successor preserves this review contract. Do not copy the fenced-payload or repository fixture harness.

## Acceptance

- [ ] Parse accepts a valid review document without a final newline (SJ32).
- [ ] ReadPlan accepts the real fenced completion-plan payload (SJ33).
- [ ] Review results match the committed baseline family (SJ34).
- [ ] Review byte limits retain their current refusal (SJ35).

## Checkpoint verification

- `bench test --package ./internal/reviewrecord`
- `bench test --package ./internal/reviewrecord/recordtest`

Record one red and one green for each covered row. Use a pre-edit red, a valid bench probe, or a recorded revert.
Keep each independent expectation only after its named behavioral mutation makes the assertion red. A compile failure proves no acceptance behavior.
The checkpoint commits only on the ticket lane after all focused checks pass. The orchestrator freezes the chunk for independent semantic review.

## Mutation proof

Replace the real review bytes with null at the shared decode call. The valid review-document assertion must fail.

`bench probe internal/reviewrecord/parse.go --swap 'jsonfile.DecodeDocument(data, value)' --with 'jsonfile.DecodeDocument([]byte("null"), value)' --package ./internal/reviewrecord --run TestReviewJSONCompatibility`

The command is a planned probe spelling. Confirm the exact production expression after the implementation exists.
If the planned expression differs, derive an equivalent compiling mutation at the same owned seam. Record its site and mutation kind.
Require bench probe to report bit and a byte-exact restore. A silent or invalid result requires a discriminating replacement probe.
