# State the slicing checks and proofs in craft-tickets

Blocked by: 1-close-anchor-registry-writes.md, 2-match-fence-to-writes.md
Writes: .agents/skills/bench-craft-tickets/SKILL.md, .agents/skills/bench-craft-tickets/references/slicing-checks.md (new), internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_retained_workflow.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/workflow-guidance-anchors/craft-tickets-slice-acceptance-row, tests/canary/workflow-guidance-anchors/delegated-serial-ticket-checkpoint, tests/canary/workflow-guidance-anchors/dg-17, tests/canary/workflow-guidance-anchors/dg-18, tests/canary/workflow-guidance-anchors/dg-19, tests/canary/workflow-guidance-anchors/dg-20, tests/canary/workflow-guidance-anchors/dg-21, tests/canary/workflow-guidance-anchors/dg-22, tests/canary/workflow-guidance-anchors/dg-23, tests/canary/workflow-guidance-anchors/ticket-executable-route-pass, tests/canary/workflow-guidance-anchors/ticket-light-path-carve-out, tests/canary/workflow-guidance-anchors/ticket-lock-passes, tests/canary/workflow-guidance-anchors/ticket-relocation-destinations, tests/canary/workflow-guidance-anchors/ticket-seam-creating-slice, tests/canary/workflow-guidance-anchors/ticket-skill-contract-anchor, tests/canary/workflow-guidance-anchors/ticket-source-clause-pass, tests/canary/workflow-guidance-anchors/ticket-template-anchor
Covers: SC24, SC25, SC26, SC27, SC28, SC29, SC30, SC31

## What to build

A ticket slicer opens the `craft-tickets` skill and follows one pointer sentence to `references/slicing-checks.md`. That reference lists every `Writes:` rule that the ticket parser and preflight enforce. The list includes the anchor closure and the fence union that tickets 1 and 2 ship. It then states six slicing rules, each as one sentence:

- A lane-check ticket proves its check through the real lane over a composed tree.
- A posture change that reds a fixture helper makes the slicer list every call site of that helper before the map locks.
- A combined behavior row belongs to the ticket that completes its final consumer.
- A retirement pass gives each sentence that grants the retired behavior its own forbid row and red-capable check.
- A cited verifier row names the exact checks it performed.
- The slicer runs build preflight again after each fence change and before review.

The skill replaces its parser-rules paragraph with the pointer sentence, so it stays within its 100-line prose budget. The ticket-slicing anchor group pins the pointer, the two enforced rules, and each slicing rule. Its test file states each needle and diagnostic independently. The other anchor files and the canary fixture directories are named for closure only; this ticket does not edit them.

## Acceptance

- [ ] `TestTicketSlicingPasses` pins the pointer, the anchor closure rule, the fence union rule, and each of the six slicing rules.
- [ ] Deleting any one pinned sentence reds the anchor harness.
- [ ] The `guidance-prose-budgets` check passes with the skill at 100 lines or fewer.
- [ ] The prose lane passes for the skill and the reference.
- [ ] `bench preflight build slicing-closure` renders `anchor-closure` and `fence-writes` green for this ticket.
