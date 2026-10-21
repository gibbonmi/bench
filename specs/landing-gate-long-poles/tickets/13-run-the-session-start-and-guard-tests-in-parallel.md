# Run the session-start and guard tests in parallel

Blocked by: 12-run-the-compatibility-and-adoption-tests-in-parallel.md
Writes: internal/systemtest/session_start_test.go, internal/systemtest/session_start_route_test.go (new), internal/systemtest/guard_projection_test.go, internal/systemtest/guard_refusal_test.go, internal/systemtest/guard_rim_test.go, internal/systemtest/guard_projection_differential_test.go, internal/systemtest/bench_follow_on_test.go, internal/systemtest/parallel_census_test.go (new), internal/systemtest/serial_table_test.go (new)
Covers: LGP34

## What to build

The session-start, guard, and follow-on hook tests run in parallel, each on private repositories.

`newSessionEnvironmentFixtureWithGoDir` (`internal/systemtest/session_start_test.go:275-277`) takes a private journey repository instead of `owner.repos[2]`. Fourteen tests reach it, all in the same file: 13 through `newSessionEnvironmentFixture` and one directly at `:146`. The test at `:343` uses `owner.repos[2]` at `:345` and also moves.

These tests take private repositories: `TestGuardProjectionNeverExecutes` (`guard_projection_test.go:49`), `TestGitGuardHookAndCoreParseRefusal` (`guard_refusal_test.go:103`), and both tests of `bench_follow_on_test.go` (`:42`, `:230`). Each other test in these six files calls `t.Parallel`.

`session_start_test.go:192-197` bounds an elapsed time (TE15). If parallel load breaks that bound, keep the test serial, and quote the failing run in its reason.

Remove each moved test's "pending migration" row from the serial table. Run `bench test --check system` twice before you call a test safe in parallel.

`session_start_test.go` has 457 lines and no grant, and it gains 16 `t.Parallel` lines. Move a coherent group of its tests to the new `session_start_route_test.go`, which the growth ratchet never punishes.

BENCH_KIT is supplied by bench test --check system through the existing sealed system owner.

## Acceptance

- [ ] No test in these files reaches `owner.repos`, and the census passes.
- [ ] Each moved test passes in parallel in two runs, with its assertions unchanged.
- [ ] `session_start_test.go` does not gain net lines.
- [ ] The probe "point newSessionEnvironmentFixtureWithGoDir back at owner.repos[2]" reds the census.
- [ ] The verification record holds the system phase time before and after the change (LGP34).
