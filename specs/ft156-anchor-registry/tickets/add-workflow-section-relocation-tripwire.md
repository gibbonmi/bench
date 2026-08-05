# Add the Workflow-section relocation tripwire

Blocked by: none
Ownership fence: `internal/conformance/registry_test.go`, `tests/canary/workflow-guidance-anchors/fix-dont-park-section-relocated`
Contracts: the mutated `.bench/BENCH.md` fixture and exact expected diagnostic cross `tests/canary/workflow-guidance-anchors/fix-dont-park-section-relocated`→the real graded-root conformance runner, while the family classification crosses that fixture→`internal/conformance/registry_test.go`; WF1 asserts both against the owning bite and complete fixture registry

## What to build

Add a `workflow-guidance-anchors` fixture that keeps the fix-don't-park sentence in `.bench/BENCH.md` but moves it outside the Workflow section, proving the generic section-scoped helpers own placement.

## Acceptance

- [ ] [WF1] The fixture's mutated graded root fails with the exact Workflow-section placement diagnostic while the sentence remains elsewhere in the file.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| WF1 | make the fixture copy retain the sentence inside the Workflow section | the canary fixture sweep | apply the fixture, run the `workflow-guidance-anchors` family, expect the vacuous fixture to be rejected because the targeted diagnostic disappears |
