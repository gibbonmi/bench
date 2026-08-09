# Require approved handoff before lifecycle

Blocked by: derive-ticket-evidence-and-handoff-ledger.md
Ownership fence: `.agents/commands/bench-implement-spec.md`, `projects/benchkit.md`, `CHANGELOG.md`, `internal/anchors/registry_data.go`, `internal/conformance/fixture_bite_test.go`, `internal/conformance/docs_workflow_helpers_test.go`, `tests/canary/workflow-guidance-anchors`
Integration surfaces: accountable ledger producer→derive-ticket-evidence-and-handoff-ledger.md; lifecycle pre-build route→`.agents/commands/bench-implement-spec.md`; profile advertisement→`projects/benchkit.md`; user-visible change→`CHANGELOG.md`; workflow predicates→`internal/anchors/registry_data.go`; fixture owner→`internal/conformance/fixture_bite_test.go`; focused helper→`internal/conformance/docs_workflow_helpers_test.go`; omission fixtures→`tests/canary/workflow-guidance-anchors`
Contracts: the approved-ledger prerequisite crosses `.agents/commands/bench-implement-spec.md`→`internal/anchors/registry_data.go`; type is a pre-lifecycle guidance predicate, domain is seam drift, fence drift, repaired approval, and complete ledger, order is confirm then derive and review before `start`, absence refuses lifecycle entry, asserted by LC1 and LC2 against the real command through the workflow-guidance owner
Closure: LC1/lifecycle-prerequisite, LC2/sh1-six-column-default, LC2/sh3-ownership-fence-section, LC2/sh4-reviewer-fence-approval, LC2/sh5-observed-red-route, LC2/sh6-already-covered-route, LC2/sh7-not-tdd-able-route, LC2/sh8-ledger-totality, LC2/sh9-fence-drift-stop, LC2/sh10-lifecycle-prerequisite-anchor

## What to build

The implementation phase treats fence drift like seam drift, requires repaired approval plus the complete handoff ledger before lifecycle start, registers every new normative clause with an independent section-sensitive mutation, advertises the current contract, and records the user-visible change.

## Acceptance

- [ ] [LC1] (covers SH10) `$bench-implement-spec` routes fence drift through repaired spec approval and refuses lifecycle start without the complete handoff ledger.
- [ ] [LC2] (covers SH11) SH1 and SH3–SH10 each have their own registered section-sensitive omission or relocation mutation and clause-specific diagnostic.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| LC1/lifecycle-prerequisite | relocate the complete-ledger prerequisite below lifecycle start while retaining its words | the registered section-sensitive workflow anchor | apply the subject relocation, run the focused workflow-anchor conformance test, require the lifecycle-prerequisite diagnostic, restore the subject |
| LC2/sh1-six-column-default | replace the template's identified header inside `Template` with the legacy header | the workflow-guidance canary fixture family | apply the SH1 subject swap through the public canary owner, require the six-column-default diagnostic, restore the subject and require green |
| LC2/sh3-ownership-fence-section | relocate the ownership-fence template outside `Template` while retaining its words | the workflow-guidance canary fixture family | apply the SH3 subject relocation through the public canary owner, require the ownership-fence-section diagnostic, restore the subject and require green |
| LC2/sh4-reviewer-fence-approval | remove the ownership-fence disposition from the approval paragraph | the workflow-guidance canary fixture family | apply the SH4 subject omission through the public canary owner, require the reviewer-fence-approval diagnostic, restore the subject and require green |
| LC2/sh5-observed-red-route | replace the distinct subject-mutation requirement with the obsolete absence probe | the workflow-guidance canary fixture family | apply the SH5 subject swap through the public canary owner, require the observed-red-route diagnostic, restore the subject and require green |
| LC2/sh6-already-covered-route | remove the changed-route mutation while retaining the existing positive control | the workflow-guidance canary fixture family | apply the SH6 subject omission through the public canary owner, require the already-covered-route diagnostic, restore the subject and require green |
| LC2/sh7-not-tdd-able-route | retain mutation exemption after the blocker clears and the seam exists | the workflow-guidance canary fixture family | apply the SH7 subject swap through the public canary owner, require the not-tdd-able-route diagnostic, restore the subject and require green |
| LC2/sh8-ledger-totality | remove approved ownership fences from the handoff ledger's domain | the workflow-guidance canary fixture family | apply the SH8 subject omission through the public canary owner, require the ledger-totality diagnostic, restore the subject and require green |
| LC2/sh9-fence-drift-stop | replace the exact return-to-spec route with ticket-local fence widening | the workflow-guidance canary fixture family | apply the SH9 subject swap through the public canary owner, require the fence-drift-stop diagnostic, restore the subject and require green |
| LC2/sh10-lifecycle-prerequisite-anchor | relocate the complete-ledger prerequisite below lifecycle start | the workflow-guidance canary fixture family | apply the SH10 subject relocation through the public canary owner, require the lifecycle-prerequisite diagnostic, restore the subject and require green |
