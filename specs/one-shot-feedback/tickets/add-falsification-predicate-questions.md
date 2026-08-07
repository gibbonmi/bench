# Add the exact-predicate and partition questions to the falsification charge

Blocked by: add-retro-repair-attribution.md
Ownership fence: `.agents/commands/bench-write-spec.md`, `internal/anchors/registry_data.go`
Integration surfaces: anchors Require rows→`internal/anchors/registry_data.go`; registry append tail→add-grill-predicate-close.md
Contracts: none crosses
Closure: AB1/predicate-question, AB2/partition-question

## What to build

`/bench-write-spec` step 9's falsification charge asks two more questions
(spec story 3), extending the existing charge sentence in place — no new files,
no new skills. Needle-first per row: land each anchors `Require` row, observe
the missing-anchor conformance red, then the prose.

- FC1 needle: `does any behavior, red signal, or decision answer name an outcome
  family instead of an exact predicate`
- FC2 needle: `could a narrower capability ship on its own gate` — phrase the
  question so the story-partition ask (do the stories partition into disjoint
  package or fence sets) contains this exact substring.

Registry rows: append at the true end of the registry slice, no `Group:` key,
diagnostic format `<file> missing acceptance coverage anchor: <needle>`.

## Acceptance

- [ ] [AB1] (covers FC1) the step-9 charge carries the exact-predicate question needle, with its Require row in the registry.
- [ ] [AB2] (covers FC2) the step-9 charge carries the story-partition question needle, with its Require row in the registry.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AB1/predicate-question | delete the exact-predicate question from step 9 (Require row retained) | conformance workflow-anchors check | `BENCH_CONFORMANCE_ROOT=$PWD go test -count=1 ./internal/conformance -run '^TestRootConformance$'`, expect the FC1 missing-anchor diagnostic |
| AB2/partition-question | delete the partition question from step 9 (Require row retained) | conformance workflow-anchors check | same command, expect the FC2 missing-anchor diagnostic |
