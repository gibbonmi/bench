# Add the canary real-runner junction test

Blocked by: none
Ownership fence: `internal/canary/runner_junction_test.go`
Contracts: the gate invocation (script path, cwd, exit code, output) crosses `internal/canary`→bash through the real `defaultRunner`, asserted by CR1 against the real subprocess
Assumptions: production is presumed correct — a red against unmodified production is a stop-and-surface finding; `resolvingRunner` and the fake-driven sweep breadth stay; the fixture is one gate-owned hermetic temp repo with a planted trivial bash gate, honoring the package's nested GOMAXPROCS=2 constraint; claims re-derived from the tree at pickup

## What to build

One junction test that drives `SweepTier` (or the narrowest public sweep
surface that reaches the runner) with the real `defaultRunner` against a
hermetic fixture: a relative gate path must resolve against the fixture cwd
with real bash, and a missing gate must surface the exit-127 diagnostic.

## Acceptance

- [ ] [CR1] the sweep with the real `defaultRunner` runs the planted bash gate via its cwd-relative path and reports the gate's real exit verdict through the sweep result.
- [ ] [CR2] pointing the fixture at a missing gate path surfaces the exit-127 no-such-file diagnostic through the sweep result.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CR1 | make the planted gate script non-executable-shaped (rename it away) | the junction test | apply the mutation, run `go test ./internal/canary -run <CR1 test>`, expect the asserted verdict to fail |
| CR2 | restore the gate at the resolved path | the junction test | apply the mutation, run the CR2 test, expect the exit-127 assertion to fail because the gate now runs |
