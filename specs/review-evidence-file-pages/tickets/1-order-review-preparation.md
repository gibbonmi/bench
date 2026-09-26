# Order the author record before the review charge

Blocked by: none
Writes: .agents/commands/bench-implement-spec.md, internal/anchors, tests/canary/workflow-guidance-anchors, internal/preflight/evidencecmd, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/conformance/injected_ports_registry_test.go, projects/benchkit.md
Covers: RE11, RE12

## What to build

Deliver chunk RE1. Verify the roadmap premise against the phase and current-binding
owner. State this complete order: commit the author record, prepare the review
charge, then dispatch the axes. Extend the existing chunk-chain rule.

RE2 receives a charge bound to the committed author record tip. Keep the stale-tip
refusal. Add a scenario that commits the record before preparation and binds
successfully. Then commit another record change and require the stale-tip refusal.

The separate debug repair owns the conflicting full-read instruction. Consume
its accepted result. Extend the existing anchor and fixture owners. Use
`internal/prose` for any sentence or paragraph analysis.

## Acceptance

- [ ] The canonical rule places the author record commit before the review charge.
- [ ] The rule places axis dispatch after the review charge.
- [ ] An order-swap mutation reds the workflow check while all command tokens remain present.
- [ ] A charge prepared after the record commit passes current binding.
- [ ] A later record commit makes that charge fail current binding.
- [ ] Existing chunk-chain requirements retain their prior coverage.

Use planned `TestReviewRecordChargeOrder` under `internal/preflight/evidencecmd`.
CE57 already checks a moved tip in build mode. This review-mode scenario adds
the author-record sequence. Reuse the shared current-binding fixture.

Extend `TestChunkChainAnchors` and the workflow fixture family for RE11.
Demonstrate the independent expectation's omission mutation through `bench probe`.
Require its restored result before interpreting the red.

Run `bench test --check docs-currency-workflow` and the focused evidence suite.
Use `-parallel 2` for Go tests. Run prose checks before the lane commit.
Commit this ticket on its own green lane before RE2 starts.
