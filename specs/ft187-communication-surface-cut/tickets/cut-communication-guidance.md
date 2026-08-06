# Cut the communication guidance

Blocked by: pin-structured-phase-clause-set.md, register-current-communication-markers.md
Ownership fence: `.bench/BENCH.md`, `internal/anchors/registry_data.go`, `internal/anchors/registry_data_test.go`, `CHANGELOG.md`
Integration surfaces: revised shared rules→`.bench/BENCH.md`; Roles marker migration and final registry tuple→`internal/anchors/registry_data.go`; independent final tuple expectation→`internal/anchors/registry_data_test.go`; canonical-source consumer→existing unchanged `internal/conformance/validity_checks_test.go` + CG3; structured-phase consumer→existing unchanged `internal/conformance/docs_workflow_helpers_test.go` + CG1; user-visible record→`CHANGELOG.md`; cold semantic consumer→opus/low cross-harness review + CG1
Contracts: the final Roles marker string with type text, exact membership `Never assume the reviewer's decisions, and never assume a claim the gate could check instead`, registry order unchanged, and old-marker absence required crosses `.bench/BENCH.md`→`internal/anchors/registry_data.go`, asserted by CG3 against the real evaluator and final tuple expectation; the five fixed case→expected-action pairs cross `.bench/BENCH.md`→the opus/low reader, asserted by CG1

## What to build

Apply stories 1–6 as six deletions or shorter replacements, atomically migrate the Roles marker constant and its independent tuple expectation to the approved final sentence, and extend the existing Unreleased communication-rules changelog entry. Preserve the spec's five-case semantic oracle for the later fresh opus/low review.

## Acceptance

- [ ] [CG1] A fresh opus/low reader with no spec or conversation context returns the specified expected action for each of the five case-table prompts.
- [ ] [CG2] The build records six story-scoped before/after byte counts; every passage shrinks or deletes text and the six passages add no seventh guidance rule.
- [ ] [CG3] The final Roles, How to talk to me, and Workflow tuples match the spec exactly, the old Roles marker is absent, and focused graded-root conformance is green.
- [ ] [CG4] The existing Unreleased communication-rules changelog entry describes the reduced, conflict-free surface without creating a second entry for the same change.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CG1 | restore one old contradiction at a time in the reader's subject guide | the fresh opus/low semantic review | apply each restoration to the subject, run the five fixed prompts, expect at least its mapped action to diverge or become unresolved |
| CG2 | append a seventh standing guidance sentence while preserving all six edited markers | the before/after passage receipt | add the sentence, recompute the six passages and added-rule count, expect the scope receipt to reject the additive surface |
| CG3 | restore `NEVER assume, always verify` in Roles while retaining the final exported marker | the real graded-root conformance check | restore only the old sentence, run `BENCH_CONFORMANCE_ROOT=$PWD go test ./internal/conformance -run '^TestRootConformance$' -count=1`, expect the final Roles-marker diagnostic |
| CG4 | leave the existing Unreleased communication-rules entry unchanged | the synthesis consistency review | compare the user-visible diff to the Unreleased Changed entry, expect the stale description to be reported |
