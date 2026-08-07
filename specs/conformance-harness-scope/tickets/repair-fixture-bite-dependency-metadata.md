# Repair fixture-bite dependency metadata

Blocked by: expose-resolved-canary-fixture-checks.md
Ownership fence: `specs/conformance-harness-scope/tickets/scope-direct-conformance-fixture-bites.md`
Integration surfaces: exported resolved `Fixture.Check`→scope-direct-conformance-fixture-bites.md + RM1; CHECK-over-family precedence→existing expose-resolved-canary-fixture-checks.md + CR1/CR2, consumed without a second derivation by RM1
Contracts: resolved check name crosses expose-resolved-canary-fixture-checks.md→`specs/conformance-harness-scope/tickets/scope-direct-conformance-fixture-bites.md`, asserted by RM1 against the dependent ticket's blocker, surface, contract, and current-state build description
Closure: RM1/blocker-edge, RM1/resolved-check-contract

## What to build

Repair the durable consumer-ticket metadata after the canary precedence repair.
The fixture-bite ticket names `expose-resolved-canary-fixture-checks.md` as its
blocker and describes its existing implementation as consuming the resolved
`Fixture.Check`. It does not rederive CHECK-over-family precedence or claim a
direct `registry.FamilyCheck` lookup at the conformance seam.

## Acceptance

- [ ] [RM1] (covers local) `scope-direct-conformance-fixture-bites.md` names the canary repair ticket as its blocker and records `Fixture.Check` as the resolved cross-fence value consumed by the refreshed assignment, while leaving precedence owned only by the canary ticket.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RM1/blocker-edge | restore `Blocked by: none` on the consumer ticket | the ticket metadata audit | apply the omission, inspect the dependent ticket, require the producer basename to be absent, then restore the edge |
| RM1/resolved-check-contract | describe the consumer as deriving its check from `Fixture.Family` through a direct `registry.FamilyCheck` lookup | the ticket metadata audit | apply the stale contract, compare it with the integrated `Fixture.Check` consumer, require the metadata to disagree with the current tree, then restore the resolved-check contract |
