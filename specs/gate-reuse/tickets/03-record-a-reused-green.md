# Record a reused green in the checkout's gate cache

Blocked by: 02-name-evidence-by-the-grade-key.md
Writes: internal/gate/run_transaction.go, internal/gate/run_durability.go, internal/gate/verdict.go, internal/gate/reused_green_test.go (new), internal/gate/durable_replacement_test.go, internal/gate/run_outcomes_test.go
Covers: GR22, GR47, GR23, GR24, GR44

## What to build

A run that reuses a store green writes the checkout's gate cache, so `bench status` and the green-marker bootstrap see the reused green.
The run evaluates the reuse at its two existing sites in `executeSubjectOwning`, before the gate lock and under it.

- Before the lock, a reuse that the gate cache already describes as this subject's fresh green answers with no write.
- When another run holds the gate lock, a reuse that held before the lock answers with no write.
- Under the lock, a reuse writes the cache unless the cache already describes this subject's fresh green.

The written record holds the tree, the checkout's oracle digest, `green`, and the recorded time of the store record.
So an absent cache, a cache of another tree, and a red, timed-out, or expired cache of the same tree each take the write.
A failed write returns an operational result, and the run reports no green. `reusedGreenResult` keeps its two callers in `executeSubjectOwning`.

## Acceptance

- [ ] A primary run with no gate cache that reuses a linked-worktree green writes the cache with the tree, `green`, and the source run's recorded time.
- [ ] A primary run whose cache holds a red verdict of the same tree replaces it with the reused green.
- [ ] A reuse that the cache already describes leaves the cache bytes unchanged.
- [ ] A reuse whose cache write fails returns an operational result with no reusable green.
- [ ] A reuse while another run holds the gate lock prints the reuse line and leaves the cache bytes unchanged.
