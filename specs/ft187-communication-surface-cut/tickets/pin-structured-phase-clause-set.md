# Pin the structured-phase clause set

Blocked by: none
Ownership fence: `internal/conformance/docs_workflow_helpers_test.go`, `internal/conformance/validity_checks_test.go`
Integration surfaces: structured-phase declaration and body parser→`internal/conformance/docs_workflow_helpers_test.go`; inactive-guidance unit cases→`internal/conformance/validity_checks_test.go`; real Progress deletion canary and `workflow-guidance-anchors` family→existing unchanged `tests/canary/workflow-guidance-anchors/structured-phase-progress-anchor` and `internal/conformance/registry_test.go` + PC2
Contracts: the ordered `[]string` clause set with membership exactly `progress`, `exit`, `omission`, and `cohesion`, declaration order preserved, and any missing/duplicate/unknown name invalid crosses `internal/conformance/docs_workflow_helpers_test.go` declaration parsing→body validation, asserted by PC1 against table-driven real guide text

## What to build

Make the bespoke structured-phase parser validate the fixed four-name contract instead of deriving its own requirements from whichever backticks remain in the document. Keep the existing active-body and canary behavior intact.

## Acceptance

- [ ] [PC1] Deleting any one required name with its body, duplicating a name, or adding an unknown name fails `checkStructuredPhaseContract` with an attributed clause diagnostic.
- [ ] [PC2] The existing commented, quoted, negated, wrong-section, and real Progress-deletion cases retain their targeted failures.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PC1 | delete `cohesion` from both the declaration and its body in the subject guide | the exact-set table test | mutate the guide text, run `go test ./internal/conformance -run '^TestStructuredPhaseContract'`, expect the missing-cohesion diagnostic |
| PC2 | restore a live Progress body to the `structured-phase-progress-anchor` subject | the existing canary fixture sweep | apply the restored subject, run the owning fixture test, expect the canary to red because its targeted diagnostic disappears |
