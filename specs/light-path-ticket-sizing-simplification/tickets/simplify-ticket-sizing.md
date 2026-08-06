# Simplify ordinary ticket sizing

Blocked by: none
Ownership fence: `.agents/skills/bench-craft-tickets/SKILL.md`, `internal/anchors/registry_data.go`, `specs/light-path-ticket-sizing-simplification/tickets/simplify-ticket-sizing.md`
Integration surfaces: ordinary ticket sizing guidance→`.agents/skills/bench-craft-tickets/SKILL.md`; Claude skill consumer→`.claude/skills/bench-craft-tickets/SKILL.md` + TS3; retired exact-prose anchor→`internal/anchors/registry_data.go`; skills-index advertisement→`.bench/BENCH-reference.md` + TS3
Contracts: the canonical ticket-sizing guidance crosses `.agents/skills/bench-craft-tickets/SKILL.md`→`.claude/skills/bench-craft-tickets/SKILL.md` through the tracked symlink, asserted by TS3; absence means the two harnesses load different guidance

## What to build

Ordinary ticket slicing has one sizing rule: keep splitting while every resulting unit can land independently green. A group stays whole only when the ticket names the specific red a thinner cut would strand so review can reproduce that claim by attempting the split. Remove the numeric proxy rules and their escape hatches without changing the wide-refactor sequence or importing tier economics from `craft-line`.

## Acceptance

- [ ] [TS1] `craft-tickets` says to split until another split would leave no independently-green landing, and keeping a group whole names the specific stranded red for review to re-derive by attempting the split rather than describing feature wholeness.
- [ ] [TS2] the fence-width and acceptance-row-count proxy paragraphs, their justification escape hatches, their two signal-table rows, and any resulting in-skill example or prose residues are absent.
- [ ] [TS3] expand–migrate–contract remains unchanged, no tier or cost policy enters `craft-tickets`, the obsolete exact-prose anchor is removed, the Claude symlink resolves to the edited source, and the skills-index row remains unchanged.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| TS1 | permit a whole group with only a description of the feature as one outcome | the semantic reviewer | attempt the thinner independently-green split, expect the grouping to name the exact red that split strands |
| TS2 | restore either numeric proxy or its keep-whole escape hatch | the consistency reviewer | sweep the skill for the retired sizing vocabulary and compare every hit with the closed sizing decision |
| TS3 | retain the exact-prose anchor or alter the wide-refactor, tier, mirror, or skills-index surfaces | the conformance gate | run the path-scoped `bench commit`, expect anchor, mirror, index, and example-agreement conformance to grade the proposed tree |
