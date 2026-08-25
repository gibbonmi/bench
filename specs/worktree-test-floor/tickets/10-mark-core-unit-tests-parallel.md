# Mark the core unit tests as parallel

Line: sonnet / low.

The spec is exact, the shape is known, and the census and `go test -race` grade
the result.

Blocked by: 02-guard-harness-state-for-parallel-siblings.md
Writes: internal/worktree/worktree_test.go, internal/worktree/exec_test.go, internal/worktree/identity_component_test.go, internal/worktree/identifier_operand_test.go, internal/worktree/path_identifier_test.go, internal/worktree/request_token_test.go, internal/worktree/snapshot_test.go, internal/worktree/classifier_shape_test.go, internal/worktree/eligibility_test.go

## What to build

The nine core unit files run their eligible tests in parallel. About 54
top-level tests sit in this set, and most of them start no repository.

Add `t.Parallel()` to every eligible test in this file set. A test that binds
environment or changes the directory stays serial, and it keeps its current
shape. Where a journey needs a private `BENCH_HOME` or a stub `PATH` for the
child alone, pass `Dir` and `Env` on the child instead. That journey then
becomes eligible.

Remove no test, and merge no test. A table-driven parent with parallel subtests
calls `t.Parallel()` too, when it binds no environment.

## Acceptance

- [ ] Every eligible test in the nine files calls `t.Parallel()`.
- [ ] The package passes `go test -race`.
- [ ] The package holds the same top-level test names as before.
