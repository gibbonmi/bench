# Refuse process-state binding in a parallel system test

Blocked by: none
Writes: internal/systemtest/parallel_census_test.go, internal/systemtest/parallel_census_cases_test.go, internal/systemtest/process_state_census_test.go, internal/systemtest/serial_census_test.go, specs/lp-census-process-state/tickets/01-refuse-process-state-binding-in-a-parallel-system-test.md
Covers: none

## What to build

The system parallel census refuses a TMPDIR binding only when the key is the literal `TMPDIR`. A parallel test can change other process state, and the census does not see it.

Widen the census in two parts:

1. The reach walk of a parallel test refuses a call to `os.Setenv`, `os.Unsetenv`, or `os.Chdir` with any argument. The walk also follows the helpers that the test can run. A serial test with a serial row can still make these calls.
2. The TMPDIR rule refuses a `Setenv` call whose key is not a string literal in any function. The census reads no values, so such a key can be TMPDIR.

Pin the receiver name. Every method with a `systemOwner` receiver names its receiver `o`. The census follows owner selectors by that name, so a different name hides the owner from the census. The census reports each `systemOwner` method with a different name or with no name.

The process-state rules move to their own file, so that each census file stays below the structure limit.

## Acceptance

- [ ] A planted parallel test that calls `os.Setenv` with the key `HOME` turns the census red by name.
- [ ] A planted parallel test that calls `os.Unsetenv` turns the census red by name.
- [ ] A planted parallel test that calls `os.Chdir` through a helper turns the census red, and the report names the helper.
- [ ] A planted serial test that calls `t.Setenv` with a variable key turns the census red with the TMPDIR report.
- [ ] A planted serial test with a serial row that calls `os.Setenv` and `os.Chdir` passes the census.
- [ ] A planted `systemOwner` method with the receiver name `s` turns the census red by name.
- [ ] Each planted red row fails before its rule exists and passes after it.
- [ ] `bench test --check system` passes.
