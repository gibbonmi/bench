# Define and split the spec-row repair cause

Blocked by: none
Writes: .agents/commands/bench-final-check.md, CHANGELOG.md
Covers: none

## What to build

The repair-attribution table of each implementation retro gives one cause for each repair round.
`/bench-final-check` owns the cause vocabulary, and it defines only `one-source` and `check-gap`.
`spec-row` is the largest cause since August, with about 178 rounds, and it has no definition.

The 2026-10-19 assessment in `docs/research/2026-10-19-spec-row-repair-rounds.md` classified the `spec-row` rounds.
About 40 percent were edges that a build review found and no row decided.
About 30 percent were clear rows that the code or its test missed.
The remaining rounds were defects in the spec text.
Each kind needs a different fix, so one undefined label cannot steer a fix.

Add `edge-found` to the vocabulary for a review-found edge that no row decided.
Define `spec-row` as a defect in the spec text.
Define `delegate-error` as a clear row that the code or its test missed, including a test that stays silent under a probe.
Remove the sentence that names an anchors-registry needle as the enforcement copy of the vocabulary, because anchor-admission-prune removed that needle on 2026-10-17.

The reviewer approved proposal S1 of the assessment on 2026-10-19.
The drain tally reads the terms from the drained tables, so it needs no change, and an older retro keeps its terms.

## Acceptance

- [ ] The repair-attribution paragraph of `/bench-final-check` lists `edge-found` and defines `spec-row`, `edge-found`, and `delegate-error`.
- [ ] No sentence in the paragraph names an anchors-registry needle for the vocabulary.
- [ ] `CHANGELOG.md` records the vocabulary change.
- [ ] The landing gate is green.
