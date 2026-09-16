# Slice independently verifiable outcomes

Blocked by: 2-guide-spec-evidence.md
Writes: .agents/skills/bench-craft-tickets/SKILL.md, internal/anchors/registry_debug_loop.go (new), internal/anchors/registry_data.go, internal/conformance/registry_test.go, tests/canary/workflow-guidance-anchors, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, reviews/debug-loop-guidance.md (new), CHANGELOG.md
Covers: DG17, DG18, DG19, DG20, DG21, DG22, DG23, DG24

## What to build

Implement DG-C3 at `.agents/skills/bench-craft-tickets/SKILL.md`.
Keep the guidance, anchors, fixtures, and fresh-session adoption task in this complete ticket.
Consume the accepted DG-C2 tip and its verified shared registry before editing.
The predecessor supplies a green integration base and complete evidence for its own phase.

Shared writes require this order but do not merge the independent phase outcomes.
Use the spec's acceptance predicates and existing conformance seam.
Reuse sufficient existing owner checks instead of adding duplicate policy sentences.
Add a missing rule or reference fixture and prove its mutation red before restoration.
Keep command registry changes limited to mechanical closure of existing inventories.

Use the spec's independent summary and export outcomes with a shared formatter.
Retain two complete serial tickets and checks usable before successors exist.
Merge the supplied test-only fragment into its complete behavior slice.
Do not implement the feature during this task.

Run the matching fresh-session task from the spec before this chunk closes.
Retain its native evidence under this chunk in `reviews/debug-loop-guidance.md`.
Record the source tip, task, first action, evidence, stop behavior, and handoff.
Leave absent or failed adoption evidence open.

## Acceptance

- [ ] Start each ticket with its delivered outcome and smallest complete behavior, tests, and integration.
- [ ] State a concrete acceptance scenario before locking the ticket.
- [ ] Name checks that establish completion while successor tickets remain unbuilt.
- [ ] Record each real dependency and the value its predecessor supplies.
- [ ] Shared writes determine serial order without automatically merging useful outcomes.
- [ ] Split independently useful outcomes and merge fragments that cannot deliver or verify anything alone.
- [ ] Ticket slicing plans evidence without requiring implementation or an existing executable red.
- [ ] A fresh-session slicing task separates useful outcomes that share writes and supplies checks usable before successors exist.
- [ ] Each owned omission fixture turns its registered check red and restores to green.
- [ ] The focused output shows executed fixture tests with no environment skip.
- [ ] The fresh-session evidence names this ticket's source tip and observable result.
- [ ] This outcome verifies while successor tickets remain unbuilt.

## Focused checks

- `bench test --check docs-currency-workflow`
- `bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner`
- `bench test --check guidance-prose-budgets`
- `bench test --check ticket-grammar`
