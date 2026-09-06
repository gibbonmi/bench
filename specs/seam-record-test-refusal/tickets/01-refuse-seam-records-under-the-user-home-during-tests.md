# Refuse seam records under the user's Bench home during tests

Blocked by: none
Writes: internal/benchhome/benchhome.go, internal/benchhome/benchhome_test.go, internal/otelrecord/provider.go, internal/otelrecord/processor_test.go, internal/systemtest/owner_test.go, internal/conformance/checks_test.go, CONTEXT.md
Covers: none

## What to build

While a test binary runs, `internal/otelrecord`'s provider writes no seam
record when the resolved Bench home is the user's own fallback home. The
refusal is a silent no-op: every span still starts and ends, but nothing
reaches disk. The refusal applies whether `BENCH_HOME` names that exact path
or the fallback supplies it. An explicit `BENCH_HOME` that names a different
directory still records.

The system owner in `internal/systemtest/owner_test.go` and the line-repo
fixture in `internal/conformance/checks_test.go` each give their bench child
a private `BENCH_HOME` at their one environment-composition seam. Neither
fixture's bench child can inherit the operator's real Bench home.

`CONTEXT.md` gains glossary entries for Bench home, fallback home, and seam
record.

## Acceptance

- [ ] While a test binary runs, the provider writes nothing when the
      resolved home equals the user's fallback home, with `BENCH_HOME` unset.
- [ ] While a test binary runs, the provider writes nothing when `BENCH_HOME`
      names the fallback home explicitly.
- [ ] While a test binary runs, the provider still records when `BENCH_HOME`
      names a directory other than the fallback home.
- [ ] The system owner's one environment-composition seam gives every bench
      child a private `BENCH_HOME`, unless the caller already names its own.
- [ ] The line-repo fixture's one environment-composition seam gives every
      bench child a private `BENCH_HOME`.
- [ ] `CONTEXT.md` defines Bench home, fallback home, and seam record.
