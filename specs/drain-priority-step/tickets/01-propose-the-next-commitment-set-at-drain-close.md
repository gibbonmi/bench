# Propose the next commitment set at drain close

Blocked by: none
Writes: .agents/commands/bench-drain.md, CONTEXT.md, internal/anchors/registry_commitment.go, internal/conformance/commitment_guidance_test.go, internal/conformance/recurrence_maintenance_contract_test.go, CHANGELOG.md
Covers: none

## What to build

The reviewer directed on 2026-10-10 that a drain ends with a priority review. Today `/bench-drain` closes with the landing, the sealed-generation retirement, the `commitment_outlook` start, and `bench handoff`. No step reviews the roadmap priorities or proposes the next commitments.

The closing steps of `.agents/commands/bench-drain.md` get one last step. After the sealed-generation retirement, the coordinator charges read-only research delegates under `craft-delegate`, `craft-line`, and `craft-research`. The step starts before or with the `commitment_outlook` start, and `bench handoff` stays last. The delegates read the landed `ROADMAP.md` and the staged specs, and they start from the rows that the reviewer already prioritized. The coordinator reports one proposed next commitment set in the exit, in the purpose order that `.bench/BENCH.md` owns.

The proposal is unapproved input, and the closing exit is its only output. Only the reviewer's direction adopts it, through the commitment rule in `.bench/BENCH.md`. The step does not write the `## Recommended sequence` section, so `## 7. Leave the sequence to the commitment` stays true.

The reviewer directed on 2026-10-10 that the drain also prompts the decisions, and named this part ticket triage. In ticket triage, the coordinator puts each reviewer decision that the proposal depends on to the reviewer as a structured selection prompt. These decisions are the proposed order and each open decision that a proposed row names. Each prompt gives one option set for one decision, with the recommended option first. The prompt uses the native question surface of the harness when one exists, and a plain numbered list otherwise. The answers are the decision reference for `bench commitment approve`, and the coordinator records each settled row decision in its `roadmap/FT<n>.md`.

Ticket triage prompts only a decision that the reviewer can settle in one selection. A row whose next action is `/bench-shape-idea` is not part of ticket triage. Such a row has `Next: shape`, or a decision tree that needs more than one selection. The closing exit lists each such row as work for a separate session and names the harness-native shape command for it. `CONTEXT.md` defines the term ticket triage once, and the drain uses that term.

The drain keeps one read-delegate cap. The cap sentence in `## Delegate the evidence` becomes a bound on the live read delegates, and the new step points to that cap. Every evidence read returns before the batch worktree opens, so the research readers fit inside the same cap. The recurrence maintenance contract check pins the reworded cap sentence.

Commitment anchors pin four sentences of the new step: the key sentence, the unapproved-input sentence, the decision-prompt sentence, and the shape-exclusion sentence. For each sentence, the commitment guidance test removes it and expects its diagnostic.

## Acceptance

- [ ] `.agents/commands/bench-drain.md` charges read-only research delegates under `craft-delegate`, `craft-line`, and `craft-research` after the sealed-generation retirement and before `bench handoff`.
- [ ] The step starts from the rows that the reviewer prioritized in `ROADMAP.md` and from the staged specs.
- [ ] The exit reports one proposed next commitment set, ordered by the purpose rule in `.bench/BENCH.md`.
- [ ] The drain states that the proposal is unapproved input that only the reviewer's direction adopts, and it restates no commitment rule.
- [ ] The drain states one read-delegate cap, and the new step points to it.
- [ ] In ticket triage, the drain puts each reviewer decision of the proposal to the reviewer as a structured selection prompt.
- [ ] Each prompt has one option set for one decision, with the recommended option first.
- [ ] The drain names the native question surface of the harness for the prompts, and a plain numbered list where no such surface exists.
- [ ] The drain states that the answers are the decision reference for `bench commitment approve`.
- [ ] The drain records each settled row decision in its `roadmap/FT<n>.md`.
- [ ] Ticket triage excludes a row that needs `/bench-shape-idea`, and the closing exit lists that row as separate-session work with its shape command.
- [ ] `CONTEXT.md` defines ticket triage, and `.agents/commands/bench-drain.md` uses the term without a second definition.
- [ ] The commitment guidance check goes red when any of the four pinned sentences is removed from `.agents/commands/bench-drain.md`.
- [ ] `CHANGELOG.md` records the change.
