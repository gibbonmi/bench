# Activate a local collection pilot

Blocked by: none
Writes: tests/canary/package-core-guard/unrouted-subcommand, tests/canary/data-handling-derivation/undocumented-passlist-var, tests/canary/docs-currency-token-diet/signal-vocabulary-drift, tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary, tests/canary/workflow-guidance-anchors/context-coverage-map-term, tests/canary/workflow-guidance-anchors/context-coverage-row-parts, tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary, tests/canary/workflow-guidance-anchors/context-decision-map-term, tests/canary/workflow-guidance-anchors/context-reader-sweep-term, tests/canary/workflow-guidance-anchors/context-ticket-vocabulary, internal/repairpilot (new), cmd/bench/main.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, DATA_HANDLING.md, CONTEXT.md, specs/repair-collection-pilot, reviews/repair-collection-pilot.md (new)
Covers: RP1, RP2, RP3, RP4, RP5, RP6, RP7, RP8, RP9, RP10, RP11, RP12, RP45, RP46

## What to build

Deliver `bench repair-pilot activate` and the inactive or active report through the production command registry.
Use the shared kit and repository identity owners to keep activation inside the Bench kit.
Create a private, versioned pilot document outside the disposable pool.
Implement the lock and atomic replacement contract, including hostile-input refusals.
A stopped document cannot restart through activation.
The command's initial record route validates input and reports the absent collection capability without claiming to collect evidence.

Move `boundaryRoot` beside the registry's root-resolution helpers before adding the pilot route.
This move pays the new route's headroom in the same commit.
Do not add a file to a crowded command or assessment directory.
The required co-named registries remain unchanged when their facts do not change.

Implement RP-C1 with its command, storage, grammar, and route tests.
`TestRepairPilotStorage` supplies filesystem failures through the consumer-facing dependency seam.
`TestRepairPilotRoute` lives in the existing help inventory test file.
The final report fields for later capabilities remain absent until their producers land.

## Acceptance

- [ ] Explicit activation stores the current time for the canonical kit repository.
- [ ] Inactive report and linked-repository refusal write no pilot state.
- [ ] Repeat activation retains the original document.
- [ ] Lock and persistence failures preserve the previous document.
- [ ] The newly introduced command rejects the hostile input inventory before mutation.
- [ ] The public dispatcher reaches the real pilot owner and preserves assessment routes.
