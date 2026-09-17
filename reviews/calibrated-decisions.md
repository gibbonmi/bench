# Calibrated decisions

## CD1 ticket 1 author evidence

The ticket adds one `Claim schema` section to the delegation discipline reference. It adds one pointer sentence to the delegate skill. Seven anchors and seven omission canary fixtures grade the seven sentences.

### Scenario rows

These two rows show the claim schema shape the ticket describes. The table is illustrative: its label cell shows the scenario the ticket states, not a coordinator probe. These rows do not count as pairs.

| surface | claim | status | confidence | label |
| --- | --- | --- | --- | --- |
| delegate return | the named check runs green on the returned tree | claimed | 7 | held |
| delegate return | the guidance reads clearly for a cold reader | abstained |  |  |

The coordinator probes the first row's named check on the exact tree and labels it `held`. The coordinator labels nothing on the second row and counts one abstention.

### Done-claim table

The author did not write the label cell. The coordinator wrote each label. The label source is the gate lane at the merge, green at tip 09f26779. One independent swap probe on the CR13 sentence also bit.

| row | status | confidence | label |
| --- | --- | --- | --- |
| CR1 | verified | 9 | held |
| CR2 | verified | 9 | held |
| CR3 | verified | 9 | held |
| CR13 | verified | 9 | held |
| CR14 | verified | 9 | held |
| CR32 | verified | 9 | held |
| CR34 | verified | 9 | held |
| CR37 | verified | 9 | held |

### Red-then-green log

Each row below except CR14 uses `bench probe` with the omission kind against `docs-currency-workflow`. The row for CR14 uses a line-growth swap against `guidance-prose-budgets`. The verb records the green baseline, then the red under the mutation, then the proved restore.

| row | red line under the mutation | green |
| --- | --- | --- |
| CR1 | `gate: calibration: a done-claim row needs its status and stated confidence` | baseline passed, restored yes |
| CR2 | `gate: calibration: an unstatable confidence needs an abstained row` | baseline passed, restored yes |
| CR3 | `gate: calibration: the coordinator's tree probe labels each done-claim row` | baseline passed, restored yes |
| CR13 | `gate: calibration: the charge must name the claim schema section` | baseline passed, restored yes |
| CR14 | `gate: prose-budget exceeded: .agents/skills/bench-craft-delegate/SKILL.md is 128 lines, over its 126-line budget` | baseline passed, restored yes |
| CR32 | `gate: calibration: a claim must carry no free-text field` | baseline passed, restored yes |
| CR34 | `gate: calibration: verified and claimed need their evidence definitions` | baseline passed, restored yes |
| CR37 | `gate: calibration: an abstention must never become a refuted claim` | baseline passed, restored yes |

`TestEveryRetainedFixtureBitesThroughRegisteredOwner` proves each of the seven new canaries bites through its registered owner. The same run proves the 525 earlier fixtures keep their planted diagnostics.

### Probe verdict

The self-probe omits the CR3 sentence from the delegation discipline reference and runs `docs-currency-workflow`. The verdict line reads `bit`, and `restored` reads `yes`.

### Verification table

Each run reports no skip. The elapsed time is the package time the verb reports.

| check | verdict | elapsed |
| --- | --- | --- |
| `bench test --check docs-currency-workflow` | pass | 890 ms |
| `bench test --check guidance-prose-budgets` | pass | 4 ms |
| `bench test --package ./internal/anchors/... --run 'TestCalibration'` | pass | 15 ms |
| `bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner` | pass | 9299 ms |
| `go vet ./...` | pass | no output |
| `bench gate-prose` on both edited files | pass | two pass rows |

## CD2 ticket 2 author evidence

The ticket adds one `What a confidence states` section to the finding discipline reference. It adds one second sentence to the review skill's pointer line and one sentence to the review-implementation pickup step. Six anchors and six omission canary fixtures grade the six sentences.

### Scenario finding

This row shows the finding shape the ticket describes. The label cell holds the scenario the ticket states, not a reviewer disposition on real work. The row does not count as a pair.

| surface | claim | citation | status | confidence | disposition | label |
| --- | --- | --- | --- | --- | --- | --- |
| review finding | the behavior under the cited line is wrong | the diff line the axis read this pass | claimed | 3 | no-op | refuted |

The finding blocks the round, because its kind and its citation decide. The confidence of 3 changes nothing about that block. The reviewer disposes the finding `no-op`, and the fixed mapping labels it `refuted`.

### Done-claim table

The author did not write a label cell. The coordinator wrote each label. The label source is the gate lane at the merge, green at tip 1ecb33ec. One independent swap probe on the CR6 mapping also bit.

| row | status | confidence | label |
| --- | --- | --- | --- |
| CR4 | verified | 9 | held |
| CR5 | verified | 9 | held |
| CR6 | verified | 9 | held |
| CR15 | verified | 9 | held |
| CR16 | verified | 9 | held |
| CR22 | verified | 9 | held |
| CR23 | verified | 8 | held |
| CR30 | verified | 9 | held |

### Red-then-green log

Each row below except CR16 and CR23 uses `bench probe` with the omission kind against `docs-currency-workflow`. The row for CR16 uses a line-growth swap against `guidance-prose-budgets`. The row for CR23 uses a diff of the package against the base commit. The verb records the green baseline, then the red under the mutation, then the proved restore.

| row | red line under the mutation | green |
| --- | --- | --- |
| CR4 | `gate: calibration: a finding needs its stated confidence` | baseline passed, restored yes |
| CR5 | `gate: calibration: the confidence must never change whether a finding blocks` | baseline passed, restored yes |
| CR6 | `gate: calibration: the dispositions need their label mapping` | baseline passed, restored yes |
| CR15 | `gate: calibration: the review skill must point at the finding confidence` | baseline passed, restored yes |
| CR16 | `gate: prose-budget exceeded: .agents/skills/bench-craft-review/SKILL.md is 123 lines, over its 122-line budget` | baseline passed, restored yes |
| CR22 | `gate: calibration: the pickup line must carry its stated confidence` | baseline passed, restored yes |
| CR23 | no mutation: `git diff --stat <base> -- internal/reviewrecord` gives empty output | the package is byte-identical to the base |
| CR30 | `gate: calibration: optional advice must carry no confidence` | baseline passed, restored yes |

`TestEveryRetainedFixtureBitesThroughRegisteredOwner` proves each of the six new canaries bites through its registered owner. A control run with one wrong `EXPECT` byte made that fixture fail, which shows the run reads the new fixtures. The same run proves every earlier fixture keeps its planted diagnostic.

### Probe verdict

The self-probe omits the CR15 sentence from the review skill and runs `docs-currency-workflow`. The verdict line reads `bit`, and `restored` reads `yes`.

### Verification table

Each run reports no skip. The elapsed time is the package time the verb reports.

| check | verdict | elapsed |
| --- | --- | --- |
| `bench test --check docs-currency-workflow` | pass | 1185 ms |
| `bench test --check guidance-prose-budgets` | pass | 5 ms |
| `bench test --package ./internal/anchors/... --run 'TestCalibration'` | pass | 29 ms |
| `bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner` | pass | 11923 ms |
| `go vet ./...` | pass | no output |
| `git diff --stat <base> -- internal/reviewrecord` | empty | no output |
| `bench gate-prose` on the three edited files and this pickup | pass | four pass rows |

## CD2b ticket 7 author evidence

The ticket adds the `RequireInStep` anchor kind and a `Step` field on `Anchor`. The kind resolves the H2 section first, then narrows that body to one numbered step. The `bench anchors` projection gains a `step` column after `section`. The pickup-confidence anchor moves from `require-in-section` to step 6 of `Process`.

### Done-claim table

The author wrote no label cell. Each status is `verified`, because the author ran the named check and kept its log. The coordinator's probe on `stepScoped` was silent against the first diff. The kind now owns the step narrowing on both sides, and that probe bites. Two repair cycles closed six more silent mutations.

| row | status | confidence | label |
| --- | --- | --- | --- |
| CR22 | verified | 9 |  |
| CR38 | verified | 9 |  |

CR38 holds at 9. The step parser and all four scope boundaries carry a biting test now. A reviewer reads the guarantee from the suite, and not from the code.

### One owner for the step narrowing

The first diff decided the narrowing twice. The evaluator read `anchor.Step != 0`, and the locator read `kind.stepScoped()`. A section-scoped anchor with a step would then narrow in the evaluator and not in the locator.

The kind is the one owner now. `resolveStep` raises its own diagnostic when a step-scoped anchor names no step. `TestRegistryBindsStepToItsKind` refuses a registry row whose kind and `Step` field disagree, so no such row reaches the evaluator.

### Repair cycle 1

Seven findings arrived together. Two name a production defect, and the tree refutes both. Four hold as test gaps or comment defects. One asks for a decision, and this section states it.

CD2b-C2 says `stepOpener` trims leading space. The parser reads `line[0]` directly and never trims. An indented `6.` line therefore opens no step. CD2b-C3 says the parser reads one digit. The scan takes every leading digit, and `strconv.Atoi` reads them together, so `10.` reads as step 10.

A scratch test at this tip printed `"   6. indented" -> (0, false)` and `"10. ten" -> (10, true)`. The author deleted that scratch file. Both findings are `no-op` on production. Each names a real test gap, because both mutations were silent. Both gaps carry a biting test now.

CD2b-C1 says the step close boundary has no biting test. The registered opener leads in each moved tree now. A body that ran past the next opener would read the moved needle as its own. CD2b-C4 says the section open boundary has none. `TestMarkdownH2SectionExcludesItsHeading` grades the body and the heading's own line.

CD2b-S1 asks for a decision on `Locate`. The exported signature stays as it is, and the doc comment states the step-0 case. The signature is read by `locate_test.go`, which sits outside this ticket's fence.

CD2b-S2 and CD2b-S3 are comment repairs. The narrating comment reads timeless now. The step-scoping rationale sits once on `stepScoped` in `match.go`, and production names no test.

| finding | fix | probe | verdict |
| --- | --- | --- | --- |
| CD2b-C2 | no production change; a harness case writes an indented numbered continuation inside the registered step | trim the line before the opener test | `bit,internal/anchors/match.go,swap,failed,1,yes` |
| CD2b-C3 | no production change; `TestStepOpenerReadsEveryLeadingDigitAtColumnZero` grades `1.`, `10.`, `06.`, and four refusals | read one digit instead of every digit | `bit,internal/anchors/match.go,swap,failed,2,yes` |
| CD2b-C1 | the extra opener leads the rules' own, and two moved cases join the table | the step `closes` test answers false | `bit,internal/anchors/locate.go,swap,failed,2,yes` |
| CD2b-C4 | `TestMarkdownH2SectionExcludesItsHeading` grades the section's open boundary | the `keepOpener` branch never runs | `bit,internal/anchors/locate.go,swap,failed,1,yes` |
| CD2b-S1 | the `Locate` doc states the step-0 answer; the signature stays | none | comment |
| CD2b-S2 | the diagnostics-test comment reads timeless | none | comment |
| CD2b-S3 | the rationale sits once on `stepScoped`; the two restatements and the test name are gone | none | comment |

### Repair cycle 2

The Coverage pass reopened CD2b-C4. The cycle-1 test grades the section arm of the shared narrowing, where the body starts under its heading. The step arm inverts that rule, and it had no test. A mutation that starts every body one line down was therefore silent on the step side.

`TestMarkdownNumberedStepsIncludesItsOpener` grades the step arm. A step's own first words sit on its opener line, so the body must start with that line. The test also reads a needle from the opener line through `Satisfied`, and it holds the next opener out of the body.

| finding | fix | probe | verdict |
| --- | --- | --- | --- |
| CD2b-C4, step arm | `TestMarkdownNumberedStepsIncludesItsOpener` grades the step body's open boundary | the `keepOpener` branch collapses to `start = i + 1` | `bit,internal/anchors/locate.go,swap,failed,1,yes` |

The mutation fails that one test, and no other. This result names the step arm as the half the earlier test could not reach.

`06.` reads as step 6. A markdown reader sees `06.` and `6.` as the same step, and the doc comment on `stepOpener` states that rule.

### Red-then-green log

`TestAnchorHarnessStepRules` came first. The package did not compile, and the compiler named `undefined: RequireInStep`, `anchor.Step undefined (type Anchor has no field or method Step)`, and `unknown field step in struct literal of type anchorRule`. The kind, the field, and the harness support turned that red green.

The CR22 step move is a live-tree red. The author copied `.agents/commands/bench-review-implementation.md` aside. The author then deleted the pinned sentence from step 6 and inserted the identical bytes into step 5. `bench test --check docs-currency-workflow` failed with `gate: calibration: the pickup line must carry its stated confidence`. The author restored the copy, and `cmp` reported no difference. `git status --short` on the path reported no change, and the check passed again.

The new `calibration-pickup-step-move` canary plants that same move. A verbose run of `TestEveryRetainedFixtureBitesThroughRegisteredOwner` shows `calibration-pickup-confidence` and `calibration-pickup-step-move` each pass as its own subtest. The omission canary keeps its planted diagnostic.

The re-pinned row reads `require-in-step,Process,6,Each actionable finding line carries its stated confidence.,178`.

### Probe verdict

The self-probe swaps the step-digit test in `internal/anchors/locate.go` for a test that accepts any opener. The verdict line reads `bit,internal/anchors/locate.go,swap,failed,1,yes`, and the failed test is a moved case of `TestAnchorHarnessStepRules`.

The coordinator's probe swaps the body of `stepScoped` in `internal/anchors/match.go` for `return false`. The verdict line reads `bit,internal/anchors/match.go,swap,failed,3,yes`. The failed tests are the moved, the no-such-step, and the duplicated-step cases of `TestAnchorHarnessStepRules`.

Repair cycle 1 adds four probes. Each verdict line sits in the table above, and each came back `bit`.

A later probe omits the period test from the opener. Digits and blank space then read as an opener, and two refusal rows hold that shape out. The verdict line reads `bit,internal/anchors/match.go,omit,failed,2,yes`.

### Verification table

Each run reports no skip. The elapsed time is the package time the verb reports.

| check | verdict | elapsed |
| --- | --- | --- |
| `bench test --package ./internal/anchors/...` | pass | 407 ms |
| `bench test --package ./cmd/bench/... --run TestAnchors` | pass | 113 ms |
| `bench test --check docs-currency-workflow` | pass | 805 ms |
| `bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner` | pass | 9298 ms |
| `go vet ./...` | pass | no output |
| `gofmt -l internal/anchors cmd/bench` | pass | no output |
| `bench gate-prose . -- reviews/calibrated-decisions.md` | pass | one row |

`internal/anchors/locate.go` first grew to 452 lines, over its 400-line budget. One narrowing walk now serves both the section resolution and the step resolution, and the file reads 400 lines.

## CD3 ticket 3 author evidence

The ticket adds `references/calibration-score.md` under the line skill. The reference owns the score rule, the abstention rule, the label sources, and the expectation label. The declaration block in `.agents/skills/bench-craft-line/SKILL.md` gains one `Expected repair rounds:` line, and the pointer to the reference rides inside it. Seven anchors and seven omission canaries pin the new sentences.

### This build's line declaration

The coordinator declared `Expected repair rounds: 1 / confidence 6` for this run. That declaration is the calibration surface of this ticket, and the retro scores it later.

| surface | claim | status | confidence | label |
| --- | --- | --- | --- | --- |
| line declaration | one repair round after the initial review | claimed | 6 |  |

The status is `claimed`, because no check labels an expected round count. The repair-attribution table of this build supplies the label.

### The reclaimed line

The skill sits at its 130-line budget, so the edit reclaims one line. The reclaimed line is the unanchored fan-out clause line:

`Declare fan-out for visibility before spend. Report an overrun like a ladder move. Derive a numeric cap from expected cycles plus one red. Price a likely shift repair higher.`

`bench anchors .agents/skills/bench-craft-line/SKILL.md` listed no needle on that line. Its four sentences now follow the anchored iteration-policy sentence on one line, so no rule and no anchored byte is lost. The file reads 130 lines.

### Done-claim table

The author wrote no label cell. Each status is `verified`, because the author ran the named check and kept its red-then-green log. The coordinator wrote each label. The label source is the gate lane at the merge, green at tip 7c3ef892. One independent swap probe on the declaration line's range also bit.

| row | status | confidence | label |
| --- | --- | --- | --- |
| CR7 | verified | 9 | held |
| CR8 | verified | 9 | held |
| CR9 | verified | 9 | held |
| CR10 | verified | 9 | held |
| CR11 | verified | 9 | held |
| CR12 | verified | 9 | held |
| CR31 | verified | 9 | held |
| CR36 | verified | 9 | held |

### Red-then-green log

Each anchored row ran `bench probe <file> --omit "<the sentence>" --check docs-currency-workflow`. Each baseline passed, and each omission failed the check with the row's own diagnostic. The probe restored every subject.

| row | omitted sentence owner | verdict | diagnostic |
| --- | --- | --- | --- |
| CR7 | `references/calibration-score.md` | `bit,...,omit,failed,1,yes` | calibration: the line declaration needs its expected rounds and stated confidence |
| CR8 | `references/calibration-score.md` | `bit,...,omit,failed,1,yes` | calibration: the Brier rule needs its score expression |
| CR9 | `references/calibration-score.md` | `bit,...,omit,failed,1,yes` | calibration: an abstention needs its own scoring rule |
| CR10 | `references/calibration-score.md` | `bit,...,omit,failed,1,yes` | calibration: the three label sources must stay named |
| CR36 | `references/calibration-score.md` | `bit,...,omit,failed,1,yes` | calibration: a model judgment must never label a claim |
| CR31 | `references/calibration-score.md` | `bit,...,omit,failed,1,yes` | calibration: the round count must label the expectation |
| CR11 | `SKILL.md` | `bit,...,omit,failed,1,yes` | calibration: the declaration must state expected repair rounds |

CR12 ran a line-growth swap against `guidance-prose-budgets`. The swap splits the joined line back into two lines. The verdict line reads `bit,.agents/skills/bench-craft-line/SKILL.md,swap,failed,1,yes`, and the check reported `prose-budget exceeded: .agents/skills/bench-craft-line/SKILL.md is 131 lines, over its 130-line budget`.

`TestEveryRetainedFixtureBitesThroughRegisteredOwner` passed with the seven new canaries. Every pre-existing fixture that pins the line skill still plants its diagnostic.

### Probe verdict

The self-probe swaps the Brier sentence for the same sentence with the log rule. The verdict line reads `bit,.agents/skills/bench-craft-line/references/calibration-score.md,swap,failed,1,yes`. The failed test is `TestRootConformance`, and it reported `calibration: the Brier rule needs its score expression`.

### Verification table

Each run reports no skip. The elapsed time is the package time the verb reports.

| check | verdict | elapsed |
| --- | --- | --- |
| `bench test --check docs-currency-workflow` | pass | 1095 ms |
| `bench test --check guidance-prose-budgets` | pass | 4 ms |
| `bench test --package ./internal/anchors/... --run 'TestCalibration'` | pass | 34 ms |
| `bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner` | pass | 9261 ms |
| `go vet ./...` | pass | no output |
| `bench gate-prose . -- three Markdown files` | pass | see below |
| `wc -l .agents/skills/bench-craft-line/SKILL.md` | 130 | no output |

## CD4 ticket 4 author evidence

The ticket exports the calibration table header from the retros package and renders it in the retro scaffold. The renderer gains one case in `scaffoldSection`, under the delegate-performance heading. The final-check command gains the table duty and the aggregate duty, each with one anchor and one omission canary.

### Scaffold scenario

`bench retro calibrated-decisions --scaffold` printed this delegate-performance section:

```markdown
## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| unknown | unknown | unknown | unknown | unknown | unknown |
```

The scaffold test repository holds no `capture/agent-performance` directory, and the test asserts that absence.

### Done-claim table

The author wrote no label cell. Each status is `verified`, because the author ran the named check and kept its red-then-green log. CR35 is `claimed`, because no check grades a second spelling; the review greps `internal/roadmap` for it.

| row | status | confidence | label |
| --- | --- | --- | --- |
| CR17 | verified | 9 |  |
| CR18 | verified | 9 |  |
| CR19 | verified | 8 |  |
| CR20 | verified | 9 |  |
| CR21 | verified | 9 |  |
| CR28 | verified | 9 |  |
| CR35 | claimed | 7 |  |

### TDD red

The test came first. `go test ./internal/roadmap/... -run TestRetroScaffoldRendersCalibrationTable` failed to build with `undefined: retros.DelegateHeading` and `undefined: retros.CalibrationHeader`. After the two constants landed, and before the renderer case existed, the same run reported `delegate-performance section holds 1 lines, want the header, the separator, and one row`. The renderer case turned it green.

### Red-then-green log

| row | red | green |
| --- | --- | --- |
| CR17, CR18, CR28 | `delegate-performance section holds 1 lines` | `bench test --package ./internal/roadmap/...` passes |
| CR19 | no red; the heading list keeps its nine members | `bench test --package ./internal/retros/...` passes |
| CR20 | `bit,...,omit,failed,1,yes` with `calibration: the retro must fill the calibration table` | `bench test --check docs-currency-workflow` passes |
| CR21 | `bit,...,omit,failed,1,yes` with `calibration: the retro must state the Brier mean and the counts` | `bench test --check docs-currency-workflow` passes |

Each omission probe ran `bench probe .agents/commands/bench-final-check.md --omit "<the sentence>" --check docs-currency-workflow`. Each baseline passed, each omission failed `TestRootConformance` with the row's own diagnostic, and the probe restored the subject. `TestEveryRetainedFixtureBitesThroughRegisteredOwner` passed with the two new canaries, so every pre-existing fixture on the final-check command still plants its diagnostic. `bench anchors .agents/commands/bench-final-check.md` lists 36 needles, which is the 34 prior needles and these two.

### Probe verdict

The self-probe swaps the rendered row list for a list that holds the row twice. The verdict line reads `bit,internal/roadmap/retro_scaffold.go,swap,failed,1,yes`. The failed test is `TestRetroScaffoldRendersCalibrationTable`, and it reported `delegate-performance section holds 4 lines`.

### Verification table

Each run reports no skip. The elapsed time is the package time the verb reports.

| check | verdict | elapsed |
| --- | --- | --- |
| `bench test --package ./internal/roadmap/...` | pass | 1927 ms |
| `bench test --package ./internal/retros/...` | pass | 11 ms |
| `bench test --check docs-currency-workflow` | pass | 946 ms |
| `bench test --package ./internal/anchors/... --run 'TestCalibration'` | pass | 38 ms |
| `bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner` | pass | 9632 ms |
| `go vet ./...` and `gofmt -l` | pass | no output |
| `bench gate-prose . -- two Markdown files` | pass | no output |
| `bench structure` | no new issue | 105 pre-existing issues |

`bench structure` lists no file this ticket edits.

## CD1 review

The frozen pair is base `9148850200714f000eec2fbf44cddea6182f95c7` and tip `09f26779f65b7938f313cff9ec877fabe9d009f5`. The reviewer directed the review line. The first pass of every chunk review runs fable / medium. Every later pass on the same chunk runs sonnet / xhigh. Each axis ran in its own read-only worktree.

The raw finding count is 2, both on the Standards axis. The de-duplicated repair target count is 1. CD1 consumed 0 repair cycles. The one accepted repair is an evidence-only prose correction inside this pickup.

### Standards

Finding count: 2 judgment calls, 0 hard violations. Worst issue: the registry test restates the seven needles with no recorded independence red.

| id | finding | citation | confidence | disposition | label |
| --- | --- | --- | --- | --- | --- |
| CD1-S1 | `internal/anchors/registry_calibration_test.go` restates the seven needles and diagnostics from the registry file, and the pickup records omission reds through `bench probe`, not a registry-row deletion the test alone catches. | AGENTS.md, one source per fact, the test-expectation exception. | 5 | no-op | refuted |
| CD1-S2 | The red-then-green paragraph in this pickup stated a universal rule, then retracted it for CR14 in the next sentence. | ASD-STE100 prose rules; a universal statement that the next sentence retracts. | 6 | auto-fix | held |

CD1-S1 is `no-op`. The exemplar pair `registry_ticket_passes.go` and its test hold the same shape. The Coverage axis showed the red route through `anchorHarness.check` and the fixture-bite run. The coordinator repaired CD1-S2 in this pickup as an evidence-only correction. The paragraph now names the CR14 exception in its first sentence.

A sonnet / xhigh reaffirmation pass in its own worktree confirmed the correction. It found CD1-S2 gone, CD1-S1 still `no-op`, no other claim changed, and no new finding.

The axis gave two items of optional advice with no confidence. Replace the derived count "the 525 earlier fixtures" with "every earlier fixture". Prefer "for every status" over "whatever its status" in the reference. Both stay open. The reference bytes are the spec's pasted needle, and the count is run evidence.

### Spec

The finding count is 0, and there is no worst issue. The axis audited every CD1 row and found each held.

The exact needle bytes sit in the named section. Each anchor is `require-in-section`. Each fixture exists under its fence name. The delegate skill holds 125 lines. The package `internal/reviewrecord` is unchanged, and every pre-existing anchored sentence keeps its bytes.

The axis gave one item of optional advice with no confidence. The scenario table's label is illustrative. The coordinator added that note above the table.

### Coverage

The finding count is 0. No worst issue holds against the tree. The axis ran eight independent probes. Five bit and three stayed silent.

The probes were a sentence moved above its heading, a duplicate heading, a fenced copy, and a suffix weakening. They were also a case-only rewrite, a renamed skill section, a duplicated sentence, and a contradicting sibling. The three silent probes are matcher-wide semantics of `internal/anchors/locate.go`. The spec routes those semantics to review, so they are not defects of this delta.

The axis gave one item of optional advice with no confidence. The CR1 why-it-catches clause could say "omitted or byte-changed under fold" instead of "reworded". A fenced copy and a case-only rewrite stay green under the section-scoped matcher.

### Pairs recorded for CD5

| surface | claim | status | confidence | label | model / effort / role |
| --- | --- | --- | --- | --- | --- |
| delegate return | CR1 | verified | 9 | held | opus / high / author |
| delegate return | CR2 | verified | 9 | held | opus / high / author |
| delegate return | CR3 | verified | 9 | held | opus / high / author |
| delegate return | CR13 | verified | 9 | held | opus / high / author |
| delegate return | CR14 | verified | 9 | held | opus / high / author |
| delegate return | CR32 | verified | 9 | held | opus / high / author |
| delegate return | CR34 | verified | 9 | held | opus / high / author |
| delegate return | CR37 | verified | 9 | held | opus / high / author |
| review finding | CD1-S1 | claimed | 5 | refuted | fable / medium / Standards |
| review finding | CD1-S2 | claimed | 6 | held | fable / medium / Standards |

## CD2 review

The first pass reviewed base `87869759bf975c1f76a6e12ccc638e8471d40dfd` and tip `1ecb33ec4cede2be8bde4bc07bb648786e906d13`. Each axis ran fable / medium in its own read-only worktree. The checkpoint requires a chunk base whose source tree equals the CD1 tip and a chunk tip whose source tree equals the graded source. The recorded CD2 pair is therefore base `7ce1266321c1a2bd8a974dd34255d161aa665cdf` and tip `8b14208c1b87b0528d54541aa9aa4ad3936fea8d`. The uncovered delta holds the plan and spec amendments, the ticket 7 file, pickup commits, and two merged main commits on `roadmap/FT287.md`.

One sonnet / xhigh later pass per axis covered that delta in its own worktree. Standards and Spec passed with no finding. Coverage found one `Writes:` gap, CD2-C3: ticket 7 lacked `internal/anchors/match_test.go`, whose kind table drives the scoped-kind predicate. The coordinator applied it as a ticket expectation expansion under the plan-expansion policy, inside the CD2b delta. The author of ticket 7 received the fence expansion before its return.

The raw finding count is 4. The de-duplicated repair target count is 1. CD2 consumed 0 repair cycles. The one accepted finding is a new seam that the reviewer routed to a new chunk, CD2b, not a repair of this delta.

### Standards

Finding count: 2 judgment calls, 0 hard violations. Worst issue: two guidance sentences carry the 0 to 10 range.

| id | finding | citation | confidence | disposition | label |
| --- | --- | --- | --- | --- | --- |
| CD2-S1 | The finding discipline reference and the review skill each state the 0 to 10 range. | AGENTS.md, one source per fact; the spec requires both sentences verbatim as CR4 and CR15. | 8 | no-op | refuted |
| CD2-S2 | The registry test adds a second function instead of extending the first. | Ticket 2, "its harness test"; the exemplar keeps one function per topic. | 4 | no-op | refuted |

Both findings are `no-op`. The spec's pasted needles decide CD2-S1, and the ticket-passes exemplar decides CD2-S2. The axis noted that the scenario finding table omits the model, effort, and role cell; that table is illustrative and not a pair.

### Spec

The finding count is 0, and there is no worst issue. The axis audited every CD2 row and found each held. It accepted the author's heading judgment: `## What a confidence states` follows the reference's one-question-per-heading shape.

The axis gave two items of optional advice with no confidence. The CR22 anchor is section-scoped, so it cannot pin step 6. The disposition-label fixture removes two lines where one sentence would match its siblings.

### Coverage

Finding count: 2. Worst issue: the CR22 sentence can leave the pickup step while the check stays green.

| id | finding | citation | confidence | disposition | label |
| --- | --- | --- | --- | --- | --- |
| CD2-C1 | Omitting the step 6 heading moves the CR22 sentence into step 5, and `docs-currency-workflow` stays green, because `RequireInSection` is heading-scoped. | Spec row CR22, "the pickup step"; probe 8, `silent`. | 7 | ask-user | held |
| CD2-C2 | A contradicting sentence appended after the CR5 needle stays green. | The known bound of every sentence anchor; no row promises it. | 3 | no-op | refuted |

The reviewer decided CD2-C1. A step-scoped anchor kind is a new seam, and the reviewer expanded this run with one ticket for it. The ticket forms chunk CD2b, re-pins CR22 on step 6, and takes its own three-axis review. CD2 closes on the section anchor, and the finding is resolved by that decision. The axis ran eight probes: six bit and two stayed silent.

### Pairs recorded for CD5

| surface | claim | status | confidence | label | model / effort / role |
| --- | --- | --- | --- | --- | --- |
| delegate return | CR4 | verified | 9 | held | opus / high / author |
| delegate return | CR5 | verified | 9 | held | opus / high / author |
| delegate return | CR6 | verified | 9 | held | opus / high / author |
| delegate return | CR15 | verified | 9 | held | opus / high / author |
| delegate return | CR16 | verified | 9 | held | opus / high / author |
| delegate return | CR22 | verified | 9 | held | opus / high / author |
| delegate return | CR23 | verified | 8 | held | opus / high / author |
| delegate return | CR30 | verified | 9 | held | opus / high / author |
| review finding | CD2-S1 | claimed | 8 | refuted | fable / medium / Standards |
| review finding | CD2-S2 | claimed | 4 | refuted | fable / medium / Standards |
| review finding | CD2-C1 | claimed | 7 | held | fable / medium / Coverage |
| review finding | CD2-C2 | claimed | 3 | refuted | fable / medium / Coverage |
| review finding | CD2-C3 | claimed | 6 | held | sonnet / xhigh / Coverage |
| review finding | CD2-C4 | claimed | 8 | refuted | sonnet / xhigh / Coverage |

The later-pass Coverage axis issued its two findings as CD2b-C1 and CD2b-C2, and the record excerpt keeps that spelling. This pickup names them CD2-C3 and CD2-C4, because the CD2b chunk review below owns the CD2b ids.

## CD2b review

The frozen pair is base `67af9d502c9f37c4853f5f72f5684115ca0154a9` and tip `779f0e68b68e97f5a313da6d6e90c4b778af4052`. The delta holds ticket 7, one coordinator-probe repair before commit, the ticket 2 fixture-closure fix, and one merged main commit outside the fence. Each axis ran fable / medium in its own read-only worktree.

The raw finding count is 8. The de-duplicated repair target count is 7. CD2b consumed both of its 2 repair cycles. Cycle 1 carried the seven `auto-fix` findings below to the ticket 7 author as one batch. Cycle 2 carried the reopened step arm of CD2b-C4.

### Standards

Finding count: 2 hard violations, 2 judgment calls. Worst issue: the exported `Locate` doc comment omits that a step-scoped call through it answers 0.

| id | finding | citation | confidence | disposition | label |
| --- | --- | --- | --- | --- | --- |
| CD2b-S1 | The `Locate` doc comment lists the zero cases and omits the step-kind case, because `Locate` forwards step 0. | craft-comments, aging: update a comment over edited code. | 8 | auto-fix | held |
| CD2b-S2 | A diagnostics test comment narrates a change with "no longer reads around". | craft-comments, no narration. | 9 | auto-fix | held |
| CD2b-S3 | The "kind alone decides step scoping" rationale appears in three files, and production names a test. | craft-comments and AGENTS.md, one source owns a fact. | 6 | auto-fix | held |
| CD2b-S4 | The shared walk takes five parameters with a flag argument. | Fowler smell baseline, long parameter list; the collapse bought the line budget. | 4 | no-op | refuted |

The axis confirmed four clean points. The section walk reproduces the old semantics. `Anchor` stays comparable. Step scoping has one owner, and the opener has one parser. The column list widened in lockstep across its six test derivations.

### Spec

The finding count is 0, and there is no worst issue. CR22 and CR38 held. The re-pinned anchor, the two-entry step-move fixture, the four harness cases, the kind table, and the `step` column match the spec. The shared narrowing walk is within the ticket. A second copy of the fence toggle and origin mapping would duplicate knowledge.

### Coverage

Finding count: 4. Worst issue: the step body's close boundary has no biting test, so a move of the sentence into step 7 would pass.

| id | finding | citation | confidence | disposition | label |
| --- | --- | --- | --- | --- | --- |
| CD2b-C1 | Replacing the step `closes` predicate with `return false` stays green on every check; the existing moved case writes the needle's step before the bare opener. | Ticket 7, "The step body runs to the next such line"; three silent probes. | 8 | auto-fix | held |
| CD2b-C2 | The opener trims leading space, so an indented numbered line opens a step. | Ticket 7, an indented continuation stays inside its step; two silent probes. | 6 | auto-fix | refuted |
| CD2b-C3 | The opener parses one digit, so `10.` reads as step 1. | Ticket 7, literal digits the reader sees; one silent probe. | 5 | auto-fix | refuted |
| CD2b-C4 | Removing the `keepOpener` branch of the shared walk stays green. | The refactor moved the heading-exclusion decision into a parameter with no assertion. | 4 | auto-fix | held |

Ten probes ran: one bit on the fence toggle, and nine stayed silent. The silent ones are the four findings above and one item of advice.

The coordinator read the pre-repair opener at tip 779f0e68 and refuted CD2b-C2 and CD2b-C3 as stated. The opener reads column zero without a trim and parses every leading digit. Their silent probes were real test gaps, and the repair closed them with new rows. The labels above score the claims as written.

After repair cycle 1, one sonnet / xhigh later pass per axis ran at tip da019659. Standards and Spec passed. Coverage reran the four mutations. Three bit. The step arm of the shared walk stayed silent, because the cycle-1 test covers only the section arm. CD2b-C4 reopened for that arm, and repair cycle 2 closes it with one test.

After repair cycle 2, each axis reaffirmed at tip d4ddc088 in its own worktree. The new test bites both the shared-walk branch and its call-site value. No finding remains open on CD2b.

## CD3 review

The frozen pair is base `b47edde2c7a5bc4f00a77ea94e77ac09a6df24ea` and tip `7c3ef89216a5797de657e77ff26670f638260d2a`. The reviewer changed the review line before this chunk. The first pass now runs opus / medium, and every later pass runs sonnet / high. Each axis ran in its own read-only worktree.

The raw finding count is 3. The de-duplicated repair target count is 0. CD3 consumed 0 repair cycles. The author's fold of the fan-out sentences onto the anchored iteration-policy line stays open for the reviewer's veto.

### Standards

Finding count: 2 judgment calls, 0 hard violations. Worst issue: the reference pointer sits inside the copy-paste declaration template line.

| id | finding | citation | confidence | disposition | label |
| --- | --- | --- | --- | --- | --- |
| CD3-S1 | The `Expected repair rounds:` template line carries the reference pointer as a second sentence. | ste-prose, one instruction per sentence; the spec's Implementation decisions fix the needle bytes and state that the pointer rides inside that line. | 6 | no-op | refuted |
| CD3-S2 | The folded line carries four sentences on one physical line against the file's one-sentence-per-line convention. | File style; paragraph structure and sentence lengths hold. | 5 | no-op | refuted |

CD3-S1 contradicts a spec keep decision, so it is `no-op`. The axis enumerated the score rule, the label-source list, and the abstention rule across the guidance tree and found no second operative spelling.

### Spec

The finding count is 0, and there is no worst issue. Every CD3 row held. The fold is within the spec's decision, which names the fan-out line as the candidate and does not require deletion. The pickup holds this build's declaration as a `claimed` row at confidence 6.

### Coverage

Finding count: 0 blocking. Worst issue: none. The axis ran eight probes: five bit, and three stayed silent.

| id | finding | citation | confidence | disposition | label |
| --- | --- | --- | --- | --- | --- |
| CD3-C1 | A fourth label source appended after the CR10 sentence stays green. | Spec Won't-handle line for a model judgment offered as a label; CR10's require anchor is the whole defense. | 3 | no-op | refuted |

The axis gave one item of optional advice with no confidence. A needle demoted into a code fence keeps its anchor green, which is a matcher-wide property that predates this chunk.

### Pairs recorded for CD5

| surface | claim | status | confidence | label | model / effort / role |
| --- | --- | --- | --- | --- | --- |
| delegate return | CR7 | verified | 9 | held | opus / high / author |
| delegate return | CR8 | verified | 9 | held | opus / high / author |
| delegate return | CR9 | verified | 9 | held | opus / high / author |
| delegate return | CR10 | verified | 9 | held | opus / high / author |
| delegate return | CR11 | verified | 9 | held | opus / high / author |
| delegate return | CR12 | verified | 9 | held | opus / high / author |
| delegate return | CR31 | verified | 9 | held | opus / high / author |
| delegate return | CR36 | verified | 9 | held | opus / high / author |
| review finding | CD3-S1 | claimed | 6 | refuted | opus / medium / Standards |
| review finding | CD3-S2 | claimed | 5 | refuted | opus / medium / Standards |
| review finding | CD3-C1 | claimed | 3 | refuted | opus / medium / Coverage |

The axis gave two items of optional advice with no confidence. The step cache key drops the step number, which one registered step anchor cannot expose. A `06.` opener parses as step 6, and the ticket's literal-digits rule leaves leading zeros undecided.

### Pairs recorded for CD5

| surface | claim | status | confidence | label | model / effort / role |
| --- | --- | --- | --- | --- | --- |
| review finding | CD2b-S1 | claimed | 8 | held | fable / medium / Standards |
| review finding | CD2b-S2 | claimed | 9 | held | fable / medium / Standards |
| review finding | CD2b-S3 | claimed | 6 | held | fable / medium / Standards |
| review finding | CD2b-S4 | claimed | 4 | refuted | fable / medium / Standards |
| review finding | CD2b-C1 | claimed | 8 | held | fable / medium / Coverage |
| review finding | CD2b-C2 | claimed | 6 | refuted | fable / medium / Coverage |
| review finding | CD2b-C3 | claimed | 5 | refuted | fable / medium / Coverage |
| review finding | CD2b-C4 | claimed | 4 | held | fable / medium / Coverage |

## Native review record

The fenced payload below retains every terminal return for the checkpoint. The first CD1 Standards occurrence lists its two findings; a later occurrence supersedes it after the evidence-only correction is reaffirmed. The first CD2 Coverage occurrence lists its actionable finding; the reaffirmation supersedes it after the reviewer decision. The amendment list maps each plan digest to its successor across the two assignment commits and the CD2b expansion.

```bench-review-record
{
  "version": 2,
  "spec": "specs/calibrated-decisions/spec.md",
  "plan_digest": "sha256:b76aff6080777f4bd95ef7ea7be9b075fb5fd511e8f54b14c92681be2ccbb6a9",
  "implementation_session": "",
  "amendments": [
    {"from": "sha256:37ef0eb88a3496e3fc048c85e66bd6b4db2ac676d1dca36243cb5e160430ffe7", "to": "sha256:b76aff6080777f4bd95ef7ea7be9b075fb5fd511e8f54b14c92681be2ccbb6a9", "chunk_ids": {"CD1": ["CD1"], "CD2": ["CD2"], "CD2b": ["CD2b"], "CD3": ["CD3"], "CD4": ["CD4"], "CD5": ["CD5"]}},
    {"from": "sha256:2a20907b4162c42f8fdc101392a17e0be6300188b009aaabf0ad40dd64f0a105", "to": "sha256:37ef0eb88a3496e3fc048c85e66bd6b4db2ac676d1dca36243cb5e160430ffe7", "chunk_ids": {"CD1": ["CD1"], "CD2": ["CD2"], "CD2b": ["CD2b"], "CD3": ["CD3"], "CD4": ["CD4"], "CD5": ["CD5"]}},
    {"from": "sha256:9a5079a914a8c51fea3eccf16a76b5dbb9a687c7eebec78a9bb84c36bfc6ff94", "to": "sha256:bd4a4ea1b9ec05c6aa378023c9c43ff8fae6021cf41f88949e8359bda3f6f601", "chunk_ids": {"CD1": ["CD1"], "CD2": ["CD2"], "CD3": ["CD3"], "CD4": ["CD4"], "CD5": ["CD5"]}},
    {"from": "sha256:bd4a4ea1b9ec05c6aa378023c9c43ff8fae6021cf41f88949e8359bda3f6f601", "to": "sha256:ba0c975ed0afb1fb695519ebdcbedb9b0335d30f4c276c1573643285ee9f9588", "chunk_ids": {"CD1": ["CD1"], "CD2": ["CD2"], "CD3": ["CD3"], "CD4": ["CD4"], "CD5": ["CD5"]}},
    {"from": "sha256:ba0c975ed0afb1fb695519ebdcbedb9b0335d30f4c276c1573643285ee9f9588", "to": "sha256:2a20907b4162c42f8fdc101392a17e0be6300188b009aaabf0ad40dd64f0a105", "chunk_ids": {"CD1": ["CD1"], "CD2": ["CD2"], "CD2b": ["CD2b"], "CD3": ["CD3"], "CD4": ["CD4"], "CD5": ["CD5"]}}
  ],
  "chunks": [
    {
      "id": "CD1",
      "base": "9148850200714f000eec2fbf44cddea6182f95c7",
      "tip": "09f26779f65b7938f313cff9ec877fabe9d009f5",
      "plan_digest": "sha256:9a5079a914a8c51fea3eccf16a76b5dbb9a687c7eebec78a9bb84c36bfc6ff94",
      "source_digest": "e4866f1d8569f350474dcf28334a23f195baedd8",
      "acceptance_rows": ["CR1", "CR2", "CR3", "CR13", "CR14", "CR32", "CR34", "CR37"],
      "verification": [
        {
          "id": "cd1-workflow",
          "performer": "claude:bench-writer/cd-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e4866f1d8569f350474dcf28334a23f195baedd8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-t1-author-20260917/cd1-workflow@09f26779",
            "digest": "sha256:22335da2557fa9f4e7ee2983872b3ccb666dda93bbc3cd85f52d14fc03f34491",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,817\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "workflow",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "cd1-budgets",
          "performer": "claude:bench-writer/cd-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e4866f1d8569f350474dcf28334a23f195baedd8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-t1-author-20260917/cd1-budgets@09f26779",
            "digest": "sha256:6f72c99a424331a3381e284e55ffa703b490f81e19b498e848e8bfc5d9ca25ed",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,5\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        },
        {
          "id": "cd1-prose",
          "performer": "claude:bench-writer/cd-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e4866f1d8569f350474dcf28334a23f195baedd8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-t1-author-20260917/cd1-prose@09f26779",
            "digest": "sha256:4094ea8549b8b3a7e5f00286119ddfa586e3ac1a55cd0a75a8d29fbe59e6bacd",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,156\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "prose",
          "command": "bench test --check prose-mechanics",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "cd1-standards",
          "performer": "claude:bench-reviewer/cd-c1-standards",
          "role": "independent-review",
          "model": "fable",
          "effort": "medium",
          "source_digest": "e4866f1d8569f350474dcf28334a23f195baedd8",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/cd-c1-standards-20260917@09f26779",
            "digest": "sha256:80c16504f159fca1ef6d0ca4e647ff755032e64b04b5dd96ddae18608c38f315",
            "excerpt": "Count: 0 hard violations, 2 judgment calls. CD1-S1 duplicated expectation without a recorded independence red, confidence 5, no-op. CD1-S2 STE contradiction inside one pickup paragraph, confidence 6, auto-fix."
          },
          "axis": "Standards",
          "base": "9148850200714f000eec2fbf44cddea6182f95c7",
          "tip": "09f26779f65b7938f313cff9ec877fabe9d009f5",
          "finding_ids": ["CD1-S1", "CD1-S2"],
          "supersedes": []
        },
        {
          "id": "cd1-standards-r2",
          "performer": "claude:bench-reviewer/cd-c1-standards-r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "xhigh",
          "source_digest": "e4866f1d8569f350474dcf28334a23f195baedd8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-c1-standards-r2-20260917@09f26779",
            "digest": "sha256:3bb29329cb87762c6aa5d13b5c81412238dc14bc68c213f8960525cf1e869255",
            "excerpt": "Verdict: pass. CD1-S2 repaired in the pickup; CD1-S1 no-op holds against the tree; no new findings."
          },
          "axis": "Standards",
          "base": "9148850200714f000eec2fbf44cddea6182f95c7",
          "tip": "09f26779f65b7938f313cff9ec877fabe9d009f5",
          "finding_ids": [],
          "supersedes": ["cd1-standards"]
        },
        {
          "id": "cd1-spec",
          "performer": "claude:bench-reviewer/cd-c1-spec",
          "role": "independent-review",
          "model": "fable",
          "effort": "medium",
          "source_digest": "e4866f1d8569f350474dcf28334a23f195baedd8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-c1-spec-20260917@09f26779",
            "digest": "sha256:1287bdc2526232d4fe9411913e22d0b3d1b65478f1220b82e6f5a2473b238175",
            "excerpt": "Finding count: 0 blocking, 0 auto-fix. Worst issue: none. Every CD1 row held: CR1, CR2, CR3, CR13, CR14, CR32, CR34, CR37."
          },
          "axis": "Spec",
          "base": "9148850200714f000eec2fbf44cddea6182f95c7",
          "tip": "09f26779f65b7938f313cff9ec877fabe9d009f5",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "cd1-coverage",
          "performer": "claude:bench-reviewer/cd-c1-coverage",
          "role": "independent-review",
          "model": "fable",
          "effort": "medium",
          "source_digest": "e4866f1d8569f350474dcf28334a23f195baedd8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-c1-coverage-20260917@09f26779",
            "digest": "sha256:86daef777047f41fb6b08515c5cc99f8172e72e245a7d80afb68a2fee2a9a4bb",
            "excerpt": "Finding count: 0 blocking. Worst issue: none that holds against the tree. Eight independent probes: five bit, three silent on matcher-wide semantics routed to review."
          },
          "axis": "Coverage",
          "base": "9148850200714f000eec2fbf44cddea6182f95c7",
          "tip": "09f26779f65b7938f313cff9ec877fabe9d009f5",
          "finding_ids": [],
          "supersedes": []
        }
      ]
    },
    {
      "id": "CD2",
      "base": "7ce1266321c1a2bd8a974dd34255d161aa665cdf",
      "tip": "8b14208c1b87b0528d54541aa9aa4ad3936fea8d",
      "plan_digest": "sha256:2a20907b4162c42f8fdc101392a17e0be6300188b009aaabf0ad40dd64f0a105",
      "source_digest": "a121c0a528a7979eb77ab6b45080d1478a8e34a1",
      "acceptance_rows": ["CR4", "CR5", "CR6", "CR15", "CR16", "CR22", "CR23", "CR30"],
      "verification": [
        {
          "id": "cd2-workflow",
          "performer": "claude:bench-writer/cd-t2-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a121c0a528a7979eb77ab6b45080d1478a8e34a1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-t2-author-20260917/cd2-workflow@8b14208c",
            "digest": "sha256:3858ce1c6c9d09a126f346f744c17c7b18b3d801505114aad7c709475ad419fb",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,838\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "workflow",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "cd2-budgets",
          "performer": "claude:bench-writer/cd-t2-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a121c0a528a7979eb77ab6b45080d1478a8e34a1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-t2-author-20260917/cd2-budgets@8b14208c",
            "digest": "sha256:6f72c99a424331a3381e284e55ffa703b490f81e19b498e848e8bfc5d9ca25ed",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,5\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        },
        {
          "id": "cd2-prose",
          "performer": "claude:bench-writer/cd-t2-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a121c0a528a7979eb77ab6b45080d1478a8e34a1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-t2-author-20260917/cd2-prose@8b14208c",
            "digest": "sha256:1fe8571edc70f18ded3c0bfc0ce47738ceab1301e324902abcddd2308df39c07",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,164\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "prose",
          "command": "bench test --check prose-mechanics",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "cd2-standards",
          "performer": "claude:bench-reviewer/cd-c2-standards",
          "role": "independent-review",
          "model": "fable",
          "effort": "medium",
          "source_digest": "f19df4a3214da33af5e54a114c1e57c3125d3e02",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-c2-standards-20260917@1ecb33ec",
            "digest": "sha256:7f502333470dcb01f651c27bc503e5baff5924c04166db347939d98b01fc0b5e",
            "excerpt": "Count: 0 hard violations, 2 judgment calls. CD2-S1 two guidance sources state the 0 to 10 range, confidence 8, no-op. CD2-S2 second test function instead of extending the first, confidence 4, no-op."
          },
          "axis": "Standards",
          "base": "87869759bf975c1f76a6e12ccc638e8471d40dfd",
          "tip": "1ecb33ec4cede2be8bde4bc07bb648786e906d13",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "cd2-spec",
          "performer": "claude:bench-reviewer/cd-c2-spec",
          "role": "independent-review",
          "model": "fable",
          "effort": "medium",
          "source_digest": "f19df4a3214da33af5e54a114c1e57c3125d3e02",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-c2-spec-20260917@1ecb33ec",
            "digest": "sha256:d843d14e024bf89611be77c314fd029db3ac267e877c223a642d98c310353eb3",
            "excerpt": "Finding count: 0. Worst issue: none. Every CD2 row held: CR4, CR5, CR6, CR15, CR16, CR22, CR23, CR30. The author's heading judgment is accepted."
          },
          "axis": "Spec",
          "base": "87869759bf975c1f76a6e12ccc638e8471d40dfd",
          "tip": "1ecb33ec4cede2be8bde4bc07bb648786e906d13",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "cd2-coverage",
          "performer": "claude:bench-reviewer/cd-c2-coverage",
          "role": "independent-review",
          "model": "fable",
          "effort": "medium",
          "source_digest": "f19df4a3214da33af5e54a114c1e57c3125d3e02",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/cd-c2-coverage-20260917@1ecb33ec",
            "digest": "sha256:990cbaa0ebaa3ab37f47e12ca45d65b63e028748182e49aca4ff5755d7ec8b94",
            "excerpt": "Finding count: 2. CD2-C1 the CR22 sentence can leave the pickup step under a section-scoped anchor, confidence 7, ask-user. CD2-C2 an appended contradiction beside CR5 stays green, confidence 3, no-op."
          },
          "axis": "Coverage",
          "base": "87869759bf975c1f76a6e12ccc638e8471d40dfd",
          "tip": "1ecb33ec4cede2be8bde4bc07bb648786e906d13",
          "finding_ids": ["CD2-C1"],
          "supersedes": []
        },
        {
          "id": "cd2-coverage-r2",
          "performer": "claude:bench-reviewer/cd-c2-coverage-r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "xhigh",
          "source_digest": "f19df4a3214da33af5e54a114c1e57c3125d3e02",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-c2-coverage-r2-20260917@1ecb33ec",
            "digest": "sha256:c86f3caa3e5b0b928f2605f254bcbcc5ca162b5d726c6f592a34436646a3eb85",
            "excerpt": "Verdict: pass. CD2-C1 routed to chunk CD2b by reviewer decision; CD2-C2 no-op; no unresolved Coverage finding at tip 1ecb33ec; no new finding."
          },
          "axis": "Coverage",
          "base": "87869759bf975c1f76a6e12ccc638e8471d40dfd",
          "tip": "1ecb33ec4cede2be8bde4bc07bb648786e906d13",
          "finding_ids": [],
          "supersedes": ["cd2-coverage"]
        },
        {
          "id": "cd2-standards-r3",
          "performer": "claude:bench-reviewer/cd-c2-standards-r3",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "xhigh",
          "source_digest": "a121c0a528a7979eb77ab6b45080d1478a8e34a1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-c2-standards-r3-20260917@8b14208c",
            "digest": "sha256:3024b632b7b92e415231a69cce2e3aba31b88dc01cef46ce028ffbc301fb0c2c",
            "excerpt": "Verdict: pass. No new finding on the uncovered delta 1ecb33ec..8b14208c; CD2-S1 and CD2-S2 still no-op at the wider tip."
          },
          "axis": "Standards",
          "base": "7ce1266321c1a2bd8a974dd34255d161aa665cdf",
          "tip": "8b14208c1b87b0528d54541aa9aa4ad3936fea8d",
          "finding_ids": [],
          "supersedes": ["cd2-standards"]
        },
        {
          "id": "cd2-spec-r3",
          "performer": "claude:bench-reviewer/cd-c2-spec-r3",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "xhigh",
          "source_digest": "a121c0a528a7979eb77ab6b45080d1478a8e34a1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-c2-spec-r3-20260917@8b14208c",
            "digest": "sha256:613e25761097e268486305116a77e6d60eea796dc874b8ba5c968bb9efc92093",
            "excerpt": "Verdict: pass. The CD2b amendment preserves coverage, dependencies, checkpoints, and checks; coverage state mapped with 37 rows; the eight first-pass rows hold at the wider tip."
          },
          "axis": "Spec",
          "base": "7ce1266321c1a2bd8a974dd34255d161aa665cdf",
          "tip": "8b14208c1b87b0528d54541aa9aa4ad3936fea8d",
          "finding_ids": [],
          "supersedes": ["cd2-spec"]
        },
        {
          "id": "cd2-coverage-r3",
          "performer": "claude:bench-reviewer/cd-c2-coverage-r3",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "xhigh",
          "source_digest": "a121c0a528a7979eb77ab6b45080d1478a8e34a1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-c2-coverage-r3-20260917@8b14208c",
            "digest": "sha256:357d6dd1c64e8916d4d2e23227db047199308555ae5ce507b74967ce1a93125e",
            "excerpt": "Verdict: pass after one Writes expansion. CD2b-C1 ticket 7 Writes lacked internal/anchors/match_test.go, confidence 6, applied under the plan-expansion policy in CD2b. CD2b-C2 no-op. No production file changed in the uncovered delta."
          },
          "axis": "Coverage",
          "base": "7ce1266321c1a2bd8a974dd34255d161aa665cdf",
          "tip": "8b14208c1b87b0528d54541aa9aa4ad3936fea8d",
          "finding_ids": [],
          "supersedes": ["cd2-coverage-r2"]
        }
      ]
    },
    {
      "id": "CD2b",
      "base": "67af9d502c9f37c4853f5f72f5684115ca0154a9",
      "tip": "d4ddc08852af81003cd5b5403d67240cb024f603",
      "plan_digest": "sha256:37ef0eb88a3496e3fc048c85e66bd6b4db2ac676d1dca36243cb5e160430ffe7",
      "source_digest": "c34a87699da42aa63a2990b58182b7a2db9b1974",
      "acceptance_rows": ["CR22", "CR38"],
      "verification": [
        {
          "id": "cd2b-workflow",
          "performer": "claude:bench-writer/cd-t7-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "c34a87699da42aa63a2990b58182b7a2db9b1974",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-t7-author-20260917/cd2b-workflow@d4ddc088",
            "digest": "sha256:0391db62814c79d08fdfe7b1bab495b9ef2dd91f9131b17010a593f31ff47b7e",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,815\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "workflow",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "cd2b-anchors",
          "performer": "claude:bench-writer/cd-t7-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "c34a87699da42aa63a2990b58182b7a2db9b1974",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-t7-author-20260917/cd2b-anchors@d4ddc088",
            "digest": "sha256:50be4c99f031aef11b80b3693cd6bbb7a309b191f780769fb790d7ab2ec8aacf",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/anchors,pass,408\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "anchors",
          "command": "bench test --package ./internal/anchors/...",
          "exit_code": 0
        },
        {
          "id": "cd2b-projection",
          "performer": "claude:bench-writer/cd-t7-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "c34a87699da42aa63a2990b58182b7a2db9b1974",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-t7-author-20260917/cd2b-projection@d4ddc088",
            "digest": "sha256:b1014e2ed4a2943d849d57e2c1fd3ff5c5e311c07ca529a4646baad5cd12f63c",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,114\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "projection",
          "command": "bench test --package ./cmd/bench/... --run TestAnchors",
          "exit_code": 0
        },
        {
          "id": "cd2b-bite",
          "performer": "claude:bench-writer/cd-t7-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "c34a87699da42aa63a2990b58182b7a2db9b1974",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-t7-author-20260917/cd2b-bite@d4ddc088",
            "digest": "sha256:45f42a6bb21e567cb797521fab992e63b41e3edb47762a3ace1d2e62ad9af1ab",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,9937\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "bite",
          "command": "bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner",
          "exit_code": 0,
          "probe": {
            "mutation": "move the pickup-confidence sentence from step 6 to step 5",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/cd-t7-author-20260917/cd2b-step-move-probe@d4ddc088",
              "digest": "sha256:24ed1c54fbeef3e273fc63b1a21e53baa2961a55ab96f0016b1e3161b7be2c52",
              "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,fail,901\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:29: gate: calibration: the pickup line must carry its stated confidence\"\nskips[0]{package,test,reason}:\nrestore: cmp exit 0, git status clean, check pass 817 ms"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "cd2b-standards",
          "performer": "claude:bench-reviewer/cd-c2b-standards",
          "role": "independent-review",
          "model": "fable",
          "effort": "medium",
          "source_digest": "c34a87699da42aa63a2990b58182b7a2db9b1974",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/cd-c2b-standards-20260917@779f0e68",
            "digest": "sha256:c7d7263297952e4b17dd0b21798975e6645b8385e7f6fa98b7b338d0dbef438f",
            "excerpt": "Count: 2 hard, 2 judgment. CD2b-S1 stale Locate doc comment, auto-fix. CD2b-S2 change narration in a test comment, auto-fix. CD2b-S3 one why stated three times, auto-fix. CD2b-S4 flag argument on the shared walk, no-op."
          },
          "axis": "Standards",
          "base": "67af9d502c9f37c4853f5f72f5684115ca0154a9",
          "tip": "779f0e68b68e97f5a313da6d6e90c4b778af4052",
          "finding_ids": ["CD2b-S1", "CD2b-S2", "CD2b-S3"],
          "supersedes": []
        },
        {
          "id": "cd2b-standards-r2",
          "performer": "claude:bench-reviewer/cd-c2b-standards-r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "xhigh",
          "source_digest": "c34a87699da42aa63a2990b58182b7a2db9b1974",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-c2b-standards-r2-20260917@da019659",
            "digest": "sha256:c78e6274407391bf9ff1d1c8c857b4e59a24d810d3020ca97e039065cfe51dd1",
            "excerpt": "Verdict: pass. S1, S2, S3 closed at da019659; S4 still no-op; no new finding."
          },
          "axis": "Standards",
          "base": "67af9d502c9f37c4853f5f72f5684115ca0154a9",
          "tip": "da019659879af9706f575e45dd6afffc45338864",
          "finding_ids": [],
          "supersedes": ["cd2b-standards"]
        },
        {
          "id": "cd2b-standards-r3",
          "performer": "claude:bench-reviewer/cd-c2b-standards-r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "xhigh",
          "source_digest": "c34a87699da42aa63a2990b58182b7a2db9b1974",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-c2b-standards-r2-20260917@d4ddc088",
            "digest": "sha256:45d7894b788cbfc04c4b5fd406eeff287daf0b138c278eabce32db4169dd0d5e",
            "excerpt": "Verdict: pass. The cycle-2 test comment obeys craft-comments; no production file changed; prior verdict stands at d4ddc088."
          },
          "axis": "Standards",
          "base": "67af9d502c9f37c4853f5f72f5684115ca0154a9",
          "tip": "d4ddc08852af81003cd5b5403d67240cb024f603",
          "finding_ids": [],
          "supersedes": ["cd2b-standards-r2"]
        },
        {
          "id": "cd2b-spec",
          "performer": "claude:bench-reviewer/cd-c2b-spec",
          "role": "independent-review",
          "model": "fable",
          "effort": "medium",
          "source_digest": "c34a87699da42aa63a2990b58182b7a2db9b1974",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-c2b-spec-20260917@779f0e68",
            "digest": "sha256:300a31885cdfa444561a7c9d534d518f724ca156285b5297bea91269cae841ce",
            "excerpt": "Finding count: 0. CR22 and CR38 held at 779f0e68; the shared narrowing walk is within the ticket."
          },
          "axis": "Spec",
          "base": "67af9d502c9f37c4853f5f72f5684115ca0154a9",
          "tip": "779f0e68b68e97f5a313da6d6e90c4b778af4052",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "cd2b-spec-r2",
          "performer": "claude:bench-reviewer/cd-c2b-spec-r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "xhigh",
          "source_digest": "c34a87699da42aa63a2990b58182b7a2db9b1974",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-c2b-spec-r2-20260917@da019659",
            "digest": "sha256:b202bf6ef667aed34e85e5389d5c8c64791f88672e309f6e5fac2e1a62110dc9",
            "excerpt": "Verdict: pass. CR22 and CR38 hold at da019659; leading-zero rule within the ticket; repair delta comment-only in production."
          },
          "axis": "Spec",
          "base": "67af9d502c9f37c4853f5f72f5684115ca0154a9",
          "tip": "da019659879af9706f575e45dd6afffc45338864",
          "finding_ids": [],
          "supersedes": ["cd2b-spec"]
        },
        {
          "id": "cd2b-spec-r3",
          "performer": "claude:bench-reviewer/cd-c2b-spec-r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "xhigh",
          "source_digest": "c34a87699da42aa63a2990b58182b7a2db9b1974",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-c2b-spec-r2-20260917@d4ddc088",
            "digest": "sha256:8d9e3b4a1a0635af6d77f71c4d871a1042c1ca6de9dd577659f4c931e009b391",
            "excerpt": "Verdict: pass. CR22 and CR38 unchanged at d4ddc088; the cycle-2 test asserts the opener line joins the step body."
          },
          "axis": "Spec",
          "base": "67af9d502c9f37c4853f5f72f5684115ca0154a9",
          "tip": "d4ddc08852af81003cd5b5403d67240cb024f603",
          "finding_ids": [],
          "supersedes": ["cd2b-spec-r2"]
        },
        {
          "id": "cd2b-coverage",
          "performer": "claude:bench-reviewer/cd-c2b-coverage",
          "role": "independent-review",
          "model": "fable",
          "effort": "medium",
          "source_digest": "c34a87699da42aa63a2990b58182b7a2db9b1974",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/cd-c2b-coverage-20260917@779f0e68",
            "digest": "sha256:08cf36e44a0209e2256f0972f6e5e3d0b50c1c3b0051a2571d26d679a50f716c",
            "excerpt": "Finding count: 4. CD2b-C1 step close boundary untested, auto-fix. CD2b-C2 opener trims leading space, auto-fix. CD2b-C3 opener parses one digit, auto-fix. CD2b-C4 keepOpener branch untested, auto-fix."
          },
          "axis": "Coverage",
          "base": "67af9d502c9f37c4853f5f72f5684115ca0154a9",
          "tip": "779f0e68b68e97f5a313da6d6e90c4b778af4052",
          "finding_ids": ["CD2b-C1", "CD2b-C2", "CD2b-C3", "CD2b-C4"],
          "supersedes": []
        },
        {
          "id": "cd2b-coverage-r2",
          "performer": "claude:bench-reviewer/cd-c2b-coverage-r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "xhigh",
          "source_digest": "c34a87699da42aa63a2990b58182b7a2db9b1974",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/cd-c2b-coverage-r2-20260917@da019659",
            "digest": "sha256:130e1f9936d629cf75bc96f8a7886123583c507a2b11de0959d593423cf55b26",
            "excerpt": "Verdict: findings. Three mutations bit; the step arm of keepOpener stayed silent, CD2b-C4 reopened. C2 and C3 confirmed false as stated."
          },
          "axis": "Coverage",
          "base": "67af9d502c9f37c4853f5f72f5684115ca0154a9",
          "tip": "da019659879af9706f575e45dd6afffc45338864",
          "finding_ids": ["CD2b-C4"],
          "supersedes": ["cd2b-coverage"]
        },
        {
          "id": "cd2b-coverage-r3",
          "performer": "claude:bench-reviewer/cd-c2b-coverage-r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "xhigh",
          "source_digest": "c34a87699da42aa63a2990b58182b7a2db9b1974",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-c2b-coverage-r2-20260917@d4ddc088",
            "digest": "sha256:ab53bdc78fddd016dd090eab746c4a52c8675f3663870f6bb18b427d48e39912",
            "excerpt": "Verdict: pass. The keepOpener swap and the call-site swap both bit TestMarkdownNumberedStepsIncludesItsOpener; no unresolved Coverage finding at d4ddc088."
          },
          "axis": "Coverage",
          "base": "67af9d502c9f37c4853f5f72f5684115ca0154a9",
          "tip": "d4ddc08852af81003cd5b5403d67240cb024f603",
          "finding_ids": [],
          "supersedes": ["cd2b-coverage-r2"]
        }
      ]
    },
    {
      "id": "CD3",
      "base": "b47edde2c7a5bc4f00a77ea94e77ac09a6df24ea",
      "tip": "7c3ef89216a5797de657e77ff26670f638260d2a",
      "plan_digest": "sha256:b76aff6080777f4bd95ef7ea7be9b075fb5fd511e8f54b14c92681be2ccbb6a9",
      "source_digest": "4552813d704f9ca1c4506ce0bb1cbb31b2c71e2c",
      "acceptance_rows": ["CR7", "CR8", "CR9", "CR10", "CR11", "CR12", "CR31", "CR36"],
      "verification": [
        {
          "id": "cd3-workflow",
          "performer": "claude:bench-writer/cd-t3-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "4552813d704f9ca1c4506ce0bb1cbb31b2c71e2c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-t3-author-20260917/cd3-workflow@7c3ef892",
            "digest": "sha256:4c0a9e0f2dc8cbfc69e6fc6a0c35605171628e1c8ffaed57583a2d3697e29494",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,1019\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "workflow",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "cd3-budgets",
          "performer": "claude:bench-writer/cd-t3-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "4552813d704f9ca1c4506ce0bb1cbb31b2c71e2c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-t3-author-20260917/cd3-budgets@7c3ef892",
            "digest": "sha256:6f72c99a424331a3381e284e55ffa703b490f81e19b498e848e8bfc5d9ca25ed",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,5\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        },
        {
          "id": "cd3-prose",
          "performer": "claude:bench-writer/cd-t3-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "4552813d704f9ca1c4506ce0bb1cbb31b2c71e2c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-t3-author-20260917/cd3-prose@7c3ef892",
            "digest": "sha256:e5d1227f4d01a10308bf558fe195080d0e0062d4e541dea2c214e52b944ab24d",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,160\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "prose",
          "command": "bench test --check prose-mechanics",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "cd3-standards",
          "performer": "claude:bench-reviewer/cd-c3-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "4552813d704f9ca1c4506ce0bb1cbb31b2c71e2c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-c3-standards-20260917@7c3ef892",
            "digest": "sha256:34e2da718790a0e2c14e43ce73b0e9cde19ee037a875ea0b57e305bba6753649",
            "excerpt": "Findings: 2 judgment calls. CD3-S1 the reference pointer sits inside the declaration template line, confidence 6, decided by the spec's pasted needle, no-op. CD3-S2 the folded line carries four sentences, confidence 5, no-op. One source per fact holds."
          },
          "axis": "Standards",
          "base": "b47edde2c7a5bc4f00a77ea94e77ac09a6df24ea",
          "tip": "7c3ef89216a5797de657e77ff26670f638260d2a",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "cd3-spec",
          "performer": "claude:bench-reviewer/cd-c3-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "4552813d704f9ca1c4506ce0bb1cbb31b2c71e2c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-c3-spec-20260917@7c3ef892",
            "digest": "sha256:0abc90453e415b789e3ec01c30296d157de6a761ec27cbaf6b8a027799098e47",
            "excerpt": "Finding count: 0. Every CD3 row held: CR7, CR8, CR9, CR10, CR11, CR12, CR31, CR36. The fold of the fan-out line is within the spec's decision."
          },
          "axis": "Spec",
          "base": "b47edde2c7a5bc4f00a77ea94e77ac09a6df24ea",
          "tip": "7c3ef89216a5797de657e77ff26670f638260d2a",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "cd3-coverage",
          "performer": "claude:bench-reviewer/cd-c3-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "4552813d704f9ca1c4506ce0bb1cbb31b2c71e2c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-c3-coverage-20260917@7c3ef892",
            "digest": "sha256:32034e96dc99643e8ec3337a929881087a78b9ea2706b841da8482109745e6fb",
            "excerpt": "Finding count: 0 blocking. CD3-C1 an appended fourth label source stays green, confidence 3, Won't handle per the spec, no-op. Eight probes: five bit, three silent on matcher-wide semantics."
          },
          "axis": "Coverage",
          "base": "b47edde2c7a5bc4f00a77ea94e77ac09a6df24ea",
          "tip": "7c3ef89216a5797de657e77ff26670f638260d2a",
          "finding_ids": [],
          "supersedes": []
        }
      ]
    }
  ],
  "completion": {
    "state": "pending",
    "source_digest": "",
    "performer": "",
    "reconciliation": {},
    "verification": []
  }
}
```
