# Wire the guard facts in the subcommand

Blocked by: resolve-the-bare-push-destination.md, rewrite-the-push-verdict-for-refspecs.md
Writes: cmd/bench/main.go, cmd/bench/main_test.go, internal/gitguard/checker_junction_test.go, tests/canary/package-core-guard/unrouted-subcommand, tests/canary/injected-ports/unregistered-port, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: PG20, PG21, PG22, PG23, PG24, PG25, PG26, PG27, PG29, PG30, PG31, PG35, PG40

## What to build

Verify the premise first: read `guardGit` in cmd/bench/main.go, which builds
the `Checker` literal from `git.RefResolves` and `git.BranchExists`. Read
`TestGuardGitBlockAllow` in cmd/bench/main_test.go, whose block case is
`git push`. Read `TestClassifyRealCheckerResolvedComposition`,
`TestClassifyRealCheckerTimeoutComposition`, and the helpers `gc1Repo`,
`runGit`, and `stubGit` in internal/gitguard/checker_junction_test.go. Read the
two sibling diffs for the shared contract: `git.BarePushDestination` and the
three `Checker` fields `DefaultBranch`, `CheckedOut`, and `BareDestination`.

Wire the three facts in `guardGit`. `DefaultBranch` calls `git.ResolvedDefault`,
`CheckedOut` calls `git.CheckedOutBranch`, and `BareDestination` calls
`git.BarePushDestination`. Each reads the process working directory, the way the
two existing probes do. The three fields share one shape, `func() (string, bool)`.
`git.CheckedOutBranch` returns `(string, error)` and the literal `HEAD` on a
detached head, so the `CheckedOut` adapter reports `false` for an error and for
that literal. `guardGit` keeps its exit codes: 0 allow, 2 block, and
3 a genuine failure.

Extend `TestGuardGitBlockAllow` so its block case is `git push --force`, and add
an allowed push. Extend the composition test over a temp repository with the
push rows. It sets `push.default`, the upstream, a detached `HEAD`, and
`bench.allowProtectedPush`. Extend the timeout test with a bare push over the
stub `git`. Both composition tests use `t.Chdir` and stay non-parallel.

## Acceptance

- [ ] With `bench.allowProtectedPush` set to `true`, `git push origin main` returns the default-branch label.
- [ ] A bare `git push` on a checked-out `topic` returns the allow verdict with `push.default` unset and with `push.default` set to `current`.
- [ ] A bare `git push` on a checked-out `main` returns the default-branch label.
- [ ] A bare `git push` on `topic` with upstream `origin/main` and `push.default` set to `upstream` returns the default-branch label.
- [ ] A bare `git push` on a detached `HEAD` returns the unresolved label.
- [ ] `git push origin HEAD` on a detached `HEAD` returns the unresolved label.
- [ ] A bare `git push` under `push.default` `matching` and under `nothing` returns the unresolved label.
- [ ] `git push origin` with no refspec on a checked-out `main` returns the default-branch label.
- [ ] A bare `git push` in a directory outside a repository returns the unresolved label.
- [ ] A bare `git push` with the stub `git` that sleeps past the probe bound returns the unresolved label.
- [ ] `guardGit` exits 2 and prints `BLOCKED:` on stderr for `git push --force`.
- [ ] `guardGit` exits 0 for a push to a topic branch.
- [ ] Self-probe: wire `BareDestination` to the checked-out branch alone, and report the detached-head and `matching` rows red.
