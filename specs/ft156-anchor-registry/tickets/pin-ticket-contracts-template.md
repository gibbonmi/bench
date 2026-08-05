# Pin the ticket Contracts template

Blocked by: none
Ownership fence: `internal/conformance/docs_workflow_helpers_test.go`, `internal/conformance/registry_test.go`, `tests/canary/workflow-guidance-anchors/ticket-contracts-template-anchor`
Contracts: the `Contracts:` template anchor tuple crosses `internal/conformance/docs_workflow_helpers_test.go`→the mutated craft-tickets fixture, while the fixture family classification crosses `tests/canary/workflow-guidance-anchors/ticket-contracts-template-anchor`→`internal/conformance/registry_test.go`; CT1 asserts both against the real graded-root diagnostic and complete fixture registry

## What to build

Require the `Contracts:` line inside craft-tickets' "Write one file per ticket" section and add a canary fixture that deletes only that line, closing the recorded false green.

## Acceptance

- [ ] [CT1] Deleting only the ticket template's `Contracts:` line fails the graded root with a specific missing-template diagnostic.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CT1 | omit the new section-scoped require while retaining the deletion fixture | the workflow-guidance canary fixture | run the owning fixture before the anchor exists, observe that its expected diagnostic is absent; add the anchor and rerun green |
