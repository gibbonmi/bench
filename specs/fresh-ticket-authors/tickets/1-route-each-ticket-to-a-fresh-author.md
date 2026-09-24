# 1. Route each ticket to a fresh author

Blocked by: none
Writes: .bench/BENCH.md, .bench/BENCH-reference.md, .agents/commands/bench-implement-spec.md, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_retained_workflow.go, internal/conformance/retained_workflow_test.go, internal/conformance/charge_evidence_guidance_test.go, internal/conformance/docs_workflow_helpers_test.go, tests/canary/workflow-guidance-anchors/, internal/anchors/anchor_harness_diagnostics_test.go, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, tests/canary/docs-currency-token-diet/, tests/canary/load-validity-metadata/, tests/canary/skills-index-command-adapters/, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: FA1, FA2, FA3, FA4, FA5, FA6, FA7, FA8, FA9, FA10, FA11, FA12, FA13, FA14

## What to build

Chunk: FA.

Reviewer approval step: the Claude Code auto-mode classifier refuses an agent edit to `.bench/BENCH.md`. Before the author starts, the reviewer makes the `.bench/BENCH.md` edit or grants a permission rule for it. The orchestrator asks for that step in its dispatch.

Rewrite the owner rule in `.bench/BENCH.md`.

Every spec-backed build gives each ticket a fresh author session, with or without `--full`, on the declared line. The build records its authors in the version 2 delegate plan with an author limit of 1. A tier move still asks the reviewer. An explicit `--delegate` adds concurrent authors and the full tier range. The orchestrator reconciles the final acceptance and integration. A post-review repair goes to a fresh session that the plan records as a new assignment with the trigger `user-directed`.

Each affected ticket takes its own repair session, and that session reruns the ticket's verification. At final reconciliation, the orchestrator does not repair. A finding there goes to a fresh repair session for the ticket that owns the path.

Rewrite `.agents/commands/bench-implement-spec.md` to that rule. The author charge carries the ticket, its coverage rows, its `Writes:` fence, the evidence identity, and the line. The author binds the evidence, fetches the metadata and ticket pages, and reads the rest with targeted reads. The build family's delivery prerequisite becomes `a narrow author read`. The orchestrator reads manifests, returns, and verdicts, not code, and refreshes `bench handoff` at each chunk checkpoint.

Before the first dispatch, the orchestrator writes the version 2 `execution` block and takes the plan amendment. The same amendment splits each chunk verification into one verification for each ticket.

Keep the two continuation-policy pointer sentences in `.bench/BENCH.md` and `.agents/commands/bench-implement-spec.md` byte for byte. The `craft-line` section name "Retained implementation continuation" stays unchanged.

State in `.bench/BENCH-reference.md` that the plan amendment declares version 2 with a delegate execution block. State there that the amendment splits each chunk verification into one verification for each ticket. In the same paragraph, replace "accepted findings return to the retained author there" with the fresh repair session rule. Add the Require and Forbid rows of FA1 to FA14. Change each pin of a retired sentence, and update each canary whose mutation names a changed sentence.

## Acceptance

- [ ] `bench test --check docs-currency-workflow` passes with the FA1 to FA14 needles in place.
- [ ] Each retired sentence of FA4, FA7, FA8, and FA13 reds its Forbid needle when a probe restores it.
- [ ] Each changed canary under `tests/canary/workflow-guidance-anchors/` still reds its mutation.
