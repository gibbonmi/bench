# Render the recovery matrix from the registry

Blocked by: 02-move-the-landing-faces-into-the-shared-registry.md
Writes: internal/refusalroute/command.go (new), internal/refusalroute/command_test.go (new), cmd/bench/main.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/package-core-guard/unrouted-subcommand, .bench/BENCH-reference.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_retained_workflow.go, tests/canary/docs-currency-token-diet/benchref-imported, tests/canary/docs-currency-token-diet/benchref-pointer-dropped, tests/canary/docs-currency-token-diet/benchref-section-duplicated, tests/canary/skills-index-command-adapters/adapter-inert-invocation-key, tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy, tests/canary/skills-index-command-adapters/dangling-index, tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted, tests/canary/skills-index-command-adapters/missing-index-field, tests/canary/skills-index-command-adapters/stale-index-wording, tests/canary/skills-index-command-adapters/unindexed-skill, tests/canary/workflow-guidance-anchors/agents-handoff-section-rule, tests/canary/workflow-guidance-anchors/reference-agent-push-rule, tests/canary/workflow-guidance-anchors/reference-bench-operational-layer, tests/canary/workflow-guidance-anchors/reference-category-context, tests/canary/workflow-guidance-anchors/reference-category-oracle, tests/canary/workflow-guidance-anchors/reference-category-setup, tests/canary/workflow-guidance-anchors/reference-category-work, tests/canary/workflow-guidance-anchors/reference-gate-authority, tests/canary/workflow-guidance-anchors/reference-kit-only-ship, tests/canary/workflow-guidance-anchors/reference-no-path-fallback, tests/canary/workflow-guidance-anchors/reference-progressive-loading-term, tests/canary/workflow-guidance-anchors/reference-refusal-route-shape, tests/canary/workflow-guidance-anchors/reference-retro-capture-owner, tests/canary/workflow-guidance-anchors/reference-retro-drain-owner, tests/canary/workflow-guidance-anchors/reference-skills-guidance, tests/canary/workflow-guidance-anchors/reference-upgrade-route
Covers: RR51, RR52, RR53, RR54

## What to build

Add `bench recovery`, a repository-scoped read verb.
The command lives in `internal/refusalroute`, and `cmd/bench` routes the verb to it.
It reads only the compiled registry.
It prints `recovery[N]{verb,face,authority,route}` with one row per registered face, in registry order, and exits 0.
The route cell holds the rendered route with its placeholders.

Register the verb in the command registry, so that `bench help` lists it.
Update the routing table and the query registry tests that the new verb forces.

The reference guide names `bench recovery` as the recovery matrix.
Update the sentence about the landing refusal shape, so that it names the shared registry.
Update each anchor registry file and canary fixture that pins the changed reference text.

The row expectations derive from the registry's declared faces, never from the renderer under test.

## Acceptance

- [ ] `bench recovery` exits 0 and prints the header `recovery[N]{verb,face,authority,route}`.
- [ ] `bench recovery` prints exactly one row for each registered face, in registry order, and N equals the face count.
- [ ] `bench help` lists `bench recovery`.
- [ ] The reference guide names `bench recovery` as the recovery matrix.
