# Pin the structured-phase clause set

Blocked by: none
Ownership fence: `internal/conformance/docs_workflow_helpers_test.go`, `internal/conformance/validity_checks_test.go`
Integration surfaces: structured-phase declaration and body parser→`internal/conformance/docs_workflow_helpers_test.go`; inactive-guidance unit cases→`internal/conformance/validity_checks_test.go`; real Progress deletion canary and `workflow-guidance-anchors` family→existing unchanged `tests/canary/workflow-guidance-anchors/structured-phase-progress-anchor` and `internal/conformance/registry_test.go` + PC2
Contracts: the ordered `[]string` clause set with membership exactly `progress`, `exit`, `omission`, and `cohesion`, declaration order preserved, and any missing/duplicate/unknown name invalid crosses `internal/conformance/docs_workflow_helpers_test.go` declaration parsing→body validation, asserted by PC1 against table-driven real guide text
Closure: PC1/missing-name, PC1/zero-names, PC1/duplicate-name, PC1/unknown-name, PC1/renamed-name, PC1/declaration-order, PC2/commented-body, PC2/quoted-body, PC2/negated-body, PC2/wrong-section-body, PC2/progress-canary

## What to build

Make the bespoke structured-phase parser validate the fixed four-name contract instead of deriving its own requirements from whichever backticks remain in the document. Keep the existing active-body and canary behavior intact.

## Acceptance

- [ ] [PC1] Deleting one or all required names with their bodies, duplicating, adding, renaming, or reordering names fails `checkStructuredPhaseContract` with an attributed clause diagnostic.
- [ ] [PC2] The existing commented, quoted, negated, wrong-section, and real Progress-deletion cases retain their targeted failures.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PC1/missing-name | delete `cohesion` from both the declaration and its body | the exact-set table test | mutate the guide text, run `go test ./internal/conformance -run '^TestStructuredPhaseContract'`, expect the missing-cohesion diagnostic |
| PC1/zero-names | delete all four names and bodies | the exact-set table test | mutate the guide text, run the focused test, expect the exact-set diagnostic rather than an empty green |
| PC1/duplicate-name | duplicate `progress` in the declaration | the exact-set table test | mutate the declaration, run the focused test, expect the duplicate-progress diagnostic |
| PC1/unknown-name | add `handoff` to the declaration and body | the exact-set table test | mutate the guide text, run the focused test, expect the unknown-handoff diagnostic |
| PC1/renamed-name | rename `progress` to `status` in the declaration and body | the exact-set table test | mutate the guide text, run the focused test, expect missing-progress and unknown-status diagnostics |
| PC1/declaration-order | swap `progress` and `exit` in the declaration | the exact-set table test | mutate the declaration, run the focused test, expect the declaration-order diagnostic |
| PC2/commented-body | move a clause body into an HTML comment | the inactive-guidance table test | mutate the guide text, run the focused test, expect the commented-body diagnostic |
| PC2/quoted-body | leave a clause only in quoted guidance | the inactive-guidance table test | mutate the guide text, run the focused test, expect the quoted-body diagnostic |
| PC2/negated-body | leave only a negated clause body | the inactive-guidance table test | mutate the guide text, run the focused test, expect the negated-body diagnostic |
| PC2/wrong-section-body | move a clause body outside the structured-phase section | the inactive-guidance table test | mutate the guide text, run the focused test, expect the wrong-section diagnostic |
| PC2/progress-canary | restore a live Progress body to the `structured-phase-progress-anchor` subject | the existing canary fixture sweep | apply the restored subject, run the owning fixture test, expect the canary to red because its targeted diagnostic disappears |
