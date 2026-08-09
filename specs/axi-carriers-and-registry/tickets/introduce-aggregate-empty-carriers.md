# Introduce aggregate and empty carriers

Blocked by: none
Ownership fence: `internal/axi`
Integration surfaces: aggregate/empty APIs→`internal/axi`; registry declarations→declare-production-axi-registry.md
Contracts: ordered name plus typed scalar and empty-class enum cross caller→`internal/axi`, scalar domain is string/number/bool/unknown, order is insertion order, and absence is explicit no-semantic-empty, asserted by AE1
Closure: AE1/name, AE1/order, AE1/type, AE1/unknown, AE1/zero, AE1/empty-class, AE1/no-empty

## What to build

aggregates and empty declarations preserve names, order, scalar types, unknowns, zeros, exact classes, and explicit no-empty.

## Acceptance

- [ ] [AE1] (covers CR4) aggregates and empty declarations preserve names, order, scalar types, unknowns, zeros, exact classes, and explicit no-empty.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AE1/name | erase one owner name | aggregate test | construct and require the name |
| AE1/order | sort facts lexically | aggregate test | supply reverse order and require insertion order |
| AE1/type | stringify a numeric fact | aggregate test | require numeric type |
| AE1/unknown | coerce unknown to zero | aggregate test | require distinct values |
| AE1/zero | omit an explicit zero | aggregate test | require zero remains present |
| AE1/empty-class | normalize two empty classes | empty test | construct both and require distinct enums |
| AE1/no-empty | default no-empty to zero-row | empty validation test | require explicit no-semantic-empty |

