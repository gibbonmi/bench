# State the comment-only correction rule in the guidance

Blocked by: 3-accept-proven-gaps-at-checkpoint.md
Writes: CHANGELOG.md, .agents/skills/bench-craft-line/references/bounded-repair-policy.md, .agents/commands/bench-implement-spec.md, internal/anchors/registry_retained_workflow.go, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/conformance/implementation_continuation_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_ft311_preparation.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/workflow-guidance-anchors/delegated-dispatch-declaration, tests/canary/workflow-guidance-anchors/delegated-entry-refusals, tests/canary/workflow-guidance-anchors/delegated-resumption-contents, tests/canary/workflow-guidance-anchors/dg-25, tests/canary/workflow-guidance-anchors/dg-26, tests/canary/workflow-guidance-anchors/dg-29, tests/canary/workflow-guidance-anchors/dg-29-verification-target, tests/canary/workflow-guidance-anchors/dg-30, tests/canary/workflow-guidance-anchors/dg-31, tests/canary/workflow-guidance-anchors/dg-31-contradiction-trigger, tests/canary/workflow-guidance-anchors/implement-spec-adoption-freshness, tests/canary/workflow-guidance-anchors/implement-spec-coverage-task-seeding, tests/canary/workflow-guidance-anchors/implement-spec-cross-harness-pointer, tests/canary/workflow-guidance-anchors/implement-spec-entry-validation, tests/canary/workflow-guidance-anchors/implement-spec-incapable-harness, tests/canary/workflow-guidance-anchors/implement-spec-inline-exception, tests/canary/workflow-guidance-anchors/implement-spec-mandatory-delegation-anchor, tests/canary/workflow-guidance-anchors/implement-spec-prose-owner-transfer, tests/canary/workflow-guidance-anchors/implement-spec-read-only-helper, tests/canary/workflow-guidance-anchors/implement-spec-red-preflight-route, tests/canary/workflow-guidance-anchors/implement-spec-review-charge-omitted, tests/canary/workflow-guidance-anchors/implement-spec-review-charge-order, tests/canary/workflow-guidance-anchors/implement-spec-review-charge-reversed, tests/canary/workflow-guidance-anchors/implement-spec-status-flip-anchor, tests/canary/workflow-guidance-anchors/implement-spec-worktree-before-preflight, tests/canary/workflow-guidance-anchors/implement-spec-write-delegation, tests/canary/workflow-guidance-anchors/line-anchor-missing, tests/canary/workflow-guidance-anchors/prepared-build-approval, tests/canary/workflow-guidance-anchors/prepared-build-freshness
Covers: CG37, CG38, CG39

## What to build

Add a concise typed `CHANGELOG.md` entry for the checkpoint behavior.

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

`Writes:` also names the closure that build preflight requires for the two
anchored guidance files. That closure is the other anchor registry files that
name them, the five command-binding files, and the canary fixtures that pin
them. No edit is expected in that closure. No canary
fixture holds the old sentence or the new sentence.

## Acceptance

- [ ] The policy section `Classification and completion` holds the new sentence directly after the review-record prose definition.
- [ ] A removal of the new sentence makes the new anchor fail with `implementation continuation: bounded repair dropped comment-only evidence`.
- [ ] A deletion of the new registry row makes `TestImplementationContinuation` fail with `implementation-continuation anchor is absent`.
- [ ] The `Land` section of the implementation phase holds `Only record commits and comment-only corrections follow the chunk tip.`
- [ ] The old sentence in the `Land` section makes `TestChunkChainAnchors` fail.
- [ ] `bench test --package ./internal/anchors` and `bench test --package ./internal/conformance` pass.
- [ ] `bench-implement-spec.md` stays at 80 lines, and `registry_retained_workflow.go` stays under 400 lines.
