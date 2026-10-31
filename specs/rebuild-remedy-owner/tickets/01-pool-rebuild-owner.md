# Name the worktree build verb for a pool root from the rebuild owner

Blocked by: none
Writes: internal/freshness/freshness_verify.go, internal/benchguard/rebuild_remedy_test.go (new), cmd/bench/commands.go, cmd/bench/commands_brief_test.go, internal/systemtest/bench_follow_on_test.go, internal/systemtest/session_start_route_test.go, internal/systemtest/owner_stale_seal_test.go, .bench/lib/resolve-bench.sh, .bench/commitment.json
Covers: RB01, RB02, RB03, RB04, RB05, RB06, RB07, RB08, RB09, RB10, RB11, RB12, RB17, RB18, RB19, RB20, RB21, RB22, RB26

## What to build

The rebuild owner in `internal/freshness` returns `bench worktree build '<root>'` for each root whose `cd` remedy the hook's pool guard refuses, and it keeps the `cd` form for every other root.
Every printer of `RebuildAction` and `Refusal.Error` then names a remedy that runs, and the shell copy keeps the `cd` form under D1.
Read these spec sections before you start: "Problem", "Implementation decisions" > "The owner", "The callers", "The shell copy", and "The system pins", the acceptance coverage map, and the edge inventory.

Make these changes:

1. In `internal/freshness/freshness_verify.go`, move the current `RebuildAction` body to `CheckoutRebuildAction(root string) string`. Add the `Rebuild` struct with the fields `Action string` and `Worktree bool`. Add `RebuildIn(home, root string) Rebuild`. `Worktree` is true when `home` is not empty and `benchguard.PoolReference(CheckoutRebuildAction(root), poolkey.Pools(home))` is not empty. Then `Action` is `usage.WorktreeBuildFor(shellQuote(root))`. Otherwise `Action` is `CheckoutRebuildAction(root)`. `RebuildAction(root)` returns `RebuildIn(benchhome.Dir(), root).Action`. Keep `PublishedExecutable` as one `filepath.Join` statement, because the `published-executable-path` check reads its spelling. Keep the `Refusal.Error` text shape. Update the `Refusal` doc comment as the spec states.
2. Add the new file `internal/benchguard/rebuild_remedy_test.go` in `package benchguard_test`. It holds `TestRebuildRemedyKind`, `TestRebuildRemedyPassesTheHooks`, `TestRebuildRemedyHostileRoot`, `TestRebuildActionReadsTheProcessHome`, and `TestRefusalCarriesThePoolRemedy` for rows RB01 to RB10. Each expected new-form remedy is a literal built from the fixture path. The tests that call `t.Setenv` do not call `t.Parallel`. Do not add a file to `internal/freshness`: `bench structure` grades that directory at its 12-file maximum, and `freshness_verify_test.go` is within 7 lines of the 400-line maximum.
3. In `cmd/bench/commands_brief_test.go`, add a pool-root subtest to `TestCommandsBriefGradesTheRunningExecutable` for RB12. Its synthetic checkout lies at `<home>/worktrees/<key>/<id>`, and its child runs with `BENCH_HOME=<home>`. The expectation is the literal `bench worktree build '<root>'`, and the output holds no `cd `. `cmd/bench/commands.go` keeps its `RebuildAction` call.
4. In `internal/systemtest/bench_follow_on_test.go`, make `TestBenchFollowOnHookDegradedRim` read its expected remedy from `freshness.CheckoutRebuildAction(owner.kit)`. Add the pool-kit case of RB17 and RB18: `BENCH_KIT` at `<home>/worktrees/k/x` with `BENCH_HOME=<home>`, a BLOCKED line for `bench gate` that holds `CheckoutRebuildAction(kit)`, and an exit-0 pass of that printed command. In `internal/systemtest/session_start_route_test.go`, read the expectation from `freshness.CheckoutRebuildAction(root)`.
5. In `internal/systemtest/owner_stale_seal_test.go`, read the expected remedy from `freshness.RebuildIn(<child home>, owner.kit).Action`, where `<child home>` is the Bench home that `owner.runAt` gives the child (`internal/systemtest/owner_test.go:96-101`).
6. In `.bench/lib/resolve-bench.sh`, rewrite only the comment above `bench_shell_quote` (RB21). It names `freshness.CheckoutRebuildAction` as the Go source and states that the copy keeps the `cd` form, because it prints only while the follow-on guard is off. Do not change a function body. The file is within 19 lines of the 400-line maximum.

`.bench/commitment.json` is on the `Writes:` line for the spec-stage binding commit only. The author does not edit it.

The slicing scratch run of this change in a pool worktree observed exactly two reds before the test edits: `TestBenchFollowOnHookDegradedRim` and `TestSystemOwnerRefusesAStaleSelectedSeal`. The ordinary Go test phase stayed green across every package. Steps 4 and 5 clear those two reds.

After each Go edit, run `bench worktree build <target>` before a `bench <verb> --in <target>` call (`docs/known-issues.md:10`, `docs/known-issues.md:12`).

## Acceptance

- [ ] `RebuildIn` returns `bench worktree build '<root>'` and `Worktree` true for a root under `<home>/worktrees/`, and the `cd` form and `Worktree` false for a root outside it, for a sibling `<home>/worktrees-old/x`, and for an empty home.
- [ ] The pool-root remedy passes `benchguard.PoolReference` and `benchguard.Classify`, and the `cd` form for the same root fails `PoolReference`.
- [ ] A hostile pool root parses to the four words `bench`, `worktree`, `build`, and the root.
- [ ] `RebuildAction` and `Refusal.Error` read `BENCH_HOME`, and `Refusal.Error` for a pool root carries exactly one `; rebuild with bench worktree build '<root>'` and no `cd `.
- [ ] `bench commands --brief` in a stale pool-path checkout prints the build remedy and no `cd `.
- [ ] The stale follow-on hook for a pool kit prints the `cd` form and passes that printed command.
- [ ] `bench test --check system` is green in this pool worktree, and `TestRefusalRebuildActionIsCopyPasteSafeForHostilePaths` and the `published-executable-path` check stay green.
- [ ] No production file under `internal/freshness` holds the literal `worktrees`.
