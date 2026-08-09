# Introduce run-scoped Bench selection

Blocked by: none
Ownership fence: `internal/runbinary/`
Integration surfaces: run owner and selected-path API→own-gate-run-binary.md; focused-test owner→own-focused-test-binary.md; ordinary phase consumer→route-ordinary-phase-plumbing.md; internal gate helper consumer→migrate-gate-helpers.md; contract/preflight consumer→migrate-contract-preflight-helpers.md; nested gate consumer→propagate-selected-binary-to-nested-gates.md; lifecycle contraction→contract-run-directory-lifecycle.md
Contracts: `Selection` (clean absolute executable path, exact source root, seal verification, ownership state) crosses `internal/runbinary/`→gate, testreport, helper, and canary consumers, membership is one selected host Bench executable per owner, ordering is create private directory then canonical subject build then validate then expose, absence is authorable only at a top-level owner and refuses at inherited consumers, asserted by RB1 against the real canonical builder seam
Closure: RB1/absolute-regular-path, RB1/exact-source-binding, RB1/canonical-host-build, RB1/private-run-lifetime, RB1/hostile-ambient-overwrite

## What to build

Add the expansion seam that creates and validates a private run-scoped Bench selection without changing every consumer at once. The API owns canonical-builder invocation and represents inherited selection distinctly from top-level authorship, so later migrations cannot accidentally treat a missing nested value as permission to build.

## Acceptance

- [ ] [RB1] (covers RS1) a top-level owner selects one canonical host executable at a cleaned absolute private path bound to the exact source, overwrites hostile ambient selection, and exposes no reusable cross-run location.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RB1/absolute-regular-path | return a relative or symlinked selected path | runbinary black-box test | create a selection, resolve it from another working directory, and require refusal before child launch |
| RB1/exact-source-binding | omit the source-root freshness verification | runbinary source-drift test | build from source A, mutate a freshness input, validate for source B, and require refusal |
| RB1/canonical-host-build | invoke `go build` directly or preserve ambient GOOS/GOARCH | canonical-builder trace test | author a selection under hostile target variables and require one subject-mode script trace for the host target |
| RB1/private-run-lifetime | place the selected executable under `dist/`, a cache, or a stable worktree path | runbinary filesystem test | author two runs and require distinct private parents outside durable repository/cache slots |
| RB1/hostile-ambient-overwrite | accept a caller-supplied top-level `BENCH_RUN_BINARY` | runbinary ambient-input test | plant an executable in the ambient path, create a top-level owner, and require a different owner-authored path |
