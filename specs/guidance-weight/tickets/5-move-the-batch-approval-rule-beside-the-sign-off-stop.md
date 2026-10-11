# Move the batch-approval rule beside the sign-off stop

Blocked by: 4-merge-the-light-path-rules-into-one-block.md
Writes: .agents/commands/bench-drain.md, .agents/commands/bench-write-spec.md, .agents/skills/bench-craft-tickets/SKILL.md, .bench/BENCH.md, CHANGELOG.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/anchors/anchor_harness_diagnostics_test.go, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_commitment.go, internal/anchors/registry_craft_research.go, internal/anchors/registry_craft_research_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_decision_maps.go, internal/anchors/registry_decision_maps_test.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/commitment_guidance_test.go, internal/conformance/subcommand_routing_table_test.go, projects/benchkit.md, tests/canary/docs-currency-token-diet/benchref-imported, tests/canary/docs-currency-token-diet/benchref-pointer-dropped, tests/canary/docs-currency-token-diet/benchref-section-duplicated, tests/canary/docs-currency-token-diet/dogfood-referent-shipped, tests/canary/docs-currency-token-diet/missing-cli-inventory, tests/canary/docs-currency-token-diet/stale-cli-doc-reference, tests/canary/guidance-prose-budgets/over-budget-skill, tests/canary/line-routing/line-binding-prose-drift, tests/canary/load-validity-metadata/readme-shared-rule-drift, tests/canary/load-validity-metadata/shared-rule-drift, tests/canary/row-next-grammar/token-table-lacks-kit-edit, tests/canary/skill-description-budgets/budget-table-missing, tests/canary/skill-description-budgets/description-folded, tests/canary/skill-description-budgets/description-missing, tests/canary/skill-description-budgets/over-budget-command, tests/canary/skill-description-budgets/over-budget-description, tests/canary/skills-index-command-adapters/adapter-inert-invocation-key, tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy, tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted, tests/canary/workflow-guidance-anchors/acceptance-coverage-anchor, tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading, tests/canary/workflow-guidance-anchors/benchkit-review-round-owner, tests/canary/workflow-guidance-anchors/benchkit-review-round-routing, tests/canary/workflow-guidance-anchors/benchkit-spec-ownership, tests/canary/workflow-guidance-anchors/benchkit-system-suite-route, tests/canary/workflow-guidance-anchors/capture-sink-anchor, tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns, tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary, tests/canary/workflow-guidance-anchors/command-handoff-anchor, tests/canary/workflow-guidance-anchors/core-absent-bench-operational-layer, tests/canary/workflow-guidance-anchors/core-absent-category-context, tests/canary/workflow-guidance-anchors/core-absent-category-oracle, tests/canary/workflow-guidance-anchors/core-absent-category-setup, tests/canary/workflow-guidance-anchors/core-absent-category-work, tests/canary/workflow-guidance-anchors/core-absent-gate-authority, tests/canary/workflow-guidance-anchors/core-absent-kit-only-ship, tests/canary/workflow-guidance-anchors/core-absent-no-path-fallback, tests/canary/workflow-guidance-anchors/core-absent-retained-integration-source, tests/canary/workflow-guidance-anchors/core-absent-retro-capture-owner, tests/canary/workflow-guidance-anchors/core-absent-retro-drain-owner, tests/canary/workflow-guidance-anchors/core-absent-skills-guidance, tests/canary/workflow-guidance-anchors/core-absent-upgrade-route, tests/canary/workflow-guidance-anchors/core-cli-bounded-complete-response, tests/canary/workflow-guidance-anchors/core-cli-call-pair, tests/canary/workflow-guidance-anchors/core-progressive-loading-term, tests/canary/workflow-guidance-anchors/craft-tickets-slice-acceptance-row, tests/canary/workflow-guidance-anchors/decision-map-asset-path, tests/canary/workflow-guidance-anchors/delegated-chunk-tip-review, tests/canary/workflow-guidance-anchors/delegated-green-predecessor, tests/canary/workflow-guidance-anchors/delegated-opt-in-entry, tests/canary/workflow-guidance-anchors/delegated-pending-predecessor-stop, tests/canary/workflow-guidance-anchors/delegated-prerequisite-checkpoint, tests/canary/workflow-guidance-anchors/delegated-repair-ownership, tests/canary/workflow-guidance-anchors/delegated-serial-ticket-checkpoint, tests/canary/workflow-guidance-anchors/dg-17, tests/canary/workflow-guidance-anchors/dg-18, tests/canary/workflow-guidance-anchors/dg-19, tests/canary/workflow-guidance-anchors/dg-20, tests/canary/workflow-guidance-anchors/dg-21, tests/canary/workflow-guidance-anchors/dg-22, tests/canary/workflow-guidance-anchors/dg-23, tests/canary/workflow-guidance-anchors/drain-anchor, tests/canary/workflow-guidance-anchors/drain-implement-now-commit, tests/canary/workflow-guidance-anchors/drain-implement-now-per-spec-exception, tests/canary/workflow-guidance-anchors/drain-implement-now-route, tests/canary/workflow-guidance-anchors/drain-implement-now-row-fallback, tests/canary/workflow-guidance-anchors/drain-implement-now-second-exception, tests/canary/workflow-guidance-anchors/drain-roadmap-context-anchor, tests/canary/workflow-guidance-anchors/drain-spec-history-anchor, tests/canary/workflow-guidance-anchors/drain-split-board-detail-owner, tests/canary/workflow-guidance-anchors/drain-split-board-retirement-pair, tests/canary/workflow-guidance-anchors/drain-split-board-row-detail-owner, tests/canary/workflow-guidance-anchors/edge-inventory-anchor, tests/canary/workflow-guidance-anchors/fix-dont-park-section-relocated, tests/canary/workflow-guidance-anchors/implementation-retro-drain-anchor, tests/canary/workflow-guidance-anchors/light-path-inline, tests/canary/workflow-guidance-anchors/slicing-in-write-spec, tests/canary/workflow-guidance-anchors/spec-retire-roadmap-row, tests/canary/workflow-guidance-anchors/staged-command-sweep-anchor, tests/canary/workflow-guidance-anchors/story-line-anchor-missing, tests/canary/workflow-guidance-anchors/structured-phase-progress-anchor, tests/canary/workflow-guidance-anchors/ticket-decision-map-lifecycle-anchor, tests/canary/workflow-guidance-anchors/ticket-executable-route-pass, tests/canary/workflow-guidance-anchors/ticket-light-path-anchor, tests/canary/workflow-guidance-anchors/ticket-light-path-carve-out, tests/canary/workflow-guidance-anchors/ticket-lock-passes, tests/canary/workflow-guidance-anchors/ticket-relocation-destinations, tests/canary/workflow-guidance-anchors/ticket-seam-creating-slice, tests/canary/workflow-guidance-anchors/ticket-skill-contract-anchor, tests/canary/workflow-guidance-anchors/ticket-source-clause-pass, tests/canary/workflow-guidance-anchors/ticket-template-anchor, tests/canary/workflow-guidance-anchors/write-spec-artifact-authorization, tests/canary/workflow-guidance-anchors/write-spec-authorization-boundary, tests/canary/workflow-guidance-anchors/write-spec-branch-split, tests/canary/workflow-guidance-anchors/write-spec-conversation-authorization, tests/canary/workflow-guidance-anchors/write-spec-conversation-fork, tests/canary/workflow-guidance-anchors/write-spec-decision-source, tests/canary/workflow-guidance-anchors/write-spec-frozen-base-and-tip-review, tests/canary/workflow-guidance-anchors/write-spec-handoff-anchor, tests/canary/workflow-guidance-anchors/write-spec-late-uncertainty, tests/canary/workflow-guidance-anchors/write-spec-map-sources, tests/canary/workflow-guidance-anchors/write-spec-phase-ownership, tests/canary/workflow-guidance-anchors/write-spec-post-slicing-handoff, tests/canary/workflow-guidance-anchors/write-spec-reader-sweep-sequence, tests/canary/workflow-guidance-anchors/write-spec-ready-map-authorization, tests/canary/workflow-guidance-anchors/write-spec-review-made-conditional, tests/canary/workflow-guidance-anchors/write-spec-review-tier-escalated, tests/canary/workflow-guidance-anchors/write-spec-slicing-sign-off, tests/canary/workflow-guidance-anchors/write-spec-ticket-approval-fields, tests/canary/workflow-guidance-anchors/write-spec-ticket-breakdown-charge, tests/canary/workflow-guidance-anchors/write-spec-verification-learning, tests/canary/workflow-guidance-anchors/write-spec-verification-learning-contents, tests/canary/workflow-guidance-anchors/write-spec-verification-log
Covers: GW32, GW33, GW34, GW35, GW36, GW37, GW51

## What to build

A spec stage reads the batch-approval rule B43 where it acts, beside the sign-off stop of `/bench-write-spec`.
A session that never signs off a spec no longer loads that rule from `.bench/BENCH.md`.
`craft-tickets` points to `/bench-write-spec` for the batch-approval carve-out.
`/bench-drain` names its own diff approval in its own words, so "batch approval" keeps one meaning.
Each anchor row that quotes an edited line takes the new text, so it stays green and still bites.

This ticket reads three contracts from the earlier chunks.
Ticket `1-move-the-build-rules-into-their-own-file.md` added the section `### Always-loaded guidance weight` to `CHANGELOG.md`, and tickets 3 and 4 added the GW2 and GW3 bullets.
Ticket `4-merge-the-light-path-rules-into-one-block.md` left the B43 paragraph of `.bench/BENCH.md` as it is, because this chunk moves it.
Ticket `2-repoint-the-build-rule-pointers.md` changed only line 50 of `.agents/commands/bench-write-spec.md` and lines 67 and 95 of `.agents/skills/bench-craft-tickets/SKILL.md`.
Line citations below are at base `bc2ea0ff`.
If a landing moves a line, find the text by its quotation.

**The move.** Remove the B43 paragraph from `## Workflow` of `.bench/BENCH.md`, with its blank line.
At base it is lines 156-159, and it reads as below:

```markdown
**A batch approval covers per-spec sign-offs when I'm unreachable.** If I
approved a batch plan and went AFK, build on rather than stall. Leave each
spec in `specs/` as post-hoc veto surface and flag contestable calls. Without
a batch approval, spec sign-off is a hard stop.
```

Put the same four physical lines, byte for byte, into `## Exit handoff` of `.agents/commands/bench-write-spec.md`.
They are their own paragraph, directly after line 16, with one blank line on each side.
Lines 15 and 16 are one paragraph that holds "Stop for sign-off." and ends at "`bench commitment start`.", and that paragraph stays byte-identical.
The rule keeps its reviewer voice ("I", "I'm"), as `.agents/commands/bench-debug.md:89-90` does.

`internal/conformance/commitment_guidance_test.go:69` requires the line-16 sentence that opens " A staged spec is planning work" once, with its leading space.
The row `internal/anchors/registry_commitment.go:42` pins the same sentence.
The new paragraph starts after line 16, so both stay green with no edit.
If the close-order block of `specs/spec-stage-binding` ticket 3 has landed after the sign-off paragraph, B43 still goes directly after that paragraph, before the block.

**The line budget.** `.agents/commands/bench-write-spec.md` has 87 lines, and its row in `projects/benchkit.md:492` is 87.
The move adds four lines and one blank line.
In the same commit, raise the row by those five lines, so at base it goes from 87 to 92.
`.bench/BENCH.md` loses five lines, so its row at `projects/benchkit.md:490` stays green.
The `guidance-prose-budgets` lane runs because the profile changes, and this ticket also runs the check itself.

**The craft-tickets pointer.** Line 37 of `.agents/skills/bench-craft-tickets/SKILL.md` becomes this one physical line:

```text
For spec-backed builds, this is the only route onto the frontier; the batch-approval AFK carve-out in `/bench-write-spec` is the sole no-round-trip exception.
```

Line 38 stays byte-identical.

**The drain line.** The second sentence of `.agents/commands/bench-drain.md:291` changes, and the first sentence stays.
The line stays one physical line, and it then reads as below:

```text
Never commit the drain without that approval. That approval is the reviewer's sign-off on this drain's proposed diff.
```

**The anchor rows.** Make these edits, and keep each row on its physical line:

| row | edit |
|---|---|
| `internal/anchors/registry_data.go:207` | The needle opens with the new line-37 sentence. A space and the unchanged line-38 text follow it. The diagnostic stays. |
| `internal/anchors/registry_retained_workflow.go:29` | The needle becomes "That approval is the reviewer's sign-off on this drain's proposed diff.". The diagnostic becomes "reference route: drain dropped the reviewer's sign-off on its own proposed diff". |
| `internal/anchors/registry_retained_workflow.go:30` | In the diagnostic, "the operating guide's batch approval rule" becomes "/bench-write-spec's batch approval rule". |
| `internal/anchors/registry_retained_workflow.go:31` | In the diagnostic, "that the operating guide owns" becomes "that /bench-write-spec owns". |

The forbid needles at `internal/anchors/registry_retained_workflow.go:28`, `:30`, and `:31` stay, and so does the forbid row `internal/anchors/registry_data.go:182`.
The diagnostic of `:28` names `AGENTS.md`, not the operating guide, so it stays.
In the doc comment of `referenceRouteAnchors` (`internal/anchors/registry_retained_workflow.go:16-25`), "the drain's batch-approval owner" becomes "the drain's own diff approval".
Each diagnostic keeps the prefix `reference route: `, and no count pin reads that prefix.

**The canaries.** No canary `EXPECT` holds a diagnostic of `internal/anchors/registry_retained_workflow.go:28-31`.
Three overlays hold old text, and each still bites its own row:

- `ticket-light-path-carve-out/files/dot-agents/skills/bench-craft-tickets/SKILL.md:6-8` lacks the new needle. Its `EXPECT` holds the diagnostic of `internal/anchors/registry_data.go:207`, which stays, and the live file clears it on restoration.
- `capture-sink-anchor/files/dot-bench/BENCH.md:154-157` holds an old B43. Its `EXPECT` names the Capture sink row, which this ticket does not change.
- `drain-anchor/files/dot-agents/commands/bench-drain.md:60-62` holds an old drain approval text. Its `EXPECT` names the "Reconcile first" anchor.

No overlay needs an edit.
If a canary reds, restore its bite with an overlay edit, and report it.

**The readers that stay.** `internal/conformance/fixture_bite_test.go:368-369` quotes lines 61-62 and 68-69 of `bench-write-spec.md` byte for byte, and the move does not change that text.
`internal/conformance/fixture_bite_test.go:371` quotes line 35 of the craft-tickets skill, which stays.

**The changelog.** Add one bullet for chunk GW4 to `### Always-loaded guidance weight` in `CHANGELOG.md`, after the GW3 bullet.
This is the proposed text:

```markdown
- Changed `.bench/BENCH.md` so that it no longer holds the batch-approval rule. The rule moved word for word into `/bench-write-spec`, directly after its stop for spec sign-off, and `craft-tickets` now names `/bench-write-spec` as the home of that carve-out. `/bench-drain` now says that its approval is the reviewer's sign-off on the drain's proposed diff. Before, every session loaded the rule, and the drain sent its own diff approval to that rule.
```

**The structure limits.** `internal/anchors/registry_data.go` has 478 lines, over the 400-line cap with no grant (`internal/structure/structure.go:308-331`).
Its row edit stays on its physical line, so it gains no line.
`internal/anchors/registry_retained_workflow.go` has 396 lines, and its edits stay in place.

**The Writes line.** The line holds each research path and the five bound files of `internal/anchors` and `cmd/bench` (`internal/tickets/registry_data.go:21-27`, `:45`).
It holds the 14 anchor registry files that hold a literal equal to an edited guidance path (`internal/anchors/references.go:31-59`).
It also holds the 113 fixture directories that pin a path of this line through `BASE`, an overlay, or a mutation.

Ticket 1 adds the GW03 row and its test to `internal/anchors/registry_chunk_chain.go` and `internal/anchors/registry_chunk_chain_test.go`.
Both files can then hold the literal `.bench/BENCH.md`, so both join the line.
`internal/conformance/commitment_guidance_test.go` stays in the line, because its case at `:69` quotes the paragraph that B43 follows.
This design edits five guidance files, `CHANGELOG.md`, and two registry files, and the other paths are the approved upper bound.
This ticket drops no research path.

## Acceptance

- [ ] `bench-write-spec.md` holds B43 byte for byte directly after the sign-off paragraph, which the GW4 review compares with the base (GW32).
- [ ] `rg -F "batch approval covers" .bench/BENCH.md` finds no line (GW33).
- [ ] Line 37 of the craft-tickets skill names `/bench-write-spec`, and the canary `ticket-light-path-carve-out` bites the row `internal/anchors/registry_data.go:207` in `TestEveryRetainedFixtureBitesThroughRegisteredOwner` (GW34).
- [ ] `.agents/commands/bench-drain.md:291` names the reviewer's sign-off on the drain's proposed diff, and the row `internal/anchors/registry_retained_workflow.go:29` passes in `TestRootConformance` with the new needle (GW35).
- [ ] The forbid rows `internal/anchors/registry_retained_workflow.go:28`, `:30`, and `:31` keep their needles, `TestRootConformance` passes, and each restored needle raises its diagnostic (GW36).
- [ ] The row of `.agents/commands/bench-write-spec.md` rises by the five lines that B43 adds, and `TestGuidanceProseBudgetsHoldOnTheLiveTree` passes (GW37).
- [ ] The diagnostics of `internal/anchors/registry_retained_workflow.go:30` and `:31` name `/bench-write-spec` as the owner, which the GW4 review reads for the owner-rename rule.
- [ ] `bench test --check prose-mechanics` passes with B43 in `.agents/commands/bench-write-spec.md`.
- [ ] The section `### Always-loaded guidance weight` of `CHANGELOG.md` holds the GW4 bullet after the GW3 bullet (GW51).

## Verification

Record each red in the return, with its command and its failure line.
Run `bench worktree build <target>` before each named-check run that follows a Go edit, because a named check compiles from the run binary's source.

- Move B43 before the row edit, and run `bench test --check guidance-prose-budgets`. The check reds with `prose-budget exceeded: .agents/commands/bench-write-spec.md is 92 lines, over its 87-line budget`. Then raise the row, and record the green run.
- Change the drain line before the row `internal/anchors/registry_retained_workflow.go:29`, and run `bench test --check docs-currency-workflow`. The check reds with the base diagnostic of that row. Then change the row, and record the green run.
- Change line 37 of the craft-tickets skill before the row `internal/anchors/registry_data.go:207`, and run `bench test --check docs-currency-workflow`. The check reds with the diagnostic of that row. Then change the row, and record the green run.

After the edit, record one observed red through `bench probe` for each executable row:

- GW35: swap the new drain sentence for "`.bench/BENCH.md` owns the batch approval rule." with `--check docs-currency-workflow`. The check reds with the new diagnostic of `:29`.
- GW36: swap "Never commit the drain without that approval." for "Never commit the drain without that approval. A batch approval covers per-spec sign-offs when I'm unreachable." with `--check docs-currency-workflow`. The check reds with the renamed diagnostic of `:30`.
- GW36: swap the same sentence for "Never commit the drain without that approval; a standing batch approval counts, with contestable calls flagged for post-hoc veto." with `--check docs-currency-workflow`. The check reds with the renamed diagnostic of `:31`.
- GW36: swap the same sentence for "Never commit the drain without that approval (the AGENTS.md rule)." with `--check docs-currency-workflow`. The check reds with the diagnostic of `:28`.
- GW34: swap "carve-out in `/bench-write-spec` is the sole" for "carve-out in `.bench/BENCH.md` is the sole" in the craft-tickets skill with `--check docs-currency-workflow`. The check reds with the diagnostic of `internal/anchors/registry_data.go:207`.
- GW37: lower the row of `.agents/commands/bench-write-spec.md` in `projects/benchkit.md` by one, with `--check guidance-prose-budgets`. The check reds with the over-budget diagnostic of that file.

Each probe reports a restored subject.
For GW32, run `rg -n -F` for each B43 line over `.agents/commands/bench-write-spec.md`.
Record that the four lines run in order after the blank line that follows line 16.
For GW33, run the `rg` search of its acceptance row, and record that it finds no line.

B43 holds four sentences, the bold label included, and its longest sentence holds 14 words, inside the bounds of `internal/prose/prose.go:16-19`.
`prose-mechanics` grades only those two bounds, so the reviewer voice does not red it.

Then run `bench test --package ./internal/anchors --package ./internal/conformance --package ./cmd/bench`.
Also run the checks `docs-currency-workflow`, `prose-mechanics`, `guidance-prose-budgets`, and `root-conformance`.
