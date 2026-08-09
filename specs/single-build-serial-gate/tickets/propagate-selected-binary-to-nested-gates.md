# Propagate the selected binary to nested gates

Blocked by: introduce-run-scoped-bench-selection.md, own-gate-run-binary.md, route-ordinary-phase-plumbing.md
Ownership fence: `internal/canary/`
Integration surfaces: selected owner environment→own-gate-run-binary.md; phase launch propagation→route-ordinary-phase-plumbing.md; selection validation→introduce-run-scoped-bench-selection.md; nested consumer closure→contract-ordinary-build-census.md; nested teardown junction→contract-run-directory-lifecycle.md
Contracts: `BENCH_RUN_BINARY` crosses canary selection and `subjectCall` in `internal/canary/`→the real nested `.bench/gate.sh`, membership is each nested gate call but not compiled canary test executables, ordering preserves the exact outer value through environment sanitization then validates it before inner phase construction, and missing or invalid inheritance refuses with no authoring fallback, asserted by NG1 against a real outer/inner process harness
Closure: NG1/unchanged-path, NG1/zero-nested-build, NG1/missing-refusal, NG1/relative-refusal, NG1/symlink-refusal, NG1/stale-refusal, NG1/source-mismatch-refusal, NG1/compiled-canary-exception

## What to build

Carry the selected executable deliberately through canary environment scrubbing and nested gate launch. An inner gate is always a consumer: it validates inheritance and never calls the canonical builder. Keep `go test -c` canary compilation because it authors test executables whose compiler behavior is the proof.

## Acceptance

- [ ] [NG1] (covers RS7) every nested gate observes the exact outer absolute path and builds zero, invalid inheritance refuses before phases, and intentional compiled-canary test executables remain independent.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| NG1/unchanged-path | strip or rewrite selection between outer and inner | nested marker harness | run one canary inner gate and require identical path, inode, and digest markers |
| NG1/zero-nested-build | let the inner gate create an owner | nested builder-trap harness | run a valid inner gate and require zero nested builder traces |
| NG1/missing-refusal | author when selection is absent | nested refusal harness | remove the variable, trap builder and phases, and require refusal before both |
| NG1/relative-refusal | accept a relative inherited value | nested refusal harness | pass `./bench`, change cwd, and require bounded refusal |
| NG1/symlink-refusal | follow a selected symlink | nested refusal harness | point a symlink at valid bytes and require refusal before execution |
| NG1/stale-refusal | accept a stale seal | nested freshness harness | mutate a selected byte after publication and require refusal |
| NG1/source-mismatch-refusal | validate against the outer source instead of the inner graded source contract | nested source harness | pair valid bytes with a mismatched source root and require refusal |
| NG1/compiled-canary-exception | route `go test -c` through the Bench selection | compiled-bite proof | execute the compile proof and require a distinct test executable that observes the intended package mutation |
