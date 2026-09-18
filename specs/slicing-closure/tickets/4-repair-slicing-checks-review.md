# Repair the SC-C3 review findings

Blocked by: 3-state-slicing-checks.md
Writes: .agents/skills/bench-craft-tickets/SKILL.md, .agents/skills/bench-craft-tickets/references/slicing-checks.md, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_retained_workflow.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/workflow-guidance-anchors/craft-tickets-slice-acceptance-row, tests/canary/workflow-guidance-anchors/delegated-serial-ticket-checkpoint, tests/canary/workflow-guidance-anchors/dg-17, tests/canary/workflow-guidance-anchors/dg-18, tests/canary/workflow-guidance-anchors/dg-19, tests/canary/workflow-guidance-anchors/dg-20, tests/canary/workflow-guidance-anchors/dg-21, tests/canary/workflow-guidance-anchors/dg-22, tests/canary/workflow-guidance-anchors/dg-23, tests/canary/workflow-guidance-anchors/ticket-executable-route-pass, tests/canary/workflow-guidance-anchors/ticket-light-path-carve-out, tests/canary/workflow-guidance-anchors/ticket-lock-passes, tests/canary/workflow-guidance-anchors/ticket-relocation-destinations, tests/canary/workflow-guidance-anchors/ticket-seam-creating-slice, tests/canary/workflow-guidance-anchors/ticket-skill-contract-anchor, tests/canary/workflow-guidance-anchors/ticket-source-clause-pass, tests/canary/workflow-guidance-anchors/ticket-template-anchor
Covers: SC24, SC25, SC26, SC27, SC28, SC29, SC30, SC31

## What to build

The SC-C3 review accepted three findings. This ticket repairs them in one cycle.

- SC-C3-ST1 and SC-C3-SP1: the skill pointer and the reference introduction claim the `Writes:` rules only. The reference adds the single required `Writes:` field and the representable-entry rule. The anchor closure rule also names a directory that holds an anchored path.
- SC-C3-CV1: map rows SC24 to SC31 name `docs-currency-workflow` as the grader of a deleted live sentence. `TestTicketSlicingPasses` grades a drifted registry row.

The pinned pointer and the pinned anchor closure rule change with their text, so the registry row and its independent expectation change with them.

## Acceptance

- [ ] The pointer and the reference introduction claim every `Writes:` rule. The reference states each `Writes:` rule that the parser and build preflight enforce.
- [ ] Deleting the new pointer sentence reds `docs-currency-workflow`.
- [ ] Map rows SC24 to SC31 name both graders, and `bench coverage --check` passes.
- [ ] `bench preflight build slicing-closure` renders `anchor-closure` and `fence-writes` green for this ticket.
