# State the comment-gap rule in the guidance

Blocked by: 01-chain-proves-comment-gap.md
Writes: .agents/skills/bench-craft-line/references/bounded-repair-policy.md, .agents/commands/bench-implement-spec.md, .agents/commands/bench-review-implementation.md, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, CHANGELOG.md, capture/agent-performance/claude-models.md, .bench/commitment.json
Covers: CG18, CG19, CG20, CG21, CG22, CG23, CG24, CG25, CG26

## What to build

The guidance states the chain rule that ticket 01 ships. A comment-only Go correction before the green checkpoint takes no refreeze, and the review chain accepts a later chunk base after it. Ticket 01 supplies that chain behavior, so this ticket starts after its chunk checkpoint.

Write the exact sentences of the spec's Implementation decisions:

- In the bounded repair policy, replace line 49 with `A comment-only Go correction before the green checkpoint stays evidence-only and takes no refreeze.` Then add `The review chain accepts a later chunk base after that correction, so the chunk keeps its verifications and its reviews.` Lines 48, 50, and 51 stay.
- In the Land paragraph of the implement phase, replace the sentence `Each source change except a record commit or a recorded fold moves the chunk tip and takes a refreeze of that chunk.` with `Each source change except a record commit, a recorded fold, or a proved comment-only Go correction moves the chunk tip and takes a refreeze.` The sentence stays on line 60, because the file is at its 97-line budget (`projects/benchkit.md:553`).
- In the review phase, add `A commit that adds only a proved comment-only Go correction after that base can also be the chunk base.` after `When recorded folds follow the predecessor tip, the base is the last of them.`

Each of the three edited paragraphs then holds 6 sentences, which is the paragraph bound (`internal/prose/prose.go:18`). Add no other sentence to them.

In `internal/anchors/registry_chunk_chain.go`, keep the Forbid row with the needle `comment-only corrections follow the chunk tip` and its mirror row. Add two Forbid rows in group `AfterImplementSpec`:

- the bounded repair policy, needle `takes a refreeze`, diagnostic `chunk chain: the bounded repair policy restored the refreeze of a comment-only correction`;
- the implement phase, needle `Each source change except a record commit or a recorded fold moves the chunk tip`, diagnostic `chunk chain: the implement phase restored the refreeze of a comment-only correction`.

Add one mirror row for each new Forbid row to `TestChunkChainAnchors` in `internal/anchors/registry_chunk_chain_test.go`. Each mirror row pins the retired sentence whole. Add one short sentence to the doc comment of `chunkChainAnchors` that names the retired refreeze sentences, and leave its current long sentence as it is. Use the existing file constants: `implementPhase` (`registry_chunk_chain.go:53`) and `boundedRepairPolicy` (`registry_chunk_chain.go:58`) in the registry, and `repairPolicy` in the test (`registry_chunk_chain_test.go:10`).

Add one `### Review chain after a comment-only correction` section under `## [Unreleased]` in `CHANGELOG.md`. Keep lines 238 and 239 as history, and keep the rows that `internal/anchors/registry_data.go:94` and `:158` hold on that file.

Change the reasons of the two decisions at `capture/agent-performance/claude-models.md:79-80`. Line 79 no longer says that the review chain refuses any source commit between two chunks. Line 80 gives the review staleness of a refreeze as its reason, in place of orphaned verifications.

The spec-stage close writes the deliverable binding to `.bench/commitment.json`. This ticket names the file so that exactly one ticket holds it.

## Acceptance

- [ ] The bounded repair policy, the implement phase, and the review phase hold the exact sentences above, and `bench test --check prose-mechanics` passes.
- [ ] `bench test --check guidance-prose-budgets` passes with `.agents/commands/bench-implement-spec.md` at 97 lines.
- [ ] `TestChunkChainAnchors` reds each retired refreeze sentence and the unscoped wording `comment-only corrections follow the chunk tip`.
- [ ] `TestRegistryNeedlesDistinct` passes with the two new needles.
- [ ] The `docs-currency-workflow` check passes on the tree and turns red when the retired implement sentence returns.
- [ ] `CHANGELOG.md` holds the new section under `## [Unreleased]`, and the two scorecard decisions state no false reason.
