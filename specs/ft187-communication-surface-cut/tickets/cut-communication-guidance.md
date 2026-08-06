# Cut the communication guidance

Blocked by: pin-structured-phase-clause-set.md, register-current-communication-markers.md
Ownership fence: `.bench/BENCH.md`, `internal/anchors/registry_data.go`, `internal/anchors/registry_data_test.go`, `CHANGELOG.md`
Integration surfaces: revised shared rules→`.bench/BENCH.md`; Roles marker migration and final registry tuple→`internal/anchors/registry_data.go`; independent final tuple expectation→`internal/anchors/registry_data_test.go`; canonical-source consumer→existing unchanged `internal/conformance/validity_checks_test.go` + CG3; structured-phase consumer→existing unchanged `internal/conformance/docs_workflow_helpers_test.go` + CG1; user-visible record→`CHANGELOG.md`; cold semantic consumer→opus/low cross-harness review + CG1
Contracts: the final Roles marker string with type text, exact membership `Never assume the reviewer's decisions, and never assume a claim the gate could check instead`, registry order unchanged, and old-marker absence required crosses `.bench/BENCH.md`→`internal/anchors/registry_data.go`, asserted by CG3 against the real evaluator and final tuple expectation; the five fixed case→expected-action pairs cross `.bench/BENCH.md`→the opus/low reader, asserted by CG1
Closure: CG1/parallel-format, CG1/dependent-format, CG1/plain-acknowledgement, CG1/pre-phase-update, CG1/active-phase-update, CG1/main-session-line, CG1/delegate-line, CG1/lighter-path-approved, CG1/lighter-path-unapproved, CG2/story-one-shrink, CG2/story-two-shrink, CG2/story-three-shrink, CG2/story-four-shrink, CG2/story-five-shrink, CG2/story-six-shrink, CG2/no-additive-rule, CG3/final-roles-tuple, CG3/old-roles-absence, CG3/clear-tuple, CG3/workflow-tuple, CG3/registry-order, CG3/graded-root, CG4/changelog-current, CG4/changelog-single-entry

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
| CG1/parallel-format | restore the rule demanding a list or table for every multi-fact update | the fresh opus/low semantic review | run the parallel-facts prompt, expect the reader to lose the colleague-voice qualification |
| CG1/dependent-format | restore the rule demanding stacked one-sentence paragraphs | the fresh opus/low semantic review | run the dependent-facts prompt, expect the reader to prescribe the restored conflicting format or report unresolved precedence |
| CG1/plain-acknowledgement | restore labels for every one-sentence progress update | the fresh opus/low semantic review | run the routine-acknowledgement prompt, expect labels where the oracle requires plain prose |
| CG1/pre-phase-update | make structured-phase rules apply before a Bench command is invoked | the fresh opus/low semantic review | run the pre-invocation prompt, expect structured labels before the observable trigger |
| CG1/active-phase-update | remove `$bench-*` from the active-phase trigger | the fresh opus/low semantic review | run the post-invocation prompt with `$bench-*`, expect the reader not to recognize the active phase |
| CG1/main-session-line | claim the declaration can switch the running main model | the fresh opus/low semantic review | run the main-session prompt, expect the reader to prescribe an impossible model switch |
| CG1/delegate-line | remove delegates and headless runs from the binding scope | the fresh opus/low semantic review | run the delegated-work prompt, expect no binding line where the oracle requires one |
| CG1/lighter-path-approved | restore unconditional approval before every skipped phase | the fresh opus/low semantic review | run the both-observables prompt, expect an unnecessary approval request |
| CG1/lighter-path-unapproved | remove the fallback approval requirement | the fresh opus/low semantic review | run the failed-observable prompt, expect canonical phases to be skipped without approval |
| CG2/story-one-shrink | restore story 1's longer formatting passage | the story-scoped byte receipt | recompute passage 1, expect its after count not to shrink |
| CG2/story-two-shrink | restore story 2's longer Progress passage | the story-scoped byte receipt | recompute passage 2, expect its after count not to shrink |
| CG2/story-three-shrink | restore story 3's longer trigger passage | the story-scoped byte receipt | recompute passage 3, expect its after count not to shrink |
| CG2/story-four-shrink | restore story 4's longer Roles passage | the story-scoped byte receipt | recompute passage 4, expect its after count not to shrink |
| CG2/story-five-shrink | restore story 5's longer invariant passage | the story-scoped byte receipt | recompute passage 5, expect its after count not to shrink |
| CG2/story-six-shrink | restore story 6's longer lighter-path passage | the story-scoped byte receipt | recompute passage 6, expect its after count not to shrink |
| CG2/no-additive-rule | append a seventh standing guidance sentence while preserving all six edits | the before/after passage receipt | recompute the six passages and added-rule count, expect the additive surface to be rejected |
| CG3/final-roles-tuple | remove the final Roles registry row | the independent final tuple expectation | run `go test ./internal/anchors`, expect the missing-final-Roles-tuple failure |
| CG3/old-roles-absence | restore `NEVER assume, always verify` while retaining the final marker | the retired-marker registry prohibition | run `BENCH_CONFORMANCE_ROOT=$PWD go test ./internal/conformance -run '^TestRootConformance$' -count=1`, expect the forbidden-old-marker diagnostic |
| CG3/clear-tuple | alter the Clear tuple during the Roles migration | the independent final tuple expectation | run `go test ./internal/anchors`, expect the Clear tuple mismatch |
| CG3/workflow-tuple | alter the Workflow tuple during the Roles migration | the independent final tuple expectation | run `go test ./internal/anchors`, expect the Workflow tuple mismatch |
| CG3/registry-order | reorder the final communication tuples | the independent final tuple expectation | run `go test ./internal/anchors`, expect the tuple-order failure |
| CG3/graded-root | move the final Roles marker out of Roles | the real graded-root conformance check | run the rooted conformance test, expect the attributed section diagnostic |
| CG4/changelog-current | leave the existing Unreleased communication-rules entry unchanged | the synthesis consistency review | compare the user-visible diff to the Unreleased Changed entry, expect the stale description to be reported |
| CG4/changelog-single-entry | add a second Unreleased entry for the same cut | the synthesis consistency review | inspect the Unreleased section, expect the duplicate communication-rules entry to be reported |
