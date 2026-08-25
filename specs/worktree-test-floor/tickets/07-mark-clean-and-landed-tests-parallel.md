# Mark the clean and landed tests as parallel

Line: sonnet / low.

The spec is exact, the shape is known, and the census and `go test -race` grade
the result.

Blocked by: 01-take-explicit-root-in-four-verbs.md, 02-guard-harness-state-for-parallel-siblings.md
Writes: internal/worktree/clean_landed_test.go, internal/worktree/clean_landed_apply_test.go, internal/worktree/clean_landed_hostile_test.go, internal/worktree/clean_branch_test.go, internal/worktree/clean_operand_test.go, internal/worktree/landed_test.go, internal/worktree/squash_landed_test.go

## What to build

The seven clean and landed files run their eligible tests in parallel. About 46
top-level tests sit in this set. The explicit root of ticket 01 removes the
directory change from the `clean` tests, so most of them become eligible.

Add `t.Parallel()` to every eligible test in this file set. A test that binds
environment or changes the directory stays serial, and it keeps its current
shape. Where a journey needs a private `BENCH_HOME` or a stub `PATH` for the
child alone, pass `Dir` and `Env` on the child instead. That journey then
becomes eligible.

Remove no test, and merge no test. A table-driven parent with parallel subtests
calls `t.Parallel()` too, when it binds no environment.

## Acceptance

- [ ] Every eligible test in the seven files calls `t.Parallel()`.
- [ ] The package passes `go test -race`.
- [ ] The package holds the same top-level test names as before.
