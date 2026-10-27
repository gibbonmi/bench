# State the repair carry in the guidance

Blocked by: 04-record-carried-entries.md
Writes: .agents/commands/bench-implement-spec.md, .bench/BENCH.md, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, CHANGELOG.md, .bench/commitment.json
Covers: CF45, CF46, CF47, CF49, CF50

## What to build

Two guidance statements still demand a full rerun after a repair. Replace each of them, keep each line count, and keep each prose bound of 25 words for each sentence and 6 sentences for each paragraph.

- In the Land section of the implement phase, replace the sentence ``A repair session runs each check that its ticket owes once at its repair commit, and its `bench record ticket` call refreezes the chunk at that commit.`` The new text is one sentence of 25 words or fewer. It states that a repair session runs each check that its repair delta reaches, and that the record call carries the rest and refreezes the chunk.
- In the operating guide, replace the clause `that session reruns the ticket's verification with current repair coverage`. The new text states that the session reruns each ticket check that its repair delta reaches, with current repair coverage. It can be a second sentence, because the paragraph then holds 6 sentences.

Add two Forbid rows to the chunk chain anchors, with the needles `runs each check that its ticket owes once at its repair commit` and `reruns the ticket's verification with current repair coverage`. Add one changelog heading that states the carry of `bench record ticket` and its test-only bound.

## Acceptance

- [ ] The Land section states that a repair session runs each check that its repair delta reaches, and the anchor ``Before a successor chunk starts, run `bench gate --checkpoint specs/<slug>/spec.md --chunk <id>`.`` still holds.
- [ ] The operating guide states that a repair session reruns each ticket check that its repair delta reaches.
- [ ] `TestChunkChainAnchors` refuses each retired needle.
- [ ] The changelog states the carry and its test-only bound.
- [ ] `bench test --check prose` passes, and each edited file keeps its line count.
