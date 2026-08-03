# Derive complete contract document inputs

Blocked by: seal-outer-conformance-selection.md
Ownership fence: `internal/gate/component_inputs.go`, `internal/gate/component_inputs_test.go`, `internal/packagesurface`, `internal/contract`, `projects/benchkit.md`
Assumptions: FT183 resolver-identity work landed at `6a3ea99`; lifecycle contracts consume managed kit assets through an authoritative inventory; `.bench/BENCH.md` is one observed contract input; claims re-derived from the tree at pickup

## What to build

The contract component derives every managed document it reads from an authoritative
consumer inventory, so changing any such document moves its identity and an unresolved
inventory runs the component.

## Acceptance

- [ ] [CI1] Mutating `.bench/BENCH.md` moves the contract component identity.
- [ ] [CI2] Every managed document read by lifecycle contracts is derived from the consumer inventory rather than a copied path family.
- [ ] [CI3] Missing, malformed, or unresolved inventory input runs the contract component instead of granting a skip.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CI1 | remove `.bench/BENCH.md` from the derived identity | component identity contract | mutate only the guide in a fixture and expect the unchanged-identity failure |
| CI2 | restore the `.agents/`-only copied enumeration | package-surface agreement test | enumerate the real contract consumers and expect the missing managed-asset diagnostic |
| CI3 | convert resolver failure into an empty path set | component decision contract | inject an unreadable inventory source and expect contract in the executed set |
