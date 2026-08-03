# Teach the executable-contract template and its examples

Blocked by: none
Ownership fence: `.agents/skills/bench-craft-tickets/SKILL.md`, `internal/conformance/docs_workflow_helpers_test.go`, `internal/conformance/fixture_bite_test.go`
Assumptions: `markdownH2Sections` is fence-blind — any column-0 `## ` inside a fenced block truncates the enclosing section body — so the five template needles scope to `Write one file per ticket` whose body then runs to end of file; the whole-file needles for `## What to build` and `## Acceptance` and `Blocked by:` stay registered because collapsed substring matching is satisfied by the `###` forms and only the placeholder needle at `docs_workflow_helpers_test.go` lines 232–233 retires; the mutation harness `TestSpecBuildCadenceAnchorsRejectDeletionSwapAndRawGitRouting` asserts an anchor count of exactly 1 against live file bytes; the `tests/canary/workflow-guidance-anchors` fixtures must stay green. Re-derive from the tree at pickup.

## What to build

FT164 stories 1 and 2: a ticket written cold from the `craft-tickets` template
assigns without hand-normalization, because the template teaches the shape the
specbuild parser actually accepts.

The template block teaches single-line `- [ ] [ID] <behavior>` rows with explicit
ticket-local ids — a short uppercase tag plus a number, unique in the ticket —
noting beside them that only `R`-prefixed ids range-expand; a one-line backticked
`Ownership fence:` enumerating every path the ticket will write; a one-line
`Assumptions:` whose clauses separate with semicolons because the parser splits
on commas and drops a wrapped continuation; a `Blocked by:` keyed on sibling
ticket file basenames; and a red-mutations table binding each row id to one
mutation, its independent owner, and the operation sequence proving the red.

Every heading inside the fenced template and inside both examples is `###`, not
`##`, because a column-0 `## ` inside a fence truncates the resolved body of the
section that owns it and would leave the new needles unscopeable.

The Good example demonstrates every one of those fields — two rows with distinct
ids, a multi-entry fence, real assumptions, a blocked-by line, a mutations row
per id. Its `<!-- ticket-example:begin -->` and `<!-- ticket-example:end -->`
markers sit at column 0 immediately outside the fenced block — begin above the
opening fence line and end below the closing one — so both fence lines fall
inside the marked region for the extractor to require and strip. The Bad example
becomes the failure the corpus actually commits: an oversized ticket of spec
fragments, taxonomy essays, and reviewer narrative, kept compact with elision
markers. The gate-checkbox prohibition in the cadence paragraph survives verbatim
and gains a needle of its own.

Enforcement lands anchor-first here: the five template facts and the prohibition
register section-scoped through the bulk helper pattern, and each new needle gets
a byte-exact mutation row proving its diagnostic fires.

## Acceptance

- [ ] [TT1] the template block teaches the labeled single-line acceptance row and the `R`-only range caveat beside it.
- [ ] [TT2] the template block teaches the one-line backticked `Ownership fence:` enumerating every path the ticket writes.
- [ ] [TT3] the template block teaches the one-line `Assumptions:` field with semicolon-separated clauses.
- [ ] [TT4] the template block teaches `Blocked by:` keyed on sibling ticket file basenames.
- [ ] [TT5] the template block carries the red-mutations table header naming criterion, mutation, owner, and operation sequence.
- [ ] [TT6] the Good example demonstrates every taught field with two distinct-id rows and one mutations row per id, its markers at column 0 outside the fence lines.
- [ ] [TT7] the Bad example is an oversized-but-credible ticket rather than a layer list.
- [ ] [TT8] the gate-checkbox prohibition sentence survives verbatim in the cadence paragraph and is pinned by a section-scoped needle.
- [ ] [TT9] the placeholder needle retires while the What-to-build and Acceptance and Blocked-by whole-file needles stay registered and green.
- [ ] [TT10] every heading inside the fenced template and both examples is `###`, so the owning section body resolves past them to end of file.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| TT1 | restore the unlabeled `- [ ] <Observable behavioral criterion>` placeholder | the `template row` mutation subtest | swap the labeled placeholder for the old one, run the anchor check, expect the labeled-row diagnostic |
| TT2 | delete the fence line from the template block | the `template fence` mutation subtest | delete the line, run the anchor check, expect the fence diagnostic |
| TT3 | delete the assumptions line from the template block | the `template assumptions` mutation subtest | delete the line, run the anchor check, expect the assumptions diagnostic |
| TT4 | swap the basename blocked-by line back to sibling titles | the `template blocked by` mutation subtest | replace the line with the title form, run the anchor check, expect the basename diagnostic |
| TT5 | delete the red-mutations table header row | the `template mutations header` mutation subtest | delete the header, run the anchor check, expect the mutations-table diagnostic |
| TT6 | shrink the Good example to one acceptance row | `review` plus the focused conformance run | `BENCH_CONFORMANCE_ROOT=$PWD BENCH_CONFORMANCE_CHECK=docs-currency-workflow go test ./internal/conformance -count=1 -run '^TestRootConformance$'`, then read the example |
| TT7 | restore the layer-list Bad example | `review` | read the Bad example against the corpus failure the spec names |
| TT8 | rewrite the prohibition sentence to permit a gate checkbox | the `gate checkbox prohibition` mutation subtest | swap the sentence, run the anchor check, expect the prohibition diagnostic |
| TT9 | leave the retired placeholder needle registered against the rewritten template | the focused conformance run | run the conformance check over the tree and expect green only once the stale needle is gone and the kept needles still pass |
| TT10 | promote one heading inside the fenced template back to `## ` | the `template heading depth` mutation subtest | promote the heading in a temp-tree copy, run the anchor check, expect the truncated-body diagnostic from a template needle |
