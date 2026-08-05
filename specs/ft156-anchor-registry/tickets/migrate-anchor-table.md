# Migrate the anchor table

Blocked by: centralize-anchor-matching.md
Ownership fence: `internal/anchors`, `internal/conformance/docs_workflow_helpers_test.go`
Contracts: ordered anchor tuples and matcher diagnostics cross `internal/anchors`→`internal/conformance/docs_workflow_helpers_test.go`, asserted by AR1 against every pre-migration tuple and by AR2 against the real canary family

## What to build

Move every uniform needle declaration into the ordered `internal/anchors` registry and make conformance evaluate that table. Bespoke non-needle checks stay behind, diagnostics remain byte-identical, and the old private declarations disappear.

## Acceptance

- [ ] [AR1] A one-time instrumented inventory proves every pre-migration tuple of file, kind, section, needle, and diagnostic appears exactly in the registry.
- [ ] [AR2] The complete existing `workflow-guidance-anchors` family remains green against the registry-backed conformance evaluator.
- [ ] [AR3] `internal/conformance/docs_workflow_helpers_test.go` is at or below its 660-line grant and the post-change structure report adds no finding.
- [ ] [AR4] No residual uniform needle table or matcher copy remains anywhere in `internal/conformance`.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AR1 | drop one enumerated anchor tuple during migration | the coordinator-owned tuple inventory | enumerate the pre-migration declarations, enumerate registry rows, compare tuple-wise, and expect the named missing tuple |
| AR2 | change one migrated diagnostic byte string | the unchanged canary fixture family | run the full owning family and expect the fixture for that diagnostic to fail |
| AR3 | leave a private needle block in the helpers file | the recorded before/after structure receipt | run `bench structure` and expect the helpers path to remain over its 660-line grant |
| AR4 | retain a local anchor matcher or uniform needle declaration in conformance | the whole-package residue sweep | enumerate matching and needle declarations under `internal/conformance` and expect the named residual site |
