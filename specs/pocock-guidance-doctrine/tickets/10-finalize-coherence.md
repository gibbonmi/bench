# Verify the composed doctrine and update the record surfaces

Blocked by: 09-prose-budget-conformance.md
Ownership fence: `README.md`, `CHANGELOG.md`, `ROADMAP.md`, `capture/session-handoff.md`, `.bench/BENCH-reference.md`, `internal/anchors/registry_data.go`, `tests/canary/workflow-guidance-anchors`, `CONTEXT.md`
Integration surfaces: skills index and hook-layer lookup→`.bench/BENCH-reference.md`; roadmap FT107 rows→`ROADMAP.md`; changelog→`CHANGELOG.md`; residual anchors→`internal/anchors/registry_data.go`
Contracts: the two new skills' index entries crossing `.agents/skills/*/SKILL.md`→`.bench/BENCH-reference.md`, asserted by FC2 against the real generator; everything else this ticket touches is record-only
Closure: FC1/whole-reread, FC2/distribution-surfaces, FC3/records-updated

## What to build

The composed-coherence verification and record pass — verification and
records only; this ticket repairs nothing outside its fence. Reread every
changed guidance surface whole (the `craft-synthesis` legibility and
consistency rereads) and sweep for contradictions among triggers, budgets,
anchors, and phase handoffs, including cross-references from unchanged files
to retired ceremony (`rg` for: handoff ledger, red-mutation table,
`Contracts:`, `Integration surfaces:`, breakdown-review delegate, receipt).
A contradiction inside this fence (a stale index row, a residual anchor, a
record) is fixed here; a contradiction in any other file is reported as a
narrowly fenced repair-ticket proposal in the done-report, never patched from
this ticket. Verify both new skills are exposed from their canonical files:
index check green, `.claude/skills` symlinks resolve, existing adoption
controls green (`.bench/skills-index.sh --check`, focused link tests if
unsure). Update `README.md` and `CHANGELOG.md` for the doctrine stack, mark
the FT107 rows in `ROADMAP.md`, sweep `internal/anchors` for residual pins on
wording no file carries, and add glossary-only `CONTEXT.md` terms for the new
doctrine vocabulary (frontier round, disposition, prose budget) if missing.
Rewrite `capture/session-handoff.md` in full per its own shape. The final
full gate is the landing itself (`bench commit --spec`), not an acceptance
row. Won't-handle (explicit unused disposition): byte-identical and divergent
foreign adapter targets keep their landed hard refusals — this ticket does
not touch adoption behavior.

## Acceptance

- [ ] [FC1] (covers PG24) the whole-artifact reread finds no live cross-reference to retired ceremony and no trigger/budget/anchor/handoff contradiction; in-fence fixes and out-of-fence repair proposals are each listed in the done-report.
- [ ] [FC2] (covers local) index check green; both `.claude/skills` symlinks resolve; adoption controls green.
- [ ] [FC3] (covers local) README, CHANGELOG, ROADMAP, CONTEXT, and the handoff reflect the landed doctrine.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| FC1/whole-reread | leave one retired-term cross-reference in place | semantic review reread | the rg sweep returns a hit; review cites PG24 |
| FC2/distribution-surfaces | break one new-skill symlink | skills/adoption checks | point it at a missing dir, run `.bench/skills-index.sh --check` and the mirror check, expect red |
| FC3/records-updated | skip the ROADMAP flip | semantic review reread | reviewer-graded: stale roadmap against the phase-close contract |
