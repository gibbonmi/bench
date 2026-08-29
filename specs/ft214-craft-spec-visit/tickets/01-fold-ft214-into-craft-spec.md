# Fold the FT214 clauses into craft-spec

Blocked by: none
Writes: .agents/skills/bench-craft-spec/SKILL.md, .agents/skills/bench-craft-spec/references/map-discipline.md (new), .agents/skills/bench-craft-tickets/SKILL.md, .agents/skills/bench-craft-domain/SKILL.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, ROADMAP.md, roadmap/FT214.md (deleted)

## What to build

A spec author who loads `craft-spec` finds every rule that `roadmap/FT214.md`
records, placed where the author acts on it. The build-fence rule, the
Writes-derived fences, the explore-pass reads, and the per-row rubric question
stay in `SKILL.md`. The map-discipline clauses move behind one context pointer
in a new `references/map-discipline.md`. The ticket-sizing sentence lands in
`craft-tickets`, and the cardinality clause lands in `craft-domain`. Each file
stays inside its `projects/benchkit.md` budget, and every anchored sentence
keeps its bytes and line breaks. The landing retires the FT214 row: the index
line, the sequence entry, and the detail file leave together.

## Acceptance

- [ ] `.agents/skills/bench-craft-spec/SKILL.md` states that a build may not edit its own spec's acceptance rows, budget targets, or ownership fences. A spec-level shortfall returns to `/bench-write-spec`.
- [ ] `.agents/skills/bench-craft-spec/SKILL.md` derives the ownership fences from the tickets' own `Writes:` lines.
- [ ] `.agents/skills/bench-craft-spec/SKILL.md` names three explore-pass reads. These are every named enforcement file, one precedent per seam, and one tree-wide sweep for readers of each changed value.
- [ ] The `Review rubric` in `.agents/skills/bench-craft-spec/SKILL.md` asks, per row, for the gate check or test that reds the row, or a review-owned mark.
- [ ] `.agents/skills/bench-craft-spec/SKILL.md` points at `references/map-discipline.md` from the acceptance-coverage-map section.
- [ ] `references/map-discipline.md` carries every remaining FT214 clause as one rule each, in ASD-STE100 prose, with no rule stated twice across the two files.
- [ ] `.agents/skills/bench-craft-tickets/SKILL.md` states that a rewrite ticket is sized by the lines the delegate must read.
- [ ] `.agents/skills/bench-craft-domain/SKILL.md` states that a new ordering or grouping requirement fixes the cardinality at the concept edge with a many-case fixture.
- [ ] `internal/anchors/registry_data.go` pins the build-fence sentence, the per-row rubric question, and the map-discipline pointer, and `internal/anchors/registry_data_test.go` shows each pin red on removal.
- [ ] `TestGuidanceProseBudgetsHoldOnTheLiveTree` stays green with no change to the budget table.
- [ ] `TestProseMechanicsHoldsOnTheLiveTree` and the fixture-bite tests stay green.
- [ ] `ROADMAP.md` carries no FT214 line and `roadmap/FT214.md` is absent.
