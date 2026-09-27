# 2. Count the classes in status

Blocked by: 1-classify-unclaimed-refs.md
Writes: internal/worktree/clean_unclaimed.go, internal/worktree/land_prunes_landed_siblings_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/status/status.go, tests/canary/docs-currency-token-diet/signal-vocabulary-drift, internal/status/status_producible_test.go, internal/systemtest/status_route_converge_test.go
Covers: RI24, RI26, RI27, RI88, RI89, RI90, RI91

## What to build

Chunk: RI-C1b.
This ticket edits a system-tagged test, so the author runs the system suite under `BENCH_KIT` through `bench test --check system`.

Replace the exported ref list of the worktree package with one exported function that returns the landed, subsumed, and unique counts through the same planner.
The landing prune test reads that function too, so no exported function serves a test alone.

The status git row reads those counts.
Its details read `<n> landed ref`, `<n> subsumed ref`, and `<n> unique ref` through the existing plural helper, and the row omits a zero class.
Any Bench-namespace ref count above zero routes the row to the plan command, and the dirty-path and unpushed-commit details stay in the row text.
The row's action is already the plan-only command from ticket 1.
A repository with no Bench-namespace ref keeps today's row and today's actions.

A faulted set shows `<n> faulted ref`, a planner failure shows `unclaimed refs unavailable`, and both route to the plan command.
When the Git state read fails, the row still counts the unclaimed rows and routes to the plan beside `git state unavailable`.
The unique-branch count subtracts every unclaimed row, faulted rows included, so a symref is counted once.
The worktree package exports the plan command spelling, and the status action table reads it.

Add a second case to the system route test over a landed ref that runs the apply command the plan prints and confirms the removal.
The landed-unclaimed status test changes posture: a landed unclaimed branch now prints `1 landed ref` with the plan route.
Record the plan time over a fixture with 43 unrecorded refs in the ticket's verification note as a number.

## Acceptance

- [ ] Status over one landed, one subsumed, and one unique ref prints `1 landed ref, 1 subsumed ref, 1 unique ref` with the action `bench worktree clean --discard-branch --unclaimed`.
- [ ] Status over one dirty path and one unique ref prints both details and routes to `bench worktree clean --discard-branch --unclaimed`.
- [ ] A repository with one unique feature branch and no Bench-namespace ref prints `1 unique branch` and `git push`.
- [ ] The system route test runs the printed apply command over a landed ref and the ref is gone.
- [ ] Status over one unique ref and one faulted ref includes `1 unique ref, 1 faulted ref` and routes to `bench worktree clean --discard-branch --unclaimed`.
- [ ] Status over one unique Bench ref, a symref to it, and one unique feature branch prints exactly `1 unique ref, 1 faulted ref, 1 unique branch`.
- [ ] Status with an unreadable ledger prints `unclaimed refs unavailable` and routes to the plan command.
- [ ] Status with only one faulted symref to `main` prints `1 faulted ref` and routes to the plan command.
- [ ] Status over a blob-tip Bench ref prints its class details and `git state unavailable` and routes to the plan command.
- [ ] The landing prune test reads the counts function and still passes.
- [ ] The verification note records the plan time over 43 unrecorded refs.
