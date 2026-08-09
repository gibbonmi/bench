# Retain the unique row-ID template contract

Blocked by: none
Ownership fence: `.agents/commands/bench-write-spec.md`, `internal/anchors/registry_data.go`, `internal/conformance/fixture_bite_test.go`, `tests/canary/workflow-guidance-anchors`
Integration surfaces: unique row-ID template cell→`internal/anchors/registry_data.go`; mutation identity→`tests/canary/workflow-guidance-anchors`; retained membership→`internal/conformance/fixture_bite_test.go`
Contracts: the unique row-ID template predicate crosses `.agents/commands/bench-write-spec.md`→`internal/anchors/registry_data.go`→`tests/canary/workflow-guidance-anchors`; type is a section-scoped Markdown table cell, domain is unique spec-local IDs, order is the leading `row` cell, absence fails the named fixture and aggregate inventory
Closure: UR1/unique-spec-local-cell, UR2/retained-membership

## What to build

Close accepted review finding `SPEC-001-unique-row-id-mutation` by giving SH1's unique spec-local row-ID template cell its own registered section-sensitive mutation fixture and retaining that fixture in the existing aggregate. Do not change parser or lifecycle authority.

## Acceptance

- [ ] [UR1] (covers SH1) the template's leading row cell independently requires a unique spec-local ID through an exact section anchor and biting mutation fixture.
- [ ] [UR2] (covers SH11) the complete handoff inventory retains that new fixture and its exact registered diagnostic.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| UR1/unique-spec-local-cell | replace `<unique spec-local ID>` with `<row ID>` while leaving the six-column header intact | the registered section-sensitive workflow anchor | apply the swap, run the focused workflow fixture, require the unique-row-ID diagnostic, restore the cell |
| UR2/retained-membership | remove only the unique-row-ID fixture member from the aggregate list | `TestSpecTicketHandoffWorkflowFixturesAreComplete` | apply the omission, run the named aggregate, require its cardinality red, restore the member |
