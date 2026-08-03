# Add the injected-port conformance check, registry, and wiring pins

Blocked by: add-releaseowner-junction.md, redrive-abandon-decayed-family.md, add-canary-runner-junction.md, add-gitguard-checker-junction.md, pin-shapeunknown-fixtures.md
Ownership fence: `internal/conformance/injected_ports_test.go`, `internal/conformance/registry/registry.go`, `internal/conformance/registry/packages.go`, `cmd/bench/wiring_pins.go`, `tests/canary/injected-ports`
Contracts: registry rows naming real-producer tests cross `internal/conformance`→every audited package and must match tests the tree actually holds, asserted by IP2 against the real test files; the check's four failure messages cross into the gate transcript, asserted by IP1-IP4
Assumptions: a port is any port-shaped value — interface, func type, or struct of func fields — injected via constructor or exported-function parameter, so the derivation sees `canary.Runner` and `gitguard.Checker`; the named-test-exists verification is a labeled tripwire, not a behavior check; registry rows carry the audit's exemption citations for the priced-out ports; the registry `Check` row position fixes execution order per the registry's order-is-contract rule; junction tests from the blocker tickets are re-derived from the tree, never from this file; claims re-derived from the tree at pickup

## What to build

A dev-tier conformance check over a single-source injected-port registry:
every derived port names its real-producer junction test or its exemption
reason; four distinct fail-closed messages; a canary fixture keeping the
unregistered-port red alive; and the three compile-time wiring pins
(`PromotionGateOwner`, `ReleaseOwner`, `AbandonOwner` against the production
owners) in cmd/bench.

## Acceptance

- [ ] [IP1] an injected port derived from the audited packages with no registry row fails the check with the unregistered-port message, and the `tests/canary/injected-ports` fixture retains that red in the canary sweep.
- [ ] [IP2] a registry row naming a test absent from the tree fails with the missing-test message (tripwire, labeled as such in the check's doc).
- [ ] [IP3] an exempt row with an empty reason fails with the empty-exemption message.
- [ ] [IP4] a derivation returning zero ports for a package known to declare them fails closed with the zero-inventory message.
- [ ] [IP5] the three compile-time pins land in cmd/bench and removing any pinned method breaks the build.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| IP1 | delete one port's registry row | the canary fixture | run the canary sweep (or the check scoped via `BENCH_CONFORMANCE_CHECK`), expect the unregistered-port message |
| IP2 | rename a registry row's test to a name the tree lacks | the check's unit fixture | run the check, expect the missing-test message |
| IP3 | blank one exempt row's reason | the check's unit fixture | run the check, expect the empty-exemption message |
| IP4 | point the derivation at a package list excluding a port-declaring package | the check's unit fixture | run the check, expect the zero-inventory message |
| IP5 | remove `Validate` from `productionGateOwner` in a scratch tree | the Go compiler | `go build ./cmd/bench`, expect the pin's compile error; revert |
