# Select one test-run Bench binary

Blocked by: none
Writes: internal/worktree/main_test.go, internal/worktree/land_test.go, internal/worktree/clean_landed_apply_test.go, internal/worktree/test_run_test.go (new)

## What to build

Add one worktree test-run executable owner. Reuse an inherited selected binary,
or build and seal one direct-run binary, then pass its identity to every public journey.

Remove every per-journey selected-binary build. Keep invalid, stale, and seal-less
executable refusals as explicit journey cases.

## Acceptance

- [ ] SB1: A direct run builds and seals one Bench executable for all public journeys.
- [ ] SB2: A gate-selected executable reaches multiple journeys unchanged with zero private builds.
- [ ] SB3: Missing, invalid, stale, and seal-less selections refuse before any journey.
