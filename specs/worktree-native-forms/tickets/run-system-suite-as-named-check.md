# Run the system suite as a named check

Blocked by: none
Writes: internal/gate/phases.go, internal/gate/phases_test.go, internal/testreport/command.go, internal/testreport/check_test.go, internal/testreport/testreport_test.go

## What to build

An agent runs `bench test --check system` and gets the gate's system phase as a
focused run. `internal/gate` exports one system-phase producer that returns the
operands `-tags=system ./internal/systemtest` and the environment
`BENCH_SYSTEM_ROOT=<root>`. `BenchkitPhases` and the named check both read that
producer, so the two argv lists cannot drift. The check runs `focusedTestArgv`
over the producer's operands under `selectedRunEnvironment`, and it sets no
conformance variable. A graded root that differs from `BENCH_KIT` refuses at
exit 1 and starts no Go.

This ticket is the only one to write `internal/gate` and `internal/testreport`,
so it carries both packages' gate invariants. The ticket
`record-new-forms-in-prose.md` names this check in the changelog and the glossary.

## Acceptance

- [ ] WF16: `bench test --check system` runs the fake `go` with `test -trimpath -count=1 -json -tags=system ./internal/systemtest`.
- [ ] WF17: that run's environment holds `BENCH_RUN_BINARY=<selected>`, `BENCH_KIT=<root>`, and `BENCH_SYSTEM_ROOT=<root>`.
- [ ] WF18: `BenchkitPhases(kit, kit)` holds a `system` phase whose argv ends with the producer's operands and whose environment holds `BENCH_SYSTEM_ROOT=<kit>`, and both consumers call the exported producer.
- [ ] WF19: that run's environment holds no `BENCH_CONFORMANCE_SCOPE`, `BENCH_CONFORMANCE_ROOT`, or `BENCH_CONFORMANCE_TIER`.
- [ ] WF20: with `BENCH_KIT` that names another directory, the check prints `system check unavailable` at exit 1 and runs no `go`.
- [ ] WF21: a fake `go` that emits one `fail` event for `checkfixture` makes the check exit 1 with a `packages[1]{package,status}` row `checkfixture,fail`.
- [ ] WF22: after the check, the gate cache and the lane record are unchanged.
- [ ] The gate `test` phase stays green for the whole `internal/gate` and `internal/testreport` packages.
