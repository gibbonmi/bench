# Record the split in CONTEXT, the ADR, and the public docs

Blocked by: rewrite-the-two-commands-and-register-the-anchors.md
Writes: CONTEXT.md, README.md, docs/field-guide.html, docs/adr/0020-a-decision-lives-in-its-ticket-file.md (new), internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, tests/canary/workflow-guidance-anchors/, tests/canary/docs-currency-token-diet/signal-vocabulary-drift, tests/canary/workflow-guidance-anchors/context-ticket-vocabulary, tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary, tests/canary/workflow-guidance-anchors/context-coverage-map-term, tests/canary/workflow-guidance-anchors/context-coverage-row-parts, tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary, tests/canary/workflow-guidance-anchors/context-decision-map-term, tests/canary/workflow-guidance-anchors/context-reader-sweep-term, tests/canary/docs-currency-token-diet/readme-command-first, tests/canary/load-validity-metadata/readme-shared-rule-drift, tests/canary/workflow-guidance-anchors/readme-shaping-skip, tests/canary/workflow-guidance-anchors/field-guide-reduced-schema-map-card, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: DS45, DS46, DS47, DS48

## What to build

`CONTEXT.md` updates the decision map entry and adds map index, decision ticket file wording, and gist, each with an Avoid list. The gist entry carries the quoted needle from the spec's Further notes, registered with a mutation-table row. ADR 0020 records three decisions: a decision lives in its ticket file, the map is an index, and the topic folder moves as one unit. It names no path.

The README describes the index and the ticket files and keeps `Decision maps are situational`. The field guide's map card describes the same shape. Update the fixture copies of `CONTEXT.md`, the README, and the field guide under their families. This ticket is the last to touch the registry. It carries the whole-registry invariant: every anchor row this spec changed has a mutation row that bites.

## Acceptance

- [ ] [D1] `CONTEXT.md` carries the entries with Avoid lists, and the gist needle's mutation row bites.
- [ ] [D2] `docs/adr/0020-a-decision-lives-in-its-ticket-file.md` records the three decisions and names no path.
- [ ] [D3] `README.md` keeps `Decision maps are situational` and describes the index and ticket files.
- [ ] [D4] The fixture copies of the three documents validate against the updated registry.
- [ ] [D5] Every anchor row this spec changed, in both tickets, has a mutation-table row that bites.
