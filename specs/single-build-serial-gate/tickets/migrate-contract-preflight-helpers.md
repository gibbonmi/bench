# Migrate contract and preflight helpers

Blocked by: introduce-run-scoped-bench-selection.md, own-focused-test-binary.md, route-ordinary-phase-plumbing.md
Ownership fence: `internal/contract/`, `internal/preflight/`
Integration surfaces: selection API→introduce-run-scoped-bench-selection.md; focused-test environment→own-focused-test-binary.md; selected suite environment→route-ordinary-phase-plumbing.md; migrated helper inventory→contract-ordinary-build-census.md
Contracts: selected executable crosses focused or gate child environment→`internal/contract/helper.go`, fixture Bench commands, wrapper/cache fixture byte installation, and preflight current-subject helpers under `internal/preflight/`, membership excludes the spec's changed-source and intentional proof files, ordering validates selection before copy, link, or execution, and absence refuses without building, asserted by CH1 against the real contract/preflight suites
Closure: CH1/require-fresh-selected, CH1/fixture-bench-selected, CH1/coverage-selected, CH1/wrapper-fixture-selected-bytes, CH1/preflight-shared-helper, CH1/preflight-inline-builders, CH1/changed-source-proof-retained, CH1/missing-selection-refusal

## What to build

Turn the contract helper into the one source for the selected unchanged-host Bench path, route ordinary fixture commands through it, and have wrapper/cache tests install selected bytes. Collapse all preflight current-subject builders onto that helper. Leave changed-source, artifact-mode, release, reproducibility, and compiler-observing proofs independent.

## Acceptance

- [ ] [CH1] (covers RS6) ordinary contract commands including coverage and preflight helpers execute or install the inherited selected bytes with no subject build, retain the named changed-source proofs, and fail closed when an unchanged-host test lacks selection.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CH1/require-fresh-selected | verify `SubjectRoot/dist/bench` instead of the selected path | contract freshness helper test | plant stale dist bytes beside a valid selection and require the selected executable to pass |
| CH1/fixture-bench-selected | route ordinary Fixture.Bench through the kit wrapper | contract fixture identity test | plant a wrapper-resolved alternate binary and require ordinary fixture execution uses selection |
| CH1/coverage-selected | execute coverage through wrapper, dist, or a private build | real AXI coverage contract | plant alternate wrapper/dist bytes, run coverage, and require the selected executable identity with zero builder traces |
| CH1/wrapper-fixture-selected-bytes | compile bytes for a cache or installation fixture | wrapper/cache fixture trace test | exercise the fixture and require zero builder traces plus byte equality with selection |
| CH1/preflight-shared-helper | make the shared preflight helper invoke the builder | preflight builder-count test | run all helper consumers under a builder trap and require zero invocations |
| CH1/preflight-inline-builders | retain any inline current-subject builder | preflight structural census | scan post-migration preflight commands and require no subject-mode builder constructor |
| CH1/changed-source-proof-retained | replace a changed-source build with selected unchanged bytes | changed-source contract test | mutate the cloned source and require its independently built output follows the mutation |
| CH1/missing-selection-refusal | fall back to dist or builder | contract/preflight refusal test | clear selection, plant dist and builder traps, and require bounded refusal before either executes |
