# Add the reviewer-selected map-only route to the spec phase

Blocked by: none
Writes: .agents/commands/bench-write-spec.md, .agents/skills/bench-craft-spec/references/map-only-route.md (new), CHANGELOG.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, docs/field-guide.html, internal/anchors/registry_craft_research.go, internal/anchors/registry_craft_research_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_decision_maps.go, internal/anchors/registry_decision_maps_test.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, projects/benchkit.md, tests/canary/guidance-prose-budgets/over-budget-skill/, tests/canary/line-routing/line-binding-prose-drift/, tests/canary/skill-description-budgets/, tests/canary/skills-index-command-adapters/, tests/canary/workflow-guidance-anchors/
Covers: MR1, MR2, MR3, MR4, MR5, MR6, MR7, MR8, MR9, MR10, MR11, MR12, MR13, MR14, MR15, MR16, MR17, MR18, MR19, MR20, MR21, MR22, MR23, MR24, MR25, MR26, MR27, MR28, MR29, MR30, MR31, MR32, MR33

## What to build

The spec phase gains an opt-in map-only route that the reviewer selects. The full route stays the default.

- `.agents/skills/bench-craft-spec/references/map-only-route.md` is new. It states the eligibility, the ten route steps, the stop and return, and the thin-spec parts in the spec's `Implementation decisions` order.
  It sends the map compile and move to the Entry contract of `/bench-write-spec`, and it sends the close to that file's `Exit handoff`. It copies neither procedure.
- `.agents/commands/bench-write-spec.md` gains the pointer sentence of the spec's `Placement` decision, verbatim, at the end of physical line 49, under `Who runs this phase`. The file stays at 87 lines.
- `internal/anchors/registry_ticket_passes.go` gains one `RequireInSection` anchor for each route sentence of rows MR1-MR4, MR6-MR8, MR10-MR18, MR20-MR23, and the pointer of MR25.
- `internal/anchors/registry_ticket_passes_test.go` gains `TestWriteSpecMapOnlyRoute`. It states each needle and diagnostic as its own literal, as `TestWriteSpecChunkedSlicing` does.
- `projects/benchkit.md` gains one sentence in the spec-authoring entry that names the route.
- `docs/field-guide.html` gains one clause in each of its two write-spec passages.
- `CHANGELOG.md` gains one Unreleased entry.

Every anchored full-route sentence of `/bench-write-spec` and its line wraps stay byte-identical.
The phrase `**Spec authoring** → **a fork of the current conversation**` in `projects/benchkit.md` stays byte-identical.
`README.md` and ADR 0022 get no edit.
The `Writes:` line above holds the closure that the slicing research found. `--propose-writes` returned no addition at e7ab8320.

## Acceptance

- [ ] The route reference holds every route sentence, and `TestWriteSpecMapOnlyRoute` passes.
- [ ] Deleting the stop-and-return anchor from the registry turns `TestWriteSpecMapOnlyRoute` red.
- [ ] `/bench-write-spec` holds the pointer under `Who runs this phase`, and `bench test --check guidance-prose-budgets` passes at 87 lines.
- [ ] `bench test --check docs-currency-workflow` passes, so every existing write-spec anchor and canary stays green.
- [ ] `projects/benchkit.md`, both field-guide passages, and `CHANGELOG.md` name the route.
- [ ] `README.md` and ADR 0022 have no change in the diff.
- [ ] The route reference sends the map compile and move to the Entry contract and the close to the `Exit handoff`, and it copies neither procedure.
- [ ] The route reference requires the thin spec's `Decision source:` line to name the compiled map.
- [ ] The route reference states that each ticket's `Covers:` line names declared rows of the thin spec.
- [ ] `TestWriteSpecMapOnlyRoute` states each needle and diagnostic as its own literal and does not read the registry.
