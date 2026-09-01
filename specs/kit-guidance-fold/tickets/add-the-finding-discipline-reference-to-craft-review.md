# Add the finding discipline reference to craft-review

Blocked by: fold-the-probe-oracle-and-fence-rules-into-delegation-discipline.md
Writes: .agents/skills/bench-craft-review/SKILL.md, .agents/skills/bench-craft-review/references/finding-discipline.md (new), projects/benchkit.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, tests/canary/workflow-guidance-anchors, tests/canary/guidance-prose-budgets/over-budget-skill, tests/canary/line-routing/line-binding-prose-drift, tests/canary/workflow-guidance-anchors/review-string-expectation-catch (new), tests/canary/workflow-guidance-anchors/review-citation-location (new), tests/canary/workflow-guidance-anchors/review-test-deletion-coverage (new), tests/canary/workflow-guidance-anchors/review-strong-finding-run (new), tests/canary/workflow-guidance-anchors/review-env-var-producer (new), tests/canary/workflow-guidance-anchors/review-seam-amendment (new), tests/canary/workflow-guidance-anchors/review-finding-discipline-pointer (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: KG30, KG31, KG32, KG33, KG34, KG35, KG36, KG37, KG38, KG39, KG40

## What to build

A review axis reads six finding rules in a new reference file, and it finds that file from
the skill. The coordinator then pays no read to dismiss a refuted finding.

Create `.agents/skills/bench-craft-review/references/finding-discipline.md`. Open it with
a charge-time lead sentence, in the shape the two existing discipline references use. End
the file with a newline. Carry no fence block and no comment block in the file, so no
unterminated delimiter can swallow it.

Six rules land in that reference file. The first rule states that a generated script's
independently authored string expectation is the mutation catch. The second rule requires
a finding to cite the line the axis read this pass, or the symbol instead. The third rule
makes a test-deleting Standards finding name the surviving assertion, or file as coverage.
The fourth rule requires a real run before an axis reports a strong finding. The fifth
rule makes an environment-variable Coverage finding cite the producer before it claims
absence.

The sixth rule disposes an unreachable row seam as an amendment of the row's seam column,
and it records the helper seam as a decision.

The mutations fix the needle words. Row KG30 weakens `is the mutation catch` to a weaker
signal. Row KG31 drops the symbol arm. Row KG32 softens `names` to `may name`. Row KG33
drops `with a real run`. Row KG34 drops the before-the-claim clause, and row KG35 swaps
`amends` for `notes`.

Hold each of those words inside its needle.

One pointer sentence lands under `What a finding must cite` in
`.agents/skills/bench-craft-review/SKILL.md`. The sentence points at
`references/finding-discipline.md`. Give it one `RequireInSection` tuple over that
section. The `Refute before you report` section keeps one refute rule, and it takes no
copy of the real-run clause; the reference alone holds that clause.

Leave `.agents/commands/bench-review-implementation.md` untouched. That file already
states that a finding cites what its axis read now, and the fixture
`review-universal-claim-bar` pins line 84. The location clause lands in the reference
alone.

Add the exact budget row `| .agents/skills/bench-craft-review/SKILL.md | 122 |` to the
guidance prose budget table in `projects/benchkit.md`. The skill file sits at 120 of 120
today, and the pointer sentence needs one more line. The new row joins the existing exact
skill rows, and an exact row beats the glob row of 120.

Give each of the six rules one tuple in `internal/anchors/registry_data.go`. Every tuple
takes `AfterImplementSpec`, the group the newest tuples in that file use. Keep each needle
on one physical line, because a needle that wraps across two lines never matches.

The reference file has four H2 sections. They are `What a string expectation proves` for
rule 1, `What a citation points at` for rule 2, `Where an axis under-reads` for rules 3,
4, and 5, and `When a seam cannot reach the state` for rule 6. Each rule tuple takes
`RequireInSection` over its section. The lead sentence sits above the first section, and
its tuple takes `Require` with no section.

Add `TestCraftReviewFindingDisciplineAnchorsRedOnRemoval` to
`internal/anchors/registry_data_test.go`, and give it rows KG30 through KG35 and row KG37.
Extend the existing `TestReferenceFileAnchorsRedOnAbsence` with the new reference file and
its lead sentence. Each function writes its needles and its diagnostics independently of
the registry. Each function proves a red on a synthetic tree that drops the needle or the
file, and silence on the live root.

Add seven fixtures under `tests/canary/workflow-guidance-anchors/` in the live-mirror
shape. `BASE` names the live file path the tuple reads. `MUTATE.json` holds one
old-and-new needle swap. `EXPECT` holds the diagnostic the tuple raises. Do not use the
older `files/` snapshot shape for any fixture.

This ticket appends its tuples and its new test function to the two registry files. It
edits nothing the previous four tickets added, except the new reference entry inside the
existing absence test function.

## Acceptance

- [ ] KG30 — the fixture `review-string-expectation-catch` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG31 — the fixture `review-citation-location` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG32 — the fixture `review-test-deletion-coverage` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG33 — the fixture `review-strong-finding-run` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG34 — the fixture `review-env-var-producer` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG35 — the fixture `review-seam-amendment` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG36 — `TestReferenceFileAnchorsRedOnAbsence` reds on a tree without the new reference file, and stays silent on the live root.
- [ ] KG37 — the fixture `review-finding-discipline-pointer` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG38 — `TestGuidanceProseBudgetsHoldOnTheLiveTree` stays green with `craft-review` SKILL.md at 122 lines under its new exact row.
- [ ] KG39 — the fixture `review-universal-claim-bar` stays green, because line 84 of the review phase file keeps its bytes.
- [ ] KG40 — the `Refute before you report` section holds one refute rule and no copy of the real-run clause.
