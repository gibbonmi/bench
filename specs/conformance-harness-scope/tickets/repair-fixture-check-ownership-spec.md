# Repair fixture-check ownership in the governing spec

Blocked by: expose-resolved-canary-fixture-checks.md, repair-fixture-bite-dependency-metadata.md
Ownership fence: `specs/conformance-harness-scope/spec.md`
Integration surfaces: fixture CHECK and family fallback→existing canary resolver + CR1/CR2; resolved `Fixture.Check`→story 1 implementation/testing decisions and CH1 + RS1; registry family rebind→CH1 consumer-propagation red signal + RS1
Contracts: raw fixture metadata and resolved canary `Fixture.Check` ownership cross the implementation and repair tickets→`specs/conformance-harness-scope/spec.md`, asserted by RS1 in the governing story, decisions, seam, and CH1 row; family registry rebind crosses canary fallback→the same governing spec path, asserted by CH1 without making conformance a second policy owner
Closure: RS1/owner-chain, RS1/coverage-row

## What to build

Repair the staged spec's governing ownership statements after the debug repair.
Every story, implementation decision, testing decision, seam description, and
coverage row describes canary as the sole CHECK-over-family policy owner and
`Fixture.Check` as the resolved value consumed by conformance. Preserve CH1's
registry-rebind mutation as evidence that a fixture using family fallback follows
the live registry; do not describe conformance as calling `registry.FamilyCheck`
to reconstruct scope.

## Acceptance

- [ ] [RS1] (covers CH1) the governing spec consistently assigns fixture CHECK/family precedence to canary, names resolved `Fixture.Check` as the conformance input, and retains registry rebind only as a fallback-propagation proof, with no statement that the conformance helper directly owns family-to-check resolution.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RS1/owner-chain | restore one statement that conformance obtains scope directly from `registry.FamilyCheck` | the spec ownership residue audit | apply the stale wording, search the governing spec for direct `registry.FamilyCheck` consumption, require the residue search to fail, then restore the owner chain |
| RS1/coverage-row | make CH1 describe only family binding and omit fixture-level CHECK override | the CH1 coverage-row audit | apply the omission, compare CH1 with the real `default-branch-refabricated` CHECK and canary repair contract, require the row to disagree, then restore the resolved-value behavior |
