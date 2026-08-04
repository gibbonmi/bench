# Retire the Assumptions field from the grammar surface and the records

Blocked by: retire-assumptions-machinery.md
Ownership fence: `.agents/skills/bench-craft-tickets/SKILL.md`, `internal/conformance/example_agreement_test.go`, `internal/conformance/fixture_bite_test.go`, `internal/conformance/docs_workflow_helpers_test.go`, `decisions/parallel-session-landings.md`
Contracts: the ticket-grammar template crosses `.agents/skills/bench-craft-tickets/SKILL.md`→the example-agreement, fixture-bite, and docs-workflow anchors in `internal/conformance`, asserted by RS1 and RS2 through the real gate rather than a spot-check of one file

## What to build

The `craft-tickets` template, the `Assumptions:` field bullet, and the taught
example drop the line. The `example-agreement` suite's independently-authored
expectation literals, and the template anchors in the fixture-bite and
docs-workflow suites, move in this same change — the grammar fact is advertised
in five files, and a red gate between any pair is not a landable intermediate
state. The two sealed-field lists in `decisions/parallel-session-landings.md`
drop the field under the reviewer's 2026-08-04 approval of the ADR edit.

Genuinely unverifiable-at-authoring-time claims — which the evidence build
produced zero of across eight tickets — belong in a ticket's What-to-build
prose, and the field bullet's replacement says so in one sentence. Checkable
preconditions stay with `Blocked by:`.

Fence spans three directories: one grammar fact, five advertisements, one
landable change.

## Acceptance

- [ ] [RS1] the template, the field bullet, and the taught example carry no `Assumptions:` line, and the example-agreement check is green.
- [ ] [RS2] the fixture-bite and docs-workflow template anchors expect the new grammar, in the same change as the template edit.
- [ ] [RS3] both sealed-field lists in `decisions/parallel-session-landings.md` drop the field, and no other list in that ADR still names it.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RS1 | drop the example's line but keep the template's | the example-agreement check | revert the template block alone, run `bench gate`, expect the conformance phase to red on the template/example disagreement |
| RS2 | leave one anchor expecting the old template | the fixture-bite check | revert the fixture-bite anchor, run `bench gate`, expect the conformance phase to red naming the stale anchor |
| RS3 | leave the field in one sealed-field list | the reviewer, at review | revert one list, re-read the ADR against the retired grammar, expect the stale list to re-teach the field to the next parallel session |
