# Author identified spec handoffs

Blocked by: none
Ownership fence: `.agents/commands/bench-write-spec.md`, `.agents/skills/bench-craft-spec/SKILL.md`, `internal/anchors/registry_data.go`, `internal/conformance/fixture_bite_test.go`, `internal/conformance/docs_workflow_helpers_test.go`, `tests/canary/workflow-guidance-anchors`
Integration surfaces: authored coverage-row and ownership-fence contract→derive-ticket-evidence-and-handoff-ledger.md; workflow predicates→`internal/anchors/registry_data.go`; fixture owner→`internal/conformance/fixture_bite_test.go`; focused helper→`internal/conformance/docs_workflow_helpers_test.go`; omission fixtures→`tests/canary/workflow-guidance-anchors`
Contracts: identified coverage rows and approved ownership fences cross `.agents/skills/bench-craft-spec/SKILL.md`→`.agents/skills/bench-craft-tickets/SKILL.md`; type is Markdown row and literal path declarations, domain is unique spec-local row IDs plus exact repo-relative files or prefixes, order is approved map and fence order, absence keeps legacy five-column specs valid but makes a new empty fence incomplete, asserted by AS1, AS3, and AS4 against the real authored template
Closure: AS1/six-column-default, AS2/legacy-five-column, AS3/ownership-fence-section, AS4/reviewer-fence-approval

## What to build

Newly authored specs expose identified coverage rows and exact ownership fences on the reviewer approval surface while the exported parser retains legacy five-column compatibility. This producer lands before ticket derivation because the consumer needs the approved row and fence vocabulary.

## Acceptance

- [ ] [AS1] (covers SH1) the authored spec template leads its acceptance map with unique row IDs.
- [ ] [AS2] (covers SH2) `TestParseSpecOptIn` remains green for six-column opt-in and legacy five-column maps.
- [ ] [AS3] (covers SH3) the authored template includes an ownership-fence section pointing to `craft-spec` as its rule owner.
- [ ] [AS4] (covers SH4) the approval table requires an explicit ownership-fence disposition.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AS1/six-column-default | replace the template's leading `row` header with the legacy five-column header | the registered section-sensitive workflow anchor | apply the subject swap, run the focused workflow-anchor conformance test, require the six-column-default diagnostic, restore the subject |
| AS2/legacy-five-column | remove the exported parser's legacy five-column header branch temporarily | `TestParseSpecOptIn` | apply the subject omission, run `go test -count=1 -run '^TestParseSpecOptIn$' ./internal/coverage`, require its legacy-path red, restore the subject |
| AS3/ownership-fence-section | remove the template's `## Ownership fences` section while leaving `craft-spec` unchanged | the registered section-sensitive workflow anchor | apply the subject omission, run the focused workflow-anchor conformance test, require the ownership-fence-section diagnostic, restore the subject |
| AS4/reviewer-fence-approval | relocate ownership fences outside the approval paragraph while preserving the words | the registered section-sensitive workflow anchor | apply the subject relocation, run the focused workflow-anchor conformance test, require the reviewer-fence-approval diagnostic, restore the subject |
