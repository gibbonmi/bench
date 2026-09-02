# State the fence order and the claim words in craft-spec

Blocked by: derive-test-wait-deadlines-from-bounds.md
Writes: .agents/skills/bench-craft-spec/SKILL.md, .agents/skills/bench-craft-spec/references/map-discipline.md, .agents/skills/bench-craft-spec/references/ste-prose.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/conformance/guidance_token_sweep_test.go, tests/canary/workflow-guidance-anchors/craft-spec-fence-after-slice (new), tests/canary/workflow-guidance-anchors/map-discipline-shipped-surface-claim (new), tests/canary/workflow-guidance-anchors/ste-prose-two-part-label-rule (new), projects/benchkit.md, tests/canary/workflow-guidance-anchors/bootstrap-authority-before-after-softening, tests/canary/workflow-guidance-anchors/bootstrap-authority-rule-deletion, tests/canary/workflow-guidance-anchors/craft-spec-empty-or-invalid-fence, tests/canary/workflow-guidance-anchors/craft-spec-exact-literal-fence, tests/canary/workflow-guidance-anchors/craft-spec-further-notes-bare, tests/canary/workflow-guidance-anchors/craft-spec-reader-sweep-name, tests/canary/workflow-guidance-anchors/craft-spec-red-signal-grammar-retired, tests/canary/workflow-guidance-anchors/craft-spec-template-verification-log, tests/canary/workflow-guidance-anchors/ticket-cross-pointers-anchor, tests/canary/workflow-guidance-anchors/write-spec-acceptance-round-exit, tests/canary/workflow-guidance-anchors/write-spec-approval-table-emit, tests/canary/workflow-guidance-anchors/write-spec-contract-rederivation, tests/canary/workflow-guidance-anchors/write-spec-degenerate-cheapness, tests/canary/workflow-guidance-anchors/write-spec-fence-approval, tests/canary/workflow-guidance-anchors/write-spec-identified-coverage-default, tests/canary/workflow-guidance-anchors/write-spec-materiality-exit, tests/canary/workflow-guidance-anchors/write-spec-ownership-fences, tests/canary/workflow-guidance-anchors/write-spec-promise-inflation-guard, tests/canary/workflow-guidance-anchors/write-spec-unflagged-addition-removal, tests/canary/workflow-guidance-anchors/write-spec-unique-row-id, tests/canary/workflow-guidance-anchors/craft-spec-transaction-failure-rows, tests/canary/workflow-guidance-anchors/craft-spec-two-audience-inventory, tests/canary/workflow-guidance-anchors/map-discipline-addition-disposition, tests/canary/workflow-guidance-anchors/map-discipline-either-side-rows, tests/canary/workflow-guidance-anchors/map-discipline-excluded-edge-caller, tests/canary/workflow-guidance-anchors/map-discipline-executed-root-trace, tests/canary/workflow-guidance-anchors/map-discipline-fixture-reachable-state, tests/canary/workflow-guidance-anchors/map-discipline-flagged-additions, tests/canary/workflow-guidance-anchors/map-discipline-moved-bytes-sweep, tests/canary/workflow-guidance-anchors/map-discipline-promise-rows, tests/canary/workflow-guidance-anchors/map-discipline-quoted-operands, tests/canary/workflow-guidance-anchors/map-discipline-source-sentence-table, tests/canary/workflow-guidance-anchors/map-discipline-sweep-depth-bound, tests/canary/workflow-guidance-anchors/map-discipline-sweep-direct-helpers, tests/canary/workflow-guidance-anchors/map-discipline-sweep-named-consumers, tests/canary/workflow-guidance-anchors/map-discipline-sweep-reader-fence, tests/canary/workflow-guidance-anchors/reader-sweep-term, tests/canary/workflow-guidance-anchors/ste-prose-paragraph-bound, tests/canary/workflow-guidance-anchors/ste-prose-sentence-bound, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go, tests/canary/guidance-prose-budgets/over-budget-skill, tests/canary/line-routing/line-binding-prose-drift, tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading, tests/canary/workflow-guidance-anchors/benchkit-review-round-owner, tests/canary/workflow-guidance-anchors/benchkit-review-round-routing, tests/canary/workflow-guidance-anchors/benchkit-spec-ownership
Covers: LP4, LP14, LP15

## What to build

Ship three sentences in the five-part precedent shape. Each sentence carries
one anchor tuple, one red-on-removal registry test, and one live-mirror fixture
with `BASE`, `MUTATE.json`, and `EXPECT`. No guard row applies. Mirror
`TestCraftGateBothEndsAnchorsRedOnRemoval` and the fixture
`tests/canary/workflow-guidance-anchors/final-check-scratch-branch-clean`.

In `SKILL.md`, under "Slicing a build for delegates", add two rules. The fence
section is written after the ticket slice from the union of the tickets'
`Writes:` lines. A Won't handle over an anchored sentence quotes the bytes it
keeps. In `map-discipline.md`, under "Before the map locks", say the reader
sweep names the shipped-surface claim words. In `ste-prose.md`, extend the
field-line sentence with the template field-name clause.

Keep each needle on one physical line. Keep `SKILL.md` inside its 152-line
budget row. Move the row in `projects/benchkit.md` only if the count must rise.

## Acceptance

- [ ] The registry test reds a synthetic tree that drops each of the three sentences and stays silent on the live root.
- [ ] Each of the three fixtures bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] The guidance token sweep and the prose budget check pass on the worktree.
- [ ] Self-probe: reword one needle in the live file, and report the registry test red with that diagnostic.
