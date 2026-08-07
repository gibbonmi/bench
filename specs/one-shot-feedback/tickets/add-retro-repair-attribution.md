# Add retro repair attribution and the drain tally

Blocked by: add-breakdown-review-moment.md
Ownership fence: `.agents/commands/bench-final-check.md`, `.agents/commands/bench-what-next.md`, `internal/retros/`, `internal/anchors/registry_data.go`
Integration surfaces: anchors Require rows→`internal/anchors/registry_data.go`; retros regression fixture→existing `internal/retros` Recommendations parser path exercised by AB2; registry append tail→add-falsification-predicate-questions.md
Contracts: the cause vocabulary crosses `.agents/commands/bench-final-check.md` (template, sole owner)→the what-next drain by named reference only, asserted by AB1's vocabulary needle against the real template and AB3's reading-only-from-tables needle — no literal value is restated across the fence
Closure: AB1/attribution-heading, AB1/vocabulary-line, AB2/parser-noninterference, AB3/tally-needle

## What to build

Every promoted build's retro records repair attribution, and the
`/bench-what-next` drain reports the running one-shot tally (spec story 2).
Grouped as one ticket because template, drain duty, and parser fixture are one
vertical value path; no thinner cut strands a red — the grouping is the spec's
per-story slice.

Needle-first per row: land each anchors `Require` row, observe the missing-anchor
conformance red, then the prose.

- RA1 (`.agents/commands/bench-final-check.md`): the retro template gains a
  `## Repair attribution` heading placed between `## Coordinator catches` and
  `## Agent-experience improvements`, with one table row per ticket: repair
  rounds landed, and per round one cause from the vocabulary line the template
  solely owns — needle text exactly `` `shaping-ambiguity`, `spec-row`,
  `ticket-slicing`, `tree-drift`, `delegate-error`, `other` `` (one line; a
  zero-round ticket records `none`). Registry rows join the existing
  `AfterImplementSpec` block for this file (with explicit
  `Group: AfterImplementSpec`), mirroring the retro-heading rows already there.
  Placement before the improvements heading is gate-invisible (no ordering
  anchors kind) — it is a reviewer-graded exception backed by AB2's fixture.
- RA2 (`internal/retros`): a regression fixture proving the parser ignores the
  new heading — a body carrying `## Repair attribution` content before
  `## Agent-experience improvements` returns unchanged recommendation units.
  Mirror `TestRecommendationsKeepsImprovementParagraphsAndListItemsSeparate` in
  `internal/retros/recommendations_test.go` (inline byte-string fixture,
  `reflect.DeepEqual`, single-case). Honestly classified: it cannot start red —
  the parser keys on the exact improvements heading — and lands as a regression
  guard, per the spec's RA2 row.
- RA3 (`.agents/commands/bench-what-next.md`, step 4): the drain gains one duty —
  when drained retros carry attribution tables, the exit reports the tally.
  Needle: `report tickets total, one-shots, and per-cause counts, reading causes
  only from the drained tables`. Name the vocabulary only as "the cause
  vocabulary the retro template owns" — never restate the values. Reporting
  only: no roadmap grammar or CLI change.

Semantic exception (explicit): the attribution table's row shape and the
placement ordering are review-graded prose; the gate proves the quoted needles
and the parser fixture proves non-interference.

## Acceptance

- [ ] [AB1] (covers RA1) the retro template carries `## Repair attribution` before the improvements heading and the one vocabulary line it solely owns, with both Require rows in the registry's AfterImplementSpec block.
- [ ] [AB2] (covers RA2) `retros.Recommendations` returns unchanged units for a body carrying `## Repair attribution` before the improvements heading, pinned by a regression fixture in `internal/retros`.
- [ ] [AB3] (covers RA3) `/bench-what-next` step 4 carries the tally-reporting needle, with its Require row in the registry.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AB1/attribution-heading | delete the `## Repair attribution` heading from the template (Require row retained) | conformance workflow-anchors check | `BENCH_CONFORMANCE_ROOT=$PWD go test -count=1 ./internal/conformance -run '^TestRootConformance$'`, expect the heading's missing-anchor diagnostic |
| AB1/vocabulary-line | delete the vocabulary line from the template (Require row retained) | conformance workflow-anchors check | same command, expect the vocabulary needle's missing-anchor diagnostic |
| AB2/parser-noninterference | mutate `retros.Recommendations` to enter improvements mode on any `## ` heading instead of the exact one | the new regression fixture in `internal/retros` | `go test ./internal/retros -run TestRecommendations`, expect the fixture's unchanged-units failure |
| AB3/tally-needle | delete the tally clause from what-next step 4 (Require row retained) | conformance workflow-anchors check | same command as AB1, expect the what-next missing-anchor diagnostic |
