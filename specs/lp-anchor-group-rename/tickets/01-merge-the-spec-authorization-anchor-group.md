# Merge the spec-authorization anchor group into the implement-spec group

Blocked by: none
Writes: internal/anchors/registry.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/conformance/docs_workflow_helpers_test.go, capture/restructure-backlog.md
Covers: none

## What to build

Drain d-359accf567c1, learning L107, first half.

The anchor group `AfterSpecAuthorization` has the name of `checkSpecAuthorizationContract`.
An earlier ticket deleted that check.
A group places registry evaluations around bespoke conformance checks.
`checkWorkflowAnchors` evaluates `AfterSpecAuthorization` immediately after `AfterImplementSpec`, and no bespoke check runs between the two evaluations.
Thus the two groups hold one evaluation position, and the old group has no function of its own.

Move each `AfterSpecAuthorization` row and test harness to `AfterImplementSpec`.
Remove the `AfterSpecAuthorization` constant.
Remove its evaluation from `checkWorkflowAnchors`.
Remove the stale group name from the restructure backlog.

Each anchor row keeps its file, kind, class, section, needle, and diagnostic.
For each tree, the gate gives the same set of distinct diagnostics.
In a report with more than one diagnostic, the order can change, because the merged group evaluates its rows in registry order.
A refused file or a missing section that has rows in both old groups now gives its diagnostic once, not twice.
This agrees with the evaluator's rule of one refusal for each file.

The bite cases for the two kept pins (the second half of L107) are not in this ticket.

## Acceptance

- [ ] Outside this ticket, no Go source, test, or tracked document names `AfterSpecAuthorization`.
- [ ] Each of the 20 moved rows is in the `AfterImplementSpec` group and is otherwise unchanged.
- [ ] `checkWorkflowAnchors` evaluates each registry group one time.
- [ ] The anchors package tests and the conformance check pass with no change to an expectation other than the group identifier.
