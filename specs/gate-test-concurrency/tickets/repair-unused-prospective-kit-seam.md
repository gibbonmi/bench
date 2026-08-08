# Remove the unused prospective kit seam

Blocked by: none
Ownership fence: `internal/gate/engine.go`, `internal/gate/prospective_test.go`, `specs/gate-test-concurrency/tickets/inject-kit-root-below-entries.md`, `specs/gate-test-concurrency/tickets/retire-fixture-kit-pins.md`
Integration surfaces: `internal/gate/engine.go` prospective entry→`internal/gate/prospective_test.go` + RPS1; `specs/gate-test-concurrency/tickets/inject-kit-root-below-entries.md` boundary wording and `specs/gate-test-concurrency/tickets/retire-fixture-kit-pins.md` exemption wording→RPS1
Contracts: prospective execution in `internal/gate/engine.go` has no private kit parameter because current prospective scoping returns empty before kit identity can be consumed
Closure: RPS1/prospective-seam

## What to build

Remove `executeTreeAtKit`: prospective evaluation returns empty scoping before
kit identity can be consumed, so its kit parameter has no present semantic
effect. Route the one fixture call back through `ExecuteTree`, and correct the
two source tickets to state that bounded prospective exemption. Do not retain
a parameter solely for a hypothetical future consumer.

## Acceptance

- [ ] [RPS1] (covers local) the prospective path carries no unused kit parameter and the source tickets describe the present boundary accurately, repairing S3.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RPS1/prospective-seam | reintroduce `executeTreeAtKit` with a kit parameter that has no consumer | the zero-unused-carrier search | apply, run the search, expect it to name the reintroduced declaration and call |
