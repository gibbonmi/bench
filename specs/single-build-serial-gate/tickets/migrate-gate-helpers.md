# Migrate gate helpers

Blocked by: introduce-run-scoped-bench-selection.md, route-ordinary-phase-plumbing.md
Ownership fence: `internal/gate/`
Integration surfaces: selected phase contract→route-ordinary-phase-plumbing.md; selection API→introduce-run-scoped-bench-selection.md; serial scheduler tests→serialize-phase-tables.md; final helper census→contract-ordinary-build-census.md
Contracts: selected executable crosses the gate run environment→unchanged-host helpers in `internal/gate/`, membership is the current-subject binary helper, phase-command fixtures, and prospective/current-subject execution fixtures discovered from the post-plumbing tree, ordering is require selection before helper execution, and absence refuses without subject-builder or dist fallback, asserted by GH1 against the real helper census
Closure: GH1/runner-serial-builder, GH1/current-subject-dist-paths, GH1/phase-fixture-selection, GH1/missing-selection-refusal

## What to build

Replace ordinary current-subject executable construction and hard-coded `dist/bench` reads under `internal/gate` with the selected helper. Keep alternate-package, changed-source, and compiler/linker proof builders intact for the final exception census.

## Acceptance

- [ ] [GH1] (covers RS5) every post-plumbing unchanged-host gate helper consumes the inherited selected executable, records identical path/bytes, and refuses missing selection without a private builder or `dist` fallback.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| GH1/runner-serial-builder | restore `currentBenchBinary`'s canonical-builder call | gate helper structural test | run the helper census and require no subject builder in `runner_serial_test.go` |
| GH1/current-subject-dist-paths | execute the current subject from `dist/bench` | gate selected-identity test | plant different dist bytes and require the selected marker at the helper consumer |
| GH1/phase-fixture-selection | copy a process-local test executable instead of selected bytes | gate phase fixture test | run the real fixture and require its Bench identity matches the owner selection |
| GH1/missing-selection-refusal | fall back to authoring or dist when selection is absent | gate helper refusal test | remove `BENCH_RUN_BINARY`, trap builder/dist execution, and require bounded refusal |
