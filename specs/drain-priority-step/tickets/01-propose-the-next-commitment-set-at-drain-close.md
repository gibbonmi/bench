# Propose the next commitment set at drain close

Blocked by: none
Writes: .agents/commands/bench-drain.md, internal/anchors/registry_commitment.go, internal/conformance/commitment_guidance_test.go, internal/conformance/recurrence_maintenance_contract_test.go, CHANGELOG.md
Covers: none

## What to build

The reviewer directed on 2026-10-10 that a drain ends with a priority review. Today `/bench-drain` closes with the landing, the sealed-generation retirement, the `commitment_outlook` start, and `bench handoff`. No step reviews the roadmap priorities or proposes the next commitments.

The closing steps of `.agents/commands/bench-drain.md` get one last step. After the sealed-generation retirement, the coordinator charges read-only research delegates under `craft-delegate`, `craft-line`, and `craft-research`. The step starts before or with the `commitment_outlook` start, and `bench handoff` stays last. The delegates read the landed `ROADMAP.md` and the staged specs, and they start from the rows that the reviewer already prioritized. The coordinator reports one proposed next commitment set in the exit, in the purpose order that `.bench/BENCH.md` owns.

The proposal is unapproved input. Only the reviewer's direction adopts it, through the commitment rule in `.bench/BENCH.md`. The step does not write the `## Recommended sequence` section, so `## 7. Leave the sequence to the commitment` stays true.

The drain keeps one read-delegate cap. The cap sentence in `## Delegate the evidence` becomes a bound on the live read delegates, and the new step points to that cap. Every evidence read returns before the batch worktree opens, so the research readers fit inside the same cap. The recurrence maintenance contract check pins the reworded cap sentence.

A commitment anchor pins the key sentence of the new step, and the commitment guidance test removes that sentence and expects its diagnostic.

## Acceptance

- [ ] `.agents/commands/bench-drain.md` charges read-only research delegates under `craft-delegate`, `craft-line`, and `craft-research` after the sealed-generation retirement and before `bench handoff`.
- [ ] The step starts from the rows that the reviewer prioritized in `ROADMAP.md` and from the staged specs.
- [ ] The exit reports one proposed next commitment set, ordered by the purpose rule in `.bench/BENCH.md`.
- [ ] The drain states that the proposal is unapproved input that only the reviewer's direction adopts, and it restates no commitment rule.
- [ ] The drain states one read-delegate cap, and the new step points to it.
- [ ] The commitment guidance check goes red when the new key sentence is removed from `.agents/commands/bench-drain.md`.
- [ ] `CHANGELOG.md` records the change.
