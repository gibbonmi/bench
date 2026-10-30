# The implement phase states the plan probe, chunk base, cost report, sibling-chunk, witness red, and docs divergence rules

Blocked by: none
Writes: .agents/commands/bench-implement-spec.md, .agents/commands/bench-review-implementation.md, .agents/skills/bench-craft-delegate/references/delegation-discipline.md, CHANGELOG.md, ROADMAP.md, roadmap/FT456.md, roadmap/FT458.md, roadmap/FT459.md, roadmap/FT470.md, roadmap/FT471.md
Covers: none

## What to build

The roadmap rows FT456, FT458, FT459, FT470, and FT471 are decided kit edits.
Each row adds one rule to the implement phase guidance, and this ticket retires the five rows.

Add each rule beside the step that it governs, inside the existing paragraphs of `.agents/commands/bench-implement-spec.md`.
That file has no line headroom under its `guidance-prose-budgets` row, so the edit adds no line to it.
A rule that another file already owns goes to that owner:

- Step 1 of `/bench-review-implementation` owns the chunk base, so the FT471 charge rule points to it.
- `/bench-review-implementation` owns the confirming round, so the FT456 confirming round rule goes there.
- `craft-delegate` owns the charge contents, so the FT456 helper owner rule goes to its delegation discipline.

Delete the five row files and their heading lines in `ROADMAP.md`.
The pinned row FT420 names FT456, and this ticket does not edit FT420.

## Acceptance

- [ ] FT470: the plan amendment lists each mutation that a ticket names as a requirement `probe` on its registered check, not only as a helper bite.
- [ ] FT471: the first-ticket charge of each chunk names `bench record ticket --base <chunk base>`, and the guidance points to the chunk base owner.
- [ ] FT471: at each chunk checkpoint, the orchestrator reports the wall time, the bookkeeping commits against the change commits, and each repeated step, and names the cut in the same update.
- [ ] FT456: until the review record holds sibling chunks, a delegated run freezes its concurrent chunks in plan order.
- [ ] FT456: until the review record holds sibling chunks, the coordinator plans a three-axis confirming round for a repair delta that holds a sibling chunk's commit.
- [ ] FT456: when two concurrent tickets add a new helper to one Go package, the charge names the owner of that helper.
- [ ] FT458: before the initial review charge, the record holds each author witness red, in the same plan commit as the next assignment.
- [ ] FT459: the orchestrator charges the docs ticket with each settled spec-versus-code divergence and corrects the spec prose before the docs review.
- [ ] The five row files and their `ROADMAP.md` heading lines are gone.
- [ ] The `bench commit` lane passes.
