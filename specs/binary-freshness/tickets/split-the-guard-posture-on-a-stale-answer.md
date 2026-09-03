# Split the guard posture on a stale answer

Blocked by: none
Writes: .bench/hooks/block-bench-follow-on.sh, .bench/lib/resolve-bench.sh, internal/benchguard/benchguard_test.go, internal/systemtest/bench_follow_on_test.go, internal/conformance/guard_classifier_table_test.go (new), tests/canary/guard-classifier-table (new), internal/conformance/registry/registry.go, internal/conformance/tier_test.go, internal/conformance/package_core_checks_test.go, internal/conformance/registry_test.go, tests/canary/package-core-guard/guard-resolver-order-drift
Covers: BF22, BF23, BF24, BF25

## What to build

Verify the premise first: the shim in .bench/hooks/block-bench-follow-on.sh
forwards the child's exit 2 without reading its stderr, so a stale binary's
`unknown subcommand` denies every call. Then capture the child's stderr. When
the exit is 2 and the stderr holds `unknown subcommand`, run a shell word test
over the envelope's command: a segment whose first word is `bench` is a Bench
call. Refuse a Bench call at exit 2 with the `RebuildAction` sentence. Pass
any other call at exit 0 with a warning that holds the sentence. Forward a
genuine refusal, whose stderr holds `BLOCKED:`, unchanged.

Add one shared fixture table of commands and expected Bench-call verdicts,
with resolver-independent rows only, because `InvokesBench` takes a resolver
the shell test lacks. A Go test reads it over `benchguard.InvokesBench`, and a
shell test reads it over the word test. Add it as a canary fixture so the table is inventoried.

Run the system-tagged tests with `BENCH_KIT` and `BENCH_RUN_BINARY` set, as
the system suite requires.

## Acceptance

- [ ] The degraded-rim test with a fake binary that answers `unknown subcommand` passes `ls` at exit 0 with the sentence on stderr.
- [ ] The same fake binary refuses `bench gate` at exit 2 with the sentence.
- [ ] A fake binary that prints `BLOCKED:` at exit 2 keeps exit 2 for `ls`.
- [ ] The shared table test passes on both sides and reds when one row's verdict is flipped on one side.
- [ ] Self-probe: widen the shim to pass every exit 2, and report the `BLOCKED:` row red.
