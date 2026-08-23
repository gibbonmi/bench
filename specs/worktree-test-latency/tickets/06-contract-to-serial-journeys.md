# Contract to serial public journeys

Blocked by: 01-select-one-test-run-binary.md, 03-extract-landing-policy.md, 04-extract-lifecycle-policy.md, 05-extract-reclaim-policy.md
Writes: internal/worktree/, specs/worktree-test-latency/evidence/coverage-ledger.md (new)

## What to build

Route every descendant start through one serial journey harness. Remove migrated
ambient compatibility helpers and redundant repository-backed policy partitions.

Create the parent test-run proof inventory for every required journey class and
one fact adapter per policy owner. A package run fails by identifier when any
required proof is absent.

Write `specs/worktree-test-latency/evidence/coverage-ledger.md` before deletion.
Map each removed repository-backed test to one retained pure, adapter, or
public-journey test. Compare every deleted test function in the ticket's
base-to-tip diff with that ledger and record the result.

## Acceptance

- [ ] RJ1: Every named Git-supplied behavior has one representative serial journey.
- [ ] FA1: Landing, lifecycle, and reclaim each record one focused adapter proof.
- [ ] CV1: Every removed repository-backed test has a surviving coverage disposition.
- [ ] EI1: Tests outside the journey harness mutate no environment or current directory.
- [ ] NP1: The first-spec diff adds no `t.Parallel` and no scheduler.
