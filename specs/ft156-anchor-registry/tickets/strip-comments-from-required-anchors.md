# Strip comments from required anchors

Blocked by: expose-anchor-query.md, migrate-anchor-table.md
Ownership fence: `internal/anchors`, `tests/canary/workflow-guidance-anchors/commented-required-anchor`
Contracts: comment-stripped require matching crosses `internal/anchors`→the mutated graded-root fixture, asserted by HC1 against the real missing-anchor diagnostic; the current registry crosses the same matcher, asserted by HC2 through the kit root

## What to build

Make both whole-file and section-scoped require kinds ignore HTML comments in the shared matcher. Use a whole-file required anchor for the canary because checkout-current section-scoped matching already strips comments.

## Acceptance

- [ ] [HC1] Wrapping a whole-file required sentence in an HTML comment fails the graded root with its precise missing-anchor diagnostic.
- [ ] [HC2] The kit root remains green after stripping comments, proving no live required anchor was satisfied only inside a comment.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| HC1 | bypass comment stripping for whole-file require matching | the commented-required-anchor canary | run the owning fixture and expect its targeted diagnostic to disappear, turning the canary red |
| HC2 | add a registry needle found only inside an HTML comment in the kit tree | the real kit conformance run | run the focused root conformance check and expect the missing-anchor diagnostic |
