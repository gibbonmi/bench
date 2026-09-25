# Merge main into a build source only before the first chunk

Blocked by: none
Writes: .agents/commands/bench-review-implementation.md, .agents/commands/bench-implement-spec.md, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/reviewrecord/coverage.go, internal/gate/review_checkpoint_commits_test.go
Covers: none

## What to build

The review chain check requires each later chunk base to hold the tree of the accepted
predecessor tip, apart from the review record. It also requires the base to descend from that
tip. Thus a `main` merge at a chunk boundary breaks the chain. A `main` merge inside a chunk
also fails, because the review preflight counts the merged paths against the ownership fences.
A `main` merge is safe only before the first chunk. The landing composes the `main` commits
that arrive during the build.

The review phase and the implement phase tell the orchestrator to merge `main` between chunks.
Change the two phase files so that they state the rule that the check enforces. Change the two
anchor rows that pin those sentences, and their independent test expectations. Change the
refusal message of the chain check so that it states the same rule, and change the test
expectation that pins that message. Do not change the behavior of the check.

## Acceptance

- [ ] `bench test --check docs-currency-workflow` fails on the base guidance with the new anchor rows, and passes after the guidance edit.
- [ ] A probe that omits the new review-phase sentence turns `docs-currency-workflow` red, and the restore returns it to green.
- [ ] The chain-gap refusal states that a default-branch merge lands only before the first chunk, and `bench test --package ./internal/gate` passes.
- [ ] No other tracked file states the old chunk-boundary merge rule.
