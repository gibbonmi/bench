# Merge main into a build source only at a chunk boundary

Blocked by: none
Writes: .agents/commands/bench-review-implementation.md, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go
Covers: none

## What to build

An orchestrator merged `main` into a build source during the repair of a chunk. The review
preflight for that chunk then counted each path from `main` against the ownership fences of
the spec, and it refused those paths. The review phase states the chunk bases. One sentence in
that step tells the orchestrator to merge `main` only between chunks, and to use the merge
commit as the next chunk base. The next sentence gives the reason. One anchor row pins the
rule, and one independent test expectation pins the anchor row.

## Acceptance

- [ ] `bench test --check docs-currency-workflow` fails on the base guidance with the new anchor row, and passes after the guidance edit.
- [ ] A probe that omits the new sentence turns `docs-currency-workflow` red, and the restore returns it to green.
- [ ] No other file states the chunk-boundary merge rule.
