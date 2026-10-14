# Repair confirmation

Status: staged

Roadmap: FT379

Decision source: Named reviewed artifact roadmap/FT379.md with its approved criteria in .bench/commitment.json. A reviewer-delegated top-tier consultation (fable / high) closed design forks F1–F9 on 2026-10-13.

Verification log: 5 iteration(s) to accept — the spec review took 3 iterations on opus / medium, and the ticket review took 2. The third spec iteration was an extension past the cap of 2, which a fable / high consultant granted. The spec round found an outcome refusal that hid the stale refusal and two derivations of the open finding set. The ticket round found gaps in the reviewed-pair contract, the confirmation route, and the guidance prose check.

## Problem

The bounded repair policy counts a correction that changes a disposition as a repair cycle.
That count applies even when the correction changes only the review record.
On 2026-10-05 the reviewer decided that a record-only change consumes no repair cycle.
The policy does not state that decision yet.

A test-only repair also costs a full confirming round today.
The chunk checkpoint requires three axis results at the current frozen pair.
A repair refreezes the chunk, so every earlier axis result goes stale.
In the FT217 build, three test-only repairs each had a biting coordinator probe.
Each repair still waited for a confirming round of all three axes.

## Solution

The policy states that a record-only correction consumes no repair cycle, even when it changes a disposition.
The issuing axis still reaffirms the changed disposition.

The review record gains a coordinator confirmation entry on a chunk.
The coordinator writes it with `bench record confirmation` from one complete `bench probe` answer.
The chunk checkpoint and the complete checkpoint accept the axis results at the reviewed pair when a valid confirmation covers the repaired pair.
A valid confirmation needs a test-only delta, a biting probe, the completion reconciler as performer, and the complete set of open findings.
Any other delta still takes the confirming round of all three axes.

## User stories

Line: opus / high.

Implementation-line reason: The hardest material chunk is RC-CARRY, which changes the oracle at the chunk checkpoint seam.
Its wrong answer weakens every later chunk review, so the leverage override applies.
The spec is precise: F1–F9 close the shape, and the comment-only gap is a direct precedent.
Seam uncertainty is low, because `checkSource` and `checkReviews` are the one owner of axis currency.
The gate tests are strong: the checkpoint tests drive real Git trees through `RunCommand`.
RC-CYCLE is a prose and anchor edit with exact twin pins.

Harder chunks: RC-CARRY, RC-WRITE.

### Record-only correction

1. As a coordinator, I want a record-only correction to consume no repair cycle, so that record bookkeeping cannot exhaust the repair allowance.
2. As a coordinator, I want that rule to hold when the correction changes a disposition, so that the reviewer's 2026-10-05 decision is the written policy.
3. As a reviewer, I want the issuing axis to reaffirm a changed disposition, so that each disposition change keeps an independent owner.
4. As a coordinator, I want the evidence-only definition to stay unchanged, so that a disposition change does not join the batched evidence-only class.
5. As a maintainer, I want an anchor and an independent twin to pin the new rule, so that a deleted rule turns the gate red.
6. As a reviewer, I want the delivered confirming-round scope to stay as it is, so that this spec does not restate a delivered rule.

### Confirmation entry

7. As a coordinator, I want a chunk to hold a coordinator confirmation, so that a test-only repair keeps a durable record of its biting probe.
8. As a coordinator, I want the confirmation to name its pairs, its closed findings, and its probe, so that the checkpoint can grade each part.
9. As a maintainer, I want the parser to refuse a malformed confirmation, so that a bad role, pair, or excerpt never reaches the checkpoint.
10. As a maintainer, I want confirmation IDs in the one evidence-ID walk, so that the parser refuses a duplicate ID across the evidence lists.
11. As a maintainer, I want a record without confirmations to keep its current bytes, so that every existing record stays valid.
12. As a linked-repository operator, I want a known consequence for an older binary, so that I upgrade before I write a confirmation.

### Checkpoint acceptance

13. As a coordinator, I want the chunk checkpoint to accept a valid confirmation, so that a test-only repair needs no confirming round.
14. As a coordinator, I want the complete checkpoint to accept the same confirmation, so that the last chunk gains the same relief.
15. As a reviewer, I want the delta to hold only added or modified Go test files, so that a production change still takes fresh axes.
16. As a reviewer, I want a path with a `testdata` element refused, so that test data that production code embeds never counts as test-only.
17. As a reviewer, I want a deleted file refused, so that a removed test never closes a finding.
18. As a reviewer, I want a type or mode change refused, so that a symbolic link or a gitlink never counts as a test.
19. As a reviewer, I want a non-Go path and a helper-package Go file refused, so that compiled helper code still takes fresh axes.
20. As a reviewer, I want a probe that did not bite refused, so that a silent mutation never closes a finding.
21. As a reviewer, I want an unrestored probe refused, so that a broken source never passes.
22. As a reviewer, I want a probe with a zero focused exit refused, so that a mutated run that passed never counts as a bite.
23. As a reviewer, I want a performer other than the completion reconciler refused, so that only the coordinator confirms.
24. As a reviewer, I want a finding list that differs from the open findings refused, so that no finding stays open behind a confirmation.
25. As a reviewer, I want a reviewed tip outside the chunk's history refused, so that a confirmation cannot borrow reviews of another line.
26. As a reviewer, I want an axis result whose digest differs from the reviewed tip's tree refused, so that a hand-edited digest cannot pass.
27. As a reviewer, I want a confirmation at a stale repaired pair refused, so that a later refreeze needs a fresh confirmation.
28. As a coordinator, I want a second test-only repair accepted with a fresh confirmation, so that cumulative test-only repairs need no confirming round.
29. As a reviewer, I want a later non-test repair to need three fresh axes, so that an old confirmation never covers production change.
30. As a reviewer, I want a plan path in the delta to fall back to three fresh axes, so that the prover stays narrow.
31. As a maintainer, I want the comment-only gap and the strict evidence refusals unchanged, so that the existing guarantees stay red-capable.
32. As a reviewer, I want each chunk verification still graded at the repaired digest, so that the confirmation waives only the review round.
33. As a reviewer, I want the review cardinality still graded on the carried axis results, so that the distinct-session rule survives.
34. As a maintainer, I want fold reviews to keep the strict rule, so that a fold resolution still takes independent review.
35. As a maintainer, I want one probe-bite predicate for verification and confirmation, so that the two graders cannot drift.

### Writer and route

36. As a coordinator, I want `bench record confirmation` to write the entry from one complete `bench probe` answer, so that I never write hand-built JSON.
37. As a coordinator, I want the writer to derive the reviewed pair from the current axis results, so that no flag restates a recorded fact.
38. As a coordinator, I want the writer to refuse axis results that sit at different pairs, so that one confirmation covers one reviewed pair.
39. As a coordinator, I want the writer to refuse a delta that is not test-only, so that a wrong confirmation fails before the checkpoint.
40. As a coordinator, I want the writer to refuse an incorrect finding set, so that I learn the complete set at write time.
41. As a coordinator, I want a new confirmation to supersede the prior one, so that the record keeps its history.
42. As a coordinator, I want the stale-axis refusal to name the confirmation form only for a test-only delta, so that the route names the repair.
43. As a coordinator, I want `bench record` help to list the confirmation form in the completion order, so that the write order stays discoverable.
44. As a reviewer, I want the probe's mutation text recorded as caller text, so that nobody reads it as verified evidence.

### Guidance

45. As a coordinator, I want the policy to state the test-only exception beside the current-results rule, so that the exception has one owner.
46. As a reviewer, I want the policy to name the confirmation's finding list as the coordinator's disposition, so that the no-op rule still governs it.
47. As a reviewer, I want the review phase to state the test-only exception to the confirming round, so that the phase obeys the checkpoint.
48. As a coordinator, I want the review phase's record step to name the confirmation form, so that the form list stays complete.
49. As a maintainer, I want the reference and the profile to name the confirmation, so that no reader states a stale three-axis universal.
50. As an adopter, I want a changelog entry, so that the release notes name the new form.
51. As a maintainer, I want the implementation phase file left unchanged, so that its exact prose budget holds.

### Carried-evidence conditions

52. As a reviewer, I want a carried axis result that never completed refused, so that only a completed review carries forward.
53. As a reviewer, I want a confirmation that never completed refused, so that a pending entry never closes a finding.
54. As a reviewer, I want a reviewed tip outside the chunk base's line refused, so that a whole chunk never closes with no review.
55. As a reviewer, I want a confirmation whose base differs from the chunk base refused, so that the confirmation binds to its own chunk.
56. As a reviewer, I want a mix of fresh and carried axis results refused, so that the writer and the checkpoint compute one finding set.

## Implementation decisions

### Record-only cycle (F8)

The policy keeps its evidence-only definition, its no-cycle sentence for evidence-only corrections, and its issuing-axis sentence unchanged.
One new sentence follows the no-cycle sentence.
It states that a correction that changes only the review record consumes no repair cycle, even when it changes a disposition.
A second new sentence states that the issuing axis reaffirms a changed disposition.
No Go code counts repair cycles or holds a disposition, so this rule takes no Go row.
`SourceDigest` already excludes the review record, so a record-only edit already leaves the checkpoint current.

The delivered build-directions rule already limits a confirming round to its named findings.
The reviewer dropped the criterion `FT379.record-only-round` as delivered on 2026-10-10.
The general optional-advice rule already makes a new preference advice.
This spec therefore adds no row for either clause.

### Confirmation shape (F1, F5, F7)

The chunk entry gains a `confirmations` list with `omitempty`.
Each confirmation carries these fields:

- the shared evidence fields
- the chunk base, the reviewed tip, and the repaired tip
- the closed finding IDs and the supersession link
- one probe

The evidence role is `coordinator-confirmation`, and the evidence source digest is the repaired digest.
The confirmation is not a review entry, so the independent-review role rule and the performer exclusion stay unchanged.

The coordinator refreezes the chunk at the repaired tip with `bench record chunk`, as today.
The chunk verification is then recorded again at the repaired digest, as today.
The confirmation waives the review round only, never the verification.

The checkpoint reads the latest confirmation of a chunk.
That confirmation is valid only when all of these conditions hold:

- Its base, repaired tip, and source digest equal the chunk pair.
- Its performer equals `completionReconciler`, and its state is completed.
- Its probe bites, under the one probe-bite predicate.
- Its reviewed tip descends from the chunk base and is an ancestor of the chunk tip.
- The test-only prover accepts the delta from the reviewed tip's digest to the chunk digest.

The checkpoint derives the reviewed digest from the reviewed tip's tree.
With a valid confirmation, the carry rule applies only when all three current axis results sit at the reviewed pair.
A result sits at the reviewed pair when its base is the chunk base and its tip is the reviewed tip.
The digest condition grades a carried result after that match.
Otherwise a mixed state refuses through the review form before the strict per-axis refusals, with a reason that names the mixed pairs.
A valid confirmation lets `checkReviews` apply the carry rule to an axis result at the reviewed pair.

A carried axis result keeps the role, the performer-exclusion, and the `completed` state conditions.
A carried axis result must carry the chunk base, the reviewed tip, and the derived reviewed digest.
Only the outcome and finding-ID condition relaxes, and a carried result may be `fail`.
The confirmation must list the union of the open finding IDs of the carried results exactly.
A carried result that passed at the reviewed pair contributes no IDs.

The distinct-session rule counts the three carried results.
Each confirmation refusal names the chunk and the failed condition.
With a valid confirmation, a carried result that fails the role, performer, state, pair, or digest condition refuses through the review form.

A later test-only repair refreezes the chunk again and takes a fresh confirmation.
The prover then grades the cumulative delta from the reviewed tip.
A later non-test repair takes three fresh axis results at the new pair, so the old confirmation no longer matches.
`CheckReviews`, the exported test entry, keeps the strict rule, because it has no tree to prove a delta.

### Test-only prover (F2)

A new package `internal/testgap` beside `internal/commentgap` proves a test-only tree delta.
It uses `TreeChangesIncludingSubmodules`, which runs with `--no-renames`, so a rename arrives as a deletion and an addition.
It accepts a change only when every condition holds:

- The status is `A` or `M`.
- The destination mode is a regular file mode, and an `M` change keeps its source mode.
- `outline.IsGoTestFile` accepts the path.
- No path element is `testdata`.

It refuses an empty delta between different trees, as the comment-only prover does.
Each refusal returns one named sentinel error:

- `ErrUnreadable` for an unreadable tree
- `ErrEmptyChanges` for an empty delta between different trees
- `ErrStatus` for a status other than `A` or `M`
- `ErrMode` for a changed or non-regular mode
- `ErrNotTest` for a path that `outline.IsGoTestFile` refuses
- `ErrTestdata` for a path with a `testdata` element

The review record is outside every source digest, so the prover never sees it.
A plan commit under `specs/` is not a Go test file, so that delta falls back to three axes.

### Probe evidence (F4)

The confirmation reuses the `--probe-result` import of the verification form.
The probe-bite predicate requires the outcome `bit`, a non-zero focused exit, and the restore `pass`.
One function owns that predicate, and the planned-probe check and the confirmation check both call it.
The planned-probe check keeps its separate mutation-name comparison.
The planned-probe check also keeps its two refusal texts, one for the bite and one for the restore.

| site | old rule | new rule |
|---|---|---|
| planned-probe check | The mutation name, the `bit` outcome, and the non-zero exit give the bite refusal. The restore gives the restore refusal. | The mutation-name comparison stays at the site. The shared predicate supplies the bite and the restore conditions, and both refusal texts stay. |
| confirmation check | none | The shared predicate grades the probe, and the refusal names the probe condition that failed. |

The confirmation records the caller's `--mutation` text without verification.
The probe reader does not prove the origin of the answer bytes, so the confirmation has the trust level of a planned probe.

### Writer and route (F6)

`bench record confirmation <slug> --chunk <id> --id <id> --performer <session> --model <model> --effort <effort> --exit-code <n> --ref <ref> --probe-result <file> [--mutant <n>] --mutation <text> [--finding <id>]...` writes one confirmation.
The writer derives the reviewed pair from the current result of each axis.
It refuses when the three results sit at different pairs, or when they already sit at the chunk pair.
It refuses a delta that the prover refuses and a finding set that differs from the open set.
It supersedes the prior confirmation of the chunk.

`FormConfirmation` joins the form constants, and the completion order places it after the review form.
The record command registers the form.

A chunk axis with no valid confirmation can take either route.
With no valid confirmation, an axis result can sit at an earlier pair: the chunk base and an ancestor of the chunk tip.
If all three current results sit at that pair and the prover accepts the delta from that tip, the refusal names the confirmation form.
That refusal comes before the outcome refusal, so a `fail` axis at the reviewed pair also routes to the confirmation form.
Otherwise the refusal names the review form, as today.

### Compatibility

The record keeps versions 1 and 2, and the new list takes no version bump.
A record without a confirmation keeps its current bytes.
The decoder refuses unknown fields, so an older binary refuses a record that holds a confirmation.
The landing gate builds the binary from the graded tree, so the kit's own checkpoint always reads the field.

### Guidance

The policy gains the test-only exception beside the current-results rule.
It names the confirmation's finding list as the coordinator's disposition of those findings.
The review phase keeps its confirming-round sentence and adds the test-only exception after it.
The review phase's record step names the new form.
The reference's checkpoint sentence and the project profile's completion-checkpoint seam name the confirmation.
The implementation phase file stays unchanged: its pointer to the narrow evidence-only exception stays correct, and it sits at its exact prose budget.

### Route flag

The roadmap row says `Next: kit-edit`, and the earlier handoff routed FT379 to a one-ticket folder.
RC-CONFIRM changes `internal/reviewrecord`, a new prover, the record command, and the gate checkpoint seam.
That work fails the light-path test of one ticket with no declared seam.
This spec follows the invoked `/bench-write-spec` and flags the route for reviewer veto.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| RC-CYCLE / 01-state-record-only-repair-cycle-rule.md | The policy states that a record-only correction consumes no repair cycle, with its twin pins. | RC1, RC2, RC3 | `TestImplementationContinuation` | no |
| RC-PROVE / 02-prove-test-only-tree-deltas.md | A new prover accepts only an added or modified regular Go test file outside `testdata`, and both prover suites share one tree fixture builder. | RC13, RC14, RC15, RC16, RC17, RC18, RC19, RC52, RC53 | the planned prover tests and `TestProveTreeGap` | no |
| RC-CARRY / 03-carry-confirmed-axis-results.md | The record holds coordinator confirmations, and the chunk and complete checkpoints carry the reviewed axis results under a valid confirmation. | RC4–RC12, RC20–RC37, RC48–RC51, RC54 | `TestReviewCheckpointKeepsStrictEvidence`, `TestReviewCheckpointCommentOnlyGap`, and the planned record and checkpoint tests | yes |
| RC-WRITE / 04-write-confirmations-with-bench-record.md, 05-document-the-test-only-confirmation.md | `bench record confirmation` writes the entry, the refusal routes to it, and the guidance states the exception. | RC38–RC47 | `TestCheckpointEntryRepairForms`, `TestCompletionHelpOrdering`, `TestImplementationContinuation`, `TestChunkChainAnchors`, and the planned writer tests | yes |

The slicing split the reviewed chunk RC-CONFIRM into RC-PROVE, RC-CARRY, and RC-WRITE.
RC-PROVE and RC-CARRY each create a seam that a later ticket consumes, so each is its own review chunk.
RC-CYCLE and RC-WRITE touch disjoint sentences of the policy file, and RC-WRITE comes last.

The central-property mutation of RC-CARRY makes the carry branch skip the test-only prover.
RC12, RC30, and RC31 turn red under that mutation.
A second mutation keeps the pass-and-no-findings condition on a carried result, and RC9 and RC10 turn red under it.

### Completion plan

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"ft379-build-20261014","orchestrator_session":"claude:ft379-orchestrator-20261014","author_limit":1,"assignments":{"01-state-record-only-repair-cycle-rule.md":[{"session":"claude:ft379_t01","assignment":"d835a91fcc200941c50e06f300cb4839","model":"opus","effort":"high","source":"db03cfdff60a7932788569151d2a598e361afd08","native_ref":"claude-agent:ft379_t01"},{"session":"claude:ft379_t01_rep1","assignment":"d835a91fcc200941c50e06f300cb4839","model":"opus","effort":"high","source":"63b6661e49b589f5e5e091b8fde55b6b4dfbd842","native_ref":"claude-agent:ft379_t01_rep1","predecessor":"claude:ft379_t01","trigger":"user-directed","stopped":"claude-agent:ft379_t01","preserved":"6dc1bb38ec18b0cdaecefb060de116b51b4a62c7"}],"02-prove-test-only-tree-deltas.md":[],"03-carry-confirmed-axis-results.md":[],"04-write-confirmations-with-bench-record.md":[],"05-document-the-test-only-confirmation.md":[]}},"chunks":[{"id":"RC-CYCLE","tickets":["01-state-record-only-repair-cycle-rule.md"],"verification":[{"id":"policy-anchors","command":"bench test --package ./internal/conformance --package ./internal/anchors","ticket":"01-state-record-only-repair-cycle-rule.md"},{"id":"policy-prose","command":"bench test --check prose-mechanics","ticket":"01-state-record-only-repair-cycle-rule.md"},{"id":"command-registries","command":"bench test --package ./cmd/bench","ticket":"01-state-record-only-repair-cycle-rule.md"},{"id":"cycle-root-conformance","command":"bench test --check root-conformance","ticket":"01-state-record-only-repair-cycle-rule.md"}]},{"id":"RC-PROVE","tickets":["02-prove-test-only-tree-deltas.md"],"verification":[{"id":"prover","command":"bench test --package ./internal/testgap --package ./internal/commentgap --package ./internal/gittest","ticket":"02-prove-test-only-tree-deltas.md"},{"id":"prover-status","command":"bench test --package ./internal/testgap","probe":"accept status D in the prover","ticket":"02-prove-test-only-tree-deltas.md"},{"id":"prove-root-conformance","command":"bench test --check root-conformance","ticket":"02-prove-test-only-tree-deltas.md"}]},{"id":"RC-CARRY","tickets":["03-carry-confirmed-axis-results.md"],"verification":[{"id":"record-and-checkpoint","command":"bench test --package ./internal/reviewrecord/... --package ./internal/gate","ticket":"03-carry-confirmed-axis-results.md"},{"id":"carry-prover","command":"bench test --package ./internal/gate","probe":"skip the test-only prover in the carry branch","ticket":"03-carry-confirmed-axis-results.md"},{"id":"carry-root-conformance","command":"bench test --check root-conformance","ticket":"03-carry-confirmed-axis-results.md"}]},{"id":"RC-WRITE","tickets":["04-write-confirmations-with-bench-record.md","05-document-the-test-only-confirmation.md"],"verification":[{"id":"writer-and-route","command":"bench test --package ./internal/reviewrecord/... --package ./internal/gate --package ./cmd/bench --package ./internal/conformance","ticket":"04-write-confirmations-with-bench-record.md"},{"id":"guidance-anchors","command":"bench test --package ./internal/conformance --package ./internal/anchors --package ./cmd/bench","ticket":"05-document-the-test-only-confirmation.md"},{"id":"route-precedence","command":"bench test --package ./internal/gate","probe":"let the outcome refusal fire before the confirmation route","ticket":"04-write-confirmations-with-bench-record.md"},{"id":"guidance-prose","command":"bench test --check prose-mechanics","ticket":"05-document-the-test-only-confirmation.md"},{"id":"write-root-conformance","command":"bench test --check root-conformance","ticket":"04-write-confirmations-with-bench-record.md"},{"id":"guidance-root-conformance","command":"bench test --check root-conformance","ticket":"05-document-the-test-only-confirmation.md"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/repair-confirmation/spec.md"},{"id":"record-checkpoint-writer","command":"bench test --package ./internal/reviewrecord/... --package ./internal/gate --package ./cmd/bench"},{"id":"prover-and-guidance","command":"bench test --package ./internal/testgap --package ./internal/commentgap --package ./internal/conformance --package ./internal/anchors"}]}
```

## Testing decisions

- A good test drives a real Git fixture through `RunCommand` with `--checkpoint`, and it reads the exit code and the refusal text.
- The tree fixture builder moves from the comment-only prover tests into `internal/gittest` as an exported helper. The comment-only and the test-only prover tests both call it, so no second fixture harness appears.
- The prover tests build two trees with that helper and assert each sentinel, as the comment-only prover tests do.
- The record-command tests import a real `bench probe` answer through the existing `answer` helper.
- Each checkpoint refusal row starts from the RC9 fixture and breaks exactly one condition. The row asserts the refusal reason that names that condition and the route that the decision gives it. A fixture that refuses for another reason, such as a stale verification, therefore fails the row.
- The guidance rows use the anchor registry and its independent conformance twin.
- The prior art is `TestReviewCheckpointCommentOnlyGap` for acceptance and `TestReviewCheckpointKeepsStrictEvidence` for refusals.
- The gate observes the feature through `bench gate --checkpoint` and the `test` phase.

### Seam diagram

    trigger: the coordinator after a test-only repair
        │
        ▼
    bench record chunk (refreeze) ──▶ bench record verification (repaired digest)
        │
        ▼
    bench record confirmation ──▶ [ reviewrecord: Parse, RecordConfirmation ]
        │
        ▼
    bench gate --checkpoint ──▶ [ checkSource → confirmation validity → test-only prover → checkReviews carry ]
                                    ◀ tests attach here: RunCommand on a fixture repository, exit code and refusal text

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| RC1 | 1, 2 | The policy states that a correction that changes only the review record consumes no repair cycle, even when it changes a disposition. | `internal/conformance/implementation_continuation_test.go` (`TestImplementationContinuation`) | The registry anchor and its independent twin red when the sentence is absent. |
| RC2 | 3 | The policy states that the issuing axis reaffirms a changed disposition. | `internal/conformance/implementation_continuation_test.go` (`TestImplementationContinuation`) | The registry anchor and its independent twin red when the sentence is absent. |
| RC3 | 4, 5 | The evidence-only definition keeps its exact needle, including the word disposition. | `internal/conformance/implementation_continuation_test.go` (`TestImplementationContinuation`) | The existing exact needle reds a widened definition. |
| RC4 | 7, 8 | The parser accepts a chunk confirmation with role `coordinator-confirmation`, a valid pair, finding IDs, and a valid probe excerpt. | planned | A parser that refuses the new list reds the round trip. |
| RC5 | 9 | The parser refuses a confirmation whose role is not `coordinator-confirmation`. | planned | A parser that skips the role check accepts the entry and the row reds. |
| RC6 | 9 | The parser refuses a confirmation whose native excerpt digest differs from its excerpt. | planned | A parser that skips native validation accepts the entry and the row reds. |
| RC7 | 10 | The parser refuses a confirmation ID that equals a review ID. | planned | An evidence-ID walk that omits confirmations accepts the duplicate and the row reds. |
| RC8 | 11 | A record without confirmations renders with no `confirmations` key. | planned | A field without `omitempty` adds the key and the byte comparison reds. |
| RC9 | 13 | The chunk checkpoint passes when a valid confirmation lists exactly F-1 for a refrozen test-only repair, with Standards `fail` with finding F-1 and Spec and Coverage `pass` at the reviewed pair. | planned | A carry branch that keeps the pass-and-no-findings condition refuses the carried `fail` result and the row reds. |
| RC10 | 14 | The complete checkpoint passes for the RC9 chunk with the same confirmation. | planned | A carry branch wired only into the chunk path reds the complete run. |
| RC11 | 15 | The checkpoint accepts a delta that adds a new `_test.go` file. | planned | A prover that accepts only status `M` refuses the addition and the row reds. |
| RC12 | 15, 19 | The checkpoint refuses a confirmation whose delta changes a non-test Go file. | planned | A carry branch without the prover accepts the delta and the row reds. |
| RC13 | 16 | The prover refuses a `_test.go` path with a `testdata` element with `ErrTestdata`. | planned | A prover that tests only the suffix accepts the path and the row reds. |
| RC14 | 17 | The prover refuses a deleted `_test.go` file with `ErrStatus`. | planned | A prover that accepts status `D` passes the deletion and the row reds. |
| RC15 | 18 | The prover refuses a `_test.go` path whose destination mode is a symbolic link with `ErrMode`. | planned | A prover that skips the mode check accepts the link and the row reds. |
| RC16 | 18 | The prover refuses a modified `_test.go` file whose mode changes with `ErrMode`. | planned | A prover that compares only the destination mode accepts the change and the row reds. |
| RC17 | 19 | The prover refuses a non-Go path with `ErrNotTest`. | planned | A prover that accepts any path under a test directory passes the file and the row reds. |
| RC18 | 19 | The prover refuses an empty change list between different trees with `ErrEmptyChanges`. | planned | A prover that returns nil on no changes passes the case and the row reds. |
| RC19 | 15 | The prover accepts a test path that contains a space. | planned | A prover that splits paths on spaces refuses the path and the row reds. |
| RC20 | 20 | The checkpoint refuses a confirmation whose probe outcome is `silent`. | planned | A carry branch that skips the bite predicate accepts the entry and the row reds. |
| RC21 | 21 | The checkpoint refuses a confirmation whose probe restore is `fail`. | planned | A bite predicate without the restore check accepts the entry and the row reds. |
| RC22 | 22 | The checkpoint refuses a confirmation whose probe focused exit is zero. | planned | A bite predicate without the exit check accepts the entry and the row reds. |
| RC23 | 23 | The checkpoint refuses a confirmation whose performer is not the completion reconciler. | planned | A carry branch without the performer check accepts an author's entry and the row reds. |
| RC24 | 24 | The checkpoint refuses a confirmation that omits one open finding ID. | planned | A subset comparison accepts the entry and the row reds. |
| RC25 | 24 | The checkpoint refuses a confirmation that lists a finding ID that no carried axis result holds. | planned | A superset comparison accepts the entry and the row reds. |
| RC26 | 25 | The checkpoint refuses a confirmation whose reviewed tip is not an ancestor of the chunk tip. | planned | A carry branch without the ancestry check accepts the entry and the row reds. |
| RC27 | 26 | The checkpoint refuses a carried axis result whose source digest differs from the reviewed tip's tree digest, and routes to `bench record 'review' --help`. | planned | A carry branch that trusts the recorded digest accepts the edit and the row reds. |
| RC28 | 27 | The checkpoint refuses a confirmation whose repaired tip differs from a later refrozen chunk tip. | planned | A carry branch that skips the pair match accepts a stale confirmation and the row reds. |
| RC29 | 28 | The checkpoint passes after a second test-only repair with a fresh confirmation that supersedes the first. | planned | A carry branch that reads the first confirmation refuses the new pair and the row reds. |
| RC30 | 29 | The checkpoint refuses a stale axis after a later non-test repair, even with an earlier confirmation. | planned | A carry branch that ignores the delta class accepts the production change and the row reds. |
| RC31 | 30 | The checkpoint refuses a confirmation whose delta adds a file under `specs/`. | planned | A prover widened to plan paths accepts the delta and the row reds. |
| RC32 | 31 | The strict-evidence statement case still refuses with stale reviewed source. | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointKeepsStrictEvidence`) | A carry branch that widens the uncovered-delta gap passes the statement case and the row reds. |
| RC33 | 31 | The strict-evidence moved-chunk case still refuses with stale Standards when no confirmation exists. | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointKeepsStrictEvidence`) | A carry branch that needs no confirmation passes the moved chunk and the row reds. |
| RC34 | 31 | The comment-only gap still passes its chunk, complete, and reaffirm cases. | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointCommentOnlyGap`) | A change that breaks the comment-only path reds the existing cases. |
| RC35 | 32 | The checkpoint refuses a confirmed chunk whose verification sits at the reviewed digest. | planned | A confirmation that also waives verification passes the stale verification and the row reds. |
| RC36 | 33 | The checkpoint refuses a confirmed chunk whose carried axis results share one session in a three-reviewer plan. | planned | A carry branch that skips the cardinality count passes and the row reds. |
| RC37 | 35 | The planned-probe check and the confirmation check call one probe-bite predicate. | review-owned | Review compares the two call sites, because a duplicated predicate passes every behavior test. |
| RC38 | 36 | `bench record confirmation` writes a confirmation from a complete `bench probe` answer with the derived pair and digests. | planned | A writer that records the chunk pair as the reviewed pair, or a digest other than the reviewed tip's tree digest, reds the field comparison. |
| RC39 | 37, 38 | The writer refuses when the current axis results sit at different pairs. | planned | A writer that reads only one axis writes the entry and the row reds. |
| RC40 | 39 | The writer refuses a delta that the prover refuses. | planned | A writer without the prover writes the entry and the row reds. |
| RC41 | 40 | The writer refuses a finding set that differs from the open set. | planned | A writer without the set comparison writes the entry and the row reds. |
| RC42 | 41 | A second confirmation names the first in its supersession link. | planned | A writer that omits the link writes an empty link, and the row's assertion of the first confirmation ID reds. |
| RC43 | 42 | With Standards `fail` and Spec and Coverage `pass` at the reviewed pair and a test-only delta, the checkpoint refusal routes to `bench record 'confirmation' --help`. | `internal/gate/record_repair_test.go` (`TestCheckpointEntryRepairForms`) | A route that lets the outcome refusal fire first names the review form and the new case reds. |
| RC44 | 42 | With Standards at the reviewed pair and a delta that changes a non-test Go file, the refusal routes to `bench record 'review' --help`. | `internal/gate/record_repair_test.go` (`TestCheckpointEntryRepairForms`) | A route that names the confirmation form for every earlier-pair axis reds the new non-test case. |
| RC45 | 43 | `bench record` help lists the confirmation form and the completion order with confirmation after review. | `internal/reviewrecord/recordcmd/ordering_test.go` (`TestCompletionHelpOrdering`) | The independent order list reds when the form or its position is absent. |
| RC46 | 45, 46 | The policy states the test-only exception and names the confirmation's finding list as the coordinator's disposition. | `internal/conformance/implementation_continuation_test.go` (`TestImplementationContinuation`) | The new anchors and their independent twins red when either sentence is absent. |
| RC47 | 47 | The review phase states that a test-only repair whose coordinator probe bites closes through `bench record confirmation` and takes no confirming round. | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | The new anchor and its independent twin red when the sentence is absent. |
| RC48 | 52 | The checkpoint refuses a carried Standards result whose state is `failed`, with a reason that names the Standards state, and routes to `bench record 'review' --help`. | planned | A carry branch that drops the `completed` state condition accepts the result and the row reds. |
| RC49 | 53 | The checkpoint refuses a confirmation whose state is `pending`, with a reason that names the confirmation state. | planned | A validity check without the state condition accepts the entry and the row reds. |
| RC50 | 54 | The checkpoint refuses a confirmation whose reviewed tip does not descend from the chunk base, with a reason that names the reviewed tip. | planned | A one-sided ancestry check against the chunk tip accepts the entry and the row reds. |
| RC51 | 55 | The checkpoint refuses a confirmation whose base differs from the chunk base, with a reason that names the confirmation pair. | planned | A pair match that compares only the repaired tip accepts the entry and the row reds. |
| RC52 | 18 | The prover refuses a `_test.go` gitlink with mode `160000` with `ErrMode`. | planned | A prover that grades only symbolic links accepts the gitlink and the row reds. |
| RC53 | 17 | The prover refuses a renamed `_test.go` file with `ErrStatus`, because the rename arrives as a deletion. | planned | A prover that drops `--no-renames` or accepts status `R` passes the rename and the row reds. |
| RC54 | 56 | The checkpoint refuses when Coverage has a fresh pass at the chunk pair and Standards and Spec sit at the reviewed pair under a valid confirmation, with a reason that names the mixed pairs, and routes to `bench record 'review' --help`. | planned | A carry rule that grades each axis alone accepts the mixed state and the row reds. |

Not covered: story 6 — the build-directions outcome delivered the confirming-round scope, and its pins stay unchanged.
Not covered: story 12 — the compatibility consequence is a stated limit with no new behavior; the Won't handle line holds it.
Not covered: story 34 — fold reviews keep the current strict rule, and the existing fold tests stay unchanged.
Not covered: story 44 — the mutation text is caller text by design; the Won't handle line holds it.

Not covered: story 48 — the record-step sentence is review-owned prose with no anchor family.
Not covered: story 49 — the reference and the profile sentences are review-owned prose with no anchor family.
Not covered: story 50 — the changelog entry is review-owned.
Not covered: story 51 — the implementation phase file stays unchanged, and its prose budget check grades it.

### Edge inventory

The audience is every repository that links the kit, because the record and the checkpoint ship in the binary.

- Paths with spaces: RC19 covers a test path with a space. The tree-change helper uses NUL framing, so no quoting applies.
- Control bytes: the existing native-excerpt validation refuses a newline in the reference, and RC6 covers the excerpt digest.
- Git C-quoted patch headers: the prover reads `--raw -z` output, so no quoted header reaches it.
- Absent versus empty list: RC8 covers the absent key; an empty list parses as no confirmation.
- Symbolic links and gitlinks: RC15 covers a link, and RC52 covers a gitlink.
- A command whose write changes a fact it reports: the writer writes only the record file, and the record is outside every source digest.
- State reloaded by a fresh process: every checkpoint row saves the record and runs `RunCommand` on the saved bytes.
- Deleted and renamed files: RC14 covers a deletion, and RC53 covers a rename.
- A hand-edited record: RC27 derives the reviewed digest from the tree, so a hand-edited digest cannot pass.

**Won't handle** an invented probe answer — the probe reader cannot prove origin, and the planned-probe check keeps the same trust level today.
**Won't handle** a check that the probe subject is a non-test path — no criterion promises it, and the review phase prose states it as guidance.
**Won't handle** an older binary that reads a confirmation — the decoder refuses unknown fields, and the landing gate builds from the graded tree.
**Won't handle** a test-only repair of a fold resolution — fold reviews keep the strict rule, and the chunk confirmation covers the FT217 occurrences.
**Won't handle** test data and helper packages as test-only — production code can embed `testdata`, and a mixed delta still closes through the confirming round.

## Ownership fences

Each entry is an exact repo-relative file or path prefix that a writer of this build may edit.
The list is the union of the ticket `Writes:` lines plus the review pickup.

- `.agents/commands/bench-review-implementation.md`
- `.agents/skills/bench-craft-line/references/bounded-repair-policy.md`
- `.bench/BENCH-reference.md`
- `.bench/commitment.json`
- `CHANGELOG.md`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/anchors/registry_calibration.go`
- `internal/anchors/registry_calibration_test.go`
- `internal/anchors/registry_chunk_chain.go`
- `internal/anchors/registry_chunk_chain_test.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_debug_loop.go`
- `internal/anchors/registry_ft311_review_dispatch.go`
- `internal/anchors/registry_retained_workflow.go`
- `internal/commentgap/commentgap_test.go`
- `internal/commentgap/trees_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/implementation_continuation_test.go`
- `internal/conformance/retained_workflow_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/gate/checkpoint.go`
- `internal/gate/complete_checkpoint_test.go`
- `internal/gate/delegated_checkpoint_test.go`
- `internal/gate/record_repair_test.go`
- `internal/gate/review_checkpoint_commits_test.go`
- `internal/gate/review_checkpoint_confirmation_test.go`
- `internal/gate/review_checkpoint_test.go`
- `internal/gittest/gittest.go`
- `internal/gittest/gittest_test.go`
- `internal/gittest/trees.go`
- `internal/reviewrecord/check.go`
- `internal/reviewrecord/completion.go`
- `internal/reviewrecord/confirmation.go`
- `internal/reviewrecord/confirmation_test.go`
- `internal/reviewrecord/coverage.go`
- `internal/reviewrecord/delegated_test.go`
- `internal/reviewrecord/export_test.go`
- `internal/reviewrecord/fold.go`
- `internal/reviewrecord/fold_test.go`
- `internal/reviewrecord/parse.go`
- `internal/reviewrecord/record.go`
- `internal/reviewrecord/record_test.go`
- `internal/reviewrecord/recordcmd/amendment_test.go`
- `internal/reviewrecord/recordcmd/command.go`
- `internal/reviewrecord/recordcmd/confirmation.go`
- `internal/reviewrecord/recordcmd/confirmation_test.go`
- `internal/reviewrecord/recordcmd/ordering_test.go`
- `internal/reviewrecord/recordcmd/refusal.go`
- `internal/reviewrecord/recordcmd/review_test.go`
- `internal/reviewrecord/recordcmd/verification.go`
- `internal/reviewrecord/recordcmd/verification_import_test.go`
- `internal/reviewrecord/recordtest/confirmation.go`
- `internal/reviewrecord/recordtest/fixture.go`
- `internal/reviewrecord/repair.go`
- `internal/reviewrecord/repair_test.go`
- `internal/reviewrecord/source_test.go`
- `internal/reviewrecord/write.go`
- `internal/testgap/testgap.go`
- `internal/testgap/testgap_test.go`
- `projects/benchkit.md`
- `reviews/repair-confirmation.md`
- `specs/repair-confirmation/spec.md`
- `tests/canary/docs-currency-token-diet/benchref-imported`
- `tests/canary/docs-currency-token-diet/benchref-pointer-dropped`
- `tests/canary/docs-currency-token-diet/benchref-section-duplicated`
- `tests/canary/guidance-prose-budgets/over-budget-skill`
- `tests/canary/line-routing/line-binding-prose-drift`
- `tests/canary/skill-description-budgets/budget-table-missing`
- `tests/canary/skill-description-budgets/description-folded`
- `tests/canary/skill-description-budgets/description-missing`
- `tests/canary/skill-description-budgets/over-budget-command`
- `tests/canary/skill-description-budgets/over-budget-description`
- `tests/canary/skills-index-command-adapters/adapter-inert-invocation-key`
- `tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy`
- `tests/canary/skills-index-command-adapters/dangling-index`
- `tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted`
- `tests/canary/skills-index-command-adapters/missing-index-field`
- `tests/canary/skills-index-command-adapters/stale-index-wording`
- `tests/canary/skills-index-command-adapters/unindexed-skill`
- `tests/canary/workflow-guidance-anchors/agents-handoff-section-rule`
- `tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-owner`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-routing`
- `tests/canary/workflow-guidance-anchors/benchkit-spec-ownership`
- `tests/canary/workflow-guidance-anchors/benchkit-system-suite-route`
- `tests/canary/workflow-guidance-anchors/calibration-pickup-confidence`
- `tests/canary/workflow-guidance-anchors/calibration-pickup-step-move`
- `tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns`
- `tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary`
- `tests/canary/workflow-guidance-anchors/coverage-axis-anchor`
- `tests/canary/workflow-guidance-anchors/delegated-axis-exclusions`
- `tests/canary/workflow-guidance-anchors/prepared-review-axis-returns`
- `tests/canary/workflow-guidance-anchors/prepared-review-blast-evidence`
- `tests/canary/workflow-guidance-anchors/prepared-review-capable-handoff`
- `tests/canary/workflow-guidance-anchors/prepared-review-inline-axis-route`
- `tests/canary/workflow-guidance-anchors/prepared-review-legacy-entry-points`
- `tests/canary/workflow-guidance-anchors/prepared-review-native-dispatch`
- `tests/canary/workflow-guidance-anchors/prepared-review-runtime-capability`
- `tests/canary/workflow-guidance-anchors/prepared-review-shared-evidence`
- `tests/canary/workflow-guidance-anchors/reference-agent-push-rule`
- `tests/canary/workflow-guidance-anchors/reference-bench-operational-layer`
- `tests/canary/workflow-guidance-anchors/reference-category-context`
- `tests/canary/workflow-guidance-anchors/reference-category-oracle`
- `tests/canary/workflow-guidance-anchors/reference-category-setup`
- `tests/canary/workflow-guidance-anchors/reference-category-work`
- `tests/canary/workflow-guidance-anchors/reference-gate-authority`
- `tests/canary/workflow-guidance-anchors/reference-kit-only-ship`
- `tests/canary/workflow-guidance-anchors/reference-no-path-fallback`
- `tests/canary/workflow-guidance-anchors/reference-progressive-loading-term`
- `tests/canary/workflow-guidance-anchors/reference-refusal-route-shape`
- `tests/canary/workflow-guidance-anchors/reference-retro-capture-owner`
- `tests/canary/workflow-guidance-anchors/reference-retro-drain-owner`
- `tests/canary/workflow-guidance-anchors/reference-skills-guidance`
- `tests/canary/workflow-guidance-anchors/reference-upgrade-route`
- `tests/canary/workflow-guidance-anchors/review-base-merged-main-tip`
- `tests/canary/workflow-guidance-anchors/review-clean-terminal-result`
- `tests/canary/workflow-guidance-anchors/review-cross-harness-opt-in`
- `tests/canary/workflow-guidance-anchors/review-falsification-accept-routing`
- `tests/canary/workflow-guidance-anchors/review-falsification-dispositions`
- `tests/canary/workflow-guidance-anchors/review-persistence-anchor`
- `tests/canary/workflow-guidance-anchors/review-preflight-explicit-base`
- `tests/canary/workflow-guidance-anchors/review-repair-ticket-covers`
- `tests/canary/workflow-guidance-anchors/review-repair-ticket-owner`
- `tests/canary/workflow-guidance-anchors/review-standing-falsification`
- `tests/canary/workflow-guidance-anchors/review-universal-claim-bar`

## Out of scope

- Authenticate the origin of a `bench probe` answer: 12 edits, 3 gate runs.
- Grade the probe subject as a non-test path: 6 edits, 2 gate runs.
- Accept `testdata` that no production package embeds: 8 edits, 2 gate runs.
- Accept a test-only repair of a fold resolution: 7 edits, 2 gate runs.

## Further notes

### Flagged additions

- The parser refuses a confirmation with a wrong role or an invalid excerpt (RC5, RC6). The criterion needs a well-formed form.
- The checkpoint derives the reviewed digest from the reviewed tip's tree (RC27). A hand-editable record otherwise passes a forged digest.
- The checkpoint requires an ancestry chain for the reviewed tip (RC26). This closes the borrowed-review hole.
- The writer refuses before the checkpoint does (RC39, RC40, RC41). It composes the same prover and set comparison, so no second fact source appears.
- One probe-bite predicate serves both graders (RC37). The one-source-per-fact rule requires it.
- The policy names the confirmation's finding list as the coordinator's disposition (RC46). The consultant raised this risk beside the no-op rule.
- The issuing axis reaffirms a changed disposition (RC2). Fork F8 closed it, so that each disposition change keeps an independent owner.
- The checkpoint compares the finding set exactly (RC24, RC25). Forks F1 and F3 closed it, because a subset leaves an unresolved finding.
- A new confirmation supersedes the prior one (RC42). Fork F7 closed it, so that the record keeps its history.
- The stale-axis refusal routes to the confirmation form (RC43, RC44). Fork F6 closed it, so that the route names the repair.
- A carried result keeps the `completed` state, the base, and the ancestry conditions (RC48–RC51). The spec review asked for them, because the cheapest carry otherwise accepts an incomplete review.
- A mixed state of fresh and carried axes refuses (RC54). The writer and the checkpoint compute one finding set over the three results at one reviewed pair. The delegated consultation closed this on 2026-10-13 after the spec review's G1 finding.

### Source sentence to row

| source sentence | rows |
|---|---|
| FT379.record-only-cycle: A record-only correction that changes a disposition consumes no repair cycle. | RC1, RC2, RC3 |
| FT379.evidence-confirmation: The review record has an evidence-only confirmation form that a biting coordinator probe satisfies. | RC4–RC8, RC20–RC23, RC38–RC42 |
| FT379.evidence-confirmation: The chunk checkpoint accepts that form in place of three fresh axes for a test-only repair delta. | RC9–RC19, RC24–RC36, RC43, RC44, RC48–RC54 |
| FT379.test-only-skip: A test-only repair whose coordinator probe bites closes through the evidence-only confirmation form and needs no confirming round. | RC9, RC10, RC29, RC45–RC47 |
| roadmap/FT379.md: A confirming round on a record-only correction checks only the findings that it names. | delivered; no row |
| roadmap/FT379.md: A new preference in that prose is advice. | delivered by the optional-advice rule; no row |

### Pre-review proof checklist

- `Cited symbols`: `checkReviews`, `CheckReviews`, `checkSource`, `checkVerification`, `checkCompletion`, `completionReconciler`, `RecordChunk`, `RecordReview`, `validateReviews`, `evidenceIDs`, `CompletionOrder`, `importedProbe`, `commentgap.Prove`, `TreeChangesIncludingSubmodules`, and `outline.IsGoTestFile` resolve in the tree.
- `Import edges`: the new prover imports `internal/git` and `internal/outline`; `internal/reviewrecord` imports the new prover. The ticket runs `go list` before it locks the edge.
- `Source-row clauses and occurrences`: the source-sentence table quotes each clause; the occurrences are DC-C9 rounds 9 to 14 and the FT217 ALC-OWNER, ALC-LINK, and ALC-UNLINK repairs.
- `Promised field labels`: `confirmations`, role `coordinator-confirmation`, form `confirmation`, flags `--probe-result`, `--mutant`, `--mutation`, and `--finding`.
- `Changed-function callers`: `bench consumers` lists each caller below.
  - `checkReviews`: `checkSource` at `coverage.go` lines 115 and 121, and `CheckReviews`.
  - `checkSource`: `CheckTrees`, `RecordCompletion`, and `CheckSource`.
  - `evidenceIDs`: `Parse` and `unclaimed`.
  - `CompletionOrder`: `completionOrder` in the record command.
- `Copy survival`: none.
- `Rendered-shape readers`: the completion-order help text in `cmd/bench/help_inventory_test.go` and the list in `TestCompletionHelpOrdering` change with the new form.
- `Pin operators`: the anchor rows use exact-needle `RequireInSection` and `Require` matches.
- `Entry reads`: none.
- `Derived expectations`: RC45 keeps the independent order list in `TestCompletionHelpOrdering`.
- `Consolidated rules`: the consolidation table under Probe evidence gives each site's old and new rule. RC37 covers the changed cells.
- `Quantified obligations`: "every changed path" in the prover takes RC11 to RC19; "the union of the open finding IDs" takes RC24 and RC25.
- `Workflow-step writes`: `bench record confirmation` writes only the record file, which every source digest excludes.

### Not verified

- The exact line counts of the new test files against any structure budget.
- The slicing research settled one item: `TestHelpInventoryIsComplete` asserts each form's usage row, so ticket 04 writes `cmd/bench/help_inventory_test.go`.

### Red record

Each RC1 and RC2 pin has an anchor and an independent twin, and each one turned red under a `bench probe` omit mutation. Every probe reported `restored=yes`.
The anchor reds ran `bench probe` with `--check docs-currency-workflow`.
The twin reds ran `bench probe` with `--package ./internal/conformance --run TestImplementationContinuation`.

| row | red | mutation site | failing test | observed diagnostic |
|---|---|---|---|---|
| RC1 | anchor | the RC1 sentence in the bounded repair policy | `TestRootConformance` | `gate: implementation continuation: bounded repair counted a record-only correction` |
| RC1 | twin | the RC1 row of `implementationContinuationAnchors` in the retained-workflow registry | `TestImplementationContinuation` | `implementation-continuation anchor is absent: implementation continuation: bounded repair counted a record-only correction` |
| RC2 | anchor | the RC2 sentence in the bounded repair policy | `TestRootConformance` | `gate: implementation continuation: bounded repair dropped the changed-disposition reaffirmation` |
| RC2 | twin | the RC2 row of `implementationContinuationAnchors` in the retained-workflow registry | `TestImplementationContinuation` | `implementation-continuation anchor is absent: implementation continuation: bounded repair dropped the changed-disposition reaffirmation` |
