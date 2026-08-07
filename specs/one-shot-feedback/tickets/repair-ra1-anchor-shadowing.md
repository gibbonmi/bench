# Repair RA1 anchor shadowing, registry grouping, and the vocabulary claim

Blocked by: none
Ownership fence: `.agents/commands/bench-final-check.md`, `internal/anchors/registry_data.go`
Integration surfaces: none
Contracts: none crosses
Closure: AB1/heading-single-occurrence, AB2/group-block-placement, AB3/claim-precision

## What to build

The composed review of candidate `9a9c996d` accepted three findings (reviewer
disposition, 2026-08-07); this ticket repairs all three.

1. **RA1 anchor shadowing (Coverage, concrete).** `## Repair attribution`
   occurs twice in `.agents/commands/bench-final-check.md`: the template
   heading and a prose restatement ("Under `## Repair attribution`, write one
   table row per ticket"). Deleting the template heading leaves the needle
   satisfied by the prose copy, so the ticketed mutation cannot red. Reword the
   prose so it does not restate the literal heading (e.g. "Under that heading,
   write one table row per ticket"), leaving exactly one occurrence — then
   prove the recovered red: delete the template heading line, run the focused
   conformance command, observe the missing-anchor diagnostic, restore.
2. **Registry group placement (Standards/Coverage, legibility).** The new
   ungrouped rows appended at the tail of `internal/anchors/registry_data.go`
   sit after the `AfterSpecAuthorization` block while belonging to
   `BeforeStructured` (zero value). Relocate every ungrouped tail row into the
   file's `BeforeStructured` block (the region ending near the
   `projects/benchkit.md` rows), preserving their relative order and bytes.
   Data-only move; `EvaluateGroup` filters by group, so behavior is unchanged.
3. **Vocabulary-claim precision (Standards, advisory accepted).** The template
   sentence claiming it is the only place the vocabulary is written down is
   absolute while the anchors Require needle necessarily restates the six
   terms as its enforcement copy. Make the claim precise: the template is the
   single guidance source; the anchors needle is the enforcement copy that
   pins it. Keep the RA1 vocabulary needle itself byte-intact.

Semantic exception (explicit, per craft-tickets): AB2 and AB3 are
legibility/precision facts with no mechanical red — no anchors kind observes
row grouping, and the claim wording is review-graded. Their mutation rows name
the honest owner (review plus the unchanged-green checks); AB1's red is
mechanical and demonstrated.

## Acceptance

- [ ] [AB1] (covers local) `## Repair attribution` occurs exactly once in the final-check command, and deleting the template heading line reds conformance with RA1's missing-anchor diagnostic.
- [ ] [AB2] (covers local) no ungrouped registry row sits after a grouped block: the `BeforeStructured` rows are contiguous in their block, all needles and diagnostics byte-identical, conformance green.
- [ ] [AB3] (covers local) the template's ownership sentence distinguishes the guidance source from the anchors enforcement copy, with the vocabulary needle unchanged.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AB1/heading-single-occurrence | delete the template's `## Repair attribution` heading line (Require row retained) | conformance workflow-anchors check | `BENCH_CONFORMANCE_ROOT=$PWD go test -count=1 ./internal/conformance -run '^TestRootConformance$'`, expect the RA1 heading missing-anchor diagnostic — the red the shadowing suppressed |
| AB2/group-block-placement | none exists — grouping is gate-invisible by design (Out of scope kept the ordering anchors kind cut) | review, plus the unchanged-green conformance and `go test ./internal/anchors` runs proving the move is behavior-neutral | run both commands before and after the move; expect green both times and byte-identical needles/diagnostics via `git diff` inspection |
| AB3/claim-precision | none exists — claim wording is review-graded prose | review, plus conformance green proving the vocabulary needle survived the rewording | run the focused conformance command; expect green with the vocabulary needle untouched |
