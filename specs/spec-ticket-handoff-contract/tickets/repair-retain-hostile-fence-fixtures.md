# Retain hostile ownership-fence fixtures in the aggregate

Blocked by: none
Ownership fence: `internal/conformance/fixture_bite_test.go`
Integration surfaces: hostile-fence fixture identities and diagnostics→`TestSpecTicketHandoffWorkflowFixturesAreComplete`
Contracts: the existing aggregate consumes registered fixture names and exact diagnostics; type is a static required-fixture member, domain adds the two accepted hostile-fence repair fixtures, order follows the authored-spec clause before its empty/invalid refinement, absence fails the aggregate
Closure: RF1/exact-literal-membership, RF2/empty-invalid-membership

## What to build

Close accepted review finding `STD-001-repair-fixture-inventory` by adding the two hostile ownership-fence repair fixtures to the existing handoff aggregate. Do not change the fixture runner, normative guidance, parser, or lifecycle authority.

## Acceptance

- [ ] [RF1] (covers SH11) the complete handoff inventory requires the `craft-spec-exact-literal-fence` fixture with its exact registered diagnostic.
- [ ] [RF2] (covers SH11) the same inventory independently requires the `craft-spec-empty-or-invalid-fence` fixture with its exact registered diagnostic.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RF1/exact-literal-membership | remove only the exact-literal fixture member from the required list | `TestSpecTicketHandoffWorkflowFixturesAreComplete` | apply the omission, run the named aggregate, require the missing-fixture red, restore the member |
| RF2/empty-invalid-membership | remove only the empty-or-invalid fixture member from the required list | `TestSpecTicketHandoffWorkflowFixturesAreComplete` | apply the omission, run the named aggregate, require the missing-fixture red, restore the member |
