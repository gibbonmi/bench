# Retain check-level conformance evidence

Blocked by: derive-complete-contract-document-inputs.md, register-conformance-meta.md
Ownership fence: `internal/gate`, `internal/conformance/registry`, `internal/contract/runtime`
Assumptions: gate remains the sole evidence authority; exact bound resolvers optimize while uncertain checks use catch-all; claims re-derived from the tree at pickup

## What to build

Each ordinary conformance check receives a gate-owned content identity and durable slot;
green and red aggregate runs author or retire only executed checks, with every inability
to attribute widening execution.

## Acceptance

- [ ] [CE1] Each ordinary check identity binds name, tier, shared implementation closure, registry and function binding, declared inputs, canary ownership, and invocation schema.
- [ ] [CE2] A green run authors only executed check slots and leaves inherited slot bytes unchanged.
- [ ] [CE3] A red or interrupted run credits no newly executed check; a red aggregate retires every executed slot, while interruption before atomic persistence leaves no partial authorship.
- [ ] [CE4] Missing, malformed, wrong-check, wrong-tier, stale, or unresolved evidence runs the affected check or all checks when attribution is impossible.
- [ ] [CE5] Meta has no reusable slot and a conformance implementation change invalidates every ordinary check.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CE1 | omit one identity field at a time | check identity contract | hold that field changed with all others fixed, expect identity movement for every enumerated field |
| CE2 | re-stamp inherited slots after a mixed green | evidence persistence contract | run a mixed partition twice, compare inherited bytes, expect the mutation timestamp difference failure |
| CE3 | author slots before aggregate success or retain an executed red slot | interruption and red-state contracts | interrupt before final persistence and drive one red, expect no new credit and executed-slot retirement |
| CE4 | accept each malformed evidence class | evidence hostile-input table | inject every enumerated class, expect the corresponding check in the executed set |
| CE5 | author a meta slot or leave ordinary slots valid across shared implementation movement | meta/closure contract | seed the forged slot and mutate a shared helper, expect meta plus every ordinary check to execute |
