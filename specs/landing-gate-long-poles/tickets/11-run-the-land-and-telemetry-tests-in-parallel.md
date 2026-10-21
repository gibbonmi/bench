# Run the land and telemetry tests in parallel

Blocked by: 10-run-the-session-start-and-guard-tests-in-parallel.md
Writes: internal/systemtest/land_route_test.go, internal/systemtest/otel_gate_test.go, internal/systemtest/otel_verbs_test.go, internal/systemtest/otel_crash_test.go, internal/systemtest/otel_hook_test.go, internal/systemtest/otel_worktree_test.go, internal/systemtest/status_route_converge_test.go, internal/systemtest/tree_target_test.go, internal/systemtest/parallel_census_test.go (new)
Covers: LGP34

## What to build

The land-route, telemetry, status-route, and tree-target tests run in parallel. None of them reaches `owner.repos`, so each one only gains `t.Parallel`.

Each test in these eight files calls `t.Parallel`. Remove each test's "pending migration" row from the serial table.

A test that fails under parallel load keeps a serial row, and its reason quotes the failing run. Look first at a child that names no `BENCH_HOME` and so shares `owner.home`.

`land_route_test.go` has 508 lines and no grant, so it may not gain net lines. Its five new `t.Parallel` lines need an offset or a split along responsibility.

BENCH_KIT is supplied by bench test --check system through the existing sealed system owner.

## Acceptance

- [ ] Each test in these eight files passes in parallel with its assertions unchanged, or holds a serial row that quotes its failing run.
- [ ] The census passes on the live package.
- [ ] The probe "remove t.Parallel from one otel test" reds the census.
- [ ] The verification record holds the system phase time before and after the change (LGP34).
