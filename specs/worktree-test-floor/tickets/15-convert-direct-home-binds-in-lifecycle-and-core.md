# Convert the direct home binds in the lifecycle and core files

Line: sonnet / low.

Every verb takes the home, so a direct bind becomes an argument. The census
and `go test -race` grade the result.

Blocked by: 13-fixtures-pass-a-private-home.md
Writes: internal/worktree/worktree_test.go, internal/worktree/ownership_test.go, internal/worktree/lifecycle_acquire_test.go, internal/worktree/lifecycle_policy_test.go, internal/worktree/lifecycle_facts_test.go, internal/worktree/resume_reconcile_test.go, internal/worktree/identifier_operand_test.go, internal/worktree/orphan_test.go, internal/worktree/orphan_sweep_test.go, internal/worktree/orphan_render_test.go, internal/worktree/exec_test.go, internal/worktree/eligibility_test.go, internal/worktree/snapshot_test.go, internal/worktree/pool_reclaim_facts_test.go, internal/worktree/path_identifier_test.go, internal/worktree/request_token_test.go

## What to build

Each test in these files that calls `bindEnv(t, "BENCH_HOME", …)` for an
in-process call passes that home instead and drops the bind. The test then
calls `t.Parallel()`. `TestPoolDefaultBenchHome` and any test that grades the
boundary default keep the bind and stay serial. A test that changes directory
stays serial. Remove no test, and change no assertion.

## Acceptance

- [ ] No `bindEnv(t, "BENCH_HOME", …)` remains in these files except the boundary-default tests.
- [ ] WF01 stays green on the live tree.
- [ ] The package passes `go test -race`.
