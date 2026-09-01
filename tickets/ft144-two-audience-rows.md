# Add the two-audience prompt and the transaction failure rows to the coverage-map discipline

Blocked by: none
Writes: .agents/skills/bench-craft-spec/references/map-discipline.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, tests/canary/workflow-guidance-anchors/craft-spec-two-audience-inventory (new), tests/canary/workflow-guidance-anchors/craft-spec-transaction-failure-rows (new)

## What to build

The reference file `map-discipline.md` gains two rules. The first rule goes
under `## In the edge inventory`. A kit spec names the audience each behavior
serves: this repository, or every repository that links the kit. The inventory
walks the absent-versus-empty pair for each directory the spec reads, because
the two audiences can want different answers. The rule does not restate the
absent-file line in `references/hostile-input-library.md`; the audience split
is the new fact.

The second rule goes under `## Per row`. A transaction-shaped spec classifies
its verification failures into three rows: persistence before the oracle runs,
interruption inside the oracle, and persistence at the terminal step.

Each new rule gets one anchor in `internal/anchors/registry_data.go`, one
red-on-removal assertion in `internal/anchors/registry_data_test.go`, and one
canary fixture under `tests/canary/workflow-guidance-anchors/`. The
post-approval-edit ruling is out of scope; it rides its own decision record.

## Acceptance

- [ ] `map-discipline.md` states the two-audience rule with the absent-versus-empty pair under `## In the edge inventory`.
- [ ] `map-discipline.md` states the three transaction failure rows under `## Per row`.
- [ ] Removal of either new needle turns the anchors test red.
- [ ] Each new canary fixture bites through its registered owner.
- [ ] The prose lanes stay green on every edited Markdown file.
