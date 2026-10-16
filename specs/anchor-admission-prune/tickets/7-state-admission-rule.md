# State the admission rule in the ADR and the guidance

Blocked by: 6-admit-by-class.md
Writes: docs/adr/0012-document-anchors-use-a-declarative-substring-registry.md, .agents/skills/bench-craft-gate/SKILL.md, .agents/skills/bench-craft-tickets/references/slicing-checks.md, .agents/skills/bench-craft-review/references/finding-discipline.md, .agents/skills/bench-craft-spec/SKILL.md, .agents/skills/bench-craft-delegate/references/delegation-discipline.md, projects/benchkit.md, internal/anchors, tests/canary/workflow-guidance-anchors, tests/canary/skill-description-budgets, tests/canary/line-routing, tests/canary/guidance-prose-budgets, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: AP39, AP40, AP41, AP42, AP43, AP44, AP45

## What to build

ADR 0012 states the admission rule as the current decision. An anchor earns a place only in one of four classes:

- It forbids a stale or wrong form.
- It names a label or heading that a parser or template uses.
- It names a routing reference.
- It carries a `.bench/BENCH.md` invariant.

 The wording of a rule sentence is review-owned. Canary coverage stays at one needle for each family. The ADR holds no file path and no code snippet.

The guidance stops asking for rule-sentence anchors:

- The gate-authoring skill asks for a Forbid row only for a stale or wrong form.
- The slicing retirement pass asks for a Forbid row only for a stale or wrong form.
- The finding discipline says that review grades rule wording directly.
- The spec skill allows a drop of a long needle that holds rule wording.
- The delegation discipline applies the probe rule to admitted anchors only.
- The project profile describes the workflow-guidance family by its admitted classes.

The gate-authoring and spec skills are at their prose line budgets, so their edits must not add a line. Keep each kept anchor needle in these files intact, or move its registry row, its test rule, and its fixture together.

## Acceptance

- [ ] ADR 0012 states the four classes, review-owned rule wording, and one needle per family.
- [ ] Each of the six guidance sites states its rule under the amended ADR.
- [ ] `bench test --check docs-currency-workflow`, `bench test --check prose-mechanics`, and `bench test --check guidance-prose-budgets` pass.
