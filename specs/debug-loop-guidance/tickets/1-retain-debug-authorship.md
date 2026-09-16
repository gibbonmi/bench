# Retain debug repair authorship

Blocked by: none
Writes: .agents/commands/bench-debug.md, .agents/skills/bench-craft-delegate/SKILL.md, tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted, tests/canary/claude-agent-definitions, internal/anchors/registry_debug_loop.go (new), internal/anchors/registry_data.go, internal/conformance/registry_test.go, tests/canary/workflow-guidance-anchors, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, reviews/debug-loop-guidance.md (new), CHANGELOG.md
Covers: DG1, DG2, DG3, DG4, DG5, DG6, DG7, DG8, DG42

## What to build

Implement DG-C1 at `.agents/commands/bench-debug.md`.
Keep the guidance, anchors, fixtures, and fresh-session adoption task in this complete ticket.

Shared writes require this order but do not merge the independent phase outcomes.
Use the spec's acceptance predicates and existing conformance seam.
Reuse sufficient existing owner checks instead of adding duplicate policy sentences.
Add a missing rule or reference fixture and prove its mutation red before restoration.
Keep command registry changes limited to mechanical closure of existing inventories.

Replace the obsolete scoped-fix and never-charges-debug routes in the debug integration section.
Point craft-delegate to that section for the diagnostic report contract.
Keep debug Phases 1 through 6 and loop-constructions local.

Use the spec's fenced empty-list defect for the repair case.
Use its separate out-of-fence cause for the diagnostic case.
Retain the original reproduction and evidence that the delegate preserves dirty work without unauthorized edits.

Run the matching fresh-session task from the spec before this chunk closes.
Retain its native evidence under this chunk in `reviews/debug-loop-guidance.md`.
Record the source tip, task, first action, evidence, stop behavior, and handoff.
Leave absent or failed adoption evidence open.

## Acceptance

- [ ] The session that owns the debug loop writes its in-scope repair.
- [ ] Debug delegates only read-only diagnostic work.
- [ ] A write delegate runs debug through Phase 6 for an in-fence defect.
- [ ] For an out-of-fence defect, run Phases 1 through 3 before the diagnostic handoff.
- [ ] At the fence, stop implementation edits and preserve the in-fence dirty work.
- [ ] The blocked report includes the command, red output digest, ranked hypotheses, failing surface, and in-fence dirty paths.
- [ ] A fresh-session in-fence debug task retains its author through a verified repair.
- [ ] A fresh-session out-of-fence debug task returns the bounded diagnostic report without a repair outside its fence.
- [ ] Each owned omission fixture turns its registered check red and restores to green.
- [ ] The focused output shows executed fixture tests with no environment skip.
- [ ] The fresh-session evidence names this ticket's source tip and observable result.
- [ ] This outcome verifies while successor tickets remain unbuilt.

- [ ] Debug keeps its Phase 1 through Phase 6 procedure and local loop-constructions reference.
- [ ] The retained byte comparison shows unchanged debug phase sections and loop-constructions content.

## Focused checks

- `bench test --check docs-currency-workflow`
- `bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner`
- `bench test --check guidance-prose-budgets`
- `bench test --check ticket-grammar`
