# 2. Count the classes in status and route to the plan-only command

Blocked by: 1-classify-unclaimed-refs.md
Writes: internal/worktree/clean_unclaimed.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/status/status.go, tests/canary/docs-currency-token-diet/signal-vocabulary-drift, internal/status/status_producible_test.go, internal/systemtest/status_route_converge_test.go
Covers: RI24, RI25, RI26, RI27, RI28

## What to build

Chunk: RI-C1b.
This ticket edits a system-tagged test, so the author runs the system suite under `BENCH_KIT` through `bench test --check system`.

Replace the exported ref list of the worktree package with one exported function that returns the landed, subsumed, and unique counts through the same planner.
The status git row reads those counts.
Its details read `<n> landed ref`, `<n> subsumed ref`, and `<n> unique ref` through the existing plural helper, and the row omits a zero class.
The dirty-path and unpushed-commit details stay.
The row's action becomes `bench worktree clean --discard-branch --unclaimed`, and the `--apply-current` action leaves the status action table.
A repository with no Bench-namespace ref keeps today's row and today's actions.

Rewrite the system route test so that the routed command over a unique ref exits 0 and leaves the ref.
Add a second case over a landed ref that runs the apply command the plan prints and confirms the removal.
Record the plan time over a fixture with 43 unrecorded refs in the ticket's verification note as a number.

## Acceptance

- [ ] Status over one landed, one subsumed, and one unique ref prints `1 landed ref, 1 subsumed ref, 1 unique ref` with the action `bench worktree clean --discard-branch --unclaimed`.
- [ ] Status over one dirty path and one unique ref prints both details and routes to `/bench-final-check`.
- [ ] A repository with one unique feature branch and no Bench-namespace ref prints `1 unique branch` and `git push`.
- [ ] The system route test runs the routed command over a unique ref, exits 0, and the ref survives.
- [ ] The system route test runs the printed apply command over a landed ref and the ref is gone.
- [ ] The verification note records the plan time over 43 unrecorded refs.
