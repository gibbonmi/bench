# Grade the live root inside the dev gate

Blocked by: none
Ownership fence: `internal/gate/phases.go`, `internal/gate/capability_skips.go`, `internal/capability/capability.go`, `internal/conformance/registry/registry.go`, `internal/preprelease/preprelease.go`, their tests, plus the files the ten current diagnostics name (`.agents/commands/bench-implement-spec.md`, `.agents/commands/bench-final-check.md`, `.bench/BENCH.md`, `.bench/BENCH-reference.md`, `decisions/spec-build-review-gate-cadence.md`, `internal/canary/canary.go`) and the profile prose that describes the shipped shape (`projects/benchkit.md`)
Contracts: the conformance registry's 29 checks reach the live tree only through `TestRootConformance`, which `capability.Environment`-skips on an unset `BENCH_CONFORMANCE_ROOT`; the gate's ordinary `test` phase sets no such variable, so the kit's own doc/adapter/guidance contracts are unenforced between releases and the loss folds into an unnamed `environment=1` footer count
Assumptions: the conformance package stays kit-only, so a linked repo's `go test ./...` never materializes the variable; `bench test` reports skips through `internal/testreport` and never reads the gate's skip log; the dev tier stays the unset default and is pinned rather than inherited

## What to build

The gate's ordinary `test` phase points `TestRootConformance` at the graded root, so the
registered mechanical conformance checks run where the doctrine already says they run —
no new phase, no separate driver, no change to `prep-release`'s ship tier. The phase
materializes the variable only where the graded root is the kit that declares the entry
test, mirroring how the race and system phases already materialize.

An environment-class skip inside the oracle becomes red rather than a footer count, and
every skip line names the test that emitted it, so the next silent loss of a check is a
red verdict with an actionable message instead of an integer.

## Acceptance

- [ ] [LR1] With the env assignment in place and the ten tree/contract diagnostics unfixed, `bench gate` exits non-zero and names all ten with the conformance checks' own messages; after the dispositions land it is green and the test phase shows `TestRootConformance` executed rather than skipped.
- [ ] [LR2] An environment-class skip observed by the gate is red with a message naming the test and its reason, and the skip row renders `capability-skips class=environment: N (TestName: reason)`.
- [ ] [LR3] A linked repo (graded root without the conformance entry test) materializes no conformance env on its test phase, and `prep-release` still runs its ship-tier conformance exactly once.
- [ ] [LR4] `bench maps` exits 0 and `projects/benchkit.md` describes the shipped shape: one ordinary driver that grades the live root.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| LR1 | delete `## Entry orientation` from a phase file under `.agents/commands/` | `docs-currency-workflow` in the conformance registry, reached through the gate's test phase | remove the heading, run `bench gate`, expect red naming that file and heading; restore and expect green |
| LR2 | drop the conformance env assignment from the test phase | the gate's environment-skip rule | remove the assignment, run `bench gate`, expect red naming `TestRootConformance` and `BENCH_CONFORMANCE_ROOT not set`; restore |
| LR3 | materialize the env unconditionally instead of behind the entry-test probe | the phase-table materialization test | drop the probe, run the gate package tests, expect the linked-repo phase-table expectation to fail |
