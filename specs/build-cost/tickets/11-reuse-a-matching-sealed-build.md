# Reuse a sealed build that matches the target tree

Blocked by: none
Writes: internal/treetarget/build.go, internal/treetarget/run.go, internal/treetarget/build_test.go, internal/treetarget/run_test.go, internal/systemtest/tree_target_test.go, internal/worktree/target_refusal_test.go
Covers: BC61, BC62, BC63, BC64, BC65

## What to build

A kit worktree target picks its child executable from three candidates in this order:

1. its own published build
2. the primary checkout's published build
3. the running executable

It runs the first candidate that `freshness.VerifyAt(root, root, candidate)` accepts against the target tree. The loop line reads `if freshness.VerifyAt(root, root, candidate) == nil {`, because the probe-11 command swaps that exact line.
When no candidate passes, the current refusals stay. An absent own build refuses `worktree build is missing`, and each other case refuses `worktree build does not match the tree`. No child starts after a refusal.
The primary checkout target and a non-kit target keep their current executable choice.
The system test pins `BENCH_KIT` through its fixture environment, as `TestTreeTargetRefusesStaleKitBuild` does now, so an ambient `BENCH_KIT` value cannot change its verdict.

## Acceptance

- [ ] A target without its own build runs the primary build when that seal matches (BC61).
- [ ] A target without its own build runs the running executable when its seal matches and the primary build's does not (BC62).
- [ ] A target whose own build matches runs its own build when the primary build also matches (BC63).
- [ ] An edited build input with no matching candidate refuses `worktree build does not match the tree` and starts no child (BC64).
- [ ] A target without its own build and with no matching candidate refuses `worktree build is missing` and starts no child (BC65).
