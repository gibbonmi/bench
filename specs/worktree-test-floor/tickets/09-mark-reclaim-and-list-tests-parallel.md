# Mark the reclaim and list tests as parallel

Line: sonnet / low.

The spec is exact, the shape is known, and the census and `go test -race` grade
the result.

Blocked by: 01-take-explicit-root-in-four-verbs.md, 02-guard-harness-state-for-parallel-siblings.md
Writes: internal/worktree/pool_reclaim_test.go, internal/worktree/pool_reclaim_facts_test.go, internal/worktree/list_actions_test.go, internal/worktree/resume_test.go, internal/worktree/resume_reconcile_test.go, internal/worktree/release_registration_test.go, internal/worktree/recovery_retry_test.go

## What to build

The seven reclaim, list, resume, and release files run their eligible tests in
parallel. About 46 top-level tests sit in this set. The explicit root of ticket
01 removes the directory change from these verbs' tests.

Add `t.Parallel()` to every eligible test in this file set. A test that binds
environment or changes the directory stays serial, and it keeps its current
shape. Where a journey needs a private `BENCH_HOME` or a stub `PATH` for the
child alone, pass `Dir` and `Env` on the child instead. That journey then
becomes eligible.

Keep `TestConcurrentCleanupRecordsOneTransaction` in the race-test registry.
Remove no test, and merge no test.

## Acceptance

- [ ] Every eligible test in the seven files calls `t.Parallel()`.
- [ ] The package passes `go test -race`.
- [ ] The package holds the same top-level test names as before.
