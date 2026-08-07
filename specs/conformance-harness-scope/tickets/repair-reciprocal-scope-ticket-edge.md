# Repair the reciprocal scope-ticket edge

Blocked by: repair-fixture-bite-dependency-metadata.md
Ownership fence: `specs/conformance-harness-scope/tickets/scope-direct-conformance-fixture-bites.md`
Integration surfaces: durable blocker graph→owned consumer-ticket metadata + RG1; resolved-check metadata producer→existing repair-fixture-bite-dependency-metadata.md + RG1
Contracts: reciprocal producer identity crosses the integrated metadata repair→`specs/conformance-harness-scope/tickets/scope-direct-conformance-fixture-bites.md`, asserted by RG1 against the current blocker list and resolved-check contract
Closure: RG1/metadata-producer-edge

## What to build

Add `repair-fixture-bite-dependency-metadata.md` to the scoped fixture-bite
ticket's `Blocked by:` list alongside the canary producer. This makes the durable
ticket graph agree with the metadata repair's dependent surface and the assign
rule that enforces reciprocal edges. Preserve all existing behavior, ownership
fences, acceptance rows, and CHECK-over-family ownership language.

## Acceptance

- [ ] [RG1] (covers local) the scoped fixture-bite ticket names both the canary resolved-check producer and the metadata repair in `Blocked by:`, so every producer that names it as a dependent has one reciprocal edge.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RG1/metadata-producer-edge | remove the metadata repair basename from the consumer's blocker list | the reciprocal-edge assignment control | apply the omission, assign the consumer ticket against the current sibling set, require the missing reciprocal-edge refusal, then restore the blocker |
