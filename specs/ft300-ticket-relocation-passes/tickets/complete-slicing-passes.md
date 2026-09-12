# Complete ticket relocation and verification passes

Blocked by: none
Writes: .agents/skills/bench-craft-tickets/SKILL.md, internal/anchors, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/workflow-guidance-anchors, internal/conformance/registry_test.go, CHANGELOG.md
Covers: none

## What to build

Require a ticket to name each destination its own change relocates.
Include destinations for snapshots and registry rows.
Require the slicer to perform the source-clause and executable-route passes before ticket lock.
Point to the existing map-discipline proof rules where they already own a requirement.
Add omission checks through the existing workflow-anchor owner.

The coordinator reconciles the roadmap after this fix lands.

## Acceptance

- [ ] The slicing guidance requires each relocation destination in Writes.
- [ ] The guidance places both verification passes before ticket lock.
- [ ] The source pass accounts for each applicable source clause.
- [ ] The route pass traces the claimed operation through its executable owner.
- [ ] Removing either pass or the destination requirement fails its owning check.
