# Point craft-spec at the ticket contracts and name the process boundary

Blocked by: charge-the-delegate-side-duties.md
Ownership fence: `.agents/skills/bench-craft-spec/SKILL.md`, `projects/benchkit.md`, `internal/conformance/docs_workflow_helpers_test.go`, `internal/conformance/fixture_bite_test.go`
Assumptions: the three owning H2 titles are `Slicing a build for delegates` and `The edge inventory` in `craft-spec` and `Hostile-input checklist (shell CLI)` in `projects/benchkit.md`; the canonical edge classes are an inline em-dash-delimited run inside one sentence of `The edge inventory` whose form is preserved when the class is added; the anchored `craft-tickets` cross-pointers in `craft-spec` are added beside and never edited. Re-derive from the tree at pickup.

## What to build

FT164 story 6: the spec-side halves of the ticket contract, so cross-seam
mismatches and serialize-then-reload defects are visible before a build starts.

`craft-spec`'s "Slicing a build for delegates" gains one sentence naming that
each fence carries value contracts, pointing at the `craft-tickets` rule by name
rather than restating it — the same shape as the existing cross-pointer pair,
which stays untouched, because a restated procedure is exactly the drift the
code standard calls a defect.

The canonical edge walk gains a process-boundary lifecycle class: defects
visible only after state is serialized and a fresh process reloads it, plus
recomposition suites that stop at first success. It joins the existing inline
em-dash-delimited class run, keeping that run's form. The profile's
hostile-input checklist gains the concrete entry naming the same two shapes, so
a benchkit spec walking the checklist hits it with a project-specific handle.
Unit-level success has hidden defects that only appeared once a fresh process
reloaded serialized state; the class makes them nameable at spec time.

Enforcement is anchor-first: three needles scoped to their owning sections — the
pointer sentence, the process-boundary class token inside the edge-walk
sentence, and the checklist entry's lead phrase — each with one mutation-table
row proving its diagnostic fires. Wording beyond the pinned phrases is
review-graded and stated as such.

## Acceptance

- [ ] [SP1] the slicing section names value contracts at each fence by pointer to the `craft-tickets` rule, with the existing cross-pointers unchanged.
- [ ] [SP2] the canonical edge-walk sentence carries the process-boundary lifecycle class within its existing inline class run.
- [ ] [SP3] the profile's hostile-input checklist carries the process-boundary entry naming serialize-then-reload defects and first-success recomposition suites.
- [ ] [SP4] every needle this ticket registers has a byte-exact entry in `TestSpecBuildCadenceAnchorsRejectDeletionSwapAndRawGitRouting` whose named diagnostic fires.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SP1 | replace the pointer sentence with a restatement of the contracts procedure | the `craft-spec contract pointer` mutation subtest | swap the sentence, run the anchor check, expect the pointer diagnostic |
| SP2 | delete the process-boundary class token from the edge-walk run | the `edge walk process boundary` mutation subtest | delete the token, run the anchor check, expect the edge-class diagnostic |
| SP3 | delete the checklist entry's lead phrase | the `profile process boundary entry` mutation subtest | delete the entry, run the anchor check, expect the checklist diagnostic |
| SP4 | register a needle and land no mutation entry for it | `review` plus the mutation harness | delete one entry, run the harness and watch the remaining rows still pass, then read this ticket's needle list against the table |
