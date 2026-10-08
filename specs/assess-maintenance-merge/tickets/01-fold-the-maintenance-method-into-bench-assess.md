# Fold the maintenance assessment method into bench-assess

Blocked by: none
Writes: .agents/commands/bench-assess.md, CHANGELOG.md, internal/anchors/registry_data.go, internal/anchors/registry_retained_workflow.go
Covers: none

## What to build

Two methods assess the Bench platform today. The kit phase `/bench-assess`
replaces the root `ASSESSMENT.md`. A skill at
`~/.codex/skills/bench-mainenance/SKILL.md` produced
`capture/bench-maintenance-assessment-2026-10-07.md`. Only Codex reads that
skill. Claude Code, the gate, and linked repos never see it, and the two
methods disagree on their output.

On 2026-10-08 the reviewer decided to merge the skill into `/bench-assess` and
to retire the Codex-only copy. The decision record is Appendix B of
`capture/skills-comparison-assessment-2026-10-08.md`.

Add the unique rules of the skill to `.agents/commands/bench-assess.md`. Each
rule appears in one place:

- At entry, record the repository, the source commit, the scope, and the
  requested stopping point.
- In the tests area, assess test cost through test purpose. A cheaper test
  keeps the same failure witness at the right seam. Package durations overlap,
  so their savings do not add. A shared helper alone proves no speedup. A speed
  gain inside the recorded run-to-run spread is no change, and the assessment
  records each reverted attempt.
- In the guidance area, separate conflicting authority from useful constraints
  and lookup detail. Fewer words alone are not a better outcome.
- For each finding, separate a source mechanism, an inference, and an executed
  reproduction. A defect names its wrong observable result and its smallest
  discriminating regression. A repair stays apart from a refactor. A stale
  roadmap premise is corrected from current code. A removal proposal names its
  consumers inside and outside the tree, its replacement, and the reason the
  surface exists.
- Map each finding to its existing owner with an overlap class: complete,
  partial, separate, or unresolved. Give each active roadmap row a proposed
  disposition, and list parked rows with their return conditions. Rank by
  consequence, evidence, recurrence, test value, effort, and prerequisites.
  Recurrence includes repeated fix commits in one area. Keep literal blockers
  apart from the proposed order. The order changes no commitment.
- Planning follow-through runs only on request. It hands off to
  `/bench-shape-idea` and `/bench-write-spec`, and it reports the exact state of
  each outcome with one contrastive pair.

Keep the root `ASSESSMENT.md` contract, both description lines, and every
anchored sentence. The phase stays read-only for the roadmap: new backlog still
reaches `ROADMAP.md` only through `/bench-drain`. Add one `CHANGELOG.md` entry
under Unreleased.

## Acceptance

- [ ] `.agents/commands/bench-assess.md` states each rule in the list above, and each rule appears in one place.
- [ ] The phase still replaces the root `ASSESSMENT.md`, parks new backlog with `bench idea`, and edits no roadmap row.
- [ ] `bench anchors .agents/commands/bench-assess.md` reports each require needle at a nonzero line and the forbid needle at line 0.
- [ ] `bench gate-prose` passes on `.agents/commands/bench-assess.md`, `CHANGELOG.md`, and this ticket.
- [ ] `CHANGELOG.md` holds one Unreleased entry for the change.
