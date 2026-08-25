# Mark the lifecycle and orphan tests as parallel

Line: sonnet / low.

The spec is exact, the shape is known, and the census and `go test -race` grade
the result.

Blocked by: 02-guard-harness-state-for-parallel-siblings.md
Writes: internal/worktree/lifecycle_test.go, internal/worktree/lifecycle_policy_test.go, internal/worktree/lifecycle_facts_test.go, internal/worktree/orphan_test.go, internal/worktree/orphan_sweep_test.go, internal/worktree/orphan_render_test.go

## What to build

The six lifecycle and orphan files run their eligible tests in parallel. About
53 top-level tests sit in this set.

Add `t.Parallel()` to every eligible test in this file set. A test that binds
environment or changes the directory stays serial, and it keeps its current
shape. Where a journey needs a private `BENCH_HOME` or a stub `PATH` for the
child alone, pass `Dir` and `Env` on the child instead. That journey then
becomes eligible.

Remove no test, and merge no test. A table-driven parent with parallel subtests
calls `t.Parallel()` too, when it binds no environment.

## Acceptance

- [ ] Every eligible test in the six files calls `t.Parallel()`.
- [ ] The package passes `go test -race`.
- [ ] The package holds the same top-level test names as before.
