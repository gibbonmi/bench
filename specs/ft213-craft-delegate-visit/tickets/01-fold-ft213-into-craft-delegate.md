# Fold the FT213 clauses into craft-delegate

Blocked by: none
Writes: .agents/skills/bench-craft-delegate/SKILL.md, .agents/skills/bench-craft-delegate/references/delegation-discipline.md (new), .agents/commands/bench-write-spec.md, internal/canary/mutation.go, internal/canary/mutation_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, ROADMAP.md, roadmap/FT213.md (deleted)

## What to build

A coordinator who loads `craft-delegate` finds every rule that
`roadmap/FT213.md` records, placed where the coordinator acts on it. Three
rules stay in `SKILL.md`: a read-only delegate that reads a tree the
coordinator will grade runs in its own worktree, a mutation probe requires a
behavioral mutation, and one pointer to the new reference. The remaining
clauses move behind that pointer into `references/delegation-discipline.md`,
one rule each, grouped by the moment the coordinator applies them. The
review-round cap sentence lands in `bench-write-spec`. `RestoreMutationFixture`
refuses a destination equal to its root instead of a silent no-op. Each file
stays inside its `projects/benchkit.md` budget, and every anchored sentence
keeps its bytes and line breaks. The landing retires the FT213 row: the index
line, the sequence entry, and the detail file leave together.

## Acceptance

- [ ] The `Isolation` section of `.agents/skills/bench-craft-delegate/SKILL.md` gives a read-only delegate its own worktree when it reads a tree the coordinator will grade. Without that worktree, the coordinator verifies the tree unchanged before the landing gate. The sentence `Read-only delegations need no worktree` is gone.
- [ ] `.agents/skills/bench-craft-delegate/SKILL.md` states that a mutation probe requires a behavioral mutation. It states that a probe which fails to compile proves nothing. It states that a probe of provably redundant code passes by construction.
- [ ] `.agents/skills/bench-craft-delegate/SKILL.md` points at `references/delegation-discipline.md` once.
- [ ] `references/delegation-discipline.md` carries every remaining clause of `roadmap/FT213.md` above its `Next:` line as one rule each, in ASD-STE100 prose. No rule is stated twice across the two files. A clause `SKILL.md` already states, such as the exec-only command form, is not restated.
- [ ] `.agents/commands/bench-write-spec.md` states that the review round declares its iteration cap before the first charge, within its 73-line budget.
- [ ] `RestoreMutationFixture` returns an error that names the refusal when `dst` and `root` resolve to the same directory, and writes nothing. `internal/canary/mutation_test.go` shows that row red before the edit and green after.
- [ ] `internal/anchors/registry_data.go` pins the read-only-worktree sentence, the behavioral-mutation sentence, and the reference pointer, and `internal/anchors/registry_data_test.go` shows each pin red on removal.
- [ ] `TestGuidanceProseBudgetsHoldOnTheLiveTree` stays green with no change to the budget table.
- [ ] `TestProseMechanicsHoldsOnTheLiveTree`, the fixture-bite tests, the `workflow-guidance-anchors` canary fixtures, and the retired-token check stay green.
- [ ] `ROADMAP.md` carries no FT213 line, its recommended sequence renumbers from 1, and `roadmap/FT213.md` is absent.
