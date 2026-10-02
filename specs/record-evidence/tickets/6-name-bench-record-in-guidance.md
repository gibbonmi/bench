# 6. Name `bench record` in the review and implement phases

Blocked by: 5-record-plan-amendment.md
Writes: .agents/commands/bench-review-implementation.md, .agents/commands/bench-implement-spec.md, internal/anchors/, tests/canary/workflow-guidance-anchors/, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: RE98, RE99, RE100, RE101, RE109, RE110, RE111

## What to build

Chunk: RE-C4.

Replace lines 208 to 214 of step 6 in `.agents/commands/bench-review-implementation.md` with the text that the spec's "The guidance" section quotes. The text names each completed form and keeps the completion entry and each failed, skipped, or pending result hand-written. It retires three sentences:

- "Preflight supplies the source and plan digests."
- "Record the performer, role, model, effort, frozen base and tip, source, state, and native result."
- "Embed the minimal native excerpt and its SHA-256 digest; local logs are supplemental evidence."

Insert these two sentences after the first sentence of the Land paragraph in `.agents/commands/bench-implement-spec.md`, on the same physical line 56:

- "After each ticket commit, the author sets the chunk tip to that commit with `bench record chunk`."
- "An author writes each verification entry with `bench record verification`."

The file keeps 80 lines by the `proseBudgetLineCount` rule.

Add six rules to `chunkChainAnchors` and six expectations to `TestChunkChainAnchors`. The require-in-step rule for step 6 of the review phase takes this needle:

    Write each completed chunk, verification, review, and amendment entry with `bench record`, not with hand-built JSON.

Three forbid rules for the review phase take the three retired sentences above. Two require-in-section rules for the Land section of the implement phase take the two inserted sentences above.

## Acceptance

- [ ] Step 6 of the review phase holds the new sentence, and a removal reds `TestChunkChainAnchors`.
- [ ] Each retired sentence is absent, and a restored copy reds `TestChunkChainAnchors`.
- [ ] The Land section holds both inserted sentences, and a removal of either reds `TestChunkChainAnchors`.
- [ ] The implement phase holds 80 lines, and `bench test --check guidance-prose-budgets` passes.
