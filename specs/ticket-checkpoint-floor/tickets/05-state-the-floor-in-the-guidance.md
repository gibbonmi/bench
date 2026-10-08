# State the check floor in the guidance

Blocked by: 01-name-the-root-conformance-check.md, 03-grade-the-plan-against-the-floor.md, 04-refuse-a-planned-cell-at-the-freeze.md
Writes: .agents/commands/bench-implement-spec.md, .agents/skills/bench-craft-delegate/references/delegation-discipline.md, .agents/skills/bench-craft-spec/references/map-discipline.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/anchors/registry_calibration.go, internal/anchors/registry_calibration_test.go, internal/anchors/registry_charge_binding.go, internal/anchors/registry_charge_binding_test.go, internal/anchors/registry_check_floor.go (new), internal/anchors/registry_check_floor_test.go (new), internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_commitment.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, internal/anchors/registry_spec_trace_test.go, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/workflow-guidance-anchors/
Covers: TF60, TF61, TF62, TF63, TF64, TF65

## What to build

The guidance states the floor, the planned-cell rewrite, and the one root conformance command. The spec section "Guidance" owns each statement, and this ticket writes them unchanged.

- The implement-spec floor sentence replaces the focused-checks sentence on line 44. It states that the completion plan carries each ticket's floor and that `check-floor` grades it.
- The implement-spec planned-cell sentence joins the four-sentence paragraph on line 48. It states that the ticket commit rewrites each covered `planned` cell into the citation form.
- Map discipline states that the `planned` marker lasts until the covering ticket commits.
- Delegation discipline names `bench test --check root-conformance` as the root conformance pass, and it keeps its anchored sentence.
- The build-entry route sentence on line 19 adds `check-floor` to the reds that the staleness pass takes.

Each statement describes a behavior that a blocker ticket produces:

- Ticket 01 produces the `bench test --check root-conformance` command that the delegation sentence names.
- Ticket 03 produces the `check-floor` row that the floor sentence and the route sentence name.
- Ticket 04 produces the chunk freeze refusal that makes the TF61 planned-cell rewrite sentence true.

`.agents/commands/bench-implement-spec.md` holds 83 lines against its budget of 83. Each new sentence joins an existing paragraph line, so the file stays at 83 lines.

Each new sentence gets one anchor and one `workflow-guidance-anchors` canary that `docs-currency-workflow` grades. `internal/anchors/registry_data.go` holds 479 lines and has no structure grant. So the four new anchors go into `registry_check_floor.go`, and the new family joins the registry on the existing join line of `registry_data.go`. The route change updates the needle in `registry_chunk_chain.go` and the `old` text of the `implement-spec-red-preflight-route` canary. The other anchor registry files and the command registry files are closure entries, and this ticket expects no edit in them.

At the ticket commit, rewrite the seam cell of TF60 to TF63 and TF65 to cite `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`).

## Acceptance

- [ ] TF60: The implement-spec guidance holds the anchored floor sentence, and its canary mutation reds `docs-currency-workflow`.
- [ ] TF61: The implement-spec guidance holds the anchored planned-cell rewrite sentence, and its canary mutation reds `docs-currency-workflow`.
- [ ] TF62: Map discipline holds the anchored marker-lifetime sentence, and its canary mutation reds `docs-currency-workflow`.
- [ ] TF63: Delegation discipline holds the anchored root-conformance command sentence, and its canary mutation reds `docs-currency-workflow`.
- [ ] TF64: The existing `delegate-root-conformance-pass` canary bites through `docs-currency-workflow` after the guidance edit. Its red comes from a mutation probe, not from absence: change one word of the `old` text in that canary's `MUTATE.json`, and `TestEveryRetainedFixtureBitesThroughRegisteredOwner` turns red.
- [ ] TF65: The build-entry route sentence names `check-floor` among the reds that the staleness pass takes, and the updated `implement-spec-red-preflight-route` canary bites.
