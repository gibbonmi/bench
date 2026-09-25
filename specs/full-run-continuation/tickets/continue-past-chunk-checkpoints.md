# Continue a full build past each chunk checkpoint

Blocked by: none
Writes: .agents/commands/bench-implement-spec.md, .agents/skills/bench-craft-line/SKILL.md, AGENTS.md, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_retained_workflow.go, internal/conformance/implementation_continuation_test.go
Covers: none

## What to build

An orchestrator ran a `--full` build and reached a green chunk checkpoint. It refreshed
`bench handoff` and then ended its turn with a phase-exit report, but no stop condition
held. The `--full` section of the implement phase did not say that a chunk checkpoint is
not a phase exit. The phase-close rule in the working agreement also read a chunk review
as a phase close.

Two sentences in the `--full` section tell the orchestrator to continue into the successor
chunk in the same turn. The orchestrator stops only on a `craft-line` stop condition. The
phase-close rule applies to a review phase only outside a `--full` build. One anchor row
pins the new sentences, and one independent test expectation pins the anchor row.

The stop conditions in `craft-line` govern each ticket author. One sentence in that section
also applies them to the orchestrator of a `--full` run between chunks. One more anchor
row and its independent test expectation pin that sentence.

## Acceptance

- [ ] `bench test --check docs-currency-workflow` fails on the base guidance with the new anchor row, and passes after the guidance edit.
- [ ] A probe that omits the new sentences turns `docs-currency-workflow` red, and the restore returns it to green.
- [ ] The phase-close rule in `AGENTS.md` does not make a chunk review inside a `--full` build a phase close.
- [ ] The `craft-line` stop conditions govern the orchestrator of a `--full` run between chunks. A probe that omits that sentence turns `docs-currency-workflow` red.
- [ ] No other tracked file states the old phase-close rule.
