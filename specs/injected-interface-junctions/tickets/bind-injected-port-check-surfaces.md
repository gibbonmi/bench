# Bind the injected-port check to its execution and advertisement surfaces

Blocked by: add-injected-port-conformance-check.md
Ownership fence: `internal/conformance/checks_test.go`, `internal/conformance/registry_test.go`, `projects/benchkit.md`
Contracts: the check name `injected-port-registry` crosses `internal/conformance/registry`→`internal/conformance` (executable binding) and →`projects/benchkit.md` (check-input advertisement), asserted by BP1 against the real scoped conformance run
Assumptions: the check, registry row, canary fixture, and pins landed in the blocker ticket; this ticket adds only the three binding lines the blocker's fence could not hold; claims re-derived from the tree at pickup

## What to build

The three bindings that make the already-landed check runnable and advertised:
the name-to-function row in the conformance executable table, the canary
fixture's family registration, and the profile's check-input advertisement row.

## Acceptance

- [ ] [BP1] the scoped root-conformance run (`BENCH_CONFORMANCE_ROOT` set, `BENCH_CONFORMANCE_CHECK=injected-port-registry`) executes the check and reports green on the composed tree.
- [ ] [BP2] the canary fixture registry test classifies `tests/canary/injected-ports` under its registered family.
- [ ] [BP3] the conformance-meta check accepts the profile's check-input row for `injected-port-registry`.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| BP1 | remove the binding row | the scoped conformance run | run the scoped check, expect "registered with no bound function" |
| BP2 | remove the family registration | the fixture-classification test | run it, expect the unclassified-fixture failure |
| BP3 | remove the profile row | the conformance-meta check | run it scoped, expect the stale check-input diagnostic |
