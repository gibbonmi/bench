# Calibrated decisions: a stated confidence on every labeled claim

Status: ready

## Destination

Bench records a stated confidence with each claim on a decision surface that
Bench already labels later. The candidate surfaces are a delegate done-claim
row, a review finding, and the expected repair rounds in a line declaration.
A bounded proper scoring rule scores each confidence-and-label pair. The retro
records the pairs. The scorecard aggregates the score per model, effort, and
role. The routing loop then rewards an honest claim over a confident guess,
and it counts an honest abstention.

The map decides eight things. They are the surfaces, the scale and rule, the
abstention treatment, and the label sources. They are also the claim schema,
where the score lands, how the score steers routing, and whether the first
spec adds a verb.

## Notes

Domain: process calibration for agent claims, transferred from RLCR,
behavioral calibration, and calibration-aware RL. Bench trains no weights; the
reward steers the routing tiers and the guidance prose.

Skills a session consults: `craft-line`, `craft-delegate`, `craft-review`,
`craft-research`, and `craft-synthesis` for the kit fold.

Standing preferences: a recommendation names the complete shape and never
scopes by implementation time. The top tier implements nothing unless the
reviewer names it. The gate and the coordinator are the only label sources.

A map-owned asset stays in the map's assets folder, decisions/calibrated-decisions/assets/.

## Decisions so far

- [What do Jev and the calibration papers establish?](calibrated-decisions/tickets/1.md): Jev names a method and publishes none. The published tactic is a bounded proper scoring rule on a stated confidence, with an abstention scored zero.
- [Which decision surfaces carry a stated confidence?](calibrated-decisions/tickets/2.md): A done-claim row, a review finding, and a line declaration's expected repair rounds each carry one confidence.
- [Which confidence scale and scoring rule apply?](calibrated-decisions/tickets/3.md): An integer from 0 to 10, normalized to `p`, scored by the Brier rule `(p - label)^2`.
- [Which sources may label an outcome?](calibrated-decisions/tickets/4.md): The gate, the coordinator's probe, and the reviewer's disposition; never a model judgment.
- [How is an abstention scored and counted?](calibrated-decisions/tickets/5.md): An abstention scores 0, stays out of the Brier mean, and is counted as its own number.
- [What fixed schema does a claim carry?](calibrated-decisions/tickets/6.md): `status` in verified, claimed, abstained; `confidence` 0 to 10; `label` held or refuted; no free text.
- [Where do the pair and the score land?](calibrated-decisions/tickets/7.md): Each retro holds a calibration table of pairs. The scorecard holds the Brier mean, the pair count, and the abstention count.
- [Does confidence change whether a review finding blocks?](calibrated-decisions/tickets/9.md): No. Kind and citation decide; the confidence is recorded and scored later.
- [How does the calibration score steer routing?](calibrated-decisions/tickets/8.md): One input to the two-run judgment rule; a threshold decision waits for ten pairs.
- [Is the first spec guidance-only, or guidance plus a scoring verb?](calibrated-decisions/tickets/10.md): Guidance-only; the Brier mean is computed by hand until ten pairs exist.
- [Which landing dogfoods the pairs first?](calibrated-decisions/tickets/11.md): The build that lands this guidance records the first calibration table in its own retro.

## Not yet specified

## Spec-writer discretion

- The exact column order of the retro calibration table, provided every
  field in ticket 6 is present.
- The exact field names in the delegate return and the review record,
  provided the enum values and the integer range in ticket 6 are exact.
- Whether a new rule lands in a skill's `SKILL.md` or in its
  `references/` file, provided the skill's prose budget holds.

## Out of scope

- Training or fine-tuning any model. Bench steers routing and prose only.
- A Bench verb that computes the calibration score. It waits for ten
  recorded pairs.
- A numeric threshold that moves a tier on its own. The reviewer decides it
  after ten pairs.
- A reference-model label as the outcome source. The gate is the oracle.
- The Jev type-safety guarantee as a Bench claim. Bench emits prose and
  cannot guarantee a schema by construction.

## Sources

- Path: `decisions/calibrated-decisions/assets/calibration-research.md`
  Supports: #1 in full, and the recommendations in #2 through #10. Produced 2026-09-16 from four arXiv papers and the TypeSafe article.
  Drift: a new primary source on process-level calibration, or a change to the scorecard contract.
- URL: https://arxiv.org/abs/2507.16806
  Supports: #3 the bounded proper scoring rule and the Brier reward; #8 the out-of-domain caveat.
  Drift: a later version that changes Theorem 1 or the reward.
- URL: https://arxiv.org/abs/2512.19920
  Supports: #5 the abstention reward and the failure of an explicit risk threshold.
  Drift: a later version that changes the reward or the failure report.
- URL: https://arxiv.org/abs/2601.13284
  Supports: #4 why a binary reward alone yields overconfidence.
  Drift: a later version that changes the diagnosis.
- URL: https://arxiv.org/abs/2503.02623
  Supports: #3 the integer 0-to-10 confidence scale.
  Drift: a later version that changes the scale.
- URL: https://typesafe.ai/blog/introducing-system-one-models-and-jev
  Supports: #1 the claim table and the rejected evaluation method.
  Drift: a published RLCD method or a third-party calibration measurement of Jev.
- Path: `capture/agent-performance/README.md`
  Supports: #7 the scorecard measures the new score joins.
  Drift: a change to the Measures table or the update contract.
- Path: `.bench/BENCH.md`
  Supports: #4 invariant 1, the gate as the only oracle.
  Drift: a change to the four invariants.
