# Serialize the primary and stripped schedule

Blocked by: serialize-phase-tables.md
Ownership fence: `internal/gate/`
Integration surfaces: serial scheduler→serialize-phase-tables.md; combined schedule→contract-ordinary-build-census.md; combined cancellation/teardown→contract-run-directory-lifecycle.md
Contracts: primary and stripped `Phase` values cross split classification in `internal/gate/stripped_worktree.go`→one serial scheduler, membership is each declared phase exactly once with its primary or stripped execution root, ordering is one stable topological sequence across both sets, and empty or unequal partitions remain valid, asserted by SS1 against the real split runner
Closure: SS1/one-global-active, SS1/each-phase-once, SS1/primary-root, SS1/stripped-root, SS1/cross-table-order, SS1/empty-primary, SS1/empty-stripped, SS1/combined-red

## What to build

Replace the primary goroutine plus stripped runner with one composed schedule. Preserve the materialized stripped worktree's lifetime and phase directories, but do not give either partition its own scheduler invocation.

## Acceptance

- [ ] [SS1] (covers SG2) primary and stripped phases enter one stable serial topological schedule, each runs once at its correct root, empty partitions work, and red aggregation spans both sets.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SS1/one-global-active | retain separate primary and stripped runners | split barrier test | make both partitions ready and require global max active one |
| SS1/each-phase-once | duplicate a phase during composition | split record test | run a mixed table and require the declared multiset exactly once |
| SS1/primary-root | assign the stripped root to a primary phase | phase cwd marker | record each primary cwd and require the graded primary root |
| SS1/stripped-root | assign the primary root to a stripped phase | phase cwd marker | record each stripped cwd and require the materialized stripped root |
| SS1/cross-table-order | concatenate results after independently scheduling tables | cross-table needs test | declare a dependency-compatible mixed order and require one stable launch record |
| SS1/empty-primary | require a primary goroutine result | empty-partition test | run an all-stripped table and require completion without blocking |
| SS1/empty-stripped | require stripped materialization or result | empty-partition test | run an all-primary table and require no stripped runner |
| SS1/combined-red | discard one partition's red | split aggregation test | red one primary and one stripped phase and require aggregate red with both results |
