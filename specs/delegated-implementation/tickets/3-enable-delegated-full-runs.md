# Enable opt-in delegated full runs

Blocked by: 1-bind-delegated-evidence.md, 2-account-for-participants.md
Writes: tests/canary/docs-currency-token-diet, tests/canary/load-validity-metadata, tests/canary/skills-index-command-adapters, tests/canary/guidance-prose-budgets/over-budget-skill, tests/canary/line-routing/line-binding-prose-drift, .bench/BENCH.md, .bench/BENCH-reference.md, .agents/commands/bench-implement-spec.md, .agents/commands/bench-review-implementation.md, .agents/commands/bench-final-check.md, .agents/skills/bench-craft-line/SKILL.md, .agents/skills/bench-craft-delegate/SKILL.md, .agents/skills/bench-craft-delegate/references/delegation-discipline.md, .agents/skills/bench-craft-tickets/SKILL.md, CONTEXT.md, projects/benchkit.md, docs/adr/0021-benchmark-workflow-orchestration.md, docs/field-guide.html, CHANGELOG.md, internal/anchors/registry_retained_workflow.go, internal/conformance/retained_workflow_test.go, internal/conformance/docs_workflow_helpers_test.go, internal/worktree/delegated_integration_test.go (new), tests/canary/workflow-guidance-anchors, reviews/delegated-implementation.md (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: DI20, DI21, DI22, DI23, DI24, DI25, DI26, DI27, DI28, DI29, DI30, DI32, DI33, DI34, DI35

## What to build

Expose the opt-in phase grammar and delegate protocol through the existing guidance owners.
Use the first ticket's identity contract and the second ticket's assessment support before sending any write charge.
Keep default retained-authorship instructions explicit and unchanged in behavior.
Fold the opt-in exception into current routing, delegation, review, continuation, and final-check guidance.
Keep within current prose budgets by editing those owners, not adding a parallel workflow manual.

Define orchestrator and chunk author in the glossary with distinct responsibilities.
Update ADR 0021 for this exception only.
Add prose tripwires through the existing anchors owner and prove their omission failures.
The tripwires do not claim to execute the harness protocol.

Use existing worktree fixtures for a synthetic journey with two independent author branches and one dependent chunk.
Fold them serially, reject branch-only verification, retain fresh integrated results, and exercise final landing.
Include a failed author, replacement, review repair, and all participant attempts in the synthetic account.
Keep an unrelated assignment unchanged.
No live model comparison or new model default is authorized.
The registry files are required co-named consumers; no new CLI command is introduced.

## Acceptance

- [ ] Only an approved full run with the explicit flag delegates authorship.
- [ ] The orchestrator declares configured model, effort, cap, and author limit before dispatch.
- [ ] Top-tier selection is authorized by the flag without changing no-flag routing.
- [ ] Concurrent authors retain separate assignments and serial integrated acceptance.
- [ ] Dependent dispatch waits for every prerequisite checkpoint.
- [ ] Two completed no-progress attempts permit replacement or model change only after reassessment.
- [ ] Terminal author failure permits replacement or model change.
- [ ] Exhausting the declared author cap permits replacement or model change.
- [ ] A lost author session permits replacement or model change.
- [ ] Every transfer waits for confirmed termination of the old writer.
- [ ] Production repairs remain with recorded chunk authors.
- [ ] The final orchestrator verification and reconciliation precede broker-owned landing.
- [ ] The handoff preserves identities, source pins, replacement history, and pending obligations.
- [ ] The account includes failed dispatches, all authors, every review axis, diagnostics, verification, and orchestration.
- [ ] The synthetic journey preserves unrelated work and launches no paid comparison.
