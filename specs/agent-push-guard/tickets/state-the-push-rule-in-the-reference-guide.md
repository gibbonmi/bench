# State the push rule in the reference guide

Blocked by: none
Writes: .bench/BENCH-reference.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, tests/canary/workflow-guidance-anchors/reference-agent-push-rule/BASE (new), tests/canary/workflow-guidance-anchors/reference-agent-push-rule/EXPECT (new), tests/canary/workflow-guidance-anchors/reference-agent-push-rule/MUTATE.json (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go, tests/canary/docs-currency-token-diet/benchref-imported, tests/canary/docs-currency-token-diet/benchref-pointer-dropped, tests/canary/docs-currency-token-diet/benchref-section-duplicated, tests/canary/skills-index-command-adapters/adapter-inert-invocation-key, tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy, tests/canary/skills-index-command-adapters/dangling-index, tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted, tests/canary/skills-index-command-adapters/missing-index-field, tests/canary/skills-index-command-adapters/stale-index-wording, tests/canary/skills-index-command-adapters/unindexed-skill, tests/canary/workflow-guidance-anchors/reference-bench-operational-layer, tests/canary/workflow-guidance-anchors/reference-category-context, tests/canary/workflow-guidance-anchors/reference-category-oracle, tests/canary/workflow-guidance-anchors/reference-category-setup, tests/canary/workflow-guidance-anchors/reference-category-work, tests/canary/workflow-guidance-anchors/reference-gate-authority, tests/canary/workflow-guidance-anchors/reference-kit-only-ship, tests/canary/workflow-guidance-anchors/reference-no-path-fallback, tests/canary/workflow-guidance-anchors/reference-progressive-loading-term, tests/canary/workflow-guidance-anchors/reference-refusal-route-shape, tests/canary/workflow-guidance-anchors/reference-retro-capture-owner, tests/canary/workflow-guidance-anchors/reference-retro-drain-owner, tests/canary/workflow-guidance-anchors/reference-skills-guidance, tests/canary/workflow-guidance-anchors/reference-upgrade-route
Covers: PG36, PG41

## What to build

Verify the premise first: read the `Hook Layers` section of
.bench/BENCH-reference.md, which today states the pre-push hook clause alone.
Read the `AfterSpecAuthorization` rows in internal/anchors/registry_data.go and
copy their shape. Read `TestSystemSuiteRouteAnchorsRedOnRemoval` and the
`anchorHarness` helper in internal/anchors/registry_data_test.go. Read the
three fixture files under
tests/canary/workflow-guidance-anchors/agents-system-suite-route/ for the
`BASE`, `EXPECT`, and `MUTATE.json` shape.

Add one sentence to the `Hook Layers` list. The sentence states that the guard
allows a push to a branch other than the default branch. It also states the
four denied forms: a force, a deletion, a broadcast, and an unresolved
destination. Keep the sentence
on one physical line, because the anchor needle reads that line. Name no
repo-only path in the sentence, because the `package-core-guard` conformance
check reds a claim word beside such a path.

Add one `AfterSpecAuthorization` anchor row that pins the sentence, with its own
diagnostic. Add the red-on-removal test beside the other anchor tests. Add the
canary fixture directory `reference-agent-push-rule`. Its `BASE` names
`.bench/BENCH-reference.md`, its `MUTATE.json` replaces the sentence, and its
`EXPECT` holds the diagnostic.

## Acceptance

- [ ] The `Hook Layers` list of .bench/BENCH-reference.md holds the push-rule sentence on one line.
- [ ] The anchor registry holds one row whose needle is that sentence.
- [ ] The new anchor test reds when the sentence leaves the file, and names the row diagnostic.
- [ ] `bench canary` runs the `reference-agent-push-rule` fixture and reports the expected diagnostic.
- [ ] The `package-core-guard` conformance check stays green over .bench/BENCH-reference.md.
- [ ] Self-probe: delete the anchor row and keep the sentence, and report the canary fixture red.
