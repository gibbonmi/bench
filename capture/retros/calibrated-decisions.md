## Outcome

The landing added a stated confidence to three decision surfaces: the delegate done-claim row, the review finding, and the line declaration.
It added the Brier score rule, a calibration table to the retro scaffold, and a calibration measure with one column per provider scorecard.
It also added the step-scoped anchor kind `require-in-step`, the `step` column of the anchors projection, and 27 calibration canary fixtures.
The build recorded its own first 61 pairs in the review pickup before the retire step removed that pickup.
The landing published commit `5139240fe02a22199f5114eff5c5eee0386970a7` from reviewed source `3c8af8dfa110c3617f52fad98775d1d4ff0f183e` on base `8e4cf48021ec5a3a576d2f60e40480059e1e8c08`.

## Gate-stage timings

- landing: commit `5139240fe02a22199f5114eff5c5eee0386970a7`, trace `7f2c538645aa1ac30ee877db0abec3ea`
- gofmt: 131 ms
- vet: 1264 ms
- test: 109184 ms
- race: 3082 ms
- system: 31475 ms
- shellcheck: 539 ms

## Ticket-versus-spec-slice and delegate performance

The run was delegated: seven Opus / high authors, one per ticket, in six chunks CD1, CD2, CD2b, CD3, CD4, and CD5.
Tickets 1, 2, and 3 landed in one pass each, and every one of their done-claim rows held under the coordinator's probes.
Ticket 7 was a plan expansion the reviewer approved after the CD2 review. The CR22 step placement needed a step-scoped anchor kind that the tree lacked.
Ticket 7 took two repair cycles, because the first return kept two owners of the step narrowing and left two predicates without a biting test.

Tickets 4 and 5 shared one CD4 repair round, and ticket 5 needed the CR25 gate expansion for the README `unknown` cell.
Ticket 6 took one repair cycle, because its first return left the line-declaration pair without a label and mis-stated a file count.
The first review pass ran Fable / medium for CD1, CD2, and CD2b, and Opus / medium from CD3 on. Every later pass ran Sonnet, at xhigh for CD2 and at high after that.

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| delegate return | CR1 | verified | 9 | held | opus / high / author |
| delegate return | CR2 | verified | 9 | held | opus / high / author |
| delegate return | CR3 | verified | 9 | held | opus / high / author |
| delegate return | CR13 | verified | 9 | held | opus / high / author |
| delegate return | CR14 | verified | 9 | held | opus / high / author |
| delegate return | CR32 | verified | 9 | held | opus / high / author |
| delegate return | CR34 | verified | 9 | held | opus / high / author |
| delegate return | CR37 | verified | 9 | held | opus / high / author |
| delegate return | CR4 | verified | 9 | held | opus / high / author |
| delegate return | CR5 | verified | 9 | held | opus / high / author |
| delegate return | CR6 | verified | 9 | held | opus / high / author |
| delegate return | CR15 | verified | 9 | held | opus / high / author |
| delegate return | CR16 | verified | 9 | held | opus / high / author |
| delegate return | CR22 | verified | 9 | held | opus / high / author |
| delegate return | CR23 | verified | 8 | held | opus / high / author |
| delegate return | CR30 | verified | 9 | held | opus / high / author |
| delegate return | CR7 | verified | 9 | held | opus / high / author |
| delegate return | CR8 | verified | 9 | held | opus / high / author |
| delegate return | CR9 | verified | 9 | held | opus / high / author |
| delegate return | CR10 | verified | 9 | held | opus / high / author |
| delegate return | CR11 | verified | 9 | held | opus / high / author |
| delegate return | CR12 | verified | 9 | held | opus / high / author |
| delegate return | CR31 | verified | 9 | held | opus / high / author |
| delegate return | CR36 | verified | 9 | held | opus / high / author |
| delegate return | CR17 | verified | 9 | held | opus / high / author |
| delegate return | CR18 | verified | 9 | held | opus / high / author |
| delegate return | CR19 | verified | 8 | held | opus / high / author |
| delegate return | CR20 | verified | 9 | held | opus / high / author |
| delegate return | CR21 | verified | 9 | held | opus / high / author |
| delegate return | CR28 | verified | 9 | held | opus / high / author |
| delegate return | CR35 | verified | 7 | held | opus / high / author |
| delegate return | CR24 | verified | 9 | held | opus / high / author |
| delegate return | CR25 | claimed | 8 | held | opus / high / author |
| delegate return | CR26 | verified | 9 | held | opus / high / author |
| delegate return | CR22 | verified | 9 | held | opus / high / author |
| delegate return | CR38 | verified | 9 | held | opus / high / author |
| review finding | CD1-S1 | claimed | 5 | refuted | fable / medium / Standards |
| review finding | CD1-S2 | claimed | 6 | held | fable / medium / Standards |
| review finding | CD2-S1 | claimed | 8 | refuted | fable / medium / Standards |
| review finding | CD2-S2 | claimed | 4 | refuted | fable / medium / Standards |
| review finding | CD2b-S1 | claimed | 8 | held | fable / medium / Standards |
| review finding | CD2b-S2 | claimed | 9 | held | fable / medium / Standards |
| review finding | CD2b-S3 | claimed | 6 | held | fable / medium / Standards |
| review finding | CD2b-S4 | claimed | 4 | refuted | fable / medium / Standards |
| review finding | CD2-C1 | claimed | 7 | held | fable / medium / Coverage |
| review finding | CD2-C2 | claimed | 3 | refuted | fable / medium / Coverage |
| review finding | CD2b-C1 | claimed | 8 | held | fable / medium / Coverage |
| review finding | CD2b-C2 | claimed | 6 | refuted | fable / medium / Coverage |
| review finding | CD2b-C3 | claimed | 5 | refuted | fable / medium / Coverage |
| review finding | CD2b-C4 | claimed | 4 | held | fable / medium / Coverage |
| review finding | CD2-C3 | claimed | 6 | held | sonnet / xhigh / Coverage |
| review finding | CD2-C4 | claimed | 8 | refuted | sonnet / xhigh / Coverage |
| review finding | CD3-S1 | claimed | 6 | refuted | opus / medium / Standards |
| review finding | CD3-S2 | claimed | 5 | refuted | opus / medium / Standards |
| review finding | CD4-S1 | claimed | 8 | held | opus / medium / Standards |
| review finding | CD4-S2 | claimed | 9 | held | opus / medium / Standards |
| review finding | CD3-C1 | claimed | 3 | refuted | opus / medium / Coverage |
| review finding | CD4-C1 | claimed | 9 | held | opus / medium / Coverage |
| review finding | CD4-C2 | claimed | 8 | refuted | opus / medium / Coverage |
| review finding | CD4-C3 | claimed | 7 | refuted | opus / medium / Coverage |
| line declaration | Expected repair rounds: 1 / confidence 6 | claimed | 6 | refuted | fable / low / orchestrator |

The Brier mean is 0.104 over 61 labeled pairs, and the abstention count is 0.
The author rows alone give 0.015 over 36 pairs, the review rows give 0.228 over 24 pairs, and the one orchestrator row gives 0.36.
The CR22 row appears twice, because the CD2 chunk and the CD2b chunk each returned it as its own claim.

## Coordinator catches

The coordinator caught the first CD2 checkpoint refusal. A plan commit and two main merges sat outside the recorded chunk pair, so the pair was widened and later passes covered the extra delta.
The coordinator derived the chain rule from that refusal. Plan commits and main merges go before the ticket merge, and only record commits follow the chunk tip.

The coordinator read the pre-repair file on CD2b and refuted two Coverage findings, CD2b-C2 and CD2b-C3, that described behavior the axis's own probes had added.
The coordinator ran silent probes on the `stepScoped` predicate and on the period check. Each silent probe sent the author back for a biting test.
The coordinator refused a post-hoc confidence raise from the ticket 4 author and froze every confidence at return time.
The coordinator caught the fixture-closure red after ticket 7, because the ticket 2 `Writes:` list did not yet name the new fixture.
The coordinator caught the first landing refusal on the stale review base and retried with the merge base, under which every refused path matched main.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-state-claim-schema-on-delegate-return.md | 0 | none |
| 2-state-finding-confidence-in-review.md | 0 | none |
| 3-declare-expected-repair-rounds-with-score.md | 0 | none |
| 4-render-calibration-table-in-retro-scaffold.md | 1 | delegate-error |
| 5-define-calibration-measure-in-scorecard.md | 1 | spec-row |
| 6-record-first-pairs-in-own-retro.md | 1 | delegate-error |
| 7-pin-the-pickup-step-with-a-step-anchor.md | 2 | delegate-error; delegate-error |

Tickets 4 and 5 shared one CD4 chunk repair round, so the build counted four rounds against the declared one.

## Agent-experience improvements

### Bench CLI

- The `calibrated-decisions-spec` census recorded 1 raw call, one `ls` head; `bench worktree path` already names the tree, so no verb change is needed.
  Feeds: none
- Make `bench worktree land` name the merge base in a fence refusal when the reviewed base predates a main merge, because every refused path matched main.
  Feeds: new
- Make the checkpoint chain-gap refusal name the rule that plan commits and main merges belong inside the chunk delta.
  Feeds: new
- Make `bench worktree release` accept a review sibling whose landing is pending, because every review worktree stayed retained until the landing.
  Feeds: new

### Skills

- State in the bounded repair policy that a stated confidence freezes at return time and that a post-hoc raise is refused.
  Feeds: new
- State in `craft-review` that a Coverage finding describes the pre-probe tree, because two CD2b findings described behavior the axis's own probes added.
  Feeds: new

### Process

- Record a checkpoint chain rule in the implement-spec phase that puts plan commits and main merges before the ticket merge.
  Feeds: new
- Count a shared chunk repair round once in the line-declaration label, because tickets 4 and 5 shared one round.
  Feeds: none
