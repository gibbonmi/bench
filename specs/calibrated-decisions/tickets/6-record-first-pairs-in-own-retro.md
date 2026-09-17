# Record the first pairs in this build's own retro

Blocked by: 5-define-calibration-measure-in-scorecard.md
Writes: reviews/calibrated-decisions.md (new), CHANGELOG.md, specs/calibrated-decisions/spec.md, tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns, tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary
Covers: CR27, CR33

Chunk: CD5.

## What to build

The build that lands this spec records the first confidence-and-label pairs on all three surfaces.
Collect the pairs the earlier tickets recorded in the review pickup. The pairs are each done-claim row with its confidence and probe label, and each review finding with its confidence and disposition label. They also include this build's line declaration with its expected rounds and actual rounds.
Group them by surface and by model, effort, and role, so the retro author can fill the calibration table from the pickup alone.
State the pair count per role and the abstention count in the pickup.

Write one concise typed changelog entry under `[Unreleased]` for the calibration guidance, the scaffold table, and the scorecard measure.
Claim no measured routing change.
Record the completion evidence in the spec's completion record.

Scenario: the review pickup lists at least one labeled row for each of the three surfaces.
The final-check retro copies those rows into the calibration table and states `unknown` or the hand-computed Brier mean below it.

This is the last ticket, so it carries the integrated invariant: every Markdown file the build edited passes the prose lane.

Preserve the two co-named changelog fixture pins without changing their planted diagnostics.

## Acceptance

- [ ] The review pickup holds at least one labeled row for each of the three surfaces, with its confidence and label.
- [ ] The review pickup states the pair count per role and the abstention count.
- [ ] The changelog entry describes the guidance, the scaffold table, and the measure without a measured routing claim.
- [ ] Every Markdown file the build edited passes `bench gate-prose`.

Run the prose lane over every Markdown file the build edited, then the whole-tree gate before the landing.
The retro duty itself runs at final-check after the landing, and the pickup is its only input.
