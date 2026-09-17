# Calibrated decisions: a stated confidence on every labeled claim

Status: staged

Decision source: ready compiled map `specs/calibrated-decisions/decisions/calibrated-decisions.md`, confirmed 2026-09-17.

Verification log: 2 iteration(s) to accept — one opus/high round found six blocking findings. They were on anchor bytes, the fence union, a conditional row, a section order, a missing definition row, and an unverifiable retro row. One author fold closed all six and nine advisories under coordinator verification.

## Problem

Bench measures whether an agent's claim was right. It never records what the agent expected before the label arrived. A delegate returns "done", a review axis returns a finding, and a line declaration names an iteration cap. The gate, the coordinator, and the reviewer then label each claim held or refuted.

Without a stated confidence beside each claim, Bench has an accuracy measure and no calibration measure. A confident wrong claim and a hesitant wrong claim score the same, and an honest abstention scores nothing. The routing loop in the scorecards cannot reward honesty over a confident guess.

## Solution

Three decision surfaces carry a stated confidence: a delegate done-claim row, a review finding, and a line declaration's expected repair-round count. The confidence is an integer from 0 to 10. Each claim carries a status of `verified`, `claimed`, or `abstained`, and later a label of `held` or `refuted`. Only the gate, the coordinator's probe of the exact tree, or the reviewer's disposition writes a label. The Brier rule scores each labeled pair. An abstention scores 0, stays out of the Brier mean, and is counted apart.

The retro holds the pairs in one calibration table, and the retro scaffold renders that table's shape. The scorecard README defines one calibration measure, and each provider routing row carries it. The measure is one input to the existing two-run routing rule. This spec changes guidance, the retro scaffold, and the scorecard contract. It adds no verb and no threshold. The build that lands this spec records the first pairs in its own retro.

## User stories

Line: opus / high.
Implementation-line reason: The hardest chunks edit two guidance files at their prose budgets, under 33 and 14 anchors, and one Go seam every retro reads. The decisions are exact and the seams are known. The anchors and the scaffold test red a missing sentence, and the semantic rules are review-owned. The leverage override therefore routes the mid tier at high effort.
Harder chunks: CD3, CD4.

### Claim schema on the delegate return

1. As a coordinator, I want each done-claim row to carry a status of `verified`, `claimed`, or `abstained`, so that I know the author's evidence.
2. As a coordinator, I want an integer confidence from 0 to 10 on every non-abstained row, so that a later label can score it.
3. As a coordinator, I want an `abstained` row to carry no confidence, so that a withheld claim cannot be scored as a guess.
4. As a retro author, I want a claim to carry no free-text field, so that I aggregate rows without prose repair.
5. As a delegate, I want the charge to name the three status values and the integer range, so that I return the exact shape.
6. As a coordinator, I want my tree probe to label each row `held` or `refuted`, so that the author never labels its own claim.
7. As a coordinator, I want a `verified` row to receive my probe's label too, so that a red-to-green log does not label itself.
9. As a delegate, I want to return a row `abstained` when I cannot state a confidence, so that honest withholding is not a guess.

### Confidence on a review finding

10. As a review axis, I want each pickup finding to carry an integer confidence from 0 to 10, so that a disposition can score it.
11. As a reviewer, I want my disposition to label a finding, so that no model judgment labels it.
12. As a coordinator, I want a finding's confidence to never change whether it blocks, so that kind and citation keep deciding the round.
13. As a coordinator, I want the review record JSON schema unchanged, so that existing records parse and the confidence rides in the pickup line.
14. As a coordinator, I want optional advice to carry no confidence, so that advice stays outside finding totals and scores.

### Confidence on a line declaration

15. As a retained author, I want the line declaration to state expected repair rounds with a confidence, so that the retro scores the expectation.
16. As a retro author, I want the repair-attribution round count to label the expectation, so that the label needs no new evidence.
17. As any author, I want one calibration score rule for every surface, so that pairs from different surfaces aggregate.
18. As a reviewer, I want the log rule excluded, so that a wrong 0 or 10 cannot score to negative infinity.
19. As a reviewer, I want a real-number confidence excluded, so that an author does not state false precision.

### Abstention

20. As a retro author, I want an abstention to score 0 outside the Brier mean, so that the mean measures only stated claims.
21. As a retro author, I want the abstention count beside the pair count, so that over-abstention is visible.
22. As a delegate, I want no later probe to turn my abstention into a refuted claim, so that abstention is never punished.

### Label sources

23. As a coordinator, I want the gate to label an acceptance row, so that a gate-visible claim never needs a hand label.
24. As a reviewer, I want a model judgment, native or cross-family, never to write a label, so that a preference signal cannot enter the score.

### The retro calibration table

25. As a retro author, I want the scaffold to render the calibration table header under the delegate-performance heading, so that the draft carries the shape.
26. As a retro author, I want the scaffold to render one `unknown` row beneath that header, so that an absent measurement never blocks the draft.
27. As a retro author, I want each row to hold six cells, so that a pair traces to its author. The cells are the surface, the claim, the status, the confidence, the label, and the model, effort, and role.
28. As a retro author, I want to state the Brier mean and both counts below the table, so that the aggregate has evidence beneath it.
29. As a linked repository, I want the required retro heading list unchanged, so that a retro authored before this spec still parses.
30. As a retro author with no scorecard directory, I want the scaffold to render the table, so that a linked repository records pairs.
31. As a coordinator, I want the final-check guidance to name the calibration table as a retro duty, so that a retro without it is incomplete.

### The scorecard measure and routing

32. As a retro author, I want the scorecard README's Measures table to define one calibration measure, so that the measure has one source.
33. As a retro author, I want each provider routing row to carry a calibration cell, so that routing reads the measure beside observed quality.
34. As a retro author, I want a provider with no labeled pair to show `unknown` in that cell, so that no invented number enters.
35. As a reviewer, I want the calibration measure as one input to the two-run rule, never a tier mover, so that routing stays a judgment.
36. As a retro author, I want the ten-assignment cap to bound the calibration aggregate, so that the measure obeys the existing aggregation rule.

### Dogfood on this build

37. As a reviewer, I want this spec's own build to state a confidence on its three surfaces, so that the first pairs exist.
38. As a reviewer, I want that build's retro to hold the first calibration table and the pair count per role, so that validation starts.

### Reviewed exclusions

39. As a reviewer, I want no scoring verb in this spec, so that the record format is observed before code fixes it.
40. As a reviewer, I want no numeric routing threshold in this spec, so that the threshold decision waits for ten pairs.
41. As a reviewer, I want no model training, so that Bench steers only routing and prose.
42. As a guidance author, I want every changed Markdown file to pass the prose lane, so that the kit's prose contract holds.
43. As a coordinator, I want `verified` and `claimed` defined by whether the author ran the named check, so that a status means one thing.
44. As a maintainer, I want the calibration header spelled once in the retros package, so that the renderer cannot drift from the parser.

### Step-scoped anchors

45. As a maintainer, I want an anchor kind that pins a sentence inside one numbered step, so that the sentence cannot move unseen. The reviewer added this story in the CD2 review.

## Implementation decisions

- The claim schema has three fields and no others. `status` is one of `verified`, `claimed`, or `abstained`. `confidence` is an integer from 0 to 10, and it is absent when abstained. `label` is `held` or `refuted`, and only a label source writes it. The glossary terms **claim**, **stated confidence**, **outcome label**, **label source**, **calibration score**, and **abstention** already exist and are used exactly.
- The delegation discipline reference gains one `Claim schema` section that owns the schema, the status definitions, the abstention rules, and the coordinator-label rule. The delegate skill gains one pointer sentence inside its two lines of headroom.
- The finding discipline reference owns the finding confidence rule, the disposition-to-label mapping, and the no-confidence rule for optional advice. The review skill keeps its existing pointer sentence byte for byte, period included. It adds one second sentence on the same line, so the file grows by zero lines. The new sentence is anchored on its own words. The review-implementation command's pickup step states that each actionable finding line carries its confidence.
- The disposition-to-label mapping is fixed: `auto-fix` and `ask-user` label a finding `held`, and `no-op` labels it `refuted`.
- The review record JSON types in `internal/reviewrecord` do not change. The confidence rides on the finding line in the pickup and in the retro table.
- The line skill gains one reference file that owns the calibration score rule, the label-source rule, the abstention scoring rule, and the expectation label rule. The declaration block in the line skill gains one `Expected repair rounds:` line. The pointer to the reference rides inside that line, so the net growth is one line. The skill reclaims one line to hold its budget: the unanchored fan-out clause line in the declaration section is the candidate.
- The expectation label rule is fixed. The actual round count in the repair-attribution table labels the expectation `held` when it equals the expected count, and `refuted` otherwise.
- The retro keeps its nine required headings. The calibration table renders under the existing delegate-performance heading. A new required heading would red every retro authored before this spec in this repository and in every linked repository.
- The retros package owns the calibration table header as one exported constant beside the two derived-section headings. The scaffold renders that header and one `unknown` row as a third case in `scaffoldSection`, under the delegate-performance heading that precedes the repair table. No second spelling of the header exists in the renderer, and CR35 grades that promise.
- The scorecard README's Measures table gains one `calibration` row that defines the Brier mean, the pair count, and the abstention count. Its update contract states the `unknown` cell rule, the ten-assignment cap, and the two-run routing rule with calibration as one input. Each provider routing table gains a `calibration` column, and this build writes `unknown` in every cell.
- Every new guidance sentence a row cites is a registry anchor with its own omission canary fixture. The anchors live in a new registry file beside the ticket-passes precedent. Its group joins the existing composition line in the registry data file with no added line.
- Each anchor needle below is a pasted operand. The build writes these bytes and the fixtures mutate them.
  - CR1: A done-claim row carries a `status` of `verified`, `claimed`, or `abstained` and a stated confidence as an integer from 0 to 10.
  - CR34: `verified` means the author ran the named check and returns its red-to-green log, and `claimed` means an assertion with no executed check.
  - CR2: A delegate that cannot state a confidence returns the row `abstained` with no confidence.
  - CR37: No later probe turns an abstention into a refuted claim.
  - CR3: The coordinator's probe of the exact tree labels a done-claim row `held` or `refuted`, whatever its status.
  - CR32: A claim carries no free-text field.
  - CR13: The charge names the `Claim schema` section of `references/delegation-discipline.md` as the return shape.
  - CR4: A finding carries a stated confidence as an integer from 0 to 10.
  - CR5: The confidence never changes whether a finding blocks.
  - CR6: `auto-fix` and `ask-user` label a finding `held`, and `no-op` labels it `refuted`.
  - CR30: Optional advice carries no confidence.
  - CR15: A finding also states its confidence as an integer from 0 to 10.
  - CR22: Each actionable finding line carries its stated confidence.
  - CR7: The line declaration states the expected repair-round count and a stated confidence as an integer from 0 to 10.
  - CR8: One claim's calibration score is `(p - label)^2`, with `p = n / 10` and label 1 for `held` or 0 for `refuted`.
  - CR9: An abstention scores 0, stays out of the Brier mean, and is counted apart.
  - CR10: A label source is the gate, the coordinator's probe of the exact tree, or the reviewer's disposition.
  - CR36: A model judgment is never a label source.
  - CR31: The repair-attribution table's actual round count labels the expectation `held` when it equals the expected count and `refuted` otherwise.
  - CR11: Expected repair rounds: <count> / confidence <0-10>. `references/calibration-score.md` owns the score.
  - CR20: The retro fills the calibration table with one row per labeled claim: surface, claim, status, confidence, label, and model, effort, and role.
  - CR21: The retro states the Brier mean, the pair count, and the abstention count below the table, with `unknown` for a mean over zero pairs.
  - CR24: | calibration | the Brier mean over labeled pairs, the pair count, and the abstention count |
  - CR26: The calibration measure is one input to the two-run routing rule, obeys the ten-assignment cap, and never moves a tier on its own.
- The Bootstrap authority rule does not apply. This spec makes no trusted-execution or refusal-before-execution claim.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| CD1 / 1-state-claim-schema-on-delegate-return.md | A delegate return carries the claim schema with its status definitions, the coordinator labels each row by probe, and an abstained row follows its rule | CR1, CR2, CR3, CR13, CR14, CR32, CR34, CR37 | docs-currency-workflow anchors with omission canaries, guidance-prose-budgets, prose lane | no |
| CD2 / 2-state-finding-confidence-in-review.md | A review finding carries a confidence that the reviewer's disposition labels and that never changes blocking | CR4, CR5, CR6, CR15, CR16, CR22, CR23, CR30 | docs-currency-workflow anchors with omission canaries, guidance-prose-budgets, prose lane, review-owned schema check | no |
| CD2b / 7-pin-the-pickup-step-with-a-step-anchor.md | A step-scoped anchor kind pins a sentence inside one numbered step, `bench anchors` prints the step, and the pickup-confidence anchor is re-pinned on step 6 | CR22, CR38 | internal/anchors harness and diagnostics tests, cmd/bench anchors projection tests, docs-currency-workflow with a step-move canary | yes |
| CD3 / 3-declare-expected-repair-rounds-with-score.md | A line declaration states expected repair rounds with a confidence, and one reference owns the score, label-source, abstention, and expectation-label rules | CR7, CR8, CR9, CR10, CR11, CR12, CR31, CR36 | docs-currency-workflow anchors with omission canaries, guidance-prose-budgets, prose lane | yes |
| CD4 / 4-render-calibration-table-in-retro-scaffold.md, 5-define-calibration-measure-in-scorecard.md | The retro scaffold renders the calibration table, the final-check guidance names the retro duty, and the scorecard defines and carries the measure | CR17, CR18, CR19, CR20, CR21, CR24, CR25, CR26, CR28, CR35 | internal/roadmap scaffold test, internal/retros parse tests, docs-currency-workflow anchors with omission canaries, prose lane, review-owned data check | yes |
| CD5 / 6-record-first-pairs-in-own-retro.md | The build that lands this spec records the first pairs in its own retro and every edited Markdown file passes the prose lane | CR27, CR33 | review-owned at final reconciliation, prose lane | no |

## Testing decisions

- A good test shows that a rule sentence is present where the guidance says it is, and that its absence reds the gate. The anchors registry plus one omission canary per sentence gives that red.
- The retro scaffold test in `internal/roadmap` receives the table assertions. Its prior art is the timings and repair-table assertions in the same file, which read the section by heading.
- The heading invariant is held by the existing `internal/retros` parse tests over `testdata/eligible.md`, which spell the nine headings independently of the renderer.
- The guidance-prose-budgets check observes the three skill files that sit at or near budget.
- The docs-currency-workflow check is the gate seam for every anchor. The fixture-bite test proves each canary bites through its registered owner.
- The semantic rules, the disposition mapping, the provider-file data edits, and the one-source header promise are review-owned, because no parser reads them.

### Seam diagram

    trigger: a retained author fills a retro, or a coordinator reads a return
        │
        ▼
    guidance sentence  ──▶  [ anchors registry + canary fixture ]  ──▶  docs-currency-workflow verdict
                                ◀ tests attach here: the fixture-bite test mutates the sentence and expects the diagnostic
    slug + tickets dir  ──▶  [ retro scaffold: retros.RequiredHeadings + derived sections ]  ──▶  draft body
                                ◀ tests attach here: the scaffold test reads the delegate-performance section and asserts the table

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| CR1 | 1, 2, 5 | The delegation discipline reference carries, in a `Claim schema` section, the sentence that a done-claim row carries a `status` of `verified`, `claimed`, or `abstained` and a stated confidence as an integer from 0 to 10 | anchor `require-in-section` plus canary `calibration-claim-schema` under docs-currency-workflow, bitten by `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | An omitted or reworded sentence fails the fixture bite and the root anchor check |
| CR2 | 3, 9 | The same section carries the sentence that a delegate that cannot state a confidence returns the row `abstained` with no confidence | anchor plus canary `calibration-abstained-row` under docs-currency-workflow, bitten by `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | Dropping the sentence leaves abstention undefined and reds the fixture bite |
| CR3 | 6, 7 | The same section carries the sentence that the coordinator's probe of the exact tree labels a done-claim row `held` or `refuted` whatever its status | anchor plus canary `calibration-coordinator-label` under docs-currency-workflow, bitten by `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A version that lets a verified row label itself omits the sentence and reds the fixture bite |
| CR4 | 10 | The finding discipline reference carries the sentence that a finding carries a stated confidence as an integer from 0 to 10 | anchor plus canary `calibration-finding-confidence` under docs-currency-workflow, bitten by `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | Omission reds the fixture bite |
| CR5 | 12 | The finding discipline reference carries the sentence that the confidence never changes whether a finding blocks | anchor plus canary `calibration-finding-blocking` under docs-currency-workflow, bitten by `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A version that demotes a low-confidence finding drops the sentence and reds the fixture bite |
| CR6 | 11 | The finding discipline reference carries the sentence that `auto-fix` and `ask-user` label a finding `held` and `no-op` labels it `refuted` | anchor plus canary `calibration-disposition-label` under docs-currency-workflow, bitten by `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | Without the mapping, a disposition cannot label, and omission reds the fixture bite |
| CR7 | 15 | The line skill's calibration reference carries the sentence that the line declaration states the expected repair-round count and a stated confidence as an integer from 0 to 10 | anchor plus canary `calibration-line-expected-rounds` under docs-currency-workflow, bitten by `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | Omission reds the fixture bite |
| CR8 | 17 | The same reference carries the sentence that the calibration score of one claim is `(p - label)^2`, where `p = n / 10` and the label is 1 for `held` and 0 for `refuted` | anchor plus canary `calibration-brier-rule` under docs-currency-workflow, bitten by `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A log-rule or a real-number rewrite changes the bytes and reds the fixture bite |
| CR9 | 20, 21 | The same reference carries the sentence that an abstention scores 0, stays out of the Brier mean, and is counted apart | anchor plus canary `calibration-abstention-score` under docs-currency-workflow, bitten by `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A version that folds abstentions into the mean drops the sentence and reds the fixture bite |
| CR10 | 23 | The same reference carries the sentence that a label source is the gate, the coordinator's probe of the exact tree, or the reviewer's disposition | anchor plus canary `calibration-label-sources` under docs-currency-workflow, bitten by `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A version that admits a model label drops the sentence and reds the fixture bite |
| CR11 | 15 | The line skill's declaration block carries the `Expected repair rounds:` line | anchor `require-in-section` on `The declaration` plus canary `calibration-declaration-line` under docs-currency-workflow, bitten by `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A declaration without the line reds the fixture bite |
| CR12 | 15 | The line skill file holds at most 130 lines after the edit | `guidance-prose-budgets` check over the profile's budget table | Adding the line without reclaiming one reds the budget check |
| CR13 | 5 | The delegate skill carries one sentence that points the charge at the claim schema in the delegation discipline reference | anchor plus canary `calibration-delegate-pointer` under docs-currency-workflow, bitten by `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | Omission reds the fixture bite |
| CR14 | 5 | The delegate skill file holds at most 126 lines after the edit | `guidance-prose-budgets` check | Growth past the two lines of headroom reds the budget check |
| CR15 | 10 | The review skill's pointer line keeps the existing needle's bytes and period and carries the second sentence that a finding also states its confidence as an integer from 0 to 10 | anchor plus canary `calibration-review-pointer` under docs-currency-workflow, bitten by `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | Omission of the new sentence reds the fixture bite while the existing needle's three pins stay intact |
| CR16 | 10 | The review skill file holds at most 122 lines after the edit | `guidance-prose-budgets` check | Any line growth reds the budget check |
| CR17 | 25, 30 | `bench retro <slug> --scaffold` renders the calibration table header row under the delegate-performance heading in a repository with no `capture/agent-performance` directory | a new scaffold test in internal/roadmap, `TestRetroScaffoldRendersCalibrationTable`, in a `newScaffoldRepo` temp repository, reading the section through `sectionOf`, beside `internal/roadmap/retro_scaffold_test.go` (`TestRetroScaffoldParses`) | A scaffold without the header fails the assertion |
| CR18 | 26 | The same scaffold renders exactly one row of six `unknown` cells beneath the header | the same new scaffold test, `TestRetroScaffoldRendersCalibrationTable` | A scaffold that renders no row or a derived value fails the assertion |
| CR19 | 29 | `retros.Parse` accepts the nine-heading body in `testdata/eligible.md` after the change | `internal/retros/retros_test.go` (`TestParseAcceptsCanonicalRetro`, `TestEligibleFixtureKeepsRequiredHeadings`) | A new required heading makes the old body fail to parse |
| CR20 | 27, 31 | The final-check command carries the sentence that the retro fills the calibration table with one row per labeled claim holding its surface, claim, status, confidence, label, and model, effort, and role | anchor plus canary `calibration-retro-table-duty` under docs-currency-workflow, bitten by `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | Omission reds the fixture bite |
| CR21 | 28 | The final-check command carries the sentence that the retro states the Brier mean, the pair count, and the abstention count below the table, with `unknown` for a mean over zero pairs | anchor plus canary `calibration-retro-aggregate-duty` under docs-currency-workflow, bitten by `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | Omission reds the fixture bite |
| CR22 | 10, 13 | The review-implementation command's pickup step carries the sentence that each actionable finding line carries its stated confidence | anchor `require-in-step` on section `Process` step 6, plus canary `calibration-pickup-confidence` for omission and canary `calibration-pickup-step-move` for a move to step 5, under docs-currency-workflow, bitten by `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | Omission or a move to another step reds the fixture bite |
| CR38 | 45 | The anchors package accepts a `RequireInStep` kind with a `Step` field, raises the anchor's diagnostic when the needle sits in another step, raises its own diagnostic for a missing or duplicate step, and `bench anchors` prints the step column | a new step-rule test `TestAnchorHarnessStepRules` beside the anchor harness, the kind table in `internal/anchors/anchor_harness_diagnostics_test.go` (`TestEvaluatePathAnchorKinds`), and `cmd/bench/anchor_help_test.go` (`TestAnchorsReportsNeedleLines`) | A kind that ignores the step passes the positive rule and fails the move rule |
| CR23 | 13 | The exported types in `internal/reviewrecord` are byte-identical to the base commit | review-owned: the Spec axis compares the package against the frozen base | A schema edit appears in the chunk diff |
| CR24 | 32 | The scorecard README's Measures table carries a `calibration` row that defines the Brier mean, the pair count, and the abstention count | anchor plus canary `calibration-scorecard-measure` under docs-currency-workflow, bitten by `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | Omission reds the fixture bite |
| CR25 | 33, 34 | Each provider routing table carries a `calibration` column and every cell reads `unknown` | review-owned: the Spec axis reads both provider files | A missing column or an invented number appears in the diff |
| CR26 | 35, 36 | The scorecard README's update contract carries the sentence that the calibration measure is one input to the two-run routing rule, obeys the ten-assignment cap, and never moves a tier on its own | anchor plus canary `calibration-routing-input` under docs-currency-workflow, bitten by `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A version that adds a threshold drops the sentence and reds the fixture bite |
| CR27 | 37, 38 | The review pickup `reviews/calibrated-decisions.md` holds at least one labeled row for each of the three surfaces and states the pair count per role | review-owned at final reconciliation: the Spec axis reads the pickup at the final tip | A pickup without the rows fails the final acceptance reconciliation, and the retro copy is the final-check duty CR20 and CR21 pin |
| CR28 | 30 | The scaffold test repository holds no `capture/agent-performance` directory when the table renders | the same new scaffold test, `TestRetroScaffoldRendersCalibrationTable`, asserts the directory is absent | A test that seeds a scorecard proves nothing about a linked repository |
| CR30 | 14 | The finding discipline reference carries the sentence that optional advice carries no confidence | anchor plus canary `calibration-advice-no-confidence` under docs-currency-workflow, bitten by `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | Omission reds the fixture bite |
| CR31 | 16 | The line skill's calibration reference carries the sentence that the repair-attribution table's actual round count labels the expectation `held` when it equals the expected count and `refuted` otherwise | anchor plus canary `calibration-expectation-label` under docs-currency-workflow, bitten by `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | Omission reds the fixture bite |
| CR32 | 4 | The delegation discipline reference carries the sentence that a claim carries no free-text field | anchor plus canary `calibration-no-free-text` under docs-currency-workflow, bitten by `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | Omission reds the fixture bite |
| CR33 | 42 | Every Markdown file the build edits passes `bench gate-prose` | the prose lane at each `bench commit` | A sentence over the bound or a long paragraph reds the lane |
| CR34 | 43 | The `Claim schema` section carries the sentence that `verified` means the author ran the named check and returns its red-to-green log, and `claimed` means an assertion with no executed check | anchor plus canary `calibration-status-definitions` under docs-currency-workflow, bitten by `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | Without the definitions a status is a label with no meaning, and omission reds the fixture bite |
| CR35 | 44 | The renderer holds no second spelling of the calibration header, and the rendered header equals the exported retros constant | review-owned: the Standards axis greps the header text across `internal/roadmap`, and `TestRetroScaffoldRendersCalibrationTable` compares the rendered header with the constant | A second spelling appears in the grep, because a byte-equal copy passes the test alone |
| CR36 | 24 | The line skill's calibration reference carries the sentence that a model judgment is never a label source | anchor plus canary `calibration-model-judgment` under docs-currency-workflow, bitten by `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A version that admits a model label drops the sentence and reds the fixture bite |
| CR37 | 22 | The `Claim schema` section carries the sentence that no later probe turns an abstention into a refuted claim | anchor plus canary `calibration-abstention-final` under docs-currency-workflow, bitten by `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A version that punishes abstention drops the sentence and reds the fixture bite |

Not covered: story 18 — the exclusion is a Won't handle; CR8 pins the Brier bytes.
Not covered: story 19 — the exclusion is a Won't handle; CR8 pins the integer normalization.
Not covered: story 39 — the exclusion is recorded under Out of scope, and the review round removes an unlisted verb.
Not covered: story 40 — the exclusion is recorded under Out of scope; CR26 pins the no-threshold sentence.
Not covered: story 41 — the exclusion is recorded under Out of scope, and no row can observe an absent training step.

### Edge inventory

In-scope edges, each with a row:

- a retro with zero pairs (CR18, CR21)
- a repository with no scorecard (CR28)
- a retro authored before this spec (CR19)
- a provider with no pairs (CR25)
- a low-confidence finding whose kind blocks (CR5)
- optional advice (CR30)

- **Won't handle** a malformed row with no status or an out-of-range confidence — CR1 and CR2 define the shape, and ordinary done-claim verification returns it.
- **Won't handle** a log-rule score — the Brier rule in CR8 is the one surviving rule for every caller.
- **Won't handle** a real-number confidence — the integer 0 to 10 in CR1 and CR8 is the one surviving scale.
- **Won't handle** a model judgment offered as a label — CR10 keeps the three label sources for every surface.
- **Won't handle** a numeric threshold that moves a tier — CR26 keeps the two-run rule as the surviving routing caller until ten pairs exist.
- **Won't handle** a scoring verb — the retro author's hand computation in CR21 is the surviving caller.
- **Won't handle** a claim with no label at retro time — the table holds labeled claims only under CR20, and the retro prose names it.
- **Won't handle** a provider row over the ten-assignment cap — the README's existing cap rule aggregates the latest ten, and CR26 restates it.
- **Won't handle** a claim cell that contains a `|` character — the retro author rewords the claim, and the surviving caller is the Markdown reader.
- **Won't handle** control bytes, numeric-looking cells, and Unicode separators from the hostile-input checklist — the prose lane reads every hand-written table, and no TOON sink exists.
- **Won't handle** a dangling or live symlink at a guidance path — the existing checks refuse special files unread, and this spec adds no reader.
- **Won't handle** the anchored sentence "Known-flaky retry stops are in `craft-delegate`'s delegation discipline." and every other current anchor — the build keeps their bytes, and a reflow runs the fixture-bite check.

## Ownership fences

- `.agents/skills/bench-craft-delegate/SKILL.md`
- `.agents/skills/bench-craft-delegate/references/delegation-discipline.md`
- `.agents/skills/bench-craft-review/SKILL.md`
- `.agents/skills/bench-craft-review/references/finding-discipline.md`
- `.agents/skills/bench-craft-line/SKILL.md`
- `.agents/skills/bench-craft-line/references/calibration-score.md`
- `.agents/commands/bench-review-implementation.md`
- `.agents/commands/bench-final-check.md`
- `capture/agent-performance/README.md`
- `capture/agent-performance/claude-models.md`
- `capture/agent-performance/open-ai-models.md`
- `internal/retros/retros.go`
- `internal/roadmap/retro_scaffold.go`
- `internal/roadmap/retro_scaffold_test.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry.go`
- `internal/anchors/match.go`
- `internal/anchors/match_test.go`
- `internal/anchors/locate.go`
- `internal/anchors/anchor_harness_test.go`
- `internal/anchors/anchor_harness_diagnostics_test.go`
- `internal/anchors/registry_data_test.go`
- `cmd/bench/anchors_command.go`
- `cmd/bench/anchor_help_test.go`
- `tests/canary/workflow-guidance-anchors/calibration-pickup-step-move`
- `specs/calibrated-decisions/tickets/7-pin-the-pickup-step-with-a-step-anchor.md`
- `internal/anchors/registry_calibration.go`
- `internal/anchors/registry_calibration_test.go`
- `tests/canary/workflow-guidance-anchors/calibration-claim-schema`
- `tests/canary/workflow-guidance-anchors/calibration-abstained-row`
- `tests/canary/workflow-guidance-anchors/calibration-coordinator-label`
- `tests/canary/workflow-guidance-anchors/calibration-no-free-text`
- `tests/canary/workflow-guidance-anchors/calibration-status-definitions`
- `tests/canary/workflow-guidance-anchors/calibration-abstention-final`
- `tests/canary/workflow-guidance-anchors/calibration-delegate-pointer`
- `tests/canary/workflow-guidance-anchors/calibration-finding-confidence`
- `tests/canary/workflow-guidance-anchors/calibration-finding-blocking`
- `tests/canary/workflow-guidance-anchors/calibration-disposition-label`
- `tests/canary/workflow-guidance-anchors/calibration-advice-no-confidence`
- `tests/canary/workflow-guidance-anchors/calibration-review-pointer`
- `tests/canary/workflow-guidance-anchors/calibration-pickup-confidence`
- `tests/canary/workflow-guidance-anchors/calibration-line-expected-rounds`
- `tests/canary/workflow-guidance-anchors/calibration-brier-rule`
- `tests/canary/workflow-guidance-anchors/calibration-abstention-score`
- `tests/canary/workflow-guidance-anchors/calibration-label-sources`
- `tests/canary/workflow-guidance-anchors/calibration-model-judgment`
- `tests/canary/workflow-guidance-anchors/calibration-expectation-label`
- `tests/canary/workflow-guidance-anchors/calibration-declaration-line`
- `tests/canary/workflow-guidance-anchors/calibration-retro-table-duty`
- `tests/canary/workflow-guidance-anchors/calibration-retro-aggregate-duty`
- `tests/canary/workflow-guidance-anchors/calibration-scorecard-measure`
- `tests/canary/workflow-guidance-anchors/calibration-routing-input`
- `tests/canary/claude-agent-definitions/agent-unnamed-in-skill`
- `tests/canary/claude-agent-definitions/skill-names-missing-agent`
- `tests/canary/claude-agent-definitions/model-declared`
- `tests/canary/claude-agent-definitions/name-mismatch`
- `tests/canary/claude-agent-definitions/shell-tool-absent`
- `tests/canary/claude-agent-definitions/spawning-tool`
- `tests/canary/claude-agent-definitions/tools-absent`
- `tests/canary/docs-currency-token-diet/introduces-undeclared-command`
- `tests/canary/docs-currency-token-diet/stale-command-reference`
- `tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns`
- `tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `tests/canary/` every canary fixture, in any family, that a ticket's `Writes:` line co-names because it pins an edited file
- `reviews/calibrated-decisions.md`
- `CHANGELOG.md`
- `specs/calibrated-decisions/spec.md`
- `decisions/calibrated-decisions.md` and `decisions/calibrated-decisions/`, the retired top-level map paths the spec commit moved under this folder

## Out of scope

- A Bench verb that reads the retro tables and prints the calibration aggregate: 6 edits, 3 gate runs. It waits for ten recorded pairs.
- A numeric threshold that moves a tier on its own: 3 edits, 1 gate run. The reviewer decides it after ten pairs.
- Training or fine-tuning any model: not a Bench capability, 0 edits.
- A reference-model label as an outcome source: 0 edits. Invariant 1 forbids it.
- A new required retro heading: 5 edits, 2 gate runs. It would red every retro authored under the current list.

## Further notes

### Flagged additions

Four edge dispositions went beyond the decision source. The review round removed the first, and each of the other three keeps its row.

1. The malformed-row rule was removed. Its story, its row CR29, its anchor, and its fixture are gone. A malformed row is a Won't handle under ordinary done-claim verification.
2. The disposition-to-label mapping (CR6): the map names the reviewer's disposition as the label source and never maps the three dispositions to `held` and `refuted`.
3. The expectation label rule (CR31): the map names the expected repair rounds as a surface and never says what labels it.
4. The omission of an unlabeled claim from the table at retro time (Won't handle line).

### Source-sentence-to-row table

| source | predicate | rows |
| --- | --- | --- |
| ticket 2 | A done-claim carries one confidence per acceptance row it claims | CR1, CR13 |
| ticket 2 | A review finding carries one confidence that the finding holds against the tree | CR4, CR22 |
| ticket 2 | A line declaration carries one confidence on its expected repair-round count | CR7, CR11 |
| ticket 2 | An acceptance row itself and a ticket `Writes:` path set carry none | Won't handle by omission; no row adds one |
| ticket 3 | An integer from 0 to 10, normalized to `p = n / 10`, scored by `(p - label)^2` | CR8 |
| ticket 3 | The log rule is excluded | CR8, story 18 |
| ticket 4 | The gate, the coordinator's probe, and the reviewer's disposition are the only label sources | CR3, CR6, CR10, CR36 |
| ticket 5 | An abstention scores 0, stays out of the Brier mean, and is counted apart | CR9, CR21 |
| ticket 5 | No later probe turns an abstention into a refuted claim | CR37 |
| ticket 6 | `status`, `confidence`, `label`, and no free text | CR1, CR32 |
| ticket 6 | `verified` means the author ran the named check and returns its red-to-green log; `claimed` means an assertion with no executed check | CR34 |
| ticket 7 | Each retro gains one calibration table with the six fields | CR17, CR18, CR20, CR35 |
| ticket 7 | The scorecard README gains the measure and each provider row carries it | CR24, CR25 |
| ticket 8 | One input to the two-run rule and no threshold | CR26 |
| ticket 9 | A confidence never changes whether a finding blocks | CR5 |
| ticket 10 | Guidance-only, the Brier mean by hand | CR21, Out of scope |
| ticket 11 | This build records the first pairs in its own retro | CR27 |
| map Out of scope | No model training, no reference-model label, no verb, no threshold | stories 39 to 41, CR10, CR26 |
| map Out of scope | The Jev type-safety guarantee is not a Bench claim | covered by omission: no story, row, or decision claims a schema guarantee |

### Pre-review proof checklist

- Cited symbols: `retros.RequiredHeadings`, `retros.RepairHeading`, `retros.TimingsHeading`, `retros.Parse`, `scaffoldSection`, `scaffoldBody`, `scaffoldTickets`, `RetroCommand`, `newScaffoldRepo`, `sectionOf`, `TestParseAcceptsCanonicalRetro`, `TestEligibleFixtureKeepsRequiredHeadings`, `TestEveryRetainedFixtureBitesThroughRegisteredOwner`, `checkDocsCurrencyAndWorkflow`, `checkWorkflowAnchors`, `anchors.Anchor`, `anchors.RequireInSection`, `anchors.AfterImplementSpec`, and `reviewrecord.Record`. Each resolved in this session by a read of its defining file.
- Import edges: `internal/roadmap/retro_scaffold.go` imports `internal/retros`; `internal/conformance/fixture_bite_test.go` imports `internal/anchors` and `internal/canary`. Both read this session.
- Source-row clauses and occurrences: the table above quotes each clause once, and each clause occurs once in its ticket file.
- Promised field labels: `status`, `confidence`, `label`, `verified`, `claimed`, `abstained`, `held`, `refuted`, `Expected repair rounds:`, and the table header cells `surface`, `claim`, `status`, `confidence`, `label`, `model / effort / role`.
- Changed-function callers: `scaffoldSection` has one caller, `scaffoldBody`. `retros.RequiredHeadings` keeps its callers `scaffoldBody` and the retros tests unchanged.
- Copy survival: the header replaces no copy. CR35 is review-owned, because a byte-equal second spelling passes the equality test alone.

### Enforcement reads

- `.agents/skills/bench-craft-line/SKILL.md` (130 lines, budget 130, 33 anchors), `.agents/skills/bench-craft-delegate/SKILL.md` (124 lines, budget 126, 50 anchors), `.agents/skills/bench-craft-review/SKILL.md` (122 lines, budget 122, 14 anchors).
- `.agents/commands/bench-review-implementation.md` (44 anchors), `.agents/commands/bench-final-check.md` (34 anchors, 30 registry rows), `capture/agent-performance/README.md` (0 anchors), `CONTEXT.md` (8 anchors), `references/delegation-discipline.md` (34 anchors), `references/bounded-repair-policy.md` (35 anchors).
- `projects/benchkit.md` prose budget table and hostile-input checklist.
- `internal/retros/retros.go`, `internal/retros/retros_test.go`, `internal/retros/testdata/eligible.md`, `internal/roadmap/retro_scaffold.go`, `internal/roadmap/retro_scaffold_test.go`, `internal/roadmap/retro.go`, `internal/roadmap/context_parse.go`.
- `internal/reviewrecord/record.go`, `internal/anchors/registry_data.go` (composition line and group counts), `internal/anchors/registry_ticket_passes.go` and its test, `internal/conformance/fixture_bite_test.go`, `internal/conformance/docs_workflow_checks_test.go`, `internal/conformance/registry/registry.go` rows 122, 147, 152, and the family map.
- Canary shapes: `tests/canary/workflow-guidance-anchors/delegate-self-probe-missing-row` and `final-check-census-read-before-land` (`BASE`, `EXPECT`, `MUTATE.json`); `tests/canary/retro-improvement-markers/retro-item-unmarked`.
- `.bench/lines.env` for the `opus` mid binding, and `.bench/BENCH-reference.md` line 44 to 48 for scorecard ownership.

### Reader sweep

- Readers of the retro heading list with a row: `retros.go` as owner, `retro_scaffold.go` (CR17), `retro_scaffold_test.go` (CR17, CR18), and `retros_test.go` with `testdata/eligible.md` (CR19).
- Readers of the retro heading list excluded: `recommendations_test.go` reads the repair table only, which does not change. `context_parse.go` and `roadmap.go` read the directory through `retros.Facts`, and no heading changes. The anchors registry needles for the final-check heading list stay valid, because the list does not change.
- Readers of the scorecard Measures table with a row: the README itself (CR24), both provider files (CR25), and the final-check retro duty (CR20, CR21).
- Readers of the scorecard Measures table excluded: the `.bench/BENCH-reference.md` ownership line, because ownership does not change. Also excluded: `internal/landing/settlepolicy/settlepolicy_test.go`, which uses the path as a blob label only.
- Readers of the delegate return shape: `craft-delegate` SKILL (CR13) and its reference (CR1), `internal/preflight/charge.go` (no return-shape text; excluded).
- Readers of the review pickup line format: `.agents/commands/bench-review-implementation.md` step 6 (CR22), `internal/reviewrecord` fenced JSON parser (CR23, unchanged).
- `.mjs` scripts and `.github/workflows`: `native-runtime.yml`, `release.yml`, and `scripts/gremlins-diff.sh` name none of the changed facts; excluded.

### Sources re-read

- All four arXiv abstracts and the TypeSafe page were re-opened on 2026-09-17 and still carry the quoted claims. The Rewarding Doubt abstract does not show the 0-to-10 scale; that fact rests on the full-text read recorded in the asset on 2026-09-16.
- Local paths in the map's Sources were read: `capture/agent-performance/README.md` and `.bench/BENCH.md` invariant 1.

### Completion plan

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"calibrated-decisions-full-20260917","orchestrator_session":"claude:session_01PAwDzo8gg18CxyXMeGdzrN","author_limit":1,"assignments":{"1-state-claim-schema-on-delegate-return.md":[{"session":"claude:bench-writer/cd-t1-author","assignment":"cd-t1-author","model":"opus","effort":"high","source":"80cb9f7cd3a1be1154f3debee03919deb4139f9d","native_ref":"claude:agent/cd-t1-author-20260917@80cb9f7cd3a1be1154f3debee03919deb4139f9d"}],"2-state-finding-confidence-in-review.md":[{"session":"claude:bench-writer/cd-t2-author","assignment":"cd-t2-author","model":"opus","effort":"high","source":"7ce1266321c1a2bd8a974dd34255d161aa665cdf","native_ref":"claude:agent/cd-t2-author-20260917@7ce1266321c1a2bd8a974dd34255d161aa665cdf"}],"3-declare-expected-repair-rounds-with-score.md":[{"session":"claude:bench-writer/cd-t3-author","assignment":"cd-t3-author","model":"opus","effort":"high","source":"b47edde2c7a5bc4f00a77ea94e77ac09a6df24ea","native_ref":"claude:agent/cd-t3-author-20260917@b47edde2c7a5bc4f00a77ea94e77ac09a6df24ea"}],"4-render-calibration-table-in-retro-scaffold.md":[],"5-define-calibration-measure-in-scorecard.md":[],"6-record-first-pairs-in-own-retro.md":[],"7-pin-the-pickup-step-with-a-step-anchor.md":[{"session":"claude:bench-writer/cd-t7-author","assignment":"cd-t7-author","model":"opus","effort":"high","source":"5435fdf63f80997b153c097de439e4a70d2018df","native_ref":"claude:agent/cd-t7-author-20260917@5435fdf63f80997b153c097de439e4a70d2018df"}]}},"chunks":[{"id":"CD1","tickets":["1-state-claim-schema-on-delegate-return.md"],"verification":[{"id":"workflow","command":"bench test --check docs-currency-workflow","ticket":"1-state-claim-schema-on-delegate-return.md"},{"id":"budgets","command":"bench test --check guidance-prose-budgets","ticket":"1-state-claim-schema-on-delegate-return.md"},{"id":"prose","command":"bench test --check prose-mechanics","ticket":"1-state-claim-schema-on-delegate-return.md"}]},{"id":"CD2","tickets":["2-state-finding-confidence-in-review.md"],"verification":[{"id":"workflow","command":"bench test --check docs-currency-workflow","ticket":"2-state-finding-confidence-in-review.md"},{"id":"budgets","command":"bench test --check guidance-prose-budgets","ticket":"2-state-finding-confidence-in-review.md"},{"id":"prose","command":"bench test --check prose-mechanics","ticket":"2-state-finding-confidence-in-review.md"}]},{"id":"CD2b","tickets":["7-pin-the-pickup-step-with-a-step-anchor.md"],"verification":[{"id":"workflow","command":"bench test --check docs-currency-workflow","ticket":"7-pin-the-pickup-step-with-a-step-anchor.md"},{"id":"anchors","command":"bench test --package ./internal/anchors/...","ticket":"7-pin-the-pickup-step-with-a-step-anchor.md"},{"id":"projection","command":"bench test --package ./cmd/bench/... --run TestAnchors","ticket":"7-pin-the-pickup-step-with-a-step-anchor.md"},{"id":"bite","command":"bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner","probe":"move the pickup-confidence sentence from step 6 to step 5","ticket":"7-pin-the-pickup-step-with-a-step-anchor.md"}]},{"id":"CD3","tickets":["3-declare-expected-repair-rounds-with-score.md"],"verification":[{"id":"workflow","command":"bench test --check docs-currency-workflow","ticket":"3-declare-expected-repair-rounds-with-score.md"},{"id":"budgets","command":"bench test --check guidance-prose-budgets","ticket":"3-declare-expected-repair-rounds-with-score.md"},{"id":"prose","command":"bench test --check prose-mechanics","ticket":"3-declare-expected-repair-rounds-with-score.md"}]},{"id":"CD4","tickets":["4-render-calibration-table-in-retro-scaffold.md","5-define-calibration-measure-in-scorecard.md"],"verification":[{"id":"workflow-4","command":"bench test --check docs-currency-workflow","ticket":"4-render-calibration-table-in-retro-scaffold.md"},{"id":"workflow-5","command":"bench test --check docs-currency-workflow","ticket":"5-define-calibration-measure-in-scorecard.md"},{"id":"budgets","command":"bench test --check guidance-prose-budgets","ticket":"5-define-calibration-measure-in-scorecard.md"},{"id":"prose","command":"bench test --check prose-mechanics","ticket":"5-define-calibration-measure-in-scorecard.md"},{"id":"scaffold","command":"bench test --package ./internal/roadmap/...","ticket":"4-render-calibration-table-in-retro-scaffold.md"},{"id":"retros","command":"bench test --package ./internal/retros/...","ticket":"4-render-calibration-table-in-retro-scaffold.md"}]},{"id":"CD5","tickets":["6-record-first-pairs-in-own-retro.md"],"verification":[{"id":"workflow","command":"bench test --check docs-currency-workflow","ticket":"6-record-first-pairs-in-own-retro.md"},{"id":"budgets","command":"bench test --check guidance-prose-budgets","ticket":"6-record-first-pairs-in-own-retro.md"},{"id":"prose","command":"bench test --check prose-mechanics","ticket":"6-record-first-pairs-in-own-retro.md"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/calibrated-decisions/spec.md"},{"id":"workflow","command":"bench test --check docs-currency-workflow"},{"id":"budgets","command":"bench test --check guidance-prose-budgets"},{"id":"prose","command":"bench test --check prose-mechanics"},{"id":"scaffold","command":"bench test --package ./internal/roadmap/..."},{"id":"retros","command":"bench test --package ./internal/retros/..."},{"id":"anchors","command":"bench test --package ./internal/anchors/..."},{"id":"projection","command":"bench test --package ./cmd/bench/... --run TestAnchors"}]}
```

### Assumptions

- The `unknown` cell in every provider routing row is a data edit this build makes so the column exists before the first retro fills it.
- The canary fixture names above are the author's proposal; the ticket fork may rename them, and the fence follows the tickets.
