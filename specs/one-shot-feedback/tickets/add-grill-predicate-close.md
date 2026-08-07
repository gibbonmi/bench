# Close each grill decision as an exact predicate

Blocked by: add-falsification-predicate-questions.md
Ownership fence: `.agents/skills/bench-craft-grill/SKILL.md`, `internal/anchors/registry_data.go`
Integration surfaces: anchors Require row→`internal/anchors/registry_data.go`; none further — the grill skill's `.claude/skills/` entry is a symlink, so this edit is the whole edit
Contracts: none crosses — the needle and its Require row land inside this ticket's own fence
Closure: AB1/predicate-close

## What to build

`craft-grill` gains one discipline bullet (spec story 4): each decision closes
with the answer restated as the exact predicate it fixes, never an outcome
label. Add the bullet at the end of the `## Discipline` list, before `## Form`.
Needle-first: land the anchors `Require` row, observe the missing-anchor
conformance red, then the prose.

- GP1 needle: `close each decision by restating the answer as the exact
  predicate it fixes — never an outcome label`

Registry row: append at the true end of the registry slice, no `Group:` key,
diagnostic format `<file> missing acceptance coverage anchor: <needle>`.

## Acceptance

- [ ] [AB1] (covers GP1) `craft-grill`'s discipline list carries the predicate-close needle, with its Require row in the registry.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AB1/predicate-close | delete the predicate-close bullet from the grill skill (Require row retained) | conformance workflow-anchors check | `BENCH_CONFORMANCE_ROOT=$PWD go test -count=1 ./internal/conformance -run '^TestRootConformance$'`, expect the GP1 missing-anchor diagnostic |
