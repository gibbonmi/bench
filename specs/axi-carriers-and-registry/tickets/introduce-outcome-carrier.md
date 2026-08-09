# Introduce the outcome carrier

Blocked by: none
Ownership fence: `internal/axi`
Integration surfaces: outcome policy API→`internal/axi`; registry declarations→declare-production-axi-registry.md
Contracts: domain kind string, exit-policy map, typed payload, aggregate list, and action list cross caller→`internal/axi`, membership is owner-declared, ordering is payload-aggregates-actions, and absence is explicit empty slices, asserted by OC1
Closure: OC1/kind, OC1/exit, OC1/payload, OC1/order, OC1/absence

## What to build

shared outcomes preserve domain kinds, payloads, ordering, absence, and exact exit policies.

## Acceptance

- [ ] [OC1] (covers CR1) shared outcomes preserve domain kinds, payloads, ordering, absence, and exact exit policies.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| OC1/kind | accept an undeclared kind | outcome policy test | construct the kind and require refusal |
| OC1/exit | replace a specialized policy with 0/1/2 | outcome policy test | construct the specialized kind and require its exact exit |
| OC1/payload | erase the typed payload | outcome round-trip test | construct and require the original typed value |
| OC1/order | render actions before aggregates | outcome round-trip test | enumerate components and require declared order |
| OC1/absence | replace absent aggregates with an implicit nil default | outcome validation test | construct explicit absence and require it remain distinguishable |

