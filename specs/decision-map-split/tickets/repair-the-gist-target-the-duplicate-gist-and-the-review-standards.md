# Repair the gist target, the duplicate gist, and the review standards findings

Blocked by: record-the-split-in-context-the-adr-and-the-public-docs.md
Writes: internal/maps/tickets.go, internal/maps/tickets_test.go, internal/maps/maps.go, internal/maps/validation.go, internal/maps/schema.go, tests/canary/decision-map-integrity/, internal/conformance/decision_map_integrity_test.go, internal/anchors/registry_decision_maps.go, internal/anchors/registry_decision_maps_test.go, tests/canary/workflow-guidance-anchors/, README.md, tests/canary/docs-currency-token-diet/readme-command-first, tests/canary/load-validity-metadata/readme-shared-rule-drift, tests/canary/workflow-guidance-anchors/readme-shaping-skip, tests/canary/workflow-guidance-anchors/field-guide-reduced-schema-map-card, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: DS51, DS52, DS53

## What to build

A gist whose link target is not this map's own `<topic>/tickets/<n>.md` reds with `Decisions so far links missing ticket #<n>`. Two gists that link one ticket red with `Decisions so far duplicate gist for ticket #<n>`. Each new diagnostic gets one fixture, `gist-wrong-folder` and `gist-duplicate`, and its inventory entry.

The README sentence that describes the index and the ticket files registers as a Require needle with its mutation-table row. The doc comment of the decision-map registry file states what the file holds and names no test. The ticket-file layout derives from one function, and `unresolvedState` is inlined.

## Acceptance

- [ ] A gist to `other/tickets/1.md` on map `split` reds with `Decisions so far links missing ticket #1`, and the fixture bites.
- [ ] Two gists to `split/tickets/1.md` red with `Decisions so far duplicate gist for ticket #1`, and the fixture bites.
- [ ] The README needle `each decision lives in one ticket file under the map's tickets folder` reds when removed, through its registry row and its test row.
- [ ] `ticketFolder` and `ticketPath` share one layout derivation, and `unresolvedState` no longer exists.
