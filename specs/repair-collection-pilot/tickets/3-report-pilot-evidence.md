# Report audited pilot evidence

Blocked by: 2-collect-repair-evidence.md
Writes: internal/repairpilot (new), docs/repair-collection-pilot.md (new), specs/repair-collection-pilot, reviews/repair-collection-pilot.md (new)
Covers: RP31, RP38, RP39, RP40, RP41, RP42, RP43, RP44, RP47, RP49, RP50, RP55, RP56, RP58, RP59, RP61

## What to build

Deliver the default and full report from RP-C2's stored observations and audits.
Keep the required example classes distinct from progress conclusions and sequence counts.
A sequence can supply more than one interval class.
A missing required class makes the terminal result inconclusive.

The full report retains every observation, unknown label, audit, and evidence gap.
The full report lists incomplete sequences, and the default report gives their count.
Reports before cutoff identify the sample as provisional.

Write the repository-only operating guide with the exact activation, record, and report commands.
Require a separate native-evidence audit before the operator counts a proposed progress label.
Show a synthetic example as a fixture, never as real pilot evidence.
Name `capture/reports/repair-collection-pilot.md` as the eventual real report destination.
Do not create an empirical report or start the pilot as an implementation step.

Carry the whole-package invariant: collection and reports do not change required checks, gate outcomes, warnings, or model routing.
The guide routes the eventual report to the reviewer for a later detector decision.
Freeze RP-C3 for review, then reconcile the complete acceptance map and final integration source.

## Acceptance

- [ ] The report distinguishes stalled repairs, productive repairs, unchanged reruns, and overlapping assignments.
- [ ] Any missing required class produces an inconclusive terminal report.
- [ ] The full report exposes every retained observation and its evidence status.
- [ ] Repeated reports over the same document and clock produce the same output without mutation.
- [ ] The guide requires a separate evidence audit and clearly marks synthetic examples.
- [ ] Implementation does not activate the real pilot or introduce a detector effect.
