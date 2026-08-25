# Mark the five slowest tests as parallel

Line: sonnet / low.

The spec is exact, the shape is known, and the census and `go test -race` grade
the result.

Blocked by: 02-guard-harness-state-for-parallel-siblings.md
Writes: internal/worktree/land_resume_refusal_test.go, internal/worktree/ownership_test.go, internal/worktree/land_journey_test.go, internal/worktree/lifecycle_acquire_test.go

## What to build

The four files that hold the five slowest tests run their eligible tests in
parallel. The priced cut lands first where it pays most.

Add `t.Parallel()` to every eligible test in this file set. A test that binds
environment or changes the directory stays serial, and it keeps its current
shape. Where a journey needs a private `BENCH_HOME` or a stub `PATH` for the
child alone, pass `Dir` and `Env` on the child instead. That journey then
becomes eligible.

Remove no test, and merge no test. A table-driven parent with parallel subtests
calls `t.Parallel()` too, when it binds no environment.

## Acceptance

- [ ] Every eligible test in the four files calls `t.Parallel()`.
- [ ] The package passes `go test -race`.
- [ ] The package holds the same top-level test names as before.
