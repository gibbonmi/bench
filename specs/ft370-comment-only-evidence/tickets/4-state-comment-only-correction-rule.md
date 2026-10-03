# State the comment-only correction rule in the guidance

Blocked by: 3-accept-proven-gaps-at-checkpoint.md
Writes: .agents/skills/bench-craft-line/references/bounded-repair-policy.md, .agents/commands/bench-implement-spec.md, internal/anchors/registry_retained_workflow.go, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/conformance/implementation_continuation_test.go
Covers: CG37, CG38, CG39

## What to build

Add one sentence to the bounded repair policy. Put it in the section
`Classification and completion`, directly after the review-record prose
definition. The sentence is:

`A comment-only Go correction that the checkpoint accepts is evidence-only despite its source change, and the orchestrator commits it with no plan assignment.`

Add one anchor row for that sentence in
`internal/anchors/registry_retained_workflow.go`, beside the three
evidence-only rows. The row requires the sentence in the same section. Its
diagnostic is `implementation continuation: bounded repair dropped comment-only evidence`.

Add one independent expectation for that diagnostic to
`TestImplementationContinuation`, with the file, the section, and the needle.
Delete the new registry row, and record the red
`implementation-continuation anchor is absent` in the commit message. Then
restore the row.

Change one sentence of `bench-implement-spec.md` line 54 in place:

- The old sentence is `Only record commits follow the chunk tip.`
- The new sentence is `Only record commits and comment-only corrections follow the chunk tip.`

The chunk-chain needle in
`registry_chunk_chain.go` and its expectation in `TestChunkChainAnchors` hold
the old sentence after the `main` merge sentence. Change that sentence in both
to the new sentence. Their diagnostic becomes
`chunk chain: only record and comment-only commits follow the chunk tip`.

Ticket 3 supplies the checkpoint that accepts the gap and the chain-gap clause
`only record commits and comment-only corrections follow a chunk tip`. This
guidance states that behavior. Do not change `.bench/BENCH.md` or
`bench-review-implementation.md`. The spec's `The guidance` decision states why
their repair sentences stay.

## Acceptance

- [ ] The policy section `Classification and completion` holds the new sentence directly after the review-record prose definition.
- [ ] A removal of the new sentence makes the new anchor fail with `implementation continuation: bounded repair dropped comment-only evidence`.
- [ ] A deletion of the new registry row makes `TestImplementationContinuation` fail with `implementation-continuation anchor is absent`.
- [ ] The `Land` section of the implementation phase holds `Only record commits and comment-only corrections follow the chunk tip.`
- [ ] The old sentence in the `Land` section makes `TestChunkChainAnchors` fail.
- [ ] `bench test --package ./internal/anchors` and `bench test --package ./internal/conformance` pass.
- [ ] `bench-implement-spec.md` stays at 80 lines, and `registry_retained_workflow.go` stays under 400 lines.
