# Anchor the realigned consumers

Blocked by: realign-the-consumers-glossary-and-docs.md
Writes: internal/anchors/registry_data.go, tests/canary/workflow-guidance-anchors, .agents/skills/bench-craft-spec/SKILL.md

## What to build

Repair ticket from the build's own probes. The realign ticket landed prose the
gate does not grade: deleting `**coverage row**` from `CONTEXT.md` stays green,
and so does swapping the `rows[N]{story,behavior,seam}` projection sentence in
`craft-review` or `bench-implement-spec.md`, or `craft-delegate`'s behavior/seam
charge — SR27 and SR28 promised enforcement that did not exist. Anchor each of
those needles with a `Require` row and a canary that plants its own deletion,
mirroring the `**decision ticket**` row and its `context-ticket-vocabulary`
fixture. Separately, `craft-spec`'s Review rubric still asks whether a "red
signal" names an outcome family — a column the map no longer carries; reword the
sentence, its registry row, and its canary together so the anchor guards a true
claim.

## Acceptance

- [ ] `(covers SR28)` `CONTEXT.md`'s `**coverage row**` and `**acceptance row**` entries are anchored; deleting either turns the gate red through a canary.
- [ ] `(covers SR27)` `craft-review`'s and `bench-implement-spec.md`'s `rows[N]{story,behavior,seam}` sentences and `craft-delegate`'s behavior/seam/why charge are anchored, each with a canary planting its deletion.
- [ ] `(covers SR21)` The rubric sentence names behavior, why-it-catches, or decision answer — not `red signal` — with needle, row, and canary reworded together; the old needle appears nowhere.
- [ ] `(covers SR26)` `bench canary` and `bench gate` are green.
