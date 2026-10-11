# Merge the light-path rules into one block under the right-size table

Blocked by: 3-remove-the-repeated-copies.md
Writes: .bench/BENCH.md, CHANGELOG.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/anchors/anchor_harness_diagnostics_test.go, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_commitment.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, internal/conformance/axi_query_registry_test.go, internal/conformance/commitment_guidance_test.go, internal/conformance/subcommand_routing_table_test.go, internal/conformance/validity_checks_test.go, tests/canary/docs-currency-token-diet/benchref-imported, tests/canary/docs-currency-token-diet/benchref-pointer-dropped, tests/canary/docs-currency-token-diet/benchref-section-duplicated, tests/canary/docs-currency-token-diet/dogfood-referent-shipped, tests/canary/docs-currency-token-diet/missing-cli-inventory, tests/canary/docs-currency-token-diet/stale-cli-doc-reference, tests/canary/load-validity-metadata/readme-shared-rule-drift, tests/canary/load-validity-metadata/shared-rule-drift, tests/canary/skills-index-command-adapters/adapter-inert-invocation-key, tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy, tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted, tests/canary/workflow-guidance-anchors/capture-sink-anchor, tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns, tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary, tests/canary/workflow-guidance-anchors/core-absent-bench-operational-layer, tests/canary/workflow-guidance-anchors/core-absent-category-context, tests/canary/workflow-guidance-anchors/core-absent-category-oracle, tests/canary/workflow-guidance-anchors/core-absent-category-setup, tests/canary/workflow-guidance-anchors/core-absent-category-work, tests/canary/workflow-guidance-anchors/core-absent-gate-authority, tests/canary/workflow-guidance-anchors/core-absent-kit-only-ship, tests/canary/workflow-guidance-anchors/core-absent-no-path-fallback, tests/canary/workflow-guidance-anchors/core-absent-retained-integration-source, tests/canary/workflow-guidance-anchors/core-absent-retro-capture-owner, tests/canary/workflow-guidance-anchors/core-absent-retro-drain-owner, tests/canary/workflow-guidance-anchors/core-absent-skills-guidance, tests/canary/workflow-guidance-anchors/core-absent-upgrade-route, tests/canary/workflow-guidance-anchors/core-cli-bounded-complete-response, tests/canary/workflow-guidance-anchors/core-cli-call-pair, tests/canary/workflow-guidance-anchors/core-progressive-loading-term, tests/canary/workflow-guidance-anchors/delegated-chunk-tip-review, tests/canary/workflow-guidance-anchors/delegated-green-predecessor, tests/canary/workflow-guidance-anchors/delegated-opt-in-entry, tests/canary/workflow-guidance-anchors/delegated-pending-predecessor-stop, tests/canary/workflow-guidance-anchors/delegated-prerequisite-checkpoint, tests/canary/workflow-guidance-anchors/delegated-repair-ownership, tests/canary/workflow-guidance-anchors/fix-dont-park-section-relocated, tests/canary/workflow-guidance-anchors/light-path-inline, tests/canary/workflow-guidance-anchors/structured-phase-progress-anchor, tests/canary/workflow-guidance-anchors/ticket-light-path-anchor
Covers: GW23, GW24, GW25, GW26, GW27, GW28, GW29, GW30, GW31, GW51

## What to build

A session that sizes a change reads every light-path rule of `.bench/BENCH.md` in one block, directly under the right-size table.
The block adds the delegated route of decision #4, so a fix from a learning or a drain goes to a write delegate.
Only a change that the reviewer asks for in the session goes inline under the table.
The fix-don't-park rule and the Capture section no longer repeat the drain dispatch, so the dispatch has one statement.
Each anchored needle stays word for word, so each anchor row and each commitment mirror stays green and still bites.

This ticket reads two contracts from the earlier chunks.
Ticket `1-move-the-build-rules-into-their-own-file.md` removed the build rules from `## Workflow` and added the section `### Always-loaded guidance weight` to `CHANGELOG.md`.
Ticket `3-remove-the-repeated-copies.md` added the GW2 bullet to that section.
Neither ticket edits the light-path text or the Capture section.
Line citations below are at base `bc2ea0ff`.
If a landing moves a line, find the text by its quotation.

**The light-path block.** In `## Workflow` of `.bench/BENCH.md`, put these four paragraphs directly after the right-size table, which ends at line 138:

1. The B39 paragraph of line 146, byte for byte. It opens with "**A light-path change needs no commitment.**".
2. The new paragraph below. It holds the first two sentences of the B40 paragraph of line 148, the two route sentences of decision #4, and the third B40 sentence.
3. The last three sentences of the B40 paragraph, byte for byte. They run from "Dispatch that fix" to "as the light-path row states.".
4. The B42 paragraph of line 154, without the clause ", and the drain delegates its light-path fix".

Paragraphs 2 and 4 then read as below:

```markdown
**Delegate each light-path fix.** Outside a build, a `bench learning` entry's light-path fix ships at once only when the active outcome needs it; `/bench-implement-spec` states the build rule. Every light-path fix from a `bench learning` entry or a drain goes to a write delegate. Only a change that I ask for in the session goes inline under the table. `/bench-drain` implements no light-path fix inline: it dispatches every light-path fix in its run, whatever its source and whatever the active commitment is.

**Fix, don't park.** A small defect that the active committed outcome needs fixed is not roadmap work: the fix lands in the active workflow as its own commit. A small defect that the active outcome does not need goes to `bench learning`. Park a fix to `capture/IDEAS.md` or `capture/learnings.md` when it needs a reviewer decision, a new seam, or spec-level design.
```

Keep each paragraph on one physical line, as each source paragraph is today.
Remove the old B39, B40, and B42 paragraphs, each with its blank line.
The paragraph that opens with "**Deliver only the committed outcome.**" then follows the block.
The right-size paragraph at lines 132-133, the table, and each other paragraph keep their text and their order.
The batch-approval paragraph at lines 156-159 stays in place, because chunk GW4 moves it.

**The Capture section.** Remove the last sentence of `## Capture` at lines 180-181: "A drain implements a light-path idea under the light-path fix rule instead."
Line 180 then ends at "as uncommitted intake.", and line 181 goes.
The sentence "Parked ideas land in `capture/IDEAS.md`" stays, because the row `internal/anchors/registry_data.go:159` pins it.
`.agents/commands/bench-drain.md` names "the light-path fix rule in `.bench/BENCH.md`", and paragraph 2 keeps that rule, so the drain needs no edit.

**The sentences that the mirrors quote.** `TestCommitmentGuidance` requires each sentence that it quotes exactly once in the live guide (`internal/conformance/commitment_guidance_test.go:92-96`).
The block sentences at `:49-52` and `:54` open with a space, so each stays in the middle of its physical line.
The B40 lead at `:53` keeps its text, and so does its row `internal/anchors/registry_commitment.go:26`.
The bold labels "**Delegate each light-path fix.**" and "**Fix, don't park.**" each occur once, because the restore cases at `:85` and `:87-89` replace them.
No new sentence holds a forbid needle of `internal/anchors/registry_commitment.go:60-64`.

**The markers and the canary.** The right-size, fix-don't-park, and commitment markers stay in `## Workflow`.
The section rows `internal/anchors/registry_data.go:161`, `:165`, and `internal/anchors/registry_commitment.go:16` red a marker outside that section.
The shared-rule check `internal/conformance/validity_checks_test.go:243-285` requires each marker in the guide.

The canary `fix-dont-park-section-relocated` mutates the B42 text.
Its two `old` strings are "fix lands in the active workflow as its own commit." and `## Capture` with its newline.
Each still occurs once in the live guide.

**No row changes.** This ticket adds, moves, and retires no anchor row, mirror, or canary.
No row pins the two route sentences, so GW25 stays review-owned, as the spec states.
`internal/anchors/registry_data.go` (478 lines) and `internal/anchors/registry_data_test.go` (1,286 lines) are over the 400-line cap, and they gain no line (`internal/structure/structure.go:308-331`).

**The prose bounds.** A bold label counts as one sentence, and a paragraph holds six sentences at most (`internal/prose/prose.go:18`).
Paragraph 2 holds five sentences, and paragraph 3 holds three.
Each new sentence holds fewer than 25 words.
The guide gains one line net: paragraph 3 and its blank line, less line 181.
Its line row in `projects/benchkit.md:490` stays green.

**The changelog.** Add one bullet for chunk GW3 to `### Always-loaded guidance weight` in `CHANGELOG.md`, after the GW2 bullet.
This is the proposed text:

```markdown
- Changed `.bench/BENCH.md` so that its light-path rules read as one block directly under the right-size table. The block now says that a light-path fix from a `bench learning` entry or a drain goes to a write delegate, and that only a change that the reviewer asks for in the session goes inline. The fix-don't-park rule and the Capture section no longer repeat the drain dispatch. Before, these rules sat apart in the Workflow section, and three sentences stated the drain dispatch.
```

**The Writes line.** The line holds each research path and the five bound files of `internal/anchors` and `cmd/bench` (`internal/tickets/registry_data.go:21-27`, `:44-45`).
It holds the six anchor registry files that hold the literal `.bench/BENCH.md` (`internal/anchors/references.go:31-59`).
It also holds the 40 fixture directories that pin `.bench/BENCH.md` or `CHANGELOG.md`.
Of these, 11 pin the guide by `BASE`, 27 pin it by an overlay, and 2 pin the changelog by `BASE`.

This design edits only `.bench/BENCH.md` and `CHANGELOG.md`, and the other paths are the approved upper bound.
This ticket drops no research path.
`internal/conformance/commitment_guidance_test.go` stays in the line, because a reviewer amendment that makes GW25 executable puts its mirror case there.
`internal/conformance/validity_checks_test.go` stays in the line, because it holds the shared-rule markers that GW31 keeps.
Ticket 1 adds the GW03 row and its test to `internal/anchors/registry_chunk_chain.go` and `internal/anchors/registry_chunk_chain_test.go`.
That row binds `.bench/BENCH.md`, so either file can hold the literal at this checkpoint, and both join the line.

## Acceptance

- [ ] The four paragraphs sit directly under the right-size table and before "**Deliver only the committed outcome.**", which the GW3 review reads (GW23).
- [ ] The table keeps its inline route, and the row `internal/anchors/registry_data.go:187` stays green in `TestRootConformance` (GW24).
- [ ] Paragraph 2 holds the two route sentences of decision #4, which the GW3 review reads (GW25).
- [ ] The block keeps "`/bench-implement-spec` states the build rule.", and the case `build rule pointer` of `TestCommitmentGuidance` passes (GW26).
- [ ] `rg -F "drain delegates its light-path fix" .bench/BENCH.md` finds no line (GW27).
- [ ] `rg -F "light-path fix rule instead" .bench/BENCH.md` finds no line (GW28).
- [ ] Each row of `internal/anchors/registry_commitment.go:22-27` and `internal/anchors/registry_retained_workflow.go:243-244` keeps its needle, and `TestRootConformance` passes (GW29).
- [ ] Each case of `TestCommitmentGuidance` finds its sentence once in the live guide and passes (GW30).
- [ ] The three markers stay in `## Workflow`, and the canary `fix-dont-park-section-relocated` bites in `TestEveryRetainedFixtureBitesThroughRegisteredOwner` (GW31).
- [ ] The section `### Always-loaded guidance weight` of `CHANGELOG.md` holds the GW3 bullet after the GW2 bullet (GW51).

## Verification

Record each red in the return, with its command and its failure line.

- Before the edit, run `bench test --package ./internal/conformance --run 'TestCommitmentGuidance|TestEveryRetainedFixtureBitesThroughRegisteredOwner'`, and record the green run.
- Put the four paragraphs under the table first, and keep the old paragraphs. Run `bench test --package ./internal/conformance --run TestCommitmentGuidance`. Each block case reds, because the guide holds its sentence twice. Then remove the old paragraphs and the Capture sentence, and record the green run.

After the edit, run `bench worktree build <target>`, because a named check compiles from the run binary's source.
Then record one observed red through `bench probe` for each executable row:

- GW30: swap "under the table. `/bench-drain`" for "under the table.`/bench-drain`" with `--package ./internal/conformance --run TestCommitmentGuidance`. The case `drain light-path dispatch` reds, because its sentence with the leading space is absent.
- GW29: omit "Dispatch that fix to a fresh write delegate on the mid tier at high effort, in its own bench worktree." with `--check docs-currency-workflow`. The check reds with the diagnostic of `internal/anchors/registry_retained_workflow.go:243`.
- GW24: swap "no breakdown-approval pause, no write-delegate." for "no breakdown-approval pause." with `--check docs-currency-workflow`. The check reds with the diagnostic of `internal/anchors/registry_data.go:187`.
- GW31: omit "the fix lands in the active workflow as its own commit." with `--package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner`. The canary `fix-dont-park-section-relocated` reds, because its `old` text is absent.

Each probe reports a restored subject.
For GW26, the case `build rule pointer` removes its sentence and requires its diagnostic, so its pass records that red (`internal/conformance/commitment_guidance_test.go:101-105`).

For GW23, run `rg -n -F` over `.bench/BENCH.md` for the last table row and the five paragraph openings.
Record that the line numbers run in block order, with one blank line between each pair.
For GW25, run `rg -c -F` for each route sentence over `.bench/BENCH.md`, and record a count of 1.
For GW27 and GW28, run the two `rg` searches of the acceptance rows, and record that each finds no line.

Then run `bench test --package ./internal/anchors --package ./internal/conformance --package ./cmd/bench`.
Also run the checks `docs-currency-workflow`, `prose-mechanics`, `guidance-prose-budgets`, and `root-conformance`.
