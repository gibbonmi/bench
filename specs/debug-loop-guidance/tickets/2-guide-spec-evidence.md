# Guide spec authoring with concrete evidence

Blocked by: 1a-enable-unified-review-trial.md
Writes: .agents/skills/bench-craft-spec/SKILL.md, internal/anchors/registry_debug_loop.go (new), internal/anchors/registry_data.go, internal/conformance/registry_test.go, tests/canary/workflow-guidance-anchors, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, reviews/debug-loop-guidance.md (new), CHANGELOG.md
Covers: DG9, DG10, DG11, DG12, DG13, DG14, DG15, DG16

## What to build

Implement DG-C2 at `.agents/skills/bench-craft-spec/SKILL.md`.
Keep the guidance, anchors, fixtures, and fresh-session adoption task in this complete ticket.
Consume the accepted DG-C1 tip and its verified shared registry before editing.
The predecessor supplies a green integration base and complete evidence for its own phase.

Shared writes require this order but do not merge the independent phase outcomes.
Use the spec's acceptance predicates and existing conformance seam.
Reuse sufficient existing owner checks instead of adding duplicate policy sentences.
Add a missing rule or reference fixture and prove its mutation red before restoration.
Keep command registry changes limited to mechanical closure of existing inventories.

Use the spec's new list-summary feature with no executable check.
In a separate variant, leave the empty-list result unspecified.
Retain the reviewer decision request before dependent design, as well as the successful evidence plan.
Keep craft-seams as the owner of the existing uncertain-seam procedure.

Run the matching fresh-session task from the spec before this chunk closes.
Retain its native evidence under this chunk in `reviews/debug-loop-guidance.md`.
Record the source tip, task, first action, evidence, stop behavior, and handoff.
Leave absent or failed adoption evidence open.

## Acceptance

- [ ] For each approved outcome, state a concrete scenario before choosing its verification seam.
- [ ] Inspect the current behavior and its relevant owner before selecting the next authoring action.
- [ ] Name evidence that exposes the cheapest wrong result for the scenario.
- [ ] Use a sufficient existing seam before exploring alternatives.
- [ ] Choose a bounded authoring action and inspect its result before continuing.
- [ ] Return unresolved intended behavior to the reviewer before dependent authoring continues.
- [ ] A new-feature spec can plan its evidence without an existing executable red.
- [ ] A fresh-session spec task produces a concrete scenario and evidence plan for a feature with no executable implementation.
- [ ] Each owned omission fixture turns its registered check red and restores to green.
- [ ] The focused output shows executed fixture tests with no environment skip.
- [ ] The fresh-session evidence names this ticket's source tip and observable result.
- [ ] This outcome verifies while successor tickets remain unbuilt.

## Focused checks

- `bench test --check docs-currency-workflow`
- `bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner`
- `bench test --check guidance-prose-budgets`
- `bench test --check ticket-grammar`
