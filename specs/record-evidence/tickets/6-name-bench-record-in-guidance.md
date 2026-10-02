# 6. Name `bench record` in the review and implement phases

Blocked by: 5-record-plan-amendment.md
Writes: .agents/commands/bench-review-implementation.md, .agents/commands/bench-implement-spec.md, internal/anchors/, tests/canary/workflow-guidance-anchors/, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: RE98, RE99, RE100, RE101

## What to build

Chunk: RE-C4.

Replace the step 6 lines of `.agents/commands/bench-review-implementation.md` with the text that the spec's "The guidance" section quotes. The text names each form and retires the sentence "Preflight supplies the source and plan digests."

Insert the sentence "An author writes each verification entry with `bench record verification`." after the first sentence of the Land paragraph in `.agents/commands/bench-implement-spec.md`. Keep it on the same physical line, so the file stays at 81 lines.

Add three rules to `chunkChainAnchors` and three expectations to `TestChunkChainAnchors`:

- a require-in-step rule for step 6 of the review phase, with the needle "Write each record entry with `bench record`, not with hand-built JSON."
- a forbid rule for the review phase, with the needle "Preflight supplies the source and plan digests."
- a require-in-section rule for the Land section of the implement phase, with the needle "An author writes each verification entry with `bench record verification`."

## Acceptance

- [ ] Step 6 of the review phase holds the new sentence, and a removal reds `TestChunkChainAnchors`.
- [ ] The retired sentence is absent, and a restored copy reds `TestChunkChainAnchors`.
- [ ] The Land section holds the author sentence, and a removal reds `TestChunkChainAnchors`.
- [ ] `bench test --check guidance-prose-budgets` passes with the implement phase at 81 lines.
