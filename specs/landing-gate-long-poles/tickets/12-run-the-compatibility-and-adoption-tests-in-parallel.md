# Run the compatibility and adoption tests in parallel

Blocked by: 11-run-the-owner-and-recovery-tests-in-parallel.md
Writes: internal/systemtest/compatibility_test.go, internal/systemtest/compatibility_writers_test.go (new), internal/systemtest/compatibility_session_test.go, internal/systemtest/lifecycle_test.go, internal/systemtest/adoption_test.go, internal/systemtest/wrapper_pool_home_test.go, internal/systemtest/parallel_census_test.go (new), internal/systemtest/serial_table_test.go (new)
Covers: LGP34

## What to build

The compatibility, lifecycle, adoption, and wrapper-home tests run in parallel, each on private repositories.

`newCompatibilitySession` (`internal/systemtest/compatibility_test.go:58-60`) takes a private journey repository instead of `owner.repos[0]`. Its callers are at `compatibility_test.go:52`, `:165`, `:177`, `:208`, `:236`, `:271`, and `:320`, at `compatibility_session_test.go:15`, `:56`, and `:74`, and at `lifecycle_test.go:138`.

`TestCompatibilityConcurrentWriters` uses `owner.repos[1]` at `:339` and `owner.repos[2]` at `:384`. It takes two private repositories instead.

`TestAdoptionSmokeJourney` (`adoption_test.go:71`) and `TestWrapperNamesTheMissingPoolHome` (`wrapper_pool_home_test.go:24`) take private repositories. Each other test in these five files calls `t.Parallel`.

Remove each moved test's "pending migration" row from the serial table. Run `bench test --check system` twice before you call a test safe in parallel. A test that fails under parallel load keeps a serial row, and its reason quotes the failing run.

`compatibility_test.go` has 387 lines, so 13 lines of headroom remain under the 400-line limit. If the change needs more, move `TestCompatibilityConcurrentWriters` to the new `compatibility_writers_test.go`.

The compatibility tests hold the PS40 witness that `specs/production-test-seam-policy/spec.md:850` cites. Keep its assertions.

BENCH_KIT is supplied by bench test --check system through the existing sealed system owner.

## Acceptance

- [ ] No test in these files reaches `owner.repos`, and the census passes.
- [ ] Each moved test passes in parallel in two runs, with its assertions unchanged.
- [ ] The probe "point newCompatibilitySession back at owner.repos[0]" reds the census.
- [ ] The verification record holds the system phase time before and after the change (LGP34).
