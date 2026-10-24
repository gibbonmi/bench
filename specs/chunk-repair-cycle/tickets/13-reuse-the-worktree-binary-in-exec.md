# Reuse the fresh worktree binary for a Bench child of exec

Blocked by: 12-measure-the-run-binary-cost.md
Writes: internal/worktree/exec.go, internal/worktree/exec_test.go, internal/worktree/exec_call_env_test.go, internal/worktree/exec_pwd_test.go, internal/worktree/effects.go, internal/runbinary/runbinary.go, internal/runbinary/runbinary_test.go, internal/freshness/freshness_verify.go, internal/usage/worktree.go, internal/treetarget/run.go, internal/treetarget/build_test.go, internal/systemtest/tree_target_test.go, internal/systemtest/owner_environment_test.go
Covers: RC77, RC78, RC79, RC80, RC81

## What to build

Review chunk: RC-C7. Ticket 12 must record a positive net saving first.

When the base name of the child's `argv[0]` is `bench`, `bench worktree exec` runs `freshness.VerifyAt(worktree, worktree, dist/bench)` (`internal/freshness/freshness_verify.go:23-40`). On a pass, it sets `BENCH_RUN_BINARY` to the worktree's `dist/bench`. On a refusal, it omits the variable, the child builds its private selection as today, and exec prints `bench worktree build` on stderr. For any other `argv[0]`, exec runs no freshness check and keeps today's environment (`internal/worktree/exec.go:244-252`).

Put the rule in the exec verb path only. `runChild` and `execEnv` (`exec.go:148-151`) also serve `RunTreeChild` (`exec.go:119-121`), which runs the `--in` tree-target child (`internal/treetarget/run.go:61-66`). That child keeps today's environment, and `internal/treetarget/build_test.go:82-86` stays green. `internal/systemtest/tree_target_test.go:80` runs `bench worktree exec alpha -- <bench> status`, so it reaches the new check. The worktree package reads the environment only in its effects file (`internal/worktree/effect_census_test.go:15,21-49`), so pass each value by parameter.

The child's `ReuseOrOwn`, `Inherit`, and `Validate` sequence checks the source digest again (`internal/runbinary/runbinary.go:155-183`). A binary that became stale after the exec check refuses, and the refusal names `bench worktree build`. Today the inherited refusal wraps `freshness.RebuildAction` (`runbinary.go:211`, `freshness_verify.go:107`), which `internal/runbinary/runbinary_test.go:169-170` pins and about twenty other sites pin on other routes. Compose `Refusal.Reason()` with `usage.WorktreeBuildFor` (`internal/usage/worktree.go:15,61`), as `internal/preflight/binary_seal.go:65-66` does, so `RebuildAction` and its other pins stay. Word the exec hint and the refusal so that no route-bypass needle of `internal/conformance/refusal_route_bypass_test.go:38-47` ("next=", "run bench ") matches outside `internal/refusalroute`. `runbinary.go` is at 387 of 400 lines.

The worktree assignment owns `dist/bench`, and `Close` does nothing for an inherited selection (`runbinary.go:232-235`), so add no process owner. Rewrite the comment at `exec.go:214-243`, because `Inherit` with a named root checks the source digest.

Rewrite `TestExecChildOwnsRatherThanInheritsItsRunBinary` and `TestExecChildDropsWrapperRouting` (`internal/worktree/exec_test.go:105-126`) to RC77 and RC79. Their fixture `worktreeWithWrapper` has no `dist/bench`, and their children are `env` and `printf`, so the RC77 case needs a fresh `dist/bench` and a `bench` child. `exec_test.go` is at 408 lines over its limit, and the commit lane grades growth, so move an existing test to `exec_call_env_test.go` (54 lines) or `exec_pwd_test.go` (99 lines) before you add RC78. `execEnv` has one test caller (`exec_pwd_test.go:84`). Do not add a file to the crowded `internal/worktree/` directory.

The tagged system suite receives `BENCH_KIT` and the sealed run binary from `bench test --check system`; never run the tagged suite directly.

RC81 is review-owned: the final reconciliation of the build measures one test-only repair cycle against the 45-to-60-minute baseline and records it in the review record. This ticket owns the row so that one ticket holds it; its author records nothing for it.

## Acceptance

- [ ] RC77: an exec call of a `bench` child with a fresh `dist/bench` sets `BENCH_RUN_BINARY` to that path.
- [ ] RC78: an exec call of a `bench` child with a stale `dist/bench` omits the variable and prints `bench worktree build`.
- [ ] RC79: an exec call of a non-Bench child runs no freshness check and keeps today's environment.
- [ ] RC80: a child that inherits a binary that became stale refuses and names `bench worktree build`.
- [ ] An `--in` tree-target child keeps today's environment.
- [ ] The probe "pass a stale worktree binary to the child" bites in `bench test --package ./internal/worktree`.
