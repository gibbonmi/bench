# Own the gate run binary

Blocked by: introduce-run-scoped-bench-selection.md
Ownership fence: `internal/gate/gate.go`, `internal/gate/engine.go`, `internal/gate/prospective.go`, `internal/gate/prospective_test.go`, `internal/gate/run_binary_test.go`, `.bench/gate.sh`, `.bench/gate-prospective.sh`, `internal/conformance/gate_entry_test.go`, `internal/canary/gate_entry_test.go`
Integration surfaces: selection producer→introduce-run-scoped-bench-selection.md; selected gate environment→route-ordinary-phase-plumbing.md; outer selection→propagate-selected-binary-to-nested-gates.md; gate-owner teardown→contract-run-directory-lifecycle.md; shell entry advertisements→owned conformance and canary entry tests
Contracts: `Selection` crosses `internal/gate/gate.go`→`.bench/gate.sh` as `BENCH_RUN_BINARY`, membership is each non-reused direct or prospective execution, ordering is accept/admit exact subject then author once then execute the gate, while exact-green reuse and pre-owner refusal author zero; prospective Bench source is the materialized runtime root and linked-project source is the captured kit, asserted by GO1 against both real entry routes
Closure: GO1/direct-single-build, GO1/prospective-single-build, GO1/shared-owner, GO1/linked-kit-source, GO1/preowner-zero-build, GO1/reused-green-zero-build, GO1/selected-self-verifier, GO1/no-prospective-builder

## What to build

Move direct and prospective execution behind the shared run owner. Keep `.bench/gate.sh` as shell, but make it require the selected executable for both freshness self-check and phase routing. Contract `.bench/gate-prospective.sh` to a no-build pass-through. The owner must remain compatible with the existing exact-subject lock and verdict flow.

## Acceptance

- [ ] [GO1] (covers RS2) each non-reused direct or prospective gate uses one shared owner and exactly one canonical build from the materialized Bench candidate or captured linked-project kit, both shell routes execute the selected binary for freshness and phases, and pre-owner refusal or exact-green reuse builds zero.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| GO1/direct-single-build | call the canonical builder a second time on the direct route | real direct-gate build-trace test | run one full direct gate and require exactly one subject-mode trace |
| GO1/prospective-single-build | retain a builder call in `.bench/gate-prospective.sh` | real prospective-gate build-trace test | execute an unpublished Bench tree and require exactly one owner trace and no root-local `dist` output |
| GO1/shared-owner | give prospective execution a separate run-directory constructor | direct/prospective identity test | execute both routes through the injected owner and require the same ownership implementation and environment contract |
| GO1/linked-kit-source | build a linked non-Bench project's runtime root instead of its captured wrapper-selected kit | linked prospective source-marker test | give project and kit distinct builder markers, execute the linked prospective tree, and require exactly the kit marker with zero project-root builds |
| GO1/preowner-zero-build | create the owner before admission or exact-subject validation | gate refusal test | force lock contention and subject drift and require zero builder traces and no run directory |
| GO1/reused-green-zero-build | author before checking exact-green reuse | gate reuse test | seed exact green, rerun, and require the reuse announcement with zero builder traces |
| GO1/selected-self-verifier | invoke `go run ./internal/freshness/check` from gate shell | gate-entry conformance test | run the real shell under a Go trap and require selected `freshness-check` before selected `gate-phases` |
| GO1/no-prospective-builder | restore subject-builder text to the prospective shell | gate-entry structural test | inspect assembled shell commands and require no builder constructor in `.bench/gate-prospective.sh` |
